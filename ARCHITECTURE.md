# Dhaavak Architecture

## Overview

Dhaavak is a multi-channel AI agent orchestration platform. Messages flow from channels (Telegram, WebSocket) through routing, session management, and queuing into an agentic runtime powered by Claude, with streaming responses broadcast back to clients.

```
                    +-------------------+
                    |   Telegram Bot    |
                    +--------+----------+
                             |
                    +--------v----------+
                    |   Route Resolver  |  7-level priority binding
                    +--------+----------+
                             |
+----------------+  +--------v----------+
| WebSocket GW   +-->  Session Manager  |  TTL, history, cleanup
+----------------+  +--------+----------+
                             |
                    +--------v----------+
                    |    Lane Queue     |  per-session serial execution
                    +--------+----------+
                             |
                    +--------v----------+
                    |   Agent Runtime   |  agentic loop (LLM -> tool -> repeat)
                    +--------+----------+
                             |
                    +--------v----------+
                    |   Claude API      |  streaming, tool use
                    +--------+----------+
                             |
              +--------------+--------------+
              |                             |
     +--------v----------+        +---------v--------+
     | WS Broadcast      |        | Telegram Reply   |
     | (delta throttled)  |        | (chunked HTML)   |
     +--------------------+        +------------------+
```

---

## Project Layout

```
cmd/dhaavak/           CLI entry point, component wiring
internal/
  config/              YAML loading, ${ENV_VAR} substitution, validation
  gateway/             HTTP + WebSocket server
  session/             Session lifecycle, key building, access policy
  queue/               Per-session serial execution lanes
  routing/             Priority-based agent resolution
  agent/               Agentic loop, conversation, stream events
  llm/                 Provider interface, Anthropic implementation
  channel/             Adapter interface, registry
    telegram/          Bot polling, access control, message delivery
pkg/protocol/          Frame types, message types, event constants
```

---

## Subsystem Details

### 1. Gateway Server

**Files:** `internal/gateway/`

The gateway exposes an HTTP server with `/ws` (WebSocket) and `/health` endpoints.

**Key types:**

| Type | Purpose |
|------|---------|
| `Server` | HTTP listener, client registry, broadcast hub |
| `Client` | Single WebSocket connection with read/write pump goroutines |
| `Authenticator` | Timing-safe token validation (Bearer header or query param) |
| `RunState` | Per-session monotonic run sequence counter |
| `ChatRunState` | Streaming delta accumulator with 150ms throttle |

**Goroutine model per client:**

```
WebSocket Accept
  |
  +-- readPump(ctx)   blocks on conn.Read(), dispatches to server
  |
  +-- writePump(ctx)  drains sendCh (cap 256), writes to conn
```

**Delta throttling:** Text deltas accumulate in a buffer. A 150ms timer fires to flush the buffer as a single `chat.delta` event. On completion, an immediate flush ensures no text is lost.

**Request/Response flow:**

```
Client sends RequestFrame  -->  readPump  -->  server.handleMessage()
                                                  |
                           chat.send  ----------->  OnChatSend callback
                           ping       ----------->  ResponseFrame{pong: ok}
```

### 2. Session Management

**Files:** `internal/session/`

Sessions track conversation history and last-access time per agent-channel-peer combination.

**Session key format:**

| Context | Key |
|---------|-----|
| WebSocket / default | `agent:{id}:main` |
| Telegram DM | `agent:{id}:telegram:user:{userID}` |
| Telegram group | `agent:{id}:telegram:group:{groupID}` |
| Group thread | `agent:{id}:telegram:group:{groupID}:{threadID}` |

**Entry struct:** Holds `Key`, `AgentID`, `CreatedAt`, `TouchedAt`, and `History []Message`. History is bounded by `maxHistory` — oldest messages are trimmed on append.

**Cleanup:** A background goroutine runs on a configurable interval, removing sessions not touched within the TTL.

**SendPolicy:** Controls per-channel access:
- **DM:** `open` (all), `allowlist` (specific user IDs), `disabled`
- **Group:** `mention` (only when bot is @mentioned), `all`, `disabled`

### 3. Lane Queue

**Files:** `internal/queue/`

Each session gets a dedicated lane — one goroutine reading from a buffered channel. This guarantees serial execution per session while allowing concurrency across sessions.

```
Session A:  [task1] -> [task2] -> [task3]   (serial)
Session B:  [task1] -> [task2]              (serial, concurrent with A)
```

**Task type:**

```go
type Task struct {
    SessionID string
    Fn        func(ctx context.Context) error
}
```

**Lane lifecycle:**
1. Created lazily on first `Enqueue()` for a session
2. Worker goroutine reads tasks sequentially
3. Idle lanes cleaned up after configurable timeout
4. `StopAll()` on shutdown cancels all lane contexts

**Back-pressure:** If the lane's buffered channel is full, `Enqueue()` returns `false` and the message is dropped (logged as warning).

### 4. Route Resolution

**Files:** `internal/routing/`

The resolver walks a 7-level priority chain to determine which agent handles a message:

| Priority | Match Criteria | Example |
|----------|---------------|---------|
| 1 | Exact peer (channel + kind + ID) | User 42 on Telegram -> `personal-bot` |
| 2 | Parent peer (channel + kind) | All Telegram users -> `dm-bot` |
| 3 | Guild (channel + group ID) | Group 100 on Telegram -> `group-bot` |
| 4 | Team (team ID) | Team X -> `team-bot` |
| 5 | Channel wildcard (channel only) | Any Telegram message -> `tg-default` |
| 6 | Account global (no filters) | Catch-all -> `global-bot` |
| 7 | Default agent | Configured fallback |

The binding store is immutable after startup — resolution is read-only with no locking.

### 5. Agent Runtime

**Files:** `internal/agent/`

**Agentic loop (`RunLoop`):**

```
for turn in 0..maxTurns:
    1. provider.Stream(systemPrompt, messages, tools)
    2. Collect text deltas + tool_use blocks
    3. Append assistant message to conversation
    4. If no tool calls -> return final text
    5. Execute each tool via ToolExecutor callback
    6. Append tool results as user message
    7. Continue loop
```

**Event flow:** Each streaming event from the LLM is mapped to an `agent.Event` and forwarded to the `EventSink` callback, which routes it to the gateway for WebSocket broadcasting.

**Conversation building:** `BuildMessages()` converts session history into `[]llm.Message`, appending the new user message at the end.

### 6. LLM Layer

**Files:** `internal/llm/`

**Provider interface:**

```go
type Provider interface {
    Stream(ctx, systemPrompt, messages, tools) (<-chan StreamEvent, error)
    Complete(ctx, systemPrompt, messages, tools) (*CompletionResult, error)
}
```

**Anthropic implementation:**
- Uses `anthropic-sdk-go` with streaming support
- Handles `content_block_start/delta/stop` and `message_stop` events
- Accumulates tool input JSON from partial deltas
- Uses `Message.Accumulate()` to track stop reason
- MaxTokens: 8192

**ContentBlock types:** `text`, `tool_use` (ID + name + input JSON), `tool_result` (ID + output text).

### 7. Channel Adapters

**Files:** `internal/channel/`, `internal/channel/telegram/`

**Adapter interface:**

```go
type Adapter interface {
    ID() string
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
    SendMessage(ctx context.Context, msg protocol.OutboundMessage) error
}
```

**MessageSink** decouples adapters from internals — adapters call the sink with `InboundMessage`, never touching sessions or queues directly.

**Telegram adapter:**
- Long-polling with 30s timeout via `go-telegram-bot-api`
- `extractContext()` parses updates: detects DMs vs groups, bot mentions, extracts text
- `checkAccess()` evaluates send policy before processing
- `sendText()` chunks output at 4000 chars, breaks at newlines, sends as HTML with plain-text fallback

**Registry** manages adapter lifecycle: `Register()`, `StartAll()`, `StopAll()`, and `SendMessage()` routing.

### 8. Protocol Types

**Files:** `pkg/protocol/`

| Type | Direction | Purpose |
|------|-----------|---------|
| `RequestFrame` | Client -> Server | Method invocation (id, method, params) |
| `ResponseFrame` | Server -> Client | Request reply (id, result/error) |
| `EventFrame` | Server -> Client | Async push (event, session_id, run_seq, data) |
| `InboundMessage` | Channel -> System | Unified incoming message |
| `OutboundMessage` | System -> Channel | Reply to deliver |

**Methods:** `chat.send`, `chat.cancel`, `session.list`, `session.get`, `ping`

**Events:** `connected`, `run.start`, `chat.delta`, `chat.tool_use`, `chat.tool_done`, `chat.complete`, `chat.error`, `run.end`

---

## End-to-End Message Flow

### WebSocket Client

```
1. Client connects to /ws, receives "connected" event
2. Client sends: {"id":"1", "method":"chat.send", "params":{"text":"hello"}}
3. Server responds: {"id":"1", "result":{"status":"queued"}}
4. Server broadcasts to session subscribers:
     run.start  ->  chat.delta (throttled)  ->  chat.complete  ->  run.end
```

### Telegram

```
1. User sends "hello @bot" in group chat
2. Bot polls update, extracts context
3. Access check: group policy + mention gate
4. MessageSink -> processMessage()
5. Route -> session -> lane queue -> agent runtime
6. Agent streams response via Claude API
7. Events broadcast to WS subscribers (if any)
8. Final text sent back via Telegram API (chunked HTML)
```

---

## Concurrency Model

| Component | Pattern | Synchronization |
|-----------|---------|-----------------|
| Gateway clients map | Read-heavy | `sync.RWMutex` |
| Client send buffer | Producer-consumer | Buffered channel (cap 256), drop on full |
| Session entries | Per-entry locking | `sync.Mutex` on each Entry |
| Sessions map | Read-heavy | `sync.RWMutex` |
| Queue lanes | One goroutine per session | Buffered task channel |
| Lane last-used time | Lock-free | `atomic.Int64` |
| Route resolution | Read-only | No locking (immutable after init) |
| Delta throttle | Timer-based flush | `sync.Mutex` on buffer map |

**Goroutine budget:** 2 per WebSocket client + 1 per active session lane + 1 Telegram poller + 2 cleanup timers + 1 HTTP listener + 1 per active LLM stream.

---

## Initialization Order

```
1. Config         load YAML, substitute env vars, validate
2. Session Mgr    create manager, start cleanup goroutine
3. Queue Mgr      create manager, start cleanup goroutine
4. Router         build binding store and resolver
5. LLM Provider   create Anthropic client
6. Agent Runtime  register agents, wire event sink + tool executor
7. Gateway        create server, wire OnChatSend handler
8. Telegram       create bot, set message sink, register in channel registry
9. Start          registry.StartAll() then gw.Start()
10. Signal wait   SIGINT/SIGTERM -> cancel context -> drain -> shutdown
```

---

## Extension Points

### Adding a new channel adapter

1. Create `internal/channel/slack/` implementing the `Adapter` interface
2. Parse channel-specific messages into `protocol.InboundMessage`
3. Implement `SendMessage()` for outbound delivery
4. Set `MessageSink` callback for inbound routing
5. Register in `main.go` via `registry.Register()`

### Adding a new LLM provider

1. Create `internal/llm/openai.go` implementing the `Provider` interface
2. Handle streaming events and tool calling in your provider's format
3. Map to `StreamEvent` types for compatibility with the agent runtime
4. Switch on `cfg.LLM.Provider` in `main.go`

### Adding agent tools

1. Define `llm.ToolDef` with name, description, and JSON schema
2. Add to agent definition in config
3. Implement `ToolExecutor` function in `main.go`
4. Register via `runtime.SetToolExecutor()`
