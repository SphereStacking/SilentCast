package notify

import (
	"context"
	"fmt"
	"time"

	"github.com/SphereStacking/silentcast/internal/updater"
)

// UpdateNotificationManager handles update-related notifications
type UpdateNotificationManager struct {
	notifier *Manager
	config   UpdateNotificationConfig
}

// UpdateNotificationConfig contains settings for update notifications
type UpdateNotificationConfig struct {
	// Enabled controls whether update notifications are shown
	Enabled bool
	// CheckInterval controls how often to check for updates
	CheckInterval time.Duration
	// ShowOnStartup controls whether to check for updates on application startup
	ShowOnStartup bool
	// RemindInterval controls how often to remind about dismissed updates
	RemindInterval time.Duration
	// AutoCheck enables background update checking
	AutoCheck bool
	// IncludePreReleases includes pre-release versions in checks
	IncludePreReleases bool
}

// DefaultUpdateNotificationConfig returns default configuration
func DefaultUpdateNotificationConfig() UpdateNotificationConfig {
	return UpdateNotificationConfig{
		Enabled:            true,
		CheckInterval:      24 * time.Hour,  // Check daily
		ShowOnStartup:      true,
		RemindInterval:     7 * 24 * time.Hour, // Remind weekly
		AutoCheck:          true,
		IncludePreReleases: false,
	}
}

// NewUpdateNotificationManager creates a new update notification manager
func NewUpdateNotificationManager(notifier *Manager) *UpdateNotificationManager {
	return &UpdateNotificationManager{
		notifier: notifier,
		config:   DefaultUpdateNotificationConfig(),
	}
}

// SetConfig updates the notification configuration
func (m *UpdateNotificationManager) SetConfig(config UpdateNotificationConfig) {
	m.config = config
}

// GetConfig returns the current configuration
func (m *UpdateNotificationManager) GetConfig() UpdateNotificationConfig {
	return m.config
}

// NotifyUpdateAvailable sends a notification about an available update
func (m *UpdateNotificationManager) NotifyUpdateAvailable(ctx context.Context, currentVersion string, updateInfo *updater.UpdateInfo) error {
	if !m.config.Enabled {
		return nil
	}

	// Format download size
	var sizeStr string
	if updateInfo.Size > 0 {
		sizeStr = formatUpdateSize(updateInfo.Size)
	}

	// Format publication date
	publishedStr := updateInfo.PublishedAt.Format("2006-01-02")

	// Create update notification
	notification := UpdateNotification{
		Notification: Notification{
			Title:   "🚀 SilentCast Update Available",
			Message: fmt.Sprintf("Version %s is now available", updateInfo.Version),
			Level:   LevelInfo,
		},
		CurrentVersion: currentVersion,
		NewVersion:     updateInfo.Version,
		ReleaseNotes:   updateInfo.ReleaseNotes,
		DownloadSize:   updateInfo.Size,
		PublishedAt:    publishedStr,
		DownloadURL:    updateInfo.DownloadURL,
		Actions:        []string{"update", "view", "remind", "dismiss"},
	}

	// Add size to message if available
	if sizeStr != "" {
		notification.Message += fmt.Sprintf(" (%s)", sizeStr)
	}

	return m.notifier.NotifyUpdate(ctx, notification)
}

// NotifyUpdateCheckFailed sends a notification about update check failure
func (m *UpdateNotificationManager) NotifyUpdateCheckFailed(ctx context.Context, err error) error {
	if !m.config.Enabled {
		return nil
	}

	notification := Notification{
		Title:   "⚠️ Update Check Failed",
		Message: fmt.Sprintf("Failed to check for updates: %v", err),
		Level:   LevelWarning,
	}

	return m.notifier.Notify(ctx, notification)
}

// NotifyUpdateStarted sends a notification when update download starts
func (m *UpdateNotificationManager) NotifyUpdateStarted(ctx context.Context, version string) error {
	if !m.config.Enabled {
		return nil
	}

	notification := Notification{
		Title:   "⬇️ Downloading Update",
		Message: fmt.Sprintf("Downloading SilentCast %s...", version),
		Level:   LevelInfo,
	}

	return m.notifier.Notify(ctx, notification)
}

// NotifyUpdateComplete sends a notification when update is complete
func (m *UpdateNotificationManager) NotifyUpdateComplete(ctx context.Context, version string) error {
	if !m.config.Enabled {
		return nil
	}

	notification := Notification{
		Title:   "✅ Update Complete",
		Message: fmt.Sprintf("SilentCast has been updated to %s", version),
		Level:   LevelSuccess,
	}

	return m.notifier.Notify(ctx, notification)
}

// NotifyUpdateFailed sends a notification when update fails
func (m *UpdateNotificationManager) NotifyUpdateFailed(ctx context.Context, version string, err error) error {
	if !m.config.Enabled {
		return nil
	}

	notification := Notification{
		Title:   "❌ Update Failed",
		Message: fmt.Sprintf("Failed to update to %s: %v", version, err),
		Level:   LevelError,
	}

	return m.notifier.Notify(ctx, notification)
}

// NotifyNoUpdatesAvailable sends a notification when no updates are found (for manual checks)
func (m *UpdateNotificationManager) NotifyNoUpdatesAvailable(ctx context.Context, currentVersion string) error {
	if !m.config.Enabled {
		return nil
	}

	notification := Notification{
		Title:   "✅ No Updates Available",
		Message: fmt.Sprintf("You're running the latest version (%s)", currentVersion),
		Level:   LevelInfo,
	}

	return m.notifier.Notify(ctx, notification)
}

// StartPeriodicChecks starts background update checking
func (m *UpdateNotificationManager) StartPeriodicChecks(ctx context.Context, updater *updater.Updater, currentVersion string) {
	if !m.config.Enabled || !m.config.AutoCheck {
		return
	}

	go func() {
		ticker := time.NewTicker(m.config.CheckInterval)
		defer ticker.Stop()

		// Initial check after a short delay (if ShowOnStartup is enabled)
		if m.config.ShowOnStartup {
			time.Sleep(30 * time.Second) // Wait for app to fully start
			m.checkAndNotify(ctx, updater, currentVersion)
		}

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				m.checkAndNotify(ctx, updater, currentVersion)
			}
		}
	}()
}

// checkAndNotify performs an update check and sends notifications if needed
func (m *UpdateNotificationManager) checkAndNotify(ctx context.Context, upd *updater.Updater, currentVersion string) {
	updateInfo, err := upd.CheckForUpdate(ctx)
	if err != nil {
		// Only notify about failures in debug mode or for manual checks
		// Auto-checks shouldn't spam users with network error notifications
		return
	}

	if updateInfo != nil {
		// Update available - send notification
		if err := m.NotifyUpdateAvailable(ctx, currentVersion, updateInfo); err != nil {
			// Log error but don't fail the check
			_ = err
		}
	}
	// No notification for "no updates available" during auto-checks
}

// HandleUpdateAction processes user actions on update notifications
func (m *UpdateNotificationManager) HandleUpdateAction(ctx context.Context, action UpdateAction, updateInfo UpdateNotification, upd *updater.Updater) error {
	switch action {
	case UpdateActionUpdate:
		return m.handleUpdateAction(ctx, updateInfo, upd)
	case UpdateActionView:
		return m.handleViewAction(ctx, updateInfo)
	case UpdateActionRemind:
		return m.handleRemindAction(ctx, updateInfo)
	case UpdateActionDismiss:
		return m.handleDismissAction(ctx, updateInfo)
	default:
		return fmt.Errorf("unknown update action: %s", action)
	}
}

// handleUpdateAction starts the update process
func (m *UpdateNotificationManager) handleUpdateAction(ctx context.Context, updateInfo UpdateNotification, upd *updater.Updater) error {
	// Notify that update is starting
	if err := m.NotifyUpdateStarted(ctx, updateInfo.NewVersion); err != nil {
		return err
	}

	// Create UpdateInfo struct for the updater
	updaterInfo := &updater.UpdateInfo{
		Version:      updateInfo.NewVersion,
		ReleaseNotes: updateInfo.ReleaseNotes,
		PublishedAt:  time.Now(), // We could parse the string, but this is sufficient
		DownloadURL:  updateInfo.DownloadURL,
		Size:         updateInfo.DownloadSize,
	}

	// Download update
	downloadPath, err := upd.DownloadUpdate(ctx, updaterInfo)
	if err != nil {
		// Best effort notification - failure doesn't prevent update process
		_ = m.NotifyUpdateFailed(ctx, updateInfo.NewVersion, err)
		return fmt.Errorf("failed to download update: %w", err)
	}

	// Apply update
	if err := upd.ApplyUpdate(downloadPath); err != nil {
		// Best effort notification - failure doesn't prevent update process
		_ = m.NotifyUpdateFailed(ctx, updateInfo.NewVersion, err)
		return fmt.Errorf("failed to apply update: %w", err)
	}

	// Notify success
	if err := m.NotifyUpdateComplete(ctx, updateInfo.NewVersion); err != nil {
		return err
	}

	return nil
}

// handleViewAction shows release notes
func (m *UpdateNotificationManager) handleViewAction(ctx context.Context, updateInfo UpdateNotification) error {
	// For now, just send a notification with full release notes
	// In the future, this could open a web browser or dedicated viewer
	notification := Notification{
		Title:   fmt.Sprintf("📋 Release Notes - %s", updateInfo.NewVersion),
		Message: updateInfo.ReleaseNotes,
		Level:   LevelInfo,
	}

	return m.notifier.Notify(ctx, notification)
}

// handleRemindAction schedules a reminder
func (m *UpdateNotificationManager) handleRemindAction(ctx context.Context, updateInfo UpdateNotification) error {
	// Schedule a reminder after the configured interval
	go func() {
		timer := time.NewTimer(m.config.RemindInterval)
		defer timer.Stop()

		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			// Send reminder notification
			reminderNotif := UpdateNotification{
				Notification: Notification{
					Title:   "🔔 Update Reminder",
					Message: fmt.Sprintf("Don't forget: SilentCast %s is available", updateInfo.NewVersion),
					Level:   LevelInfo,
				},
				CurrentVersion: updateInfo.CurrentVersion,
				NewVersion:     updateInfo.NewVersion,
				ReleaseNotes:   updateInfo.ReleaseNotes,
				DownloadSize:   updateInfo.DownloadSize,
				PublishedAt:    updateInfo.PublishedAt,
				DownloadURL:    updateInfo.DownloadURL,
				Actions:        []string{"update", "view", "dismiss"},
			}
			// Best effort reminder notification
			_ = m.notifier.NotifyUpdate(ctx, reminderNotif)
		}
	}()

	// Send confirmation
	notification := Notification{
		Title:   "⏰ Reminder Set",
		Message: fmt.Sprintf("You'll be reminded about this update in %v", m.config.RemindInterval),
		Level:   LevelInfo,
	}

	return m.notifier.Notify(ctx, notification)
}

// handleDismissAction dismisses the update notification
func (m *UpdateNotificationManager) handleDismissAction(ctx context.Context, updateInfo UpdateNotification) error {
	// For now, just send a confirmation
	// In the future, this could store dismissed versions to avoid showing them again
	notification := Notification{
		Title:   "🚫 Update Dismissed",
		Message: fmt.Sprintf("Update to %s has been dismissed", updateInfo.NewVersion),
		Level:   LevelInfo,
	}

	return m.notifier.Notify(ctx, notification)
}

// formatUpdateSize formats byte size to human-readable string
func formatUpdateSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	
	units := []string{"KB", "MB", "GB", "TB"}
	return fmt.Sprintf("%.1f %s", float64(size)/float64(div), units[exp])
}