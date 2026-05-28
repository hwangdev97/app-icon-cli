package generator

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hwang/app-icon-cli/internal/provider"
)

type Result struct {
	Files   []string `json:"files"`
	Prompt  string   `json:"prompt"`
	Preview string   `json:"preview,omitempty"`
}

func SaveImages(outDir, appName, templateID string, images []provider.Image) ([]string, error) {
	if outDir == "" {
		outDir = "appicons"
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return nil, err
	}
	slug := slugify(appName)
	if slug == "" {
		slug = "app"
	}
	stamp := time.Now().Format("20060102-150405")
	files := make([]string, 0, len(images))
	for i, img := range images {
		data := img.Data
		if len(data) == 0 && img.URL != "" {
			fetched, err := fetch(img.URL)
			if err != nil {
				return nil, err
			}
			data = fetched
		}
		if len(data) == 0 {
			return nil, fmt.Errorf("image %d has no data or url", i+1)
		}
		name := fmt.Sprintf("%s-%s-%s-%02d.png", slug, templateID, stamp, i+1)
		path := filepath.Join(outDir, name)
		if err := os.WriteFile(path, data, 0o644); err != nil {
			return nil, err
		}
		files = append(files, path)
	}
	return files, nil
}

func SaveManifest(path string, result Result) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

func fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("download %s returned %s", url, resp.Status)
	}
	return io.ReadAll(resp.Body)
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	lastDash := false
	for _, r := range s {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			b.WriteByte('-')
			lastDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}
