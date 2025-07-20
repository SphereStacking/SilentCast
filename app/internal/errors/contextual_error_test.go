package errors

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

// TestContextualErrorFormatting tests that errors include contextual information
func TestContextualErrorFormatting(t *testing.T) {
	tests := []struct {
		name     string
		errType  ErrorType
		message  string
		context  map[string]interface{}
		want     string
	}{
		{
			name:    "config error with file path",
			errType: ErrorTypeConfig,
			message: "invalid configuration",
			context: map[string]interface{}{
				"file": "/path/to/spellbook.yml",
				"line": 42,
			},
			want: "invalid configuration", // Check base message, context order varies
		},
		{
			name:    "hotkey error with key sequence",
			errType: ErrorTypeHotkey,
			message: "key sequence failed",
			context: map[string]interface{}{
				"sequence": "alt+space,g,s",
				"step":     2,
			},
			want: "key sequence failed", // Check base message, context order varies
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := New(tt.errType, tt.message)
			for key, value := range tt.context {
				err = err.WithContext(key, value)
			}

			// This should format the error with context
			got := err.ErrorWithContext()
			
			// Check that base message is included
			if !strings.Contains(got, tt.want) {
				t.Errorf("ErrorWithContext() = %q, should contain %q", got, tt.want)
			}
			
			// Check that context is included  
			for key, value := range tt.context {
				expectedContext := fmt.Sprintf("%s: %v", key, value)
				if !strings.Contains(got, expectedContext) {
					t.Errorf("ErrorWithContext() = %q, should contain context %q", got, expectedContext)
				}
			}
		})
	}
}

// TestErrorChaining tests proper error chaining with context preservation
func TestErrorChaining(t *testing.T) {
	originalErr := errors.New("original system error")
	
	// Chain errors while preserving context
	err1 := Wrap(ErrorTypeSystem, "file operation failed", originalErr).
		WithContext("operation", "read").
		WithContext("file", "/tmp/test.txt")
	
	err2 := Wrap(ErrorTypeConfig, "configuration loading failed", err1).
		WithContext("config_type", "user")
	
	// Should be able to unwrap to original error
	if !errors.Is(err2, originalErr) {
		t.Error("error chain should preserve original error")
	}
	
	// Should preserve all context from chain
	contextStr := err2.FullContextString()
	wantContains := []string{
		"config_type: user",
		"operation: read", 
		"file: /tmp/test.txt",
	}
	
	for _, want := range wantContains {
		if !contains(contextStr, want) {
			t.Errorf("FullContextString() missing %q, got: %q", want, contextStr)
		}
	}
}

// TestStructuredLogging tests error formatting for structured logging
func TestStructuredLogging(t *testing.T) {
	err := New(ErrorTypeExecution, "spell execution failed").
		WithContext("spell", "editor").
		WithContext("command", "/usr/bin/code").
		WithContext("exit_code", 1)
	
	// Should provide structured data for logging
	fields := err.LogFields()
	
	expectedFields := map[string]interface{}{
		"error_type": "execution",
		"message":    "spell execution failed",
		"spell":      "editor", 
		"command":    "/usr/bin/code",
		"exit_code":  1,
	}
	
	for key, expectedValue := range expectedFields {
		if fields[key] != expectedValue {
			t.Errorf("LogFields()[%q] = %v, want %v", key, fields[key], expectedValue)
		}
	}
}

// Helper function for string containment check  
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}