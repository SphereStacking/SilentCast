//go:build linux

package config

import (
	"os"
	"path/filepath"
)

func init() {
	platformFactory = func() PlatformResolver {
		return &linuxResolver{}
	}
}

type linuxResolver struct{}

func (l *linuxResolver) GetPlatformConfigFile() string {
	return "spellbook.linux.yml"
}

func (l *linuxResolver) GetDefaultConfigPath() string {
	// First try XDG_CONFIG_HOME
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return filepath.Join(xdgConfig, "silentcast")
	}

	// Fallback to ~/.config
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homeDir, ".config", "silentcast")
}
