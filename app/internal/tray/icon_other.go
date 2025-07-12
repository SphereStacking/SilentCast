//go:build !windows && !darwin && !notray

package tray

import "github.com/SphereStacking/silentcast/assets"

// getIcon returns the tray icon data for other platforms
func getIcon() []byte {
	// Use PNG format for Linux and other platforms
	return assets.Icon
}