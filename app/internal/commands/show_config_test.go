package commands

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestShowConfigCommand(t *testing.T) {
	// Create temp directory for test configs
	tempDir := t.TempDir()

	// Test config
	testConfig := `
daemon:
  log_level: debug
  auto_start: true
hotkeys:
  prefix: "alt+space"
  timeout: 1000
spells:
  e: "editor"
grimoire:
  editor:
    type: app
    command: vim
    description: "Open vim editor"
logger:
  level: info
  max_size: 10
`

	// Create config file
	configPath := filepath.Join(tempDir, "test")
	os.MkdirAll(configPath, 0o755)
	err := os.WriteFile(filepath.Join(configPath, "spellbook.yml"), []byte(testConfig), 0o644)
	if err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	getConfigPath := func() string { return configPath }
	getConfigSearchPaths := func() []string { return []string{configPath} }

	tests := []struct {
		name         string
		flags        *Flags
		wantActive   bool
		wantContains []string
	}{
		{
			name: "show config in human format",
			flags: &Flags{
				ShowConfig: true,
				ShowFormat: "human",
			},
			wantActive: true,
			wantContains: []string{
				"Daemon Settings:",
				"Auto Start: true",
				"Hotkey Settings:",
				"Prefix: alt+space",
				"Spells (1 defined):",
				"e → editor",
			},
		},
		{
			name: "show config in json format",
			flags: &Flags{
				ShowConfig: true,
				ShowFormat: "json",
			},
			wantActive:   true,
			wantContains: []string{`"prefix": "alt+space"`},
		},
		{
			name: "show config in yaml format",
			flags: &Flags{
				ShowConfig: true,
				ShowFormat: "yaml",
			},
			wantActive:   true,
			wantContains: []string{"prefix: alt+space"},
		},
		{
			name: "show config with paths",
			flags: &Flags{
				ShowConfig: true,
				ShowFormat: "human",
				ShowPaths:  true,
			},
			wantActive: true,
			wantContains: []string{
				"Configuration Search Paths:",
				"spellbook.yml",
			},
		},
		{
			name:       "command not active",
			flags:      &Flags{ShowConfig: false},
			wantActive: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewShowConfigCommand(getConfigPath, getConfigSearchPaths)

			// Test metadata
			if cmd.Name() != "Show Config" {
				t.Errorf("Name() = %v, want %v", cmd.Name(), "Show Config")
			}

			if cmd.FlagName() != "show-config" {
				t.Errorf("FlagName() = %v, want %v", cmd.FlagName(), "show-config")
			}

			if cmd.Group() != "config" {
				t.Errorf("Group() = %v, want %v", cmd.Group(), "config")
			}

			if cmd.HasOptions() != true {
				t.Errorf("HasOptions() = %v, want %v", cmd.HasOptions(), true)
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
						t.Errorf("Full output:\n%s", output)
					}
				}

				// Additional format-specific checks
				switch tt.flags.ShowFormat {
				case "json":
					var jsonData map[string]interface{}
					if err := json.Unmarshal([]byte(output), &jsonData); err != nil {
						t.Errorf("Invalid JSON output: %v", err)
					}
				case "yaml":
					var yamlData map[string]interface{}
					if err := yaml.Unmarshal([]byte(output), &yamlData); err != nil {
						t.Errorf("Invalid YAML output: %v", err)
					}
				}
			}
		})
	}
}

func TestShowConfigCommand_InvalidFlags(t *testing.T) {
	getConfigPath := func() string { return "." }
	getConfigSearchPaths := func() []string { return []string{"."} }
	cmd := NewShowConfigCommand(getConfigPath, getConfigSearchPaths)

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
