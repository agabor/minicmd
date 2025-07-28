# MiniCmd - AI-Powered Code Generation Tool

MiniCmd is a command-line tool that uses AI (Ollama, Claude, or DeepSeek) to generate code files based on prompts.

## Project Structure

The project has been refactored into modular components for better maintainability:

### Core Modules

- **`minicmd.py`** - Main application entry point and argument parsing
- **`config.py`** - Configuration management (loading/saving settings)
- **`api_clients.py`** - API client implementations for Ollama, Claude, and DeepSeek
- **`file_processor.py`** - Code block processing and file creation
- **`prompt_manager.py`** - Prompt file management and file reference resolution
- **`commands.py`** - Command handlers for different operations

### Legacy File

- **`minicmd.py`** - Original monolithic implementation (kept for reference)

## Usage

Run the tool using the new modular entry point:

```bash
# Use the new modular version
python3 minicmd.py [options] [command] [args]

# Examples:
python3 minicmd.py --help
python3 minicmd.py config
python3 minicmd.py edit
python3 minicmd.py run "create a hello world function"
python3 minicmd.py --claude
python3 minicmd.py --deepseek
```

## Commands

- **`run [prompt]`** - Generate code using AI (with optional custom prompt)
- **`edit`** - Edit the prompt file using vim
- **`add <file>`** - Add file reference to prompt
- **`config [key] [value]`** - Show or set configuration values
- **`help`** - Show help message

## Options

- **`--claude`** - Use Claude API
- **`--ollama`** - Use Ollama API (default)
- **`--deepseek`** - Use DeepSeek API

## Configuration

Configuration is stored in `~/.minicmd/config` and includes:

- `default_provider` - Default AI provider (claude, ollama, or deepseek)
- `anthropic_api_key` - Claude API key
- `deepseek_api_key` - DeepSeek API key
- `claude_model` - Claude model name
- `ollama_url` - Ollama API URL
- `ollama_model` - Ollama model name
- `deepseek_url` - DeepSeek API URL
- `deepseek_model` - DeepSeek model name

## Dependencies

- `requests` - For HTTP API calls
- `anthropic` - For Claude API (optional, only needed if using Claude)

Install dependencies:
```bash
pip install requests
pip install anthropic  # Only if using Claude
```

## Migration from Original

The new modular structure maintains full backward compatibility. You can:

1. Continue using `python3 minicmd.py` (original version)
2. Switch to `python3 minicmd.py` (new modular version)

Both versions provide identical functionality.
