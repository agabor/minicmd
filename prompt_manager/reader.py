import sys
from pathlib import Path

def get_prompt_from_file():
    """Read raw prompt from the prompt file without resolving references"""
    prompt_file = Path.home() / ".minicmd" / "prompt"
    
    if not prompt_file.exists():
        print("Error: Prompt file does not exist.")
        print("Please run 'python3 minicmd.py edit' to create and edit your prompt.")
        sys.exit(1)
    
    try:
        with open(prompt_file, 'r', encoding='utf-8') as f:
            content = f.read().strip()
        
        if not content:
            print("Error: Prompt file is empty.")
            print("Please run 'python3 minicmd.py edit' to add content to your prompt.")
            sys.exit(1)
            
        return content
        
    except IOError as e:
        print(f"Error reading prompt file: {e}")
        sys.exit(1)
