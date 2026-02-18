package commands

import (
	"fmt"
	"strings"
	"yact/logic"
)

func HandleReload() ([]logic.Message, error) {
	messages, err := logic.LoadContext()
	if err != nil {
		return nil, err
	}

	var newMessages []logic.Message
	seenPaths := make(map[string]bool)
	var reloadErrors []string

	for _, message := range messages {
		if message.Type == logic.MessageTypeFile {
			if seenPaths[message.Path] {
				continue
			}

			content, err := logic.ReadAsCodeBlock(message.Path)
			if err != nil {
				reloadErrors = append(reloadErrors, fmt.Sprintf("could not reload %s: %v", message.Path, err))
				continue
			}

			newMessages = append(newMessages, logic.Message{Type: logic.MessageTypeFile, Path: message.Path, Content: content})
			seenPaths[message.Path] = true
		} else if message.Type == logic.MessageTypeAction {
			for _, block := range logic.ParseCodeBlocks(message.Content) {
				if seenPaths[block.Path] {
					continue
				}

				newMessages = append(newMessages, logic.Message{Type: logic.MessageTypeFile, Path: block.Path, Content: logic.AsCodeBlock(block.Path, block.Content)})
				seenPaths[block.Path] = true
			}
		} else {
			newMessages = append(newMessages, message)
		}
	}

	if err := logic.SaveContext(newMessages); err != nil {
		return nil, err
	}

	if len(reloadErrors) > 0 {
		return nil, fmt.Errorf("reloaded context with errors: %s", strings.Join(reloadErrors, "; "))
	}
	fmt.Println("Context files reloaded")
	return newMessages, nil
}
