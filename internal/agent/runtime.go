package agent

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/harshadpatil/dhaavak/internal/llm"
	"github.com/harshadpatil/dhaavak/internal/session"
)

// Runtime orchestrates agent execution.
type Runtime struct {
	provider  llm.Provider
	agents    map[string]AgentDef
	maxTurns  int
	toolExec  ToolExecutor
	eventSink EventSink
}

// NewRuntime creates an agent runtime.
func NewRuntime(provider llm.Provider, maxTurns int) *Runtime {
	return &Runtime{
		provider: provider,
		agents:   make(map[string]AgentDef),
		maxTurns: maxTurns,
	}
}

// RegisterAgent adds an agent definition.
func (rt *Runtime) RegisterAgent(def AgentDef) {
	rt.agents[def.ID] = def
}

// SetToolExecutor configures the tool executor callback.
func (rt *Runtime) SetToolExecutor(fn ToolExecutor) {
	rt.toolExec = fn
}

// SetEventSink configures where agent events are sent.
func (rt *Runtime) SetEventSink(sink EventSink) {
	rt.eventSink = sink
}

// Run executes an agent for a given session and user message.
func (rt *Runtime) Run(ctx context.Context, agentID string, entry *session.Entry, userText string, runSeq int) (*RunResult, error) {
	def, ok := rt.agents[agentID]
	if !ok {
		return nil, fmt.Errorf("agent not found: %s", agentID)
	}

	slog.Info("agent run start", "agent", agentID, "session", entry.Key, "run_seq", runSeq)

	messages := BuildMessages(entry, userText)

	toolExec := rt.toolExec
	if toolExec == nil {
		toolExec = func(name, input string) (string, error) {
			return "", fmt.Errorf("no tool executor configured for tool: %s", name)
		}
	}

	result, _, err := RunLoop(
		ctx,
		rt.provider,
		def.SystemPrompt,
		messages,
		def.Tools,
		toolExec,
		rt.eventSink,
		entry.Key,
		runSeq,
		rt.maxTurns,
	)
	if err != nil {
		slog.Error("agent run error", "agent", agentID, "session", entry.Key, "err", err)
		return nil, err
	}

	slog.Info("agent run complete", "agent", agentID, "session", entry.Key, "tool_calls", result.ToolCalls)
	return result, nil
}
