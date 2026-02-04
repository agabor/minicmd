# YACT

Yet Another AI Coding Tool

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
go build -o ya
```

This will create a `ya` binary in the current directory.

## Installation

To install the binary to your system:

```bash
cd go
go build -o ya
sudo mv ya /usr/local/bin/
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
ya --help

# Configure API key
ya config anthropic_api_key YOUR_API_KEY

# Run code generation
ya act
ya act "create a hello world function"

# Edit prompt
ya edit

# Add files to context
ya read file.go

# List attachments
ya list

# Clear prompt and attachments
ya clear

# Show configuration
ya config
```

## Features

All features from the Python implementation are supported:

- Code generation from prompts
- File attachment support
- Configuration management
- Progress indicators
- Verbose and debug modes
- Safe mode (adds .new suffix to generated files)

## Configuration

Configuration is stored in `~/.yact/config` (same location as Python version).

The Go implementation is compatible with the Python version's configuration.

## Dependencies

- `github.com/anthropics/anthropic-sdk-go` - Claude API client
- `github.com/spf13/pflag` - Command-line flag parsing

Dependencies are managed via Go modules and will be automatically downloaded when you run `go mod download` or `go build`.
