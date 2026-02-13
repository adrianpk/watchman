package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/adrianpk/watchman/internal/cli"
	"github.com/adrianpk/watchman/internal/config"
	"github.com/adrianpk/watchman/internal/hook"
)

const logFile = "/tmp/watchman.log"

func main() {
	// Handle CLI commands
	if len(os.Args) > 1 {
		if err := runCommand(os.Args[1]); err != nil {
			fatal("%v", err)
		}
		return
	}

	// Run hook evaluation
	if err := runHook(); err != nil {
		fatal("%v", err)
	}
}

func runCommand(cmd string) error {
	switch cmd {
	case "init":
		local := len(os.Args) > 2 && os.Args[2] == "--local"
		return cli.RunInit(local)
	case "setup":
		return cli.RunSetup()
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
}

func runHook() error {
	cfg, err := config.Load()
	if err != nil {
		reason := "watchman config error: " + err.Error()
		logDeny(hookInput{}, reason)
		deny(reason)
		return nil
	}

	evaluator := hook.NewEvaluator(cfg)

	rawInput, _ := io.ReadAll(os.Stdin)

	var input hookInput
	if err := json.Unmarshal(rawInput, &input); err != nil {
		reason := "watchman input error: " + err.Error()
		logDeny(hookInput{}, reason)
		deny(reason)
		return nil
	}

	evalInput := hook.Input{
		HookType:  input.HookType,
		ToolName:  input.ToolName,
		ToolInput: input.ToolInput,
		CWD:       input.CWD,
	}

	result := evaluator.Evaluate(evalInput)

	if !result.Allowed {
		logDeny(input, result.Reason)
		deny(result.Reason)
		return nil
	}

	allow(result.Warning)
	return nil
}

func logDeny(input hookInput, reason string) {
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	ts := time.Now().Format("2006-01-02 15:04:05")

	fmt.Fprintf(f, "[%s] DENY\n", ts)
	fmt.Fprintf(f, "  tool:   %s\n", input.ToolName)
	fmt.Fprintf(f, "  cwd:    %s\n", input.CWD)
	fmt.Fprintf(f, "  reason: %s\n", reason)

	switch input.ToolName {
	case "Bash":
		if cmd, ok := input.ToolInput["command"].(string); ok {
			if len(cmd) > 200 {
				cmd = cmd[:200] + "..."
			}
			fmt.Fprintf(f, "  cmd:    %s\n", cmd)
		}
	case "Read", "Write", "Edit":
		if fp, ok := input.ToolInput["file_path"].(string); ok {
			fmt.Fprintf(f, "  path:   %s\n", fp)
		}
	case "Glob":
		if p, ok := input.ToolInput["pattern"].(string); ok {
			fmt.Fprintf(f, "  pattern: %s\n", p)
		}
		if p, ok := input.ToolInput["path"].(string); ok {
			fmt.Fprintf(f, "  path:   %s\n", p)
		}
	case "Grep":
		if p, ok := input.ToolInput["pattern"].(string); ok {
			fmt.Fprintf(f, "  pattern: %s\n", p)
		}
		if p, ok := input.ToolInput["path"].(string); ok {
			fmt.Fprintf(f, "  path:   %s\n", p)
		}
	}

	fmt.Fprintln(f, "")
}

type hookInput struct {
	HookType  string                 `json:"hook_type"`
	ToolName  string                 `json:"tool_name"`
	ToolInput map[string]interface{} `json:"tool_input"`
	CWD       string                 `json:"cwd"`
}

type hookOutput struct {
	HookSpecificOutput *hookSpecificOutput `json:"hookSpecificOutput,omitempty"`
}

type hookSpecificOutput struct {
	HookEventName      string `json:"hookEventName"`
	PermissionDecision string `json:"permissionDecision"`
	AdditionalContext  string `json:"additionalContext,omitempty"`
	Reason             string `json:"reason,omitempty"`
}

func allow(additionalContext string) {
	out := hookOutput{
		HookSpecificOutput: &hookSpecificOutput{
			HookEventName:      "PreToolUse",
			PermissionDecision: "allow",
			AdditionalContext:  additionalContext,
		},
	}
	json.NewEncoder(os.Stdout).Encode(out)
	os.Exit(0)
}

func deny(reason string) {
	out := hookOutput{
		HookSpecificOutput: &hookSpecificOutput{
			HookEventName:      "PreToolUse",
			PermissionDecision: "deny",
			Reason:             reason,
		},
	}
	json.NewEncoder(os.Stdout).Encode(out)
	ts := time.Now().Format("15:04:05")
	fmt.Fprintf(os.Stderr, "[%s] %s\n", ts, reason)
	os.Exit(2)
}

func fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
