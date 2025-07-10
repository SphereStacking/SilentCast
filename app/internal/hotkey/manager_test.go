package hotkey

import (
	"sync"
	"testing"
	"time"

	"github.com/SphereStacking/silentcast/internal/config"
)

func TestDefaultManager_StartStop(t *testing.T) {
	t.Skip("Skipping test that requires X11 libraries")
	
	cfg := &config.HotkeyConfig{
		Prefix:          "alt+space",
		Timeout:         config.Duration(1000 * time.Millisecond),
		SequenceTimeout: config.Duration(2000 * time.Millisecond),
	}
	
	manager, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	
	// Test starting
	err = manager.Start()
	if err != nil {
		t.Errorf("Start() error = %v", err)
	}
	
	if !manager.IsRunning() {
		t.Error("Expected manager to be running after Start()")
	}
	
	// Test double start
	err = manager.Start()
	if err == nil {
		t.Error("Expected error when starting already running manager")
	}
	
	// Test stopping
	err = manager.Stop()
	if err != nil {
		t.Errorf("Stop() error = %v", err)
	}
	
	if manager.IsRunning() {
		t.Error("Expected manager to not be running after Stop()")
	}
	
	// Test double stop
	err = manager.Stop()
	if err == nil {
		t.Error("Expected error when stopping already stopped manager")
	}
}

func TestDefaultManager_Register(t *testing.T) {
	t.Skip("Skipping test that requires X11 libraries")
	
	cfg := &config.HotkeyConfig{
		Prefix:          "alt+space",
		Timeout:         config.Duration(1000 * time.Millisecond),
		SequenceTimeout: config.Duration(2000 * time.Millisecond),
	}
	
	manager, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	
	// Test registration
	err = manager.Register("ctrl+a", "select_all")
	if err != nil {
		t.Errorf("Register() error = %v", err)
	}
	
	// Test duplicate registration
	err = manager.Register("ctrl+a", "different_spell")
	if err == nil {
		t.Error("Expected error for duplicate registration")
	}
	
	// Test invalid sequence
	err = manager.Register("invalid_key", "spell")
	if err == nil {
		t.Error("Expected error for invalid key sequence")
	}
	
	// Test unregister
	err = manager.Unregister("ctrl+a")
	if err != nil {
		t.Errorf("Unregister() error = %v", err)
	}
	
	// Should be able to register again
	err = manager.Register("ctrl+a", "new_spell")
	if err != nil {
		t.Errorf("Register() after unregister error = %v", err)
	}
}

func TestDefaultManager_Handler(t *testing.T) {
	t.Skip("Skipping test that requires X11 libraries")
	
	cfg := &config.HotkeyConfig{
		Prefix:          "alt+space",
		Timeout:         config.Duration(1000 * time.Millisecond),
		SequenceTimeout: config.Duration(2000 * time.Millisecond),
	}
	
	manager, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	
	// Test with no handler - skip internal state check
	
	// Set handler
	var mu sync.Mutex
	handler := HandlerFunc(func(event Event) error {
		mu.Lock()
		defer mu.Unlock()
		return nil
	})
	
	manager.SetHandler(handler)
}

func TestMockManager(t *testing.T) {
	mock := NewMockManager()
	
	// Test handler
	events := []Event{}
	var mu sync.Mutex
	
	handler := HandlerFunc(func(event Event) error {
		mu.Lock()
		defer mu.Unlock()
		events = append(events, event)
		return nil
	})
	
	mock.SetHandler(handler)
	
	// Register sequences
	mock.Register("ctrl+a", "select_all")
	mock.Register("g,s", "git_status")
	
	// Simulate key presses
	mock.SimulateKeyPress("ctrl+a")
	mock.SimulateKeyPress("g,s")
	mock.SimulateKeyPress("unknown") // Should not trigger
	
	// Wait a bit for handlers
	time.Sleep(10 * time.Millisecond)
	
	// Check events
	mu.Lock()
	defer mu.Unlock()
	
	if len(events) != 2 {
		t.Errorf("Expected 2 events, got %d", len(events))
	}
	
	if len(events) > 0 && events[0].SpellName != "select_all" {
		t.Errorf("Expected first spell to be 'select_all', got '%s'", events[0].SpellName)
	}
	
	if len(events) > 1 && events[1].SpellName != "git_status" {
		t.Errorf("Expected second spell to be 'git_status', got '%s'", events[1].SpellName)
	}
}

func TestNewManager_InvalidPrefix(t *testing.T) {
	t.Skip("Skipping test that requires X11 libraries")
	
	cfg := &config.HotkeyConfig{
		Prefix:          "invalid_key",
		Timeout:         config.Duration(1000 * time.Millisecond),
		SequenceTimeout: config.Duration(2000 * time.Millisecond),
	}
	
	_, err := NewManager(cfg)
	if err == nil {
		t.Error("Expected error for invalid prefix key")
	}
}