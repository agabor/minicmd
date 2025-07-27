#!/usr/bin/env python3

import requests
import json
import re
import sys
import subprocess
from pathlib import Path

# Configuration
OLLAMA_URL = "http://localhost:11434/api/generate"
MODEL = "deepseek-coder-v2:16b"
SYSTEM_PROMPT = "IMPORTANT: answer with one or more code blocks only without explanation. The first line should be a comment containing the file path and name."

def call_ollama(user_prompt):
    """Make API call to Ollama"""
    payload = {
        "model": MODEL,
        "prompt": user_prompt,
        "system": SYSTEM_PROMPT,
        "stream": False
    }
    
    try:
        response = requests.post(OLLAMA_URL, json=payload, timeout=60)
        response.raise_for_status()
        data = response.json()
        return data.get('response', '')
    except requests.exceptions.RequestException as e:
        print(f"Error calling Ollama API: {e}")
        return None
    except json.JSONDecodeError as e:
        print(f"Error parsing JSON response: {e}")
        return None

def extract_filename_from_comment(line):
    """Extract filename from comment line"""
    # Match various comment styles: //, #, /* */, etc.
    patterns = [
        r'^\s*//\s*(.+?)(?:\s*//.*)?$',  # // filename
        r'^\s*#\s*(.+?)(?:\s*#.*)?$',    # # filename
        r'^\s*/\*\s*(.+?)\s*\*/$',       # /* filename */
        r'^\s*--\s*(.+?)(?:\s*--.*)?$',  # -- filename (SQL)
        r'^\s*<!--\s*(.+?)\s*-->$',      # <!-- filename --> (HTML)
    ]
    
    for pattern in patterns:
        match = re.match(pattern, line)
        if match:
            filename = match.group(1).strip()
            # Remove any trailing comment markers
            filename = re.sub(r'\s*\*+/$', '', filename)
            return filename
    return None

def process_code_blocks(response):
    """Extract code blocks and create files"""
    lines = response.split('\n')
    
    # Check if response has markdown code blocks
    if '```' in response:
        process_markdown_blocks(lines)
    else:
        process_raw_code(lines)

def process_markdown_blocks(lines):
    """Process markdown code blocks"""
    in_code_block = False
    file_path = ""
    content_lines = []
    
    for line in lines:
        if '```' in line:
            if in_code_block:
                # End of code block - create file
                if file_path and content_lines:
                    create_file(file_path, '\n'.join(content_lines))
                in_code_block = False
                file_path = ""
                content_lines = []
            else:
                # Start of code block
                in_code_block = True
        elif in_code_block:
            if not file_path:
                # Check if this line contains the filename
                extracted_path = extract_filename_from_comment(line)
                if extracted_path:
                    file_path = extracted_path
                    continue
            
            # Add to content (skip the filename comment line)
            if file_path and not extract_filename_from_comment(line):
                content_lines.append(line)

def process_raw_code(lines):
    """Process raw code (no markdown blocks)"""
    file_path = ""
    content_lines = []
    
    for i, line in enumerate(lines):
        if i == 0 or not file_path:
            # Check if this line contains the filename
            extracted_path = extract_filename_from_comment(line)
            if extracted_path:
                file_path = extracted_path
                continue
        
        # Add all other lines to content
        content_lines.append(line)
    
    if file_path and content_lines:
        # Remove empty lines from the end
        while content_lines and not content_lines[-1].strip():
            content_lines.pop()
        create_file(file_path, '\n'.join(content_lines))

def create_file(file_path, content):
    """Create file with given content"""
    try:
        # Create directory if it doesn't exist
        path_obj = Path(file_path)
        path_obj.parent.mkdir(parents=True, exist_ok=True)
        
        # Write content to file
        with open(file_path, 'w', encoding='utf-8') as f:
            f.write(content)
            if content and not content.endswith('\n'):
                f.write('\n')
        
        print(f"Created/Updated: {file_path}")
        
    except IOError as e:
        print(f"Error creating file {file_path}: {e}")

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

def get_prompt_from_file():
    """Read prompt from the prompt file"""
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

def main():
    # Check if user wants to edit the prompt file
    if len(sys.argv) > 1 and sys.argv[1] == "edit":
        edit_prompt_file()
        return
    
    # Check for unexpected arguments
    if len(sys.argv) > 1:
        print("Usage:")
        print("  python3 minicmd.py       # Use prompt from ~/.minicmd/prompt")
        print("  python3 minicmd.py edit  # Edit the prompt file")
        sys.exit(1)
    
    # Get prompt from file
    prompt = get_prompt_from_file()
    
    print("Sending request to Ollama...")
    print(f"Model: {MODEL}")
    print(f"Prompt: {prompt}")
    print("---")
    
    # Get response from Ollama
    response = call_ollama(prompt)
    
    if response is None:
        print("Error: No response from Ollama API")
        sys.exit(1)
    
    if not response.strip():
        print("Error: Empty response from Ollama API")
        sys.exit(1)
    
    # Echo the response to see what we got
    print("Raw response:")
    print("==============")
    print(response)
    print("==============")
    print()
    
    # Process the response and create files
    print("Processing response...")
    process_code_blocks(response)
    
    print("Done!")

if __name__ == "__main__":
    # Check if requests is available
    try:
        import requests
    except ImportError:
        print("Error: requests library is required. Install with: pip install requests")
        sys.exit(1)
    
    main()
