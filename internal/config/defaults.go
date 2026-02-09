package config

import "time"

// Default returns a Config populated with default values.
func Default() Config {
	return Config{
		Server: ServerConfig{
			Port: 18789,
			Host: "127.0.0.1",
		},
		LLM: LLMConfig{
			Provider: "anthropic",
			Model:    "claude-sonnet-4-5-20250929",
			MaxTurns: 25,
		},
		Channels: ChannelsConfig{
			Telegram: TelegramConfig{
				DMPolicy:    "open",
				GroupPolicy: "mention",
			},
		},
		Session: SessionConfig{
			TTL:             30 * time.Minute,
			CleanupInterval: 5 * time.Minute,
			MaxHistory:      100,
		},
		Queue: QueueConfig{
			BufferSize:      64,
			IdleTimeout:     10 * time.Minute,
			CleanupInterval: 2 * time.Minute,
		},
	}
}
