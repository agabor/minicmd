package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	OllamaURL     = "http://localhost:11434/api/generate"
	OllamaModel   = "deepseek-coder-v2:16b"
	ClaudeModel   = "claude-sonnet-4-20250514"
	DeepSeekURL   = "https://api.deepseek.com/v1/chat/completions"
	DeepSeekModel = "deepseek-coder"
	SystemPrompt  = "IMPORTANT: answer with one or more code blocks only without explanation. The first line should be a comment containing the file path and name. " +
		"When updating an existing source file, leave comments, identation and white spaces unchanged. Always respond with the complete file content. " +
		"Code blocks should always be delimited by triple backticks (```). Do not use any other formatting or text outside of code blocks. " +
		"Each file content should be placed in a separate code block. If you need to delete a file return a bash script named minicmd_rm.sh with the necesarry rm commands."
)

type Config struct {
	DefaultProvider   string `json:"default_provider"`
	AnthropicAPIKey   string `json:"anthropic_api_key"`
	DeepSeekAPIKey    string `json:"deepseek_api_key"`
	OllamaURL         string `json:"ollama_url"`
	OllamaModel       string `json:"ollama_model"`
	ClaudeModel       string `json:"claude_model"`
	DeepSeekURL       string `json:"deepseek_url"`
	DeepSeekModel     string `json:"deepseek_model"`
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
