//go:build darwin

package app

import (
	"context"
	"os/exec"
	"strings"
)

func init() {
	launcherFactory = func() AppLauncher {
		return &darwinLauncher{}
	}
}

type darwinLauncher struct{}

func (d *darwinLauncher) PrepareCommand(ctx context.Context, path string, args []string) *exec.Cmd {
	// Use 'open' for .app bundles
	if strings.HasSuffix(path, ".app") {
		openArgs := []string{}
		if len(args) > 0 {
			openArgs = append(openArgs, "-a", path, "--args")
			openArgs = append(openArgs, args...)
		} else {
			openArgs = append(openArgs, path)
		}
		return exec.CommandContext(ctx, "open", openArgs...)
	}

	// Regular executable
	if len(args) > 0 {
		return exec.CommandContext(ctx, path, args...)
	}
	return exec.CommandContext(ctx, path)
}

func (d *darwinLauncher) IsSpecialPath(path string) bool {
	// System applications and bundles
	return strings.HasPrefix(path, "/System/") ||
		strings.HasPrefix(path, "/Applications/") ||
		strings.HasSuffix(path, ".app")
}

func (d *darwinLauncher) RequiresShell(path string) bool {
	// .app bundles need to be opened via 'open' command
	return strings.HasSuffix(path, ".app")
}
