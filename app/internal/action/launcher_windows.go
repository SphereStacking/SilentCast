//go:build windows

package action

import (
	"context"
	"os/exec"
	"strings"
)

func init() {
	launcherFactory = func() AppLauncher {
		return &windowsLauncher{}
	}
}

type windowsLauncher struct{}

func (w *windowsLauncher) PrepareCommand(ctx context.Context, path string, args []string) *exec.Cmd {
	// Use 'start' for special paths
	if w.RequiresShell(path) {
		return exec.CommandContext(ctx, "cmd", "/c", "start", "", path)
	}
	
	// Regular executable
	if len(args) > 0 {
		return exec.CommandContext(ctx, path, args...)
	}
	return exec.CommandContext(ctx, path)
}

func (w *windowsLauncher) IsSpecialPath(path string) bool {
	// Windows Store apps and URLs
	return strings.HasPrefix(path, "ms-") ||
		strings.HasPrefix(path, "http://") ||
		strings.HasPrefix(path, "https://")
}

func (w *windowsLauncher) RequiresShell(path string) bool {
	// Windows Store apps, URLs, and certain file types need shell
	return w.IsSpecialPath(path) ||
		strings.HasSuffix(strings.ToLower(path), ".url") ||
		strings.HasSuffix(strings.ToLower(path), ".lnk")
}