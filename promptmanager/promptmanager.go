package promptmanager

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func getConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".yact"), nil
}

func getPromptFile() (string, error) {
	dir, err := getConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "prompt"), nil
}

func getAttachmentsFile() (string, error) {
	dir, err := getConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "attachments.json"), nil
}

func EditPromptFile() error {
	promptFile, err := getPromptFile()
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(promptFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Create empty file if it doesn't exist
	if _, err := os.Stat(promptFile); os.IsNotExist(err) {
		if err := os.WriteFile(promptFile, []byte{}, 0644); err != nil {
			return err
		}
	}

	// Launch vim to edit the file
	cmd := exec.Command("vim", promptFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error launching vim: %w", err)
	}

	fmt.Printf("Prompt file edited: %s\n", promptFile)
	return nil
}

func GetPromptFromFile() (string, error) {
	promptFile, err := getPromptFile()
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(promptFile); os.IsNotExist(err) {
		return "", fmt.Errorf("prompt file does not exist. Please run 'y edit' to create and edit your prompt")
	}

	content, err := os.ReadFile(promptFile)
	if err != nil {
		return "", fmt.Errorf("error reading prompt file: %w", err)
	}

	contentStr := strings.TrimSpace(string(content))
	if contentStr == "" {
		return "", fmt.Errorf("prompt file is empty. Please run 'y edit' to add content to your prompt")
	}

	return contentStr, nil
}

func AddFileToPrompt(filePath string) error {
	attachmentsFile, err := getAttachmentsFile()
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(attachmentsFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Load existing attachments
	var attachments []string
	if data, err := os.ReadFile(attachmentsFile); err == nil {
		json.Unmarshal(data, &attachments)
	}

	// Add file if not already present
	found := false
	for _, path := range attachments {
		if path == filePath {
			found = true
			break
		}
	}

	if !found {
		attachments = append(attachments, filePath)

		// Save updated attachments
		data, err := json.MarshalIndent(attachments, "", "  ")
		if err != nil {
			return err
		}

		if err := os.WriteFile(attachmentsFile, data, 0644); err != nil {
			return err
		}

		fmt.Printf("Added file to attachments: %s\n", filePath)
	} else {
		fmt.Printf("File already in attachments: %s\n", filePath)
	}

	return nil
}

func GetAttachments() ([]string, error) {
	attachmentsFile, err := getAttachmentsFile()
	if err != nil {
		return nil, err
	}

	// Read attachments
	var attachments []string
	if data, err := os.ReadFile(attachmentsFile); err == nil {
		json.Unmarshal(data, &attachments)
	}

	// Build file contents list
	var fileContents []string

	for _, filePath := range attachments {
		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fileContents = append(fileContents, fmt.Sprintf("// %s\n// Error: File not found", filePath))
			continue
		}

		// Read file content
		content, err := os.ReadFile(filePath)
		if err != nil {
			fileContents = append(fileContents, fmt.Sprintf("// %s\n// Error reading file: %v", filePath, err))
			continue
		}

		// Add formatted content
		contentStr := strings.TrimRight(string(content), "\n")
		fileContents = append(fileContents, fmt.Sprintf("```\n// %s\n%s\n```", filePath, contentStr))
	}

	return fileContents, nil
}

func ListAttachments() error {
	attachmentsFile, err := getAttachmentsFile()
	if err != nil {
		return err
	}

	// Check if attachments file exists
	if _, err := os.Stat(attachmentsFile); os.IsNotExist(err) {
		fmt.Println("No attachments found.")
		return nil
	}

	// Load attachments
	var attachments []string
	data, err := os.ReadFile(attachmentsFile)
	if err != nil {
		return fmt.Errorf("error reading attachments file: %w", err)
	}

	if err := json.Unmarshal(data, &attachments); err != nil {
		return fmt.Errorf("error parsing attachments file: %w", err)
	}

	// Display attachments
	if len(attachments) == 0 {
		fmt.Println("No attachments found.")
	} else {
		fmt.Println("Current attachments:")
		for i, filePath := range attachments {
			// Check if file exists
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				fmt.Printf("  %d. %s (file not found)\n", i+1, filePath)
			} else {
				fmt.Printf("  %d. %s\n", i+1, filePath)
			}
		}
	}

	return nil
}

func ClearPrompt() error {
	promptFile, err := getPromptFile()
	if err != nil {
		return err
	}

	attachmentsFile, err := getAttachmentsFile()
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(promptFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Clear prompt file
	if err := os.WriteFile(promptFile, []byte{}, 0644); err != nil {
		return err
	}

	// Remove attachments file if it exists
	if _, err := os.Stat(attachmentsFile); err == nil {
		if err := os.Remove(attachmentsFile); err != nil {
			return err
		}
	}

	fmt.Println("Cleared prompt and attachments")
	return nil
}
