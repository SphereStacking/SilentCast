package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/SphereStacking/silentcast/internal/config"
)

// ShowConfigPathCommand shows configuration search paths
type ShowConfigPathCommand struct {
	getConfigPath        func() string
	getConfigSearchPaths func() []string
}

// NewShowConfigPathCommand creates a new show config path command
func NewShowConfigPathCommand(getConfigPath func() string, getConfigSearchPaths func() []string) Command {
	return &ShowConfigPathCommand{
		getConfigPath:        getConfigPath,
		getConfigSearchPaths: getConfigSearchPaths,
	}
}

// Name returns the command name
func (c *ShowConfigPathCommand) Name() string {
	return "Show Config Path"
}

// Description returns the command description
func (c *ShowConfigPathCommand) Description() string {
	return "Show configuration file search paths"
}

// FlagName returns the flag name
func (c *ShowConfigPathCommand) FlagName() string {
	return "show-config-path"
}

// IsActive checks if the command should run
func (c *ShowConfigPathCommand) IsActive(flags interface{}) bool {
	f, ok := flags.(*Flags)
	if !ok {
		return false
	}
	return f.ShowConfigPath
}

// Execute runs the command
func (c *ShowConfigPathCommand) Execute(flags interface{}) error {
	fmt.Println("📍 SilentCast Configuration Search Paths")
	fmt.Println("=========================================")
	fmt.Println()
	fmt.Println("Configuration files are searched in the following order:")
	fmt.Println()

	paths := c.getConfigSearchPaths()
	currentPath := c.getConfigPath()

	for i, path := range paths {
		fmt.Printf("%d. %s", i+1, path)

		// Check if this is the current path
		if path == currentPath {
			fmt.Print(" ⭐ (active)")
		}

		// Check what files exist
		var found []string

		// Base config
		baseConfig := filepath.Join(path, config.ConfigName+".yml")
		if _, err := os.Stat(baseConfig); err == nil {
			found = append(found, config.ConfigName+".yml")
		}

		// OS-specific config
		osConfig := ""
		switch runtime.GOOS {
		case "darwin":
			osConfig = config.ConfigName + ".mac.yml"
		case "windows":
			osConfig = config.ConfigName + ".windows.yml"
		case "linux":
			osConfig = config.ConfigName + ".linux.yml"
		}

		if osConfig != "" {
			osConfigPath := filepath.Join(path, osConfig)
			if _, err := os.Stat(osConfigPath); err == nil {
				found = append(found, osConfig)
			}
		}

		// Show found files
		if len(found) > 0 {
			fmt.Printf("\n   ✅ Found: %s", strings.Join(found, ", "))
		} else {
			fmt.Print("\n   ❌ No config files found")
		}

		// Show path details based on actual position
		if os.Getenv("SILENTCAST_CONFIG") != "" {
			if i == 0 {
				fmt.Print(" (from SILENTCAST_CONFIG env)")
			} else if i == 1 {
				fmt.Print(" (current directory)")
			} else if i == 2 {
				fmt.Print(" (user config directory)")
			} else if i == 3 && runtime.GOOS != "windows" {
				fmt.Print(" (system config directory)")
			}
		} else {
			if i == 0 {
				fmt.Print(" (current directory)")
			} else if i == 1 {
				fmt.Print(" (user config directory)")
			} else if i == 2 && runtime.GOOS != "windows" {
				fmt.Print(" (system config directory)")
			}
		}

		fmt.Println()
		fmt.Println()
	}

	fmt.Println("📝 Configuration Loading Order:")
	fmt.Println("   1. Base config (spellbook.yml)")
	fmt.Printf("   2. OS-specific override (spellbook.%s.yml)\n", runtime.GOOS)
	fmt.Println("   3. Settings are merged with OS-specific values taking precedence")

	return nil
}

// Group returns the command group
func (c *ShowConfigPathCommand) Group() string {
	return "config"
}

// HasOptions returns if this command has additional options
func (c *ShowConfigPathCommand) HasOptions() bool {
	return false
}
