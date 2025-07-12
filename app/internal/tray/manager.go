//go:build !notray

package tray

import (
	"context"

	"github.com/getlantern/systray"

	"github.com/SphereStacking/silentcast/internal/config"
	"github.com/SphereStacking/silentcast/pkg/logger"
)

// Manager manages the system tray integration
type Manager struct {
	title     string
	tooltip   string
	onReady   func()
	onExit    func()
	menuItems []*MenuItem
	quitItem  *systray.MenuItem
}

// MenuItem represents a tray menu item
type MenuItem struct {
	Title   string
	Tooltip string
	Handler func()
	item    *systray.MenuItem
}

// Config represents tray configuration
type Config struct {
	Title   string
	Tooltip string
}

// NewManager creates a new tray manager
func NewManager(ctx context.Context, cfg *config.Config) (*Manager, error) {
	trayConfig := Config{
		Title:   config.AppDisplayName,
		Tooltip: config.AppDescription,
	}
	return &Manager{
		title:     trayConfig.Title,
		tooltip:   trayConfig.Tooltip,
		menuItems: make([]*MenuItem, 0),
	}, nil
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
	// Add a nil item to represent separator
	m.menuItems = append(m.menuItems, nil)
}

// Start starts the system tray
func (m *Manager) Start() {
	m.onReady = m.setupMenu
	m.onExit = func() {
		logger.Info("System tray exiting")
	}

	// Run systray in the main thread
	systray.Run(m.onReady, m.onExit)
}

// Stop stops the system tray
func (m *Manager) Stop() {
	systray.Quit()
}

// setupMenu sets up the tray menu
func (m *Manager) setupMenu() {
	// Set tray icon and tooltip
	systray.SetIcon(getIcon())
	systray.SetTitle(m.title)
	systray.SetTooltip(m.tooltip)

	// Add menu items
	for _, item := range m.menuItems {
		if item == nil {
			systray.AddSeparator()
			continue
		}

		menuItem := systray.AddMenuItem(item.Title, item.Tooltip)
		item.item = menuItem

		// Handle click in goroutine
		go func(mi *MenuItem) {
			for range mi.item.ClickedCh {
				if mi.Handler != nil {
					mi.Handler()
				}
			}
		}(item)
	}

	// Add separator before quit
	systray.AddSeparator()

	// Add quit item
	m.quitItem = systray.AddMenuItem("Quit "+config.AppDisplayName, "Exit the application")
	go func() {
		<-m.quitItem.ClickedCh
		logger.Info("Quit requested from system tray")
		systray.Quit()
	}()
}

// UpdateTooltip updates the tray tooltip
func (m *Manager) UpdateTooltip(tooltip string) {
	m.tooltip = tooltip
	systray.SetTooltip(tooltip)
}

// ShowNotification shows a system notification
func (m *Manager) ShowNotification(title, message string) {
	// Note: systray doesn't support notifications directly
	// This would need platform-specific implementation
	logger.Info("Tray notification: %s - %s", title, message)
}

// getIcon returns the tray icon data
func getIcon() []byte {
	// This is a minimal 16x16 PNG icon (magic wand emoji style)
	// In production, embed a proper icon file
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
