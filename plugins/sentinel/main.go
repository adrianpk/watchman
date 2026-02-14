package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/adrianpk/watchman/plugins/sentinel/internal/client"
	"github.com/adrianpk/watchman/plugins/sentinel/internal/config"
	"github.com/adrianpk/watchman/plugins/sentinel/internal/evaluator"
	"github.com/adrianpk/watchman/plugins/sentinel/internal/processor"
	"github.com/adrianpk/watchman/plugins/sentinel/internal/standards"
	"github.com/adrianpk/watchman/plugins/sentinel/internal/types"
)

// protectedPaths are paths that Sentinel needs protected from agent access.
var protectedPaths = []string{
	"~/.config/sentinel/",
	".sentinel.yml",
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--protected-paths" {
		json.NewEncoder(os.Stdout).Encode(protectedPaths)
		return
	}

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "sentinel: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// DEBUG: Log that sentinel was called
	if f, err := os.OpenFile("/tmp/sentinel-debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		f.WriteString("sentinel called at " + time.Now().Format(time.RFC3339) + "\n")
		f.Close()
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	aiClient, err := client.New(cfg)
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	loader := standards.NewFileLoader(cfg.Standards.File, cfg.Standards.CacheTTL)
	eval := evaluator.New(loader, aiClient)
	proc := processor.New(eval, cfg)

	var input types.HookInput
	if err := json.NewDecoder(os.Stdin).Decode(&input); err != nil {
		writeOutput(types.HookOutput{
			Decision: cfg.Evaluation.DefaultDecision,
			Warning:  "cannot decode input: " + err.Error(),
		})
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Evaluation.Timeout)
	defer cancel()

	output, err := proc.Process(ctx, input)
	if err != nil {
		output = types.HookOutput{
			Decision: cfg.Evaluation.DefaultDecision,
			Warning:  "process error: " + err.Error(),
		}
	}

	writeOutput(output)
	return nil
}

func writeOutput(output types.HookOutput) {
	json.NewEncoder(os.Stdout).Encode(output)
}
