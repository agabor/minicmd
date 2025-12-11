package commands

import (
	"fmt"
	"os"
	"strings"

	"minicmd/config"
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

	prompt := buildFimPrompt(prefix, suffix)
	response, err := client.Call(prompt, "", []string{})

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

func buildFimPrompt(prefix string, suffix string) string {
	return fmt.Sprintf("<｜fim_prefix｜>%s<｜fim_middle｜><｜fim_suffix｜>%s<｜fim_middle｜>", prefix, suffix)
}