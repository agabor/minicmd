package commands

import (
	"fmt"
	"os"
	"strings"

	"minicmd/config"
	"minicmd/promptmanager"
)

func HandleFimCommand(args []string, provider string, safe bool, cfg *config.Config) error {
	if len(args) == 0 {
		return fmt.Errorf("fim command requires a file path argument")
	}

	filePath := args[0]

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	content := string(fileContent)
	parts := strings.Split(content, cfg.FimToken)

	if len(parts) != 2 {
		return fmt.Errorf("file must contain exactly one '%s' token", cfg.FimToken)
	}

	prefix := parts[0]
	suffix := parts[1]

	if provider == "" {
		provider = cfg.DefaultProvider
	}

	fmt.Printf("Sending fill-in-the-middle request to %s...\n", strings.Title(provider))

	client := getAPIClient(provider)
	client.Init(cfg)

	fmt.Printf("Model: %s\n", client.GetModelName())

	done := make(chan bool)
	go showProgress(done)

	attachments, err := promptmanager.GetAttachments()
	if err != nil {
		done <- true
		close(done)
		return err
	}

	prompt := buildFimPromptWithAttachments(prefix, suffix, attachments)
	response, err := client.FIM(prompt)

	done <- true
	close(done)

	if err != nil {
		return err
	}

	if response == "" {
		return fmt.Errorf("error: no response from %s API", strings.Title(provider))
	}

	if strings.TrimSpace(response) == "" {
		return fmt.Errorf("error: empty response from %s API", strings.Title(provider))
	}

	response = trimCodeBlockDelimiters(response)

	if err := saveLastResponse(response); err != nil {
		fmt.Printf("Warning: could not save response to last_response file: %v\n", err)
	}

	result := prefix + response + suffix

	fmt.Println("Writing result to file...")
	if err := os.WriteFile(filePath, []byte(result), 0644); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	fmt.Println("Done!")
	return nil
}

func buildFimPromptWithAttachments(prefix string, suffix string, attachments []string) string {
	var promptParts []string

	for _, attachment := range attachments {
		lines := strings.Split(attachment, "\n")
		if len(lines) > 0 && strings.HasPrefix(lines[0], "```") {
			filePath := ""
			fileContent := ""
			inContent := false

			for _, line := range lines {
				if strings.HasPrefix(line, "// ") && !inContent {
					filePath = strings.TrimPrefix(line, "// ")
					inContent = true
					continue
				}
				if strings.HasPrefix(line, "```") && inContent {
					break
				}
				if inContent && filePath != "" {
					fileContent += line + "\n"
				}
			}

			if filePath != "" {
				fileContent = strings.TrimRight(fileContent, "\n")
				promptParts = append(promptParts, fmt.Sprintf("<|file_sep|>%s\n%s", filePath, fileContent))
			}
		}
	}

	fimPart := fmt.Sprintf("<|fim_prefix|>%s<|fim_suffix|>%s<|fim_middle|>", prefix, suffix)
	promptParts = append(promptParts, fimPart)

	return strings.Join(promptParts, "\n\n")
}

func trimCodeBlockDelimiters(response string) string {
	response = strings.TrimSpace(response)

	if strings.HasPrefix(response, "```") {
		response = strings.TrimPrefix(response, "```")
		response = strings.TrimLeft(response, "\n")
	}

	if strings.HasSuffix(response, "```") {
		response = strings.TrimSuffix(response, "```")
		response = strings.TrimRight(response, "\n")
	}

	return response
}