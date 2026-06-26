package main

import (
	"context"
	"fmt"
	"ignite/internal/history"
	"ignite/internal/settings"
	"os"
	"path/filepath"
)

type App struct {
	ctx   context.Context
	cfg   *settings.Config
	store *history.Store
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
func (a *App) AddMessage(m history.Message) error { return a.store.AddMessage(m) }
func (a *App) GetMessages(projectID string) ([]history.Message, error) { return a.store.GetMessages(projectID) }
