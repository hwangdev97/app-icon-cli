package preview

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWritePreviewUsesRelativeIconPaths(t *testing.T) {
	dir := t.TempDir()
	icon := filepath.Join(dir, "icons", "icon.png")
	if err := os.MkdirAll(filepath.Dir(icon), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(icon, []byte("png"), 0o644); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(dir, "preview.html")
	written, err := Write(out, Data{
		AppName:    "Focus",
		Platform:   "macOS",
		TemplateID: "test",
		Icons:      IconsFromFiles([]string{icon}),
	})
	if err != nil {
		t.Fatal(err)
	}
	if written != out {
		t.Fatalf("written = %s", written)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	html := string(data)
	for _, want := range []string{"Focus App Icon Preview", "icons/icon.png", "Source images are shown unmodified"} {
		if !strings.Contains(html, want) {
			t.Fatalf("preview missing %q:\n%s", want, html)
		}
	}
	if strings.Contains(html, ".icon-source {\n      border-radius") || strings.Contains(html, "clip-path") {
		t.Fatalf("preview should not mask icon source:\n%s", html)
	}
}

func TestValidateRequiresIcons(t *testing.T) {
	if err := Validate(Data{}); err == nil {
		t.Fatal("expected validation error")
	}
}
