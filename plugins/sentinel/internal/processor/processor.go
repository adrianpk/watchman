package processor

import (
	"context"

	"github.com/adrianpk/watchman/plugins/sentinel/internal/interfaces"
	"github.com/adrianpk/watchman/plugins/sentinel/internal/types"
)

type Processor struct {
	evaluator       interfaces.Evaluator
	defaultDecision string
}

func New(evaluator interfaces.Evaluator, defaultDecision string) *Processor {
	return &Processor{
		evaluator:       evaluator,
		defaultDecision: defaultDecision,
	}
}

func (p *Processor) Process(ctx context.Context, input types.HookInput) (types.HookOutput, error) {
	output, err := p.evaluator.Evaluate(ctx, input)
	if err != nil {
		return types.HookOutput{
			Decision: p.defaultDecision,
			Warning:  "evaluation error: " + err.Error(),
		}, nil
	}
	return output, nil
}
