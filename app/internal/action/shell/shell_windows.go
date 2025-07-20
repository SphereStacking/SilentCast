//go:build windows

package shell

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
	// Windows: use cmd with start, default to keep open
	return w.WrapInTerminalWithOptions(ctx, cmd, true)
}

func (w *windowsShellExecutor) WrapInTerminalWithOptions(ctx context.Context, cmd *exec.Cmd, keepOpen bool) *exec.Cmd {
	// Windows: use cmd with start
	// /k keeps window open, /c closes after execution
	flag := "/k"
	if !keepOpen {
		flag = "/c"
	}
	return exec.CommandContext(ctx, "cmd", "/c", "start", "cmd", flag, cmd.String())
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
