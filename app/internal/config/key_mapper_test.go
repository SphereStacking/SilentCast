package config

import (
	"testing"
	"gopkg.in/yaml.v3"
)

func TestMapCustomKeys(t *testing.T) {
	// Save original key values
	origDaemon := KeyDaemon
	origHotkeys := KeyHotkeys
	origSpells := KeyShortcuts
	origGrimoire := KeyActions
	origLogger := KeyLogger
	origUpdater := KeyUpdater
	
	// Test cases
	tests := []struct {
		name        string
		setupKeys   func()
		input       string
		wantErr     bool
		checkOutput func(t *testing.T, output []byte)
	}{
		{
			name: "standard keys - no mapping needed",
			setupKeys: func() {
				// Reset to standard keys
				KeyDaemon = "daemon"
				KeyHotkeys = "hotkeys"
				KeyShortcuts = "spells"
				KeyActions = "grimoire"
				KeyLogger = "logger"
				KeyUpdater = "updater"
			},
			input: `
daemon:
  auto_start: true
hotkeys:
  prefix: alt+space
spells:
  test: test-action
`,
			wantErr: false,
			checkOutput: func(t *testing.T, output []byte) {
				// Should return the same content
				var result map[string]interface{}
				if err := yaml.Unmarshal(output, &result); err != nil {
					t.Fatalf("Failed to unmarshal output: %v", err)
				}
				if _, ok := result["daemon"]; !ok {
					t.Error("Expected 'daemon' key in output")
				}
			},
		},
		{
			name: "Japanese keys",
			setupKeys: func() {
				KeyDaemon = "デーモン"
				KeyHotkeys = "ホットキー"
				KeyShortcuts = "呪文"
				KeyActions = "魔導書"
				KeyLogger = "ログ"
				KeyUpdater = "更新"
			},
			input: `
デーモン:
  auto_start: true
  log_level: info
ホットキー:
  prefix: alt+space
  timeout: 1000
呪文:
  gs: git-status
魔導書:
  git-status:
    type: script
    command: git status
`,
			wantErr: false,
			checkOutput: func(t *testing.T, output []byte) {
				var result map[string]interface{}
				if err := yaml.Unmarshal(output, &result); err != nil {
					t.Fatalf("Failed to unmarshal output: %v", err)
				}
				
				// Check that Japanese keys were mapped to standard keys
				if _, ok := result["daemon"]; !ok {
					t.Error("Expected 'daemon' key in output")
				}
				if _, ok := result["hotkeys"]; !ok {
					t.Error("Expected 'hotkeys' key in output")
				}
				if _, ok := result["spells"]; !ok {
					t.Error("Expected 'spells' key in output")
				}
				if _, ok := result["grimoire"]; !ok {
					t.Error("Expected 'grimoire' key in output")
				}
				
				// Check that Japanese keys are NOT in output
				if _, ok := result["デーモン"]; ok {
					t.Error("Japanese key 'デーモン' should not be in output")
				}
			},
		},
		{
			name: "corporate naming",
			setupKeys: func() {
				KeyDaemon = "service"
				KeyHotkeys = "shortcuts"
				KeyShortcuts = "mappings"
				KeyActions = "commands"
				KeyLogger = "logging"
				KeyUpdater = "updates"
			},
			input: `
service:
  auto_start: false
shortcuts:
  prefix: ctrl+shift+x
mappings:
  open: open-app
commands:
  open-app:
    type: app
    command: /usr/bin/app
logging:
  level: debug
updates:
  enabled: true
`,
			wantErr: false,
			checkOutput: func(t *testing.T, output []byte) {
				var result map[string]interface{}
				if err := yaml.Unmarshal(output, &result); err != nil {
					t.Fatalf("Failed to unmarshal output: %v", err)
				}
				
				// All custom keys should be mapped to standard keys
				expectedKeys := []string{"daemon", "hotkeys", "spells", "grimoire", "logger", "updater"}
				for _, key := range expectedKeys {
					if _, ok := result[key]; !ok {
						t.Errorf("Expected '%s' key in output", key)
					}
				}
				
				// Custom keys should not exist in output
				customKeys := []string{"service", "shortcuts", "mappings", "commands", "logging", "updates"}
				for _, key := range customKeys {
					if _, ok := result[key]; ok {
						t.Errorf("Custom key '%s' should not be in output", key)
					}
				}
			},
		},
		{
			name: "mixed standard and custom keys",
			setupKeys: func() {
				KeyDaemon = "service"
				KeyHotkeys = "hotkeys" // standard
				KeyShortcuts = "呪文"
				KeyActions = "grimoire" // standard
			},
			input: `
service:
  auto_start: true
hotkeys:
  prefix: alt+x
呪文:
  test: test-cmd
grimoire:
  test-cmd:
    type: script
    command: echo test
`,
			wantErr: false,
			checkOutput: func(t *testing.T, output []byte) {
				var result map[string]interface{}
				if err := yaml.Unmarshal(output, &result); err != nil {
					t.Fatalf("Failed to unmarshal output: %v", err)
				}
				
				// All should be mapped to standard keys
				if _, ok := result["daemon"]; !ok {
					t.Error("Expected 'daemon' key in output")
				}
				if _, ok := result["hotkeys"]; !ok {
					t.Error("Expected 'hotkeys' key in output")
				}
				if _, ok := result["spells"]; !ok {
					t.Error("Expected 'spells' key in output")
				}
			},
		},
		{
			name: "invalid YAML",
			setupKeys: func() {
				KeyDaemon = "daemon"
			},
			input:   `invalid: yaml: content`,
			wantErr: true,
		},
		{
			name: "empty input",
			setupKeys: func() {
				KeyDaemon = "daemon"
			},
			input:   ``,
			wantErr: false,
			checkOutput: func(t *testing.T, output []byte) {
				var result map[string]interface{}
				if err := yaml.Unmarshal(output, &result); err != nil {
					t.Fatalf("Failed to unmarshal output: %v", err)
				}
				if len(result) != 0 {
					t.Error("Expected empty map for empty input")
				}
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup custom keys
			tt.setupKeys()
			
			// Run the mapping
			output, _, err := MapCustomKeys([]byte(tt.input))
			
			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("MapCustomKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			// Check output if no error expected
			if !tt.wantErr && tt.checkOutput != nil {
				tt.checkOutput(t, output)
			}
		})
	}
	
	// Restore original keys
	KeyDaemon = origDaemon
	KeyHotkeys = origHotkeys
	KeyShortcuts = origSpells
	KeyActions = origGrimoire
	KeyLogger = origLogger
	KeyUpdater = origUpdater
}