//go:build !darwin && !windows && !linux

package shell

import (
	"context"
	"os"
	"os/exec"
	"strings"
)

func init() {
	shellFactory = func() ShellExecutor {
		return &stubShellExecutor{}
	}
}

type stubShellExecutor struct{}

func (s *stubShellExecutor) GetShell() (string, string) {
	// Check for custom shell
	if customShell := os.Getenv("SHELL"); customShell != "" {
		return customShell, "-c"
	}
	return "sh", "-c"
}

func (s *stubShellExecutor) WrapInTerminal(ctx context.Context, cmd *exec.Cmd) *exec.Cmd {
	// Stub implementation: just return the command as-is
	return cmd
}

func (s *stubShellExecutor) IsInteractiveCommand(command string) bool {
	// Common interactive commands
	interactiveCommands := []string{"vim", "nano", "emacs", "htop", "top", "less", "more"}
	for _, icmd := range interactiveCommands {
		if strings.Contains(command, icmd) {
			return true
		}
	}
	return false
}
