package glob

import "testing"

func TestMatch(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		pattern string
		want    bool
	}{
		{"exact match", "main.go", "main.go", true},
		{"no match", "main.go", "other.go", false},
		{"star pattern", "main.go", "*.go", true},
		{"star pattern no match", "main.go", "*.js", false},
		{"directory pattern", "src/main.go", "src/*.go", true},
		{"deep path matches filename", "src/pkg/main.go", "*.go", true},
		{"doublestar any", "src/pkg/main.go", "**/*.go", true},
		{"doublestar prefix", "src/pkg/main.go", "src/**/*.go", true},
		{"doublestar wrong prefix", "vendor/pkg/main.go", "src/**/*.go", false},
		{"doublestar all", "anything/goes/here", "**", true},
		{"relative path", "./src/main.go", "src/*.go", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Match(tt.path, tt.pattern)
			if got != tt.want {
				t.Errorf("Match(%q, %q) = %v, want %v", tt.path, tt.pattern, got, tt.want)
			}
		})
	}
}

func TestMatchDoublestar(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		pattern string
		want    bool
	}{
		{"match all go files", "src/pkg/sub/main.go", "**/*.go", true},
		{"match with prefix", "internal/hook/eval.go", "internal/**/*.go", true},
		{"no match wrong prefix", "external/hook/eval.go", "internal/**/*.go", false},
		{"match everything", "any/path/here", "**", true},
		{"match directory", "vendor/lib/code.go", "vendor/**", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchDoublestar(tt.path, tt.pattern)
			if got != tt.want {
				t.Errorf("matchDoublestar(%q, %q) = %v, want %v", tt.path, tt.pattern, got, tt.want)
			}
		})
	}
}

func TestMatchAny(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		patterns []string
		want     bool
	}{
		{"empty patterns", "main.go", []string{}, false},
		{"single match", "main.go", []string{"*.go"}, true},
		{"multiple first match", "main.go", []string{"*.go", "*.js"}, true},
		{"multiple second match", "app.js", []string{"*.go", "*.js"}, true},
		{"multiple no match", "style.css", []string{"*.go", "*.js"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MatchAny(tt.path, tt.patterns)
			if got != tt.want {
				t.Errorf("MatchAny(%q, %v) = %v, want %v", tt.path, tt.patterns, got, tt.want)
			}
		})
	}
}
