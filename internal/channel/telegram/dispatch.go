package telegram

import (
	"context"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/harshadpatil/dhaavak/internal/channel"
	"github.com/harshadpatil/dhaavak/internal/session"
)

// dispatch processes an incoming Telegram update.
func (b *Bot) dispatch(ctx context.Context, update tgbotapi.Update) {
	mc := extractContext(update, b.api.Self.UserName)
	if mc == nil {
		return
	}

	policy := &session.SendPolicy{
		DMPolicy:      b.cfg.DMPolicy,
		GroupPolicy:   b.cfg.GroupPolicy,
		AllowedUsers:  b.cfg.AllowedUsers,
		AllowedGroups: b.cfg.AllowedGroups,
	}

	if !checkAccess(mc, policy, b.cfg.GroupPolicy) {
		slog.Debug("telegram access denied", "user", mc.UserID, "chat", mc.ChatID)
		return
	}

	msg := toInboundMessage(mc)

	if b.sink != nil {
		if err := b.sink(ctx, msg); err != nil {
			slog.Error("telegram message sink error", "err", err)
		}
	}
}

// SetSink sets the message sink callback.
func (b *Bot) SetSink(sink channel.MessageSink) {
	b.sink = sink
}
