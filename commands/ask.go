package commands

import (
	"fmt"
	"yact/config"
)

func HandleAskCommand(args []string, cfg *config.Config, systemPrompt string) error {
	responseContent, err := HandleCall(args, cfg, systemPrompt)
	if err != nil {
		return err
	}

	fmt.Println("\n" + responseContent)
	return nil
}
