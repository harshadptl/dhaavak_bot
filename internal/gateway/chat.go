package gateway

import (
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/harshadpatil/dhaavak/pkg/protocol"
)

const throttleInterval = 150 * time.Millisecond

// ChatRunState tracks per-session streaming delta throttle buffers.
type ChatRunState struct {
	mu      sync.Mutex
	buffers map[string]*deltaBuffer // session_id -> buffer
	server  *Server
}

type deltaBuffer struct {
	text    string
	timer   *time.Timer
	session string
	runSeq  int
}

// NewChatRunState creates a new ChatRunState.
func NewChatRunState(srv *Server) *ChatRunState {
	return &ChatRunState{
		buffers: make(map[string]*deltaBuffer),
		server:  srv,
	}
}

// AccumulateDelta adds streaming text and flushes on a 150ms throttle.
func (cs *ChatRunState) AccumulateDelta(sessionID string, runSeq int, text string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	buf, ok := cs.buffers[sessionID]
	if !ok {
		buf = &deltaBuffer{session: sessionID, runSeq: runSeq}
		cs.buffers[sessionID] = buf
	}
	buf.text += text
	buf.runSeq = runSeq

	if buf.timer == nil {
		buf.timer = time.AfterFunc(throttleInterval, func() {
			cs.flush(sessionID)
		})
	}
}

// Flush immediately sends any buffered delta for a session.
func (cs *ChatRunState) Flush(sessionID string) {
	cs.flush(sessionID)
}

func (cs *ChatRunState) flush(sessionID string) {
	cs.mu.Lock()
	buf, ok := cs.buffers[sessionID]
	if !ok || buf.text == "" {
		cs.mu.Unlock()
		return
	}
	text := buf.text
	runSeq := buf.runSeq
	buf.text = ""
	if buf.timer != nil {
		buf.timer.Stop()
		buf.timer = nil
	}
	cs.mu.Unlock()

	data, err := json.Marshal(map[string]string{"text": text})
	if err != nil {
		slog.Error("flush marshal error", "err", err)
		return
	}
	cs.server.BroadcastSession(sessionID, protocol.EventFrame{
		Event:     protocol.EventChatDelta,
		SessionID: sessionID,
		RunSeq:    runSeq,
		Data:      data,
	})
}

// Clear removes the buffer for a session.
func (cs *ChatRunState) Clear(sessionID string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	if buf, ok := cs.buffers[sessionID]; ok {
		if buf.timer != nil {
			buf.timer.Stop()
		}
		delete(cs.buffers, sessionID)
	}
}
