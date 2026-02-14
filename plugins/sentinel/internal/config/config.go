package config

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the complete sentinel configuration.
type Config struct {
	Provider   string          `yaml:"provider"`
	Providers  []string        `yaml:"providers"`
	Anthropic  AnthropicConfig `yaml:"anthropic"`
	OpenAI     OpenAIConfig    `yaml:"openai"`
	Ollama     OllamaConfig    `yaml:"ollama"`
	Standards  StandardsConfig `yaml:"standards"`
	Evaluation EvalConfig      `yaml:"evaluation"`
}

// AnthropicConfig holds Anthropic API settings.
type AnthropicConfig struct {
	APIKey    string `yaml:"api_key"`
	Model     string `yaml:"model"`
	MaxTokens int    `yaml:"max_tokens"`
}

// OpenAIConfig holds OpenAI API settings.
type OpenAIConfig struct {
	APIKey    string `yaml:"api_key"`
	Model     string `yaml:"model"`
	MaxTokens int    `yaml:"max_tokens"`
}

// OllamaConfig holds Ollama local inference settings.
type OllamaConfig struct {
	Host  string `yaml:"host"`
	Model string `yaml:"model"`
}

// StandardsConfig holds standards file location and caching settings.
type StandardsConfig struct {
	File     string        `yaml:"file"`
	CacheTTL time.Duration `yaml:"cache_ttl"`
}

// EvalConfig holds evaluation behavior settings.
type EvalConfig struct {
	DefaultDecision string        `yaml:"default_decision"`
	MaxContentSize  int           `yaml:"max_content_size"`
	Timeout         time.Duration `yaml:"timeout"`
	Mode            string        `yaml:"mode"`
	Batch           BatchConfig   `yaml:"batch"`
}

// BatchConfig controls how staged diffs are batched for evaluation.
type BatchConfig struct {
	MaxFilesPerRequest   int `yaml:"max_files_per_request"`
	MaxContentPerRequest int `yaml:"max_content_per_request"`
	Concurrency          int `yaml:"concurrency"`
}

func Default() *Config {
	return &Config{
		Provider:  "anthropic",
		Providers: nil,
		Anthropic: AnthropicConfig{
			Model:     "claude-sonnet-4-20250514",
			MaxTokens: 1024,
		},
		OpenAI: OpenAIConfig{
			Model:     "gpt-4o-mini",
			MaxTokens: 1024,
		},
		Ollama: OllamaConfig{
			Host:  "http://localhost:11434",
			Model: "llama3",
		},
		Standards: StandardsConfig{
			File:     "AGENTS.md",
			CacheTTL: 5 * time.Minute,
		},
		Evaluation: EvalConfig{
			DefaultDecision: "allow",
			MaxContentSize:  50000,
			Timeout:         25 * time.Second,
			Mode:            "all",
			Batch: BatchConfig{
				MaxFilesPerRequest:   10,
				MaxContentPerRequest: 40000,
				Concurrency:          3,
			},
		},
	}
}

// GetProviders returns the list of providers to try in order.
// If Providers is set, use that. Otherwise, use single Provider.
func (c *Config) GetProviders() []string {
	if len(c.Providers) > 0 {
		return c.Providers
	}
	return []string{c.Provider}
}

func Load() (*Config, error) {
	cfg := Default()

	paths := []string{
		configPath(),
		".sentinel.yml",
	}

	for _, p := range paths {
		if p == "" {
			continue
		}
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
		break
	}

	cfg.expandEnv()
	return cfg, nil
}

func (c *Config) expandEnv() {
	c.Anthropic.APIKey = os.ExpandEnv(c.Anthropic.APIKey)
	c.OpenAI.APIKey = os.ExpandEnv(c.OpenAI.APIKey)
}

func configPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "sentinel", "config.yml")
}
