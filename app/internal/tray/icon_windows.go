//go:build windows && !notray

package tray

import "github.com/SphereStacking/silentcast/assets"

// getIcon returns the tray icon data for Windows
func getIcon() []byte {
	// Use ICO format for Windows
	if len(assets.IconWindows) > 0 {
		return assets.IconWindows
	}
	// Fallback to PNG if ICO is not available
	return assets.Icon
}