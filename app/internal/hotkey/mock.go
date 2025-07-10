package hotkey

import (
	"sync"
	"time"
)

// MockManager is a mock implementation of Manager for testing
type MockManager struct {
	mu         sync.RWMutex
	running    bool
	handler    Handler
	sequences  map[string]string
	startErr   error
	stopErr    error
}

// NewMockManager creates a new mock manager
func NewMockManager() *MockManager {
	return &MockManager{
		sequences: make(map[string]string),
	}
}

// Start begins listening for hotkeys
func (m *MockManager) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if m.startErr != nil {
		return m.startErr
	}
	
	m.running = true
	return nil
}

// Stop stops listening for hotkeys
func (m *MockManager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if m.stopErr != nil {
		return m.stopErr
	}
	
	m.running = false
	return nil
}

// Register registers a hotkey sequence with a spell name
func (m *MockManager) Register(sequence string, spellName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.sequences[sequence] = spellName
	return nil
}

// Unregister removes a hotkey registration
func (m *MockManager) Unregister(sequence string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	delete(m.sequences, sequence)
	return nil
}

// SetHandler sets the event handler
func (m *MockManager) SetHandler(handler Handler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handler = handler
}

// IsRunning returns whether the manager is actively listening
func (m *MockManager) IsRunning() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.running
}

// SimulateKeyPress simulates a key press for testing
func (m *MockManager) SimulateKeyPress(sequence string) error {
	m.mu.RLock()
	handler := m.handler
	spellName, exists := m.sequences[sequence]
	m.mu.RUnlock()
	
	if !exists {
		return nil
	}
	
	if handler != nil {
		parser := NewParser()
		keySeq, err := parser.Parse(sequence)
		if err != nil {
			return err
		}
		
		event := Event{
			Sequence:  keySeq,
			SpellName: spellName,
			Timestamp: time.Now(),
		}
		
		return handler.Handle(event)
	}
	
	return nil
}

// SetStartError sets an error to be returned by Start
func (m *MockManager) SetStartError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.startErr = err
}

// SetStopError sets an error to be returned by Stop
func (m *MockManager) SetStopError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stopErr = err
}