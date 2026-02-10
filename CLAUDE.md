# CLAUDE.md

This file provides guidance to Claude Code when working on this codebase.

## Project Overview

Dhaavak is a multi-channel AI assistant platform written in Go. It routes messages from Telegram and WebSocket clients through session management and lane queues into an agentic runtime powered by Claude.

**Module:** `github.com/harshadpatil/dhaavak`

## Build & Test

```bash
make build        # Build binary to bin/dhaavak
make test         # Run all tests (go test ./...)
make run          # Build and run with dhaavak.yaml
make tidy         # go mod tidy
make clean        # Remove bin/
```

Direct commands:
```bash
go build ./...              # Compile all packages
go test ./...               # Run all tests
go vet ./...                # Static analysis
go build -trimpath -o bin/dhaavak ./cmd/dhaavak  # Production build
```

## Project Structure

```
cmd/dhaavak/main.go         Entry point + component wiring
internal/config/             YAML config, env substitution, validation
internal/gateway/            WebSocket server, client mgmt, broadcasting, delta throttle
internal/session/            Session lifecycle, key building, send policy
internal/queue/              Per-session serial execution lanes
internal/routing/            7-level priority route resolution
internal/agent/              Agentic loop, conversation history, stream events
internal/llm/                Provider interface, Anthropic Claude implementation
internal/channel/            Adapter interface + registry
internal/channel/telegram/   Telegram bot polling, access control, message delivery
pkg/protocol/                WebSocket frame types, message types, event constants
```

## Key Architecture Patterns

- **Interfaces at boundaries:** `llm.Provider`, `channel.Adapter`, `channel.MessageSink` — always code to these interfaces
- **Per-session serial execution:** Lane queues ensure one goroutine per session processes tasks sequentially
- **Goroutine pairs per WS client:** `readPump` + `writePump` with a buffered `sendCh` (cap 256)
- **150ms delta throttle:** Streaming text accumulates and flushes periodically to reduce WS traffic
- **Immutable routing:** Binding store and resolver are read-only after init — no locking needed
- **Context propagation:** All operations accept `context.Context` for cancellation and timeouts

## Key Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/coder/websocket` | WebSocket (context-aware, concurrent-safe) |
| `github.com/go-telegram-bot-api/telegram-bot-api/v5` | Telegram Bot API |
| `github.com/anthropics/anthropic-sdk-go` | Claude API with streaming + tool use |
| `github.com/knadh/koanf/v2` | Config loading (YAML) |
| `github.com/google/uuid` | Request/session IDs |

## Configuration

Config loaded from YAML with `${ENV_VAR}` substitution (optional default: `${VAR:default}`).

Required env vars for runtime:
- `ANTHROPIC_API_KEY` — Claude API key
- `TELEGRAM_BOT_TOKEN` — Telegram bot token (if telegram enabled)

## Code Conventions

- Use `log/slog` for all logging (structured, stdlib)
- Errors wrap with `fmt.Errorf("context: %w", err)`
- Session keys follow format: `agent:{id}:{channel}:{peerKind}:{peerId}`
- Tests live alongside source files (`*_test.go`)
- No external test frameworks — stdlib `testing` only

## Test Coverage

Tests exist for:
- `internal/session/` — key building/parsing, send policy
- `internal/routing/` — 7-level priority resolution, fallback
- `internal/queue/` — serial execution guarantee, cross-session parallelism
