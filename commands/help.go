
package commands

import (
	"fmt"
)

func ShowHelp() {
	fmt.Println("ya - AI-powered code generation tool")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  ya run [prompt_content] [--claude|--ollama|--deepseek]  # Generate code with optional custom prompt content")
	fmt.Println("  ya bash [script_content] [--claude|--ollama|--deepseek]  # Generate a bash script file")
	fmt.Println("  ya edit                      # Edit the prompt file")
	fmt.Println("  ya read <file>               # Add file reference to prompt")
	fmt.Println("  ya list                      # List current attachments")
	fmt.Println("  ya clear                     # Clear prompt and attachments")
	fmt.Println("  ya showlast                  # Show last AI response")
	fmt.Println("  ya config                    # Show current configuration")
	fmt.Println("  ya config <key> <value>      # Set configuration value")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --claude    Use Claude API (requires API key)")
	fmt.Println("  --ollama    Use Ollama API (requires local Ollama)")
	fmt.Println("  --deepseek  Use DeepSeek API (requires API key)")
	fmt.Println("  --safe, -s       Add .new suffix to generated files")
	fmt.Println()
	fmt.Println("Configuration keys:")
	fmt.Println("  default_provider    Default AI provider (claude, ollama, or deepseek)")
	fmt.Println("  anthropic_api_key   Claude API key")
	fmt.Println("  deepseek_api_key    DeepSeek API key")
	fmt.Println("  claude_model        Claude model name")
	fmt.Println("  ollama_url          Ollama API URL")
	fmt.Println("  ollama_model        Ollama model name")
	fmt.Println("  deepseek_url        DeepSeek API URL")
	fmt.Println("  deepseek_model      DeepSeek model name")
}