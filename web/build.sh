#!/bin/bash

# Build script for Cloudflare Pages deployment

echo "Building WebAssembly..."

# Build WebAssembly
GOOS=js GOARCH=wasm go build -o web/static/imagekit.wasm cmd/wasm/main.go

if [ $? -ne 0 ]; then
    echo "WebAssembly build failed"
    exit 1
fi

echo "WebAssembly build complete"
ls -lh web/static/imagekit.wasm

echo "Build complete! Ready for deployment."