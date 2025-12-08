package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"minicmd/apiclient"
	"minicmd/config"
	"minicmd/fileprocessor"
	"minicmd/promptmanager"
)

func showProgress(done chan bool) {
	chars := "⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏"
	idx := 0
	for {
		select {
		case <-done:
			fmt.Print("\r \r")
			return
		default:
			fmt.Printf("\r%c", chars[idx%len(chars)])
			idx++
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func HandleRunCommand(args []string, claudeFlag, ollamaFlag, deepseekFlag, verbose, debug, safe bool) error {
	// Check for conflicting provider options
	providerFlags := 0
	if claudeFlag {
		providerFlags++
	}
	if ollamaFlag {
		providerFlags++
	}
	if deepseekFlag {
		providerFlags++
	}
	if providerFlags > 1 {
		return fmt.Errorf("cannot specify multiple provider flags")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	// Determine which provider to use
	provider := cfg.DefaultProvider
	if claudeFlag {
		provider = "claude"
	} else if ollamaFlag {
		provider = "ollama"
	} else if deepseekFlag {
		provider = "deepseek"
	}

	// Get prompt content from args if provided, otherwise use default prompt file
	var prompt string
	if len(args) > 0 {
		// Use provided prompt content directly
		prompt = strings.Join(args, " ")
		fmt.Println("Using provided prompt content")
	} else {
		// Use default prompt file
		if err := promptmanager.EditPromptFile(); err != nil {
			return err
		}
		prompt, err = promptmanager.GetPromptFromFile()
		if err != nil {
			return err
		}
		fmt.Println("Using default prompt file")
	}

	if verbose {
		fmt.Printf("Prompt: %s\n", prompt)
		fmt.Println("---")
	}

	fmt.Printf("Sending request to %s...\n", strings.Title(provider))
	switch provider {
	case "claude":
		fmt.Printf("Model: %s\n", cfg.ClaudeModel)
	case "deepseek":
		fmt.Printf("Model: %s\n", cfg.DeepSeekModel)
	default:
		fmt.Printf("Model: %s\n", cfg.OllamaModel)
	}

	// Get attachments
	attachments, err := promptmanager.GetAttachments()
	if err != nil {
		return fmt.Errorf("error getting attachments: %w", err)
	}

	// Start progress indicator
	done := make(chan bool)
	go showProgress(done)

	// Call API
	var response, rawResponse string
	systemPrompt := config.SystemPrompt
	
	switch provider {
	case "claude":
		response, rawResponse, err = apiclient.CallClaude(prompt, cfg, systemPrompt, debug, attachments)
	case "deepseek":
		response, rawResponse, err = apiclient.CallDeepSeek(prompt, cfg, systemPrompt, debug, attachments)
	default:
		response, rawResponse, err = apiclient.CallOllama(prompt, cfg, systemPrompt, debug, attachments)
	}

	// Stop progress indicator
	done <- true
	close(done)

	if err != nil {
		return err
	}

	// Clear prompt
	if err := promptmanager.ClearPrompt(); err != nil {
		return fmt.Errorf("error clearing prompt: %w", err)
	}

	if response == "" {
		return fmt.Errorf("error: no response from %s API", strings.Title(provider))
	}

	if strings.TrimSpace(response) == "" {
		return fmt.Errorf("error: empty response from %s API", strings.Title(provider))
	}

	// Echo the response to see what we got
	if verbose {
		fmt.Println("Raw response:")
		fmt.Println("==============")
		fmt.Println(response)
		fmt.Println("==============")
		fmt.Println()
	}

	// In debug mode, also show the complete API response
	if debug && rawResponse != "" {
		fmt.Println("Complete API response (JSON):")
		fmt.Println("=============================")
		fmt.Println(rawResponse)
		fmt.Println("=============================")
		fmt.Println()
	}

	// Process the response and create files
	fmt.Println("Processing response...")
	if err := fileprocessor.ProcessCodeBlocks(response, safe); err != nil {
		return fmt.Errorf("error processing code blocks: %w", err)
	}

	fmt.Println("Done!")
	return nil
}

func HandleConfigCommand(args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	if len(args) == 0 {
		// Show current config
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

func ShowHelp() {
	fmt.Println("minicmd - AI-powered code generation tool")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  minicmd run [prompt_content] [--claude|--ollama|--deepseek]  # Generate code with optional custom prompt content")
	fmt.Println("  minicmd edit                      # Edit the prompt file")
	fmt.Println("  minicmd read <file>                # Add file reference to prompt")
	fmt.Println("  minicmd list                      # List current attachments")
	fmt.Println("  minicmd config                    # Show current configuration")
	fmt.Println("  minicmd config <key> <value>      # Set configuration value")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --claude    Use Claude API (requires API key)")
	fmt.Println("  --ollama    Use Ollama API (requires local Ollama)")
	fmt.Println("  --deepseek  Use DeepSeek API (requires API key)")
	fmt.Println("  --verbose, -v    Print verbose output")
	fmt.Println("  --debug, -d      Print debug output (includes verbose and raw API response)")
	fmt.Println("  --safe, -s       Add .new suffix to generated files")
	fmt.Println()
	fmt.Println("Configuration keys:")
	fmt.Println("  default_provider    Default AI provider (claude, ollama, or deepseek)")
	fmt.Println("  anthropic_api_key   Claude API key")
	fmt.Println("  deepseek_api_key    DeepSeek API key")
	fmt.Println("  claude_model        Claude model name")
	fmt.Println("  ollama_url          Ollama API URL")
	fmt.Println("  ollama_model        Ollama model name")
	fmt.Println("  deepseek_url        DeepSeek API URL")
	fmt.Println("  deepseek_model      DeepSeek model name")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  minicmd config anthropic_api_key sk-ant-...")
	fmt.Println("  minicmd config deepseek_api_key sk-...")
	fmt.Println("  minicmd config default_provider deepseek")
	fmt.Println("  minicmd --deepseek")
	fmt.Println("  minicmd run")
	fmt.Println("  minicmd run \"create a hello world function\"")
	fmt.Println("  minicmd run \"write a Python calculator\" --deepseek")
	fmt.Println("  minicmd run --verbose")
	fmt.Println("  minicmd run --debug")
	fmt.Println("  minicmd run --safe")
	fmt.Println("  minicmd read minicmd.go")
	fmt.Println("  minicmd list")
}
