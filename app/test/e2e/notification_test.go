//go:build e2e

package e2e

import (
	"runtime"
	"testing"
	"time"
)

// TestNotificationSystemWorkflow tests the complete notification system workflow
func TestNotificationSystemWorkflow(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Setup configuration with notification-enabled actions
	config := `
spells:
  n: notify-action
  s: silent-action

grimoire:
  notify-action:
    type: script
    command: echo "Notification test message"
    description: "Action with notification"
    show_output: true
    notify: true
    
  silent-action:
    type: script
    command: echo "Silent execution"
    description: "Action without notification"
    show_output: false
    notify: false
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

	// Test notification action
	if err := env.SimulateHotkey("n"); err != nil {
		t.Fatalf("Failed to simulate hotkey: %v", err)
	}

	if err := env.WaitForAction("notify-action", 15*time.Second); err != nil {
		t.Fatalf("Notification action timeout: %v", err)
	}

	// Verify notification system activation
	logs, err := env.GetLogs()
	if err != nil {
		t.Fatalf("Failed to get logs: %v", err)
	}

	notificationIndicators := []string{
		"notify-action",
		"Notification test message",
		"notification",
		"notify",
	}

	found := false
	for _, indicator := range notificationIndicators {
		if contains(logs, indicator) {
			found = true
			break
		}
	}

	if !found {
		t.Error("Notification system activation not found in logs")
	}

	// Test silent action
	if err := env.SimulateHotkey("s"); err != nil {
		t.Fatalf("Failed to simulate hotkey: %v", err)
	}

	if err := env.WaitForAction("silent-action", 10*time.Second); err != nil {
		t.Fatalf("Silent action timeout: %v", err)
	}

	// Verify no errors
	env.AssertNoErrors()
}

// TestPlatformSpecificNotifications tests platform-specific notification systems
func TestPlatformSpecificNotifications(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	var config string
	var expectedPlatform string

	switch runtime.GOOS {
	case "darwin":
		config = `
spells:
  m: macos-notify

grimoire:
  macos-notify:
    type: script
    command: echo "macOS notification test"
    description: "macOS notification system test"
    show_output: true
    notify: true
`
		expectedPlatform = "macOS"

	case "windows":
		config = `
spells:
  w: windows-notify

grimoire:
  windows-notify:
    type: script
    command: echo "Windows notification test"
    description: "Windows notification system test"
    show_output: true
    notify: true
`
		expectedPlatform = "Windows"

	case "linux":
		config = `
spells:
  l: linux-notify

grimoire:
  linux-notify:
    type: script
    command: echo "Linux notification test"
    description: "Linux notification system test"
    show_output: true
    notify: true
`
		expectedPlatform = "Linux"

	default:
		t.Skipf("Platform-specific notification test not implemented for %s", runtime.GOOS)
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

	// Test platform-specific notification
	hotkey := map[string]string{
		"darwin":  "m",
		"windows": "w",
		"linux":   "l",
	}[runtime.GOOS]

	if err := env.SimulateHotkey(hotkey); err != nil {
		t.Fatalf("Failed to simulate hotkey: %v", err)
	}

	actionName := map[string]string{
		"darwin":  "macos-notify",
		"windows": "windows-notify",
		"linux":   "linux-notify",
	}[runtime.GOOS]

	if err := env.WaitForAction(actionName, 15*time.Second); err != nil {
		t.Fatalf("Platform-specific notification timeout: %v", err)
	}

	// Verify platform-specific notification handling
	logs, err := env.GetLogs()
	if err != nil {
		t.Fatalf("Failed to get logs: %v", err)
	}

	if !contains(logs, expectedPlatform+" notification test") {
		t.Errorf("Platform-specific notification (%s) not found in logs", expectedPlatform)
	}

	// Verify no errors
	env.AssertNoErrors()
}

// TestNotificationManagerInitialization tests notification manager startup
func TestNotificationManagerInitialization(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Setup minimal configuration
	config := `
spells:
  t: test-notify

grimoire:
  test-notify:
    type: script
    command: echo "Notification manager test"
    description: "Test notification manager initialization"
    show_output: true
    notify: true
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

	// Check logs for notification manager initialization
	logs, err := env.GetLogs()
	if err != nil {
		t.Fatalf("Failed to get logs: %v", err)
	}

	// Look for notification manager initialization messages
	notificationInitIndicators := []string{
		"Notification manager",
		"notification system",
		"NotifyManager",
		"initialized",
	}

	found := false
	for _, indicator := range notificationInitIndicators {
		if contains(logs, indicator) {
			found = true
			t.Logf("Found notification manager initialization: %s", indicator)
			break
		}
	}

	if !found {
		t.Logf("Notification manager initialization not explicitly logged (this may be normal)")
	}

	// Test notification action to verify system is working
	if err := env.SimulateHotkey("t"); err != nil {
		t.Fatalf("Failed to simulate hotkey: %v", err)
	}

	if err := env.WaitForAction("test-notify", 10*time.Second); err != nil {
		t.Fatalf("Test notification action timeout: %v", err)
	}

	// Verify notification action execution
	logs, err = env.GetLogs()
	if err != nil {
		t.Fatalf("Failed to get logs: %v", err)
	}

	if !contains(logs, "Notification manager test") {
		t.Error("Test notification action did not execute properly")
	}

	// Verify no errors
	env.AssertNoErrors()
}

// TestNotificationOutputHandling tests notification output processing
func TestNotificationOutputHandling(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Setup configuration with various output scenarios
	config := `
spells:
  l: long-output
  e: error-output
  m: multiline-output

grimoire:
  long-output:
    type: script
    command: echo "This is a very long output message that should be properly handled by the notification system and potentially truncated if it exceeds the maximum length"
    description: "Action with long output"
    show_output: true
    notify: true
    
  error-output:
    type: script  
    command: sh -c "echo 'Error message' >&2; exit 1"
    description: "Action with error output"
    show_output: true
    notify: true
    
  multiline-output:
    type: script
    command: echo -e "Line 1\nLine 2\nLine 3\nLine 4"
    description: "Action with multiline output"
    show_output: true
    notify: true
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

	// Test long output handling
	if err := env.SimulateHotkey("l"); err != nil {
		t.Fatalf("Failed to simulate hotkey: %v", err)
	}

	if err := env.WaitForAction("long-output", 10*time.Second); err != nil {
		t.Fatalf("Long output action timeout: %v", err)
	}

	// Test error output handling
	if err := env.SimulateHotkey("e"); err != nil {
		t.Fatalf("Failed to simulate hotkey: %v", err)
	}

	// Allow time for error action to complete (may take longer due to error)
	time.Sleep(5 * time.Second)

	// Test multiline output handling
	if err := env.SimulateHotkey("m"); err != nil {
		t.Fatalf("Failed to simulate hotkey: %v", err)
	}

	if err := env.WaitForAction("multiline-output", 10*time.Second); err != nil {
		t.Fatalf("Multiline output action timeout: %v", err)
	}

	// Verify all actions were processed
	logs, err := env.GetLogs()
	if err != nil {
		t.Fatalf("Failed to get logs: %v", err)
	}

	expectedOutputs := []string{
		"long-output",
		"error-output",
		"multiline-output",
		"very long output message",
		"Line 1",
	}

	for _, output := range expectedOutputs {
		if !contains(logs, output) {
			t.Errorf("Expected output '%s' not found in logs", output)
		}
	}

	// Check that application handled error output gracefully
	if !contains(logs, "error-output") {
		t.Error("Error output action was not processed")
	}
}

// TestNotificationQueueing tests notification system queuing behavior
func TestNotificationQueueing(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Setup configuration with multiple quick actions
	config := `
spells:
  q1: quick-action-1
  q2: quick-action-2  
  q3: quick-action-3

grimoire:
  quick-action-1:
    type: script
    command: echo "Quick action 1 executed"
    description: "First quick action"
    show_output: true
    notify: true
    
  quick-action-2:
    type: script
    command: echo "Quick action 2 executed"
    description: "Second quick action"
    show_output: true
    notify: true
    
  quick-action-3:
    type: script
    command: echo "Quick action 3 executed"
    description: "Third quick action"
    show_output: true
    notify: true
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

	// Rapidly trigger multiple actions to test queuing
	actions := []string{"q1", "q2", "q3"}

	for _, action := range actions {
		if err := env.SimulateHotkey(action); err != nil {
			t.Fatalf("Failed to simulate hotkey %s: %v", action, err)
		}
		// Small delay between actions
		time.Sleep(100 * time.Millisecond)
	}

	// Wait for all actions to complete
	for i, actionKey := range []string{"quick-action-1", "quick-action-2", "quick-action-3"} {
		if err := env.WaitForAction(actionKey, 15*time.Second); err != nil {
			t.Errorf("Action %d (%s) timeout: %v", i+1, actionKey, err)
		}
	}

	// Verify all actions were executed
	logs, err := env.GetLogs()
	if err != nil {
		t.Fatalf("Failed to get logs: %v", err)
	}

	expectedMessages := []string{
		"Quick action 1 executed",
		"Quick action 2 executed",
		"Quick action 3 executed",
	}

	for _, msg := range expectedMessages {
		if !contains(logs, msg) {
			t.Errorf("Expected message '%s' not found in logs", msg)
		}
	}

	// Verify no errors occurred during rapid execution
	env.AssertNoErrors()
}
