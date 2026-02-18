package api

import (
	"yact/config"
	"yact/logic"
)

type Client interface {
	Init(cfg *config.Config)
	GetModelName() string
	Call(messages []logic.Message, systemPrompt string) (logic.Message, error)
}
