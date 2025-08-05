
import requests
import json
import time
from config import SYSTEM_PROMPT

def call_deepseek(user_prompt, config, debug=False, attachments=None):
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
    
    # Build messages array
    messages = [
        {"role": "system", "content": SYSTEM_PROMPT}
    ]
    
    # Add attachment files as separate messages
    if attachments:
        for attachment in attachments:
            messages.append({"role": "user", "content": attachment})
    
    # Add main user prompt
    messages.append({"role": "user", "content": user_prompt})
    
    payload = {
        "model": config["deepseek_model"],
        "messages": messages,
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
        print(f"Token usage - Input: {data['usage']['prompt_tokens']}, Output: {data['usage']['completion_tokens']}, Cached: {data['usage']['prompt_tokens_details']['cached_tokens']}")
        
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
