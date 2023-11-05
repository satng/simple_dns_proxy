#!/bin/bash

# Set the main package name
MAIN_PACKAGE="simple_dns_proxy"

# Set the output binary name
OUTPUT_BINARY="dns_proxy"

# Build for Windows
echo "Building for Windows..."
GOOS=windows GOARCH=amd64 go build -o "${OUTPUT_BINARY}_windows_amd64.exe" "${MAIN_PACKAGE}"

# Build for macOS
echo "Building for macOS..."
GOOS=darwin GOARCH=amd64 go build -o "${OUTPUT_BINARY}_darwin_amd64" "${MAIN_PACKAGE}"

# Build for Linux
echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o "${OUTPUT_BINARY}_linux_amd64" "${MAIN_PACKAGE}"

# Package the binaries (optional)
echo "Packaging the binaries..."
mkdir -p dist
mv "${OUTPUT_BINARY}_windows_amd64.exe" dist/
mv "${OUTPUT_BINARY}_darwin_amd64" dist/
mv "${OUTPUT_BINARY}_linux_amd64" dist/

echo "Build and packaging completed. Binaries saved in the 'dist' directory."