package llm

// Role constants.
const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
)

// Message represents a conversation turn.
type Message struct {
	Role    string         `json:"role"`
	Content []ContentBlock `json:"content"`
}

// ContentBlock is a piece of message content.
type ContentBlock struct {
	Type  string `json:"type"` // "text", "tool_use", "tool_result"
	Text  string `json:"text,omitempty"`
	ID    string `json:"id,omitempty"`    // tool_use ID
	Name  string `json:"name,omitempty"`  // tool name
	Input string `json:"input,omitempty"` // tool input JSON
}

// ToolDef defines a tool the LLM can call.
type ToolDef struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"input_schema"`
}

// StreamEvent is emitted during streaming.
type StreamEvent struct {
	Type string // "delta", "tool_use", "tool_done", "complete", "error"

	// Delta fields
	Text string

	// Tool use fields
	ToolUseID string
	ToolName  string
	ToolInput string

	// Complete fields
	StopReason string

	// Error fields
	Err error
}

// CompletionResult is the outcome of a non-streaming call.
type CompletionResult struct {
	Content    []ContentBlock
	StopReason string
}
