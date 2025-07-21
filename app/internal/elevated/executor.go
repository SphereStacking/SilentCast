package elevated

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/SphereStacking/silentcast/pkg/logger"
)

// Executor is the interface for executing actions
type Executor interface {
	// Execute runs the action
	Execute(ctx context.Context) error

	// String returns a string representation of the action
	String() string
}

// ElevatedExecutor handles commands that require elevated privileges
type ElevatedExecutor struct {
	baseExecutor Executor
	isAdmin      bool
}

// NewElevatedExecutor wraps an executor to run with elevated privileges
func NewElevatedExecutor(baseExecutor Executor, needsAdmin bool) Executor {
	if !needsAdmin {
		// No elevation needed
		return baseExecutor
	}

	return &ElevatedExecutor{
		baseExecutor: baseExecutor,
		isAdmin:      needsAdmin,
	}
}

// Execute runs the command with elevated privileges
func (e *ElevatedExecutor) Execute(ctx context.Context) error {
	// Check if already running as admin
	if e.isRunningAsAdmin() {
		// Already elevated, execute directly
		return e.baseExecutor.Execute(ctx)
	}

	// Need to elevate
	return e.executeElevated(ctx)
}

// String returns a string representation
func (e *ElevatedExecutor) String() string {
	return fmt.Sprintf("[Admin] %s", e.baseExecutor.String())
}

// isRunningAsAdmin checks if the current process has admin privileges
func (e *ElevatedExecutor) isRunningAsAdmin() bool {
	switch runtime.GOOS {
	case "windows":
		// Check if running as Administrator on Windows
		// This is a simplified check
		_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
		if err != nil {
			return false
		}
		return true
	case "darwin", "linux":
		// Check if running as root
		return os.Geteuid() == 0
	default:
		return false
	}
}

// executeElevated runs the command with elevated privileges
func (e *ElevatedExecutor) executeElevated(ctx context.Context) error {
	switch runtime.GOOS {
	case "windows":
		return e.executeElevatedWindows(ctx)
	case "darwin":
		return e.executeElevatedDarwin(ctx)
	case "linux":
		return e.executeElevatedLinux(ctx)
	default:
		return fmt.Errorf("admin execution not supported on %s", runtime.GOOS)
	}
}

// executeElevatedWindows runs command as Administrator on Windows
func (e *ElevatedExecutor) executeElevatedWindows(ctx context.Context) error {
	// Get the command to execute
	cmdStr := e.getCommandString()
	if cmdStr == "" {
		return fmt.Errorf("cannot determine command to elevate")
	}

	logger.Debug("Requesting Administrator privileges for: %s", cmdStr)

	// Use PowerShell to run as Administrator
	psScript := fmt.Sprintf(`Start-Process cmd -ArgumentList '/c %s' -Verb RunAs -Wait`,
		escapeForPowerShell(cmdStr))

	cmd := exec.CommandContext(ctx, "powershell", "-Command", psScript)
	return cmd.Run()
}

// executeElevatedDarwin runs command with sudo on macOS
func (e *ElevatedExecutor) executeElevatedDarwin(ctx context.Context) error {
	// Get the command to execute
	cmdStr := e.getCommandString()
	if cmdStr == "" {
		return fmt.Errorf("cannot determine command to elevate")
	}

	logger.Debug("Requesting sudo privileges for: %s", cmdStr)

	// Use osascript to prompt for admin privileges
	script := fmt.Sprintf(`do shell script %q with administrator privileges`,
		escapeForAppleScript(cmdStr))

	cmd := exec.CommandContext(ctx, "osascript", "-e", script)
	return cmd.Run()
}

// executeElevatedLinux runs command with elevated privileges on Linux
func (e *ElevatedExecutor) executeElevatedLinux(ctx context.Context) error {
	// Get the command to execute
	cmdStr := e.getCommandString()
	if cmdStr == "" {
		return fmt.Errorf("cannot determine command to elevate")
	}

	logger.Debug("Requesting elevated privileges for: %s", cmdStr)

	// Try different elevation methods
	elevationTools := []struct {
		name string
		args func(string) []string
	}{
		{
			name: "pkexec",
			args: func(cmd string) []string {
				return []string{"sh", "-c", cmd}
			},
		},
		{
			name: "gksudo",
			args: func(cmd string) []string {
				return []string{"--", "sh", "-c", cmd}
			},
		},
		{
			name: "kdesudo",
			args: func(cmd string) []string {
				return []string{"--", "sh", "-c", cmd}
			},
		},
		{
			name: "sudo",
			args: func(cmd string) []string {
				// Use graphical sudo if available
				return []string{"-A", "sh", "-c", cmd}
			},
		},
	}

	for _, tool := range elevationTools {
		if _, err := exec.LookPath(tool.name); err != nil {
			continue
		}
		args := tool.args(cmdStr)
		//nolint:gosec // tool.name is from predefined trusted list
		cmd := exec.CommandContext(ctx, tool.name, args...)

		// Set SUDO_ASKPASS for graphical password prompt
		if tool.name == "sudo" {
			cmd.Env = append(os.Environ(), "SUDO_ASKPASS=/usr/bin/ssh-askpass")
		}

		if err := cmd.Run(); err == nil {
			return nil
		}
		logger.Debug("Elevation with %s failed, trying next method", tool.name)
	}

	return fmt.Errorf("no suitable elevation tool found")
}

// getCommandString extracts the command string from the base executor
func (e *ElevatedExecutor) getCommandString() string {
	// Since we can't access the internal structure of executors,
	// we'll need to use the String() method which should provide
	// a meaningful representation
	return e.baseExecutor.String()
}

// Helper function to escape strings for PowerShell
func escapeForPowerShell(s string) string {
	s = strings.ReplaceAll(s, "'", "''")
	s = strings.ReplaceAll(s, "\n", "`n")
	s = strings.ReplaceAll(s, "\r", "`r")
	s = strings.ReplaceAll(s, "\"", "`\"")
	return s
}

// Helper function to escape strings for AppleScript
func escapeForAppleScript(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	return s
}
