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
	chars := []rune("⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏")
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

func saveLastResponse(response string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	
	minicmdDir := filepath.Join(homeDir, ".minicmd")
	if err := os.MkdirAll(minicmdDir, 0755); err != nil {
		return err
	}
	
	responseFile := filepath.Join(minicmdDir, "last_response")
	return os.WriteFile(responseFile, []byte(response), 0644)
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

	// Save the response to last_response file
	if err := saveLastResponse(response); err != nil {
		fmt.Printf("Warning: could not save response to last_response file: %v\n", err)
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
