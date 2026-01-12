package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunInitGlobal(t *testing.T) {
	// Save original HOME and restore after test
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)

	// Use temp dir as HOME
	tmpHome := t.TempDir()
	os.Setenv("HOME", tmpHome)

	// Run init
	err := RunInit(false)
	if err != nil {
		t.Fatalf("RunInit(false) failed: %v", err)
	}

	// Verify config was created
	configPath := filepath.Join(tmpHome, ".config", "watchman", "config.yml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("config file was not created")
	}

	// Verify content
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("cannot read config: %v", err)
	}

	if !strings.Contains(string(content), "version: 1") {
		t.Error("config missing version")
	}

	if !strings.Contains(string(content), "workspace: true") {
		t.Error("config missing workspace rule")
	}
}

func TestRunInitGlobalAlreadyExists(t *testing.T) {
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)

	tmpHome := t.TempDir()
	os.Setenv("HOME", tmpHome)

	// Create config dir and file
	configDir := filepath.Join(tmpHome, ".config", "watchman")
	os.MkdirAll(configDir, 0755)
	configPath := filepath.Join(configDir, "config.yml")
	os.WriteFile(configPath, []byte("existing"), 0644)

	// Run init should not overwrite
	err := RunInit(false)
	if err != nil {
		t.Fatalf("RunInit(false) failed: %v", err)
	}

	content, _ := os.ReadFile(configPath)
	if string(content) != "existing" {
		t.Error("existing config was overwritten")
	}
}

func TestRunInitLocal(t *testing.T) {
	// Save and change working directory
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)

	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	// Run init --local
	err := RunInit(true)
	if err != nil {
		t.Fatalf("RunInit(true) failed: %v", err)
	}

	// Verify .watchman.yml was created
	configPath := filepath.Join(tmpDir, ".watchman.yml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("local config file was not created")
	}
}

func TestRunInitLocalAlreadyExists(t *testing.T) {
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)

	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	// Create existing config
	configPath := filepath.Join(tmpDir, ".watchman.yml")
	os.WriteFile(configPath, []byte("existing"), 0644)

	// Run init should not overwrite
	err := RunInit(true)
	if err != nil {
		t.Fatalf("RunInit(true) failed: %v", err)
	}

	content, _ := os.ReadFile(configPath)
	if string(content) != "existing" {
		t.Error("existing local config was overwritten")
	}
}
