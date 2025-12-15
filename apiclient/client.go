package apiclient

import (
	"yact/config"
)

type APIClient interface {
	Init(cfg *config.Config)
	GetModelName() string
	Call(userPrompt string, systemPrompt string, attachments []string) (string, error)
	FIM(prefix string, suffix string, attachments []string) (string, error)
}