package config

import (
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestValidator_ValidateHotkeys(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr []string
		noErr   bool
	}{
		{
			name: "valid hotkeys",
			config: Config{
				Hotkeys: HotkeyConfig{
					Prefix:          "alt+space",
					Timeout:         Duration(1000 * time.Millisecond),
					SequenceTimeout: Duration(2000 * time.Millisecond),
				},
				prefixExplicitlySet: true,
			},
			noErr: true,
		},
		{
			name: "missing prefix",
			config: Config{
				Hotkeys: HotkeyConfig{
					Timeout:         Duration(1000 * time.Millisecond),
					SequenceTimeout: Duration(2000 * time.Millisecond),
				},
			},
			wantErr: []string{"prefix key is required"},
		},
		{
			name: "negative timeout",
			config: Config{
				Hotkeys: HotkeyConfig{
					Prefix:          "alt+space",
					Timeout:         Duration(-100 * time.Millisecond),
					SequenceTimeout: Duration(2000 * time.Millisecond),
				},
				prefixExplicitlySet: true,
			},
			wantErr: []string{"timeout must be non-negative"},
		},
		{
			name: "very short timeout warning",
			config: Config{
				Hotkeys: HotkeyConfig{
					Prefix:          "alt+space",
					Timeout:         Duration(50 * time.Millisecond),
					SequenceTimeout: Duration(2000 * time.Millisecond),
				},
				prefixExplicitlySet: true,
			},
			wantErr: []string{"timeout is very short"},
		},
		{
			name: "sequence timeout less than timeout",
			config: Config{
				Hotkeys: HotkeyConfig{
					Prefix:          "alt+space",
					Timeout:         Duration(2000 * time.Millisecond),
					SequenceTimeout: Duration(1000 * time.Millisecond),
				},
				prefixExplicitlySet: true,
			},
			wantErr: []string{"sequence timeout should be greater"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator()
			errors := v.Validate(&tt.config)

			if tt.noErr {
				if len(errors) > 0 {
					t.Errorf("Expected no errors, got: %v", errors)
				}
				return
			}

			for _, wantErr := range tt.wantErr {
				found := false
				for _, err := range errors {
					if strings.Contains(err.Message, wantErr) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error containing '%s' not found in %v", wantErr, errors)
				}
			}
		})
	}
}

func TestValidator_ValidateDaemon(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr string
		noErr   bool
	}{
		{
			name: "valid log level",
			config: Config{
				Daemon: DaemonConfig{
					LogLevel: "info",
				},
			},
			noErr: true,
		},
		{
			name: "invalid log level",
			config: Config{
				Daemon: DaemonConfig{
					LogLevel: "verbose",
				},
			},
			wantErr: "invalid log level 'verbose'",
		},
		{
			name: "empty log level is valid",
			config: Config{
				Daemon: DaemonConfig{
					LogLevel: "",
				},
			},
			noErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator()
			// Set valid hotkeys to avoid that error
			tt.config.Hotkeys.Prefix = "alt+space"
			tt.config.prefixExplicitlySet = true

			errors := v.Validate(&tt.config)

			if tt.noErr {
				if len(errors) > 0 {
					t.Errorf("Expected no errors, got: %v", errors)
				}
				return
			}

			found := false
			for _, err := range errors {
				if strings.Contains(err.Message, tt.wantErr) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected error '%s' not found in %v", tt.wantErr, errors)
			}
		})
	}
}

func TestValidator_ValidateSpells(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr []string
		noErr   bool
	}{
		{
			name: "valid spells",
			config: Config{
				Hotkeys: HotkeyConfig{
					Prefix: "alt+space",
				},
				Shortcuts: map[string]string{
					"e":   "editor",
					"g,s": "git_status",
				},
				Actions: map[string]ActionConfig{
					"editor":     {Type: "app", Command: "vim"},
					"git_status": {Type: "script", Command: "git status"},
				},
				prefixExplicitlySet: true,
			},
			noErr: true,
		},
		{
			name: "spell conflicts with prefix key",
			config: Config{
				Hotkeys: HotkeyConfig{
					Prefix: "alt+space",
				},
				Shortcuts: map[string]string{
					"space": "something",
				},
				Actions: map[string]ActionConfig{
					"something": {Type: "app", Command: "app"},
				},
				prefixExplicitlySet: true,
			},
			wantErr: []string{"conflicts with prefix key"},
		},
		{
			name: "references non-existent action",
			config: Config{
				Hotkeys: HotkeyConfig{
					Prefix: "alt+space",
				},
				Shortcuts: map[string]string{
					"x": "missing_action",
				},
				Actions:             map[string]ActionConfig{},
				prefixExplicitlySet: true,
			},
			wantErr: []string{"references non-existent grimoire action"},
		},
		{
			name: "invalid spell key format",
			config: Config{
				Hotkeys: HotkeyConfig{
					Prefix: "alt+space",
				},
				Shortcuts: map[string]string{
					"": "action",
				},
				Actions: map[string]ActionConfig{
					"action": {Type: "app", Command: "app"},
				},
				prefixExplicitlySet: true,
			},
			wantErr: []string{"empty key"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator()
			errors := v.Validate(&tt.config)

			if tt.noErr {
				if len(errors) > 0 {
					t.Errorf("Expected no errors, got: %v", errors)
				}
				return
			}

			for _, wantErr := range tt.wantErr {
				found := false
				for _, err := range errors {
					if strings.Contains(err.Message, wantErr) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error containing '%s' not found in %v", wantErr, errors)
				}
			}
		})
	}
}

func TestValidator_ValidateGrimoire(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr []string
		noErr   bool
	}{
		{
			name: "valid actions",
			config: Config{
				Hotkeys: HotkeyConfig{
					Prefix: "alt+space",
				},
				Actions: map[string]ActionConfig{
					"app_action": {
						Type:    "app",
						Command: "/bin/sh",
					},
					"script_action": {
						Type:    "script",
						Command: "echo hello",
					},
					"url_action": {
						Type:    "url",
						Command: "https://example.com",
					},
				},
				prefixExplicitlySet: true,
			},
			noErr: true,
		},
		{
			name: "missing type",
			config: Config{
				Hotkeys: HotkeyConfig{
					Prefix: "alt+space",
				},
				Actions: map[string]ActionConfig{
					"bad_action": {
						Command: "something",
					},
				},
				prefixExplicitlySet: true,
			},
			wantErr: []string{"type is required"},
		},
		{
			name: "invalid type",
			config: Config{
				Hotkeys: HotkeyConfig{
					Prefix: "alt+space",
				},
				Actions: map[string]ActionConfig{
					"bad_action": {
						Type:    "invalid",
						Command: "something",
					},
				},
				prefixExplicitlySet: true,
			},
			wantErr: []string{"invalid type 'invalid'"},
		},
		{
			name: "missing command",
			config: Config{
				Hotkeys: HotkeyConfig{
					Prefix: "alt+space",
				},
				Actions: map[string]ActionConfig{
					"bad_action": {
						Type: "app",
					},
				},
				prefixExplicitlySet: true,
			},
			wantErr: []string{"command is required"},
		},
		{
			name: "mutually exclusive options",
			config: Config{
				Hotkeys: HotkeyConfig{
					Prefix: "alt+space",
				},
				Actions: map[string]ActionConfig{
					"conflict": {
						Type:       "script",
						Command:    "echo test",
						ShowOutput: true,
						Terminal:   true,
					},
				},
				prefixExplicitlySet: true,
			},
			wantErr: []string{"show_output and terminal are mutually exclusive"},
		},
		{
			name: "keep_open requires terminal",
			config: Config{
				Hotkeys: HotkeyConfig{
					Prefix: "alt+space",
				},
				Actions: map[string]ActionConfig{
					"bad_keep": {
						Type:     "script",
						Command:  "echo test",
						KeepOpen: true,
						Terminal: false,
					},
				},
				prefixExplicitlySet: true,
			},
			wantErr: []string{"keep_open requires terminal"},
		},
		{
			name: "invalid URL format",
			config: Config{
				Hotkeys: HotkeyConfig{
					Prefix: "alt+space",
				},
				Actions: map[string]ActionConfig{
					"bad_url": {
						Type:    "url",
						Command: "ht!tp://bad url",
					},
				},
				prefixExplicitlySet: true,
			},
			wantErr: []string{"invalid URL format"},
		},
		{
			name: "negative timeout",
			config: Config{
				Hotkeys: HotkeyConfig{
					Prefix: "alt+space",
				},
				Actions: map[string]ActionConfig{
					"bad_timeout": {
						Type:    "script",
						Command: "sleep 10",
						Timeout: -5,
					},
				},
				prefixExplicitlySet: true,
			},
			wantErr: []string{"timeout must be non-negative"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator()
			errors := v.Validate(&tt.config)

			if tt.noErr {
				if len(errors) > 0 {
					t.Errorf("Expected no errors, got: %v", errors)
				}
				return
			}

			for _, wantErr := range tt.wantErr {
				found := false
				for _, err := range errors {
					if strings.Contains(err.Message, wantErr) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error containing '%s' not found in %v", wantErr, errors)
				}
			}
		})
	}
}

func TestValidator_ValidateWithYAML(t *testing.T) {
	yamlContent := `
hotkeys:
  prefix: "alt+space"
  timeout: -100
spells:
  e: "editor"
grimoire:
  editor:
    type: "app"
    command: "vim"
`

	cfg := &Config{
		Hotkeys: HotkeyConfig{
			Prefix:  "alt+space",
			Timeout: Duration(-100 * time.Millisecond),
		},
		Shortcuts: map[string]string{
			"e": "editor",
		},
		Actions: map[string]ActionConfig{
			"editor": {
				Type:    "app",
				Command: "vim",
			},
		},
		prefixExplicitlySet: true,
	}

	v := NewValidator()
	errors := v.ValidateWithYAML(cfg, []byte(yamlContent))

	// Should have error for negative timeout
	found := false
	for _, err := range errors {
		if strings.Contains(err.Message, "timeout must be non-negative") {
			found = true
			// Check if line number is included
			if err.Line == 0 {
				t.Error("Expected line number to be set")
			}
			break
		}
	}

	if !found {
		t.Error("Expected timeout error not found")
	}
}

func TestValidateSpellKey(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		wantErr bool
		errMsg  string
	}{
		{"single key", "e", false, ""},
		{"key with modifier", "ctrl+e", false, ""},
		{"multiple modifiers", "ctrl+shift+e", false, ""},
		{"sequence", "g,s", false, ""},
		{"long sequence", "g,i,t", false, ""},
		{"empty key", "", true, "empty key"},
		{"empty in sequence", "g,,s", true, "empty key in sequence"},
		{"modifier in sequence", "g,ctrl+s", true, "modifiers not allowed"},
		{"invalid modifier", "foo+e", true, "invalid modifier"},
		{"incomplete sequence", "g,", true, "empty key in sequence"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSpellKey(tt.key)
			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.errMsg, err.Error())
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestValidator_LoggerValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr string
		noErr   bool
	}{
		{
			name: "valid logger config",
			config: Config{
				Hotkeys: HotkeyConfig{
					Prefix: "alt+space",
				},
				Logger: LoggerConfig{
					File:       "/tmp/test.log",
					MaxSize:    100,
					MaxBackups: 5,
					MaxAge:     30,
				},
				prefixExplicitlySet: true,
			},
			noErr: true,
		},
		{
			name: "negative max size",
			config: Config{
				Hotkeys: HotkeyConfig{
					Prefix: "alt+space",
				},
				Logger: LoggerConfig{
					MaxSize: -10,
				},
				prefixExplicitlySet: true,
			},
			wantErr: "max size must be non-negative",
		},
		{
			name: "negative max backups",
			config: Config{
				Hotkeys: HotkeyConfig{
					Prefix: "alt+space",
				},
				Logger: LoggerConfig{
					MaxBackups: -1,
				},
				prefixExplicitlySet: true,
			},
			wantErr: "max backups must be non-negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator()
			errors := v.Validate(&tt.config)

			if tt.noErr {
				if len(errors) > 0 {
					t.Errorf("Expected no errors, got: %v", errors)
				}
				return
			}

			found := false
			for _, err := range errors {
				if strings.Contains(err.Message, tt.wantErr) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected error '%s' not found in %v", tt.wantErr, errors)
			}
		})
	}
}

func TestValidator_UpdaterValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr string
		noErr   bool
	}{
		{
			name: "disabled updater skips validation",
			config: Config{
				Hotkeys: HotkeyConfig{
					Prefix: "alt+space",
				},
				Updater: UpdaterConfig{
					Enabled:       false,
					CheckInterval: "invalid",
				},
				prefixExplicitlySet: true,
			},
			noErr: true,
		},
		{
			name: "valid check interval",
			config: Config{
				Hotkeys: HotkeyConfig{
					Prefix: "alt+space",
				},
				Updater: UpdaterConfig{
					Enabled:       true,
					CheckInterval: "24h",
				},
				prefixExplicitlySet: true,
			},
			noErr: true,
		},
		{
			name: "invalid check interval",
			config: Config{
				Hotkeys: HotkeyConfig{
					Prefix: "alt+space",
				},
				Updater: UpdaterConfig{
					Enabled:       true,
					CheckInterval: "24 hours",
				},
				prefixExplicitlySet: true,
			},
			wantErr: "invalid duration format",
		},
		{
			name: "very short check interval",
			config: Config{
				Hotkeys: HotkeyConfig{
					Prefix: "alt+space",
				},
				Updater: UpdaterConfig{
					Enabled:       true,
					CheckInterval: "30m",
				},
				prefixExplicitlySet: true,
			},
			wantErr: "check interval is very short",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator()
			errors := v.Validate(&tt.config)

			if tt.noErr {
				if len(errors) > 0 {
					t.Errorf("Expected no errors, got: %v", errors)
				}
				return
			}

			found := false
			for _, err := range errors {
				if strings.Contains(err.Message, tt.wantErr) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected error '%s' not found in %v", tt.wantErr, errors)
			}
		})
	}
}

func TestValidator_PlatformSpecific(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Platform-specific test for Linux")
	}

	// Test admin validation on Linux without elevation tools
	// This test might pass or fail depending on the system
	config := Config{
		Hotkeys: HotkeyConfig{
			Prefix: "alt+space",
		},
		Actions: map[string]ActionConfig{
			"admin_action": {
				Type:    "script",
				Command: "echo test",
				Admin:   true,
			},
		},
		prefixExplicitlySet: true,
	}

	v := NewValidator()
	errors := v.Validate(&config)

	// Just verify the validation runs without panic
	// The actual result depends on system configuration
	_ = errors
}

func TestValidationError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  ValidationError
		want string
	}{
		{
			name: "with line number",
			err: ValidationError{
				Field:   "hotkeys.timeout",
				Message: "must be positive",
				Line:    10,
			},
			want: "line 10: hotkeys.timeout: must be positive",
		},
		{
			name: "without line number",
			err: ValidationError{
				Field:   "spells.x",
				Message: "invalid key",
			},
			want: "spells.x: invalid key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}
