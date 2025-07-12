package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/SphereStacking/silentcast/internal/action"
	"github.com/SphereStacking/silentcast/internal/config"
	"github.com/SphereStacking/silentcast/internal/errors"
	"github.com/SphereStacking/silentcast/internal/hotkey"
	"github.com/SphereStacking/silentcast/internal/notify"
	"github.com/SphereStacking/silentcast/internal/permission"
	"github.com/SphereStacking/silentcast/internal/tray"
	"github.com/SphereStacking/silentcast/pkg/logger"
)

var Version = "0.1.0-dev"

func main() {
	// Parse command line flags
	var (
		noTray  bool
		version bool
	)
	flag.BoolVar(&noTray, "no-tray", false, "Disable system tray integration")
	flag.BoolVar(&version, "version", false, "Print version and exit")
	flag.Parse()

	if version {
		fmt.Printf("SilentCast v%s\n", Version)
		os.Exit(0)
	}

	if err := run(noTray); err != nil {
		// Print user-friendly error message
		fmt.Fprintf(os.Stderr, "❌ %s\n", errors.GetUserMessage(err))

		// Log detailed error for debugging
		log.Printf("Error: %+v", err)
		os.Exit(1)
	}
}

func run(noTray bool) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// WaitGroup to coordinate shutdown
	var wg sync.WaitGroup

	// Print banner
	fmt.Printf("🪄 %s - %s v%s\n", config.AppDisplayName, config.AppDescription, Version)
	fmt.Println("Press Ctrl+C to exit")
	fmt.Println()

	// Load configuration first for logger settings
	configPath := getConfigPath()
	loader := config.NewLoader(configPath)

	cfg, err := loader.Load()
	if err != nil {
		return errors.Wrap(errors.ErrorTypeConfig, "failed to load configuration", err)
	}

	// Initialize logger
	loggerConfig := logger.Config{
		Level:      cfg.Logger.Level,
		File:       cfg.Logger.File,
		MaxSize:    cfg.Logger.MaxSize,
		MaxBackups: cfg.Logger.MaxBackups,
		MaxAge:     cfg.Logger.MaxAge,
		Compress:   cfg.Logger.Compress,
		Console:    true, // Always enable console output
	}
	if initErr := logger.Initialize(loggerConfig); initErr != nil {
		return errors.Wrap(errors.ErrorTypeSystem, "failed to initialize logger", initErr)
	}

	logger.Info("%s starting up v%s", config.AppDisplayName, Version)
	logger.Info("Configuration loaded from %s", configPath)

	// Initialize notification manager
	notifier := notify.NewManager()

	// Check permissions
	logger.Info("Checking permissions...")
	permManager, err := permission.NewManager()
	if err != nil {
		return errors.Wrap(errors.ErrorTypeSystem, "failed to create permission manager", err)
	}

	permissions, err := permManager.Check(ctx)
	if err != nil {
		return errors.Wrap(errors.ErrorTypePermission, "failed to check permissions", err)
	}

	// Display permission status
	for _, perm := range permissions {
		if !perm.Required || perm.Status == permission.StatusGranted {
			continue
		}

		logger.Warn("Permission required: %s - %s", perm.Type, perm.Description)
		if notifyErr := notifier.Warning(ctx, "Permission Required",
			fmt.Sprintf("%s: %s", perm.Type, perm.Description)); notifyErr != nil {
			logger.Error("Failed to send warning notification: %v", notifyErr)
		}

		instructions := permManager.GetInstructions(perm.Type)
		fmt.Println(instructions)
		fmt.Println()
	}

	// Initialize action manager
	logger.Info("Initializing action manager...")
	actionManager := action.NewManager(cfg.Actions)

	// Initialize hotkey manager
	logger.Info("Initializing hotkey manager...")
	hotkeyManager, err := hotkey.NewManager(&cfg.Hotkeys)
	if err != nil {
		return errors.Wrap(errors.ErrorTypeHotkey, "failed to create hotkey manager", err)
	}

	// Set up hotkey handler
	hotkeyManager.SetHandler(hotkey.HandlerFunc(func(event hotkey.Event) error {
		logger.Info("Spell cast: %s → %s", event.Sequence.String(), event.SpellName)
		if err := notifier.Info(ctx, "Spell Cast",
			fmt.Sprintf("🎯 %s → %s", event.Sequence.String(), event.SpellName)); err != nil {
			logger.Error("Failed to send info notification: %v", err)
		}

		// Execute the action
		if err := actionManager.Execute(ctx, event.SpellName); err != nil {
			logger.Error("Failed to execute spell %s: %v", event.SpellName, err)
			if notifyErr := notifier.Error(ctx, "Spell Failed", err.Error()); notifyErr != nil {
				logger.Error("Failed to send error notification: %v", notifyErr)
			}
			return err
		}

		logger.Info("Successfully executed spell: %s", event.SpellName)
		return nil
	}))

	// Register all hotkeys
	for sequence, spellName := range cfg.Shortcuts {
		if err := hotkeyManager.Register(sequence, spellName); err != nil {
			logger.Warn("Failed to register hotkey %s: %v", sequence, err)
			if notifyErr := notifier.Warning(ctx, "Registration Failed",
				fmt.Sprintf("Could not register %s: %v", sequence, err)); notifyErr != nil {
				logger.Error("Failed to send warning notification: %v", notifyErr)
			}
			continue
		}
		logger.Info("Registered hotkey: %s → %s", sequence, spellName)
		fmt.Printf("  ✨ %s → %s\n", sequence, spellName)
	}
	fmt.Println()

	// Start hotkey manager
	logger.Info("Starting hotkey manager...")
	if err := hotkeyManager.Start(); err != nil {
		return errors.Wrap(errors.ErrorTypeHotkey, "failed to start hotkey manager", err)
	}
	defer func() {
		if err := hotkeyManager.Stop(); err != nil {
			logger.Error("Failed to stop hotkey manager: %v", err)
		}
	}()

	logger.Info("%s is active! Listening with prefix: %s", config.AppDisplayName, cfg.Hotkeys.Prefix)
	if err := notifier.Success(ctx, config.AppDisplayName+" Active",
		fmt.Sprintf("Listening with prefix: %s", cfg.Hotkeys.Prefix)); err != nil {
		logger.Error("Failed to send success notification: %v", err)
	}

	// Setup shutdown handling
	shutdownCh := make(chan struct{})

	// Handle OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logger.Info("Received signal: %v", sig)
		close(shutdownCh)
	}()

	if !noTray {
		// Initialize system tray
		trayManager, err := tray.NewManager(ctx, cfg)
		if err != nil {
			return errors.Wrap(errors.ErrorTypeSystem, "failed to initialize tray manager", err)
		}

		// Add menu items
		trayManager.AddMenuItem("Show Hotkeys", "Display configured hotkeys", func() {
			logger.Info("Show hotkeys requested")
			fmt.Println("\n🗿 Configured Hotkeys:")
			for sequence, spellName := range cfg.Shortcuts {
				fmt.Printf("  ✨ %s → %s\n", sequence, spellName)
			}
		})

		trayManager.AddMenuItem("Reload Config", "Reload configuration file", func() {
			logger.Info("Config reload requested")
			if err := notifier.Info(ctx, "Config Reload", "Feature coming soon"); err != nil {
				logger.Error("Failed to send info notification: %v", err)
			}
		})

		trayManager.AddSeparator()

		trayManager.AddMenuItem("About", "About "+config.AppDisplayName, func() {
			logger.Info("About requested")
			fmt.Printf("\n🪄 %s v%s\n", config.AppDisplayName, Version)
			fmt.Println(config.AppDescription)
			fmt.Printf("https://github.com/%s/%s\n", config.AppOrg, config.AppRepo)
		})

		// Run tray in a goroutine since Start blocks
		wg.Add(1)
		go func() {
			defer wg.Done()
			trayManager.Start()
			logger.Info("System tray stopped")
			close(shutdownCh) // Trigger shutdown if tray exits
		}()

		// Give tray time to initialize
		// In a real implementation, we'd use a better synchronization method
		logger.Info("Starting system tray...")
	}

	// Wait for shutdown signal
	<-shutdownCh

	fmt.Println("\n👋 Shutting down " + config.AppDisplayName + "...")
	logger.Info("Shutting down %s...", config.AppDisplayName)

	// Cancel context to stop all components
	cancel()

	// Wait for all goroutines
	wg.Wait()

	return nil
}

// getConfigPath returns the configuration directory path
func getConfigPath() string {
	// Check for environment variable
	if path := os.Getenv("SILENTCAST_CONFIG"); path != "" {
		return path
	}

	// Check current directory
	if _, err := os.Stat(config.ConfigName + ".yml"); err == nil {
		return "."
	}

	// Use user config directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		// Fallback to current directory
		return "."
	}

	return filepath.Join(configDir, config.AppName)
}
