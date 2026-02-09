package queue

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

// Manager handles lane lifecycle: lazy creation, idle cleanup.
type Manager struct {
	lanes           map[string]*Lane
	mu              sync.Mutex
	bufferSize      int
	idleTimeout     time.Duration
	ctx             context.Context
}

// NewManager creates a queue manager.
func NewManager(ctx context.Context, bufferSize int, idleTimeout time.Duration) *Manager {
	return &Manager{
		lanes:       make(map[string]*Lane),
		bufferSize:  bufferSize,
		idleTimeout: idleTimeout,
		ctx:         ctx,
	}
}

// Enqueue adds a task to the lane for the given session, creating it lazily if needed.
func (m *Manager) Enqueue(t Task) bool {
	m.mu.Lock()
	l, ok := m.lanes[t.SessionID]
	if !ok {
		l = newLane(t.SessionID, m.bufferSize, m.ctx)
		m.lanes[t.SessionID] = l
		slog.Debug("lane created", "session", t.SessionID)
	}
	m.mu.Unlock()
	return l.Enqueue(t)
}

// StartCleanup launches a goroutine that removes idle lanes.
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

	cutoff := time.Now().Add(-m.idleTimeout)
	for id, l := range m.lanes {
		if l.idleSince().Before(cutoff) {
			l.Stop()
			delete(m.lanes, id)
			slog.Debug("lane removed (idle)", "session", id)
		}
	}
}

// StopAll stops all lanes.
func (m *Manager) StopAll() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for id, l := range m.lanes {
		l.Stop()
		delete(m.lanes, id)
	}
}
