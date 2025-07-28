#!/usr/bin/env python3

import sys
from config import load_config, save_config
from api_clients import call_claude, call_ollama
from file_processor import process_code_blocks
from prompt_manager import edit_prompt_file, add_file_to_prompt, get_prompt_from_file

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

def handle_edit_command():
    """Handle edit command"""
    edit_prompt_file()

def handle_add_command(args):
    """Handle add command"""
    if len(args) >= 1:
        file_path = args[0]
        add_file_to_prompt(file_path)
    else:
        print("Usage: python3 minicmd.py add <file>")
        sys.exit(1)

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
    print("  python3 minicmd.py add minicmd.py")
