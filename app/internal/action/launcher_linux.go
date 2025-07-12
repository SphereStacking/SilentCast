//go:build linux

package action

import (
	"context"
	"os/exec"
	"strings"
)

func init() {
	launcherFactory = func() AppLauncher {
		return &linuxLauncher{}
	}
}

type linuxLauncher struct{}

func (l *linuxLauncher) PrepareCommand(ctx context.Context, path string, args []string) *exec.Cmd {
	// Check if it's a special path that should use xdg-open
	if l.IsSpecialPath(path) {
		return exec.CommandContext(ctx, "xdg-open", path)
	}

	// Normal application launch
	return exec.CommandContext(ctx, path, args...)
}

func (l *linuxLauncher) IsSpecialPath(path string) bool {
	// Check for URL schemes
	schemes := []string{"http://", "https://", "file://", "mailto:"}
	for _, scheme := range schemes {
		if strings.HasPrefix(path, scheme) {
			return true
		}
	}

	// Check for common file extensions that should open with default app
	extensions := []string{".pdf", ".html", ".png", ".jpg", ".jpeg", ".gif"}
	lowerPath := strings.ToLower(path)
	for _, ext := range extensions {
		if strings.HasSuffix(lowerPath, ext) {
			return true
		}
	}

	return false
}

func (l *linuxLauncher) RequiresShell(path string) bool {
	// On Linux, .desktop files might need special handling
	return strings.HasSuffix(path, ".desktop")
}
