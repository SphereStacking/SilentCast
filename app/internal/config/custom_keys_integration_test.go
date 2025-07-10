package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoaderWithCustomKeys(t *testing.T) {
	// Save original key values
	origDaemon := KeyDaemon
	origHotkeys := KeyHotkeys
	origSpells := KeyShortcuts
	origGrimoire := KeyActions
	defer func() {
		// Restore original keys
		KeyDaemon = origDaemon
		KeyHotkeys = origHotkeys
		KeyShortcuts = origSpells
		KeyActions = origGrimoire
	}()
	
	// Create temp directory for test configs
	
	tests := []struct {
		name      string
		setupKeys func()
		config    string
		filename  string
		validate  func(t *testing.T, cfg *Config)
	}{
		{
			name: "Japanese keys full config",
			setupKeys: func() {
				KeyDaemon = "デーモン"
				KeyHotkeys = "ホットキー"
				KeyShortcuts = "呪文"
				KeyActions = "魔導書"
			},
			config: `
デーモン:
  auto_start: true
  log_level: debug
  config_watch: true

ホットキー:
  prefix: ctrl+alt+space
  timeout: 2000
  sequence_timeout: 3000

呪文:
  gs: git-status
  gb: git-branch
  gc: git-commit

魔導書:
  git-status:
    type: script
    command: git status
    description: Show git status
  git-branch:
    type: script
    command: git branch -a
    args: ["--color=always"]
  git-commit:
    type: script
    command: git commit
    working_dir: /tmp
`,
			filename: ConfigName + ".yml",
			validate: func(t *testing.T, cfg *Config) {
				// Check daemon config
				if !cfg.Daemon.AutoStart {
					t.Error("Expected daemon.auto_start to be true")
				}
				if cfg.Daemon.LogLevel != "debug" {
					t.Errorf("Expected log_level 'debug', got '%s'", cfg.Daemon.LogLevel)
				}
				
				// Check hotkeys config
				if cfg.Hotkeys.Prefix != "ctrl+alt+space" {
					t.Errorf("Expected prefix 'ctrl+alt+space', got '%s'", cfg.Hotkeys.Prefix)
				}
				if cfg.Hotkeys.Timeout.ToDuration().Milliseconds() != 2000 {
					t.Errorf("Expected timeout 2000ms, got %d", cfg.Hotkeys.Timeout.ToDuration().Milliseconds())
				}
				
				// Check spells
				if len(cfg.Shortcuts) != 3 {
					t.Errorf("Expected 3 spells, got %d", len(cfg.Shortcuts))
				}
				if cfg.Shortcuts["gs"] != "git-status" {
					t.Error("Expected spell 'gs' to map to 'git-status'")
				}
				
				// Check grimoire
				if len(cfg.Actions) != 3 {
					t.Errorf("Expected 3 grimoire entries, got %d", len(cfg.Actions))
				}
				if action, ok := cfg.Actions["git-status"]; ok {
					if action.Type != "script" {
						t.Errorf("Expected type 'script', got '%s'", action.Type)
					}
					if action.Description != "Show git status" {
						t.Errorf("Expected description 'Show git status', got '%s'", action.Description)
					}
				} else {
					t.Error("Expected 'git-status' in grimoire")
				}
			},
		},
		{
			name: "Corporate naming simple",
			setupKeys: func() {
				KeyDaemon = "service"
				KeyHotkeys = "shortcuts"
				KeyShortcuts = "mappings"
				KeyActions = "commands"
			},
			config: `
service:
  auto_start: false
  log_level: info

shortcuts:
  prefix: "win+shift+x"
  timeout: 1500

mappings:
  f1: help
  f2: save

commands:
  help:
    type: app
    command: help.exe
  save:
    type: script
    command: save.bat
`,
			filename: ConfigName + ".yml",
			validate: func(t *testing.T, cfg *Config) {
				// Check prefix
				if cfg.Hotkeys.Prefix != "win+shift+x" {
					t.Errorf("Expected prefix 'win+shift+x', got '%s'", cfg.Hotkeys.Prefix)
				}
				
				// Check mappings
				if cfg.Shortcuts["f1"] != "help" {
					t.Error("Expected mapping 'f1' -> 'help'")
				}
				if cfg.Shortcuts["f2"] != "save" {
					t.Error("Expected mapping 'f2' -> 'save'")
				}
				
				// Check commands  
				if cmd, ok := cfg.Actions["help"]; ok {
					if cmd.Type != "app" {
						t.Errorf("Expected type 'app', got '%s'", cmd.Type)
					}
					if cmd.Command != "help.exe" {
						t.Errorf("Expected command 'help.exe', got '%s'", cmd.Command)
					}
				} else {
					t.Error("Expected 'help' in commands")
				}
			},
		},
		{
			name: "Mixed languages config",
			setupKeys: func() {
				KeyDaemon = "service"
				KeyHotkeys = "ホットキー"
				KeyShortcuts = "mappings"
				KeyActions = "命令"
			},
			config: `
service:
  auto_start: true

ホットキー:
  prefix: win+x

mappings:
  test: test-cmd

命令:
  test-cmd:
    type: script
    command: echo "Mixed language config works!"
`,
			filename: ConfigName + ".yml", 
			validate: func(t *testing.T, cfg *Config) {
				if !cfg.Daemon.AutoStart {
					t.Error("Expected auto_start to be true")
				}
				if cfg.Hotkeys.Prefix != "win+x" {
					t.Errorf("Expected prefix 'win+x', got '%s'", cfg.Hotkeys.Prefix)
				}
				if cfg.Shortcuts["test"] != "test-cmd" {
					t.Error("Expected mapping 'test' -> 'test-cmd'")
				}
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory for this test
			tmpDir, err := os.MkdirTemp("", "spellbook-test-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)
			
			// Setup custom keys
			tt.setupKeys()
			
			// Write config file
			configPath := filepath.Join(tmpDir, tt.filename)
			if err := os.WriteFile(configPath, []byte(tt.config), 0644); err != nil {
				t.Fatalf("Failed to write config: %v", err)
			}
			
			// Load config
			loader := NewLoader(tmpDir)
			cfg, err := loader.Load()
			if err != nil {
				t.Fatalf("Failed to load config: %v", err)
			}
			
			// Validate
			tt.validate(t, cfg)
		})
	}
}

func TestPrefixExplicitlySetWithCustomKeys(t *testing.T) {
	// Save original
	origHotkeys := KeyHotkeys
	defer func() {
		KeyHotkeys = origHotkeys
	}()
	
	// Use Japanese key
	KeyHotkeys = "ホットキー"
	
	tmpDir, err := os.MkdirTemp("", "spellbook-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	
	// Test with explicitly empty prefix
	config := `
ホットキー:
  prefix: ""
`
	
	configPath := filepath.Join(tmpDir, ConfigName+".yml")
	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}
	
	loader := NewLoader(tmpDir)
	cfg, err := loader.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	// Should have empty prefix (not default)
	if cfg.Hotkeys.Prefix != "" {
		t.Errorf("Expected empty prefix, got '%s'", cfg.Hotkeys.Prefix)
	}
	
	// prefixExplicitlySet should be true
	if !cfg.prefixExplicitlySet {
		t.Error("Expected prefixExplicitlySet to be true")
	}
}