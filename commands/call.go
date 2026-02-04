package commands

import (
	"fmt"
	"strings"
	"time"

	"yact/apiclient"
	"yact/config"
	"yact/fileprocessor"
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

func HandleActCommand(args []string, safe bool, cfg *config.Config, systemPrompt string) error {

	responseContent, err := HandleCall(args, cfg, systemPrompt)
	if err != nil {
		return err
	}

	fmt.Println("Processing response...")
	if err := fileprocessor.ProcessCodeBlocks(responseContent, safe); err != nil {
		return fmt.Errorf("error processing code blocks: %w", err)
	}

	fmt.Println("Done!")
	return nil
}

func HandleAskCommand(args []string, cfg *config.Config, systemPrompt string) error {
	responseContent, err := HandleCall(args, cfg, systemPrompt)
	if err != nil {
		return err
	}

	fmt.Println("\n" + responseContent)
	return nil
}

func HandleNewCommand() error {
	if err := ClearContext(); err != nil {
		return fmt.Errorf("error clearing context: %w", err)
	}

	fmt.Println("New context created")
	return nil
}

func HandleCall(args []string, cfg *config.Config, systemPrompt string) (string, error) {
	prompt := strings.Join(args, " ")

	fmt.Printf("Sending request to Claude...\n")

	var client apiclient.APIClient
	client = &apiclient.ClaudeClient{}
	client.Init(cfg)

	fmt.Printf("Model: %s\n", client.GetModelName())

	messages, err := BuildMessages(prompt)

	if err != nil {
		return "", err
	}

	done := make(chan bool)
	go showProgress(done)

	response, err := client.Call(messages, systemPrompt)

	done <- true
	close(done)

	if err != nil {
		return "", err
	}

	responseContent := response.Content

	if strings.TrimSpace(responseContent) == "" {
		return "", fmt.Errorf("error: empty response from Claude API")
	}

	updatedMessages := append(messages, response)
	if err := SaveContext(updatedMessages); err != nil {
		fmt.Printf("Warning: could not save context: %v\n", err)
	}
	return responseContent, nil
}
