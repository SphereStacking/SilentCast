//go:build darwin

package service

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"text/template"
)

func TestDarwinManager_GetPlistPath(t *testing.T) {
	mgr := &DarwinManager{
		executable:    "/usr/local/bin/silentcast",
		isSystemLevel: false,
	}

	path, err := mgr.getPlistPath()
	if err != nil {
		t.Fatalf("getPlistPath() error = %v", err)
	}

	// Should be in user's LaunchAgents
	if !strings.Contains(path, "Library/LaunchAgents") {
		t.Errorf("getPlistPath() = %v, want path containing Library/LaunchAgents", path)
	}

	if !strings.HasSuffix(path, ".plist") {
		t.Errorf("getPlistPath() = %v, want path ending with .plist", path)
	}

	expectedName := serviceName + ".plist"
	if !strings.Contains(path, expectedName) {
		t.Errorf("getPlistPath() = %v, want path containing %s", path, expectedName)
	}
}

func TestDarwinManager_Status_NotInstalled(t *testing.T) {
	// Create a manager with a non-existent plist path
	mgr := &DarwinManager{
		executable:    "/tmp/test/silentcast",
		isSystemLevel: false,
	}

	status, err := mgr.Status()
	if err != nil {
		t.Fatalf("Status() error = %v", err)
	}

	if status.Installed {
		t.Error("Status().Installed = true, want false for non-existent service")
	}

	if status.Running {
		t.Error("Status().Running = true, want false for non-existent service")
	}

	if status.Message != "Service not installed" {
		t.Errorf("Status().Message = %v, want 'Service not installed'", status.Message)
	}
}

func TestDarwinManager_SystemLevelPath(t *testing.T) {
	mgr := &DarwinManager{
		executable:    "/usr/local/bin/silentcast",
		isSystemLevel: true,
	}

	path, err := mgr.getPlistPath()
	if err != nil {
		t.Fatalf("getPlistPath() error = %v", err)
	}

	// Should be in system LaunchAgents
	if !strings.Contains(path, "/Library/LaunchAgents") {
		t.Errorf("System-level getPlistPath() = %v, want path containing /Library/LaunchAgents", path)
	}
}

func TestDarwinManager_Run(t *testing.T) {
	called := false
	onRun := func() error {
		called = true
		return nil
	}

	mgr := &DarwinManager{
		executable: "/usr/local/bin/silentcast",
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

func TestDarwinManager_RunError(t *testing.T) {
	expectedErr := "test error"
	onRun := func() error {
		return fmt.Errorf(expectedErr)
	}

	mgr := &DarwinManager{
		executable: "/usr/local/bin/silentcast",
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

func TestNewManager_Darwin(t *testing.T) {
	onRun := func() error { return nil }
	manager := NewManager(onRun)

	// Should return DarwinManager
	darwinMgr, ok := manager.(*DarwinManager)
	if !ok {
		t.Errorf("NewManager() should return *DarwinManager, got %T", manager)
	}

	if darwinMgr.executable == "" {
		t.Error("NewManager() should set executable")
	}

	if darwinMgr.isSystemLevel {
		t.Error("NewManager() should default to user-level service")
	}
}

func TestDarwinManager_InstallUninstallMethods(t *testing.T) {
	mgr := &DarwinManager{
		executable:    "/nonexistent/silentcast",
		onRun:         func() error { return nil },
		isSystemLevel: false,
	}

	// Test Install - should fail because directories don't exist and launchctl unavailable
	err := mgr.Install()
	// Don't assert on error as it depends on environment, just ensure it doesn't panic
	_ = err

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

	// Test Stop - might succeed or fail depending on launchctl
	err = mgr.Stop()
	_ = err // Don't assert as it's environment dependent
}

func TestPlistTemplate(t *testing.T) {
	// Test that the plist template is valid
	data := struct {
		Label            string
		Executable       string
		RunAtLoad        string
		LogPath          string
		WorkingDirectory string
		IsSystemLevel    bool
		UserName         string
	}{
		Label:            "com.test.service",
		Executable:       "/usr/local/bin/test",
		RunAtLoad:        "true",
		LogPath:          "/tmp",
		WorkingDirectory: "/tmp",
		IsSystemLevel:    false,
		UserName:         "testuser",
	}

	// Just verify template can be parsed and executed
	tmpl, err := template.New("plist").Parse(plistTemplate)
	if err != nil {
		t.Fatalf("Failed to parse plist template: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		t.Fatalf("Failed to execute plist template: %v", err)
	}

	output := buf.String()

	// Check for required elements
	requiredElements := []string{
		"<?xml version",
		"<!DOCTYPE plist",
		"<plist version=\"1.0\">",
		"<key>Label</key>",
		"<string>com.test.service</string>",
		"<key>ProgramArguments</key>",
		"<key>RunAtLoad</key>",
		"<key>KeepAlive</key>",
		"</plist>",
	}

	for _, elem := range requiredElements {
		if !strings.Contains(output, elem) {
			t.Errorf("Plist output missing required element: %s", elem)
		}
	}
}
