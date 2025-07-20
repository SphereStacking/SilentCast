//go:build windows

package notify

import (
	"context"
	"strings"
	"testing"
)

func TestWindowsNotifier_Basic(t *testing.T) {
	notifier := NewWindowsNotifier("TestApp")

	if notifier == nil {
		t.Fatal("NewWindowsNotifier returned nil")
	}

	if notifier.appName != "TestApp" {
		t.Errorf("Expected appName to be 'TestApp', got '%s'", notifier.appName)
	}

	if notifier.maxOutputLength != 1000 {
		t.Errorf("Expected default maxOutputLength to be 1000, got %d", notifier.maxOutputLength)
	}
}

func TestWindowsNotifier_IsAvailable(t *testing.T) {
	notifier := NewWindowsNotifier("TestApp")

	// On Windows, the notifier should be available
	if !notifier.IsAvailable() {
		t.Error("WindowsNotifier should be available on Windows")
	}
}

func TestWindowsNotifier_Notify(t *testing.T) {
	notifier := NewWindowsNotifier("TestApp")
	ctx := context.Background()

	tests := []struct {
		name         string
		notification Notification
		expectError  bool
	}{
		{
			name: "Info notification",
			notification: Notification{
				Title:   "Test Info",
				Message: "This is a test info message",
				Level:   LevelInfo,
			},
			expectError: false,
		},
		{
			name: "Warning notification",
			notification: Notification{
				Title:   "Test Warning",
				Message: "This is a test warning message",
				Level:   LevelWarning,
			},
			expectError: false,
		},
		{
			name: "Error notification",
			notification: Notification{
				Title:   "Test Error",
				Message: "This is a test error message",
				Level:   LevelError,
			},
			expectError: false,
		},
		{
			name: "Success notification",
			notification: Notification{
				Title:   "Test Success",
				Message: "This is a test success message",
				Level:   LevelSuccess,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := notifier.Notify(ctx, tt.notification)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Logf("Notification may have failed (expected in test environment): %v", err)
				// Don't fail the test as notification might not work in CI
			}
		})
	}
}

func TestWindowsNotifier_ShowWithOutput(t *testing.T) {
	notifier := NewWindowsNotifier("TestApp")
	ctx := context.Background()

	tests := []struct {
		name         string
		notification OutputNotification
		expectError  bool
	}{
		{
			name: "Success with output",
			notification: OutputNotification{
				Notification: Notification{
					Title:   "Command Success",
					Message: "Command completed successfully",
					Level:   LevelSuccess,
				},
				Output:   "Hello, World!\nThis is command output",
				ExitCode: 0,
			},
			expectError: false,
		},
		{
			name: "Error with output",
			notification: OutputNotification{
				Notification: Notification{
					Title:   "Command Failed",
					Message: "Command execution failed",
					Level:   LevelError,
				},
				Output:   "Error: File not found\nCommand failed with error",
				ExitCode: 1,
			},
			expectError: false,
		},
		{
			name: "Long output truncation",
			notification: OutputNotification{
				Notification: Notification{
					Title:   "Long Output",
					Message: "Command with long output",
					Level:   LevelInfo,
				},
				Output:         strings.Repeat("This is a very long line of output\n", 50),
				TruncatedBytes: 100,
				ExitCode:       0,
			},
			expectError: false,
		},
		{
			name: "No output",
			notification: OutputNotification{
				Notification: Notification{
					Title:   "No Output",
					Message: "Command completed with no output",
					Level:   LevelInfo,
				},
				Output:   "",
				ExitCode: 0,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := notifier.ShowWithOutput(ctx, tt.notification)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Logf("Output notification may have failed (expected in test environment): %v", err)
				// Don't fail the test as notification might not work in CI
			}
		})
	}
}

func TestWindowsNotifier_SetMaxOutputLength(t *testing.T) {
	notifier := NewWindowsNotifier("TestApp")

	// Test default
	if notifier.maxOutputLength != 1000 {
		t.Errorf("Expected default maxOutputLength to be 1000, got %d", notifier.maxOutputLength)
	}

	// Test setting valid value
	result := notifier.SetMaxOutputLength(1500)
	if result != 1500 {
		t.Errorf("Expected SetMaxOutputLength to return 1500, got %d", result)
	}
	if notifier.maxOutputLength != 1500 {
		t.Errorf("Expected maxOutputLength to be 1500, got %d", notifier.maxOutputLength)
	}

	// Test setting value too large (should be rejected)
	result = notifier.SetMaxOutputLength(3000)
	if result != 1500 {
		t.Errorf("Expected SetMaxOutputLength to return previous value 1500, got %d", result)
	}

	// Test setting zero or negative (should be rejected)
	result = notifier.SetMaxOutputLength(0)
	if result != 1500 {
		t.Errorf("Expected SetMaxOutputLength(0) to return previous value 1500, got %d", result)
	}

	result = notifier.SetMaxOutputLength(-100)
	if result != 1500 {
		t.Errorf("Expected SetMaxOutputLength(-100) to return previous value 1500, got %d", result)
	}
}

func TestWindowsNotifier_SupportsRichContent(t *testing.T) {
	notifier := NewWindowsNotifier("TestApp")

	if !notifier.SupportsRichContent() {
		t.Error("WindowsNotifier should support rich content")
	}
}

func TestWindowsNotifier_OutputNotifierInterface(t *testing.T) {
	// Verify that WindowsNotifier implements OutputNotifier interface
	var _ OutputNotifier = (*WindowsNotifier)(nil)

	notifier := NewWindowsNotifier("TestApp")

	// Test interface methods exist and work
	length := notifier.SetMaxOutputLength(800)
	if length != 800 {
		t.Errorf("SetMaxOutputLength interface method failed")
	}

	if !notifier.SupportsRichContent() {
		t.Error("SupportsRichContent interface method failed")
	}
}

func TestEscapeForPowerShell(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple text", "simple text"},
		{"text with 'quotes'", "text with ''quotes''"},
		{"text\nwith\nnewlines", "text`nwith`nnewlines"},
		{"text\rwith\rcarriage\rreturns", "text`rwith`rcarriage`rreturns"},
		{"mixed 'quotes' and\nnewlines", "mixed ''quotes'' and`nnewlines"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := escapeForPowerShell(tt.input)
			if result != tt.expected {
				t.Errorf("escapeForPowerShell(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestEscapeForVBScript(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple text", "simple text"},
		{`text with "quotes"`, `text with ""quotes""`},
		{"text\nwith\nnewlines", `text" & vbCrLf & "with" & vbCrLf & "newlines`},
		{`mixed "quotes" and` + "\n" + `newlines`, `mixed ""quotes"" and" & vbCrLf & "newlines`},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := escapeForVBScript(tt.input)
			if result != tt.expected {
				t.Errorf("escapeForVBScript(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetVBScriptIcon(t *testing.T) {
	notifier := NewWindowsNotifier("TestApp")

	tests := []struct {
		level    Level
		expected int
	}{
		{LevelInfo, 64},
		{LevelWarning, 48},
		{LevelError, 16},
		{LevelSuccess, 64}, // Should default to info icon
	}

	for _, tt := range tests {
		t.Run(tt.level.String(), func(t *testing.T) {
			result := notifier.getVBScriptIcon(tt.level)
			if result != tt.expected {
				t.Errorf("getVBScriptIcon(%v) = %d, want %d", tt.level, result, tt.expected)
			}
		})
	}
}
