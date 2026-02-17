package commands

import (
	"fmt"
	"strconv"
	"yact/logic"
)

func HandleDelete(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("index required for del subcommand")
	}

	idx, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid index: %s", args[0])
	}

	messages, err := logic.LoadContext()
	if err != nil {
		return err
	}

	if idx < 0 || idx >= len(messages) {
		return fmt.Errorf("index out of range: %d", idx)
	}

	messages = append(messages[:idx], messages[idx+1:]...)

	if err := logic.SaveContext(messages); err != nil {
		return err
	}

	fmt.Printf("Removed message at index %d\n", idx)
	return nil
}
