package main

import (
	"context"
	"fmt"
	"ignite/internal/history"
	"ignite/internal/providers"
	"ignite/internal/scanner"
	"ignite/internal/settings"
	"ignite/internal/templates"
	"os"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx       context.Context
	cfg       *settings.Config
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
	knownProviders := map[string]string{
		"opencode-go":  "https://opencode.ai/zen/go/v1",
		"opencode-zen": "https://opencode.ai/zen/v1",
		"claude":       "https://api.anthropic.com/v1/messages",
		"deepseek":     "https://api.deepseek.com/v1",
	}
	for pid := range a.cfg.Providers {
		if _, ok := knownProviders[pid]; !ok {
			delete(a.cfg.Providers, pid)
		}
	}
	for pid, endpoint := range knownProviders {
		if _, ok := a.cfg.Providers[pid]; !ok {
			a.cfg.Providers[pid] = settings.ProviderConfig{Endpoint: endpoint}
		}
	}
	settings.SaveConfig(a.cfg)
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
	for _, name := range []string{"opencode-go", "opencode-zen", "claude", "deepseek"} {
		key, _ := settings.GetAPIKey(name)
		var p providers.LLMProvider
		switch name {
		case "opencode-go":
			p = providers.NewOpenCodeProvider(key, "https://opencode.ai/zen/go/v1")
		case "opencode-zen":
			p = providers.NewOpenCodeProvider(key, "https://opencode.ai/zen/v1")
		case "claude":
			p = providers.NewClaudeProvider(key)
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

func (a *App) GetSettings() *settings.Config { return a.cfg }

func (a *App) SaveSettings(cfg settings.Config) error {
	a.cfg = &cfg
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
	case "claude":
		return providers.NewClaudeProvider(key), nil
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
	case "claude":
		p = providers.NewClaudeProvider(key)
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
	case "claude":
		p = providers.NewClaudeProvider(key)
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

func (a *App) SelectDirectory() string {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:       "Select project folder",
		DefaultDirectory: a.cfg.DefaultProjectDir,
	})
	if err != nil || dir == "" {
		return a.cfg.DefaultProjectDir
	}
	a.cfg.DefaultProjectDir = dir
	settings.SaveConfig(a.cfg)
	return dir
}
