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

def get_attachments():
    """Get attachment file contents as a list of strings"""
    attachments_file = Path.home() / ".minicmd" / "attachments.json"
    
    # Read attachments
    attachments = []
    if attachments_file.exists():
        try:
            with open(attachments_file, 'r', encoding='utf-8') as f:
                attachments = json.load(f)
        except (json.JSONDecodeError, IOError):
            attachments = []
    
    # Build file contents list
    file_contents = []
    
    for file_path in attachments:
        try:
            # Check if file exists
            if not Path(file_path).exists():
                file_contents.append(f"// {file_path}\n// Error: File not found")
                continue
            
            # Read file content
            with open(file_path, 'r', encoding='utf-8') as f:
                file_content = f.read().rstrip()
            
            # Add formatted content
            file_contents.append(f"// {file_path}\n{file_content}")
            
        except IOError as e:
            file_contents.append(f"// {file_path}\n// Error reading file: {e}")
    
    return file_contents
