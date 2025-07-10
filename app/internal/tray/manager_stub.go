//go:build notray

package tray

import (
	"github.com/SphereStacking/silentcast/pkg/logger"
)

// Manager manages the system tray integration (stub version)
type Manager struct {
	title       string
	tooltip     string
	menuItems   []*MenuItem
}

// MenuItem represents a tray menu item
type MenuItem struct {
	Title    string
	Tooltip  string
	Handler  func()
}

// Config represents tray configuration
type Config struct {
	Title   string
	Tooltip string
}

// NewManager creates a new tray manager
func NewManager(cfg Config) *Manager {
	logger.Info("System tray disabled (notray build tag)")
	return &Manager{
		title:     cfg.Title,
		tooltip:   cfg.Tooltip,
		menuItems: make([]*MenuItem, 0),
	}
}

// AddMenuItem adds a menu item to the tray
func (m *Manager) AddMenuItem(title, tooltip string, handler func()) *MenuItem {
	item := &MenuItem{
		Title:   title,
		Tooltip: tooltip,
		Handler: handler,
	}
	m.menuItems = append(m.menuItems, item)
	return item
}

// AddSeparator adds a separator to the menu
func (m *Manager) AddSeparator() {
	m.menuItems = append(m.menuItems, nil)
}

// Start starts the system tray (no-op in stub)
func (m *Manager) Start() {
	logger.Info("System tray Start() called (no-op in stub)")
}

// Stop stops the system tray (no-op in stub)
func (m *Manager) Stop() {
	logger.Info("System tray Stop() called (no-op in stub)")
}

// UpdateTooltip updates the tray tooltip
func (m *Manager) UpdateTooltip(tooltip string) {
	m.tooltip = tooltip
}

// ShowNotification shows a system notification
func (m *Manager) ShowNotification(title, message string) {
	logger.Info("Tray notification: %s - %s", title, message)
}

// getIcon returns the tray icon data
func getIcon() []byte {
	// Minimal icon for stub
	return []byte{
		0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d,
		0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x10,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0xf3, 0xff, 0x61, 0x00, 0x00, 0x00,
		0x19, 0x74, 0x45, 0x58, 0x74, 0x53, 0x6f, 0x66, 0x74, 0x77, 0x61, 0x72,
		0x65, 0x00, 0x41, 0x64, 0x6f, 0x62, 0x65, 0x20, 0x49, 0x6d, 0x61, 0x67,
		0x65, 0x52, 0x65, 0x61, 0x64, 0x79, 0x71, 0xc9, 0x65, 0x3c, 0x00, 0x00,
		0x01, 0x4e, 0x49, 0x44, 0x41, 0x54, 0x78, 0xda, 0x62, 0xfc, 0xff, 0xff,
		0x3f, 0x03, 0x3e, 0x00, 0x08, 0x00, 0x01, 0x00, 0x00, 0x05, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae,
		0x42, 0x60, 0x82,
	}
}