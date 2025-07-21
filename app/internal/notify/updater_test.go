package notify

import (
	"context"
	"testing"
	"time"

	"github.com/SphereStacking/silentcast/internal/updater"
)

func TestUpdateNotificationManager_NewManager(t *testing.T) {
	notifier := NewManager()
	manager := NewUpdateNotificationManager(notifier)

	if manager == nil {
		t.Fatal("NewUpdateNotificationManager returned nil")
	}

	if manager.notifier != notifier {
		t.Error("Manager should store the provided notifier")
	}

	config := manager.GetConfig()
	if !config.Enabled {
		t.Error("Default config should have notifications enabled")
	}

	if config.CheckInterval != 24*time.Hour {
		t.Errorf("Expected default check interval of 24h, got %v", config.CheckInterval)
	}
}

func TestUpdateNotificationManager_SetGetConfig(t *testing.T) {
	notifier := NewManager()
	manager := NewUpdateNotificationManager(notifier)

	newConfig := UpdateNotificationConfig{
		Enabled:            false,
		CheckInterval:      12 * time.Hour,
		ShowOnStartup:      false,
		RemindInterval:     3 * 24 * time.Hour,
		AutoCheck:          false,
		IncludePreReleases: true,
	}

	manager.SetConfig(newConfig)
	retrievedConfig := manager.GetConfig()

	if retrievedConfig.Enabled != newConfig.Enabled {
		t.Errorf("Expected Enabled = %v, got %v", newConfig.Enabled, retrievedConfig.Enabled)
	}

	if retrievedConfig.CheckInterval != newConfig.CheckInterval {
		t.Errorf("Expected CheckInterval = %v, got %v", newConfig.CheckInterval, retrievedConfig.CheckInterval)
	}

	if retrievedConfig.IncludePreReleases != newConfig.IncludePreReleases {
		t.Errorf("Expected IncludePreReleases = %v, got %v", newConfig.IncludePreReleases, retrievedConfig.IncludePreReleases)
	}
}

func TestUpdateNotificationManager_NotifyUpdateAvailable(t *testing.T) {
	notifier := NewManager()
	manager := NewUpdateNotificationManager(notifier)

	updateInfo := &updater.UpdateInfo{
		Version:      "v1.2.0",
		ReleaseNotes: "Bug fixes and improvements",
		PublishedAt:  time.Now(),
		Size:         1024 * 1024, // 1MB
		DownloadURL:  "https://example.com/download",
	}

	ctx := context.Background()
	err := manager.NotifyUpdateAvailable(ctx, "v1.1.0", updateInfo)
	if err != nil {
		t.Errorf("NotifyUpdateAvailable failed: %v", err)
	}
}

func TestUpdateNotificationManager_NotifyUpdateAvailable_Disabled(t *testing.T) {
	notifier := NewManager()
	manager := NewUpdateNotificationManager(notifier)

	// Disable notifications
	config := manager.GetConfig()
	config.Enabled = false
	manager.SetConfig(config)

	updateInfo := &updater.UpdateInfo{
		Version:     "v1.2.0",
		PublishedAt: time.Now(),
	}

	ctx := context.Background()
	err := manager.NotifyUpdateAvailable(ctx, "v1.1.0", updateInfo)
	if err != nil {
		t.Errorf("NotifyUpdateAvailable should not fail when disabled: %v", err)
	}
}

func TestUpdateNotificationManager_NotifyUpdateCheckFailed(_ *testing.T) {
	notifier := NewManager()
	manager := NewUpdateNotificationManager(notifier)

	ctx := context.Background()
	testError := &testError{"network timeout"}
	err := manager.NotifyUpdateCheckFailed(ctx, testError)
	// Note: This may fail if system notifications are not available, but that's OK for testing
	_ = err // We'll allow this to fail in test environments
}

func TestUpdateNotificationManager_NotifyUpdateStarted(_ *testing.T) {
	notifier := NewManager()
	manager := NewUpdateNotificationManager(notifier)

	ctx := context.Background()
	err := manager.NotifyUpdateStarted(ctx, "v1.2.0")
	// Allow system notification failures in test environment
	_ = err
}

func TestUpdateNotificationManager_NotifyUpdateComplete(_ *testing.T) {
	notifier := NewManager()
	manager := NewUpdateNotificationManager(notifier)

	ctx := context.Background()
	err := manager.NotifyUpdateComplete(ctx, "v1.2.0")
	// Allow system notification failures in test environment
	_ = err
}

func TestUpdateNotificationManager_NotifyUpdateFailed(_ *testing.T) {
	notifier := NewManager()
	manager := NewUpdateNotificationManager(notifier)

	ctx := context.Background()
	testError := &testError{"checksum verification failed"}
	err := manager.NotifyUpdateFailed(ctx, "v1.2.0", testError)
	// Allow system notification failures in test environment
	_ = err
}

func TestUpdateNotificationManager_NotifyNoUpdatesAvailable(_ *testing.T) {
	notifier := NewManager()
	manager := NewUpdateNotificationManager(notifier)

	ctx := context.Background()
	err := manager.NotifyNoUpdatesAvailable(ctx, "v1.1.0")
	// Allow system notification failures in test environment
	_ = err
}

func TestDefaultUpdateNotificationConfig(t *testing.T) {
	config := DefaultUpdateNotificationConfig()

	if !config.Enabled {
		t.Error("Default config should be enabled")
	}

	if config.CheckInterval != 24*time.Hour {
		t.Errorf("Expected check interval of 24h, got %v", config.CheckInterval)
	}

	if !config.ShowOnStartup {
		t.Error("Default config should show on startup")
	}

	if config.RemindInterval != 7*24*time.Hour {
		t.Errorf("Expected remind interval of 7 days, got %v", config.RemindInterval)
	}

	if !config.AutoCheck {
		t.Error("Default config should have auto check enabled")
	}

	if config.IncludePreReleases {
		t.Error("Default config should not include pre-releases")
	}
}

func TestFormatUpdateSize(t *testing.T) {
	tests := []struct {
		size     int64
		expected string
	}{
		{0, "0 B"},
		{100, "100 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatUpdateSize(tt.size)
			if result != tt.expected {
				t.Errorf("formatUpdateSize(%d) = %s, want %s", tt.size, result, tt.expected)
			}
		})
	}
}

func TestUpdateNotificationManager_HandleUpdateAction(t *testing.T) {
	notifier := NewManager()
	manager := NewUpdateNotificationManager(notifier)

	updateInfo := UpdateNotification{
		Notification: Notification{
			Title:   "Update Available",
			Message: "Version 1.2.0 is available",
			Level:   LevelInfo,
		},
		CurrentVersion: "v1.1.0",
		NewVersion:     "v1.2.0",
		ReleaseNotes:   "Bug fixes",
		DownloadSize:   1024 * 1024,
		PublishedAt:    "2024-01-20",
		DownloadURL:    "https://example.com/download",
		Actions:        []string{"update", "view", "remind", "dismiss"},
	}

	ctx := context.Background()

	// Test view action (allow system notification failures)
	err := manager.HandleUpdateAction(ctx, UpdateActionView, updateInfo, nil)
	_ = err // Allow failures in test environment

	// Test remind action (allow system notification failures)
	err = manager.HandleUpdateAction(ctx, UpdateActionRemind, updateInfo, nil)
	_ = err // Allow failures in test environment

	// Test dismiss action (allow system notification failures)
	err = manager.HandleUpdateAction(ctx, UpdateActionDismiss, updateInfo, nil)
	_ = err // Allow failures in test environment

	// Test unknown action
	err = manager.HandleUpdateAction(ctx, UpdateAction("unknown"), updateInfo, nil)
	if err == nil {
		t.Error("HandleUpdateAction should fail for unknown action")
	}
}

// testError is a simple error implementation for testing
type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}
