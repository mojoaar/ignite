package main

import (
	"context"
	"fmt"
	"ignite/internal/settings"
)

type App struct {
	ctx context.Context
	cfg *settings.Config
}

func NewApp() *App { return &App{} }

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	cfg, err := settings.LoadConfig()
	if err != nil {
		cfg = settings.DefaultConfig()
	}
	a.cfg = cfg
}

func (a *App) shutdown(ctx context.Context) {}

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
