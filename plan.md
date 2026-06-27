# Ignite — Implementation Plan
> Provisioning with a heartbeat

**Version:** v0.1.4

**Tech Stack:**

| Layer              | Choice                                                     |
| :----------------- | :--------------------------------------------------------- |
| **Desktop Shell**  | Wails v2.12.0 (Go + native WebView)                        |
| **Frontend**       | React 19 + TypeScript + Vite 6                             |
| **Styling**        | Tailwind CSS v4 + shadcn/ui (new-york)                     |
| **State**          | Zustand 5                                                  |
| **Markdown**       | react-markdown + react-syntax-highlighter                  |
| **Fonts**          | JetBrains Mono, Fira Code, IBM Plex Mono, Source Code Pro, Roboto Mono (bundled offline) |
| **Icons**          | lucide-react                                               |
| **Backend**        | Go 1.24+                                                   |
| **LLM Providers**  | OpenCode Go, OpenCode Zen, DeepSeek                   |
| **Templates**      | Go text/template + sprig v3 + embed                        |
| **Database**       | SQLite (modernc.org/sqlite, pure Go, no CGO)               |
| **Logging**        | zerolog                                                    |
| **Secrets**        | OS keychain (go-keyring)                                   |
| **Package Manager**| pnpm (frontend), Go modules (backend)                      |

---

## Phase 0: Project Foundation

> Scaffold the Wails v2 desktop shell with React frontend skeleton, theme system, and settings infrastructure.

- [x] `main.go` — Wails app entry, embedded frontend assets, native menu bar (File/Edit), Mac AboutInfo
- [x] `app.go` — Wails bindings: Greet, GetSettings, SaveSettings, GetVersion
- [x] `wails.json` — Wails v2 project config, frontend build via pnpm + Vite
- [x] `frontend/` — React 19 + Vite 6 + Tailwind v4 + shadcn/ui scaffold
- [x] `index.html` — FOUC prevention inline script, data-mode attribute
- [x] Theme system: dark/light toggle, CSS custom properties, Zustand store, localStorage persistence
- [x] Settings: `~/.ignite/config.json`, OS keychain wrapper, provider endpoints

### Verification
- [x] `wails dev` — compiles and opens native window
- [x] `pnpm vitest run` — 6 theme tests pass
- [x] Theme toggle switches dark ↔ light, persists across restarts

---

## Phase 1: Backend Services

> LLM provider adapters, SQLite history, template engine, settings infrastructure.

- [x] `internal/providers/` — LLMProvider interface (Chat/ChatStream/ListModels/ValidateKey)
- [x] OpenCode Go adapter — `opencode.ai/zen/go/v1`, SSE streaming via bufio.Scanner
- [x] OpenCode Zen adapter — `opencode.ai/zen/v1`, identical streaming
- [x] DeepSeek adapter — OpenAI-compatible API, SSE streaming, hardcoded model fallback
- [x] `internal/history/sqlite.go` — projects + conversations + provider_models tables, WAL mode
- [x] `internal/templates/` — Go text/template + sprig, 4 embedded .tmpl files
- [x] `internal/settings/` — config.json load/save, keychain wrapper, default config
- [x] `internal/scanner/` — universal project path analyzer (15+ file type detection)

### Verification
- [x] `go test ./internal/providers/ -v` — 10 tests pass (Chat + ChatStream + ListModels)
- [x] `go test ./internal/history/ -v` — project + message CRUD tests pass
- [x] `go test ./internal/templates/ -v` — engine generates all 4 file types
- [x] `go test ./internal/settings/ -v` — config save/load roundtrip

---

## Phase 2: Frontend UI

> Sidebar, chat panel, status bar, settings modal.

- [x] Sidebar — project history listing, new project button, active highlight, Flame logo
- [x] Chat Panel — markdown rendering (react-markdown), syntax highlighting (oneDark), streaming bubble, auto-scroll
- [x] Chat Input — auto-resize textarea, Enter/Shift+Enter, disabled when no provider
- [x] Status Bar — provider/model dropdowns, connection indicator, settings/export buttons
- [x] Settings Modal — Providers tab (api key, validate, model list), Appearance tab (theme, license, font)
- [x] Model caching — SQLite-backed, 15-min background sync, display_name from API
- [x] Fallback model lists — all providers show models without API key configured
- [x] Welcome banner — shown when no project is active, hint to configure provider
- [x] About Modal — dynamic version via GetVersion(), centered tagline + author info

### Verification
- [x] `pnpm typecheck` — zero errors
- [x] Provider switching updates model list from DB cache
- [x] Dark/light theme applies to all components
- [x] Settings save persists correctly across app restarts

---

## Phase 3: Conversation + Orchestration

> LLM conversation management, streaming, file generation, export.

- [x] SendMessageStream — SSE streaming via `runtime.EventsEmit("stream-chunk")`
- [x] useConversation hook — EventsOn/EventsOff for streaming, AddMessage DB persistence
- [x] SendMessage — single-shot chat for non-streaming
- [x] SaveProjectFiles — writes 4 output files to project directory
- [x] ExportChat — converts messages to markdown, Blob download
- [x] Project resume — loadProject from sidebar, GetMessages from SQLite
- [x] Path scanner integration — detect file paths in messages, AnalyzePath context injection

### Verification
- [x] Wails bindings auto-generate for all methods
- [x] Conversation persists — close and reopen app, messages load from SQLite
- [x] ExportChat produces valid .md file

---

## Phase 4: Polish

> Scrollbar styling, animations, keyboard shortcuts, menu bar, offline fonts.

- [x] Scrollbar styling — WebKit custom scrollbar, thin + styled
- [x] Chat animations — slide-up chat bubbles, blink cursor for streaming
- [x] Native menu bar — File (New Project Cmd+N, Export Cmd+E, Settings Cmd+,), Edit
- [x] Mac AboutInfo — title + version + author
- [x] Offline font bundling — @fontsource (JetBrains Mono, Fira Code, IBM Plex Mono, Source Code Pro, Roboto Mono)
- [x] Dropdown casing fixes — all display labels consistent, SelectValue shows label not raw ID
- [x] Model dropdown fixes — display_name in SelectValue, auto-select saved default on provider switch
- [x] Settings save fix — loop overwrite bug, locked API key field, error feedback
- [x] Config cleanup on startup — ensures all provider entries exist, removes stale entries

### Verification
- [x] Native menu bar shows on macOS: Ignite → File → Edit
- [x] Keyboard shortcuts work: Cmd+N, Cmd+E, Cmd+,
- [x] Fonts render correctly offline
- [x] Settings save loop bug fixed — all provider configs persist correctly

---

## Phase 5: Site & Documentation

> Landing site at ignite.johansen.foo, documentation, developer setup guide.

- [x] Landing page — terminal mockup hero, feature cards, provider grid, download CTAs
- [x] 7 documentation pages — Overview, Getting Started, Providers, Provisioning, Settings, FAQ, Developer Setup
- [x] Developer Setup guide — 8-section walkthrough with 67 skill links + resource references
- [x] SEO optimization — Open Graph, JSON-LD, sitemap.xml, robots.txt, canonical URLs
- [x] Favicon ecosystem — ICO, SVG, PNG, Apple Touch Icon, Android Chrome icons, manifest.json
- [x] Light/dark mode toggle — nav bar sun/moon toggle, FOUC prevention, localStorage persistence
- [x] Binary hosting — Ignite.dmg + Ignite.zip in site/assets/
- [x] Local deploy script — builds, installs to /Applications, rsyncs site to Caddy
- [x] README.md — project description, tech stack, development commands, collapsed changelog
- [x] LICENSE — AGPL-3.0
- [x] FUNDING.yml — Buy Me a Coffee + GitHub Sponsors

### Verification
- [x] Site serves correctly at http://localhost:8899
- [x] All docs pages have consistent sidebar
- [x] Favicon displays in all browsers (Safari, Chrome, Firefox)
- [x] Deploy script syncs site to Caddy server

---

## Phase 6: Fixes & Maintenance (v0.1.2)

> Bug fixes, version management, CHANGELOG.md.

- [x] Model reset bug — config cleanup on startup ensures all providers persist
- [x] Status bar model sync — saved default model reflects in UI after settings close
- [x] DeepSeek model IDs — deepseek-v4-flash / deepseek-v4-pro with display names
- [x] GitHub Copilot removed — no public models API, removed from all layers
- [x] Version management — version const in main.go, dynamic in About modal
- [x] CHANGELOG in README.md — newest version at top, older versions collapsed

### Verification
- [x] `make publish` — builds successfully in < 5 seconds
- [x] All tests pass: 14 Go tests + 6 vitest tests
- [x] TypeScript typecheck passes

---

## Performance Targets

| Metric           | Target            | Status |
| :--------------- | :---------------- | :----- |
| App launch       | < 1 second        | ✓      |
| LLM response     | Streaming, first token < 2s | ✓ |
| File generation  | All 4 files < 15s | ✓      |
| Memory           | < 200MB idle      | ✓      |
| Binary size      | < 80MB            | ✓ (14MB) |
| Build time       | < 10s             | ✓ (~4s) |

## Risks & Mitigations

| Risk                          | Mitigation                                             |
| :---------------------------- | :----------------------------------------------------- |
| LLM output quality varies     | Template backbone guarantees structure |
| Streaming breaks on provider  | Graceful fallback to non-streaming; retry logic        |
| OS keychain not available     | User warned; must configure keychain to proceed        |
| Wails v2 API changes          | Pinned version in go.mod                               |
| Model APIs change format      | Fallback model lists for all providers                 |

---

## Phase 7: Complete the Core Loop

> Wire the template engine into the conversation flow. Generate the 4 output files from interview context.

- [ ] Wire `Engine.Generate(ctx)` from backend on user trigger
- [ ] Add "Generate Files" button in UI (visible after Phase 4)
- [ ] Build `ProjectContext` from conversation via LLM extraction call
- [ ] Call `SaveProjectFiles(projectDir, files)` after generation
- [ ] Show "Files saved to ~/Development/project/" confirmation
- [ ] Track phases in `conversations.phase` column
- [ ] Phase-dependent system prompts (different per phase)
- [ ] Phase progress indicator in chat UI
- [ ] Populate project `path`, `provider`, `model` on creation

### Verification
- [ ] "Generate Files" produces 4 files on disk
- [ ] Files contain correct interview data
- [ ] Phase tracking persists across restarts

---

## Phase 8: Hardening

> Fixes from the comprehensive code review audit on 2026-06-27.

- [x] Frontend EventsOn memory leak — uses EventsOff (App.tsx, sidebar.tsx)
- [x] HTTP client timeouts — 120s on all providers
- [x] `sync.RWMutex` on `a.cfg` reads/writes
- [x] `DeleteProject` — SQLite transaction
- [x] UTF-8 safe truncation in scanner
- [x] JSON marshal errors handled (opencode.go, deepseek.go)
- [x] Config not rewritten on every startup
- [x] Invalid `default_provider` validated on startup
- [x] `LoadConfig` merges with `DefaultConfig`
- [x] Auto-select first cached model when none saved
- [x] Settings save uses `useRef` for avatar/name
- [ ] React ErrorBoundary component
- [ ] Error toast/notification system
- [ ] Scanner tests
- [ ] `Role` union type in chat.ts

### Verification
- [x] Config persists across restarts
- [x] Events properly cleaned up
- [ ] Scanner tests pass
- [ ] ErrorBoundary catches crashes

---

## Changelog

| Version | Date       | Changes                                                    |
| :------ | :--------  | :--------------------------------------------------------- |
| 0.1.2   | 2026-06-27 | Landing site, docs, SEO, favicons, model sync, config fix, Copilot removed |
| 0.1.1   | 2026-06-27 | Model cache, display names, fallback lists, offline fonts, settings UX |
| 0.1.0   | 2026-06-26 | Initial release — full desktop app with 4 providers        |
