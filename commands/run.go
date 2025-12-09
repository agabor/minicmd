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

func HandleRunCommand(args []string, provider string, safe bool) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	// If no provider specified via flags, use default from config
	if provider == "" {
		provider = cfg.DefaultProvider
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
	var response string
	systemPrompt := config.SystemPrompt
	
	switch provider {
	case "claude":
		response, err = apiclient.CallClaude(prompt, cfg, systemPrompt, attachments)
	case "deepseek":
		response, err = apiclient.CallDeepSeek(prompt, cfg, systemPrompt, attachments)
	default:
		response, err = apiclient.CallOllama(prompt, cfg, systemPrompt, attachments)
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

	// Process the response and create files
	fmt.Println("Processing response...")
	if err := fileprocessor.ProcessCodeBlocks(response, safe); err != nil {
		return fmt.Errorf("error processing code blocks: %w", err)
	}

	fmt.Println("Done!")
	return nil
}
