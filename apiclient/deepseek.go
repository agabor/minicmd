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

type DeepSeekMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type DeepSeekRequest struct {
	Model       string            `json:"model"`
	Messages    []DeepSeekMessage `json:"messages"`
	MaxTokens   int               `json:"max_tokens"`
	Temperature float64           `json:"temperature"`
	Stream      bool              `json:"stream"`
}

type DeepSeekUsage struct {
	PromptTokens            int                        `json:"prompt_tokens"`
	CompletionTokens        int                        `json:"completion_tokens"`
	TotalTokens             int                        `json:"total_tokens"`
	PromptTokensDetails     DeepSeekPromptTokenDetails `json:"prompt_tokens_details"`
}

type DeepSeekPromptTokenDetails struct {
	CachedTokens int `json:"cached_tokens"`
}

type DeepSeekChoice struct {
	Index   int             `json:"index"`
	Message DeepSeekMessage `json:"message"`
}

type DeepSeekResponse struct {
	Choices []DeepSeekChoice `json:"choices"`
	Usage   DeepSeekUsage    `json:"usage"`
}

func CallDeepSeek(userPrompt string, cfg *config.Config, systemPrompt string, debug bool, attachments []string) (string, string, error) {
	if cfg.DeepSeekAPIKey == "" {
		return "", "", fmt.Errorf("DeepSeek API key not configured. Please set your API key with: minicmd config deepseek_api_key YOUR_API_KEY")
	}

	startTime := time.Now()

	// Build messages array
	messages := []DeepSeekMessage{
		{Role: "system", Content: systemPrompt},
	}

	// Add attachment files as separate messages
	for _, attachment := range attachments {
		messages = append(messages, DeepSeekMessage{
			Role:    "user",
			Content: attachment,
		})
	}

	// Add main user prompt
	messages = append(messages, DeepSeekMessage{
		Role:    "user",
		Content: userPrompt,
	})

	payload := DeepSeekRequest{
		Model:       cfg.DeepSeekModel,
		Messages:    messages,
		MaxTokens:   cfg.MaxOutputTokens,
		Temperature: 0.1,
		Stream:      false,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", "", fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", cfg.DeepSeekURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+cfg.DeepSeekAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("error calling DeepSeek API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("DeepSeek API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("error reading response: %w", err)
	}

	var deepSeekResp DeepSeekResponse
	if err := json.Unmarshal(body, &deepSeekResp); err != nil {
		return "", "", fmt.Errorf("error parsing JSON response: %w", err)
	}

	duration := time.Since(startTime)
	fmt.Printf("DeepSeek API call took %.2f seconds\n", duration.Seconds())
	fmt.Printf("Token usage - Input: %d, Output: %d, Cached: %d\n",
		deepSeekResp.Usage.PromptTokens,
		deepSeekResp.Usage.CompletionTokens,
		deepSeekResp.Usage.PromptTokensDetails.CachedTokens)

	// Check if maximum output tokens reached
	if deepSeekResp.Usage.CompletionTokens >= cfg.MaxOutputTokens {
		fmt.Printf("⚠️  WARNING: Maximum output tokens (%d) reached. Response may be incomplete.\n", cfg.MaxOutputTokens)
	}

	if len(deepSeekResp.Choices) == 0 {
		return "", "", fmt.Errorf("unexpected response format from DeepSeek API: no choices")
	}

	rawResponse := ""
	if debug {
		rawResponse = string(body)
	}

	return deepSeekResp.Choices[0].Message.Content, rawResponse, nil
}
