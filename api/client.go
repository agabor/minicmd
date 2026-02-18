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
	MessageTypeRevision  MessageType = "Revision"
)

func ResponseType(messageType MessageType) MessageType {

	switch messageType {
	case MessageTypeCommand:
		return MessageTypeAction
	case MessageTypeQuestion:
		return MessageTypeAnswer
	case MessageTypeObjective:
		return MessageTypePlan
	default:
		return messageType
	}
}

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
