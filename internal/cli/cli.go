package cli

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/hwang/app-icon-cli/internal/config"
	"github.com/hwang/app-icon-cli/internal/generator"
	"github.com/hwang/app-icon-cli/internal/provider"
	appicon "github.com/hwang/app-icon-cli/internal/template"
)

func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		return wizard(stdin, stdout, stderr)
	}
	switch args[0] {
	case "init", "login":
		return configure(args[1:], stdin, stdout, stderr)
	case "config":
		return configCmd(args[1:], stdout, stderr)
	case "templates":
		return templatesCmd(args[1:], stdout, stderr)
	case "generate", "gen":
		return generateCmd(args[1:], stdout, stderr, provider.NewHTTP())
	case "help", "-h", "--help":
		usage(stdout)
		return 0
	default:
		fmt.Fprintf(stderr, "unknown command %q\n\n", args[0])
		usage(stderr)
		return 2
	}
}

func usage(w io.Writer) {
	fmt.Fprintln(w, `appicon creates iPhone App Icons from AI image providers.

Usage:
  appicon init
  appicon templates list
  appicon templates show <id>
  appicon generate --app "My App" --idea "habit tracker" --template ios-clay-symbol --count 3
  appicon config show

Environment:
  APPICON_CLI_HOME  Override config directory. Defaults to ~/.appicon-cli`)
}

func configure(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	fs.SetOutput(stderr)
	providerName := fs.String("provider", "", "provider name: nanobanana or image2")
	apiKey := fs.String("api-key", "", "provider API key")
	baseURL := fs.String("base-url", "", "provider API endpoint")
	model := fs.String("model", "", "provider model")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	cfg, path, err := config.Load()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	reader := bufio.NewReader(stdin)
	name := valueOrPrompt(*providerName, "Provider [nanobanana/image2]", cfg.DefaultProvider, reader, stdout)
	if name == "" {
		name = "nanobanana"
	}
	pc := cfg.Providers[name]
	pc.APIKey = valueOrPromptSecret(*apiKey, "API key", pc.APIKey, reader, stdout)
	pc.BaseURL = valueOrPrompt(*baseURL, "Base URL", pc.BaseURL, reader, stdout)
	pc.Model = valueOrPrompt(*model, "Model", pc.Model, reader, stdout)
	cfg.DefaultProvider = name
	cfg.Providers[name] = pc
	if saved, err := config.Save(cfg); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	} else {
		path = saved
	}
	fmt.Fprintf(stdout, "Saved config to %s\n", path)
	return 0
}

func valueOrPrompt(value, label, current string, reader *bufio.Reader, out io.Writer) string {
	if value != "" {
		return value
	}
	if current != "" {
		fmt.Fprintf(out, "%s [%s]: ", label, current)
	} else {
		fmt.Fprintf(out, "%s: ", label)
	}
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return current
	}
	return line
}

func valueOrPromptSecret(value, label, current string, reader *bufio.Reader, out io.Writer) string {
	if value != "" {
		return value
	}
	mask := ""
	if current != "" {
		mask = " [configured]"
	}
	fmt.Fprintf(out, "%s%s: ", label, mask)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return current
	}
	return line
}

func configCmd(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || args[0] == "show" {
		cfg, path, err := config.Load()
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		type safeProvider struct {
			APIKey  string `json:"api_key,omitempty"`
			BaseURL string `json:"base_url,omitempty"`
			Model   string `json:"model,omitempty"`
		}
		safe := struct {
			Path            string                  `json:"path"`
			Version         int                     `json:"version"`
			DefaultProvider string                  `json:"default_provider"`
			DefaultTemplate string                  `json:"default_template,omitempty"`
			OutputDirectory string                  `json:"output_directory,omitempty"`
			TemplatesDir    string                  `json:"templates_dir,omitempty"`
			Providers       map[string]safeProvider `json:"providers"`
		}{
			Path:            path,
			Version:         cfg.Version,
			DefaultProvider: cfg.DefaultProvider,
			DefaultTemplate: cfg.DefaultTemplate,
			OutputDirectory: cfg.OutputDirectory,
			TemplatesDir:    cfg.TemplatesDir,
			Providers:       map[string]safeProvider{},
		}
		for name, pc := range cfg.Providers {
			key := ""
			if pc.APIKey != "" {
				key = "configured"
			}
			safe.Providers[name] = safeProvider{APIKey: key, BaseURL: pc.BaseURL, Model: pc.Model}
		}
		printJSON(stdout, safe)
		return 0
	}
	fmt.Fprintf(stderr, "unknown config command %q\n", args[0])
	return 2
}

func templatesCmd(args []string, stdout, stderr io.Writer) int {
	cfg, _, err := config.Load()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	templates, err := appicon.LoadAll(cfg.TemplatesDir)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if len(args) == 0 || args[0] == "list" {
		for _, tpl := range templates {
			fmt.Fprintf(stdout, "%s\t%s\t%s\n", tpl.ID, tpl.Name, tpl.Description)
		}
		return 0
	}
	if args[0] == "show" && len(args) == 2 {
		tpl, ok := appicon.Find(templates, args[1])
		if !ok {
			fmt.Fprintf(stderr, "template %q not found\n", args[1])
			return 1
		}
		printJSON(stdout, tpl)
		return 0
	}
	fmt.Fprintln(stderr, "usage: appicon templates list | appicon templates show <id>")
	return 2
}

func generateCmd(args []string, stdout, stderr io.Writer, client provider.Client) int {
	fs := flag.NewFlagSet("generate", flag.ContinueOnError)
	fs.SetOutput(stderr)
	appName := fs.String("app", "", "app name")
	idea := fs.String("idea", "", "app concept")
	platform := fs.String("platform", "iPhone", "target platform, for example iPhone, iPad, macOS, watchOS")
	richness := fs.String("richness", generator.DefaultRichness, "visual richness: minimal, simple, balanced, rich, maximal")
	templateID := fs.String("template", "", "template id")
	providerName := fs.String("provider", "", "provider name")
	count := fs.Int("count", 1, "number of icons")
	size := fs.Int("size", 1024, "icon size")
	outDir := fs.String("out", "", "output directory")
	dryRun := fs.Bool("dry-run", false, "print request without calling provider")
	jsonOut := fs.Bool("json", false, "print JSON result")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *count < 1 || *count > 20 {
		fmt.Fprintln(stderr, "--count must be between 1 and 20")
		return 2
	}
	if *size < 256 || *size > 2048 {
		fmt.Fprintln(stderr, "--size must be between 256 and 2048")
		return 2
	}
	normalizedRichness, err := generator.NormalizeRichness(*richness)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	cfg, _, err := config.Load()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if *providerName == "" {
		*providerName = cfg.DefaultProvider
	}
	if *templateID == "" {
		*templateID = cfg.DefaultTemplate
	}
	if *templateID == "" {
		*templateID = "ios-clay-symbol"
	}
	if *outDir == "" {
		*outDir = cfg.OutputDirectory
	}
	if *appName == "" {
		fmt.Fprintln(stderr, "--app is required")
		return 2
	}
	templates, err := appicon.LoadAll(cfg.TemplatesDir)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	tpl, ok := appicon.Find(templates, *templateID)
	if !ok {
		fmt.Fprintf(stderr, "template %q not found\n", *templateID)
		return 1
	}
	pc, ok := cfg.Providers[*providerName]
	if !ok {
		fmt.Fprintf(stderr, "provider %q is not configured\n", *providerName)
		return 1
	}
	req := generator.Request{AppName: *appName, Idea: *idea, Platform: *platform, Richness: normalizedRichness, Count: *count, Size: *size, Template: tpl}
	prompt := generator.BuildPrompt(req)
	providerReq := provider.GenerateRequest{Provider: *providerName, Prompt: prompt, Count: *count, Size: *size, Config: pc}
	if *dryRun {
		printJSON(stdout, providerReq)
		return 0
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	images, err := client.Generate(ctx, providerReq)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	files, err := generator.SaveImages(*outDir, *appName, tpl.ID, images)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	result := generator.Result{Files: files, Prompt: prompt}
	manifest := filepath.Join(outputDir(*outDir), "last-run.json")
	_ = generator.SaveManifest(manifest, result)
	if *jsonOut {
		printJSON(stdout, result)
		return 0
	}
	for _, file := range files {
		fmt.Fprintf(stdout, "Created %s\n", file)
	}
	fmt.Fprintf(stdout, "Saved run manifest %s\n", manifest)
	return 0
}

func wizard(stdin io.Reader, stdout, stderr io.Writer) int {
	cfg, _, err := config.Load()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	reader := bufio.NewReader(stdin)
	fmt.Fprintln(stdout, "AppIcon CLI")
	if _, ok := cfg.Providers[cfg.DefaultProvider]; !ok {
		fmt.Fprintln(stdout, "No provider configured. Run appicon init first.")
		return 1
	}
	appName := valueOrPrompt("", "App name", "", reader, stdout)
	idea := valueOrPrompt("", "App concept", "", reader, stdout)
	templates, err := appicon.LoadAll(cfg.TemplatesDir)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintln(stdout, "Templates:")
	for i, tpl := range templates {
		fmt.Fprintf(stdout, "  %d. %s - %s\n", i+1, tpl.ID, tpl.Description)
	}
	choice := valueOrPrompt("", "Template number", "1", reader, stdout)
	idx, _ := strconv.Atoi(choice)
	if idx < 1 || idx > len(templates) {
		idx = 1
	}
	countText := valueOrPrompt("", "How many icons", "1", reader, stdout)
	count, _ := strconv.Atoi(countText)
	args := []string{"--app", appName, "--idea", idea, "--template", templates[idx-1].ID, "--count", strconv.Itoa(count)}
	return generateCmd(args, stdout, stderr, provider.NewHTTP())
}

func outputDir(v string) string {
	if v == "" {
		return "appicons"
	}
	return v
}

func printJSON(w io.Writer, v any) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}
