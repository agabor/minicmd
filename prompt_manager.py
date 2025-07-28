#!/usr/bin/env python3

import re
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

def add_file_to_prompt(file_path):
    """Add a file reference to the prompt file"""
    prompt_file = Path.home() / ".minicmd" / "prompt"
    
    # Create directory if it doesn't exist
    prompt_file.parent.mkdir(parents=True, exist_ok=True)
    
    # Create empty file if it doesn't exist
    if not prompt_file.exists():
        prompt_file.touch()
    
    # Add the file reference
    file_reference = f"[[ {file_path} ]]"
    
    try:
        # Read existing content
        existing_content = ""
        if prompt_file.exists():
            with open(prompt_file, 'r', encoding='utf-8') as f:
                existing_content = f.read().rstrip()
        
        # Append the file reference
        if existing_content:
            new_content = existing_content + "\n" + file_reference
        else:
            new_content = file_reference
        
        # Write back to file
        with open(prompt_file, 'w', encoding='utf-8') as f:
            f.write(new_content + "\n")
        
        print(f"Added file reference to prompt: {file_reference}")
        
    except IOError as e:
        print(f"Error updating prompt file: {e}")
        sys.exit(1)

def get_prompt_from_file(prompt_file_path=None):
    """Read prompt from the prompt file and resolve file references"""
    if prompt_file_path:
        prompt_file = Path(prompt_file_path)
    else:
        prompt_file = Path.home() / ".minicmd" / "prompt"
    
    if not prompt_file.exists():
        if prompt_file_path:
            print(f"Error: Prompt file '{prompt_file_path}' does not exist.")
        else:
            print("Error: Prompt file does not exist.")
            print("Please run 'python3 minicmd.py edit' to create and edit your prompt.")
        sys.exit(1)
    
    try:
        with open(prompt_file, 'r', encoding='utf-8') as f:
            content = f.read().strip()
        
        if not content:
            if prompt_file_path:
                print(f"Error: Prompt file '{prompt_file_path}' is empty.")
            else:
                print("Error: Prompt file is empty.")
                print("Please run 'python3 minicmd.py edit' to add content to your prompt.")
            sys.exit(1)
        
        # Resolve file references
        resolved_content = resolve_file_references(content)
        return resolved_content
        
    except IOError as e:
        print(f"Error reading prompt file: {e}")
        sys.exit(1)

def resolve_file_references(content):
    """Resolve [[ file_path ]] references in the content"""
    def replace_file_reference(match):
        file_path = match.group(1).strip()
        try:
            # Check if file exists
            if not Path(file_path).exists():
                return f"// {file_path}\n// Error: File not found"
            
            # Read file content
            with open(file_path, 'r', encoding='utf-8') as f:
                file_content = f.read().rstrip()
            
            # Return formatted content
            return f"// {file_path}\n{file_content}"
            
        except IOError as e:
            return f"// {file_path}\n// Error reading file: {e}"
    
    # Replace all [[ file_path ]] patterns
    pattern = r'\[\[\s*([^\]]+)\s*\]\]'
    resolved = re.sub(pattern, replace_file_reference, content)
    
    return resolved
