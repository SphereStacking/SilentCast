package errors

import (
	"errors"
	"testing"
)

func TestSpellbookError(t *testing.T) {
	// Test New
	err := New(ErrorTypeConfig, "config error")
	if err.Type != ErrorTypeConfig {
		t.Errorf("Expected type %v, got %v", ErrorTypeConfig, err.Type)
	}
	if err.Message != "config error" {
		t.Errorf("Expected message 'config error', got '%s'", err.Message)
	}

	// Test Wrap
	cause := errors.New("underlying error")
	wrapped := Wrap(ErrorTypeSystem, "system error", cause)
	if wrapped.Cause != cause { //nolint:errorlint // Testing exact object reference
		t.Error("Expected wrapped error to contain cause")
	}
	if !errors.Is(wrapped, cause) {
		t.Error("errors.Is should work with wrapped error")
	}

	// Test WithContext
	err = New(ErrorTypeExecution, "exec error")
	err.WithContext("spell", "test_spell").WithContext("key", "ctrl+a")

	if err.Context["spell"] != "test_spell" {
		t.Errorf("Expected context spell='test_spell', got '%v'", err.Context["spell"])
	}
	if err.Context["key"] != "ctrl+a" {
		t.Errorf("Expected context key='ctrl+a', got '%v'", err.Context["key"])
	}
}

func TestIsType(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		errType  ErrorType
		expected bool
	}{
		{
			name:     "Matching type",
			err:      New(ErrorTypeConfig, "config error"),
			errType:  ErrorTypeConfig,
			expected: true,
		},
		{
			name:     "Non-matching type",
			err:      New(ErrorTypeConfig, "config error"),
			errType:  ErrorTypePermission,
			expected: false,
		},
		{
			name:     "Nil error",
			err:      nil,
			errType:  ErrorTypeConfig,
			expected: false,
		},
		{
			name:     "Non-SpellbookError",
			err:      errors.New("regular error"),
			errType:  ErrorTypeConfig,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsType(tt.err, tt.errType); got != tt.expected {
				t.Errorf("IsType() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetUserMessage(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "Config error",
			err:  New(ErrorTypeConfig, "invalid YAML"),
			want: "Configuration error: invalid YAML",
		},
		{
			name: "Permission error",
			err:  New(ErrorTypePermission, "accessibility not granted"),
			want: "Permission error: accessibility not granted. Please check the permissions guide.",
		},
		{
			name: "Hotkey error",
			err:  New(ErrorTypeHotkey, "invalid key sequence"),
			want: "Hotkey error: invalid key sequence",
		},
		{
			name: "Execution error",
			err:  New(ErrorTypeExecution, "app not found"),
			want: "Execution error: app not found",
		},
		{
			name: "System error",
			err:  New(ErrorTypeSystem, "out of memory"),
			want: "System error: out of memory",
		},
		{
			name: "Unknown error type",
			err:  New(ErrorTypeUnknown, "something went wrong"),
			want: "Error: something went wrong",
		},
		{
			name: "Nil error",
			err:  nil,
			want: "",
		},
		{
			name: "Non-SpellbookError",
			err:  errors.New("regular error"),
			want: "An unexpected error occurred",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetUserMessage(tt.err); got != tt.want {
				t.Errorf("GetUserMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Error(t *testing.T) {
	// Test without cause
	err := New(ErrorTypeConfig, "config error")
	if err.Error() != "config error" {
		t.Errorf("Error() = %v, want %v", err.Error(), "config error")
	}

	// Test with cause
	cause := errors.New("underlying issue")
	wrapped := Wrap(ErrorTypeConfig, "config error", cause)
	expected := "config error: underlying issue"
	if wrapped.Error() != expected {
		t.Errorf("Error() = %v, want %v", wrapped.Error(), expected)
	}
}
