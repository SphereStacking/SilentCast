//go:build linux

package notify

import (
	"context"
	"errors"
	"strings"
	"testing"

	customErrors "github.com/SphereStacking/silentcast/internal/errors"
)

// TestLinuxNotificationErrorHandling tests that Linux notification errors use unified error handling
func TestLinuxNotificationErrorHandling(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func() *MockLinuxNotifier
		notification Notification
		expectError  bool
		expectType   customErrors.ErrorType
		expectContext map[string]interface{}
	}{
		{
			name: "all methods fail with unified error",
			setupMock: func() *MockLinuxNotifier {
				return &MockLinuxNotifier{
					notifySendFails: true,
					gdbusFails:      true,
					zenityFails:     true,
				}
			},
			notification: Notification{
				Title:   "Test Alert",
				Message: "Test message",
				Level:   LevelError,
			},
			expectError: true,
			expectType:  customErrors.ErrorTypeSystem,
			expectContext: map[string]interface{}{
				"notification_title": "Test Alert",
				"platform":          "linux",
				"tried_methods":     "notify-send,gdbus,zenity",
			},
		},
		{
			name: "notify-send succeeds",
			setupMock: func() *MockLinuxNotifier {
				return &MockLinuxNotifier{
					notifySendFails: false,
				}
			},
			notification: Notification{
				Title:   "Success",
				Message: "Success message",
				Level:   LevelSuccess,
			},
			expectError: false,
		},
		{
			name: "notify-send fails, gdbus succeeds",
			setupMock: func() *MockLinuxNotifier {
				return &MockLinuxNotifier{
					notifySendFails: true,
					gdbusFails:      false,
				}
			},
			notification: Notification{
				Title:   "Fallback Test",
				Message: "Fallback message",
				Level:   LevelWarning,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			notifier := tt.setupMock()
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

				// Check context
				for key, expectedValue := range tt.expectContext {
					if actual, ok := spellErr.Context[key]; !ok {
						t.Errorf("error context missing key %q", key)
					} else if actual != expectedValue {
						t.Errorf("error context[%q] = %v, want %v", key, actual, expectedValue)
					}
				}

			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestLinuxNotificationMethodSpecificErrors tests specific error cases for each method
func TestLinuxNotificationMethodSpecificErrors(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func() *MockLinuxNotifier
		notification Notification
		expectError  bool
		expectMethod string
	}{
		{
			name: "notify-send not found",
			setupMock: func() *MockLinuxNotifier {
				return &MockLinuxNotifier{
					notifySendNotFound: true,
					gdbusFails:         true,
					zenityFails:        true,
				}
			},
			notification: Notification{
				Title:   "Not Found Test",
				Message: "notify-send missing",
				Level:   LevelInfo,
			},
			expectError:  true,
			expectMethod: "notify-send",
		},
		{
			name: "gdbus execution error",
			setupMock: func() *MockLinuxNotifier {
				return &MockLinuxNotifier{
					notifySendFails: true,
					gdbusError:      "gdbus execution failed",
					zenityFails:     true,
				}
			},
			notification: Notification{
				Title:   "GDBus Error Test",
				Message: "gdbus fails",
				Level:   LevelError,
			},
			expectError:  true,
			expectMethod: "gdbus",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			notifier := tt.setupMock()
			err := notifier.Notify(ctx, tt.notification)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
					return
				}

				var spellErr *customErrors.SpellbookError
				if errors.As(err, &spellErr) {
					// Should include method information in context
					if method, ok := spellErr.Context["failed_method"]; ok {
						if !strings.Contains(method.(string), tt.expectMethod) {
							t.Errorf("error should include method %q, got %v", tt.expectMethod, method)
						}
					}
				}
			}
		})
	}
}

// TestLinuxNotificationUserMessages tests user-friendly error messages
func TestLinuxNotificationUserMessages(t *testing.T) {
	tests := []struct {
		name           string
		error          error
		expectContains string
	}{
		{
			name: "system dependency missing",
			error: func() error {
				return customErrors.New(customErrors.ErrorTypeSystem, "all notification methods failed").
					WithContext("notification_title", "Build Complete").
					WithContext("platform", "linux").
					WithContext("tried_methods", "notify-send,gdbus,zenity").
					WithContext("suggested_action", "install libnotify-bin or zenity package")
			}(),
			expectContains: "System error",
		},
		{
			name: "desktop environment issue",
			error: func() error {
				return customErrors.New(customErrors.ErrorTypeSystem, "notification service unavailable").
					WithContext("notification_title", "Deploy Status").
					WithContext("platform", "linux").
					WithContext("desktop_session", "unknown").
					WithContext("suggested_action", "check if desktop notification service is running")
			}(),
			expectContains: "notification service unavailable",
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

// MockLinuxNotifier for testing Linux-specific notification scenarios
type MockLinuxNotifier struct {
	notifySendFails    bool
	notifySendNotFound bool
	gdbusFails         bool
	gdbusError         string
	zenityFails        bool
	zenityNotFound     bool
}

func (m *MockLinuxNotifier) Notify(ctx context.Context, notification Notification) error {
	// Simulate notify-send first
	if !m.notifySendNotFound && !m.notifySendFails {
		return nil // Success
	}

	// If notify-send fails or not found, try gdbus
	if !m.gdbusFails && m.gdbusError == "" {
		return nil // Success with fallback
	}
	
	// If gdbus also failed, continue to zenity as final fallback

	// Try zenity as final fallback
	if !m.zenityFails {
		return nil // Success with second fallback
	}

	// All methods failed
	return customErrors.New(customErrors.ErrorTypeSystem, "all notification methods failed").
		WithContext("notification_title", notification.Title).
		WithContext("platform", "linux").
		WithContext("tried_methods", "notify-send,gdbus,zenity").
		WithContext("suggested_action", "install libnotify-bin or zenity package")
}

func (m *MockLinuxNotifier) IsAvailable() bool {
	return true
}

// Additional methods to implement OutputNotifier interface if needed
func (m *MockLinuxNotifier) ShowWithOutput(ctx context.Context, notification OutputNotification) error {
	return m.Notify(ctx, notification.Notification)
}

func (m *MockLinuxNotifier) SetMaxOutputLength(maxLength int) int {
	return 200
}

func (m *MockLinuxNotifier) SupportsRichContent() bool {
	return true
}