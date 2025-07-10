//go:build nogohook
// +build nogohook

package hotkey

import (
	"github.com/SphereStacking/silentcast/internal/config"
)

// NewManager creates a new hotkey manager (stub version)
func NewManager(cfg *config.HotkeyConfig) (*MockManager, error) {
	// Return a mock manager when gohook is not available
	return NewMockManager(), nil
}