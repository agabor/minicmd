import json
from pathlib import Path

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
