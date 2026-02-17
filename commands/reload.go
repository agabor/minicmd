package commands

import (
	"fmt"
	"strings"
	"yact/api"
	"yact/logic"
)

func HandleReload() error {
	messages, err := logic.LoadContext()
	if err != nil {
		return err
	}

	var newMessages []api.Message
	seenPaths := make(map[string]bool)
	var reloadErrors []string

	for _, message := range messages {
		if message.Type == "file" {
			if seenPaths[message.Path] {
				continue
			}

			content, err := logic.ReadAsCodeBlock(message.Path)
			if err != nil {
				reloadErrors = append(reloadErrors, fmt.Sprintf("could not reload %s: %v", message.Path, err))
				continue
			}

			newMessages = append(newMessages, api.Message{Role: "user", Type: "file", Path: message.Path, Content: content})
			seenPaths[message.Path] = true
		} else if message.Type == "act" {
			for _, block := range logic.ParseCodeBlocks(message.Content) {
				if seenPaths[block.Path] {
					continue
				}

				newMessages = append(newMessages, api.Message{Role: "user", Type: "file", Path: block.Path, Content: logic.AsCodeBlock(block.Path, block.Content)})
				seenPaths[block.Path] = true
			}
		} else {
			newMessages = append(newMessages, message)
		}
	}

	if err := logic.SaveContext(newMessages); err != nil {
		return err
	}

	if len(reloadErrors) > 0 {
		return fmt.Errorf("reloaded context with errors: %s", strings.Join(reloadErrors, "; "))
	}
	fmt.Println("Context files reloaded")
	return nil
}
