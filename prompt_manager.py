#!/usr/bin/env python3

import re
import sys
import json
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

def get_resolved_prompt_from_file():
    """Read prompt from file and resolve all file references"""
    content = get_prompt_from_file()
    return resolve_file_references(content)

def resolve_file_references(content):
    """Read attachments.json and add file contents to the beginning of the prompt"""
    attachments_file = Path.home() / ".minicmd" / "attachments.json"
    
    # Read attachments
    attachments = []
    if attachments_file.exists():
        try:
            with open(attachments_file, 'r', encoding='utf-8') as f:
                attachments = json.load(f)
        except (json.JSONDecodeError, IOError):
            attachments = []
    
    # Build file contents section
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
    
    # Combine file contents with original prompt
    if file_contents:
        return "\n\n".join(file_contents) + "\n\n" + content
    else:
        return content
