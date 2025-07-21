//go:build linux

package service

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestLinuxManager_Install tests service installation
func TestLinuxManager_Install(t *testing.T) {
	// Skip if not running as user (can't install services in CI)
	if os.Getuid() == 0 {
		t.Skip("Skipping user service test when running as root")
	}

	// Create test manager
	mgr := &LinuxManager{
		executable:   "/usr/local/bin/silentcast",
		onRun:        func() error { return nil },
		useSystemd:   false, // Don't use systemd in tests
		isSystemWide: false,
	}

	// Test XDG autostart creation
	t.Run("XDG Autostart", func(t *testing.T) {
		// Create temporary home directory for testing
		tmpHome := t.TempDir()
		
		// Save and restore HOME
		originalHome := os.Getenv("HOME")
		defer func() {
			os.Setenv("HOME", originalHome)
		}()
		os.Setenv("HOME", tmpHome)

		// Test installation
		err := mgr.installXDGAutostart()
		if err != nil {
			t.Fatalf("Failed to install XDG autostart: %v", err)
		}

		// Check if desktop file was created
		desktopPath := filepath.Join(tmpHome, ".config/autostart/silentcast.desktop")
		if _, statErr := os.Stat(desktopPath); os.IsNotExist(statErr) {
			t.Error("Desktop file was not created")
		}

		// Verify desktop file contents
		content, err := os.ReadFile(desktopPath)
		if err != nil {
			t.Fatalf("Failed to read desktop file: %v", err)
		}

		expectedContents := []string{
			"[Desktop Entry]",
			"Type=Application",
			"Name=SilentCast",
			"Exec=/usr/local/bin/silentcast --no-tray",
			"Terminal=false",
			"X-GNOME-Autostart-enabled=true",
		}

		for _, expected := range expectedContents {
			if !strings.Contains(string(content), expected) {
				t.Errorf("Desktop file missing expected content: %s", expected)
			}
		}

		// Test duplicate installation
		err = mgr.installXDGAutostart()
		if err == nil {
			t.Error("Expected error when installing duplicate XDG autostart")
		}
	})
}

// TestLinuxManager_Uninstall tests service uninstallation
func TestLinuxManager_Uninstall(t *testing.T) {
	// Skip if not running as user
	if os.Getuid() == 0 {
		t.Skip("Skipping user service test when running as root")
	}

	// Create test manager
	mgr := &LinuxManager{
		executable:   "/usr/local/bin/silentcast",
		onRun:        func() error { return nil },
		useSystemd:   false,
		isSystemWide: false,
	}

	t.Run("XDG Autostart Removal", func(t *testing.T) {
		// Create temporary home directory
		tmpHome := t.TempDir()
		
		// Save and restore HOME
		originalHome := os.Getenv("HOME")
		defer func() {
			os.Setenv("HOME", originalHome)
		}()
		os.Setenv("HOME", tmpHome)

		// Install first
		err := mgr.installXDGAutostart()
		if err != nil {
			t.Fatalf("Failed to install for uninstall test: %v", err)
		}

		// Test uninstallation
		err = mgr.uninstallXDGAutostart()
		if err != nil {
			t.Fatalf("Failed to uninstall XDG autostart: %v", err)
		}

		// Check if desktop file was removed
		desktopPath := filepath.Join(tmpHome, ".config/autostart/silentcast.desktop")
		if _, statErr := os.Stat(desktopPath); statErr == nil {
			t.Error("Desktop file was not removed")
		}

		// Test uninstalling non-existent
		err = mgr.uninstallXDGAutostart()
		if err == nil {
			t.Error("Expected error when uninstalling non-existent XDG autostart")
		}
	})
}

// TestSystemdTemplate tests systemd service file generation
func TestSystemdTemplate(t *testing.T) {
	// Test service file would contain expected content
	expectedContents := []string{
		"Description=Silent hotkey-driven task runner",
		"After=graphical-session.target",
		"Type=simple",
		"ExecStart=/usr/local/bin/silentcast --no-tray",
		"Restart=on-failure",
		"WantedBy=default.target",
	}

	// This is a template validation test
	for _, expected := range expectedContents {
		if !strings.Contains(systemdTemplate, strings.Split(expected, "=")[0]) {
			t.Errorf("Systemd template missing expected directive: %s", expected)
		}
	}
}

// TestDesktopTemplate tests XDG desktop file generation
func TestDesktopTemplate(t *testing.T) {
	// Test desktop template contains required fields
	expectedFields := []string{
		"[Desktop Entry]",
		"Type=",
		"Name=",
		"Comment=",
		"Exec=",
		"Terminal=",
		"Categories=",
		"X-GNOME-Autostart-enabled=",
	}

	for _, field := range expectedFields {
		if !strings.Contains(desktopTemplate, field) {
			t.Errorf("Desktop template missing required field: %s", field)
		}
	}
}

// TestLinuxManager_Status tests service status checking
func TestLinuxManager_Status(t *testing.T) {
	// Create temporary home directory
	tmpHome := t.TempDir()
	
	// Save and restore HOME
	originalHome := os.Getenv("HOME")
	defer func() {
		os.Setenv("HOME", originalHome)
	}()
	os.Setenv("HOME", tmpHome)
	
	// Create test manager
	mgr := &LinuxManager{
		executable:   "/usr/local/bin/silentcast",
		onRun:        func() error { return nil },
		useSystemd:   false,
		isSystemWide: false,
	}

	// Test status when not installed
	status, err := mgr.Status()
	if err != nil {
		t.Fatalf("Failed to get status: %v", err)
	}

	if status.Installed {
		t.Error("Expected service to not be installed")
	}

	if status.Running {
		t.Error("Expected service to not be running")
	}

	if status.Message != "Service not installed" {
		t.Errorf("Unexpected status message: %s", status.Message)
	}
}

// TestHasSystemd tests systemd detection
func TestHasSystemd(t *testing.T) {
	// This test just ensures the function doesn't panic
	_ = hasSystemd()
}

// TestLinuxManager_Run tests the Run method
func TestLinuxManager_Run(t *testing.T) {
	runCalled := false
	mgr := &LinuxManager{
		executable: "/usr/local/bin/silentcast",
		onRun: func() error {
			runCalled = true
			return nil
		},
	}

	err := mgr.Run()
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if !runCalled {
		t.Error("onRun callback was not called")
	}
}

// TestLinuxManager_StartStopWithoutSystemd tests Start/Stop when systemd is not available
func TestLinuxManager_StartStopWithoutSystemd(t *testing.T) {
	mgr := &LinuxManager{
		executable:   "/usr/local/bin/silentcast",
		onRun:        func() error { return nil },
		useSystemd:   false,
		isSystemWide: false,
	}

	// Test Start without systemd
	err := mgr.Start()
	if err == nil {
		t.Error("Expected error when starting without systemd")
	}
	if !strings.Contains(err.Error(), "systemd not available") {
		t.Errorf("Expected systemd error, got: %v", err)
	}

	// Test Stop without systemd
	err = mgr.Stop()
	if err == nil {
		t.Error("Expected error when stopping without systemd")
	}
	if !strings.Contains(err.Error(), "systemd not available") {
		t.Errorf("Expected systemd error, got: %v", err)
	}
}

// TestLinuxManager_StatusWithXDGOnly tests status when only XDG autostart is installed
func TestLinuxManager_StatusWithXDGOnly(t *testing.T) {
	// Create temporary home directory
	tmpHome := t.TempDir()
	
	// Save and restore HOME
	originalHome := os.Getenv("HOME")
	defer func() {
		os.Setenv("HOME", originalHome)
	}()
	os.Setenv("HOME", tmpHome)

	mgr := &LinuxManager{
		executable:   "/usr/local/bin/silentcast",
		onRun:        func() error { return nil },
		useSystemd:   false,
		isSystemWide: false,
	}

	// Install XDG autostart
	err := mgr.installXDGAutostart()
	if err != nil {
		t.Fatalf("Failed to install XDG autostart: %v", err)
	}

	// Check status
	status, err := mgr.Status()
	if err != nil {
		t.Fatalf("Failed to get status: %v", err)
	}

	if !status.Installed {
		t.Error("Expected service to be installed via XDG autostart")
	}

	if !strings.Contains(status.Message, "XDG autostart installed (no systemd)") {
		t.Errorf("Expected XDG autostart message, got: %s", status.Message)
	}
}

// TestLinuxManager_NewManagerExecutableResolution tests executable path resolution
func TestLinuxManager_NewManagerExecutableResolution(t *testing.T) {
	mgr := NewManager(func() error { return nil }).(*LinuxManager)

	// Should have resolved an executable path
	if mgr.executable == "" {
		t.Error("NewManager should resolve an executable path")
	}

	// Should be using systemd detection
	expectedSystemd := hasSystemd()
	if mgr.useSystemd != expectedSystemd {
		t.Errorf("Expected useSystemd=%v, got %v", expectedSystemd, mgr.useSystemd)
	}

	// Should default to user installation
	if mgr.isSystemWide {
		t.Error("Should default to user installation")
	}
}

// TestLinuxManager_HomeDirectoryFallback tests fallback when HOME is not set
func TestLinuxManager_HomeDirectoryFallback(t *testing.T) {
	// Save and restore HOME
	originalHome := os.Getenv("HOME")
	defer func() {
		os.Setenv("HOME", originalHome)
	}()
	
	// Unset HOME environment variable
	os.Unsetenv("HOME")

	mgr := &LinuxManager{
		executable:   "/usr/local/bin/silentcast",
		onRun:        func() error { return nil },
		useSystemd:   false,
		isSystemWide: false,
	}

	// Should not panic and should handle gracefully
	status, err := mgr.Status()
	if err != nil {
		t.Fatalf("Status should not fail when HOME is unset: %v", err)
	}

	// The status depends on whether user.Current() succeeds
	// When HOME is unset, user.Current() might still work and return a home directory
	// This test ensures the function doesn't panic, behavior may vary by system
	_ = status // Don't assume specific behavior when HOME is unset
}

// TestLinuxManager_InstallUninstallErrors tests error conditions
func TestLinuxManager_InstallUninstallErrors(t *testing.T) {
	// Create temporary home directory
	tmpHome := t.TempDir()
	
	// Save and restore HOME
	originalHome := os.Getenv("HOME")
	defer func() {
		os.Setenv("HOME", originalHome)
	}()
	os.Setenv("HOME", tmpHome)

	mgr := &LinuxManager{
		executable:   "/usr/local/bin/silentcast",
		onRun:        func() error { return nil },
		useSystemd:   false,
		isSystemWide: false,
	}

	// Test full install/uninstall cycle
	err := mgr.Install()
	if err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	// Try to install again (should fail)
	err = mgr.Install()
	if err == nil {
		t.Error("Install should fail when already installed")
	}

	// Uninstall
	err = mgr.Uninstall()
	if err != nil {
		t.Fatalf("Uninstall failed: %v", err)
	}

	// Try to uninstall again (should fail)
	err = mgr.Uninstall()
	if err == nil {
		t.Error("Uninstall should fail when not installed")
	}
}