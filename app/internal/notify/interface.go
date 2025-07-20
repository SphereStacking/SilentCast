package notify

import (
	"context"
	"fmt"
)

// Level represents the notification level
type Level int

const (
	// LevelInfo is for informational messages
	LevelInfo Level = iota
	// LevelWarning is for warning messages
	LevelWarning
	// LevelError is for error messages
	LevelError
	// LevelSuccess is for success messages
	LevelSuccess
)

// String returns the string representation of the Level
func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "Info"
	case LevelWarning:
		return "Warning"
	case LevelError:
		return "Error"
	case LevelSuccess:
		return "Success"
	default:
		return "Unknown"
	}
}

// Notification represents a notification to be displayed
type Notification struct {
	Title   string
	Message string
	Level   Level
}

// OutputNotification represents a notification with command output
type OutputNotification struct {
	Notification
	Output         string // Raw command output
	TruncatedBytes int    // Number of bytes truncated (0 if not truncated)
	ExitCode       int    // Command exit code (-1 if not applicable)
}

// TimeoutNotification represents a timeout-specific notification
type TimeoutNotification struct {
	Notification
	ActionName      string // Name of the action that timed out
	TimeoutDuration int    // Timeout duration in seconds
	ElapsedTime     int    // Actual elapsed time in seconds
	WasGraceful     bool   // Whether process exited gracefully during grace period
	Output          string // Partial output captured before timeout
}

// UpdateNotification represents an update-related notification
type UpdateNotification struct {
	Notification
	CurrentVersion string    // Current application version
	NewVersion     string    // Available update version
	ReleaseNotes   string    // Release notes/changelog
	DownloadSize   int64     // Download size in bytes
	PublishedAt    string    // Release publication date
	DownloadURL    string    // Direct download URL
	Actions        []string  // Available actions: "update", "dismiss", "remind"
}

// UpdateAction represents an action that can be taken on an update notification
type UpdateAction string

const (
	UpdateActionUpdate  UpdateAction = "update"  // Start update process
	UpdateActionDismiss UpdateAction = "dismiss" // Dismiss notification permanently
	UpdateActionRemind  UpdateAction = "remind"  // Remind later (snooze)
	UpdateActionView    UpdateAction = "view"    // View release notes
)

// NotificationOptions contains platform-specific notification options
type NotificationOptions struct {
	// MaxOutputLength limits the output size (0 for platform default)
	MaxOutputLength int
	// FormatAsCode formats output in monospace/code style if supported
	FormatAsCode bool
	// Priority affects notification display (platform-specific)
	Priority string // "low", "normal", "high", "critical"
	// Sound plays notification sound if supported
	Sound bool
}

// Notifier is the interface for sending notifications
type Notifier interface {
	// Notify sends a notification
	Notify(ctx context.Context, notification Notification) error

	// IsAvailable checks if the notifier is available on this system
	IsAvailable() bool
}

// OutputNotifier extends Notifier with support for rich output notifications
type OutputNotifier interface {
	Notifier

	// ShowWithOutput displays a notification with command output
	ShowWithOutput(ctx context.Context, notification OutputNotification) error

	// SetMaxOutputLength sets the maximum output length for notifications
	// Returns the actual limit applied (may be less than requested)
	SetMaxOutputLength(maxLength int) int

	// SupportsRichContent returns true if the notifier supports formatted output
	SupportsRichContent() bool
}

// UpdateNotifier extends Notifier with support for update notifications
type UpdateNotifier interface {
	Notifier

	// ShowUpdateNotification displays an update notification with actions
	ShowUpdateNotification(ctx context.Context, notification UpdateNotification) error

	// SupportsUpdateActions returns true if the notifier supports interactive actions
	SupportsUpdateActions() bool

	// OnUpdateAction handles user actions on update notifications
	OnUpdateAction(action UpdateAction, updateInfo UpdateNotification) error
}

// Manager manages multiple notifiers
type Manager struct {
	notifiers []Notifier
	options   NotificationOptions
}

// NewManager creates a new notification manager
func NewManager() *Manager {
	manager := &Manager{
		notifiers: make([]Notifier, 0),
		options: NotificationOptions{
			MaxOutputLength: 1024, // Default 1KB
			FormatAsCode:    true,
			Sound:           true,
		},
	}

	// Add available notifiers
	// Console notifier is always available
	manager.AddNotifier(NewConsoleNotifier())

	// Add platform-specific system notifier
	if systemNotifier := GetSystemNotifier(); systemNotifier != nil {
		manager.AddNotifier(systemNotifier)
	}

	return manager
}

// AddNotifier adds a notifier to the manager
func (m *Manager) AddNotifier(notifier Notifier) {
	if notifier != nil && notifier.IsAvailable() {
		m.notifiers = append(m.notifiers, notifier)
	}
}

// Notify sends a notification through all available notifiers
func (m *Manager) Notify(ctx context.Context, notification Notification) error {
	var lastError error

	for _, notifier := range m.notifiers {
		if err := notifier.Notify(ctx, notification); err != nil {
			lastError = err
		}
	}

	return lastError
}

// Info sends an info notification
func (m *Manager) Info(ctx context.Context, title, message string) error {
	return m.Notify(ctx, Notification{
		Title:   title,
		Message: message,
		Level:   LevelInfo,
	})
}

// Warning sends a warning notification
func (m *Manager) Warning(ctx context.Context, title, message string) error {
	return m.Notify(ctx, Notification{
		Title:   title,
		Message: message,
		Level:   LevelWarning,
	})
}

// Error sends an error notification
func (m *Manager) Error(ctx context.Context, title, message string) error {
	return m.Notify(ctx, Notification{
		Title:   title,
		Message: message,
		Level:   LevelError,
	})
}

// Success sends a success notification
func (m *Manager) Success(ctx context.Context, title, message string) error {
	return m.Notify(ctx, Notification{
		Title:   title,
		Message: message,
		Level:   LevelSuccess,
	})
}

// GetAvailableNotifiers returns a list of available notifier names
func (m *Manager) GetAvailableNotifiers() []string {
	var names []string
	for _, notifier := range m.notifiers {
		if notifier.IsAvailable() {
			// Get type name based on actual type
			switch n := notifier.(type) {
			case *ConsoleNotifier:
				names = append(names, "console")
			default:
				// For platform-specific notifiers, use a generic name
				_ = n
				names = append(names, "system")
			}
		}
	}
	return names
}

// GetSystemNotifier returns a platform-specific system notifier
// This function is implemented in platform-specific files
func GetSystemNotifier() Notifier {
	return getSystemNotifier()
}

// getSystemNotifier is implemented per platform
var getSystemNotifier func() Notifier

// SetOptions updates the notification options
func (m *Manager) SetOptions(options NotificationOptions) {
	m.options = options
}

// GetOptions returns the current notification options
func (m *Manager) GetOptions() NotificationOptions {
	return m.options
}

// NotifyWithOutput sends an output notification through all capable notifiers
func (m *Manager) NotifyWithOutput(ctx context.Context, notification OutputNotification) error {
	var lastError error
	notified := false

	// Truncate output if needed
	if m.options.MaxOutputLength > 0 && len(notification.Output) > m.options.MaxOutputLength {
		notification.TruncatedBytes = len(notification.Output) - m.options.MaxOutputLength
		notification.Output = notification.Output[:m.options.MaxOutputLength]
	}

	for _, notifier := range m.notifiers {
		// Try to use OutputNotifier interface if available
		if outputNotifier, ok := notifier.(OutputNotifier); ok {
			if err := outputNotifier.ShowWithOutput(ctx, notification); err != nil {
				lastError = err
			} else {
				notified = true
			}
		} else {
			// Fallback to regular notification with output in message
			fallback := notification.Notification
			if notification.Output != "" {
				fallback.Message += "\n\nOutput:\n" + notification.Output
				if notification.TruncatedBytes > 0 {
					fallback.Message += fmt.Sprintf("\n... (%d bytes truncated)", notification.TruncatedBytes)
				}
			}
			if err := notifier.Notify(ctx, fallback); err != nil {
				lastError = err
			} else {
				notified = true
			}
		}
	}

	if !notified && lastError != nil {
		return lastError
	}
	return nil
}

// SupportsOutputNotifications returns true if any notifier supports output notifications
func (m *Manager) SupportsOutputNotifications() bool {
	for _, notifier := range m.notifiers {
		if _, ok := notifier.(OutputNotifier); ok {
			return true
		}
	}
	return false
}

// GetOutputNotifiers returns all notifiers that support output notifications
func (m *Manager) GetOutputNotifiers() []OutputNotifier {
	var outputNotifiers []OutputNotifier
	for _, notifier := range m.notifiers {
		if on, ok := notifier.(OutputNotifier); ok {
			outputNotifiers = append(outputNotifiers, on)
		}
	}
	return outputNotifiers
}

// NotifyTimeout sends a timeout-specific notification
func (m *Manager) NotifyTimeout(ctx context.Context, notification TimeoutNotification) error {
	// Format a detailed timeout message
	message := fmt.Sprintf("Script execution timed out after %d seconds", notification.TimeoutDuration)

	if notification.WasGraceful {
		message += " (terminated gracefully)"
	} else {
		message += " (force terminated)"
	}

	// Add partial output if available
	if notification.Output != "" {
		maxLen := 200 // Show limited output in timeout notification
		output := notification.Output
		if len(output) > maxLen {
			output = output[:maxLen] + "..."
		}
		message += fmt.Sprintf("\n\nPartial output:\n%s", output)
	}

	// Override the base notification message
	notification.Notification.Message = message

	// Use OutputNotification for notifiers that support it
	if m.SupportsOutputNotifications() && notification.Output != "" {
		outputNotif := OutputNotification{
			Notification: notification.Notification,
			Output:       notification.Output,
			ExitCode:     -1, // Timeout exit code
		}
		return m.NotifyWithOutput(ctx, outputNotif)
	}

	// Otherwise use regular notification
	return m.Notify(ctx, notification.Notification)
}

// NotifyUpdate sends an update notification through all capable notifiers
func (m *Manager) NotifyUpdate(ctx context.Context, notification UpdateNotification) error {
	var lastError error
	notified := false

	// Format basic message if not provided
	if notification.Message == "" {
		notification.Message = fmt.Sprintf("Update available: %s → %s", 
			notification.CurrentVersion, notification.NewVersion)
	}

	for _, notifier := range m.notifiers {
		// Try to use UpdateNotifier interface if available
		if updateNotifier, ok := notifier.(UpdateNotifier); ok {
			if err := updateNotifier.ShowUpdateNotification(ctx, notification); err != nil {
				lastError = err
			} else {
				notified = true
			}
		} else {
			// Fallback to regular notification
			fallbackNotif := notification.Notification
			if notification.ReleaseNotes != "" {
				// Truncate release notes for regular notifications
				notes := notification.ReleaseNotes
				if len(notes) > 200 {
					notes = notes[:197] + "..."
				}
				fallbackNotif.Message += "\n\n" + notes
			}
			if err := notifier.Notify(ctx, fallbackNotif); err != nil {
				lastError = err
			} else {
				notified = true
			}
		}
	}

	if !notified && lastError != nil {
		return lastError
	}
	return nil
}

// SupportsUpdateNotifications returns true if any notifier supports update notifications
func (m *Manager) SupportsUpdateNotifications() bool {
	for _, notifier := range m.notifiers {
		if _, ok := notifier.(UpdateNotifier); ok {
			return true
		}
	}
	return false
}

// GetUpdateNotifiers returns all notifiers that support update notifications
func (m *Manager) GetUpdateNotifiers() []UpdateNotifier {
	var updateNotifiers []UpdateNotifier
	for _, notifier := range m.notifiers {
		if un, ok := notifier.(UpdateNotifier); ok {
			updateNotifiers = append(updateNotifiers, un)
		}
	}
	return updateNotifiers
}
