# Go Implementation Notes

This document describes the Go reimplementation of the MiniCmd project and highlights key differences from the Python version.

## Implementation Status

✅ **Complete** - All features from the Python version have been reimplemented in Go.

## Project Structure Comparison

### Python Version
```
minicmd/
├── minicmd (main entry point)
├── config.py
├── file_processor.py
├── help.py
├── api_clients/
│   ├── claude_client.py
│   ├── deepseek_client.py
│   └── ollama_client.py
├── commands/
│   ├── add.py
│   ├── clear.py
│   ├── config.py
│   ├── edit.py
│   ├── list.py
│   └── run.py
└── prompt_manager/
    ├── add_file_to_prompt.py
    ├── edit_prompt_file.py
    ├── get_attachments.py
    └── get_prompt_from_file.py
```

### Go Version
```
go/
├── main.go
├── config/
│   └── config.go
├── fileprocessor/
│   └── fileprocessor.go
├── apiclient/
│   ├── claude.go
│   ├── deepseek.go
│   └── ollama.go
├── commands/
│   └── commands.go
└── promptmanager/
    └── promptmanager.go
```

## Key Differences

### 1. Module Organization
- **Python**: Uses separate files for each command handler
- **Go**: Consolidates command handlers in a single `commands.go` file for better cohesion

### 2. Error Handling
- **Python**: Uses `sys.exit()` and prints errors to stdout
- **Go**: Returns errors through the call stack and uses proper error types

### 3. Concurrency
- **Python**: Uses threading for progress indicator
- **Go**: Uses goroutines and channels for progress indicator (more idiomatic)

### 4. Dependencies
- **Python**: 
  - requests (HTTP client)
  - anthropic (Claude SDK)
- **Go**: 
  - github.com/anthropics/anthropic-sdk-go (Claude SDK)
  - github.com/spf13/pflag (command-line flags)
  - Standard library for HTTP (no external HTTP client needed)

### 5. Configuration
- Both versions store configuration in `~/.minicmd/config`
- Both use JSON format
- **Fully compatible** - same configuration file works for both implementations

### 6. Binary Distribution
- **Python**: Requires Python interpreter + pip packages
- **Go**: Single static binary, no dependencies required at runtime

## Feature Parity

All features from the Python version are implemented:

| Feature | Python | Go | Notes |
|---------|--------|-----|-------|
| Multiple AI providers | ✅ | ✅ | Claude, Ollama, DeepSeek |
| Configuration management | ✅ | ✅ | Compatible config format |
| File attachments | ✅ | ✅ | |
| Prompt editing (vim) | ✅ | ✅ | |
| Code block extraction | ✅ | ✅ | |
| Progress indicator | ✅ | ✅ | |
| Verbose mode | ✅ | ✅ | |
| Debug mode | ✅ | ✅ | |
| Safe mode | ✅ | ✅ | Adds .new suffix |
| Glob patterns for add | ✅ | ✅ | |
| Token usage reporting | ✅ | ✅ | |

## Advantages of Go Implementation

1. **Performance**: Faster startup and execution
2. **Deployment**: Single binary with no dependencies
3. **Memory**: Lower memory footprint
4. **Concurrency**: Better handling of concurrent operations with goroutines
5. **Type Safety**: Compile-time type checking prevents many runtime errors
6. **Cross-compilation**: Easy to build for different platforms

## Building and Installing

### Prerequisites
```bash
# Ubuntu/Debian
sudo snap install go
# or
sudo apt install golang-go
```

### Build
```bash
cd go
go mod download
go build -o minicmd
```

### Install
```bash
cd go
chmod +x install.sh
./install.sh
```

Or manually:
```bash
cd go
go build -o minicmd
sudo mv minicmd /usr/local/bin/
```

## Testing

To test the Go implementation:

```bash
# Build
cd go
go build -o minicmd

# Test help
./minicmd --help

# Test config
./minicmd config

# Test with a simple prompt (requires Ollama or API keys)
./minicmd run "create a hello world function"
```

## Compatibility Notes

1. **Configuration File**: The Go version uses the same config file format and location as Python (`~/.minicmd/config`)
2. **Prompt File**: Uses the same prompt file location (`~/.minicmd/prompt`)
3. **Attachments**: Uses the same attachments file (`~/.minicmd/attachments.json`)
4. **API Compatibility**: Same API endpoints and request formats

You can use both implementations interchangeably - they share the same configuration and data files.

## Future Enhancements

Potential improvements for the Go version:

1. Add unit tests
2. Add integration tests
3. Support for additional AI providers
4. WebSocket support for streaming responses
5. Configuration file hot-reloading
6. Plugin system for custom processors
