package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"minicmd/config"
)

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	System string `json:"system"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

func CallOllama(userPrompt string, cfg *config.Config, systemPrompt string, debug bool, attachments []string) (string, string, error) {
	startTime := time.Now()

	// For Ollama, combine attachments with user prompt
	fullPrompt := userPrompt
	if len(attachments) > 0 {
		parts := append(attachments, userPrompt)
		fullPrompt = joinStrings(parts, "\n\n")
	}

	payload := OllamaRequest{
		Model:  cfg.OllamaModel,
		Prompt: fullPrompt,
		System: systemPrompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", "", fmt.Errorf("error marshaling request: %w", err)
	}

	resp, err := http.Post(cfg.OllamaURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", fmt.Errorf("error calling Ollama API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("Ollama API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("error reading response: %w", err)
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", "", fmt.Errorf("error parsing JSON response: %w", err)
	}

	duration := time.Since(startTime)
	fmt.Printf("Ollama API call took %.2f seconds\n", duration.Seconds())

	rawResponse := ""
	if debug {
		rawResponse = string(body)
	}

	return ollamaResp.Response, rawResponse, nil
}
