//go:build integration

package integration

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/SphereStacking/silentcast/internal/config"
	"github.com/SphereStacking/silentcast/internal/action"
	"github.com/SphereStacking/silentcast/internal/hotkey"
	"github.com/SphereStacking/silentcast/internal/notify"
	"github.com/SphereStacking/silentcast/internal/permission"
	"github.com/SphereStacking/silentcast/internal/output"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEnvironment represents a complete integration test environment
type TestEnvironment struct {
	t *testing.T
	
	// Test directories
	TempDir    string
	ConfigDir  string
	
	// Components
	Config         *config.Config
	ConfigLoader   *config.Loader
	ActionManager  *action.Manager
	HotkeyManager  hotkey.Manager
	NotifyManager  *notify.Manager
	OutputManager  output.Manager
	PermManager    permission.Manager
	
	// Test control
	ctx    context.Context
	cancel context.CancelFunc
}

// SetupTestEnvironment creates a new test environment with all components
func SetupTestEnvironment(t *testing.T) *TestEnvironment {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "silentcast-integration-*")
	require.NoError(t, err)
	
	configDir := filepath.Join(tempDir, "config")
	require.NoError(t, os.MkdirAll(configDir, 0o755))
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	
	env := &TestEnvironment{
		t:         t,
		TempDir:   tempDir,
		ConfigDir: configDir,
		ctx:       ctx,
		cancel:    cancel,
	}
	
	// Setup cleanup
	t.Cleanup(func() {
		env.Cleanup()
	})
	
	return env
}

// WriteConfigFile writes a configuration file to the test environment
func (e *TestEnvironment) WriteConfigFile(filename, content string) {
	configPath := filepath.Join(e.ConfigDir, filename)
	err := os.WriteFile(configPath, []byte(content), 0o644)
	require.NoError(e.t, err)
}

// LoadConfig loads configuration from the test environment
func (e *TestEnvironment) LoadConfig() {
	e.ConfigLoader = config.NewLoader(e.ConfigDir)
	cfg, err := e.ConfigLoader.Load()
	require.NoError(e.t, err)
	e.Config = cfg
}

// InitializeComponents initializes all components with the loaded configuration
func (e *TestEnvironment) InitializeComponents() {
	require.NotNil(e.t, e.Config, "Config must be loaded before initializing components")
	
	// Initialize output manager
	e.OutputManager = output.NewBufferedManager(output.DefaultOptions())
	
	// Initialize notification manager with console notifier for testing
	e.NotifyManager = notify.NewManager()
	
	// Initialize permission manager
	permManager, err := permission.NewManager()
	require.NoError(e.t, err)
	e.PermManager = permManager
	
	// Initialize action manager
	e.ActionManager = action.NewManager(e.Config.Actions)
	
	// Initialize hotkey manager (using mock for testing)
	e.HotkeyManager = hotkey.NewMockManager()
}

// RegisterSpells registers spells from the configuration with the hotkey manager
func (e *TestEnvironment) RegisterSpells() error {
	require.NotNil(e.t, e.HotkeyManager, "HotkeyManager must be initialized")
	require.NotNil(e.t, e.Config, "Config must be loaded")
	
	// Register spell shortcuts
	for spell, actionName := range e.Config.Shortcuts {
		if _, exists := e.Config.Actions[actionName]; !exists {
			continue // Skip spells with no corresponding actions
		}
		
		err := e.HotkeyManager.Register(spell, actionName)
		if err != nil {
			return err
		}
	}
	
	return nil
}

// ExecuteAction executes an action by name and returns the result
func (e *TestEnvironment) ExecuteAction(actionName string) error {
	require.NotNil(e.t, e.ActionManager, "ActionManager must be initialized")
	return e.ActionManager.Execute(e.ctx, actionName)
}

// SimulateSpell simulates a spell sequence and returns the result
func (e *TestEnvironment) SimulateSpell(sequence string) error {
	require.NotNil(e.t, e.HotkeyManager, "HotkeyManager must be initialized")
	
	// For mock manager, we'll directly execute the action
	if actionName, exists := e.Config.Shortcuts[sequence]; exists {
		return e.ExecuteAction(actionName)
	}
	
	return assert.AnError // Spell not found
}

// UpdateConfig reloads configuration from files
func (e *TestEnvironment) UpdateConfig() error {
	require.NotNil(e.t, e.ConfigLoader, "ConfigLoader must be initialized")
	
	cfg, err := e.ConfigLoader.Load()
	if err != nil {
		return err
	}
	
	e.Config = cfg
	
	// Update action manager with new configuration
	if e.ActionManager != nil {
		e.ActionManager.UpdateActions(cfg.Actions)
	}
	
	return nil
}

// WaitForNotification waits for a notification to be sent
func (e *TestEnvironment) WaitForNotification(timeout time.Duration) bool {
	// For console notifier, we can't easily wait for notifications
	// This would be implemented for a mock notifier
	time.Sleep(100 * time.Millisecond) // Small delay to allow processing
	return true
}

// AssertConfigValid validates that the configuration is properly loaded
func (e *TestEnvironment) AssertConfigValid() {
	assert.NotNil(e.t, e.Config)
	assert.NotEmpty(e.t, e.Config.Actions, "Expected at least one action")
}

// AssertActionExists validates that an action exists in the configuration
func (e *TestEnvironment) AssertActionExists(actionName string) {
	assert.Contains(e.t, e.Config.Actions, actionName, "Action %s should exist", actionName)
}

// AssertSpellMapped validates that a spell is mapped to an action
func (e *TestEnvironment) AssertSpellMapped(spell, expectedAction string) {
	action, exists := e.Config.Shortcuts[spell]
	assert.True(e.t, exists, "Spell %s should be mapped", spell)
	assert.Equal(e.t, expectedAction, action, "Spell %s should map to %s", spell, expectedAction)
}

// Cleanup cleans up the test environment
func (e *TestEnvironment) Cleanup() {
	if e.cancel != nil {
		e.cancel()
	}
	if e.TempDir != "" {
		os.RemoveAll(e.TempDir)
	}
}

// GetTestConfig returns a sample configuration for testing
func GetTestConfig() string {
	return `
daemon:
  auto_start: false
  log_level: info
  config_watch: false

logger:
  level: info
  file: ""
  max_size: 10
  max_backups: 3
  max_age: 7
  compress: false

hotkeys:
  prefix: "alt+space"
  timeout: 1000
  sequence_timeout: 2000

spells:
  e: "editor"
  t: "terminal"
  b: "browser"
  "g,s": "git_status"
  "g,p": "git_pull"
  "o,g": "open_github"

grimoire:
  editor:
    type: app
    command: "echo"
    args: ["Editor opened"]
    description: "Open editor"
    
  terminal:
    type: app
    command: "echo"
    args: ["Terminal opened"]
    description: "Open terminal"
    
  browser:
    type: url
    command: "https://example.com"
    description: "Open browser"
    
  git_status:
    type: script
    command: "echo"
    args: ["Git status output"]
    show_output: true
    description: "Show git status"
    
  git_pull:
    type: script
    command: "echo"
    args: ["Git pull completed"]
    show_output: true
    description: "Git pull"
    
  open_github:
    type: url
    command: "https://github.com"
    description: "Open GitHub"
`
}

// GetMinimalConfig returns a minimal configuration for basic testing
func GetMinimalConfig() string {
	return `
spells:
  t: "test_action"

grimoire:
  test_action:
    type: script
    command: "echo"
    args: ["test output"]
    description: "Test action"
`
}

// GetPlatformSpecificConfig returns platform-specific overrides
func GetPlatformSpecificConfig() string {
	return `
grimoire:
  editor:
    type: app
    command: "notepad"
    description: "Platform-specific editor"
`
}