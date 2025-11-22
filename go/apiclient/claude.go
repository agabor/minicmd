package apiclient

import (
	"context"
	"fmt"
	"time"

	"minicmd/config"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func CallClaude(userPrompt string, cfg *config.Config, systemPrompt string, debug bool, attachments []string) (string, string, error) {
	if cfg.AnthropicAPIKey == "" {
		return "", "", fmt.Errorf("Claude API key not configured. Please set your API key with: minicmd config anthropic_api_key YOUR_API_KEY")
	}

	startTime := time.Now()

	client := anthropic.NewClient(option.WithAPIKey(cfg.AnthropicAPIKey))

	// Build content array
	var content []anthropic.MessageParamContentUnion

	// Add attachment files as separate messages with cache control
	for _, attachment := range attachments {
		content = append(content, anthropic.NewTextBlock(attachment, anthropic.EphemeralCacheControlParam()))
	}

	// Add main user prompt
	content = append(content, anthropic.NewTextBlock(userPrompt))

	// Create message
	message, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.F(cfg.ClaudeModel),
		MaxTokens: anthropic.F(int64(4000)),
		System: anthropic.F([]anthropic.TextBlockParam{
			anthropic.NewTextBlock(systemPrompt),
		}),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(content...),
		}),
	})

	if err != nil {
		return "", "", fmt.Errorf("error calling Claude API: %w", err)
	}

	duration := time.Since(startTime)
	fmt.Printf("Claude API call took %.2f seconds\n", duration.Seconds())
	fmt.Printf("Token usage - Input: %d, Output: %d, Cache Create: %d, Cache Read: %d\n",
		message.Usage.InputTokens,
		message.Usage.OutputTokens,
		message.Usage.CacheCreationInputTokens,
		message.Usage.CacheReadInputTokens)

	// Extract text from response
	var responseText string
	for _, block := range message.Content {
		if textBlock, ok := block.AsUnion().(anthropic.ContentBlockText); ok {
			responseText += textBlock.Text
		}
	}

	rawResponse := ""
	if debug {
		rawResponse = fmt.Sprintf("%+v", message)
	}

	return responseText, rawResponse, nil
}
