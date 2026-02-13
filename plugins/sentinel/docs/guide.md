# Sentinel User Guide

Sentinel is an AI-powered code standards evaluation plugin for Watchman. It provides semantic validation that goes beyond pattern matching, using LLMs to evaluate code changes against natural language standards.

## Why Sentinel?

Watchman's built-in rules are fast and deterministic but limited to patterns:
- Can block `rm -rf /` but not "functions should have meaningful names"
- Can enforce file naming but not "error messages should be actionable"

Sentinel fills this gap by sending code to an AI provider that understands context and intent.

```
Claude Code
    │
    ▼
Watchman ──────► Layer 1: Deterministic rules
    │                    - workspace boundaries
    │                    - scope restrictions
    │                    - git policies
    │
    ▼
Sentinel ──────► Layer 2: AI evaluation
                         - semantic standards
                         - style guidelines
                         - domain conventions
```

## Quick Start

### 1. Install

```bash
# From watchman root
make install-all

# Or just sentinel
cd plugins/sentinel && make install
```

### 2. Configure Provider

Create `~/.config/sentinel/config.yml`:

```yaml
provider: openai

openai:
  api_key: ${OPENAI_API_KEY}
  model: gpt-4o-mini
  max_tokens: 1024
```

Or use environment variable directly:

```bash
export OPENAI_API_KEY=sk-...
```

### 3. Create Standards

Create `AGENTS.md` in your project root:

```markdown
# Code Standards

## Documentation
- All exported functions MUST have a doc comment

## Naming
- No single-letter variables except loop indices
- Function names must be verbs (Get, Set, Create, Handle)

## Magic Numbers
- No literal numbers except 0 and 1
- Use named constants
```

### 4. Add to Watchman

Add to `.watchman.yml`:

```yaml
hooks:
  - name: sentinel
    command: sentinel
    tools: [Write, Edit]
    paths: ["**/*.go"]
    timeout: 30s
    on_error: allow
```

### 5. Test

```bash
# Should be denied (no doc comment)
echo '{"tool_name":"Write","tool_input":{"file_path":"test.go","content":"package main\n\nfunc Foo() {}"},"paths":["test.go"],"working_dir":"/tmp"}' | sentinel
```

## Providers

Sentinel supports multiple AI providers with automatic fallback.

### Anthropic (Claude)

```yaml
provider: anthropic

anthropic:
  api_key: ${ANTHROPIC_API_KEY}
  model: claude-sonnet-4-20250514
  max_tokens: 1024
```

**Cost**: ~$0.003/1K input, $0.015/1K output. Light usage: ~$1-2/day.

### OpenAI (GPT)

```yaml
provider: openai

openai:
  api_key: ${OPENAI_API_KEY}
  model: gpt-4o-mini
  max_tokens: 1024
```

**Cost**: ~$0.00015/1K input, $0.0006/1K output. Very economical.

### Ollama (Local)

```yaml
provider: ollama

ollama:
  host: http://localhost:11434
  model: llama3
```

**Cost**: Free (runs locally). Requires Ollama installed.

```bash
# Install Ollama
curl -fsSL https://ollama.com/install.sh | sh

# Pull a model
ollama pull llama3
```

## Provider Fallback

Configure multiple providers with automatic fallback. Sentinel tries each in order until one succeeds.

```yaml
providers:
  - ollama      # Try local first (free)
  - openai      # Then OpenAI (cheap)
  - anthropic   # Finally Anthropic (best quality)

ollama:
  host: http://localhost:11434
  model: llama3

openai:
  api_key: ${OPENAI_API_KEY}
  model: gpt-4o-mini

anthropic:
  api_key: ${ANTHROPIC_API_KEY}
  model: claude-sonnet-4-20250514
```

Use cases:
- **Cost optimization**: Ollama (free) → OpenAI (cheap) → Anthropic (quality)
- **Reliability**: Multiple paid providers as backup
- **Development**: Ollama locally, paid providers in CI

## Configuration Reference

Full config with all options:

```yaml
# Primary provider (if providers list is empty)
provider: anthropic

# Fallback chain (overrides provider if set)
providers:
  - ollama
  - openai
  - anthropic

anthropic:
  api_key: ${ANTHROPIC_API_KEY}
  model: claude-sonnet-4-20250514
  max_tokens: 1024

openai:
  api_key: ${OPENAI_API_KEY}
  model: gpt-4o-mini
  max_tokens: 1024

ollama:
  host: http://localhost:11434
  model: llama3

standards:
  file: AGENTS.md          # Standards document path
  cache_ttl: 5m            # Cache duration

evaluation:
  default_decision: allow  # If evaluation fails
  max_content_size: 50000  # Max bytes to evaluate
  timeout: 25s             # Per-evaluation timeout
```

### Config File Locations

Sentinel looks for config in order:
1. `~/.config/sentinel/config.yml`
2. `.sentinel.yml` (project root)

### Environment Variables

API keys support `${VAR}` expansion:

```yaml
anthropic:
  api_key: ${ANTHROPIC_API_KEY}
```

## Writing Effective Standards

### Be Specific

```markdown
# Bad
- Code should be clean

# Good
- Functions must not exceed 50 lines
- No more than 3 levels of nesting
```

### Provide Context

```markdown
# Good
## Naming
- Prefer short names in small scopes: `cfg` over `configuration`
- Never use "Helper" or "Utils" suffixes (lazy naming)
```

### Group by Category

```markdown
# Documentation
- Exported functions need doc comments
- Doc comments start with function name

# Naming
- Use camelCase for functions
- Use PascalCase for types

# Error Handling
- Always return errors, never panic
- Wrap errors with context
```

## Decisions

Sentinel returns one of three decisions:

| Decision | Meaning | Action |
|----------|---------|--------|
| `allow` | Compliant | Watchman permits |
| `advise` | Minor issues | Warning shown, action permitted |
| `deny` | Violates standards | Watchman blocks |

## Troubleshooting

### "cannot decode input"

Sentinel receives malformed JSON. Test with:

```bash
echo '{"tool_name":"Write","tool_input":{},"paths":[],"working_dir":"/tmp"}' | sentinel
```

### "401 Unauthorized"

API key missing or invalid:

```bash
# Check key is set
echo ${ANTHROPIC_API_KEY:+set}
echo ${OPENAI_API_KEY:+set}
```

### "cannot load standards"

`AGENTS.md` not found. Either:
- Create it in project root
- Set absolute path in config: `standards.file: /path/to/AGENTS.md`

### "all providers failed"

When using fallback, all providers errored. Check:
- Network connectivity
- API keys for each provider
- Ollama is running (`ollama serve`)

### Evaluation too slow

- Increase timeout: `timeout: 45s`
- Use faster model: `gpt-4o-mini` instead of `gpt-4o`
- Use local Ollama with smaller model

## Integration with Watchman

### Hook Configuration

```yaml
hooks:
  - name: sentinel
    command: sentinel
    tools: [Write, Edit]      # Which tools trigger evaluation
    paths: ["**/*.go"]        # Which files to evaluate
    timeout: 30s              # Max time for hook
    on_error: allow           # What to do if hook fails
```

### Selective Evaluation

Only evaluate certain files:

```yaml
hooks:
  - name: sentinel-go
    command: sentinel
    tools: [Write, Edit]
    paths: ["**/*.go", "!**/*_test.go"]  # Skip tests
    timeout: 30s

  - name: sentinel-ts
    command: sentinel
    tools: [Write, Edit]
    paths: ["**/*.ts", "**/*.tsx"]
    timeout: 30s
```

### Per-Project Standards

Use different AGENTS.md per project:

```yaml
# In .sentinel.yml at project root
standards:
  file: ./docs/CODING_STANDARDS.md
```

## Examples

### Go Project

`AGENTS.md`:
```markdown
# Go Standards

## Documentation
- Exported functions must have doc comments
- Doc comment starts with function name

## Error Handling
- Return errors, don't panic
- Wrap with fmt.Errorf and %w
- Error messages lowercase, no period

## Naming
- Receivers: single letter of type (s for Server)
- Interfaces: -er suffix (Reader, Writer)
```

### TypeScript Project

`AGENTS.md`:
```markdown
# TypeScript Standards

## Types
- No `any` - use `unknown` and narrow
- Prefer interfaces over type aliases
- Export types from index.ts

## Functions
- Max 3 parameters, use options object otherwise
- Async functions return Promise<T>, not T | Promise<T>

## React
- Functional components only
- Custom hooks start with "use"
- Props interface named ComponentNameProps
```
