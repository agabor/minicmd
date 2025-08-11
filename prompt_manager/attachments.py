import sys
import json
from pathlib import Path

def add_file_to_prompt(file_path):
    """Add a file reference to attachments.json"""
    attachments_file = Path.home() / ".minicmd" / "attachments.json"
    
    # Create directory if it doesn't exist
    attachments_file.parent.mkdir(parents=True, exist_ok=True)
    
    # Load existing attachments
    attachments = []
    if attachments_file.exists():
        try:
            with open(attachments_file, 'r', encoding='utf-8') as f:
                attachments = json.load(f)
        except (json.JSONDecodeError, IOError):
            attachments = []
    
    # Add file if not already present
    if file_path not in attachments:
        attachments.append(file_path)
        
        # Save updated attachments
        try:
            with open(attachments_file, 'w', encoding='utf-8') as f:
                json.dump(attachments, f, indent=2)
            
            print(f"Added file to attachments: {file_path}")
            
        except IOError as e:
            print(f"Error updating attachments file: {e}")
            sys.exit(1)
    else:
        print(f"File already in attachments: {file_path}")
