#!/usr/bin/env python3

import re
from pathlib import Path

def extract_filename_from_comment(line):
    """Extract filename from comment line"""
    # Match various comment styles: //, #, /* */, etc.
    patterns = [
        r'^\s*//\s*(.+?)(?:\s*//.*)?$',  # // filename
        r'^\s*#\s*(.+?)(?:\s*#.*)?$',    # # filename
        r'^\s*#\s*//\s*(.+?)(?:\s*#.*)?$',    # # // filename
        r'^\s*/\*\s*(.+?)\s*\*/$',       # /* filename */
        r'^\s*--\s*(.+?)(?:\s*--.*)?$',  # -- filename (SQL)
        r'^\s*<!--\s*(.+?)\s*-->$',      # <!-- filename --> (HTML)
    ]
    
    for pattern in patterns:
        match = re.match(pattern, line)
        if match:
            filename = match.group(1).strip()
            if "!" in filename:
                continue
            # Remove any trailing comment markers
            filename = re.sub(r'\s*\*+/$', '', filename)
            return filename
    return None

def process_code_blocks(response, safe=False):
    """Extract code blocks and create files"""
    lines = response.split('\n')
    
    process_markdown_blocks(lines, safe)

def process_markdown_blocks(lines, safe):
    """Process markdown code blocks"""
    in_code_block = False
    file_path = ""
    content_lines = []
    
    for line in lines:
        if '```' in line:
            if in_code_block:
                # End of code block - create file
                if file_path and content_lines:
                    create_file(file_path + '.new' if safe else file_path, '\n'.join(content_lines))
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
            content_lines.append(line)

def create_file(file_path, content):
    """Create file with given content"""
    try:
        # Handle absolute paths
        if file_path.startswith('/'):
            # First check if file exists as absolute path within working directory
            rel_path = file_path[1:]  # Remove leading slash
            if Path(rel_path).exists():
                file_path = rel_path
            else:
                # Try using absolute path if file exists
                abs_path = Path(file_path)
                if abs_path.exists():
                    file_path = str(abs_path)
                else:
                    # Default to relative path if file doesn't exist
                    file_path = rel_path
        
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
