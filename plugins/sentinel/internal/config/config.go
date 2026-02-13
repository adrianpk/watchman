package config

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Provider   string          `yaml:"provider"`
	Anthropic  AnthropicConfig `yaml:"anthropic"`
	OpenAI     OpenAIConfig    `yaml:"openai"`
	Ollama     OllamaConfig    `yaml:"ollama"`
	Standards  StandardsConfig `yaml:"standards"`
	Evaluation EvalConfig      `yaml:"evaluation"`
}

type AnthropicConfig struct {
	APIKey    string `yaml:"api_key"`
	Model     string `yaml:"model"`
	MaxTokens int    `yaml:"max_tokens"`
}

type OpenAIConfig struct {
	APIKey    string `yaml:"api_key"`
	Model     string `yaml:"model"`
	MaxTokens int    `yaml:"max_tokens"`
}

type OllamaConfig struct {
	Host  string `yaml:"host"`
	Model string `yaml:"model"`
}

type StandardsConfig struct {
	File     string        `yaml:"file"`
	CacheTTL time.Duration `yaml:"cache_ttl"`
}

type EvalConfig struct {
	DefaultDecision string `yaml:"default_decision"`
	MaxContentSize  int    `yaml:"max_content_size"`
	Timeout         time.Duration `yaml:"timeout"`
}

func Default() *Config {
	return &Config{
		Provider: "anthropic",
		Anthropic: AnthropicConfig{
			Model:     "claude-sonnet-4-20250514",
			MaxTokens: 1024,
		},
		Standards: StandardsConfig{
			File:     "AGENTS.md",
			CacheTTL: 5 * time.Minute,
		},
		Evaluation: EvalConfig{
			DefaultDecision: "allow",
			MaxContentSize:  50000,
			Timeout:         25 * time.Second,
		},
	}
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
