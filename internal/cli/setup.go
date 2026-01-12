package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// RunSetup configures the Claude Code hook.
func RunSetup() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot get home directory: %w", err)
	}

	claudeDir := filepath.Join(home, ".claude")
	settingsPath := filepath.Join(claudeDir, "settings.json")
	watchmanPath := filepath.Join(home, "go", "bin", "watchman")

	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		return fmt.Errorf("cannot create .claude directory: %w", err)
	}

	settings := make(map[string]interface{})

	data, err := os.ReadFile(settingsPath)
	if err == nil && len(data) > 0 {
		if err := json.Unmarshal(data, &settings); err != nil {
			return fmt.Errorf("cannot parse settings.json: %w", err)
		}
	}

	hooks, ok := settings["hooks"].(map[string]interface{})
	if !ok {
		hooks = make(map[string]interface{})
		settings["hooks"] = hooks
	}

	preToolUse, ok := hooks["PreToolUse"].([]interface{})
	if !ok {
		preToolUse = []interface{}{}
	}

	if hasWatchmanHook(preToolUse, watchmanPath) {
		fmt.Println("Watchman hook already configured")
		return nil
	}

	watchmanHook := map[string]interface{}{
		"matcher": "*",
		"hooks": []interface{}{
			map[string]interface{}{
				"type":    "command",
				"command": watchmanPath,
			},
		},
	}

	hooks["PreToolUse"] = []interface{}{watchmanHook}

	output, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot marshal settings: %w", err)
	}

	if err := os.WriteFile(settingsPath, output, 0644); err != nil {
		return fmt.Errorf("cannot write settings.json: %w", err)
	}

	fmt.Printf("Configured hook: %s\n", settingsPath)
	fmt.Println("Run 'watchman init' to create watchman config")
	return nil
}

func hasWatchmanHook(preToolUse []interface{}, watchmanPath string) bool {
	for _, entry := range preToolUse {
		e, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}
		hooksList, ok := e["hooks"].([]interface{})
		if !ok {
			continue
		}
		for _, h := range hooksList {
			if h == "watchman" {
				return true
			}
			if hm, ok := h.(map[string]interface{}); ok {
				if cmd, ok := hm["command"].(string); ok {
					if strings.Contains(cmd, "watchman") {
						return true
					}
				}
			}
		}
	}
	return false
}
