package permission

import "context"

// PermissionType represents different types of permissions required
type PermissionType string

const (
	// PermissionTypeAccessibility is required for global hotkey detection
	PermissionTypeAccessibility PermissionType = "accessibility"
	// PermissionTypeNotification is required for system notifications
	PermissionTypeNotification PermissionType = "notification"
	// PermissionTypeAutoStart is required for system startup
	PermissionTypeAutoStart PermissionType = "autostart"
)

// Status represents the current status of a permission
type Status string

const (
	// StatusGranted means the permission is granted
	StatusGranted Status = "granted"
	// StatusDenied means the permission is denied
	StatusDenied Status = "denied"
	// StatusNotDetermined means the permission hasn't been requested yet
	StatusNotDetermined Status = "not_determined"
	// StatusNotApplicable means the permission is not applicable on this OS
	StatusNotApplicable Status = "not_applicable"
)

// Permission represents a single permission requirement
type Permission struct {
	Type        PermissionType
	Status      Status
	Description string
	Required    bool
}

// Manager defines the interface for OS-specific permission management
type Manager interface {
	// Check returns the current status of all permissions
	Check(ctx context.Context) ([]Permission, error)

	// Request attempts to request the specified permission
	Request(ctx context.Context, permType PermissionType) error

	// GetInstructions returns user-friendly instructions for granting a permission
	GetInstructions(permType PermissionType) string

	// OpenSettings opens the system settings for the specified permission
	OpenSettings(permType PermissionType) error

	// IsSupported checks if a permission type is supported on this OS
	IsSupported(permType PermissionType) bool
}
