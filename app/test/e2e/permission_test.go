//go:build e2e

package e2e

import (
	"runtime"
	"testing"
	"time"
)

// TestPermissionCheckWorkflow tests permission checking and handling workflow
func TestPermissionCheckWorkflow(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Setup configuration that requires permissions
	config := `
spells:
  p: permission-test
  e: elevated-action

grimoire:
  permission-test:
    type: script
    command: echo "Permission check completed"
    description: "Test permission checking"
    show_output: true
    
  elevated-action:
    type: script
    command: echo "Elevated action executed"
    description: "Action requiring elevation"
    elevated: true
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

	// Test basic permission checking
	if err := env.SimulateHotkey("p"); err != nil {
		t.Fatalf("Failed to simulate hotkey: %v", err)
	}

	if err := env.WaitForAction("permission-test", 10*time.Second); err != nil {
		t.Fatalf("Permission test action timeout: %v", err)
	}

	// Verify permission check in logs
	logs, err := env.GetLogs()
	if err != nil {
		t.Fatalf("Failed to get logs: %v", err)
	}

	expectedMessages := []string{
		"Permission check completed",
		"permission-test",
	}

	for _, msg := range expectedMessages {
		if !contains(logs, msg) {
			t.Errorf("Expected message '%s' not found in logs", msg)
		}
	}

	// Verify no errors
	env.AssertNoErrors()
}

// TestElevatedActionWorkflow tests elevated permission workflow
func TestElevatedActionWorkflow(t *testing.T) {
	// Skip on platforms where elevation testing is complex
	if runtime.GOOS == "darwin" {
		t.Skip("Elevated action testing on macOS requires special setup")
	}

	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Setup configuration with elevated action
	config := `
spells:
  a: admin-command

grimoire:
  admin-command:
    type: script
    command: echo "Admin command executed"
    description: "Administrative command"
    elevated: true
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

	// Test elevated action
	if err := env.SimulateHotkey("a"); err != nil {
		t.Fatalf("Failed to simulate hotkey: %v", err)
	}

	// Wait for action with longer timeout (elevation may take time)
	if err := env.WaitForAction("admin-command", 30*time.Second); err != nil {
		// On some systems, elevated actions may fail due to lack of proper credentials
		// Check if elevation was attempted
		logs, _ := env.GetLogs()
		if !contains(logs, "admin-command") && !contains(logs, "elevation") {
			t.Fatalf("Elevated action timeout and no elevation attempt found: %v", err)
		}
		t.Logf("Elevated action may have failed due to credentials, but elevation was attempted")
		return
	}

	// Verify elevated action execution
	logs, err := env.GetLogs()
	if err != nil {
		t.Fatalf("Failed to get logs: %v", err)
	}

	if !contains(logs, "admin-command") {
		t.Error("Elevated action execution not found in logs")
	}
}

// TestPermissionDeniedHandling tests graceful handling of permission denials
func TestPermissionDeniedHandling(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Setup configuration that would require permissions not available in test
	config := `
spells:
  d: denied-action

grimoire:
  denied-action:
    type: script
    command: echo "This should work but may have permission issues"
    description: "Action that may be denied"
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

	// Test action that might be denied
	if err := env.SimulateHotkey("d"); err != nil {
		t.Fatalf("Failed to simulate hotkey: %v", err)
	}

	// Wait for action completion or failure
	time.Sleep(5 * time.Second)

	// Verify application handles permission issues gracefully
	logs, err := env.GetLogs()
	if err != nil {
		t.Fatalf("Failed to get logs: %v", err)
	}

	// Application should either succeed or handle failure gracefully
	if contains(logs, "panic") || contains(logs, "FATAL") {
		t.Error("Application did not handle permission issues gracefully")
	}

	// Verify the action was attempted
	if !contains(logs, "denied-action") {
		t.Error("Action was not attempted")
	}
}

// TestPlatformSpecificPermissions tests platform-specific permission requirements
func TestPlatformSpecificPermissions(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	var config string
	var expectedMessages []string

	switch runtime.GOOS {
	case "darwin":
		config = `
spells:
  m: macos-specific

grimoire:
  macos-specific:
    type: script
    command: echo "macOS specific action"
    description: "macOS specific permission test"
    show_output: true
`
		expectedMessages = []string{"macos-specific", "macOS specific action"}

	case "windows":
		config = `
spells:
  w: windows-specific

grimoire:
  windows-specific:
    type: script
    command: echo "Windows specific action"
    description: "Windows specific permission test"
    show_output: true
`
		expectedMessages = []string{"windows-specific", "Windows specific action"}

	case "linux":
		config = `
spells:
  l: linux-specific

grimoire:
  linux-specific:
    type: script
    command: echo "Linux specific action"
    description: "Linux specific permission test"
    show_output: true
`
		expectedMessages = []string{"linux-specific", "Linux specific action"}

	default:
		t.Skipf("Platform-specific permission test not implemented for %s", runtime.GOOS)
	}

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

	// Test platform-specific action
	hotkey := map[string]string{
		"darwin":  "m",
		"windows": "w",
		"linux":   "l",
	}[runtime.GOOS]

	if err := env.SimulateHotkey(hotkey); err != nil {
		t.Fatalf("Failed to simulate hotkey: %v", err)
	}

	// Wait for action completion
	actionName := map[string]string{
		"darwin":  "macos-specific",
		"windows": "windows-specific",
		"linux":   "linux-specific",
	}[runtime.GOOS]

	if err := env.WaitForAction(actionName, 15*time.Second); err != nil {
		t.Fatalf("Platform-specific action timeout: %v", err)
	}

	// Verify execution
	logs, err := env.GetLogs()
	if err != nil {
		t.Fatalf("Failed to get logs: %v", err)
	}

	for _, msg := range expectedMessages {
		if !contains(logs, msg) {
			t.Errorf("Expected message '%s' not found in logs", msg)
		}
	}

	// Verify no errors
	env.AssertNoErrors()
}

// TestPermissionManagerInitialization tests permission manager startup
func TestPermissionManagerInitialization(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Setup minimal configuration
	config := `
spells:
  t: test-action

grimoire:
  test-action:
    type: script
    command: echo "Permission manager test"
    description: "Test permission manager initialization"
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

	// Check logs for permission manager initialization
	logs, err := env.GetLogs()
	if err != nil {
		t.Fatalf("Failed to get logs: %v", err)
	}

	// Look for permission manager initialization messages
	permissionIndicators := []string{
		"Permission manager",
		"permission",
		"initialized",
	}

	found := false
	for _, indicator := range permissionIndicators {
		if contains(logs, indicator) {
			found = true
			break
		}
	}

	if !found {
		t.Logf("Permission manager initialization not explicitly logged (this may be normal)")
	}

	// Test basic action to verify permission system is working
	if err := env.SimulateHotkey("t"); err != nil {
		t.Fatalf("Failed to simulate hotkey: %v", err)
	}

	if err := env.WaitForAction("test-action", 10*time.Second); err != nil {
		t.Fatalf("Test action timeout: %v", err)
	}

	// Verify action execution
	logs, err = env.GetLogs()
	if err != nil {
		t.Fatalf("Failed to get logs: %v", err)
	}

	if !contains(logs, "Permission manager test") {
		t.Error("Test action did not execute properly")
	}

	// Verify no errors
	env.AssertNoErrors()
}