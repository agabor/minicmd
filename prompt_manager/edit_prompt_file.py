import sys
import subprocess
from pathlib import Path

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
