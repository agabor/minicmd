package commands

import (
	"fmt"
	"os"
	"yact/logic"
)

func HandleLastCommand(filePath string) error {
	contextMessages, err := logic.LoadContext()
	if err != nil {
		return fmt.Errorf("error loading context: %w", err)
	}

	lastIndex := len(contextMessages) - 1

	if lastIndex == -1 {
		return fmt.Errorf("no previous assistant message found")
	}

	if filePath == "" {
		fmt.Print(contextMessages[lastIndex].Content)
		return nil
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	contextMessages[lastIndex].Content = string(content)

	if err := logic.SaveContext(contextMessages); err != nil {
		return fmt.Errorf("error saving context: %w", err)
	}

	return nil
}
