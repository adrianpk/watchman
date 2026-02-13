# Watchman Plugins

External hooks that extend Watchman's validation capabilities.

## Available Plugins

| Plugin | Description | Status |
|--------|-------------|--------|
| [sentinel](./sentinel/) | AI-powered standards evaluation | Active |

## Creating a Plugin

Plugins are executables that:

1. Receive JSON input via stdin (HookInput)
2. Return JSON output via stdout (HookOutput)

### Input Format

```json
{
  "tool_name": "Write",
  "tool_input": {
    "file_path": "/path/to/file.go",
    "content": "package main..."
  },
  "paths": ["file.go"],
  "working_dir": "/project/root"
}
```

### Output Format

```json
{
  "decision": "allow|advise|deny",
  "reason": "explanation for deny",
  "warning": "advisory message for advise"
}
```

### Integration

Add to `.watchman.yml`:

```yaml
hooks:
  - name: my-plugin
    command: /path/to/plugin
    tools: [Write, Edit, Bash]
    paths: ["src/**"]
    timeout: 10s
    on_error: allow
```

## Plugin Guidelines

- Respect timeout constraints
- Return valid JSON always
- Use `on_error: allow` for non-critical plugins
- Keep plugins focused and fast
