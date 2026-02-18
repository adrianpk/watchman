package main

import (
	cryptoRand "crypto/rand"
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
	case "log":
		return cli.RunLog(os.Args[2:])
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
}

func runHook() error {
	cfg, err := config.Load()
	if err != nil {
		reason := "watchman config error: " + err.Error()
		id := logDeny(hookInput{}, reason)
		deny(reason, id)
		return nil
	}

	evaluator := hook.NewEvaluator(cfg)

	rawInput, _ := io.ReadAll(os.Stdin)

	var input hookInput
	if err := json.Unmarshal(rawInput, &input); err != nil {
		reason := "watchman input error: " + err.Error()
		id := logDeny(hookInput{}, reason)
		deny(reason, id)
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
		id := logDeny(input, result.Reason)
		deny(result.Reason, id)
		return nil
	}

	allow(result.Warning)
	return nil
}

type logEntry struct {
	Timestamp string                 `json:"ts"`
	ID        string                 `json:"id"`
	Decision  string                 `json:"decision"`
	Tool      string                 `json:"tool"`
	CWD       string                 `json:"cwd"`
	Reason    string                 `json:"reason"`
	Command   string                 `json:"cmd,omitempty"`
	Path      string                 `json:"path,omitempty"`
	Pattern   string                 `json:"pattern,omitempty"`
	Input     map[string]interface{} `json:"input,omitempty"`
}

func generateID() string {
	b := make([]byte, 4)
	if _, err := io.ReadFull(cryptoRand.Reader, b); err != nil {
		return fmt.Sprintf("%08x", time.Now().UnixNano()&0xFFFFFFFF)
	}
	return fmt.Sprintf("%08x", b)
}

func logDeny(input hookInput, reason string) string {
	id := generateID()

	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return id
	}
	defer f.Close()

	entry := logEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		ID:        id,
		Decision:  "deny",
		Tool:      input.ToolName,
		CWD:       input.CWD,
		Reason:    reason,
	}

	switch input.ToolName {
	case "Bash":
		if cmd, ok := input.ToolInput["command"].(string); ok {
			entry.Command = cmd
		}
	case "Read", "Write", "Edit":
		if fp, ok := input.ToolInput["file_path"].(string); ok {
			entry.Path = fp
		}
	case "Glob", "Grep":
		if p, ok := input.ToolInput["pattern"].(string); ok {
			entry.Pattern = p
		}
		if p, ok := input.ToolInput["path"].(string); ok {
			entry.Path = p
		}
	default:
		entry.Input = input.ToolInput
	}

	json.NewEncoder(f).Encode(entry)
	return id
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

func deny(reason, id string) {
	reasonWithID := fmt.Sprintf("[%s] %s", id, reason)
	hint := fmt.Sprintf("Run `watchman log` for details (id: %s)", id)
	out := hookOutput{
		HookSpecificOutput: &hookSpecificOutput{
			HookEventName:      "PreToolUse",
			PermissionDecision: "deny",
			Reason:             reasonWithID,
			AdditionalContext:  hint,
		},
	}
	json.NewEncoder(os.Stdout).Encode(out)
	ts := time.Now().Format("15:04:05")
	fmt.Fprintf(os.Stderr, "[%s] [%s] %s\n", ts, id, reason)
	os.Exit(2)
}

func fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
