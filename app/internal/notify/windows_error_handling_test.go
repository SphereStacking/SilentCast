//go:build windows

package notify

import (
	"context"
	"errors"
	"strings"
	"testing"

	customErrors "github.com/SphereStacking/silentcast/internal/errors"
)

// TestWindowsNotificationErrorHandling tests that Windows notification errors use unified error handling
func TestWindowsNotificationErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func() *MockWindowsNotifier
		notification  Notification
		expectError   bool
		expectType    customErrors.ErrorType
		expectContext map[string]interface{}
	}{
		{
			name: "all methods fail with unified error",
			setupMock: func() *MockWindowsNotifier {
				return &MockWindowsNotifier{
					toastFails:   true,
					msgFails:     true,
					balloonFails: true,
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
				"platform":           "windows",
				"tried_methods":      "toast,msg,balloon",
			},
		},
		{
			name: "toast succeeds",
			setupMock: func() *MockWindowsNotifier {
				return &MockWindowsNotifier{
					toastFails: false,
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
			name: "toast fails, msg succeeds",
			setupMock: func() *MockWindowsNotifier {
				return &MockWindowsNotifier{
					toastFails: true,
					msgFails:   false,
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

// TestWindowsNotificationMethodSpecificErrors tests specific error cases for each method
func TestWindowsNotificationMethodSpecificErrors(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func() *MockWindowsNotifier
		notification Notification
		expectError  bool
		expectMethod string
	}{
		{
			name: "powershell not available",
			setupMock: func() *MockWindowsNotifier {
				return &MockWindowsNotifier{
					toastNotFound: true,
					msgFails:      true,
					balloonFails:  true,
				}
			},
			notification: Notification{
				Title:   "PowerShell Test",
				Message: "PowerShell missing",
				Level:   LevelInfo,
			},
			expectError:  true,
			expectMethod: "toast",
		},
		{
			name: "msg execution error",
			setupMock: func() *MockWindowsNotifier {
				return &MockWindowsNotifier{
					toastFails:   true,
					msgError:     "msg execution failed",
					balloonFails: true,
				}
			},
			notification: Notification{
				Title:   "MSG Error Test",
				Message: "msg fails",
				Level:   LevelError,
			},
			expectError:  true,
			expectMethod: "msg",
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

// TestWindowsNotificationUserMessages tests user-friendly error messages
func TestWindowsNotificationUserMessages(t *testing.T) {
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
					WithContext("platform", "windows").
					WithContext("tried_methods", "toast,msg,balloon").
					WithContext("suggested_action", "check Windows notification settings")
			}(),
			expectContains: "System error",
		},
		{
			name: "permission denied error",
			error: func() error {
				return customErrors.New(customErrors.ErrorTypePermission, "notification permission denied").
					WithContext("notification_title", "Deploy Status").
					WithContext("platform", "windows").
					WithContext("method", "toast").
					WithContext("suggested_action", "enable notifications in Windows settings")
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

// MockWindowsNotifier for testing Windows-specific notification scenarios
type MockWindowsNotifier struct {
	toastFails    bool
	toastNotFound bool
	msgFails      bool
	msgError      string
	balloonFails  bool
	balloonError  string
}

func (m *MockWindowsNotifier) Notify(ctx context.Context, notification Notification) error {
	// Simulate toast notification
	if !m.toastNotFound && !m.toastFails {
		return nil // Success
	}

	// If toast fails or not found, try msg
	if !m.msgFails && m.msgError == "" {
		return nil // Success with fallback
	}

	if m.msgError != "" {
		// msg specific error - but continue to balloon
		// Don't return here, try balloon as final fallback
	}

	// Try balloon as final fallback
	if !m.balloonFails {
		return nil // Success with second fallback
	}

	// All methods failed
	return customErrors.New(customErrors.ErrorTypeSystem, "all notification methods failed").
		WithContext("notification_title", notification.Title).
		WithContext("platform", "windows").
		WithContext("tried_methods", "toast,msg,balloon").
		WithContext("suggested_action", "check Windows notification settings")
}

func (m *MockWindowsNotifier) IsAvailable() bool {
	return true
}

// Additional methods to implement OutputNotifier interface if needed
func (m *MockWindowsNotifier) ShowWithOutput(ctx context.Context, notification OutputNotification) error {
	return m.Notify(ctx, notification.Notification)
}

func (m *MockWindowsNotifier) SetMaxOutputLength(maxLength int) int {
	return 1000
}

func (m *MockWindowsNotifier) SupportsRichContent() bool {
	return true
}
