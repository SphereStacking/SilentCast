package notify

import (
	"context"
	"strings"
	"testing"
)

// mockNotifier is a test implementation of Notifier
type mockNotifier struct {
	name      string
	available bool
	notified  bool
	lastNotif Notification
}

func (m *mockNotifier) Notify(ctx context.Context, notification Notification) error {
	m.notified = true
	m.lastNotif = notification
	return nil
}

func (m *mockNotifier) IsAvailable() bool {
	return m.available
}

func TestManager_GetAvailableNotifiers(t *testing.T) {
	tests := []struct {
		name           string
		setupManager   func() *Manager
		wantContains   []string
		wantNotContain []string
	}{
		{
			name: "default manager with console",
			setupManager: func() *Manager {
				return NewManager()
			},
			wantContains: []string{"console"},
		},
		{
			name: "manager with custom notifiers",
			setupManager: func() *Manager {
				m := &Manager{notifiers: make([]Notifier, 0)}
				m.AddNotifier(&ConsoleNotifier{})
				m.AddNotifier(&mockNotifier{name: "mock1", available: true})
				m.AddNotifier(&mockNotifier{name: "mock2", available: true})
				return m
			},
			wantContains: []string{"console", "system"},
		},
		{
			name: "manager with unavailable notifiers",
			setupManager: func() *Manager {
				m := &Manager{notifiers: make([]Notifier, 0)}
				m.AddNotifier(&ConsoleNotifier{})
				m.AddNotifier(&mockNotifier{name: "mock1", available: false})
				return m
			},
			wantContains:   []string{"console"},
			wantNotContain: []string{"system"},
		},
		{
			name: "empty manager",
			setupManager: func() *Manager {
				return &Manager{notifiers: make([]Notifier, 0)}
			},
			wantContains: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.setupManager()
			got := m.GetAvailableNotifiers()

			// Check expected notifiers are present
			for _, want := range tt.wantContains {
				found := false
				for _, n := range got {
					if n == want {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("GetAvailableNotifiers() missing %q", want)
					t.Errorf("Got: %v", got)
				}
			}

			// Check unwanted notifiers are not present
			for _, notWant := range tt.wantNotContain {
				for _, n := range got {
					if n == notWant {
						t.Errorf("GetAvailableNotifiers() should not contain %q", notWant)
					}
				}
			}
		})
	}
}

func TestManager_AddNotifier(t *testing.T) {
	m := &Manager{notifiers: make([]Notifier, 0)}

	// Add available notifier
	n1 := &mockNotifier{name: "test1", available: true}
	m.AddNotifier(n1)

	if len(m.notifiers) != 1 {
		t.Errorf("AddNotifier() should add available notifier, got %d notifiers", len(m.notifiers))
	}

	// Add unavailable notifier
	n2 := &mockNotifier{name: "test2", available: false}
	m.AddNotifier(n2)

	if len(m.notifiers) != 1 {
		t.Errorf("AddNotifier() should not add unavailable notifier, got %d notifiers", len(m.notifiers))
	}

	// Add nil notifier
	m.AddNotifier(nil)

	if len(m.notifiers) != 1 {
		t.Errorf("AddNotifier() should not add nil notifier, got %d notifiers", len(m.notifiers))
	}
}

func TestManager_Notifications(t *testing.T) {
	ctx := context.Background()

	// Create manager with mock notifiers
	m := &Manager{notifiers: make([]Notifier, 0)}
	n1 := &mockNotifier{name: "test1", available: true}
	n2 := &mockNotifier{name: "test2", available: true}
	m.AddNotifier(n1)
	m.AddNotifier(n2)

	// Test Info notification
	err := m.Info(ctx, "Test Title", "Test Message")
	if err != nil {
		t.Errorf("Info() error = %v", err)
	}

	if !n1.notified || n1.lastNotif.Level != LevelInfo {
		t.Error("Info() should notify with LevelInfo")
	}

	// Test Warning notification
	n1.notified = false
	err = m.Warning(ctx, "Warning Title", "Warning Message")
	if err != nil {
		t.Errorf("Warning() error = %v", err)
	}

	if !n1.notified || n1.lastNotif.Level != LevelWarning {
		t.Error("Warning() should notify with LevelWarning")
	}

	// Test Error notification
	n1.notified = false
	err = m.Error(ctx, "Error Title", "Error Message")
	if err != nil {
		t.Errorf("Error() error = %v", err)
	}

	if !n1.notified || n1.lastNotif.Level != LevelError {
		t.Error("Error() should notify with LevelError")
	}

	// Test Success notification
	n1.notified = false
	err = m.Success(ctx, "Success Title", "Success Message")
	if err != nil {
		t.Errorf("Success() error = %v", err)
	}

	if !n1.notified || n1.lastNotif.Level != LevelSuccess {
		t.Error("Success() should notify with LevelSuccess")
	}
}

// TestInterfaceCompliance verifies that all notifiers implement the required interfaces
func TestInterfaceCompliance(t *testing.T) {
	// Only test the console notifier which is available on all platforms
	notifier := NewConsoleNotifier()

	// Test Notifier interface
	var _ Notifier = notifier

	// Test OutputNotifier interface
	var _ OutputNotifier = notifier

	// If we reach here, both interface checks passed
	t.Log("ConsoleNotifier implements both Notifier and OutputNotifier interfaces")
}

// TestOutputNotifierInterfaceMethods verifies OutputNotifier methods exist
func TestOutputNotifierInterfaceMethods(t *testing.T) {
	// Only test console notifier which is available on all platforms
	notifier := NewConsoleNotifier()

	// Test SetMaxOutputLength
	oldMax := notifier.SetMaxOutputLength(500)
	if oldMax < 0 {
		t.Errorf("SetMaxOutputLength returned negative value: %d", oldMax)
	}

	// Test SupportsRichContent
	_ = notifier.SupportsRichContent() // Should not panic
}

func TestManager_NotifyTimeout(t *testing.T) {
	ctx := context.Background()
	
	// Create manager with mock notifier
	m := NewManager()
	mockNotif := &mockNotifier{name: "test", available: true}
	m.AddNotifier(mockNotif)

	tests := []struct {
		name         string
		notification TimeoutNotification
	}{
		{
			name: "graceful timeout",
			notification: TimeoutNotification{
				ActionName:      "test script",
				TimeoutDuration: 30,
				ElapsedTime:     25,
				WasGraceful:     true,
				Output:          "Partial output here",
			},
		},
		{
			name: "force timeout",
			notification: TimeoutNotification{
				ActionName:      "long running script",
				TimeoutDuration: 60,
				ElapsedTime:     60,
				WasGraceful:     false,
				Output:          "Some output before timeout",
			},
		},
		{
			name: "timeout with long output",
			notification: TimeoutNotification{
				ActionName:      "verbose script",
				TimeoutDuration: 10,
				ElapsedTime:     10,
				WasGraceful:     true,
				Output:          strings.Repeat("Long output line\n", 50),
			},
		},
		{
			name: "timeout with no output",
			notification: TimeoutNotification{
				ActionName:      "silent script",
				TimeoutDuration: 5,
				ElapsedTime:     5,
				WasGraceful:     false,
				Output:          "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockNotif.notified = false
			err := m.NotifyTimeout(ctx, tt.notification)
			if err != nil {
				t.Errorf("NotifyTimeout failed: %v", err)
			}
			
			// Verify notification was sent
			if !mockNotif.notified {
				t.Error("NotifyTimeout should have sent notification")
			}
			
			// Verify notification level
			if mockNotif.lastNotif.Level != LevelWarning {
				t.Errorf("NotifyTimeout should use LevelWarning, got %v", mockNotif.lastNotif.Level)
			}
			
			// Verify title contains action name
			if !strings.Contains(mockNotif.lastNotif.Title, tt.notification.ActionName) {
				t.Errorf("Notification title should contain action name %q", tt.notification.ActionName)
			}
		})
	}
}

func TestManager_SupportsUpdateNotifications(t *testing.T) {
	tests := []struct {
		name     string
		notifier Notifier
		want     bool
	}{
		{
			name:     "console notifier",
			notifier: NewConsoleNotifier(),
			want:     true,
		},
		{
			name:     "mock notifier",
			notifier: &mockNotifier{available: true},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{notifiers: []Notifier{tt.notifier}}
			got := m.SupportsUpdateNotifications()
			if got != tt.want {
				t.Errorf("SupportsUpdateNotifications() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_GetUpdateNotifiers(t *testing.T) {
	// Create manager with different notifier types
	m := &Manager{notifiers: make([]Notifier, 0)}
	
	// Add console notifier (implements UpdateNotifier)
	console := NewConsoleNotifier()
	m.AddNotifier(console)
	
	// Add mock notifier (doesn't implement UpdateNotifier)
	mock := &mockNotifier{available: true}
	m.AddNotifier(mock)
	
	// Get update notifiers
	updateNotifiers := m.GetUpdateNotifiers()
	
	// Should only contain console notifier
	if len(updateNotifiers) != 1 {
		t.Errorf("GetUpdateNotifiers() returned %d notifiers, want 1", len(updateNotifiers))
	}
	
	// Verify it's the console notifier
	if _, ok := updateNotifiers[0].(*ConsoleNotifier); !ok {
		t.Error("GetUpdateNotifiers() should return ConsoleNotifier")
	}
}
