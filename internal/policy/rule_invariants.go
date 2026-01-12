package policy

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/adrianpk/watchman/internal/config"
	"github.com/adrianpk/watchman/internal/glob"
)

// InvariantsRule enforces declarative structural checks.
type InvariantsRule struct {
	cfg *config.InvariantsConfig
}

// NewInvariantsRule creates an invariants rule from config.
func NewInvariantsRule(cfg *config.InvariantsConfig) *InvariantsRule {
	if cfg == nil {
		return &InvariantsRule{cfg: &config.InvariantsConfig{}}
	}
	return &InvariantsRule{cfg: cfg}
}

// Evaluate checks if the file modification violates any invariants.
// Only applies to modification tools (Write, Edit, NotebookEdit).
func (r *InvariantsRule) Evaluate(toolName, filePath, content string) Decision {
	if !writeTools[toolName] {
		return Decision{Allowed: true}
	}

	// Check coexistence rules
	if decision := r.checkCoexistence(filePath); !decision.Allowed {
		return decision
	}

	// Check content rules
	if decision := r.checkContent(filePath, content); !decision.Allowed {
		return decision
	}

	// Check import rules
	if decision := r.checkImports(filePath, content); !decision.Allowed {
		return decision
	}

	// Check naming rules
	if decision := r.checkNaming(filePath); !decision.Allowed {
		return decision
	}

	// Check required files rules
	if decision := r.checkRequired(filePath); !decision.Allowed {
		return decision
	}

	return Decision{Allowed: true}
}

// checkCoexistence ensures related files exist together.
func (r *InvariantsRule) checkCoexistence(filePath string) Decision {
	for _, check := range r.cfg.Coexistence {
		if !glob.Match(filePath, check.If) {
			continue
		}

		requiredPath := expandPlaceholders(check.Require, filePath)
		if _, err := os.Stat(requiredPath); os.IsNotExist(err) {
			msg := check.Message
			if msg == "" {
				msg = "coexistence check failed: " + check.Name + " requires " + requiredPath
			}
			return Decision{Allowed: false, Reason: msg}
		}
	}
	return Decision{Allowed: true}
}

// checkContent validates file content against patterns.
func (r *InvariantsRule) checkContent(filePath, content string) Decision {
	for _, check := range r.cfg.Content {
		if !matchesPathPatterns(filePath, check.Paths) {
			continue
		}

		// Check forbidden patterns
		if check.Forbid != "" {
			re, err := regexp.Compile(check.Forbid)
			if err != nil {
				continue // Skip invalid regex
			}
			if re.MatchString(content) {
				msg := check.Message
				if msg == "" {
					msg = "content check failed: " + check.Name + " forbids pattern: " + check.Forbid
				}
				return Decision{Allowed: false, Reason: msg}
			}
		}

		// Check required patterns
		if check.Require != "" {
			re, err := regexp.Compile(check.Require)
			if err != nil {
				continue // Skip invalid regex
			}
			if !re.MatchString(content) {
				msg := check.Message
				if msg == "" {
					msg = "content check failed: " + check.Name + " requires pattern: " + check.Require
				}
				return Decision{Allowed: false, Reason: msg}
			}
		}
	}
	return Decision{Allowed: true}
}

// checkImports validates import statements (regex-based).
func (r *InvariantsRule) checkImports(filePath, content string) Decision {
	for _, check := range r.cfg.Imports {
		if !matchesPathPatterns(filePath, check.Paths) {
			continue
		}

		re, err := regexp.Compile(check.Forbid)
		if err != nil {
			continue // Skip invalid regex
		}
		if re.MatchString(content) {
			msg := check.Message
			if msg == "" {
				msg = "import check failed: " + check.Name + " forbids import matching: " + check.Forbid
			}
			return Decision{Allowed: false, Reason: msg}
		}
	}
	return Decision{Allowed: true}
}

// checkNaming validates file naming conventions.
func (r *InvariantsRule) checkNaming(filePath string) Decision {
	for _, check := range r.cfg.Naming {
		if !matchesPathPatterns(filePath, check.Paths) {
			continue
		}

		filename := filepath.Base(filePath)
		re, err := regexp.Compile(check.Pattern)
		if err != nil {
			continue // Skip invalid regex
		}
		if !re.MatchString(filename) {
			msg := check.Message
			if msg == "" {
				msg = "naming check failed: " + check.Name + " requires pattern: " + check.Pattern
			}
			return Decision{Allowed: false, Reason: msg}
		}
	}
	return Decision{Allowed: true}
}

// checkRequired ensures certain files exist in directories.
func (r *InvariantsRule) checkRequired(filePath string) Decision {
	dir := filepath.Dir(filePath)

	for _, check := range r.cfg.Required {
		if !glob.Match(dir, check.Dirs) {
			continue
		}

		// Check "when" condition if specified
		if check.When != "" {
			entries, err := os.ReadDir(dir)
			if err != nil {
				continue
			}
			found := false
			for _, entry := range entries {
				if glob.Match(entry.Name(), check.When) {
					found = true
					break
				}
			}
			if !found {
				continue // "when" condition not met, skip this check
			}
		}

		requiredFile := filepath.Join(dir, check.Require)
		if _, err := os.Stat(requiredFile); os.IsNotExist(err) {
			msg := check.Message
			if msg == "" {
				msg = "required check failed: " + check.Name + " requires " + check.Require + " in " + dir
			}
			return Decision{Allowed: false, Reason: msg}
		}
	}
	return Decision{Allowed: true}
}

// expandPlaceholders replaces ${name}, ${base}, ${ext} in a pattern.
func expandPlaceholders(pattern, filePath string) string {
	dir := filepath.Dir(filePath)
	base := filepath.Base(filePath)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	// ${base} for test files: user_test.go -> user
	baseName := name
	if strings.HasSuffix(name, "_test") {
		baseName = strings.TrimSuffix(name, "_test")
	}

	result := pattern
	result = strings.ReplaceAll(result, "${name}", name)
	result = strings.ReplaceAll(result, "${base}", baseName)
	result = strings.ReplaceAll(result, "${ext}", ext)

	// If result is relative, join with directory
	if !filepath.IsAbs(result) && !strings.HasPrefix(result, ".") {
		result = filepath.Join(dir, result)
	}

	return result
}

// matchesPathPatterns checks if a path matches any pattern in the list.
// Supports exclusion patterns with ! prefix.
func matchesPathPatterns(filePath string, patterns []string) bool {
	if len(patterns) == 0 {
		return true
	}

	included := false
	excluded := false

	for _, pattern := range patterns {
		if strings.HasPrefix(pattern, "!") {
			// Exclusion pattern
			if glob.Match(filePath, strings.TrimPrefix(pattern, "!")) {
				excluded = true
			}
		} else {
			// Inclusion pattern
			if glob.Match(filePath, pattern) {
				included = true
			}
		}
	}

	return included && !excluded
}
