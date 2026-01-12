package hook

import (
	"strings"

	"github.com/adrianpk/watchman/internal/config"
	"github.com/adrianpk/watchman/internal/glob"
)

// HookMatcher determines if a hook should be triggered.
type HookMatcher struct{}

// NewHookMatcher creates a new matcher.
func NewHookMatcher() *HookMatcher {
	return &HookMatcher{}
}

// Matches checks if a hook should be triggered for the given tool and paths.
// Both tool AND path (if patterns defined) must match.
func (m *HookMatcher) Matches(hookCfg *config.HookConfig, toolName string, paths []string) bool {
	if !m.matchesTool(hookCfg.Tools, toolName) {
		return false
	}

	if len(hookCfg.Paths) == 0 {
		return true
	}

	return m.matchesAnyPath(hookCfg.Paths, paths)
}

func (m *HookMatcher) matchesTool(tools []string, toolName string) bool {
	for _, t := range tools {
		if strings.EqualFold(t, toolName) {
			return true
		}
	}
	return false
}

func (m *HookMatcher) matchesAnyPath(patterns []string, paths []string) bool {
	for _, path := range paths {
		if glob.MatchAny(path, patterns) {
			return true
		}
	}
	return false
}
