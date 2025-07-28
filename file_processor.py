#!/usr/bin/env python3

import re
from pathlib import Path

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
        # Treat paths starting with '/' as relative to current working directory
        if file_path.startswith('/'):
            file_path = file_path[1:]  # Remove leading slash
        
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
