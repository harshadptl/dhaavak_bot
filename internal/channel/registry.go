package channel

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/harshadpatil/dhaavak/pkg/protocol"
)

// Registry manages adapter lifecycle.
type Registry struct {
	adapters map[string]Adapter
	mu       sync.RWMutex
}

// NewRegistry creates a channel registry.
func NewRegistry() *Registry {
	return &Registry{adapters: make(map[string]Adapter)}
}

// Register adds an adapter.
func (r *Registry) Register(a Adapter) {
	r.mu.Lock()
	r.adapters[a.ID()] = a
	r.mu.Unlock()
	slog.Info("channel registered", "channel", a.ID())
}

// Get returns an adapter by ID.
func (r *Registry) Get(id string) (Adapter, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.adapters[id]
	return a, ok
}

// StartAll starts all registered adapters.
func (r *Registry) StartAll(ctx context.Context) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, a := range r.adapters {
		if err := a.Start(ctx); err != nil {
			return fmt.Errorf("start channel %s: %w", a.ID(), err)
		}
		slog.Info("channel started", "channel", a.ID())
	}
	return nil
}

// StopAll stops all registered adapters.
func (r *Registry) StopAll(ctx context.Context) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, a := range r.adapters {
		if err := a.Stop(ctx); err != nil {
			slog.Error("stop channel error", "channel", a.ID(), "err", err)
		}
	}
}

// SendMessage routes an outbound message to the correct adapter.
func (r *Registry) SendMessage(ctx context.Context, msg protocol.OutboundMessage) error {
	a, ok := r.Get(msg.Channel)
	if !ok {
		return fmt.Errorf("channel not found: %s", msg.Channel)
	}
	return a.SendMessage(ctx, msg)
}
