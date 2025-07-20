package terminal

import (
	"context"
	"fmt"
	"os/exec"
	"sync"
)

// baseManager provides common functionality for terminal managers
type baseManager struct {
	detector Detector
	builder  CommandBuilder

	// Cache for available terminals
	terminalCache []Terminal
	cacheMu       sync.RWMutex
	cacheValid    bool
}

// newBaseManager creates a new base manager with the given detector and builder
func newBaseManager(detector Detector, builder CommandBuilder) *baseManager {
	return &baseManager{
		detector:   detector,
		builder:    builder,
		cacheValid: false,
	}
}

// ExecuteInTerminal implements the Manager interface
func (m *baseManager) ExecuteInTerminal(ctx context.Context, cmd *exec.Cmd, options Options) error {
	// Get the terminal to use
	terminal, err := m.selectTerminal(options)
	if err != nil {
		return err
	}

	// Build the terminal command
	args, err := m.builder.BuildCommand(terminal, cmd, options)
	if err != nil {
		return &Error{
			Op:  "ExecuteInTerminal",
			Err: err,
			Msg: fmt.Sprintf("failed to build command for terminal %s", terminal.Name),
		}
	}

	// Create the terminal command
	termCmd := exec.CommandContext(ctx, terminal.Command, args...)

	// Set working directory if specified
	if options.WorkingDir != "" {
		termCmd.Dir = options.WorkingDir
	} else if cmd.Dir != "" {
		termCmd.Dir = cmd.Dir
	}

	// Start the terminal
	if err := termCmd.Start(); err != nil {
		return &Error{
			Op:  "ExecuteInTerminal",
			Err: err,
			Msg: fmt.Sprintf("failed to start terminal %s", terminal.Name),
		}
	}

	// Detach from the process so it continues running independently
	if err := termCmd.Process.Release(); err != nil {
		// Non-fatal error, log but don't fail
		_ = err
	}

	return nil
}

// GetAvailableTerminals returns cached terminals or detects them
func (m *baseManager) GetAvailableTerminals() []Terminal {
	m.cacheMu.RLock()
	if m.cacheValid {
		terminals := make([]Terminal, len(m.terminalCache))
		copy(terminals, m.terminalCache)
		m.cacheMu.RUnlock()
		return terminals
	}
	m.cacheMu.RUnlock()

	// Need to detect terminals
	m.cacheMu.Lock()
	defer m.cacheMu.Unlock()

	// Double-check in case another goroutine already updated
	if m.cacheValid {
		terminals := make([]Terminal, len(m.terminalCache))
		copy(terminals, m.terminalCache)
		return terminals
	}

	// Detect terminals
	m.terminalCache = m.detector.DetectTerminals()
	m.cacheValid = true

	terminals := make([]Terminal, len(m.terminalCache))
	copy(terminals, m.terminalCache)
	return terminals
}

// GetDefaultTerminal returns the highest priority terminal
func (m *baseManager) GetDefaultTerminal() (Terminal, error) {
	terminals := m.GetAvailableTerminals()
	if len(terminals) == 0 {
		return Terminal{}, ErrNoTerminalFound
	}

	// Find the default terminal or the highest priority one
	var defaultTerminal Terminal
	highestPriority := -1

	for _, t := range terminals {
		if t.IsDefault {
			return t, nil
		}
		if t.Priority > highestPriority {
			defaultTerminal = t
			highestPriority = t.Priority
		}
	}

	if highestPriority == -1 {
		return Terminal{}, ErrNoTerminalFound
	}

	return defaultTerminal, nil
}

// IsTerminalAvailable checks if a specific terminal is available
func (m *baseManager) IsTerminalAvailable(terminal Terminal) bool {
	terminals := m.GetAvailableTerminals()
	for _, t := range terminals {
		if t.Command == terminal.Command || t.Name == terminal.Name {
			return true
		}
	}
	return false
}

// selectTerminal chooses the appropriate terminal based on options
func (m *baseManager) selectTerminal(options Options) (Terminal, error) {
	// If a preferred terminal is specified, try to use it
	if options.PreferredTerminal.Command != "" {
		if m.IsTerminalAvailable(options.PreferredTerminal) {
			return options.PreferredTerminal, nil
		}
		// Fall back to default if preferred is not available
	}

	// Use the default terminal
	return m.GetDefaultTerminal()
}

// invalidateCache forces re-detection of terminals on next access
func (m *baseManager) invalidateCache() {
	m.cacheMu.Lock()
	m.cacheValid = false
	m.cacheMu.Unlock()
}
