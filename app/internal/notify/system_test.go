package notify

import (
	"context"
	"runtime"
	"testing"
)

func TestSystemNotifier(t *testing.T) {
	// Get system notifier
	notifier := GetSystemNotifier()

	// System notifier should be available on major platforms
	switch runtime.GOOS {
	case "windows", "darwin", "linux":
		if notifier == nil {
			t.Error("System notifier should be available on", runtime.GOOS)
		}

		// Test availability check
		if notifier != nil && !notifier.IsAvailable() {
			t.Log("System notifier not available on this system (may be running in CI)")
		}
	default:
		if notifier != nil {
			t.Error("System notifier should be nil on unsupported platform")
		}
	}
}

func TestManager_WithSystemNotifier(t *testing.T) {
	manager := NewManager()

	// Manager should have at least console notifier
	if len(manager.notifiers) < 1 {
		t.Error("Manager should have at least console notifier")
	}

	// On supported platforms, check if we have more than just console notifier
	switch runtime.GOOS {
	case "windows", "darwin", "linux":
		systemNotifier := GetSystemNotifier()
		if systemNotifier != nil && systemNotifier.IsAvailable() && len(manager.notifiers) < 2 {
			t.Error("System notifier should be added to manager on", runtime.GOOS)
		}
	}
}

func TestNotification_AllLevels(t *testing.T) {
	// Skip if no system notifier available
	notifier := GetSystemNotifier()
	if notifier == nil || !notifier.IsAvailable() {
		t.Skip("System notifier not available")
	}

	ctx := context.Background()

	tests := []struct {
		name  string
		level Level
		title string
		msg   string
	}{
		{
			name:  "Info notification",
			level: LevelInfo,
			title: "Test Info",
			msg:   "This is an info notification from SilentCast tests",
		},
		{
			name:  "Warning notification",
			level: LevelWarning,
			title: "Test Warning",
			msg:   "This is a warning notification from SilentCast tests",
		},
		{
			name:  "Error notification",
			level: LevelError,
			title: "Test Error",
			msg:   "This is an error notification from SilentCast tests",
		},
		{
			name:  "Success notification",
			level: LevelSuccess,
			title: "Test Success",
			msg:   "This is a success notification from SilentCast tests",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notification := Notification{
				Title:   tt.title,
				Message: tt.msg,
				Level:   tt.level,
			}

			err := notifier.Notify(ctx, notification)
			if err != nil {
				t.Logf("Notification failed (may be expected in test environment): %v", err)
			}
		})
	}
}
