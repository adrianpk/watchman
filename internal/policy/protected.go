package policy

import (
	"os"
	"path/filepath"
	"strings"
)

// alwaysProtected contains paths that are NEVER accessible, regardless of config.
// This is a hardcoded security boundary that cannot be overridden.
var alwaysProtected = []string{
	// Claude Code credentials and settings (not operational dirs like plans/, todos/)
	"~/.claude/.credentials.json",
	"~/.claude/settings.json",
	"~/.claude/settings.local.json",
	// SSH and cloud credentials
	"~/.ssh/",
	"~/.aws/",
	"~/.gnupg/",
	"~/.gpg/",
	"~/.config/gh/",
	"~/.config/watchman/",
	"~/.netrc",
	"~/.git-credentials",
	"~/go/bin/watchman",
}

// protectedFilenames are filenames that are protected in any directory.
var protectedFilenames = []string{
	".watchman.yml",
}

// IsAlwaysProtected checks if a path matches any hardcoded protected path.
// This check cannot be overridden by configuration.
func IsAlwaysProtected(p string) bool {
	if p == "" {
		return false
	}

	absPath := resolvePath(p)

	filename := filepath.Base(absPath)
	for _, protected := range protectedFilenames {
		if filename == protected {
			return true
		}
	}

	for _, pattern := range alwaysProtected {
		isDir := strings.HasSuffix(pattern, "/")

		expandedPattern := strings.TrimSuffix(pattern, "/")
		if strings.HasPrefix(expandedPattern, "~/") {
			if home, err := os.UserHomeDir(); err == nil {
				expandedPattern = filepath.Join(home, expandedPattern[2:])
			}
		}

		if isDir {
			if absPath == expandedPattern || strings.HasPrefix(absPath, expandedPattern+string(filepath.Separator)) {
				return true
			}
		} else if absPath == expandedPattern {
			return true
		}
	}

	return false
}

// resolvePath converts a path to absolute form.
func resolvePath(p string) string {
	if strings.HasPrefix(p, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			p = filepath.Join(home, p[2:])
		}
	}
	if filepath.IsAbs(p) {
		return filepath.Clean(p)
	}
	if cwd, err := os.Getwd(); err == nil {
		return filepath.Clean(filepath.Join(cwd, p))
	}
	return filepath.Clean(p)
}

// MatchProtectedPath checks if a path matches a protected pattern.
// Supports ~/ expansion and directory patterns (ending with /).
func MatchProtectedPath(path, pattern string) bool {
	absPath := resolvePath(path)

	// Check if pattern is a filename (no path separators)
	if !strings.Contains(pattern, "/") {
		filename := filepath.Base(absPath)
		return filename == pattern
	}

	return matchPath(absPath, pattern)
}

// matchPath checks if a path matches a pattern.
// Supports exact match and prefix match (pattern ending with /).
func matchPath(path, pattern string) bool {
	if strings.HasPrefix(pattern, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			pattern = filepath.Join(home, pattern[2:])
		}
	}

	if path == pattern {
		return true
	}

	if strings.HasSuffix(pattern, "/") {
		return strings.HasPrefix(path, pattern) || path == strings.TrimSuffix(pattern, "/")
	}

	if strings.HasPrefix(path, pattern+"/") || strings.HasPrefix(path, pattern+string(filepath.Separator)) {
		return true
	}

	return false
}
