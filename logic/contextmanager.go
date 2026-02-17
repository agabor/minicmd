package logic

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

func getFileContentAsCodeBlock(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return asCodeBlock(filePath, string(content)), nil
}

func AddFileToPrompt(filePath string) error {
	messages, err := LoadContext()
	if err != nil {
		return err
	}
	content, err := getFileContentAsCodeBlock(filePath)
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

func ReloadContextFiles() error {
	messages, err := LoadContext()
	if err != nil {
		return err
	}

	newMessages := []api.Message{}
	seenPaths := make(map[string]bool)
	var reloadErrors []string

	for _, message := range messages {
		if message.Type == "file" {
			if seenPaths[message.Path] {
				continue
			}

			content, err := getFileContentAsCodeBlock(message.Path)
			if err != nil {
				reloadErrors = append(reloadErrors, fmt.Sprintf("could not reload %s: %v", message.Path, err))
				continue
			}

			newMessages = append(newMessages, api.Message{Role: "user", Type: "file", Path: message.Path, Content: content})
			seenPaths[message.Path] = true
		} else if message.Type == "act" {
			for _, block := range ParseCodeBlocks(message.Content) {
				if seenPaths[block.path] {
					continue
				}

				newMessages = append(newMessages, api.Message{Role: "user", Type: "file", Path: block.path, Content: block.serialize()})
				seenPaths[block.path] = true
			}
		} else {
			newMessages = append(newMessages, message)
		}
	}

	if err := SaveContext(newMessages); err != nil {
		return err
	}

	if len(reloadErrors) > 0 {
		return fmt.Errorf("reloaded context with errors: %s", strings.Join(reloadErrors, "; "))
	}

	return nil
}
