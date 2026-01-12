package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/adrianpk/watchman/internal/cli"
	"github.com/adrianpk/watchman/internal/config"
	"github.com/adrianpk/watchman/internal/hook"
)

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
	// Constructor
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("cannot load config: %w", err)
	}

	evaluator := hook.NewEvaluator(cfg)

	// Setup: parse input
	var input hookInput
	if err := json.NewDecoder(os.Stdin).Decode(&input); err != nil {
		return fmt.Errorf("cannot decode input: %w", err)
	}

	// Start: evaluate
	result := evaluator.Evaluate(hook.Input{
		HookType:  input.HookType,
		ToolName:  input.ToolName,
		ToolInput: input.ToolInput,
	})

	// Output result
	if !result.Allowed {
		deny(result.Reason)
		return nil
	}

	if result.Warning != "" {
		warn(result.Warning)
	}

	allow()
	return nil
}

type hookInput struct {
	HookType  string                 `json:"hook_type"`
	ToolName  string                 `json:"tool_name"`
	ToolInput map[string]interface{} `json:"tool_input"`
}

type hookOutput struct {
	Decision string `json:"decision"`
}

func allow() {
	json.NewEncoder(os.Stdout).Encode(hookOutput{Decision: "allow"})
	os.Exit(0)
}

func deny(reason string) {
	fmt.Fprintln(os.Stderr, reason)
	os.Exit(2)
}

func warn(message string) {
	fmt.Fprintln(os.Stderr, "warning: "+message)
}

func fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
