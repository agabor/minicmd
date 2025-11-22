pip install requests
pip install anthropic
if [ -n "$ZSH_VERSION" ]; then
    echo "expoPATH=\"$(pwd):\$PATH\"" >> ~/.zshrc
elif [ -n "$BASH_VERSION" ]; then
    echo "export PATH=\"$(pwd):\$PATH\"" >> ~/.bashrc
fi
