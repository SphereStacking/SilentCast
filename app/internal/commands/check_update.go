package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/SphereStacking/silentcast/internal/updater"
	"github.com/SphereStacking/silentcast/internal/version"
)

// CheckUpdateCommand implements the check-update command
type CheckUpdateCommand struct {
	configPath string
	force      bool
}

// NewCheckUpdateCommand creates a new check update command
func NewCheckUpdateCommand(configPathFunc func() string) Command {
	return &CheckUpdateCommand{
		configPath: configPathFunc(),
	}
}

// Name returns the command name
func (c *CheckUpdateCommand) Name() string {
	return "CheckUpdate"
}

// Description returns the command description
func (c *CheckUpdateCommand) Description() string {
	return "Check for available updates"
}

// FlagName returns the flag name
func (c *CheckUpdateCommand) FlagName() string {
	return "check-update"
}

// IsActive checks if the command should run
func (c *CheckUpdateCommand) IsActive(flags interface{}) bool {
	f, ok := flags.(*Flags)
	if !ok {
		return false
	}
	// Store force flag if present
	c.force = f.ForceUpdateCheck
	return f.CheckUpdate
}

// Execute runs the command
func (c *CheckUpdateCommand) Execute(flags interface{}) error {
	fmt.Println("🔍 Checking for updates...")

	// Create updater config
	cfg := updater.Config{
		CurrentVersion: version.GetVersionString(),
		RepoOwner:      "SphereStacking",
		RepoName:       "SilentCast",
		CheckInterval:  24 * time.Hour,
		ConfigDir:      c.configPath,
		CacheDuration:  1 * time.Hour,
	}

	// Create updater
	upd := updater.NewUpdater(cfg)

	// Check for updates
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var updateInfo *updater.UpdateInfo
	var err error

	if c.force {
		fmt.Println("  Forcing fresh update check (ignoring cache)...")
		updateInfo, err = upd.ForceCheck(ctx)
	} else {
		updateInfo, err = upd.CheckForUpdate(ctx)
	}

	if err != nil {
		return fmt.Errorf("update check failed: %w", err)
	}

	// Display results
	if updateInfo == nil {
		fmt.Printf("✅ You're running the latest version (%s)\n", cfg.CurrentVersion)
		return nil
	}

	// New version available
	fmt.Println("\n🎉 New version available!")
	fmt.Printf("  Current version: %s\n", cfg.CurrentVersion)
	fmt.Printf("  Latest version:  %s\n", updateInfo.Version)
	fmt.Printf("  Published:       %s\n", updateInfo.PublishedAt.Format("2006-01-02"))
	
	if updateInfo.Size > 0 {
		fmt.Printf("  Download size:   %s\n", formatBytes(updateInfo.Size))
	}

	fmt.Println("\n📋 Release Notes:")
	if updateInfo.ReleaseNotes != "" {
		// Limit release notes to first 500 chars
		notes := updateInfo.ReleaseNotes
		if len(notes) > 500 {
			notes = notes[:497] + "..."
		}
		fmt.Println(notes)
	} else {
		fmt.Println("  No release notes available")
	}

	fmt.Println("\n💡 To update, visit: https://github.com/SphereStacking/SilentCast/releases")

	return nil
}

// Group returns the command group
func (c *CheckUpdateCommand) Group() string {
	return "utility"
}

// HasOptions returns if this command has additional options
func (c *CheckUpdateCommand) HasOptions() bool {
	return true
}

// formatBytes formats bytes to human readable format
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}