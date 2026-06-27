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

### v0.1.4 (current)

- System prompt injected into AI conversations (5-phase interview guide)
- JSON project name parsing — AI auto-renames "Untitled" projects
- Inline project rename via pencil icon on hover (sidebar)
- Editable API key field + error feedback on settings save
- Connection indicator refreshes immediately after saving API key
- Open Folder (Cmd+O) with native dialog + tilde expansion fix
- Current project directory displayed in status bar
- Share API key between OpenCode Go and Zen automatically
- Delete project from sidebar (Trash2 icon on hover)
- Window size presets: Small, Medium, Large, 4K, Full Screen
- User profile: display name + avatar upload in Preferences
- Bot lucide icon for AI in chat, user avatar for your messages
- Removed extra fonts (JetBrains Mono only), DMG: 7.2M → 6.8M
- Path scanner reads file contents for referenced projects
- Corrected OpenCode Go/Zen model lists from their actual APIs
- Removed Claude as standalone provider
- Preferences tab (renamed from Appearance) with folder picker
- Providers docs page shows separate Go and Zen model lists
- Deploy target in Makefile + local deploy script for Caddy

### v0.1.3

- Removed Claude provider (3 providers: OpenCode Go, OpenCode Zen, DeepSeek)
- Open Folder (Cmd+O) with native macOS dialog
- Current project directory displayed in status bar
- Delete project from sidebar (Trash2 icon on hover)
- About modal centering fixes

<details>
<summary>v0.1.2</summary>

- Landing site at ignite.johansen.foo with terminal mockup hero + light/dark mode
- 7 documentation pages (Overview, Getting Started, Providers, Provisioning, Settings, FAQ, Developer Setup)
- Developer setup guide with 67 skill links and resource references
- SEO optimization (Open Graph, JSON-LD, sitemap, robots.txt, canonical URLs)
- Favicon ecosystem (ICO, SVG, PNG, Apple Touch Icon, Android Chrome icons)
- Saved default model now auto-selects on provider switch (fixes all providers)
- Settings: save overwrite bug fixed, locked API key field when key is in keychain
- Settings: model dropdown shows display names, auto-selects first model
- Settings: theme labels "Dark"/"Light" (capitalized)
- About modal: centered tagline and version, removed Kvasir reference
- Consistent docs sidebar with all 7 links on every page
- App name casing: Ignite (capital I) throughout .app bundle, binary, DMG
- FUNDING.yml with Buy Me a Coffee support
- AGPL-3.0 LICENSE file

<details>
<summary>v0.1.1</summary>

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

<details>
<summary>v0.1.0</summary>

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

</details>
