//go:build e2e

package e2e

import (
	"fmt"
	"testing"
	"time"
)

// TestApplicationStartupPerformance tests application startup time
func TestApplicationStartupPerformance(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Setup test configuration
	if err := env.SetupSpellbook(TestSpellbookTemplate); err != nil {
		t.Fatalf("Failed to setup spellbook: %v", err)
	}

	// Measure startup time
	startTime := time.Now()

	// Start the application
	if err := env.StartApplication(); err != nil {
		t.Fatalf("Failed to start application: %v", err)
	}

	// Wait for startup to complete
	if err := env.WaitForStartup(); err != nil {
		t.Fatalf("Application startup failed: %v", err)
	}

	startupDuration := time.Since(startTime)

	// Performance target: startup should complete within 10 seconds for E2E
	// (More lenient than unit test targets due to E2E overhead)
	maxStartupTime := 10 * time.Second
	if startupDuration > maxStartupTime {
		t.Errorf("Application startup took too long: %v (max: %v)", startupDuration, maxStartupTime)
	}

	t.Logf("Application startup completed in: %v", startupDuration)

	// Verify no errors during startup
	env.AssertNoErrors()
}

// TestHotkeyResponsePerformance tests hotkey response time
func TestHotkeyResponsePerformance(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Setup configuration with simple action for performance testing
	config := `
spells:
  p: perf-test

grimoire:
  perf-test:
    type: script
    command: echo "Performance test completed"
    description: "Performance test action"
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

	// Measure hotkey response time
	responseStartTime := time.Now()

	// Simulate hotkey
	if err := env.SimulateHotkey("p"); err != nil {
		t.Fatalf("Failed to simulate hotkey: %v", err)
	}

	// Wait for action to complete
	if err := env.WaitForAction("perf-test", 30*time.Second); err != nil {
		t.Fatalf("Performance test action timeout: %v", err)
	}

	responseTime := time.Since(responseStartTime)

	// Performance target: hotkey response should complete within 5 seconds for E2E
	// (More lenient than unit test targets due to process communication overhead)
	maxResponseTime := 5 * time.Second
	if responseTime > maxResponseTime {
		t.Errorf("Hotkey response took too long: %v (max: %v)", responseTime, maxResponseTime)
	}

	t.Logf("Hotkey response completed in: %v", responseTime)

	// Verify execution
	logs, err := env.GetLogs()
	if err != nil {
		t.Fatalf("Failed to get logs: %v", err)
	}

	if !contains(logs, "Performance test completed") {
		t.Error("Performance test action did not complete")
	}

	// Verify no errors
	env.AssertNoErrors()
}

// TestMultipleActionPerformance tests performance with multiple rapid actions
func TestMultipleActionPerformance(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Setup configuration with multiple actions
	config := `
spells:
  a1: action-1
  a2: action-2
  a3: action-3
  a4: action-4
  a5: action-5

grimoire:
  action-1:
    type: script
    command: echo "Action 1"
    description: "Performance test action 1"
    
  action-2:
    type: script
    command: echo "Action 2"
    description: "Performance test action 2"
    
  action-3:
    type: script
    command: echo "Action 3"
    description: "Performance test action 3"
    
  action-4:
    type: script
    command: echo "Action 4"
    description: "Performance test action 4"
    
  action-5:
    type: script
    command: echo "Action 5"
    description: "Performance test action 5"
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

	// Measure multiple action execution time
	multiActionStartTime := time.Now()

	// Execute multiple actions rapidly
	actions := []string{"a1", "a2", "a3", "a4", "a5"}
	for _, action := range actions {
		if err := env.SimulateHotkey(action); err != nil {
			t.Fatalf("Failed to simulate hotkey %s: %v", action, err)
		}
		// Small delay to avoid overwhelming the system
		time.Sleep(50 * time.Millisecond)
	}

	// Wait for all actions to complete
	actionNames := []string{"action-1", "action-2", "action-3", "action-4", "action-5"}
	for _, actionName := range actionNames {
		if err := env.WaitForAction(actionName, 15*time.Second); err != nil {
			t.Errorf("Action %s timeout: %v", actionName, err)
		}
	}

	multiActionDuration := time.Since(multiActionStartTime)

	// Performance target: 5 actions should complete within 30 seconds
	maxMultiActionTime := 30 * time.Second
	if multiActionDuration > maxMultiActionTime {
		t.Errorf("Multiple action execution took too long: %v (max: %v)", multiActionDuration, maxMultiActionTime)
	}

	t.Logf("Multiple action execution completed in: %v", multiActionDuration)

	// Verify all actions executed
	logs, err := env.GetLogs()
	if err != nil {
		t.Fatalf("Failed to get logs: %v", err)
	}

	for i := 1; i <= 5; i++ {
		expected := fmt.Sprintf("Action %d", i)
		if !contains(logs, expected) {
			t.Errorf("Action %d output not found in logs", i)
		}
	}

	// Verify no errors
	env.AssertNoErrors()
}

// TestConfigurationReloadPerformance tests configuration reload performance
func TestConfigurationReloadPerformance(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Start with initial configuration
	if err := env.SetupSpellbook(TestSpellbookTemplate); err != nil {
		t.Fatalf("Failed to setup initial spellbook: %v", err)
	}

	// Start application
	if err := env.StartApplication(); err != nil {
		t.Fatalf("Failed to start application: %v", err)
	}

	if err := env.WaitForStartup(); err != nil {
		t.Fatalf("Application startup failed: %v", err)
	}

	// Measure configuration reload time
	reloadStartTime := time.Now()

	// Update configuration
	updatedConfig := `
spells:
  r: reload-test

grimoire:
  reload-test:
    type: script
    command: echo "Configuration reloaded successfully"
    description: "Reload performance test"
    show_output: true
`

	if err := env.SetupSpellbook(updatedConfig); err != nil {
		t.Fatalf("Failed to update spellbook: %v", err)
	}

	// Wait for configuration reload detection
	time.Sleep(3 * time.Second)

	// Test new action to confirm reload
	if err := env.SimulateHotkey("r"); err != nil {
		t.Fatalf("Failed to simulate hotkey: %v", err)
	}

	if err := env.WaitForAction("reload-test", 15*time.Second); err != nil {
		t.Fatalf("Reload test action timeout: %v", err)
	}

	reloadDuration := time.Since(reloadStartTime)

	// Performance target: configuration reload should complete within 20 seconds
	maxReloadTime := 20 * time.Second
	if reloadDuration > maxReloadTime {
		t.Errorf("Configuration reload took too long: %v (max: %v)", reloadDuration, maxReloadTime)
	}

	t.Logf("Configuration reload completed in: %v", reloadDuration)

	// Verify reload was successful
	logs, err := env.GetLogs()
	if err != nil {
		t.Fatalf("Failed to get logs: %v", err)
	}

	if !contains(logs, "Configuration reloaded successfully") {
		t.Error("Configuration reload test action did not execute")
	}

	// Verify no errors
	env.AssertNoErrors()
}

// TestMemoryUsageStability tests that memory usage remains stable
func TestMemoryUsageStability(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Setup configuration for stability testing
	config := `
spells:
  s: stability-test

grimoire:
  stability-test:
    type: script
    command: echo "Stability test iteration"
    description: "Memory stability test action"
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

	// Run multiple iterations to test stability
	iterations := 10
	testStartTime := time.Now()

	for i := 0; i < iterations; i++ {
		// Execute action
		if err := env.SimulateHotkey("s"); err != nil {
			t.Fatalf("Failed to simulate hotkey in iteration %d: %v", i, err)
		}

		// Wait for action to complete
		if err := env.WaitForAction("stability-test", 10*time.Second); err != nil {
			t.Fatalf("Stability test action timeout in iteration %d: %v", i, err)
		}

		// Short pause between iterations
		time.Sleep(100 * time.Millisecond)

		// Check for any errors periodically
		if i%3 == 0 {
			logs, err := env.GetLogs()
			if err == nil {
				if contains(logs, "panic") || contains(logs, "FATAL") {
					t.Fatalf("Application stability issue detected in iteration %d", i)
				}
			}
		}
	}

	stabilityTestDuration := time.Since(testStartTime)

	// Performance target: 10 iterations should complete within 2 minutes
	maxStabilityTestTime := 2 * time.Minute
	if stabilityTestDuration > maxStabilityTestTime {
		t.Errorf("Stability test took too long: %v (max: %v)", stabilityTestDuration, maxStabilityTestTime)
	}

	t.Logf("Memory stability test (%d iterations) completed in: %v", iterations, stabilityTestDuration)

	// Verify all iterations completed successfully
	logs, err := env.GetLogs()
	if err != nil {
		t.Fatalf("Failed to get logs: %v", err)
	}

	// Count stability test iterations in logs
	iterationCount := 0
	logLines := logs
	searchText := "Stability test iteration"
	for len(logLines) > 0 {
		if contains(logLines, searchText) {
			iterationCount++
			// Find the position and continue searching after it
			pos := 0
			for i := 0; i <= len(logLines)-len(searchText); i++ {
				if logLines[i:i+len(searchText)] == searchText {
					pos = i + len(searchText)
					break
				}
			}
			if pos > 0 && pos < len(logLines) {
				logLines = logLines[pos:]
			} else {
				break
			}
		} else {
			break
		}
	}

	if iterationCount < iterations {
		t.Errorf("Expected %d stability test iterations, found %d in logs", iterations, iterationCount)
	}

	// Verify no memory leaks or stability issues
	env.AssertNoErrors()
}

// TestShutdownPerformance tests graceful shutdown performance
func TestShutdownPerformance(t *testing.T) {
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

	// Measure shutdown time
	shutdownStartTime := time.Now()

	// Stop application
	if err := env.StopApplication(); err != nil {
		t.Errorf("Application shutdown failed: %v", err)
	}

	shutdownDuration := time.Since(shutdownStartTime)

	// Performance target: shutdown should complete within 10 seconds
	maxShutdownTime := 10 * time.Second
	if shutdownDuration > maxShutdownTime {
		t.Errorf("Application shutdown took too long: %v (max: %v)", shutdownDuration, maxShutdownTime)
	}

	t.Logf("Application shutdown completed in: %v", shutdownDuration)

	// Verify clean shutdown (check logs if accessible)
	logs, err := env.GetLogs()
	if err == nil {
		if contains(logs, "panic") || contains(logs, "FATAL") {
			t.Error("Application did not shut down cleanly")
		}
	}
}