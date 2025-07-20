package notify

import (
	"context"
	"testing"
)

func TestMockNotifier_IsAvailable(t *testing.T) {
	tests := []struct {
		name      string
		available bool
	}{
		{
			name:      "available notifier",
			available: true,
		},
		{
			name:      "unavailable notifier",
			available: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockNotifier{
				available: tt.available,
			}
			
			got := m.IsAvailable()
			if got != tt.available {
				t.Errorf("IsAvailable() = %v, want %v", got, tt.available)
			}
		})
	}
}

func TestMockNotifier_Clear(t *testing.T) {
	// Create mock notifier and add some notifications
	m := &MockNotifier{
		available: true,
	}
	
	// Send some notifications
	ctx := context.Background()
	notifications := []Notification{
		{Title: "Test 1", Message: "Message 1", Level: LevelInfo},
		{Title: "Test 2", Message: "Message 2", Level: LevelWarning},
		{Title: "Test 3", Message: "Message 3", Level: LevelError},
	}
	
	for _, n := range notifications {
		err := m.Notify(ctx, n)
		if err != nil {
			t.Fatalf("Notify() error = %v", err)
		}
	}
	
	// Verify notifications were recorded
	recorded := m.GetNotifications()
	if len(recorded) != len(notifications) {
		t.Errorf("Expected %d notifications, got %d", len(notifications), len(recorded))
	}
	
	// Clear notifications
	m.Clear()
	
	// Verify notifications were cleared
	recorded = m.GetNotifications()
	if len(recorded) != 0 {
		t.Errorf("Clear() should remove all notifications, but %d remain", len(recorded))
	}
}

func TestMockNotifier_ConcurrentAccess(t *testing.T) {
	m := &MockNotifier{
		available: true,
	}
	
	ctx := context.Background()
	
	// Test concurrent Notify and GetNotifications
	done := make(chan bool)
	
	// Goroutine 1: Send notifications
	go func() {
		for i := 0; i < 100; i++ {
			n := Notification{
				Title:   "Test",
				Message: "Message",
				Level:   LevelInfo,
			}
			_ = m.Notify(ctx, n)
		}
		done <- true
	}()
	
	// Goroutine 2: Read notifications
	go func() {
		for i := 0; i < 100; i++ {
			_ = m.GetNotifications()
		}
		done <- true
	}()
	
	// Goroutine 3: Clear notifications
	go func() {
		for i := 0; i < 10; i++ {
			m.Clear()
		}
		done <- true
	}()
	
	// Wait for all goroutines
	for i := 0; i < 3; i++ {
		<-done
	}
	
	// If we reach here without deadlock or panic, the test passes
	t.Log("Concurrent access test passed")
}

