package errors

import (
	"errors"
	"fmt"
)

// Sentinel errors for common conditions
var (
	// ErrConfigNotFound indicates no configuration file was found
	ErrConfigNotFound = New(ErrorTypeConfig, "configuration file not found")
	
	// ErrInvalidConfig indicates the configuration is invalid
	ErrInvalidConfig = New(ErrorTypeConfig, "invalid configuration")
	
	// ErrSpellNotFound indicates a spell was not found
	ErrSpellNotFound = New(ErrorTypeConfig, "spell not found")
	
	// ErrActionNotFound indicates an action was not found
	ErrActionNotFound = New(ErrorTypeConfig, "action not found")
	
	// ErrPermissionDenied indicates insufficient permissions
	ErrPermissionDenied = New(ErrorTypePermission, "permission denied")
	
	// ErrHotkeyAlreadyRegistered indicates hotkey is already in use
	ErrHotkeyAlreadyRegistered = New(ErrorTypeHotkey, "hotkey already registered")
	
	// ErrHotkeyManagerNotInitialized indicates hotkey manager is not initialized
	ErrHotkeyManagerNotInitialized = New(ErrorTypeHotkey, "hotkey manager not initialized")
	
	// ErrExecutionTimeout indicates execution timed out
	ErrExecutionTimeout = New(ErrorTypeTimeout, "execution timeout")
	
	// ErrApplicationNotFound indicates application was not found
	ErrApplicationNotFound = New(ErrorTypeNotFound, "application not found")
	
	// ErrBrowserNotFound indicates no browser was found
	ErrBrowserNotFound = New(ErrorTypeNotFound, "no browser found")
	
	// ErrFileNotFound indicates file was not found
	ErrFileNotFound = New(ErrorTypeIO, "file not found")
	
	// ErrNetworkUnavailable indicates network is unavailable
	ErrNetworkUnavailable = New(ErrorTypeNetwork, "network unavailable")
)

// ErrorType represents the type of error
type ErrorType int

const (
	// ErrorTypeUnknown is for unknown errors
	ErrorTypeUnknown ErrorType = iota
	// ErrorTypeConfig is for configuration errors
	ErrorTypeConfig
	// ErrorTypePermission is for permission errors
	ErrorTypePermission
	// ErrorTypeHotkey is for hotkey errors
	ErrorTypeHotkey
	// ErrorTypeExecution is for execution errors
	ErrorTypeExecution
	// ErrorTypeSystem is for system errors
	ErrorTypeSystem
	// ErrorTypeValidation is for validation errors
	ErrorTypeValidation
	// ErrorTypeNetwork is for network and HTTP errors
	ErrorTypeNetwork
	// ErrorTypeIO is for file I/O errors
	ErrorTypeIO
	// ErrorTypePlatform is for platform-specific errors
	ErrorTypePlatform
	// ErrorTypeTimeout is for timeout errors
	ErrorTypeTimeout
	// ErrorTypeNotFound is for resource not found errors
	ErrorTypeNotFound
)

// SpellbookError is the base error type for all Spellbook errors
type SpellbookError struct {
	Type    ErrorType
	Message string
	Cause   error
	Context map[string]interface{}
}

// Error implements the error interface
func (e *SpellbookError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap implements the errors.Unwrap interface
func (e *SpellbookError) Unwrap() error {
	return e.Cause
}

// WithContext adds context to the error
func (e *SpellbookError) WithContext(key string, value interface{}) *SpellbookError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// New creates a new SpellbookError
func New(errType ErrorType, message string) *SpellbookError {
	return &SpellbookError{
		Type:    errType,
		Message: message,
	}
}

// Wrap wraps an error with a SpellbookError
func Wrap(errType ErrorType, message string, cause error) *SpellbookError {
	return &SpellbookError{
		Type:    errType,
		Message: message,
		Cause:   cause,
	}
}

// IsType checks if an error is of a specific type
func IsType(err error, errType ErrorType) bool {
	if err == nil {
		return false
	}

	var spellErr *SpellbookError
	if !errors.As(err, &spellErr) {
		return false
	}

	return spellErr.Type == errType
}

// ErrorWithContext returns error message with context formatted as key-value pairs
func (e *SpellbookError) ErrorWithContext() string {
	if len(e.Context) == 0 {
		return e.Message
	}
	
	contextParts := make([]string, 0, len(e.Context))
	for key, value := range e.Context {
		contextParts = append(contextParts, fmt.Sprintf("%s: %v", key, value))
	}
	
	return fmt.Sprintf("%s (%s)", e.Message, joinContextParts(contextParts))
}

// FullContextString returns all context information from the error chain
func (e *SpellbookError) FullContextString() string {
	contextMap := make(map[string]interface{})
	
	// Collect context from entire error chain
	current := e
	for current != nil {
		for key, value := range current.Context {
			if _, exists := contextMap[key]; !exists {
				contextMap[key] = value
			}
		}
		
		// Continue with wrapped error if it's also a SpellbookError
		var nextErr *SpellbookError
		if errors.As(current.Cause, &nextErr) {
			current = nextErr
		} else {
			break
		}
	}
	
	if len(contextMap) == 0 {
		return ""
	}
	
	contextParts := make([]string, 0, len(contextMap))
	for key, value := range contextMap {
		contextParts = append(contextParts, fmt.Sprintf("%s: %v", key, value))
	}
	
	return joinContextParts(contextParts)
}

// LogFields returns structured fields for logging
func (e *SpellbookError) LogFields() map[string]interface{} {
	fields := make(map[string]interface{})
	
	// Add error type
	switch e.Type {
	case ErrorTypeConfig:
		fields["error_type"] = "config"
	case ErrorTypePermission:
		fields["error_type"] = "permission"
	case ErrorTypeHotkey:
		fields["error_type"] = "hotkey"
	case ErrorTypeExecution:
		fields["error_type"] = "execution"
	case ErrorTypeSystem:
		fields["error_type"] = "system"
	case ErrorTypeValidation:
		fields["error_type"] = "validation"
	case ErrorTypeNetwork:
		fields["error_type"] = "network"
	case ErrorTypeIO:
		fields["error_type"] = "io"
	case ErrorTypePlatform:
		fields["error_type"] = "platform"
	case ErrorTypeTimeout:
		fields["error_type"] = "timeout"
	case ErrorTypeNotFound:
		fields["error_type"] = "not_found"
	default:
		fields["error_type"] = "unknown"
	}
	
	// Add message
	fields["message"] = e.Message
	
	// Add context fields
	for key, value := range e.Context {
		fields[key] = value
	}
	
	return fields
}

// Helper function to join context parts consistently
func joinContextParts(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	if len(parts) == 1 {
		return parts[0]
	}
	
	result := parts[0]
	for i := 1; i < len(parts); i++ {
		result += ", " + parts[i]
	}
	return result
}

// GetUserMessage returns a user-friendly error message
func GetUserMessage(err error) string {
	if err == nil {
		return ""
	}

	var spellErr *SpellbookError
	if !errors.As(err, &spellErr) {
		return "An unexpected error occurred"
	}

	switch spellErr.Type {
	case ErrorTypeConfig:
		return fmt.Sprintf("Configuration error: %s", spellErr.Message)
	case ErrorTypePermission:
		return fmt.Sprintf("Permission error: %s. Please check the permissions guide.", spellErr.Message)
	case ErrorTypeHotkey:
		return fmt.Sprintf("Hotkey error: %s", spellErr.Message)
	case ErrorTypeExecution:
		return fmt.Sprintf("Execution error: %s", spellErr.Message)
	case ErrorTypeSystem:
		return fmt.Sprintf("System error: %s", spellErr.Message)
	case ErrorTypeValidation:
		return fmt.Sprintf("Validation error: %s", spellErr.Message)
	case ErrorTypeNetwork:
		return fmt.Sprintf("Network error: %s. Please check your connection.", spellErr.Message)
	case ErrorTypeIO:
		return fmt.Sprintf("File error: %s", spellErr.Message)
	case ErrorTypePlatform:
		return fmt.Sprintf("Platform error: %s. This may be platform-specific.", spellErr.Message)
	case ErrorTypeTimeout:
		return fmt.Sprintf("Timeout error: %s. The operation took too long.", spellErr.Message)
	case ErrorTypeNotFound:
		return fmt.Sprintf("Not found: %s", spellErr.Message)
	default:
		return fmt.Sprintf("Error: %s", spellErr.Message)
	}
}

// Helper functions for common error patterns

// NewConfigError creates a new configuration error with context
func NewConfigError(message string) *SpellbookError {
	return New(ErrorTypeConfig, message)
}

// NewValidationError creates a new validation error with field context
func NewValidationError(field, message string) *SpellbookError {
	return New(ErrorTypeValidation, message).
		WithContext("field", field)
}

// NewPermissionError creates a new permission error with resource context
func NewPermissionError(resource, message string) *SpellbookError {
	return New(ErrorTypePermission, message).
		WithContext("resource", resource)
}

// NewExecutionError creates a new execution error with command context
func NewExecutionError(command, message string) *SpellbookError {
	return New(ErrorTypeExecution, message).
		WithContext("command", command)
}

// NewIOError creates a new I/O error with file path context
func NewIOError(path, message string) *SpellbookError {
	return New(ErrorTypeIO, message).
		WithContext("path", path)
}

// NewNetworkError creates a new network error with URL context
func NewNetworkError(url, message string) *SpellbookError {
	return New(ErrorTypeNetwork, message).
		WithContext("url", url)
}

// NewTimeoutError creates a new timeout error with duration context
func NewTimeoutError(operation string, timeout interface{}) *SpellbookError {
	return New(ErrorTypeTimeout, fmt.Sprintf("%s timed out", operation)).
		WithContext("operation", operation).
		WithContext("timeout", timeout)
}

// NewNotFoundError creates a new not found error with resource context
func NewNotFoundError(resource, name string) *SpellbookError {
	return New(ErrorTypeNotFound, fmt.Sprintf("%s not found", resource)).
		WithContext("resource", resource).
		WithContext("name", name)
}

// WrapWithContext wraps an error with additional context
func WrapWithContext(errType ErrorType, message string, cause error, context map[string]interface{}) *SpellbookError {
	err := Wrap(errType, message, cause)
	for key, value := range context {
		err = err.WithContext(key, value)
	}
	return err
}

// FromStandardError converts a standard error to SpellbookError if it isn't already
func FromStandardError(err error, errType ErrorType, message string) *SpellbookError {
	if err == nil {
		return nil
	}
	
	var spellErr *SpellbookError
	if errors.As(err, &spellErr) {
		return spellErr
	}
	
	return Wrap(errType, message, err)
}
