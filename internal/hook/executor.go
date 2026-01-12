package hook

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"time"

	"github.com/adrianpk/watchman/internal/config"
)

const defaultTimeout = 5 * time.Second

// HookInput is the JSON structure sent to external hooks via stdin.
type HookInput struct {
	ToolName   string                 `json:"tool_name"`
	ToolInput  map[string]interface{} `json:"tool_input"`
	Paths      []string               `json:"paths"`
	WorkingDir string                 `json:"working_dir"`
}

// HookOutput is the JSON structure expected from hook stdout.
type HookOutput struct {
	Decision string `json:"decision"`
	Reason   string `json:"reason,omitempty"`
	Warning  string `json:"warning,omitempty"`
}

// HookExecutor runs external hook commands.
type HookExecutor struct {
	defaultTimeout time.Duration
}

// NewHookExecutor creates a new executor with default settings.
func NewHookExecutor() *HookExecutor {
	return &HookExecutor{
		defaultTimeout: defaultTimeout,
	}
}

// Execute runs a hook and returns its decision.
func (e *HookExecutor) Execute(hookCfg *config.HookConfig, input HookInput) Result {
	timeout := e.defaultTimeout
	if hookCfg.Timeout > 0 {
		timeout = hookCfg.Timeout
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, hookCfg.Command, hookCfg.Args...)

	inputJSON, err := json.Marshal(input)
	if err != nil {
		return e.handleError(hookCfg, "failed to encode input: "+err.Error())
	}
	cmd.Stdin = bytes.NewReader(inputJSON)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	if ctx.Err() == context.DeadlineExceeded {
		return e.handleError(hookCfg, "hook timed out after "+timeout.String())
	}

	if err != nil {
		if isCommandNotFound(err) {
			return e.handleError(hookCfg, "command not found: "+hookCfg.Command)
		}
	}

	if stdout.Len() > 0 {
		var output HookOutput
		if jsonErr := json.Unmarshal(stdout.Bytes(), &output); jsonErr == nil {
			return e.outputToResult(output)
		}
	}

	if err != nil {
		reason := stderr.String()
		if reason == "" {
			reason = "hook denied (exit code non-zero)"
		}
		return Result{Allowed: false, Reason: reason}
	}

	return Result{Allowed: true}
}

func isCommandNotFound(err error) bool {
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode() == 127
	}
	if _, ok := err.(*exec.Error); ok {
		return true
	}
	if pathErr, ok := err.(*os.PathError); ok {
		return os.IsNotExist(pathErr)
	}
	return false
}

func (e *HookExecutor) outputToResult(output HookOutput) Result {
	switch output.Decision {
	case "deny":
		return Result{Allowed: false, Reason: output.Reason}
	case "advise":
		return Result{Allowed: true, Warning: output.Warning}
	default:
		return Result{Allowed: true}
	}
}

func (e *HookExecutor) handleError(hookCfg *config.HookConfig, errMsg string) Result {
	if hookCfg.OnError == "deny" {
		return Result{Allowed: false, Reason: "hook error: " + errMsg}
	}
	return Result{Allowed: true, Warning: "hook error (allowed): " + errMsg}
}
