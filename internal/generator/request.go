package generator

import (
	"fmt"
	"strings"

	appicon "github.com/hwang/app-icon-cli/internal/template"
)

type Request struct {
	AppName     string           `json:"app_name"`
	Idea        string           `json:"idea"`
	Platform    string           `json:"platform,omitempty"`
	Richness    string           `json:"richness,omitempty"`
	Count       int              `json:"count"`
	Size        int              `json:"size"`
	Template    appicon.Template `json:"template"`
	ExtraParams map[string]any   `json:"extra_params,omitempty"`
}

func BuildPrompt(req Request) string {
	var b strings.Builder
	platform := req.Platform
	if platform == "" {
		platform = "iPhone"
	}
	fmt.Fprintf(&b, "Create an app-store-ready %s app icon for %q.\n", platform, req.AppName)
	if req.Idea != "" {
		fmt.Fprintf(&b, "App concept: %s.\n", req.Idea)
	}
	fmt.Fprintf(&b, "%s\n", RichnessGuidance(req.Richness))
	fmt.Fprintf(&b, "Template: %s. %s\n", req.Template.Name, req.Template.Description)
	fmt.Fprintf(&b, "Style: %s.", req.Template.Style.Aesthetic)
	if len(req.Template.Style.Palette) > 0 {
		fmt.Fprintf(&b, " Palette: %s.", strings.Join(req.Template.Style.Palette, ", "))
	}
	if req.Template.Style.Material != "" {
		fmt.Fprintf(&b, " Material: %s.", req.Template.Style.Material)
	}
	if req.Template.Style.Lighting != "" {
		fmt.Fprintf(&b, " Lighting: %s.", req.Template.Style.Lighting)
	}
	if len(req.Template.BestFor) > 0 {
		fmt.Fprintf(&b, "\nBest use cases for this template: %s.", strings.Join(req.Template.BestFor, ", "))
	}
	if len(req.Template.AvoidFor) > 0 {
		fmt.Fprintf(&b, "\nAvoid drifting toward: %s.", strings.Join(req.Template.AvoidFor, ", "))
	}
	b.WriteString("\nLayers:\n")
	for _, layer := range req.Template.Layers {
		required := "optional"
		if layer.Required {
			required = "required"
		}
		fmt.Fprintf(&b, "- %s (%s, %s): %s\n", layer.Name, layer.Role, required, layer.Description)
	}
	if len(req.Template.PromptHints) > 0 {
		fmt.Fprintf(&b, "Constraints: %s.\n", strings.Join(req.Template.PromptHints, ", "))
	}
	if len(req.Template.NegativePrompt) > 0 {
		fmt.Fprintf(&b, "Avoid: %s.\n", strings.Join(req.Template.NegativePrompt, ", "))
	}
	format := req.Template.Output.Format
	if format == "" {
		format = "PNG"
	}
	alpha := "no transparency unless the template asks for it"
	if req.Template.Output.Alpha != nil {
		if *req.Template.Output.Alpha {
			alpha = "transparent areas are allowed only where useful"
		} else {
			alpha = "no alpha channel"
		}
	}
	fmt.Fprintf(&b, "Render as a square %dx%d %s icon, centered composition, %s.", req.Size, req.Size, strings.ToUpper(format), alpha)
	return b.String()
}
