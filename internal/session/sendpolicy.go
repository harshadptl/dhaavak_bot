package session

// SendPolicy controls whether messages are allowed for a session context.
type SendPolicy struct {
	DMPolicy      string  // "open", "allowlist", "disabled"
	GroupPolicy   string  // "mention", "all", "disabled"
	AllowedUsers  []int64
	AllowedGroups []int64
}

// AllowDM checks if a DM from userID is permitted.
func (p *SendPolicy) AllowDM(userID int64) bool {
	switch p.DMPolicy {
	case "open":
		return true
	case "allowlist":
		for _, id := range p.AllowedUsers {
			if id == userID {
				return true
			}
		}
		return false
	default:
		return false
	}
}

// AllowGroup checks if a group message from groupID is permitted.
func (p *SendPolicy) AllowGroup(groupID int64) bool {
	switch p.GroupPolicy {
	case "disabled":
		return false
	case "all", "mention":
		if len(p.AllowedGroups) == 0 {
			return true
		}
		for _, id := range p.AllowedGroups {
			if id == groupID {
				return true
			}
		}
		return false
	default:
		return false
	}
}
