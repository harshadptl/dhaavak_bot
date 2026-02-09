package agent

import (
	"github.com/harshadpatil/dhaavak/internal/llm"
	"github.com/harshadpatil/dhaavak/internal/session"
)

// BuildMessages converts session history into LLM messages and appends the new user message.
func BuildMessages(entry *session.Entry, userText string) []llm.Message {
	history := entry.GetHistory()

	var msgs []llm.Message
	for _, h := range history {
		msgs = append(msgs, llm.Message{
			Role: h.Role,
			Content: []llm.ContentBlock{
				{Type: "text", Text: h.Content},
			},
		})
	}

	msgs = append(msgs, llm.Message{
		Role: llm.RoleUser,
		Content: []llm.ContentBlock{
			{Type: "text", Text: userText},
		},
	})

	return msgs
}
