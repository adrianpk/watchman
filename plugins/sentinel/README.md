# Sentinel

AI-powered standards evaluation plugin for Watchman.

## Overview

Sentinel evaluates code actions against standards defined in a document (e.g., `AGENTS.md`) using an AI provider. It acts as a second layer of validation within Watchman's hook system.

## Installation

```bash
cd plugins/sentinel
go build -o sentinel .
```

Move the binary to your PATH or reference it directly in Watchman config.

## Configuration

Create `~/.config/sentinel/config.yml` or `.sentinel.yml` in your project:

```yaml
provider: anthropic

anthropic:
  api_key: ${ANTHROPIC_API_KEY}
  model: claude-sonnet-4-20250514
  max_tokens: 1024

standards:
  file: AGENTS.md
  cache_ttl: 5m

evaluation:
  default_decision: allow
  timeout: 25s
```

## Watchman Integration

Add to your `.watchman.yml`:

```yaml
hooks:
  - name: sentinel
    command: sentinel
    tools: [Write, Edit]
    paths: ["src/**", "internal/**"]
    timeout: 30s
    on_error: allow
```

## Supported Providers

- `anthropic` - Claude models (implemented)
- `openai` - GPT models (planned)
- `ollama` - Local models (planned)

## How It Works

1. Watchman intercepts a tool call (e.g., Write)
2. Sends hook input to Sentinel via stdin
3. Sentinel loads standards document
4. Sends evaluation request to AI provider
5. AI returns structured decision (allow/advise/deny)
6. Sentinel returns decision to Watchman

## Standards Document

The standards document should contain clear, evaluable guidelines. Example:

```markdown
# Code Standards

## Naming
- Use camelCase for functions
- Use PascalCase for types

## Error Handling
- Always return errors, never panic
- Use descriptive error messages
```

## Environment Variables

- `ANTHROPIC_API_KEY` - Required for Anthropic provider
- `OPENAI_API_KEY` - Required for OpenAI provider
