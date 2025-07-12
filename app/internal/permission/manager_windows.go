//go:build windows
// +build windows

package permission

import (
	"context"
	"fmt"
	"os/exec"
	"syscall"
)

// windowsManager implements the Manager interface for Windows
type windowsManager struct{}

// newWindowsManager creates a new Windows permission manager
func newWindowsManager() Manager {
	return &windowsManager{}
}

// NewManager creates a new permission manager for Windows
func NewManager() (Manager, error) {
	return newWindowsManager(), nil
}

// Check returns the current status of all permissions
func (m *windowsManager) Check(ctx context.Context) ([]Permission, error) {
	permissions := []Permission{
		{
			Type:        PermissionTypeAccessibility,
			Status:      m.checkAccessibility(),
			Description: "Administrative privileges may be required for global hotkeys",
			Required:    false, // Not always required on Windows
		},
		{
			Type:        PermissionTypeNotification,
			Status:      StatusGranted, // Windows 10+ allows notifications by default
			Description: "Required for showing system notifications",
			Required:    false,
		},
		{
			Type:        PermissionTypeAutoStart,
			Status:      m.checkAutoStart(),
			Description: "Required for starting with Windows",
			Required:    false,
		},
	}

	return permissions, nil
}

// Request attempts to request the specified permission
func (m *windowsManager) Request(ctx context.Context, permType PermissionType) error {
	switch permType {
	case PermissionTypeAccessibility:
		// On Windows, we might need to restart with admin privileges
		return m.requestAdminPrivileges()
	case PermissionTypeNotification:
		// Notifications work by default on Windows 10+
		return nil
	case PermissionTypeAutoStart:
		// This is handled by creating a registry entry or scheduled task
		return m.setupAutoStart()
	default:
		return fmt.Errorf("unknown permission type: %s", permType)
	}
}

// GetInstructions returns user-friendly instructions for granting a permission
func (m *windowsManager) GetInstructions(permType PermissionType) string {
	switch permType {
	case PermissionTypeAccessibility:
		return `To run with administrative privileges:
1. Right-click on Spellbook
2. Select "Run as administrator"
3. Click "Yes" in the UAC prompt

Or to disable UAC prompts for Spellbook:
1. Open Task Scheduler
2. Create a new task with highest privileges
3. Set it to run Spellbook at login`

	case PermissionTypeNotification:
		return `Notifications are enabled by default on Windows 10 and later.
If you don't see notifications:
1. Open Settings (Win+I)
2. Go to System > Notifications & actions
3. Make sure notifications are turned on
4. Find Spellbook in the app list and enable it`

	case PermissionTypeAutoStart:
		return `To enable auto-start:
1. Press Win+R and type: shell:startup
2. Copy Spellbook shortcut to the Startup folder

Or use Task Manager:
1. Open Task Manager (Ctrl+Shift+Esc)
2. Go to the Startup tab
3. If Spellbook is listed, right-click and Enable`

	default:
		return "No instructions available for this permission type"
	}
}

// OpenSettings opens the system settings for the specified permission
func (m *windowsManager) OpenSettings(permType PermissionType) error {
	var args []string

	switch permType {
	case PermissionTypeAccessibility:
		// Open UAC settings
		args = []string{"UserAccountControlSettings.exe"}
	case PermissionTypeNotification:
		// Open notification settings
		args = []string{"start", "ms-settings:notifications"}
	case PermissionTypeAutoStart:
		// Open startup folder
		args = []string{"explorer.exe", "shell:startup"}
	default:
		return fmt.Errorf("no settings URL for permission type: %s", permType)
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}

// IsSupported checks if a permission type is supported on Windows
func (m *windowsManager) IsSupported(permType PermissionType) bool {
	switch permType {
	case PermissionTypeAccessibility, PermissionTypeNotification, PermissionTypeAutoStart:
		return true
	default:
		return false
	}
}

// checkAccessibility checks if we have admin privileges
func (m *windowsManager) checkAccessibility() Status {
	// Check if running as administrator
	// This is a simplified check - in production would use Windows API
	return StatusNotDetermined
}

// checkAutoStart checks if auto-start is enabled
func (m *windowsManager) checkAutoStart() Status {
	// Would check registry or Task Scheduler
	return StatusNotDetermined
}

// requestAdminPrivileges attempts to restart with admin privileges
func (m *windowsManager) requestAdminPrivileges() error {
	// This would use ShellExecute with "runas" verb
	return fmt.Errorf("admin privilege elevation not implemented yet")
}

// setupAutoStart sets up auto-start via registry or scheduled task
func (m *windowsManager) setupAutoStart() error {
	// This would create registry entry in:
	// HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Run
	return fmt.Errorf("auto-start setup not implemented yet")
}
