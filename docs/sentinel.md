# Sentinel Plugin

Sentinel is an AI-powered code standards evaluation plugin for Watchman. It provides semantic validation using LLMs to evaluate code changes against natural language standards defined in `AGENTS.md`.

## When to Use Sentinel

Use Sentinel when you need validation that goes beyond pattern matching:

| Watchman (Deterministic) | Sentinel (AI) |
|--------------------------|---------------|
| Block `rm -rf /` | "Functions should have meaningful names" |
| Enforce file naming | "Error messages should be actionable" |
| Forbid `fmt.Print` | "Code should follow domain conventions" |

## Quick Start

```yaml
# .watchman.yml
hooks:
  - name: sentinel
    command: sentinel
    tools: [Write, Edit]
    paths: ["**/*.go"]
    timeout: 30s
```

## Documentation

Full documentation is in the plugin module:

- **[User Guide](../plugins/sentinel/docs/guide.md)** - Installation, configuration, providers, examples
- **[README](../plugins/sentinel/README.md)** - Quick reference
- **[Config Example](../plugins/sentinel/config.example.yml)** - Full configuration template

## Supported Providers

| Provider | Cost | Notes |
|----------|------|-------|
| Ollama | Free | Local, requires setup |
| OpenAI | ~$0.50/day | Good balance |
| Anthropic | ~$1-2/day | Best quality |

Supports automatic fallback chain: `providers: [ollama, openai, anthropic]`
