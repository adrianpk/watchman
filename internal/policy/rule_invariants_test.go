package policy

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/adrianpk/watchman/internal/config"
)

func TestInvariantsRuleNonModificationTool(t *testing.T) {
	cfg := &config.InvariantsConfig{
		Content: []config.ContentCheck{
			{
				Name:   "forbid-test",
				Paths:  []string{"**/*.go"},
				Forbid: "FORBIDDEN",
			},
		},
	}
	rule := NewInvariantsRule(cfg)

	// Read tool should always pass
	decision := rule.Evaluate("Read", "test.go", "FORBIDDEN content")
	if !decision.Allowed {
		t.Error("expected Read tool to be allowed regardless of content")
	}
}

func TestInvariantsContentForbid(t *testing.T) {
	cfg := &config.InvariantsConfig{
		Content: []config.ContentCheck{
			{
				Name:   "no-todos",
				Paths:  []string{"**/*.go"},
				Forbid: "TODO|FIXME",
			},
		},
	}
	rule := NewInvariantsRule(cfg)

	tests := []struct {
		name    string
		path    string
		content string
		allowed bool
	}{
		{"clean content", "src/main.go", "package main", true},
		{"has TODO", "src/main.go", "// TODO: fix this", false},
		{"has FIXME", "src/main.go", "// FIXME: broken", false},
		{"non-matching path", "src/main.js", "// TODO: fix", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decision := rule.Evaluate("Write", tt.path, tt.content)
			if decision.Allowed != tt.allowed {
				t.Errorf("Evaluate(%q, %q) = %v, want %v: %s",
					tt.path, tt.content, decision.Allowed, tt.allowed, decision.Reason)
			}
		})
	}
}

func TestInvariantsContentRequire(t *testing.T) {
	cfg := &config.InvariantsConfig{
		Content: []config.ContentCheck{
			{
				Name:    "copyright",
				Paths:   []string{"**/*.go"},
				Require: "^// Copyright",
			},
		},
	}
	rule := NewInvariantsRule(cfg)

	tests := []struct {
		name    string
		content string
		allowed bool
	}{
		{"has copyright", "// Copyright 2024\npackage main", true},
		{"no copyright", "package main\n\nfunc main() {}", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decision := rule.Evaluate("Write", "src/main.go", tt.content)
			if decision.Allowed != tt.allowed {
				t.Errorf("Evaluate() = %v, want %v: %s",
					decision.Allowed, tt.allowed, decision.Reason)
			}
		})
	}
}

func TestInvariantsPathExclusion(t *testing.T) {
	cfg := &config.InvariantsConfig{
		Content: []config.ContentCheck{
			{
				Name:   "no-todos-except-tests",
				Paths:  []string{"**/*.go", "!**/*_test.go"},
				Forbid: "TODO",
			},
		},
	}
	rule := NewInvariantsRule(cfg)

	tests := []struct {
		name    string
		path    string
		content string
		allowed bool
	}{
		{"regular file with TODO", "src/main.go", "// TODO", false},
		{"test file with TODO", "src/main_test.go", "// TODO", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decision := rule.Evaluate("Write", tt.path, tt.content)
			if decision.Allowed != tt.allowed {
				t.Errorf("Evaluate(%q) = %v, want %v: %s",
					tt.path, decision.Allowed, tt.allowed, decision.Reason)
			}
		})
	}
}

func TestInvariantsImports(t *testing.T) {
	cfg := &config.InvariantsConfig{
		Imports: []config.ImportCheck{
			{
				Name:   "no-unsafe",
				Paths:  []string{"**/*.go"},
				Forbid: `"unsafe"`,
			},
		},
	}
	rule := NewInvariantsRule(cfg)

	tests := []struct {
		name    string
		content string
		allowed bool
	}{
		{"no unsafe import", `import "fmt"`, true},
		{"has unsafe import", `import "unsafe"`, false},
		{"unsafe in multi-import", "import (\n\t\"fmt\"\n\t\"unsafe\"\n)", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decision := rule.Evaluate("Write", "src/main.go", tt.content)
			if decision.Allowed != tt.allowed {
				t.Errorf("Evaluate() = %v, want %v: %s",
					decision.Allowed, tt.allowed, decision.Reason)
			}
		})
	}
}

func TestInvariantsNaming(t *testing.T) {
	cfg := &config.InvariantsConfig{
		Naming: []config.NamingCheck{
			{
				Name:    "snake-case",
				Paths:   []string{"internal/**/*.go"},
				Pattern: "^[a-z][a-z0-9_]*\\.go$",
			},
		},
	}
	rule := NewInvariantsRule(cfg)

	tests := []struct {
		name    string
		path    string
		allowed bool
	}{
		{"valid snake_case", "internal/pkg/my_file.go", true},
		{"invalid CamelCase", "internal/pkg/MyFile.go", false},
		{"outside internal", "src/MyFile.go", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decision := rule.Evaluate("Write", tt.path, "content")
			if decision.Allowed != tt.allowed {
				t.Errorf("Evaluate(%q) = %v, want %v: %s",
					tt.path, decision.Allowed, tt.allowed, decision.Reason)
			}
		})
	}
}

func TestInvariantsCoexistence(t *testing.T) {
	// Create temp directory with test files
	tmpDir := t.TempDir()
	implFile := filepath.Join(tmpDir, "user.go")
	testFile := filepath.Join(tmpDir, "user_test.go")

	if err := os.WriteFile(implFile, []byte("package test"), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := &config.InvariantsConfig{
		Coexistence: []config.CoexistenceCheck{
			{
				Name:    "test-requires-impl",
				If:      "**/*_test.go",
				Require: "${base}.go",
			},
		},
	}
	rule := NewInvariantsRule(cfg)

	// Test file where impl exists
	decision := rule.Evaluate("Write", testFile, "package test")
	if !decision.Allowed {
		t.Errorf("expected test file to be allowed when impl exists: %s", decision.Reason)
	}

	// Test file where impl doesn't exist
	otherTestFile := filepath.Join(tmpDir, "other_test.go")
	decision = rule.Evaluate("Write", otherTestFile, "package test")
	if decision.Allowed {
		t.Error("expected test file to be blocked when impl doesn't exist")
	}
}

func TestInvariantsRequired(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()
	pkgDir := filepath.Join(tmpDir, "pkg")
	if err := os.MkdirAll(pkgDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a .go file to trigger "when" condition
	goFile := filepath.Join(pkgDir, "main.go")
	if err := os.WriteFile(goFile, []byte("package pkg"), 0644); err != nil {
		t.Fatal(err)
	}

	// Use glob pattern that matches the temp directory
	cfg := &config.InvariantsConfig{
		Required: []config.RequiredCheck{
			{
				Name:    "doc-required",
				Dirs:    tmpDir + "/**",
				When:    "*.go",
				Require: "doc.go",
			},
		},
	}
	rule := NewInvariantsRule(cfg)

	// Writing to pkg without doc.go should fail
	newFile := filepath.Join(pkgDir, "other.go")
	decision := rule.Evaluate("Write", newFile, "package pkg")
	if decision.Allowed {
		t.Error("expected to block when doc.go is missing")
	}

	// Create doc.go
	docFile := filepath.Join(pkgDir, "doc.go")
	if err := os.WriteFile(docFile, []byte("// Package pkg docs"), 0644); err != nil {
		t.Fatal(err)
	}

	// Now should pass
	decision = rule.Evaluate("Write", newFile, "package pkg")
	if !decision.Allowed {
		t.Errorf("expected to allow when doc.go exists: %s", decision.Reason)
	}
}

func TestExpandPlaceholders(t *testing.T) {
	tests := []struct {
		pattern  string
		filePath string
		expected string
	}{
		{"${name}.go", "/tmp/user.go", "/tmp/user.go"},
		{"${base}.go", "/tmp/user_test.go", "/tmp/user.go"},
		{"${ext}", "/tmp/user.go", ".go"},
		{"${name}${ext}", "/tmp/main.go", "/tmp/main.go"},
	}

	for _, tt := range tests {
		t.Run(tt.pattern, func(t *testing.T) {
			result := expandPlaceholders(tt.pattern, tt.filePath)
			if result != tt.expected {
				t.Errorf("expandPlaceholders(%q, %q) = %q, want %q",
					tt.pattern, tt.filePath, result, tt.expected)
			}
		})
	}
}

func TestMatchesPathPatterns(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		patterns []string
		expected bool
	}{
		{"empty patterns", "src/main.go", []string{}, true},
		{"single match", "src/main.go", []string{"**/*.go"}, true},
		{"no match", "src/main.js", []string{"**/*.go"}, false},
		{"exclusion", "src/main_test.go", []string{"**/*.go", "!**/*_test.go"}, false},
		{"included not excluded", "src/main.go", []string{"**/*.go", "!**/*_test.go"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchesPathPatterns(tt.path, tt.patterns)
			if result != tt.expected {
				t.Errorf("matchesPathPatterns(%q, %v) = %v, want %v",
					tt.path, tt.patterns, result, tt.expected)
			}
		})
	}
}

func TestInvariantsNilConfig(t *testing.T) {
	rule := NewInvariantsRule(nil)
	decision := rule.Evaluate("Write", "test.go", "any content")
	if !decision.Allowed {
		t.Error("expected nil config to allow all")
	}
}

func TestInvariantsCustomMessage(t *testing.T) {
	cfg := &config.InvariantsConfig{
		Content: []config.ContentCheck{
			{
				Name:    "custom-msg",
				Paths:   []string{"**/*.go"},
				Forbid:  "FORBIDDEN",
				Message: "Custom error: forbidden pattern found",
			},
		},
	}
	rule := NewInvariantsRule(cfg)

	decision := rule.Evaluate("Write", "test.go", "FORBIDDEN")
	if decision.Allowed {
		t.Error("expected to be blocked")
	}
	if decision.Reason != "Custom error: forbidden pattern found" {
		t.Errorf("unexpected reason: %s", decision.Reason)
	}
}

func TestInvariantsNoFailedPrefix(t *testing.T) {
	cfg := &config.InvariantsConfig{
		Content: []config.ContentCheck{
			{
				Name:    "no-failed-prefix",
				Paths:   []string{"**/*.go"},
				Forbid:  `(fmt\.Errorf|errors\.New)\s*\(\s*"[Ff]ailed`,
				Message: "Use 'Cannot' for limitations or 'Error:' for failures, not 'Failed to'",
			},
		},
	}
	rule := NewInvariantsRule(cfg)

	tests := []struct {
		name    string
		content string
		allowed bool
	}{
		{"fmt.Errorf with Failed", `fmt.Errorf("Failed to open file")`, false},
		{"errors.New with failed", `errors.New("failed to connect")`, false},
		{"fmt.Errorf with Cannot", `fmt.Errorf("Cannot open file: %w", err)`, true},
		{"errors.New with Error", `errors.New("Error: connection refused")`, true},
		{"failed in message not prefix", `fmt.Errorf("operation failed")`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decision := rule.Evaluate("Write", "src/main.go", tt.content)
			if decision.Allowed != tt.allowed {
				t.Errorf("Evaluate(%q) = %v, want %v: %s",
					tt.content, decision.Allowed, tt.allowed, decision.Reason)
			}
		})
	}
}

func TestInvariantsErrorLowercase(t *testing.T) {
	cfg := &config.InvariantsConfig{
		Content: []config.ContentCheck{
			{
				Name:    "error-lowercase",
				Paths:   []string{"**/*.go"},
				Forbid:  `(fmt\.Errorf|errors\.New)\s*\(\s*"[A-Z]`,
				Message: "Error messages should not be capitalized (Go style)",
			},
		},
	}
	rule := NewInvariantsRule(cfg)

	tests := []struct {
		name    string
		content string
		allowed bool
	}{
		{"uppercase start", `errors.New("Invalid input")`, false},
		{"lowercase start", `errors.New("invalid input")`, true},
		{"fmt.Errorf uppercase", `fmt.Errorf("Connection refused: %w", err)`, false},
		{"fmt.Errorf lowercase", `fmt.Errorf("connection refused: %w", err)`, true},
		{"format verb start", `fmt.Errorf("%s: connection failed", host)`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decision := rule.Evaluate("Write", "src/main.go", tt.content)
			if decision.Allowed != tt.allowed {
				t.Errorf("Evaluate(%q) = %v, want %v: %s",
					tt.content, decision.Allowed, tt.allowed, decision.Reason)
			}
		})
	}
}
