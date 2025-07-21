package notify

import (
	"context"
	"sync"

	customErrors "github.com/SphereStacking/silentcast/internal/errors"
)

// MockNotifier is a mock implementation of the Notifier interface for testing
type MockNotifier struct {
	mu            sync.Mutex
	notifications []Notification
	available     bool
	simulateError error
	maxOutputLen  int
	supportsRich  bool
	notifyFunc    func(context.Context, Notification) error // Custom notify function for testing
}

// NewMockNotifier creates a new mock notifier
func NewMockNotifier(available bool) *MockNotifier {
	return &MockNotifier{
		notifications: make([]Notification, 0),
		available:     available,
		maxOutputLen:  2048, // 2KB default
		supportsRich:  true,
	}
}

// Notify implements the Notifier interface
func (m *MockNotifier) Notify(ctx context.Context, notification Notification) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Use custom function if set
	if m.notifyFunc != nil {
		return m.notifyFunc(ctx, notification)
	}

	if m.simulateError != nil {
		// Check if the error already has context (is already a SpellbookError)
		if spellErr, ok := m.simulateError.(interface {
			WithContext(key string, value interface{}) *customErrors.SpellbookError
		}); ok {
			// Add notification context to the existing SpellbookError
			return spellErr.WithContext("notification_title", notification.Title)
		}
		// If it's not a SpellbookError, return as-is
		return m.simulateError
	}

	m.notifications = append(m.notifications, notification)
	return nil
}

// IsAvailable implements the Notifier interface
func (m *MockNotifier) IsAvailable() bool {
	return m.available
}

// GetNotifications returns all recorded notifications
func (m *MockNotifier) GetNotifications() []Notification {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Return a copy to avoid race conditions
	result := make([]Notification, len(m.notifications))
	copy(result, m.notifications)
	return result
}

// Clear clears all recorded notifications
func (m *MockNotifier) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.notifications = m.notifications[:0]
}

// SetError sets an error to simulate for the next notification
func (m *MockNotifier) SetError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.simulateError = err
}

// SetNotifyFunc sets a custom notify function for testing
func (m *MockNotifier) SetNotifyFunc(fn func(context.Context, Notification) error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.notifyFunc = fn
}

// MockOutputNotifier extends MockNotifier with OutputNotifier capabilities
type MockOutputNotifier struct {
	*MockNotifier
	outputNotifications []OutputNotification
}

// NewMockOutputNotifier creates a new mock output notifier
func NewMockOutputNotifier(available bool) *MockOutputNotifier {
	return &MockOutputNotifier{
		MockNotifier:        NewMockNotifier(available),
		outputNotifications: make([]OutputNotification, 0),
	}
}

// ShowWithOutput implements the OutputNotifier interface
func (m *MockOutputNotifier) ShowWithOutput(ctx context.Context, notification OutputNotification) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.simulateError != nil {
		return m.simulateError
	}

	// Also record as regular notification
	m.notifications = append(m.notifications, notification.Notification)
	m.outputNotifications = append(m.outputNotifications, notification)
	return nil
}

// SetMaxOutputLength implements the OutputNotifier interface
func (m *MockOutputNotifier) SetMaxOutputLength(maxLength int) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	if maxLength > 0 {
		m.maxOutputLen = maxLength
	}
	return m.maxOutputLen
}

// SupportsRichContent implements the OutputNotifier interface
func (m *MockOutputNotifier) SupportsRichContent() bool {
	return m.supportsRich
}

// GetOutputNotifications returns all recorded output notifications
func (m *MockOutputNotifier) GetOutputNotifications() []OutputNotification {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Return a copy to avoid race conditions
	result := make([]OutputNotification, len(m.outputNotifications))
	copy(result, m.outputNotifications)
	return result
}

// SetSupportsRichContent sets whether rich content is supported
func (m *MockOutputNotifier) SetSupportsRichContent(supports bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.supportsRich = supports
}
