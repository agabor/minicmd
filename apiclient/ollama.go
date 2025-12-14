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

func (oc *OllamaClient) FIM(prefix string, suffix string, attachments []string) (string, error) {
	startTime := time.Now()

	prompt := buildFimPrompt(prefix, suffix, attachments)

	payload := OllamaRequest{
		Model:  oc.model,
		Prompt: prompt,
		Stream: false,
		Options: OllamaOptions{
			Temperature: 0.1,
		},
	}

	if err := saveLastRequestJSON(payload, "ollama"); err != nil {
		fmt.Printf("Warning: could not save request to last_request file: %v\n", err)
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

	if err := saveLastRequestJSON(payload, "ollama"); err != nil {
		fmt.Printf("Warning: could not save request to last_request file: %v\n", err)
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

func buildFimPrompt(prefix string, suffix string, attachments []string) string {
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