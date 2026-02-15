package commands

import (
	"fmt"
	"yact/api"
	"yact/logic"
)

func validatePlanAndMessage(messages []api.Message) error {
	if len(messages) < 2 {
		return fmt.Errorf("not enough messages in context (need at least 2)")
	}

	lastIdx := len(messages) - 1
	secondLastIdx := len(messages) - 2

	if messages[lastIdx].Type != "plan" {
		return fmt.Errorf("last message is not a plan (type: %s)", messages[lastIdx].Type)
	}

	if messages[secondLastIdx].Type != "message" {
		return fmt.Errorf("second last message is not a message (type: %s)", messages[secondLastIdx].Type)
	}

	return nil
}

func convertPlanToMessage(messages []api.Message) []api.Message {
	lastIdx := len(messages) - 1
	secondLastIdx := len(messages) - 2

	messages[lastIdx].Role = "user"
	messages[lastIdx].Type = "message"

	return append(messages[:secondLastIdx], messages[secondLastIdx+1:]...)
}

func HandleAcceptCommand() error {
	messages, err := logic.LoadContext()
	if err != nil {
		return err
	}

	if err := validatePlanAndMessage(messages); err != nil {
		return err
	}

	messages = convertPlanToMessage(messages)

	if err := logic.SaveContext(messages); err != nil {
		return err
	}

	fmt.Println("Plan accepted and converted to message")
	return nil
}
