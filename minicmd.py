#!/usr/bin/env python3

import requests
import json
import re
import sys
import subprocess
import argparse
from pathlib import Path

# Configuration
OLLAMA_URL = "http://localhost:11434/api/generate"
OLLAMA_MODEL = "deepseek-coder-v2:16b"
CLAUDE_MODEL = "claude-3-5-sonnet-20241022"
SYSTEM_PROMPT = "IMPORTANT: answer with one or more code blocks only without explanation. The first line should be a comment containing the file path and name."
CONFIG_DIR = Path.home() / ".minicmd"
CONFIG_FILE = CONFIG_DIR / "config"

def load_config():
    """Load configuration from config file"""
    default_config = {
        "default_provider": "ollama",
        "claude_api_key": "",
        "ollama_url": OLLAMA_URL,
        "ollama_model": OLLAMA_MODEL,
        "claude_model": CLAUDE_MODEL
    }
    
    if not CONFIG_FILE.exists():
        return default_config
    
    try:
        with open(CONFIG_FILE, 'r', encoding='utf-8') as f:
            config = json.load(f)
        # Merge with defaults to handle missing keys
        for key, value in default_config.items():
            if key not in config:
                config[key] = value
        return config
    except (IOError, json.JSONDecodeError) as e:
        print(f"Error loading config: {e}")
        return default_config

def save_config(config):
    """Save configuration to config file"""
    try:
        CONFIG_DIR.mkdir(parents=True, exist_ok=True)
        with open(CONFIG_FILE, 'w', encoding='utf-8') as f:
            json.dump(config, f, indent=2)
    except IOError as e:
        print(f"Error saving config: {e}")

def call_ollama(user_prompt, config):
    """Make API call to Ollama"""
    payload = {
        "model": config["ollama_model"],
        "prompt": user_prompt,
        "system": SYSTEM_PROMPT,
        "stream": False
    }
    
    try:
        response = requests.post(config["ollama_url"], json=payload, timeout=60)
        response.raise_for_status()
        data = response.json()
        return data.get('response', '')
    except requests.exceptions.RequestException as e:
        print(f"Error calling Ollama API: {e}")
        return None
    except json.JSONDecodeError as e:
        print(f"Error parsing JSON response: {e}")
        return None

def call_claude(user_prompt, config):
    """Make API call to Claude"""
    if not config["claude_api_key"]:
        print("Error: Claude API key not configured.")
        print("Please set your API key with: python3 minicmd.py config claude_api_key YOUR_API_KEY")
        return None
    
    try:
        import anthropic
    except ImportError:
        print("Error: anthropic library is required for Claude API. Install with: pip install anthropic")
        return None
    
    try:
        client = anthropic.Anthropic(api_key=config["claude_api_key"])
        
        response = client.messages.create(
            model=config["claude_model"],
            max_tokens=4000,
            system=SYSTEM_PROMPT,
            messages=[
                {"role": "user", "content": user_prompt}
            ]
        )
        
        return response.content[0].text
    except Exception as e:
        print(f"Error calling Claude API: {e}")
        return None

def extract_filename_from_comment(line):
    """Extract filename from comment line"""
    # Match various comment styles: //, #, /* */, etc.
    patterns = [
        r'^\s*//\s*(.+?)(?:\s*//.*)?$',  # // filename
        r'^\s*#\s*(.+?)(?:\s*#.*)?$',    # # filename
        r'^\s*/\*\s*(.+?)\s*\*/$',       # /* filename */
        r'^\s*--\s*(.+?)(?:\s*--.*)?$',  # -- filename (SQL)
        r'^\s*<!--\s*(.+?)\s*-->$',      # <!-- filename --> (HTML)
    ]
    
    for pattern in patterns:
        match = re.match(pattern, line)
        if match:
            filename = match.group(1).strip()
            # Remove any trailing comment markers
            filename = re.sub(r'\s*\*+/$', '', filename)
            return filename
    return None

def process_code_blocks(response):
    """Extract code blocks and create files"""
    lines = response.split('\n')
    
    # Check if response has markdown code blocks
    if '```' in response:
        process_markdown_blocks(lines)
    else:
        process_raw_code(lines)

def process_markdown_blocks(lines):
    """Process markdown code blocks"""
    in_code_block = False
    file_path = ""
    content_lines = []
    
    for line in lines:
        if '```' in line:
            if in_code_block:
                # End of code block - create file
                if file_path and content_lines:
                    create_file(file_path, '\n'.join(content_lines))
                in_code_block = False
                file_path = ""
                content_lines = []
            else:
                # Start of code block
                in_code_block = True
        elif in_code_block:
            if not file_path:
                # Check if this line contains the filename
                extracted_path = extract_filename_from_comment(line)
                if extracted_path:
                    file_path = extracted_path
                    continue
            
            # Add to content (skip the filename comment line)
            if file_path and not extract_filename_from_comment(line):
                content_lines.append(line)

def process_raw_code(lines):
    """Process raw code (no markdown blocks)"""
    file_path = ""
    content_lines = []
    
    for i, line in enumerate(lines):
        if i == 0 or not file_path:
            # Check if this line contains the filename
            extracted_path = extract_filename_from_comment(line)
            if extracted_path:
                file_path = extracted_path
                continue
        
        # Add all other lines to content
        content_lines.append(line)
    
    if file_path and content_lines:
        # Remove empty lines from the end
        while content_lines and not content_lines[-1].strip():
            content_lines.pop()
        create_file(file_path, '\n'.join(content_lines))

def create_file(file_path, content):
    """Create file with given content"""
    try:
        # Create directory if it doesn't exist
        path_obj = Path(file_path)
        path_obj.parent.mkdir(parents=True, exist_ok=True)
        
        # Write content to file
        with open(file_path, 'w', encoding='utf-8') as f:
            f.write(content)
            if content and not content.endswith('\n'):
                f.write('\n')
        
        print(f"Created/Updated: {file_path}")
        
    except IOError as e:
        print(f"Error creating file {file_path}: {e}")

def edit_prompt_file():
    """Launch vim to edit the prompt file"""
    prompt_file = Path.home() / ".minicmd" / "prompt"
    
    # Create directory if it doesn't exist
    prompt_file.parent.mkdir(parents=True, exist_ok=True)
    
    # Create empty file if it doesn't exist
    if not prompt_file.exists():
        prompt_file.touch()
    
    # Launch vim to edit the file
    try:
        subprocess.run(["vim", str(prompt_file)], check=True)
        print(f"Prompt file edited: {prompt_file}")
    except subprocess.CalledProcessError as e:
        print(f"Error launching vim: {e}")
        sys.exit(1)
    except FileNotFoundError:
        print("Error: vim not found. Please install vim or ensure it's in your PATH.")
        sys.exit(1)

def add_file_to_prompt(file_path):
    """Add a file reference to the prompt file"""
    prompt_file = Path.home() / ".minicmd" / "prompt"
    
    # Create directory if it doesn't exist
    prompt_file.parent.mkdir(parents=True, exist_ok=True)
    
    # Create empty file if it doesn't exist
    if not prompt_file.exists():
        prompt_file.touch()
    
    # Add the file reference
    file_reference = f"[[ {file_path} ]]"
    
    try:
        # Read existing content
        existing_content = ""
        if prompt_file.exists():
            with open(prompt_file, 'r', encoding='utf-8') as f:
                existing_content = f.read().rstrip()
        
        # Append the file reference
        if existing_content:
            new_content = existing_content + "\n" + file_reference
        else:
            new_content = file_reference
        
        # Write back to file
        with open(prompt_file, 'w', encoding='utf-8') as f:
            f.write(new_content + "\n")
        
        print(f"Added file reference to prompt: {file_reference}")
        
    except IOError as e:
        print(f"Error updating prompt file: {e}")
        sys.exit(1)

def get_prompt_from_file(prompt_file_path=None):
    """Read prompt from the prompt file and resolve file references"""
    if prompt_file_path:
        prompt_file = Path(prompt_file_path)
    else:
        prompt_file = Path.home() / ".minicmd" / "prompt"
    
    if not prompt_file.exists():
        if prompt_file_path:
            print(f"Error: Prompt file '{prompt_file_path}' does not exist.")
        else:
            print("Error: Prompt file does not exist.")
            print("Please run 'python3 minicmd.py edit' to create and edit your prompt.")
        sys.exit(1)
    
    try:
        with open(prompt_file, 'r', encoding='utf-8') as f:
            content = f.read().strip()
        
        if not content:
            if prompt_file_path:
                print(f"Error: Prompt file '{prompt_file_path}' is empty.")
            else:
                print("Error: Prompt file is empty.")
                print("Please run 'python3 minicmd.py edit' to add content to your prompt.")
            sys.exit(1)
        
        # Resolve file references
        resolved_content = resolve_file_references(content)
        return resolved_content
        
    except IOError as e:
        print(f"Error reading prompt file: {e}")
        sys.exit(1)

def resolve_file_references(content):
    """Resolve [[ file_path ]] references in the content"""
    import re
    
    def replace_file_reference(match):
        file_path = match.group(1).strip()
        try:
            # Check if file exists
            if not Path(file_path).exists():
                return f"// {file_path}\n// Error: File not found"
            
            # Read file content
            with open(file_path, 'r', encoding='utf-8') as f:
                file_content = f.read().rstrip()
            
            # Return formatted content
            return f"// {file_path}\n{file_content}"
            
        except IOError as e:
            return f"// {file_path}\n// Error reading file: {e}"
    
    # Replace all [[ file_path ]] patterns
    pattern = r'\[\[\s*([^\]]+)\s*\]\]'
    resolved = re.sub(pattern, replace_file_reference, content)
    
    return resolved

def handle_run_command(args, claude_flag, ollama_flag):
    """Handle run command with optional prompt content parameter"""
    # Check for conflicting provider options
    if claude_flag and ollama_flag:
        print("Error: Cannot specify both --claude and --ollama")
        sys.exit(1)
    
    # Load configuration
    config = load_config()
    
    # Determine which provider to use
    if claude_flag:
        provider = "claude"
    elif ollama_flag:
        provider = "ollama"
    else:
        provider = config["default_provider"]
    
    # Get prompt content from args if provided, otherwise use default prompt file
    if len(args) > 0:
        # Use provided prompt content directly
        prompt = " ".join(args)
        print("Using provided prompt content")
    else:
        # Use default prompt file
        prompt = get_prompt_from_file()
        print("Using default prompt file")
    
    print(f"Sending request to {provider.title()}...")
    if provider == "claude":
        print(f"Model: {config['claude_model']}")
        response = call_claude(prompt, config)
    else:
        print(f"Model: {config['ollama_model']}")
        response = call_ollama(prompt, config)
    
    print(f"Prompt: {prompt}")
    print("---")
    
    if response is None:
        print(f"Error: No response from {provider.title()} API")
        sys.exit(1)
    
    if not response.strip():
        print(f"Error: Empty response from {provider.title()} API")
        sys.exit(1)
    
    # Echo the response to see what we got
    print("Raw response:")
    print("==============")
    print(response)
    print("==============")
    print()
    
    # Process the response and create files
    print("Processing response...")
    process_code_blocks(response)
    
    print("Done!")

def handle_config_command(args):
    """Handle config command"""
    config = load_config()
    
    if len(args) == 0:
        # Show current config
        print("Current configuration:")
        for key, value in config.items():
            if key == "claude_api_key" and value:
                print(f"  {key}: {'*' * len(value)}")  # Hide API key
            else:
                print(f"  {key}: {value}")
        return
    
    if len(args) == 2:
        key, value = args
        if key in config:
            config[key] = value
            save_config(config)
            if key == "claude_api_key":
                print(f"Set {key} to {'*' * len(value)}")
            else:
                print(f"Set {key} to {value}")
        else:
            print(f"Error: Unknown config key '{key}'")
            print("Available keys:", ", ".join(config.keys()))
    else:
        print("Usage:")
        print("  python3 minicmd.py config                    # Show current config")
        print("  python3 minicmd.py config <key> <value>      # Set config value")

def show_help():
    """Show help message"""
    print("minicmd - AI-powered code generation tool")
    print()
    print("Usage:")
    print("  python3 minicmd.py [--claude|--ollama]       # Generate code using AI (default prompt)")
    print("  python3 minicmd.py run [prompt_content] [--claude|--ollama]  # Generate code with optional custom prompt content")
    print("  python3 minicmd.py edit                      # Edit the prompt file")
    print("  python3 minicmd.py add <file>                # Add file reference to prompt")
    print("  python3 minicmd.py config                    # Show current configuration")
    print("  python3 minicmd.py config <key> <value>      # Set configuration value")
    print()
    print("Options:")
    print("  --claude    Use Claude API (requires API key)")
    print("  --ollama    Use Ollama API (requires local Ollama)")
    print()
    print("Configuration keys:")
    print("  default_provider    Default AI provider (claude or ollama)")
    print("  claude_api_key      Claude API key")
    print("  claude_model        Claude model name")
    print("  ollama_url          Ollama API URL")
    print("  ollama_model        Ollama model name")
    print()
    print("Examples:")
    print("  python3 minicmd.py config claude_api_key sk-ant-...")
    print("  python3 minicmd.py config default_provider claude")
    print("  python3 minicmd.py --claude")
    print("  python3 minicmd.py run")
    print("  python3 minicmd.py run \"create a hello world function\"")
    print("  python3 minicmd.py run \"write a Python calculator\" --claude")
    print("  python3 minicmd.py add main.py")

def main():
    parser = argparse.ArgumentParser(description='AI-powered code generation tool', add_help=False)
    parser.add_argument('--claude', action='store_true', help='Use Claude API')
    parser.add_argument('--ollama', action='store_true', help='Use Ollama API')
    parser.add_argument('--help', '-h', action='store_true', help='Show help message')
    parser.add_argument('command', nargs='?', help='Command to execute')
    parser.add_argument('args', nargs='*', help='Command arguments')
    
    args = parser.parse_args()
    
    # Show help if requested
    if args.help or (args.command == "help"):
        show_help()
        return
    
    # Handle special commands
    if args.command == "edit":
        edit_prompt_file()
        return
    
    if args.command == "add" and len(args.args) >= 1:
        file_path = args.args[0]
        add_file_to_prompt(file_path)
        return
    
    if args.command == "config":
        handle_config_command(args.args)
        return
    
    if args.command == "run":
        handle_run_command(args.args, args.claude, args.ollama)
        return
    
    # If no command specified, use the default behavior (backward compatibility)
    if args.command is None:
        handle_run_command([], args.claude, args.ollama)
        return
    
    # Unknown command
    print(f"Error: Unknown command '{args.command}'")
    print("Run 'python3 minicmd.py --help' for usage information.")
    sys.exit(1)

if __name__ == "__main__":
    # Check if requests is available
    try:
        import requests
    except ImportError:
        print("Error: requests library is required. Install with: pip install requests")
        sys.exit(1)
    
    main()
