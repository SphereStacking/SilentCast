package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
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
	err = err.WithContext("spell", "test_spell").WithContext("key", "ctrl+a")

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

// Test new error types added in unified error handling
func TestNewErrorTypes(t *testing.T) {
	tests := []struct {
		name     string
		errType  ErrorType
		expected string
	}{
		{"Validation error", ErrorTypeValidation, "validation"},
		{"Network error", ErrorTypeNetwork, "network"},
		{"IO error", ErrorTypeIO, "io"},
		{"Platform error", ErrorTypePlatform, "platform"},
		{"Timeout error", ErrorTypeTimeout, "timeout"},
		{"Not found error", ErrorTypeNotFound, "not_found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := New(tt.errType, "test message")
			fields := err.LogFields()

			assert.Equal(t, tt.expected, fields["error_type"])
			assert.Equal(t, "test message", fields["message"])
		})
	}
}

// Test sentinel errors
func TestSentinelErrors(t *testing.T) {
	t.Run("ErrConfigNotFound", func(t *testing.T) {
		err := ErrConfigNotFound.WithContext("path", "/test/config.yml")

		assert.True(t, errors.Is(err, ErrConfigNotFound))
		assert.True(t, IsType(err, ErrorTypeConfig))
		assert.Equal(t, "/test/config.yml", err.Context["path"])
	})

	t.Run("ErrSpellNotFound", func(t *testing.T) {
		err := ErrSpellNotFound.WithContext("spell", "unknown")

		assert.True(t, errors.Is(err, ErrSpellNotFound))
		assert.True(t, IsType(err, ErrorTypeConfig))
	})

	t.Run("ErrBrowserNotFound", func(t *testing.T) {
		err := ErrBrowserNotFound.WithContext("browsers_tried", []string{"chrome", "firefox"})

		assert.True(t, errors.Is(err, ErrBrowserNotFound))
		assert.True(t, IsType(err, ErrorTypeNotFound))
	})
}

// Test helper functions
func TestHelperFunctions(t *testing.T) {
	t.Run("NewValidationError", func(t *testing.T) {
		err := NewValidationError("command", "command cannot be empty")

		assert.True(t, IsType(err, ErrorTypeValidation))
		assert.Equal(t, "command cannot be empty", err.Message)
		assert.Equal(t, "command", err.Context["field"])
	})

	t.Run("NewIOError", func(t *testing.T) {
		err := NewIOError("/test/file.txt", "file not found")

		assert.True(t, IsType(err, ErrorTypeIO))
		assert.Equal(t, "file not found", err.Message)
		assert.Equal(t, "/test/file.txt", err.Context["path"])
	})

	t.Run("NewTimeoutError", func(t *testing.T) {
		err := NewTimeoutError("download", "30s")

		assert.True(t, IsType(err, ErrorTypeTimeout))
		assert.Contains(t, err.Message, "download timed out")
		assert.Equal(t, "download", err.Context["operation"])
		assert.Equal(t, "30s", err.Context["timeout"])
	})

	t.Run("NewNetworkError", func(t *testing.T) {
		err := NewNetworkError("https://example.com", "connection timeout")

		assert.True(t, IsType(err, ErrorTypeNetwork))
		assert.Equal(t, "connection timeout", err.Message)
		assert.Equal(t, "https://example.com", err.Context["url"])
	})
}

// Test WrapWithContext
func TestWrapWithContext(t *testing.T) {
	originalErr := errors.New("original error")
	context := map[string]interface{}{
		"file":      "test.txt",
		"operation": "read",
		"size":      1024,
	}

	err := WrapWithContext(ErrorTypeIO, "failed to read file", originalErr, context)

	assert.True(t, IsType(err, ErrorTypeIO))
	assert.Equal(t, "failed to read file", err.Message)
	assert.True(t, errors.Is(err, originalErr))
	assert.Equal(t, "test.txt", err.Context["file"])
	assert.Equal(t, "read", err.Context["operation"])
	assert.Equal(t, 1024, err.Context["size"])
}

// Test FromStandardError
func TestFromStandardError(t *testing.T) {
	t.Run("Standard error conversion", func(t *testing.T) {
		originalErr := errors.New("standard error")
		err := FromStandardError(originalErr, ErrorTypeSystem, "converted error")

		assert.True(t, IsType(err, ErrorTypeSystem))
		assert.Equal(t, "converted error", err.Message)
		assert.True(t, errors.Is(err, originalErr))
	})

	t.Run("SpellbookError passthrough", func(t *testing.T) {
		originalErr := New(ErrorTypeConfig, "original")
		err := FromStandardError(originalErr, ErrorTypeSystem, "converted error")

		// Should return the original SpellbookError unchanged
		assert.Equal(t, originalErr, err)
	})

	t.Run("Nil error", func(t *testing.T) {
		err := FromStandardError(nil, ErrorTypeSystem, "converted error")
		assert.Nil(t, err)
	})
}
