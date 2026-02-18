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
	responseContent, err := HandleCall(args, cfg, systemPrompt, logic.MessageTypeCommand)
	if err != nil {
		return err
	}

	return processCodeBlocks(responseContent, safe)
}

func HandleVerbalCommand(args []string, cfg *config.Config, systemPrompt string, messageType logic.MessageType) error {
	responseContent, err := HandleCall(args, cfg, systemPrompt, messageType)
	if err != nil {
		return err
	}

	fmt.Println("\n" + responseContent)
	return nil
}

func HandleGoCommand(cfg *config.Config, systemPrompt string) error {

	messages, err := logic.LoadContextForMessageType(logic.MessageTypeCommand)
	if err != nil {
		fmt.Printf("Warning: could not load context: %v\n", err)
		messages = []logic.Message{}
	}

	responseContent, err := callClaudeAPI(messages, cfg, systemPrompt, logic.MessageTypeCommand)
	if err != nil {
		return err
	}

	return processCodeBlocks(responseContent, false)
}

func HandleCall(args []string, cfg *config.Config, systemPrompt string, messageType logic.MessageType) (string, error) {
	prompt := strings.Join(args, " ")

	contextMessages, err := logic.LoadContextForMessageType(messageType)
	if err != nil {
		fmt.Printf("Warning: could not load context: %v\n", err)
		contextMessages = []logic.Message{}
	}

	userMessage := logic.Message{
		Type:    messageType,
		Content: prompt,
	}

	messages := append(contextMessages, userMessage)

	responseContent, err := callClaudeAPI(messages, cfg, systemPrompt, messageType)
	if err != nil {
		return "", err
	}

	return responseContent, nil
}

func callClaudeAPI(messages []logic.Message, cfg *config.Config, systemPrompt string, messageType logic.MessageType) (string, error) {
	fmt.Printf("Sending request to Claude...\n")

	var client api.Client
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

	message := logic.Message{
		Content: responseContent,
		Type:    logic.ResponseType(messageType),
	}

	updatedMessages := append(messages, message)
	if err := logic.SaveContext(updatedMessages); err != nil {
		fmt.Printf("Warning: could not save context: %v\n", err)
	}

	return responseContent, nil
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
