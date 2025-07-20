//go:build integration

package integration

import (
	"context"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActionExecution_AllTypes(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := SetupTestEnvironment(t)
	
	config := `
spells:
  app: "test_app"
  script: "test_script"
  url: "test_url"

grimoire:
  test_app:
    type: app
    command: "echo"
    args: ["App executed successfully"]
    description: "Test application action"
    
  test_script:
    type: script
    command: "echo"
    args: ["Script executed successfully"]
    show_output: true
    description: "Test script action"
    
  test_url:
    type: url
    command: "https://example.com/test"
    description: "Test URL action"
`
	
	env.WriteConfigFile("spellbook.yml", config)
	env.LoadConfig()
	env.InitializeComponents()
	
	tests := []struct {
		name       string
		actionName string
		actionType string
	}{
		{"App action", "test_app", "app"},
		{"Script action", "test_script", "script"},
		{"URL action", "test_url", "url"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify action configuration
			action := env.Config.Actions[tt.actionName]
			assert.Equal(t, tt.actionType, action.Type)
			
			// Execute action
			err := env.ExecuteAction(tt.actionName)
			
			// App and script actions using echo should succeed
			// URL actions might fail in test environment, but that's expected
			if tt.actionType == "url" {
				// URL actions may fail in headless environments
				t.Logf("URL action result: %v (failure expected in test environment)", err)
			} else {
				assert.NoError(t, err, "%s should execute successfully", tt.name)
			}
		})
	}
}

func TestScriptActionExecution_OutputHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := SetupTestEnvironment(t)
	
	config := `
spells:
  output_true: "script_with_output"
  output_false: "script_no_output"
  long_output: "script_long_output"

grimoire:
  script_with_output:
    type: script
    command: "echo"
    args: ["Output should be shown"]
    show_output: true
    description: "Script with output display"
    
  script_no_output:
    type: script
    command: "echo"
    args: ["Output should not be shown"]
    show_output: false
    description: "Script without output display"
    
  script_long_output:
    type: script
    command: "sh"
    args: ["-c", "for i in $(seq 1 5); do echo 'Line $i'; done"]
    show_output: true
    description: "Script with multiple lines of output"
`
	
	env.WriteConfigFile("spellbook.yml", config)
	env.LoadConfig()
	env.InitializeComponents()
	
	tests := []struct {
		name       string
		actionName string
		showOutput bool
	}{
		{"With output display", "script_with_output", true},
		{"Without output display", "script_no_output", false},
		{"Long output", "script_long_output", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action := env.Config.Actions[tt.actionName]
			assert.Equal(t, tt.showOutput, action.ShowOutput)
			
			err := env.ExecuteAction(tt.actionName)
			assert.NoError(t, err, "Script action should execute successfully")
			
			// In a real test, we would verify that notifications were sent
			// when show_output is true, but for this integration test,
			// we just verify execution succeeded
		})
	}
}

func TestScriptActionExecution_WorkingDirectory(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := SetupTestEnvironment(t)
	
	// Create a test subdirectory
	testDir := env.TempDir + "/testdir"
	err := os.MkdirAll(testDir, 0755)
	require.NoError(t, err)
	
	config := `
spells:
  pwd_test: "working_dir_test"

grimoire:
  working_dir_test:
    type: script
    command: "pwd"
    working_dir: "` + testDir + `"
    show_output: true
    description: "Test working directory setting"
`
	
	env.WriteConfigFile("spellbook.yml", config)
	env.LoadConfig()
	env.InitializeComponents()
	
	action := env.Config.Actions["working_dir_test"]
	assert.Equal(t, testDir, action.WorkingDir)
	
	err = env.ExecuteAction("working_dir_test")
	assert.NoError(t, err, "Action with working directory should execute successfully")
}

func TestScriptActionExecution_Environment(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := SetupTestEnvironment(t)
	
	var envCommand string
	var envArgs []string
	
	switch runtime.GOOS {
	case "windows":
		envCommand = "cmd"
		envArgs = []string{"/c", "echo %TEST_VAR%"}
	default:
		envCommand = "sh"
		envArgs = []string{"-c", "echo $TEST_VAR"}
	}
	
	config := `
spells:
  env_test: "environment_test"

grimoire:
  environment_test:
    type: script
    command: "` + envCommand + `"
    args: ` + formatArgsForYAML(envArgs) + `
    env:
      TEST_VAR: "integration_test_value"
    show_output: true
    description: "Test environment variable setting"
`
	
	env.WriteConfigFile("spellbook.yml", config)
	env.LoadConfig()
	env.InitializeComponents()
	
	action := env.Config.Actions["environment_test"]
	assert.Contains(t, action.Env, "TEST_VAR")
	assert.Equal(t, "integration_test_value", action.Env["TEST_VAR"])
	
	err := env.ExecuteAction("environment_test")
	assert.NoError(t, err, "Action with environment variables should execute successfully")
}

func TestScriptActionExecution_Timeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := SetupTestEnvironment(t)
	
	var sleepCommand string
	var sleepArgs []string
	
	switch runtime.GOOS {
	case "windows":
		sleepCommand = "timeout"
		sleepArgs = []string{"/t", "2", "/nobreak"}
	default:
		sleepCommand = "sleep"
		sleepArgs = []string{"2"}
	}
	
	config := `
spells:
  quick_timeout: "quick_timeout_test"
  long_timeout: "long_timeout_test"

grimoire:
  quick_timeout_test:
    type: script
    command: "` + sleepCommand + `"
    args: ` + formatArgsForYAML(sleepArgs) + `
    timeout: 1
    description: "Test timeout (should timeout)"
    
  long_timeout_test:
    type: script
    command: "echo"
    args: ["Quick execution"]
    timeout: 10
    description: "Test with generous timeout"
`
	
	env.WriteConfigFile("spellbook.yml", config)
	env.LoadConfig()
	env.InitializeComponents()
	
	// Test quick timeout (should timeout)
	t.Run("Should timeout", func(t *testing.T) {
		start := time.Now()
		err := env.ExecuteAction("quick_timeout_test")
		elapsed := time.Since(start)
		
		// Should fail due to timeout
		assert.Error(t, err, "Action should timeout")
		// Should complete within reasonable time (not wait for full sleep)
		assert.Less(t, elapsed.Seconds(), 3.0, "Should timeout quickly")
	})
	
	// Test generous timeout (should succeed)
	t.Run("Should succeed", func(t *testing.T) {
		err := env.ExecuteAction("long_timeout_test")
		assert.NoError(t, err, "Action with generous timeout should succeed")
	})
}

func TestAppActionExecution_PlatformSpecific(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := SetupTestEnvironment(t)
	
	// Use platform-specific commands that should exist
	var testCommand string
	var testArgs []string
	
	switch runtime.GOOS {
	case "windows":
		testCommand = "cmd"
		testArgs = []string{"/c", "echo Windows app test"}
	case "darwin":
		testCommand = "echo"
		testArgs = []string{"macOS app test"}
	default: // linux
		testCommand = "echo"
		testArgs = []string{"Linux app test"}
	}
	
	config := `
spells:
  platform_app: "platform_specific_app"

grimoire:
  platform_specific_app:
    type: app
    command: "` + testCommand + `"
    args: ` + formatArgsForYAML(testArgs) + `
    description: "Platform-specific application test"
`
	
	env.WriteConfigFile("spellbook.yml", config)
	env.LoadConfig()
	env.InitializeComponents()
	
	err := env.ExecuteAction("platform_specific_app")
	assert.NoError(t, err, "Platform-specific app should execute successfully")
}

func TestActionExecution_ErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := SetupTestEnvironment(t)
	
	config := `
spells:
  nonexistent_command: "nonexistent_command_action"
  failing_script: "failing_script_action"

grimoire:
  nonexistent_command_action:
    type: app
    command: "nonexistent_command_12345"
    description: "This command should not exist"
    
  failing_script_action:
    type: script
    command: "sh"
    args: ["-c", "exit 1"]
    description: "Script that exits with error"
`
	
	env.WriteConfigFile("spellbook.yml", config)
	env.LoadConfig()
	env.InitializeComponents()
	
	// Test nonexistent command
	t.Run("Nonexistent command", func(t *testing.T) {
		err := env.ExecuteAction("nonexistent_command_action")
		assert.Error(t, err, "Nonexistent command should fail")
	})
	
	// Test failing script
	t.Run("Failing script", func(t *testing.T) {
		err := env.ExecuteAction("failing_script_action")
		assert.Error(t, err, "Script with exit code 1 should fail")
	})
	
	// Test nonexistent action
	t.Run("Nonexistent action", func(t *testing.T) {
		err := env.ExecuteAction("nonexistent_action")
		assert.Error(t, err, "Nonexistent action should fail")
	})
}

func TestActionExecution_PermissionRequirements(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := SetupTestEnvironment(t)
	
	config := `
spells:
  admin_action: "admin_required_action"
  normal_action: "normal_action"

grimoire:
  admin_required_action:
    type: script
    command: "echo"
    args: ["Admin action"]
    admin: true
    description: "Action requiring admin privileges"
    
  normal_action:
    type: script
    command: "echo"
    args: ["Normal action"]
    admin: false
    description: "Normal action"
`
	
	env.WriteConfigFile("spellbook.yml", config)
	env.LoadConfig()
	env.InitializeComponents()
	
	// Test normal action (should work)
	t.Run("Normal action", func(t *testing.T) {
		err := env.ExecuteAction("normal_action")
		assert.NoError(t, err, "Normal action should execute successfully")
	})
	
	// Test admin action (behavior depends on platform and test environment)
	t.Run("Admin action", func(t *testing.T) {
		err := env.ExecuteAction("admin_required_action")
		// In test environment, this might fail due to permission requirements
		// or succeed if running with appropriate privileges
		t.Logf("Admin action result: %v", err)
		// We don't assert success/failure as it depends on test environment
	})
}

func TestActionExecution_ConcurrentExecution(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := SetupTestEnvironment(t)
	
	config := `
spells:
  concurrent1: "concurrent_action_1"
  concurrent2: "concurrent_action_2"
  concurrent3: "concurrent_action_3"

grimoire:
  concurrent_action_1:
    type: script
    command: "echo"
    args: ["Concurrent 1"]
    
  concurrent_action_2:
    type: script
    command: "echo"
    args: ["Concurrent 2"]
    
  concurrent_action_3:
    type: script
    command: "echo"
    args: ["Concurrent 3"]
`
	
	env.WriteConfigFile("spellbook.yml", config)
	env.LoadConfig()
	env.InitializeComponents()
	
	// Execute multiple actions concurrently
	results := make(chan error, 3)
	actions := []string{"concurrent_action_1", "concurrent_action_2", "concurrent_action_3"}
	
	for _, action := range actions {
		go func(actionName string) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			
			// Use action manager directly with context
			err := env.ActionManager.Execute(ctx, actionName)
			results <- err
		}(action)
	}
	
	// Wait for all executions to complete
	for i := 0; i < len(actions); i++ {
		select {
		case err := <-results:
			assert.NoError(t, err, "Concurrent action should execute successfully")
		case <-time.After(10 * time.Second):
			t.Error("Concurrent execution timed out")
		}
	}
}

// Helper function to format args array for YAML
func formatArgsForYAML(args []string) string {
	if len(args) == 0 {
		return "[]"
	}
	
	result := "["
	for i, arg := range args {
		if i > 0 {
			result += ", "
		}
		result += `"` + arg + `"`
	}
	result += "]"
	return result
}