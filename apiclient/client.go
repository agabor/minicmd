package apiclient

import (
	"minicmd/config"
)

// APIClient defines the interface that all API clients must implement
type APIClient interface {
	// Call sends a prompt to the API and returns the response
	Call(userPrompt string, cfg *config.Config, systemPrompt string, debug bool, attachments []string) (string, string, error)
}
