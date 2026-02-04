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

	client := apiclient.ClaudeClient{}
	client.Init(cfg)

	fmt.Printf("Model: %s\n", client.GetModelName())

	attachments, err := promptmanager.GetAttachments()
	if err != nil {
		return fmt.Errorf("error getting attachments: %w", err)
	}

	done := make(chan bool)
	go showProgress(done)

	response, err := client.Call(prompt, config.SystemPromptAsk, attachments)

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
