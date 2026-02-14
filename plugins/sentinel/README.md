# Sentinel

<p align="center">
  <img src="docs/img/sentinel.png" alt="Sentinel" width="400">
</p>

AI-powered code standards evaluation plugin for Watchman.

**[User Guide](docs/guide.md)**

## Overview

Sentinel evaluates code changes against natural language standards defined in `AGENTS.md`, `CLAUDE.md`, or any specification file you configure. It acts as a semantic validation layer that catches issues pattern matching cannot:

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

| Provider | Status | Notes |
|----------|--------|-------|
| Anthropic | Ready | Claude models |
| OpenAI | Ready | GPT models |
| Ollama | Ready | Local, free, requires setup |

### Fallback Chain

Providers are tried in the order you define until one succeeds:

```yaml
providers:
  - ollama      # First choice
  - openai      # Second choice
  - anthropic   # Third choice
```

Configure the order based on your needs (cost, quality, availability).

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
  mode: all  # or "commits_only"
EOF
```

## Evaluation Modes

### Mode: all (default)

Evaluates every Write/Edit operation in real-time. Provides immediate feedback but incurs more API calls.

### Mode: commits_only

Only evaluates when a git commit is attempted. Evaluates the entire staged diff in a single request.

```yaml
evaluation:
  mode: commits_only
  batch:
    max_files_per_request: 10
    max_content_per_request: 40000
    concurrency: 3
```

**Benefits:**
- Reduces API costs significantly (1 call per commit vs N calls per change)
- Full context of all changes in single evaluation
- Encourages granular commits (works well with Watchman's incremental rule)

**Hook configuration for commits_only:**

```yaml
# .watchman.yml - Git only
hooks:
  - name: sentinel
    command: sentinel
    tools: [Bash]
    match_command: "git.*(add|commit)"
    timeout: 60s
    on_error: allow
```

```yaml
# .watchman.yml - Git and jj (jujutsu)
hooks:
  - name: sentinel
    command: sentinel
    tools: [Bash]
    match_command: "(git.*(add|commit))|(jj.*(new|commit|squash))"
    timeout: 60s
    on_error: allow
```

**Git commands:**
- `git add`: Evaluates files before staging
- `git commit`: Evaluates the staged diff

**jj commands:**
- `jj new`: Evaluates working copy before creating new commit
- `jj commit`: Evaluates working copy before finalizing
- `jj squash`: Evaluates before squashing changes

## Troubleshooting

| Error | Cause | Fix |
|-------|-------|-----|
| `401 Unauthorized` | Missing API key | Set `OPENAI_API_KEY` or `ANTHROPIC_API_KEY` |
| `cannot load standards` | No AGENTS.md | Create in project root |
| `all providers failed` | All providers errored | Check keys, network, ollama running |

See [User Guide](docs/guide.md) for detailed troubleshooting.

