//go:build darwin
// +build darwin

package permission

import (
	"context"
	"fmt"
	"os/exec"
)

// macManager implements the Manager interface for macOS
type macManager struct{}

// newMacManager creates a new macOS permission manager
func newMacManager() Manager {
	return &macManager{}
}

// NewManager creates a new permission manager for macOS
func NewManager() (Manager, error) {
	return newMacManager(), nil
}

// Check returns the current status of all permissions
func (m *macManager) Check(ctx context.Context) ([]Permission, error) {
	permissions := []Permission{
		{
			Type:        PermissionTypeAccessibility,
			Status:      m.checkAccessibility(),
			Description: "Required for detecting global hotkeys",
			Required:    true,
		},
		{
			Type:        PermissionTypeNotification,
			Status:      m.checkNotifications(),
			Description: "Required for showing system notifications",
			Required:    false,
		},
		{
			Type:        PermissionTypeAutoStart,
			Status:      m.checkAutoStart(),
			Description: "Required for starting with system",
			Required:    false,
		},
	}

	return permissions, nil
}

// Request attempts to request the specified permission
func (m *macManager) Request(ctx context.Context, permType PermissionType) error {
	switch permType {
	case PermissionTypeAccessibility:
		// On macOS, we can't directly request accessibility permission
		// We can only guide the user to enable it
		return fmt.Errorf("accessibility permission must be granted manually in System Preferences")
	case PermissionTypeNotification:
		// Notification permission is usually requested when first notification is sent
		return nil
	case PermissionTypeAutoStart:
		// This is handled by creating a LaunchAgent
		return m.createLaunchAgent()
	default:
		return fmt.Errorf("unknown permission type: %s", permType)
	}
}

// GetInstructions returns user-friendly instructions for granting a permission
func (m *macManager) GetInstructions(permType PermissionType) string {
	switch permType {
	case PermissionTypeAccessibility:
		return `To grant accessibility permission:
1. Open System Preferences
2. Go to Security & Privacy
3. Click the Privacy tab
4. Select Accessibility from the left sidebar
5. Click the lock icon and enter your password
6. Check the box next to Spellbook
7. Restart Spellbook`

	case PermissionTypeNotification:
		return `To enable notifications:
1. Open System Preferences
2. Go to Notifications & Focus
3. Find Spellbook in the list
4. Turn on "Allow Notifications"`

	case PermissionTypeAutoStart:
		return `To enable auto-start:
1. Open System Preferences
2. Go to Users & Groups
3. Select your user account
4. Click the Login Items tab
5. Click the + button
6. Select Spellbook from Applications
7. Click Add`

	default:
		return "No instructions available for this permission type"
	}
}

// OpenSettings opens the system settings for the specified permission
func (m *macManager) OpenSettings(permType PermissionType) error {
	var url string

	switch permType {
	case PermissionTypeAccessibility:
		url = "x-apple.systempreferences:com.apple.preference.security?Privacy_Accessibility"
	case PermissionTypeNotification:
		url = "x-apple.systempreferences:com.apple.preference.notifications"
	case PermissionTypeAutoStart:
		url = "x-apple.systempreferences:com.apple.preference.users"
	default:
		return fmt.Errorf("no settings URL for permission type: %s", permType)
	}

	cmd := exec.Command("open", url)
	return cmd.Run()
}

// IsSupported checks if a permission type is supported on macOS
func (m *macManager) IsSupported(permType PermissionType) bool {
	switch permType {
	case PermissionTypeAccessibility, PermissionTypeNotification, PermissionTypeAutoStart:
		return true
	default:
		return false
	}
}

// checkAccessibility checks if accessibility permission is granted
func (m *macManager) checkAccessibility() Status {
	// This would need to use CGO or a helper binary to check actual status
	// For now, return NotDetermined
	// In production, this would check:
	// AXIsProcessTrustedWithOptions or similar API
	return StatusNotDetermined
}

// checkNotifications checks if notification permission is granted
func (m *macManager) checkNotifications() Status {
	// This would check notification center permissions
	// For now, return NotDetermined
	return StatusNotDetermined
}

// checkAutoStart checks if auto-start is enabled
func (m *macManager) checkAutoStart() Status {
	// This would check if LaunchAgent exists
	// For now, return NotDetermined
	return StatusNotDetermined
}

// createLaunchAgent creates a LaunchAgent for auto-start
func (m *macManager) createLaunchAgent() error {
	// This would create a plist file in ~/Library/LaunchAgents
	// For now, return not implemented
	return fmt.Errorf("auto-start setup not implemented yet")
}
