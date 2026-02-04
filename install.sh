#!/bin/bash

# Build the Go binary
echo "Building yact..."
go build -o y

if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi

# Install to /usr/local/bin
echo "Installing y to /usr/local/bin..."
sudo mv y /usr/local/bin/

if [ $? -eq 0 ]; then
    echo "Installation successful!"
    echo "You can now run 'y' from anywhere."
else
    echo "Installation failed. You may need to run this script with sudo."
    exit 1
fi