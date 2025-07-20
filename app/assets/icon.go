package assets

import _ "embed"

// Icon is the application icon in PNG format
//
//go:embed icon.png
var Icon []byte

// IconWindows is the application icon in ICO format for Windows
//
//go:embed icon.ico
var IconWindows []byte

// IconMacOS is the application icon in ICNS format for macOS
//
//go:embed icon.icns
var IconMacOS []byte
