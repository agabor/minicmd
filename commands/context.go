package commands

import (
	"fmt"
	"strconv"
	"strings"

	"yact/logic"
)

func HandleContextCommand(args []string) error {
	if len(args) == 0 {
		return handleContextList()
	}

	subcommand := args[0]

	switch subcommand {
	case "pop":
		return handleContextPop(args[1:])
	case "popto":
		return handleContextPopto(args[1:])
	case "del":
		return handleContextDel(args[1:])
	case "reload":
		return handleContextReload()
	default:
		return fmt.Errorf("unknown context subcommand: %s", subcommand)
	}
}

func handleContextList() error {
	messages, err := logic.LoadContext()
	if err != nil {
		return err
	}

	if len(messages) == 0 {
		fmt.Println("Context is empty")
		return nil
	}

	for i, message := range messages {
		fmt.Printf("[%d] %s (%s)", i, message.Role, message.Type)
		if message.Path != "" {
			fmt.Printf(" - %s", message.Path)
		} else {
			truncatedContent := message.Content
			if len(truncatedContent) > 200 {
				truncatedContent = truncatedContent[:200] + "..."
			}

			truncatedContent = strings.ReplaceAll(truncatedContent, "\n", " ")
			fmt.Printf(" - %s", truncatedContent)
		}
		fmt.Println()
	}

	return nil
}

func handleContextPop(args []string) error {
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

func handleContextPopto(args []string) error {
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

func handleContextDel(args []string) error {
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

func handleContextReload() error {
	if err := logic.ReloadContextFiles(); err != nil {
		return err
	}

	fmt.Println("Context files reloaded")
	return nil
}
