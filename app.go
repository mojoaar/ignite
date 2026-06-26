package main

import (
	"context"
	"fmt"
	"ignite/internal/history"
	"ignite/internal/providers"
	"ignite/internal/settings"
	"ignite/internal/templates"
	"os"
	"path/filepath"

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
}

func (a *App) shutdown(ctx context.Context) {
	if a.store != nil {
		a.store.Close()
	}
}

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, Ignite is alive!", name)
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
	return err == nil
}

func (a *App) CreateProject(p history.Project) error { return a.store.CreateProject(p) }
func (a *App) UpdateProject(p history.Project) error { return a.store.UpdateProject(p) }
func (a *App) ListProjects() ([]history.Project, error) { return a.store.ListProjects() }
func (a *App) GetProject(id string) (*history.Project, error) { return a.store.GetProject(id) }
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
	case "opencode-go", "opencode-zen":
		return providers.NewOpenCodeProvider(key), nil
	case "claude":
		return providers.NewClaudeProvider(key), nil
	case "deepseek":
		return providers.NewDeepSeekProvider(key), nil
	case "github-copilot":
		return providers.NewGitHubCopilotProvider(key), nil
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
	case "opencode-go", "opencode-zen":
		p = providers.NewOpenCodeProvider(key)
	case "claude":
		p = providers.NewClaudeProvider(key)
	case "deepseek":
		p = providers.NewDeepSeekProvider(key)
	case "github-copilot":
		p = providers.NewGitHubCopilotProvider(key)
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

func (a *App) ExportChat(messages []history.Message) string {
	var md string
	for _, m := range messages {
		role := string(m.Role)
		md += fmt.Sprintf("### %s\n\n%s\n\n---\n\n", role, m.Content)
	}
	return md
}
