# Dhaavak → OpenClaw Feature Parity Gap Analysis

A comprehensive list of features Dhaavak needs to reach parity with [OpenClaw](https://github.com/openclaw/openclaw).

> **Legend:** Items marked with `[PARTIAL]` indicate Dhaavak has some foundation but needs significant expansion. Items marked with `[MISSING]` are entirely absent.

---

## 1. Channel Adapters

Dhaavak currently supports **Telegram** and **WebSocket (WebChat)**. OpenClaw supports 14+ channels.

| Channel | Status |
|---------|--------|
| WhatsApp (Baileys) | `[MISSING]` |
| Slack (Bolt) | `[MISSING]` |
| Discord (discord.js) | `[MISSING]` |
| Signal (signal-cli) | `[MISSING]` |
| iMessage / BlueBubbles | `[MISSING]` |
| Microsoft Teams | `[MISSING]` |
| Google Chat | `[MISSING]` |
| Matrix | `[MISSING]` |
| Nostr | `[MISSING]` |
| Nextcloud Talk | `[MISSING]` |
| Tlon Messenger | `[MISSING]` |
| Zalo (personal & bot) | `[MISSING]` |

---

## 2. AI Model Providers

Dhaavak supports **Anthropic Claude only**. OpenClaw supports 11+ providers with failover.

| Feature | Status |
|---------|--------|
| OpenAI (GPT series) provider | `[MISSING]` |
| Google Gemini provider | `[MISSING]` |
| xAI Grok provider | `[MISSING]` |
| OpenRouter gateway | `[MISSING]` |
| Mistral provider | `[MISSING]` |
| DeepSeek provider | `[MISSING]` |
| Perplexity provider | `[MISSING]` |
| Hugging Face provider | `[MISSING]` |
| Local models (Ollama / LM Studio) | `[MISSING]` |
| Model failover & rotation | `[MISSING]` |
| Per-session model selection | `[PARTIAL]` — per-agent model exists, but not per-session switching |
| Auth profile rotation (OAuth vs API keys) | `[MISSING]` |

---

## 3. Tool & Skill System

Dhaavak has a basic `ToolExecutor` callback. OpenClaw has a full-fledged skill platform with groups, permissions, and a registry.

| Feature | Status |
|---------|--------|
| Built-in shell execution tool (`exec`, `bash`) | `[MISSING]` |
| Background process management (`process`) | `[MISSING]` |
| File operation tools (`read`, `write`, `edit`) | `[MISSING]` |
| Multi-file patch application (`apply_patch`) | `[MISSING]` |
| Web search tool (Brave Search) | `[MISSING]` |
| Web fetch / URL extraction tool | `[MISSING]` |
| Image analysis tool | `[MISSING]` |
| Tool group system (`group:fs`, `group:runtime`, etc.) | `[MISSING]` |
| Tool permission allowlists per agent | `[MISSING]` |
| Skill install/management platform | `[MISSING]` |
| ClawHub-style skill registry & discovery | `[MISSING]` |
| Bundled, managed, and workspace skill tiers | `[MISSING]` |
| Skill install gating & approval | `[MISSING]` |
| MCP (Model Context Protocol) server support | `[MISSING]` |
| 700+ community skills ecosystem | `[MISSING]` |

---

## 4. Browser Automation

Entirely absent from Dhaavak.

| Feature | Status |
|---------|--------|
| Dedicated Chrome/Chromium instance with CDP | `[MISSING]` |
| Page snapshots & screenshots | `[MISSING]` |
| Browser actions (click, type, drag, navigate) | `[MISSING]` |
| Browser profiles | `[MISSING]` |

---

## 5. Canvas & A2UI (Agent-to-UI)

Entirely absent from Dhaavak.

| Feature | Status |
|---------|--------|
| Live Canvas visual workspace | `[MISSING]` |
| Agent-driven UI rendering | `[MISSING]` |
| Canvas evaluation & snapshots | `[MISSING]` |
| A2UI protocol | `[MISSING]` |

---

## 6. Voice & Speech

Entirely absent from Dhaavak.

| Feature | Status |
|---------|--------|
| Voice Wake (always-on speech recognition) | `[MISSING]` |
| Talk Mode (continuous conversation overlay) | `[MISSING]` |
| ElevenLabs TTS integration | `[MISSING]` |
| Audio transcription hooks | `[MISSING]` |
| Audio pipeline with size caps & lifecycle management | `[MISSING]` |

---

## 7. Memory & Persistence

Dhaavak has in-memory session history only. OpenClaw has persistent, searchable memory.

| Feature | Status |
|---------|--------|
| Persistent memory backend | `[MISSING]` |
| Memory search (`memory_search`) | `[MISSING]` |
| Memory retrieval (`memory_get`) | `[MISSING]` |
| Voyage AI embeddings for memory | `[MISSING]` |
| Pluggable memory backends | `[MISSING]` |
| Cross-session memory persistence | `[MISSING]` |
| Context compaction (`/compact` command) | `[MISSING]` |

---

## 8. Session Management (Advanced)

Dhaavak has basic session lifecycle. OpenClaw has inter-session coordination.

| Feature | Status |
|---------|--------|
| `sessions_list` — query active sessions | `[MISSING]` |
| `sessions_history` — retrieve session history | `[MISSING]` |
| `sessions_send` — inter-session messaging | `[MISSING]` |
| `sessions_spawn` — spawn sub-agent sessions | `[MISSING]` |
| Session pruning with context management | `[MISSING]` |
| Session status queries | `[MISSING]` |
| Docker sandbox for non-main sessions | `[MISSING]` |

---

## 9. Scheduling & Automation

Entirely absent from Dhaavak.

| Feature | Status |
|---------|--------|
| Cron job scheduling | `[MISSING]` |
| Webhook triggers | `[MISSING]` |
| Gmail Pub/Sub integration | `[MISSING]` |
| Wakeup triggers | `[MISSING]` |

---

## 10. Companion Applications & Device Nodes

Entirely absent from Dhaavak.

| Feature | Status |
|---------|--------|
| macOS menu bar app | `[MISSING]` |
| iOS node (Canvas, Voice Wake, camera) | `[MISSING]` |
| Android node (Canvas, Talk Mode, camera, SMS) | `[MISSING]` |
| Device node protocol with advertised capabilities | `[MISSING]` |
| Bonjour pairing for local nodes | `[MISSING]` |
| Camera snap/clip capture | `[MISSING]` |
| Screen recording | `[MISSING]` |
| Location retrieval | `[MISSING]` |
| Push notifications to nodes | `[MISSING]` |

---

## 11. Cross-Channel Messaging

Dhaavak routes inbound messages but cannot send cross-channel.

| Feature | Status |
|---------|--------|
| `message` tool — send to any channel from agent | `[MISSING]` |
| Cross-channel message routing from within tools | `[MISSING]` |
| Agent-initiated outbound messages | `[MISSING]` |

---

## 12. Productivity Integrations

Entirely absent from Dhaavak.

| Feature | Status |
|---------|--------|
| Apple Notes | `[MISSING]` |
| Apple Reminders | `[MISSING]` |
| Things 3 | `[MISSING]` |
| Notion | `[MISSING]` |
| Obsidian | `[MISSING]` |
| Bear Notes | `[MISSING]` |
| Trello | `[MISSING]` |
| GitHub | `[MISSING]` |
| Email management | `[MISSING]` |
| Twitter/X posting | `[MISSING]` |

---

## 13. Entertainment & Smart Home

Entirely absent from Dhaavak.

| Feature | Status |
|---------|--------|
| Spotify playback control | `[MISSING]` |
| Sonos multi-room audio | `[MISSING]` |
| Shazam song recognition | `[MISSING]` |
| Philips Hue lighting | `[MISSING]` |
| 8Sleep mattress control | `[MISSING]` |
| Home Assistant hub | `[MISSING]` |
| Weather data | `[MISSING]` |

---

## 14. Media & Creative

Entirely absent from Dhaavak.

| Feature | Status |
|---------|--------|
| AI image generation | `[MISSING]` |
| GIF search | `[MISSING]` |
| Screen capture (Peekaboo) | `[MISSING]` |

---

## 15. Security & Permissions (Advanced)

Dhaavak has basic token auth and allowlists. OpenClaw has a deeper permission model.

| Feature | Status |
|---------|--------|
| DM pairing codes for unknown senders | `[MISSING]` |
| Elevated bash access toggle (per-session) | `[MISSING]` |
| macOS TCC permission integration | `[MISSING]` |
| Tool permission profiles & per-agent allowlists | `[MISSING]` |
| 1Password credential management integration | `[MISSING]` |
| Safety scanner | `[MISSING]` |

---

## 16. Deployment & Infrastructure

Dhaavak has a single binary. OpenClaw has multi-mode deployment.

| Feature | Status |
|---------|--------|
| Docker sandbox support | `[MISSING]` |
| Nix declarative configuration | `[MISSING]` |
| Tailscale Serve / Funnel exposure | `[MISSING]` |
| Remote gateway via SSH tunnel | `[MISSING]` |
| Daemon installation (launchd / systemd) | `[MISSING]` |
| Interactive CLI setup wizard | `[MISSING]` |
| Release channels (stable / beta / dev) | `[MISSING]` |

---

## 17. Chat Commands

Dhaavak has `ping` and `chat.send`. OpenClaw has a richer command set.

| Feature | Status |
|---------|--------|
| `/status` — session info | `[MISSING]` |
| `/new` / `/reset` — new session | `[MISSING]` |
| `/compact` — context summary/compression | `[MISSING]` |
| `/think` — reasoning level control | `[MISSING]` |
| `/verbose` — toggle verbose mode | `[MISSING]` |
| `/usage` — token/cost tracking | `[MISSING]` |
| `/restart` — owner-only restart | `[MISSING]` |
| `/activation` — group toggle | `[MISSING]` |

---

## 18. Observability & Usage Tracking

Dhaavak has structured logging. OpenClaw adds usage metrics.

| Feature | Status |
|---------|--------|
| Token count tracking | `[MISSING]` |
| Cost monitoring | `[MISSING]` |
| Typing indicators | `[MISSING]` |
| Debug tools (macOS app) | `[MISSING]` |

---

## 19. Agent Workspace & Prompt Management

Dhaavak has per-agent system prompts. OpenClaw has a file-based prompt injection system.

| Feature | Status |
|---------|--------|
| `AGENTS.md` workspace prompt injection | `[MISSING]` |
| `SOUL.md` personality prompt file | `[MISSING]` |
| `TOOLS.md` tool instruction file | `[MISSING]` |
| Multi-agent workspace isolation | `[PARTIAL]` — multi-agent routing exists, but no workspace isolation |

---

## 20. Multi-Agent Coordination

Dhaavak has per-agent routing. OpenClaw has full multi-agent orchestration.

| Feature | Status |
|---------|--------|
| `agents_list` — discover available sub-agents | `[MISSING]` |
| Sub-agent spawning from within a session | `[MISSING]` |
| Inter-agent message passing | `[MISSING]` |
| Per-workspace agent assignment | `[PARTIAL]` — per-channel bindings exist |

---

## Summary: Priority Tiers

### Tier 1 — Core Platform Gaps (High Impact)
1. Additional LLM providers (OpenAI, Gemini, Ollama) + failover
2. Persistent memory backend with search
3. Built-in tool suite (shell exec, file ops, web search, web fetch)
4. Skill/plugin system with registry
5. MCP server support
6. Chat commands (`/status`, `/new`, `/compact`, `/usage`)
7. Token/cost tracking
8. Context compaction for long conversations

### Tier 2 — Channel Expansion (Reach)
9. WhatsApp adapter
10. Slack adapter
11. Discord adapter
12. Signal adapter
13. iMessage/BlueBubbles adapter
14. Microsoft Teams adapter

### Tier 3 — Advanced Agent Features
15. Inter-session messaging & sub-agent spawning
16. Cross-channel outbound messaging from tools
17. Browser automation (CDP)
18. Cron scheduling & webhook triggers
19. Agent workspace prompt files (AGENTS.md, SOUL.md)
20. Advanced permission profiles & tool allowlists

### Tier 4 — Platform & UX
21. Voice Wake & Talk Mode
22. Canvas & A2UI
23. Companion apps (macOS, iOS, Android)
24. Device node protocol
25. Docker sandbox for sessions
26. Deployment tooling (Tailscale, Nix, daemon install)
27. Interactive setup wizard

### Tier 5 — Integrations Ecosystem
28. Productivity (Notion, Obsidian, GitHub, Apple Notes, etc.)
29. Smart home (Hue, Home Assistant, 8Sleep)
30. Entertainment (Spotify, Sonos, Shazam)
31. Media (AI image gen, GIF search, screen capture)
32. Security (1Password, safety scanner)

---

*Generated 2026-02-09. Based on OpenClaw v2026.2.x feature set.*
