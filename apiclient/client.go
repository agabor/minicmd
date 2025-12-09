package apiclient

import (
	"minicmd/config"
)

type APIClient interface {
	Init(cfg *config.Config)
	Call(userPrompt string, systemPrompt string, attachments []string) (string, error)
}
