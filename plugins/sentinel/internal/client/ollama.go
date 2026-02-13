package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
- "reason": explanation if deny (empty string if allow)
- "warning": advisory note if advise (empty string otherwise)
- "violations": array of specific standard violations found

Example response:
{"decision": "deny", "reason": "Missing doc comment on exported function", "warning": "", "violations": ["Exported function Foo has no doc comment"]}

Respond ONLY with the JSON object, no other text.`, req.Standards, req.ToolName, req.FilePath, req.Content)
}

func parseOllamaResponse(response string) (types.EvalResult, error) {
	var result struct {
		Decision   string   `json:"decision"`
		Reason     string   `json:"reason"`
		Warning    string   `json:"warning"`
		Violations []string `json:"violations"`
	}

	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return types.EvalResult{}, fmt.Errorf("cannot parse ollama response: %w (response: %s)", err, response)
	}

	if result.Decision != "allow" && result.Decision != "advise" && result.Decision != "deny" {
		return types.EvalResult{}, fmt.Errorf("invalid decision: %s", result.Decision)
	}

	return types.EvalResult{
		Decision:   result.Decision,
		Reason:     result.Reason,
		Warning:    result.Warning,
		Violations: result.Violations,
	}, nil
}
