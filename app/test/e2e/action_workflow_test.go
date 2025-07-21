//go:build e2e

package e2e

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestScriptActionExecution tests end-to-end script action execution
func TestScriptActionExecution(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Create test script
	testScript := filepath.Join(env.TempDir, "test_script.sh")
	scriptContent := `#!/bin/bash
echo "Test script executed successfully"
echo "Current directory: $(pwd)"
echo "Arguments: $@"
`
	if err := os.WriteFile(testScript, []byte(scriptContent), 0o755); err != nil {
		t.Fatalf("Failed to create test script: %v", err)
	}

	// Setup configuration with script action
	config := `
spells:
  s: test-script

grimoire:
  test-script:
    type: script
    command: ` + testScript + `
    description: "Test script execution"
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

	// Simulate hotkey press for script action
	if err := env.SimulateHotkey("s"); err != nil {
		t.Fatalf("Failed to simulate hotkey: %v", err)
	}

	// Wait for action to complete
	if err := env.WaitForAction("test-script", 10*time.Second); err != nil {
		t.Fatalf("Action execution timeout: %v", err)
	}

	// Verify script execution in logs
	logs, err := env.GetLogs()
	if err != nil {
		t.Fatalf("Failed to get logs: %v", err)
	}

	expectedOutputs := []string{
		"Test script executed successfully",
		"test-script",
	}

	for _, output := range expectedOutputs {
		if !contains(logs, output) {
			t.Errorf("Expected output '%s' not found in logs", output)
		}
	}

	// Verify no errors
	env.AssertNoErrors()
}

// TestAppActionExecution tests application launching workflow
func TestAppActionExecution(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Setup configuration with app action (using echo as test app)
	config := `
spells:
  a: test-app

grimoire:
  test-app:
    type: app
    command: echo
    args: ["Hello", "from", "app"]
    description: "Test application launch"
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

	// Simulate hotkey press for app action
	if err := env.SimulateHotkey("a"); err != nil {
		t.Fatalf("Failed to simulate hotkey: %v", err)
	}

	// Wait for action to complete
	if err := env.WaitForAction("test-app", 10*time.Second); err != nil {
		t.Fatalf("Action execution timeout: %v", err)
	}

	// Verify app execution in logs
	logs, err := env.GetLogs()
	if err != nil {
		t.Fatalf("Failed to get logs: %v", err)
	}

	if !contains(logs, "test-app") {
		t.Error("App action execution not found in logs")
	}

	// Verify no errors
	env.AssertNoErrors()
}

// TestURLActionExecution tests URL opening workflow
func TestURLActionExecution(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Setup configuration with URL action
	config := `
spells:
  u: test-url

grimoire:
  test-url:
    type: url
    command: "https://example.com"
    description: "Test URL opening"
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

	// Simulate hotkey press for URL action
	if err := env.SimulateHotkey("u"); err != nil {
		t.Fatalf("Failed to simulate hotkey: %v", err)
	}

	// Wait for action to complete
	if err := env.WaitForAction("test-url", 10*time.Second); err != nil {
		t.Fatalf("Action execution timeout: %v", err)
	}

	// Verify URL action execution in logs
	logs, err := env.GetLogs()
	if err != nil {
		t.Fatalf("Failed to get logs: %v", err)
	}

	if !contains(logs, "test-url") || !contains(logs, "example.com") {
		t.Error("URL action execution not found in logs")
	}

	// Verify no errors
	env.AssertNoErrors()
}

// TestSequentialHotkeyWorkflow tests complex hotkey sequences
func TestSequentialHotkeyWorkflow(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Setup configuration with sequential hotkeys
	config := `
spells:
  g,s: git-status
  g,l: git-log
  d,b: debug-build

grimoire:
  git-status:
    type: script
    command: echo "Git status executed"
    description: "Git status check"
    show_output: true
    
  git-log:
    type: script
    command: echo "Git log executed"
    description: "Git log display"
    show_output: true
    
  debug-build:
    type: script
    command: echo "Debug build started"
    description: "Debug build process"
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

	// Test sequential hotkey sequences
	sequences := []struct {
		sequence string
		action   string
		expected string
	}{
		{"g,s", "git-status", "Git status executed"},
		{"g,l", "git-log", "Git log executed"},
		{"d,b", "debug-build", "Debug build started"},
	}

	for _, seq := range sequences {
		t.Run(seq.sequence, func(t *testing.T) {
			// Simulate hotkey sequence
			if err := env.SimulateHotkey(seq.sequence); err != nil {
				t.Fatalf("Failed to simulate hotkey sequence '%s': %v", seq.sequence, err)
			}

			// Wait for action to complete
			if err := env.WaitForAction(seq.action, 10*time.Second); err != nil {
				t.Fatalf("Action '%s' execution timeout: %v", seq.action, err)
			}

			// Verify execution in logs
			logs, err := env.GetLogs()
			if err != nil {
				t.Fatalf("Failed to get logs: %v", err)
			}

			if !contains(logs, seq.expected) {
				t.Errorf("Expected output '%s' not found in logs for sequence '%s'", seq.expected, seq.sequence)
			}
		})
	}

	// Verify no errors
	env.AssertNoErrors()
}

// TestErrorHandlingWorkflow tests error scenarios and recovery
func TestErrorHandlingWorkflow(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Setup configuration with actions that will fail
	config := `
spells:
  f: failing-script
  n: nonexistent-app

grimoire:
  failing-script:
    type: script
    command: exit 1
    description: "Script that fails"
    show_output: true
    
  nonexistent-app:
    type: app
    command: nonexistent-application-xyz
    description: "App that doesn't exist"
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

	// Test failing script action
	if err := env.SimulateHotkey("f"); err != nil {
		t.Fatalf("Failed to simulate hotkey: %v", err)
	}

	// Wait for action attempt
	time.Sleep(3 * time.Second)

	// Test nonexistent app action
	if err := env.SimulateHotkey("n"); err != nil {
		t.Fatalf("Failed to simulate hotkey: %v", err)
	}

	// Wait for action attempt
	time.Sleep(3 * time.Second)

	// Verify error handling in logs
	logs, err := env.GetLogs()
	if err != nil {
		t.Fatalf("Failed to get logs: %v", err)
	}

	// Should contain error information but application should still be running
	expectedErrors := []string{
		"failing-script",
		"nonexistent-app",
	}

	for _, expectedError := range expectedErrors {
		if !contains(logs, expectedError) {
			t.Errorf("Expected error handling for '%s' not found in logs", expectedError)
		}
	}

	// Application should still be responsive (not crashed)
	// Test with a valid action
	validConfig := `
spells:
  v: valid-action

grimoire:
  valid-action:
    type: script
    command: echo "Recovery successful"
    description: "Valid action after errors"
    show_output: true
`

	if err := env.SetupSpellbook(validConfig); err != nil {
		t.Fatalf("Failed to update spellbook: %v", err)
	}

	// Give time for config reload
	time.Sleep(2 * time.Second)

	// Test recovery
	if err := env.SimulateHotkey("v"); err != nil {
		t.Fatalf("Failed to simulate recovery hotkey: %v", err)
	}

	if err := env.WaitForAction("valid-action", 10*time.Second); err != nil {
		t.Fatalf("Recovery action timeout: %v", err)
	}

	// Verify recovery
	logs, err = env.GetLogs()
	if err != nil {
		t.Fatalf("Failed to get logs: %v", err)
	}

	if !contains(logs, "Recovery successful") {
		t.Error("Application failed to recover from errors")
	}
}
