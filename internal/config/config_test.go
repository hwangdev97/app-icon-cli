package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaultAndSaveUsesAppHome(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("APPICON_CLI_HOME", dir)

	cfg, path, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if path != filepath.Join(dir, FileName) {
		t.Fatalf("path = %s", path)
	}
	if cfg.DefaultProvider != "nanobanana" {
		t.Fatalf("default provider = %s", cfg.DefaultProvider)
	}

	cfg.Providers["nanobanana"] = ProviderConfig{APIKey: "secret", BaseURL: "http://example.test", Model: "m"}
	saved, err := Save(cfg)
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(saved)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("config mode = %v", info.Mode().Perm())
	}
}
