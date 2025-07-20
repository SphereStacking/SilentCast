package terminal

import (
	"testing"
)

func TestNewManager(t *testing.T) {
	manager := NewManager()
	
	if manager == nil {
		t.Fatal("NewManager() returned nil")
	}
	
	// Test that manager implements the Manager interface
	var _ Manager = manager
	
	// Test basic functionality
	terminals := manager.GetAvailableTerminals()
	if terminals == nil {
		t.Error("GetAvailableTerminals() returned nil")
	}
	
	// Should not panic even if no terminals are available
	_, err := manager.GetDefaultTerminal()
	if err != nil {
		t.Logf("GetDefaultTerminal() returned error (expected in test environment): %v", err)
	}
}

func TestManagerImplementsInterface(t *testing.T) {
	manager := NewManager()
	
	// Verify all interface methods are implemented
	terminals := manager.GetAvailableTerminals()
	_ = terminals
	
	_, err := manager.GetDefaultTerminal()
	_ = err // Error is acceptable in test environment
	
	// Test with a dummy terminal
	dummyTerminal := Terminal{
		Name:    "Dummy",
		Command: "dummy",
	}
	
	available := manager.IsTerminalAvailable(dummyTerminal)
	_ = available // Should return false for non-existent terminal
}