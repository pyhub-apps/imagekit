#!/bin/bash

# Cloudflare Pages Build Script
# This script is called by Cloudflare Pages during build

set -e

echo "Starting WebAssembly build..."

# Install Go if not available (Cloudflare Pages has it pre-installed)
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Please ensure Go is available in the build environment."
    exit 1
fi

echo "Go version: $(go version)"

# Build WebAssembly
echo "Building WebAssembly..."
GOOS=js GOARCH=wasm go build -o web/static/imagekit.wasm cmd/wasm/main.go

if [ -f "web/static/imagekit.wasm" ]; then
    echo "✅ WebAssembly build successful!"
    echo "File size: $(ls -lh web/static/imagekit.wasm | awk '{print $5}')"
else
    echo "❌ WebAssembly build failed!"
    exit 1
fi

echo "Build complete!"