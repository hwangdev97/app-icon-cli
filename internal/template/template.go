package template

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Template struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	Categories     []string       `json:"categories,omitempty"`
	BestFor        []string       `json:"best_for,omitempty"`
	AvoidFor       []string       `json:"avoid_for,omitempty"`
	Style          Style          `json:"style"`
	Layers         []Layer        `json:"layers"`
	Parameters     map[string]any `json:"parameters,omitempty"`
	PromptHints    []string       `json:"prompt_hints,omitempty"`
	NegativePrompt []string       `json:"negative_prompt,omitempty"`
	Output         Output         `json:"output,omitempty"`
}

type Style struct {
	Aesthetic string   `json:"aesthetic"`
	Palette   []string `json:"palette,omitempty"`
	Lighting  string   `json:"lighting,omitempty"`
	Material  string   `json:"material,omitempty"`
}

type Layer struct {
	Name        string         `json:"name"`
	Role        string         `json:"role"`
	Description string         `json:"description"`
	Required    bool           `json:"required"`
	Parameters  map[string]any `json:"parameters,omitempty"`
}

type Output struct {
	Size   int    `json:"size,omitempty"`
	Format string `json:"format,omitempty"`
	Alpha  *bool  `json:"alpha,omitempty"`
}

func BuiltIns() []Template {
	noAlpha := false
	return []Template{
		{
			ID:          "ios-clay-symbol",
			Name:        "iOS Clay Symbol",
			Description: "Soft dimensional iPhone icon with a tactile clay-like symbol and clean background.",
			Categories:  []string{"productivity", "wellness", "education", "finance", "developer-tools"},
			BestFor:     []string{"productivity apps", "wellness apps", "education apps", "finance tools", "developer tools"},
			AvoidFor:    []string{"apps that need hard-edged industrial visuals", "photo-heavy brands"},
			Style: Style{
				Aesthetic: "modern iOS app icon, soft 3D clay, minimal, high legibility at small sizes",
				Palette:   []string{"adaptive brand color", "warm white", "deep accent"},
				Lighting:  "large softbox, subtle ambient occlusion",
				Material:  "matte clay",
			},
			Layers: []Layer{
				{Name: "background", Role: "base", Description: "rounded-square safe-area background with gentle depth", Required: true},
				{Name: "symbol", Role: "primary", Description: "single centered metaphor for the app purpose", Required: true},
				{Name: "highlight", Role: "detail", Description: "small light-catching accent without text", Required: false},
			},
			PromptHints:    []string{"matte clay", "soft depth", "centered object", "high small-size legibility", "no words"},
			NegativePrompt: []string{"letters", "tiny details", "screenshots", "Apple hardware replicas", "brand logos"},
			Output:         Output{Size: 1024, Format: "png", Alpha: &noAlpha},
		},
		{
			ID:          "ios-glass-layer",
			Name:        "iOS Glass Layer",
			Description: "Layered translucent icon with glass panels, crisp foreground symbol, and high-contrast depth.",
			Categories:  []string{"ai", "media", "premium-utility", "dashboard"},
			BestFor:     []string{"AI tools", "media apps", "premium utilities", "analytics dashboards"},
			AvoidFor:    []string{"icons that must stay flat", "low-contrast brand palettes"},
			Style: Style{
				Aesthetic: "premium iOS app icon, translucent layered glass, crisp geometric foreground",
				Palette:   []string{"brand gradient", "white highlights", "dark contrast"},
				Lighting:  "specular highlights, soft refraction",
				Material:  "frosted glass",
			},
			Layers: []Layer{
				{Name: "gradient", Role: "base", Description: "vibrant background gradient inside iOS icon safe area", Required: true},
				{Name: "glass-panel", Role: "middle", Description: "one or two translucent panels framing the symbol", Required: true},
				{Name: "glyph", Role: "primary", Description: "bold simple glyph readable at 60 px", Required: true},
				{Name: "specular-highlight", Role: "detail", Description: "controlled glass highlight that does not obscure the glyph", Required: false},
			},
			PromptHints:    []string{"layered glass", "clear foreground silhouette", "high contrast", "avoid blurry edges", "app-store-ready"},
			NegativePrompt: []string{"muddy transparency", "low contrast", "blurred glyph", "busy reflections", "text", "screenshots"},
			Output:         Output{Size: 1024, Format: "png", Alpha: &noAlpha},
		},
		{
			ID:          "ios-flat-gradient",
			Name:        "iOS Flat Gradient",
			Description: "Modern flat gradient icon with one bold geometric symbol and strong small-size readability.",
			Categories:  []string{"mvp", "saas", "utility", "internal-tool"},
			BestFor:     []string{"MVPs", "SaaS products", "internal tools", "simple utilities"},
			AvoidFor:    []string{"apps that need realism", "highly playful game brands"},
			Style: Style{
				Aesthetic: "flat modern iOS app icon, vector-like, bold geometric symbol, clean gradient field",
				Palette:   []string{"two brand colors", "neutral highlight"},
				Lighting:  "none, flat graphic",
				Material:  "clean vector",
			},
			Layers: []Layer{
				{Name: "gradient-field", Role: "base", Description: "smooth two-color or three-color gradient background", Required: true},
				{Name: "geometric-glyph", Role: "primary", Description: "single bold geometric glyph with clear silhouette", Required: true},
				{Name: "inner-depth", Role: "detail", Description: "very subtle inner shadow or highlight for separation", Required: false},
			},
			PromptHints:    []string{"vector-like", "bold silhouette", "two to three colors", "no texture clutter", "large centered shape"},
			NegativePrompt: []string{"photo realism", "thin lines", "typography", "complex background", "small decorations"},
			Output:         Output{Size: 1024, Format: "png", Alpha: &noAlpha},
		},
		{
			ID:          "ios-monoline-badge",
			Name:        "iOS Monoline Badge",
			Description: "Clean line-art icon with thick rounded strokes, restrained color, and a clear central metaphor.",
			Categories:  []string{"notes", "habits", "calendar", "privacy", "security"},
			BestFor:     []string{"note apps", "habit trackers", "calendar apps", "privacy tools", "security tools"},
			AvoidFor:    []string{"visual brands needing rich 3D depth", "icons with many tiny concepts"},
			Style: Style{
				Aesthetic: "minimal iOS app icon, rounded monoline glyph, balanced badge composition",
				Palette:   []string{"solid brand background", "light line symbol", "small accent color"},
				Lighting:  "flat with minimal shadow",
				Material:  "clean vector linework",
			},
			Layers: []Layer{
				{Name: "badge-base", Role: "base", Description: "solid or very subtle gradient background with strong contrast", Required: true},
				{Name: "monoline-symbol", Role: "primary", Description: "thick rounded line symbol optically centered and readable at 60 px", Required: true},
				{Name: "accent-dot", Role: "detail", Description: "optional small accent that supports meaning without adding clutter", Required: false},
			},
			PromptHints:    []string{"thick rounded stroke", "optical center", "large simple shape", "no tiny details", "consistent stroke weight"},
			NegativePrompt: []string{"thin hairline strokes", "uneven stroke widths", "text", "complex illustration", "photographic detail"},
			Output:         Output{Size: 1024, Format: "png", Alpha: &noAlpha},
		},
		{
			ID:          "ios-neon-depth",
			Name:        "iOS Neon Depth",
			Description: "Dark energetic icon with a luminous foreground form, controlled glow, and crisp depth.",
			Categories:  []string{"game", "music", "creator-tool", "developer-tool", "crypto"},
			BestFor:     []string{"games", "music apps", "creator tools", "developer utilities", "high-energy products"},
			AvoidFor:    []string{"calm wellness apps", "formal enterprise tools", "low-power accessibility-first icons"},
			Style: Style{
				Aesthetic: "dramatic iOS app icon, controlled neon glow, dark dimensional base, crisp luminous symbol",
				Palette:   []string{"charcoal or black base", "electric accent", "secondary warm or cool glow"},
				Lighting:  "rim light, controlled bloom, subtle depth shadow",
				Material:  "glowing glass or polished synthetic material",
			},
			Layers: []Layer{
				{Name: "dark-base", Role: "base", Description: "dark high-contrast background with restrained depth", Required: true},
				{Name: "glow-core", Role: "primary", Description: "luminous central symbol with crisp edges", Required: true},
				{Name: "rim-light", Role: "detail", Description: "controlled rim light that separates the symbol from the base", Required: false},
				{Name: "depth-shadow", Role: "detail", Description: "faint shadow grounding the foreground form", Required: false},
			},
			PromptHints:    []string{"controlled neon glow", "crisp silhouette", "avoid cyberpunk clutter", "high contrast", "minimal bloom"},
			NegativePrompt: []string{"excessive glow", "purple-only palette", "busy cyberpunk city", "blur", "text", "small particles"},
			Output:         Output{Size: 1024, Format: "png", Alpha: &noAlpha},
		},
		{
			ID:          "ios-paper-cut",
			Name:        "iOS Paper Cut",
			Description: "Layered paper-cut icon with stacked shapes, soft shadows, and a friendly graphic symbol.",
			Categories:  []string{"education", "reading", "maps", "planning", "family"},
			BestFor:     []string{"education apps", "reading apps", "planning tools", "maps", "family apps"},
			AvoidFor:    []string{"hard tech brands", "dark gaming apps", "photorealistic product icons"},
			Style: Style{
				Aesthetic: "layered paper-cut iOS app icon, clean geometric shapes, friendly depth, simple composition",
				Palette:   []string{"fresh brand color", "paper white", "contrasting accent"},
				Lighting:  "soft overhead light with delicate cast shadows",
				Material:  "cut paper layers",
			},
			Layers: []Layer{
				{Name: "paper-base", Role: "base", Description: "flat paper-like background inside iOS safe area", Required: true},
				{Name: "cutout-symbol", Role: "primary", Description: "large cutout or raised paper symbol representing the app", Required: true},
				{Name: "stacked-accent", Role: "middle", Description: "one or two supporting paper layers that add depth and meaning", Required: false},
				{Name: "soft-shadow", Role: "detail", Description: "soft shadow between layers without making the icon muddy", Required: false},
			},
			PromptHints:    []string{"clean paper layers", "simple geometry", "soft shadow", "avoid beige-heavy palette", "high small-size legibility"},
			NegativePrompt: []string{"messy craft texture", "low contrast", "too many paper scraps", "text", "photographic paper scene"},
			Output:         Output{Size: 1024, Format: "png", Alpha: &noAlpha},
		},
		{
			ID:          "nelson-precision",
			Name:        "Nelson Precision",
			Description: "Brand-grade iOS/macOS icon with a single refined metaphor, precise geometry, tactile depth, and restrained color.",
			Categories:  []string{"premium", "macos", "ios", "productivity", "developer-tool", "utility"},
			BestFor:     []string{"polished macOS apps", "premium iOS apps", "developer tools", "productivity tools", "brand-defining app icons"},
			AvoidFor:    []string{"busy concepts requiring multiple objects", "icons that depend on text", "photographic product shots"},
			Style: Style{
				Aesthetic: "premium Apple-platform app icon, refined geometric metaphor, editorial precision, tactile but restrained depth, brand-grade finish",
				Palette:   []string{"one dominant brand color", "one neutral highlight", "one deep shadow tone"},
				Lighting:  "careful studio lighting, soft contact shadows, subtle rim highlights, no dramatic glare",
				Material:  "polished vector-like surfaces with selective 3D depth",
			},
			Layers: []Layer{
				{Name: "platform-base", Role: "base", Description: "clean rounded-square base with subtle depth and strong silhouette", Required: true},
				{Name: "core-metaphor", Role: "primary", Description: "one centered geometric symbol that communicates the app idea without secondary objects", Required: true},
				{Name: "material-edge", Role: "middle", Description: "precise bevel, inset, or overlap that makes the symbol feel intentional and crafted", Required: false},
				{Name: "controlled-shadow", Role: "detail", Description: "soft shadow or highlight used only to separate planes and preserve legibility", Required: false},
			},
			PromptHints: []string{
				"single iconic metaphor",
				"precise geometry",
				"large centered symbol",
				"premium macOS and iOS icon craft",
				"restrained palette",
				"crisp edges",
				"small-size legibility",
				"directly usable final icon",
			},
			NegativePrompt: []string{
				"text",
				"letters",
				"multiple competing objects",
				"decorative clutter",
				"busy background",
				"stock illustration",
				"photorealistic scene",
				"excessive glow",
				"random accent dots",
				"brand logo replica",
			},
			Output: Output{Size: 1024, Format: "png", Alpha: &noAlpha},
		},
		{
			ID:          "pixelresort-object",
			Name:        "Pixelresort Object",
			Description: "Playful object-based Apple-platform icon with saturated color, physical metaphor, and polished dimensional craft.",
			Categories:  []string{"playful", "macos", "ios", "utility", "creative-tool", "consumer-app"},
			BestFor:     []string{"macOS utilities", "creative tools", "consumer apps", "apps with a tangible object metaphor", "friendly productivity apps"},
			AvoidFor:    []string{"strictly minimal monochrome brands", "apps that need flat corporate styling", "icons that cannot use a physical metaphor"},
			Style: Style{
				Aesthetic: "playful premium app icon, tangible object metaphor, saturated color, polished 3D illustration, friendly dimensional realism",
				Palette:   []string{"bright dominant color", "contrasting accent", "soft highlight", "deep ambient shadow"},
				Lighting:  "studio lighting with soft shadows, glossy highlights, and clear object separation",
				Material:  "toy-like polished surfaces, soft plastic, enamel, glass, paper, or crafted physical materials",
			},
			Layers: []Layer{
				{Name: "rounded-base", Role: "base", Description: "vibrant rounded-square platform base with subtle depth and shadow", Required: true},
				{Name: "physical-object", Role: "primary", Description: "one large tangible object that acts as the app metaphor", Required: true},
				{Name: "object-detail", Role: "middle", Description: "a small functional detail that makes the object recognizable without clutter", Required: false},
				{Name: "cast-shadow", Role: "detail", Description: "soft contact shadow and highlight that make the object feel placed in the icon", Required: false},
			},
			PromptHints: []string{
				"single physical metaphor",
				"playful but professional",
				"bright saturated color",
				"polished dimensional object",
				"large centered composition",
				"friendly tactile finish",
				"clear silhouette",
				"small-size legibility",
				"directly usable final icon",
			},
			NegativePrompt: []string{
				"text",
				"letters",
				"multiple props",
				"busy scene",
				"photographic realism",
				"flat generic glyph",
				"random decorative elements",
				"messy texture",
				"brand logo replica",
				"tiny unreadable details",
			},
			Output: Output{Size: 1024, Format: "png", Alpha: &noAlpha},
		},
		{
			ID:          "classic-mac-skeuo",
			Name:        "Classic Mac Skeuo",
			Description: "Classic macOS-style skeuomorphic icon with an isolated desktop object, angled perspective, rich material detail, and soft shadow.",
			Categories:  []string{"classic-macos", "skeuomorphic", "desktop", "utility", "creative-tool", "productivity"},
			BestFor:     []string{"macOS utilities", "desktop tools", "creative apps", "file or document apps", "apps with a strong real-world object metaphor"},
			AvoidFor:    []string{"strictly flat modern icons", "icons that need only an abstract glyph", "brands that require a rounded-square iOS base"},
			Style: Style{
				Aesthetic: "classic pre-flat macOS app icon, rich skeuomorphic desktop object, detailed materials, polished realism, slightly playful Apple Aqua-era craft",
				Palette:   []string{"natural object colors", "metallic grays", "glass highlights", "small saturated accent"},
				Lighting:  "top-left studio light, glossy highlights, soft ambient occlusion, realistic cast shadow on a transparent or plain background",
				Material:  "brushed metal, glass, paper, plastic, leather, mesh, or glossy enamel depending on the object metaphor",
			},
			Layers: []Layer{
				{Name: "object-silhouette", Role: "primary", Description: "one recognizable standalone object with a strong silhouette, not a flat glyph", Required: true},
				{Name: "perspective-volume", Role: "middle", Description: "subtle three-quarter perspective, thickness, bevels, or curved surfaces", Required: true},
				{Name: "material-detail", Role: "detail", Description: "selective texture and functional details that make the object feel tangible", Required: false},
				{Name: "cast-shadow", Role: "base", Description: "soft realistic shadow grounding the object without adding a separate badge or background scene", Required: false},
			},
			PromptHints: []string{
				"standalone desktop object",
				"classic macOS skeuomorphic icon",
				"three-quarter perspective",
				"rich material highlights",
				"soft realistic shadow",
				"plain or transparent-looking background",
				"no rounded-square badge unless the object itself is a document or device",
				"directly usable final icon",
			},
			NegativePrompt: []string{
				"text",
				"letters",
				"modern flat glyph",
				"generic rounded-square app tile",
				"busy scene",
				"multiple unrelated objects",
				"photo collage",
				"low-detail clip art",
				"brand logo replica",
				"tiny unreadable labels",
			},
			Output: Output{Size: 1024, Format: "png", Alpha: &noAlpha},
		},
		{
			ID:          "glossy-toy-3d",
			Name:        "Glossy Toy 3D",
			Description: "Bright rounded app icon with one glossy toy-like 3D object, soft cartoon volume, and high-saturation consumer appeal.",
			Categories:  []string{"glossy", "toy", "3d", "consumer-app", "game", "lifestyle"},
			BestFor:     []string{"consumer apps", "games", "lifestyle apps", "kids and family apps", "playful utilities", "social or entertainment products"},
			AvoidFor:    []string{"serious enterprise tools", "minimal monochrome brands", "apps requiring technical precision", "flat icon systems"},
			Style: Style{
				Aesthetic: "glossy toy-like 3D app icon, rounded soft forms, high-saturation color, cute consumer-grade polish, simple App Store composition",
				Palette:   []string{"bright candy-like base color", "contrasting glossy object color", "white specular highlights", "soft shadow tone"},
				Lighting:  "large softbox highlights, glossy reflections, soft ambient shadow, gentle top-front lighting",
				Material:  "shiny plastic, enamel, soft rubber, polished toy material, or glossy candy-like surface",
			},
			Layers: []Layer{
				{Name: "rounded-color-base", Role: "base", Description: "simple rounded-square base in a bright solid or subtle gradient color", Required: true},
				{Name: "glossy-object", Role: "primary", Description: "one large centered toy-like 3D object, animal head, vehicle, heart, tool, or product metaphor", Required: true},
				{Name: "specular-highlight", Role: "detail", Description: "clear glossy highlight that sells the shiny material without obscuring the silhouette", Required: false},
				{Name: "soft-ground-shadow", Role: "detail", Description: "subtle contact shadow behind or below the object for depth", Required: false},
			},
			PromptHints: []string{
				"one object only",
				"glossy toy-like 3D render",
				"bright saturated colors",
				"soft rounded forms",
				"large centered subject",
				"cute but clean",
				"simple rounded-square base",
				"strong specular highlights",
				"directly usable final icon",
			},
			NegativePrompt: []string{
				"text",
				"letters",
				"multiple objects",
				"busy background",
				"realistic photo",
				"flat vector glyph",
				"sharp industrial edges",
				"tiny details",
				"extra badges",
				"random decorations",
				"brand logo replica",
			},
			Output: Output{Size: 1024, Format: "png", Alpha: &noAlpha},
		},
		{
			ID:          "lance-brand-mark",
			Name:        "Lance Brand Mark",
			Description: "Clean brand-mark icon with a clever geometric symbol, friendly curves, mascot potential, and strong logo-system readability.",
			Categories:  []string{"brand", "logo", "mark", "startup", "saas", "mascot"},
			BestFor:     []string{"startup apps", "SaaS products", "AI tools", "creator tools", "social products", "apps needing a memorable brand mark"},
			AvoidFor:    []string{"photorealistic icons", "heavy skeuomorphic objects", "icons that require detailed illustration", "templates that depend on text"},
			Style: Style{
				Aesthetic: "modern logo-driven app icon, clever geometric mark, friendly rounded curves, simple mascot or symbol, highly readable brand system",
				Palette:   []string{"one strong brand color", "one supporting accent", "neutral negative space"},
				Lighting:  "flat or nearly flat, optional very subtle depth only for app icon polish",
				Material:  "clean vector brand mark, smooth filled shapes, crisp edges",
			},
			Layers: []Layer{
				{Name: "brand-field", Role: "base", Description: "simple solid or subtle-gradient field that supports strong negative space", Required: true},
				{Name: "logomark", Role: "primary", Description: "one clever geometric mark or mascot silhouette built from a small number of bold shapes", Required: true},
				{Name: "negative-space-hook", Role: "middle", Description: "intentional cutout, overlap, or visual pun that makes the mark memorable", Required: false},
				{Name: "icon-polish", Role: "detail", Description: "very subtle rounding, shadow, or highlight only if needed for app icon presence", Required: false},
			},
			PromptHints: []string{
				"logo-first app icon",
				"clever geometric mark",
				"friendly rounded shapes",
				"bold silhouette",
				"simple mascot allowed",
				"intentional negative space",
				"minimal color count",
				"works as a standalone brand mark",
				"directly usable final icon",
			},
			NegativePrompt: []string{
				"text",
				"letters",
				"wordmark",
				"complex illustration",
				"photorealistic render",
				"busy background",
				"many small shapes",
				"generic clip art",
				"brand logo replica",
				"random decorations",
			},
			Output: Output{Size: 1024, Format: "png", Alpha: &noAlpha},
		},
	}
}

func LoadAll(customDir string) ([]Template, error) {
	templates := BuiltIns()
	if customDir == "" {
		return sorted(templates), nil
	}
	entries, err := os.ReadDir(customDir)
	if err != nil {
		if os.IsNotExist(err) {
			return sorted(templates), nil
		}
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		tpl, err := LoadFile(filepath.Join(customDir, entry.Name()))
		if err != nil {
			return nil, err
		}
		templates = upsert(templates, tpl)
	}
	return sorted(templates), nil
}

func LoadFile(path string) (Template, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Template{}, err
	}
	var tpl Template
	if err := json.Unmarshal(data, &tpl); err != nil {
		return Template{}, fmt.Errorf("read template %s: %w", path, err)
	}
	if err := tpl.Validate(); err != nil {
		return Template{}, err
	}
	return tpl, nil
}

func (t Template) Validate() error {
	if strings.TrimSpace(t.ID) == "" {
		return fmt.Errorf("template id is required")
	}
	if strings.TrimSpace(t.Name) == "" {
		return fmt.Errorf("template %s name is required", t.ID)
	}
	if strings.TrimSpace(t.Style.Aesthetic) == "" {
		return fmt.Errorf("template %s style.aesthetic is required", t.ID)
	}
	requiredLayer := false
	for _, layer := range t.Layers {
		if strings.TrimSpace(layer.Name) == "" {
			return fmt.Errorf("template %s contains a layer without name", t.ID)
		}
		if layer.Required {
			requiredLayer = true
		}
	}
	if !requiredLayer {
		return fmt.Errorf("template %s must include at least one required layer", t.ID)
	}
	return nil
}

func Find(templates []Template, id string) (Template, bool) {
	for _, tpl := range templates {
		if tpl.ID == id {
			return tpl, true
		}
	}
	return Template{}, false
}

func sorted(templates []Template) []Template {
	sort.Slice(templates, func(i, j int) bool { return templates[i].ID < templates[j].ID })
	return templates
}

func upsert(templates []Template, next Template) []Template {
	for i := range templates {
		if templates[i].ID == next.ID {
			templates[i] = next
			return templates
		}
	}
	return append(templates, next)
}
