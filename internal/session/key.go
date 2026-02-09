package session

import (
	"fmt"
	"strings"
)

// Key builds a session key from its components.
//
// Formats:
//
//	agent:{agentID}:main                                   — default/WebSocket
//	agent:{agentID}:{channel}:{peerKind}:{peerID}         — DM
//	agent:{agentID}:{channel}:group:{guildID}              — group
//	agent:{agentID}:{channel}:group:{guildID}:{threadID}   — thread
func Key(agentID, channel, peerKind, peerID, guildID, threadID string) string {
	if channel == "" || channel == "websocket" {
		return fmt.Sprintf("agent:%s:main", agentID)
	}
	if guildID != "" {
		base := fmt.Sprintf("agent:%s:%s:group:%s", agentID, channel, guildID)
		if threadID != "" {
			return base + ":" + threadID
		}
		return base
	}
	return fmt.Sprintf("agent:%s:%s:%s:%s", agentID, channel, peerKind, peerID)
}

// ParsedKey holds the decomposed parts of a session key.
type ParsedKey struct {
	AgentID  string
	Channel  string
	PeerKind string
	PeerID   string
	GuildID  string
	ThreadID string
}

// ParseKey decomposes a session key string.
func ParseKey(key string) (ParsedKey, error) {
	parts := strings.Split(key, ":")
	if len(parts) < 3 || parts[0] != "agent" {
		return ParsedKey{}, fmt.Errorf("invalid session key: %s", key)
	}

	pk := ParsedKey{AgentID: parts[1]}

	if parts[2] == "main" {
		pk.Channel = "websocket"
		return pk, nil
	}

	pk.Channel = parts[2]
	if len(parts) < 5 {
		return ParsedKey{}, fmt.Errorf("invalid session key: %s", key)
	}
	pk.PeerKind = parts[3]
	pk.PeerID = parts[4]

	if pk.PeerKind == "group" {
		pk.GuildID = parts[4]
		pk.PeerID = ""
		if len(parts) >= 6 {
			pk.ThreadID = parts[5]
		}
	}

	return pk, nil
}
