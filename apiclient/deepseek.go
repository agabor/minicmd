package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"yact/config"
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

type DeepSeekCompletionRequest struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	Suffix      string  `json:"suffix"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
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
	Text    string          `json:"text"`
}

type DeepSeekResponse struct {
	Choices []DeepSeekChoice `json:"choices"`
	Usage   DeepSeekUsage    `json:"usage"`
}

type DeepSeekClient struct {
	apiKey    string
	model     string
	maxTokens int
}

func (c *DeepSeekClient) Init(cfg *config.Config) {
	c.apiKey = cfg.DeepSeekAPIKey
	c.model = cfg.DeepSeekModel
	c.maxTokens = cfg.MaxOutputTokens
}

func (c *DeepSeekClient) GetModelName() string {
	return c.model
}

func (c *DeepSeekClient) FIM(prefix string, suffix string, attachments []string) (string, error) {
	if c.apiKey == "" {
		return "", fmt.Errorf("DeepSeek API key not configured. Please set your API key with: ya config deepseek_api_key YOUR_API_KEY")
	}

	startTime := time.Now()

	payload := DeepSeekCompletionRequest{
		Model:       c.model,
		Prompt:      prefix,
		Suffix:      suffix,
		MaxTokens:   c.maxTokens,
		Temperature: 0.1,
	}

	if err := saveLastRequestJSON(payload, "deepseek"); err != nil {
		fmt.Printf("Warning: could not save request to last_request file: %v\n", err)
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.deepseek.com/beta", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error calling DeepSeek API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("DeepSeek API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	var deepSeekResp DeepSeekResponse
	if err := json.Unmarshal(body, &deepSeekResp); err != nil {
		return "", fmt.Errorf("error parsing JSON response: %w", err)
	}

	duration := time.Since(startTime)
	fmt.Printf("DeepSeek API call took %.2f seconds\n", duration.Seconds())
	fmt.Printf("Token usage - Input: %d, Output: %d, Cached: %d\n",
		deepSeekResp.Usage.PromptTokens,
		deepSeekResp.Usage.CompletionTokens,
		deepSeekResp.Usage.PromptTokensDetails.CachedTokens)

	if deepSeekResp.Usage.CompletionTokens >= c.maxTokens {
		fmt.Printf("⚠️  WARNING: Maximum output tokens (%d) reached. Response may be incomplete.\n", c.maxTokens)
	}

	if len(deepSeekResp.Choices) == 0 {
		return "", fmt.Errorf("unexpected response format from DeepSeek API: no choices")
	}

	return deepSeekResp.Choices[0].Text, nil
}

func (c *DeepSeekClient) Call(userPrompt string, systemPrompt string, attachments []string) (string, error) {
	if c.apiKey == "" {
		return "", fmt.Errorf("DeepSeek API key not configured. Please set your API key with: ya config deepseek_api_key YOUR_API_KEY")
	}

	startTime := time.Now()

	messages := []DeepSeekMessage{
		{Role: "system", Content: systemPrompt},
	}

	for _, attachment := range attachments {
		messages = append(messages, DeepSeekMessage{
			Role:    "user",
			Content: attachment,
		})
	}

	messages = append(messages, DeepSeekMessage{
		Role:    "user",
		Content: userPrompt,
	})

	payload := DeepSeekRequest{
		Model:       c.model,
		Messages:    messages,
		MaxTokens:   c.maxTokens,
		Temperature: 0.1,
		Stream:      false,
	}

	if err := saveLastRequestJSON(payload, "deepseek"); err != nil {
		fmt.Printf("Warning: could not save request to last_request file: %v\n", err)
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.deepseek.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error calling DeepSeek API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("DeepSeek API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	var deepSeekResp DeepSeekResponse
	if err := json.Unmarshal(body, &deepSeekResp); err != nil {
		return "", fmt.Errorf("error parsing JSON response: %w", err)
	}

	duration := time.Since(startTime)
	fmt.Printf("DeepSeek API call took %.2f seconds\n", duration.Seconds())
	fmt.Printf("Token usage - Input: %d, Output: %d, Cached: %d\n",
		deepSeekResp.Usage.PromptTokens,
		deepSeekResp.Usage.CompletionTokens,
		deepSeekResp.Usage.PromptTokensDetails.CachedTokens)

	if deepSeekResp.Usage.CompletionTokens >= c.maxTokens {
		fmt.Printf("⚠️  WARNING: Maximum output tokens (%d) reached. Response may be incomplete.\n", c.maxTokens)
	}

	if len(deepSeekResp.Choices) == 0 {
		return "", fmt.Errorf("unexpected response format from DeepSeek API: no choices")
	}

	return deepSeekResp.Choices[0].Message.Content, nil
}

func saveLastRequestJSON(data interface{}, provider string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	yactDir := filepath.Join(homeDir, ".yact")
	if err := os.MkdirAll(yactDir, 0755); err != nil {
		return err
	}

	requestFile := filepath.Join(yactDir, "last_request_"+provider+".json")
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(requestFile, jsonData, 0644)
}