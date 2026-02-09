package telegram

import (
	"github.com/harshadpatil/dhaavak/internal/channel"
	"github.com/harshadpatil/dhaavak/internal/config"
)

// BotConfig holds Telegram adapter configuration.
type BotConfig struct {
	Token         string
	DefaultAgent  string
	DMPolicy      string
	GroupPolicy   string
	AllowedUsers  []int64
	AllowedGroups []int64
}

// ConfigFromApp extracts Telegram config from the app config.
func ConfigFromApp(cfg config.TelegramConfig) BotConfig {
	return BotConfig{
		Token:         cfg.BotToken,
		DefaultAgent:  cfg.DefaultAgent,
		DMPolicy:      cfg.DMPolicy,
		GroupPolicy:   cfg.GroupPolicy,
		AllowedUsers:  cfg.AllowedUsers,
		AllowedGroups: cfg.AllowedGroups,
	}
}

// ensure Bot implements Adapter.
var _ channel.Adapter = (*Bot)(nil)
