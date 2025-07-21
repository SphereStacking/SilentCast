package commands

import (
	"fmt"

	"github.com/SphereStacking/silentcast/internal/version"
)

// VersionCommand displays version information
type VersionCommand struct {
	version string
}

// NewVersionCommand creates a new version command
func NewVersionCommand(version string) Command {
	return &VersionCommand{
		version: version,
	}
}

// Name returns the command name
func (c *VersionCommand) Name() string {
	return "Version"
}

// Description returns the command description
func (c *VersionCommand) Description() string {
	return "Print detailed version and build information"
}

// FlagName returns the flag name
func (c *VersionCommand) FlagName() string {
	return "version"
}

// IsActive checks if the command should run
func (c *VersionCommand) IsActive(flags interface{}) bool {
	f, ok := flags.(*Flags)
	if !ok {
		return false
	}
	return f.Version
}

// Execute runs the command
func (c *VersionCommand) Execute(flags interface{}) error {
	f, ok := flags.(*Flags)
	if !ok {
		return fmt.Errorf("invalid flags type")
	}

	buildInfo := version.GetBuildInfo()

	var output string
	var err error

	switch f.VersionFormat {
	case "json":
		output, err = buildInfo.FormatJSON()
		if err != nil {
			return fmt.Errorf("failed to format version as JSON: %w", err)
		}
	case "compact":
		output = buildInfo.FormatCompact()
	case "human":
		fallthrough
	default:
		output = buildInfo.FormatHuman()
	}

	fmt.Print(output)
	if f.VersionFormat != "human" {
		fmt.Println() // Add newline for non-human formats
	}

	return nil
}

// Group returns the command group
func (c *VersionCommand) Group() string {
	return "core"
}

// HasOptions returns if this command has additional options
func (c *VersionCommand) HasOptions() bool {
	return true
}
