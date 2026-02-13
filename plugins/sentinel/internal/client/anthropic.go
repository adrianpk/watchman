package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/adrianpk/watchman/plugins/sentinel/internal/types"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type AnthropicClient struct {
	client    anthropic.Client
	model     string
	maxTokens int
}

func NewAnthropicClient(apiKey, model string, maxTokens int) *AnthropicClient {
	client := anthropic.NewClient(option.WithAPIKey(apiKey))
	return &AnthropicClient{
		client:    client,
		model:     model,
		maxTokens: maxTokens,
	}
}

func (a *AnthropicClient) Evaluate(ctx context.Context, req types.EvalRequest) (types.EvalResult, error) {
	tool := anthropic.ToolUnionParamOfTool(evaluationSchema(), "evaluate_action")
	tool.OfTool.Description = anthropic.String("Report evaluation result for a code action against standards")

	prompt := buildPrompt(req)

	msg, err := a.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     a.model,
		MaxTokens: int64(a.maxTokens),
		System: []anthropic.TextBlockParam{
			{Text: "You are a code standards evaluator. Evaluate the action against the provided standards. You MUST use the evaluate_action tool to respond."},
		},
		Tools:      []anthropic.ToolUnionParam{tool},
		ToolChoice: anthropic.ToolChoiceParamOfToolChoiceTool("evaluate_action"),
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})
	if err != nil {
		return types.EvalResult{}, fmt.Errorf("anthropic api: %w", err)
	}

	return parseResponse(msg)
}

func evaluationSchema() anthropic.ToolInputSchemaParam {
	return anthropic.ToolInputSchemaParam{
		Properties: map[string]interface{}{
			"decision": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"allow", "advise", "deny"},
				"description": "allow=compliant, advise=minor issues, deny=violates standards",
			},
			"reason": map[string]interface{}{
				"type":        "string",
				"description": "Explanation if deny. Empty if allow.",
			},
			"warning": map[string]interface{}{
				"type":        "string",
				"description": "Advisory note if advise. Empty otherwise.",
			},
			"violations": map[string]interface{}{
				"type":        "array",
				"items":       map[string]string{"type": "string"},
				"description": "List of specific standard violations found",
			},
		},
		ExtraFields: map[string]interface{}{
			"required": []string{"decision"},
		},
	}
}

func buildPrompt(req types.EvalRequest) string {
	return fmt.Sprintf(`## Standards

%s

## Action to Evaluate

Tool: %s
File: %s

Content:
%s

Evaluate this action against the standards.`, req.Standards, req.ToolName, req.FilePath, req.Content)
}

func parseResponse(msg *anthropic.Message) (types.EvalResult, error) {
	for _, block := range msg.Content {
		if block.Type == "tool_use" {
			var result struct {
				Decision   string   `json:"decision"`
				Reason     string   `json:"reason"`
				Warning    string   `json:"warning"`
				Violations []string `json:"violations"`
			}
			if err := json.Unmarshal(block.Input, &result); err != nil {
				return types.EvalResult{}, fmt.Errorf("cannot parse tool response: %w", err)
			}
			return types.EvalResult{
				Decision:   result.Decision,
				Reason:     result.Reason,
				Warning:    result.Warning,
				Violations: result.Violations,
			}, nil
		}
	}
	return types.EvalResult{}, fmt.Errorf("error: no tool use block in response")
}
