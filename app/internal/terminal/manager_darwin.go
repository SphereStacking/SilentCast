//go:build darwin

package terminal

import (
	"context"
	"os/exec"
	"strings"
)

// DarwinManager implements the Manager interface for macOS
type DarwinManager struct {
	*baseManager
}

// NewDarwinManager creates a new macOS terminal manager
func NewDarwinManager() Manager {
	detector := NewMacOSDetector()
	builder := NewMacOSCommandBuilder()
	return &DarwinManager{
		baseManager: newBaseManager(detector, builder),
	}
}

// ExecuteInTerminal overrides to handle macOS-specific behavior
func (m *DarwinManager) ExecuteInTerminal(ctx context.Context, cmd *exec.Cmd, options Options) error {
	// Get the terminal to use
	terminal, err := m.selectTerminal(options)
	if err != nil {
		return err
	}

	// For Terminal.app and iTerm2, we use osascript
	if strings.Contains(terminal.Command, "Terminal.app") || strings.Contains(terminal.Command, "iTerm") {
		args, err := m.builder.BuildCommand(terminal, cmd, options)
		if err != nil {
			return &Error{
				Op:  "ExecuteInTerminal",
				Err: err,
				Msg: "failed to build AppleScript command",
			}
		}

		// Execute via osascript
		osascriptCmd := exec.CommandContext(ctx, "osascript", args...)
		if options.WorkingDir != "" {
			osascriptCmd.Dir = options.WorkingDir
		} else if cmd.Dir != "" {
			osascriptCmd.Dir = cmd.Dir
		}

		if err := osascriptCmd.Start(); err != nil {
			return &Error{
				Op:  "ExecuteInTerminal",
				Err: err,
				Msg: "failed to execute AppleScript",
			}
		}

		// Detach from process
		if err := osascriptCmd.Process.Release(); err != nil {
			// Non-fatal
			_ = err
		}

		return nil
	}

	// For other terminals, use the base implementation
	return m.baseManager.ExecuteInTerminal(ctx, cmd, options)
}
