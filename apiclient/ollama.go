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

type OllamaClient struct {
	model string
	url   string
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	System string `json:"system"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

func (oc *OllamaClient) Init(cfg *config.Config) {
	oc.model = cfg.OllamaModel
	oc.url = cfg.OllamaURL
}

func (oc *OllamaClient) Call(userPrompt string, systemPrompt string, attachments []string) (string, error) {
	startTime := time.Now()

	fullPrompt := userPrompt
	if len(attachments) > 0 {
		parts := append(attachments, userPrompt)
		fullPrompt = joinStrings(parts, "\n\n")
	}

	payload := OllamaRequest{
		Model:  oc.model,
		Prompt: fullPrompt,
		System: systemPrompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	resp, err := http.Post(oc.url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error calling Ollama API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Ollama API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("error parsing JSON response: %w", err)
	}

	duration := time.Since(startTime)
	fmt.Printf("Ollama API call took %.2f seconds\n", duration.Seconds())

	return ollamaResp.Response, nil
}
