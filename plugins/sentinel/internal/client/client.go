package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/adrianpk/watchman/plugins/sentinel/internal/config"
	"github.com/adrianpk/watchman/plugins/sentinel/internal/interfaces"
	"github.com/adrianpk/watchman/plugins/sentinel/internal/types"
)

// New creates an AIClient based on config.
// If multiple providers are configured, returns a FallbackClient.
func New(cfg *config.Config) (interfaces.AIClient, error) {
	providers := cfg.GetProviders()

	if len(providers) == 1 {
		return newSingleClient(cfg, providers[0])
	}

	return newFallbackClient(cfg, providers)
}

func newSingleClient(cfg *config.Config, provider string) (interfaces.AIClient, error) {
	switch provider {
	case "anthropic":
		return NewAnthropicClient(
			cfg.Anthropic.APIKey,
			cfg.Anthropic.Model,
			cfg.Anthropic.MaxTokens,
		), nil
	case "openai":
		return NewOpenAIClient(
			cfg.OpenAI.APIKey,
			cfg.OpenAI.Model,
			cfg.OpenAI.MaxTokens,
		), nil
	case "ollama":
		return NewOllamaClient(
			cfg.Ollama.Host,
			cfg.Ollama.Model,
		), nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", provider)
	}
}

func newFallbackClient(cfg *config.Config, providers []string) (*FallbackClient, error) {
	var clients []namedClient
	for _, p := range providers {
		c, err := newSingleClient(cfg, p)
		if err != nil {
			return nil, fmt.Errorf("cannot create %s client: %w", p, err)
		}
		clients = append(clients, namedClient{name: p, client: c})
	}
	return &FallbackClient{clients: clients}, nil
}

type namedClient struct {
	name   string
	client interfaces.AIClient
}

// FallbackClient tries multiple providers in order until one succeeds.
type FallbackClient struct {
	clients []namedClient
}

// Evaluate tries each provider in order, returning the first successful result.
func (f *FallbackClient) Evaluate(ctx context.Context, req types.EvalRequest) (types.EvalResult, error) {
	var errs []string

	for _, nc := range f.clients {
		result, err := nc.client.Evaluate(ctx, req)
		if err == nil {
			return result, nil
		}
		errs = append(errs, fmt.Sprintf("%s: %v", nc.name, err))
	}

	return types.EvalResult{}, fmt.Errorf("all providers failed: %s", strings.Join(errs, "; "))
}
