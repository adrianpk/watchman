package policy

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/adrianpk/watchman/internal/config"
	"github.com/adrianpk/watchman/internal/parser"
)

func TestNewScopeToFiles(t *testing.T) {
	tests := []struct {
		name string
		cfg  *config.ScopeConfig
		want *ScopeToFiles
	}{
		{
			name: "nil config",
			cfg:  nil,
			want: &ScopeToFiles{},
		},
		{
			name: "with allow and block",
			cfg: &config.ScopeConfig{
				Allow: []string{"src/**/*.go"},
				Block: []string{"vendor/**"},
			},
			want: &ScopeToFiles{
				Allow: []string{"src/**/*.go"},
				Block: []string{"vendor/**"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewScopeToFiles(tt.cfg)
			if tt.cfg == nil {
				if got.Allow != nil || got.Block != nil {
					t.Errorf("expected empty rule for nil config")
				}
				return
			}
			if len(got.Allow) != len(tt.want.Allow) {
				t.Errorf("Allow = %v, want %v", got.Allow, tt.want.Allow)
			}
			if len(got.Block) != len(tt.want.Block) {
				t.Errorf("Block = %v, want %v", got.Block, tt.want.Block)
			}
		})
	}
}

func TestScopeToFilesEvaluate(t *testing.T) {
	tests := []struct {
		name        string
		rule        *ScopeToFiles
		toolName    string
		cmd         parser.Command
		wantAllowed bool
	}{
		{
			name:        "read tool always allowed",
			rule:        &ScopeToFiles{Allow: []string{"src/**"}},
			toolName:    "Read",
			cmd:         parser.Command{Args: []string{"/etc/passwd"}},
			wantAllowed: true,
		},
		{
			name:        "grep tool always allowed",
			rule:        &ScopeToFiles{Block: []string{"**"}},
			toolName:    "Grep",
			cmd:         parser.Command{Args: []string{"vendor/lib.go"}},
			wantAllowed: true,
		},
		{
			name:        "write tool no scope allows all",
			rule:        &ScopeToFiles{},
			toolName:    "Write",
			cmd:         parser.Command{Args: []string{"any/file.go"}},
			wantAllowed: true,
		},
		{
			name:        "write tool in scope allowed",
			rule:        &ScopeToFiles{Allow: []string{"src/**/*.go"}},
			toolName:    "Write",
			cmd:         parser.Command{Args: []string{"src/main.go"}},
			wantAllowed: true,
		},
		{
			name:        "write tool out of scope blocked",
			rule:        &ScopeToFiles{Allow: []string{"src/**/*.go"}},
			toolName:    "Write",
			cmd:         parser.Command{Args: []string{"vendor/lib.go"}},
			wantAllowed: false,
		},
		{
			name:        "edit tool blocked by block list",
			rule:        &ScopeToFiles{Block: []string{"vendor/**"}},
			toolName:    "Edit",
			cmd:         parser.Command{Args: []string{"vendor/lib.go"}},
			wantAllowed: false,
		},
		{
			name:        "block takes precedence over allow",
			rule:        &ScopeToFiles{Allow: []string{"**/*.go"}, Block: []string{"vendor/**"}},
			toolName:    "Edit",
			cmd:         parser.Command{Args: []string{"vendor/lib.go"}},
			wantAllowed: false,
		},
		{
			name:        "notebook edit blocked",
			rule:        &ScopeToFiles{Allow: []string{"notebooks/*.ipynb"}},
			toolName:    "NotebookEdit",
			cmd:         parser.Command{Args: []string{"analysis/test.ipynb"}},
			wantAllowed: false,
		},
		{
			name:        "notebook edit allowed",
			rule:        &ScopeToFiles{Allow: []string{"notebooks/*.ipynb"}},
			toolName:    "NotebookEdit",
			cmd:         parser.Command{Args: []string{"notebooks/test.ipynb"}},
			wantAllowed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rule.Evaluate(tt.toolName, tt.cmd)
			if got.Allowed != tt.wantAllowed {
				t.Errorf("Evaluate() = %v, want %v, reason: %s", got.Allowed, tt.wantAllowed, got.Reason)
			}
		})
	}
}

func TestScopeIsBlocked(t *testing.T) {
	rule := &ScopeToFiles{
		Block: []string{"vendor/**", "**/*_generated.go", ".env"},
	}

	tests := []struct {
		path    string
		blocked bool
	}{
		{"vendor/lib/file.go", true},
		{"vendor/other.go", true},
		{"src/types_generated.go", true},
		{"internal/api_generated.go", true},
		{".env", true},
		{"src/main.go", false},
		{"internal/api.go", false},
		{"README.md", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := rule.isBlocked(tt.path)
			if got != tt.blocked {
				t.Errorf("isBlocked(%q) = %v, want %v", tt.path, got, tt.blocked)
			}
		})
	}
}

func TestScopeIsInScope(t *testing.T) {
	tests := []struct {
		name    string
		rule    *ScopeToFiles
		path    string
		inScope bool
	}{
		{
			name:    "empty allow list allows all",
			rule:    &ScopeToFiles{},
			path:    "any/path/file.go",
			inScope: true,
		},
		{
			name:    "path matches allow pattern",
			rule:    &ScopeToFiles{Allow: []string{"src/**/*.go"}},
			path:    "src/pkg/file.go",
			inScope: true,
		},
		{
			name:    "path does not match allow pattern",
			rule:    &ScopeToFiles{Allow: []string{"src/**/*.go"}},
			path:    "vendor/lib.go",
			inScope: false,
		},
		{
			name:    "multiple allow patterns",
			rule:    &ScopeToFiles{Allow: []string{"src/**", "internal/**"}},
			path:    "internal/pkg/file.go",
			inScope: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rule.isInScope(tt.path)
			if got != tt.inScope {
				t.Errorf("isInScope(%q) = %v, want %v", tt.path, got, tt.inScope)
			}
		})
	}
}

func TestScopeAbsolutePathNormalization(t *testing.T) {
	// This test verifies that absolute paths within cwd are normalized
	// to relative paths for glob matching
	rule := &ScopeToFiles{Allow: []string{"assets/**/*.html", "src/**/*.go"}}

	cwd, err := os.Getwd()
	if err != nil {
		t.Skip("cannot get cwd")
	}

	tests := []struct {
		name    string
		path    string
		inScope bool
	}{
		// Relative paths should work as before
		{"relative html", "assets/templates/page.html", true},
		{"relative go", "src/main.go", true},
		{"relative out of scope", "vendor/lib.go", false},

		// Absolute paths within cwd should be normalized and match
		{"absolute html", filepath.Join(cwd, "assets/templates/page.html"), true},
		{"absolute go", filepath.Join(cwd, "src/main.go"), true},
		{"absolute out of scope", filepath.Join(cwd, "vendor/lib.go"), false},

		// Absolute paths outside cwd should not match relative patterns
		{"absolute outside cwd", "/tmp/assets/templates/page.html", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rule.isInScope(tt.path)
			if got != tt.inScope {
				t.Errorf("isInScope(%q) = %v, want %v", tt.path, got, tt.inScope)
			}
		})
	}
}

func TestToRelativePath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string // empty means should return same path
	}{
		{"relative stays relative", "src/main.go", "src/main.go"},
		{"dot relative", "./src/main.go", "./src/main.go"},
		{"outside cwd stays absolute", "/etc/passwd", "/etc/passwd"},
		{"tmp stays absolute", "/tmp/file.txt", "/tmp/file.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toRelativePath(tt.path)
			want := tt.want
			if want == "" {
				want = tt.path
			}
			if got != want {
				t.Errorf("toRelativePath(%q) = %q, want %q", tt.path, got, want)
			}
		})
	}
}

// Note: matchGlob and matchDoublestar tests are now in internal/glob/glob_test.go
