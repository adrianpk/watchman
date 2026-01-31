package policy

import (
	"os"
	"testing"

	"github.com/adrianpk/watchman/internal/config"
	"github.com/adrianpk/watchman/internal/parser"
)

func TestConfineToWorkspaceEvaluate(t *testing.T) {
	tests := []struct {
		name        string
		cmd         string
		wantAllowed bool
	}{
		// Allowed cases
		{
			name:        "go test relative",
			cmd:         "go test ./...",
			wantAllowed: true,
		},
		{
			name:        "go test package",
			cmd:         "go test ./pkg/...",
			wantAllowed: true,
		},
		{
			name:        "make test",
			cmd:         "make test",
			wantAllowed: true,
		},
		{
			name:        "go build current dir",
			cmd:         "go build .",
			wantAllowed: true,
		},
		{
			name:        "simple command no args",
			cmd:         "ls",
			wantAllowed: true,
		},
		{
			name:        "flags without values",
			cmd:         "go test -race -v ./...",
			wantAllowed: true,
		},
		{
			name:        "empty command",
			cmd:         "",
			wantAllowed: true,
		},

		// Blocked cases - absolute paths
		{
			name:        "rm absolute path",
			cmd:         "rm -rf /",
			wantAllowed: false,
		},
		{
			name:        "cat absolute path",
			cmd:         "cat /etc/passwd",
			wantAllowed: false,
		},
		{
			name:        "cp to absolute",
			cmd:         "cp file.txt /tmp/file.txt",
			wantAllowed: false,
		},
		{
			name:        "flag with absolute path",
			cmd:         "go test -coverprofile=/tmp/cover.out ./...",
			wantAllowed: false,
		},

		// Blocked cases - traversal
		{
			name:        "parent dir only",
			cmd:         "cat ..",
			wantAllowed: false,
		},
		{
			name:        "parent dir traversal",
			cmd:         "cat ../secrets",
			wantAllowed: false,
		},
		{
			name:        "deep traversal",
			cmd:         "cp ../../other/file .",
			wantAllowed: false,
		},
		{
			name:        "traversal in flag value",
			cmd:         "go test -coverprofile=../cover.out ./...",
			wantAllowed: false,
		},

		// Blocked cases - env vars with absolute paths
		{
			name:        "env var absolute path",
			cmd:         "GOMODCACHE=/tmp/mod go test ./...",
			wantAllowed: false,
		},
		{
			name:        "multiple env vars one absolute",
			cmd:         "FOO=bar GOBIN=/usr/local/bin go install ./...",
			wantAllowed: false,
		},
	}

	rule := ConfineToWorkspace{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := parser.Parse(tt.cmd)
			got := rule.Evaluate(cmd, "")

			if got.Allowed != tt.wantAllowed {
				t.Errorf("Evaluate(%q) Allowed = %v, want %v", tt.cmd, got.Allowed, tt.wantAllowed)
			}

			if !tt.wantAllowed && got.Reason == "" {
				t.Error("blocked decision should have a reason")
			}
		})
	}
}

func TestCollectPathCandidates(t *testing.T) {
	cmd := parser.Parse("GOBIN=/tmp go test -coverprofile=cover.out ./pkg ./internal")
	candidates := collectPathCandidates(cmd)

	// Should contain: args (./pkg, ./internal), flag values (cover.out), env values (/tmp)
	if len(candidates) < 3 {
		t.Errorf("expected at least 3 candidates, got %d: %v", len(candidates), candidates)
	}

	hasAbsolute := false
	for _, c := range candidates {
		if c == "/tmp" {
			hasAbsolute = true
			break
		}
	}
	if !hasAbsolute {
		t.Error("expected /tmp from env var in candidates")
	}
}

func TestViolatesWorkspaceBoundary(t *testing.T) {
	tests := []struct {
		path     string
		violates bool
	}{
		{"", false},
		{".", false},
		{"./foo", false},
		{"foo/bar", false},
		{"./...", false},

		{"/", true},
		{"/tmp", true},
		{"/etc/passwd", true},

		{"..", true},
		{"../foo", true},
		{"../../bar", true},
		{"foo/../..", true},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := ViolatesWorkspaceBoundary(tt.path)
			if got != tt.violates {
				t.Errorf("ViolatesWorkspaceBoundary(%q) = %v, want %v", tt.path, got, tt.violates)
			}
		})
	}
}

func TestNewConfineToWorkspace(t *testing.T) {
	cfg := &config.WorkspaceConfig{
		Allow: []string{"/tmp"},
		Block: []string{".env"},
	}

	rule := NewConfineToWorkspace(cfg)

	if len(rule.Allow) != 1 || rule.Allow[0] != "/tmp" {
		t.Errorf("Allow = %v, want [/tmp]", rule.Allow)
	}
	if len(rule.Block) != 1 || rule.Block[0] != ".env" {
		t.Errorf("Block = %v, want [.env]", rule.Block)
	}
}

func TestNewConfineToWorkspaceNil(t *testing.T) {
	rule := NewConfineToWorkspace(nil)

	if rule.Allow != nil && len(rule.Allow) != 0 {
		t.Errorf("Allow should be empty for nil config")
	}
	if rule.Block != nil && len(rule.Block) != 0 {
		t.Errorf("Block should be empty for nil config")
	}
}

func TestIsBlocked(t *testing.T) {
	rule := &ConfineToWorkspace{
		Block: []string{".env", "secrets/"},
	}

	tests := []struct {
		path    string
		blocked bool
	}{
		{".env", true},
		{"secrets/key.pem", true},
		{"config.yml", false},
		{"src/main.go", false},
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

func TestIsAllowed(t *testing.T) {
	rule := &ConfineToWorkspace{
		Allow: []string{"/tmp/", "/var/cache/"},
	}

	tests := []struct {
		path    string
		allowed bool
	}{
		{"/tmp/test.txt", true},
		{"/var/cache/data", true},
		{"/etc/passwd", false},
		{"/root/.ssh", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := rule.isAllowed(tt.path)
			if got != tt.allowed {
				t.Errorf("isAllowed(%q) = %v, want %v", tt.path, got, tt.allowed)
			}
		})
	}
}

func TestEvaluateWithBlockList(t *testing.T) {
	rule := &ConfineToWorkspace{
		Block: []string{".env"},
	}

	cmd := parser.Command{Args: []string{".env"}}
	decision := rule.Evaluate(cmd, "")

	if decision.Allowed {
		t.Error("should block .env file")
	}
}

func TestEvaluateWithAllowList(t *testing.T) {
	rule := &ConfineToWorkspace{
		Allow: []string{"/tmp/"},
	}

	cmd := parser.Command{Args: []string{"/tmp/test.txt"}}
	decision := rule.Evaluate(cmd, "")

	if !decision.Allowed {
		t.Error("should allow /tmp/test.txt")
	}
}

func TestIsClaudeOperationalPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		// Operational paths should be allowed
		{"plans directory", "~/.claude/plans/plan.md", true},
		{"todos directory", "~/.claude/todos/session.json", true},
		{"file-history", "~/.claude/file-history/file.json", true},
		{"cache", "~/.claude/cache/data", true},
		{"claude root", "~/.claude", true},

		// Sensitive files should NOT be allowed (protected separately)
		{"credentials", "~/.claude/.credentials.json", false},
		{"settings", "~/.claude/settings.json", false},
		{"settings local", "~/.claude/settings.local.json", false},

		// Non-claude paths should not match
		{"random path", "/tmp/something", false},
		{"home file", "~/somefile.txt", false},
		{"similar name", "~/.claudex/file", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Expand ~ to home directory
			path := tt.path
			if len(path) > 0 && path[0] == '~' {
				if home, err := os.UserHomeDir(); err == nil {
					path = home + path[1:]
				}
			}
			got := isClaudeOperationalPath(path)
			if got != tt.want {
				t.Errorf("isClaudeOperationalPath(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestWorkspaceAllowsClaudeOperationalPaths(t *testing.T) {
	rule := &ConfineToWorkspace{} // No special config, but should allow ~/.claude/plans

	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("cannot get home dir")
	}

	// Plans should be allowed even without explicit config
	plansPath := home + "/.claude/plans/plan.md"
	cmd := parser.Command{Args: []string{plansPath}}
	decision := rule.Evaluate(cmd, "")

	if !decision.Allowed {
		t.Errorf("should allow Claude plans path %s: %s", plansPath, decision.Reason)
	}

	// But credentials should still be blocked (by IsAlwaysProtected)
	credsPath := home + "/.claude/.credentials.json"
	cmd = parser.Command{Args: []string{credsPath}}
	decision = rule.Evaluate(cmd, "")

	if decision.Allowed {
		t.Error("should block Claude credentials path")
	}
}
