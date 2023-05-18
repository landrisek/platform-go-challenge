#!/bin/bash

# Set the root directory
ROOT_DIR=$(dirname "$0")

# Load environment variables from .env file
ENV_FILE="artifacts/.asset.env"
if [ -f "$ENV_FILE" ]; then
  while IFS= read -r line; do
    # comments and white spaces lines
    if [[ $line != \#* ]] && [[ -n "$line" ]] && [[ ! "$line" =~ ^[[:space:]]*$ ]]; then
      export "$line"
    fi
  done < "$ENV_FILE"
fi

# Export other required environment variables
export GO111MODULE="on"
export GOPATH="$HOME/go"

# Run the Go application
go run ./cmd/asset/main.go