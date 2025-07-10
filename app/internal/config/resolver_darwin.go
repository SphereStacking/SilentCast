//go:build darwin

package config

import (
	"os"
	"path/filepath"
)

func init() {
	platformFactory = func() PlatformResolver {
		return &darwinPlatform{}
	}
}

type darwinPlatform struct{}

func (d *darwinPlatform) GetPlatformConfigFile() string {
	return ConfigName + ".mac.yml"
}

func (d *darwinPlatform) GetDefaultConfigPath() string {
	// On macOS, use ~/Library/Application Support
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, "Library", "Application Support", AppName)
	}
	// Fallback to user config dir
	if configDir, err := os.UserConfigDir(); err == nil {
		return filepath.Join(configDir, AppName)
	}
	return "."
}