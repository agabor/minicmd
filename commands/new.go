package commands

import (
	"fmt"
	"yact/logic"
)

func HandleNewCommand() error {
	err := logic.SaveContext(make([]logic.Message, 0))
	fmt.Println("New context created")
	return err
}
