//go:build !linux

package terminal

// NewLinuxManager creates a new Linux terminal manager (stub for non-Linux)
func NewLinuxManager() Manager {
	detector := NewLinuxDetector()
	builder := NewLinuxCommandBuilder()
	return newBaseManager(detector, builder)
}
