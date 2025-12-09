package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"minicmd/config"
	"minicmd/promptmanager"
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
		fmt.Printf("  deepseek_url: %s\n", cfg.DeepSeekURL)
		fmt.Printf("  deepseek_model: %s\n", cfg.DeepSeekModel)
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
		case "deepseek_url":
			cfg.DeepSeekURL = value
		case "deepseek_model":
			cfg.DeepSeekModel = value
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

func HandleEditCommand() error {
	return promptmanager.EditPromptFile()
}

func HandleAddCommand(args []string) error {
	if len(args) < 1 {
		fmt.Println("Usage: minicmd read <file> [<file2> ...]")
		return fmt.Errorf("missing file argument")
	}

	for _, pattern := range args {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return fmt.Errorf("error matching pattern %s: %w", pattern, err)
		}

		if len(matches) == 0 {
			fmt.Printf("No files found matching pattern: %s\n", pattern)
			continue
		}

		for _, filePath := range matches {
			info, err := os.Stat(filePath)
			if err != nil {
				fmt.Printf("Error accessing %s: %v\n", filePath, err)
				continue
			}

			if info.IsDir() {
				fmt.Printf("Skipping directory: %s\n", filePath)
				continue
			}

			if err := promptmanager.AddFileToPrompt(filePath); err != nil {
				return err
			}
		}
	}

	return nil
}

func HandleListCommand() error {
	return promptmanager.ListAttachments()
}

func HandleClearCommand() error {
	return promptmanager.ClearPrompt()
}

func HandleShowLastCommand() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting home directory: %w", err)
	}
	
	responseFile := filepath.Join(homeDir, ".minicmd", "last_response")
	content, err := os.ReadFile(responseFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no previous response found")
		}
		return fmt.Errorf("error reading last response: %w", err)
	}
	
	fmt.Print(string(content))
	return nil
}
