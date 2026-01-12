# Rules

Watchman enforces semantic rules that apply to all tools. Rules block intent, not specific tools.

## Workspace

**Status**: Implemented

Confines the AI agent to the project directory. Blocks access to paths outside the workspace.

### Purpose

Prevents the agent from:
- Reading sensitive files outside the project (`/etc/passwd`, `~/.ssh/`)
- Writing to system directories
- Escaping the project boundary via path traversal (`../`)

### Configuration

```yaml
rules:
  workspace: true

workspace:
  allow:
    - /tmp/
    - ~/.cache/go-build/
  block:
    - .env
    - secrets/
```

### Behavior

| Path | Result |
|------|--------|
| `./src/main.go` | Allowed |
| `/etc/passwd` | Blocked |
| `../other-project/` | Blocked |
| `/tmp/test.txt` | Allowed (if in allow list) |
| `.env` | Blocked (if in block list) |

### Protected Paths

Some paths are always protected regardless of configuration:

- `~/.claude/` - Claude settings and hooks
- `~/.ssh/` - SSH keys
- `~/.aws/` - AWS credentials
- `~/.gnupg/` - GPG keys
- `~/.config/watchman/` - Watchman global config
- `.watchman.yml` - Local config (any directory)

These cannot be overridden.

### All Options Reference

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `allow` | []string | [] | Paths allowed outside workspace |
| `block` | []string | [] | Paths blocked inside workspace |

---

## Scope

**Status**: Implemented

Limits modifications to explicitly declared files within the workspace.

### Purpose

Even within the workspace, restrict which files can be modified. Useful for:
- Focusing changes on specific files
- Protecting generated code
- Preventing accidental edits to unrelated files

### Configuration

```yaml
rules:
  scope: true

scope:
  allow:
    - src/**/*.go
    - go.mod
  block:
    - vendor/**
    - **/*_generated.go
```

### Behavior

| Tool | Scope Applied |
|------|--------------|
| `Read` | No (read-only) |
| `Glob` | No (read-only) |
| `Grep` | No (read-only) |
| `Bash` | No |
| `Write` | Yes |
| `Edit` | Yes |
| `NotebookEdit` | Yes |

### Pattern Matching

Supports glob patterns:
- `*` - matches any characters except path separator
- `**` - matches any characters including path separator (recursive)
- `?` - matches any single character
- `[abc]` - matches character class

| Path | Pattern | Match |
|------|---------|-------|
| `src/main.go` | `src/**/*.go` | Yes |
| `src/pkg/util.go` | `src/**/*.go` | Yes |
| `vendor/lib.go` | `src/**/*.go` | No |
| `types_generated.go` | `**/*_generated.go` | Yes |
| `.env` | `.env` | Yes |

### Rules

1. If no `allow` patterns defined, all paths are allowed
2. If `allow` patterns defined, path must match at least one
3. `block` patterns take precedence over `allow`

### All Options Reference

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `allow` | []string | [] | Glob patterns for allowed files |
| `block` | []string | [] | Glob patterns for blocked files |

---

## Versioning

**Status**: Implemented

Controls commit message format and branch protection for git/jj workflows.

### Purpose

Enforces consistent commit practices:
- Commit message formatting (length, case, punctuation)
- Branch protection (prevent commits to main/master)
- Tool preference (prefer jj over git)

### Configuration

```yaml
rules:
  versioning: true

versioning:
  commit:
    max_length: 72
    require_uppercase: true
    no_period: true
    prefix_pattern: ""  # e.g., "\[JIRA-\d+\]"
  branches:
    protected:
      - main
      - master
  tool: ""  # "jj" to prefer jj over git
```

### Commit Message Rules

| Rule | Description |
|------|-------------|
| `max_length` | Maximum characters (0 = unlimited) |
| `require_uppercase` | First character must be uppercase |
| `no_period` | Must not end with period |
| `prefix_pattern` | Regex pattern message must match |

### Branch Protection

Protected branches block direct commits:

```yaml
versioning:
  branches:
    protected:
      - main
      - master
      - release/*
```

### Operations Block

Block specific git operations:

```yaml
versioning:
  operations:
    block:
      - push --force
      - push -f
```

### Workflow

Enforce a git workflow style:

| Value | Effect |
|-------|--------|
| `""` | No restriction (default) |
| `linear` | Blocks merge, prefer rebase |
| `merge` | Blocks rebase, prefer merge |

```yaml
versioning:
  workflow: linear
```

### Tool Preference

When `tool: jj` is set, blocks `git commit` and suggests `jj commit`:

```yaml
versioning:
  tool: jj
```

Other VCS (mercurial, etc.) not yet supported.

### All Options Reference

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `commit.max_length` | int | 0 | Max message length (0 = unlimited) |
| `commit.require_uppercase` | bool | false | First char must be uppercase |
| `commit.no_period` | bool | false | Must not end with period |
| `commit.prefix_pattern` | string | "" | Regex prefix pattern |
| `branches.protected` | []string | [] | Branches that block commits |
| `operations.block` | []string | [] | Git operations to block |
| `workflow` | string | "" | Workflow style: linear, merge |
| `tool` | string | "" | Preferred VCS: jj |

---

## Incremental

**Status**: Implemented

Limits the number of files modified before requiring a commit or review.

### Purpose

Prevents large-scale rewrites by:
- Tracking modified files via `git status`
- Warning when approaching the limit
- Blocking when the limit is reached

This encourages small, reviewable changes and prevents runaway modifications.

### Configuration

```yaml
rules:
  incremental: true

incremental:
  max_files: 10
  warn_ratio: 0.7
```

### Behavior

| Modified Files | Action |
|----------------|--------|
| 0-6 (under 70%) | Silent, allowed |
| 7-9 (70-99%) | Warning: "approaching file limit" |
| 10+ (at limit) | Blocked: "commit or review before continuing" |

### Warning vs Block

The `warn_ratio` determines when warnings start:
- `0.7` = warn at 70% of max (default)
- `0.5` = warn at 50% of max
- `0` = use default 70%

Warnings give the agent runway to finish current work gracefully. Blocking forces a decision: commit, revert, or adjust the threshold.

### Tools Affected

| Tool | Incremental Applied |
|------|---------------------|
| `Read` | No (read-only) |
| `Glob` | No (read-only) |
| `Grep` | No (read-only) |
| `Bash` | No |
| `Write` | Yes |
| `Edit` | Yes |
| `NotebookEdit` | Yes |

### All Options Reference

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `max_files` | int | 0 | Maximum modified files (0 = unlimited) |
| `warn_ratio` | float | 0.7 | Ratio at which to start warning (0-1) |

---

## Hooks (External Hooks)

**Status**: Implemented

Execute external programs to enforce custom rules. This is the extensibility mechanism for defining project-specific invariants, patterns, and boundaries.

### Purpose

Watchman provides built-in rules, but every project has unique requirements. External hooks allow you to:
- Define custom validation logic in any language
- Integrate with project-specific tools (linters, analyzers)
- Enforce domain-specific constraints
- Build reusable rule packages

### Configuration

```yaml
hooks:
  - name: "no-direct-db-access"
    command: "check-db-access"
    tools: ["Write", "Edit"]
    paths: ["src/handlers/**/*.go"]
    timeout: 5s
    on_error: allow
```

### Hook Protocol

**Input (stdin, JSON):**
```json
{
  "tool_name": "Write",
  "tool_input": {"file_path": "src/main.go", "content": "..."},
  "paths": ["src/main.go"],
  "working_dir": "/path/to/project"
}
```

**Output (stdout, JSON):**
```json
{"decision": "deny", "reason": "Direct database access in handler layer"}
```

Or use exit codes:
- Exit 0 = allow
- Exit 1 = deny (stderr = reason)

**Decision Values:**
| Decision | Effect |
|----------|--------|
| `allow` | Permits the action |
| `deny` | Blocks the action (reason shown to user) |
| `advise` | Permits but shows warning |

### Matching

Hooks are triggered when BOTH conditions match:

| Field | Description |
|-------|-------------|
| `tools` | Tools that trigger this hook (required) |
| `paths` | File patterns that trigger this hook (optional) |

If `paths` is empty, the hook triggers for any path with a matching tool.

### Error Handling

| `on_error` Value | Behavior When Hook Fails |
|------------------|--------------------------|
| `allow` (default) | Permit the action, show warning |
| `deny` | Block the action |

### Example: Vendor Readonly

```bash
#!/bin/bash
# hooks/vendor-readonly.sh
echo '{"decision":"deny","reason":"vendor directory is readonly"}'
```

```yaml
hooks:
  - name: "vendor-readonly"
    command: "./hooks/vendor-readonly.sh"
    tools: ["Write", "Edit"]
    paths: ["vendor/**"]
```

### Example: Custom Linter

```bash
#!/bin/bash
# hooks/lint-check.sh
# Receives JSON on stdin, can parse with jq
input=$(cat)
file_path=$(echo "$input" | jq -r '.paths[0]')

if ! golint "$file_path" > /dev/null 2>&1; then
    echo '{"decision":"advise","warning":"Lint issues detected"}'
else
    echo '{"decision":"allow"}'
fi
```

### All Options Reference

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `name` | string | required | Unique identifier for the hook |
| `command` | string | required | Path to executable |
| `args` | []string | [] | Arguments to pass to command |
| `tools` | []string | required | Tools that trigger this hook |
| `paths` | []string | [] | Glob patterns (empty = all paths) |
| `timeout` | duration | 5s | Max execution time |
| `on_error` | string | allow | Behavior on failure: allow, deny |

---

## Invariants

**Status**: Implemented via Hooks

Structural rules are now implemented through the external hooks system. Define custom hooks that validate architectural constraints.

See [Hooks](#hooks-external-hooks) for implementation details.

---

## Patterns

**Status**: Planned (can be implemented via Hooks)

Ensures new code follows existing conventions and patterns.

Can be partially implemented using external hooks that analyze code against patterns.

---

## Boundaries

**Status**: Planned (can be implemented via Hooks)

Respects module boundaries and dependency rules.

Can be implemented using external hooks that analyze import graphs and dependencies.
