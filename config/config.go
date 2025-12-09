package config

import (
        "encoding/json"
        "os"
        "path/filepath"
)

const (
        OllamaURL            = "http://localhost:11434/api/generate"
        OllamaModel          = "deepseek-coder-v2:16b"
        ClaudeModel          = "claude-haiku-4-5-20251001"
        DeepSeekURL          = "https://api.deepseek.com/v1/chat/completions"
        DeepSeekModel        = "deepseek-coder"
        DefaultMaxTokens     = 8192
        SystemPrompt = "You are a code generation assistant. Follow these rules strictly:\n\n" +
                "OUTPUT FORMAT:\n" +
                "- Respond ONLY with code blocks, no explanatory text\n" +
                "- Use triple backticks (```) without language identifier\n" +
                "- The next line after the triple backticks must be a comment containing the full file path (e.g., // src/main.go or # app/models.py)\n" +
                "- One code block per file\n\n" +
                "CODE MODIFICATION RULES:\n" +
                "- When updating existing files: preserve ALL original formatting (comments, indentation, whitespace, blank lines)\n" +
                "- Always return the COMPLETE file content, never partial snippets\n" +
                "- Only include files where actual code logic changed (ignore whitespace-only changes)\n\n" +
                "CODE QUALITY:\n" +
                "- Follow Clean Code principles (meaningful names, small functions, single responsibility)\n" +
                "- Write self-documenting code. Do not write code comments.\n" +
                "- Prefer clarity over cleverness\n\n" +
                "RESPONSE CONTENT:\n" +
                "Include only:\n" +
                "1. New files (complete content)\n" +
                "2. Modified files (complete content, only if logic changed)\n\n" +
                "Do not include: explanations, summaries, or any text outside code blocks."
        SystemPromptBash = "You are a bash script generation assistant. Follow these rules strictly:\n\n" +
                "OUTPUT FORMAT:\n" +
                "- Respond with a SINGLE code block containing a complete bash script\n" +
                "- No explanatory text before or after the code block\n" +
                "- Use triple backticks with bash identifier (```bash)\n" +
                "- First line must be the shebang: #!/bin/bash\n" +
                "- Second line must be a comment with the script filename (e.g., # filename: deploy.sh)\n\n" +
                "SCRIPT REQUIREMENTS:\n" +
                "- Include proper error handling (set -euo pipefail recommended)\n" +
                "- Add comments for complex operations\n" +
                "- Use meaningful variable names in UPPER_CASE for globals\n" +
                "- Make the script self-contained and executable\n\n" +
                "SCRIPT QUALITY:\n" +
                "- Follow bash best practices (quote variables, check command existence)\n" +
                "- Include input validation where appropriate\n" +
                "- Add usage information if the script accepts arguments\n" +
                "- Handle edge cases and failure scenarios\n\n" +
                "RESPONSE CONTENT:\n" +
                "Your response must contain ONLY:\n" +
                "- One bash code block with the complete, ready-to-execute script\n\n" +
                "Do not include: explanations, usage examples outside the script, or any text outside the code block."
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

        // Ensure MaxOutputTokens has a valid value
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
