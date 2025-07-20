package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/SphereStacking/silentcast/cmd/silentcast/commands"
	"github.com/SphereStacking/silentcast/internal/version"
)

// CommandRegistry manages command pattern integration
type CommandRegistry struct {
	registry *commands.Registry
	flags    *CommandFlags
}

// NewCommandRegistry creates a new command registry
func NewCommandRegistry(flags *CommandFlags) *CommandRegistry {
	return NewCommandRegistryWithService(flags, nil)
}

// NewCommandRegistryWithService creates a new command registry with service support
func NewCommandRegistryWithService(flags *CommandFlags, onRun func() error) *CommandRegistry {
	registry := commands.NewRegistry()

	// Register all commands
	registry.RegisterAll(
		commands.NewVersionCommand(version.GetVersionString()),
		commands.NewValidateConfigCommand(getConfigPath),
		commands.NewShowConfigCommand(getConfigPath, getConfigSearchPaths),
		commands.NewShowConfigPathCommand(getConfigPath, getConfigSearchPaths),
		commands.NewListSpellsCommand(getConfigPath),
		commands.NewTestHotkeyCommand(getConfigPath),
		commands.NewExportConfigCommand(getConfigPath, getConfigSearchPaths),
		commands.NewImportConfigCommand(getConfigPath, getConfigSearchPaths),
		commands.NewCheckUpdateCommand(getConfigPath),
		commands.NewSelfUpdateCommand(getConfigPath),
		commands.NewUpdateStatusCommand(getConfigPath),
		// Service commands (Windows only)
		commands.NewServiceInstallCommand(onRun),
		commands.NewServiceUninstallCommand(),
		commands.NewServiceStartCommand(),
		commands.NewServiceStopCommand(),
		commands.NewServiceStatusCommand(),
	)

	return &CommandRegistry{
		registry: registry,
		flags:    flags,
	}
}

// ExecuteCommands runs commands using the command pattern
func (cr *CommandRegistry) ExecuteCommands() bool {
	// No conversion needed - CommandFlags is now an alias to commands.Flags
	executed, err := cr.registry.Execute(cr.flags)

	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ %s\n", err)
		os.Exit(1)
	}

	return executed
}

// GenerateHelp creates comprehensive help text with usage examples
func (cr *CommandRegistry) GenerateHelp() string {
	var sb strings.Builder

	// Header and introduction
	sb.WriteString("🪄 SilentCast - Silent hotkey-driven task runner\n\n")
	sb.WriteString("SilentCast executes tasks via keyboard shortcuts. Press your prefix key\n")
	sb.WriteString("(default: Alt+Space) followed by configured spells to trigger actions.\n\n")
	
	sb.WriteString(fmt.Sprintf("Usage: %s [options]\n\n", os.Args[0]))

	// Quick start section
	sb.WriteString("📚 Quick Start:\n")
	sb.WriteString(fmt.Sprintf("  %s                    # Run with system tray\n", os.Args[0]))
	sb.WriteString(fmt.Sprintf("  %s --no-tray          # Run without tray (terminal mode)\n", os.Args[0]))
	sb.WriteString(fmt.Sprintf("  %s --list-spells      # See all available spells\n", os.Args[0]))
	sb.WriteString(fmt.Sprintf("  %s --test-hotkey      # Test hotkey detection\n", os.Args[0]))
	sb.WriteString("\n")

	// Core options
	sb.WriteString("🔧 Core Options:\n")
	sb.WriteString("  -no-tray              Disable system tray integration\n")
	sb.WriteString("  -debug                Enable debug logging for troubleshooting\n")
	sb.WriteString("  -help                 Show this comprehensive help\n")
	sb.WriteString("\n")

	// Get groups from registry and display organized commands
	groups := cr.registry.GetGroups()
	for _, group := range groups {
		title := strings.Title(group.Name)
		sb.WriteString(fmt.Sprintf("📋 %s Commands - %s:\n", title, group.Description))

		for _, cmd := range group.Commands {
			sb.WriteString(fmt.Sprintf("  -%s", cmd.FlagName()))
			if cmd.HasOptions() {
				sb.WriteString(" [options]")
			}
			sb.WriteString(fmt.Sprintf("\n        %s\n", cmd.Description()))
		}
		sb.WriteString("\n")
	}

	// Single execution mode section
	sb.WriteString("⚡ Single Execution Mode (Execute spells directly):\n")
	sb.WriteString("  -once                 Execute a spell once and exit (requires -spell)\n")
	sb.WriteString("  -spell=<spell>        Spell to execute (e.g., 'e', 'g,s', 'vs,code')\n")
	sb.WriteString("  -test-spell           Test a spell with detailed debug information\n")
	sb.WriteString("  -dry-run              Show what would be executed without running it\n")
	sb.WriteString("\n")

	// Performance and diagnostics
	sb.WriteString("📊 Performance & Diagnostics:\n")
	sb.WriteString("  -benchmark            Run comprehensive performance benchmarks\n")
	sb.WriteString("  -test-hotkey          Test hotkey detection and registration\n")
	sb.WriteString("  -duration=<seconds>   Test duration for hotkey testing (0 = until Ctrl+C)\n")
	sb.WriteString("\n")

	// Output formatting options
	sb.WriteString("🎨 Output Formatting:\n")
	sb.WriteString("  -format=<format>      Output format: human, json, yaml (for show-config)\n")
	sb.WriteString("  -version-format=<fmt> Version format: human, json, compact\n")
	sb.WriteString("  -show-paths           Show configuration search paths\n")
	sb.WriteString("  -filter=<text>        Filter spells by sequence, name, or description\n")
	sb.WriteString("\n")
	
	// Export/Import options
	sb.WriteString("💾 Backup & Restore:\n")
	sb.WriteString("  -export-config=<file> Export configuration (use '-' for stdout)\n")
	sb.WriteString("  -export-format=<fmt>  Export format: yaml, tar.gz (default: yaml)\n")
	sb.WriteString("  -import-config=<file> Import configuration (use '-' for stdin)\n")
	sb.WriteString("\n")
	
	// Update management
	sb.WriteString("🔄 Update Management:\n")
	sb.WriteString("  -check-update         Check for available updates\n")
	sb.WriteString("  -force-update-check   Force update check (ignore cache)\n")
	sb.WriteString("\n")
	
	// Service management 
	sb.WriteString("🔧 Service Management:\n")
	sb.WriteString("  -service-install      Install as system service\n")
	sb.WriteString("  -service-uninstall    Remove system service\n")
	sb.WriteString("  -service-start        Start service\n")
	sb.WriteString("  -service-stop         Stop service\n")
	sb.WriteString("  -service-status       Show service status\n")
	sb.WriteString("\n")

	// Common workflows section
	sb.WriteString("🔄 Common Workflows:\n")
	sb.WriteString("\n")
	
	sb.WriteString("  Getting Started:\n")
	sb.WriteString(fmt.Sprintf("    %s --validate-config     # Check your spellbook.yml\n", os.Args[0]))
	sb.WriteString(fmt.Sprintf("    %s --show-config         # View merged configuration\n", os.Args[0]))
	sb.WriteString(fmt.Sprintf("    %s --list-spells         # See all available spells\n", os.Args[0]))
	sb.WriteString(fmt.Sprintf("    %s --test-hotkey         # Test hotkey detection\n", os.Args[0]))
	sb.WriteString("\n")
	
	sb.WriteString("  Development & Testing:\n")
	sb.WriteString(fmt.Sprintf("    %s --debug --no-tray     # Debug mode without tray\n", os.Args[0]))
	sb.WriteString(fmt.Sprintf("    %s --test-spell --spell \"e\" # Test specific spell\n", os.Args[0]))
	sb.WriteString(fmt.Sprintf("    %s --dry-run --spell \"g,s\" # Preview spell execution\n", os.Args[0]))
	sb.WriteString(fmt.Sprintf("    %s --once --spell \"build\" # Execute spell once\n", os.Args[0]))
	sb.WriteString("\n")
	
	sb.WriteString("  Configuration Management:\n")
	sb.WriteString(fmt.Sprintf("    %s --show-config --format json # Export config as JSON\n", os.Args[0]))
	sb.WriteString(fmt.Sprintf("    %s --show-config-path      # Find config file location\n", os.Args[0]))
	sb.WriteString(fmt.Sprintf("    %s --list-spells --filter git # Find git-related spells\n", os.Args[0]))
	sb.WriteString(fmt.Sprintf("    %s --export-config backup.yml # Backup configuration\n", os.Args[0]))
	sb.WriteString(fmt.Sprintf("    %s --export-config - --export-format yaml # Export to stdout\n", os.Args[0]))
	sb.WriteString(fmt.Sprintf("    %s --export-config backup.tar.gz --export-format tar.gz # Archive\n", os.Args[0]))
	sb.WriteString(fmt.Sprintf("    %s --import-config backup.yml # Restore from backup\n", os.Args[0]))
	sb.WriteString(fmt.Sprintf("    %s --import-config backup.tar.gz # Import from archive\n", os.Args[0]))
	sb.WriteString(fmt.Sprintf("    cat config.yml | %s --import-config - # Import from stdin\n", os.Args[0]))
	sb.WriteString("\n")
	
	sb.WriteString("  Performance Analysis:\n")
	sb.WriteString(fmt.Sprintf("    %s --benchmark              # Run performance tests\n", os.Args[0]))
	sb.WriteString(fmt.Sprintf("    %s --version --version-format json # Detailed build info\n", os.Args[0]))
	sb.WriteString("\n")

	// Platform-specific examples
	sb.WriteString("🖥️  Platform-Specific Examples:\n")
	sb.WriteString("\n")
	
	sb.WriteString("  Windows:\n")
	sb.WriteString(fmt.Sprintf("    %s --once --spell \"notepad\"   # Open Notepad\n", os.Args[0]))
	sb.WriteString(fmt.Sprintf("    %s --once --spell \"cmd\"       # Open Command Prompt\n", os.Args[0]))
	sb.WriteString("    # Spells: \"explorer\", \"powershell\", \"taskmgr\"\n")
	sb.WriteString("\n")
	sb.WriteString("    # Service management (run as Administrator):\n")
	sb.WriteString(fmt.Sprintf("    %s --service-install    # Install as service\n", os.Args[0]))
	sb.WriteString(fmt.Sprintf("    %s --service-start      # Start service\n", os.Args[0]))
	sb.WriteString(fmt.Sprintf("    %s --service-status     # Check status\n", os.Args[0]))
	sb.WriteString("\n")
	
	sb.WriteString("  macOS:\n")
	sb.WriteString(fmt.Sprintf("    %s --once --spell \"finder\"    # Open Finder\n", os.Args[0]))
	sb.WriteString(fmt.Sprintf("    %s --once --spell \"terminal\"  # Open Terminal\n", os.Args[0]))
	sb.WriteString("    # Spells: \"safari\", \"activity\", \"system\"\n")
	sb.WriteString("\n")
	sb.WriteString("    # Service management (LaunchAgent):\n")
	sb.WriteString(fmt.Sprintf("    %s --service-install    # Install LaunchAgent\n", os.Args[0]))
	sb.WriteString(fmt.Sprintf("    %s --service-start      # Start service\n", os.Args[0]))
	sb.WriteString(fmt.Sprintf("    %s --service-status     # Check status\n", os.Args[0]))
	sb.WriteString("\n")
	
	sb.WriteString("  Cross-Platform:\n")
	sb.WriteString(fmt.Sprintf("    %s --once --spell \"e\"         # Editor (VS Code/configured)\n", os.Args[0]))
	sb.WriteString(fmt.Sprintf("    %s --once --spell \"g,s\"       # Git status\n", os.Args[0]))
	sb.WriteString(fmt.Sprintf("    %s --once --spell \"browser\"   # Default browser\n", os.Args[0]))
	sb.WriteString("    # Common spells: \"t\" (terminal), \"c\" (calculator)\n")
	sb.WriteString("\n")

	// Spell configuration examples
	sb.WriteString("📖 Spell Configuration Examples:\n")
	sb.WriteString("\n")
	sb.WriteString("  Basic spellbook.yml structure:\n")
	sb.WriteString("    spells:\n")
	sb.WriteString("      e: editor        # Single key spell\n")
	sb.WriteString("      \"g,s\": git_status # Sequential key spell\n")
	sb.WriteString("      \"vs,code\": vscode # Multi-key spell\n")
	sb.WriteString("\n")
	sb.WriteString("    grimoire:\n")
	sb.WriteString("      editor:\n")
	sb.WriteString("        type: app\n")
	sb.WriteString("        app: \"code\"\n")
	sb.WriteString("      git_status:\n")
	sb.WriteString("        type: script\n")
	sb.WriteString("        script: \"git status\"\n")
	sb.WriteString("\n")

	// Troubleshooting section
	sb.WriteString("🔧 Troubleshooting:\n")
	sb.WriteString("  • Hotkeys not working: Check permissions, try --test-hotkey\n")
	sb.WriteString("  • Config errors: Use --validate-config to check syntax\n")
	sb.WriteString("  • Spells not found: Use --list-spells to see available spells\n")
	sb.WriteString("  • Performance issues: Run --benchmark to analyze\n")
	sb.WriteString("  • Debug mode: Add --debug flag for detailed logging\n")
	sb.WriteString("\n")

	// Footer
	sb.WriteString("📚 For more information:\n")
	sb.WriteString("  • Documentation: docs/\n")
	sb.WriteString("  • Configuration examples: examples/config/\n")
	sb.WriteString("  • Issue reporting: GitHub issues\n")

	return sb.String()
}
