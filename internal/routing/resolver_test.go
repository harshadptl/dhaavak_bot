package routing

import "testing"

func TestResolve(t *testing.T) {
	bindings := []Binding{
		{Channel: "telegram", PeerKind: "user", PeerID: "42", AgentID: "personal"},
		{Channel: "telegram", PeerKind: "user", AgentID: "dm-default"},
		{Channel: "telegram", GuildID: "100", AgentID: "group-bot"},
		{Channel: "telegram", AgentID: "tg-wildcard"},
		{AgentID: "global"},
	}
	store := NewBindingStore(bindings, "fallback")
	r := NewResolver(store)

	tests := []struct {
		name   string
		params ResolveParams
		want   string
	}{
		{"exact peer match", ResolveParams{Channel: "telegram", PeerKind: "user", PeerID: "42"}, "personal"},
		{"parent peer match", ResolveParams{Channel: "telegram", PeerKind: "user", PeerID: "99"}, "dm-default"},
		{"guild match", ResolveParams{Channel: "telegram", PeerKind: "group", PeerID: "100", GuildID: "100"}, "group-bot"},
		{"channel wildcard", ResolveParams{Channel: "telegram", PeerKind: "group", PeerID: "200", GuildID: "200"}, "tg-wildcard"},
		{"account global", ResolveParams{Channel: "slack", PeerKind: "user", PeerID: "1"}, "global"},
		{"account global no peer", ResolveParams{Channel: "discord"}, "global"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := r.Resolve(tt.params)
			if got != tt.want {
				t.Errorf("Resolve() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestResolveDefaultFallback(t *testing.T) {
	s := NewBindingStore([]Binding{
		{Channel: "telegram", AgentID: "tg-only"},
	}, "fallback")
	r := NewResolver(s)

	got := r.Resolve(ResolveParams{Channel: "slack"})
	if got != "fallback" {
		t.Errorf("Resolve() = %q, want %q", got, "fallback")
	}
}
