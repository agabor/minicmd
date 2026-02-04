package commands

import (
	"fmt"
	"strings"

	"yact/apiclient"
	"yact/config"
	"yact/promptmanager"
)

func HandleAskCommand(args []string, cfg *config.Config) error {
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

	messages := []apiclient.Message{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	done := make(chan bool)
	go showProgress(done)

	response, err := client.Call(messages, config.SystemPromptAsk)

	done <- true
	close(done)

	if err != nil {
		return err
	}

	if response == "" {
		return fmt.Errorf("error: no response from Claude API")
	}

	if strings.TrimSpace(response) == "" {
		return fmt.Errorf("error: empty response from Claude API")
	}

	if err := promptmanager.ClearPrompt(); err != nil {
		return fmt.Errorf("error clearing prompt: %w", err)
	}

	if err := saveLastResponse(response); err != nil {
		fmt.Printf("Warning: could not save response to last_response file: %v\n", err)
	}

	fmt.Println("\n" + response)
	return nil
}
