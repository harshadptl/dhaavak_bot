package llm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// AnthropicProvider implements Provider using the Anthropic SDK.
type AnthropicProvider struct {
	client *anthropic.Client
	model  anthropic.Model
}

// NewAnthropicProvider creates a provider for Claude models.
func NewAnthropicProvider(apiKey, model string) *AnthropicProvider {
	client := anthropic.NewClient(option.WithAPIKey(apiKey))
	return &AnthropicProvider{
		client: &client,
		model:  anthropic.Model(model),
	}
}

func (a *AnthropicProvider) Complete(ctx context.Context, systemPrompt string, messages []Message, tools []ToolDef) (*CompletionResult, error) {
	params := a.buildParams(systemPrompt, messages, tools)

	resp, err := a.client.Messages.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("anthropic complete: %w", err)
	}

	var blocks []ContentBlock
	for _, b := range resp.Content {
		switch b.Type {
		case "text":
			blocks = append(blocks, ContentBlock{Type: "text", Text: b.Text})
		case "tool_use":
			blocks = append(blocks, ContentBlock{
				Type:  "tool_use",
				ID:    b.ID,
				Name:  b.Name,
				Input: string(b.Input),
			})
		}
	}

	return &CompletionResult{
		Content:    blocks,
		StopReason: string(resp.StopReason),
	}, nil
}

func (a *AnthropicProvider) Stream(ctx context.Context, systemPrompt string, messages []Message, tools []ToolDef) (<-chan StreamEvent, error) {
	params := a.buildParams(systemPrompt, messages, tools)

	stream := a.client.Messages.NewStreaming(ctx, params)

	ch := make(chan StreamEvent, 64)
	go func() {
		defer close(ch)
		defer stream.Close()

		// Use Accumulate to track the full message for stop reason.
		accumulated := &anthropic.Message{}
		var currentToolID, currentToolName string
		var toolInputBuf string

		for stream.Next() {
			evt := stream.Current()
			accumulated.Accumulate(evt)

			switch evt.Type {
			case "content_block_start":
				if evt.ContentBlock.Type == "tool_use" {
					currentToolID = evt.ContentBlock.ID
					currentToolName = evt.ContentBlock.Name
					toolInputBuf = ""
					ch <- StreamEvent{
						Type:      "tool_use",
						ToolUseID: currentToolID,
						ToolName:  currentToolName,
					}
				}

			case "content_block_delta":
				if evt.Delta.Type == "text_delta" {
					ch <- StreamEvent{Type: "delta", Text: evt.Delta.Text}
				} else if evt.Delta.Type == "input_json_delta" {
					toolInputBuf += evt.Delta.PartialJSON
				}

			case "content_block_stop":
				if currentToolID != "" {
					ch <- StreamEvent{
						Type:      "tool_done",
						ToolUseID: currentToolID,
						ToolName:  currentToolName,
						ToolInput: toolInputBuf,
					}
					currentToolID = ""
					currentToolName = ""
					toolInputBuf = ""
				}

			case "message_stop":
				ch <- StreamEvent{
					Type:       "complete",
					StopReason: string(accumulated.StopReason),
				}
			}
		}

		if err := stream.Err(); err != nil {
			ch <- StreamEvent{Type: "error", Err: err}
		}
	}()

	return ch, nil
}

func (a *AnthropicProvider) buildParams(systemPrompt string, messages []Message, tools []ToolDef) anthropic.MessageNewParams {
	params := anthropic.MessageNewParams{
		Model:     a.model,
		MaxTokens: 8192,
	}

	if systemPrompt != "" {
		params.System = []anthropic.TextBlockParam{
			{Text: systemPrompt},
		}
	}

	for _, m := range messages {
		var blocks []anthropic.ContentBlockParamUnion
		for _, b := range m.Content {
			switch b.Type {
			case "text":
				blocks = append(blocks, anthropic.NewTextBlock(b.Text))
			case "tool_use":
				var input interface{}
				json.Unmarshal([]byte(b.Input), &input)
				blocks = append(blocks, anthropic.NewToolUseBlock(b.ID, input, b.Name))
			case "tool_result":
				blocks = append(blocks, anthropic.NewToolResultBlock(b.ID, b.Text, false))
			}
		}

		switch m.Role {
		case RoleUser:
			params.Messages = append(params.Messages, anthropic.NewUserMessage(blocks...))
		case RoleAssistant:
			params.Messages = append(params.Messages, anthropic.NewAssistantMessage(blocks...))
		}
	}

	for _, t := range tools {
		params.Tools = append(params.Tools, anthropic.ToolUnionParam{
			OfTool: &anthropic.ToolParam{
				Name:        t.Name,
				Description: anthropic.String(t.Description),
				InputSchema: anthropic.ToolInputSchemaParam{
					Properties: t.InputSchema,
				},
			},
		})
	}

	return params
}
