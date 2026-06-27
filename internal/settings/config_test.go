package settings

import (
	"os"
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
	cfg.DefaultProvider = "deepseek"
	cfg.Providers["deepseek"] = ProviderConfig{
		Endpoint: "https://api.deepseek.com/v1", DefaultModel: "deepseek-v4-pro",
	}
	if err := SaveConfig(cfg); err != nil {
		t.Fatalf("SaveConfig: %v", err)
	}
	loaded, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if loaded.DefaultProvider != "deepseek" {
		t.Errorf("provider mismatch")
	}
	if loaded.WindowWidth != 1024 {
		t.Errorf("WindowWidth should be 1024, got %d", loaded.WindowWidth)
	}
}

func TestDefaultConfigAllFields(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Appearance != "dark" {
		t.Error("expected dark appearance")
	}
	if cfg.DefaultLicense != "AGPL-3.0" {
		t.Error("expected AGPL-3.0")
	}
	if cfg.Font != "JetBrains Mono" {
		t.Error("expected JetBrains Mono")
	}
	if cfg.WindowWidth != 1024 {
		t.Errorf("expected WindowWidth 1024, got %d", cfg.WindowWidth)
	}
	if cfg.WindowHeight != 768 {
		t.Errorf("expected WindowHeight 768, got %d", cfg.WindowHeight)
	}
	if cfg.Name != "" {
		t.Error("expected empty Name")
	}
	if cfg.Avatar != "" {
		t.Error("expected empty Avatar")
	}
	if cfg.DefaultProjectDir == "" {
		t.Error("expected non-empty DefaultProjectDir")
	}
	if cfg.Providers == nil {
		t.Error("expected non-nil Providers map")
	}
}

func TestLoadConfigCorruptFile(t *testing.T) {
	dir := t.TempDir()
	SetConfigDir(dir)
	defer func() { SetConfigDir("") }()
	cfg := DefaultConfig()
	SaveConfig(cfg)
	if err := os.WriteFile(dir+"/config.json", []byte("this is not json"), 0600); err != nil {
		t.Fatalf("write corrupt: %v", err)
	}
	_, err := LoadConfig()
	if err == nil {
		t.Error("expected error for corrupt config")
	}
}

func TestLoadConfigMissingFile(t *testing.T) {
	dir := t.TempDir()
	os.RemoveAll(dir)
	SetConfigDir(dir)
	defer func() { SetConfigDir("") }()
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig should return DefaultConfig when file missing: %v", err)
	}
	if cfg.Appearance != "dark" {
		t.Error("expected dark appearance from DefaultConfig")
	}
}
