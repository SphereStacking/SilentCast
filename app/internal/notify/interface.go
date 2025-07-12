package notify

import "context"

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

// Notification represents a notification to be displayed
type Notification struct {
	Title   string
	Message string
	Level   Level
}

// Notifier is the interface for sending notifications
type Notifier interface {
	// Notify sends a notification
	Notify(ctx context.Context, notification Notification) error

	// IsAvailable checks if the notifier is available on this system
	IsAvailable() bool
}

// Manager manages multiple notifiers
type Manager struct {
	notifiers []Notifier
}

// NewManager creates a new notification manager
func NewManager() *Manager {
	manager := &Manager{
		notifiers: make([]Notifier, 0),
	}

	// Add available notifiers
	// Console notifier is always available
	manager.AddNotifier(NewConsoleNotifier())

	// TODO: Add system notifiers (desktop notifications)

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
