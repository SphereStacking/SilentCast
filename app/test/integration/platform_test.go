//go:build integration

package integration

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrossPlatformIntegration_AllPlatforms(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This test runs the core integration scenarios on all platforms
	t.Run("Configuration Loading", func(t *testing.T) {
		testConfigurationLoadingCrossPlatform(t)
	})
	
	t.Run("Action Execution", func(t *testing.T) {
		testActionExecutionCrossPlatform(t)
	})
	
	t.Run("Notification System", func(t *testing.T) {
		testNotificationSystemCrossPlatform(t)
	})
	
	t.Run("Error Handling", func(t *testing.T) {
		testErrorHandlingCrossPlatform(t)
	})
}

func testConfigurationLoadingCrossPlatform(t *testing.T) {
	env := SetupTestEnvironment(t)
	
	// Create platform-specific configuration
	baseConfig := GetTestConfig()
	
	var platformConfig string
	switch runtime.GOOS {
	case "windows":
		platformConfig = `
grimoire:
  editor:
    type: app
    command: "notepad"
    description: "Windows Notepad"
  terminal:
    type: app
    command: "cmd"
    description: "Command Prompt"
`
	case "darwin":
		platformConfig = `
grimoire:
  editor:
    type: app
    command: "TextEdit"
    description: "macOS TextEdit"
  terminal:
    type: app
    command: "Terminal"
    description: "macOS Terminal"
`
	default: // linux
		platformConfig = `
grimoire:
  editor:
    type: app
    command: "nano"
    description: "Nano Editor"
  terminal:
    type: app
    command: "bash"
    description: "Bash Terminal"
`
	}
	
	// Write base and platform-specific configs
	env.WriteConfigFile("spellbook.yml", baseConfig)
	env.WriteConfigFile("spellbook."+runtime.GOOS+".yml", platformConfig)
	
	// Load configuration
	env.LoadConfig()
	env.AssertConfigValid()
	
	// Verify platform-specific overrides
	editor := env.Config.Actions["editor"]
	terminal := env.Config.Actions["terminal"]
	
	switch runtime.GOOS {
	case "windows":
		assert.Equal(t, "notepad", editor.Command)
		assert.Equal(t, "cmd", terminal.Command)
	case "darwin":
		assert.Equal(t, "TextEdit", editor.Command)
		assert.Equal(t, "Terminal", terminal.Command)
	default: // linux
		assert.Equal(t, "nano", editor.Command)
		assert.Equal(t, "bash", terminal.Command)
	}
}

func testActionExecutionCrossPlatform(t *testing.T) {
	env := SetupTestEnvironment(t)
	
	// Create configuration with platform-appropriate commands
	var config string
	switch runtime.GOOS {
	case "windows":
		config = `
spells:
  echo_test: "windows_echo"
  dir_test: "windows_dir"

grimoire:
  windows_echo:
    type: script
    command: "cmd"
    args: ["/c", "echo Windows echo test"]
    show_output: true
    
  windows_dir:
    type: script
    command: "cmd"
    args: ["/c", "dir", "/b"]
    working_dir: "C:\\Windows\\System32"
    show_output: true
`
	case "darwin":
		config = `
spells:
  echo_test: "macos_echo"
  ls_test: "macos_ls"

grimoire:
  macos_echo:
    type: script
    command: "echo"
    args: ["macOS echo test"]
    show_output: true
    
  macos_ls:
    type: script
    command: "ls"
    args: ["-la"]
    working_dir: "/Applications"
    show_output: true
`
	default: // linux
		config = `
spells:
  echo_test: "linux_echo"
  ls_test: "linux_ls"

grimoire:
  linux_echo:
    type: script
    command: "echo"
    args: ["Linux echo test"]
    show_output: true
    
  linux_ls:
    type: script
    command: "ls"
    args: ["-la"]
    working_dir: "/usr/bin"
    show_output: true
`
	}
	
	env.WriteConfigFile("spellbook.yml", config)
	env.LoadConfig()
	env.InitializeComponents()
	
	// Test platform-appropriate commands
	switch runtime.GOOS {
	case "windows":
		err := env.ExecuteAction("windows_echo")
		assert.NoError(t, err, "Windows echo should work")
		
		err = env.ExecuteAction("windows_dir")
		assert.NoError(t, err, "Windows dir should work")
		
	case "darwin":
		err := env.ExecuteAction("macos_echo")
		assert.NoError(t, err, "macOS echo should work")
		
		err = env.ExecuteAction("macos_ls")
		assert.NoError(t, err, "macOS ls should work")
		
	default: // linux
		err := env.ExecuteAction("linux_echo")
		assert.NoError(t, err, "Linux echo should work")
		
		err = env.ExecuteAction("linux_ls")
		assert.NoError(t, err, "Linux ls should work")
	}
}

func testNotificationSystemCrossPlatform(t *testing.T) {
	env := SetupTestEnvironment(t)
	
	config := `
spells:
  notify_test: "platform_notification"

grimoire:
  platform_notification:
    type: script
    command: "echo"
    args: ["Platform notification test for ` + runtime.GOOS + `"]
    show_output: true
    description: "Cross-platform notification test"
`
	
	env.WriteConfigFile("spellbook.yml", config)
	env.LoadConfig()
	env.InitializeComponents()
	
	err := env.ExecuteAction("platform_notification")
	assert.NoError(t, err, "Platform notification should work on "+runtime.GOOS)
	
	// Wait for notification processing
	env.WaitForNotification(500)
}

func testErrorHandlingCrossPlatform(t *testing.T) {
	env := SetupTestEnvironment(t)
	
	// Create configuration with platform-specific failing commands
	var config string
	switch runtime.GOOS {
	case "windows":
		config = `
spells:
  fail_test: "windows_fail"

grimoire:
  windows_fail:
    type: script
    command: "cmd"
    args: ["/c", "echo Error test && exit 1"]
    show_output: true
`
	default: // unix-like
		config = `
spells:
  fail_test: "unix_fail"

grimoire:
  unix_fail:
    type: script
    command: "sh"
    args: ["-c", "echo 'Error test' && exit 1"]
    show_output: true
`
	}
	
	env.WriteConfigFile("spellbook.yml", config)
	env.LoadConfig()
	env.InitializeComponents()
	
	err := env.ExecuteAction("windows_fail")
	if runtime.GOOS == "windows" {
		assert.Error(t, err, "Windows failing command should fail")
	}
	
	err = env.ExecuteAction("unix_fail")
	if runtime.GOOS != "windows" {
		assert.Error(t, err, "Unix failing command should fail")
	}
}

func TestPlatformSpecificFeatures(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	switch runtime.GOOS {
	case "windows":
		testWindowsSpecificFeatures(t)
	case "darwin":
		testMacOSSpecificFeatures(t)
	case "linux":
		testLinuxSpecificFeatures(t)
	default:
		t.Skipf("No platform-specific tests for %s", runtime.GOOS)
	}
}

func testWindowsSpecificFeatures(t *testing.T) {
	env := SetupTestEnvironment(t)
	
	config := `
spells:
  powershell_test: "powershell_action"
  admin_test: "admin_action"

grimoire:
  powershell_action:
    type: script
    command: "powershell"
    args: ["-Command", "Get-Location"]
    show_output: true
    
  admin_action:
    type: script
    command: "cmd"
    args: ["/c", "echo Admin test"]
    admin: true
    show_output: true
`
	
	env.WriteConfigFile("spellbook.yml", config)
	env.LoadConfig()
	env.InitializeComponents()
	
	// Test PowerShell execution
	err := env.ExecuteAction("powershell_action")
	if err != nil {
		t.Logf("PowerShell test failed (may not be available): %v", err)
	}
	
	// Test admin elevation (may fail in test environment)
	err = env.ExecuteAction("admin_action")
	t.Logf("Admin action result: %v", err)
}

func testMacOSSpecificFeatures(t *testing.T) {
	env := SetupTestEnvironment(t)
	
	config := `
spells:
  osascript_test: "osascript_action"
  admin_test: "admin_action"

grimoire:
  osascript_action:
    type: script
    command: "osascript"
    args: ["-e", "return (current date) as string"]
    show_output: true
    
  admin_action:
    type: script
    command: "echo"
    args: ["Admin test"]
    admin: true
    show_output: true
`
	
	env.WriteConfigFile("spellbook.yml", config)
	env.LoadConfig()
	env.InitializeComponents()
	
	// Test osascript execution
	err := env.ExecuteAction("osascript_action")
	if err != nil {
		t.Logf("osascript test failed: %v", err)
	}
	
	// Test admin elevation (may require user interaction)
	err = env.ExecuteAction("admin_action")
	t.Logf("Admin action result: %v", err)
}

func testLinuxSpecificFeatures(t *testing.T) {
	env := SetupTestEnvironment(t)
	
	config := `
spells:
  shell_test: "shell_variety"
  admin_test: "admin_action"

grimoire:
  shell_variety:
    type: script
    command: "bash"
    args: ["-c", "echo $SHELL"]
    show_output: true
    
  admin_action:
    type: script
    command: "echo"
    args: ["Admin test"]
    admin: true
    show_output: true
`
	
	env.WriteConfigFile("spellbook.yml", config)
	env.LoadConfig()
	env.InitializeComponents()
	
	// Test shell execution
	err := env.ExecuteAction("shell_variety")
	assert.NoError(t, err, "Shell command should work on Linux")
	
	// Test admin elevation (may require sudo/pkexec)
	err = env.ExecuteAction("admin_action")
	t.Logf("Admin action result: %v", err)
}

func TestFullWorkflowIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Test complete workflow: config load -> spell registration -> execution -> notification
	env := SetupTestEnvironment(t)
	
	// Create comprehensive configuration
	config := `
daemon:
  auto_start: false
  log_level: info

hotkeys:
  prefix: "alt+space"
  timeout: 1000

spells:
  w: "workflow_test"
  "w,s": "workflow_sequence"

grimoire:
  workflow_test:
    type: script
    command: "echo"
    args: ["Workflow test successful"]
    show_output: true
    description: "Complete workflow test"
    
  workflow_sequence:
    type: script
    command: "echo"
    args: ["Sequence workflow successful"]
    show_output: true
    description: "Sequence workflow test"
`
	
	env.WriteConfigFile("spellbook.yml", config)
	
	// Step 1: Load configuration
	env.LoadConfig()
	env.AssertConfigValid()
	
	// Step 2: Initialize all components
	env.InitializeComponents()
	
	// Step 3: Register spells
	err := env.RegisterSpells()
	assert.NoError(t, err, "Spell registration should succeed")
	
	// Step 4: Execute spells
	err = env.SimulateSpell("w")
	assert.NoError(t, err, "Single spell should execute")
	
	err = env.SimulateSpell("w,s")
	assert.NoError(t, err, "Sequence spell should execute")
	
	// Step 5: Verify notifications were triggered
	env.WaitForNotification(500)
	
	t.Log("Complete workflow integration test passed")
}