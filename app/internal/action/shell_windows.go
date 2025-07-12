//go:build windows

package action

import (
	"context"
	"os/exec"
	"strings"
)

func init() {
	shellFactory = func() ShellExecutor {
		return &windowsShellExecutor{}
	}
}

type windowsShellExecutor struct{}

func (w *windowsShellExecutor) GetShell() (string, string) {
	return "cmd", "/c"
}

func (w *windowsShellExecutor) WrapInTerminal(ctx context.Context, cmd *exec.Cmd) *exec.Cmd {
	// Windows: use cmd with start
	return exec.CommandContext(ctx, "cmd", "/c", "start", "cmd", "/k", cmd.String())
}

func (w *windowsShellExecutor) IsInteractiveCommand(command string) bool {
	// Windows interactive commands
	interactiveCommands := []string{"vim", "nano", "more", "edit"}
	cmdLower := strings.ToLower(command)
	for _, icmd := range interactiveCommands {
		if strings.Contains(cmdLower, icmd) {
			return true
		}
	}
	return false
}
