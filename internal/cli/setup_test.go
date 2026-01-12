package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRunSetup(t *testing.T) {
	// Save original HOME and restore after test
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)

	tmpHome := t.TempDir()
	os.Setenv("HOME", tmpHome)

	// Create go/bin directory for watchman path
	os.MkdirAll(filepath.Join(tmpHome, "go", "bin"), 0755)

	// Run setup
	err := RunSetup()
	if err != nil {
		t.Fatalf("RunSetup() failed: %v", err)
	}

	// Verify settings.json was created
	settingsPath := filepath.Join(tmpHome, ".claude", "settings.json")
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		t.Error("settings.json was not created")
	}

	// Verify content
	content, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("cannot read settings: %v", err)
	}

	var settings map[string]interface{}
	if err := json.Unmarshal(content, &settings); err != nil {
		t.Fatalf("cannot parse settings: %v", err)
	}

	hooks, ok := settings["hooks"].(map[string]interface{})
	if !ok {
		t.Fatal("settings missing hooks")
	}

	if _, ok := hooks["PreToolUse"]; !ok {
		t.Error("settings missing PreToolUse hook")
	}
}

func TestRunSetupAlreadyConfigured(t *testing.T) {
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)

	tmpHome := t.TempDir()
	os.Setenv("HOME", tmpHome)

	// Create .claude dir and existing settings with watchman hook
	claudeDir := filepath.Join(tmpHome, ".claude")
	os.MkdirAll(claudeDir, 0755)

	existingSettings := map[string]interface{}{
		"hooks": map[string]interface{}{
			"PreToolUse": []interface{}{
				map[string]interface{}{
					"matcher": "*",
					"hooks": []interface{}{
						map[string]interface{}{
							"type":    "command",
							"command": filepath.Join(tmpHome, "go", "bin", "watchman"),
						},
					},
				},
			},
		},
	}
	data, _ := json.Marshal(existingSettings)
	settingsPath := filepath.Join(claudeDir, "settings.json")
	os.WriteFile(settingsPath, data, 0644)

	// Run setup should not fail
	err := RunSetup()
	if err != nil {
		t.Fatalf("RunSetup() failed: %v", err)
	}
}

func TestRunSetupExistingSettings(t *testing.T) {
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)

	tmpHome := t.TempDir()
	os.Setenv("HOME", tmpHome)

	// Create .claude dir with existing settings (no hooks)
	claudeDir := filepath.Join(tmpHome, ".claude")
	os.MkdirAll(claudeDir, 0755)

	existingSettings := map[string]interface{}{
		"other": "setting",
	}
	data, _ := json.Marshal(existingSettings)
	settingsPath := filepath.Join(claudeDir, "settings.json")
	os.WriteFile(settingsPath, data, 0644)

	// Run setup
	err := RunSetup()
	if err != nil {
		t.Fatalf("RunSetup() failed: %v", err)
	}

	// Verify existing settings preserved
	content, _ := os.ReadFile(settingsPath)
	var settings map[string]interface{}
	json.Unmarshal(content, &settings)

	if settings["other"] != "setting" {
		t.Error("existing settings were not preserved")
	}
}

func TestHasWatchmanHook(t *testing.T) {
	tests := []struct {
		name       string
		preToolUse []interface{}
		want       bool
	}{
		{
			name:       "empty list",
			preToolUse: []interface{}{},
			want:       false,
		},
		{
			name: "short form watchman",
			preToolUse: []interface{}{
				map[string]interface{}{
					"matcher": "*",
					"hooks":   []interface{}{"watchman"},
				},
			},
			want: true,
		},
		{
			name: "long form watchman",
			preToolUse: []interface{}{
				map[string]interface{}{
					"matcher": "Bash",
					"hooks": []interface{}{
						map[string]interface{}{
							"type":    "command",
							"command": "/home/user/go/bin/watchman",
						},
					},
				},
			},
			want: true,
		},
		{
			name: "other hook",
			preToolUse: []interface{}{
				map[string]interface{}{
					"matcher": "*",
					"hooks":   []interface{}{"other-tool"},
				},
			},
			want: false,
		},
		{
			name: "invalid entry",
			preToolUse: []interface{}{
				"not a map",
			},
			want: false,
		},
		{
			name: "invalid hooks list",
			preToolUse: []interface{}{
				map[string]interface{}{
					"matcher": "*",
					"hooks":   "not a list",
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasWatchmanHook(tt.preToolUse, "/any/path")
			if got != tt.want {
				t.Errorf("hasWatchmanHook() = %v, want %v", got, tt.want)
			}
		})
	}
}
