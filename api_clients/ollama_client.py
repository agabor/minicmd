
import requests
import json
import time
from config import SYSTEM_PROMPT

def call_ollama(user_prompt, config, debug=False, attachments=None):
    """Make API call to Ollama"""
    start_time = time.time()
    
    # For Ollama, combine attachments with user prompt as it doesn't support multiple messages
    full_prompt = user_prompt
    if attachments:
        full_prompt = "\n\n".join(attachments) + "\n\n" + user_prompt
    
    payload = {
        "model": config["ollama_model"],
        "prompt": full_prompt,
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
