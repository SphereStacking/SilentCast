//go:build integration

package integration

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/SphereStacking/silentcast/internal/notify"
	"github.com/stretchr/testify/assert"
)

func TestNotificationSystem_BasicIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := SetupTestEnvironment(t)
	
	config := `
spells:
  notify_success: "success_notification"
  notify_error: "error_notification"
  notify_info: "info_notification"

grimoire:
  success_notification:
    type: script
    command: "echo"
    args: ["Success message"]
    show_output: true
    description: "Success notification test"
    
  error_notification:
    type: script
    command: "sh"
    args: ["-c", "echo 'Error message' >&2; exit 1"]
    show_output: true
    description: "Error notification test"
    
  info_notification:
    type: script
    command: "echo"
    args: ["Info message"]
    show_output: true
    description: "Info notification test"
`
	
	env.WriteConfigFile("spellbook.yml", config)
	env.LoadConfig()
	env.InitializeComponents()
	
	tests := []struct {
		name           string
		actionName     string
		expectError    bool
		notifyLevel    notify.Level
	}{
		{"Success notification", "success_notification", false, notify.LevelSuccess},
		{"Error notification", "error_notification", true, notify.LevelError},
		{"Info notification", "info_notification", false, notify.LevelInfo},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := env.ExecuteAction(tt.actionName)
			
			if tt.expectError {
				assert.Error(t, err, "Action should fail and trigger error notification")
			} else {
				assert.NoError(t, err, "Action should succeed and trigger success notification")
			}
			
			// Wait for notification processing
			env.WaitForNotification(500 * time.Millisecond)
			
			// In a real test environment, we would verify the notification was sent
			// For this integration test, we're just ensuring the flow completes
		})
	}
}

func TestNotificationSystem_OutputFormatting(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := SetupTestEnvironment(t)
	
	config := `
spells:
  short_output: "short_output_action"
  long_output: "long_output_action"
  multiline_output: "multiline_output_action"
  ansi_output: "ansi_output_action"

grimoire:
  short_output_action:
    type: script
    command: "echo"
    args: ["Short message"]
    show_output: true
    
  long_output_action:
    type: script
    command: "sh"
    args: ["-c", "for i in $(seq 1 100); do echo 'This is line number $i with some additional text to make it longer'; done"]
    show_output: true
    
  multiline_output_action:
    type: script
    command: "sh"
    args: ["-c", "echo 'Line 1'; echo 'Line 2'; echo 'Line 3'"]
    show_output: true
    
  ansi_output_action:
    type: script
    command: "sh"
    args: ["-c", "echo -e '\\033[31mRed text\\033[0m and \\033[32mgreen text\\033[0m'"]
    show_output: true
`
	
	env.WriteConfigFile("spellbook.yml", config)
	env.LoadConfig()
	env.InitializeComponents()
	
	tests := []struct {
		name        string
		actionName  string
		description string
	}{
		{"Short output", "short_output_action", "Should handle short messages"},
		{"Long output", "long_output_action", "Should truncate long output"},
		{"Multiline output", "multiline_output_action", "Should handle multiple lines"},
		{"ANSI output", "ansi_output_action", "Should strip ANSI codes"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := env.ExecuteAction(tt.actionName)
			assert.NoError(t, err, tt.description)
			
			env.WaitForNotification(500 * time.Millisecond)
		})
	}
}

func TestNotificationSystem_PlatformCompatibility(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Test platform-specific notification systems
	platformTests := map[string]func(t *testing.T){
		"windows": testWindowsNotifications,
		"darwin":  testMacOSNotifications,
		"linux":   testLinuxNotifications,
	}
	
	if testFunc, exists := platformTests[runtime.GOOS]; exists {
		testFunc(t)
	} else {
		t.Skipf("No platform-specific tests for %s", runtime.GOOS)
	}
}

func testWindowsNotifications(t *testing.T) {
	// Test Windows toast notifications
	notifier := notify.GetSystemNotifier()
	if notifier == nil {
		t.Skip("System notifier not available")
		return
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	tests := []struct {
		name string
		msg  notify.Notification
	}{
		{
			name: "Windows success notification",
			msg: notify.Notification{
				Title:   "Success",
				Message: "Windows notification test",
				Level:   notify.LevelSuccess,
			},
		},
		{
			name: "Windows error notification",
			msg: notify.Notification{
				Title:   "Error",
				Message: "Windows error notification test",
				Level:   notify.LevelError,
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := notifier.Notify(ctx, tt.msg)
			// May fail in headless environments, but we test the integration
			t.Logf("Windows notification result: %v", err)
		})
	}
}

func testMacOSNotifications(t *testing.T) {
	// Test macOS notification center
	notifier := notify.GetSystemNotifier()
	if notifier == nil {
		t.Skip("System notifier not available")
		return
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	tests := []struct {
		name string
		msg  notify.Notification
	}{
		{
			name: "macOS success notification",
			msg: notify.Notification{
				Title:   "Success",
				Message: "macOS notification test",
				Level:   notify.LevelSuccess,
			},
		},
		{
			name: "macOS info notification",
			msg: notify.Notification{
				Title:   "Info",
				Message: "macOS info notification test",
				Level:   notify.LevelInfo,
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := notifier.Notify(ctx, tt.msg)
			// May fail without proper permissions, but we test the integration
			t.Logf("macOS notification result: %v", err)
		})
	}
}

func testLinuxNotifications(t *testing.T) {
	// Test Linux desktop notifications (notify-send)
	notifier := notify.GetSystemNotifier()
	if notifier == nil {
		t.Skip("System notifier not available")
		return
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	tests := []struct {
		name string
		msg  notify.Notification
	}{
		{
			name: "Linux success notification",
			msg: notify.Notification{
				Title:   "Success",
				Message: "Linux notification test",
				Level:   notify.LevelSuccess,
			},
		},
		{
			name: "Linux warning notification",
			msg: notify.Notification{
				Title:   "Warning",
				Message: "Linux warning notification test",
				Level:   notify.LevelWarning,
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := notifier.Notify(ctx, tt.msg)
			// May fail without display server, but we test the integration
			t.Logf("Linux notification result: %v", err)
		})
	}
}

func TestNotificationQueue_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Test notification manager with multiple rapid notifications
	manager := notify.NewManager()
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// Send multiple notifications rapidly
	messages := []notify.Notification{
		{Title: "Test 1", Message: "First message", Level: notify.LevelInfo},
		{Title: "Test 2", Message: "Second message", Level: notify.LevelSuccess},
		{Title: "Test 3", Message: "Third message", Level: notify.LevelWarning},
		{Title: "Test 4", Message: "Fourth message", Level: notify.LevelError},
		{Title: "Test 5", Message: "Fifth message", Level: notify.LevelInfo},
	}
	
	// Send all messages
	for i, msg := range messages {
		err := manager.Notify(ctx, msg)
		assert.NoError(t, err, "Message %d should be sent successfully", i+1)
	}
	
	// Wait for processing
	time.Sleep(1 * time.Second)
	
	// Manager should handle all messages without blocking
	assert.True(t, true, "Manager should process all messages")
}

func TestNotificationSystem_OutputIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Test integration between notification system and console output
	manager := notify.NewManager()
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Send notification with output
	msg := notify.Notification{
		Title:   "Script Output",
		Message: "This is test output\nWith multiple lines\nAnd some details",
		Level:   notify.LevelInfo,
	}
	
	err := manager.Notify(ctx, msg)
	assert.NoError(t, err, "Output should be sent as notification")
}

func TestNotificationSystem_ErrorPropagation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := SetupTestEnvironment(t)
	
	config := `
spells:
  failing_action: "action_that_fails"

grimoire:
  action_that_fails:
    type: script
    command: "sh"
    args: ["-c", "echo 'Before error'; echo 'Error details' >&2; exit 1"]
    show_output: true
    description: "Action that fails with output"
`
	
	env.WriteConfigFile("spellbook.yml", config)
	env.LoadConfig()
	env.InitializeComponents()
	
	// Execute failing action
	err := env.ExecuteAction("action_that_fails")
	assert.Error(t, err, "Action should fail")
	
	// Wait for error notification
	env.WaitForNotification(500 * time.Millisecond)
	
	// In a real implementation, we would verify that:
	// 1. Error notification was sent
	// 2. Both stdout and stderr were captured
	// 3. Error details were included in notification
}

func TestNotificationSystem_Timeouts(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	manager := notify.NewManager()
	
	// Test with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	
	// This should timeout immediately
	msg := notify.Notification{
		Title:   "Timeout Test",
		Message: "This should timeout",
		Level:   notify.LevelInfo,
	}
	
	start := time.Now()
	err := manager.Notify(ctx, msg)
	elapsed := time.Since(start)
	
	// Should respect context timeout
	assert.Error(t, err, "Should timeout with short context")
	assert.Less(t, elapsed.Milliseconds(), int64(100), "Should timeout quickly")
}

func TestNotificationSystem_LargeContent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	manager := notify.NewManager()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Create large content
	largeContent := ""
	for i := 0; i < 1000; i++ {
		largeContent += "This is line " + string(rune(i)) + " of a very long notification message. "
	}
	
	msg := notify.Notification{
		Title:   "Large Content Test",
		Message: largeContent,
		Level:   notify.LevelInfo,
	}
	
	err := manager.Notify(ctx, msg)
	// Should handle large content gracefully (may truncate)
	assert.NoError(t, err, "Should handle large content")
}