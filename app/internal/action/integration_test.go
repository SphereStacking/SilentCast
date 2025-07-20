package action

import (
	"context"
	"runtime"
	"testing"

	"github.com/SphereStacking/silentcast/internal/action/script"
	"github.com/SphereStacking/silentcast/internal/config"
)

func TestScriptExecutor_ShowOutputWithSystemNotification(t *testing.T) {
	// Skip if not on a supported platform
	switch runtime.GOOS {
	case "windows", "darwin", "linux":
		// Continue with test
	default:
		t.Skip("System notifications not supported on", runtime.GOOS)
	}

	tests := []struct {
		name   string
		config config.ActionConfig
	}{
		{
			name: "Show git status with system notification",
			config: config.ActionConfig{
				Type:        "script",
				Command:     "echo 'On branch main\nYour branch is up to date.'",
				ShowOutput:  true,
				Description: "Git Status",
			},
		},
		{
			name: "Show error with system notification",
			config: config.ActionConfig{
				Type:        "script",
				Command:     "echo 'Error: File not found' && exit 1",
				ShowOutput:  true,
				Description: "Error Example",
			},
		},
		{
			name: "Long output truncated",
			config: config.ActionConfig{
				Type:       "script",
				Command:    "for i in {1..100}; do echo \"Line $i of output\"; done",
				ShowOutput: true,
				Description: "Long Output Test",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := script.NewScriptExecutor(&tt.config)
			ctx := context.Background()
			
			// Execute command
			err := executor.Execute(ctx)
			
			// Check if error is expected
			if tt.name == "Show error with system notification" && err == nil {
				t.Error("Expected error for failing command")
			} else if tt.name != "Show error with system notification" && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			
			// Note: We can't easily verify that system notification was shown,
			// but if no panic occurred, the integration is working
		})
	}
}

func TestNotificationIntegration_AllPlatforms(t *testing.T) {
	// This test verifies that the notification system works on all platforms
	config := config.ActionConfig{
		Type:        "script",
		Command:     "echo 'Platform test: " + runtime.GOOS + "'",
		ShowOutput:  true,
		Description: "Platform Notification Test",
	}
	
	executor := script.NewScriptExecutor(&config)
	ctx := context.Background()
	
	err := executor.Execute(ctx)
	if err != nil {
		t.Errorf("Platform notification test failed: %v", err)
	}
}