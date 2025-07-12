package tray

import (
	"context"
	"testing"

	"github.com/SphereStacking/silentcast/internal/config"
)

func TestNewManager(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}

	manager, err := NewManager(ctx, cfg)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	if manager == nil {
		t.Fatal("Expected manager to be created")
	}

	if manager.title != config.AppDisplayName {
		t.Errorf("Expected title %s, got %s", config.AppDisplayName, manager.title)
	}

	if manager.tooltip != config.AppDescription {
		t.Errorf("Expected tooltip %s, got %s", config.AppDescription, manager.tooltip)
	}

	if len(manager.menuItems) != 0 {
		t.Errorf("Expected no menu items, got %d", len(manager.menuItems))
	}
}

func TestAddMenuItem(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}
	manager, err := NewManager(ctx, cfg)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	called := false
	handler := func() {
		called = true
	}

	item := manager.AddMenuItem("Test Item", "Test tooltip", handler)

	if item == nil {
		t.Fatal("Expected menu item to be created")
	}

	if item.Title != "Test Item" {
		t.Errorf("Expected title 'Test Item', got %s", item.Title)
	}

	if item.Tooltip != "Test tooltip" {
		t.Errorf("Expected tooltip 'Test tooltip', got %s", item.Tooltip)
	}

	if len(manager.menuItems) != 1 {
		t.Errorf("Expected 1 menu item, got %d", len(manager.menuItems))
	}

	// Test handler
	item.Handler()
	if !called {
		t.Error("Expected handler to be called")
	}
}

func TestAddSeparator(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}
	manager, err := NewManager(ctx, cfg)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	manager.AddMenuItem("Item 1", "", nil)
	manager.AddSeparator()
	manager.AddMenuItem("Item 2", "", nil)

	if len(manager.menuItems) != 3 {
		t.Errorf("Expected 3 menu items (including separator), got %d", len(manager.menuItems))
	}

	// Check that separator is nil
	if manager.menuItems[1] != nil {
		t.Error("Expected separator to be nil")
	}
}

func TestUpdateTooltip(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}
	manager, err := NewManager(ctx, cfg)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	manager.UpdateTooltip("Updated tooltip")

	if manager.tooltip != "Updated tooltip" {
		t.Errorf("Expected tooltip 'Updated tooltip', got %s", manager.tooltip)
	}
}

func TestGetIcon(t *testing.T) {
	icon := getIcon()

	if len(icon) == 0 {
		t.Error("Expected icon data, got empty")
	}

	// Check PNG header
	pngHeader := []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}
	if len(icon) < len(pngHeader) {
		t.Error("Icon data too short to be valid PNG")
	}

	for i, b := range pngHeader {
		if icon[i] != b {
			t.Error("Invalid PNG header")
			break
		}
	}
}
