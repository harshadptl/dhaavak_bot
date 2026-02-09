package queue

import (
	"context"
	"log/slog"
	"sync/atomic"
	"time"
)

// Lane is a per-session serial execution queue.
// One goroutine processes tasks sequentially from a buffered channel.
type Lane struct {
	sessionID string
	tasks     chan Task
	lastUsed  atomic.Int64 // unix nanos
	cancel    context.CancelFunc
}

func newLane(sessionID string, bufferSize int, ctx context.Context) *Lane {
	ctx, cancel := context.WithCancel(ctx)
	l := &Lane{
		sessionID: sessionID,
		tasks:     make(chan Task, bufferSize),
		cancel:    cancel,
	}
	l.touch()
	go l.run(ctx)
	return l
}

// Enqueue adds a task to the lane. Returns false if the buffer is full.
func (l *Lane) Enqueue(t Task) bool {
	l.touch()
	select {
	case l.tasks <- t:
		return true
	default:
		slog.Warn("lane queue full", "session", l.sessionID)
		return false
	}
}

// Stop signals the lane goroutine to exit.
func (l *Lane) Stop() {
	l.cancel()
}

func (l *Lane) touch() {
	l.lastUsed.Store(time.Now().UnixNano())
}

func (l *Lane) idleSince() time.Time {
	return time.Unix(0, l.lastUsed.Load())
}

func (l *Lane) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case t, ok := <-l.tasks:
			if !ok {
				return
			}
			l.touch()
			if err := t.Fn(ctx); err != nil {
				slog.Error("lane task error", "session", l.sessionID, "err", err)
			}
			l.touch()
		}
	}
}
