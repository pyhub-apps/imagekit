#!/bin/bash

# Script to update version numbers across the project
# Usage: ./scripts/update-version.sh <new-version>

if [ -z "$1" ]; then
    echo "Usage: $0 <new-version>"
    echo "Example: $0 1.2535.18"
    exit 1
fi

NEW_VERSION=$1
OLD_VERSION=$(grep 'const WASM_VERSION' web/static/app.js | sed -E "s/.*'([0-9.]+)'.*/\1/")

if [ -z "$OLD_VERSION" ]; then
    echo "Could not find current version"
    exit 1
fi

echo "Updating version from $OLD_VERSION to $NEW_VERSION"

# Update app.js
sed -i.bak "s/const WASM_VERSION = '$OLD_VERSION'/const WASM_VERSION = '$NEW_VERSION'/" web/static/app.js

# Update index.html - script tags
sed -i.bak "s/\?v=$OLD_VERSION/\?v=$NEW_VERSION/g" web/index.html

# Update index.html - version display
sed -i.bak "s/>$OLD_VERSION</>$NEW_VERSION</" web/index.html

# Update cmd/wasm/main.go
sed -i.bak "s/var Version = \"$OLD_VERSION\"/var Version = \"$NEW_VERSION\"/" cmd/wasm/main.go

# Update cmd/imagekit/main.go (if version is not "dev")
if grep -q "var Version = \"$OLD_VERSION\"" cmd/imagekit/main.go; then
    sed -i.bak "s/var Version = \"$OLD_VERSION\"/var Version = \"$NEW_VERSION\"/" cmd/imagekit/main.go
fi

# Remove backup files
find . -name "*.bak" -type f -delete

echo "Version updated to $NEW_VERSION"
echo "Don't forget to:"
echo "1. Rebuild WASM: make build-wasm"
echo "2. Commit changes"
echo "3. Create git tag: git tag v$NEW_VERSION"