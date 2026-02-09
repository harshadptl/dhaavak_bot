package session

import (
	"sync"
	"time"
)

// Entry holds one active session's state.
type Entry struct {
	Key       string
	AgentID   string
	CreatedAt time.Time
	TouchedAt time.Time
	History   []Message
	mu        sync.Mutex
}

// Message is a single turn in the conversation history.
type Message struct {
	Role    string `json:"role"` // "user", "assistant"
	Content string `json:"content"`
}

// Touch updates the last-access timestamp.
func (e *Entry) Touch() {
	e.mu.Lock()
	e.TouchedAt = time.Now()
	e.mu.Unlock()
}

// AppendHistory adds a message, enforcing max history length.
func (e *Entry) AppendHistory(msg Message, maxHistory int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.History = append(e.History, msg)
	if maxHistory > 0 && len(e.History) > maxHistory {
		e.History = e.History[len(e.History)-maxHistory:]
	}
	e.TouchedAt = time.Now()
}

// GetHistory returns a copy of the conversation history.
func (e *Entry) GetHistory() []Message {
	e.mu.Lock()
	defer e.mu.Unlock()
	cp := make([]Message, len(e.History))
	copy(cp, e.History)
	return cp
}
