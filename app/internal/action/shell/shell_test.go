package shell

import (
	"context"
	"os/exec"
	"strings"
	"testing"
)

// MockShellExecutor for testing
type mockShellExecutor struct {
	keepOpenCalled bool
	keepOpenValue  bool
}

func (m *mockShellExecutor) GetShell() (string, string) {
	return "sh", "-c"
}

func (m *mockShellExecutor) WrapInTerminal(ctx context.Context, cmd *exec.Cmd) *exec.Cmd {
	return m.WrapInTerminalWithOptions(ctx, cmd, true)
}

func (m *mockShellExecutor) WrapInTerminalWithOptions(ctx context.Context, cmd *exec.Cmd, keepOpen bool) *exec.Cmd {
	m.keepOpenCalled = true
	m.keepOpenValue = keepOpen
	// Return a mock command that we can verify
	return exec.CommandContext(ctx, "echo", "terminal-wrapped")
}

func (m *mockShellExecutor) IsInteractiveCommand(command string) bool {
	return false
}

func TestShellExecutor_KeepOpen(t *testing.T) {
	// Save original shell factory
	originalFactory := shellFactory
	defer func() { shellFactory = originalFactory }()

	// Create mock shell executor
	mock := &mockShellExecutor{}
	shellFactory = func() ShellExecutor { return mock }

	tests := []struct {
		name           string
		keepOpen       bool
		expectWrapped  bool
		expectKeepOpen bool
	}{
		{
			name:           "keep_open true",
			keepOpen:       true,
			expectWrapped:  true,
			expectKeepOpen: true,
		},
		{
			name:           "keep_open false",
			keepOpen:       false,
			expectWrapped:  false,
			expectKeepOpen: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock state
			mock.keepOpenCalled = false
			mock.keepOpenValue = false

			// Run a simple command with keep_open setting
			ctx := context.Background()
			cmd := exec.CommandContext(ctx, "echo", "test")

			executor := GetShellExecutor()
			if tt.keepOpen {
				wrappedCmd := executor.WrapInTerminalWithOptions(ctx, cmd, tt.keepOpen)
				// Verify command was wrapped
				if !strings.Contains(wrappedCmd.String(), "terminal-wrapped") {
					t.Error("Command was not wrapped for terminal execution")
				}
			}

			// Verify expectations
			if tt.expectWrapped && !mock.keepOpenCalled {
				t.Error("Expected WrapInTerminalWithOptions to be called")
			}
			if tt.expectWrapped && mock.keepOpenValue != tt.expectKeepOpen {
				t.Errorf("Expected keepOpen=%v, got %v", tt.expectKeepOpen, mock.keepOpenValue)
			}
		})
	}
}

// Test platform-specific implementations
func TestPlatformShellExecutor_Interface(t *testing.T) {
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, "echo", "test")

	// Get the current platform's shell executor
	executor := GetShellExecutor()

	// Test WrapInTerminalWithOptions exists and doesn't panic
	t.Run("WrapInTerminalWithOptions_keepOpen_true", func(t *testing.T) {
		wrappedCmd := executor.WrapInTerminalWithOptions(ctx, cmd, true)
		if wrappedCmd == nil {
			t.Error("WrapInTerminalWithOptions returned nil")
		}
	})

	t.Run("WrapInTerminalWithOptions_keepOpen_false", func(t *testing.T) {
		wrappedCmd := executor.WrapInTerminalWithOptions(ctx, cmd, false)
		if wrappedCmd == nil {
			t.Error("WrapInTerminalWithOptions with keepOpen=false returned nil")
		}
	})

	t.Run("WrapInTerminal_compatibility", func(t *testing.T) {
		// Test that WrapInTerminal still works
		wrappedCmd := executor.WrapInTerminal(ctx, cmd)
		if wrappedCmd == nil {
			t.Error("WrapInTerminal returned nil")
		}
	})
}
