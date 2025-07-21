package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestShowConfigPathCommand(t *testing.T) {
	// Create temp directory for test configs
	tempDir := t.TempDir()

	// Create some test config files
	configPaths := []string{
		filepath.Join(tempDir, "path1"),
		filepath.Join(tempDir, "path2"),
		filepath.Join(tempDir, "path3"),
	}

	// Create config in path2
	os.MkdirAll(configPaths[1], 0o755)
	err := os.WriteFile(filepath.Join(configPaths[1], "spellbook.yml"), []byte("test: config"), 0o644)
	if err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	getConfigPath := func() string { return configPaths[1] }
	getConfigSearchPaths := func() []string { return configPaths }

	tests := []struct {
		name         string
		flags        *Flags
		wantActive   bool
		wantContains []string
	}{
		{
			name: "show config paths",
			flags: &Flags{
				ShowConfigPath: true,
			},
			wantActive: true,
			wantContains: []string{
				"Configuration Search Paths",
				"1. " + configPaths[0],
				"2. " + configPaths[1],
				"3. " + configPaths[2],
				"⭐ (active)",
				"✅ Found: spellbook.yml",
				"❌ No config files found",
				"Configuration Loading Order:",
			},
		},
		{
			name:       "command not active",
			flags:      &Flags{ShowConfigPath: false},
			wantActive: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewShowConfigPathCommand(getConfigPath, getConfigSearchPaths)

			// Test metadata
			if cmd.Name() != "Show Config Path" {
				t.Errorf("Name() = %v, want %v", cmd.Name(), "Show Config Path")
			}

			if cmd.FlagName() != "show-config-path" {
				t.Errorf("FlagName() = %v, want %v", cmd.FlagName(), "show-config-path")
			}

			if cmd.Group() != "config" {
				t.Errorf("Group() = %v, want %v", cmd.Group(), "config")
			}

			if cmd.HasOptions() != false {
				t.Errorf("HasOptions() = %v, want %v", cmd.HasOptions(), false)
			}

			// Test IsActive
			if got := cmd.IsActive(tt.flags); got != tt.wantActive {
				t.Errorf("IsActive() = %v, want %v", got, tt.wantActive)
			}

			// Test Execute if active
			if tt.wantActive {
				// Capture stdout
				old := os.Stdout
				r, w, _ := os.Pipe()
				os.Stdout = w

				err := cmd.Execute(tt.flags)

				w.Close()
				os.Stdout = old

				if err != nil {
					t.Errorf("Execute() error = %v", err)
				}

				// Read captured output
				var buf bytes.Buffer
				buf.ReadFrom(r)
				output := buf.String()

				// Check output contains expected strings
				for _, want := range tt.wantContains {
					if !strings.Contains(output, want) {
						t.Errorf("Execute() output missing %q", want)
					}
				}
			}
		})
	}
}

func TestShowConfigPathCommand_InvalidFlags(t *testing.T) {
	getConfigPath := func() string { return "." }
	getConfigSearchPaths := func() []string { return []string{"."} }
	cmd := NewShowConfigPathCommand(getConfigPath, getConfigSearchPaths)

	// Test with invalid flags type
	if cmd.IsActive("invalid") {
		t.Error("IsActive() should return false for invalid flags type")
	}

	// Test Execute with invalid flags type - this command doesn't check flags type in Execute
	err := cmd.Execute("invalid")
	if err != nil {
		t.Errorf("Execute() should not return error for this command, got %v", err)
	}
}
