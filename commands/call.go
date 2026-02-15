package commands

import (
	"fmt"
	"strings"
	"time"

	"yact/api"
	"yact/config"
	"yact/logic"
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

	responseContent, err := HandleCall(args, cfg, systemPrompt, "act")
	if err != nil {
		return err
	}

	fmt.Println("Processing response...")
	if err := logic.ProcessCodeBlocks(responseContent, safe); err != nil {
		return fmt.Errorf("error processing code blocks: %w", err)
	}

	fmt.Println("Done!")
	return nil
}

func HandleVerbalCommand(args []string, cfg *config.Config, systemPrompt string, promptType string) error {
	responseContent, err := HandleCall(args, cfg, systemPrompt, promptType)
	if err != nil {
		return err
	}

	fmt.Println("\n" + responseContent)
	return nil
}

func HandleNewCommand() error {
	if err := logic.ClearContext(); err != nil {
		return fmt.Errorf("error clearing context: %w", err)
	}

	fmt.Println("New context created")
	return nil
}

func HandleGoCommand(cfg *config.Config, systemPrompt string) error {
	messages, err := logic.LoadContext()
	if err != nil {
		return err
	}

	if err := validatePlanAndMessage(messages); err != nil {
		return err
	}

	messages = convertPlanToMessage(messages)

	if err := logic.SaveContext(messages); err != nil {
		return err
	}

	fmt.Printf("Sending request to Claude...\n")

	var client api.APIClient
	client = &api.ClaudeClient{}
	client.Init(cfg)

	fmt.Printf("Model: %s\n", client.GetModelName())

	done := make(chan bool)
	go showProgress(done)

	response, err := client.Call(messages, systemPrompt)

	done <- true
	close(done)

	if err != nil {
		return err
	}

	responseContent := response.Content

	if strings.TrimSpace(responseContent) == "" {
		return fmt.Errorf("error: empty response from Claude API")
	}

	response.Type = "act"

	updatedMessages := append(messages, response)
	if err := logic.SaveContext(updatedMessages); err != nil {
		fmt.Printf("Warning: could not save context: %v\n", err)
	}

	fmt.Println("Processing response...")
	if err := logic.ProcessCodeBlocks(responseContent, false); err != nil {
		return fmt.Errorf("error processing code blocks: %w", err)
	}

	fmt.Println("Done!")
	return nil
}

func HandleCall(args []string, cfg *config.Config, systemPrompt string, promptType string) (string, error) {
	prompt := strings.Join(args, " ")

	fmt.Printf("Sending request to Claude...\n")

	var client api.APIClient
	client = &api.ClaudeClient{}
	client.Init(cfg)

	fmt.Printf("Model: %s\n", client.GetModelName())

	messages, err := logic.BuildMessages(prompt)

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

	response.Type = promptType

	updatedMessages := append(messages, response)
	if err := logic.SaveContext(updatedMessages); err != nil {
		fmt.Printf("Warning: could not save context: %v\n", err)
	}
	return responseContent, nil
}
