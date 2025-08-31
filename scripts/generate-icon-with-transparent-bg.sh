#!/bin/bash

# Generate PWA icons with transparent background (keeping white elements)
# The icon has a purple gradient rounded square with white content

ICON_DIR="web/static/icons"
mkdir -p $ICON_DIR

# Check if ImageMagick is installed
if ! command -v convert &> /dev/null && ! command -v magick &> /dev/null; then
    echo "ImageMagick is not installed. Please install it first:"
    echo "  macOS: brew install imagemagick"
    echo "  Ubuntu: sudo apt-get install imagemagick"
    exit 1
fi

# Use magick if available (ImageMagick 7), otherwise use convert
if command -v magick &> /dev/null; then
    CMD="magick"
else
    CMD="convert"
fi

echo "Generating PWA icons with transparent background..."

# Create the main 512x512 icon
# 1. Start with transparent canvas
# 2. Create gradient background
# 3. Apply to rounded rectangle
# 4. Add white icon/text on top
$CMD -size 512x512 xc:transparent \
    \( -size 400x400 \
       gradient:'#667eea-#764ba2' \
       -resize 400x400 \
    \) \
    -gravity center -geometry +0+0 -compose over -composite \
    \( -size 512x512 xc:transparent \
       -fill black \
       -draw "roundrectangle 56,56,456,456,40,40" \
    \) \
    -gravity center -compose DstIn -composite \
    -gravity center \
    -fill white \
    -font Helvetica-Bold \
    -pointsize 180 \
    -annotate +0-10 '🖼' \
    -gravity center \
    -fill white \
    -font Helvetica \
    -pointsize 40 \
    -annotate +0+100 'ImageKit' \
    $ICON_DIR/icon-512x512.png

echo "Generated icon-512x512.png"

# Alternative simpler design (just icon, no text)
$CMD -size 512x512 xc:transparent \
    \( -size 400x400 \
       gradient:'#667eea-#764ba2' \
       -resize 400x400 \
    \) \
    -gravity center -geometry +0+0 -compose over -composite \
    \( -size 512x512 xc:transparent \
       -fill black \
       -draw "roundrectangle 56,56,456,456,40,40" \
    \) \
    -gravity center -compose DstIn -composite \
    -gravity center \
    -fill white \
    -stroke white \
    -strokewidth 8 \
    -draw "rectangle 156,176,356,336" \
    -fill none \
    -stroke white \
    -strokewidth 6 \
    -draw "polyline 226,256 276,306 376,206" \
    $ICON_DIR/icon-512x512-alt.png

echo "Generated icon-512x512-alt.png (alternative design)"

# Generate all required sizes from the main icon
sizes=(16 32 72 96 128 144 152 180 192 384)
for size in "${sizes[@]}"; do
    $CMD $ICON_DIR/icon-512x512.png \
        -resize ${size}x${size} \
        -background transparent \
        -gravity center \
        -extent ${size}x${size} \
        $ICON_DIR/icon-${size}x${size}.png
    echo "Generated icon-${size}x${size}.png"
done

echo ""
echo "✅ All icons generated with transparent background!"
echo "   Purple gradient background with white content preserved."
echo ""
echo "Two designs available:"
echo "  - icon-512x512.png: With emoji and text"
echo "  - icon-512x512-alt.png: Simple geometric design"
echo ""
echo "Choose the one you prefer and use:"
echo "  ./scripts/generate-all-icon-sizes.sh web/static/icons/icon-512x512.png"