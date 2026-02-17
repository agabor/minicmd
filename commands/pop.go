package commands

import (
	"fmt"
	"strconv"
	"yact/logic"
)

func HandlePop(args []string) error {
	messages, err := logic.LoadContext()
	if err != nil {
		return err
	}

	numToPop := 1
	if len(args) > 0 {
		num, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid number: %s", args[0])
		}
		numToPop = num
	}

	if numToPop > len(messages) {
		numToPop = len(messages)
	}

	messages = messages[:len(messages)-numToPop]

	if err := logic.SaveContext(messages); err != nil {
		return err
	}

	fmt.Printf("Removed %d message(s)\n", numToPop)
	return nil
}
