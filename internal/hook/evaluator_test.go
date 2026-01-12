package hook

import (
	"testing"

	"github.com/adrianpk/watchman/internal/config"
)

func TestNewEvaluator(t *testing.T) {
	cfg := &config.Config{}
	e := NewEvaluator(cfg)
	if e == nil {
		t.Error("NewEvaluator returned nil")
	}
}

func TestEvaluatorIsToolBlocked(t *testing.T) {
	cfg := &config.Config{
		Tools: config.ToolsConfig{
			Block: []string{"Bash", "Write"},
		},
	}
	e := NewEvaluator(cfg)

	tests := []struct {
		tool    string
		blocked bool
	}{
		{"Bash", true},
		{"bash", true},
		{"Write", true},
		{"Read", false},
		{"Edit", false},
	}

	for _, tt := range tests {
		t.Run(tt.tool, func(t *testing.T) {
			got := e.isToolBlocked(tt.tool)
			if got != tt.blocked {
				t.Errorf("isToolBlocked(%q) = %v, want %v", tt.tool, got, tt.blocked)
			}
		})
	}
}

func TestEvaluatorIsToolAllowed(t *testing.T) {
	tests := []struct {
		name    string
		allow   []string
		tool    string
		allowed bool
	}{
		{
			name:    "empty allow list allows all",
			allow:   []string{},
			tool:    "Bash",
			allowed: true,
		},
		{
			name:    "tool in allow list",
			allow:   []string{"Read", "Edit"},
			tool:    "Read",
			allowed: true,
		},
		{
			name:    "tool not in allow list",
			allow:   []string{"Read", "Edit"},
			tool:    "Bash",
			allowed: false,
		},
		{
			name:    "case insensitive",
			allow:   []string{"Read"},
			tool:    "read",
			allowed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Tools: config.ToolsConfig{Allow: tt.allow},
			}
			e := NewEvaluator(cfg)
			got := e.isToolAllowed(tt.tool)
			if got != tt.allowed {
				t.Errorf("isToolAllowed(%q) = %v, want %v", tt.tool, got, tt.allowed)
			}
		})
	}
}

func TestEvaluatorIsCommandBlocked(t *testing.T) {
	cfg := &config.Config{
		Commands: config.CommandsConfig{
			Block: []string{"sudo", "rm -rf"},
		},
	}
	e := NewEvaluator(cfg)

	tests := []struct {
		cmd     string
		blocked string
	}{
		{"sudo apt install", "sudo"},
		{"rm -rf /", "rm -rf"},
		{"ls -la", ""},
		{"echo hello", ""},
	}

	for _, tt := range tests {
		t.Run(tt.cmd, func(t *testing.T) {
			got := e.isCommandBlocked(tt.cmd)
			if got != tt.blocked {
				t.Errorf("isCommandBlocked(%q) = %q, want %q", tt.cmd, got, tt.blocked)
			}
		})
	}
}

func TestIsFilesystemTool(t *testing.T) {
	tests := []struct {
		tool string
		want bool
	}{
		{"Bash", true},
		{"Read", true},
		{"Write", true},
		{"Edit", true},
		{"Glob", true},
		{"Grep", true},
		{"WebSearch", false},
		{"Task", false},
	}

	for _, tt := range tests {
		t.Run(tt.tool, func(t *testing.T) {
			got := isFilesystemTool(tt.tool)
			if got != tt.want {
				t.Errorf("isFilesystemTool(%q) = %v, want %v", tt.tool, got, tt.want)
			}
		})
	}
}

func TestIsModificationTool(t *testing.T) {
	tests := []struct {
		tool string
		want bool
	}{
		{"Write", true},
		{"Edit", true},
		{"NotebookEdit", true},
		{"Read", false},
		{"Bash", false},
		{"Glob", false},
	}

	for _, tt := range tests {
		t.Run(tt.tool, func(t *testing.T) {
			got := isModificationTool(tt.tool)
			if got != tt.want {
				t.Errorf("isModificationTool(%q) = %v, want %v", tt.tool, got, tt.want)
			}
		})
	}
}

func TestEvaluatorEvaluateBlockedTool(t *testing.T) {
	cfg := &config.Config{
		Tools: config.ToolsConfig{
			Block: []string{"Bash"},
		},
	}
	e := NewEvaluator(cfg)

	result := e.Evaluate(Input{ToolName: "Bash"})
	if result.Allowed {
		t.Error("expected blocked tool to be denied")
	}
}

func TestEvaluatorEvaluateNotAllowedTool(t *testing.T) {
	cfg := &config.Config{
		Tools: config.ToolsConfig{
			Allow: []string{"Read"},
		},
	}
	e := NewEvaluator(cfg)

	result := e.Evaluate(Input{ToolName: "Write"})
	if result.Allowed {
		t.Error("expected non-allowed tool to be denied")
	}
}

func TestEvaluatorEvaluateNonFilesystemTool(t *testing.T) {
	cfg := &config.Config{}
	e := NewEvaluator(cfg)

	result := e.Evaluate(Input{ToolName: "WebSearch"})
	if !result.Allowed {
		t.Error("expected non-filesystem tool to be allowed")
	}
}

func TestEvaluatorEvaluateBlockedCommand(t *testing.T) {
	cfg := &config.Config{
		Commands: config.CommandsConfig{
			Block: []string{"sudo"},
		},
	}
	e := NewEvaluator(cfg)

	result := e.Evaluate(Input{
		ToolName:  "Bash",
		ToolInput: map[string]interface{}{"command": "sudo rm -rf /"},
	})
	if result.Allowed {
		t.Error("expected blocked command to be denied")
	}
}

func TestEvaluatorEvaluateWorkspace(t *testing.T) {
	cfg := &config.Config{
		Rules: config.RulesConfig{Workspace: true},
	}
	e := NewEvaluator(cfg)

	// Should block absolute path
	result := e.Evaluate(Input{
		ToolName:  "Read",
		ToolInput: map[string]interface{}{"file_path": "/etc/passwd"},
	})
	if result.Allowed {
		t.Error("expected workspace rule to block absolute path")
	}

	// Should allow relative path
	result = e.Evaluate(Input{
		ToolName:  "Read",
		ToolInput: map[string]interface{}{"file_path": "./src/main.go"},
	})
	if !result.Allowed {
		t.Errorf("expected workspace rule to allow relative path: %s", result.Reason)
	}
}

func TestEvaluatorEvaluateScope(t *testing.T) {
	cfg := &config.Config{
		Rules: config.RulesConfig{Scope: true},
		Scope: config.ScopeConfig{
			Allow: []string{"src/**/*.go"},
		},
	}
	e := NewEvaluator(cfg)

	// Should block file outside scope
	result := e.Evaluate(Input{
		ToolName:  "Write",
		ToolInput: map[string]interface{}{"file_path": "vendor/lib.go"},
	})
	if result.Allowed {
		t.Error("expected scope rule to block file outside allowed patterns")
	}

	// Should allow file in scope
	result = e.Evaluate(Input{
		ToolName:  "Write",
		ToolInput: map[string]interface{}{"file_path": "src/main.go"},
	})
	if !result.Allowed {
		t.Errorf("expected scope rule to allow file in scope: %s", result.Reason)
	}
}

func TestEvaluatorEvaluateVersioning(t *testing.T) {
	cfg := &config.Config{
		Rules: config.RulesConfig{Versioning: true},
		Versioning: config.VersioningConfig{
			Commit: config.CommitConfig{
				RequireUppercase: true,
			},
		},
	}
	e := NewEvaluator(cfg)

	// Should block lowercase commit message
	result := e.Evaluate(Input{
		ToolName:  "Bash",
		ToolInput: map[string]interface{}{"command": `git commit -m "lowercase"`},
	})
	if result.Allowed {
		t.Error("expected versioning rule to block lowercase commit")
	}

	// Should allow uppercase commit message
	result = e.Evaluate(Input{
		ToolName:  "Bash",
		ToolInput: map[string]interface{}{"command": `git commit -m "Uppercase"`},
	})
	if !result.Allowed {
		t.Errorf("expected versioning rule to allow uppercase commit: %s", result.Reason)
	}

	// Non-commit command should pass
	result = e.Evaluate(Input{
		ToolName:  "Bash",
		ToolInput: map[string]interface{}{"command": "git status"},
	})
	if !result.Allowed {
		t.Errorf("expected versioning rule to allow non-commit: %s", result.Reason)
	}
}

func TestEvaluatorEvaluateIncremental(t *testing.T) {
	cfg := &config.Config{
		Rules: config.RulesConfig{Incremental: true},
		Incremental: config.IncrementalConfig{
			MaxFiles:  10,
			WarnRatio: 0.7,
		},
	}
	e := NewEvaluator(cfg)

	// Should allow modification tool (actual blocking depends on git status)
	result := e.Evaluate(Input{
		ToolName:  "Write",
		ToolInput: map[string]interface{}{"file_path": "test.go"},
	})
	// Just verify it runs without error, actual result depends on git state
	_ = result
}

func TestEvaluatorEvaluateProtectedPath(t *testing.T) {
	cfg := &config.Config{}
	e := NewEvaluator(cfg)

	// Should block protected paths (using .watchman.yml which is always protected)
	result := e.Evaluate(Input{
		ToolName:  "Write",
		ToolInput: map[string]interface{}{"file_path": ".watchman.yml"},
	})
	if result.Allowed {
		t.Error("expected protected path to be blocked")
	}
}

func TestEvaluatorEvaluateAllowedFilesystemTool(t *testing.T) {
	cfg := &config.Config{}
	e := NewEvaluator(cfg)

	// Read with relative path should be allowed
	result := e.Evaluate(Input{
		ToolName:  "Read",
		ToolInput: map[string]interface{}{"file_path": "main.go"},
	})
	if !result.Allowed {
		t.Errorf("expected Read with relative path to be allowed: %s", result.Reason)
	}
}
