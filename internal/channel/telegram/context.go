package telegram

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/harshadpatil/dhaavak/internal/session"
	"github.com/harshadpatil/dhaavak/pkg/protocol"
)

// messageContext extracts routing information from a Telegram update.
type messageContext struct {
	ChatID    int64
	ThreadID  int
	UserID    int64
	Text      string
	PeerKind  string // "user" or "group"
	PeerID    string
	GuildID   string
	IsMention bool
}

func extractContext(update tgbotapi.Update, botUsername string) *messageContext {
	msg := update.Message
	if msg == nil {
		return nil
	}

	text := msg.Text
	if text == "" {
		text = msg.Caption
	}
	if text == "" {
		return nil
	}

	mc := &messageContext{
		ChatID: msg.Chat.ID,
		UserID: msg.From.ID,
		Text:   text,
	}

	if msg.Chat.IsPrivate() {
		mc.PeerKind = "user"
		mc.PeerID = fmt.Sprintf("%d", msg.From.ID)
	} else {
		mc.PeerKind = "group"
		mc.PeerID = fmt.Sprintf("%d", msg.Chat.ID)
		mc.GuildID = fmt.Sprintf("%d", msg.Chat.ID)
	}

	// Check for bot mention.
	if botUsername != "" {
		mention := "@" + botUsername
		if strings.Contains(text, mention) {
			mc.IsMention = true
			mc.Text = strings.TrimSpace(strings.ReplaceAll(text, mention, ""))
		}
	}

	// Check entities for mention.
	if msg.Entities != nil {
		for _, e := range msg.Entities {
			if e.Type == "mention" {
				mentioned := text[e.Offset : e.Offset+e.Length]
				if strings.EqualFold(mentioned, "@"+botUsername) {
					mc.IsMention = true
				}
			}
		}
	}

	return mc
}

// checkAccess verifies whether this message should be processed.
func checkAccess(mc *messageContext, policy *session.SendPolicy, groupPolicy string) bool {
	if mc.PeerKind == "user" {
		return policy.AllowDM(mc.UserID)
	}
	// Group message.
	if !policy.AllowGroup(mc.ChatID) {
		return false
	}
	if groupPolicy == "mention" && !mc.IsMention {
		return false
	}
	return true
}

// toInboundMessage converts a message context to a protocol message.
func toInboundMessage(mc *messageContext) protocol.InboundMessage {
	threadID := ""
	if mc.ThreadID != 0 {
		threadID = fmt.Sprintf("%d", mc.ThreadID)
	}
	return protocol.InboundMessage{
		Channel:  "telegram",
		PeerKind: mc.PeerKind,
		PeerID:   mc.PeerID,
		GuildID:  mc.GuildID,
		ThreadID: threadID,
		Text:     mc.Text,
	}
}
