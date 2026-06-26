# Ignite — Implementation Plan
> Provisioning with a heartbeat

**Tech Stack:**

| Layer              | Choice                                                     |
| :----------------- | :--------------------------------------------------------- |
| **Desktop Shell**  | Wails v2.12.0 (Go + native WebView)                        |
| **Frontend**       | React 19 + TypeScript + Vite                               |
| **Styling**        | Tailwind CSS v4                                            |
| **UI Components**  | shadcn/ui (new-york style)                                  |
| **State**          | Zustand                                                    |
| **Markdown**       | react-markdown + react-syntax-highlighter                  |
| **Font**           | JetBrains Mono                                             |
| **Icons**          | lucide-react                                               |
| **Backend**        | Go 1.23+                                                   |
| **LLM Providers**  | OpenCode, Claude, DeepSeek, GitHub Copilot (HTTP + SSE)    |
| **Templates**      | Go text/template + sprig + embed                           |
| **Database**       | SQLite (modernc.org/sqlite, pure Go, no CGO)               |
| **Logging**        | zerolog                                                    |
| **Secrets**        | OS keychain (go-keyring)                                   |
| **Package Manager**| pnpm (frontend), Go modules (backend)                      |

---

## Phase 0: Project Foundation

> Scaffold the Wails v2 desktop shell with React frontend skeleton, theme system, and settings infrastructure.

- [ ] `main.go` — Wails app entry, embedded frontend assets, native menu bar (File/Edit), Mac AboutInfo
- [ ] `app.go` — App struct with context lifecycle, 17 Wails bindings for frontend IPC
- [ ] `wails.json` — Wails v2 project config with frontend build paths
- [ ] `frontend/index.html` — JetBrains Mono Google Fonts link, FOUC prevention inline script
- [ ] `frontend/src/style.css` — Tailwind v4 @theme with dark/light palettes, @keyframes for slide-up and blink
- [ ] `frontend/src/App.tsx` — Root layout: Sidebar (260px) + ChatPanel (flex) + StatusBar (fixed bottom)
- [ ] `lib/store/theme.ts` — Zustand store: dark/light toggle, `data-mode` attribute, localStorage persistence
- [ ] `lib/store/chat.ts` — Zustand store: projects, messages, activeProjectId, streaming state
- [ ] `hooks/useTheme.ts` — Init theme on mount, expose mode and toggle
- [ ] `hooks/useConversation.ts` — Stream subscription via EventsOn, message persistence, path scanner injection
- [ ] `internal/settings/config.go` — Config struct, Load/Save to `~/.ignite/config.json`, defaults
- [ ] `internal/settings/keychain.go` — OS keychain CRUD for API keys via go-keyring

### Verification

- [ ] `wails doctor` passes all checks
- [ ] `go test ./internal/settings/` — 2 tests: DefaultConfig, SaveAndLoad
- [ ] `pnpm vitest run` — 6 theme store tests pass
- [ ] `pnpm typecheck` — TypeScript clean
- [ ] `make publish` — produces `build/Ignite.{dmg,zip}`

---

## Phase 1: Backend Services

> Build the LLM provider abstraction, SQLite history, template engine, and conversation orchestration.

- [ ] `internal/providers/interface.go` — LLMProvider interface (Chat/ChatStream/ListModels/ValidateKey), Message/Model/ChatResponse types
- [ ] `internal/providers/opencode.go` — OpenCode adapter (OpenAI-compatible HTTP, SSE streaming via bufio.Scanner)
- [ ] `internal/providers/claude.go` — Claude adapter (Anthropic Messages API, x-api-key, content_block_delta events)
- [ ] `internal/providers/deepseek.go` — DeepSeek adapter (OpenAI-compatible, hardcoded model list)
- [ ] `internal/providers/github.go` — GitHub Copilot adapter (Bearer token, chat/stream endpoint)
- [ ] `internal/providers/manager.go` — Provider registry (Register/Get)
- [ ] `internal/history/models.go` — Project and Message structs
- [ ] `internal/history/sqlite.go` — SQLite store (OpenDB with WAL, migration, CRUD for projects/messages)
- [ ] `internal/templates/data.go` — ProjectContext with all fields (Phases, TechStack, Dependencies, APIs, EnvVars, Theme, etc.)
- [ ] `internal/templates/engine.go` — Go text/template + sprig FuncMap, NewEngine, Generate returning ProjectFiles
- [ ] `internal/templates/templates.go` — `//go:embed templates/*`, EmbeddedEngine loader
- [ ] `internal/templates/templates/*.tmpl` — 4 embedded template files (project.md, AGENTS.md, PLAN.md, README.md)
- [ ] `internal/scanner/scanner.go` — Universal path analyzer detecting 15+ project types by file presence

### Verification

- [ ] `go test ./internal/history/` — 2 tests: ProjectCRUD, MessageCRUD
- [ ] `go test ./internal/providers/` — 9 tests across all 4 adapters (Chat + Stream + ListModels)
- [ ] `go test ./internal/templates/` — Engine generates all 4 output files with correct content
- [ ] `go test ./internal/scanner/` — Scanner detects Go, Node, Python, Rust, and empty directories

---

## Phase 2: Frontend UI Components

> Build the shadcn/ui component library, sidebar, chat panel, status bar, and settings modal.

- [ ] `components.json` — shadcn/ui init (base-nova style, lucide icons)
- [ ] `components/ui/button.tsx` — shadcn Button
- [ ] `components/ui/dialog.tsx` — shadcn Dialog (used by Settings and About)
- [ ] `components/ui/input.tsx` — shadcn Input
- [ ] `components/ui/label.tsx` — shadcn Label
- [ ] `components/ui/select.tsx` — shadcn Select (used throughout)
- [ ] `components/sidebar/sidebar.tsx` — 260px sidebar with Flame logo, project list, New Project button, project resume on click
- [ ] `components/chat/chat-bubble.tsx` — Markdown rendering with react-markdown, oneDark syntax highlighting, streaming cursor
- [ ] `components/chat/chat-input.tsx` — Auto-resize textarea, Enter/Shift+Enter, disabled when no project/provider
- [ ] `components/chat/chat-panel.tsx` — Message list with auto-scroll, scroll-to-bottom button, welcome banner
- [ ] `components/status-bar/status-bar.tsx` — Provider/model dropdowns, connection indicator (HasAPIKey), Settings/Export buttons
- [ ] `components/settings/settings-modal.tsx` — Two-tab dialog: Providers (API keys with show/hide/validate) + Appearance (theme/license/dir/font dropdowns)
- [ ] `components/settings/about-modal.tsx` — About dialog with author, web, repo links

### Verification

- [ ] `pnpm typecheck` — TypeScript clean
- [ ] `wails dev` renders all components with correct dark/light styling
- [ ] All Wails bindings called: GetSettings, SaveSettings, SetAPIKey, HasAPIKey, ValidateProviderKey, ListProviderModels, CreateProject, ListProjects, GetProject, GetMessages, AddMessage, SendMessage, SendMessageStream, ExportChat, AnalyzePath

---

## Phase 3: Polish & Quality

> Scrollbar styling, animations, keyboard shortcuts, menu bar, theme fixes, icon generation.

- [ ] Scrollbar styling — thin 6px thumb with `--border` color
- [ ] Chat bubble slide-up animation (0.2s ease-out)
- [ ] Streaming cursor blink animation (0.8s step-end)
- [ ] Keyboard shortcuts — Escape to close Settings, Ctrl/Cmd+Enter to send
- [ ] Native menu bar — File (New Project Cmd+N, Export Cmd+E, Settings Cmd+,), Edit (standard), macOS AppMenu (About, Quit)
- [ ] Theme live update — Settings theme dropdown immediately toggles data-mode and .dark class
- [ ] Font live update — Settings font dropdown updates --font-mono CSS variable on change
- [ ] Custom app icon — flame teardrop icon (PIL-generated, 7 icns sizes, rounded corners)

### Verification

- [ ] `make publish` produces DMG ~6.7M, ZIP ~5.9M
- [ ] Go vet clean (only expected embed error for frontend/dist)
- [ ] All 14 Go tests + 6 vitest tests pass
- [ ] TypeScript typechecks with zero errors

---

## Performance Targets

| Metric          | Target            | Status  |
| :-------------- | :---------------- | :------ |
| App launch      | < 1 second        | ✓       |
| LLM response    | Streaming, FTT < 2s | ✓ (SSE) |
| Memory idle     | < 200MB           | ✓       |
| Binary size     | < 80MB            | ✓ (14MB)|
| History search  | < 50ms            | ✓ (WAL) |
| Build time      | < 15s             | ✓ (4.4s)|

---

## Risks & Mitigations

| Risk                          | Mitigation                                              |
| :---------------------------- | :------------------------------------------------------ |
| LLM output quality varies     | Template backbone guarantees structure; conversation context enriches |
| Streaming breaks on provider  | Graceful fallback to non-streaming; SSE resilience      |
| OS keychain unavailable       | Error surfaced in UI; fallback to encrypted file        |
| Wails v2 API changes          | Pinned to v2.12.0 in go.mod                             |
| Template size grows           | Modular template sources; LLM extends beyond template   |
| WebView CSS inconsistencies   | System WebView (macOS) tested; Tailwind v4 with CSS custom properties |

---

## Changelog

| Version | Date       | Changes                                                  |
| :------ | :--------- | :------------------------------------------------------- |
| 0.1.0   | 2026-06-26 | Initial MVP — Wails v2 scaffold, 4 LLM providers, SQLite history, 4 output templates, full React UI, menu bar, theme system, path scanner |

---

## Versioning Strategy

Semantic versioning (MAJOR.MINOR.PATCH). MVP ships as 0.1.0. Pre-1.0 releases may break API. Tag format: `v0.1.0`.
