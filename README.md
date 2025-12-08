# MiniCmd Go Implementation

This is a Go reimplementation of the MiniCmd AI-powered code generation tool.

## Prerequisites

You need to have Go installed on your system. To install Go:

### Ubuntu/Debian
```bash
sudo snap install go
# or
sudo apt install golang-go
```

### Other Systems
Download and install from: https://golang.org/dl/

## Building

```bash
cd go
go mod download
go build -o minicmd
```

This will create a `minicmd` binary in the current directory.

## Installation

To install the binary to your system:

```bash
cd go
go build -o minicmd
sudo mv minicmd /usr/local/bin/
```

Or use the provided install script:

```bash
cd go
chmod +x install.sh
./install.sh
```

## Usage

The Go implementation has the same interface as the Python version:

```bash
# Show help
minicmd --help

# Configure API keys
minicmd config anthropic_api_key YOUR_API_KEY
minicmd config deepseek_api_key YOUR_API_KEY

# Run code generation
minicmd run
minicmd run "create a hello world function"

# Edit prompt
minicmd edit

# Add files to context
minicmd read file.go

# List attachments
minicmd list

# Clear prompt and attachments
minicmd clear

# Show configuration
minicmd config
```

## Project Structure

```
go/
├── main.go                 # Main entry point
├── config/                 # Configuration management
│   └── config.go
├── apiclient/             # API client implementations
│   ├── claude.go
│   ├── deepseek.go
│   └── ollama.go
├── commands/              # Command handlers
│   └── commands.go
├── fileprocessor/         # Code block processing
│   └── fileprocessor.go
├── promptmanager/         # Prompt and attachment management
│   └── promptmanager.go
├── go.mod                 # Go module definition
└── README.md             # This file
```

## Features

All features from the Python implementation are supported:

- Multiple AI providers (Claude, Ollama, DeepSeek)
- Code generation from prompts
- File attachment support
- Configuration management
- Progress indicators
- Verbose and debug modes
- Safe mode (adds .new suffix to generated files)

## Configuration

Configuration is stored in `~/.minicmd/config` (same location as Python version).

The Go implementation is compatible with the Python version's configuration.

## Dependencies

- `github.com/anthropics/anthropic-sdk-go` - Claude API client
- `github.com/spf13/pflag` - Command-line flag parsing

Dependencies are managed via Go modules and will be automatically downloaded when you run `go mod download` or `go build`.
