package client

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/adrianpk/watchman/plugins/sentinel/internal/types"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// OpenAIClient implements AIClient using the OpenAI API.
type OpenAIClient struct {
	client    openai.Client
	model     string
	maxTokens int
}

// NewOpenAIClient creates a new OpenAI client.
func NewOpenAIClient(apiKey, model string, maxTokens int) *OpenAIClient {
	client := openai.NewClient(option.WithAPIKey(apiKey))
	return &OpenAIClient{
		client:    client,
		model:     model,
		maxTokens: maxTokens,
	}
}

// Evaluate evaluates the given request against standards using OpenAI.
func (o *OpenAIClient) Evaluate(ctx context.Context, req types.EvalRequest) (types.EvalResult, error) {
	tool := openai.ChatCompletionToolParam{
		Type: "function",
		Function: openai.FunctionDefinitionParam{
			Name:        "evaluate_action",
			Description: openai.String("Report evaluation result for a code action against standards"),
			Parameters:  openaiEvaluationSchema(),
		},
	}

	prompt := buildPrompt(req)

	resp, err := o.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:     o.model,
		MaxTokens: openai.Int(int64(o.maxTokens)),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a code standards evaluator. Evaluate the action against the provided standards. You MUST use the evaluate_action function to respond."),
			openai.UserMessage(prompt),
		},
		Tools: []openai.ChatCompletionToolParam{tool},
		ToolChoice: openai.ChatCompletionToolChoiceOptionUnionParam{
			OfChatCompletionNamedToolChoice: &openai.ChatCompletionNamedToolChoiceParam{
				Type: "function",
				Function: openai.ChatCompletionNamedToolChoiceFunctionParam{
					Name: "evaluate_action",
				},
			},
		},
	})
	if err != nil {
		return types.EvalResult{}, fmt.Errorf("openai api: %w", err)
	}

	return parseOpenAIResponse(resp)
}

func openaiEvaluationSchema() openai.FunctionParameters {
	return openai.FunctionParameters{
		"type": "object",
		"properties": map[string]any{
			"decision": map[string]any{
				"type":        "string",
				"enum":        []string{"allow", "advise", "deny"},
				"description": "allow=compliant, advise=minor issues, deny=violates standards",
			},
			"warning": map[string]any{
				"type":        "string",
				"description": "Advisory note if advise. Empty otherwise.",
			},
			"violations": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"rule": map[string]any{
							"type":        "string",
							"description": "Name of the violated rule from standards",
						},
						"file": map[string]any{
							"type":        "string",
							"description": "Filename where violation occurs",
						},
						"lines": map[string]any{
							"type":        "string",
							"description": "Line number or range (e.g., '3', '3-4', 'N/A')",
						},
						"detail": map[string]any{
							"type":        "string",
							"description": "Specific explanation of what is wrong and how to fix it",
						},
					},
					"required": []string{"rule", "file", "lines", "detail"},
				},
				"description": "Structured list of violations. Required if decision is deny or advise.",
			},
		},
		"required": []string{"decision", "violations"},
	}
}

func parseOpenAIResponse(resp *openai.ChatCompletion) (types.EvalResult, error) {
	if len(resp.Choices) == 0 {
		return types.EvalResult{}, fmt.Errorf("error: no choices in response")
	}

	choice := resp.Choices[0]
	if len(choice.Message.ToolCalls) == 0 {
		return types.EvalResult{}, fmt.Errorf("error: no tool calls in response")
	}

	toolCall := choice.Message.ToolCalls[0]
	if toolCall.Function.Name != "evaluate_action" {
		return types.EvalResult{}, fmt.Errorf("error: unexpected function: %s", toolCall.Function.Name)
	}

	var result struct {
		Decision   string            `json:"decision"`
		Warning    string            `json:"warning"`
		Violations []types.Violation `json:"violations"`
	}
	if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &result); err != nil {
		return types.EvalResult{}, fmt.Errorf("cannot parse function arguments: %w", err)
	}

	return types.EvalResult{
		Decision:   result.Decision,
		Reason:     buildReasonFromViolations(result.Violations),
		Warning:    result.Warning,
		Violations: result.Violations,
	}, nil
}

func buildReasonFromViolations(violations []types.Violation) string {
	if len(violations) == 0 {
		return ""
	}

	if len(violations) == 1 {
		v := violations[0]
		return fmt.Sprintf("%s in %s:%s. %s", v.Rule, v.File, v.Lines, v.Detail)
	}

	var lines []string
	lines = append(lines, "Multiple violations found:")
	for i, v := range violations {
		lines = append(lines, fmt.Sprintf("%d. %s in %s:%s - %s", i+1, v.Rule, v.File, v.Lines, v.Detail))
	}
	return strings.Join(lines, "\n")
}
