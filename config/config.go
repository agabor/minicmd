package config

import (
        "encoding/json"
        "os"
        "path/filepath"
)

const (
        OllamaURL            = "http://localhost:11434/api/generate"
        OllamaModel          = "qwen2.5-coder:7b"
        ClaudeModel          = "claude-haiku-4-5-20251001"
        DeepSeekURL          = "https://api.deepseek.com/v1/chat/completions"
        DeepSeekModel        = "deepseek-coder"
        DefaultMaxTokens     = 8192
        DefaultFimToken      = "//FIM"
        SystemPrompt = "CODE GENERATION ASSISTANT\n\n" +
                "====================\n" +
                "STRICT OUTPUT RULES:\n" +
                "====================\n\n" +
                "1. OUTPUT STRUCTURE (REQUIRED):\n" +
                "   - Only output code blocks\n" +
                "   - No explanations before code blocks\n" +
                "   - No explanations after code blocks\n" +
                "   - No summaries\n" +
                "   - No descriptions\n\n" +
                "2. CODE BLOCK FORMAT (REQUIRED):\n" +
                "   ```\n" +
                "   // full/path/to/file.ext\n" +
                "   [complete file content here]\n" +
                "   ```\n\n" +
                "   Rules:\n" +
                "   - Start with: ```\n" +
                "   - Next line: comment with full file path\n" +
                "   - Then: complete file content\n" +
                "   - End with: ```\n" +
                "   - One code block = one file\n" +
                "   - Do NOT add language identifier after ```\n\n" +
                "3. FILE MODIFICATION RULES:\n" +
                "   When editing existing files:\n" +
                "   - Return COMPLETE file (not partial)\n" +
                "   - Keep all original comments\n" +
                "   - Keep all original indentation\n" +
                "   - Keep all original blank lines\n" +
                "   - Keep all original whitespace\n" +
                "   - Only include if code logic changed\n" +
                "   - Do NOT include if only whitespace changed\n\n" +
                "4. WHAT TO INCLUDE:\n" +
                "   Include these files:\n" +
                "   - New files you created (complete content)\n" +
                "   - Files where you changed code logic (complete content)\n\n" +
                "   Do NOT include:\n" +
                "   - Files with no changes\n" +
                "   - Files with only whitespace changes\n\n" +
                "5. CODE QUALITY REQUIREMENTS:\n" +
                "   - Use descriptive variable names\n" +
                "   - Use descriptive function names\n" +
                "   - Keep functions small (one purpose per function)\n" +
                "   - Write clear, readable code\n" +
                "   - Do NOT write code comments\n" +
                "   - Make code self-explanatory\n\n" +
                "EXAMPLE CORRECT OUTPUT:\n" +
                "```\n" +
                "// src/handlers/user.go\n" +
                "[complete file content]\n" +
                "```\n\n" +
                "```\n" +
                "// src/models/user.go\n" +
                "[complete file content]\n" +
                "```\n\n" +
                "INVALID OUTPUT EXAMPLES (DO NOT DO THIS):\n" +
                "- Text before code blocks\n" +
                "- Text after code blocks\n" +
                "- \"Here's the code...\"\n" +
                "- \"I've updated...\"\n" +
                "- Explanations of changes\n" +
                "- Partial file content\n" +
                "- Language identifier: ```go (WRONG)\n\n" +
                "\n\nBEFORE RESPONDING CHECK:\n" +
                "✓ Check: Using ``` without language identifier?\n" +
                "✓ Check: File path comment on line 2?\n" +
                "✓ Check: Complete file content?\n" +
                "✓ Check: No text outside code blocks?\n"
                "REMEMBER: Only code blocks. Nothing else."
        SystemPromptBash = "BASH SCRIPT GENERATION ASSISTANT\n\n" +
                "====================\n" +
                "STRICT OUTPUT RULES:\n" +
                "====================\n\n" +
                "1. OUTPUT STRUCTURE (REQUIRED):\n" +
                "   - Only output ONE code block\n" +
                "   - No explanations before code block\n" +
                "   - No explanations after code block\n" +
                "   - No usage examples outside script\n" +
                "   - No descriptions\n\n" +
                "2. CODE BLOCK FORMAT (REQUIRED):\n" +
                "   ```\n" +
                "   #!/bin/bash\n" +
                "   # filename.sh\n" +
                "   [complete script content here]\n" +
                "   ```\n\n" +
                "   Rules:\n" +
                "   - Start with: ``` (no language identifier)\n" +
                "   - Line 1: #!/bin/bash (shebang)\n" +
                "   - Line 2: # filename.sh (script name)\n" +
                "   - Then: complete script content\n" +
                "   - End with: ```\n" +
                "   - Only ONE code block per response\n" +
                "   - Do NOT use ```bash (wrong)\n\n" +
                "3. SCRIPT STRUCTURE REQUIREMENTS:\n" +
                "   Every script must include:\n" +
                "   - Shebang: #!/bin/bash (first line)\n" +
                "   - Filename comment: # scriptname.sh (second line)\n" +
                "   - Error handling: set -euo pipefail (recommended)\n" +
                "   - Main script logic\n\n" +
                "4. ERROR HANDLING:\n" +
                "   - Add: set -euo pipefail near the top\n" +
                "   - Exit on errors\n" +
                "   - Handle command failures\n" +
                "   - Check for required commands/files\n\n" +
                "5. CODE QUALITY REQUIREMENTS:\n" +
                "   Variables:\n" +
                "   - Global variables: UPPER_CASE\n" +
                "   - Local variables: lower_case\n" +
                "   - Always quote variables: \"$VARIABLE\"\n\n" +
                "   Comments:\n" +
                "   - Add comments for complex operations\n" +
                "   - Explain non-obvious logic\n" +
                "   - Document expected inputs\n\n" +
                "   Best practices:\n" +
                "   - Quote all variable expansions\n" +
                "   - Check if commands exist before using\n" +
                "   - Validate inputs\n" +
                "   - Handle edge cases\n\n" +
                "6. IF SCRIPT ACCEPTS ARGUMENTS:\n" +
                "   Include inside the script:\n" +
                "   - Usage function showing syntax\n" +
                "   - Argument validation\n" +
                "   - Help message (if -h or --help)\n\n" +
                "EXAMPLE CORRECT OUTPUT:\n" +
                "```\n" +
                "#!/bin/bash\n" +
                "# deploy.sh\n" +
                "set -euo pipefail\n\n" +
                "# Script content here\n" +
                "```\n\n" +
                "INVALID OUTPUT EXAMPLES (DO NOT DO THIS):\n" +
                "- Text before code block\n" +
                "- Text after code block\n" +
                "- \"Here's the script...\"\n" +
                "- \"This script does...\"\n" +
                "- Multiple code blocks\n" +
                "- Usage examples outside script\n" +
                "- Missing shebang\n" +
                "- Missing filename comment\n" +
                "- Using ```bash (WRONG - use ``` only)\n\n" +
                "BEFORE RESPONDING CHECK:\n" +
                "✓ Using ``` without language identifier?\n" +
                "✓ Shebang on line 1?\n" +
                "✓ Filename comment on line 2?\n" +
                "✓ Error handling included?\n" +
                "✓ No text outside code block?\n\n" +
                "REMEMBER: Only ONE code block with ```. Nothing else."
)

type Config struct {
        DefaultProvider  string `json:"default_provider"`
        AnthropicAPIKey  string `json:"anthropic_api_key"`
        DeepSeekAPIKey   string `json:"deepseek_api_key"`
        OllamaURL        string `json:"ollama_url"`
        OllamaModel      string `json:"ollama_model"`
        ClaudeModel      string `json:"claude_model"`
        DeepSeekURL      string `json:"deepseek_url"`
        DeepSeekModel    string `json:"deepseek_model"`
        MaxOutputTokens  int    `json:"max_output_tokens"`
        FimToken         string `json:"fim_token"`
}

func getConfigDir() (string, error) {
        home, err := os.UserHomeDir()
        if err != nil {
                return "", err
        }
        return filepath.Join(home, ".minicmd"), nil
}

func getConfigFile() (string, error) {
        dir, err := getConfigDir()
        if err != nil {
                return "", err
        }
        return filepath.Join(dir, "config"), nil
}

func DefaultConfig() *Config {
        return &Config{
                DefaultProvider: "ollama",
                AnthropicAPIKey: "",
                DeepSeekAPIKey:  "",
                OllamaURL:       OllamaURL,
                OllamaModel:     OllamaModel,
                ClaudeModel:     ClaudeModel,
                DeepSeekURL:     DeepSeekURL,
                DeepSeekModel:   DeepSeekModel,
                MaxOutputTokens: DefaultMaxTokens,
                FimToken:        DefaultFimToken,
        }
}

func Load() (*Config, error) {
        cfg := DefaultConfig()

        configFile, err := getConfigFile()
        if err != nil {
                return cfg, err
        }

        data, err := os.ReadFile(configFile)
        if err != nil {
                if os.IsNotExist(err) {
                        return cfg, nil
                }
                return cfg, err
        }

        if err := json.Unmarshal(data, cfg); err != nil {
                return cfg, err
        }

        // Ensure MaxOutputTokens has a valid value
        if cfg.MaxOutputTokens <= 0 {
                cfg.MaxOutputTokens = DefaultMaxTokens
        }

        // Ensure FimToken has a valid value
        if cfg.FimToken == "" {
                cfg.FimToken = DefaultFimToken
        }

        return cfg, nil
}

func (c *Config) Save() error {
        configFile, err := getConfigFile()
        if err != nil {
                return err
        }

        dir, err := getConfigDir()
        if err != nil {
                return err
        }

        if err := os.MkdirAll(dir, 0755); err != nil {
                return err
        }

        data, err := json.MarshalIndent(c, "", "  ")
        if err != nil {
                return err
        }

        return os.WriteFile(configFile, data, 0644)
}