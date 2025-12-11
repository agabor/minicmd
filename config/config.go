package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	OllamaURL           = "http://localhost:11434/api/generate"
	OllamaModel         = "qwen2.5-coder:7b"
	ClaudeModel         = "claude-haiku-4-5-20251001"
	DeepSeekURL         = "https://api.deepseek.com/v1/chat/completions"
	DeepSeekModel       = "deepseek-coder"
	DefaultMaxTokens    = 8192
	DefaultFimToken     = "//FIM"
)

type Config struct {
	DefaultProvider  string `json:"default_provider"`
	AnthropicAPIKey  string `json:"anthropic_api_key"`
	DeepSeekAPIKey   string `json:"deepseek_api_key"`
	OllamaURL        string `json:"ollama_url"`
	OllamaModel      string `json:"ollama_model"`
	ClaudeModel      string `json:"claude_model"`
	DeepSeekURL      string `json:"deepseek_url"`
	DeepSeekModel    string `json:"deepseek_model"`
	MaxOutputTokens  int    `json:"max_output_tokens"`
	FimToken         string `json:"fim_token"`
}

func getConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".minicmd"), nil
}

func getConfigFile() (string, error) {
	dir, err := getConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config"), nil
}

func DefaultConfig() *Config {
	return &Config{
		DefaultProvider: "ollama",
		AnthropicAPIKey: "",
		DeepSeekAPIKey:  "",
		OllamaURL:       OllamaURL,
		OllamaModel:     OllamaModel,
		ClaudeModel:     ClaudeModel,
		DeepSeekURL:     DeepSeekURL,
		DeepSeekModel:   DeepSeekModel,
		MaxOutputTokens: DefaultMaxTokens,
		FimToken:        DefaultFimToken,
	}
}

func Load() (*Config, error) {
	cfg := DefaultConfig()

	configFile, err := getConfigFile()
	if err != nil {
		return cfg, err
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return cfg, err
	}

	if cfg.MaxOutputTokens <= 0 {
		cfg.MaxOutputTokens = DefaultMaxTokens
	}

	if cfg.FimToken == "" {
		cfg.FimToken = DefaultFimToken
	}

	return cfg, nil
}

func (c *Config) Save() error {
	configFile, err := getConfigFile()
	if err != nil {
		return err
	}

	dir, err := getConfigDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configFile, data, 0644)
}