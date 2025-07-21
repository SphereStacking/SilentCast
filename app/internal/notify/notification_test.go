package notify

import (
	"testing"
)

func TestNotification(t *testing.T) {
	// Test basic notification creation
	notif := Notification{
		Title:   "Test Title",
		Message: "Test Message",
		Level:   LevelInfo,
	}

	if notif.Title != "Test Title" {
		t.Errorf("Notification.Title = %v, want %v", notif.Title, "Test Title")
	}
	if notif.Message != "Test Message" {
		t.Errorf("Notification.Message = %v, want %v", notif.Message, "Test Message")
	}
	if notif.Level != LevelInfo {
		t.Errorf("Notification.Level = %v, want %v", notif.Level, LevelInfo)
	}
}

func TestOutputNotification(t *testing.T) {
	// Test notification with output
	notif := OutputNotification{
		Notification: Notification{
			Title:   "Command Output",
			Message: "Command completed",
			Level:   LevelSuccess,
		},
		Output:         "Hello World\nLine 2",
		ExitCode:       0,
		TruncatedBytes: 0,
	}

	if notif.Output != "Hello World\nLine 2" {
		t.Errorf("OutputNotification.Output = %v, want %v", notif.Output, "Hello World\nLine 2")
	}
	if notif.ExitCode != 0 {
		t.Errorf("OutputNotification.ExitCode = %v, want %v", notif.ExitCode, 0)
	}
	if notif.TruncatedBytes != 0 {
		t.Errorf("OutputNotification.TruncatedBytes = %v, want %v", notif.TruncatedBytes, 0)
	}
}

func TestTimeoutNotification(t *testing.T) {
	// Test timeout notification
	notif := TimeoutNotification{
		Notification: Notification{
			Title:   "Command Timeout",
			Message: "Command timed out",
			Level:   LevelError,
		},
		ActionName:      "long_script",
		TimeoutDuration: 30,
		ElapsedTime:     35,
		WasGraceful:     false,
		Output:          "Partial output...",
	}

	if notif.ActionName != "long_script" {
		t.Errorf("TimeoutNotification.ActionName = %v, want %v", notif.ActionName, "long_script")
	}
	if notif.TimeoutDuration != 30 {
		t.Errorf("TimeoutNotification.TimeoutDuration = %v, want %v", notif.TimeoutDuration, 30)
	}
	if notif.WasGraceful {
		t.Error("TimeoutNotification.WasGraceful should be false")
	}
}

func TestUpdateNotification(t *testing.T) {
	// Test update notification
	notif := UpdateNotification{
		Notification: Notification{
			Title:   "Update Available",
			Message: "New version is available",
			Level:   LevelInfo,
		},
		CurrentVersion: "1.0.0",
		NewVersion:     "1.1.0",
		ReleaseNotes:   "Bug fixes and improvements",
		DownloadSize:   1024 * 1024, // 1MB
		PublishedAt:    "2023-01-01",
		DownloadURL:    "https://example.com/download",
		Actions:        []string{"update", "dismiss", "remind"},
	}

	if notif.CurrentVersion != "1.0.0" {
		t.Errorf("UpdateNotification.CurrentVersion = %v, want %v", notif.CurrentVersion, "1.0.0")
	}
	if notif.NewVersion != "1.1.0" {
		t.Errorf("UpdateNotification.NewVersion = %v, want %v", notif.NewVersion, "1.1.0")
	}
	if notif.DownloadSize != 1024*1024 {
		t.Errorf("UpdateNotification.DownloadSize = %v, want %v", notif.DownloadSize, 1024*1024)
	}
	if len(notif.Actions) != 3 {
		t.Errorf("UpdateNotification.Actions length = %v, want %v", len(notif.Actions), 3)
	}
}

func TestNotificationLevels(t *testing.T) {
	levels := []Level{
		LevelInfo,
		LevelWarning,
		LevelError,
		LevelSuccess,
	}

	expectedStrings := []string{
		"Info",
		"Warning",
		"Error",
		"Success",
	}

	for i, level := range levels {
		notif := Notification{
			Title:   "Test",
			Message: "Test",
			Level:   level,
		}

		if notif.Level.String() != expectedStrings[i] {
			t.Errorf("Level %v should be %v", notif.Level.String(), expectedStrings[i])
		}
	}
}

func TestLevelString(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{LevelInfo, "Info"},
		{LevelWarning, "Warning"},
		{LevelError, "Error"},
		{LevelSuccess, "Success"},
		{Level(999), "Unknown"}, // Test unknown level
	}

	for _, tt := range tests {
		if got := tt.level.String(); got != tt.expected {
			t.Errorf("Level(%d).String() = %v, want %v", tt.level, got, tt.expected)
		}
	}
}

func TestUpdateAction(t *testing.T) {
	// Test update action constants
	actions := []UpdateAction{
		UpdateActionUpdate,
		UpdateActionDismiss,
		UpdateActionRemind,
		UpdateActionView,
	}

	expectedStrings := []string{
		"update",
		"dismiss",
		"remind",
		"view",
	}

	for i, action := range actions {
		if string(action) != expectedStrings[i] {
			t.Errorf("UpdateAction %v should be %v", action, expectedStrings[i])
		}
	}
}

func TestNotificationOptions(t *testing.T) {
	// Test notification options
	options := NotificationOptions{
		MaxOutputLength: 2048,
		FormatAsCode:    true,
		Priority:        "high",
		Sound:           false,
	}

	if options.MaxOutputLength != 2048 {
		t.Errorf("NotificationOptions.MaxOutputLength = %v, want %v", options.MaxOutputLength, 2048)
	}
	if !options.FormatAsCode {
		t.Error("NotificationOptions.FormatAsCode should be true")
	}
	if options.Priority != "high" {
		t.Errorf("NotificationOptions.Priority = %v, want %v", options.Priority, "high")
	}
	if options.Sound {
		t.Error("NotificationOptions.Sound should be false")
	}
}

func TestNotificationCopy(t *testing.T) {
	// Test that notification fields can be copied
	original := Notification{
		Title:   "Original Title",
		Message: "Original Message",
		Level:   LevelInfo,
	}

	// Create a copy by value
	notifCopy := original
	notifCopy.Title = "Modified Title"

	// Original should remain unchanged
	if original.Title != "Original Title" {
		t.Error("Original notification was modified when copying by value")
	}
	if notifCopy.Title != "Modified Title" {
		t.Error("Copy was not modified correctly")
	}
}

func TestNotificationValidation(t *testing.T) {
	// Test various notification field combinations
	tests := []struct {
		name  string
		notif Notification
		valid bool
	}{
		{
			name: "valid complete notification",
			notif: Notification{
				Title:   "Title",
				Message: "Message",
				Level:   LevelInfo,
			},
			valid: true,
		},
		{
			name: "empty title",
			notif: Notification{
				Title:   "",
				Message: "Message",
				Level:   LevelInfo,
			},
			valid: true, // Empty title might be allowed
		},
		{
			name: "empty message",
			notif: Notification{
				Title:   "Title",
				Message: "",
				Level:   LevelInfo,
			},
			valid: true, // Empty message might be allowed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This is mainly checking that the struct can be created
			// Real validation would be in the notifier implementations
			if tt.notif.Level == 0 && tt.notif.Title == "" && tt.notif.Message == "" {
				t.Log("Empty notification created")
			}

			// Basic sanity check - notification should have at least title or message
			hasContent := tt.notif.Title != "" || tt.notif.Message != ""
			if !hasContent && tt.valid {
				t.Error("Valid notification should have some content")
			}
		})
	}
}

func TestNotificationWithLongContent(t *testing.T) {
	// Test notification with very long content
	longTitle := "This is a very long title that might need to be truncated " +
		"depending on the notification system being used"
	longMessage := "This is a very long message that contains a lot of text " +
		"and might need to be handled specially by some notification systems " +
		"to avoid overwhelming the user with too much information at once"

	notif := Notification{
		Title:   longTitle,
		Message: longMessage,
		Level:   LevelInfo,
	}

	if len(notif.Title) == 0 {
		t.Error("Long title should be preserved")
	}
	if len(notif.Message) == 0 {
		t.Error("Long message should be preserved")
	}
}

func TestNotificationEquality(t *testing.T) {
	// Test notification equality
	notif1 := Notification{
		Title:   "Test",
		Message: "Message",
		Level:   LevelInfo,
	}

	notif2 := Notification{
		Title:   "Test",
		Message: "Message",
		Level:   LevelInfo,
	}

	notif3 := Notification{
		Title:   "Different",
		Message: "Message",
		Level:   LevelInfo,
	}

	// Manual equality check
	if notif1.Title != notif2.Title || notif1.Message != notif2.Message || notif1.Level != notif2.Level {
		t.Error("notif1 and notif2 should be equal")
	}

	if notif1.Title == notif3.Title {
		t.Error("notif1 and notif3 should not be equal")
	}
}

func TestNotificationEmbedding(t *testing.T) {
	// Test that embedded notifications work correctly
	output := OutputNotification{
		Notification: Notification{
			Title:   "Output Test",
			Message: "Command executed",
			Level:   LevelSuccess,
		},
		Output:   "command output",
		ExitCode: 0,
	}

	// Should be able to access embedded fields
	if output.Title != "Output Test" {
		t.Error("Should be able to access embedded Notification.Title")
	}
	if output.Level != LevelSuccess {
		t.Error("Should be able to access embedded Notification.Level")
	}

	timeout := TimeoutNotification{
		Notification: Notification{
			Title:   "Timeout Test",
			Message: "Command timed out",
			Level:   LevelError,
		},
		ActionName:      "test_action",
		TimeoutDuration: 30,
	}

	// Should be able to access embedded fields
	if timeout.Title != "Timeout Test" {
		t.Error("Should be able to access embedded Notification.Title")
	}
	if timeout.Level != LevelError {
		t.Error("Should be able to access embedded Notification.Level")
	}
}
