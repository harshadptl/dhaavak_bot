# Dhaavak

A personal AI assistant platform written in Go. Features a WebSocket gateway, multi-channel messaging (Telegram), session management, lane-based task queues, and an agentic runtime powered by Claude.

## Architecture

```
Telegram / WebSocket
        |
    Route Resolver  (7-level priority binding)
        |
    Session Manager (TTL, history, cleanup)
        |
    Lane Queue      (per-session serial execution)
        |
    Agent Runtime   (agentic loop: LLM -> tool -> repeat)
        |
    Claude API      (streaming, tool use)
        |
    Reply           (WS broadcast / Telegram message)
```

### Project Layout

```
cmd/dhaavak/       CLI entry point
internal/
  config/          YAML config with ${ENV_VAR} substitution
  gateway/         HTTP + WebSocket server, client mgmt, 150ms delta throttle
  session/         Session lifecycle, key building, send policy
  queue/           Per-session serial execution lanes
  routing/         7-level priority route resolution
  agent/           Agentic loop, conversation history, stream events
  llm/             Provider interface, Anthropic Claude implementation
  channel/         Adapter interface, registry
    telegram/      Bot polling, access control, message chunking
pkg/protocol/      WebSocket frame types, message types, event constants
```

## Quick Start

### Prerequisites

- Go 1.21+
- A [Telegram Bot Token](https://core.telegram.org/bots#botfather)
- An [Anthropic API Key](https://console.anthropic.com/)

### Setup

```bash
git clone https://github.com/harshadptl/dhaavak_bot.git
cd dhaavak_bot

export ANTHROPIC_API_KEY=sk-ant-...
export TELEGRAM_BOT_TOKEN=123456:ABC-...

make run
```

This builds the binary and starts the bot with the default config (`dhaavak.yaml`).

### Build Only

```bash
make build        # outputs bin/dhaavak
make test         # run tests
make clean        # remove build artifacts
```

## Configuration

Edit `dhaavak.yaml` or provide a custom path:

```bash
./bin/dhaavak --config /path/to/config.yaml
```

Environment variables are substituted using `${VAR_NAME}` syntax, with optional defaults via `${VAR_NAME:default}`.

### Key Settings

| Section | Field | Default | Description |
|---------|-------|---------|-------------|
| `server.port` | int | `18789` | WebSocket server port |
| `llm.model` | string | `claude-sonnet-4-5-20250929` | Claude model ID |
| `llm.max_turns` | int | `25` | Max agentic loop iterations |
| `channels.telegram.dm_policy` | string | `open` | `open`, `allowlist`, or `disabled` |
| `channels.telegram.group_policy` | string | `mention` | `mention`, `all`, or `disabled` |
| `session.ttl` | duration | `30m` | Session inactivity timeout |
| `session.max_history` | int | `100` | Max conversation turns kept |

## WebSocket API

Connect to `ws://127.0.0.1:18789/ws` (add `?token=...` if auth is configured).

### Send a message

```json
{
  "id": "req-1",
  "method": "chat.send",
  "params": {
    "session_id": "agent:default:main",
    "text": "Hello!",
    "agent_id": "default"
  }
}
```

### Events

| Event | Description |
|-------|-------------|
| `connected` | Connection established, includes `client_id` |
| `run.start` | Agent run begins |
| `chat.delta` | Streaming text chunk (throttled to 150ms) |
| `chat.tool_use` | Agent is calling a tool |
| `chat.tool_done` | Tool execution completed |
| `chat.complete` | Agent response finished |
| `chat.error` | Error during agent run |
| `run.end` | Agent run finished |

## Route Resolution

Messages are routed to agents using a 7-level priority chain:

1. **Exact peer** - specific user/group binding
2. **Parent peer** - channel + peer kind (e.g. all Telegram users)
3. **Guild** - group/server binding
4. **Team** - team-level binding
5. **Channel wildcard** - any message on a channel
6. **Account global** - catch-all binding
7. **Default agent** - configured fallback

## License

MIT

