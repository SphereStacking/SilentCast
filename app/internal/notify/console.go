package notify

import (
	"context"
	"fmt"
	"os"
	"time"
)

// ConsoleNotifier sends notifications to the console
type ConsoleNotifier struct{}

// NewConsoleNotifier creates a new console notifier
func NewConsoleNotifier() *ConsoleNotifier {
	return &ConsoleNotifier{}
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
	if notification.Title != "" && notification.Message != "" {
		fmt.Fprintf(os.Stderr, "%s[%s] %s%s: %s%s\n", 
			color, timestamp, prefix, reset, notification.Title, reset)
		fmt.Fprintf(os.Stderr, "%s        %s%s\n", 
			color, notification.Message, reset)
	} else if notification.Title != "" {
		fmt.Fprintf(os.Stderr, "%s[%s] %s: %s%s\n", 
			color, timestamp, prefix, notification.Title, reset)
	} else {
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