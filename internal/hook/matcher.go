package hook

import (
	"regexp"
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

// Matches checks if a hook should be triggered for the given tool, paths, and command.
// Tool must match. If match_command is defined, command must match the regex.
// If paths are defined, at least one path must match.
func (m *HookMatcher) Matches(hookCfg *config.HookConfig, toolName string, paths []string, command string) bool {
	if !m.matchesTool(hookCfg.Tools, toolName) {
		return false
	}

	if hookCfg.MatchCommand != "" {
		if !m.matchesCommand(hookCfg.MatchCommand, command) {
			return false
		}
	}

	if len(hookCfg.Paths) == 0 {
		return true
	}

	return m.matchesAnyPath(hookCfg.Paths, paths)
}

func (m *HookMatcher) matchesCommand(pattern, command string) bool {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return re.MatchString(command)
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
