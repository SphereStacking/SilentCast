package notify

import (
	"context"
	"errors"
	"strings"
	"testing"

	customErrors "github.com/SphereStacking/silentcast/internal/errors"
)

// TestNotificationErrorHandling tests that notification errors use the unified error handling pattern
func TestNotificationErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		setupNotifier func() Notifier
		notification  Notification
		expectError   bool
		expectType    customErrors.ErrorType
	}{
		{
			name: "mock notifier with error",
			setupNotifier: func() Notifier {
				mock := NewMockNotifier(true)
				mock.SetError(customErrors.New(customErrors.ErrorTypeSystem, "mock notification failed"))
				return mock
			},
			notification: Notification{
				Title:   "Test",
				Message: "Test message",
				Level:   LevelInfo,
			},
			expectError: true,
			expectType:  customErrors.ErrorTypeSystem,
		},
		{
			name: "mock notifier success",
			setupNotifier: func() Notifier {
				return NewMockNotifier(true)
			},
			notification: Notification{
				Title:   "Test",
				Message: "Test message",
				Level:   LevelInfo,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			notifier := tt.setupNotifier()
			err := notifier.Notify(ctx, tt.notification)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
					return
				}

				// Check if error uses unified error pattern
				var spellErr *customErrors.SpellbookError
				if !errors.As(err, &spellErr) {
					t.Errorf("error should be SpellbookError, got %T", err)
					return
				}

				// Check error type
				if spellErr.Type != tt.expectType {
					t.Errorf("expected error type %v, got %v", tt.expectType, spellErr.Type)
				}

				// Check context is included
				if spellErr.Context == nil || len(spellErr.Context) == 0 {
					t.Error("error should include context information")
				}

				// Check specific context fields
				if title, ok := spellErr.Context["notification_title"]; !ok || title != tt.notification.Title {
					t.Errorf("error context should include notification_title=%v", tt.notification.Title)
				}

			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestNotificationRetryWithContext tests retry logic with contextual errors
func TestNotificationRetryWithContext(t *testing.T) {
	ctx := context.Background()

	// Create a notifier that fails multiple times then succeeds
	mockNotifier := NewRetryMockNotifier(3)

	notification := Notification{
		Title:   "Retry Test",
		Message: "Testing retry logic",
		Level:   LevelInfo,
	}

	// Use an enhanced notifier with retry logic
	enhancedNotifier := NewNotifierWithRetry(mockNotifier, 3)
	err := enhancedNotifier.Notify(ctx, notification)

	if err != nil {
		var spellErr *customErrors.SpellbookError
		if errors.As(err, &spellErr) {
			// Should have retry context
			if attempts, ok := spellErr.Context["retry_attempts"]; !ok {
				t.Error("error should include retry_attempts in context")
			} else if attempts != 3 {
				t.Errorf("expected 3 retry attempts, got %v", attempts)
			}
		} else {
			t.Errorf("retry error should be SpellbookError, got %T", err)
		}
	}
}

// TestNotificationErrorUserMessages tests user-friendly error messages
func TestNotificationErrorUserMessages(t *testing.T) {
	tests := []struct {
		name           string
		error          error
		expectContains string
	}{
		{
			name: "notification failure with context",
			error: func() error {
				return customErrors.New(customErrors.ErrorTypeSystem, "notification failed").
					WithContext("notification_title", "Important Alert").
					WithContext("platform", "darwin").
					WithContext("method", "osascript")
			}(),
			expectContains: "notification failed",
		},
		{
			name: "permission denied error",
			error: func() error {
				return customErrors.New(customErrors.ErrorTypePermission, "notification permission denied").
					WithContext("notification_title", "System Update").
					WithContext("suggested_action", "enable notifications in system preferences")
			}(),
			expectContains: "Permission error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userMsg := customErrors.GetUserMessage(tt.error)
			if !strings.Contains(userMsg, tt.expectContains) {
				t.Errorf("user message %q should contain %q", userMsg, tt.expectContains)
			}
		})
	}
}

// TestNotificationErrorLogging tests structured logging integration
func TestNotificationErrorLogging(t *testing.T) {
	err := func() error {
		return customErrors.New(customErrors.ErrorTypeSystem, "notification delivery failed").
			WithContext("notification_title", "Build Failed").
			WithContext("platform", "darwin").
			WithContext("fallback_tried", "terminal-notifier").
			WithContext("exit_code", 1)
	}()

	var spellErr *customErrors.SpellbookError
	if !errors.As(err, &spellErr) {
		t.Fatal("error should be SpellbookError")
	}

	fields := spellErr.LogFields()

	expectedFields := map[string]interface{}{
		"error_type":         "system",
		"message":            "notification delivery failed",
		"notification_title": "Build Failed",
		"platform":           "darwin",
		"fallback_tried":     "terminal-notifier",
		"exit_code":          1,
	}

	for key, expectedValue := range expectedFields {
		if actual, ok := fields[key]; !ok {
			t.Errorf("LogFields() missing key %q", key)
		} else if actual != expectedValue {
			t.Errorf("LogFields()[%q] = %v, want %v", key, actual, expectedValue)
		}
	}
}

// RetryMockNotifier simulates retry scenarios
type RetryMockNotifier struct {
	maxRetries int
	attempts   int
}

func NewRetryMockNotifier(maxRetries int) *RetryMockNotifier {
	return &RetryMockNotifier{
		maxRetries: maxRetries,
		attempts:   0,
	}
}

func (r *RetryMockNotifier) Notify(ctx context.Context, notification Notification) error {
	r.attempts++
	if r.attempts <= r.maxRetries {
		return customErrors.New(customErrors.ErrorTypeSystem, "temporary notification failure").
			WithContext("notification_title", notification.Title).
			WithContext("attempt", r.attempts)
	}
	return nil // Success after retries
}

func (r *RetryMockNotifier) IsAvailable() bool {
	return true
}

// NotifierWithRetry adds retry logic with unified error handling
type NotifierWithRetry struct {
	notifier   Notifier
	maxRetries int
}

func NewNotifierWithRetry(notifier Notifier, maxRetries int) *NotifierWithRetry {
	return &NotifierWithRetry{
		notifier:   notifier,
		maxRetries: maxRetries,
	}
}

func (n *NotifierWithRetry) Notify(ctx context.Context, notification Notification) error {
	var lastErr error

	for attempt := 1; attempt <= n.maxRetries; attempt++ {
		err := n.notifier.Notify(ctx, notification)
		if err == nil {
			return nil
		}

		lastErr = err

		// Add retry context to the error
		var spellErr *customErrors.SpellbookError
		if errors.As(err, &spellErr) {
			spellErr = spellErr.WithContext("retry_attempt", attempt)
		}
	}

	// Wrap final error with retry information
	var spellErr *customErrors.SpellbookError
	if errors.As(lastErr, &spellErr) {
		return spellErr.WithContext("retry_attempts", n.maxRetries).
			WithContext("max_retries_exceeded", true)
	}

	return customErrors.Wrap(customErrors.ErrorTypeSystem, "notification failed after retries", lastErr).
		WithContext("retry_attempts", n.maxRetries)
}

func (n *NotifierWithRetry) IsAvailable() bool {
	return n.notifier.IsAvailable()
}
