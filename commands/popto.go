package commands

import (
	"fmt"
	"strconv"
	"yact/logic"
)

func HandlePopto(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("index required for popto subcommand")
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

	numRemoved := len(messages) - idx - 1
	messages = messages[:idx+1]

	if err := logic.SaveContext(messages); err != nil {
		return err
	}

	fmt.Printf("Removed %d message(s)\n", numRemoved)
	return nil
}
