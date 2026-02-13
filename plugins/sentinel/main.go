package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/adrianpk/watchman/plugins/sentinel/internal/client"
	"github.com/adrianpk/watchman/plugins/sentinel/internal/config"
	"github.com/adrianpk/watchman/plugins/sentinel/internal/evaluator"
	"github.com/adrianpk/watchman/plugins/sentinel/internal/processor"
	"github.com/adrianpk/watchman/plugins/sentinel/internal/standards"
	"github.com/adrianpk/watchman/plugins/sentinel/internal/types"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "sentinel: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
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
	proc := processor.New(eval, cfg.Evaluation.DefaultDecision)

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
