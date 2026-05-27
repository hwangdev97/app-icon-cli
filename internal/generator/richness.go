package generator

import (
	"fmt"
	"strings"
)

const DefaultRichness = "balanced"

var richnessGuidance = map[string]string{
	"minimal":  "Richness: minimal. Use one unmistakable subject, one to two colors, flat or nearly flat treatment, no optional detail layers, no decoration, maximum small-size clarity.",
	"simple":   "Richness: simple. Use one primary subject with very restrained depth, two to three colors, at most one supporting detail, and a clean uncluttered silhouette.",
	"balanced": "Richness: balanced. Use one primary subject with controlled material, lighting, and depth; include only details that improve recognition and keep the icon readable at small sizes.",
	"rich":     "Richness: rich. Allow layered material, stronger lighting, tactile depth, and a few carefully chosen recognition details while preserving one clear focal subject.",
	"maximal":  "Richness: maximal. Allow the template's full visual craft: richer materials, stronger dimensionality, refined micro-details, and dramatic finish, but still avoid multiple competing subjects or clutter.",
}

func NormalizeRichness(value string) (string, error) {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return DefaultRichness, nil
	}
	if _, ok := richnessGuidance[value]; !ok {
		return "", fmt.Errorf("richness must be one of: minimal, simple, balanced, rich, maximal")
	}
	return value, nil
}

func RichnessGuidance(value string) string {
	normalized, err := NormalizeRichness(value)
	if err != nil {
		normalized = DefaultRichness
	}
	return richnessGuidance[normalized]
}
