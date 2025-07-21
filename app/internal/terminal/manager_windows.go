//go:build windows

package terminal

import (
	"context"
	"os/exec"
)

// WindowsManager implements the Manager interface for Windows
type WindowsManager struct {
	*baseManager
}

// NewWindowsManager creates a new Windows terminal manager
func NewWindowsManager() Manager {
	detector := NewWindowsDetector()
	builder := NewWindowsCommandBuilder()
	return &WindowsManager{
		baseManager: newBaseManager(detector, builder),
	}
}

// ExecuteInTerminal overrides to handle Windows-specific behavior
func (m *WindowsManager) ExecuteInTerminal(ctx context.Context, cmd *exec.Cmd, options *Options) error {
	// For Windows Terminal (wt.exe), we need to use it directly
	terminal, err := m.selectTerminal(options)
	if err != nil {
		return err
	}

	// Special handling for Windows Terminal
	if terminal.Command == "wt.exe" {
		// Windows Terminal requires special handling
		args, err := m.builder.BuildCommand(terminal, cmd, options)
		if err != nil {
			return &Error{
				Op:  "ExecuteInTerminal",
				Err: err,
				Msg: "failed to build Windows Terminal command",
			}
		}

		wtCmd := exec.CommandContext(ctx, "wt.exe", args...)
		if options.WorkingDir != "" {
			wtCmd.Dir = options.WorkingDir
		} else if cmd.Dir != "" {
			wtCmd.Dir = cmd.Dir
		}

		return wtCmd.Start()
	}

	// For other terminals, use the base implementation
	return m.baseManager.ExecuteInTerminal(ctx, cmd, options)
}
