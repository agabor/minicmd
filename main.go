package main

import (
	"fmt"
	"os"

	"minicmd/commands"
	"minicmd/config"

	flag "github.com/spf13/pflag"
)

func main() {
	claudeFlag := flag.Bool("claude", false, "Use Claude API")
	ollamaFlag := flag.Bool("ollama", false, "Use Ollama API")
	deepseekFlag := flag.Bool("deepseek", false, "Use DeepSeek API")
	helpFlag := flag.BoolP("help", "h", false, "Show help message")
	safeFlag := flag.BoolP("safe", "s", false, "Add .new suffix to generated files")

	flag.Parse()

	args := flag.Args()

	if *helpFlag || len(args) == 0 || (len(args) > 0 && args[0] == "help") {
		commands.ShowHelp()
		return
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
	case "showlast":
		commandErr = commands.HandleShowLastCommand()
	default:
		fmt.Printf("Error: Unknown command '%s'\n", command)
		fmt.Println("Run 'minicmd --help' for usage information.")
		os.Exit(1)
	}

	if commandErr != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", commandErr)
		os.Exit(1)
	}
}