//go:build linux && integration

package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestLinuxServiceCommands tests the Linux service management commands
func TestLinuxServiceCommands(t *testing.T) {
	// Skip if not running as regular user
	if os.Getuid() == 0 {
		t.Skip("Skipping service test when running as root")
	}

	// Build the test binary
	testBinary := filepath.Join(t.TempDir(), "silentcast-test")
	buildCmd := exec.Command("go", "build", "-tags", "nogohook notray", "-o", testBinary, "../../cmd/silentcast")
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build test binary: %v\nOutput: %s", err, output)
	}

	// Create temporary HOME for isolated testing
	tmpHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", originalHome)

	// Test service status when not installed
	t.Run("Status Not Installed", func(t *testing.T) {
		cmd := exec.Command(testBinary, "--service-status")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("service-status failed: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "Installed: false") {
			t.Errorf("Expected 'Installed: false', got: %s", outputStr)
		}
		if !strings.Contains(outputStr, "Service not installed") {
			t.Errorf("Expected 'Service not installed', got: %s", outputStr)
		}
	})

	// Test service installation
	t.Run("Service Install", func(t *testing.T) {
		cmd := exec.Command(testBinary, "--service-install")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("service-install failed: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "Installing SilentCast service") {
			t.Errorf("Missing installation message in: %s", outputStr)
		}
		if !strings.Contains(outputStr, "Service installed successfully") {
			t.Errorf("Missing success message in: %s", outputStr)
		}

		// Verify files were created
		systemdPath := filepath.Join(tmpHome, ".config/systemd/user/silentcast.service")
		if _, err := os.Stat(systemdPath); os.IsNotExist(err) {
			t.Error("Systemd service file was not created")
		}

		xdgPath := filepath.Join(tmpHome, ".config/autostart/silentcast.desktop")
		if _, err := os.Stat(xdgPath); os.IsNotExist(err) {
			t.Error("XDG autostart file was not created")
		}
	})

	// Test service status when installed
	t.Run("Status Installed", func(t *testing.T) {
		// Ensure service is installed from previous test
		// or install it if running this test in isolation
		installCmd := exec.Command(testBinary, "--service-install")
		installCmd.Run()

		cmd := exec.Command(testBinary, "--service-status")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("service-status failed: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "Installed: true") {
			t.Errorf("Expected 'Installed: true', got: %s", outputStr)
		}
	})

	// Test duplicate installation
	t.Run("Duplicate Install", func(t *testing.T) {
		// First install
		installCmd := exec.Command(testBinary, "--service-install")
		installCmd.Run()

		// Try to install again
		cmd := exec.Command(testBinary, "--service-install")
		output, err := cmd.CombinedOutput()
		if err == nil {
			t.Error("Expected error when installing duplicate service")
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "service already installed") {
			t.Errorf("Expected 'service already installed' error, got: %s", outputStr)
		}
	})

	// Test service uninstall
	t.Run("Service Uninstall", func(t *testing.T) {
		// Ensure service is installed
		installCmd := exec.Command(testBinary, "--service-install")
		installCmd.Run()

		// Uninstall
		cmd := exec.Command(testBinary, "--service-uninstall")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("service-uninstall failed: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "Uninstalling SilentCast service") {
			t.Errorf("Missing uninstallation message in: %s", outputStr)
		}
		if !strings.Contains(outputStr, "Service uninstalled successfully") {
			t.Errorf("Missing success message in: %s", outputStr)
		}

		// Verify files were removed
		systemdPath := filepath.Join(tmpHome, ".config/systemd/user/silentcast.service")
		if _, err := os.Stat(systemdPath); err == nil {
			t.Error("Systemd service file was not removed")
		}

		xdgPath := filepath.Join(tmpHome, ".config/autostart/silentcast.desktop")
		if _, err := os.Stat(xdgPath); err == nil {
			t.Error("XDG autostart file was not removed")
		}
	})

	// Test uninstalling non-existent service
	t.Run("Uninstall Non-Existent", func(t *testing.T) {
		// Make sure service is not installed
		uninstallCmd := exec.Command(testBinary, "--service-uninstall")
		uninstallCmd.Run()
		time.Sleep(100 * time.Millisecond)

		// Try to uninstall again
		cmd := exec.Command(testBinary, "--service-uninstall")
		output, err := cmd.CombinedOutput()
		if err == nil {
			t.Error("Expected error when uninstalling non-existent service")
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "failed to uninstall service") {
			t.Errorf("Expected uninstall error, got: %s", outputStr)
		}
	})

	// Test systemd-specific commands (if systemd is available)
	if hasSystemd(t) {
		t.Run("Systemd Start/Stop", func(t *testing.T) {
			// Install service first
			installCmd := exec.Command(testBinary, "--service-install")
			installCmd.Run()

			// Test start (will likely fail in test environment)
			startCmd := exec.Command(testBinary, "--service-start")
			startOutput, _ := startCmd.CombinedOutput()
			// We expect this to fail in most test environments
			// Just check that the command runs
			if !strings.Contains(string(startOutput), "Starting SilentCast service") {
				t.Logf("Start command output: %s", startOutput)
			}

			// Test stop
			stopCmd := exec.Command(testBinary, "--service-stop")
			stopOutput, _ := stopCmd.CombinedOutput()
			if !strings.Contains(string(stopOutput), "Stopping SilentCast service") {
				t.Logf("Stop command output: %s", stopOutput)
			}
		})
	}
}

// hasSystemd checks if systemd is available
func hasSystemd(t *testing.T) bool {
	cmd := exec.Command("systemctl", "--version")
	err := cmd.Run()
	return err == nil
}