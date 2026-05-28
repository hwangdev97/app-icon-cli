# appicon

`appicon` is a Go CLI for generating iPhone App Icons from AI image APIs such as Image2 or Nanobanana. It stores user configuration as JSON under `~/.appicon-cli/config.json` by default and uses flexible JSON templates to build detailed prompts.

## Install

```sh
go install ./cmd/appicon
```

## Configure

Interactive:

```sh
appicon init
```

Scriptable:

```sh
appicon init \
  --provider nanobanana \
  --api-key "$NANOBANANA_API_KEY" \
  --base-url "https://api.nanobanana.example/v1/images" \
  --model "app-icon"
```

To avoid storing raw keys in `~/.appicon-cli/config.json`, use an environment reference:

```sh
appicon init \
  --provider nanobanana \
  --api-key "env:OPENAI_API_KEY" \
  --base-url "https://api.openai.com/v1/images/generations" \
  --model "gpt-image-1.5"
```

Use `APPICON_CLI_HOME` to move config and templates, for example in tests or agents:

```sh
APPICON_CLI_HOME=/tmp/appicon appicon config show
```

## Generate

```sh
appicon generate \
  --app "Focus Timer" \
  --idea "a calm productivity timer for deep work" \
  --template ios-clay-symbol \
  --richness balanced \
  --count 4 \
  --out ./icons
```

Use `--dry-run` to inspect the provider request without calling an API:

```sh
appicon generate --app "Focus Timer" --idea "deep work timer" --dry-run
```

Use `--json` for machine-readable result output:

```sh
appicon generate --app "Focus Timer" --json
```

## Templates

Built-in templates:

```sh
appicon templates list
appicon templates show lance-brand-mark
```

Custom templates live in the directory configured by `templates_dir` in `~/.appicon-cli/config.json`. A custom template with the same `id` overrides the built-in template.

Template shape:

```json
{
  "id": "ios-custom",
  "name": "iOS Custom",
  "description": "A clear iPhone app icon style.",
  "categories": ["productivity", "utility"],
  "best_for": ["simple apps", "brandable tools"],
  "avoid_for": ["photo-heavy apps"],
  "style": {
    "aesthetic": "modern iOS app icon, simple, high legibility",
    "palette": ["brand color", "white", "dark accent"],
    "lighting": "soft studio lighting",
    "material": "matte 3D"
  },
  "layers": [
    {
      "name": "background",
      "role": "base",
      "description": "rounded-square safe-area background",
      "required": true
    },
    {
      "name": "symbol",
      "role": "primary",
      "description": "single centered metaphor for the app",
      "required": true
    }
  ],
  "prompt_hints": ["no words", "no letters", "centered composition"],
  "negative_prompt": ["text", "letters", "screenshots", "Apple hardware replicas"],
  "output": {
    "size": 1024,
    "format": "png",
    "alpha": false
  }
}
```

## Commands

```text
appicon init
appicon config show
appicon templates list
appicon templates show <id>
appicon generate --app <name> [--idea <text>] [--platform <name>] [--richness minimal|simple|balanced|rich|maximal] [--template <id>] [--provider <name>] [--count <n>] [--size <px>] [--out <dir>] [--dry-run] [--json]
appicon preview --app <name> --out preview.html <icon.png> [more.png]
appicon version
```

## Preview

Generate a browser preview alongside new icons:

```sh
appicon generate --app "Focus Timer" --idea "deep work timer" --preview
```

Or preview existing PNG files:

```sh
appicon preview --app "Focus Timer" --platform macOS --out ./appicons/preview.html ./appicons/icon.png
```

The preview uses the original 1024 x 1024 image directly. It does not add a CSS corner radius or mask to the icon image; rounded platform treatment should be applied later in Icon Composer or the target platform pipeline.

## Homebrew

After a tagged release is published, install from the tap:

```sh
brew install hwangdev97/tap/appicon
```
