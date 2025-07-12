# Application Assets

This directory contains embedded assets for the application.

## Icons

- `icon.png` - Default icon in PNG format (used for system tray)
- `icon.ico` - Windows application icon
- `icon.icns` - macOS application icon

## Generating Icons

To generate icons from the SVG logo:

1. Install required tools:
   ```bash
   # Ubuntu/Debian
   sudo apt-get install inkscape imagemagick

   # macOS
   brew install librsvg imagemagick
   ```

2. Run the generation script:
   ```bash
   ../scripts/generate-icons.sh
   ```

This will convert the SVG logo into various PNG sizes and platform-specific formats.

## Embedding Icons in Executables

### Windows
Use a resource compiler like `windres` or `goversioninfo`:
```bash
go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo
goversioninfo -icon=assets/icon.ico
go build -ldflags="-H windowsgui"
```

### macOS
The icon is set via Info.plist when creating an .app bundle.

### Linux
The PNG icon is used directly by the system tray.