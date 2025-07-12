//go:build !darwin && !windows
// +build !darwin,!windows

package permission

import (
	"context"
	"fmt"
)

// stubManager implements the Manager interface for unsupported platforms
//
//nolint:unused // Used in build tag for unsupported platforms
type stubManager struct{}

// newStubManager creates a new stub permission manager
//
//nolint:unused // Used in build tag for unsupported platforms
func newStubManager() Manager {
	return &stubManager{}
}

// Check returns the current status of all permissions
//
//nolint:unused // Used in build tag for unsupported platforms
func (m *stubManager) Check(ctx context.Context) ([]Permission, error) {
	return []Permission{
		{
			Type:        PermissionTypeAccessibility,
			Status:      StatusNotApplicable,
			Description: "Not supported on this platform",
			Required:    false,
		},
	}, nil
}

// Request attempts to request the specified permission
//
//nolint:unused // Used in build tag for unsupported platforms
func (m *stubManager) Request(ctx context.Context, permType PermissionType) error {
	return fmt.Errorf("permission management not supported on this platform")
}

// GetInstructions returns user-friendly instructions for granting a permission
//
//nolint:unused // Used in build tag for unsupported platforms
func (m *stubManager) GetInstructions(permType PermissionType) string {
	return "Permission management is not supported on this platform"
}

// OpenSettings opens the system settings for the specified permission
//
//nolint:unused // Used in build tag for unsupported platforms
func (m *stubManager) OpenSettings(permType PermissionType) error {
	return fmt.Errorf("permission settings not available on this platform")
}

// IsSupported checks if a permission type is supported
//
//nolint:unused // Used in build tag for unsupported platforms
func (m *stubManager) IsSupported(permType PermissionType) bool {
	return false
}
