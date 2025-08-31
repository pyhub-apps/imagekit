#!/bin/bash

# WebAssembly build script with version injection
# This script is used by Cloudflare Pages CI/CD and local development

set -e

# Get version from git tag or use default
# Try multiple methods to get version info
if [ -n "$VERSION" ]; then
    # Use environment variable if set
    echo "Using VERSION from environment: $VERSION"
elif git describe --tags --always --long 2>/dev/null > /dev/null; then
    # Full version with tags
    VERSION=$(git describe --tags --always --long 2>/dev/null)
    echo "Using git describe with tags: $VERSION"
elif git rev-parse --short HEAD 2>/dev/null > /dev/null; then
    # Fallback to short commit hash with default version prefix
    COMMIT=$(git rev-parse --short HEAD 2>/dev/null)
    VERSION="1.2535.17-${COMMIT}"
    echo "Using fallback version with commit: $VERSION"
else
    # Final fallback
    VERSION="1.2535.17-dev"
    echo "Using default version: $VERSION"
fi

# Remove leading 'v' if present to avoid duplication in HTML
VERSION_CLEAN=${VERSION#v}
echo "Building WebAssembly with version: $VERSION_CLEAN"

# Create template files with version placeholders
echo "Creating versioned files..."

# Create app.js from template
if [ -f "web/static/app.template.js" ]; then
    sed "s/{{VERSION}}/$VERSION_CLEAN/g" web/static/app.template.js > web/static/app.js
else
    # If template doesn't exist, update existing file
    sed -i.bak "s/const WASM_VERSION = '[^']*'/const WASM_VERSION = '$VERSION_CLEAN'/" web/static/app.js
    rm -f web/static/app.js.bak
fi

# Create index.html from template
if [ -f "web/index.template.html" ]; then
    sed "s/{{VERSION}}/$VERSION_CLEAN/g" web/index.template.html > web/index.html
else
    # If template doesn't exist, update existing file
    sed -i.bak "s/\?v=[0-9.]*/?v=$VERSION_CLEAN/g" web/index.html
    sed -i.bak "s/imagekit-version\">[^<]*/imagekit-version\">$VERSION_CLEAN/" web/index.html
    rm -f web/index.html.bak
fi

# Build WASM with version
echo "Building WASM..."
GOOS=js GOARCH=wasm go build \
    -ldflags="-X main.Version=$VERSION_CLEAN" \
    -o web/static/imagekit.wasm \
    cmd/wasm/main.go

echo "WASM build complete!"
ls -lh web/static/imagekit.wasm

echo "Build completed successfully with version $VERSION_CLEAN"