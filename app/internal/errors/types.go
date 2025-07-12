package errors

import (
	"errors"
	"fmt"
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
	default:
		return fmt.Sprintf("Error: %s", spellErr.Message)
	}
}
