package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoader_ValidateMethod(t *testing.T) {
	tests := []struct {
		name           string
		config         string
		wantErrors     []string
		wantErrLoading bool
	}{
		{
			name: "valid config",
			config: `
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
`,
			wantErrors: []string{},
		},
		{
			name: "invalid timeout",
			config: `
hotkeys:
  prefix: "alt+space"
  timeout: -1
  sequence_timeout: 2000
`,
			wantErrors: []string{"hotkeys.timeout: timeout must be non-negative"},
		},
		{
			name: "invalid sequence timeout",
			config: `
hotkeys:
  prefix: "alt+space"
  timeout: 1000
  sequence_timeout: -100
`,
			wantErrors: []string{"hotkeys.sequence_timeout: sequence timeout must be non-negative"},
		},
		{
			name: "action missing type",
			config: `
hotkeys:
  prefix: "alt+space"
  timeout: 1000
  sequence_timeout: 2000
grimoire:
  bad_action:
    command: echo
`,
			wantErrors: []string{"grimoire.bad_action.type: type is required"},
		},
		{
			name: "action missing command",
			config: `
hotkeys:
  prefix: "alt+space"
  timeout: 1000
  sequence_timeout: 2000
grimoire:
  bad_action:
    type: script
`,
			wantErrors: []string{"grimoire.bad_action.command: command is required"},
		},
		{
			name: "spell references non-existent action",
			config: `
hotkeys:
  prefix: "alt+space"
  timeout: 1000
  sequence_timeout: 2000
spells:
  x: "missing"
grimoire:
  other:
    type: app
    command: vim
`,
			wantErrors: []string{"spells.x: references non-existent grimoire action 'missing'"},
		},
		{
			name: "multiple errors",
			config: `
hotkeys:
  prefix: "alt+space"
  timeout: -1
spells:
  x: "missing1"
  y: "missing2"
grimoire:
  incomplete:
    type: app
`,
			wantErrors: []string{
				"hotkeys.timeout: timeout must be non-negative",
				"grimoire.incomplete.command: command is required",
				"spells.x: references non-existent grimoire action 'missing1'",
				"spells.y: references non-existent grimoire action 'missing2'",
			},
		},
		{
			name: "invalid yaml",
			config: `
invalid: yaml: with: too: many: colons:
  - unclosed [bracket
  - invalid "quote
`,
			wantErrLoading: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tempDir := t.TempDir()

			// Write config
			configPath := filepath.Join(tempDir, "spellbook.yml")
			err := os.WriteFile(configPath, []byte(tt.config), 0644)
			if err != nil {
				t.Fatalf("Failed to write config: %v", err)
			}

			// Create loader
			loader := NewLoader(tempDir)

			// Validate
			errors, err := loader.Validate()

			if tt.wantErrLoading {
				if err == nil {
					t.Error("Validate() should return error for invalid config")
				}
				return
			}

			if err != nil {
				t.Errorf("Validate() unexpected error: %v", err)
				return
			}

			// Check error count
			if len(errors) != len(tt.wantErrors) {
				t.Errorf("Validate() returned %d errors, want %d", len(errors), len(tt.wantErrors))
				t.Errorf("Got: %v", errors)
				t.Errorf("Want: %v", tt.wantErrors)
			}

			// Check each expected error is present
			for _, want := range tt.wantErrors {
				found := false
				for _, got := range errors {
					if got == want {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Validate() missing expected error: %s", want)
				}
			}
		})
	}
}

func TestLoader_Validate_EmptyConfig(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()

	// Create empty config
	configPath := filepath.Join(tempDir, "spellbook.yml")
	err := os.WriteFile(configPath, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Create loader
	loader := NewLoader(tempDir)

	// Validate - empty config should have required field errors
	errors, err := loader.Validate()
	if err != nil {
		t.Errorf("Validate() unexpected error: %v", err)
	}

	// Should have errors for required fields
	if len(errors) == 0 {
		t.Error("Validate() should return errors for empty config")
	}

	// Check for specific required field errors
	foundPrefixError := false
	for _, e := range errors {
		if e == "hotkeys.prefix: prefix key is required" {
			foundPrefixError = true
			break
		}
	}

	if !foundPrefixError {
		t.Error("Validate() should include error for missing prefix")
		t.Errorf("Got errors: %v", errors)
	}
}

func TestLoader_Validate_NoConfigFile(t *testing.T) {
	// Create temp directory but no config file
	tempDir := t.TempDir()

	// Create loader
	loader := NewLoader(tempDir)

	// Validate - should handle missing file gracefully
	_, err := loader.Validate()
	if err == nil {
		t.Error("Validate() should return error when no config file exists")
	}
}

func TestValidateURLActionType(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()

	// Test URL action type validation
	config := `
daemon:
  log_level: info
hotkeys:
  prefix: "alt+space"
  timeout: 1000
  sequence_timeout: 2000
spells:
  g: "github"
grimoire:
  github:
    type: url
    command: "https://github.com"
    description: "Open GitHub"
`

	configPath := filepath.Join(tempDir, "spellbook.yml")
	err := os.WriteFile(configPath, []byte(config), 0644)
	if err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	loader := NewLoader(tempDir)
	errors, err := loader.Validate()
	if err != nil {
		t.Errorf("URL action should load without error, got: %v", err)
	}

	if len(errors) > 0 {
		t.Errorf("URL action should be valid, got validation errors: %v", errors)
	}
}

func TestValidateInvalidActionType(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()

	// Test invalid action type
	config := `
daemon:
  log_level: info
hotkeys:
  prefix: "alt+space"
  timeout: 1000
  sequence_timeout: 2000
spells:
  x: "invalid_action"
grimoire:
  invalid_action:
    type: invalid
    command: "test"
    description: "Invalid type"
`

	configPath := filepath.Join(tempDir, "spellbook.yml")
	err := os.WriteFile(configPath, []byte(config), 0644)
	if err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	loader := NewLoader(tempDir)
	errors, err := loader.Validate()
	if err != nil {
		t.Errorf("Validate() unexpected error: %v", err)
	}

	if len(errors) == 0 {
		t.Error("Invalid action type should return validation error")
	}

	foundError := false
	expected := "must be 'app', 'script', or 'url'"
	for _, e := range errors {
		if len(e) >= len(expected) {
			for i := 0; i <= len(e)-len(expected); i++ {
				if e[i:i+len(expected)] == expected {
					foundError = true
					break
				}
			}
		}
		if foundError {
			break
		}
	}

	if !foundError {
		t.Errorf("Error should mention valid types, got errors: %v", errors)
	}
}

func TestValidateAllValidActionTypes(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()

	// Test all valid action types
	config := `
daemon:
  log_level: info
hotkeys:
  prefix: "alt+space"
  timeout: 1000
  sequence_timeout: 2000
spells:
  a: "app_action"
  s: "script_action" 
  u: "url_action"
grimoire:
  app_action:
    type: app
    command: "/bin/sh"
    description: "Open shell"
  script_action:
    type: script
    command: "echo hello"
    description: "Echo hello"
  url_action:
    type: url
    command: "https://example.com"
    description: "Open example.com"
`

	configPath := filepath.Join(tempDir, "spellbook.yml")
	err := os.WriteFile(configPath, []byte(config), 0644)
	if err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	loader := NewLoader(tempDir)
	errors, err := loader.Validate()
	if err != nil {
		t.Errorf("Validate() unexpected error: %v", err)
	}

	if len(errors) > 0 {
		t.Errorf("All valid action types should pass validation, got errors: %v", errors)
	}
}
