//go:build !darwin

package terminal

// NewDarwinManager creates a new macOS terminal manager (stub for non-macOS)
func NewDarwinManager() Manager {
	detector := NewMacOSDetector()
	builder := NewMacOSCommandBuilder()
	return newBaseManager(detector, builder)
}
