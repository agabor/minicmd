#!/usr/bin/env python3

import sys
import time
import threading
import json
from pathlib import Path
from glob import glob
from config import load_config, save_config, SYSTEM_PROMPT
from api_clients import call_claude, call_ollama, call_deepseek
from file_processor import process_code_blocks
from prompt_manager import edit_prompt_file, add_file_to_prompt, get_prompt_from_file, get_attachments

def show_progress():
    """Show a simple progress indicator"""
    chars = "⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏"
    idx = 0
    while True:
        print(f"\r{chars[idx % len(chars)]}", end="", flush=True)
        idx += 1
        time.sleep(0.1)

def handle_run_command(args, claude_flag, ollama_flag, deepseek_flag, verbose, debug):
    """Handle run command with optional prompt content parameter"""
    # Check for conflicting provider options
    provider_flags = [claude_flag, ollama_flag, deepseek_flag]
    if sum(provider_flags) > 1:
        print("Error: Cannot specify multiple provider flags")
        sys.exit(1)
    
    # Load configuration
    config = load_config()
    
    # Determine which provider to use
    if claude_flag:
        provider = "claude"
    elif ollama_flag:
        provider = "ollama"
    elif deepseek_flag:
        provider = "deepseek"
    else:
        provider = config["default_provider"]
    
    # Get prompt content from args if provided, otherwise use default prompt file
    if len(args) > 0:
        # Use provided prompt content directly
        prompt = " ".join(args)
        print("Using provided prompt content")
    else:
        # Use default prompt file
        edit_prompt_file()
        prompt = get_prompt_from_file()
        print("Using default prompt file")
    
    if verbose:
        print(f"Prompt: {prompt}")
        print("---")
    
    print(f"Sending request to {provider.title()}...")
    if provider == "claude":
        print(f"Model: {config['claude_model']}")
    elif provider == "deepseek":
        print(f"Model: {config['deepseek_model']}")
    else:
        print(f"Model: {config['ollama_model']}")
    
    # Start progress indicator
    progress_thread = threading.Thread(target=show_progress, daemon=True)
    progress_thread.start()
    
    try:
        if provider == "claude":
            response, raw_response = call_claude(prompt, config, SYSTEM_PROMPT, debug, get_attachments())
        elif provider == "deepseek":
            response, raw_response = call_deepseek(prompt, config, SYSTEM_PROMPT, debug, get_attachments())
        else:
            response, raw_response = call_ollama(prompt, config, SYSTEM_PROMPT, debug, get_attachments())
    finally:
        # Clear progress indicator
        print("\r ", end="", flush=True)
        print("\r", end="", flush=True)
    
    handle_clear_command()

    if response is None:
        print(f"Error: No response from {provider.title()} API")
        sys.exit(1)
    
    if not response.strip():
        print(f"Error: Empty response from {provider.title()} API")
        sys.exit(1)
    
    # Echo the response to see what we got
    if verbose:
        print("Raw response:")
        print("==============")
        print(response)
        print("==============")
        print()
    
    # In debug mode, also show the complete API response as JSON
    if debug and raw_response is not None:
        print("Complete API response (JSON):")
        print("=============================")
        print(raw_response)
        print("=============================")
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
            if key in ["anthropic_api_key", "deepseek_api_key"] and value:
                print(f"  {key}: {'*' * len(value)}")  # Hide API key
            else:
                print(f"  {key}: {value}")
        return
    
    if len(args) == 2:
        key, value = args
        if key in config:
            config[key] = value
            save_config(config)
            if key in ["anthropic_api_key", "deepseek_api_key"]:
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
        for pattern in args:
            files = glob(pattern)
            if not files:
                print(f"No files found matching pattern: {pattern}")
                continue
            for file_path in files:
                if Path(file_path).is_file():
                    add_file_to_prompt(file_path)
                else:
                    print(f"Skipping directory: {file_path}")
    else:
        print("Usage: python3 minicmd.py add <file> [<file2> ...]")
        sys.exit(1)

def handle_list_command():
    """Handle list command"""
    attachments_file = Path.home() / ".minicmd" / "attachments.json"
    
    # Check if attachments file exists
    if not attachments_file.exists():
        print("No attachments found.")
        return
    
    # Load attachments
    try:
        with open(attachments_file, 'r', encoding='utf-8') as f:
            attachments = json.load(f)
    except (json.JSONDecodeError, IOError) as e:
        print(f"Error reading attachments file: {e}")
        sys.exit(1)
    
    # Display attachments
    if not attachments:
        print("No attachments found.")
    else:
        print("Current attachments:")
        for i, file_path in enumerate(attachments, 1):
            # Check if file exists
            if Path(file_path).exists():
                print(f"  {i}. {file_path}")
            else:
                print(f"  {i}. {file_path} (file not found)")

def handle_clear_command():
    """Handle clear command"""
    prompt_file = Path.home() / ".minicmd" / "prompt"
    attachments_file = Path.home() / ".minicmd" / "attachments.json"
    
    try:
        prompt_file.parent.mkdir(parents=True, exist_ok=True)
        prompt_file.write_text('')
        if attachments_file.exists():
            attachments_file.unlink()
        print("Cleared prompt and attachments")
    except IOError as e:
        print(f"Error clearing files: {e}")
        sys.exit(1)
