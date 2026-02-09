package agent

import (
	"context"
	"fmt"
	"strings"

	"github.com/harshadpatil/dhaavak/internal/llm"
)

// RunLoop executes the agentic loop: call LLM, execute tools if needed, repeat.
func RunLoop(
	ctx context.Context,
	provider llm.Provider,
	systemPrompt string,
	messages []llm.Message,
	tools []llm.ToolDef,
	toolExec ToolExecutor,
	sink EventSink,
	sessionID string,
	runSeq int,
	maxTurns int,
) (*RunResult, []llm.Message, error) {
	result := &RunResult{}

	for turn := 0; turn < maxTurns; turn++ {
		select {
		case <-ctx.Done():
			return nil, messages, ctx.Err()
		default:
		}

		stream, err := provider.Stream(ctx, systemPrompt, messages, tools)
		if err != nil {
			return nil, messages, fmt.Errorf("llm stream: %w", err)
		}

		var textBuf strings.Builder
		var toolCalls []llm.ContentBlock

		for evt := range stream {
			if sink != nil {
				sink(MapStreamEvent(evt, sessionID, runSeq))
			}

			switch evt.Type {
			case "delta":
				textBuf.WriteString(evt.Text)
			case "tool_done":
				toolCalls = append(toolCalls, llm.ContentBlock{
					Type:  "tool_use",
					ID:    evt.ToolUseID,
					Name:  evt.ToolName,
					Input: evt.ToolInput,
				})
			case "error":
				return nil, messages, fmt.Errorf("stream error: %w", evt.Err)
			case "complete":
				result.StopReason = evt.StopReason
			}
		}

		// Build the assistant message from this turn.
		var assistantBlocks []llm.ContentBlock
		if textBuf.Len() > 0 {
			assistantBlocks = append(assistantBlocks, llm.ContentBlock{
				Type: "text",
				Text: textBuf.String(),
			})
		}
		assistantBlocks = append(assistantBlocks, toolCalls...)

		messages = append(messages, llm.Message{
			Role:    llm.RoleAssistant,
			Content: assistantBlocks,
		})

		// If no tool calls, we're done.
		if len(toolCalls) == 0 {
			result.Text = textBuf.String()
			return result, messages, nil
		}

		// Execute tools and add results.
		result.ToolCalls += len(toolCalls)
		var toolResults []llm.ContentBlock
		for _, tc := range toolCalls {
			output, execErr := toolExec(tc.Name, tc.Input)
			if execErr != nil {
				output = fmt.Sprintf("Error: %s", execErr.Error())
			}
			toolResults = append(toolResults, llm.ContentBlock{
				Type: "tool_result",
				ID:   tc.ID,
				Text: output,
			})
		}

		messages = append(messages, llm.Message{
			Role:    llm.RoleUser,
			Content: toolResults,
		})
	}

	return result, messages, fmt.Errorf("max turns (%d) exceeded", maxTurns)
}
