//go:build linux

package terminal

// LinuxManager implements the Manager interface for Linux
type LinuxManager struct {
	*baseManager
}

// NewLinuxManager creates a new Linux terminal manager
func NewLinuxManager() Manager {
	detector := NewLinuxDetector()
	builder := NewLinuxCommandBuilder()
	return &LinuxManager{
		baseManager: newBaseManager(detector, builder),
	}
}

// Linux uses the base implementation for all terminals
