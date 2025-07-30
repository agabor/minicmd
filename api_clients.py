#!/usr/bin/env python3

import requests
import json
import time
from config import SYSTEM_PROMPT

def call_ollama(user_prompt, config, debug=False):
    """Make API call to Ollama"""
    start_time = time.time()
    
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
        
        end_time = time.time()
        print(f"Ollama API call took {end_time - start_time:.2f} seconds")
        
        raw_response = json.dumps(data, indent=2) if debug else None
        return data.get('response', ''), raw_response
    except requests.exceptions.RequestException as e:
        print(f"Error calling Ollama API: {e}")
        return None, None
    except json.JSONDecodeError as e:
        print(f"Error parsing JSON response: {e}")
        return None, None

def call_claude(user_prompt, config, debug=False):
    """Make API call to Claude"""
    if not config["anthropic_api_key"]:
        print("Error: Claude API key not configured.")
        print("Please set your API key with: python3 minicmd.py config anthropic_api_key YOUR_API_KEY")
        return None, None
    
    try:
        import anthropic
    except ImportError:
        print("Error: anthropic library is required for Claude API.")
        print("Install with: pip install anthropic")
        return None, None
    
    start_time = time.time()
    
    try:
        client = anthropic.Anthropic(api_key=config["anthropic_api_key"])
        
        response = client.messages.create(
            model=config["claude_model"],
            max_tokens=4000,
            system=SYSTEM_PROMPT,
            messages=[
                {"role": "user", "content": user_prompt}
            ]
        )
        
        end_time = time.time()
        print(f"Claude API call took {end_time - start_time:.2f} seconds")
        print(f"Token usage - Input: {response.usage.input_tokens}, Output: {response.usage.output_tokens}")
        
        # For debug mode, convert the response object to a JSON-serializable format
        raw_response = None
        if debug:
            raw_response = json.dumps(response, indent=2)
        
        return response.content[0].text, raw_response
    except Exception as e:
        print(f"Error calling Claude API: {e}")
        return None, None

def call_deepseek(user_prompt, config, debug=False):
    """Make API call to DeepSeek"""
    if not config["deepseek_api_key"]:
        print("Error: DeepSeek API key not configured.")
        print("Please set your API key with: python3 minicmd.py config deepseek_api_key YOUR_API_KEY")
        return None, None
    
    start_time = time.time()
    
    headers = {
        "Authorization": f"Bearer {config['deepseek_api_key']}",
        "Content-Type": "application/json"
    }
    
    payload = {
        "model": config["deepseek_model"],
        "messages": [
            {"role": "system", "content": SYSTEM_PROMPT},
            {"role": "user", "content": user_prompt}
        ],
        "max_tokens": 4000,
        "temperature": 0.1,
        "stream": False
    }
    
    try:
        response = requests.post(config["deepseek_url"], json=payload, headers=headers, timeout=60)
        response.raise_for_status()
        data = response.json()
        
        end_time = time.time()
        print(f"DeepSeek API call took {end_time - start_time:.2f} seconds")
        print(f"Token usage - Total: {data['usage']['total_tokens']}")
        
        raw_response = json.dumps(data, indent=2) if debug else None
        
        if 'choices' in data and len(data['choices']) > 0:
            return data['choices'][0]['message']['content'], raw_response
        else:
            print("Error: Unexpected response format from DeepSeek API")
            return None, None
            
    except requests.exceptions.RequestException as e:
        print(f"Error calling DeepSeek API: {e}")
        return None, None
    except json.JSONDecodeError as e:
        print(f"Error parsing JSON response: {e}")
        return None, None
    except KeyError as e:
        print(f"Error parsing response structure: {e}")
        return None, None
