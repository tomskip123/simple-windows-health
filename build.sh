#!/bin/bash

# Exit on error
set -e

# Set Windows as target OS and architecture
export GOOS=windows
export GOARCH=amd64

# Output file name with .exe extension
OUTPUT="bin/wincleaner.exe"

# Create bin directory if it doesn't exist
mkdir -p bin

echo "Building for Windows (${GOARCH})..."
go build -o ${OUTPUT} ./cmd/wincleaner

echo "Build complete: ${OUTPUT}"
