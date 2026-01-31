package policy

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/adrianpk/watchman/internal/config"
	"github.com/adrianpk/watchman/internal/glob"
	"github.com/adrianpk/watchman/internal/parser"
)

// writeTools are tools that modify files.
var writeTools = map[string]bool{
	"Write":        true,
	"Edit":         true,
	"NotebookEdit": true,
}

// ScopeToFiles restricts modifications to declared file patterns.
type ScopeToFiles struct {
	Allow []string
	Block []string
}

// NewScopeToFiles creates a scope rule from config.
func NewScopeToFiles(cfg *config.ScopeConfig) *ScopeToFiles {
	if cfg == nil {
		return &ScopeToFiles{}
	}
	return &ScopeToFiles{
		Allow: cfg.Allow,
		Block: cfg.Block,
	}
}

// Evaluate checks if the command modifies files within the defined scope.
func (r *ScopeToFiles) Evaluate(toolName string, cmd parser.Command, cwd string) Decision {
	if !writeTools[toolName] {
		return Decision{Allowed: true}
	}

	paths := collectPathCandidates(cmd)
	for _, p := range paths {
		if r.isBlocked(p) {
			return Decision{
				Allowed: false,
				Reason:  "scope.block: " + p + " matches blocked pattern",
			}
		}
		if !r.isInScope(p, cwd) {
			return Decision{
				Allowed: false,
				Reason:  "scope.allow: " + p + " does not match any allowed pattern " + r.summarizeAllow(),
			}
		}
	}

	return Decision{Allowed: true}
}

// summarizeAllow returns a short summary of allowed patterns for error messages.
func (r *ScopeToFiles) summarizeAllow() string {
	if len(r.Allow) == 0 {
		return "(none configured)"
	}
	if len(r.Allow) <= 5 {
		return "(" + strings.Join(r.Allow, ", ") + ")"
	}
	return "(" + strings.Join(r.Allow[:5], ", ") + ", ...)"
}

// isBlocked checks if a path matches any block pattern.
func (r *ScopeToFiles) isBlocked(p string) bool {
	return glob.MatchAny(p, r.Block)
}

// isInScope checks if a path is within the allowed scope.
// If no allow patterns are defined, all paths are in scope.
func (r *ScopeToFiles) isInScope(p string, cwd string) bool {
	if len(r.Allow) == 0 {
		return true
	}

	// Normalize path to relative for glob matching
	// This allows patterns like "src/**/*.go" to match absolute paths
	relPath := toRelativePath(p, cwd)

	// Try both the original path and the relative version
	return glob.MatchAny(p, r.Allow) || glob.MatchAny(relPath, r.Allow)
}

// toRelativePath converts an absolute path to relative (if within cwd).
func toRelativePath(p string, cwd string) string {
	if !filepath.IsAbs(p) {
		return p
	}

	// Use provided cwd, fallback to os.Getwd() if empty
	if cwd == "" {
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			return p
		}
	}

	// Check if path is within cwd
	if strings.HasPrefix(p, cwd+string(filepath.Separator)) {
		rel, err := filepath.Rel(cwd, p)
		if err == nil {
			return rel
		}
	}

	return p
}
