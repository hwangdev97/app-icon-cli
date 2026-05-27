package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hwang/app-icon-cli/internal/provider"
)

func TestSaveImagesCreatesSluggedPNGFiles(t *testing.T) {
	dir := t.TempDir()
	files, err := SaveImages(dir, "Focus Timer!", "ios-clay", []provider.Image{{Data: []byte("png")}})
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Fatalf("files = %v", files)
	}
	if !strings.HasPrefix(filepath.Base(files[0]), "focus-timer-ios-clay-") {
		t.Fatalf("unexpected file name: %s", files[0])
	}
	data, err := os.ReadFile(files[0])
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "png" {
		t.Fatalf("data = %q", data)
	}
}

func TestSaveImagesRejectsEmptyImage(t *testing.T) {
	_, err := SaveImages(t.TempDir(), "App", "tpl", []provider.Image{{}})
	if err == nil {
		t.Fatal("expected error")
	}
}
