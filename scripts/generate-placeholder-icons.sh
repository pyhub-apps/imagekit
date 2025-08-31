#!/bin/bash

# Generate placeholder PWA icons using ImageMagick
# These are temporary until proper icons are created with ChatGPT

ICON_DIR="web/static/icons"
mkdir -p $ICON_DIR

# Check if ImageMagick is installed
if ! command -v convert &> /dev/null; then
    echo "ImageMagick is not installed. Please install it first:"
    echo "  macOS: brew install imagemagick"
    echo "  Ubuntu: sudo apt-get install imagemagick"
    exit 1
fi

echo "Generating placeholder PWA icons..."

# Create a simple gradient placeholder icon (512x512)
convert -size 512x512 \
    -define gradient:angle=45 \
    gradient:'#667eea-#764ba2' \
    -gravity center \
    -fill white \
    -pointsize 200 \
    -annotate +0+0 '🖼' \
    $ICON_DIR/icon-512x512.png

# Generate all required sizes
sizes=(16 32 72 96 128 144 152 180 192 384)
for size in "${sizes[@]}"; do
    convert $ICON_DIR/icon-512x512.png \
        -resize ${size}x${size} \
        $ICON_DIR/icon-${size}x${size}.png
    echo "Generated icon-${size}x${size}.png"
done

echo "Placeholder icons generated successfully!"
echo "Note: These are temporary. Use ChatGPT to generate proper icons."