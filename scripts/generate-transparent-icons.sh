#!/bin/bash

# Generate PWA icons with transparent background
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

echo "Generating PWA icons with transparent background..."

# Create a rounded rectangle icon with gradient and transparent background (512x512)
convert -size 512x512 xc:transparent \
    \( -size 412x412 \
       -define gradient:angle=45 \
       gradient:'#667eea-#764ba2' \
       -gravity center \
    \) \
    -gravity center -geometry +0+0 -composite \
    \( -size 512x512 xc:transparent \
       -fill white \
       -draw "roundrectangle 50,50,462,462,60,60" \
    \) \
    -compose DstIn -composite \
    -gravity center \
    -fill white \
    -font Arial-Bold \
    -pointsize 200 \
    -annotate +0+0 '📷' \
    -background transparent \
    $ICON_DIR/icon-512x512.png

echo "Generated icon-512x512.png with transparent background"

# Generate all required sizes maintaining transparency
sizes=(16 32 72 96 128 144 152 180 192 384)
for size in "${sizes[@]}"; do
    convert $ICON_DIR/icon-512x512.png \
        -resize ${size}x${size} \
        -background transparent \
        -gravity center \
        -extent ${size}x${size} \
        $ICON_DIR/icon-${size}x${size}.png
    echo "Generated icon-${size}x${size}.png"
done

echo ""
echo "✅ All icons generated with transparent backgrounds!"
echo ""
echo "Note: These are placeholder icons. Use ChatGPT to generate proper icons."
echo "The generated icons will work well on both light and dark backgrounds."