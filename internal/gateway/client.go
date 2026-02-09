package gateway

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

const (
	sendChCap    = 256
	writeTimeout = 10 * time.Second
	readLimit    = 1 << 20 // 1 MB
)

// Client represents a single WebSocket connection.
type Client struct {
	ID        string
	conn      *websocket.Conn
	sendCh    chan []byte
	server    *Server
	sessions  map[string]bool // subscribed session IDs
	mu        sync.RWMutex
	cancelCtx context.CancelFunc
}

func newClient(conn *websocket.Conn, srv *Server) *Client {
	return &Client{
		ID:       uuid.New().String(),
		conn:     conn,
		sendCh:   make(chan []byte, sendChCap),
		server:   srv,
		sessions: make(map[string]bool),
	}
}

// Subscribe registers this client for events on the given session.
func (c *Client) Subscribe(sessionID string) {
	c.mu.Lock()
	c.sessions[sessionID] = true
	c.mu.Unlock()
}

// IsSubscribed checks if client listens to a session.
func (c *Client) IsSubscribed(sessionID string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.sessions[sessionID]
}

// Send enqueues a message for delivery. Drops if buffer is full.
func (c *Client) Send(data []byte) {
	select {
	case c.sendCh <- data:
	default:
		slog.Warn("client send buffer full, dropping message", "client", c.ID)
	}
}

// readPump reads frames from the WebSocket and dispatches them.
func (c *Client) readPump(ctx context.Context) {
	defer func() {
		c.server.unregister(c)
		c.conn.Close(websocket.StatusNormalClosure, "")
	}()

	c.conn.SetReadLimit(readLimit)

	for {
		_, data, err := c.conn.Read(ctx)
		if err != nil {
			if websocket.CloseStatus(err) != websocket.StatusNormalClosure {
				slog.Debug("client read error", "client", c.ID, "err", err)
			}
			return
		}
		c.server.handleMessage(c, data)
	}
}

// writePump sends queued messages to the WebSocket.
func (c *Client) writePump(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-c.sendCh:
			if !ok {
				return
			}
			writeCtx, cancel := context.WithTimeout(ctx, writeTimeout)
			err := c.conn.Write(writeCtx, websocket.MessageText, msg)
			cancel()
			if err != nil {
				slog.Debug("client write error", "client", c.ID, "err", err)
				return
			}
		}
	}
}

// sendJSON marshals v and enqueues it for sending.
func (c *Client) sendJSON(v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		slog.Error("marshal error", "err", err)
		return
	}
	c.Send(data)
}
