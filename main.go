package main

import (
	"fmt"
	"io"
	"os"
	"yact/config/systemprompt"

	"yact/commands"
	"yact/config"

	flag "github.com/spf13/pflag"
)

func isStdinPiped() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func getPromptFromStdin() (string, error) {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func main() {
	helpFlag := flag.BoolP("help", "h", false, "Show help message")
	safeFlag := flag.BoolP("safe", "s", false, "Add .new suffix to generated files")

	flag.Parse()

	args := flag.Args()

	if *helpFlag {
		commands.ShowHelp()
		return
	}

	if len(args) == 0 {
		if isStdinPiped() {
			stdinContent, err := getPromptFromStdin()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(stdinContent)
			args = []string{"act", stdinContent}
		} else {
			fmt.Fprintf(os.Stderr, "Error: no command provided\n")
			fmt.Println("Run 'y --help' for usage information.")
			os.Exit(1)
		}
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	command := args[0]
	commandArgs := []string{}
	if len(args) > 1 {
		commandArgs = args[1:]
	}

	var commandErr error

	switch command {
	case "help":
		commands.ShowHelp()
		return
	case "read":
		commandErr = commands.HandleReadCommand(commandArgs)
	case "config":
		commandErr = commands.HandleConfigCommand(commandArgs, cfg)
	case "context":
		commandErr = commands.HandleContextCommand(commandArgs)
	case "act":
		commandErr = commands.HandleActCommand(commandArgs, *safeFlag, cfg, systemprompt.Act)
	case "bash":
		commandErr = commands.HandleActCommand(commandArgs, *safeFlag, cfg, systemprompt.Bash)
	case "ask":
		commandErr = commands.HandleAskCommand(commandArgs, cfg, systemprompt.Ask)
	case "plan":
		commandErr = commands.HandleAskCommand(commandArgs, cfg, systemprompt.Plan)
	case "new":
		commandErr = commands.HandleNewCommand()
	case "last":
		commandErr = commands.HandleLastCommand()
	default:
		fmt.Printf("Error: Unknown command '%s'\n", command)
		fmt.Println("Run 'y --help' for usage information.")
		os.Exit(1)
	}

	if commandErr != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", commandErr)
		os.Exit(1)
	}
}
