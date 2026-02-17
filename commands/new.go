package commands

import (
	"fmt"
	"os"
	"yact/logic"
)

func HandleNewCommand() error {
	contextPath, err := logic.GetContextFilePath()
	if err != nil {
		return err
	}

	if err := os.Remove(contextPath); err != nil && !os.IsNotExist(err) {
		return err
	}

	fmt.Println("New context created")
	return nil
}
