import sys
import json
from pathlib import Path

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
