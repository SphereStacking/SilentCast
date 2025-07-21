package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListSpellsCommand(t *testing.T) {
	// Create temp directory for test configs
	tempDir := t.TempDir()

	// Test config with various spells
	testConfig := `
hotkeys:
  prefix: "alt+space"
spells:
  e: "editor"
  t: "terminal"
  "g,s": "git_status"
  "g,c": "git_commit"
  x: "missing_action"
grimoire:
  missing_action:
    type: app
    command: /bin/ls
    description: "Missing action"
  editor:
    type: app
    command: /bin/sh
    description: "Open shell"
  terminal:
    type: app
    command: /bin/echo
    description: "Echo command"
  git_status:
    type: script
    command: "git status"
    description: "Show git status"
  git_commit:
    type: script
    command: "git commit"
    description: "Create git commit"
`

	// Create config file
	configPath := filepath.Join(tempDir, "test")
	os.MkdirAll(configPath, 0o755)
	err := os.WriteFile(filepath.Join(configPath, "spellbook.yml"), []byte(testConfig), 0o644)
	if err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	getConfigPath := func() string { return configPath }

	tests := []struct {
		name         string
		flags        *Flags
		wantActive   bool
		wantContains []string
		wantMissing  []string
	}{
		{
			name: "list all spells",
			flags: &Flags{
				ListSpells: true,
			},
			wantActive: true,
			wantContains: []string{
				"SilentCast Spells",
				"Prefix: alt+space",
				"SEQUENCE",
				"alt+space → e",
				"editor",
				"/bin/sh",
				"Open shell",
				"alt+space → t",
				"terminal",
				"alt+space → g,s",
				"git_status",
				"alt+space → g,c",
				"git_commit",
				"alt+space → x",
				"missing_action",
				"Missing action",
				"Total: 5 spells",
			},
		},
		{
			name: "filter by sequence",
			flags: &Flags{
				ListSpells: true,
				ListFilter: "g",
			},
			wantActive: true,
			wantContains: []string{
				"(filtered: g)",
				"git_status",
				"git_commit",
			},
			wantMissing: []string{
				"editor",
				"terminal",
			},
		},
		{
			name: "filter by name",
			flags: &Flags{
				ListSpells: true,
				ListFilter: "editor",
			},
			wantActive: true,
			wantContains: []string{
				"(filtered: editor)",
				"editor",
			},
			wantMissing: []string{
				"terminal",
				"git_status",
			},
		},
		{
			name: "filter by description",
			flags: &Flags{
				ListSpells: true,
				ListFilter: "git",
			},
			wantActive: true,
			wantContains: []string{
				"git_status",
				"git_commit",
			},
			wantMissing: []string{
				"editor",
				"terminal",
			},
		},
		{
			name: "no matching filter",
			flags: &Flags{
				ListSpells: true,
				ListFilter: "nonexistent",
			},
			wantActive: true,
			wantContains: []string{
				"No spells found matching filter: nonexistent",
			},
		},
		{
			name:       "command not active",
			flags:      &Flags{ListSpells: false},
			wantActive: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewListSpellsCommand(getConfigPath)

			// Test metadata
			if cmd.Name() != "List Spells" {
				t.Errorf("Name() = %v, want %v", cmd.Name(), "List Spells")
			}

			if cmd.FlagName() != "list-spells" {
				t.Errorf("FlagName() = %v, want %v", cmd.FlagName(), "list-spells")
			}

			if cmd.Group() != "spell" {
				t.Errorf("Group() = %v, want %v", cmd.Group(), "spell")
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
						t.Logf("Full output:\n%s", output)
					}
				}

				// Check output doesn't contain unwanted strings
				for _, missing := range tt.wantMissing {
					if strings.Contains(output, missing) {
						t.Errorf("Execute() output should not contain %q", missing)
					}
				}
			}
		})
	}
}

func TestListSpellsCommand_EmptyConfig(t *testing.T) {
	// Create temp directory for test configs
	tempDir := t.TempDir()

	// Empty config
	emptyConfig := `
hotkeys:
  prefix: "alt+space"
`

	// Create config file
	configPath := filepath.Join(tempDir, "empty")
	os.MkdirAll(configPath, 0o755)
	err := os.WriteFile(filepath.Join(configPath, "spellbook.yml"), []byte(emptyConfig), 0o644)
	if err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	getConfigPath := func() string { return configPath }
	cmd := NewListSpellsCommand(getConfigPath)

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err = cmd.Execute(&Flags{ListSpells: true})

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}

	// Read captured output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "No spells configured") {
		t.Errorf("Execute() should show 'No spells configured' for empty config")
	}
}

func TestListSpellsCommand_InvalidFlags(t *testing.T) {
	getConfigPath := func() string { return "." }
	cmd := NewListSpellsCommand(getConfigPath)

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
