package config

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/SphereStacking/silentcast/pkg/logger"
)

// Watcher watches configuration files for changes
type Watcher struct {
	loader      *Loader
	watcher     *fsnotify.Watcher
	onChange    func(*Config)
	debounce    time.Duration
	configPaths []string
}

// WatcherConfig holds watcher configuration
type WatcherConfig struct {
	ConfigPath string
	OnChange   func(*Config)
	Debounce   time.Duration
}

// NewWatcher creates a new configuration watcher
func NewWatcher(cfg WatcherConfig) (*Watcher, error) {
	// Create fsnotify watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create file watcher: %w", err)
	}

	// Create loader
	loader := NewLoader(cfg.ConfigPath)

	// Set default debounce if not specified
	debounce := cfg.Debounce
	if debounce == 0 {
		debounce = 500 * time.Millisecond
	}

	// Get platform-specific config file
	platform := GetPlatformResolver()
	platformConfigFile := platform.GetPlatformConfigFile()

	w := &Watcher{
		loader:   loader,
		watcher:  watcher,
		onChange: cfg.OnChange,
		debounce: debounce,
		configPaths: []string{
			filepath.Join(cfg.ConfigPath, ConfigName+".yml"),
			filepath.Join(cfg.ConfigPath, platformConfigFile),
		},
	}

	// Watch all potential config files
	for _, path := range w.configPaths {
		// Add the file to watcher even if it doesn't exist
		// fsnotify will watch for creation
		if err := watcher.Add(path); err != nil {
			logger.Warn("Failed to watch %s: %v", path, err)
		} else {
			logger.Info("Watching config file: %s", path)
		}
	}

	// Also watch the directory for new files
	if err := watcher.Add(cfg.ConfigPath); err != nil {
		logger.Warn("Failed to watch config directory %s: %v", cfg.ConfigPath, err)
	}

	return w, nil
}

// Start starts watching for configuration changes
func (w *Watcher) Start(ctx context.Context) {
	// Debounce timer
	var timer *time.Timer

	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Info("Config watcher stopping...")
				return

			case event, ok := <-w.watcher.Events:
				if !ok {
					return
				}

				// Check if this is one of our config files
				isConfigFile := false
				for _, path := range w.configPaths {
					if event.Name == path {
						isConfigFile = true
						break
					}
				}

				if !isConfigFile {
					continue
				}

				// Handle different event types
				switch {
				case event.Op&fsnotify.Write == fsnotify.Write:
					logger.Info("Config file modified: %s", event.Name)
				case event.Op&fsnotify.Create == fsnotify.Create:
					logger.Info("Config file created: %s", event.Name)
					// Re-add to watcher in case it's a new file
					if err := w.watcher.Add(event.Name); err != nil {
						logger.Warn("Failed to add %s to watcher: %v", event.Name, err)
					}
				case event.Op&fsnotify.Remove == fsnotify.Remove:
					logger.Info("Config file removed: %s", event.Name)
					continue
				case event.Op&fsnotify.Rename == fsnotify.Rename:
					logger.Info("Config file renamed: %s", event.Name)
					continue
				default:
					continue
				}

				// Debounce - reset timer on each event
				if timer != nil {
					timer.Stop()
				}

				timer = time.AfterFunc(w.debounce, func() {
					w.reload()
				})

			case err, ok := <-w.watcher.Errors:
				if !ok {
					return
				}
				logger.Error("Config watcher error: %v", err)
			}
		}
	}()
}

// Stop stops the watcher
func (w *Watcher) Stop() error {
	return w.watcher.Close()
}

// reload reloads the configuration and calls the callback
func (w *Watcher) reload() {
	logger.Info("Reloading configuration...")

	cfg, err := w.loader.Load()
	if err != nil {
		logger.Error("Failed to reload configuration: %v", err)
		return
	}

	// Call the callback with the new configuration
	if w.onChange != nil {
		w.onChange(cfg)
	}

	logger.Info("Configuration reloaded successfully")
}
