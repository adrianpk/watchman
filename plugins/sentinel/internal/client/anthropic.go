package client

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

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
			"warning": map[string]interface{}{
				"type":        "string",
				"description": "Advisory note if advise. Empty otherwise.",
			},
			"violations": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"rule": map[string]interface{}{
							"type":        "string",
							"description": "Name of the violated rule from standards (e.g., 'Comment violation', 'Test naming', 'Commit message')",
						},
						"file": map[string]interface{}{
							"type":        "string",
							"description": "Filename where violation occurs (e.g., 'namespace.go')",
						},
						"lines": map[string]interface{}{
							"type":        "string",
							"description": "Line number or range (e.g., '3', '3-4', 'N/A' for commit messages)",
						},
						"detail": map[string]interface{}{
							"type":        "string",
							"description": "Specific explanation of what is wrong and how to fix it",
						},
					},
					"required": []string{"rule", "file", "lines", "detail"},
				},
				"description": "Structured list of violations. Required if decision is deny or advise.",
			},
		},
		ExtraFields: map[string]interface{}{
			"required": []string{"decision", "violations"},
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

Evaluate this action against the standards.

## Response Format Requirements

The reason field MUST be actionable. The agent receiving this needs to:
1. Understand which specific rule was violated
2. Locate exactly where (file:line)
3. Decide whether to fix or report as false positive

### Single violation format:
[RULE_NAME] in [FILE:LINE]. [SPECIFIC_DETAIL]

Example: "Comment violation in namespace.go:3-4. Multi-line comment on exported type may not add value per CONVENTIONS.md"

### Multiple violations format:
Multiple violations found:
1. [RULE] in [FILE:LINE] - [detail]
2. [RULE] in [FILE:LINE] - [detail]

### If ambiguous:
Indicate uncertainty so the agent can decide:
"Possible violation in namespace.go:3-4. Comment may or may not add value - REVIEW RECOMMENDED"

### Bad examples (DO NOT use these vague formats):
- "The action violates several code standards"
- "Code style violation detected"
- "Does not follow conventions"

Be specific. Reference the exact section of standards violated.`, req.Standards, req.ToolName, req.FilePath, req.Content)
}

func parseResponse(msg *anthropic.Message) (types.EvalResult, error) {
	for _, block := range msg.Content {
		if block.Type == "tool_use" {
			var result struct {
				Decision   string            `json:"decision"`
				Warning    string            `json:"warning"`
				Violations []types.Violation `json:"violations"`
			}
			if err := json.Unmarshal(block.Input, &result); err != nil {
				return types.EvalResult{}, fmt.Errorf("cannot parse tool response: %w", err)
			}
			return types.EvalResult{
				Decision:   result.Decision,
				Reason:     buildReason(result.Violations),
				Warning:    result.Warning,
				Violations: result.Violations,
			}, nil
		}
	}
	return types.EvalResult{}, fmt.Errorf("error: no tool use block in response")
}

func buildReason(violations []types.Violation) string {
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
