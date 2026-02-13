# Sentinel

AI-powered code standards evaluation plugin for Watchman.

**[Full User Guide](docs/guide.md)**

## Overview

Sentinel evaluates code changes against natural language standards (e.g., `AGENTS.md`) using AI. It acts as a semantic validation layer that catches issues pattern matching cannot:

- "Exported functions must have doc comments"
- "No magic numbers"
- "Error messages should use 'cannot', not 'Failed to'"

## Quick Start

```bash
# 1. Install
make install-all

# 2. Configure (~/.config/sentinel/config.yml)
provider: openai
openai:
  api_key: ${OPENAI_API_KEY}
  model: gpt-4o-mini

# 3. Create AGENTS.md in project root
# 4. Add hook to .watchman.yml
hooks:
  - name: sentinel
    command: sentinel
    tools: [Write, Edit]
    paths: ["**/*.go"]
    timeout: 30s
```

## Providers

| Provider | Status | Cost | Notes |
|----------|--------|------|-------|
| Anthropic | ✅ Ready | ~$1-2/day | Best quality |
| OpenAI | ✅ Ready | ~$0.50/day | Good balance |
| Ollama | ✅ Ready | Free | Local, requires setup |

### Fallback Chain

Try providers in order until one succeeds:

```yaml
providers:
  - ollama      # Free, local
  - openai      # Cheap
  - anthropic   # Quality fallback
```

## Configuration

### Option A: Environment Variable

```bash
export OPENAI_API_KEY=sk-...
# or
export ANTHROPIC_API_KEY=sk-ant-...
```

### Option B: Config File

```bash
mkdir -p ~/.config/sentinel
cat > ~/.config/sentinel/config.yml << 'EOF'
provider: openai

openai:
  api_key: ${OPENAI_API_KEY}
  model: gpt-4o-mini
  max_tokens: 1024

standards:
  file: AGENTS.md
  cache_ttl: 5m

evaluation:
  default_decision: allow
  timeout: 25s
EOF
```

## Testing

### Standalone

```bash
echo '{"tool_name":"Write","tool_input":{"file_path":"test.go","content":"package main\n\nfunc Foo() {}"},"paths":["test.go"],"working_dir":"."}' | sentinel
```

### Expected Output

```json
{"decision":"deny","reason":"Missing doc comment on exported function Foo"}
```

## Troubleshooting

| Error | Cause | Fix |
|-------|-------|-----|
| `401 Unauthorized` | Missing API key | Set `OPENAI_API_KEY` or `ANTHROPIC_API_KEY` |
| `cannot load standards` | No AGENTS.md | Create in project root |
| `all providers failed` | All providers errored | Check keys, network, ollama running |

See [User Guide](../../docs/sentinel.md) for detailed troubleshooting.
