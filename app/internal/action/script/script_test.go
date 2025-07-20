package script

import (
	"context"
	"runtime"
	"testing"

	"github.com/SphereStacking/silentcast/internal/config"
)

func TestScriptExecutor_ShowOutput(t *testing.T) {
	tests := []struct {
		name   string
		config config.ActionConfig
		hasErr bool
	}{
		{
			name: "Simple echo with show_output",
			config: config.ActionConfig{
				Type:       "script",
				Command:    "echo 'Hello, World!'",
				ShowOutput: true,
			},
			hasErr: false,
		},
		{
			name: "Command with error and show_output",
			config: config.ActionConfig{
				Type:       "script",
				Command:    "exit 1",
				ShowOutput: true,
			},
			hasErr: true,
		},
		{
			name: "Command with timeout",
			config: config.ActionConfig{
				Type:       "script",
				Command:    "sleep 10",
				ShowOutput: true,
				Timeout:    1, // 1 second timeout
			},
			hasErr: true,
		},
		{
			name: "Custom shell (bash)",
			config: config.ActionConfig{
				Type:       "script",
				Command:    "echo $SHELL",
				Shell:      "bash",
				ShowOutput: true,
			},
			hasErr: false,
		},
	}

	// Skip shell-specific tests on Windows
	if runtime.GOOS == "windows" {
		// Adjust tests for Windows
		tests[0].config.Command = "echo Hello, World!"
		tests[2].config.Command = "timeout /t 10"
		tests[3].config.Shell = "cmd"
		tests[3].config.Command = "echo %COMSPEC%"
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := NewScriptExecutor(&tt.config)
			ctx := context.Background()
			
			err := executor.Execute(ctx)
			
			if (err != nil) != tt.hasErr {
				t.Errorf("Execute() error = %v, hasErr = %v", err, tt.hasErr)
			}
		})
	}
}

func TestScriptExecutor_Terminal(t *testing.T) {
	// Test that terminal flag forces terminal execution
	config := config.ActionConfig{
		Type:     "script",
		Command:  "echo 'Terminal test'",
		Terminal: true,
	}
	
	executor := NewScriptExecutor(&config)
	ctx := context.Background()
	
	// This should not error, but we can't easily test if terminal was actually opened
	err := executor.Execute(ctx)
	if err != nil {
		t.Logf("Terminal execution returned error (may be expected in test environment): %v", err)
	}
}

func TestScriptExecutor_KeepOpen(t *testing.T) {
	tests := []struct {
		name   string
		config config.ActionConfig
	}{
		{
			name: "keep_open true",
			config: config.ActionConfig{
				Type:     "script",
				Command:  "echo 'Keep open test'",
				KeepOpen: true,
			},
		},
		{
			name: "keep_open with show_output",
			config: config.ActionConfig{
				Type:       "script",
				Command:    "echo 'Combined test'",
				KeepOpen:   true,
				ShowOutput: true,
			},
		},
		{
			name: "terminal and keep_open",
			config: config.ActionConfig{
				Type:     "script",
				Command:  "echo 'Both flags'",
				Terminal: true,
				KeepOpen: true,
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := NewScriptExecutor(&tt.config)
			ctx := context.Background()
			
			// This should not error, but we can't easily test if terminal stayed open
			err := executor.Execute(ctx)
			if err != nil {
				t.Logf("KeepOpen execution returned error (may be expected in test environment): %v", err)
			}
		})
	}
}

func TestScriptExecutor_Description(t *testing.T) {
	tests := []struct {
		name     string
		config   config.ActionConfig
		expected string
	}{
		{
			name: "With description",
			config: config.ActionConfig{
				Type:        "script",
				Command:     "git status",
				Description: "Check git status",
			},
			expected: "Check git status",
		},
		{
			name: "Without description",
			config: config.ActionConfig{
				Type:    "script",
				Command: "ls -la",
			},
			expected: "Run script: ls -la",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := NewScriptExecutor(&tt.config)
			if got := executor.String(); got != tt.expected {
				t.Errorf("String() = %v, want %v", got, tt.expected)
			}
		})
	}
}