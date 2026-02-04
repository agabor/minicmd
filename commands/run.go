package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"yact/apiclient"
	"yact/config"
	"yact/fileprocessor"
	"yact/promptmanager"
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

	yactDir := filepath.Join(homeDir, ".yact")
	if err := os.MkdirAll(yactDir, 0755); err != nil {
		return err
	}

	responseFile := filepath.Join(yactDir, "last_response")
	return os.WriteFile(responseFile, []byte(response), 0644)
}

func HandleRunCommand(args []string, safe bool, cfg *config.Config, systemPrompt string) error {

	var prompt string

	if len(args) > 0 {
		prompt = strings.Join(args, " ")
	} else {
		if err := promptmanager.EditPromptFile(); err != nil {
			return err
		}
		promptContent, err := promptmanager.GetPromptFromFile()
		if err != nil {
			return err
		}
		prompt = promptContent
		fmt.Println("Using default prompt file")
	}

	fmt.Printf("Sending request to Claude...\n")

	client := apiclient.ClaudeClient{}
	client.Init(cfg)

	fmt.Printf("Model: %s\n", client.GetModelName())

	attachments, err := promptmanager.GetAttachments()
	if err != nil {
		return fmt.Errorf("error getting attachments: %w", err)
	}

	done := make(chan bool)
	go showProgress(done)

	response, err := client.Call(prompt, systemPrompt, attachments)

	done <- true
	close(done)

	if err != nil {
		return err
	}

	if err := promptmanager.ClearPrompt(); err != nil {
		return fmt.Errorf("error clearing prompt: %w", err)
	}

	if response == "" {
		return fmt.Errorf("error: no response from Claude API")
	}

	if strings.TrimSpace(response) == "" {
		return fmt.Errorf("error: empty response from Claude API")
	}

	if err := saveLastResponse(response); err != nil {
		fmt.Printf("Warning: could not save response to last_response file: %v\n", err)
	}

	fmt.Println("Processing response...")
	if err := fileprocessor.ProcessCodeBlocks(response, safe); err != nil {
		return fmt.Errorf("error processing code blocks: %w", err)
	}

	fmt.Println("Done!")
	return nil
}
