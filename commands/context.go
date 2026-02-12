package commands

import (
	"fmt"
	"strings"

	"yact/logic"
)

func HandleContextCommand(args []string) error {
	messages, err := logic.LoadContext()
	if err != nil {
		return err
	}

	if len(messages) == 0 {
		fmt.Println("No messages in context")
		return nil
	}

	fmt.Println("Context messages:")
	fmt.Println()

	for i, msg := range messages {
		truncatedContent := msg.Content
		if len(truncatedContent) > 200 {
			truncatedContent = truncatedContent[:200] + "..."
		}

		truncatedContent = strings.ReplaceAll(truncatedContent, "\n", " ")

		fmt.Printf("[%d] %s: %s\n", i, msg.Role, truncatedContent)
	}

	return nil
}
