package glob

import (
	"path/filepath"
	"strings"
)

// Match matches a path against a glob pattern.
// Supports ** for recursive directory matching.
func Match(path, pattern string) bool {
	path = filepath.Clean(path)
	pattern = filepath.Clean(pattern)

	if strings.Contains(pattern, "**") {
		return matchDoublestar(path, pattern)
	}

	matched, _ := filepath.Match(pattern, path)
	if matched {
		return true
	}

	matched, _ = filepath.Match(pattern, filepath.Base(path))
	return matched
}

// matchDoublestar handles ** glob patterns.
func matchDoublestar(path, pattern string) bool {
	parts := strings.Split(pattern, "**")
	if len(parts) != 2 {
		return false
	}

	prefix := strings.TrimSuffix(parts[0], string(filepath.Separator))
	suffix := strings.TrimPrefix(parts[1], string(filepath.Separator))

	if prefix != "" && !strings.HasPrefix(path, prefix) {
		return false
	}

	if suffix == "" {
		return true
	}

	remaining := path
	if prefix != "" {
		remaining = strings.TrimPrefix(path, prefix)
		remaining = strings.TrimPrefix(remaining, string(filepath.Separator))
	}

	if suffix == "" {
		return true
	}

	pathParts := strings.Split(remaining, string(filepath.Separator))
	for i := range pathParts {
		candidate := strings.Join(pathParts[i:], string(filepath.Separator))
		matched, _ := filepath.Match(suffix, candidate)
		if matched {
			return true
		}
		if len(pathParts[i:]) == 1 {
			matched, _ = filepath.Match(suffix, pathParts[len(pathParts)-1])
			if matched {
				return true
			}
		}
	}

	return false
}

// MatchAny returns true if the path matches any of the patterns.
func MatchAny(path string, patterns []string) bool {
	for _, pattern := range patterns {
		if Match(path, pattern) {
			return true
		}
	}
	return false
}
