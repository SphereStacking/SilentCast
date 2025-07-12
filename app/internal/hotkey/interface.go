package hotkey

// Handler processes hotkey events
type Handler interface {
	// Handle processes a hotkey event
	Handle(event Event) error
}

// Manager manages hotkey registration and detection
type Manager interface {
	// Start begins listening for hotkeys
	Start() error

	// Stop stops listening for hotkeys
	Stop() error

	// Register registers a hotkey sequence with a spell name
	Register(sequence string, spellName string) error

	// Unregister removes a hotkey registration
	Unregister(sequence string) error

	// SetHandler sets the event handler
	SetHandler(handler Handler)

	// IsRunning returns whether the manager is actively listening
	IsRunning() bool
}
