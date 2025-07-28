#!/usr/bin/env python3

import json
from pathlib import Path

# Default configuration values
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
