package client

import (
	"fmt"

	"github.com/adrianpk/watchman/plugins/sentinel/internal/config"
	"github.com/adrianpk/watchman/plugins/sentinel/internal/interfaces"
)

func New(cfg *config.Config) (interfaces.AIClient, error) {
	switch cfg.Provider {
	case "anthropic":
		return NewAnthropicClient(
			cfg.Anthropic.APIKey,
			cfg.Anthropic.Model,
			cfg.Anthropic.MaxTokens,
		), nil
	case "openai":
		return nil, fmt.Errorf("openai provider not yet implemented")
	case "ollama":
		return nil, fmt.Errorf("ollama provider not yet implemented")
	default:
		return nil, fmt.Errorf("unknown provider: %s", cfg.Provider)
	}
}
