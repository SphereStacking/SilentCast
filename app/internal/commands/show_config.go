package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/SphereStacking/silentcast/internal/config"
)

// ShowConfigCommand shows the merged configuration
type ShowConfigCommand struct {
	getConfigPath        func() string
	getConfigSearchPaths func() []string
}

// NewShowConfigCommand creates a new show config command
func NewShowConfigCommand(getConfigPath func() string, getConfigSearchPaths func() []string) Command {
	return &ShowConfigCommand{
		getConfigPath:        getConfigPath,
		getConfigSearchPaths: getConfigSearchPaths,
	}
}

// Name returns the command name
func (c *ShowConfigCommand) Name() string {
	return "Show Config"
}

// Description returns the command description
func (c *ShowConfigCommand) Description() string {
	return "Show merged configuration and exit"
}

// FlagName returns the flag name
func (c *ShowConfigCommand) FlagName() string {
	return "show-config"
}

// IsActive checks if the command should run
func (c *ShowConfigCommand) IsActive(flags interface{}) bool {
	f, ok := flags.(*Flags)
	if !ok {
		return false
	}
	return f.ShowConfig
}

// Execute runs the command
func (c *ShowConfigCommand) Execute(flags interface{}) error {
	f, ok := flags.(*Flags)
	if !ok {
		return fmt.Errorf("invalid flags type")
	}

	configPath := c.getConfigPath()

	// Only show header for human format
	if f.ShowFormat == "human" || f.ShowFormat == "" {
		fmt.Printf("🔍 Loading configuration from: %s\n", configPath)
	}

	// Show config search paths if requested
	if f.ShowPaths && (f.ShowFormat == "human" || f.ShowFormat == "") {
		fmt.Println("\n📍 Configuration Search Paths:")
		paths := c.getConfigSearchPaths()
		for i, path := range paths {
			fmt.Printf("   %d. %s", i+1, path)

			// Check if files exist
			var found []string
			baseConfig := filepath.Join(path, config.ConfigName+".yml")
			if _, err := os.Stat(baseConfig); err == nil {
				found = append(found, config.ConfigName+".yml")
			}

			// Check OS-specific config
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

			if len(found) > 0 {
				fmt.Printf(" ✅ [%s]", strings.Join(found, ", "))
			}
			fmt.Println()
		}
		fmt.Println()
	}

	// Create loader
	loader := config.NewLoader(configPath)

	// Load configuration
	cfg, err := loader.Load()
	if err != nil {
		fmt.Printf("❌ Configuration loading failed: %v\n", err)
		return err
	}

	// Only show header for human format
	if f.ShowFormat == "human" || f.ShowFormat == "" {
		fmt.Println("\n📋 Merged Configuration:")
		fmt.Println(strings.Repeat("-", 50))
	}

	// Display based on format
	switch f.ShowFormat {
	case "json":
		return c.showConfigJSON(cfg)
	case "yaml":
		return c.showConfigYAML(cfg)
	default:
		return c.showConfigHuman(cfg)
	}
}

// Group returns the command group
func (c *ShowConfigCommand) Group() string {
	return "config"
}

// HasOptions returns if this command has additional options
func (c *ShowConfigCommand) HasOptions() bool {
	return true
}

// showConfigJSON displays configuration in JSON format
func (c *ShowConfigCommand) showConfigJSON(cfg *config.Config) error {
	// Create a clean structure for JSON output
	output := map[string]interface{}{
		"daemon": map[string]interface{}{
			"auto_start":   cfg.Daemon.AutoStart,
			"log_level":    cfg.Daemon.LogLevel,
			"config_watch": cfg.Daemon.ConfigWatch,
		},
		"hotkeys": map[string]interface{}{
			"prefix":           cfg.Hotkeys.Prefix,
			"timeout":          cfg.Hotkeys.Timeout.ToDuration().Milliseconds(),
			"sequence_timeout": cfg.Hotkeys.SequenceTimeout.ToDuration().Milliseconds(),
		},
		"spells":   cfg.Shortcuts,
		"grimoire": cfg.Actions,
		"logger": map[string]interface{}{
			"level":       cfg.Logger.Level,
			"file":        cfg.Logger.File,
			"max_size":    cfg.Logger.MaxSize,
			"max_backups": cfg.Logger.MaxBackups,
			"max_age":     cfg.Logger.MaxAge,
			"compress":    cfg.Logger.Compress,
		},
		"updater": map[string]interface{}{
			"enabled":        cfg.Updater.Enabled,
			"check_interval": cfg.Updater.CheckInterval,
			"auto_install":   cfg.Updater.AutoInstall,
			"prerelease":     cfg.Updater.Prerelease,
		},
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}

// showConfigYAML displays configuration in YAML format
func (c *ShowConfigCommand) showConfigYAML(cfg *config.Config) error {
	encoder := yaml.NewEncoder(os.Stdout)
	encoder.SetIndent(2)
	return encoder.Encode(cfg)
}

// showConfigHuman displays configuration in human-readable format
func (c *ShowConfigCommand) showConfigHuman(cfg *config.Config) error {
	// Daemon settings
	fmt.Println("🔧 Daemon Settings:")
	fmt.Printf("   Auto Start: %v\n", cfg.Daemon.AutoStart)
	fmt.Printf("   Log Level: %s\n", cfg.Daemon.LogLevel)
	fmt.Printf("   Config Watch: %v\n", cfg.Daemon.ConfigWatch)

	// Hotkey settings
	fmt.Println("\n⌨️  Hotkey Settings:")
	fmt.Printf("   Prefix: %s\n", cfg.Hotkeys.Prefix)
	fmt.Printf("   Timeout: %d ms\n", cfg.Hotkeys.Timeout.ToDuration().Milliseconds())
	fmt.Printf("   Sequence Timeout: %d ms\n", cfg.Hotkeys.SequenceTimeout.ToDuration().Milliseconds())

	// Spells (sorted)
	fmt.Printf("\n✨ Spells (%d defined):\n", len(cfg.Shortcuts))
	var spells []string
	for sequence := range cfg.Shortcuts {
		spells = append(spells, sequence)
	}
	sort.Strings(spells)

	for _, sequence := range spells {
		spellName := cfg.Shortcuts[sequence]
		fmt.Printf("   %s → %s", sequence, spellName)

		// Show action description if available
		if action, exists := cfg.Actions[spellName]; exists && action.Description != "" {
			fmt.Printf(" (%s)", action.Description)
		}
		fmt.Println()
	}

	// Grimoire actions (sorted)
	fmt.Printf("\n📚 Grimoire Actions (%d defined):\n", len(cfg.Actions))
	var actionNames []string
	for name := range cfg.Actions {
		actionNames = append(actionNames, name)
	}
	sort.Strings(actionNames)

	for _, name := range actionNames {
		action := cfg.Actions[name]
		fmt.Printf("\n   %s:\n", name)
		fmt.Printf("      Type: %s\n", action.Type)
		fmt.Printf("      Command: %s\n", action.Command)

		if len(action.Args) > 0 {
			fmt.Printf("      Args: %v\n", action.Args)
		}
		if action.WorkingDir != "" {
			fmt.Printf("      Working Dir: %s\n", action.WorkingDir)
		}
		if len(action.Env) > 0 {
			fmt.Printf("      Environment:\n")
			for k, v := range action.Env {
				fmt.Printf("         %s: %s\n", k, v)
			}
		}

		// Show options
		var options []string
		if action.ShowOutput {
			options = append(options, "show_output")
		}
		if action.KeepOpen {
			options = append(options, "keep_open")
		}
		if action.Admin {
			options = append(options, "admin")
		}
		if action.Terminal {
			options = append(options, "terminal")
		}
		if action.Timeout > 0 {
			options = append(options, fmt.Sprintf("timeout=%ds", action.Timeout))
		}
		if action.Shell != "" {
			options = append(options, fmt.Sprintf("shell=%s", action.Shell))
		}

		if len(options) > 0 {
			fmt.Printf("      Options: %s\n", strings.Join(options, ", "))
		}

		if action.Description != "" {
			fmt.Printf("      Description: %s\n", action.Description)
		}
	}

	// Logger settings
	fmt.Println("\n📝 Logger Settings:")
	fmt.Printf("   Level: %s\n", cfg.Logger.Level)
	if cfg.Logger.File != "" {
		fmt.Printf("   File: %s\n", cfg.Logger.File)
	}
	fmt.Printf("   Max Size: %d MB\n", cfg.Logger.MaxSize)
	fmt.Printf("   Max Backups: %d\n", cfg.Logger.MaxBackups)
	fmt.Printf("   Max Age: %d days\n", cfg.Logger.MaxAge)
	fmt.Printf("   Compress: %v\n", cfg.Logger.Compress)

	// Updater settings
	fmt.Println("\n🔄 Updater Settings:")
	fmt.Printf("   Enabled: %v\n", cfg.Updater.Enabled)
	if cfg.Updater.Enabled {
		fmt.Printf("   Check Interval: %s\n", cfg.Updater.CheckInterval)
		fmt.Printf("   Auto Install: %v\n", cfg.Updater.AutoInstall)
		fmt.Printf("   Include Prereleases: %v\n", cfg.Updater.Prerelease)
	}

	fmt.Println("\n✅ Configuration loaded successfully!")
	return nil
}
