//go:build darwin

package notify

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/SphereStacking/silentcast/internal/errors"
	"github.com/SphereStacking/silentcast/pkg/logger"
)

// DarwinNotifier sends system notifications on macOS
type DarwinNotifier struct {
	appName         string
	maxOutputLength int
}

// NewDarwinNotifier creates a new macOS notifier
func NewDarwinNotifier(appName string) *DarwinNotifier {
	return &DarwinNotifier{
		appName:         appName,
		maxOutputLength: 300, // Default maximum output length for macOS notifications
	}
}

// Notify sends a notification using macOS notification center
func (n *DarwinNotifier) Notify(ctx context.Context, notification Notification) error {
	// Try osascript first (most reliable)
	if err := n.sendOSAScriptNotification(ctx, notification); err != nil {
		logger.Debug("osascript notification failed: %v", err)

		// Try terminal-notifier if available
		if err := n.sendTerminalNotification(ctx, notification); err != nil {
			logger.Debug("terminal-notifier failed: %v", err)
			
			// Return unified error with context
			return errors.New(errors.ErrorTypeSystem, "all notification methods failed").
				WithContext("notification_title", notification.Title).
				WithContext("platform", "darwin").
				WithContext("tried_methods", "osascript,terminal-notifier").
				WithContext("suggested_action", "check notification permissions in System Preferences")
		}
	}

	return nil
}

// ShowWithOutput implements OutputNotifier interface for macOS
func (n *DarwinNotifier) ShowWithOutput(ctx context.Context, notification OutputNotification) error {
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
			message += fmt.Sprintf("\n\nExit code: %d", notification.ExitCode)
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
func (n *DarwinNotifier) SetMaxOutputLength(maxLength int) int {
	oldMax := n.maxOutputLength
	n.maxOutputLength = maxLength
	return oldMax
}

// SupportsRichContent returns whether macOS notifications support rich content
func (n *DarwinNotifier) SupportsRichContent() bool {
	// macOS notifications support basic formatting but not full rich content
	return true
}

// sendOSAScriptNotification sends a notification using osascript (AppleScript)
func (n *DarwinNotifier) sendOSAScriptNotification(ctx context.Context, notification Notification) error {
	// Build AppleScript command
	title := escapeForAppleScript(notification.Title)
	message := escapeForAppleScript(notification.Message)

	// Determine sound based on level
	sound := ""
	switch notification.Level {
	case LevelError:
		sound = "Basso"
	case LevelWarning:
		sound = "Hero"
	case LevelSuccess:
		sound = "Glass"
	default:
		sound = "Purr"
	}

	script := fmt.Sprintf(`display notification "%s" with title "%s" subtitle "%s" sound name "%s"`,
		message, n.appName, title, sound)

	cmd := exec.CommandContext(ctx, "osascript", "-e", script)
	if err := cmd.Run(); err != nil {
		return errors.Wrap(errors.ErrorTypeSystem, "osascript notification failed", err).
			WithContext("notification_title", notification.Title).
			WithContext("method", "osascript").
			WithContext("script_length", len(script))
	}
	return nil
}

// sendTerminalNotification sends a notification using terminal-notifier if installed
func (n *DarwinNotifier) sendTerminalNotification(ctx context.Context, notification Notification) error {
	// Check if terminal-notifier is available
	if _, err := exec.LookPath("terminal-notifier"); err != nil {
		return errors.Wrap(errors.ErrorTypeSystem, "terminal-notifier not found", err).
			WithContext("notification_title", notification.Title).
			WithContext("method", "terminal-notifier").
			WithContext("suggested_action", "install terminal-notifier via brew")
	}

	args := []string{
		"-title", n.appName,
		"-subtitle", notification.Title,
		"-message", notification.Message,
		"-group", "silentcast",
	}

	// Add sound based on level
	switch notification.Level {
	case LevelError:
		args = append(args, "-sound", "Basso")
	case LevelWarning:
		args = append(args, "-sound", "Hero")
	case LevelSuccess:
		args = append(args, "-sound", "Glass")
	}

	cmd := exec.CommandContext(ctx, "terminal-notifier", args...)
	if err := cmd.Run(); err != nil {
		return errors.Wrap(errors.ErrorTypeSystem, "terminal-notifier execution failed", err).
			WithContext("notification_title", notification.Title).
			WithContext("method", "terminal-notifier").
			WithContext("args_count", len(args))
	}
	return nil
}

// IsAvailable checks if macOS notifications are available
func (n *DarwinNotifier) IsAvailable() bool {
	// osascript should always be available on macOS
	if _, err := exec.LookPath("osascript"); err == nil {
		return true
	}

	// Check for terminal-notifier as fallback
	if _, err := exec.LookPath("terminal-notifier"); err == nil {
		return true
	}

	return false
}

// escapeForAppleScript escapes strings for AppleScript
func escapeForAppleScript(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	return s
}

func init() {
	// Register macOS notifier factory
	getSystemNotifier = func() Notifier {
		return NewDarwinNotifier("SilentCast")
	}
}
