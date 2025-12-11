package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"minicmd/config"
)

type OllamaClient struct {
	model string
	url   string
}

type OllamaRequest struct {
	Model       string        `json:"model"`
	Prompt      string        `json:"prompt"`
	System      string        `json:"system"`
	Stream      bool          `json:"stream"`
	Options     OllamaOptions `json:"options"`
}

type OllamaOptions struct {
	Temperature float64 `json:"temperature"`
}

type OllamaResponse struct {
	Response        string `json:"response"`
	PromptEvalCount int    `json:"prompt_eval_count,omitempty"`
	EvalCount       int    `json:"eval_count,omitempty"`
	Done            bool   `json:"done"`
	DoneReason      string `json:"done_reason"`
}

func (oc *OllamaClient) Init(cfg *config.Config) {
	oc.model = cfg.OllamaModel
	oc.url = cfg.OllamaURL
}

func (c *OllamaClient) GetModelName() string {
	return c.model
}

func (c *OllamaClient) GetFIMSystemPrompt() string {
	return ""
}

func (oc *OllamaClient) Call(userPrompt string, systemPrompt string, attachments []string) (string, error) {
	startTime := time.Now()

	fullPrompt := userPrompt
	if len(attachments) > 0 {
		parts := append(attachments, userPrompt)
		fullPrompt = strings.Join(parts, "\n\n")
	}

	payload := OllamaRequest{
		Model:  oc.model,
		Prompt: fullPrompt,
		System: systemPrompt,
		Stream: false,
		Options: OllamaOptions{
			Temperature: 0.1,
		},
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
	fmt.Printf("Token usage - Input: %d, Output: %d\n",
		ollamaResp.PromptEvalCount,
		ollamaResp.EvalCount)

	fmt.Printf("Done: %t, Done Reason: %s\n", ollamaResp.Done, ollamaResp.DoneReason)

	return ollamaResp.Response, nil
}