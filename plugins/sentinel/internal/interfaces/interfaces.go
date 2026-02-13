package interfaces

import (
	"context"

	"github.com/adrianpk/watchman/plugins/sentinel/internal/types"
)

type StandardsLoader interface {
	Load(ctx context.Context) (string, error)
	Reload(ctx context.Context) (string, error)
}

type AIClient interface {
	Evaluate(ctx context.Context, req types.EvalRequest) (types.EvalResult, error)
}

type Evaluator interface {
	Evaluate(ctx context.Context, input types.HookInput) (types.HookOutput, error)
}

type Processor interface {
	Process(ctx context.Context, input types.HookInput) (types.HookOutput, error)
}
