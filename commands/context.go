package commands

import (
	"fmt"
	"strings"

	"yact/logic"
)

func HandleContextCommand() error {
	messages, err := logic.LoadContext()
	if err != nil {
		return err
	}

	if len(messages) == 0 {
		fmt.Println("Context is empty")
		return nil
	}

	for i, message := range messages {
		fmt.Printf("[%d] %s", i, message.Type)
		if message.Path != "" {
			fmt.Printf(" - %s", message.Path)
		} else {
			truncatedContent := message.Content
			if len(truncatedContent) > 200 {
				truncatedContent = truncatedContent[:200] + "..."
			}

			truncatedContent = strings.ReplaceAll(truncatedContent, "\n", " ")
			fmt.Printf(" - %s", truncatedContent)
		}
		fmt.Println()
	}

	return nil
}
