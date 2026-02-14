// Package processor handles the main evaluation logic for sentinel.
package processor

import (
	"context"

	"github.com/adrianpk/watchman/plugins/sentinel/internal/config"
	"github.com/adrianpk/watchman/plugins/sentinel/internal/git"
	"github.com/adrianpk/watchman/plugins/sentinel/internal/interfaces"
	"github.com/adrianpk/watchman/plugins/sentinel/internal/types"
)

// Processor orchestrates code evaluation against standards.
type Processor struct {
	evaluator interfaces.Evaluator
	cfg       *config.Config
}

// New creates a new Processor with the given evaluator and config.
func New(evaluator interfaces.Evaluator, cfg *config.Config) *Processor {
	return &Processor{
		evaluator: evaluator,
		cfg:       cfg,
	}
}

// Process evaluates the input based on the configured mode.
func (p *Processor) Process(ctx context.Context, input types.HookInput) (types.HookOutput, error) {
	if p.cfg.Evaluation.Mode == "commits_only" {
		return p.processCommitsOnly(ctx, input)
	}
	return p.processAll(ctx, input)
}

// processAll evaluates every Write/Edit operation.
func (p *Processor) processAll(ctx context.Context, input types.HookInput) (types.HookOutput, error) {
	output, err := p.evaluator.Evaluate(ctx, input)
	if err != nil {
		return types.HookOutput{
			Decision: p.cfg.Evaluation.DefaultDecision,
			Warning:  "evaluation error: " + err.Error(),
		}, nil
	}
	return output, nil
}

// processCommitsOnly evaluates on git add/commit, skipping other operations.
func (p *Processor) processCommitsOnly(ctx context.Context, input types.HookInput) (types.HookOutput, error) {
	if input.ToolName != "Bash" {
		return types.HookOutput{Decision: "allow"}, nil
	}

	command, ok := input.ToolInput["command"].(string)
	if !ok {
		return types.HookOutput{Decision: "allow"}, nil
	}

	// Handle git add: parse files from command and evaluate them directly
	if git.IsAddCommand(command) {
		return p.processGitAdd(ctx, input, command)
	}

	// Handle git commit: evaluate staged diff
	if git.IsCommitCommand(command) {
		return p.processGitCommit(ctx, input)
	}

	return types.HookOutput{Decision: "allow"}, nil
}

// processGitAdd extracts files from git add command and evaluates them.
func (p *Processor) processGitAdd(ctx context.Context, input types.HookInput, command string) (types.HookOutput, error) {
	files := git.ExtractAddFiles(command)
	if len(files) == 0 {
		return types.HookOutput{Decision: "allow"}, nil
	}

	content, err := git.ReadFiles(input.WorkingDir, files)
	if err != nil {
		return types.HookOutput{
			Decision: p.cfg.Evaluation.DefaultDecision,
			Warning:  "failed to read files: " + err.Error(),
		}, nil
	}

	addInput := types.HookInput{
		ToolName:   "GitAdd",
		ToolInput:  map[string]any{"content": content, "files": files},
		Paths:      files,
		WorkingDir: input.WorkingDir,
	}

	output, err := p.evaluator.Evaluate(ctx, addInput)
	if err != nil {
		return types.HookOutput{
			Decision: p.cfg.Evaluation.DefaultDecision,
			Warning:  "evaluation error: " + err.Error(),
		}, nil
	}
	return output, nil
}

// processGitCommit evaluates the staged diff.
func (p *Processor) processGitCommit(ctx context.Context, input types.HookInput) (types.HookOutput, error) {
	diff, err := git.GetStagedDiff(input.WorkingDir)
	if err != nil {
		return types.HookOutput{
			Decision: p.cfg.Evaluation.DefaultDecision,
			Warning:  "failed to get staged diff: " + err.Error(),
		}, nil
	}

	if len(diff.Files) == 0 {
		return types.HookOutput{Decision: "allow"}, nil
	}

	commitInput := types.HookInput{
		ToolName:   "Commit",
		ToolInput:  map[string]any{"content": diff.Content, "files": diff.Files},
		Paths:      diff.Files,
		WorkingDir: input.WorkingDir,
	}

	output, err := p.evaluator.Evaluate(ctx, commitInput)
	if err != nil {
		return types.HookOutput{
			Decision: p.cfg.Evaluation.DefaultDecision,
			Warning:  "evaluation error: " + err.Error(),
		}, nil
	}
	return output, nil
}
