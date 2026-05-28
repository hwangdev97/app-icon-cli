package preview

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"
)

//go:embed template.html
var templates embed.FS

type Icon struct {
	File string
	Name string
}

type Data struct {
	AppName     string
	Platform    string
	TemplateID  string
	GeneratedAt string
	Icons       []Icon
}

func Write(outPath string, data Data) (string, error) {
	if outPath == "" {
		outPath = "preview.html"
	}
	if data.GeneratedAt == "" {
		data.GeneratedAt = time.Now().Format(time.RFC3339)
	}
	baseDir := filepath.Dir(outPath)
	icons := make([]Icon, 0, len(data.Icons))
	for _, icon := range data.Icons {
		file := icon.File
		if file != "" {
			if rel, err := filepath.Rel(baseDir, file); err == nil {
				file = filepath.ToSlash(rel)
			} else {
				file = filepath.ToSlash(file)
			}
		}
		name := icon.Name
		if name == "" {
			name = filepath.Base(icon.File)
		}
		icons = append(icons, Icon{File: file, Name: name})
	}
	data.Icons = icons
	tpl, err := template.ParseFS(templates, "template.html")
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(outPath, buf.Bytes(), 0o644); err != nil {
		return "", err
	}
	return outPath, nil
}

func IconsFromFiles(files []string) []Icon {
	icons := make([]Icon, 0, len(files))
	for _, file := range files {
		icons = append(icons, Icon{File: file, Name: filepath.Base(file)})
	}
	return icons
}

func DefaultPath(outDir string) string {
	if outDir == "" {
		outDir = "appicons"
	}
	return filepath.Join(outDir, "preview.html")
}

func Validate(data Data) error {
	if len(data.Icons) == 0 {
		return fmt.Errorf("preview requires at least one icon")
	}
	return nil
}
