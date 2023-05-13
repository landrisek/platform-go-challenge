#!/bin/bash

# Detect the operating system
OS=$(uname | tr '[:upper:]' '[:lower:]')

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Installing now..."
    if [ "$OS" == "darwin" ]; then
        # Install Go on macOS using Homebrew
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
        brew install go
    elif [ "$OS" == "linux" ]; then
        # Install Go on Linux using apt-get
        sudo apt-get update
        sudo apt-get install -y golang
    else
        echo "Unsupported operating system: $OS"
        exit 1
    fi
fi
