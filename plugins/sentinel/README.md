# Sentinel

<p align="center">
  <img src="docs/img/sentinel.png" alt="Sentinel" width="400">
</p>

AI-powered code standards evaluation plugin for Watchman.

**[User Guide](docs/guide.md)**

## Overview

Sentinel evaluates code changes against natural language standards defined in `AGENTS.md` using AI. It acts as a semantic validation layer that catches issues pattern matching cannot:

- "Exported functions must have doc comments"
- "No magic numbers"
- "No nested callbacks deeper than 3 levels"

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
| Anthropic | Ready | ~$1-2/day | Best quality |
| OpenAI | Ready | ~$0.50/day | Good balance |
| Ollama | Ready | Free | Local, requires setup |

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

## Troubleshooting

| Error | Cause | Fix |
|-------|-------|-----|
| `401 Unauthorized` | Missing API key | Set `OPENAI_API_KEY` or `ANTHROPIC_API_KEY` |
| `cannot load standards` | No AGENTS.md | Create in project root |
| `all providers failed` | All providers errored | Check keys, network, ollama running |

See [User Guide](docs/guide.md) for detailed troubleshooting.

## Future Considerations

Optimizations under analysis to improve speed, reduce resource usage, and minimize API costs:

- **Skip on heuristic failure**: If Watchman's deterministic rules already deny, skip the AI roundtrip
- **Selective evaluation**: Only invoke sentinel for certain change types (new files, large diffs, specific patterns)
- **Batch mode**: Instead of evaluating every write/edit, evaluate at commit time with the staged diff
- **Chunked evaluation**: Split large diffs into smaller requests (by file or size) to work within provider limits and enable parallelization
