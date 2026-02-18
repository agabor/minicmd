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
			messages, err := logic.LoadContext()
			if err != nil {
				return err
			}
			content, err := logic.ReadAsCodeBlock(filePath)
			if err != nil {
				return err
			}

			if hasMessageWithPath(messages, filePath) {
				fmt.Printf("Skipping: %s\n", filePath)
				continue
			} else {
				fmt.Printf("Reading: %s\n", filePath)
			}

			messages = append(messages, logic.Message{Type: logic.MessageTypeFile, Path: filePath, Content: content})
			err2 := logic.SaveContext(messages)
			if err2 != nil {
				return err2
			}
		}
	}

	return nil
}

func hasMessageWithPath(messages []logic.Message, path string) bool {

	for _, message := range messages {
		if message.Path == path {
			return true
		}
	}
	return false
}
