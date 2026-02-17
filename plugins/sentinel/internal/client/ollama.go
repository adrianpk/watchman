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
	return fmt.Sprintf(`You are a code standards evaluator. Evaluate the action against the provided standards.

## Standards

%s

## Action to Evaluate

Tool: %s
File: %s

Content:
%s

Respond with a JSON object containing:
- "decision": one of "allow", "advise", or "deny"
  - "allow" = code is compliant with standards
  - "advise" = minor issues, warning only
  - "deny" = violates standards, must be blocked
- "warning": advisory note if advise (empty string otherwise)
- "violations": array of violation objects, each with:
  - "rule": name of the violated rule (e.g., "Comment violation", "Test naming")
  - "file": filename where violation occurs
  - "lines": line number or range (e.g., "3", "3-4", "N/A")
  - "detail": specific explanation of what is wrong and how to fix it

Example response:
{"decision": "deny", "warning": "", "violations": [{"rule": "Comment violation", "file": "namespace.go", "lines": "3-4", "detail": "Multi-line comment on exported type may not add value per standards"}]}

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
