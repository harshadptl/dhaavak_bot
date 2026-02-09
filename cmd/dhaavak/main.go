package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/harshadpatil/dhaavak/internal/agent"
	"github.com/harshadpatil/dhaavak/internal/channel"
	"github.com/harshadpatil/dhaavak/internal/channel/telegram"
	"github.com/harshadpatil/dhaavak/internal/config"
	"github.com/harshadpatil/dhaavak/internal/gateway"
	"github.com/harshadpatil/dhaavak/internal/llm"
	"github.com/harshadpatil/dhaavak/internal/queue"
	"github.com/harshadpatil/dhaavak/internal/routing"
	"github.com/harshadpatil/dhaavak/internal/session"
	"github.com/harshadpatil/dhaavak/pkg/protocol"
)

func main() {
	configPath := flag.String("config", "dhaavak.yaml", "path to config file")
	flag.Parse()

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	cfg, err := config.Load(*configPath)
	if err != nil {
		slog.Error("failed to load config", "err", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// --- Session Manager ---
	sessionMgr := session.NewManager(cfg.Session.TTL, cfg.Session.MaxHistory)
	sessionMgr.StartCleanup(ctx, cfg.Session.CleanupInterval)

	// --- Queue Manager ---
	queueMgr := queue.NewManager(ctx, cfg.Queue.BufferSize, cfg.Queue.IdleTimeout)
	queueMgr.StartCleanup(ctx, cfg.Queue.CleanupInterval)

	// --- Router ---
	var bindings []routing.Binding
	for _, b := range cfg.Channels.Telegram.Bindings {
		bindings = append(bindings, routing.Binding{
			Channel:  "telegram",
			PeerKind: b.PeerKind,
			PeerID:   b.PeerID,
			AgentID:  b.AgentID,
		})
	}
	defaultAgent := ""
	if len(cfg.Agents) > 0 {
		defaultAgent = cfg.Agents[0].ID
	}
	if cfg.Channels.Telegram.DefaultAgent != "" {
		defaultAgent = cfg.Channels.Telegram.DefaultAgent
	}
	store := routing.NewBindingStore(bindings, defaultAgent)
	router := routing.NewResolver(store)

	// --- LLM Provider ---
	var provider llm.Provider
	switch cfg.LLM.Provider {
	case "anthropic":
		provider = llm.NewAnthropicProvider(cfg.LLM.APIKey, cfg.LLM.Model)
	default:
		slog.Error("unsupported LLM provider", "provider", cfg.LLM.Provider)
		os.Exit(1)
	}

	// --- Agent Runtime ---
	runtime := agent.NewRuntime(provider, cfg.LLM.MaxTurns)
	for _, a := range cfg.Agents {
		runtime.RegisterAgent(agent.AgentDef{
			ID:           a.ID,
			Name:         a.Name,
			SystemPrompt: a.SystemPrompt,
			Model:        a.Model,
		})
	}

	// --- Gateway Server ---
	gw := gateway.New(cfg.Server, cfg.Auth.Token)

	// --- Channel Registry ---
	registry := channel.NewRegistry()

	// Wire agent event sink -> gateway broadcast.
	runtime.SetEventSink(func(evt agent.Event) {
		switch evt.Type {
		case "delta":
			gw.ChatState.AccumulateDelta(evt.SessionID, evt.RunSeq, evt.Text)
		case "tool_use":
			data, _ := json.Marshal(map[string]string{
				"tool_use_id": evt.ToolUseID,
				"tool_name":   evt.ToolName,
			})
			gw.BroadcastSession(evt.SessionID, protocol.EventFrame{
				Event:     protocol.EventChatToolUse,
				SessionID: evt.SessionID,
				RunSeq:    evt.RunSeq,
				Data:      data,
			})
		case "tool_done":
			data, _ := json.Marshal(map[string]string{
				"tool_use_id": evt.ToolUseID,
				"tool_name":   evt.ToolName,
			})
			gw.BroadcastSession(evt.SessionID, protocol.EventFrame{
				Event:     protocol.EventChatToolDone,
				SessionID: evt.SessionID,
				RunSeq:    evt.RunSeq,
				Data:      data,
			})
		case "complete":
			gw.ChatState.Flush(evt.SessionID)
			gw.BroadcastSession(evt.SessionID, protocol.EventFrame{
				Event:     protocol.EventChatComplete,
				SessionID: evt.SessionID,
				RunSeq:    evt.RunSeq,
			})
		case "error":
			errMsg := "unknown error"
			if evt.Err != nil {
				errMsg = evt.Err.Error()
			}
			data, _ := json.Marshal(map[string]string{"error": errMsg})
			gw.BroadcastSession(evt.SessionID, protocol.EventFrame{
				Event:     protocol.EventChatError,
				SessionID: evt.SessionID,
				RunSeq:    evt.RunSeq,
				Data:      data,
			})
		}
	})

	// processMessage is the unified message handler for both WS and channel messages.
	processMessage := func(ctx context.Context, msg protocol.InboundMessage) error {
		// Resolve agent.
		agentID := msg.AgentID
		if agentID == "" {
			agentID = router.Resolve(routing.ResolveParams{
				Channel:  msg.Channel,
				PeerKind: msg.PeerKind,
				PeerID:   msg.PeerID,
				GuildID:  msg.GuildID,
			})
		}
		msg.AgentID = agentID

		// Build session key.
		sessKey := session.Key(agentID, msg.Channel, msg.PeerKind, msg.PeerID, msg.GuildID, msg.ThreadID)
		if msg.SessionID != "" {
			sessKey = msg.SessionID
		}
		msg.SessionID = sessKey

		entry := sessionMgr.GetOrCreate(sessKey, agentID)

		// Enqueue task for serial execution.
		ok := queueMgr.Enqueue(queue.Task{
			SessionID: sessKey,
			Fn: func(ctx context.Context) error {
				runSeq := gw.RunState.Next(sessKey)

				// Broadcast run start.
				gw.BroadcastSession(sessKey, protocol.EventFrame{
					Event:     protocol.EventRunStart,
					SessionID: sessKey,
					RunSeq:    runSeq,
				})

				result, err := runtime.Run(ctx, agentID, entry, msg.Text, runSeq)
				if err != nil {
					return err
				}

				// Save to session history.
				entry.AppendHistory(session.Message{Role: "user", Content: msg.Text}, sessionMgr.MaxHistory())
				entry.AppendHistory(session.Message{Role: "assistant", Content: result.Text}, sessionMgr.MaxHistory())

				// Broadcast run end.
				gw.BroadcastSession(sessKey, protocol.EventFrame{
					Event:     protocol.EventRunEnd,
					SessionID: sessKey,
					RunSeq:    runSeq,
				})

				// Send reply back through the originating channel.
				if msg.Channel != "websocket" {
					return registry.SendMessage(ctx, protocol.OutboundMessage{
						SessionID: sessKey,
						Channel:   msg.Channel,
						PeerID:    msg.PeerID,
						ThreadID:  msg.ThreadID,
						Text:      result.Text,
						Format:    "markdown",
					})
				}
				return nil
			},
		})
		if !ok {
			return fmt.Errorf("queue full for session %s", sessKey)
		}
		return nil
	}

	// Wire WebSocket chat.send -> processMessage.
	gw.OnChatSend = func(ctx context.Context, clientID string, msg protocol.InboundMessage) error {
		return processMessage(ctx, msg)
	}

	// --- Telegram Adapter ---
	if cfg.Channels.Telegram.Enabled {
		tgCfg := telegram.ConfigFromApp(cfg.Channels.Telegram)
		bot, err := telegram.NewBot(tgCfg)
		if err != nil {
			slog.Error("failed to create telegram bot", "err", err)
			os.Exit(1)
		}
		bot.SetSink(func(ctx context.Context, msg protocol.InboundMessage) error {
			return processMessage(ctx, msg)
		})
		registry.Register(bot)
	}

	// --- Start ---
	if err := registry.StartAll(ctx); err != nil {
		slog.Error("failed to start channels", "err", err)
		os.Exit(1)
	}

	// Start gateway in a goroutine.
	go func() {
		if err := gw.Start(ctx); err != nil {
			slog.Error("gateway error", "err", err)
			cancel()
		}
	}()

	slog.Info("dhaavak started", "port", cfg.Server.Port)

	// --- Signal Handling ---
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	slog.Info("shutting down...")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	queueMgr.StopAll()
	registry.StopAll(shutdownCtx)
	if err := gw.Stop(shutdownCtx); err != nil {
		slog.Error("gateway shutdown error", "err", err)
	}

	slog.Info("dhaavak stopped")
}
