package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var envPattern = regexp.MustCompile(`\$\{([^}]+)\}`)

// Load reads a YAML config file, applies env substitution, and returns Config.
func Load(path string) (*Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}

	// Substitute ${ENV_VAR} with environment values.
	expanded := envPattern.ReplaceAllStringFunc(string(raw), func(m string) string {
		key := envPattern.FindStringSubmatch(m)[1]
		parts := strings.SplitN(key, ":", 2)
		if val, ok := os.LookupEnv(parts[0]); ok {
			return val
		}
		if len(parts) == 2 {
			return parts[1] // default value after colon
		}
		return m
	})

	k := koanf.New(".")
	if err := k.Load(file.Provider(path), yaml.Parser(), koanf.WithMergeFunc(func(src, dest map[string]interface{}) error {
		// We won't use the file provider directly; we'll parse the expanded YAML.
		return nil
	})); err != nil {
		// Ignore file provider error since we parse manually below.
		_ = err
	}

	// Parse expanded YAML manually.
	k = koanf.New(".")
	if err := k.Load(rawBytesProvider([]byte(expanded)), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	cfg := Default()

	// Server
	if k.Exists("server.port") {
		cfg.Server.Port = k.Int("server.port")
	}
	if k.Exists("server.host") {
		cfg.Server.Host = k.String("server.host")
	}

	// Auth
	cfg.Auth.Token = k.String("auth.token")

	// LLM
	if k.Exists("llm.provider") {
		cfg.LLM.Provider = k.String("llm.provider")
	}
	if k.Exists("llm.api_key") {
		cfg.LLM.APIKey = k.String("llm.api_key")
	}
	if k.Exists("llm.model") {
		cfg.LLM.Model = k.String("llm.model")
	}
	if k.Exists("llm.max_turns") {
		cfg.LLM.MaxTurns = k.Int("llm.max_turns")
	}

	// Agents
	if k.Exists("agents") {
		var agents []AgentConfig
		for _, raw := range k.Slices("agents") {
			agents = append(agents, AgentConfig{
				ID:           raw.String("id"),
				Name:         raw.String("name"),
				SystemPrompt: raw.String("system_prompt"),
				Model:        raw.String("model"),
			})
		}
		cfg.Agents = agents
	}

	// Channels - Telegram
	if k.Exists("channels.telegram") {
		tg := &cfg.Channels.Telegram
		if k.Exists("channels.telegram.enabled") {
			tg.Enabled = k.Bool("channels.telegram.enabled")
		}
		tg.BotToken = k.String("channels.telegram.bot_token")
		if k.Exists("channels.telegram.default_agent") {
			tg.DefaultAgent = k.String("channels.telegram.default_agent")
		}
		if k.Exists("channels.telegram.dm_policy") {
			tg.DMPolicy = k.String("channels.telegram.dm_policy")
		}
		if k.Exists("channels.telegram.group_policy") {
			tg.GroupPolicy = k.String("channels.telegram.group_policy")
		}
		if k.Exists("channels.telegram.allowed_users") {
			tg.AllowedUsers = k.Int64s("channels.telegram.allowed_users")
		}
		if k.Exists("channels.telegram.allowed_groups") {
			tg.AllowedGroups = k.Int64s("channels.telegram.allowed_groups")
		}
	}

	// Session
	if k.Exists("session.ttl") {
		cfg.Session.TTL = k.Duration("session.ttl")
	}
	if k.Exists("session.cleanup_interval") {
		cfg.Session.CleanupInterval = k.Duration("session.cleanup_interval")
	}
	if k.Exists("session.max_history") {
		cfg.Session.MaxHistory = k.Int("session.max_history")
	}

	// Queue
	if k.Exists("queue.buffer_size") {
		cfg.Queue.BufferSize = k.Int("queue.buffer_size")
	}
	if k.Exists("queue.idle_timeout") {
		cfg.Queue.IdleTimeout = k.Duration("queue.idle_timeout")
	}
	if k.Exists("queue.cleanup_interval") {
		cfg.Queue.CleanupInterval = k.Duration("queue.cleanup_interval")
	}

	if err := validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func validate(cfg *Config) error {
	if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
		return fmt.Errorf("config: server.port must be 1-65535, got %d", cfg.Server.Port)
	}
	if cfg.LLM.APIKey == "" {
		return fmt.Errorf("config: llm.api_key is required")
	}
	if len(cfg.Agents) == 0 {
		return fmt.Errorf("config: at least one agent must be defined")
	}
	if cfg.Channels.Telegram.Enabled && cfg.Channels.Telegram.BotToken == "" {
		return fmt.Errorf("config: channels.telegram.bot_token is required when telegram is enabled")
	}
	return nil
}

// rawBytesProvider implements koanf.Provider for raw bytes.
type rawBytesProvider []byte

func (r rawBytesProvider) ReadBytes() ([]byte, error) { return r, nil }
func (r rawBytesProvider) Read() (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}
