//go:build e2e

package e2e

import (
	"testing"
	"time"
)

// TestApplicationStartup tests the complete application startup workflow
func TestApplicationStartup(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()
	
	// Setup test configuration
	if err := env.SetupSpellbook(TestSpellbookTemplate); err != nil {
		t.Fatalf("Failed to setup spellbook: %v", err)
	}
	
	// Start the application
	if err := env.StartApplication(); err != nil {
		t.Fatalf("Failed to start application: %v", err)
	}
	
	// Wait for startup to complete
	if err := env.WaitForStartup(); err != nil {
		t.Fatalf("Application startup failed: %v", err)
	}
	
	// Verify no errors occurred during startup
	env.AssertNoErrors()
	
	// Verify logs contain expected startup messages
	logs, err := env.GetLogs()
	if err != nil {
		t.Fatalf("Failed to get logs: %v", err)
	}
	
	expectedMessages := []string{
		"Configuration loaded",
		"Hotkey manager",
	}
	
	for _, msg := range expectedMessages {
		if !contains(logs, msg) {
			t.Errorf("Expected message '%s' not found in logs", msg)
		}
	}
}

// TestApplicationShutdown tests graceful application shutdown
func TestApplicationShutdown(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()
	
	// Setup and start application
	if err := env.SetupSpellbook(TestSpellbookTemplate); err != nil {
		t.Fatalf("Failed to setup spellbook: %v", err)
	}
	
	if err := env.StartApplication(); err != nil {
		t.Fatalf("Failed to start application: %v", err)
	}
	
	if err := env.WaitForStartup(); err != nil {
		t.Fatalf("Application startup failed: %v", err)
	}
	
	// Test graceful shutdown
	if err := env.StopApplication(); err != nil {
		t.Errorf("Application shutdown failed: %v", err)
	}
	
	// Verify clean shutdown (no error in final logs)
	logs, err := env.GetLogs()
	if err == nil {
		if contains(logs, "panic") || contains(logs, "FATAL") {
			t.Error("Application did not shut down gracefully")
		}
	}
}

// TestConfigurationLoading tests configuration file loading and validation
func TestConfigurationLoading(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()
	
	tests := []struct {
		name           string
		config         string
		expectStartup  bool
		expectError    string
	}{
		{
			name:          "valid configuration",
			config:        TestSpellbookTemplate,
			expectStartup: true,
		},
		{
			name: "invalid YAML",
			config: `
spells:
  e: editor
grimoire:
  editor:
    type: app
    command: echo
    invalid_yaml: [unclosed
`,
			expectStartup: false,
			expectError:   "yaml",
		},
		{
			name: "missing grimoire entry",
			config: `
spells:
  e: nonexistent
grimoire:
  editor:
    type: app
    command: echo
`,
			expectStartup: true, // Should start but log warning
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup configuration
			if err := env.SetupSpellbook(tt.config); err != nil {
				t.Fatalf("Failed to setup spellbook: %v", err)
			}
			
			// Start application
			startErr := env.StartApplication()
			
			if tt.expectStartup {
				if startErr != nil {
					t.Fatalf("Expected successful startup, got error: %v", startErr)
				}
				
				// Wait for startup
				if err := env.WaitForStartup(); err != nil {
					t.Fatalf("Startup timeout: %v", err)
				}
				
				// Check for expected errors in logs if any
				if tt.expectError != "" {
					logs, _ := env.GetLogs()
					if !contains(logs, tt.expectError) {
						t.Errorf("Expected error '%s' not found in logs", tt.expectError)
					}
				}
			} else {
				if startErr == nil {
					// Application started, wait a bit then check if it exits
					time.Sleep(2 * time.Second)
				}
				
				// Application should either fail to start or exit quickly
				logs, _ := env.GetLogs()
				if tt.expectError != "" && !contains(logs, tt.expectError) {
					t.Errorf("Expected error '%s' not found in logs", tt.expectError)
				}
			}
		})
	}
}

// TestConfigurationReloading tests configuration file change detection and reloading
func TestConfigurationReloading(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()
	
	// Start with initial configuration
	if err := env.SetupSpellbook(TestSpellbookTemplate); err != nil {
		t.Fatalf("Failed to setup initial spellbook: %v", err)
	}
	
	if err := env.StartApplication(); err != nil {
		t.Fatalf("Failed to start application: %v", err)
	}
	
	if err := env.WaitForStartup(); err != nil {
		t.Fatalf("Application startup failed: %v", err)
	}
	
	// Update configuration
	updatedConfig := `
spells:
  e: editor
  n: new-action

grimoire:
  editor:
    type: app
    command: echo
    description: "Updated editor action"
    
  new-action:
    type: script
    command: echo "New action executed"
    description: "Newly added action"
`
	
	if err := env.SetupSpellbook(updatedConfig); err != nil {
		t.Fatalf("Failed to update spellbook: %v", err)
	}
	
	// Wait for configuration reload
	time.Sleep(2 * time.Second)
	
	// Check logs for reload confirmation
	logs, err := env.GetLogs()
	if err != nil {
		t.Fatalf("Failed to get logs: %v", err)
	}
	
	reloadIndicators := []string{
		"Configuration reloaded",
		"Config file changed",
		"Reloading configuration",
	}
	
	found := false
	for _, indicator := range reloadIndicators {
		if contains(logs, indicator) {
			found = true
			break
		}
	}
	
	if !found {
		t.Error("Configuration reload not detected in logs")
	}
}