package agent

import "github.com/harshadpatil/dhaavak/internal/llm"

// Event represents an agent runtime event broadcast to observers.
type Event struct {
	Type      string // "run_start", "delta", "tool_use", "tool_done", "complete", "error"
	SessionID string
	RunSeq    int
	Text      string
	ToolUseID string
	ToolName  string
	ToolInput string
	Err       error
}

// RunResult is the final outcome of an agent run.
type RunResult struct {
	Text       string
	ToolCalls  int
	StopReason string
}

// ToolExecutor runs a tool and returns its output.
type ToolExecutor func(name, input string) (string, error)

// EventSink receives agent events for broadcasting.
type EventSink func(Event)

// AgentDef holds the static definition of an agent.
type AgentDef struct {
	ID           string
	Name         string
	SystemPrompt string
	Model        string
	Tools        []llm.ToolDef
}
