# Provider Notes

Date: 2026-05-27

## Usable Provider From The Supplied Config

The config blob includes many credentials that are not relevant to AppIcon generation. The practical candidate for this CLI is the OpenAI-compatible provider entry:

- Provider type: OpenAI image generation
- Base URL for this CLI: `https://api.openai.com/v1/images/generations`
- Suggested model: `gpt-image-1.5`, or `gpt-image-1` if the account does not have `gpt-image-1.5` access
- Key storage: use `env:OPENAI_API_KEY` in config instead of writing the raw key to disk

OpenAI's official image docs describe the Image API `generations` endpoint for one-shot image generation and GPT Image models including `gpt-image-1.5`, `gpt-image-1`, and `gpt-image-1-mini`.

Sources:
- [OpenAI Image generation guide](https://platform.openai.com/docs/guides/image-generation)
- [OpenAI Images API reference](https://platform.openai.com/docs/api-reference/images/create)

## Recommended Test Setup

```sh
export OPENAI_API_KEY="..."

./dist/appicon init \
  --provider nanobanana \
  --api-key "env:OPENAI_API_KEY" \
  --base-url "https://api.openai.com/v1/images/generations" \
  --model "gpt-image-1.5"

./dist/appicon generate \
  --app "Focus Timer" \
  --idea "a calm deep work timer" \
  --platform iPhone \
  --template ios-clay-symbol \
  --count 1 \
  --out ./appicons \
  --json
```

If `gpt-image-1.5` is not available for the account, retry with:

```sh
./dist/appicon init \
  --provider nanobanana \
  --api-key "env:OPENAI_API_KEY" \
  --base-url "https://api.openai.com/v1/images/generations" \
  --model "gpt-image-1"
```

## Security Note

The supplied config contains live-looking database, OAuth, mail, LLM, storage, payment, ClickHouse, and service-account secrets. Treat it as compromised and rotate it before committing, sharing, or using it outside a private staging environment.
