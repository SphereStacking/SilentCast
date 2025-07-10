package tray

import (
	"testing"
)

func TestNewManager(t *testing.T) {
	cfg := Config{
		Title:   "Test",
		Tooltip: "Test tooltip",
	}
	
	manager := NewManager(cfg)
	
	if manager == nil {
		t.Fatal("Expected manager to be created")
	}
	
	if manager.title != cfg.Title {
		t.Errorf("Expected title %s, got %s", cfg.Title, manager.title)
	}
	
	if manager.tooltip != cfg.Tooltip {
		t.Errorf("Expected tooltip %s, got %s", cfg.Tooltip, manager.tooltip)
	}
	
	if len(manager.menuItems) != 0 {
		t.Errorf("Expected no menu items, got %d", len(manager.menuItems))
	}
}

func TestAddMenuItem(t *testing.T) {
	manager := NewManager(Config{})
	
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
	manager := NewManager(Config{})
	
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
	manager := NewManager(Config{
		Tooltip: "Initial tooltip",
	})
	
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