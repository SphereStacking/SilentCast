//go:build linux

package shell

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
	// Linux: default to keep open
	return l.WrapInTerminalWithOptions(ctx, cmd, true)
}

func (l *linuxShellExecutor) WrapInTerminalWithOptions(ctx context.Context, cmd *exec.Cmd, keepOpen bool) *exec.Cmd {
	// Linux: try common terminal emulators
	terminals := []string{"gnome-terminal", "konsole", "xterm", "xfce4-terminal"}

	for _, term := range terminals {
		if _, err := exec.LookPath(term); err == nil {
			switch term {
			case "gnome-terminal":
				if keepOpen {
					// Use bash to execute command and wait for input
					return exec.CommandContext(ctx, term, "--", "bash", "-c",
						cmd.String()+`; echo ""; echo "Press Enter to close..."; read`)
				}
				args := []string{"--"}
				args = append(args, cmd.Path)
				args = append(args, cmd.Args[1:]...)
				return exec.CommandContext(ctx, term, args...)
			case "konsole":
				if keepOpen {
					// Use --hold to keep konsole open
					return exec.CommandContext(ctx, term, "--hold", "-e", cmd.String())
				}
				return exec.CommandContext(ctx, term, "-e", cmd.String())
			case "xterm":
				if keepOpen {
					// Use -hold to keep xterm open
					return exec.CommandContext(ctx, term, "-hold", "-e", cmd.String())
				}
				return exec.CommandContext(ctx, term, "-e", cmd.String())
			case "xfce4-terminal":
				if keepOpen {
					// Use --hold to keep terminal open
					return exec.CommandContext(ctx, term, "--hold", "-e", cmd.String())
				}
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
