package service

import (
	"runtime"
	"testing"
)

func TestServiceManager(t *testing.T) {
	onRunCalled := false
	onRun := func() error {
		onRunCalled = true
		return nil
	}

	mgr := NewManager(onRun)

	// Test that the manager is created
	if mgr == nil {
		t.Fatal("NewManager returned nil")
	}

	// Platform-specific tests
	switch runtime.GOOS {
	case "windows", "darwin", "linux":
		// These platforms have service management support
		// Basic smoke test - Status should work
		status, err := mgr.Status()
		if err != nil {
			t.Errorf("Status returned error on %s: %v", runtime.GOOS, err)
		}
		// Service likely not installed in test environment
		if status.Installed && runtime.GOOS != "linux" {
			t.Log("Warning: Service appears to be installed in test environment")
		}

		// Run should work on supported platforms
		if runtime.GOOS == "linux" {
			// Linux Run() just calls onRun
			if err := mgr.Run(); err != nil {
				t.Errorf("Run returned error on Linux: %v", err)
			}
			if !onRunCalled {
				t.Error("onRun should be called on Linux")
			}
		}
	default:
		// Unsupported platforms should return errors
		if err := mgr.Install(); err == nil {
			t.Errorf("Install should return error on %s", runtime.GOOS)
		}
		if err := mgr.Uninstall(); err == nil {
			t.Errorf("Uninstall should return error on %s", runtime.GOOS)
		}
		if err := mgr.Start(); err == nil {
			t.Errorf("Start should return error on %s", runtime.GOOS)
		}
		if err := mgr.Stop(); err == nil {
			t.Errorf("Stop should return error on %s", runtime.GOOS)
		}

		status, err := mgr.Status()
		if err == nil {
			t.Errorf("Status should return error on %s", runtime.GOOS)
		}
		if status.Installed {
			t.Errorf("Status.Installed should be false on %s", runtime.GOOS)
		}
	}
}

func TestServiceStatusDefaults(t *testing.T) {
	// Test default ServiceStatus
	status := ServiceStatus{}

	if status.Installed {
		t.Error("Default status should have Installed = false")
	}

	if status.Running {
		t.Error("Default status should have Running = false")
	}

	if status.StartType != "" {
		t.Error("Default status should have empty StartType")
	}

	if status.Message != "" {
		t.Error("Default status should have empty Message")
	}
}
