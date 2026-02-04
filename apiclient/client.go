package apiclient

import (
	"yact/config"
)

type Message struct {
	Role    string
	Content string
}

type APIClient interface {
	Init(cfg *config.Config)
	GetModelName() string
	Call(messages []Message, systemPrompt string) (Message, error)
}
