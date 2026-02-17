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

	if err := processCodeBlocks(responseContent, safe); err != nil {
		return err
	}

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

	messages, err := HandleAcceptCommand()
	if err != nil {
		return err
	}

	responseContent, err := callClaudeAPI(messages, cfg, systemPrompt)
	if err != nil {
		return err
	}

	saveContext(messages, responseContent, "act")

	return processCodeBlocks(responseContent, false)
}

func HandleCall(args []string, cfg *config.Config, systemPrompt string, promptType string) (string, error) {
	prompt := strings.Join(args, " ")

	messages, err := logic.BuildMessages(prompt)
	if err != nil {
		return "", err
	}

	responseContent, err := callClaudeAPI(messages, cfg, systemPrompt)
	if err != nil {
		return "", err
	}

	saveContext(messages, responseContent, promptType)

	return responseContent, nil
}

func callClaudeAPI(messages []api.Message, cfg *config.Config, systemPrompt string) (string, error) {
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
		return "", err
	}

	responseContent := response.Content

	if strings.TrimSpace(responseContent) == "" {
		return "", fmt.Errorf("error: empty response from Claude API")
	}

	return responseContent, nil
}

func saveContext(messages []api.Message, content string, messageType string) {
	response := api.Message{
		Content: content,
		Role:    "assistant",
		Type:    messageType,
	}

	updatedMessages := append(messages, response)
	if err := logic.SaveContext(updatedMessages); err != nil {
		fmt.Printf("Warning: could not save context: %v\n", err)
	}
}

func processCodeBlocks(content string, safe bool) error {
	fmt.Println("Processing response...")
	var parseErrors []string
	for _, codeBlock := range logic.ParseCodeBlocks(content) {
		err := codeBlock.Write(safe)
		if err != nil {
			parseErrors = append(parseErrors, fmt.Sprintf("%v", err))
		}
	}

	if len(parseErrors) > 0 {
		return fmt.Errorf("error processing code blocks: %s", strings.Join(parseErrors, "; "))
	}

	fmt.Println("Done!")
	return nil
}
