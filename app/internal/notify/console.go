package notify

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	appErrors "github.com/SphereStacking/silentcast/internal/errors"
)

// ConsoleNotifier sends notifications to the console
type ConsoleNotifier struct {
	maxOutputLength int
}

// NewConsoleNotifier creates a new console notifier
func NewConsoleNotifier() *ConsoleNotifier {
	return &ConsoleNotifier{
		maxOutputLength: 2048, // 2KB default
	}
}

// Notify sends a notification to the console
func (n *ConsoleNotifier) Notify(ctx context.Context, notification Notification) error {
	// Format timestamp
	timestamp := time.Now().Format("15:04:05")

	// Choose prefix and color based on level
	var prefix, color, reset string

	// Check if we're in a terminal that supports colors
	if isTerminal() {
		reset = "\033[0m"
		switch notification.Level {
		case LevelInfo:
			prefix = "ℹ️  INFO"
			color = "\033[36m" // Cyan
		case LevelWarning:
			prefix = "⚠️  WARN"
			color = "\033[33m" // Yellow
		case LevelError:
			prefix = "❌ ERROR"
			color = "\033[31m" // Red
		case LevelSuccess:
			prefix = "✅ SUCCESS"
			color = "\033[32m" // Green
		default:
			prefix = "📢 NOTIFY"
			color = ""
		}
	} else {
		// No colors for non-terminal output
		switch notification.Level {
		case LevelInfo:
			prefix = "[INFO]"
		case LevelWarning:
			prefix = "[WARN]"
		case LevelError:
			prefix = "[ERROR]"
		case LevelSuccess:
			prefix = "[SUCCESS]"
		default:
			prefix = "[NOTIFY]"
		}
	}

	// Format and print the notification
	switch {
	case notification.Title != "" && notification.Message != "":
		fmt.Fprintf(os.Stderr, "%s[%s] %s%s: %s%s\n",
			color, timestamp, prefix, reset, notification.Title, reset)
		fmt.Fprintf(os.Stderr, "%s        %s%s\n",
			color, notification.Message, reset)
	case notification.Title != "":
		fmt.Fprintf(os.Stderr, "%s[%s] %s: %s%s\n",
			color, timestamp, prefix, notification.Title, reset)
	default:
		fmt.Fprintf(os.Stderr, "%s[%s] %s: %s%s\n",
			color, timestamp, prefix, notification.Message, reset)
	}

	return nil
}

// IsAvailable checks if the console notifier is available
func (n *ConsoleNotifier) IsAvailable() bool {
	// Console notifier is always available
	return true
}

// isTerminal checks if we're running in a terminal
func isTerminal() bool {
	fileInfo, err := os.Stderr.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// ShowWithOutput displays a notification with command output
func (n *ConsoleNotifier) ShowWithOutput(ctx context.Context, notification OutputNotification) error {
	// First show the regular notification
	if err := n.Notify(ctx, notification.Notification); err != nil {
		return err
	}

	// Then show the output if present
	if notification.Output != "" {
		timestamp := time.Now().Format("15:04:05")
		var color, reset string

		if isTerminal() {
			color = "\033[90m" // Dark gray for output
			reset = "\033[0m"
		}

		// Print output header
		fmt.Fprintf(os.Stderr, "%s[%s] Output:%s\n", color, timestamp, reset)

		// Print the output with indentation
		fmt.Fprintf(os.Stderr, "%s%s%s\n", color, indentOutput(notification.Output), reset)

		// Show truncation info if needed
		if notification.TruncatedBytes > 0 {
			fmt.Fprintf(os.Stderr, "%s... (%d bytes truncated)%s\n",
				color, notification.TruncatedBytes, reset)
		}

		// Show exit code if non-zero
		if notification.ExitCode != 0 && notification.ExitCode != -1 {
			exitColor := color
			if isTerminal() && notification.ExitCode != 0 {
				exitColor = "\033[31m" // Red for non-zero exit codes
			}
			fmt.Fprintf(os.Stderr, "%s[%s] Exit code: %d%s\n",
				exitColor, timestamp, notification.ExitCode, reset)
		}
	}

	return nil
}

// SetMaxOutputLength sets the maximum output length for notifications
func (n *ConsoleNotifier) SetMaxOutputLength(maxLength int) int {
	if maxLength > 0 {
		n.maxOutputLength = maxLength
	}
	return n.maxOutputLength
}

// SupportsRichContent returns true if the notifier supports formatted output
func (n *ConsoleNotifier) SupportsRichContent() bool {
	// Console supports rich content if in a terminal
	return isTerminal()
}

// ShowUpdateNotification displays an update notification with actions
func (n *ConsoleNotifier) ShowUpdateNotification(ctx context.Context, notification UpdateNotification) error {
	// Display the basic notification first
	if err := n.Notify(ctx, notification.Notification); err != nil {
		return err
	}

	timestamp := time.Now().Format("15:04:05")
	var color, reset string

	if isTerminal() {
		color = "\033[36m" // Cyan for update info
		reset = "\033[0m"
	}

	// Show detailed update information
	fmt.Fprintf(os.Stderr, "%s[%s] Update Details:%s\n", color, timestamp, reset)
	fmt.Fprintf(os.Stderr, "%s  Current: %s%s\n", color, notification.CurrentVersion, reset)
	fmt.Fprintf(os.Stderr, "%s  Latest:  %s%s\n", color, notification.NewVersion, reset)

	if notification.PublishedAt != "" {
		fmt.Fprintf(os.Stderr, "%s  Published: %s%s\n", color, notification.PublishedAt, reset)
	}

	if notification.DownloadSize > 0 {
		sizeStr := formatUpdateSize(notification.DownloadSize)
		fmt.Fprintf(os.Stderr, "%s  Size: %s%s\n", color, sizeStr, reset)
	}

	// Show truncated release notes
	if notification.ReleaseNotes != "" {
		fmt.Fprintf(os.Stderr, "%s[%s] Release Notes:%s\n", color, timestamp, reset)
		notes := notification.ReleaseNotes
		if len(notes) > 300 {
			notes = notes[:297] + "..."
		}
		fmt.Fprintf(os.Stderr, "%s%s%s\n", color, indentOutput(notes), reset)
	}

	// Show available actions
	if len(notification.Actions) > 0 {
		fmt.Fprintf(os.Stderr, "%s[%s] Available actions: %s%s\n",
			color, timestamp, strings.Join(notification.Actions, ", "), reset)
		fmt.Fprintf(os.Stderr, "%s  Run: ./silentcast --self-update%s\n", color, reset)
	}

	return nil
}

// SupportsUpdateActions returns true if the notifier supports interactive actions
func (n *ConsoleNotifier) SupportsUpdateActions() bool {
	// Console doesn't support interactive actions (it's display-only)
	return false
}

// OnUpdateAction handles user actions on update notifications
func (n *ConsoleNotifier) OnUpdateAction(action UpdateAction, updateInfo UpdateNotification) error {
	// Console notifier doesn't handle actions - they must be handled externally
	return appErrors.New(appErrors.ErrorTypeSystem, "console notifier does not support interactive actions").
		WithContext("notifier_type", "console").
		WithContext("action", action).
		WithContext("supported_actions", []string{})
}

// indentOutput adds indentation to each line of the output
func indentOutput(output string) string {
	lines := strings.Split(output, "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = "        " + line
		}
	}
	return strings.Join(lines, "\n")
}
