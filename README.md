# Ignite

> Provisioning with a heartbeat

A desktop GUI for provisioning new software
documentation quality. Ignite conducts AI-guided conversational interviews and
generates project specs, agent guides, implementation plans, and READMEs.

## Features

- **5-Phase AI Interview** — Identity → Tech Stack → Features → Architecture → Generation
- **Multi-Provider** — OpenCode Go, OpenCode Zen, Claude, DeepSeek
- **Template-Driven Output** — project.md, agents.md, plan.md, README.md
- **Dark/Light Theme** — Nordic-inspired dark palette, system sans-serif, offline fonts
- **Project History** — SQLite persistence, resume mid-project from sidebar
- **Model Cache** — Background sync every 15 minutes, instant dropdowns
- **Native macOS** — Menu bar, keyboard shortcuts, signed .app bundle, flame icon

## Quick Start

```sh
git clone git@github.com:mojoaar/ignite.git
cd ignite

# Frontend
cd frontend && pnpm install

# Backend
cd ..
go mod download

# Development
wails dev

# Production build
make publish
```

## Tech Stack

| Layer              | Choice                                          |
| ------------------ | ----------------------------------------------- |
| Desktop Shell      | Wails v2.12.0 (Go + native WebView)             |
| Frontend           | React 19, TypeScript, Vite 6                    |
| Styling            | Tailwind CSS v4, shadcn/ui (new-york)           |
| State              | Zustand 5                                       |
| Markdown Rendering | react-markdown + react-syntax-highlighter       |
| Icons              | lucide-react (Flame logo)                       |
| Backend            | Go 1.24+                                        |
| Template Engine    | Go text/template + sprig v3 + embed             |
| Database           | SQLite (modernc.org/sqlite, pure Go, no CGO)    |
| Secret Storage     | OS Keychain (go-keyring)                        |
| Logging            | zerolog                                         |
| Fonts              | JetBrains Mono, Fira Code, IBM Plex Mono, Source Code Pro, Roboto Mono (bundled) |

## Development

| Command      | Description           |
| ------------ | --------------------- |
| `wails dev`    | Dev server with HMR   |
| `make test`    | Run all tests         |
| `make publish` | Build + DMG + ZIP     |
| `make clean`   | Remove build artifacts |

## License

AGPL-3.0 — see [LICENSE](LICENSE)

## Changelog

### v0.1.1 (current)

- Model cache with 15-minute background sync
- Display names in model dropdowns (e.g. "Claude Sonnet 4" not "claude-sonnet-4")
- Fallback model lists for all providers (no API key needed to browse)
- Unified flame icon — matching SVG in-app and .app bundle icon
- Offline font bundling via @fontsource (JetBrains Mono, Fira Code, IBM Plex Mono, Source Code Pro, Roboto Mono)
- Dynamic version in About dialog
- "Set API key" hint for empty model dropdowns
- Welcome banner when no project selected
- Disabled chat input until provider configured
- Native macOS menu bar (File > Settings, Cmd+,)
- App name casing: Ignite (capital I) everywhere
- Removed GitHub Copilot provider (no public models API)
- Settings UX: dropdown selectors for model, license, font
- Universal project path scanner for conversation context

### v0.1.0

- Initial release
- Wails v2 desktop shell + React 19 frontend
- 5-phase AI interview: Identity → Tech Stack → Features → Architecture → Generation
- 4 LLM providers: OpenCode Go, OpenCode Zen, Claude, DeepSeek
- Go template engine with 4 embedded output templates (project.md, agents.md, plan.md, README.md)
- Dark/Light theme system with CSS custom properties
- SQLite project history with conversation resume
- Settings system with OS keychain for API keys
- Provider validation with green/red connection indicator
- Streaming chat with SSE, markdown rendering, syntax highlighting
- Status bar with live provider/model switching
- Project sidebar with one-click resume
- Chat export as markdown
- Keyboard shortcuts (Enter/Shift+Enter to send)
- Scroll-to-bottom button, slide-up animations, blinking cursor
