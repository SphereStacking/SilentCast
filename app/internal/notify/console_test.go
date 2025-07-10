package notify

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"
)

func TestConsoleNotifier_Notify(t *testing.T) {
	notifier := NewConsoleNotifier()
	ctx := context.Background()
	
	tests := []struct {
		name         string
		notification Notification
		wantContains []string
	}{
		{
			name: "Info notification",
			notification: Notification{
				Title:   "Test Info",
				Message: "This is an info message",
				Level:   LevelInfo,
			},
			wantContains: []string{"INFO", "Test Info", "This is an info message"},
		},
		{
			name: "Warning notification",
			notification: Notification{
				Title:   "Test Warning",
				Message: "This is a warning",
				Level:   LevelWarning,
			},
			wantContains: []string{"WARN", "Test Warning", "This is a warning"},
		},
		{
			name: "Error notification",
			notification: Notification{
				Title:   "Test Error",
				Message: "This is an error",
				Level:   LevelError,
			},
			wantContains: []string{"ERROR", "Test Error", "This is an error"},
		},
		{
			name: "Success notification",
			notification: Notification{
				Title:   "Test Success",
				Message: "Operation completed",
				Level:   LevelSuccess,
			},
			wantContains: []string{"SUCCESS", "Test Success", "Operation completed"},
		},
		{
			name: "Title only",
			notification: Notification{
				Title: "Title Only",
				Level: LevelInfo,
			},
			wantContains: []string{"INFO", "Title Only"},
		},
		{
			name: "Message only",
			notification: Notification{
				Message: "Message Only",
				Level:   LevelInfo,
			},
			wantContains: []string{"INFO", "Message Only"},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stderr output
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w
			
			// Send notification
			err := notifier.Notify(ctx, tt.notification)
			if err != nil {
				t.Errorf("Notify() error = %v", err)
			}
			
			// Restore stderr and read output
			w.Close()
			os.Stderr = oldStderr
			
			var buf bytes.Buffer
			buf.ReadFrom(r)
			output := buf.String()
			
			// Check output contains expected strings
			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("Output missing '%s', got: %s", want, output)
				}
			}
		})
	}
}

func TestConsoleNotifier_IsAvailable(t *testing.T) {
	notifier := NewConsoleNotifier()
	
	if !notifier.IsAvailable() {
		t.Error("ConsoleNotifier should always be available")
	}
}

func TestManager_Notify(t *testing.T) {
	manager := NewManager()
	ctx := context.Background()
	
	// Test all convenience methods
	tests := []struct {
		name   string
		method func() error
	}{
		{
			name: "Info",
			method: func() error {
				return manager.Info(ctx, "Info Title", "Info message")
			},
		},
		{
			name: "Warning",
			method: func() error {
				return manager.Warning(ctx, "Warning Title", "Warning message")
			},
		},
		{
			name: "Error",
			method: func() error {
				return manager.Error(ctx, "Error Title", "Error message")
			},
		},
		{
			name: "Success",
			method: func() error {
				return manager.Success(ctx, "Success Title", "Success message")
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stderr output
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w
			
			// Call method
			err := tt.method()
			if err != nil {
				t.Errorf("%s() error = %v", tt.name, err)
			}
			
			// Restore stderr and read output
			w.Close()
			os.Stderr = oldStderr
			
			var buf bytes.Buffer
			buf.ReadFrom(r)
			output := buf.String()
			
			// Verify output is not empty
			if output == "" {
				t.Errorf("%s() produced no output", tt.name)
			}
		})
	}
}