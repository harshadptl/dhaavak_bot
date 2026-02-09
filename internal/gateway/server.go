package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/harshadpatil/dhaavak/internal/config"
	"github.com/harshadpatil/dhaavak/pkg/protocol"
)

// MessageHandler is called when a chat.send request arrives via WebSocket.
type MessageHandler func(ctx context.Context, clientID string, msg protocol.InboundMessage) error

// Server is the WebSocket gateway.
type Server struct {
	cfg          config.ServerConfig
	auth         *Authenticator
	clients      map[string]*Client
	mu           sync.RWMutex
	httpServer   *http.Server
	RunState     *RunState
	ChatState    *ChatRunState
	OnChatSend   MessageHandler
}

// New creates a new gateway server.
func New(cfg config.ServerConfig, authToken string) *Server {
	s := &Server{
		cfg:      cfg,
		auth:     NewAuthenticator(authToken),
		clients:  make(map[string]*Client),
		RunState: NewRunState(),
	}
	s.ChatState = NewChatRunState(s)
	return s
}

// Start begins listening for HTTP/WebSocket connections.
func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.handleWebSocket)
	mux.HandleFunc("/health", s.handleHealth)

	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	s.httpServer = &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		BaseContext:       func(l net.Listener) context.Context { return ctx },
	}

	slog.Info("gateway listening", "addr", addr)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("gateway listen: %w", err)
	}
	return nil
}

// Stop gracefully shuts down the gateway.
func (s *Server) Stop(ctx context.Context) error {
	s.mu.Lock()
	for _, c := range s.clients {
		close(c.sendCh)
	}
	s.clients = make(map[string]*Client)
	s.mu.Unlock()

	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	if !s.auth.Check(r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // Allow any origin for dev; configure in production.
	})
	if err != nil {
		slog.Error("websocket accept error", "err", err)
		return
	}

	ctx, cancel := context.WithCancel(r.Context())
	client := newClient(conn, s)
	client.cancelCtx = cancel

	s.register(client)

	// Send connected event.
	client.sendJSON(protocol.EventFrame{
		Event: protocol.EventConnected,
		Data:  mustJSON(map[string]string{"client_id": client.ID}),
	})

	go client.writePump(ctx)
	client.readPump(ctx) // blocks until disconnect
}

func (s *Server) register(c *Client) {
	s.mu.Lock()
	s.clients[c.ID] = c
	s.mu.Unlock()
	slog.Info("client connected", "client", c.ID)
}

func (s *Server) unregister(c *Client) {
	s.mu.Lock()
	if _, ok := s.clients[c.ID]; ok {
		close(c.sendCh)
		delete(s.clients, c.ID)
	}
	s.mu.Unlock()
	if c.cancelCtx != nil {
		c.cancelCtx()
	}
	slog.Info("client disconnected", "client", c.ID)
}

func (s *Server) handleMessage(c *Client, data []byte) {
	var req protocol.RequestFrame
	if err := json.Unmarshal(data, &req); err != nil {
		c.sendJSON(protocol.ResponseFrame{
			Error: &protocol.ErrorDetail{Code: 400, Message: "invalid request frame"},
		})
		return
	}

	switch req.Method {
	case protocol.MethodPing:
		c.sendJSON(protocol.ResponseFrame{
			ID:     req.ID,
			Result: mustJSON(map[string]string{"pong": "ok"}),
		})

	case protocol.MethodChatSend:
		s.handleChatSend(c, req)

	default:
		c.sendJSON(protocol.ResponseFrame{
			ID:    req.ID,
			Error: &protocol.ErrorDetail{Code: 404, Message: "unknown method: " + req.Method},
		})
	}
}

func (s *Server) handleChatSend(c *Client, req protocol.RequestFrame) {
	var params struct {
		SessionID string `json:"session_id"`
		Text      string `json:"text"`
		AgentID   string `json:"agent_id"`
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		c.sendJSON(protocol.ResponseFrame{
			ID:    req.ID,
			Error: &protocol.ErrorDetail{Code: 400, Message: "invalid chat.send params"},
		})
		return
	}

	c.Subscribe(params.SessionID)

	// Acknowledge the request.
	c.sendJSON(protocol.ResponseFrame{
		ID:     req.ID,
		Result: mustJSON(map[string]string{"status": "queued"}),
	})

	msg := protocol.InboundMessage{
		SessionID: params.SessionID,
		Channel:   "websocket",
		PeerKind:  "user",
		PeerID:    c.ID,
		Text:      params.Text,
		AgentID:   params.AgentID,
	}

	if s.OnChatSend != nil {
		go func() {
			if err := s.OnChatSend(context.Background(), c.ID, msg); err != nil {
				slog.Error("chat.send handler error", "err", err)
			}
		}()
	}
}

func mustJSON(v interface{}) json.RawMessage {
	data, _ := json.Marshal(v)
	return data
}
