package notify

import (
	"bytes"
	"context"
	"io"
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
			if _, err := buf.ReadFrom(r); err != nil {
				t.Errorf("Failed to read output: %v", err)
			}
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

func TestConsoleNotifier_OutputNotifier(t *testing.T) {
	// Redirect stderr to capture output
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	notifier := NewConsoleNotifier()
	ctx := context.Background()

	// Test ShowWithOutput
	notification := OutputNotification{
		Notification: Notification{
			Title:   "Test Command",
			Message: "Command completed",
			Level:   LevelSuccess,
		},
		Output:   "Hello\nWorld",
		ExitCode: 0,
	}

	err := notifier.ShowWithOutput(ctx, notification)
	if err != nil {
		t.Errorf("ShowWithOutput failed: %v", err)
	}

	// Test with truncated output
	notification.TruncatedBytes = 100
	notification.ExitCode = 1

	err = notifier.ShowWithOutput(ctx, notification)
	if err != nil {
		t.Errorf("ShowWithOutput with truncation failed: %v", err)
	}

	// Restore stderr
	w.Close()
	os.Stderr = oldStderr

	// Read captured output
	out, _ := io.ReadAll(r)
	output := string(out)

	// Verify output contains expected elements
	if !strings.Contains(output, "Test Command") {
		t.Error("Output should contain notification title")
	}
	if !strings.Contains(output, "Output:") {
		t.Error("Output should contain Output header")
	}
	if !strings.Contains(output, "Hello") {
		t.Error("Output should contain command output")
	}
	if !strings.Contains(output, "bytes truncated") {
		t.Error("Output should contain truncation info")
	}
	if !strings.Contains(output, "Exit code: 1") {
		t.Error("Output should contain exit code")
	}
}

func TestConsoleNotifier_SetMaxOutputLength(t *testing.T) {
	notifier := NewConsoleNotifier()

	// Test default
	if notifier.maxOutputLength != 2048 {
		t.Errorf("Expected default max output length to be 2048, got %d", notifier.maxOutputLength)
	}

	// Test setting new value
	result := notifier.SetMaxOutputLength(1024)
	if result != 1024 {
		t.Errorf("Expected SetMaxOutputLength to return 1024, got %d", result)
	}
	if notifier.maxOutputLength != 1024 {
		t.Errorf("Expected max output length to be 1024, got %d", notifier.maxOutputLength)
	}

	// Test with zero (should not change)
	result = notifier.SetMaxOutputLength(0)
	if result != 1024 {
		t.Errorf("Expected SetMaxOutputLength(0) to return current value 1024, got %d", result)
	}
}

func TestConsoleNotifier_SupportsRichContent(_ *testing.T) {
	notifier := NewConsoleNotifier()

	// This will depend on whether tests are run in a terminal
	// Just verify the method exists and returns a boolean
	_ = notifier.SupportsRichContent()
}

func TestConsoleNotifier_ShowUpdateNotification(t *testing.T) {
	notifier := NewConsoleNotifier()
	ctx := context.Background()

	// Capture stderr output
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	notification := UpdateNotification{
		Notification: Notification{
			Title:   "Update Available",
			Message: "New version is available",
			Level:   LevelInfo,
		},
		CurrentVersion: "1.0.0",
		NewVersion:     "1.1.0",
		PublishedAt:    "2024-01-01",
		DownloadSize:   1024 * 1024, // 1MB
		ReleaseNotes:   "Bug fixes and improvements",
		Actions:        []string{"update", "skip"},
	}

	err := notifier.ShowUpdateNotification(ctx, notification)
	if err != nil {
		t.Errorf("ShowUpdateNotification() error = %v", err)
	}

	// Restore stderr and read output
	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Check expected content
	expectedStrings := []string{
		"Update Available",
		"Current: 1.0.0",
		"Latest:  1.1.0",
		"Published: 2024-01-01",
		"Size:",
		"Release Notes:",
		"Bug fixes",
		"Available actions:",
		"./silentcast --self-update",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Output missing '%s', got: %s", expected, output)
		}
	}
}

func TestConsoleNotifier_ShowUpdateNotification_LongReleaseNotes(t *testing.T) {
	notifier := NewConsoleNotifier()
	ctx := context.Background()

	// Capture stderr output
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	longNotes := strings.Repeat("This is a very long release note. ", 20) // > 300 chars
	notification := UpdateNotification{
		Notification: Notification{
			Title: "Update Available",
			Level: LevelInfo,
		},
		CurrentVersion: "1.0.0",
		NewVersion:     "1.1.0",
		ReleaseNotes:   longNotes,
	}

	err := notifier.ShowUpdateNotification(ctx, notification)
	if err != nil {
		t.Errorf("ShowUpdateNotification() error = %v", err)
	}

	// Restore stderr and read output
	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Check that release notes are truncated
	if !strings.Contains(output, "...") {
		t.Error("Long release notes should be truncated with '...'")
	}
}

func TestConsoleNotifier_SupportsUpdateActions(t *testing.T) {
	notifier := NewConsoleNotifier()

	if notifier.SupportsUpdateActions() {
		t.Error("ConsoleNotifier should not support interactive update actions")
	}
}

func TestConsoleNotifier_OnUpdateAction(t *testing.T) {
	notifier := NewConsoleNotifier()

	updateInfo := UpdateNotification{}
	err := notifier.OnUpdateAction(UpdateActionUpdate, updateInfo)

	if err == nil {
		t.Error("OnUpdateAction should return error for console notifier")
	}

	if !strings.Contains(err.Error(), "does not support interactive actions") {
		t.Errorf("Error should mention interactive actions, got: %v", err)
	}
}

func TestIndentOutput(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "single line",
			input:  "Hello",
			output: "        Hello",
		},
		{
			name:   "multiple lines",
			input:  "Line 1\nLine 2",
			output: "        Line 1\n        Line 2",
		},
		{
			name:   "empty lines",
			input:  "Line 1\n\nLine 3",
			output: "        Line 1\n\n        Line 3",
		},
		{
			name:   "empty string",
			input:  "",
			output: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := indentOutput(tt.input)
			if result != tt.output {
				t.Errorf("indentOutput() = %q, want %q", result, tt.output)
			}
		})
	}
}

func TestConsoleNotifier_NotifyUnknownLevel(t *testing.T) {
	notifier := NewConsoleNotifier()
	ctx := context.Background()

	// Capture stderr output
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	notification := Notification{
		Title:   "Unknown Level",
		Message: "Test message",
		Level:   Level(999), // Unknown level
	}

	err := notifier.Notify(ctx, notification)
	if err != nil {
		t.Errorf("Notify() error = %v", err)
	}

	// Restore stderr and read output
	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Should use default NOTIFY prefix
	if !strings.Contains(output, "NOTIFY") || !strings.Contains(output, "Unknown Level") {
		t.Errorf("Output should contain NOTIFY and title, got: %s", output)
	}
}

func TestConsoleNotifier_FormatUpdateSize(t *testing.T) {
	// This tests the formatUpdateSize function indirectly via ShowUpdateNotification
	notifier := NewConsoleNotifier()
	ctx := context.Background()

	tests := []struct {
		name         string
		downloadSize int64
		wantContains string
	}{
		{
			name:         "bytes",
			downloadSize: 512,
			wantContains: "512 B",
		},
		{
			name:         "kilobytes",
			downloadSize: 1536, // 1.5 KB
			wantContains: "1.5 KB",
		},
		{
			name:         "megabytes",
			downloadSize: 2 * 1024 * 1024, // 2 MB
			wantContains: "2.0 MB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stderr output
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			notification := UpdateNotification{
				Notification: Notification{
					Title: "Test",
					Level: LevelInfo,
				},
				CurrentVersion: "1.0.0",
				NewVersion:     "1.1.0",
				DownloadSize:   tt.downloadSize,
			}

			_ = notifier.ShowUpdateNotification(ctx, notification)

			// Restore stderr and read output
			w.Close()
			os.Stderr = oldStderr

			var buf bytes.Buffer
			if _, err := buf.ReadFrom(r); err != nil {
				t.Errorf("Failed to read output: %v", err)
			}
			output := buf.String()

			if !strings.Contains(output, tt.wantContains) {
				t.Errorf("Output should contain '%s', got: %s", tt.wantContains, output)
			}
		})
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
			_ = tt.method()
			// In test environment, system notifier may fail, but console should work
			// so we don't treat this as a fatal error

			// Restore stderr and read output
			w.Close()
			os.Stderr = oldStderr

			var buf bytes.Buffer
			if _, err := buf.ReadFrom(r); err != nil {
				t.Errorf("Failed to read output: %v", err)
			}
			output := buf.String()

			// Verify output is not empty
			if output == "" {
				t.Errorf("%s() produced no output", tt.name)
			}
		})
	}
}
