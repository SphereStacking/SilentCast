package action

import (
	"context"
	"os/exec"
)

// AppLauncher defines the interface for platform-specific application launching
type AppLauncher interface {
	// PrepareCommand prepares the command for launching an application
	PrepareCommand(ctx context.Context, path string, args []string) *exec.Cmd

	// IsSpecialPath checks if the path requires special handling
	IsSpecialPath(path string) bool

	// RequiresShell checks if the application should be launched via shell
	RequiresShell(path string) bool
}

// launcherFactory creates the appropriate app launcher for the current platform
var launcherFactory func() AppLauncher

// GetAppLauncher returns the platform-specific app launcher
func GetAppLauncher() AppLauncher {
	if launcherFactory == nil {
		panic("launcher factory not initialized")
	}
	return launcherFactory()
}
