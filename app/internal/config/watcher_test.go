package config

import (
	"context"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewWatcher(t *testing.T) {
	tempDir := t.TempDir()

	cfg := WatcherConfig{
		ConfigPath: tempDir,
		OnChange:   func(*Config) {},
		Debounce:   100 * time.Millisecond,
	}

	watcher, err := NewWatcher(cfg)
	if err != nil {
		t.Fatalf("Failed to create watcher: %v", err)
	}
	defer watcher.Stop()

	if watcher.loader == nil {
		t.Error("Expected loader to be created")
	}

	if watcher.watcher == nil {
		t.Error("Expected fsnotify watcher to be created")
	}

	if watcher.debounce != cfg.Debounce {
		t.Errorf("Expected debounce %v, got %v", cfg.Debounce, watcher.debounce)
	}
}

func TestWatcherDefaultDebounce(t *testing.T) {
	tempDir := t.TempDir()

	cfg := WatcherConfig{
		ConfigPath: tempDir,
		OnChange:   func(*Config) {},
		// No debounce specified
	}

	watcher, err := NewWatcher(cfg)
	if err != nil {
		t.Fatalf("Failed to create watcher: %v", err)
	}
	defer watcher.Stop()

	expectedDebounce := 500 * time.Millisecond
	if watcher.debounce != expectedDebounce {
		t.Errorf("Expected default debounce %v, got %v", expectedDebounce, watcher.debounce)
	}
}

func TestWatcherFileChange(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, ConfigName+".yml")

	// Create initial config
	initialConfig := []byte(`
daemon:
  log_level: info
hotkeys:
  prefix: "alt+space"
spells:
  e: "editor"
grimoire:
  editor:
    type: app
    command: nano
`)

	if err := os.WriteFile(configFile, initialConfig, 0o600); err != nil {
		t.Fatalf("Failed to write initial config: %v", err)
	}

	// Create a channel to receive config changes
	configChanged := make(chan *Config, 1)

	cfg := WatcherConfig{
		ConfigPath: tempDir,
		OnChange: func(c *Config) {
			configChanged <- c
		},
		Debounce: 100 * time.Millisecond,
	}

	watcher, err := NewWatcher(cfg)
	if err != nil {
		t.Fatalf("Failed to create watcher: %v", err)
	}
	defer watcher.Stop()

	// Start watcher
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	watcher.Start(ctx)

	// Give watcher time to start
	time.Sleep(200 * time.Millisecond)

	// Modify config
	updatedConfig := []byte(`
daemon:
  log_level: debug  # Changed
hotkeys:
  prefix: "alt+space"
spells:
  e: "editor"
  t: "terminal"  # Added
grimoire:
  editor:
    type: app
    command: nano
  terminal:
    type: app
    command: bash
`)

	if err := os.WriteFile(configFile, updatedConfig, 0o600); err != nil {
		t.Fatalf("Failed to write updated config: %v", err)
	}

	// Wait for config change event
	select {
	case newConfig := <-configChanged:
		// Verify the config was updated
		if newConfig.Daemon.LogLevel != "debug" {
			t.Errorf("Expected log level 'debug', got '%s'", newConfig.Daemon.LogLevel)
		}

		if _, exists := newConfig.Shortcuts["t"]; !exists {
			t.Error("Expected new spell 't' to exist")
		}

		if _, exists := newConfig.Actions["terminal"]; !exists {
			t.Error("Expected new grimoire entry 'terminal' to exist")
		}

	case <-time.After(1 * time.Second):
		t.Error("Timeout waiting for config change")
	}
}

func TestWatcherDebounce(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, ConfigName+".yml")

	// Create initial config
	initialConfig := []byte(`
daemon:
  log_level: info
hotkeys:
  prefix: "alt+space"
spells: {}
grimoire: {}
`)

	if err := os.WriteFile(configFile, initialConfig, 0o600); err != nil {
		t.Fatalf("Failed to write initial config: %v", err)
	}

	// Count config changes
	var changeCount int32

	cfg := WatcherConfig{
		ConfigPath: tempDir,
		OnChange: func(c *Config) {
			atomic.AddInt32(&changeCount, 1)
		},
		Debounce: 200 * time.Millisecond,
	}

	watcher, err := NewWatcher(cfg)
	if err != nil {
		t.Fatalf("Failed to create watcher: %v", err)
	}
	defer watcher.Stop()

	// Start watcher
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	watcher.Start(ctx)

	// Give watcher time to start
	time.Sleep(100 * time.Millisecond)

	// Make multiple rapid changes
	for i := 0; i < 5; i++ {
		config := []byte(`
daemon:
  log_level: info
hotkeys:
  prefix: "alt+space"
spells: {}
grimoire: {}
# Change ` + string(rune('0'+i)) + `
`)
		if err := os.WriteFile(configFile, config, 0o600); err != nil {
			t.Fatalf("Failed to write config: %v", err)
		}
		time.Sleep(50 * time.Millisecond) // Less than debounce time
	}

	// Wait for debounce to complete
	time.Sleep(300 * time.Millisecond)

	// Should only have received one change due to debouncing
	if atomic.LoadInt32(&changeCount) != 1 {
		t.Errorf("Expected 1 config change due to debouncing, got %d", atomic.LoadInt32(&changeCount))
	}
}
