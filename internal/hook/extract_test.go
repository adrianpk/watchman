package hook

import "testing"

func TestExtractPaths(t *testing.T) {
	tests := []struct {
		name      string
		toolName  string
		toolInput map[string]interface{}
		wantLen   int
	}{
		{
			name:      "bash command",
			toolName:  "Bash",
			toolInput: map[string]interface{}{"command": "cat file.txt"},
			wantLen:   1,
		},
		{
			name:      "read file_path",
			toolName:  "Read",
			toolInput: map[string]interface{}{"file_path": "src/main.go"},
			wantLen:   1,
		},
		{
			name:      "write file_path",
			toolName:  "Write",
			toolInput: map[string]interface{}{"file_path": "output.txt", "content": "data"},
			wantLen:   1,
		},
		{
			name:      "edit file_path",
			toolName:  "Edit",
			toolInput: map[string]interface{}{"file_path": "main.go"},
			wantLen:   1,
		},
		{
			name:      "glob with path and pattern",
			toolName:  "Glob",
			toolInput: map[string]interface{}{"path": "src", "pattern": "*.go"},
			wantLen:   2,
		},
		{
			name:      "glob with only pattern",
			toolName:  "Glob",
			toolInput: map[string]interface{}{"pattern": "*.go"},
			wantLen:   1,
		},
		{
			name:      "grep with path",
			toolName:  "Grep",
			toolInput: map[string]interface{}{"pattern": "TODO", "path": "src"},
			wantLen:   1,
		},
		{
			name:      "grep without path",
			toolName:  "Grep",
			toolInput: map[string]interface{}{"pattern": "TODO"},
			wantLen:   0,
		},
		{
			name:      "unknown tool",
			toolName:  "WebSearch",
			toolInput: map[string]interface{}{"query": "test"},
			wantLen:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractPaths(tt.toolName, tt.toolInput)
			if len(got) != tt.wantLen {
				t.Errorf("ExtractPaths() returned %d paths, want %d: %v", len(got), tt.wantLen, got)
			}
		})
	}
}

func TestExtractBashPaths(t *testing.T) {
	tests := []struct {
		name    string
		input   map[string]interface{}
		wantLen int
	}{
		{
			name:    "simple command",
			input:   map[string]interface{}{"command": "cat file.txt"},
			wantLen: 1,
		},
		{
			name:    "no command",
			input:   map[string]interface{}{},
			wantLen: 0,
		},
		{
			name:    "command with flags",
			input:   map[string]interface{}{"command": "ls -la src/"},
			wantLen: 1, // src/ is arg, -la is flag without value
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractBashPaths(tt.input)
			if len(got) < tt.wantLen {
				t.Errorf("extractBashPaths() returned %d paths, want at least %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestExtractFilePath(t *testing.T) {
	tests := []struct {
		name    string
		input   map[string]interface{}
		want    []string
		wantLen int
	}{
		{
			name:    "has file_path",
			input:   map[string]interface{}{"file_path": "main.go"},
			wantLen: 1,
		},
		{
			name:    "no file_path",
			input:   map[string]interface{}{},
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractFilePath(tt.input)
			if len(got) != tt.wantLen {
				t.Errorf("extractFilePath() returned %d paths, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestExtractGlobPaths(t *testing.T) {
	tests := []struct {
		name    string
		input   map[string]interface{}
		wantLen int
	}{
		{
			name:    "path and pattern",
			input:   map[string]interface{}{"path": "src", "pattern": "*.go"},
			wantLen: 2,
		},
		{
			name:    "only pattern",
			input:   map[string]interface{}{"pattern": "*.go"},
			wantLen: 1,
		},
		{
			name:    "only path",
			input:   map[string]interface{}{"path": "src"},
			wantLen: 1,
		},
		{
			name:    "empty",
			input:   map[string]interface{}{},
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractGlobPaths(tt.input)
			if len(got) != tt.wantLen {
				t.Errorf("extractGlobPaths() returned %d paths, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestExtractGrepPaths(t *testing.T) {
	tests := []struct {
		name    string
		input   map[string]interface{}
		wantLen int
	}{
		{
			name:    "has path",
			input:   map[string]interface{}{"path": "src", "pattern": "TODO"},
			wantLen: 1,
		},
		{
			name:    "no path",
			input:   map[string]interface{}{"pattern": "TODO"},
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractGrepPaths(tt.input)
			if len(got) != tt.wantLen {
				t.Errorf("extractGrepPaths() returned %d paths, want %d", len(got), tt.wantLen)
			}
		})
	}
}
