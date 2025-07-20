//go:build windows && !darwin && !linux

package service

import (
	"fmt"
	"strings"
	"testing"
)

func TestWindowsManager_NewManager(t *testing.T) {
	onRun := func() error { return nil }
	manager := NewManager(onRun)

	// Should return WindowsManager
	winMgr, ok := manager.(*WindowsManager)
	if !ok {
		t.Errorf("NewManager() should return *WindowsManager, got %T", manager)
	}

	if winMgr.executable == "" {
		t.Error("NewManager() should set executable")
	}

	if winMgr.onRun == nil {
		t.Error("NewManager() should set onRun function")
	}
}

func TestWindowsManager_Run(t *testing.T) {
	called := false
	onRun := func() error {
		called = true
		return nil
	}

	mgr := &WindowsManager{
		executable: "C:\\test\\silentcast.exe",
		onRun:      onRun,
	}

	err := mgr.Run()
	if err != nil {
		t.Errorf("Run() error = %v", err)
	}
	if !called {
		t.Error("onRun function should be called")
	}
}

func TestWindowsManager_RunError(t *testing.T) {
	expectedErr := "test error"
	onRun := func() error {
		return fmt.Errorf(expectedErr)
	}

	mgr := &WindowsManager{
		executable: "C:\\test\\silentcast.exe",
		onRun:      onRun,
	}

	err := mgr.Run()
	if err == nil {
		t.Error("Run() should return error from onRun")
	}
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("Error should contain '%s', got: %v", expectedErr, err)
	}
}

func TestWindowsManager_ServiceMethods(t *testing.T) {
	mgr := &WindowsManager{
		executable: "C:\\nonexistent\\silentcast.exe",
		onRun:      func() error { return nil },
	}

	// Test Install - will fail due to permissions/environment
	err := mgr.Install()
	// Don't assert on specific error as it depends on environment
	_ = err

	// Test Status - should work even if service doesn't exist
	status, err := mgr.Status()
	if err != nil {
		// Status might fail due to permissions, which is expected in test environment
		t.Logf("Status() failed as expected in test environment: %v", err)
	} else {
		// If it succeeds, service should not be installed
		if status.Installed {
			t.Error("Status should show not installed for non-existent service")
		}
	}

	// Test Uninstall - should fail because service doesn't exist
	err = mgr.Uninstall()
	if err == nil {
		t.Error("Uninstall should fail when service not installed")
	}

	// Test Start - should fail because service doesn't exist
	err = mgr.Start()
	if err == nil {
		t.Error("Start should fail when service not installed")
	}

	// Test Stop - should fail because service doesn't exist
	err = mgr.Stop()
	if err == nil {
		t.Error("Stop should fail when service not installed")
	}
}

func TestWindowsManager_ServiceConstants(t *testing.T) {
	// Test that service constants are defined correctly
	if serviceName == "" {
		t.Error("serviceName should not be empty")
	}
	if serviceDisplayName == "" {
		t.Error("serviceDisplayName should not be empty")
	}
	if serviceDescription == "" {
		t.Error("serviceDescription should not be empty")
	}

	// Test expected values
	if serviceName != "SilentCast" {
		t.Errorf("serviceName = %s, want SilentCast", serviceName)
	}
	if !strings.Contains(serviceDisplayName, "SilentCast") {
		t.Errorf("serviceDisplayName should contain 'SilentCast', got: %s", serviceDisplayName)
	}
	if !strings.Contains(serviceDescription, "hotkey") {
		t.Errorf("serviceDescription should contain 'hotkey', got: %s", serviceDescription)
	}
}

func TestWindowsManager_ExecutablePath(t *testing.T) {
	mgr := &WindowsManager{
		executable: "C:\\Program Files\\SilentCast\\silentcast.exe",
		onRun:      func() error { return nil },
	}

	// Test that executable path is stored correctly
	if mgr.executable != "C:\\Program Files\\SilentCast\\silentcast.exe" {
		t.Errorf("executable = %s, want 'C:\\Program Files\\SilentCast\\silentcast.exe'", mgr.executable)
	}
}