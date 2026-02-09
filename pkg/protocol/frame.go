package protocol

import "encoding/json"

// RequestFrame is sent by clients to invoke a method.
type RequestFrame struct {
	ID     string          `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
}

// ResponseFrame is sent to clients in reply to a request.
type ResponseFrame struct {
	ID     string          `json:"id"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *ErrorDetail    `json:"error,omitempty"`
}

// EventFrame is a server-initiated push to clients.
type EventFrame struct {
	Event     string          `json:"event"`
	SessionID string          `json:"session_id,omitempty"`
	RunSeq    int             `json:"run_seq,omitempty"`
	Data      json.RawMessage `json:"data,omitempty"`
}

// ErrorDetail carries error info inside a ResponseFrame.
type ErrorDetail struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
