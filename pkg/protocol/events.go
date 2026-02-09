package protocol

// Methods (client -> server requests)
const (
	MethodChatSend    = "chat.send"
	MethodChatCancel  = "chat.cancel"
	MethodSessionList = "session.list"
	MethodSessionGet  = "session.get"
	MethodPing        = "ping"
)

// Events (server -> client pushes)
const (
	EventChatDelta    = "chat.delta"
	EventChatComplete = "chat.complete"
	EventChatError    = "chat.error"
	EventChatToolUse  = "chat.tool_use"
	EventChatToolDone = "chat.tool_done"
	EventRunStart     = "run.start"
	EventRunEnd       = "run.end"
	EventConnected    = "connected"
)
