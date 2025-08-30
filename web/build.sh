#!/bin/bash

# Build script for WebAssembly deployment with versioning

echo "Building versioned WebAssembly..."

# Generate version
VERSION=$(./scripts/headver.sh)
if [ $? -ne 0 ]; then
    echo "Failed to generate version"
    exit 1
fi

echo "Building version: $VERSION"

# Build WebAssembly with version
WASM_FILE="web/static/imagekit-${VERSION}.wasm"
GOOS=js GOARCH=wasm go build -o "$WASM_FILE" cmd/wasm/main.go

if [ $? -ne 0 ]; then
    echo "WebAssembly build failed"
    exit 1
fi

# Generate version info file
cat > web/version.json << EOF
{
  "version": "$VERSION",
  "wasmFile": "static/imagekit-${VERSION}.wasm",
  "buildTime": "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
}
EOF

echo "WebAssembly build complete"
ls -lh "$WASM_FILE" web/version.json

echo "Build complete! Ready for deployment."