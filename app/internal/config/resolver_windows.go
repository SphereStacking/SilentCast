//go:build windows

package config

import (
	"os"
	"path/filepath"
)

func init() {
	platformFactory = func() PlatformResolver {
		return &windowsPlatform{}
	}
}

type windowsPlatform struct{}

func (w *windowsPlatform) GetPlatformConfigFile() string {
	return ConfigName + ".windows.yml"
}

func (w *windowsPlatform) GetDefaultConfigPath() string {
	// On Windows, use %APPDATA%
	if appData := os.Getenv("APPDATA"); appData != "" {
		return filepath.Join(appData, AppName)
	}
	// Fallback to user config dir
	if configDir, err := os.UserConfigDir(); err == nil {
		return filepath.Join(configDir, AppName)
	}
	return "."
}