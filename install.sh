#!/bin/bash

# Build the Go binary
echo "Building minicmd..."
go build -o minicmd

if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi

# Install to /usr/local/bin
echo "Installing minicmd to /usr/local/bin..."
sudo mv minicmd /usr/local/bin/

if [ $? -eq 0 ]; then
    echo "Installation successful!"
    echo "You can now run 'minicmd' from anywhere."
else
    echo "Installation failed. You may need to run this script with sudo."
    exit 1
fi
