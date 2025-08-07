import sys
import json
from pathlib import Path

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
