//go:build windows

package notify

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/SphereStacking/silentcast/internal/errors"
	"github.com/SphereStacking/silentcast/pkg/logger"
)

// WindowsNotifier sends system notifications on Windows
type WindowsNotifier struct {
	appName         string
	useFallback     bool
	maxOutputLength int
}

// NewWindowsNotifier creates a new Windows notifier
func NewWindowsNotifier(appName string) *WindowsNotifier {
	return &WindowsNotifier{
		appName:         appName,
		maxOutputLength: 1000, // Default limit for Windows toast notifications
	}
}

// Notify sends a notification using Windows toast notifications
func (n *WindowsNotifier) Notify(ctx context.Context, notification Notification) error {
	// Try PowerShell toast notification first (Windows 10/11)
	if err := n.sendToastNotification(ctx, notification); err != nil {
		logger.Debug("Toast notification failed, trying fallback: %v", err)

		// Fallback to msg command
		if err := n.sendMsgNotification(ctx, notification); err != nil {
			logger.Debug("Msg notification failed: %v", err)

			// Final fallback to balloon notification
			if err := n.sendBalloonNotification(ctx, notification); err != nil {
				logger.Debug("Balloon notification failed: %v", err)

				// Return unified error with context
				return errors.New(errors.ErrorTypeSystem, "all notification methods failed").
					WithContext("notification_title", notification.Title).
					WithContext("platform", "windows").
					WithContext("tried_methods", "toast,msg,balloon").
					WithContext("suggested_action", "check Windows notification settings")
			}
		}
	}

	return nil
}

// sendToastNotification sends a Windows 10/11 toast notification using PowerShell
func (n *WindowsNotifier) sendToastNotification(ctx context.Context, notification Notification) error {
	// Build PowerShell script for toast notification
	title := escapeForPowerShell(notification.Title)
	message := escapeForPowerShell(notification.Message)

	// Icon could be used in future versions
	// Currently using default Windows toast notification icon

	// PowerShell script to show toast notification
	script := fmt.Sprintf(`
[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
[Windows.UI.Notifications.ToastNotification, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
[Windows.Data.Xml.Dom.XmlDocument, Windows.Data.Xml.Dom.XmlDocument, ContentType = WindowsRuntime] | Out-Null

$APP_ID = '%s'

$template = @"
<toast>
    <visual>
        <binding template="ToastGeneric">
            <text>%s</text>
            <text>%s</text>
        </binding>
    </visual>
</toast>
"@

$xml = New-Object Windows.Data.Xml.Dom.XmlDocument
$xml.LoadXml($template)
$toast = New-Object Windows.UI.Notifications.ToastNotification $xml
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier($APP_ID).Show($toast)
`, n.appName, title, message)

	cmd := exec.CommandContext(ctx, "powershell", "-NoProfile", "-NonInteractive", "-Command", script)
	if err := cmd.Run(); err != nil {
		return errors.Wrap(errors.ErrorTypeSystem, "toast notification failed", err).
			WithContext("notification_title", notification.Title).
			WithContext("method", "toast").
			WithContext("script_length", len(script))
	}
	return nil
}

// sendBalloonNotification sends a system tray balloon notification (older Windows)
func (n *WindowsNotifier) sendBalloonNotification(ctx context.Context, notification Notification) error {
	// Create a simple VBScript for balloon notification
	vbscript := fmt.Sprintf(`
Set objShell = CreateObject("WScript.Shell")
objShell.Popup "%s" & vbCrLf & "%s", 10, "%s", %d
`,
		escapeForVBScript(notification.Title),
		escapeForVBScript(notification.Message),
		n.appName,
		n.getVBScriptIcon(notification.Level))

	// Write VBScript to temp file
	tempDir := os.TempDir()
	scriptPath := filepath.Join(tempDir, "silentcast_notify.vbs")

	if err := os.WriteFile(scriptPath, []byte(vbscript), 0o600); err != nil {
		return errors.Wrap(errors.ErrorTypeSystem, "failed to write VBScript", err).
			WithContext("notification_title", notification.Title).
			WithContext("method", "balloon").
			WithContext("script_path", scriptPath)
	}
	defer os.Remove(scriptPath)

	// Execute VBScript
	cmd := exec.CommandContext(ctx, "cscript", "//NoLogo", scriptPath)
	if err := cmd.Run(); err != nil {
		return errors.Wrap(errors.ErrorTypeSystem, "balloon notification failed", err).
			WithContext("notification_title", notification.Title).
			WithContext("method", "balloon").
			WithContext("script_path", scriptPath)
	}
	return nil
}

// sendMsgNotification sends a message using the msg command (requires appropriate permissions)
func (n *WindowsNotifier) sendMsgNotification(ctx context.Context, notification Notification) error {
	// msg command is limited but works on most Windows systems
	message := fmt.Sprintf("%s\n%s", notification.Title, notification.Message)

	// Try to send to current user
	cmd := exec.CommandContext(ctx, "msg", "*", "/TIME:10", message)
	if err := cmd.Run(); err != nil {
		return errors.Wrap(errors.ErrorTypeSystem, "msg notification failed", err).
			WithContext("notification_title", notification.Title).
			WithContext("method", "msg").
			WithContext("message_length", len(message))
	}
	return nil
}

// IsAvailable checks if Windows notifications are available
func (n *WindowsNotifier) IsAvailable() bool {
	// Check if we're on Windows 10 or later for toast notifications
	cmd := exec.Command("powershell", "-Command", "[Environment]::OSVersion.Version.Major")
	output, err := cmd.Output()
	if err == nil {
		version := strings.TrimSpace(string(output))
		if version >= "10" {
			return true
		}
	}

	// Check for older notification methods
	if _, err := exec.LookPath("msg"); err == nil {
		return true
	}

	if _, err := exec.LookPath("cscript"); err == nil {
		return true
	}

	return false
}

// escapeForPowerShell escapes strings for PowerShell
func escapeForPowerShell(s string) string {
	s = strings.ReplaceAll(s, "'", "''")
	s = strings.ReplaceAll(s, "\n", "`n")
	s = strings.ReplaceAll(s, "\r", "`r")
	return s
}

// escapeForVBScript escapes strings for VBScript
func escapeForVBScript(s string) string {
	s = strings.ReplaceAll(s, `"`, `""`)
	s = strings.ReplaceAll(s, "\n", `" & vbCrLf & "`)
	return s
}

// getVBScriptIcon returns the icon code for VBScript MsgBox
func (n *WindowsNotifier) getVBScriptIcon(level Level) int {
	switch level {
	case LevelError:
		return 16 // vbCritical
	case LevelWarning:
		return 48 // vbExclamation
	case LevelInfo:
		return 64 // vbInformation
	default:
		return 64 // vbInformation
	}
}

// ShowWithOutput displays a notification with command output
func (n *WindowsNotifier) ShowWithOutput(ctx context.Context, notification OutputNotification) error {
	// Format output for Windows toast notification
	message := notification.Message

	if notification.Output != "" {
		// Truncate output if needed
		output := notification.Output
		if len(output) > n.maxOutputLength {
			output = output[:n.maxOutputLength] + "..."
		}

		// Add output to message
		message += "\n\nOutput:\n" + output

		// Add truncation info if needed
		if notification.TruncatedBytes > 0 {
			message += fmt.Sprintf("\n... (%d bytes truncated)", notification.TruncatedBytes)
		}
	}

	// Add exit code for failed commands
	if notification.ExitCode != 0 && notification.ExitCode != -1 {
		message += fmt.Sprintf("\nExit code: %d", notification.ExitCode)
	}

	// Create enhanced notification
	enhancedNotification := Notification{
		Title:   notification.Title,
		Message: message,
		Level:   notification.Level,
	}

	return n.Notify(ctx, enhancedNotification)
}

// SetMaxOutputLength sets the maximum output length for notifications
func (n *WindowsNotifier) SetMaxOutputLength(maxLength int) int {
	if maxLength > 0 && maxLength <= 2000 { // Windows toast limit
		n.maxOutputLength = maxLength
	}
	return n.maxOutputLength
}

// SupportsRichContent returns true if the notifier supports formatted output
func (n *WindowsNotifier) SupportsRichContent() bool {
	// Windows toast notifications support basic formatting
	return true
}

func init() {
	// Register Windows notifier factory
	getSystemNotifier = func() Notifier {
		return NewWindowsNotifier("SilentCast")
	}
}
