package llm

import "context"

// Provider is the interface for LLM backends.
type Provider interface {
	// Stream sends messages to the LLM and returns a channel of streaming events.
	Stream(ctx context.Context, systemPrompt string, messages []Message, tools []ToolDef) (<-chan StreamEvent, error)

	// Complete sends messages and returns the full response (non-streaming).
	Complete(ctx context.Context, systemPrompt string, messages []Message, tools []ToolDef) (*CompletionResult, error)
}
