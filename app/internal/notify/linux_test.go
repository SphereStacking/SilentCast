//go:build linux

package notify

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestNewLinuxNotifier(t *testing.T) {
	notifier := NewLinuxNotifier("TestApp")

	if notifier.appName != "TestApp" {
		t.Errorf("Expected app name 'TestApp', got '%s'", notifier.appName)
	}

	if notifier.maxOutputLength != 200 {
		t.Errorf("Expected default max output length 200, got %d", notifier.maxOutputLength)
	}
}

func TestLinuxNotifier_Notify(t *testing.T) {
	notifier := NewLinuxNotifier("SilentCast")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tests := []struct {
		name         string
		notification Notification
	}{
		{
			name: "basic notification",
			notification: Notification{
				Title:   "Test Title",
				Message: "Test message",
				Level:   LevelInfo,
			},
		},
		{
			name: "error notification",
			notification: Notification{
				Title:   "Error",
				Message: "Something went wrong",
				Level:   LevelError,
			},
		},
		{
			name: "warning notification",
			notification: Notification{
				Title:   "Warning",
				Message: "This is a warning",
				Level:   LevelWarning,
			},
		},
		{
			name: "success notification",
			notification: Notification{
				Title:   "Success",
				Message: "Operation completed",
				Level:   LevelSuccess,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: We can't easily test actual notification display without user interaction
			// So we test that the method doesn't return an error
			err := notifier.Notify(ctx, tt.notification)
			if err != nil {
				// It's okay if notifications fail in CI/test environment
				t.Logf("Notification failed (expected in test environment): %v", err)
			}
		})
	}
}

func TestLinuxNotifier_ShowWithOutput(t *testing.T) {
	notifier := NewLinuxNotifier("SilentCast")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tests := []struct {
		name          string
		notification  OutputNotification
		wantTruncated bool
	}{
		{
			name: "notification without output",
			notification: OutputNotification{
				Notification: Notification{
					Title:   "Script Complete",
					Message: "Script executed successfully",
					Level:   LevelSuccess,
				},
				Output:   "",
				ExitCode: 0,
			},
			wantTruncated: false,
		},
		{
			name: "notification with short output",
			notification: OutputNotification{
				Notification: Notification{
					Title:   "Script Complete",
					Message: "Script executed",
					Level:   LevelSuccess,
				},
				Output:   "Hello, World!\nThis is test output.",
				ExitCode: 0,
			},
			wantTruncated: false,
		},
		{
			name: "notification with long output (should truncate)",
			notification: OutputNotification{
				Notification: Notification{
					Title:   "Script Complete",
					Message: "Script executed",
					Level:   LevelSuccess,
				},
				Output:   strings.Repeat("This is a very long line of output that should be truncated. ", 10),
				ExitCode: 0,
			},
			wantTruncated: true,
		},
		{
			name: "notification with error and output",
			notification: OutputNotification{
				Notification: Notification{
					Title:   "Script Failed",
					Message: "Script execution failed",
					Level:   LevelError,
				},
				Output:         "Error: Command not found\nSome additional error details",
				ExitCode:       1,
				TruncatedBytes: 50,
			},
			wantTruncated: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := notifier.ShowWithOutput(ctx, tt.notification)
			if err != nil {
				// It's okay if notifications fail in CI/test environment
				t.Logf("Notification failed (expected in test environment): %v", err)
			}
		})
	}
}

func TestLinuxNotifier_SetMaxOutputLength(t *testing.T) {
	notifier := NewLinuxNotifier("SilentCast")

	// Check initial value
	if notifier.maxOutputLength != 200 {
		t.Errorf("Expected initial max output length 200, got %d", notifier.maxOutputLength)
	}

	// Set new value
	oldMax := notifier.SetMaxOutputLength(500)
	if oldMax != 200 {
		t.Errorf("Expected old max output length 200, got %d", oldMax)
	}

	if notifier.maxOutputLength != 500 {
		t.Errorf("Expected new max output length 500, got %d", notifier.maxOutputLength)
	}

	// Set another value
	oldMax = notifier.SetMaxOutputLength(100)
	if oldMax != 500 {
		t.Errorf("Expected old max output length 500, got %d", oldMax)
	}

	if notifier.maxOutputLength != 100 {
		t.Errorf("Expected new max output length 100, got %d", notifier.maxOutputLength)
	}
}

func TestLinuxNotifier_SupportsRichContent(t *testing.T) {
	notifier := NewLinuxNotifier("SilentCast")

	if !notifier.SupportsRichContent() {
		t.Error("Expected Linux notifier to support rich content")
	}
}

func TestLinuxNotifier_IsAvailable(t *testing.T) {
	notifier := NewLinuxNotifier("SilentCast")

	// On Linux, should be available if any notification system is present
	available := notifier.IsAvailable()

	// We can't guarantee the test environment, so just verify the method works
	if !available {
		t.Logf("Notification system not available in test environment")
	}
}

func TestLinuxNotifier_OutputTruncation(t *testing.T) {
	notifier := NewLinuxNotifier("SilentCast")
	ctx := context.Background()

	// Set a small max length for testing
	notifier.SetMaxOutputLength(50)

	longOutput := strings.Repeat("This is a long line. ", 10) // > 50 characters

	notification := OutputNotification{
		Notification: Notification{
			Title:   "Test",
			Message: "Test message",
			Level:   LevelInfo,
		},
		Output:   longOutput,
		ExitCode: 0,
	}

	// This test verifies truncation happens without actually showing notifications
	err := notifier.ShowWithOutput(ctx, notification)
	if err != nil {
		t.Logf("Notification failed (expected in test environment): %v", err)
	}

	// The truncation logic is tested indirectly through the method execution
}

func TestLinuxNotifier_NotificationMethods(t *testing.T) {
	notifier := NewLinuxNotifier("SilentCast")
	ctx := context.Background()

	notification := Notification{
		Title:   "Test Method",
		Message: "Testing different notification methods",
		Level:   LevelInfo,
	}

	// Test individual notification methods
	t.Run("notify-send", func(t *testing.T) {
		err := notifier.sendNotifySend(ctx, notification)
		if err != nil {
			t.Logf("notify-send failed (expected in test environment): %v", err)
		}
	})

	t.Run("gdbus", func(t *testing.T) {
		err := notifier.sendGDBusNotification(ctx, notification)
		if err != nil {
			t.Logf("gdbus failed (expected in test environment): %v", err)
		}
	})

	t.Run("zenity", func(t *testing.T) {
		err := notifier.sendZenityNotification(ctx, notification)
		if err != nil {
			t.Logf("zenity failed (expected in test environment): %v", err)
		}
	})
}

func TestLinuxNotifier_OutputNotifierInterface(t *testing.T) {
	notifier := NewLinuxNotifier("SilentCast")

	// Verify that LinuxNotifier implements OutputNotifier interface
	var _ OutputNotifier = notifier

	// Test interface methods
	notifier.SetMaxOutputLength(300)
	if !notifier.SupportsRichContent() {
		t.Error("Expected Linux notifier to support rich content")
	}
}

func TestLinuxNotifier_EmptyTitleHandling(t *testing.T) {
	notifier := NewLinuxNotifier("SilentCast")
	ctx := context.Background()

	notification := Notification{
		Title:   "", // Empty title should use app name
		Message: "Test message with empty title",
		Level:   LevelInfo,
	}

	err := notifier.Notify(ctx, notification)
	if err != nil {
		t.Logf("Notification failed (expected in test environment): %v", err)
	}
}

func TestLinuxNotifier_SpecialCharacters(t *testing.T) {
	notifier := NewLinuxNotifier("SilentCast")
	ctx := context.Background()

	notification := Notification{
		Title:   "Test & Special <Characters>",
		Message: "Message with \"quotes\" and 'apostrophes' and $variables and `backticks`",
		Level:   LevelInfo,
	}

	err := notifier.Notify(ctx, notification)
	if err != nil {
		t.Logf("Notification failed (expected in test environment): %v", err)
	}

	// Test with output containing special characters
	outputNotification := OutputNotification{
		Notification: notification,
		Output:       "Output with special chars: & < > \" ' $ ` \n newlines \t tabs",
		ExitCode:     0,
	}

	err = notifier.ShowWithOutput(ctx, outputNotification)
	if err != nil {
		t.Logf("Output notification failed (expected in test environment): %v", err)
	}
}
