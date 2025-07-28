#!/usr/bin/env python3

import sys
import argparse
from commands import handle_run_command, handle_config_command, handle_edit_command, handle_add_command, handle_clear_command
from help import show_help

def main():
    parser = argparse.ArgumentParser(description='AI-powered code generation tool', add_help=False)
    parser.add_argument('--claude', action='store_true', help='Use Claude API')
    parser.add_argument('--ollama', action='store_true', help='Use Ollama API')
    parser.add_argument('--help', '-h', action='store_true', help='Show help message')
    parser.add_argument('command', nargs='?', help='Command to execute')
    parser.add_argument('args', nargs='*', help='Command arguments')
    
    args = parser.parse_args()
    
    if args.help or (args.command == "help"):
        show_help()
        return
    
    # Handle special commands
    if args.command == "edit":
        handle_edit_command()
        return
    
    if args.command == "add":
        handle_add_command(args.args)
        return
    
    if args.command == "config":
        handle_config_command(args.args)
        return
    
    if args.command == "run":
        handle_run_command(args.args, args.claude, args.ollama)
        return

    if args.command == "clear":
        handle_clear_command()
        return
        
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
        print("Error: requests library is required.")
        print("Install with: pip install requests")
        sys.exit(1)
    
    main()
