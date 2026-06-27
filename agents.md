# Ignite — Agent Guide
> **Project:** Ignite — Provisioning with a heartbeat

## Quick Commands

| Action     | Command                                          |
| :--------- | :----------------------------------------------- |
| Dev        | `wails dev`                                        |
| Build      | `make publish`                                     |
| Test (Go)  | `go test ./internal/...`                           |
| Test (FE)  | `cd frontend && pnpm vitest run`                   |
| Typecheck  | `cd frontend && pnpm typecheck`                    |
| Deploy site| `bash scripts/local-deploy.sh`                     |
| Clean      | `rm -rf build frontend/dist frontend/node_modules` |

---

## Tech Stack

| Layer              | Choice                                                |
| ------------------ | ----------------------------------------------------- |
| Desktop Shell      | Wails v2.12.0 (Go + native WebView)                   |
| Frontend           | React 19 + TypeScript + Vite 6                        |
| Styling            | Tailwind CSS v4 + shadcn/ui (new-york)                |
| State              | Zustand 5                                             |
| Markdown Rendering | react-markdown + react-syntax-highlighter (oneDark)   |
| Icons              | lucide-react (Flame logo)                              |
| Fonts              | JetBrains Mono, Fira Code, IBM Plex Mono, Roboto Mono, Source Code Pro (bundled via @fontsource) |
| Backend            | Go 1.24+                                              |
| LLM Providers      | OpenCode Go, OpenCode Zen, DeepSeek               |
| Template Engine    | Go text/template + sprig v3 + embed                   |
| Database           | SQLite (modernc.org/sqlite, pure Go)                  |
| Secret Storage     | OS keychain (go-keyring)                              |
| Logging            | zerolog                                                |

---

## Code Rules

### Always

- Run `pnpm typecheck` and `go vet ./internal/...` before committing
- Run `go test ./internal/...` and `cd frontend && pnpm vitest run` before committing
- Push to GitHub after every work session
- Bump version in `main.go` when making meaningful changes
- Update `CHANGELOG` in `README.md` when bumping version
- Use snake_case for Wails binding types (they match Go JSON tags)
- Check existing patterns before writing new code — mimic existing file conventions
- Use `cp -R build/bin/Ignite.app /Applications/Ignite.app` after every build

### Never

- Commit secrets or API keys (stored in OS keychain, never in code)
- Use CGO (always use `modernc.org/sqlite`, not `mattn/go-sqlite3`)
- Hardcode API keys in config files
- Add comments without reason — code should be self-documenting
- Use `//go:embed` on non-existent paths (Wails requires `frontend/dist/` to exist for build)
- Push `scripts/local-deploy.sh` — it's gitignored (contains server config)
- Read API keys back from keychain into frontend (keys stay in keychain — show `••••••••`)

---

## Package Audit

Before adding any new dependency:

1. **Check stdlib or existing dep.** Does Go's standard library or an already-imported package provide the functionality?
2. **Check maintenance.** Last commit within 6 months, active issue resolution, no abandoned forks.
3. **Check license.** Must be compatible with AGPL-3.0. MIT, Apache-2.0, BSD, MPL-2.0 are safe. Avoid GPL-incompatible and proprietary dependencies.
4. **Check binary size impact.** Each new package grows the binary. Prefer pure Go, no CGO.
5. **Verify not in banned list.** See below.

### Banned Packages

- `github.com/mattn/go-sqlite3` — requires CGO. Use `modernc.org/sqlite` instead.
- Any package requiring CGO compilation

---

## Project Structure

```
ignite/
├── main.go                         # Wails entry, menu bar, version const
├── app.go                          # Wails bindings (all IPC methods)
├── wails.json                      # Wails v2 config
├── Makefile                        # publish/dev/test/clean targets
├── internal/
│   ├── providers/                  # LLM adapters
│   │   ├── interface.go            # LLMProvider interface
│   │   ├── opencode.go             # OpenCode (Go + Zen)
│   │   ├── deepseek.go             # DeepSeek
│   ├── settings/
│   │   ├── config.go               # config.json load/save
│   │   └── keychain.go             # OS keychain wrapper
│   ├── history/
│   │   ├── models.go               # Project/Message/ProviderModel structs
│   │   └── sqlite.go               # SQLite store + migrations
│   ├── templates/
│   │   ├── data.go                 # ProjectContext struct
│   │   ├── engine.go               # Template engine (sprig + text/template)
│   │   ├── templates.go            # //go:embed loader
│   │   └── templates/*.tmpl        # 4 embedded template files
│   └── scanner/
│       └── scanner.go              # Universal project path analyzer
├── frontend/
│   └── src/
│       ├── App.tsx                 # Main app layout
│       ├── main.tsx                # React entry
│       ├── style.css               # Tailwind v4 + @fontsource imports
│       ├── hooks/
│       │   ├── useTheme.ts         # Theme hook
│       │   └── useConversation.ts  # LLM streaming hook
│       ├── lib/
│       │   ├── store/              # Zustand stores
│       │   │   ├── theme.ts        # Theme state
│       │   │   └── chat.ts         # Chat/messages/projects state
│       │   └── utils.ts            # cn() utility
│       └── components/
│           ├── ui/                 # shadcn/ui components
│           ├── sidebar/            # Project sidebar
│           ├── chat/               # Chat panel + bubbles + input
│           ├── status-bar/         # Provider/model selector + controls
│           └── settings/           # Settings modal + About modal
└── site/                           # Static site (ignite.johansen.foo)
    ├── index.html                  # Landing page
    ├── style.css                   # Site theme
    ├── docs/                       # 7 documentation pages
    └── assets/                     # Downloadable binaries
```

---

## Architecture

```
┌──────────────────────────────────────────────┐
│                Wails Desktop Shell            │
│                                              │
│  ┌────────────────┐   IPC   ┌──────────────┐│
│  │ React Frontend  │◄──────►│  Go Backend   ││
│  │                │         │              ││
│  │ Sidebar        │         │ providers/   ││
│  │ Chat Panel     │         │ settings/    ││
│  │ Status Bar     │         │ history/     ││
│  │ Settings Modal │         │ templates/   ││
│  └────────────────┘         │ scanner/     ││
│                              └──────────────┘│
└──────────────────────────────────────────────┘
```

---

## Database Schema

### projects
| Column     | Type | Description                  |
| :--------- | :--- | :--------------------------- |
| id         | TEXT | UUID primary key             |
| name       | TEXT | Project name                 |
| tagline    | TEXT | Project tagline              |
| path       | TEXT | Output directory             |
| provider   | TEXT | AI provider used             |
| model      | TEXT | AI model used                |
| created_at | TEXT | ISO 8601 timestamp           |
| updated_at | TEXT | ISO 8601 timestamp           |

### conversations
| Column     | Type | Description                  |
| :--------- | :--- | :--------------------------- |
| id         | TEXT | UUID primary key             |
| project_id | TEXT | FK → projects.id             |
| phase      | TEXT | Conversation phase           |
| role       | TEXT | user / assistant / system    |
| content    | TEXT | Message content              |
| created_at | TEXT | ISO 8601 timestamp           |

### provider_models
| Column       | Type | Description                  |
| :----------- | :--- | :--------------------------- |
| provider     | TEXT | Provider name (PK)           |
| model_id     | TEXT | Model ID (PK)                |
| display_name | TEXT | Human-readable name          |
| cached_at    | TEXT | Last cache refresh timestamp |

---

## Environment Variables

No environment variables needed. All configuration is via:

- `~/.ignite/config.json` — provider endpoints, appearance, defaults
- OS Keychain (`com.ignite.app`) — API keys per provider
- `~/.ignite/history.db` — SQLite WAL database (auto-created)

---

## Build Pipeline

1. **Frontend:** Vite builds TypeScript + React → `frontend/dist/`
2. **Embed:** `//go:embed all:frontend/dist` bundles frontend into Go binary
3. **Bindings:** Wails auto-generates TypeScript bindings from Go structs
4. **Compile:** `go build` produces native binary
5. **Package:** Wails creates `.app` bundle + Info.plist
6. **Icon:** `scripts/seticon.sh` generates `iconfile.icns` from `appicon.png`

---

## Testing

| Layer     | Framework | Command                          | Count |
| :-------- | :-------- | :------------------------------- | :---- |
| Go unit   | testing   | `go test ./internal/...`           | 14    |
| Frontend  | vitest    | `cd frontend && pnpm vitest run`   | 6     |
| Typecheck | tsc       | `cd frontend && pnpm typecheck`   | N/A   |

---

## Gotchas

- **Wails bindings regenerate on every build.** After changing Go struct fields, always rebuild to update `frontend/src/lib/wailsjs/`.
- **Wails types use snake_case.** Frontend must use `project_id` not `projectId`, `created_at` not `createdAt`, etc.
- **Embedded fonts.** `@fontsource` imports in `style.css` bundle font files into the Vite build. Don't add Google Fonts `<link>` tags.
- **Keychain can't read back.** `SetAPIKey` writes to keychain. `HasAPIKey` confirms existence. But you cannot read the key back — the input field shows `••••••••` when key exists.
- **Config cleanup.** `ensureProviderConfigs()` runs on every startup, removing stale provider entries and adding missing ones.
- **Model cache refresh.** A background goroutine refreshes model lists every 15 minutes. On startup, the sync runs synchronously before the UI loads.
- **Native `<select>` caveat.** When `value` doesn't match any `<option>`, the browser auto-selects the first option. App.tsx syncs the model from saved config to prevent this.
- **Tailwind v4 CSS-first config.** No `tailwind.config.js` — all customization is in `style.css` via `@theme` and CSS custom properties.
