import sys
from pathlib import Path
from glob import glob
from prompt_manager import add_file_to_prompt

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
