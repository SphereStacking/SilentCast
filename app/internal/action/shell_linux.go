//go:build linux

package action

import (
	"context"
	"os"
	"os/exec"
	"strings"
)

func init() {
	shellFactory = func() ShellExecutor {
		return &linuxShellExecutor{}
	}
}

type linuxShellExecutor struct{}

func (l *linuxShellExecutor) GetShell() (shell, flag string) {
	// Check for custom shell
	if customShell := os.Getenv("SHELL"); customShell != "" {
		return customShell, "-c"
	}
	return "sh", "-c"
}

func (l *linuxShellExecutor) WrapInTerminal(ctx context.Context, cmd *exec.Cmd) *exec.Cmd {
	// Linux: try common terminal emulators
	terminals := []string{"gnome-terminal", "konsole", "xterm", "xfce4-terminal"}

	for _, term := range terminals {
		if _, err := exec.LookPath(term); err == nil {
			switch term {
			case "gnome-terminal":
				args := []string{"--"}
				args = append(args, cmd.Path)
				args = append(args, cmd.Args[1:]...)
				return exec.CommandContext(ctx, term, args...)
			case "konsole":
				return exec.CommandContext(ctx, term, "-e", cmd.String())
			case "xterm", "xfce4-terminal":
				return exec.CommandContext(ctx, term, "-e", cmd.String())
			}
		}
	}

	// Fallback: run without terminal
	return cmd
}

func (l *linuxShellExecutor) IsInteractiveCommand(command string) bool {
	// Common interactive commands
	interactiveCommands := []string{"vim", "nano", "emacs", "htop", "top", "less", "more"}
	for _, icmd := range interactiveCommands {
		if strings.Contains(command, icmd) {
			return true
		}
	}
	return false
}
