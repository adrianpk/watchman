package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/adrianpk/watchman/plugins/sentinel/internal/types"
)

// OllamaClient implements AIClient using a local Ollama instance.
type OllamaClient struct {
	host  string
	model string
}

// NewOllamaClient creates a new Ollama client.
func NewOllamaClient(host, model string) *OllamaClient {
	if host == "" {
		host = "http://localhost:11434"
	}
	return &OllamaClient{
		host:  host,
		model: model,
	}
}

// Evaluate evaluates the given request against standards using Ollama.
func (o *OllamaClient) Evaluate(ctx context.Context, req types.EvalRequest) (types.EvalResult, error) {
	prompt := buildOllamaPrompt(req)

	ollamaReq := ollamaRequest{
		Model:  o.model,
		Prompt: prompt,
		Stream: false,
		Format: "json",
	}

	body, err := json.Marshal(ollamaReq)
	if err != nil {
		return types.EvalResult{}, fmt.Errorf("cannot marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", o.host+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return types.EvalResult{}, fmt.Errorf("cannot create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return types.EvalResult{}, fmt.Errorf("ollama api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return types.EvalResult{}, fmt.Errorf("ollama api: %d %s", resp.StatusCode, string(respBody))
	}

	var ollamaResp ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return types.EvalResult{}, fmt.Errorf("cannot decode response: %w", err)
	}

	return parseOllamaResponse(ollamaResp.Response)
}

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	Format string `json:"format"`
}

type ollamaResponse struct {
	Response string `json:"response"`
}

func buildOllamaPrompt(req types.EvalRequest) string {
	return fmt.Sprintf(`You are a conservative code standards evaluator.

DECISION GUIDE:
- DENY: Only for clear, explicit violations (confidence >= 0.9)
- ADVISE: For ambiguous cases worth reviewing (confidence 0.5-0.9)
- ALLOW: When compliant or rule doesn't apply (confidence < 0.5)

CRITICAL RULES:
- Never deny based on interpretation - use advise instead
- Never extrapolate rules beyond their literal meaning
- The standards document is the ONLY source of truth
- When in doubt, advise (warn) rather than deny (block)

## Standards

%s

## Action to Evaluate

Tool: %s
File: %s

Content:
%s

Respond with a JSON object containing:
- "decision": one of "allow", "advise", or "deny"
  - "allow" = compliant or rule doesn't apply (DEFAULT)
  - "advise" = ambiguous, worth reviewing but don't block
  - "deny" = ONLY for clear, explicit violations
- "warning": advisory note if advise (empty string otherwise)
- "violations": array of violation objects, each with:
  - "rule": name of the violated rule
  - "file": filename where violation occurs
  - "lines": line number or range (e.g., "3", "3-4", "N/A")
  - "detail": specific explanation
  - "confidence": number 0.0-1.0 (1.0=explicit violation, 0.7=likely, 0.5=ambiguous, <0.5=not a violation)

Example response:
{"decision": "advise", "warning": "Consider reviewing", "violations": [{"rule": "Naming", "file": "foo.go", "lines": "10", "detail": "Variable name could be clearer", "confidence": 0.6}]}

Respond ONLY with the JSON object, no other text.`, req.Standards, req.ToolName, req.FilePath, req.Content)
}

func parseOllamaResponse(response string) (types.EvalResult, error) {
	var result struct {
		Decision   string            `json:"decision"`
		Warning    string            `json:"warning"`
		Violations []types.Violation `json:"violations"`
	}

	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return types.EvalResult{}, fmt.Errorf("cannot parse ollama response: %w (response: %s)", err, response)
	}

	if result.Decision != "allow" && result.Decision != "advise" && result.Decision != "deny" {
		return types.EvalResult{}, fmt.Errorf("invalid decision: %s", result.Decision)
	}

	return types.EvalResult{
		Decision:   result.Decision,
		Reason:     ollamaBuildReason(result.Violations),
		Warning:    result.Warning,
		Violations: result.Violations,
	}, nil
}

func ollamaBuildReason(violations []types.Violation) string {
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
