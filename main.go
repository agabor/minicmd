package main

import (
	"fmt"
	"os"

	"minicmd/commands"

	flag "github.com/spf13/pflag"
)

func main() {
	// Define flags
	claudeFlag := flag.Bool("claude", false, "Use Claude API")
	ollamaFlag := flag.Bool("ollama", false, "Use Ollama API")
	deepseekFlag := flag.Bool("deepseek", false, "Use DeepSeek API")
	helpFlag := flag.BoolP("help", "h", false, "Show help message")
	verboseFlag := flag.BoolP("verbose", "v", false, "Print verbose output")
	debugFlag := flag.BoolP("debug", "d", false, "Print debug output (includes verbose and raw API response)")
	safeFlag := flag.BoolP("safe", "s", false, "Add .new suffix to generated files")

	flag.Parse()

	// Debug mode implies verbose mode
	if *debugFlag {
		*verboseFlag = true
	}

	args := flag.Args()

	// Handle help flag or no command
	if *helpFlag || len(args) == 0 || (len(args) > 0 && args[0] == "help") {
		commands.ShowHelp()
		return
	}

	// Determine provider from flags
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

	var err error

	switch command {
	case "edit":
		err = commands.HandleEditCommand()
	case "read":
		err = commands.HandleAddCommand(commandArgs)
	case "list":
		err = commands.HandleListCommand()
	case "config":
		err = commands.HandleConfigCommand(commandArgs)
	case "run":
		err = commands.HandleRunCommand(commandArgs, provider, *verboseFlag, *debugFlag, *safeFlag)
	case "clear":
		err = commands.HandleClearCommand()
	case "showlast":
		err = commands.HandleShowLastCommand()
	default:
		fmt.Printf("Error: Unknown command '%s'\n", command)
		fmt.Println("Run 'minicmd --help' for usage information.")
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
