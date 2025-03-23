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
# Add trimpath to remove file system paths from binary
# Add -ldflags to set binary information that helps legitimize the executable
go build -trimpath -ldflags="-s -w -H=windowsgui -X 'main.version=1.0.1' -X 'main.company=TomPeacock'" -o ${OUTPUT} ./cmd/wincleaner

echo "Build complete: ${OUTPUT}"

# Suggest using Windows signtool to sign the executable
echo "Note: For maximum security compatibility, sign your executable with:"
echo "  signtool sign /tr http://timestamp.digicert.com /td sha256 /fd sha256 /a ${OUTPUT}"
