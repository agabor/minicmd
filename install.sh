#!/bin/bash

# Get the current working directory
SCRIPT_DIR="$(pwd)"

# Create the minicmd script at /usr/local/bin/minicmd
sudo tee /usr/local/bin/minicmd > /dev/null << EOF
#!/bin/bash
python3 $SCRIPT_DIR/minicmd.py "\$@"
EOF

# Make the script executable
sudo chmod +x /usr/local/bin/minicmd

echo "minicmd installed successfully!"