//go:build darwin

package shell

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
	// macOS: use Terminal.app via AppleScript, default to keep open
	return d.WrapInTerminalWithOptions(ctx, cmd, true)
}

func (d *darwinShellExecutor) WrapInTerminalWithOptions(ctx context.Context, cmd *exec.Cmd, keepOpen bool) *exec.Cmd {
	// macOS: use Terminal.app via AppleScript
	cmdStr := strings.ReplaceAll(cmd.String(), `"`, `\"`)

	if keepOpen {
		// Add a read command to keep terminal open
		cmdStr = fmt.Sprintf(`%s; echo ""; echo "Press Enter to close..."; read`, cmdStr)
	}

	script := fmt.Sprintf(`tell application "Terminal" to do script "%s"`, cmdStr)
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
