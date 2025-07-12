#!/bin/bash

# Generate icons from SVG for all platforms
# Requires: inkscape or rsvg-convert

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
SVG_FILE="$ROOT_DIR/app/assets/icons/logo.svg"
ICONS_DIR="$ROOT_DIR/app/assets/icons"

if [ ! -f "$SVG_FILE" ]; then
    echo "Error: $SVG_FILE not found"
    exit 1
fi

echo "Generating icons from $SVG_FILE..."

# Check for required tools
if command -v inkscape &> /dev/null; then
    CONVERTER="inkscape"
elif command -v rsvg-convert &> /dev/null; then
    CONVERTER="rsvg-convert"
else
    echo "Error: Neither inkscape nor rsvg-convert found. Please install one of them."
    echo "  Ubuntu/Debian: sudo apt-get install inkscape"
    echo "  macOS: brew install librsvg"
    exit 1
fi

# Generate PNG files for various sizes
SIZES=(16 32 48 64 128 256 512)
for size in "${SIZES[@]}"; do
    echo "Generating ${size}x${size} PNG..."
    if [ "$CONVERTER" = "inkscape" ]; then
        inkscape -w $size -h $size "$SVG_FILE" -o "$ICONS_DIR/icon_${size}.png"
    else
        rsvg-convert -w $size -h $size "$SVG_FILE" -o "$ICONS_DIR/icon_${size}.png"
    fi
done

# Generate Windows ICO (requires ImageMagick)
if command -v convert &> /dev/null; then
    echo "Generating Windows ICO..."
    convert "$ICONS_DIR/icon_16.png" "$ICONS_DIR/icon_32.png" "$ICONS_DIR/icon_48.png" "$ICONS_DIR/icon_256.png" "$ICONS_DIR/icon.ico"
else
    echo "Warning: ImageMagick not found. Skipping ICO generation."
    echo "  Install with: sudo apt-get install imagemagick"
fi

# Generate macOS ICNS (requires png2icns or iconutil)
if command -v png2icns &> /dev/null; then
    echo "Generating macOS ICNS..."
    png2icns "$ICONS_DIR/icon.icns" "$ICONS_DIR/icon_16.png" "$ICONS_DIR/icon_32.png" "$ICONS_DIR/icon_128.png" "$ICONS_DIR/icon_256.png" "$ICONS_DIR/icon_512.png"
elif command -v iconutil &> /dev/null; then
    echo "Generating macOS ICNS using iconutil..."
    # Create iconset directory
    ICONSET_DIR="$ICONS_DIR/icon.iconset"
    mkdir -p "$ICONSET_DIR"
    
    # Copy and rename files for iconutil
    cp "$ICONS_DIR/icon_16.png" "$ICONSET_DIR/icon_16x16.png"
    cp "$ICONS_DIR/icon_32.png" "$ICONSET_DIR/icon_16x16@2x.png"
    cp "$ICONS_DIR/icon_32.png" "$ICONSET_DIR/icon_32x32.png"
    cp "$ICONS_DIR/icon_64.png" "$ICONSET_DIR/icon_32x32@2x.png"
    cp "$ICONS_DIR/icon_128.png" "$ICONSET_DIR/icon_128x128.png"
    cp "$ICONS_DIR/icon_256.png" "$ICONSET_DIR/icon_128x128@2x.png"
    cp "$ICONS_DIR/icon_256.png" "$ICONSET_DIR/icon_256x256.png"
    cp "$ICONS_DIR/icon_512.png" "$ICONSET_DIR/icon_256x256@2x.png"
    cp "$ICONS_DIR/icon_512.png" "$ICONSET_DIR/icon_512x512.png"
    
    # Generate ICNS
    iconutil -c icns "$ICONSET_DIR" -o "$ICONS_DIR/icon.icns"
    
    # Clean up
    rm -rf "$ICONSET_DIR"
else
    echo "Warning: Neither png2icns nor iconutil found. Skipping ICNS generation."
    echo "  macOS: iconutil is included with Xcode"
    echo "  Linux: sudo apt-get install icnsutils"
fi

echo "Icon generation complete!"
echo "Generated files in: $ICONS_DIR"
ls -la "$ICONS_DIR"