package evaluator

import (
	"context"

	"github.com/adrianpk/watchman/plugins/sentinel/internal/interfaces"
	"github.com/adrianpk/watchman/plugins/sentinel/internal/types"
)

type Evaluator struct {
	loader interfaces.StandardsLoader
	client interfaces.AIClient
}

func New(loader interfaces.StandardsLoader, client interfaces.AIClient) *Evaluator {
	return &Evaluator{
		loader: loader,
		client: client,
	}
}

func (e *Evaluator) Evaluate(ctx context.Context, input types.HookInput) (types.HookOutput, error) {
	standards, err := e.loader.Load(ctx)
	if err != nil {
		return types.HookOutput{}, err
	}

	filePath, _ := input.ToolInput["file_path"].(string)
	content, _ := input.ToolInput["content"].(string)
	if content == "" {
		if cmd, ok := input.ToolInput["command"].(string); ok {
			content = cmd
		}
	}

	req := types.EvalRequest{
		ToolName:  input.ToolName,
		FilePath:  filePath,
		Content:   content,
		Standards: standards,
	}

	result, err := e.client.Evaluate(ctx, req)
	if err != nil {
		return types.HookOutput{}, err
	}

	return types.HookOutput{
		Decision: result.Decision,
		Reason:   result.Reason,
		Warning:  result.Warning,
	}, nil
}
