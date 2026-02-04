package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	ClaudeModel      = "claude-haiku-4-5-20251001"
	DefaultMaxTokens = 8192
)

type Config struct {
	AnthropicAPIKey string `json:"anthropic_api_key"`
	ClaudeModel     string `json:"claude_model"`
	MaxOutputTokens int    `json:"max_output_tokens"`
}

func getConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".yact"), nil
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
		AnthropicAPIKey: "",
		ClaudeModel:     ClaudeModel,
		MaxOutputTokens: DefaultMaxTokens,
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
