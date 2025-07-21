package commands

// Flags holds all command-line flags
// This is the single source of truth for all flag definitions
type Flags struct {
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
	Install   bool
	Uninstall bool
}
