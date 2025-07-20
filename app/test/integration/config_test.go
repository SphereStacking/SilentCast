//go:build integration

package integration

import (
	"os"
	"testing"
	"time"

	"github.com/SphereStacking/silentcast/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigurationLoading_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tests := []struct {
		name           string
		baseConfig     string
		platformConfig string
		platform       string
		validate       func(t *testing.T, cfg *config.Config)
	}{
		{
			name:       "Base configuration only",
			baseConfig: GetTestConfig(),
			validate: func(t *testing.T, cfg *config.Config) {
				assert.Equal(t, "alt+space", cfg.Hotkeys.Prefix)
				assert.Equal(t, 1000, int(cfg.Hotkeys.Timeout.ToDuration().Milliseconds()))
				assert.Len(t, cfg.Shortcuts, 6)
				assert.Len(t, cfg.Actions, 6)
				assert.Contains(t, cfg.Shortcuts, "e")
				assert.Contains(t, cfg.Shortcuts, "g,s")
			},
		},
		{
			name:       "Platform-specific overrides",
			baseConfig: GetTestConfig(),
			platformConfig: `
grimoire:
  editor:
    type: app
    command: "platform-editor"
    description: "Platform-specific editor"
`,
			platform: "windows",
			validate: func(t *testing.T, cfg *config.Config) {
				// Should have platform-specific editor command
				editor := cfg.Actions["editor"]
				assert.Equal(t, "platform-editor", editor.Command)
				assert.Equal(t, "Platform-specific editor", editor.Description)
				
				// Other actions should remain unchanged
				terminal := cfg.Actions["terminal"]
				assert.Equal(t, "echo", terminal.Command)
			},
		},
		{
			name: "Minimal configuration",
			baseConfig: `
spells:
  t: "test"
grimoire:
  test:
    type: script
    command: "echo test"
`,
			validate: func(t *testing.T, cfg *config.Config) {
				assert.Len(t, cfg.Shortcuts, 1)
				assert.Len(t, cfg.Actions, 1)
				assert.Equal(t, "test", cfg.Shortcuts["t"])
				
				// Should use default values
				assert.Equal(t, "alt+space", cfg.Hotkeys.Prefix)
				assert.Equal(t, time.Second, cfg.Hotkeys.Timeout.ToDuration())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := SetupTestEnvironment(t)
			
			// Write base configuration
			env.WriteConfigFile("spellbook.yml", tt.baseConfig)
			
			// Write platform-specific configuration if provided
			if tt.platformConfig != "" && tt.platform != "" {
				env.WriteConfigFile("spellbook."+tt.platform+".yml", tt.platformConfig)
			}
			
			// Load configuration
			env.LoadConfig()
			
			// Validate configuration
			env.AssertConfigValid()
			tt.validate(t, env.Config)
		})
	}
}

func TestConfigurationValidation_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tests := []struct {
		name        string
		config      string
		expectError bool
		errorType   string
	}{
		{
			name:        "Valid configuration",
			config:      GetTestConfig(),
			expectError: false,
		},
		{
			name: "Invalid YAML syntax",
			config: `
spells:
  e: "editor"
  invalid: yaml: syntax
`,
			expectError: true,
			errorType:   "yaml",
		},
		{
			name: "Missing action for spell",
			config: `
spells:
  e: "nonexistent_action"
grimoire:
  terminal:
    type: app
    command: "echo"
`,
			expectError: true,
			errorType:   "validation",
		},
		{
			name: "Invalid action type",
			config: `
spells:
  e: "editor"
grimoire:
  editor:
    type: "invalid_type"
    command: "echo"
`,
			expectError: true,
			errorType:   "validation",
		},
		{
			name: "Missing required command",
			config: `
spells:
  e: "editor"
grimoire:
  editor:
    type: "app"
    description: "Editor without command"
`,
			expectError: true,
			errorType:   "validation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := SetupTestEnvironment(t)
			env.WriteConfigFile("spellbook.yml", tt.config)
			
			loader := config.NewLoader(env.ConfigDir)
			cfg, err := loader.Load()
			
			if tt.expectError {
				assert.Error(t, err, "Expected error for invalid configuration")
				assert.Nil(t, cfg, "Config should be nil on error")
			} else {
				assert.NoError(t, err, "Expected no error for valid configuration")
				assert.NotNil(t, cfg, "Config should not be nil on success")
			}
		})
	}
}

func TestConfigurationWatcher_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := SetupTestEnvironment(t)
	
	// Write initial configuration
	initialConfig := `
spells:
  t: "test"
grimoire:
  test:
    type: script
    command: "echo initial"
`
	env.WriteConfigFile("spellbook.yml", initialConfig)
	env.LoadConfig()
	env.InitializeComponents()
	
	// Verify initial configuration
	assert.Equal(t, "echo", env.Config.Actions["test"].Command)
	assert.Equal(t, []string{"initial"}, env.Config.Actions["test"].Args)
	
	// Channel to receive reload events
	reloadChan := make(chan bool, 1)
	
	// Setup file watcher
	watcherConfig := config.WatcherConfig{
		ConfigPath: env.ConfigDir,
		OnChange: func(cfg *config.Config) {
			env.Config = cfg
			select {
			case reloadChan <- true:
			default:
			}
		},
		Debounce: 100 * time.Millisecond,
	}
	
	watcher, err := config.NewWatcher(watcherConfig)
	require.NoError(t, err)
	defer func() {
		if watcher != nil {
			watcher.Stop()
		}
	}()
	
	// Start watching
	watcher.Start(env.ctx)
	
	// Give watcher time to initialize
	time.Sleep(100 * time.Millisecond)
	
	// Update configuration file
	updatedConfig := `
spells:
  t: "test"
  n: "new_action"
grimoire:
  test:
    type: script
    command: "echo updated"
  new_action:
    type: app
    command: "new-command"
`
	env.WriteConfigFile("spellbook.yml", updatedConfig)
	
	// Wait for reload event
	select {
	case <-reloadChan:
		// Configuration should be updated
		assert.Equal(t, "echo", env.Config.Actions["test"].Command)
		assert.Equal(t, []string{"updated"}, env.Config.Actions["test"].Args)
		assert.Contains(t, env.Config.Actions, "new_action")
		assert.Equal(t, "new-command", env.Config.Actions["new_action"].Command)
	case <-time.After(2 * time.Second):
		t.Error("Configuration reload event not received within timeout")
	}
}

func TestConfigurationCascade_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := SetupTestEnvironment(t)
	
	// Write base configuration
	baseConfig := `
daemon:
  auto_start: true
  log_level: info

hotkeys:
  prefix: "ctrl+space"
  timeout: 1500

spells:
  e: "editor"
  t: "terminal"

grimoire:
  editor:
    type: app
    command: "base-editor"
    description: "Base editor"
  terminal:
    type: app
    command: "base-terminal"
    description: "Base terminal"
`
	
	// Write platform-specific overrides
	platformConfig := `
daemon:
  log_level: debug

hotkeys:
  timeout: 2000

grimoire:
  editor:
    command: "platform-editor"
    args: ["--platform-flag"]
  new_platform_action:
    type: script
    command: "platform-specific-script"
`
	
	env.WriteConfigFile("spellbook.yml", baseConfig)
	env.WriteConfigFile("spellbook.linux.yml", platformConfig)
	
	// Load cascaded configuration
	env.LoadConfig()
	
	// Verify cascade behavior
	// Daemon settings: log_level overridden, auto_start preserved
	assert.True(t, env.Config.Daemon.AutoStart)
	assert.Equal(t, "debug", env.Config.Daemon.LogLevel)
	
	// Hotkeys: timeout overridden, prefix preserved
	assert.Equal(t, "ctrl+space", env.Config.Hotkeys.Prefix)
	assert.Equal(t, 2000, int(env.Config.Hotkeys.Timeout.ToDuration().Milliseconds()))
	
	// Spells: preserved from base
	assert.Len(t, env.Config.Shortcuts, 2)
	
	// Actions: editor overridden, terminal preserved, new action added
	editor := env.Config.Actions["editor"]
	assert.Equal(t, "platform-editor", editor.Command)
	assert.Equal(t, []string{"--platform-flag"}, editor.Args)
	assert.Equal(t, "Base editor", editor.Description) // Description preserved
	
	terminal := env.Config.Actions["terminal"]
	assert.Equal(t, "base-terminal", terminal.Command)
	
	assert.Contains(t, env.Config.Actions, "new_platform_action")
}

func TestConfigurationEnvironmentExpansion_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Set test environment variables
	os.Setenv("TEST_EDITOR", "test-editor")
	os.Setenv("TEST_HOME", "/test/home")
	defer func() {
		os.Unsetenv("TEST_EDITOR")
		os.Unsetenv("TEST_HOME")
	}()

	env := SetupTestEnvironment(t)
	
	configWithEnvVars := `
spells:
  e: "editor"
  h: "home_action"

grimoire:
  editor:
    type: app
    command: "${TEST_EDITOR}"
    working_dir: "${TEST_HOME}/projects"
    env:
      EDITOR_CONFIG: "${TEST_HOME}/.editor"
  home_action:
    type: script
    command: "ls"
    args: ["${TEST_HOME}"]
`
	
	env.WriteConfigFile("spellbook.yml", configWithEnvVars)
	env.LoadConfig()
	
	// Verify environment variable expansion
	editor := env.Config.Actions["editor"]
	assert.Equal(t, "test-editor", editor.Command)
	assert.Equal(t, "/test/home/projects", editor.WorkingDir)
	assert.Equal(t, "/test/home/.editor", editor.Env["EDITOR_CONFIG"])
	
	homeAction := env.Config.Actions["home_action"]
	assert.Equal(t, []string{"/test/home"}, homeAction.Args)
}

func TestConfigurationErrorReporting_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := SetupTestEnvironment(t)
	
	// Write configuration with multiple errors
	configWithErrors := `
spells:
  e: "editor"
  t: "nonexistent"  # Error: action doesn't exist
  invalid-spell-key!: "terminal"  # Error: invalid spell key

grimoire:
  editor:
    # Missing type field
    command: "echo"
  terminal:
    type: "invalid_type"  # Error: invalid action type
    command: "echo"
`
	
	env.WriteConfigFile("spellbook.yml", configWithErrors)
	
	loader := config.NewLoader(env.ConfigDir)
	_, err := loader.Load()
	
	// Should have validation errors
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation", "Error should mention validation issues")
}