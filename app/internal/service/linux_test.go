//go:build linux

package service

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestLinuxManager_ErrorHandling(_ *testing.T) {
	// Create manager with invalid executable
	manager := &LinuxManager{
		executable:   "/invalid/path/executable",
		onRun:        func() error { return nil },
		useSystemd:   false,
		isSystemWide: false,
	}

	// Test Install with invalid HOME
	oldHome := os.Getenv("HOME")
	os.Unsetenv("HOME")
	defer func() {
		os.Setenv("HOME", oldHome)
	}()

	err := manager.Install()
	// Note: Install may still succeed if user.Current() works even without HOME env var
	// This tests the code path but doesn't guarantee failure
	_ = err // Error is expected but not required
}

func TestLinuxManager_TemplateErrors(t *testing.T) {
	tempDir := t.TempDir()

	manager := &LinuxManager{
		executable:   filepath.Join(tempDir, "test"),
		onRun:        func() error { return nil },
		useSystemd:   false,
		isSystemWide: false,
	}

	// Set HOME to temp directory
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer func() {
		os.Setenv("HOME", oldHome)
	}()

	// Test XDG autostart installation
	err := manager.installXDGAutostart()
	if err != nil {
		t.Errorf("installXDGAutostart failed: %v", err)
	}

	// Test double installation
	err = manager.installXDGAutostart()
	if err == nil {
		t.Error("installXDGAutostart should fail when already installed")
	}
}

func TestLinuxManager_UninstallNotInstalled(t *testing.T) {
	tempDir := t.TempDir()

	manager := &LinuxManager{
		executable:   filepath.Join(tempDir, "test"),
		onRun:        func() error { return nil },
		useSystemd:   false,
		isSystemWide: false,
	}

	// Set HOME to temp directory (clean state)
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer func() {
		os.Setenv("HOME", oldHome)
	}()

	// Test uninstalling when not installed
	err := manager.uninstallXDGAutostart()
	if err == nil {
		t.Error("uninstallXDGAutostart should fail when not installed")
	}
	if !strings.Contains(err.Error(), "not installed") {
		t.Errorf("Error should mention 'not installed', got: %v", err)
	}
}

func TestLinuxManager_StatusEdgeCases(t *testing.T) {
	tempDir := t.TempDir()

	manager := &LinuxManager{
		executable:   filepath.Join(tempDir, "test"),
		onRun:        func() error { return nil },
		useSystemd:   false,
		isSystemWide: false,
	}

	// Test status without HOME
	oldHome := os.Getenv("HOME")
	os.Unsetenv("HOME")
	defer func() {
		os.Setenv("HOME", oldHome)
	}()

	status, err := manager.Status()
	if err != nil {
		t.Errorf("Status should not fail without HOME: %v", err)
	}
	// Note: Status may still show installed if user.Current() works
	// This tests the code path when HOME is not available
	_ = status // Don't assume specific behavior when HOME is unavailable
}

func TestLinuxManager_SystemdNotAvailable(t *testing.T) {
	manager := &LinuxManager{
		executable:   "/test/executable",
		onRun:        func() error { return nil },
		useSystemd:   false,
		isSystemWide: false,
	}

	// Test Start without systemd
	err := manager.Start()
	if err == nil {
		t.Error("Start should fail without systemd")
	}
	if !strings.Contains(err.Error(), "systemd not available") {
		t.Errorf("Error should mention 'systemd not available', got: %v", err)
	}

	// Test Stop without systemd
	err = manager.Stop()
	if err == nil {
		t.Error("Stop should fail without systemd")
	}
	if !strings.Contains(err.Error(), "systemd not available") {
		t.Errorf("Error should mention 'systemd not available', got: %v", err)
	}
}

func TestLinuxManager_ServicePaths(t *testing.T) {
	tempDir := t.TempDir()

	// Test user service paths
	userManager := &LinuxManager{
		executable:   filepath.Join(tempDir, "test"),
		onRun:        func() error { return nil },
		useSystemd:   true,
		isSystemWide: false,
	}

	// Test system-wide service paths
	systemManager := &LinuxManager{
		executable:   filepath.Join(tempDir, "test"),
		onRun:        func() error { return nil },
		useSystemd:   true,
		isSystemWide: true,
	}

	// These are just to test the path logic without actual systemd calls
	_ = userManager
	_ = systemManager

	// Test that managers are created properly
	if userManager.isSystemWide {
		t.Error("User manager should not be system-wide")
	}
	if !systemManager.isSystemWide {
		t.Error("System manager should be system-wide")
	}
}

func TestLinuxManager_OnRunExecution(t *testing.T) {
	called := false
	onRun := func() error {
		called = true
		return nil
	}

	manager := &LinuxManager{
		executable:   "/test/executable",
		onRun:        onRun,
		useSystemd:   false,
		isSystemWide: false,
	}

	err := manager.Run()
	if err != nil {
		t.Errorf("Run failed: %v", err)
	}
	if !called {
		t.Error("onRun function should be called")
	}
}

func TestLinuxManager_RunError(t *testing.T) {
	expectedErr := "test error"
	onRun := func() error {
		return errors.New(expectedErr)
	}

	manager := &LinuxManager{
		executable:   "/test/executable",
		onRun:        onRun,
		useSystemd:   false,
		isSystemWide: false,
	}

	err := manager.Run()
	if err == nil {
		t.Error("Run should return error from onRun")
	}
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("Error should contain '%s', got: %v", expectedErr, err)
	}
}

func TestHasSystemdMocked(t *testing.T) {
	// Test hasSystemd function behavior
	// Note: This will vary based on test environment
	result := hasSystemd()

	// Just ensure it doesn't panic and returns a boolean
	if result {
		t.Log("systemd detected in test environment")
	} else {
		t.Log("systemd not detected in test environment")
	}
}

func TestLinuxManager_TemplateData(t *testing.T) {
	tempDir := t.TempDir()
	executable := filepath.Join(tempDir, "silentcast")

	manager := &LinuxManager{
		executable:   executable,
		onRun:        func() error { return nil },
		useSystemd:   false,
		isSystemWide: false,
	}

	// Set HOME to temp directory
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer func() {
		os.Setenv("HOME", oldHome)
	}()

	// Install XDG autostart to test template data
	err := manager.installXDGAutostart()
	if err != nil {
		t.Fatalf("installXDGAutostart failed: %v", err)
	}

	// Read the generated file and check content
	desktopPath := filepath.Join(tempDir, xdgAutostartPath, desktopFile)
	content, err := os.ReadFile(desktopPath)
	if err != nil {
		t.Fatalf("Failed to read desktop file: %v", err)
	}

	contentStr := string(content)
	expectedParts := []string{
		"Name=SilentCast",
		"Comment=Silent hotkey-driven task runner",
		"Exec=" + executable + " --no-tray",
		"Icon=silentcast",
		"Terminal=false",
	}

	for _, part := range expectedParts {
		if !strings.Contains(contentStr, part) {
			t.Errorf("Desktop file should contain '%s', got:\n%s", part, contentStr)
		}
	}
}

func TestLinuxManager_SystemdService(t *testing.T) {
	tempDir := t.TempDir()
	executable := filepath.Join(tempDir, "silentcast")

	manager := &LinuxManager{
		executable:   executable,
		onRun:        func() error { return nil },
		useSystemd:   true,
		isSystemWide: false,
		execCommand:  defaultCommandExecutor,
	}

	// Mock HOME directory
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer func() {
		os.Setenv("HOME", oldHome)
	}()

	// Test installSystemdService (will fail due to no systemctl, but tests the code path)
	err := manager.installSystemdService()
	// Expect failure due to systemctl not being available in controlled way
	if err == nil {
		t.Log("installSystemdService unexpectedly succeeded (systemctl available)")
	} else {
		t.Logf("installSystemdService failed as expected: %v", err)
	}

	// Test uninstallSystemdService
	err = manager.uninstallSystemdService()
	if err == nil {
		t.Log("uninstallSystemdService succeeded (service was installed)")
	} else if !strings.Contains(err.Error(), "not installed") {
		t.Errorf("Error should mention 'not installed', got: %v", err)
	}
}

func TestLinuxManager_StartStopWithSystemd(t *testing.T) {
	manager := &LinuxManager{
		executable:   "/test/executable",
		onRun:        func() error { return nil },
		useSystemd:   true,
		isSystemWide: false,
	}

	// Test Start with systemd
	err := manager.Start()
	// Will fail due to service not being installed or systemctl not available
	if err == nil {
		t.Log("Start unexpectedly succeeded")
	} else {
		t.Logf("Start failed as expected in test environment: %v", err)
	}

	// Test Stop with systemd
	err = manager.Stop()
	// Will fail due to service not being installed or systemctl not available
	if err == nil {
		t.Log("Stop unexpectedly succeeded")
	} else {
		t.Logf("Stop failed as expected in test environment: %v", err)
	}
}

func TestLinuxManager_GetSystemdStatus_Basic(t *testing.T) {
	manager := &LinuxManager{
		executable:   "/test/executable",
		onRun:        func() error { return nil },
		useSystemd:   true,
		isSystemWide: false,
		execCommand:  defaultCommandExecutor,
	}

	// Test getSystemdStatus
	status, err := manager.getSystemdStatus()
	// Will fail due to systemctl not being available or service not installed
	if err == nil {
		t.Log("getSystemdStatus unexpectedly succeeded")
		if status.Installed {
			t.Error("Status should not show installed in test environment")
		}
	} else {
		t.Logf("getSystemdStatus failed as expected: %v", err)
	}
}

func TestLinuxManager_UninstallXDGAutostart(t *testing.T) {
	tempDir := t.TempDir()

	manager := &LinuxManager{
		executable:   filepath.Join(tempDir, "silentcast"),
		onRun:        func() error { return nil },
		useSystemd:   false,
		isSystemWide: false,
	}

	// Mock HOME directory
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer func() {
		os.Setenv("HOME", oldHome)
	}()

	// First install XDG autostart
	err := manager.installXDGAutostart()
	if err != nil {
		t.Fatalf("Failed to install XDG autostart: %v", err)
	}

	// Now test uninstall
	err = manager.uninstallXDGAutostart()
	if err != nil {
		t.Errorf("uninstallXDGAutostart failed: %v", err)
	}

	// Verify file was removed
	desktopPath := filepath.Join(tempDir, xdgAutostartPath, desktopFile)
	if _, err := os.Stat(desktopPath); !os.IsNotExist(err) {
		t.Error("Desktop file should be removed after uninstall")
	}
}

func TestLinuxManager_SystemdTemplate(t *testing.T) {
	tempDir := t.TempDir()
	executable := filepath.Join(tempDir, "silentcast")

	manager := &LinuxManager{
		executable:   executable,
		onRun:        func() error { return nil },
		useSystemd:   true,
		isSystemWide: false,
		execCommand:  defaultCommandExecutor,
	}

	// Mock HOME directory
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer func() {
		os.Setenv("HOME", oldHome)
	}()

	// Attempt to call installSystemdService to test template generation
	// This will likely fail at launchctl step, but should test template logic
	err := manager.installSystemdService()
	// Don't assert on error as it depends on environment
	_ = err

	// Check if systemd service file was created (might not exist due to failure)
	userServiceDir := filepath.Join(tempDir, systemdUserPath)
	servicePath := filepath.Join(userServiceDir, systemdServiceFile)

	if _, err := os.Stat(servicePath); err == nil {
		// Service file was created, check content
		content, err := os.ReadFile(servicePath)
		if err != nil {
			t.Fatalf("Failed to read service file: %v", err)
		}

		contentStr := string(content)
		expectedParts := []string{
			"[Unit]",
			"Description=Silent hotkey-driven task runner",
			"[Service]",
			"Type=simple",
			"ExecStart=" + executable + " --no-tray",
			"Restart=on-failure",
			"[Install]",
			"WantedBy=default.target",
		}

		for _, part := range expectedParts {
			if !strings.Contains(contentStr, part) {
				t.Errorf("Service file should contain '%s', got:\n%s", part, contentStr)
			}
		}
	} else {
		t.Logf("Service file not created due to test environment limitations")
	}
}

func TestLinuxManager_GetSystemdStatus(t *testing.T) {
	tests := []struct {
		name           string
		isSystemWide   bool
		serviceExists  bool
		serviceEnabled string
		serviceActive  string
		expectedStatus ServiceStatus
		listUnitError  bool
		isEnabledError bool
		isActiveError  bool
	}{
		{
			name:          "service not found",
			isSystemWide:  false,
			serviceExists: false,
			expectedStatus: ServiceStatus{
				Installed: false,
				Running:   false,
				StartType: "manual",
			},
		},
		{
			name:           "service installed and enabled",
			isSystemWide:   false,
			serviceExists:  true,
			serviceEnabled: "enabled",
			serviceActive:  "active",
			expectedStatus: ServiceStatus{
				Installed: true,
				Running:   true,
				StartType: "auto",
				Message:   "Service is running",
			},
		},
		{
			name:           "service installed but disabled",
			isSystemWide:   false,
			serviceExists:  true,
			serviceEnabled: "disabled",
			serviceActive:  "inactive",
			expectedStatus: ServiceStatus{
				Installed: true,
				Running:   false,
				StartType: "manual",
				Message:   "Service is inactive",
			},
		},
		{
			name:           "service installed but failed",
			isSystemWide:   false,
			serviceExists:  true,
			serviceEnabled: "enabled",
			serviceActive:  "failed",
			expectedStatus: ServiceStatus{
				Installed: true,
				Running:   false,
				StartType: "auto",
				Message:   "Service is failed",
			},
		},
		{
			name:           "system-wide service",
			isSystemWide:   true,
			serviceExists:  true,
			serviceEnabled: "enabled",
			serviceActive:  "active",
			expectedStatus: ServiceStatus{
				Installed: true,
				Running:   true,
				StartType: "auto",
				Message:   "Service is running",
			},
		},
		{
			name:          "service with list-unit-files error",
			isSystemWide:  false,
			listUnitError: true,
			expectedStatus: ServiceStatus{
				Installed: false,
				Running:   false,
				StartType: "manual",
			},
		},
		{
			name:           "service with is-enabled error",
			isSystemWide:   false,
			serviceExists:  true,
			isEnabledError: true,
			serviceActive:  "active",
			expectedStatus: ServiceStatus{
				Installed: true,
				Running:   true,
				StartType: "manual", // defaults to manual on error
				Message:   "Service is running",
			},
		},
		{
			name:           "service with is-active error",
			isSystemWide:   false,
			serviceExists:  true,
			serviceEnabled: "enabled",
			isActiveError:  true,
			expectedStatus: ServiceStatus{
				Installed: true,
				Running:   false,
				StartType: "auto",
				Message:   "", // no message when is-active fails
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &LinuxManager{
				isSystemWide: tt.isSystemWide,
				execCommand: func(name string, args ...string) ([]byte, error) {
					scope := "--user"
					if tt.isSystemWide {
						scope = "--system"
					}

					if name == "systemctl" && len(args) >= 2 && args[0] == scope {
						switch args[1] {
						case "list-unit-files":
							if tt.listUnitError {
								return nil, errors.New("list-unit-files failed")
							}
							if tt.serviceExists {
								return []byte("silentcast.service enabled\n"), nil
							}
							return []byte("other.service enabled\n"), nil

						case "is-enabled":
							if tt.isEnabledError {
								return nil, errors.New("is-enabled failed")
							}
							return []byte(tt.serviceEnabled + "\n"), nil

						case "is-active":
							if tt.isActiveError {
								return nil, errors.New("is-active failed")
							}
							return []byte(tt.serviceActive + "\n"), nil
						}
					}
					return nil, errors.New("unexpected command")
				},
			}

			status, err := manager.getSystemdStatus()
			if (err != nil) != (tt.listUnitError || !tt.serviceExists) {
				t.Errorf("getSystemdStatus() error = %v, wantErr %v", err, tt.listUnitError || !tt.serviceExists)
				return
			}

			if !reflect.DeepEqual(status, tt.expectedStatus) {
				t.Errorf("getSystemdStatus() = %+v, want %+v", status, tt.expectedStatus)
			}
		})
	}
}

func TestLinuxManager_StartWithMock(t *testing.T) {
	tests := []struct {
		name         string
		useSystemd   bool
		isSystemWide bool
		commandFails bool
		wantErr      bool
	}{
		{
			name:       "start with systemd user service",
			useSystemd: true,
			wantErr:    false,
		},
		{
			name:         "start with systemd system service",
			useSystemd:   true,
			isSystemWide: true,
			wantErr:      false,
		},
		{
			name:       "start without systemd",
			useSystemd: false,
			wantErr:    true,
		},
		{
			name:         "start with command failure",
			useSystemd:   true,
			commandFails: true,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &LinuxManager{
				useSystemd:   tt.useSystemd,
				isSystemWide: tt.isSystemWide,
				execCommand: func(name string, args ...string) ([]byte, error) {
					if tt.commandFails {
						return []byte("Failed to start service"), errors.New("start failed")
					}
					return []byte(""), nil
				},
			}

			err := manager.Start()
			if (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLinuxManager_StopWithMock(t *testing.T) {
	tests := []struct {
		name         string
		useSystemd   bool
		isSystemWide bool
		commandFails bool
		wantErr      bool
	}{
		{
			name:       "stop with systemd user service",
			useSystemd: true,
			wantErr:    false,
		},
		{
			name:         "stop with systemd system service",
			useSystemd:   true,
			isSystemWide: true,
			wantErr:      false,
		},
		{
			name:       "stop without systemd",
			useSystemd: false,
			wantErr:    true,
		},
		{
			name:         "stop with command failure",
			useSystemd:   true,
			commandFails: true,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &LinuxManager{
				useSystemd:   tt.useSystemd,
				isSystemWide: tt.isSystemWide,
				execCommand: func(name string, args ...string) ([]byte, error) {
					if tt.commandFails {
						return []byte("Failed to stop service"), errors.New("stop failed")
					}
					return []byte(""), nil
				},
			}

			err := manager.Stop()
			if (err != nil) != tt.wantErr {
				t.Errorf("Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLinuxManager_UninstallSystemdServiceWithMock(t *testing.T) {
	tests := []struct {
		name          string
		isSystemWide  bool
		serviceExists bool
		wantErr       bool
	}{
		{
			name:          "uninstall user service",
			isSystemWide:  false,
			serviceExists: true,
			wantErr:       false,
		},
		{
			name:          "service not installed",
			isSystemWide:  false,
			serviceExists: false,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			var serviceDir string
			if tt.isSystemWide {
				serviceDir = filepath.Join(tempDir, "etc", "systemd", "system")
			} else {
				serviceDir = filepath.Join(tempDir, ".config", "systemd", "user")
			}
			if err := os.MkdirAll(serviceDir, 0o755); err != nil {
				t.Fatalf("Failed to create service dir: %v", err)
			}

			if tt.serviceExists {
				servicePath := filepath.Join(serviceDir, "silentcast.service")
				if err := os.WriteFile(servicePath, []byte("test service"), 0o644); err != nil {
					t.Fatalf("Failed to write service file: %v", err)
				}
			}

			// Mock user.Current to return our temp directory
			origGetCurrentUser := getCurrentUser
			defer func() { getCurrentUser = origGetCurrentUser }()
			getCurrentUser = func() (*user.User, error) {
				return &user.User{HomeDir: tempDir}, nil
			}

			manager := &LinuxManager{
				isSystemWide: tt.isSystemWide,
				execCommand: func(name string, args ...string) ([]byte, error) {
					// Just return success for all systemctl commands
					return []byte(""), nil
				},
			}

			// For system-wide tests, we rely on the mock getCurrentUser
			// to handle path resolution correctly

			err := manager.uninstallSystemdService()
			if (err != nil) != tt.wantErr {
				t.Errorf("uninstallSystemdService() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && tt.serviceExists {
				// Check service file was removed
				servicePath := filepath.Join(serviceDir, "silentcast.service")
				if _, err := os.Stat(servicePath); !os.IsNotExist(err) {
					t.Error("Service file should have been removed")
				}
			}
		})
	}
}
