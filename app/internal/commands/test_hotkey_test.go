package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestTestHotkeyCommand(t *testing.T) {
	// Skip if running in CI or without hotkey support
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping hotkey test in CI environment")
	}

	// Create temp directory for test configs
	tempDir := t.TempDir()

	// Test config
	testConfig := `
hotkeys:
  prefix: "alt+space"
  timeout: 500
  sequence_timeout: 1000
`

	// Create config file
	configPath := filepath.Join(tempDir, "test")
	os.MkdirAll(configPath, 0755)
	err := os.WriteFile(filepath.Join(configPath, "spellbook.yml"), []byte(testConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	getConfigPath := func() string { return configPath }

	tests := []struct {
		name       string
		flags      *Flags
		wantActive bool
	}{
		{
			name: "test hotkey active",
			flags: &Flags{
				TestHotkey:   true,
				TestDuration: 1, // 1 second test
			},
			wantActive: true,
		},
		{
			name: "test hotkey with no duration",
			flags: &Flags{
				TestHotkey:   true,
				TestDuration: 0,
			},
			wantActive: true,
		},
		{
			name:       "command not active",
			flags:      &Flags{TestHotkey: false},
			wantActive: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewTestHotkeyCommand(getConfigPath)

			// Test metadata
			if cmd.Name() != "Test Hotkey" {
				t.Errorf("Name() = %v, want %v", cmd.Name(), "Test Hotkey")
			}

			if cmd.FlagName() != "test-hotkey" {
				t.Errorf("FlagName() = %v, want %v", cmd.FlagName(), "test-hotkey")
			}

			if cmd.Group() != "debug" {
				t.Errorf("Group() = %v, want %v", cmd.Group(), "debug")
			}

			if cmd.HasOptions() != true {
				t.Errorf("HasOptions() = %v, want %v", cmd.HasOptions(), true)
			}

			// Test IsActive
			if got := cmd.IsActive(tt.flags); got != tt.wantActive {
				t.Errorf("IsActive() = %v, want %v", got, tt.wantActive)
			}

			// Test Execute if active and duration is set
			if tt.wantActive && tt.flags.TestDuration > 0 {
				// Run test in goroutine with timeout
				done := make(chan error, 1)
				go func() {
					done <- cmd.Execute(tt.flags)
				}()

				// Wait for completion or timeout
				select {
				case err := <-done:
					if err != nil {
						// This is expected in test environment without proper hotkey support
						if !strings.Contains(err.Error(), "failed to create hotkey manager") &&
							!strings.Contains(err.Error(), "failed to start hotkey manager") {
							t.Errorf("Execute() unexpected error = %v", err)
						}
					}
				case <-time.After(time.Duration(tt.flags.TestDuration+2) * time.Second):
					t.Error("Execute() timed out")
				}
			}
		})
	}
}

func TestTestHotkeyCommand_InvalidFlags(t *testing.T) {
	getConfigPath := func() string { return "." }
	cmd := NewTestHotkeyCommand(getConfigPath)

	// Test with invalid flags type
	if cmd.IsActive("invalid") {
		t.Error("IsActive() should return false for invalid flags type")
	}

	// Test Execute with invalid flags type
	err := cmd.Execute("invalid")
	if err == nil || !strings.Contains(err.Error(), "invalid flags type") {
		t.Errorf("Execute() should return error for invalid flags type, got %v", err)
	}
}

func TestTestHotkeyCommand_NoConfig(t *testing.T) {
	// Use non-existent config path
	getConfigPath := func() string { return "/non/existent/path" }
	cmd := NewTestHotkeyCommand(getConfigPath)

	flags := &Flags{
		TestHotkey:   true,
		TestDuration: 1,
	}

	// This should still work but use default config
	done := make(chan error, 1)
	go func() {
		done <- cmd.Execute(flags)
	}()

	select {
	case err := <-done:
		// Should fail or succeed depending on platform support
		if err != nil && !strings.Contains(err.Error(), "failed to") {
			t.Errorf("Execute() unexpected error = %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Error("Execute() timed out")
	}
}
