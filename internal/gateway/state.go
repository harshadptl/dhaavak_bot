package gateway

import "sync"

// RunState tracks the current run sequence per session.
type RunState struct {
	mu   sync.RWMutex
	seqs map[string]int // session_id -> current run sequence
}

// NewRunState creates a new RunState tracker.
func NewRunState() *RunState {
	return &RunState{seqs: make(map[string]int)}
}

// Next increments and returns the next run sequence for a session.
func (rs *RunState) Next(sessionID string) int {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.seqs[sessionID]++
	return rs.seqs[sessionID]
}

// Current returns the current run sequence for a session.
func (rs *RunState) Current(sessionID string) int {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	return rs.seqs[sessionID]
}
