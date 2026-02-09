package telegram

import (
	"log/slog"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const maxChunkSize = 4000

// sendText sends text back to a Telegram chat, chunking if necessary.
func sendText(api *tgbotapi.BotAPI, chatID int64, threadID string, text string) {
	tid := 0
	if threadID != "" {
		tid, _ = strconv.Atoi(threadID)
	}

	chunks := chunkText(text, maxChunkSize)
	for _, chunk := range chunks {
		msg := tgbotapi.NewMessage(chatID, chunk)
		msg.ParseMode = "HTML"
		msg.ReplyToMessageID = tid

		if _, err := api.Send(msg); err != nil {
			// Retry without HTML parse mode if it fails.
			slog.Warn("telegram send HTML failed, retrying plain", "err", err)
			msg.ParseMode = ""
			if _, err := api.Send(msg); err != nil {
				slog.Error("telegram send error", "chat_id", chatID, "err", err)
			}
		}
	}
}

// markdownToHTML converts basic markdown to Telegram HTML.
func markdownToHTML(text string) string {
	// Simple conversions for common patterns.
	r := strings.NewReplacer(
		"**", "<b>", // Bold (will need pairing)
		"__", "<i>", // Italic
		"```", "<pre>", // Code blocks
		"`", "<code>",
	)
	// This is a simplified conversion. For production, use a proper parser.
	// For now, just escape HTML entities and pass through.
	text = strings.ReplaceAll(text, "&", "&amp;")
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")
	_ = r // Skip the replacer since we're escaping HTML entities.
	return text
}

// chunkText splits text into chunks of at most maxSize bytes,
// trying to break at newlines.
func chunkText(text string, maxSize int) []string {
	if len(text) <= maxSize {
		return []string{text}
	}

	var chunks []string
	for len(text) > 0 {
		if len(text) <= maxSize {
			chunks = append(chunks, text)
			break
		}

		// Try to break at a newline within the limit.
		cutAt := maxSize
		if idx := strings.LastIndex(text[:maxSize], "\n"); idx > 0 {
			cutAt = idx + 1
		}

		chunks = append(chunks, text[:cutAt])
		text = text[cutAt:]
	}

	return chunks
}
