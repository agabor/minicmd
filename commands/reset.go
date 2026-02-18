package commands

import (
	"fmt"
	"yact/logic"
)

func HandleResetCommand() error {
	messages, err := HandleReload()
	if err != nil {
		return err
	}

	var fileMessages []logic.Message
	for _, message := range messages {
		if message.Type == logic.MessageTypeFile {
			fileMessages = append(fileMessages, message)
		}
	}

	err = logic.SaveContext(fileMessages)

	fmt.Println("Context reset and files reloaded")
	return err
}
