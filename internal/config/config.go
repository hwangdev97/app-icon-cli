package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	AppDirName = ".appicon-cli"
	FileName   = "config.json"
)

type Config struct {
	Version         int                       `json:"version"`
	DefaultProvider string                    `json:"default_provider"`
	DefaultTemplate string                    `json:"default_template,omitempty"`
	OutputDirectory string                    `json:"output_directory,omitempty"`
	Providers       map[string]ProviderConfig `json:"providers"`
	TemplatesDir    string                    `json:"templates_dir,omitempty"`
}

type ProviderConfig struct {
	APIKey  string `json:"api_key,omitempty"`
	BaseURL string `json:"base_url,omitempty"`
	Model   string `json:"model,omitempty"`
}

func Default() Config {
	return Config{
		Version:         1,
		DefaultProvider: "nanobanana",
		Providers: map[string]ProviderConfig{
			"nanobanana": {BaseURL: "https://api.nanobanana.example/v1/images", Model: "app-icon"},
			"image2":     {BaseURL: "https://api.image2.example/v1/images", Model: "app-icon"},
		},
	}
}

func HomeDir() (string, error) {
	if v := os.Getenv("APPICON_CLI_HOME"); v != "" {
		return v, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, AppDirName), nil
}

func Path() (string, error) {
	home, err := HomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, FileName), nil
}

func Load() (Config, string, error) {
	path, err := Path()
	if err != nil {
		return Config{}, "", err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			cfg := Default()
			return cfg, path, nil
		}
		return Config{}, path, err
	}
	cfg := Default()
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, path, fmt.Errorf("read config %s: %w", path, err)
	}
	if cfg.Version == 0 {
		cfg.Version = 1
	}
	if cfg.Providers == nil {
		cfg.Providers = map[string]ProviderConfig{}
	}
	return cfg, path, nil
}

func Save(cfg Config) (string, error) {
	path, err := Path()
	if err != nil {
		return "", err
	}
	if cfg.Version == 0 {
		cfg.Version = 1
	}
	if cfg.Providers == nil {
		cfg.Providers = map[string]ProviderConfig{}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return path, err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return path, err
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return path, err
	}
	return path, nil
}

func ResolveAPIKey(value string) (string, error) {
	const prefix = "env:"
	if !strings.HasPrefix(value, prefix) {
		return value, nil
	}
	name := strings.TrimSpace(strings.TrimPrefix(value, prefix))
	if name == "" {
		return "", fmt.Errorf("api_key env reference is empty")
	}
	resolved := os.Getenv(name)
	if resolved == "" {
		return "", fmt.Errorf("environment variable %s is not set", name)
	}
	return resolved, nil
}
