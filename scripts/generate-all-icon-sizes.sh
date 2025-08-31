#!/bin/bash

# Generate all PWA icon sizes from a 512x512 source image
# Usage: ./generate-all-icon-sizes.sh source-icon.png

if [ -z "$1" ]; then
    echo "Usage: $0 <source-icon-512x512.png>"
    echo "Source icon should be 512x512 pixels"
    exit 1
fi

SOURCE_ICON=$1
ICON_DIR="web/static/icons"

# Check if source file exists
if [ ! -f "$SOURCE_ICON" ]; then
    echo "Error: Source file '$SOURCE_ICON' not found"
    exit 1
fi

# Check if ImageMagick is installed
if ! command -v convert &> /dev/null; then
    echo "ImageMagick is not installed. Please install it first:"
    echo "  macOS: brew install imagemagick"
    echo "  Ubuntu: sudo apt-get install imagemagick"
    exit 1
fi

# Create icons directory
mkdir -p $ICON_DIR

echo "Generating PWA icons from $SOURCE_ICON..."

# Copy source as 512x512
cp "$SOURCE_ICON" "$ICON_DIR/icon-512x512.png"
echo "✓ icon-512x512.png"

# Generate all required sizes (maintaining transparency)
sizes=(16 32 72 96 128 144 152 180 192 384)
for size in "${sizes[@]}"; do
    convert "$SOURCE_ICON" \
        -resize ${size}x${size} \
        -background transparent \
        -gravity center \
        -extent ${size}x${size} \
        -unsharp 0.5x0.5+0.5+0.008 \
        "$ICON_DIR/icon-${size}x${size}.png"
    echo "✓ icon-${size}x${size}.png"
done

# Generate Apple Touch Icon with padding (180x180)
# Note: Apple Touch Icon typically needs a solid background
convert "$SOURCE_ICON" \
    -resize 160x160 \
    -gravity center \
    -background transparent \
    -extent 180x180 \
    "$ICON_DIR/apple-touch-icon.png"
echo "✓ apple-touch-icon.png (transparent background)"
echo "  Note: iOS may add its own background for Apple Touch Icons"

# Generate favicon.ico with multiple sizes
convert "$SOURCE_ICON" \
    -resize 16x16 "$ICON_DIR/icon-16.png"
convert "$SOURCE_ICON" \
    -resize 32x32 "$ICON_DIR/icon-32.png"
convert "$SOURCE_ICON" \
    -resize 48x48 "$ICON_DIR/icon-48.png"
convert "$ICON_DIR/icon-16.png" "$ICON_DIR/icon-32.png" "$ICON_DIR/icon-48.png" \
    "$ICON_DIR/favicon.ico"
rm "$ICON_DIR/icon-16.png" "$ICON_DIR/icon-32.png" "$ICON_DIR/icon-48.png"
echo "✓ favicon.ico"

echo ""
echo "✅ All PWA icons generated successfully in $ICON_DIR/"
echo ""
echo "Icons created:"
ls -la $ICON_DIR/*.png $ICON_DIR/*.ico | awk '{print "  - " $NF}'