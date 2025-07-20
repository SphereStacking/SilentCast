//go:build integration

package terminal

import (
	"context"
	"os/exec"
	"testing"
	"time"
)

func TestTerminalManager_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	manager := NewManager()

	// Test detecting available terminals
	terminals := manager.GetAvailableTerminals()
	if len(terminals) == 0 {
		t.Skip("No terminals available on this system")
	}

	t.Logf("Found %d terminals:", len(terminals))
	for _, term := range terminals {
		t.Logf("  - %s (%s) Priority: %d Default: %v",
			term.Name, term.Command, term.Priority, term.IsDefault)
	}

	// Test getting default terminal
	defaultTerm, err := manager.GetDefaultTerminal()
	if err != nil {
		t.Fatalf("Failed to get default terminal: %v", err)
	}
	t.Logf("Default terminal: %s", defaultTerm.Name)

	// Test executing a simple command
	cmd := exec.Command("echo", "Hello from terminal!")
	options := Options{
		KeepOpen: false,
		Title:    "Test Terminal",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = manager.ExecuteInTerminal(ctx, cmd, options)
	if err != nil {
		t.Logf("ExecuteInTerminal failed (may be expected in CI): %v", err)
	}
}

func TestTerminalManager_KeepOpen(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	manager := NewManager()
	terminals := manager.GetAvailableTerminals()
	if len(terminals) == 0 {
		t.Skip("No terminals available on this system")
	}

	// Test with keep_open
	cmd := exec.Command("echo", "This terminal should stay open")
	options := Options{
		KeepOpen: true,
		Title:    "Keep Open Test",
	}

	ctx := context.Background()
	err := manager.ExecuteInTerminal(ctx, cmd, options)
	if err != nil {
		t.Logf("ExecuteInTerminal with keep_open failed (may be expected in CI): %v", err)
	}
}
