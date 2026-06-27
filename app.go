package main

import (
	"context"
	"encoding/json"
	"fmt"
	"ignite/internal/history"
	"ignite/internal/providers"
	"ignite/internal/scanner"
	"ignite/internal/settings"
	"ignite/internal/templates"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx       context.Context
	cfg       *settings.Config
	cfgMu     sync.RWMutex
	store     *history.Store
	providers *providers.Manager
}

func NewApp() *App { return &App{} }

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	cfg, err := settings.LoadConfig()
	if err != nil {
		cfg = settings.DefaultConfig()
	}
	a.cfg = cfg

	if cfg.WindowWidth == 0 {
		cfg.WindowWidth = 1024
	}
	if cfg.WindowHeight == 0 {
		cfg.WindowHeight = 768
	}
	runtime.WindowSetSize(ctx, cfg.WindowWidth, cfg.WindowHeight)

	home, _ := os.UserHomeDir()
	dbPath := filepath.Join(home, ".ignite", "history.db")
	store, err := history.OpenDB(dbPath)
	if err != nil {
		panic(fmt.Sprintf("failed to open history db: %v", err))
	}
	a.store = store

	a.providers = providers.NewManager()

	a.ensureProviderConfigs()
	a.startModelSync()
}

func (a *App) ensureProviderConfigs() {
	a.cfgMu.Lock()
	defer a.cfgMu.Unlock()

	knownProviders := map[string]string{
		"opencode-go":  "https://opencode.ai/zen/go/v1",
		"opencode-zen": "https://opencode.ai/zen/v1",
		"deepseek":     "https://api.deepseek.com/v1",
	}
	changed := false
	for pid := range a.cfg.Providers {
		if _, ok := knownProviders[pid]; !ok {
			delete(a.cfg.Providers, pid)
			changed = true
		}
	}
	for pid, endpoint := range knownProviders {
		if _, ok := a.cfg.Providers[pid]; !ok {
			a.cfg.Providers[pid] = settings.ProviderConfig{Endpoint: endpoint}
			changed = true
		}
	}
	if _, ok := a.cfg.Providers[a.cfg.DefaultProvider]; !ok {
		a.cfg.DefaultProvider = "opencode-go"
		changed = true
	}
	if changed {
		settings.SaveConfig(a.cfg)
	}
}

func (a *App) shutdown(ctx context.Context) {
	if a.store != nil {
		a.store.Close()
	}
}

func (a *App) GetCachedModels(providerName string) ([]history.ProviderModel, error) {
	return a.store.ListCachedModels(providerName)
}

func (a *App) refreshProviderModels() {
	for _, name := range []string{"opencode-go", "opencode-zen", "deepseek"} {
		key, _ := settings.GetAPIKey(name)
		var p providers.LLMProvider
		switch name {
		case "opencode-go":
			p = providers.NewOpenCodeProvider(key, "https://opencode.ai/zen/go/v1")
		case "opencode-zen":
			p = providers.NewOpenCodeProvider(key, "https://opencode.ai/zen/v1")
		case "deepseek":
			p = providers.NewDeepSeekProvider(key)
		}
		models, err := p.ListModels(a.ctx)
		if err != nil {
			continue
		}
		for _, m := range models {
			a.store.UpsertProviderModel(name, m.ID, m.DisplayName)
		}
	}
}

func (a *App) startModelSync() {
	a.refreshProviderModels()
	go func() {
		ticker := time.NewTicker(15 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			a.refreshProviderModels()
		}
	}()
}

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, Ignite is alive!", name)
}

func (a *App) GetVersion() string {
	return version
}

func (a *App) GetSettings() *settings.Config {
	a.cfgMu.RLock()
	defer a.cfgMu.RUnlock()
	return a.cfg
}

func (a *App) SaveSettings(cfg settings.Config) error {
	a.cfgMu.Lock()
	a.cfg = &cfg
	a.cfgMu.Unlock()
	return settings.SaveConfig(&cfg)
}

func (a *App) SetAPIKey(provider string, key string) error {
	return settings.SetAPIKey(provider, key)
}

func (a *App) HasAPIKey(provider string) bool {
	_, err := settings.GetAPIKey(provider)
	if err == nil {
		return true
	}
	if provider == "opencode-go" {
		_, err = settings.GetAPIKey("opencode-zen")
		return err == nil
	}
	if provider == "opencode-zen" {
		_, err = settings.GetAPIKey("opencode-go")
		return err == nil
	}
	return false
}

func (a *App) CreateProject(p history.Project) error { return a.store.CreateProject(p) }
func (a *App) UpdateProject(p history.Project) error { return a.store.UpdateProject(p) }
func (a *App) SetProjectMeta(id, name, tagline string) error {
	return a.store.UpdateProject(history.Project{ID: id, Name: name, Tagline: tagline})
}
func (a *App) ListProjects() ([]history.Project, error) { return a.store.ListProjects() }
func (a *App) GetProject(id string) (*history.Project, error) { return a.store.GetProject(id) }
func (a *App) DeleteProject(id string) error              { return a.store.DeleteProject(id) }
func (a *App) AddMessage(m history.Message) error          { return a.store.AddMessage(m) }
func (a *App) GetMessages(projectID string) ([]history.Message, error) {
	return a.store.GetMessages(projectID)
}

func (a *App) GetProvider(name string) (providers.LLMProvider, error) {
	key, err := settings.GetAPIKey(name)
	if err != nil {
		return nil, fmt.Errorf("no API key for %s: %w", name, err)
	}

	switch name {
	case "opencode-go":
		return providers.NewOpenCodeProvider(key, "https://opencode.ai/zen/go/v1"), nil
	case "opencode-zen":
		return providers.NewOpenCodeProvider(key, "https://opencode.ai/zen/v1"), nil
	case "deepseek":
		return providers.NewDeepSeekProvider(key), nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
}

func (a *App) SendMessage(providerName, model string, messages []providers.Message) (*providers.ChatResponse, error) {
	p, err := a.GetProvider(providerName)
	if err != nil {
		return nil, err
	}
	return p.Chat(a.ctx, model, messages)
}

func (a *App) SendMessageStream(providerName, model string, messages []providers.Message) error {
	p, err := a.GetProvider(providerName)
	if err != nil {
		return err
	}

	return p.ChatStream(a.ctx, model, messages, func(chunk string) error {
		runtime.EventsEmit(a.ctx, "stream-chunk", chunk)
		return nil
	})
}

func (a *App) ValidateProviderKey(providerName, key string) error {
	var p providers.LLMProvider
	switch providerName {
	case "opencode-go":
		p = providers.NewOpenCodeProvider(key, "https://opencode.ai/zen/go/v1")
	case "opencode-zen":
		p = providers.NewOpenCodeProvider(key, "https://opencode.ai/zen/v1")
	case "deepseek":
		p = providers.NewDeepSeekProvider(key)
	default:
		return fmt.Errorf("unknown provider: %s", providerName)
	}
	return p.ValidateKey(a.ctx)
}

func (a *App) SaveProjectFiles(projectDir string, files *templates.ProjectFiles) error {
	dir := filepath.Join(a.cfg.DefaultProjectDir, projectDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	writeFile := func(name, content string) error {
		return os.WriteFile(filepath.Join(dir, name), []byte(content), 0644)
	}

	if err := writeFile(projectDir+".md", files.ProjectMD); err != nil {
		return err
	}
	if err := writeFile("agents.md", files.AgentsMD); err != nil {
		return err
	}
	if err := writeFile("plan.md", files.PlanMD); err != nil {
		return err
	}
	if err := writeFile("README.md", files.ReadmeMD); err != nil {
		return err
	}

	return nil
}

func (a *App) ListProviderModels(providerName string) ([]providers.Model, error) {
	var p providers.LLMProvider
	key, err := settings.GetAPIKey(providerName)
	switch providerName {
	case "deepseek":
		p = providers.NewDeepSeekProvider(key)
	default:
		if err != nil {
			return nil, fmt.Errorf("no API key for %s", providerName)
		}
		switch providerName {
		case "opencode-go":
			p = providers.NewOpenCodeProvider(key, "https://opencode.ai/zen/go/v1")
		case "opencode-zen":
			p = providers.NewOpenCodeProvider(key, "https://opencode.ai/zen/v1")
		default:
			return nil, fmt.Errorf("unknown provider: %s", providerName)
		}
	}
	return p.ListModels(a.ctx)
}

func (a *App) ExportChat(messages []history.Message) string {
	var md string
	for _, m := range messages {
		role := string(m.Role)
		md += fmt.Sprintf("### %s\n\n%s\n\n---\n\n", role, m.Content)
	}
	return md
}

func (a *App) AnalyzePath(path string) string {
	return scanner.AnalyzePath(path)
}

func (a *App) AnalyzePathContent(path string) string {
	return scanner.AnalyzePathContent(path)
}

func (a *App) FetchURL(url string) string {
	if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
		return ""
	}
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 50000))
	if err != nil {
		return ""
	}
	text := string(body)
	text = stripHTML(text)
	if len(text) > 5000 {
		text = text[:5000] + "..."
	}
	return text
}

func stripHTML(s string) string {
	var b strings.Builder
	inTag := false
	for _, r := range s {
		if r == '<' {
			inTag = true
		} else if r == '>' {
			inTag = false
			b.WriteRune(' ')
		} else if !inTag {
			b.WriteRune(r)
		}
	}
	return strings.Join(strings.Fields(b.String()), " ")
}

func (a *App) GenerateProjectFiles(providerName, model, projectName string, messages []providers.Message) (*templates.ProjectFiles, error) {
	p, err := a.GetProvider(providerName)
	if err != nil {
		return nil, err
	}

	extractPrompt := `Extract the project context from the conversation above as JSON. Return ONLY valid JSON, no markdown:

{
  "name": "project-name",
  "tagline": "short tagline",
  "description": "project description",
  "license": "AGPL-3.0",
  "features": ["feature1", "feature2"],
  "phases": [{"name": "Phase 0", "description": "...", "tasks": ["task1"]}],
  "techStack": [{"category": "Frontend", "choice": "React", "version": "19"}],
  "dependencies": [{"name": "react", "version": "^19.0.0", "why": "UI framework"}],
  "apis": [{"method": "GET", "path": "/api/users", "desc": "List users", "auth": "Bearer"}],
  "dbTables": [{"name": "users", "columns": [{"name": "id", "type": "TEXT", "desc": "UUID"}]}],
  "performance": [{"metric": "App launch", "target": "<1s"}],
  "risks": [{"risk": "Data loss", "mitigation": "Backups"}],
  "bannedPackages": [],
  "envVars": [{"name": "DATABASE_URL", "desc": "Connection string", "default": ""}],
  "devWorkflow": {"setup": [], "dev": [], "build": [], "test": [], "lint": [], "typeCheck": []}
}`

	extractMsg := []providers.Message{
		{Role: "system", Content: extractPrompt},
		{Role: "user", Content: "Extract the project context from these messages."},
	}

	resp, err := p.Chat(a.ctx, model, append(extractMsg, messages...))
	if err != nil {
		return nil, fmt.Errorf("extract: %w", err)
	}

	var ctx templates.ProjectContext
	if err := json.Unmarshal([]byte(resp.Content), &ctx); err != nil {
		return nil, fmt.Errorf("parse project context: %w", err)
	}

	engine, err := templates.EmbeddedEngine()
	if err != nil {
		return nil, fmt.Errorf("template engine: %w", err)
	}

	files, err := engine.Generate(ctx)
	if err != nil {
		return nil, fmt.Errorf("generate files: %w", err)
	}

	dirPath := a.cfg.DefaultProjectDir
	if strings.HasPrefix(dirPath, "~/") {
		home, _ := os.UserHomeDir()
		dirPath = filepath.Join(home, dirPath[2:])
	}
	projectDir := filepath.Join(dirPath, projectName)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return nil, fmt.Errorf("create dir: %w", err)
	}

	os.WriteFile(filepath.Join(projectDir, projectName+".md"), []byte(files.ProjectMD), 0644)
	os.WriteFile(filepath.Join(projectDir, "agents.md"), []byte(files.AgentsMD), 0644)
	os.WriteFile(filepath.Join(projectDir, "plan.md"), []byte(files.PlanMD), 0644)
	os.WriteFile(filepath.Join(projectDir, "README.md"), []byte(files.ReadmeMD), 0644)

	return files, nil
}

func (a *App) ResizeWindow(width, height int) {
	runtime.WindowSetSize(a.ctx, width, height)
}

func (a *App) SelectDirectory() string {
	dirPath := a.cfg.DefaultProjectDir
	if strings.HasPrefix(dirPath, "~/") {
		home, _ := os.UserHomeDir()
		dirPath = filepath.Join(home, dirPath[2:])
	}
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            "Select project folder",
		DefaultDirectory: dirPath,
	})
	if err != nil || dir == "" {
		return a.cfg.DefaultProjectDir
	}
	a.cfg.DefaultProjectDir = dir
	settings.SaveConfig(a.cfg)
	return dir
}
