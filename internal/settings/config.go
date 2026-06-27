package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type ProviderConfig struct {
	Endpoint     string `json:"endpoint"`
	DefaultModel string `json:"default_model"`
}

type Config struct {
	Providers         map[string]ProviderConfig `json:"providers"`
	DefaultProvider   string                    `json:"default_provider"`
	Appearance        string                    `json:"appearance"`
	DefaultLicense    string                    `json:"default_license"`
	DefaultProjectDir string                    `json:"default_project_dir"`
	Font              string                    `json:"font"`
	Name              string                    `json:"name"`
	Avatar            string                    `json:"avatar"`
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".ignite")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

func LoadConfig() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func SaveConfig(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func DefaultConfig() *Config {
	return &Config{
		Providers:         make(map[string]ProviderConfig),
		DefaultProvider:   "",
		Appearance:        "dark",
		DefaultLicense:    "AGPL-3.0",
		DefaultProjectDir: "~/Development",
		Font:              "JetBrains Mono",
	}
}
