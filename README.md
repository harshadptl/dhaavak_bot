# Dhaavak_bot
Truly local AI life automation bot

Modules
- User prompt listener
- Messaging app listener
- Channel Adapter
- Gateway Server
- - Session Router
  - Lane Queue
- Agent Runner
- LLM API
- Agentic Loop
- Response Path
  - Channel Adapter

A high-performance Go implementation of OpenClaw - a personal AI assistant framework.

## Architecture Overview

Dhaavak_bot is a production-ready AI agent framework built with Go, featuring:

- **WebSocket Gateway**: Central control plane for session management and message routing
- **Multi-Channel Support**: WhatsApp, Telegram, Discord, Slack, and more
- **Lane Queue System**: Serial execution by default to prevent race conditions
- **AI Model Abstraction**: Support for Anthropic, OpenAI, and local models
- **Plugin Architecture**: Extensible tool and channel system
- **Memory Management**: JSONL transcripts + semantic search
- **Security First**: Sandboxing, approval workflows, and allowlists

## Project Structure

```
openclaw-go/
├── cmd/
│   ├── gateway/          # Gateway server executable
│   ├── cli/              # CLI tool
│   └── agent/            # Agent runner
├── internal/
│   ├── gateway/          # Gateway server implementation
│   │   ├── server.go
│   │   ├── session.go
│   │   └── queue.go
│   ├── channels/         # Channel adapters
│   │   ├── telegram/
│   │   ├── discord/
│   │   ├── whatsapp/
│   │   └── slack/
│   ├── agent/            # Agent runtime
│   │   ├── runner.go
│   │   ├── loop.go
│   │   └── context.go
│   ├── tools/            # Tool implementations
│   │   ├── browser/
│   │   ├── filesystem/
│   │   └── exec/
│   ├── models/           # AI model providers
│   │   ├── anthropic/
│   │   ├── openai/
│   │   └── local/
│   ├── memory/           # Memory management
│   │   ├── transcript.go
│   │   ├── semantic.go
│   │   └── storage.go
│   └── plugins/          # Plugin system
│       ├── loader.go
│       └── registry.go
├── pkg/
│   ├── protocol/         # WebSocket protocol
│   ├── config/           # Configuration management
│   └── utils/            # Shared utilities
├── go.mod
└── go.sum
```

## Quick Start

```bash
# Install
go install github.com/yourusername/openclaw-go/cmd/gateway@latest

# Run gateway
openclaw-gateway --config ~/.openclaw/config.yaml

# Run CLI
openclaw-cli message send --to "+1234567890" --text "Hello"
```

## Key Features

### 1. Gateway Server (Port 18789)
- WebSocket-based control plane
- Session coordination and routing
- Lane queue management
- Real-time client connections

### 2. Channel Adapters
- Unified message format across platforms
- Attachment extraction and processing
- Platform-specific authentication
- Error handling and retries

### 3. Agent Runtime
- 6-stage processing pipeline
- Model failover and rate limiting
- Tool call execution
- Streaming responses

### 4. Lane Queue System
- Serial execution prevents state drift
- Parallel execution for safe operations
- Configurable concurrency limits
- Dead letter queue for failures

### 5. Memory System
- JSONL transcripts for audit trails
- Markdown memory files
- Vector search (semantic)
- FTS5 keyword search
- Automatic index updates

### 6. Security
- Sandboxed tool execution
- Approval workflows
- Channel allowlists
- DM pairing system

## Configuration Example

```yaml
gateway:
  port: 18789
  bind: "127.0.0.1"

agent:
  model: "anthropic/claude-opus-4"
  max_tokens: 4096

channels:
  telegram:
    bot_token: "${TELEGRAM_BOT_TOKEN}"
    allowed_users:
      - "user123"
  
  discord:
    token: "${DISCORD_TOKEN}"
    dm_policy: "pairing"

tools:
  browser:
    enabled: true
  exec:
    enabled: true
    sandbox: true

memory:
  storage_path: "~/.openclaw/memory"
  semantic_search: true
```

## Development

```bash
# Clone repository
git clone https://github.com/yourusername/openclaw-go
cd openclaw-go

# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build -o bin/gateway cmd/gateway/main.go

# Run
./bin/gateway --config config.yaml
```

## Comparison with TypeScript OpenClaw

| Feature | TypeScript | Go |
|---------|-----------|-----|
| Performance | Good | Excellent |
| Memory Usage | Higher | Lower |
| Concurrency | Event loop | Goroutines |
| Type Safety | Good | Excellent |
| Deployment | Node runtime | Single binary |
| Startup Time | ~1-2s | ~100ms |

## License

MIT
