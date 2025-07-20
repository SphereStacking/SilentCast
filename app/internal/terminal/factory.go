package terminal

import "runtime"

// NewManager creates a platform-specific terminal manager
func NewManager() Manager {
	switch runtime.GOOS {
	case "windows":
		return NewWindowsManager()
	case "darwin":
		return NewDarwinManager()
	case "linux":
		return NewLinuxManager()
	default:
		// Return a manager with empty detector for unsupported platforms
		return &baseManager{
			detector: &baseDetector{terminals: []Terminal{}},
			builder:  &baseCommandBuilder{},
		}
	}
}
