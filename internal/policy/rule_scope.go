package policy

import (
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
func (r *ScopeToFiles) Evaluate(toolName string, cmd parser.Command) Decision {
	if !writeTools[toolName] {
		return Decision{Allowed: true}
	}

	paths := collectPathCandidates(cmd)
	for _, p := range paths {
		if r.isBlocked(p) {
			return Decision{
				Allowed: false,
				Reason:  "path is blocked by scope configuration: " + p,
			}
		}
		if !r.isInScope(p) {
			return Decision{
				Allowed: false,
				Reason:  "path is outside allowed scope: " + p,
			}
		}
	}

	return Decision{Allowed: true}
}

// isBlocked checks if a path matches any block pattern.
func (r *ScopeToFiles) isBlocked(p string) bool {
	return glob.MatchAny(p, r.Block)
}

// isInScope checks if a path is within the allowed scope.
// If no allow patterns are defined, all paths are in scope.
func (r *ScopeToFiles) isInScope(p string) bool {
	if len(r.Allow) == 0 {
		return true
	}
	return glob.MatchAny(p, r.Allow)
}
