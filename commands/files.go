package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"yact/logic"
)

func HandleReadCommand(args []string) error {
	if len(args) < 1 {
		fmt.Println("Usage: y read <file> [<file2> ...]")
		return fmt.Errorf("missing file argument")
	}

	for _, pattern := range args {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return fmt.Errorf("error matching pattern %s: %w", pattern, err)
		}

		if len(matches) == 0 {
			fmt.Printf("No files found matching pattern: %s\n", pattern)
			continue
		}

		for _, filePath := range matches {
			info, err := os.Stat(filePath)
			if err != nil {
				fmt.Printf("Error accessing %s: %v\n", filePath, err)
				continue
			}

			if info.IsDir() {
				fmt.Printf("Skipping directory: %s\n", filePath)
				continue
			}

			if err := logic.AddFileToPrompt(filePath); err != nil {
				return err
			}
		}
	}

	return nil
}

func HandleLastCommand() error {
	contextMessages, err := logic.LoadContext()
	if err != nil {
		return fmt.Errorf("error loading context: %w", err)
	}

	if len(contextMessages) == 0 {
		return fmt.Errorf("no previous messages found")
	}

	for i := len(contextMessages) - 1; i >= 0; i-- {
		if contextMessages[i].Role == "assistant" {
			fmt.Print(contextMessages[i].Content)
			return nil
		}
	}

	return fmt.Errorf("no previous assistant message found")
}
