package commands

import (
	"fmt"
)

func ShowHelp() {
	fmt.Println("ya - AI-powered code generation tool")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  ya run [prompt_content]  # Generate code with optional custom prompt content")
	fmt.Println("  ya bash [script_content] # Generate a bash script file")
	fmt.Println("  ya edit                      # Edit the prompt file")
	fmt.Println("  ya read <file>               # Add file reference to prompt")
	fmt.Println("  ya list                      # List current attachments")
	fmt.Println("  ya clear                     # Clear prompt and attachments")
	fmt.Println("  ya last                  # Show last AI response")
	fmt.Println("  ya config                    # Show current configuration")
	fmt.Println("  ya config <key> <value>      # Set configuration value")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --safe, -s       Add .new suffix to generated files")
	fmt.Println()
	fmt.Println("Configuration keys:")
	fmt.Println("  anthropic_api_key   Claude API key")
	fmt.Println("  claude_model        Claude model name")
}
