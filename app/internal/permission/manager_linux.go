//go:build linux

package permission

import "context"

// NewManager creates a new permission manager for Linux
func NewManager() (Manager, error) {
	return &linuxManager{}, nil
}

type linuxManager struct{}

func (m *linuxManager) Check(ctx context.Context) ([]Permission, error) {
	// Linux doesn't require special permissions for global hotkeys
	// when running under X11/Wayland with proper user permissions
	return []Permission{
		{
			Type:        PermissionTypeAccessibility,
			Status:      StatusGranted,
			Description: "Global hotkey access",
			Required:    true,
		},
		{
			Type:        PermissionTypeNotification,
			Status:      StatusGranted,
			Description: "System notifications",
			Required:    false,
		},
		{
			Type:        PermissionTypeAutoStart,
			Status:      StatusGranted,
			Description: "Launch at startup",
			Required:    false,
		},
	}, nil
}

func (m *linuxManager) Request(ctx context.Context, permType PermissionType) error {
	// No special permission request needed on Linux
	return nil
}

func (m *linuxManager) GetInstructions(permType PermissionType) string {
	switch permType {
	case PermissionTypeAccessibility:
		return "No special permissions needed on Linux for global hotkeys."
	default:
		return "This permission type is not applicable to Linux."
	}
}

func (m *linuxManager) OpenSettings(permType PermissionType) error {
	// No settings to open on Linux
	return nil
}

func (m *linuxManager) IsSupported(permType PermissionType) bool {
	// Linux supports all permission types
	switch permType {
	case PermissionTypeAccessibility, PermissionTypeNotification, PermissionTypeAutoStart:
		return true
	default:
		return false
	}
}
