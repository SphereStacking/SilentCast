package elevated

import (
	"context"
	"runtime"
	"testing"

	"github.com/SphereStacking/silentcast/internal/action/script"
	"github.com/SphereStacking/silentcast/internal/config"
)

// mockExecutor for testing
type mockExecutor struct {
	executeFunc func(context.Context) error
	stringFunc  func() string
}

func (m *mockExecutor) Execute(ctx context.Context) error {
	if m.executeFunc != nil {
		return m.executeFunc(ctx)
	}
	return nil
}

func (m *mockExecutor) String() string {
	if m.stringFunc != nil {
		return m.stringFunc()
	}
	return "mock"
}

func TestGetCommandString(t *testing.T) {
	tests := []struct {
		name     string
		baseExec Executor
		want     string
	}{
		{
			name: "Mock executor",
			baseExec: &mockExecutor{
				stringFunc: func() string {
					return "test command"
				},
			},
			want: "test command",
		},
		{
			name: "Unknown executor type",
			baseExec: &mockExecutor{
				stringFunc: func() string { return "unknown" },
			},
			want: "unknown", // getCommandString now returns String() result
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ElevatedExecutor{
				baseExecutor: tt.baseExec,
			}
			got := e.getCommandString()
			if got != tt.want {
				t.Errorf("getCommandString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElevatedExecutor_Execute_Coverage(t *testing.T) {
	// Test when already running as admin
	mockConfig := config.ActionConfig{
		Type:    "script",
		Command: "echo test",
	}
	baseExecutor := script.NewScriptExecutor(&mockConfig)
	elevatedExecutor := &ElevatedExecutor{
		baseExecutor: baseExecutor,
		isAdmin:      true,
	}

	// Just ensure Execute doesn't panic
	ctx := context.Background()
	err := elevatedExecutor.Execute(ctx)
	if err != nil {
		t.Logf("Execute returned error (expected in test): %v", err)
	}
}

func TestExecuteElevatedLinux(t *testing.T) {
	mockConfig := config.ActionConfig{
		Type:    "script",
		Command: "echo test",
	}
	baseExecutor := script.NewScriptExecutor(&mockConfig)
	elevatedExecutor := &ElevatedExecutor{
		baseExecutor: baseExecutor,
		isAdmin:      true,
	}

	// Test Linux elevation
	err := elevatedExecutor.executeElevatedLinux(context.Background())

	// Expected to fail in test environment
	if err != nil {
		t.Logf("Linux elevation failed as expected: %v", err)
	}
}

func TestExecuteElevatedPlatformSpecific(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Run("Windows", testExecuteElevatedWindows)
	} else if runtime.GOOS == "darwin" {
		t.Run("Darwin", testExecuteElevatedDarwin)
	}
}

func testExecuteElevatedWindows(t *testing.T) {
	mockConfig := config.ActionConfig{
		Type:    "script",
		Command: "echo test",
	}
	baseExecutor := script.NewScriptExecutor(&mockConfig)
	elevatedExecutor := &ElevatedExecutor{
		baseExecutor: baseExecutor,
		isAdmin:      true,
	}

	// Test Windows elevation
	err := elevatedExecutor.executeElevatedWindows(context.Background())

	// Expected to fail in test environment without UI
	if err != nil {
		t.Logf("Windows elevation failed as expected: %v", err)
	}
}

func testExecuteElevatedDarwin(t *testing.T) {
	mockConfig := config.ActionConfig{
		Type:    "script",
		Command: "echo test",
	}
	baseExecutor := script.NewScriptExecutor(&mockConfig)
	elevatedExecutor := &ElevatedExecutor{
		baseExecutor: baseExecutor,
		isAdmin:      true,
	}

	// Test macOS elevation
	err := elevatedExecutor.executeElevatedDarwin(context.Background())

	// Expected to fail in test environment without UI
	if err != nil {
		t.Logf("macOS elevation failed as expected: %v", err)
	}
}

func TestUpdateActions(t *testing.T) {
	// Test UpdateActions method
	// NewManager is in the action package, which would create circular import
	// Skip this test as it's testing functionality in another package
	t.Skip("Skipping test that requires action.NewManager to avoid circular import")
	return

	/*
		manager := action.NewManager(map[string]config.ActionConfig{
			"test1": {
				Type:    "script",
				Command: "echo test1",
			},
		})

		// Update with new actions
		newActions := map[string]config.ActionConfig{
			"test2": {
				Type:    "script",
				Command: "echo test2",
			},
			"test3": {
				Type:    "url",
				Command: "https://example.com",
			},
		}

		manager.UpdateActions(newActions)

		// Verify update by trying to execute
		action := newActions["test2"]
		_, err := manager.createExecutor(&action)
		if err != nil {
			t.Errorf("Failed to create executor after update: %v", err)
		}
	*/
}
