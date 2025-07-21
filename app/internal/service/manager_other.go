//go:build !windows && !darwin && !linux

package service

import (
	"fmt"
	"runtime"
)

// StubManager implements service management for non-Windows platforms
type StubManager struct {
	onRun func() error
}

// NewManager creates a new service manager (stub for non-Windows)
func NewManager(onRun func() error) Manager {
	return &StubManager{
		onRun: onRun,
	}
}

// Install is not supported on this platform
func (m *StubManager) Install() error {
	return fmt.Errorf("service installation is not supported on %s", runtime.GOOS)
}

// Uninstall is not supported on this platform
func (m *StubManager) Uninstall() error {
	return fmt.Errorf("service uninstallation is not supported on %s", runtime.GOOS)
}

// Start is not supported on this platform
func (m *StubManager) Start() error {
	return fmt.Errorf("service start is not supported on %s", runtime.GOOS)
}

// Stop is not supported on this platform
func (m *StubManager) Stop() error {
	return fmt.Errorf("service stop is not supported on %s", runtime.GOOS)
}

// Status is not supported on this platform
func (m *StubManager) Status() (ServiceStatus, error) {
	return ServiceStatus{
		Message: fmt.Sprintf("Service management is not supported on %s", runtime.GOOS),
	}, fmt.Errorf("service status is not supported on %s", runtime.GOOS)
}

// Run returns an error on non-Windows platforms to indicate normal execution should continue
func (m *StubManager) Run() error {
	// Return error to indicate we're not running as a service
	// The main function will continue with normal execution
	return fmt.Errorf("service mode not supported on %s", runtime.GOOS)
}
