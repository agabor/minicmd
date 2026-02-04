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
go build -o y
```

This will create a `y` binary in the current directory.

## Installation

To install the binary to your system:

```bash
cd go
go build -o y
sudo mv y /usr/local/bin/
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
y --help

# Configure API key
y config anthropic_api_key YOUR_API_KEY

# Run code generation
y act
y act "create a hello world function"

# Edit prompt
y edit

# Add files to context
y read file.go

# List attachments
y list

# Clear prompt and attachments
y clear

# Show configuration
y config
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
