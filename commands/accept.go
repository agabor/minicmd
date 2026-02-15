package commands

import (
	"fmt"
	"yact/api"
	"yact/logic"
)

func HandleAcceptCommand() ([]api.Message, error) {
	messages, err := logic.LoadContext()
	if err != nil {
		return nil, err
	}

	if len(messages) < 2 {
		return nil, fmt.Errorf("not enough messages in context (need at least 2)")
	}

	lastIdx := len(messages) - 1
	secondLastIdx := len(messages) - 2

	if messages[lastIdx].Type != "plan" {
		return nil, fmt.Errorf("last message is not a plan (type: %s)", messages[lastIdx].Type)
	}

	if messages[secondLastIdx].Type != "message" {
		return nil, fmt.Errorf("second last message is not a message (type: %s)", messages[secondLastIdx].Type)
	}

	messages[lastIdx].Role = "user"
	messages[lastIdx].Type = "message"

	messages = append(messages[:secondLastIdx], messages[secondLastIdx+1:]...)

	if err := logic.SaveContext(messages); err != nil {
		return nil, err
	}

	fmt.Println("Plan accepted and converted to message")
	return messages, nil
}
