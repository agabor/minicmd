package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"minicmd/config"
	"minicmd/promptmanager"
)

func HandleEditCommand() error {
	return promptmanager.EditPromptFile()
}

func HandleReadCommand(args []string) error {
	if len(args) < 1 {
		fmt.Println("Usage: minicmd read <file> [<file2> ...]")
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

			if err := promptmanager.AddFileToPrompt(filePath); err != nil {
				return err
			}
		}
	}

	return nil
}

func HandleListCommand() error {
	return promptmanager.ListAttachments()
}

func HandleClearCommand() error {
	return promptmanager.ClearPrompt()
}

func HandleShowLastCommand() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting home directory: %w", err)
	}
	
	responseFile := filepath.Join(homeDir, ".minicmd", "last_response")
	content, err := os.ReadFile(responseFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no previous response found")
		}
		return fmt.Errorf("error reading last response: %w", err)
	}
	
	fmt.Print(string(content))
	return nil
}
