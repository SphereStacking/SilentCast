package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestLoader_Load(t *testing.T) {
	// Create a temporary directory for test configs
	tempDir := t.TempDir()

	// Create test configuration files
	commonConfig := `
daemon:
  log_level: debug
  config_watch: true

hotkeys:
  prefix: "alt+space"
  timeout: 500

spells:
  e: "editor"
  t: "terminal"
  "g,s": "git_status"

grimoire:
  editor:
    type: app
    command: /bin/sh
    description: "Launch shell"
  terminal:
    type: app
    command: /bin/echo
  git_status:
    type: script
    command: "git status"
`

	osSpecificConfig := `
hotkeys:
  prefix: "cmd+space"

spells:
  e: "vscode"  # Override editor

grimoire:
  vscode:
    type: app
    command: /bin/cat
    description: "Test command"
`

	// Write common config
	err := os.WriteFile(filepath.Join(tempDir, "spellbook.yml"), []byte(commonConfig), 0o600)
	if err != nil {
		t.Fatalf("Failed to write common config: %v", err)
	}

	// Write OS-specific config based on current OS
	var osConfigFile string
	switch runtime.GOOS {
	case "darwin":
		osConfigFile = "spellbook.mac.yml"
	case "windows":
		osConfigFile = "spellbook.windows.yml"
	case "linux":
		osConfigFile = "spellbook.linux.yml"
	default:
		// For other platforms, skip OS-specific config
		osConfigFile = ""
	}

	if osConfigFile != "" {
		err = os.WriteFile(filepath.Join(tempDir, osConfigFile), []byte(osSpecificConfig), 0o600)
		if err != nil {
			t.Fatalf("Failed to write OS config: %v", err)
		}
	}

	// Test loading
	loader := NewLoader(tempDir)
	cfg, err := loader.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify merged configuration
	tests := []struct {
		name     string
		check    func() bool
		expected bool
	}{
		{
			name: "OS-specific prefix overrides common",
			check: func() bool {
				// Only expect override if OS-specific config was written
				if osConfigFile != "" {
					return cfg.Hotkeys.Prefix == "cmd+space"
				}
				return cfg.Hotkeys.Prefix == "alt+space"
			},
			expected: true,
		},
		{
			name:     "Common timeout is preserved",
			check:    func() bool { return cfg.Hotkeys.Timeout.ToDuration() == 500*time.Millisecond },
			expected: true,
		},
		{
			name: "OS-specific spell overrides common",
			check: func() bool {
				// Only expect override if OS-specific config was written
				if osConfigFile != "" {
					return cfg.Shortcuts["e"] == "vscode"
				}
				return cfg.Shortcuts["e"] == "editor"
			},
			expected: true,
		},
		{
			name:     "Common spell is preserved",
			check:    func() bool { return cfg.Shortcuts["t"] == "terminal" },
			expected: true,
		},
		{
			name:     "Multi-key spell is preserved",
			check:    func() bool { return cfg.Shortcuts["g,s"] == "git_status" },
			expected: true,
		},
		{
			name: "New grimoire entry is added",
			check: func() bool {
				// Only expect vscode entry if OS-specific config was written
				if osConfigFile != "" {
					_, exists := cfg.Actions["vscode"]
					return exists
				}
				// Otherwise, editor entry should exist
				_, exists := cfg.Actions["editor"]
				return exists
			},
			expected: true,
		},
		{
			name:     "Common grimoire entry is preserved",
			check:    func() bool { _, exists := cfg.Actions["terminal"]; return exists },
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.check(); got != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestActionConfig_NewFields(t *testing.T) {
	tests := []struct {
		name   string
		config string
		check  func(*Config) bool
	}{
		{
			name: "show_output field",
			config: `
spells:
  g: "git_status"

grimoire:
  git_status:
    type: script
    command: "git status"
    show_output: true
`,
			check: func(cfg *Config) bool {
				return cfg.Actions["git_status"].ShowOutput == true
			},
		},
		{
			name: "keep_open field",
			config: `
spells:
  l: "log_tail"

grimoire:
  log_tail:
    type: script
    command: "tail -f /var/log/system.log"
    keep_open: true
    terminal: true
`,
			check: func(cfg *Config) bool {
				return cfg.Actions["log_tail"].KeepOpen == true
			},
		},
		{
			name: "timeout field",
			config: `
spells:
  s: "slow_script"

grimoire:
  slow_script:
    type: script
    command: "sleep 100"
    timeout: 30
`,
			check: func(cfg *Config) bool {
				return cfg.Actions["slow_script"].Timeout == 30
			},
		},
		{
			name: "shell field",
			config: `
spells:
  p: "powershell"

grimoire:
  powershell:
    type: script
    command: "Get-Process"
    shell: "sh"
`,
			check: func(cfg *Config) bool {
				return cfg.Actions["powershell"].Shell == "sh"
			},
		},
		{
			name: "admin field",
			config: `
spells:
  a: "admin_task"

grimoire:
  admin_task:
    type: script
    command: "netstat -an"
    admin: true
`,
			check: func(cfg *Config) bool {
				return cfg.Actions["admin_task"].Admin == true
			},
		},
		{
			name: "terminal field",
			config: `
spells:
  i: "interactive"

grimoire:
  interactive:
    type: script
    command: "python"
    terminal: true
`,
			check: func(cfg *Config) bool {
				return cfg.Actions["interactive"].Terminal == true
			},
		},
		{
			name: "all fields combined",
			config: `
spells:
  f: "full_action"

grimoire:
  full_action:
    type: script
    command: "complex_script.sh"
    args: ["--verbose"]
    env:
      DEBUG: "true"
    working_dir: "/tmp"
    keep_open: true
    timeout: 60
    shell: "sh"
    admin: true
    terminal: true
`,
			check: func(cfg *Config) bool {
				action := cfg.Actions["full_action"]
				return action.ShowOutput == false &&
					action.KeepOpen == true &&
					action.Timeout == 60 &&
					action.Shell == "sh" &&
					action.Admin == true &&
					action.Terminal == true &&
					len(action.Args) == 1 &&
					action.Args[0] == "--verbose" &&
					action.Env["DEBUG"] == "true" &&
					action.WorkingDir == "/tmp"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			err := os.WriteFile(filepath.Join(tempDir, ConfigName+".yml"), []byte(tt.config), 0o600)
			if err != nil {
				t.Fatalf("Failed to write config: %v", err)
			}

			loader := NewLoader(tempDir)
			cfg, err := loader.Load()
			if err != nil {
				t.Fatalf("Failed to load config: %v", err)
			}

			if !tt.check(cfg) {
				t.Errorf("Field check failed for %s", tt.name)
			}
		})
	}
}

func TestLoader_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  string
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid configuration",
			config: `
hotkeys:
  prefix: "alt+space"

spells:
  e: "editor"

grimoire:
  editor:
    type: app
    command: /usr/bin/vim
`,
			wantErr: false,
		},
		{
			name: "Invalid spell reference",
			config: `
hotkeys:
  prefix: "alt+space"

spells:
  e: "nonexistent"

grimoire:
  editor:
    type: app
    command: /usr/bin/vim
`,
			wantErr: true,
			errMsg:  "non-existent grimoire action",
		},
		{
			name: "Invalid action type",
			config: `
hotkeys:
  prefix: "alt+space"

spells:
  e: "editor"

grimoire:
  editor:
    type: invalid
    command: /usr/bin/vim
`,
			wantErr: true,
			errMsg:  "invalid type",
		},
		{
			name: "Missing command",
			config: `
hotkeys:
  prefix: "alt+space"

spells:
  e: "editor"

grimoire:
  editor:
    type: app
`,
			wantErr: true,
			errMsg:  "command is required",
		},
		{
			name: "Invalid log level",
			config: `
daemon:
  log_level: invalid

hotkeys:
  prefix: "alt+space"

spells:
  e: "editor"

grimoire:
  editor:
    type: app
    command: /usr/bin/vim
`,
			wantErr: true,
			errMsg:  "invalid log level",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			err := os.WriteFile(filepath.Join(tempDir, ConfigName+".yml"), []byte(tt.config), 0o600)
			if err != nil {
				t.Fatalf("Failed to write config: %v", err)
			}

			loader := NewLoader(tempDir)
			_, err = loader.Load()

			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.errMsg, err.Error())
				}
			}
		})
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && (s[0:len(substr)] == substr || contains(s[1:], substr)))
}
