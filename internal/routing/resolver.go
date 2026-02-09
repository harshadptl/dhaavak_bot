package routing

// ResolveParams are the inputs to route resolution.
type ResolveParams struct {
	Channel  string
	PeerKind string // "user", "group"
	PeerID   string
	GuildID  string
	TeamID   string
}

// Resolver determines which agent handles a given message context.
type Resolver struct {
	store *BindingStore
}

// NewResolver creates a route resolver.
func NewResolver(store *BindingStore) *Resolver {
	return &Resolver{store: store}
}

// Resolve walks a 7-level priority chain and returns the best matching agent ID.
//
// Priority (lowest number wins):
//  1. Exact peer binding (channel + peer_kind + peer_id)
//  2. Parent peer binding (channel + peer_kind, no peer_id — e.g. "all users on telegram")
//  3. Guild/group binding (channel + guild_id)
//  4. Team binding (team_id)
//  5. Account/global binding (no channel filter)
//  6. Channel wildcard (channel only, no peer/guild)
//  7. Default agent
func (r *Resolver) Resolve(p ResolveParams) string {
	var (
		peerMatch    string // priority 1
		parentPeer   string // priority 2
		guildMatch   string // priority 3
		teamMatch    string // priority 4
		accountMatch string // priority 5
		channelWild  string // priority 6
	)

	for _, b := range r.store.Bindings() {
		// Priority 1: exact peer
		if b.Channel == p.Channel && b.PeerKind == p.PeerKind && b.PeerID == p.PeerID && b.PeerID != "" {
			peerMatch = b.AgentID
		}
		// Priority 2: parent peer (same channel + kind, no specific peer)
		if b.Channel == p.Channel && b.PeerKind == p.PeerKind && b.PeerID == "" && b.GuildID == "" {
			parentPeer = b.AgentID
		}
		// Priority 3: guild
		if b.Channel == p.Channel && b.GuildID == p.GuildID && p.GuildID != "" && b.PeerKind == "" {
			guildMatch = b.AgentID
		}
		// Priority 4: team
		if b.TeamID == p.TeamID && p.TeamID != "" && b.Channel == "" {
			teamMatch = b.AgentID
		}
		// Priority 5: account/global (no channel, no peer, no guild, no team)
		if b.Channel == "" && b.PeerKind == "" && b.PeerID == "" && b.GuildID == "" && b.TeamID == "" && b.AgentID != "" {
			accountMatch = b.AgentID
		}
		// Priority 6: channel wildcard (channel only)
		if b.Channel == p.Channel && b.PeerKind == "" && b.PeerID == "" && b.GuildID == "" && b.TeamID == "" {
			channelWild = b.AgentID
		}
	}

	switch {
	case peerMatch != "":
		return peerMatch
	case parentPeer != "":
		return parentPeer
	case guildMatch != "":
		return guildMatch
	case teamMatch != "":
		return teamMatch
	case channelWild != "":
		return channelWild
	case accountMatch != "":
		return accountMatch
	default:
		return r.store.DefaultAgent()
	}
}
