#!/bin/bash

set -e

export CGO_ENABLED=0
export GO111MODULE="on"
export GOPATH="$HOME/go"

# build asset binary
GOOS=linux GOARCH=amd64 go build -o artifacts/asset/asset cmd/asset/*.go