# App Icon Template Research

Date: 2026-05-27

This document collects candidate App Icon templates for `appicon`. The goal is not to copy a specific competitor's artwork, but to turn common, useful icon design directions into flexible JSON templates that work well with Image2, Nano Banana, or similar image APIs.

## Design Constraints To Bake Into Every Template

Apple's current App Icon guidance should be treated as the baseline:

- Use a square source. The system applies the final rounded mask on iOS, iPadOS, and macOS.
- Keep primary content centered so masking and corner treatment do not crop important content.
- Prefer a simple graphic representation over photos or app screenshots.
- Avoid nonessential text because it is hard to read, hard to localize, and often redundant beside the app name.
- For current Apple platforms, think in layers. Apple references layered icons, Icon Composer, and appearance variants such as default, dark, clear, and tinted appearances.
- Use 1024 x 1024 as the master output target; downstream export can generate Xcode/AppIcon.appiconset sizes.

Sources:
- [Apple Human Interface Guidelines: App icons](https://developer.apple.com/design/human-interface-guidelines/app-icons)
- [Apple Design Resources](https://developer.apple.com/design/resources/)
- [iKit iOS App Store Icon Generator Sizes](https://ikit.app/blog/ios-app-store-icon-generator-sizes?lang=sw)

## Product Pattern Findings

Existing AI app icon tools converge on a few useful product ideas:

- Style presets matter. IconBundlr exposes styles such as minimal, flat, gradient, 3D, neon, line art, and photorealistic.
- Batch generation matters. Several tools generate multiple variations per prompt and let users pick a direction.
- Export structure matters. Tools emphasize Xcode-ready icon sets, `Contents.json`, and correct platform sizes.
- Prompt transparency matters. Icon Lab exposes an editable final prompt and supports prompt/layer/background customization.
- Provider control matters. Icon Lab and TAU both mention bring-your-own-key/local key storage patterns.
- Consistent packs are a useful advanced mode. TAU highlights shared visual language across 20-30 icons.

Sources:
- [IconBundlr](https://iconbundlr.com/)
- [IconCraft](https://www.iconcraft.ai/)
- [Icon Lab](https://iconlab.app/)
- [TAU Icon Generator](https://www.tauicongenerator.com/index.html)

## Recommended Template Fields

Our current JSON template shape is good, but these fields would make templates more expressive:

```json
{
  "id": "ios-template-id",
  "name": "Template Name",
  "description": "Short user-facing explanation.",
  "categories": ["productivity", "utility"],
  "best_for": ["simple apps", "brandable tools"],
  "avoid_for": ["photo-heavy apps"],
  "style": {
    "aesthetic": "visual direction",
    "palette": ["brand color", "white", "dark accent"],
    "lighting": "lighting direction",
    "material": "material direction"
  },
  "layers": [
    {
      "name": "background",
      "role": "base",
      "description": "safe-area background",
      "required": true
    }
  ],
  "prompt_hints": ["no words", "centered composition"],
  "negative_prompt": ["text", "letters", "screenshots", "Apple hardware replicas"],
  "output": {
    "size": 1024,
    "format": "png",
    "alpha": false
  }
}
```

## Richness Parameter

`appicon generate` supports `--richness` as a cross-template control for visual complexity:

- `minimal`: one subject, one to two colors, no optional details, maximum legibility.
- `simple`: one subject, restrained depth, two to three colors, at most one supporting detail.
- `balanced`: default; controlled material, lighting, and recognition details.
- `rich`: more material, stronger lighting, tactile depth, and several crafted details.
- `maximal`: full visual craft and micro-detail while still avoiding multiple competing subjects.

This parameter should not change the app concept. It only controls layer density, material detail, lighting intensity, and how much optional template detail is allowed.

## Candidate Template Pool

### 1. iOS Clay Symbol

Soft 3D icon with a simple tactile object or metaphor. This is one of the safest default styles because it reads clearly at small sizes and works for many app categories.

- Best for: productivity, wellness, education, finance, developer tools.
- Layers: background, soft object, highlight/shadow detail.
- Prompt hints: matte clay, soft depth, centered object, high small-size legibility, no words.
- Risks: can become generic if the symbol is not specific enough.
- Recommendation: keep as a default template.

### 2. iOS Glass / Liquid Layer

Layered translucent panels, refractive highlights, and a crisp foreground glyph. This aligns with newer Apple visual language, but needs stronger constraints to avoid blur and low contrast.

- Best for: AI tools, media apps, premium utilities, dashboards.
- Layers: gradient base, translucent panel, crisp glyph, specular highlight.
- Prompt hints: layered glass, clear foreground silhouette, high contrast, avoid blurry edges.
- Risks: generated icons can look muddy, over-glossy, or illegible under dark/tinted appearances.
- Recommendation: include, but mark as "premium/experimental" and add contrast checks later.

### 3. Flat Gradient Glyph

Modern flat background gradient with one bold geometric symbol. This is fast, robust, and works well for early-stage apps.

- Best for: MVPs, SaaS tools, internal apps, simple utilities.
- Layers: gradient field, geometric glyph, optional subtle inner shadow.
- Prompt hints: vector-like, bold silhouette, two to three colors, no texture clutter.
- Risks: may feel common without strong palette and metaphor.
- Recommendation: include as a practical default.

### 4. Monoline Badge

Line-art symbol on a restrained background. Useful when users want a clean icon that feels less like an AI render.

- Best for: note apps, habit trackers, calendar apps, privacy/security tools.
- Layers: solid or subtle base, rounded monoline symbol, optional small accent.
- Prompt hints: thick rounded stroke, optical center, large simple shape, no tiny details.
- Risks: thin lines fail at small sizes; generated strokes can be inconsistent.
- Recommendation: include with strict "thick stroke" constraints.

### 5. Neon Depth

Dark base with luminous foreground form, rim lighting, and controlled glow. This is visually strong and good for games/tools where energy matters.

- Best for: games, music, creator tools, crypto/web3, developer utilities.
- Layers: dark base, glowing core symbol, rim light, faint depth shadow.
- Prompt hints: controlled neon glow, crisp silhouette, avoid cyberpunk clutter.
- Risks: overused dark-blue/purple look, small-size blooming, poor accessibility.
- Recommendation: include, but constrain palette diversity and glow intensity.

### 6. Skeuomorphic Object

Detailed object rendered as a polished, realistic-but-simplified app icon. This can produce memorable icons, especially for apps tied to a concrete object.

- Best for: camera, audio, cooking, travel, craft, tools.
- Layers: material background, realistic simplified object, highlight, contact shadow.
- Prompt hints: iconic object, simplified realism, no photo, no real brand marks.
- Risks: too much detail, outdated look, or accidental replica of real-world products.
- Recommendation: include as an advanced template.

### 7. Paper Cut / Layered Shape

Stacked paper-like shapes with shadows and a clear symbol. It gives depth without full 3D rendering.

- Best for: education, reading, maps, planning, kids/family apps.
- Layers: paper base, cutout symbol, layered accent, soft shadow.
- Prompt hints: clean paper layers, simple geometry, warm but not beige-heavy.
- Risks: can look like craft stationery if color/contrast is weak.
- Recommendation: include.

### 8. Pixel / Retro

Pixel-art or retro-console icon with strict grid and limited palette. Useful for games and playful apps, but not a general default.

- Best for: games, hobby apps, retro utilities.
- Layers: pixel background, main sprite, border/accent.
- Prompt hints: crisp pixel grid, no antialias blur, limited palette, app icon composition.
- Risks: model may blur pixels; needs post-processing or provider support.
- Recommendation: include only if we add template-specific post-process guidance.

### 9. Minimal Lettermark

A single letter or monogram in a designed shape. Apple allows text when essential to brand, but this should be opt-in because text is usually poor for icons.

- Best for: established brands, initial-based products, simple B2B tools.
- Layers: background, monogram, optional brand accent.
- Prompt hints: one large letter only, no extra words, high contrast.
- Risks: AI text accuracy, localization, generic brand feel.
- Recommendation: include as "brand-lettermark", disabled from default suggestions unless user asks for brand initials.

### 10. Consistent Icon Pack Style

Not a single icon style, but a meta-template for generating multiple related icons with shared geometry, palette, and lighting. TAU's product positioning suggests this is valuable for teams producing variants or icon packs.

- Best for: multiple apps, seasonal variants, A/B tests, product suites.
- Layers: shared base grammar, shared symbol grammar, per-icon metaphor.
- Prompt hints: same visual language, same palette, same camera/lighting, distinct metaphor per app.
- Risks: hard to maintain consistency across providers without reference images.
- Recommendation: add later as a batch mode after single-icon templates are stable.

### 11. Nelson Precision

Inspired by Gavin Nelson's app icon portfolio, this template aims for brand-grade Apple-platform icon craft rather than a generic AI style. The portfolio page describes a selection of iOS and macOS app icons for clients including OpenAI, Linear, 1Password, Flighty, GitHub, Readwise, VS Code, and others.

- Best for: polished macOS apps, premium iOS apps, developer tools, productivity tools.
- Layers: rounded platform base, one centered metaphor, precise bevel/inset/overlap, controlled shadow.
- Prompt hints: single iconic metaphor, precise geometry, restrained palette, crisp edges, tactile depth, small-size legibility.
- Risks: too much specificity could imitate a known brand icon; keep prompts focused on general craft principles and avoid brand logo replicas.
- Recommendation: include as an advanced premium template.

Source:
- [Gavin Nelson: App icon design](https://nelson.co/icon-design)

### 12. Pixelresort Object

Inspired by Michael Flarup's Pixelresort portfolio. The homepage describes the studio as making icons, logos, illustrations, and interfaces, and the work grid includes many App Icon projects across macOS and iOS. The OK Json project is a useful example of the pattern: a Mac JSON formatter represented by a physical tool metaphor, described as a bright orange cutter.

- Best for: macOS utilities, creative tools, consumer apps, friendly productivity apps.
- Layers: rounded platform base, one large physical object metaphor, one recognition detail, contact shadow/highlight.
- Prompt hints: playful but professional, bright saturated color, polished dimensional object, tactile finish, clear silhouette.
- Risks: can become too toy-like or too detailed; constrain to one object and avoid busy scenes.
- Recommendation: include as an advanced object-metaphor template.

Sources:
- [Pixelresort homepage](https://www.pixelresort.com/)
- [Pixelresort OK Json project](https://www.pixelresort.com/project/okjson)

### 13. Classic Mac Skeuo

Based on the older macOS icon style shown in the supplied reference image: standalone desktop objects instead of flat rounded-square tiles. The style uses tangible metaphors such as compasses, folders, disks, calendars, stamps, drives, gauges, documents, cameras, and tools.

- Best for: macOS utilities, desktop tools, creative apps, document apps, object-metaphor products.
- Layers: object silhouette, perspective volume, material detail, cast shadow.
- Prompt hints: three-quarter perspective, rich materials, brushed metal/glass/paper/plastic, soft shadow, plain background.
- Risks: can become too busy or dated; keep one object and avoid tiny labels.
- Recommendation: include as a legacy macOS / skeuomorphic template.

### 14. Glossy Toy 3D

Based on the supplied reference image: a grid of bright, rounded app icons with single glossy 3D subjects such as a heart, airplane, animal head, shoe, car, dumbbell, and apple. The style is more playful and consumer-facing than classic macOS skeuomorphism.

- Best for: consumer apps, games, lifestyle apps, family apps, entertainment products.
- Layers: rounded color base, one glossy object, specular highlight, soft shadow.
- Prompt hints: toy-like 3D render, bright saturated colors, soft rounded forms, large centered subject.
- Risks: can become childish or overdecorated; constrain to one object and avoid extra badges.
- Recommendation: include for playful App Store style icons.

### 15. Lance Brand Mark

Based on Lance Draws' portfolio. The homepage emphasizes logo and mark work, with examples labeled as logomarks, combination marks, wordmarks, and a mascot. The useful app-icon abstraction is not skeuomorphic rendering but a logo-first icon: clever geometry, strong negative space, friendly curves, and a mark that can survive as a brand system.

- Best for: startup apps, SaaS products, AI tools, creator tools, social products.
- Layers: brand field, one logomark or mascot silhouette, negative-space hook, subtle icon polish.
- Prompt hints: logo-first, clever geometric mark, bold silhouette, minimal color count, friendly rounded shapes.
- Risks: AI may add letters or wordmarks; keep `text`, `letters`, and `wordmark` in negative prompts.
- Recommendation: include as a brand identity / logomark template.

Source:
- [Lance Draws](https://lancedraws.com/)

## Shortlist For First Template Release

Recommended first set:

1. `ios-clay-symbol` (implemented)
2. `ios-glass-layer` (implemented)
3. `ios-flat-gradient` (implemented)
4. `ios-monoline-badge` (implemented)
5. `ios-neon-depth` (implemented)
6. `ios-paper-cut` (implemented)
7. `nelson-precision` (implemented)
8. `pixelresort-object` (implemented)
9. `classic-mac-skeuo` (implemented)
10. `glossy-toy-3d` (implemented)
11. `lance-brand-mark` (implemented)

These cover the broadest user needs while keeping generation predictable. I would delay `skeuomorphic-object`, `pixel-retro`, `lettermark`, and `consistent-pack` until we add stronger prompt controls, reference-image support, or post-processing.

## Selection Criteria

When choosing templates for implementation, score each candidate from 1-5:

- Legibility at 60 px.
- Distinctiveness from common AI icon output.
- Works without text.
- Works with arbitrary app categories.
- Low risk of policy/IP issues.
- Easy to express in JSON layers.
- Easy for an AI model to follow consistently.

## Implementation Notes

Near-term changes that would improve template quality:

- Add `negative_prompt` to templates.
- Add `best_for`, `avoid_for`, and `categories` for interactive template recommendations.
- Add `variants` so a template can generate default/dark/tinted prompt variants.
- Add `quality_checks` for generated-image review later, such as centered content, no text, no screenshots, high contrast.
- Add `exports` for Xcode-ready `.appiconset` generation after the base 1024 PNG is created.
