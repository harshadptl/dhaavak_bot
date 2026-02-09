package agent

import "github.com/harshadpatil/dhaavak/internal/llm"

// MapStreamEvent converts an LLM stream event to an agent event.
func MapStreamEvent(evt llm.StreamEvent, sessionID string, runSeq int) Event {
	base := Event{
		SessionID: sessionID,
		RunSeq:    runSeq,
	}

	switch evt.Type {
	case "delta":
		base.Type = "delta"
		base.Text = evt.Text
	case "tool_use":
		base.Type = "tool_use"
		base.ToolUseID = evt.ToolUseID
		base.ToolName = evt.ToolName
	case "tool_done":
		base.Type = "tool_done"
		base.ToolUseID = evt.ToolUseID
		base.ToolName = evt.ToolName
		base.ToolInput = evt.ToolInput
	case "complete":
		base.Type = "complete"
	case "error":
		base.Type = "error"
		base.Err = evt.Err
	}

	return base
}
