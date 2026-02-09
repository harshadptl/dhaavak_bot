package session

import "testing"

func TestSendPolicyAllowDM(t *testing.T) {
	tests := []struct {
		name   string
		policy SendPolicy
		userID int64
		want   bool
	}{
		{"open allows all", SendPolicy{DMPolicy: "open"}, 42, true},
		{"disabled denies all", SendPolicy{DMPolicy: "disabled"}, 42, false},
		{"allowlist match", SendPolicy{DMPolicy: "allowlist", AllowedUsers: []int64{42, 99}}, 42, true},
		{"allowlist no match", SendPolicy{DMPolicy: "allowlist", AllowedUsers: []int64{99}}, 42, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.policy.AllowDM(tt.userID); got != tt.want {
				t.Errorf("AllowDM(%d) = %v, want %v", tt.userID, got, tt.want)
			}
		})
	}
}

func TestSendPolicyAllowGroup(t *testing.T) {
	tests := []struct {
		name    string
		policy  SendPolicy
		groupID int64
		want    bool
	}{
		{"disabled denies", SendPolicy{GroupPolicy: "disabled"}, 100, false},
		{"all no filter", SendPolicy{GroupPolicy: "all"}, 100, true},
		{"all with filter match", SendPolicy{GroupPolicy: "all", AllowedGroups: []int64{100}}, 100, true},
		{"all with filter no match", SendPolicy{GroupPolicy: "all", AllowedGroups: []int64{200}}, 100, false},
		{"mention no filter", SendPolicy{GroupPolicy: "mention"}, 100, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.policy.AllowGroup(tt.groupID); got != tt.want {
				t.Errorf("AllowGroup(%d) = %v, want %v", tt.groupID, got, tt.want)
			}
		})
	}
}
