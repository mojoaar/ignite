package settings

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Appearance != "dark" {
		t.Errorf("expected dark, got %s", cfg.Appearance)
	}
	if cfg.DefaultLicense != "AGPL-3.0" {
		t.Errorf("expected AGPL-3.0")
	}
	if cfg.Font != "JetBrains Mono" {
		t.Errorf("expected JetBrains Mono")
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	cfg := DefaultConfig()
	cfg.DefaultProvider = "claude"
	cfg.Providers["claude"] = ProviderConfig{
		Endpoint: "https://api.anthropic.com/v1/messages", DefaultModel: "claude-opus-4.5",
	}
	if err := SaveConfig(cfg); err != nil {
		t.Fatalf("SaveConfig: %v", err)
	}
	loaded, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if loaded.DefaultProvider != "claude" {
		t.Errorf("provider mismatch")
	}
}
