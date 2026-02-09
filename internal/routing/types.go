package routing

// Route holds a resolved agent binding.
type Route struct {
	AgentID  string
	Priority int // lower = higher priority
	Source   string
}

// Binding maps a channel context to an agent.
type Binding struct {
	Channel  string `json:"channel"   yaml:"channel"`
	PeerKind string `json:"peer_kind" yaml:"peer_kind"` // "user", "group", ""
	PeerID   string `json:"peer_id"   yaml:"peer_id"`
	GuildID  string `json:"guild_id"  yaml:"guild_id"`
	TeamID   string `json:"team_id"   yaml:"team_id"`
	AgentID  string `json:"agent_id"  yaml:"agent_id"`
}
