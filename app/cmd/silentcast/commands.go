package main

import (
	"flag"
	"fmt"
	"os"
)

// CommandFlags holds all command-line flags
type CommandFlags struct {
	// Core flags
	NoTray  bool
	Version bool
	Debug   bool

	// Config commands
	ValidateConfig bool
	ShowConfig     bool
	ShowConfigPath bool
	ShowFormat     string
	ShowPaths      bool
	
	// Version options
	VersionFormat string

	// Spell commands
	ListSpells bool
	ListFilter string

	// Debug commands
	TestHotkey   bool
	TestDuration int

	// Export/Import commands
	ExportConfig string
	ExportFormat string
	ImportConfig string
	
	// Service management
	ServiceInstall   bool
	ServiceUninstall bool
	ServiceStart     bool
	ServiceStop      bool
	ServiceStatus    bool
	
	// Update commands
	CheckUpdate      bool
	ForceUpdateCheck bool
	SelfUpdate       bool
	ForceSelfUpdate  bool
	UpdateStatus     bool
	
	// Future commands (ready for implementation)
	DryRun    bool
	Once      bool
	SpellName string
	TestSpell bool
	Benchmark bool
}

// ParseFlags parses command-line flags and returns the flags struct
func ParseFlags() *CommandFlags {
	flags := &CommandFlags{}

	// Define custom usage function using command registry
	flag.Usage = func() {
		// Create a temporary registry to generate help
		tmpFlags := &CommandFlags{}
		registry := NewCommandRegistry(tmpFlags)
		fmt.Fprint(os.Stderr, registry.GenerateHelp())
	}

	// Core flags
	flag.BoolVar(&flags.NoTray, "no-tray", false, "Disable system tray integration")
	flag.BoolVar(&flags.Version, "version", false, "Print version and exit")
	flag.BoolVar(&flags.Debug, "debug", false, "Enable debug logging")
	flag.StringVar(&flags.VersionFormat, "version-format", "human", "Version output format: human, json, compact")

	// Config commands
	flag.BoolVar(&flags.ValidateConfig, "validate-config", false, "Validate configuration and exit")
	flag.BoolVar(&flags.ShowConfig, "show-config", false, "Show merged configuration and exit")
	flag.BoolVar(&flags.ShowConfigPath, "show-config-path", false, "Show configuration file search paths")
	flag.StringVar(&flags.ShowFormat, "format", "human", "Output format for show-config: human, json, yaml")
	flag.BoolVar(&flags.ShowPaths, "show-paths", false, "Show configuration search paths with show-config")

	// Spell commands
	flag.BoolVar(&flags.ListSpells, "list-spells", false, "List all configured spells")
	flag.StringVar(&flags.ListFilter, "filter", "", "Filter spells by sequence, name, or description")

	// Debug commands
	flag.BoolVar(&flags.TestHotkey, "test-hotkey", false, "Test hotkey detection")
	flag.IntVar(&flags.TestDuration, "duration", 0, "Test duration in seconds (0 = until Ctrl+C)")

	// Single execution mode
	flag.BoolVar(&flags.Once, "once", false, "Execute a spell once and exit")
	flag.StringVar(&flags.SpellName, "spell", "", "Spell to execute in once mode")
	flag.BoolVar(&flags.TestSpell, "test-spell", false, "Test a spell with detailed debug information")
	flag.BoolVar(&flags.DryRun, "dry-run", false, "Show what would be executed without actually running it")
	flag.BoolVar(&flags.Benchmark, "benchmark", false, "Run performance benchmarks and exit")
	
	// Export/Import commands
	flag.StringVar(&flags.ExportConfig, "export-config", "", "Export configuration to file (or stdout if empty)")
	flag.StringVar(&flags.ExportFormat, "export-format", "yaml", "Export format: yaml, tar.gz")
	flag.StringVar(&flags.ImportConfig, "import-config", "", "Import configuration from file")
	
	// Service management (Windows only)
	flag.BoolVar(&flags.ServiceInstall, "service-install", false, "Install SilentCast as system service (Windows)")
	flag.BoolVar(&flags.ServiceUninstall, "service-uninstall", false, "Uninstall SilentCast service (Windows)")
	flag.BoolVar(&flags.ServiceStart, "service-start", false, "Start SilentCast service (Windows)")
	flag.BoolVar(&flags.ServiceStop, "service-stop", false, "Stop SilentCast service (Windows)")
	flag.BoolVar(&flags.ServiceStatus, "service-status", false, "Show SilentCast service status (Windows)")

	// Update commands
	flag.BoolVar(&flags.CheckUpdate, "check-update", false, "Check for available updates")
	flag.BoolVar(&flags.ForceUpdateCheck, "force-update-check", false, "Force update check (ignore cache)")
	flag.BoolVar(&flags.SelfUpdate, "self-update", false, "Update SilentCast to the latest version")
	flag.BoolVar(&flags.ForceSelfUpdate, "force", false, "Skip confirmation prompts for self-update")
	flag.BoolVar(&flags.UpdateStatus, "update-status", false, "Show current update status and available updates")

	flag.Parse()
	return flags
}
