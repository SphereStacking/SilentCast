//go:build !windows

package terminal

// NewWindowsManager creates a new Windows terminal manager (stub for non-Windows)
func NewWindowsManager() Manager {
	detector := NewWindowsDetector()
	builder := NewWindowsCommandBuilder()
	return newBaseManager(detector, builder)
}
