package logic

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"yact/promptmanager"

	"yact/apiclient"
)

func GetContextFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".yact", "context.json"), nil
}

func LoadContext() ([]apiclient.Message, error) {
	contextPath, err := GetContextFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(contextPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []apiclient.Message{}, nil
		}
		return nil, err
	}

	var messages []apiclient.Message
	if err := json.Unmarshal(data, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

func SaveContext(messages []apiclient.Message) error {
	contextPath, err := GetContextFilePath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(contextPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(messages, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(contextPath, data, 0644)
}

func ClearContext() error {
	contextPath, err := GetContextFilePath()
	if err != nil {
		return err
	}

	if err := os.Remove(contextPath); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}
func BuildMessages(prompt string) ([]apiclient.Message, error) {
	contextMessages, err := LoadContext()
	if err != nil {
		fmt.Printf("Warning: could not load context: %v\n", err)
		contextMessages = []apiclient.Message{}
	}

	attachments, err := promptmanager.GetAttachments()
	if err != nil {
		return nil, fmt.Errorf("error getting attachments: %w", err)
	}

	var content string
	if len(attachments) > 0 {
		content = prompt + "\n\n" + strings.Join(attachments, "\n")
	} else {
		content = prompt
	}

	userMessage := apiclient.Message{
		Role:    "user",
		Content: content,
	}

	return append(contextMessages, userMessage), nil
}
