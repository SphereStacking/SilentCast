package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/SphereStacking/silentcast/internal/config"
)

func TestValidateConfigCommand(t *testing.T) {
	// Create temp directory for test configs
	tempDir := t.TempDir()

	// Valid config
	validConfig := `
daemon:
  log_level: info
hotkeys:
  prefix: "alt+space"
  timeout: 1000
  sequence_timeout: 2000
spells:
  e: "editor"
grimoire:
  editor:
    type: app
    command: vim
    description: "Open editor"
`

	// Invalid config (missing action)
	invalidConfig := `
daemon:
  log_level: info
hotkeys:
  prefix: "alt+space"
spells:
  e: "missing_action"
  t: "terminal"
grimoire:
  terminal:
    type: app
    command: term
`

	tests := []struct {
		name            string
		configContent   string
		flags           *Flags
		wantActive      bool
		wantErr         bool
		wantErrContains string
	}{
		{
			name:          "valid config",
			configContent: validConfig,
			flags:         &Flags{ValidateConfig: true},
			wantActive:    true,
			wantErr:       false,
		},
		{
			name:            "invalid config - missing action",
			configContent:   invalidConfig,
			flags:           &Flags{ValidateConfig: true},
			wantActive:      true,
			wantErr:         true,
			wantErrContains: "configuration validation failed",
		},
		{
			name:          "command not active",
			configContent: validConfig,
			flags:         &Flags{ValidateConfig: false},
			wantActive:    false,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create config file
			configPath := filepath.Join(tempDir, tt.name)
			os.MkdirAll(configPath, 0o755)
			err := os.WriteFile(filepath.Join(configPath, "spellbook.yml"), []byte(tt.configContent), 0o644)
			if err != nil {
				t.Fatalf("Failed to write config: %v", err)
			}

			getConfigPath := func() string { return configPath }
			cmd := NewValidateConfigCommand(getConfigPath)

			// Test metadata
			if cmd.Name() != "Validate Config" {
				t.Errorf("Name() = %v, want %v", cmd.Name(), "Validate Config")
			}

			if cmd.FlagName() != "validate-config" {
				t.Errorf("FlagName() = %v, want %v", cmd.FlagName(), "validate-config")
			}

			if cmd.Group() != "config" {
				t.Errorf("Group() = %v, want %v", cmd.Group(), "config")
			}

			// Test IsActive
			if got := cmd.IsActive(tt.flags); got != tt.wantActive {
				t.Errorf("IsActive() = %v, want %v", got, tt.wantActive)
			}

			// Test Execute if active
			if tt.wantActive {
				err := cmd.Execute(tt.flags)

				if (err != nil) != tt.wantErr {
					t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				}

				if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("Execute() error = %v, want to contain %v", err.Error(), tt.wantErrContains)
				}
			}
		})
	}
}

func TestValidateConfigCommand_InvalidFlags(t *testing.T) {
	getConfigPath := func() string { return "." }
	cmd := NewValidateConfigCommand(getConfigPath)

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

func TestValidateAction(t *testing.T) {
	tests := []struct {
		name       string
		actionName string
		action     config.ActionConfig
		wantIssues []string
	}{
		{
			name:       "valid app action",
			actionName: "editor",
			action: config.ActionConfig{
				Type:    "app",
				Command: "vim",
			},
			wantIssues: []string{},
		},
		{
			name:       "missing type",
			actionName: "bad",
			action: config.ActionConfig{
				Command: "vim",
			},
			wantIssues: []string{"missing type"},
		},
		{
			name:       "missing command",
			actionName: "bad",
			action: config.ActionConfig{
				Type: "app",
			},
			wantIssues: []string{"missing command"},
		},
		{
			name:       "app with show_output",
			actionName: "bad",
			action: config.ActionConfig{
				Type:       "app",
				Command:    "vim",
				ShowOutput: true,
			},
			wantIssues: []string{"show_output is not applicable for app type"},
		},
		{
			name:       "unknown type",
			actionName: "bad",
			action: config.ActionConfig{
				Type:    "unknown",
				Command: "something",
			},
			wantIssues: []string{"unknown type: unknown"},
		},
		{
			name:       "negative timeout",
			actionName: "bad",
			action: config.ActionConfig{
				Type:    "script",
				Command: "echo",
				Timeout: -1,
			},
			wantIssues: []string{"timeout cannot be negative"},
		},
		{
			name:       "conflicting options",
			actionName: "bad",
			action: config.ActionConfig{
				Type:       "script",
				Command:    "echo",
				ShowOutput: true,
				Terminal:   true,
			},
			wantIssues: []string{"show_output and terminal options may conflict"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues := validateAction(tt.actionName, &tt.action)

			if len(issues) != len(tt.wantIssues) {
				t.Errorf("validateAction() returned %d issues, want %d", len(issues), len(tt.wantIssues))
				t.Errorf("Got: %v", issues)
				t.Errorf("Want: %v", tt.wantIssues)
			}

			// Check each expected issue is present
			for _, want := range tt.wantIssues {
				found := false
				for _, got := range issues {
					if got == want {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("validateAction() missing issue: %s", want)
				}
			}
		})
	}
}
