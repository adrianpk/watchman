// Package state handles persistent state for reminders.
package state

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/adrianpk/watchman/internal/config"
)

const stateFileName = ".watchman-state"

// State represents the persistent state for reminders.
type State struct {
	TaskCount   int                    `json:"task_count"`
	LastChecked map[string]time.Time   `json:"last_checked"` // Per-reminder last trigger time
	TaskCounts  map[string]int         `json:"task_counts"`  // Per-reminder task count since last trigger
}

// Manager handles state persistence and reminder checks.
type Manager struct {
	state    *State
	statePath string
}

// NewManager creates a new state manager for the current directory.
func NewManager() *Manager {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	return &Manager{
		statePath: filepath.Join(cwd, stateFileName),
	}
}

// Load loads the state from disk, or initializes a new state if none exists.
func (m *Manager) Load() error {
	m.state = &State{
		LastChecked: make(map[string]time.Time),
		TaskCounts:  make(map[string]int),
	}

	data, err := os.ReadFile(m.statePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Fresh state
		}
		return err
	}

	return json.Unmarshal(data, m.state)
}

// Save persists the state to disk.
func (m *Manager) Save() error {
	data, err := json.MarshalIndent(m.state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.statePath, data, 0644)
}

// IncrementTaskCount increments the global task counter and per-reminder counters.
func (m *Manager) IncrementTaskCount() {
	m.state.TaskCount++
	for name := range m.state.TaskCounts {
		m.state.TaskCounts[name]++
	}
}

// CheckReminders checks all configured reminders and returns any triggered messages.
func (m *Manager) CheckReminders(reminders []config.ReminderConfig) []string {
	var triggered []string
	now := time.Now()

	for _, r := range reminders {
		// Initialize tracking for new reminders
		if _, ok := m.state.TaskCounts[r.Name]; !ok {
			m.state.TaskCounts[r.Name] = 0
			m.state.LastChecked[r.Name] = now
		}

		shouldTrigger := false

		// Check task count trigger
		if r.EveryTasks > 0 && m.state.TaskCounts[r.Name] >= r.EveryTasks {
			shouldTrigger = true
		}

		// Check time trigger
		if r.EveryMinutes > 0 {
			lastCheck := m.state.LastChecked[r.Name]
			elapsed := now.Sub(lastCheck)
			if elapsed >= time.Duration(r.EveryMinutes)*time.Minute {
				shouldTrigger = true
			}
		}

		if shouldTrigger {
			triggered = append(triggered, r.Message)
			// Reset counters for this reminder
			m.state.TaskCounts[r.Name] = 0
			m.state.LastChecked[r.Name] = now
		}
	}

	return triggered
}

// StatePath returns the path to the state file.
func (m *Manager) StatePath() string {
	return m.statePath
}
