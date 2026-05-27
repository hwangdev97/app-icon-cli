package template

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAllAllowsCustomTemplateOverride(t *testing.T) {
	dir := t.TempDir()
	custom := `{
  "id": "ios-clay-symbol",
  "name": "Custom Clay",
  "description": "custom",
  "style": {"aesthetic": "custom style"},
  "layers": [{"name": "base", "role": "base", "description": "base", "required": true}]
}`
	if err := os.WriteFile(filepath.Join(dir, "custom.json"), []byte(custom), 0o644); err != nil {
		t.Fatal(err)
	}
	templates, err := LoadAll(dir)
	if err != nil {
		t.Fatal(err)
	}
	tpl, ok := Find(templates, "ios-clay-symbol")
	if !ok {
		t.Fatal("missing template")
	}
	if tpl.Name != "Custom Clay" {
		t.Fatalf("template name = %s", tpl.Name)
	}
}

func TestValidateRequiresRequiredLayer(t *testing.T) {
	tpl := Template{ID: "x", Name: "X", Style: Style{Aesthetic: "flat"}}
	if err := tpl.Validate(); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestBuiltInsIncludeShortlistedTemplates(t *testing.T) {
	templates := BuiltIns()
	wantIDs := []string{
		"ios-clay-symbol",
		"ios-glass-layer",
		"ios-flat-gradient",
		"ios-monoline-badge",
		"ios-neon-depth",
		"ios-paper-cut",
		"nelson-precision",
		"pixelresort-object",
		"classic-mac-skeuo",
		"glossy-toy-3d",
		"lance-brand-mark",
	}
	for _, id := range wantIDs {
		tpl, ok := Find(templates, id)
		if !ok {
			t.Fatalf("missing built-in template %s", id)
		}
		if err := tpl.Validate(); err != nil {
			t.Fatalf("template %s is invalid: %v", id, err)
		}
		if len(tpl.BestFor) == 0 || len(tpl.NegativePrompt) == 0 || tpl.Output.Size != 1024 {
			t.Fatalf("template %s is missing metadata: %+v", id, tpl)
		}
	}
}
