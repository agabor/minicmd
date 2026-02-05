# Project File Structure

```
.
├── main.go
│   Entry point for the application. Handles CLI flag parsing, stdin input processing,
│   and command routing. Supports commands like act, bash, ask, plan, read, list, config,
│   clear, new, and last.
│
├── api/
│   ├── client.go
│   │   Defines the APIClient interface and Message struct. Provides abstraction for
│   │   different AI provider implementations.
│   │
│   └── claude.go
│       Implementation of Claude API client. Handles authentication, message formatting,
│       token calculation, cost estimation, and API communication with Anthropic's SDK.
│
├── commands/
│   ├── call.go
│   │   Handles the main API call flow. Implements HandleActCommand, HandleAskCommand,
│   │   HandleNewCommand, and HandleCall. Shows progress spinner during API requests.
│   │
│   ├── config.go
│   │   Manages configuration commands. Displays current config and allows users to
│   │   set anthropic_api_key and claude_model values.
│   │
│   ├── files.go
│   │   Handles file operations: reading files (HandleReadCommand), listing attachments
│   │   (HandleListCommand), clearing attachments (HandleClearCommand), and retrieving
│   │   last AI response (HandleLastCommand).
│   │
│   └── help.go
│       Displays help message showing all available commands, options, and configuration keys.
│
├── config/
│   ├── config.go
│   │   Configuration management. Defines Config struct, handles loading/saving from
│   │   ~/.yact/config file, and provides default configuration values.
│   │
│   └── systemprompt/
│       ├── act.go
│       │   System prompt for code generation. Defines strict output rules, code block
│       │   format requirements, and code quality guidelines.
│       │
│       ├── ask.go
│       │   System prompt for codebase question answering. Guides AI to provide
│       │   explanations with code references and architectural context.
│       │
│       ├── bash.go
│       │   System prompt for bash script generation. Specifies shebang, error handling,
│       │   and script structure requirements.
│       │
│       └── plan.go
│           System prompt for planning and analysis. Helps break requirements into
│           actionable steps and component specifications.
│
└── logic/
    ├── codeblock.go
    │   Code block parsing and file writing. Extracts filenames from code block comments,
    │   manages file paths with safe mode (.new suffix), and writes generated code to disk.
    │
    ├── contextmanager.go
    │   Manages conversation context stored in ~/.yact/context.json. Handles loading,
    │   saving, and clearing message history. Builds messages with attachments.
    │
    ├── fileprocessor.go
    │   Processes markdown code blocks from AI responses. Parses code blocks and
    │   delegates file writing to CodeBlock struct.
    │
    └── promptmanager.go
        Manages file attachments. Stores/loads attachment list from ~/.yact/attachments.json.
        Retrieves file contents for inclusion in prompts and displays attachment status.
```

## Key Features

- **Code Generation**: `act` command generates code based on prompts
- **Bash Scripts**: `bash` command generates shell scripts
- **Q&A**: `ask` command answers questions about codebase
- **Planning**: `plan` command breaks down requirements
- **File Management**: Attach files to prompts with `read` and `list` commands
- **Context**: Maintains conversation history for multi-turn interactions
- **Configuration**: Set API keys and model selection
- **Safe Mode**: Optional `.new` suffix for generated files to prevent overwrites