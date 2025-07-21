//go:build e2e

package e2e

import (
	"os"
	"runtime"
	"testing"
	"time"
)

// TestTrayDisabledMode tests application running without system tray
func TestTrayDisabledMode(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Setup test configuration
	if err := env.SetupSpellbook(TestSpellbookTemplate); err != nil {
		t.Fatalf("Failed to setup spellbook: %v", err)
	}

	// Start application with --no-tray flag (this is default in our E2E tests)
	if err := env.StartApplication("--no-tray"); err != nil {
		t.Fatalf("Failed to start application: %v", err)
	}

	if err := env.WaitForStartup(); err != nil {
		t.Fatalf("Application startup failed: %v", err)
	}

	// Verify application runs without tray
	logs, err := env.GetLogs()
	if err != nil {
		t.Fatalf("Failed to get logs: %v", err)
	}

	// Look for indicators that tray is disabled
	trayDisabledIndicators := []string{
		"no-tray",
		"tray disabled",
		"without tray",
		"System tray disabled",
	}

	found := false
	for _, indicator := range trayDisabledIndicators {
		if contains(logs, indicator) {
			found = true
			t.Logf("Found tray disabled indicator: %s", indicator)
			break
		}
	}

	if !found {
		t.Logf("Tray disabled indicator not found in logs (may be normal)")
	}

	// Test that hotkeys still work without tray
	if err := env.SimulateHotkey("e"); err != nil {
		t.Fatalf("Failed to simulate hotkey: %v", err)
	}

	if err := env.WaitForAction("editor", 10*time.Second); err != nil {
		t.Fatalf("Action execution timeout without tray: %v", err)
	}

	// Verify no tray-related errors
	if contains(logs, "tray error") || contains(logs, "tray failed") {
		t.Error("Tray-related errors found even with --no-tray")
	}

	// Verify no errors
	env.AssertNoErrors()
}

// TestTrayEnabledMode tests application with system tray (if supported)
func TestTrayEnabledMode(t *testing.T) {
	// Skip tray tests in headless environments or unsupported platforms
	if skipTrayTest() {
		t.Skip("Skipping tray test in headless environment or unsupported platform")
	}

	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Setup test configuration with tray-specific actions
	config := `
spells:
  t: tray-test
  s: status-check

grimoire:
  tray-test:
    type: script
    command: echo "Tray test executed"
    description: "Test tray functionality"
    show_output: true
    
  status-check:
    type: script
    command: echo "Status checked from tray"
    description: "Status check via tray"
    show_output: true
`

	if err := env.SetupSpellbook(config); err != nil {
		t.Fatalf("Failed to setup spellbook: %v", err)
	}

	// Start application with tray enabled
	if err := env.StartApplication(); err != nil {
		t.Fatalf("Failed to start application: %v", err)
	}

	// Wait for startup with longer timeout for tray initialization
	ctx, cancel := env.ctx, env.cancel
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- env.WaitForStartup()
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Application startup failed: %v", err)
		}
	case <-time.After(20 * time.Second):
		t.Fatalf("Application startup timeout with tray")
	case <-ctx.Done():
		t.Fatalf("Context cancelled during startup")
	}

	// Check for tray initialization
	logs, err := env.GetLogs()
	if err != nil {
		t.Fatalf("Failed to get logs: %v", err)
	}

	// Look for tray initialization indicators
	trayIndicators := []string{
		"tray",
		"system tray",
		"tray initialized",
		"tray manager",
	}

	trayFound := false
	for _, indicator := range trayIndicators {
		if contains(logs, indicator) {
			trayFound = true
			t.Logf("Found tray indicator: %s", indicator)
			break
		}
	}

	if !trayFound {
		t.Logf("Tray initialization not explicitly logged")
	}

	// Test that actions still work with tray enabled
	if err := env.SimulateHotkey("t"); err != nil {
		t.Fatalf("Failed to simulate hotkey: %v", err)
	}

	if err := env.WaitForAction("tray-test", 15*time.Second); err != nil {
		t.Fatalf("Tray test action timeout: %v", err)
	}

	// Verify action execution
	logs, err = env.GetLogs()
	if err != nil {
		t.Fatalf("Failed to get logs: %v", err)
	}

	if !contains(logs, "Tray test executed") {
		t.Error("Tray test action did not execute")
	}

	// Verify no tray-related errors
	if contains(logs, "tray error") || contains(logs, "tray panic") {
		t.Error("Tray-related errors found")
	}

	// Verify no errors
	env.AssertNoErrors()
}

// TestTrayConfigurationHandling tests tray configuration options
func TestTrayConfigurationHandling(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Setup configuration with tray-related settings
	config := `
daemon:
  tray: false
  
spells:
  c: config-test

grimoire:
  config-test:
    type: script
    command: echo "Configuration test"
    description: "Test tray configuration"
    show_output: true
`

	if err := env.SetupSpellbook(config); err != nil {
		t.Fatalf("Failed to setup spellbook: %v", err)
	}

	// Start application
	if err := env.StartApplication(); err != nil {
		t.Fatalf("Failed to start application: %v", err)
	}

	if err := env.WaitForStartup(); err != nil {
		t.Fatalf("Application startup failed: %v", err)
	}

	// Test configuration-based tray disabling
	if err := env.SimulateHotkey("c"); err != nil {
		t.Fatalf("Failed to simulate hotkey: %v", err)
	}

	if err := env.WaitForAction("config-test", 10*time.Second); err != nil {
		t.Fatalf("Config test action timeout: %v", err)
	}

	// Verify application respects tray configuration
	logs, err := env.GetLogs()
	if err != nil {
		t.Fatalf("Failed to get logs: %v", err)
	}

	if !contains(logs, "Configuration test") {
		t.Error("Configuration test action did not execute")
	}

	// Verify no errors
	env.AssertNoErrors()
}

// TestTrayFallbackBehavior tests tray fallback when tray is unavailable
func TestTrayFallbackBehavior(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Setup test configuration
	if err := env.SetupSpellbook(TestSpellbookTemplate); err != nil {
		t.Fatalf("Failed to setup spellbook: %v", err)
	}

	// Try to start with tray, but expect graceful fallback in headless environment
	if err := env.StartApplication(); err != nil {
		t.Fatalf("Failed to start application: %v", err)
	}

	if err := env.WaitForStartup(); err != nil {
		t.Fatalf("Application startup failed: %v", err)
	}

	// Application should start successfully even if tray is unavailable
	logs, err := env.GetLogs()
	if err != nil {
		t.Fatalf("Failed to get logs: %v", err)
	}

	// Check for graceful tray fallback messages
	fallbackIndicators := []string{
		"tray unavailable",
		"tray fallback",
		"no display",
		"headless",
		"tray disabled",
	}

	for _, indicator := range fallbackIndicators {
		if contains(logs, indicator) {
			t.Logf("Found tray fallback indicator: %s", indicator)
		}
	}

	// Test that core functionality works regardless of tray status
	if err := env.SimulateHotkey("e"); err != nil {
		t.Fatalf("Failed to simulate hotkey: %v", err)
	}

	if err := env.WaitForAction("editor", 10*time.Second); err != nil {
		t.Fatalf("Action execution failed during tray fallback: %v", err)
	}

	// Verify application didn't crash due to tray issues
	if contains(logs, "panic") || contains(logs, "FATAL") {
		t.Error("Application crashed during tray fallback")
	}

	// Verify no errors
	env.AssertNoErrors()
}

// TestTrayResourceCleanup tests proper tray resource cleanup
func TestTrayResourceCleanup(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Setup test configuration
	if err := env.SetupSpellbook(TestSpellbookTemplate); err != nil {
		t.Fatalf("Failed to setup spellbook: %v", err)
	}

	// Start application
	if err := env.StartApplication(); err != nil {
		t.Fatalf("Failed to start application: %v", err)
	}

	if err := env.WaitForStartup(); err != nil {
		t.Fatalf("Application startup failed: %v", err)
	}

	// Let application run for a moment
	time.Sleep(2 * time.Second)

	// Execute an action to ensure tray system is active
	if err := env.SimulateHotkey("e"); err != nil {
		t.Fatalf("Failed to simulate hotkey: %v", err)
	}

	if err := env.WaitForAction("editor", 10*time.Second); err != nil {
		t.Fatalf("Action execution timeout: %v", err)
	}

	// Test graceful shutdown
	if err := env.StopApplication(); err != nil {
		t.Errorf("Application shutdown failed: %v", err)
	}

	// Verify clean shutdown without tray-related issues
	logs, err := env.GetLogs()
	if err == nil {
		if contains(logs, "tray panic") || contains(logs, "tray error") {
			t.Error("Tray-related errors during shutdown")
		}

		// Look for clean shutdown indicators
		cleanupIndicators := []string{
			"shutdown",
			"cleanup",
			"stopped",
			"exiting",
		}

		found := false
		for _, indicator := range cleanupIndicators {
			if contains(logs, indicator) {
				found = true
				t.Logf("Found cleanup indicator: %s", indicator)
				break
			}
		}

		if !found {
			t.Logf("Clean shutdown indicator not found (may be normal)")
		}
	}
}

// skipTrayTest determines if tray tests should be skipped
func skipTrayTest() bool {
	// Skip on platforms where tray testing is problematic
	switch runtime.GOOS {
	case "linux":
		// Skip if no DISPLAY environment variable (headless)
		display := os.Getenv("DISPLAY")
		if display == "" {
			return true
		}
	case "darwin":
		// macOS tray testing requires GUI session
		// In CI environments, this is typically not available
		return true
	case "windows":
		// Windows tray testing in CI environments can be problematic
		return true
	}

	return false
}

// TestTrayStubMode tests tray stub implementation
func TestTrayStubMode(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Setup test configuration
	if err := env.SetupSpellbook(TestSpellbookTemplate); err != nil {
		t.Fatalf("Failed to setup spellbook: %v", err)
	}

	// Start application with no-tray to force stub mode
	if err := env.StartApplication("--no-tray"); err != nil {
		t.Fatalf("Failed to start application: %v", err)
	}

	if err := env.WaitForStartup(); err != nil {
		t.Fatalf("Application startup failed: %v", err)
	}

	// Verify stub mode works correctly
	logs, err := env.GetLogs()
	if err != nil {
		t.Fatalf("Failed to get logs: %v", err)
	}

	// Test core functionality in stub mode
	if err := env.SimulateHotkey("t"); err != nil {
		t.Fatalf("Failed to simulate hotkey: %v", err)
	}

	if err := env.WaitForAction("terminal", 10*time.Second); err != nil {
		t.Fatalf("Action execution timeout in stub mode: %v", err)
	}

	// Verify action execution in stub mode
	logs, err = env.GetLogs()
	if err != nil {
		t.Fatalf("Failed to get logs: %v", err)
	}

	if !contains(logs, "terminal") {
		t.Error("Action did not execute in tray stub mode")
	}

	// Verify no tray-related errors in stub mode
	if contains(logs, "tray error") {
		t.Error("Tray errors found in stub mode")
	}

	// Verify no errors
	env.AssertNoErrors()
}
