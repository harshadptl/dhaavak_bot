package session

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

// Manager handles session lifecycle.
type Manager struct {
	sessions   map[string]*Entry
	mu         sync.RWMutex
	ttl        time.Duration
	maxHistory int
}

// NewManager creates a session manager.
func NewManager(ttl time.Duration, maxHistory int) *Manager {
	return &Manager{
		sessions:   make(map[string]*Entry),
		ttl:        ttl,
		maxHistory: maxHistory,
	}
}

// GetOrCreate returns an existing session or creates a new one.
func (m *Manager) GetOrCreate(key, agentID string) *Entry {
	m.mu.Lock()
	defer m.mu.Unlock()

	if e, ok := m.sessions[key]; ok {
		e.Touch()
		return e
	}

	now := time.Now()
	e := &Entry{
		Key:       key,
		AgentID:   agentID,
		CreatedAt: now,
		TouchedAt: now,
	}
	m.sessions[key] = e
	slog.Debug("session created", "key", key, "agent", agentID)
	return e
}

// Get returns a session if it exists.
func (m *Manager) Get(key string) (*Entry, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.sessions[key]
	if ok {
		e.Touch()
	}
	return e, ok
}

// MaxHistory returns the configured max history.
func (m *Manager) MaxHistory() int {
	return m.maxHistory
}

// StartCleanup launches a goroutine that removes expired sessions.
func (m *Manager) StartCleanup(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				m.cleanup()
			}
		}
	}()
}

func (m *Manager) cleanup() {
	m.mu.Lock()
	defer m.mu.Unlock()

	cutoff := time.Now().Add(-m.ttl)
	for key, e := range m.sessions {
		e.mu.Lock()
		expired := e.TouchedAt.Before(cutoff)
		e.mu.Unlock()
		if expired {
			delete(m.sessions, key)
			slog.Debug("session expired", "key", key)
		}
	}
}
