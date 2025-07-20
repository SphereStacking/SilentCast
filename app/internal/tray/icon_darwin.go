//go:build darwin && !notray

package tray

import "github.com/SphereStacking/silentcast/assets"

// getIcon returns the tray icon data for macOS
func getIcon() []byte {
	// Use ICNS format for macOS if available
	if len(assets.IconMacOS) > 0 {
		return assets.IconMacOS
	}
	// Fallback to PNG
	return assets.Icon
}
