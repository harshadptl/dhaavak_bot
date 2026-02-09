package gateway

import (
	"encoding/json"
	"log/slog"

	"github.com/harshadpatil/dhaavak/pkg/protocol"
)

// BroadcastAll sends an event to every connected client.
func (s *Server) BroadcastAll(event protocol.EventFrame) {
	data, err := json.Marshal(event)
	if err != nil {
		slog.Error("broadcast marshal error", "err", err)
		return
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range s.clients {
		c.Send(data)
	}
}

// BroadcastSession sends an event to clients subscribed to a session.
func (s *Server) BroadcastSession(sessionID string, event protocol.EventFrame) {
	data, err := json.Marshal(event)
	if err != nil {
		slog.Error("broadcast marshal error", "err", err)
		return
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range s.clients {
		if c.IsSubscribed(sessionID) {
			c.Send(data)
		}
	}
}

// BroadcastClient sends an event to a specific client.
func (s *Server) BroadcastClient(clientID string, event protocol.EventFrame) {
	data, err := json.Marshal(event)
	if err != nil {
		slog.Error("broadcast marshal error", "err", err)
		return
	}
	s.mu.RLock()
	c, ok := s.clients[clientID]
	s.mu.RUnlock()
	if ok {
		c.Send(data)
	}
}
