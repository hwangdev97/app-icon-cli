package cli

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hwang/app-icon-cli/internal/config"
	"github.com/hwang/app-icon-cli/internal/provider"
)

type fakeProvider struct{}

func (fakeProvider) Generate(context.Context, provider.GenerateRequest) ([]provider.Image, error) {
	data, _ := base64.StdEncoding.DecodeString("cG5n")
	return []provider.Image{{Data: data}}, nil
}

func TestInitWritesConfigNonInteractive(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("APPICON_CLI_HOME", dir)
	var out, errOut bytes.Buffer

	code := Run([]string{"init", "--provider", "image2", "--api-key", "k", "--base-url", "http://example.test", "--model", "m"}, strings.NewReader(""), &out, &errOut)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, errOut.String())
	}
	cfg, _, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.DefaultProvider != "image2" || cfg.Providers["image2"].APIKey != "k" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestGenerateDryRunPrintsProviderRequest(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("APPICON_CLI_HOME", dir)
	cfg := config.Default()
	cfg.Providers["nanobanana"] = config.ProviderConfig{APIKey: "k", BaseURL: "http://example.test", Model: "m"}
	if _, err := config.Save(cfg); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	code := Run([]string{"generate", "--app", "Focus", "--idea", "timer", "--richness", "minimal", "--dry-run"}, strings.NewReader(""), &out, &errOut)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, errOut.String())
	}
	var req provider.GenerateRequest
	if err := json.Unmarshal(out.Bytes(), &req); err != nil {
		t.Fatal(err)
	}
	if req.Provider != "nanobanana" || !strings.Contains(req.Prompt, "Focus") || !strings.Contains(req.Prompt, "Richness: minimal") {
		t.Fatalf("unexpected request: %+v", req)
	}
}

func TestGenerateSavesImageWithInjectedProvider(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("APPICON_CLI_HOME", dir)
	cfg := config.Default()
	cfg.Providers["nanobanana"] = config.ProviderConfig{APIKey: "k", BaseURL: "http://example.test", Model: "m"}
	if _, err := config.Save(cfg); err != nil {
		t.Fatal(err)
	}
	outDir := filepath.Join(dir, "out")
	var out, errOut bytes.Buffer
	code := generateCmd([]string{"--app", "Focus", "--out", outDir, "--json"}, &out, &errOut, fakeProvider{})
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, errOut.String())
	}
	files, err := filepath.Glob(filepath.Join(outDir, "*.png"))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Fatalf("files = %v", files)
	}
	data, err := os.ReadFile(files[0])
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "png" {
		t.Fatalf("data = %q", string(data))
	}
}

func TestConfigShowMasksAPIKey(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("APPICON_CLI_HOME", dir)
	cfg := config.Default()
	cfg.Providers["nanobanana"] = config.ProviderConfig{APIKey: "secret", BaseURL: "http://example.test", Model: "m"}
	if _, err := config.Save(cfg); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	code := Run([]string{"config", "show"}, strings.NewReader(""), &out, &errOut)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, errOut.String())
	}
	if strings.Contains(out.String(), "secret") {
		t.Fatalf("config leaked API key: %s", out.String())
	}
	if !strings.Contains(out.String(), "configured") {
		t.Fatalf("config did not show key status: %s", out.String())
	}
}

func TestTemplatesShowBuiltIn(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("APPICON_CLI_HOME", dir)
	var out, errOut bytes.Buffer
	code := Run([]string{"templates", "show", "ios-glass-layer"}, strings.NewReader(""), &out, &errOut)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, errOut.String())
	}
	if !strings.Contains(out.String(), "iOS Glass Layer") {
		t.Fatalf("unexpected output: %s", out.String())
	}
}

func TestGenerateRejectsInvalidCount(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("APPICON_CLI_HOME", dir)
	var out, errOut bytes.Buffer
	code := Run([]string{"generate", "--app", "Focus", "--count", "0"}, strings.NewReader(""), &out, &errOut)
	if code != 2 {
		t.Fatalf("code=%d stderr=%s", code, errOut.String())
	}
}

func TestGenerateRejectsInvalidRichness(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("APPICON_CLI_HOME", dir)
	var out, errOut bytes.Buffer
	code := Run([]string{"generate", "--app", "Focus", "--richness", "busy"}, strings.NewReader(""), &out, &errOut)
	if code != 2 {
		t.Fatalf("code=%d stderr=%s", code, errOut.String())
	}
}
