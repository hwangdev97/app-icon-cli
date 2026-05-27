package generator

import (
	"strings"
	"testing"

	appicon "github.com/hwang/app-icon-cli/internal/template"
)

func TestBuildPromptIncludesTemplateLayersAndConstraints(t *testing.T) {
	tpl := appicon.Template{
		ID:             "test",
		Name:           "Test Template",
		Description:    "desc",
		BestFor:        []string{"timers"},
		AvoidFor:       []string{"screenshots"},
		Style:          appicon.Style{Aesthetic: "soft 3D", Palette: []string{"red", "white"}},
		Layers:         []appicon.Layer{{Name: "symbol", Role: "primary", Description: "clear glyph", Required: true}},
		PromptHints:    []string{"no text"},
		NegativePrompt: []string{"letters"},
	}
	prompt := BuildPrompt(Request{AppName: "Focus", Idea: "timer", Platform: "macOS", Richness: "rich", Size: 1024, Template: tpl})
	for _, want := range []string{"macOS", "Focus", "timer", "Richness: rich", "tactile depth", "soft 3D", "symbol", "no text", "letters", "timers", "screenshots", "1024x1024"} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("prompt missing %q:\n%s", want, prompt)
		}
	}
}
