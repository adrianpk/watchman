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

Create a standards file (e.g., `CONVENTIONS.md`) in your project root:

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
  file: CONVENTIONS.md     # Technical standards only (not behavior rules)
  cache_ttl: 5m            # Cache duration

evaluation:
  default_decision: allow  # If evaluation fails
  max_content_size: 50000  # Max bytes to evaluate
  timeout: 25s             # Per-evaluation timeout
  threshold: 0.85          # Confidence threshold (see below)
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

## Writing Standards

Sentinel works with any format - plain text, markdown, bullet points, prose. More structured rules tend to produce fewer false positives, but use whatever works for your project.

### Separate Behavior from Code Rules (Recommended)

Standards files often serve two different audiences:

1. **Agents** need behavior guidance: when to ask, what not to touch, workflow rules
2. **Code evaluation** needs technical rules: naming, structure, patterns to avoid

When Sentinel evaluates a diff against a file containing both, the LLM has no way to distinguish "this rule applies to agent behavior" from "this rule applies to code." It treats everything as potentially relevant to the code under review.

For example:

```markdown
# AGENTS.md (mixed - problematic)

## Behavior
- Don't manage the dev server
- Ask before deleting files

## Code
- No magic numbers
- Functions must have doc comments
```

The LLM might interpret "Don't manage the dev server" as a code rule and reject changes containing "dev" in variable names or comments.

**Recommended approach**: Separate concerns into two files:

```markdown
# AGENTS.md (behavior only)
- Don't manage the dev server
- Ask before deleting files
- All code rules are in CONVENTIONS.md
```

```markdown
# CONVENTIONS.md (code rules only)
- No magic numbers
- Functions must have doc comments
```

Then configure Sentinel to point only to the technical file:

```yaml
standards:
  file: CONVENTIONS.md
```

This way, agents read both files for context, but Sentinel evaluates code only against technical rules.

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

## Confidence Threshold

Sentinel uses a confidence-based evaluation system to reduce false positives. For each potential violation, the AI reports its confidence level (0.0-1.0):

| Confidence | Meaning |
|------------|---------|
| `1.0` | Explicitly prohibited by standards |
| `0.7` | Likely violation, requires some interpretation |
| `0.5` | Ambiguous, could go either way |
| `<0.5` | Probably not a violation |

### How It Works

1. AI evaluates code against standards
2. For each violation found, AI reports confidence
3. If all violations have confidence below threshold → downgrade to `advise`
4. Only block (`deny`) when at least one violation meets threshold

### Configuration

```yaml
evaluation:
  threshold: 0.85  # Default
```

| Threshold | Use Case |
|-----------|----------|
| `0.9` | Very permissive - only block explicit, unambiguous violations |
| `0.85` | Balanced (default) - block clear violations, warn on interpretations |
| `0.7` | Stricter - block likely violations even with some interpretation |
| `0` | Disable confidence - block on any violation the AI finds |

### Example

Standards say: "Assets must live in `./assets/`"

| File | AI Interpretation | Confidence | Result (0.85 threshold) |
|------|-------------------|------------|-------------------------|
| `src/logo.png` | Clear violation | 0.95 | **deny** |
| `assets/migration/001.sql` | Ambiguous - is it IN assets? | 0.6 | advise (warning) |
| `internal/db/sqlc/gen.go` | "DO NOT EDIT" ≠ "don't commit" | 0.4 | advise (warning) |

The threshold prevents the AI from blocking on forced interpretations of rules.

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

Standards file not found. Either:
- Create `CONVENTIONS.md` in project root
- Set path in config: `standards.file: /path/to/your-standards.md`

### "all providers failed"

When using fallback, all providers errored. Check:
- Network connectivity
- API keys for each provider
- Ollama is running (`ollama serve`)

### Evaluation too slow

- Increase timeout: `timeout: 45s`
- Use faster model: `gpt-4o-mini` instead of `gpt-4o`
- Use local Ollama with smaller model

### Too many false positives

The AI may over-interpret rules. Solutions:

1. **Lower the confidence threshold**:
   ```yaml
   evaluation:
     threshold: 0.7  # More permissive (default: 0.85)
   ```

2. **Separate behavior rules from code rules**: Keep agent instructions in `AGENTS.md`, code standards in `CONVENTIONS.md`. Point Sentinel only at the code standards file.

3. **Be more explicit in standards**: Instead of "assets should not be alongside code", write "files in `assets/**` are correctly located"

4. **Check the violations**: When blocked, look at the confidence reported. Low confidence violations indicate the AI is unsure - consider adjusting your standards to be clearer.

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

Use different standards files per project:

```yaml
# In .sentinel.yml at project root
standards:
  file: ./docs/CODING_STANDARDS.md
```

## Examples

### Go Project

`CONVENTIONS.md`:
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

`CONVENTIONS.md`:
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
