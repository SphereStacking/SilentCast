package shell

import (
	"context"
	"os/exec"
)

// ShellExecutor defines the interface for platform-specific shell execution
type ShellExecutor interface {
	// GetShell returns the shell command and flag for the platform
	GetShell() (shell string, flag string)

	// WrapInTerminal wraps a command to run in a new terminal window
	WrapInTerminal(ctx context.Context, cmd *exec.Cmd) *exec.Cmd

	// WrapInTerminalWithOptions wraps a command with terminal options
	WrapInTerminalWithOptions(ctx context.Context, cmd *exec.Cmd, keepOpen bool) *exec.Cmd

	// IsInteractiveCommand checks if a command should run in a terminal
	IsInteractiveCommand(command string) bool
}

// shellFactory creates the appropriate shell executor for the current platform
var shellFactory func() ShellExecutor

// GetShellExecutor returns the platform-specific shell executor
func GetShellExecutor() ShellExecutor {
	if shellFactory == nil {
		panic("shell factory not initialized")
	}
	return shellFactory()
}
