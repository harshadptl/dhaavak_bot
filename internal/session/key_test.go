package session

import "testing"

func TestKey(t *testing.T) {
	tests := []struct {
		name     string
		agentID  string
		channel  string
		peerKind string
		peerID   string
		guildID  string
		threadID string
		want     string
	}{
		{"websocket default", "default", "websocket", "", "", "", "", "agent:default:main"},
		{"empty channel", "default", "", "", "", "", "", "agent:default:main"},
		{"telegram DM", "default", "telegram", "user", "12345", "", "", "agent:default:telegram:user:12345"},
		{"telegram group", "default", "telegram", "group", "99", "99", "", "agent:default:telegram:group:99"},
		{"telegram thread", "default", "telegram", "group", "99", "99", "42", "agent:default:telegram:group:99:42"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Key(tt.agentID, tt.channel, tt.peerKind, tt.peerID, tt.guildID, tt.threadID)
			if got != tt.want {
				t.Errorf("Key() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseKey(t *testing.T) {
	tests := []struct {
		key     string
		wantErr bool
		check   func(ParsedKey) bool
	}{
		{"agent:default:main", false, func(pk ParsedKey) bool {
			return pk.AgentID == "default" && pk.Channel == "websocket"
		}},
		{"agent:bot1:telegram:user:12345", false, func(pk ParsedKey) bool {
			return pk.AgentID == "bot1" && pk.Channel == "telegram" && pk.PeerKind == "user" && pk.PeerID == "12345"
		}},
		{"agent:bot1:telegram:group:99", false, func(pk ParsedKey) bool {
			return pk.AgentID == "bot1" && pk.GuildID == "99" && pk.PeerKind == "group"
		}},
		{"agent:bot1:telegram:group:99:42", false, func(pk ParsedKey) bool {
			return pk.ThreadID == "42" && pk.GuildID == "99"
		}},
		{"invalid", true, nil},
		{"bad:key", true, nil},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			pk, err := ParseKey(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseKey(%q) error = %v, wantErr %v", tt.key, err, tt.wantErr)
				return
			}
			if tt.check != nil && !tt.check(pk) {
				t.Errorf("ParseKey(%q) = %+v, check failed", tt.key, pk)
			}
		})
	}
}
