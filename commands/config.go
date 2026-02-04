package commands

import (
	"fmt"
	"strings"

	"yact/config"
)

func HandleConfigCommand(args []string, cfg *config.Config) error {
	if len(args) == 0 {
		fmt.Println("Current configuration:")
		if cfg.AnthropicAPIKey != "" {
			fmt.Printf("  anthropic_api_key: %s\n", strings.Repeat("*", len(cfg.AnthropicAPIKey)))
		} else {
			fmt.Printf("  anthropic_api_key: %s\n", cfg.AnthropicAPIKey)
		}
		fmt.Printf("  claude_model: %s\n", cfg.ClaudeModel)
		return nil
	}

	if len(args) == 2 {
		key := args[0]
		value := args[1]

		switch key {
		case "anthropic_api_key":
			cfg.AnthropicAPIKey = value
		case "claude_model":
			cfg.ClaudeModel = value
		default:
			return fmt.Errorf("unknown config key '%s'", key)
		}

		if err := cfg.Save(); err != nil {
			return fmt.Errorf("error saving config: %w", err)
		}

		if key == "anthropic_api_key" {
			fmt.Printf("Set %s to %s\n", key, strings.Repeat("*", len(value)))
		} else {
			fmt.Printf("Set %s to %s\n", key, value)
		}
		return nil
	}

	fmt.Println("Usage:")
	fmt.Println("  ya config                    # Show current config")
	fmt.Println("  ya config <key> <value>      # Set config value")
	return nil
}
