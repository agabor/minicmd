package apiclient

import (
	"context"
	"fmt"
	"strings"
	"time"

	"minicmd/config"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type ClaudeClient struct {
	apiKey         string
	model          string
	maxOutputTokens int
}

func (c *ClaudeClient) Init(cfg *config.Config) {
	c.apiKey = cfg.AnthropicAPIKey
	c.model = cfg.ClaudeModel
	c.maxOutputTokens = cfg.MaxOutputTokens
}

func (c *ClaudeClient) GetModelName() string {
	return c.model
}

func (c *ClaudeClient) FIM(prefix string, suffix string, attachments []string) (string, error) {
	if c.apiKey == "" {
		return "", fmt.Errorf("Claude API key not configured. Please set your API key with: minicmd config anthropic_api_key YOUR_API_KEY")
	}

	startTime := time.Now()

	client := anthropic.NewClient(option.WithAPIKey(c.apiKey))

	systemPrompt := `You are a code completion assistant. Complete the code based on the context provided.
Only return the code completion without any explanations or markdown formatting.
Maintain the same indentation and style as the existing code.`

	prompt := buildFimPrompt(prefix, suffix, attachments)

	params := anthropic.MessageNewParams{
		Model:     anthropic.F(c.model),
		MaxTokens: anthropic.F(int64(c.maxOutputTokens)),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		}),
	}

	params.System = anthropic.F([]anthropic.TextBlockParam{
		anthropic.NewTextBlock(systemPrompt),
	})

	message, err := client.Messages.New(context.Background(), params)

	if err != nil {
		return "", fmt.Errorf("error calling Claude API: %w", err)
	}

	duration := time.Since(startTime)
	fmt.Printf("Claude API call took %.2f seconds\n", duration.Seconds())
	fmt.Printf("Token usage - Input: %d, Output: %d\n",
		message.Usage.InputTokens,
		message.Usage.OutputTokens)

	if message.Usage.OutputTokens >= int64(c.maxOutputTokens) {
		fmt.Printf("⚠️  WARNING: Maximum output tokens (%d) reached. Response may be incomplete.\n", c.maxOutputTokens)
	}

	var responseText string
	for _, block := range message.Content {
		responseText += block.Text
	}

	return responseText, nil
}

func (c *ClaudeClient) Call(userPrompt string, systemPrompt string, attachments []string) (string, error) {
	if c.apiKey == "" {
		return "", fmt.Errorf("Claude API key not configured. Please set your API key with: minicmd config anthropic_api_key YOUR_API_KEY")
	}

	startTime := time.Now()

	client := anthropic.NewClient(option.WithAPIKey(c.apiKey))

	fullPrompt := userPrompt
	if len(attachments) > 0 {
		parts := append(attachments, userPrompt)
		fullPrompt = strings.Join(parts, "\n\n")
	}

	params := anthropic.MessageNewParams{
		Model:     anthropic.F(c.model),
		MaxTokens: anthropic.F(int64(c.maxOutputTokens)),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(fullPrompt)),
		}),
	}

	if systemPrompt != "" {
		params.System = anthropic.F([]anthropic.TextBlockParam{
			anthropic.NewTextBlock(systemPrompt),
		})
	}

	message, err := client.Messages.New(context.Background(), params)

	if err != nil {
		return "", fmt.Errorf("error calling Claude API: %w", err)
	}

	duration := time.Since(startTime)
	fmt.Printf("Claude API call took %.2f seconds\n", duration.Seconds())
	fmt.Printf("Token usage - Input: %d, Output: %d\n",
		message.Usage.InputTokens,
		message.Usage.OutputTokens)

	if message.Usage.OutputTokens >= int64(c.maxOutputTokens) {
		fmt.Printf("⚠️  WARNING: Maximum output tokens (%d) reached. Response may be incomplete.\n", c.maxOutputTokens)
	}

	var responseText string
	for _, block := range message.Content {
		responseText += block.Text
	}

	return responseText, nil
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