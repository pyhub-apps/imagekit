#!/bin/bash

# WebAssembly build script with version injection
# This script is used by Cloudflare Pages CI/CD

set -e

# Get version from git tag or use default
VERSION=${VERSION:-$(git describe --tags --always 2>/dev/null || echo "dev")}
echo "Building WebAssembly with version: $VERSION"

# Create template files with version placeholders
echo "Creating versioned files..."

# Create app.js from template
if [ -f "web/static/app.template.js" ]; then
    sed "s/{{VERSION}}/$VERSION/g" web/static/app.template.js > web/static/app.js
else
    # If template doesn't exist, update existing file
    sed -i.bak "s/const WASM_VERSION = '[^']*'/const WASM_VERSION = '$VERSION'/" web/static/app.js
    rm -f web/static/app.js.bak
fi

# Create index.html from template
if [ -f "web/index.template.html" ]; then
    sed "s/{{VERSION}}/$VERSION/g" web/index.template.html > web/index.html
else
    # If template doesn't exist, update existing file
    sed -i.bak "s/\?v=[0-9.]*/?v=$VERSION/g" web/index.html
    sed -i.bak "s/imagekit-version\">[^<]*/imagekit-version\">$VERSION/" web/index.html
    rm -f web/index.html.bak
fi

# Build WASM with version
echo "Building WASM..."
GOOS=js GOARCH=wasm go build \
    -ldflags="-X main.Version=$VERSION" \
    -o web/static/imagekit.wasm \
    cmd/wasm/main.go

echo "WASM build complete!"
ls -lh web/static/imagekit.wasm

echo "Build completed successfully with version $VERSION"