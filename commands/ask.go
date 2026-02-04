package commands

import (
	"fmt"
	"strings"

	"yact/apiclient"
	"yact/config"
	"yact/promptmanager"
)

func HandleAskCommand(args []string, cfg *config.Config, systemPrompt string) error {
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

	var client apiclient.APIClient
	client = &apiclient.ClaudeClient{}
	client.Init(cfg)

	fmt.Printf("Model: %s\n", client.GetModelName())

	contextMessages, err := LoadContext()
	if err != nil {
		fmt.Printf("Warning: could not load context: %v\n", err)
		contextMessages = []apiclient.Message{}
	}

	attachments, err := promptmanager.GetAttachments()
	if err != nil {
		return fmt.Errorf("error getting attachments: %w", err)
	}

	messages := buildMessages(contextMessages, prompt, attachments)

	done := make(chan bool)
	go showProgress(done)

	response, err := client.Call(messages, systemPrompt)

	done <- true
	close(done)

	if err != nil {
		return err
	}

	if response.Content == "" {
		return fmt.Errorf("error: no response from Claude API")
	}

	if strings.TrimSpace(response.Content) == "" {
		return fmt.Errorf("error: empty response from Claude API")
	}

	if err := promptmanager.ClearPrompt(); err != nil {
		return fmt.Errorf("error clearing prompt: %w", err)
	}

	updatedMessages := append(messages, response)
	if err := SaveContext(updatedMessages); err != nil {
		fmt.Printf("Warning: could not save context: %v\n", err)
	}

	fmt.Println("\n" + response.Content)
	return nil
}
