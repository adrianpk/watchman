package hook

import (
	"testing"

	"github.com/adrianpk/watchman/internal/config"
)

func TestNewHookMatcher(t *testing.T) {
	m := NewHookMatcher()
	if m == nil {
		t.Error("NewHookMatcher returned nil")
	}
}

func TestHookMatcherMatchesTool(t *testing.T) {
	m := NewHookMatcher()

	tests := []struct {
		name     string
		tools    []string
		toolName string
		want     bool
	}{
		{"exact match", []string{"Write"}, "Write", true},
		{"case insensitive", []string{"write"}, "Write", true},
		{"case insensitive reverse", []string{"Write"}, "write", true},
		{"multiple tools first", []string{"Read", "Write"}, "Read", true},
		{"multiple tools second", []string{"Read", "Write"}, "Write", true},
		{"no match", []string{"Read"}, "Write", false},
		{"empty tools", []string{}, "Write", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := m.matchesTool(tt.tools, tt.toolName)
			if got != tt.want {
				t.Errorf("matchesTool(%v, %q) = %v, want %v", tt.tools, tt.toolName, got, tt.want)
			}
		})
	}
}

func TestHookMatcherMatchesAnyPath(t *testing.T) {
	m := NewHookMatcher()

	tests := []struct {
		name     string
		patterns []string
		paths    []string
		want     bool
	}{
		{"exact path", []string{"src/main.go"}, []string{"src/main.go"}, true},
		{"glob star", []string{"*.go"}, []string{"main.go"}, true},
		{"double star", []string{"**/*.go"}, []string{"src/pkg/main.go"}, true},
		{"multiple patterns first", []string{"*.go", "*.js"}, []string{"main.go"}, true},
		{"multiple patterns second", []string{"*.go", "*.js"}, []string{"app.js"}, true},
		{"multiple paths first", []string{"*.go"}, []string{"main.go", "style.css"}, true},
		{"multiple paths second", []string{"*.go"}, []string{"style.css", "main.go"}, true},
		{"no match", []string{"*.js"}, []string{"main.go"}, false},
		{"empty paths", []string{"*.go"}, []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := m.matchesAnyPath(tt.patterns, tt.paths)
			if got != tt.want {
				t.Errorf("matchesAnyPath(%v, %v) = %v, want %v", tt.patterns, tt.paths, got, tt.want)
			}
		})
	}
}

func TestHookMatcherMatches(t *testing.T) {
	m := NewHookMatcher()

	tests := []struct {
		name     string
		hook     *config.HookConfig
		toolName string
		paths    []string
		want     bool
	}{
		{
			name:     "tool and path match",
			hook:     &config.HookConfig{Tools: []string{"Write"}, Paths: []string{"**/*.go"}},
			toolName: "Write",
			paths:    []string{"src/main.go"},
			want:     true,
		},
		{
			name:     "tool matches no path patterns",
			hook:     &config.HookConfig{Tools: []string{"Write"}},
			toolName: "Write",
			paths:    []string{"anything.txt"},
			want:     true,
		},
		{
			name:     "tool no match",
			hook:     &config.HookConfig{Tools: []string{"Write"}, Paths: []string{"**/*.go"}},
			toolName: "Read",
			paths:    []string{"src/main.go"},
			want:     false,
		},
		{
			name:     "path no match",
			hook:     &config.HookConfig{Tools: []string{"Write"}, Paths: []string{"**/*.go"}},
			toolName: "Write",
			paths:    []string{"style.css"},
			want:     false,
		},
		{
			name:     "empty tools never matches",
			hook:     &config.HookConfig{Tools: []string{}, Paths: []string{"**"}},
			toolName: "Write",
			paths:    []string{"anything"},
			want:     false,
		},
		{
			name:     "multiple tools multiple paths",
			hook:     &config.HookConfig{Tools: []string{"Write", "Edit"}, Paths: []string{"src/**", "internal/**"}},
			toolName: "Edit",
			paths:    []string{"internal/hook/eval.go"},
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := m.Matches(tt.hook, tt.toolName, tt.paths, "")
			if got != tt.want {
				t.Errorf("Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}
