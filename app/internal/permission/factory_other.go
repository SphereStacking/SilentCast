//go:build !darwin && !windows
// +build !darwin,!windows

package permission

import "fmt"

// NewManager creates a new permission manager for unsupported platforms
func NewManager() (Manager, error) {
	return nil, fmt.Errorf("unsupported operating system")
}