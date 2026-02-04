package main

import (
	"fmt"
	"io"
	"os"

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
	claudeFlag := flag.Bool("claude", false, "Use Claude API")
	ollamaFlag := flag.Bool("ollama", false, "Use Ollama API")
	deepseekFlag := flag.Bool("deepseek", false, "Use DeepSeek API")
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
			args = []string{"run", stdinContent}
		} else {
			fmt.Fprintf(os.Stderr, "Error: no command provided\n")
			fmt.Println("Run 'ya --help' for usage information.")
			os.Exit(1)
		}
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	provider := ""
	providerFlags := 0
	if *claudeFlag {
		provider = "claude"
		providerFlags++
	}
	if *ollamaFlag {
		provider = "ollama"
		providerFlags++
	}
	if *deepseekFlag {
		provider = "deepseek"
		providerFlags++
	}
	if providerFlags > 1 {
		fmt.Fprintf(os.Stderr, "Error: cannot specify multiple provider flags\n")
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
	case "edit":
		commandErr = commands.HandleEditCommand()
	case "read":
		commandErr = commands.HandleReadCommand(commandArgs)
	case "list":
		commandErr = commands.HandleListCommand()
	case "config":
		commandErr = commands.HandleConfigCommand(commandArgs, cfg)
	case "run":
		commandErr = commands.HandleRunCommand(commandArgs, provider, *safeFlag, cfg, config.SystemPrompt)
	case "bash":
		commandErr = commands.HandleRunCommand(commandArgs, provider, *safeFlag, cfg, config.SystemPromptBash)
	case "fim":
		commandErr = commands.HandleFimCommand(commandArgs, provider, *safeFlag, cfg)
	case "clear":
		commandErr = commands.HandleClearCommand()
	case "last":
		commandErr = commands.HandleLastCommand()
	default:
		fmt.Printf("Error: Unknown command '%s'\n", command)
		fmt.Println("Run 'ya --help' for usage information.")
		os.Exit(1)
	}

	if commandErr != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", commandErr)
		os.Exit(1)
	}
}
