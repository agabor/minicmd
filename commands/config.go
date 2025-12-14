package commands

import (
	"fmt"
	"strings"

	"minicmd/config"
)

func HandleConfigCommand(args []string, cfg *config.Config) error {
	if len(args) == 0 {
		fmt.Println("Current configuration:")
		fmt.Printf("  default_provider: %s\n", cfg.DefaultProvider)
		if cfg.AnthropicAPIKey != "" {
			fmt.Printf("  anthropic_api_key: %s\n", strings.Repeat("*", len(cfg.AnthropicAPIKey)))
		} else {
			fmt.Printf("  anthropic_api_key: %s\n", cfg.AnthropicAPIKey)
		}
		if cfg.DeepSeekAPIKey != "" {
			fmt.Printf("  deepseek_api_key: %s\n", strings.Repeat("*", len(cfg.DeepSeekAPIKey)))
		} else {
			fmt.Printf("  deepseek_api_key: %s\n", cfg.DeepSeekAPIKey)
		}
		fmt.Printf("  claude_model: %s\n", cfg.ClaudeModel)
		fmt.Printf("  ollama_url: %s\n", cfg.OllamaURL)
		fmt.Printf("  ollama_model: %s\n", cfg.OllamaModel)
		fmt.Printf("  deepseek_model: %s\n", cfg.DeepSeekModel)
		fmt.Printf("  fim_token: %s\n", cfg.FimToken)
		return nil
	}

	if len(args) == 2 {
		key := args[0]
		value := args[1]

		switch key {
		case "default_provider":
			cfg.DefaultProvider = value
		case "anthropic_api_key":
			cfg.AnthropicAPIKey = value
		case "deepseek_api_key":
			cfg.DeepSeekAPIKey = value
		case "claude_model":
			cfg.ClaudeModel = value
		case "ollama_url":
			cfg.OllamaURL = value
		case "ollama_model":
			cfg.OllamaModel = value
		case "deepseek_model":
			cfg.DeepSeekModel = value
		case "fim_token":
			cfg.FimToken = value
		default:
			return fmt.Errorf("unknown config key '%s'", key)
		}

		if err := cfg.Save(); err != nil {
			return fmt.Errorf("error saving config: %w", err)
		}

		if key == "anthropic_api_key" || key == "deepseek_api_key" {
			fmt.Printf("Set %s to %s\n", key, strings.Repeat("*", len(value)))
		} else {
			fmt.Printf("Set %s to %s\n", key, value)
		}
		return nil
	}

	fmt.Println("Usage:")
	fmt.Println("  minicmd config                    # Show current config")
	fmt.Println("  minicmd config <key> <value>      # Set config value")
	return nil
}