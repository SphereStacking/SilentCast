//go:build darwin

package action

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func init() {
	shellFactory = func() ShellExecutor {
		return &darwinShellExecutor{}
	}
}

type darwinShellExecutor struct{}

func (d *darwinShellExecutor) GetShell() (string, string) {
	// Check for custom shell
	if customShell := os.Getenv("SHELL"); customShell != "" {
		return customShell, "-c"
	}
	return "sh", "-c"
}

func (d *darwinShellExecutor) WrapInTerminal(ctx context.Context, cmd *exec.Cmd) *exec.Cmd {
	// macOS: use Terminal.app via AppleScript
	script := fmt.Sprintf(`tell application "Terminal" to do script "%s"`,
		strings.ReplaceAll(cmd.String(), `"`, `\"`))
	return exec.CommandContext(ctx, "osascript", "-e", script)
}

func (d *darwinShellExecutor) IsInteractiveCommand(command string) bool {
	// Common interactive commands
	interactiveCommands := []string{"vim", "nano", "emacs", "htop", "top", "less", "more"}
	for _, icmd := range interactiveCommands {
		if strings.Contains(command, icmd) {
			return true
		}
	}
	return false
}
