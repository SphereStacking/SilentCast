//go:build integration

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHotkeyActionIntegration_Basic(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := SetupTestEnvironment(t)
	env.WriteConfigFile("spellbook.yml", GetTestConfig())
	env.LoadConfig()
	env.InitializeComponents()
	
	// Register spells with hotkey manager
	err := env.RegisterSpells()
	require.NoError(t, err)
	
	// Test single key spells
	singleKeyTests := []struct {
		spell      string
		actionName string
	}{
		{"e", "editor"},
		{"t", "terminal"},
		{"b", "browser"},
	}
	
	for _, tt := range singleKeyTests {
		t.Run("Single key spell: "+tt.spell, func(t *testing.T) {
			env.AssertSpellMapped(tt.spell, tt.actionName)
			env.AssertActionExists(tt.actionName)
			
			// Simulate spell execution
			err := env.SimulateSpell(tt.spell)
			assert.NoError(t, err, "Spell %s should execute successfully", tt.spell)
		})
	}
	
	// Test sequence spells
	sequenceTests := []struct {
		sequence   string
		actionName string
	}{
		{"g,s", "git_status"},
		{"g,p", "git_pull"},
		{"o,g", "open_github"},
	}
	
	for _, tt := range sequenceTests {
		t.Run("Sequence spell: "+tt.sequence, func(t *testing.T) {
			env.AssertSpellMapped(tt.sequence, tt.actionName)
			env.AssertActionExists(tt.actionName)
			
			// Simulate sequence execution
			err := env.SimulateSpell(tt.sequence)
			assert.NoError(t, err, "Sequence %s should execute successfully", tt.sequence)
		})
	}
}

func TestHotkeyManager_PrefixRegistration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tests := []struct {
		name        string
		prefix      string
		expectError bool
	}{
		{
			name:        "Valid prefix key",
			prefix:      "alt+space",
			expectError: false,
		},
		{
			name:        "Single key prefix",
			prefix:      "f1",
			expectError: false,
		},
		{
			name:        "Complex modifier combination",
			prefix:      "ctrl+alt+shift+f12",
			expectError: false,
		},
		{
			name:        "Empty prefix",
			prefix:      "",
			expectError: false, // Should handle gracefully
		},
		{
			name:        "Invalid key name",
			prefix:      "invalid_key_name",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := SetupTestEnvironment(t)
			
			config := `
hotkeys:
  prefix: "` + tt.prefix + `"
  timeout: 1000

spells:
  t: "test"

grimoire:
  test:
    type: script
    command: "echo test"
`
			env.WriteConfigFile("spellbook.yml", config)
			env.LoadConfig()
			env.InitializeComponents()
			
			err := env.RegisterSpells()
			
			if tt.expectError {
				assert.Error(t, err, "Should fail with invalid prefix: %s", tt.prefix)
			} else {
				assert.NoError(t, err, "Should succeed with valid prefix: %s", tt.prefix)
			}
		})
	}
}

func TestHotkeyManager_SequenceHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := SetupTestEnvironment(t)
	
	config := `
hotkeys:
  prefix: "alt+space"
  timeout: 1000
  sequence_timeout: 2000

spells:
  # Single keys
  e: "editor"
  g: "git_menu"
  
  # Sequences
  "g,s": "git_status"
  "g,l": "git_log"
  "g,p": "git_push"
  
  # Long sequences
  "a,b,c": "long_sequence"
  
  # Overlapping sequences
  "t": "terminal"
  "t,e": "terminal_editor"

grimoire:
  editor:
    type: app
    command: "echo editor"
  git_menu:
    type: script
    command: "echo git menu"
  git_status:
    type: script
    command: "echo git status"
  git_log:
    type: script
    command: "echo git log"
  git_push:
    type: script
    command: "echo git push"
  long_sequence:
    type: script
    command: "echo long sequence executed"
  terminal:
    type: app
    command: "echo terminal"
  terminal_editor:
    type: script
    command: "echo terminal with editor"
`
	
	env.WriteConfigFile("spellbook.yml", config)
	env.LoadConfig()
	env.InitializeComponents()
	
	err := env.RegisterSpells()
	require.NoError(t, err)
	
	// Test that all spells are properly registered
	sequences := []string{"e", "g", "g,s", "g,l", "g,p", "a,b,c", "t", "t,e"}
	
	for _, seq := range sequences {
		t.Run("Sequence: "+seq, func(t *testing.T) {
			err := env.SimulateSpell(seq)
			assert.NoError(t, err, "Sequence %s should be registered and executable", seq)
		})
	}
}

func TestHotkeyManager_ErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := SetupTestEnvironment(t)
	
	config := `
hotkeys:
  prefix: "alt+space"

spells:
  valid: "valid_action"
  invalid: "nonexistent_action"

grimoire:
  valid_action:
    type: script
    command: "echo valid"
`
	
	env.WriteConfigFile("spellbook.yml", config)
	env.LoadConfig()
	env.InitializeComponents()
	
	// Register spells - this should succeed even with invalid mappings
	err := env.RegisterSpells()
	assert.NoError(t, err, "Registration should succeed even with invalid spell mappings")
	
	// Test valid spell
	err = env.SimulateSpell("valid")
	assert.NoError(t, err, "Valid spell should execute successfully")
	
	// Test invalid spell - this will fail in our test environment
	err = env.SimulateSpell("invalid")
	assert.Error(t, err, "Invalid spell should fail to execute")
	
	// Test nonexistent spell
	err = env.SimulateSpell("nonexistent")
	assert.Error(t, err, "Nonexistent spell should fail")
}

func TestHotkeyManager_ConfigurationReload(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := SetupTestEnvironment(t)
	
	// Initial configuration
	initialConfig := `
hotkeys:
  prefix: "alt+space"

spells:
  t: "test1"

grimoire:
  test1:
    type: script
    command: "echo test1"
`
	
	env.WriteConfigFile("spellbook.yml", initialConfig)
	env.LoadConfig()
	env.InitializeComponents()
	err := env.RegisterSpells()
	require.NoError(t, err)
	
	// Test initial spell
	err = env.SimulateSpell("t")
	assert.NoError(t, err)
	
	// Update configuration
	updatedConfig := `
hotkeys:
  prefix: "alt+space"

spells:
  t: "test2"
  n: "new_test"

grimoire:
  test2:
    type: script
    command: "echo test2"
  new_test:
    type: app
    command: "echo new test"
`
	
	env.WriteConfigFile("spellbook.yml", updatedConfig)
	err = env.UpdateConfig()
	require.NoError(t, err)
	
	// Re-register spells with updated configuration
	err = env.RegisterSpells()
	require.NoError(t, err)
	
	// Test updated spells
	err = env.SimulateSpell("t")
	assert.NoError(t, err, "Updated spell should work")
	
	err = env.SimulateSpell("n")
	assert.NoError(t, err, "New spell should work")
}

func TestHotkeyManager_Timeouts(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := SetupTestEnvironment(t)
	
	config := `
hotkeys:
  prefix: "alt+space"
  timeout: 100
  sequence_timeout: 200

spells:
  "q,w": "quick_sequence"

grimoire:
  quick_sequence:
    type: script
    command: "echo quick"
`
	
	env.WriteConfigFile("spellbook.yml", config)
	env.LoadConfig()
	env.InitializeComponents()
	
	err := env.RegisterSpells()
	require.NoError(t, err)
	
	// Verify timeout configuration is loaded
	assert.Equal(t, 100, int(env.Config.Hotkeys.Timeout.ToDuration().Milliseconds()))
	assert.Equal(t, 200, int(env.Config.Hotkeys.SequenceTimeout.ToDuration().Milliseconds()))
	
	// Test sequence execution (timing behavior would be tested in real hotkey manager)
	err = env.SimulateSpell("q,w")
	assert.NoError(t, err, "Sequence should execute within timeout")
}

func TestHotkeyManager_KeyValidation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tests := []struct {
		name        string
		spells      map[string]string
		expectError bool
		description string
	}{
		{
			name: "Valid single keys",
			spells: map[string]string{
				"a": "action_a",
				"1": "action_1",
				"f1": "action_f1",
				"space": "action_space",
			},
			expectError: false,
			description: "Standard keys should be valid",
		},
		{
			name: "Valid sequences",
			spells: map[string]string{
				"a,b": "action_ab",
				"g,s": "action_gs",
				"1,2,3": "action_123",
			},
			expectError: false,
			description: "Comma-separated sequences should be valid",
		},
		{
			name: "Mixed valid keys",
			spells: map[string]string{
				"ctrl+c": "copy",
				"alt+f4": "close",
				"shift+tab": "reverse_tab",
			},
			expectError: false,
			description: "Modifier combinations should be valid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := SetupTestEnvironment(t)
			
			// Build configuration
			configStr := `
hotkeys:
  prefix: "alt+space"

spells:
`
			for spell, action := range tt.spells {
				configStr += "  \"" + spell + "\": \"" + action + "\"\n"
			}
			
			configStr += "\ngrimoire:\n"
			for _, action := range tt.spells {
				configStr += "  " + action + ":\n"
				configStr += "    type: script\n"
				configStr += "    command: \"echo " + action + "\"\n"
			}
			
			env.WriteConfigFile("spellbook.yml", configStr)
			env.LoadConfig()
			env.InitializeComponents()
			
			err := env.RegisterSpells()
			
			if tt.expectError {
				assert.Error(t, err, tt.description)
			} else {
				assert.NoError(t, err, tt.description)
				
				// Test that all spells can be executed
				for spell := range tt.spells {
					err := env.SimulateSpell(spell)
					assert.NoError(t, err, "Spell %s should execute", spell)
				}
			}
		})
	}
}