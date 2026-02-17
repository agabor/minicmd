package logic

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"yact/api"
)

func GetContextFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".yact", "context.json"), nil
}

func LoadContext() ([]api.Message, error) {
	contextPath, err := GetContextFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(contextPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []api.Message{}, nil
		}
		return nil, err
	}

	var messages []api.Message
	if err := json.Unmarshal(data, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

func SaveContext(messages []api.Message) error {
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

func ReadAsCodeBlock(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return AsCodeBlock(filePath, string(content)), nil
}

func AddFileToPrompt(filePath string) error {
	messages, err := LoadContext()
	if err != nil {
		return err
	}
	content, err := ReadAsCodeBlock(filePath)
	if err != nil {
		return err
	}
	messages = append(messages, api.Message{Role: "user", Type: "file", Path: filePath, Content: content})
	return SaveContext(messages)
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

func BuildMessages(prompt string) ([]api.Message, error) {
	contextMessages, err := LoadContext()
	if err != nil {
		fmt.Printf("Warning: could not load context: %v\n", err)
		contextMessages = []api.Message{}
	}

	userMessage := api.Message{
		Role:    "user",
		Type:    "message",
		Content: prompt,
	}

	return append(contextMessages, userMessage), nil
}
