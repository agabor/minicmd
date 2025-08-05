
import time

def call_claude(user_prompt, config, system_prompt, debug=False, attachments=None):
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
        
        # Build content array
        content = []
        
        # Add attachment files as separate messages
        if attachments:
            for attachment in attachments:
                content.append({"type": "text", "text": attachment, "cache_control": {"type": "ephemeral"}})
        
        # Add main user prompt
        content.append({"type": "text", "text": user_prompt})
        
        response = client.messages.create(
            model=config["claude_model"],
            max_tokens=4000,
            system=system_prompt,
            messages=[{"role": "user", "content": content}]
        )
        
        end_time = time.time()
        print(f"Claude API call took {end_time - start_time:.2f} seconds")
        print(f"Token usage - Input: {response.usage.input_tokens}, Output: {response.usage.output_tokens}, Cache Create: {response.usage.cache_creation_input_tokens}, Cache Read: {response.usage.cache_read_input_tokens}")
        
        # For debug mode, convert the response object to a JSON-serializable format
        raw_response = None
        if debug:
            raw_response = vars(response)
        
        return response.content[0].text, raw_response
    except Exception as e:
        print(f"Error calling Claude API: {e}")
        return None, None
