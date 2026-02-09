package channel

import (
	"context"

	"github.com/harshadpatil/dhaavak/pkg/protocol"
)

// Adapter represents a messaging channel (Telegram, Slack, etc).
type Adapter interface {
	// ID returns the adapter identifier (e.g. "telegram").
	ID() string

	// Start begins listening for messages.
	Start(ctx context.Context) error

	// Stop gracefully shuts down the adapter.
	Stop(ctx context.Context) error

	// SendMessage delivers a message back to the channel.
	SendMessage(ctx context.Context, msg protocol.OutboundMessage) error
}

// MessageSink is called by adapters when they receive a message.
// It decouples the adapter from gateway internals.
type MessageSink func(ctx context.Context, msg protocol.InboundMessage) error
