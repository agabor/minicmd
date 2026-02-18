package logic

import (
	"encoding/json"
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

func LoadContextForMessageType(messageType api.MessageType) ([]api.Message, error) {
	messages, err := LoadContext()
	if err != nil {
		return nil, err
	}

	return filterMessagesForMessageType(messages, messageType), nil
}

func filterMessagesForMessageType(messages []api.Message, messageType api.MessageType) []api.Message {
	var allowedTypes []api.MessageType

	switch messageType {
	case api.MessageTypeCommand:
		allowedTypes = []api.MessageType{api.MessageTypeFile, api.MessageTypeCommand, api.MessageTypeAction}
	case api.MessageTypeObjective:
		allowedTypes = []api.MessageType{api.MessageTypeFile, api.MessageTypeQuestion, api.MessageTypeAnswer, api.MessageTypeObjective, api.MessageTypePlan, api.MessageTypeRevision}
	case api.MessageTypeQuestion:
		allowedTypes = []api.MessageType{api.MessageTypeFile, api.MessageTypeQuestion, api.MessageTypeAnswer, api.MessageTypeObjective, api.MessageTypePlan}
	default:
		return make([]api.Message, 0)
	}

	var filtered []api.Message
	for _, msg := range messages {
		for _, allowed := range allowedTypes {
			if msg.Type == allowed {
				if messageType == api.MessageTypeObjective && msg.Type == api.MessageTypePlan {
					filtered = append(filtered, api.Message{Type: api.MessageTypeRevision, Content: msg.Content})
				} else {
					filtered = append(filtered, msg)
				}
				break
			}
		}
	}

	return filtered
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
	messages, err := LoadContextForMessageType(api.MessageTypeFile)
	if err != nil {
		return err
	}
	content, err := ReadAsCodeBlock(filePath)
	if err != nil {
		return err
	}
	messages = append(messages, api.Message{Type: api.MessageTypeFile, Path: filePath, Content: content})
	return SaveContext(messages)
}
