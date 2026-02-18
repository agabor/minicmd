package commands

import (
	"fmt"
	"yact/logic"
)

func HandleResetCommand() error {
	messages, err := logic.LoadContext()
	if err != nil {
		return err
	}

	var fileMessages []logic.Message
	for _, message := range messages {
		if message.Type == logic.MessageTypeFile {
			fileMessages = append(fileMessages, message)
		}
	}

	if err := logic.SaveContext(fileMessages); err != nil {
		return err
	}

	var reloadedMessages []logic.Message
	seenPaths := make(map[string]bool)

	for _, message := range fileMessages {
		if message.Type == logic.MessageTypeFile {
			content, err := logic.ReadAsCodeBlock(message.Path)
			if err != nil {
				continue
			}

			reloadedMessages = append(reloadedMessages, logic.Message{
				Type:    logic.MessageTypeFile,
				Path:    message.Path,
				Content: content,
			})
			seenPaths[message.Path] = true
		}
	}

	if err := logic.SaveContext(reloadedMessages); err != nil {
		return err
	}

	fmt.Println("Context reset and files reloaded")
	return nil
}
