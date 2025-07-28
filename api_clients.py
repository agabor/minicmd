#!/usr/bin/env python3

import requests
import json
from config import SYSTEM_PROMPT

def call_ollama(user_prompt, config):
    """Make API call to Ollama"""
    payload = {
        "model": config["ollama_model"],
        "prompt": user_prompt,
        "system": SYSTEM_PROMPT,
        "stream": False
    }
    
    try:
        response = requests.post(config["ollama_url"], json=payload, timeout=60)
        response.raise_for_status()
        data = response.json()
        return data.get('response', '')
    except requests.exceptions.RequestException as e:
        print(f"Error calling Ollama API: {e}")
        return None
    except json.JSONDecodeError as e:
        print(f"Error parsing JSON response: {e}")
        return None

def call_claude(user_prompt, config):
    """Make API call to Claude"""
    if not config["claude_api_key"]:
        print("Error: Claude API key not configured.")
        print("Please set your API key with: python3 minicmd.py config claude_api_key YOUR_API_KEY")
        return None
    
    try:
        import anthropic
    except ImportError:
        print("Error: anthropic library is required for Claude API.")
        print("Install with: pip install anthropic")
        return None
    
    try:
        client = anthropic.Anthropic(api_key=config["claude_api_key"])
        
        response = client.messages.create(
            model=config["claude_model"],
            max_tokens=4000,
            system=SYSTEM_PROMPT,
            messages=[
                {"role": "user", "content": user_prompt}
            ]
        )
        
        return response.content[0].text
    except Exception as e:
        print(f"Error calling Claude API: {e}")
        return None
