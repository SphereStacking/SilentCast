package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/SphereStacking/silentcast/internal/updater"
	"github.com/SphereStacking/silentcast/internal/version"
)

// SelfUpdateCommand implements the self-update command
type SelfUpdateCommand struct {
	configPath string
	force      bool
}

// NewSelfUpdateCommand creates a new self update command
func NewSelfUpdateCommand(configPathFunc func() string) Command {
	return &SelfUpdateCommand{
		configPath: configPathFunc(),
	}
}

// Name returns the command name
func (c *SelfUpdateCommand) Name() string {
	return "SelfUpdate"
}

// Description returns the command description
func (c *SelfUpdateCommand) Description() string {
	return "Update SilentCast to the latest version"
}

// FlagName returns the flag name
func (c *SelfUpdateCommand) FlagName() string {
	return "self-update"
}

// IsActive checks if the command should run
func (c *SelfUpdateCommand) IsActive(flags interface{}) bool {
	f, ok := flags.(*Flags)
	if !ok {
		return false
	}
	// Store force flag if present
	c.force = f.ForceSelfUpdate
	return f.SelfUpdate
}

// Execute runs the command
func (c *SelfUpdateCommand) Execute(flags interface{}) error {
	// Handle test flags (map) for unit tests
	var checkOnly bool
	var forceUpdate bool
	
	if flagsMap, ok := flags.(map[string]interface{}); ok {
		checkOnly = flagsMap["check-only"] == true
		forceUpdate = flagsMap["force"] == true
	} else if f, ok := flags.(*Flags); ok {
		forceUpdate = f.ForceSelfUpdate
	} else if flags != nil {
		return fmt.Errorf("invalid flags type: expected *Flags or map[string]interface{}, got %T", flags)
	}
	
	fmt.Println("🚀 SilentCast Self-Update")
	fmt.Println("========================")
	
	// Check for updates first
	fmt.Println("\n🔍 Checking for updates...")
	
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

	updateInfo, err := upd.CheckForUpdate(ctx)
	if err != nil {
		return fmt.Errorf("update check failed: %w", err)
	}

	if updateInfo == nil {
		fmt.Printf("\n✅ You're already running the latest version (%s)\n", cfg.CurrentVersion)
		return nil
	}

	// Show update information
	fmt.Println("\n📦 Update Available!")
	fmt.Printf("  Current version: %s\n", cfg.CurrentVersion)
	fmt.Printf("  Latest version:  %s\n", updateInfo.Version)
	fmt.Printf("  Published:       %s\n", updateInfo.PublishedAt.Format("2006-01-02"))
	
	if updateInfo.Size > 0 {
		fmt.Printf("  Download size:   %s\n", formatBytes(updateInfo.Size))
	}

	// Show release notes
	fmt.Println("\n📋 Release Notes:")
	if updateInfo.ReleaseNotes != "" {
		notes := updateInfo.ReleaseNotes
		if len(notes) > 500 {
			notes = notes[:497] + "..."
		}
		fmt.Println(notes)
	} else {
		fmt.Println("No release notes available.")
	}

	// If check-only mode, return here
	if checkOnly {
		return nil
	}

	// Confirm update
	if !c.force && !forceUpdate {
		fmt.Print("\n❓ Do you want to update now? (y/N): ")
		reader := bufio.NewReader(os.Stdin)
		response, readErr := reader.ReadString('\n')
		if readErr != nil {
			return fmt.Errorf("failed to read response: %w", readErr)
		}
		
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("❌ Update cancelled")
			return nil
		}
	}

	// Download update with progress
	fmt.Println("\n📥 Downloading update...")
	
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	
	downloadPath, err := upd.DownloadUpdateWithProgress(ctx, updateInfo, os.Stdout)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	
	// Apply update
	fmt.Println("\n🔧 Applying update...")
	
	// Warn about restart
	fmt.Println("\n⚠️  IMPORTANT: SilentCast will restart after the update")
	fmt.Println("   Any running spells will be interrupted")
	
	if !c.force {
		fmt.Print("\n❓ Continue with update? (y/N): ")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}
		
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			// Clean up downloaded file
			os.Remove(downloadPath)
			fmt.Println("❌ Update cancelled")
			return nil
		}
	}
	
	// Apply the update
	if err := upd.ApplyUpdate(downloadPath); err != nil {
		return fmt.Errorf("failed to apply update: %w", err)
	}
	
	fmt.Println("\n✅ Update applied successfully!")
	fmt.Println("   SilentCast will now restart...")
	
	// Restart the application
	// The platform-specific updater should handle the restart
	platform := updater.GetPlatformUpdater()
	if err := platform.RestartApplication(); err != nil {
		fmt.Printf("\n⚠️  Could not restart automatically: %v\n", err)
		fmt.Println("   Please restart SilentCast manually")
	}
	
	return nil
}

// Group returns the command group
func (c *SelfUpdateCommand) Group() string {
	return "utility"
}

// HasOptions returns if this command has additional options
func (c *SelfUpdateCommand) HasOptions() bool {
	return true
}