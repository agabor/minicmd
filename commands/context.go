package commands

import (
	"fmt"
	"strconv"
	"strings"

	"yact/logic"
)

func HandleContextCommand(args []string) error {
	if len(args) == 0 {
		return listContext()
	}

	subcommand := args[0]

	switch subcommand {
	case "pop":
		num := 1
		if len(args) > 1 {
			parsedNum, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid number: %s", args[1])
			}
			num = parsedNum
		}
		return popContext(num)
	case "del":
		if len(args) < 2 {
			return fmt.Errorf("missing index argument")
		}
		idx, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid index: %s", args[1])
		}
		return delContext(idx)
	default:
		return fmt.Errorf("unknown subcommand: %s", subcommand)
	}
}

func listContext() error {
	messages, err := logic.LoadContext()
	if err != nil {
		return err
	}

	if len(messages) == 0 {
		fmt.Println("No messages in context")
		return nil
	}

	fmt.Println("Context messages:")
	fmt.Println()

	for i, msg := range messages {
		truncatedContent := msg.Content
		if len(truncatedContent) > 200 {
			truncatedContent = truncatedContent[:200] + "..."
		}

		truncatedContent = strings.ReplaceAll(truncatedContent, "\n", " ")

		fmt.Printf("[%d] %s: %s\n", i, msg.Role, truncatedContent)
	}

	return nil
}

func popContext(num int) error {
	messages, err := logic.LoadContext()
	if err != nil {
		return err
	}

	if len(messages) == 0 {
		return fmt.Errorf("no messages in context")
	}

	if num > len(messages) {
		return fmt.Errorf("cannot pop %d messages, only %d messages in context", num, len(messages))
	}

	messages = messages[:len(messages)-num]

	if err := logic.SaveContext(messages); err != nil {
		return fmt.Errorf("error saving context: %w", err)
	}

	fmt.Printf("Removed %d message(s) from context\n", num)
	return nil
}

func delContext(idx int) error {
	messages, err := logic.LoadContext()
	if err != nil {
		return err
	}

	if len(messages) == 0 {
		return fmt.Errorf("no messages in context")
	}

	if idx < 0 || idx >= len(messages) {
		return fmt.Errorf("invalid index: %d (context has %d messages)", idx, len(messages))
	}

	messages = append(messages[:idx], messages[idx+1:]...)

	if err := logic.SaveContext(messages); err != nil {
		return fmt.Errorf("error saving context: %w", err)
	}

	fmt.Printf("Removed message at index %d\n", idx)
	return nil
}
