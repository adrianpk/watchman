package hook

import (
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/adrianpk/watchman/internal/config"
)

func testdataPath(name string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "testdata", name)
}

func TestNewHookExecutor(t *testing.T) {
	e := NewHookExecutor()
	if e == nil {
		t.Error("NewHookExecutor returned nil")
	}
	if e.defaultTimeout != defaultTimeout {
		t.Errorf("defaultTimeout = %v, want %v", e.defaultTimeout, defaultTimeout)
	}
}

func TestHookExecutorExecuteAllow(t *testing.T) {
	e := NewHookExecutor()
	hookCfg := &config.HookConfig{
		Name:    "test-allow",
		Command: testdataPath("allow.sh"),
	}

	result := e.Execute(hookCfg, HookInput{})
	if !result.Allowed {
		t.Errorf("Execute() allowed = false, want true")
	}
}

func TestHookExecutorExecuteDeny(t *testing.T) {
	e := NewHookExecutor()
	hookCfg := &config.HookConfig{
		Name:    "test-deny",
		Command: testdataPath("deny.sh"),
	}

	result := e.Execute(hookCfg, HookInput{})
	if result.Allowed {
		t.Errorf("Execute() allowed = true, want false")
	}
	if result.Reason != "test denial" {
		t.Errorf("Execute() reason = %q, want %q", result.Reason, "test denial")
	}
}

func TestHookExecutorExecuteAdvise(t *testing.T) {
	e := NewHookExecutor()
	hookCfg := &config.HookConfig{
		Name:    "test-advise",
		Command: testdataPath("advise.sh"),
	}

	result := e.Execute(hookCfg, HookInput{})
	if !result.Allowed {
		t.Errorf("Execute() allowed = false, want true")
	}
	if result.Warning != "consider this" {
		t.Errorf("Execute() warning = %q, want %q", result.Warning, "consider this")
	}
}

func TestHookExecutorExecuteExitCodeFallback(t *testing.T) {
	e := NewHookExecutor()
	hookCfg := &config.HookConfig{
		Name:    "test-exitcode",
		Command: testdataPath("exitcode.sh"),
	}

	result := e.Execute(hookCfg, HookInput{})
	if result.Allowed {
		t.Errorf("Execute() allowed = true, want false")
	}
	if result.Reason != "error message\n" {
		t.Errorf("Execute() reason = %q, want %q", result.Reason, "error message\n")
	}
}

func TestHookExecutorExecuteTimeout(t *testing.T) {
	e := NewHookExecutor()
	hookCfg := &config.HookConfig{
		Name:    "test-slow",
		Command: testdataPath("slow.sh"),
		Timeout: 100 * time.Millisecond,
	}

	result := e.Execute(hookCfg, HookInput{})
	if !result.Allowed {
		t.Errorf("Execute() allowed = false, want true (on_error default is allow)")
	}
	if result.Warning == "" {
		t.Error("Execute() should have warning about timeout")
	}
}

func TestHookExecutorExecuteTimeoutDeny(t *testing.T) {
	e := NewHookExecutor()
	hookCfg := &config.HookConfig{
		Name:    "test-slow-deny",
		Command: testdataPath("slow.sh"),
		Timeout: 100 * time.Millisecond,
		OnError: "deny",
	}

	result := e.Execute(hookCfg, HookInput{})
	if result.Allowed {
		t.Errorf("Execute() allowed = true, want false (on_error is deny)")
	}
}

func TestHookExecutorExecuteNotFound(t *testing.T) {
	e := NewHookExecutor()
	hookCfg := &config.HookConfig{
		Name:    "test-notfound",
		Command: "/nonexistent/command",
	}

	result := e.Execute(hookCfg, HookInput{})
	if !result.Allowed {
		t.Errorf("Execute() allowed = false, want true (default on_error is allow)")
	}
}

func TestHookExecutorExecuteNotFoundDeny(t *testing.T) {
	e := NewHookExecutor()
	hookCfg := &config.HookConfig{
		Name:    "test-notfound-deny",
		Command: "/nonexistent/command",
		OnError: "deny",
	}

	result := e.Execute(hookCfg, HookInput{})
	if result.Allowed {
		t.Errorf("Execute() allowed = true, want false (on_error is deny)")
	}
}

func TestHookExecutorOutputToResult(t *testing.T) {
	e := NewHookExecutor()

	tests := []struct {
		name        string
		output      HookOutput
		wantAllowed bool
		wantReason  string
		wantWarning string
	}{
		{
			name:        "allow",
			output:      HookOutput{Decision: "allow"},
			wantAllowed: true,
		},
		{
			name:        "deny",
			output:      HookOutput{Decision: "deny", Reason: "blocked"},
			wantAllowed: false,
			wantReason:  "blocked",
		},
		{
			name:        "advise",
			output:      HookOutput{Decision: "advise", Warning: "be careful"},
			wantAllowed: true,
			wantWarning: "be careful",
		},
		{
			name:        "unknown decision treated as allow",
			output:      HookOutput{Decision: "unknown"},
			wantAllowed: true,
		},
		{
			name:        "empty decision treated as allow",
			output:      HookOutput{},
			wantAllowed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := e.outputToResult(tt.output)
			if got.Allowed != tt.wantAllowed {
				t.Errorf("outputToResult() allowed = %v, want %v", got.Allowed, tt.wantAllowed)
			}
			if got.Reason != tt.wantReason {
				t.Errorf("outputToResult() reason = %q, want %q", got.Reason, tt.wantReason)
			}
			if got.Warning != tt.wantWarning {
				t.Errorf("outputToResult() warning = %q, want %q", got.Warning, tt.wantWarning)
			}
		})
	}
}

func TestHookExecutorHandleError(t *testing.T) {
	e := NewHookExecutor()

	tests := []struct {
		name        string
		onError     string
		wantAllowed bool
	}{
		{"default allow", "", true},
		{"explicit allow", "allow", true},
		{"deny", "deny", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hookCfg := &config.HookConfig{OnError: tt.onError}
			got := e.handleError(hookCfg, "test error")
			if got.Allowed != tt.wantAllowed {
				t.Errorf("handleError() allowed = %v, want %v", got.Allowed, tt.wantAllowed)
			}
		})
	}
}
