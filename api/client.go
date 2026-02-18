package api

import (
	"yact/config"
)

type MessageType string

const (
	MessageTypeFile      MessageType = "File"
	MessageTypeQuestion  MessageType = "Question"
	MessageTypeAnswer    MessageType = "Answer"
	MessageTypeCommand   MessageType = "Command"
	MessageTypeAction    MessageType = "Action"
	MessageTypeObjective MessageType = "Objective"
	MessageTypePlan      MessageType = "Plan"
)

type Message struct {
	Type    MessageType
	Path    string
	Content string
}

type APIClient interface {
	Init(cfg *config.Config)
	GetModelName() string
	Call(messages []Message, systemPrompt string) (Message, error)
}
