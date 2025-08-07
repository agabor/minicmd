import sys
from config import load_config, save_config

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
