//go:build !darwin && !windows && !linux
// +build !darwin,!windows,!linux

package permission

// NewManager creates a new permission manager for unsupported platforms
func NewManager() (Manager, error) {
	return newStubManager(), nil
}
