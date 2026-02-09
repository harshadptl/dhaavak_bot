package config

import "time"

// Config is the top-level application configuration.
type Config struct {
	Server   ServerConfig   `json:"server"   yaml:"server"`
	Auth     AuthConfig     `json:"auth"     yaml:"auth"`
	LLM      LLMConfig      `json:"llm"      yaml:"llm"`
	Agents   []AgentConfig  `json:"agents"   yaml:"agents"`
	Channels ChannelsConfig `json:"channels" yaml:"channels"`
	Session  SessionConfig  `json:"session"  yaml:"session"`
	Queue    QueueConfig    `json:"queue"    yaml:"queue"`
}

type ServerConfig struct {
	Port int    `json:"port" yaml:"port"`
	Host string `json:"host" yaml:"host"`
}

type AuthConfig struct {
	Token string `json:"token" yaml:"token"`
}

type LLMConfig struct {
	Provider string `json:"provider" yaml:"provider"`
	APIKey   string `json:"api_key"  yaml:"api_key"`
	Model    string `json:"model"    yaml:"model"`
	MaxTurns int    `json:"max_turns" yaml:"max_turns"`
}

type AgentConfig struct {
	ID           string       `json:"id"            yaml:"id"`
	Name         string       `json:"name"          yaml:"name"`
	SystemPrompt string       `json:"system_prompt" yaml:"system_prompt"`
	Model        string       `json:"model,omitempty" yaml:"model,omitempty"`
	Tools        []ToolConfig `json:"tools,omitempty" yaml:"tools,omitempty"`
}

type ToolConfig struct {
	Name        string `json:"name"        yaml:"name"`
	Description string `json:"description" yaml:"description"`
}

type ChannelsConfig struct {
	Telegram TelegramConfig `json:"telegram" yaml:"telegram"`
}

type TelegramConfig struct {
	Enabled       bool          `json:"enabled"        yaml:"enabled"`
	BotToken      string        `json:"bot_token"      yaml:"bot_token"`
	DefaultAgent  string        `json:"default_agent"  yaml:"default_agent"`
	DMPolicy      string        `json:"dm_policy"      yaml:"dm_policy"`      // "open", "allowlist", "disabled"
	GroupPolicy   string        `json:"group_policy"   yaml:"group_policy"`   // "mention", "all", "disabled"
	AllowedUsers  []int64       `json:"allowed_users"  yaml:"allowed_users"`
	AllowedGroups []int64       `json:"allowed_groups" yaml:"allowed_groups"`
	Bindings      []BindingRule `json:"bindings"       yaml:"bindings"`
}

type BindingRule struct {
	PeerKind string `json:"peer_kind" yaml:"peer_kind"` // "user", "group"
	PeerID   string `json:"peer_id"   yaml:"peer_id"`
	AgentID  string `json:"agent_id"  yaml:"agent_id"`
}

type SessionConfig struct {
	TTL             time.Duration `json:"ttl"               yaml:"ttl"`
	CleanupInterval time.Duration `json:"cleanup_interval"  yaml:"cleanup_interval"`
	MaxHistory      int           `json:"max_history"       yaml:"max_history"`
}

type QueueConfig struct {
	BufferSize      int           `json:"buffer_size"       yaml:"buffer_size"`
	IdleTimeout     time.Duration `json:"idle_timeout"      yaml:"idle_timeout"`
	CleanupInterval time.Duration `json:"cleanup_interval"  yaml:"cleanup_interval"`
}
