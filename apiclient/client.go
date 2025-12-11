package apiclient

import (
	"minicmd/config"
)

type APIClient interface {
	Init(cfg *config.Config)
	GetModelName() string
	Call(userPrompt string, systemPrompt string, attachments []string) (string, error)
	FIM(prompt string) (string, error)
}