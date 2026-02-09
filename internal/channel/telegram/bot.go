package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/harshadpatil/dhaavak/internal/channel"
	"github.com/harshadpatil/dhaavak/pkg/protocol"
)

// Bot is the Telegram adapter.
type Bot struct {
	cfg    BotConfig
	api    *tgbotapi.BotAPI
	sink   channel.MessageSink
	cancel context.CancelFunc
}

// NewBot creates a Telegram bot adapter.
func NewBot(cfg BotConfig) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, fmt.Errorf("telegram bot init: %w", err)
	}
	slog.Info("telegram bot authorized", "username", api.Self.UserName)

	return &Bot{
		cfg: cfg,
		api: api,
	}, nil
}

func (b *Bot) ID() string { return "telegram" }

func (b *Bot) Start(ctx context.Context) error {
	ctx, b.cancel = context.WithCancel(ctx)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := b.api.GetUpdatesChan(u)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case update, ok := <-updates:
				if !ok {
					return
				}
				b.dispatch(ctx, update)
			}
		}
	}()

	slog.Info("telegram polling started")
	return nil
}

func (b *Bot) Stop(_ context.Context) error {
	if b.cancel != nil {
		b.cancel()
	}
	b.api.StopReceivingUpdates()
	slog.Info("telegram bot stopped")
	return nil
}

func (b *Bot) SendMessage(_ context.Context, msg protocol.OutboundMessage) error {
	chatID, err := strconv.ParseInt(msg.PeerID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid telegram chat ID: %s", msg.PeerID)
	}

	text := msg.Text
	if msg.Format == "markdown" {
		text = markdownToHTML(text)
	}

	sendText(b.api, chatID, msg.ThreadID, text)
	return nil
}
