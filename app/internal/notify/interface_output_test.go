package notify

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestManager_NotifyWithOutput(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name              string
		notification      OutputNotification
		options           NotificationOptions
		setupNotifiers    func() []Notifier
		expectError       bool
		validateNotifiers func([]Notifier) error
	}{
		{
			name: "output notifier receives full notification",
			notification: OutputNotification{
				Notification: Notification{
					Title:   "Test Command",
					Message: "Command completed",
					Level:   LevelSuccess,
				},
				Output:   "Hello, World!",
				ExitCode: 0,
			},
			setupNotifiers: func() []Notifier {
				return []Notifier{
					NewMockOutputNotifier(true),
				}
			},
			validateNotifiers: func(notifiers []Notifier) error {
				mock := notifiers[0].(*MockOutputNotifier)
				notifications := mock.GetOutputNotifications()
				if len(notifications) != 1 {
					return errors.New("expected 1 output notification")
				}
				if notifications[0].Output != "Hello, World!" {
					return errors.New("output not preserved")
				}
				return nil
			},
		},
		{
			name: "fallback to regular notifier",
			notification: OutputNotification{
				Notification: Notification{
					Title:   "Test Command",
					Message: "Command completed",
					Level:   LevelInfo,
				},
				Output:   "Some output",
				ExitCode: 0,
			},
			setupNotifiers: func() []Notifier {
				return []Notifier{
					NewMockNotifier(true),
				}
			},
			validateNotifiers: func(notifiers []Notifier) error {
				mock := notifiers[0].(*MockNotifier)
				notifications := mock.GetNotifications()
				if len(notifications) != 1 {
					return errors.New("expected 1 notification")
				}
				if !strings.Contains(notifications[0].Message, "Some output") {
					return errors.New("output not included in message")
				}
				return nil
			},
		},
		{
			name: "output truncation",
			notification: OutputNotification{
				Notification: Notification{
					Title:   "Long Output",
					Message: "Command with truncated output",
					Level:   LevelWarning,
				},
				Output:   strings.Repeat("A", 100),
				ExitCode: 0,
			},
			options: NotificationOptions{
				MaxOutputLength: 50,
			},
			setupNotifiers: func() []Notifier {
				return []Notifier{
					NewMockOutputNotifier(true),
				}
			},
			validateNotifiers: func(notifiers []Notifier) error {
				mock := notifiers[0].(*MockOutputNotifier)
				notifications := mock.GetOutputNotifications()
				if len(notifications) != 1 {
					return errors.New("expected 1 output notification")
				}
				if len(notifications[0].Output) != 50 {
					return errors.New("output not truncated to 50 chars")
				}
				if notifications[0].TruncatedBytes != 50 {
					return errors.New("truncated bytes not set correctly")
				}
				return nil
			},
		},
		{
			name: "multiple notifiers with mixed capabilities",
			notification: OutputNotification{
				Notification: Notification{
					Title:   "Mixed Test",
					Message: "Testing mixed notifiers",
					Level:   LevelInfo,
				},
				Output:   "Test output",
				ExitCode: 0,
			},
			setupNotifiers: func() []Notifier {
				return []Notifier{
					NewMockNotifier(true),
					NewMockOutputNotifier(true),
				}
			},
			validateNotifiers: func(notifiers []Notifier) error {
				// Check regular notifier
				mock1 := notifiers[0].(*MockNotifier)
				if len(mock1.GetNotifications()) != 1 {
					return errors.New("regular notifier should have received notification")
				}

				// Check output notifier
				mock2 := notifiers[1].(*MockOutputNotifier)
				if len(mock2.GetOutputNotifications()) != 1 {
					return errors.New("output notifier should have received output notification")
				}
				return nil
			},
		},
		{
			name: "error handling",
			notification: OutputNotification{
				Notification: Notification{
					Title:   "Error Test",
					Message: "This should fail",
					Level:   LevelError,
				},
			},
			setupNotifiers: func() []Notifier {
				mock := NewMockOutputNotifier(true)
				mock.SetError(errors.New("simulated error"))
				return []Notifier{mock}
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &Manager{
				notifiers: tt.setupNotifiers(),
				options:   tt.options,
			}

			if manager.options.MaxOutputLength == 0 {
				manager.options.MaxOutputLength = 1024 // Default
			}

			err := manager.NotifyWithOutput(ctx, tt.notification)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if tt.validateNotifiers != nil {
				if err := tt.validateNotifiers(manager.notifiers); err != nil {
					t.Error(err)
				}
			}
		})
	}
}

func TestManager_SupportsOutputNotifications(t *testing.T) {
	tests := []struct {
		name           string
		setupNotifiers func() []Notifier
		want           bool
	}{
		{
			name: "no output notifiers",
			setupNotifiers: func() []Notifier {
				return []Notifier{
					NewMockNotifier(true),
					NewMockNotifier(true),
				}
			},
			want: false,
		},
		{
			name: "has output notifier",
			setupNotifiers: func() []Notifier {
				return []Notifier{
					NewMockNotifier(true),
					NewMockOutputNotifier(true),
				}
			},
			want: true,
		},
		{
			name: "all output notifiers",
			setupNotifiers: func() []Notifier {
				return []Notifier{
					NewMockOutputNotifier(true),
					NewMockOutputNotifier(true),
				}
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &Manager{
				notifiers: tt.setupNotifiers(),
			}

			got := manager.SupportsOutputNotifications()
			if got != tt.want {
				t.Errorf("SupportsOutputNotifications() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_GetOutputNotifiers(t *testing.T) {
	mock1 := NewMockNotifier(true)
	mock2 := NewMockOutputNotifier(true)
	mock3 := NewMockOutputNotifier(true)

	manager := &Manager{
		notifiers: []Notifier{mock1, mock2, mock3},
	}

	outputNotifiers := manager.GetOutputNotifiers()

	if len(outputNotifiers) != 2 {
		t.Errorf("Expected 2 output notifiers, got %d", len(outputNotifiers))
	}

	// Verify they are the correct instances
	for i, on := range outputNotifiers {
		if on == nil {
			t.Errorf("Output notifier %d is nil", i)
		}
	}
}

func TestManager_SetGetOptions(t *testing.T) {
	manager := NewManager()

	// Test default options
	defaultOpts := manager.GetOptions()
	if defaultOpts.MaxOutputLength != 1024 {
		t.Errorf("Expected default MaxOutputLength to be 1024, got %d", defaultOpts.MaxOutputLength)
	}

	// Test setting new options
	newOpts := NotificationOptions{
		MaxOutputLength: 2048,
		FormatAsCode:    false,
		Priority:        "high",
		Sound:           false,
	}

	manager.SetOptions(newOpts)

	gotOpts := manager.GetOptions()
	if gotOpts.MaxOutputLength != 2048 {
		t.Errorf("Expected MaxOutputLength to be 2048, got %d", gotOpts.MaxOutputLength)
	}
	if gotOpts.FormatAsCode {
		t.Error("Expected FormatAsCode to be false")
	}
	if gotOpts.Priority != "high" {
		t.Errorf("Expected Priority to be 'high', got %s", gotOpts.Priority)
	}
	if gotOpts.Sound {
		t.Error("Expected Sound to be false")
	}
}

func TestOutputNotifier_Interface(t *testing.T) {
	// Verify MockOutputNotifier implements OutputNotifier
	var _ OutputNotifier = (*MockOutputNotifier)(nil)

	mock := NewMockOutputNotifier(true)

	// Test SetMaxOutputLength
	result := mock.SetMaxOutputLength(512)
	if result != 512 {
		t.Errorf("Expected SetMaxOutputLength to return 512, got %d", result)
	}

	// Test SupportsRichContent
	if !mock.SupportsRichContent() {
		t.Error("Expected SupportsRichContent to return true by default")
	}

	mock.SetSupportsRichContent(false)
	if mock.SupportsRichContent() {
		t.Error("Expected SupportsRichContent to return false after setting")
	}
}
