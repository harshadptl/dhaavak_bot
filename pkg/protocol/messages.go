package protocol

// InboundMessage represents a message arriving from any channel.
type InboundMessage struct {
	SessionID string `json:"session_id"`
	Channel   string `json:"channel"`   // e.g. "telegram", "websocket"
	PeerKind  string `json:"peer_kind"` // "user", "group", "channel"
	PeerID    string `json:"peer_id"`
	GuildID   string `json:"guild_id,omitempty"` // for group contexts
	ThreadID  string `json:"thread_id,omitempty"`
	Text      string `json:"text"`
	AgentID   string `json:"agent_id,omitempty"` // resolved by router
}

// OutboundMessage represents a message to be sent back to a channel.
type OutboundMessage struct {
	SessionID string `json:"session_id"`
	Channel   string `json:"channel"`
	PeerID    string `json:"peer_id"`
	ThreadID  string `json:"thread_id,omitempty"`
	Text      string `json:"text"`
	Format    string `json:"format,omitempty"` // "text", "markdown", "html"
}
