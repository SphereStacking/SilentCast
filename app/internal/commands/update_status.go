package commands

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/SphereStacking/silentcast/internal/updater"
	"github.com/SphereStacking/silentcast/internal/version"
)

// UpdateStatusCommand implements the update status command
type UpdateStatusCommand struct {
	configPath string
}

// NewUpdateStatusCommand creates a new update status command
func NewUpdateStatusCommand(configPathFunc func() string) Command {
	return &UpdateStatusCommand{
		configPath: configPathFunc(),
	}
}

// Name returns the command name
func (c *UpdateStatusCommand) Name() string {
	return "UpdateStatus"
}

// Description returns the command description
func (c *UpdateStatusCommand) Description() string {
	return "Show current update status and available updates"
}

// HasOptions returns whether this command supports options
func (c *UpdateStatusCommand) HasOptions() bool {
	return true
}

// FlagName returns the flag name for this command
func (c *UpdateStatusCommand) FlagName() string {
	return "update-status"
}

// Group returns the command group
func (c *UpdateStatusCommand) Group() string {
	return "utility"
}

// IsActive checks if the command should be executed based on flags
func (c *UpdateStatusCommand) IsActive(flags interface{}) bool {
	if f, ok := flags.(*Flags); ok {
		return f.UpdateStatus
	}

	// Handle test flags (map) for unit tests
	if flagsMap, ok := flags.(map[string]interface{}); ok {
		return flagsMap["update-status"] == true
	}

	return false
}

// Execute runs the command
func (c *UpdateStatusCommand) Execute(flags interface{}) error {
	// Handle test flags (map) for unit tests
	var verbose bool
	var forceCheck bool

	switch v := flags.(type) {
	case map[string]interface{}:
		verbose = v["verbose"] == true
		forceCheck = v["force-check"] == true
	case *Flags:
		verbose = v.Debug
		forceCheck = v.ForceSelfUpdate
	case nil:
		// defaults are already set
	default:
		return fmt.Errorf("invalid flags type: expected *Flags or map[string]interface{}, got %T", flags)
	}

	fmt.Println("🔄 SilentCast Update Status")
	fmt.Println("==========================")

	// Show current version
	currentVersion := version.GetVersionString()
	fmt.Printf("\n📍 Current Version: %s\n", currentVersion)

	// Show build information if verbose
	if verbose {
		fmt.Printf("   Go Version: %s\n", runtime.Version())
		fmt.Printf("   Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	}

	// Create updater
	cfg := updater.Config{
		CurrentVersion: currentVersion,
		RepoOwner:      "SphereStacking",
		RepoName:       "SilentCast",
		CheckInterval:  24 * time.Hour,
		ConfigDir:      c.configPath,
		CacheDuration:  1 * time.Hour,
	}

	upd := updater.NewUpdater(&cfg)

	// Check for updates
	fmt.Println("\n🔍 Checking for updates...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var updateInfo *updater.UpdateInfo
	var err error

	if forceCheck {
		fmt.Println("   (forcing fresh check, ignoring cache)")
		updateInfo, err = upd.ForceCheck(ctx)
	} else {
		updateInfo, err = upd.CheckForUpdate(ctx)
	}

	if err != nil {
		fmt.Printf("❌ Failed to check for updates: %v\n", err)
		return nil // Don't fail the command for network errors
	}

	if updateInfo == nil {
		fmt.Printf("✅ You're running the latest version (%s)\n", currentVersion)

		// Show cache information if verbose
		if verbose {
			c.showCacheInfo(upd)
		}

		return nil
	}

	// Update available
	fmt.Println("\n📦 Update Available!")
	fmt.Printf("   Current: %s\n", currentVersion)
	fmt.Printf("   Latest:  %s\n", updateInfo.Version)
	fmt.Printf("   Published: %s\n", updateInfo.PublishedAt.Format("2006-01-02 15:04"))

	if updateInfo.Size > 0 {
		fmt.Printf("   Size: %s\n", formatUpdateSize(updateInfo.Size))
	}

	// Show release notes
	if updateInfo.ReleaseNotes != "" {
		fmt.Println("\n📋 Release Notes:")
		notes := updateInfo.ReleaseNotes
		if len(notes) > 500 && !verbose {
			notes = notes[:497] + "..."
			fmt.Println("   (Use --verbose to see full release notes)")
		}
		fmt.Println(indentText(notes, "   "))
	}

	// Show update commands
	fmt.Println("\n🚀 To Update:")
	fmt.Println("   ./silentcast --self-update")
	fmt.Println("   ./silentcast --self-update --force-self-update  # No confirmation")

	// Show cache information if verbose
	if verbose {
		c.showCacheInfo(upd)
	}

	return nil
}

// showCacheInfo displays cache information
func (c *UpdateStatusCommand) showCacheInfo(_ *updater.Updater) {
	fmt.Println("\n🗂️  Cache Information:")
	fmt.Println("   Cache directory: " + c.configPath)
	fmt.Println("   Cache duration: 1 hour")
	fmt.Println("   Clear cache: ./silentcast --check-update --force-update-check")
}

// indentText adds indentation to each line of text
func indentText(text, indent string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = indent + line
		}
	}
	return strings.Join(lines, "\n")
}

// formatUpdateSize formats byte size to human-readable string
func formatUpdateSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}

	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := []string{"KB", "MB", "GB", "TB"}
	return fmt.Sprintf("%.1f %s", float64(size)/float64(div), units[exp])
}
