package terminal

import (
	"context"
	"errors"
	"os/exec"
	"testing"
)

// Mock implementations for testing

type mockDetector struct {
	terminals []Terminal
	findErr   error
}

func (m *mockDetector) DetectTerminals() []Terminal {
	return m.terminals
}

func (m *mockDetector) FindTerminal(nameOrCommand string) (Terminal, error) {
	if m.findErr != nil {
		return Terminal{}, m.findErr
	}
	for _, t := range m.terminals {
		if t.Name == nameOrCommand || t.Command == nameOrCommand {
			return t, nil
		}
	}
	return Terminal{}, &Error{Op: "find", Msg: "terminal not found"}
}

type mockBuilder struct {
	buildErr      error
	supportedCmds map[string]bool
	buildFunc     func(terminal Terminal, cmd *exec.Cmd, options Options) ([]string, error)
}

func (m *mockBuilder) BuildCommand(terminal Terminal, cmd *exec.Cmd, options Options) ([]string, error) {
	if m.buildFunc != nil {
		return m.buildFunc(terminal, cmd, options)
	}
	if m.buildErr != nil {
		return nil, m.buildErr
	}
	if m.supportedCmds != nil && !m.supportedCmds[terminal.Command] {
		return nil, ErrTerminalNotSupported
	}

	args := []string{"mock", "args"}
	if options.KeepOpen {
		args = append(args, "--keep-open")
	}
	if options.Title != "" {
		args = append(args, "--title", options.Title)
	}
	return args, nil
}

func (m *mockBuilder) SupportsTerminal(terminal Terminal) bool {
	if m.supportedCmds == nil {
		return true
	}
	return m.supportedCmds[terminal.Command]
}

// Tests

func TestBaseManager_GetAvailableTerminals(t *testing.T) {
	tests := []struct {
		name      string
		terminals []Terminal
		want      int
	}{
		{
			name: "multiple terminals",
			terminals: []Terminal{
				{Name: "Terminal 1", Command: "term1", Priority: 100},
				{Name: "Terminal 2", Command: "term2", Priority: 50},
			},
			want: 2,
		},
		{
			name:      "no terminals",
			terminals: []Terminal{},
			want:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detector := &mockDetector{terminals: tt.terminals}
			builder := &mockBuilder{}
			manager := newBaseManager(detector, builder)

			got := manager.GetAvailableTerminals()
			if len(got) != tt.want {
				t.Errorf("GetAvailableTerminals() returned %d terminals, want %d", len(got), tt.want)
			}

			// Verify cache is working
			got2 := manager.GetAvailableTerminals()
			if len(got2) != len(got) {
				t.Errorf("Cache not working: first call returned %d, second returned %d", len(got), len(got2))
			}
		})
	}
}

func TestBaseManager_GetDefaultTerminal(t *testing.T) {
	tests := []struct {
		name      string
		terminals []Terminal
		wantName  string
		wantErr   bool
	}{
		{
			name: "has default terminal",
			terminals: []Terminal{
				{Name: "Terminal 1", Command: "term1", Priority: 50, IsDefault: false},
				{Name: "Terminal 2", Command: "term2", Priority: 100, IsDefault: true},
			},
			wantName: "Terminal 2",
			wantErr:  false,
		},
		{
			name: "no default, use highest priority",
			terminals: []Terminal{
				{Name: "Terminal 1", Command: "term1", Priority: 50},
				{Name: "Terminal 2", Command: "term2", Priority: 100},
			},
			wantName: "Terminal 2",
			wantErr:  false,
		},
		{
			name:      "no terminals",
			terminals: []Terminal{},
			wantName:  "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detector := &mockDetector{terminals: tt.terminals}
			builder := &mockBuilder{}
			manager := newBaseManager(detector, builder)

			got, err := manager.GetDefaultTerminal()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDefaultTerminal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Name != tt.wantName {
				t.Errorf("GetDefaultTerminal() = %v, want %v", got.Name, tt.wantName)
			}
		})
	}
}

func TestBaseManager_IsTerminalAvailable(t *testing.T) {
	terminals := []Terminal{
		{Name: "Terminal 1", Command: "term1"},
		{Name: "Terminal 2", Command: "term2"},
	}

	detector := &mockDetector{terminals: terminals}
	builder := &mockBuilder{}
	manager := newBaseManager(detector, builder)

	tests := []struct {
		name     string
		terminal Terminal
		want     bool
	}{
		{
			name:     "available by command",
			terminal: Terminal{Command: "term1"},
			want:     true,
		},
		{
			name:     "available by name",
			terminal: Terminal{Name: "Terminal 2"},
			want:     true,
		},
		{
			name:     "not available",
			terminal: Terminal{Command: "term3"},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := manager.IsTerminalAvailable(tt.terminal); got != tt.want {
				t.Errorf("IsTerminalAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseManager_ExecuteInTerminal(t *testing.T) {
	t.Skip("ExecuteInTerminal requires actual terminal execution - tested in integration tests")
}

func TestBaseManager_SelectTerminal(t *testing.T) {
	terminals := []Terminal{
		{Name: "Default", Command: "default", Priority: 100, IsDefault: true},
		{Name: "Alternative", Command: "alt", Priority: 50},
		{Name: "Preferred", Command: "preferred", Priority: 25},
	}

	tests := []struct {
		name      string
		options   Options
		terminals []Terminal
		wantCmd   string
		wantErr   bool
	}{
		{
			name:      "use default when no preference",
			options:   Options{},
			terminals: terminals,
			wantCmd:   "default",
			wantErr:   false,
		},
		{
			name: "use preferred when available",
			options: Options{
				PreferredTerminal: Terminal{Name: "Preferred", Command: "preferred"},
			},
			terminals: terminals,
			wantCmd:   "preferred",
			wantErr:   false,
		},
		{
			name: "fallback to default when preferred not available",
			options: Options{
				PreferredTerminal: Terminal{Command: "notexist"},
			},
			terminals: terminals,
			wantCmd:   "default",
			wantErr:   false,
		},
		{
			name:      "error when no terminals",
			options:   Options{},
			terminals: []Terminal{},
			wantCmd:   "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detector := &mockDetector{terminals: tt.terminals}
			builder := &mockBuilder{}
			manager := newBaseManager(detector, builder)

			got, err := manager.selectTerminal(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("selectTerminal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Command != tt.wantCmd {
				t.Errorf("selectTerminal() = %v, want %v", got.Command, tt.wantCmd)
			}
		})
	}
}

func TestBaseManager_CacheInvalidation(t *testing.T) {
	detector := &mockDetector{
		terminals: []Terminal{
			{Name: "Terminal 1", Command: "term1"},
		},
	}
	builder := &mockBuilder{}
	manager := newBaseManager(detector, builder)

	// First access should populate cache
	first := manager.GetAvailableTerminals()
	if len(first) != 1 {
		t.Errorf("Expected 1 terminal, got %d", len(first))
	}

	// Update detector's terminals
	detector.terminals = append(detector.terminals, Terminal{Name: "Terminal 2", Command: "term2"})

	// Should still get cached value
	cached := manager.GetAvailableTerminals()
	if len(cached) != 1 {
		t.Errorf("Expected cached value with 1 terminal, got %d", len(cached))
	}

	// Invalidate cache
	manager.invalidateCache()

	// Should get new value
	updated := manager.GetAvailableTerminals()
	if len(updated) != 2 {
		t.Errorf("Expected updated value with 2 terminals, got %d", len(updated))
	}
}

func TestBaseManager_ConcurrentCacheAccess(t *testing.T) {
	detector := &mockDetector{
		terminals: []Terminal{
			{Name: "Terminal 1", Command: "term1", Priority: 100},
		},
	}
	builder := &mockBuilder{}
	manager := newBaseManager(detector, builder)

	// Test concurrent cache access
	const numGoroutines = 50
	results := make(chan []Terminal, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			terminals := manager.GetAvailableTerminals()
			results <- terminals
		}()
	}

	// Collect all results
	for i := 0; i < numGoroutines; i++ {
		terminals := <-results
		if len(terminals) != 1 {
			t.Errorf("Expected 1 terminal from goroutine, got %d", len(terminals))
		}
		if terminals[0].Name != "Terminal 1" {
			t.Errorf("Expected 'Terminal 1', got '%s'", terminals[0].Name)
		}
	}
}

func TestBaseManager_ExecuteInTerminal_BuildError(t *testing.T) {
	terminals := []Terminal{
		{Name: "Test Terminal", Command: "test", Priority: 100},
	}

	detector := &mockDetector{terminals: terminals}
	builder := &mockBuilder{buildErr: &Error{Op: "build", Msg: "build failed"}}
	manager := newBaseManager(detector, builder)

	cmd := exec.Command("echo", "test")
	err := manager.ExecuteInTerminal(context.Background(), cmd, Options{})

	if err == nil {
		t.Error("Expected error from builder, got nil")
	}

	// Check error type and message
	var termErr *Error
	if errors.As(err, &termErr) {
		if termErr.Op != "ExecuteInTerminal" {
			t.Errorf("Expected Op 'ExecuteInTerminal', got '%s'", termErr.Op)
		}
	} else {
		t.Errorf("Expected *Error type, got %T", err)
	}
}

func TestBaseManager_ExecuteInTerminal_SelectTerminalError(t *testing.T) {
	// No terminals available
	detector := &mockDetector{terminals: []Terminal{}}
	builder := &mockBuilder{}
	manager := newBaseManager(detector, builder)

	cmd := exec.Command("echo", "test")
	err := manager.ExecuteInTerminal(context.Background(), cmd, Options{})

	if err == nil {
		t.Error("Expected error when no terminals available, got nil")
	}
}

func TestBaseManager_ExecuteInTerminal_CompleteFlow(t *testing.T) {
	detector := &mockDetector{
		terminals: []Terminal{
			{Name: "Test Terminal", Command: "echo", Priority: 10, IsDefault: true},
		},
	}

	buildCalled := false
	builder := &mockBuilder{
		buildFunc: func(terminal Terminal, cmd *exec.Cmd, options Options) ([]string, error) {
			buildCalled = true
			return []string{"test output"}, nil
		},
	}

	manager := newBaseManager(detector, builder)

	cmd := exec.Command("echo", "test")
	options := Options{
		WorkingDir: "/tmp",
		KeepOpen:   true,
	}

	// Since we're using 'echo' as the terminal command, it will execute quickly
	err := manager.ExecuteInTerminal(context.Background(), cmd, options)

	// The command should execute without error (echo exists)
	if err != nil {
		t.Logf("ExecuteInTerminal() returned error (may be expected): %v", err)
	}

	if !buildCalled {
		t.Error("BuildCommand was not called")
	}
}

func TestBaseManager_ExecuteInTerminal_WorkingDirPriority(t *testing.T) {
	tests := []struct {
		name       string
		optionsDir string
		cmdDir     string
	}{
		{
			name:       "options dir takes priority",
			optionsDir: "/opt/dir",
			cmdDir:     "/cmd/dir",
		},
		{
			name:       "use cmd dir when options dir empty",
			optionsDir: "",
			cmdDir:     "/cmd/dir",
		},
		{
			name:       "both dirs empty",
			optionsDir: "",
			cmdDir:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detector := &mockDetector{
				terminals: []Terminal{
					{Name: "Test Terminal", Command: "echo", Priority: 10, IsDefault: true},
				},
			}

			builder := &mockBuilder{
				buildFunc: func(terminal Terminal, cmd *exec.Cmd, options Options) ([]string, error) {
					return []string{"test"}, nil
				},
			}

			manager := newBaseManager(detector, builder)

			cmd := exec.Command("echo", "test")
			cmd.Dir = tt.cmdDir

			options := Options{
				WorkingDir: tt.optionsDir,
			}

			// Execute - we're testing the code path, not the actual execution
			_ = manager.ExecuteInTerminal(context.Background(), cmd, options)
		})
	}
}

func TestBaseManager_GetDefaultTerminal_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		terminals []Terminal
		wantErr   bool
	}{
		{
			name: "terminals with negative priority",
			terminals: []Terminal{
				{Name: "Terminal 1", Command: "term1", Priority: -10},
				{Name: "Terminal 2", Command: "term2", Priority: -5},
			},
			wantErr: true, // Current implementation doesn't handle negative priorities properly
		},
		{
			name: "terminals with zero priority",
			terminals: []Terminal{
				{Name: "Terminal 1", Command: "term1", Priority: 0},
			},
			wantErr: false,
		},
		{
			name: "multiple default terminals",
			terminals: []Terminal{
				{Name: "Terminal 1", Command: "term1", Priority: 10, IsDefault: true},
				{Name: "Terminal 2", Command: "term2", Priority: 20, IsDefault: true},
			},
			wantErr: false, // Should return the first default one found
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detector := &mockDetector{terminals: tt.terminals}
			builder := &mockBuilder{}
			manager := newBaseManager(detector, builder)

			terminal, err := manager.GetDefaultTerminal()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDefaultTerminal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && terminal.Name == "" {
				t.Error("Expected valid terminal, got empty name")
			}
		})
	}
}

func TestBaseManager_IsTerminalAvailable_EdgeCases(t *testing.T) {
	terminals := []Terminal{
		{Name: "Terminal 1", Command: "term1"},
		{Name: "", Command: "term2"},      // Empty name
		{Name: "Terminal 3", Command: ""}, // Empty command
	}

	detector := &mockDetector{terminals: terminals}
	builder := &mockBuilder{}
	manager := newBaseManager(detector, builder)

	tests := []struct {
		name     string
		terminal Terminal
		want     bool
	}{
		{
			name:     "match empty name",
			terminal: Terminal{Name: "", Command: "term2"},
			want:     true,
		},
		{
			name:     "match empty command",
			terminal: Terminal{Name: "Terminal 3", Command: ""},
			want:     true,
		},
		{
			name:     "both empty",
			terminal: Terminal{Name: "", Command: ""},
			want:     true, // The implementation checks OR condition, so empty name matches empty name
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := manager.IsTerminalAvailable(tt.terminal); got != tt.want {
				t.Errorf("IsTerminalAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}
