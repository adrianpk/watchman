# Invariants

Declarative structural checks using regex and glob patterns. Language-agnostic, no AST parsing.

## Overview

The invariants rule provides built-in declarative checks for common structural constraints:

- **Coexistence**: Ensure related files exist together (e.g., test requires implementation)
- **Content**: Validate file content matches or excludes patterns
- **Imports**: Restrict import statements using regex
- **Naming**: Enforce file naming conventions
- **Required**: Ensure certain files exist in directories

For AST-based or language-specific checks, use [Hooks](rules.md#hooks-external-hooks).

## Quick Start

```yaml
rules:
  invariants: true

invariants:
  content:
    - name: "no-todos"
      paths: ["**/*.go"]
      forbid: "TODO|FIXME"
```

## Check Types

### Coexistence

Ensures related files exist together.

```yaml
invariants:
  coexistence:
    - name: "test-requires-impl"
      if: "**/*_test.go"
      require: "${base}.go"
      message: "Test file requires corresponding implementation"
```

| Field | Description |
|-------|-------------|
| `if` | Glob pattern that triggers the check |
| `require` | File that must exist (supports placeholders) |
| `message` | Custom error message (optional) |

**Placeholders:**
- `${name}` - Filename without extension (`user.go` → `user`)
- `${ext}` - Extension with dot (`user.go` → `.go`)
- `${base}` - For `_test` files, strips the suffix (`user_test.go` → `user`)

### Content

Validates file content against patterns.

```yaml
invariants:
  content:
    # Forbid patterns
    - name: "no-todos"
      paths: ["**/*.go", "!**/*_test.go"]
      forbid: "TODO|FIXME|HACK"

    # Require patterns
    - name: "copyright-header"
      paths: ["**/*.go"]
      require: "^// Copyright"

    # Both in one check
    - name: "strict-mode"
      paths: ["**/*.ts"]
      require: "\"use strict\""
      forbid: "eval\\("
```

| Field | Description |
|-------|-------------|
| `paths` | Glob patterns (supports `!` for exclusion) |
| `require` | Regex pattern that must match |
| `forbid` | Regex pattern that must not match |
| `message` | Custom error message (optional) |

### Imports

Restricts import statements using regex (not AST).

```yaml
invariants:
  imports:
    - name: "no-internal-in-adapters"
      paths: ["adapters/**/*.go"]
      forbid: '".*internal/core"'

    - name: "no-test-imports"
      paths: ["**/*.go", "!**/*_test.go"]
      forbid: '"testing"'
```

| Field | Description |
|-------|-------------|
| `paths` | Files to check |
| `forbid` | Regex pattern for forbidden imports |
| `message` | Custom error message (optional) |

### Naming

Validates file naming conventions.

```yaml
invariants:
  naming:
    - name: "cmd-main-only"
      paths: ["cmd/**/*.go"]
      pattern: "main\\.go$"

    - name: "snake-case"
      paths: ["internal/**/*.go"]
      pattern: "^[a-z][a-z0-9_]*\\.go$"
```

| Field | Description |
|-------|-------------|
| `paths` | Directories/patterns to check |
| `pattern` | Regex pattern filenames must match |
| `message` | Custom error message (optional) |

### Required

Ensures certain files exist in directories.

```yaml
invariants:
  required:
    - name: "doc-required"
      dirs: "internal/**"
      when: "*.go"
      require: "doc.go"

    - name: "readme-in-packages"
      dirs: "pkg/*"
      require: "README.md"
```

| Field | Description |
|-------|-------------|
| `dirs` | Glob pattern for directories to check |
| `when` | Only check when this pattern exists (optional) |
| `require` | File that must exist |
| `message` | Custom error message (optional) |

## Path Patterns

All path patterns support:

- `**` - Recursive directory matching
- `*` - Single directory/file matching
- `!pattern` - Exclusion (prefix with `!`)

Examples:
- `**/*.go` - All Go files recursively
- `!**/*_test.go` - Exclude test files
- `src/**` - Everything under src/
- `cmd/*/main.go` - main.go in direct subdirectories of cmd/

## Invariants vs Hooks

| Invariants | Hooks |
|------------|-------|
| Built-in, no external dependencies | Requires external scripts |
| Regex/glob patterns only | Full language capabilities |
| Fast, synchronous | Subprocess overhead |
| Language-agnostic | Language-specific AST possible |

**Use Invariants when:**
- Simple pattern matching is sufficient
- No external tools needed
- Performance is critical

**Use Hooks when:**
- AST analysis required
- Complex validation logic
- Integration with existing tools

## Custom Hooks for Complex Checks

When invariants are not sufficient, create custom hooks. Below are examples for different stacks.

### Shell Hook (Any Language)

```bash
#!/bin/bash
# hooks/my-check.sh
INPUT=$(cat)
FILE_PATH=$(echo "$INPUT" | jq -r '.paths[0]')
CONTENT=$(echo "$INPUT" | jq -r '.tool_input.content // empty')

# Your logic here
if [[ "$CONTENT" == *"FORBIDDEN"* ]]; then
    echo '{"decision":"deny","reason":"Forbidden content found"}'
else
    echo '{"decision":"allow"}'
fi
```

### Python Hook (with AST)

```python
#!/usr/bin/env python3
# hooks/python-check.py
import json
import sys
import ast

input_data = json.load(sys.stdin)
file_path = input_data['paths'][0]
content = input_data['tool_input'].get('content', '')

if file_path.endswith('.py'):
    try:
        tree = ast.parse(content)
        for node in ast.walk(tree):
            if isinstance(node, ast.Import):
                for alias in node.names:
                    if alias.name.startswith('forbidden_'):
                        print(json.dumps({
                            "decision": "deny",
                            "reason": f"Forbidden import: {alias.name}"
                        }))
                        sys.exit(0)
    except SyntaxError:
        pass

print(json.dumps({"decision": "allow"}))
```

### Go Hook (with AST)

```go
// hooks/go-ast-check/main.go
package main

import (
    "encoding/json"
    "go/parser"
    "go/token"
    "os"
)

type Input struct {
    ToolName  string                 `json:"tool_name"`
    ToolInput map[string]interface{} `json:"tool_input"`
    Paths     []string               `json:"paths"`
}

type Output struct {
    Decision string `json:"decision"`
    Reason   string `json:"reason,omitempty"`
}

func main() {
    var input Input
    json.NewDecoder(os.Stdin).Decode(&input)

    content, _ := input.ToolInput["content"].(string)
    fset := token.NewFileSet()

    f, err := parser.ParseFile(fset, "", content, parser.ImportsOnly)
    if err != nil {
        json.NewEncoder(os.Stdout).Encode(Output{Decision: "allow"})
        return
    }

    for _, imp := range f.Imports {
        if imp.Path.Value == `"unsafe"` {
            json.NewEncoder(os.Stdout).Encode(Output{
                Decision: "deny",
                Reason:   "unsafe package not allowed",
            })
            return
        }
    }

    json.NewEncoder(os.Stdout).Encode(Output{Decision: "allow"})
}
```

Build and configure:
```bash
go build -o hooks/go-ast-check ./hooks/go-ast-check
```

### TypeScript/Node Hook

```javascript
#!/usr/bin/env node
// hooks/ts-check.js
const ts = require('typescript');

let input = '';
process.stdin.on('data', chunk => input += chunk);
process.stdin.on('end', () => {
    const data = JSON.parse(input);
    const filePath = data.paths[0];
    const content = data.tool_input?.content || '';

    if (!filePath.endsWith('.ts')) {
        console.log(JSON.stringify({ decision: 'allow' }));
        return;
    }

    const sourceFile = ts.createSourceFile(
        filePath,
        content,
        ts.ScriptTarget.Latest,
        true
    );

    // Check for 'any' type usage
    let hasAny = false;
    function visit(node) {
        if (node.kind === ts.SyntaxKind.AnyKeyword) {
            hasAny = true;
        }
        ts.forEachChild(node, visit);
    }
    visit(sourceFile);

    if (hasAny) {
        console.log(JSON.stringify({
            decision: 'deny',
            reason: 'Usage of "any" type is forbidden'
        }));
    } else {
        console.log(JSON.stringify({ decision: 'allow' }));
    }
});
```

### Hook Configuration

```yaml
hooks:
  # Python AST check
  - name: "python-imports"
    command: "./hooks/python-check.py"
    tools: ["Write", "Edit"]
    paths: ["**/*.py"]

  # Go AST check
  - name: "go-unsafe"
    command: "./hooks/go-ast-check"
    tools: ["Write", "Edit"]
    paths: ["**/*.go"]

  # TypeScript AST check
  - name: "ts-no-any"
    command: "./hooks/ts-check.js"
    tools: ["Write", "Edit"]
    paths: ["**/*.ts"]

  # Generic shell check
  - name: "no-secrets"
    command: "./hooks/secret-scan.sh"
    tools: ["Write", "Edit"]
```

## Examples

### Go Project

```yaml
rules:
  invariants: true

invariants:
  coexistence:
    - name: "test-requires-impl"
      if: "**/*_test.go"
      require: "${base}.go"

  content:
    - name: "no-todos"
      paths: ["**/*.go", "!**/*_test.go"]
      forbid: "TODO|FIXME"
    - name: "copyright"
      paths: ["**/*.go"]
      require: "^// Copyright"

  imports:
    - name: "no-internal-in-cmd"
      paths: ["cmd/**/*.go"]
      forbid: '".*internal/.*"'

  naming:
    - name: "snake-case"
      paths: ["internal/**/*.go"]
      pattern: "^[a-z][a-z0-9_]*\\.go$"
```

### TypeScript Project

```yaml
rules:
  invariants: true

invariants:
  content:
    - name: "strict-mode"
      paths: ["**/*.ts", "!**/*.test.ts"]
      forbid: "eval\\(|Function\\("
    - name: "no-console"
      paths: ["src/**/*.ts", "!src/**/*.test.ts"]
      forbid: "console\\.(log|debug|info)"

  naming:
    - name: "component-pascal-case"
      paths: ["src/components/**/*.tsx"]
      pattern: "^[A-Z][a-zA-Z0-9]*\\.tsx$"
```

### Python Project

```yaml
rules:
  invariants: true

invariants:
  content:
    - name: "no-print-debug"
      paths: ["**/*.py", "!tests/**"]
      forbid: "print\\(.*debug"
    - name: "docstring-required"
      paths: ["src/**/*.py"]
      require: '"""'

  imports:
    - name: "no-wildcard"
      paths: ["**/*.py"]
      forbid: "from .* import \\*"

  required:
    - name: "init-required"
      dirs: "src/**"
      when: "*.py"
      require: "__init__.py"
```
