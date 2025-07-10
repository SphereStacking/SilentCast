package hotkey

import (
	"fmt"
	"strings"
	"time"
)

// Key represents a single key in a hotkey combination
type Key struct {
	Code      uint16   // Key code (platform specific)
	Modifiers []string // Modifier keys: ctrl, alt, shift, cmd/win
	Name      string   // Human-readable key name
}

// KeySequence represents a sequence of keys (e.g., "g,s" for git status)
type KeySequence struct {
	Keys []Key
}

// String returns the string representation of a key sequence
func (ks KeySequence) String() string {
	parts := make([]string, len(ks.Keys))
	for i, key := range ks.Keys {
		parts[i] = key.String()
	}
	return strings.Join(parts, ",")
}

// String returns the string representation of a key
func (k Key) String() string {
	if len(k.Modifiers) == 0 {
		return k.Name
	}
	return strings.Join(k.Modifiers, "+") + "+" + k.Name
}

// Event represents a hotkey event
type Event struct {
	Sequence  KeySequence
	SpellName string    // The spell to execute
	Timestamp time.Time
}


// HandlerFunc is a function adapter for Handler
type HandlerFunc func(Event) error

// Handle implements Handler interface
func (f HandlerFunc) Handle(event Event) error {
	return f(event)
}


// ParseError represents an error parsing a key sequence
type ParseError struct {
	Input   string
	Message string
}

func (e ParseError) Error() string {
	return fmt.Sprintf("failed to parse key sequence '%s': %s", e.Input, e.Message)
}

// ValidationError represents a validation error
type ValidationError struct {
	Sequence string
	Message  string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("invalid key sequence '%s': %s", e.Sequence, e.Message)
}