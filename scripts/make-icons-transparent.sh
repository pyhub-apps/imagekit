#!/bin/bash

# Make existing PNG icons transparent by removing white/solid backgrounds
# Usage: ./make-icons-transparent.sh [color-to-remove]

ICON_DIR="web/static/icons"
REMOVE_COLOR=${1:-white}  # Default to removing white background

# Check if ImageMagick is installed
if ! command -v convert &> /dev/null; then
    echo "ImageMagick is not installed. Please install it first:"
    echo "  macOS: brew install imagemagick"
    echo "  Ubuntu: sudo apt-get install imagemagick"
    exit 1
fi

# Check if icons directory exists
if [ ! -d "$ICON_DIR" ]; then
    echo "Icons directory not found: $ICON_DIR"
    exit 1
fi

echo "Making existing icons transparent by removing $REMOVE_COLOR background..."
echo ""

# Process each PNG file
for icon in $ICON_DIR/icon-*.png; do
    if [ -f "$icon" ]; then
        filename=$(basename "$icon")
        
        # Create backup
        cp "$icon" "${icon}.backup"
        
        # Remove background and make transparent
        convert "$icon" \
            -fuzz 10% \
            -transparent "$REMOVE_COLOR" \
            "$icon"
        
        echo "✓ Processed $filename"
    fi
done

echo ""
echo "✅ All icons processed!"
echo "Backup files created with .backup extension"
echo ""
echo "Tips:"
echo "- If the result is not satisfactory, try with different colors:"
echo "  ./make-icons-transparent.sh '#ffffff'  # Remove white"
echo "  ./make-icons-transparent.sh '#f0f0f0'  # Remove light gray"
echo "- Restore from backup: mv icon-512x512.png.backup icon-512x512.png"