package commands

import (
	"fmt"
)

func ShowHelp() {
	fmt.Println("y - Yet Another Coding Tool")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  y act [prompt]          # Generate code with prompt")
	fmt.Println("  y bash [prompt]         # Generate a bash script file")
	fmt.Println("  y ask [question]        # Ask questions about the codebase")
	fmt.Println("  y read <file>           # Add file reference to prompt")
	fmt.Println("  y list                  # List current attachments")
	fmt.Println("  y clear                 # Clear prompt and attachments")
	fmt.Println("  y last                  # Show last AI response")
	fmt.Println("  y config                # Show current configuration")
	fmt.Println("  y config <key> <value>  # Set configuration value")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --safe, -s       Add .new suffix to generated files")
	fmt.Println()
	fmt.Println("Configuration keys:")
	fmt.Println("  anthropic_api_key   Claude API key")
	fmt.Println("  claude_model        Claude model name")
}
