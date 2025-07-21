//go:build linux

package notify

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/SphereStacking/silentcast/internal/errors"
	"github.com/SphereStacking/silentcast/pkg/logger"
)

// LinuxNotifier sends system notifications on Linux
type LinuxNotifier struct {
	appName         string
	maxOutputLength int
}

// NewLinuxNotifier creates a new Linux notifier
func NewLinuxNotifier(appName string) *LinuxNotifier {
	return &LinuxNotifier{
		appName:         appName,
		maxOutputLength: 200, // Default maximum output length for Linux notifications
	}
}

// Notify sends a notification using Linux desktop notification system
func (n *LinuxNotifier) Notify(ctx context.Context, notification Notification) error {
	// Try notify-send first (most common)
	if err := n.sendNotifySend(ctx, notification); err != nil {
		logger.Debug("notify-send failed: %v", err)

		// Try gdbus for GNOME
		if err := n.sendGDBusNotification(ctx, notification); err != nil {
			logger.Debug("gdbus notification failed: %v", err)

			// Try zenity as last resort
			if err := n.sendZenityNotification(ctx, notification); err != nil {
				logger.Debug("zenity notification failed: %v", err)

				// Return unified error with context
				return errors.New(errors.ErrorTypeSystem, "all notification methods failed").
					WithContext("notification_title", notification.Title).
					WithContext("platform", "linux").
					WithContext("tried_methods", "notify-send,gdbus,zenity").
					WithContext("suggested_action", "install libnotify-bin or zenity package")
			}
		}
	}

	return nil
}

// ShowWithOutput implements OutputNotifier interface for Linux
func (n *LinuxNotifier) ShowWithOutput(ctx context.Context, notification OutputNotification) error {
	// Format the notification message
	message := notification.Message

	if notification.Output != "" {
		// Add command output
		output := notification.Output
		if len(output) > n.maxOutputLength {
			output = output[:n.maxOutputLength] + "..."
		}

		// Create a formatted message with output
		if message != "" {
			message += "\n\n"
		}
		message += "Output:\n" + output

		// Show exit code if command failed
		if notification.ExitCode != 0 {
			message += fmt.Sprintf("\nExit code: %d", notification.ExitCode)
		}

		// Show truncation info if applicable
		if notification.TruncatedBytes > 0 {
			message += fmt.Sprintf(" (%d bytes truncated)", notification.TruncatedBytes)
		}
	}

	// Create a regular notification with formatted message
	regularNotification := Notification{
		Title:   notification.Title,
		Message: message,
		Level:   notification.Level,
	}

	return n.Notify(ctx, regularNotification)
}

// SetMaxOutputLength sets the maximum output length for notifications
func (n *LinuxNotifier) SetMaxOutputLength(maxLength int) int {
	oldMax := n.maxOutputLength
	n.maxOutputLength = maxLength
	return oldMax
}

// SupportsRichContent returns whether Linux notifications support rich content
func (n *LinuxNotifier) SupportsRichContent() bool {
	// Linux desktop notifications support basic formatting
	return true
}

// sendNotifySend sends a notification using notify-send
func (n *LinuxNotifier) sendNotifySend(ctx context.Context, notification Notification) error {
	// Check if notify-send is available
	if _, err := exec.LookPath("notify-send"); err != nil {
		return errors.Wrap(errors.ErrorTypeSystem, "notify-send not found", err).
			WithContext("notification_title", notification.Title).
			WithContext("method", "notify-send").
			WithContext("suggested_action", "install libnotify-bin package")
	}

	args := []string{}

	// Set urgency based on level
	switch notification.Level {
	case LevelError:
		args = append(args, "-u", "critical")
	case LevelWarning:
		args = append(args, "-u", "normal")
	default:
		args = append(args, "-u", "low")
	}

	// Set icon based on level
	icon := "dialog-information"
	switch notification.Level {
	case LevelError:
		icon = "dialog-error"
	case LevelWarning:
		icon = "dialog-warning"
	case LevelSuccess:
		icon = "emblem-default"
	}
	args = append(args, "-i", icon)

	// Set app name
	args = append(args, "-a", n.appName)

	// Set expire time (10 seconds)
	args = append(args, "-t", "10000")

	// Add title and message
	title := notification.Title
	if title == "" {
		title = n.appName
	}
	args = append(args, title, notification.Message)

	cmd := exec.CommandContext(ctx, "notify-send", args...)
	if err := cmd.Run(); err != nil {
		return errors.Wrap(errors.ErrorTypeSystem, "notify-send execution failed", err).
			WithContext("notification_title", notification.Title).
			WithContext("method", "notify-send").
			WithContext("args_count", len(args))
	}
	return nil
}

// sendGDBusNotification sends a notification using gdbus (GNOME)
func (n *LinuxNotifier) sendGDBusNotification(ctx context.Context, notification Notification) error {
	// Check if gdbus is available
	if _, err := exec.LookPath("gdbus"); err != nil {
		return errors.Wrap(errors.ErrorTypeSystem, "gdbus not found", err).
			WithContext("notification_title", notification.Title).
			WithContext("method", "gdbus").
			WithContext("suggested_action", "install glib2.0-bin package or check GNOME installation")
	}

	// Determine icon based on level
	icon := "dialog-information"
	switch notification.Level {
	case LevelError:
		icon = "dialog-error"
	case LevelWarning:
		icon = "dialog-warning"
	case LevelSuccess:
		icon = "emblem-default"
	}

	// Build the notification ID (use timestamp)
	notificationID := uint32(time.Now().Unix())

	title := notification.Title
	if title == "" {
		title = n.appName
	}

	cmd := exec.CommandContext(ctx, "gdbus", "call",
		"--session",
		"--dest=org.freedesktop.Notifications",
		"--object-path=/org/freedesktop/Notifications",
		"--method=org.freedesktop.Notifications.Notify",
		n.appName,
		fmt.Sprintf("%d", notificationID),
		icon,
		title,
		notification.Message,
		"[]",    // actions
		"{}",    // hints
		"10000", // timeout in milliseconds
	)
	if err := cmd.Run(); err != nil {
		return errors.Wrap(errors.ErrorTypeSystem, "gdbus notification failed", err).
			WithContext("notification_title", notification.Title).
			WithContext("method", "gdbus").
			WithContext("notification_id", notificationID).
			WithContext("suggested_action", "check if notification service is running")
	}
	return nil
}

// sendZenityNotification sends a notification using zenity
func (n *LinuxNotifier) sendZenityNotification(ctx context.Context, notification Notification) error {
	// Check if zenity is available
	if _, err := exec.LookPath("zenity"); err != nil {
		return errors.Wrap(errors.ErrorTypeSystem, "zenity not found", err).
			WithContext("notification_title", notification.Title).
			WithContext("method", "zenity").
			WithContext("suggested_action", "install zenity package")
	}

	args := []string{"--notification"}

	// Set icon based on level
	switch notification.Level {
	case LevelError:
		args = append(args, "--window-icon=error")
	case LevelWarning:
		args = append(args, "--window-icon=warning")
	case LevelSuccess:
		args = append(args, "--window-icon=info")
	}

	// Build text
	text := notification.Title
	if notification.Message != "" {
		if text != "" {
			text += "\n"
		}
		text += notification.Message
	}
	args = append(args, "--text="+text)

	cmd := exec.CommandContext(ctx, "zenity", args...)
	if err := cmd.Run(); err != nil {
		return errors.Wrap(errors.ErrorTypeSystem, "zenity notification failed", err).
			WithContext("notification_title", notification.Title).
			WithContext("method", "zenity").
			WithContext("text_length", len(text))
	}
	return nil
}

// IsAvailable checks if Linux notifications are available
func (n *LinuxNotifier) IsAvailable() bool {
	// Check for various notification systems
	notificationSystems := []string{
		"notify-send", // Most common
		"gdbus",       // GNOME
		"zenity",      // GTK+ dialogs
	}

	for _, system := range notificationSystems {
		if _, err := exec.LookPath(system); err == nil {
			return true
		}
	}

	return false
}

func init() {
	// Register Linux notifier factory
	getSystemNotifier = func() Notifier {
		return NewLinuxNotifier("SilentCast")
	}
}
