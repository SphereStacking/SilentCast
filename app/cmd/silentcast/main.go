package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/SphereStacking/silentcast/internal/action"
	"github.com/SphereStacking/silentcast/internal/config"
	"github.com/SphereStacking/silentcast/internal/errors"
	"github.com/SphereStacking/silentcast/internal/hotkey"
	"github.com/SphereStacking/silentcast/internal/notify"
	"github.com/SphereStacking/silentcast/internal/permission"
	"github.com/SphereStacking/silentcast/internal/service"
	"github.com/SphereStacking/silentcast/internal/tray"
	"github.com/SphereStacking/silentcast/internal/version"
	"github.com/SphereStacking/silentcast/pkg/logger"
)

// Helper functions inlined to ensure they're available during GoReleaser build

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

// getConfigSearchPaths returns all paths where config files are searched
func getConfigSearchPaths() []string {
	var paths []string

	// 1. Environment variable
	if envPath := os.Getenv("SILENTCAST_CONFIG"); envPath != "" {
		paths = append(paths, envPath)
	}

	// 2. Current directory
	paths = append(paths, ".")

	// 3. User config directory
	if configDir, err := os.UserConfigDir(); err == nil {
		paths = append(paths, filepath.Join(configDir, config.AppName))
	}

	// 4. System config directory (Unix-like systems)
	if runtime.GOOS != "windows" {
		paths = append(paths, "/etc/"+config.AppName)
	}

	return paths
}

// Version information is now managed by the version package

func main() {
	// Parse command line flags
	flags := ParseFlags()

	// Create main run function for service mode
	mainRun := func() error {
		return run(flags.NoTray, flags.Debug)
	}

	// Create command registry with service support
	registry := NewCommandRegistryWithService(flags, mainRun)
	if registry.ExecuteCommands() {
		os.Exit(0)
	}
	
	// Check if we should run as a service (Windows)
	svcManager := service.NewManager(mainRun)
	if err := svcManager.Run(); err != nil {
		// If not running as service, error will be returned
		// Continue with normal execution
		log.Printf("Not running as service: %v", err)
	} else {
		// Running as service, exit when service stops
		os.Exit(0)
	}

	// Check for once mode
	if flags.Once {
		if err := runOnce(flags.SpellName, flags.Debug); err != nil {
			// Print user-friendly error message
			fmt.Fprintf(os.Stderr, "❌ %s\n", errors.GetUserMessage(err))
			// Log detailed error for debugging
			log.Printf("Error: %+v", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Check for test-spell mode
	if flags.TestSpell {
		if err := testSpell(flags.SpellName, flags.Debug); err != nil {
			// Print user-friendly error message
			fmt.Fprintf(os.Stderr, "❌ %s\n", errors.GetUserMessage(err))
			// Log detailed error for debugging
			log.Printf("Error: %+v", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Check for dry-run mode
	if flags.DryRun {
		if err := dryRun(flags.SpellName, flags.Debug); err != nil {
			// Print user-friendly error message
			fmt.Fprintf(os.Stderr, "❌ %s\n", errors.GetUserMessage(err))
			// Log detailed error for debugging
			log.Printf("Error: %+v", err)
			os.Exit(1)
		}
		os.Exit(0)
	}


	// No command specified, run the main application
	if err := run(flags.NoTray, flags.Debug); err != nil {
		// Print user-friendly error message
		fmt.Fprintf(os.Stderr, "❌ %s\n", errors.GetUserMessage(err))

		// Log detailed error for debugging
		log.Printf("Error: %+v", err)
		os.Exit(1)
	}
}

func run(noTray bool, debug bool) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// WaitGroup to coordinate shutdown
	var wg sync.WaitGroup

	// Print banner
	fmt.Printf("🪄 %s - %s v%s\n", config.AppDisplayName, config.AppDescription, version.GetVersionString())
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
	logFile := cfg.Logger.File
	if logFile == "" {
		// Use default log file path if not specified
		logFile = filepath.Join(configPath, "silentcast.log")
	}

	// Set debug level if debug flag is enabled
	logLevel := cfg.Logger.Level
	if debug {
		logLevel = "debug"
	}

	loggerConfig := logger.Config{
		Level:      logLevel,
		File:       logFile,
		MaxSize:    cfg.Logger.MaxSize,
		MaxBackups: cfg.Logger.MaxBackups,
		MaxAge:     cfg.Logger.MaxAge,
		Compress:   cfg.Logger.Compress,
		Console:    true, // Always enable console output
	}
	if initErr := logger.Initialize(loggerConfig); initErr != nil {
		return errors.Wrap(errors.ErrorTypeSystem, "failed to initialize logger", initErr)
	}

	logger.Info("%s starting up v%s", config.AppDisplayName, version.GetVersionString())
	logger.Info("Configuration loaded from %s", configPath)
	
	if debug {
		logger.Debug("Debug logging enabled")
		logger.Debug("Logger configuration: level=%s, file=%s", logLevel, logFile)
	}

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

	// Start configuration file watcher
	logger.Info("Starting configuration file watcher...")
	watcher, err := config.NewWatcher(config.WatcherConfig{
		ConfigPath: configPath,
		OnChange: func(newCfg *config.Config) {
			logger.Info("Configuration changed, reloading...")

			// Update action manager
			actionManager.UpdateActions(newCfg.Actions)

			// Update hotkey manager if hotkeys changed
			if !hotkeyConfigEqual(&cfg.Hotkeys, &newCfg.Hotkeys) || !shortcutsEqual(cfg.Shortcuts, newCfg.Shortcuts) {
				logger.Info("Hotkeys changed, reregistering...")

				// Stop current hotkeys
				if stopErr := hotkeyManager.Stop(); stopErr != nil {
					logger.Error("Failed to stop hotkey manager: %v", stopErr)
				}

				// Create new hotkey manager
				newHotkeyManager, newHotkeyErr := hotkey.NewManager(&newCfg.Hotkeys)
				if newHotkeyErr != nil {
					logger.Error("Failed to create new hotkey manager: %v", newHotkeyErr)
					if notifyErr := notifier.Error(ctx, "Hotkey Reload Failed", 
						"Failed to reload hotkey configuration"); notifyErr != nil {
						logger.Error("Failed to send error notification: %v", notifyErr)
					}
					return
				}

				// Set up handler with same logic
				newHotkeyManager.SetHandler(hotkey.HandlerFunc(func(event hotkey.Event) error {
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

				// Register all new hotkeys
				for sequence, spellName := range newCfg.Shortcuts {
					if err := newHotkeyManager.Register(sequence, spellName); err != nil {
						logger.Warn("Failed to register hotkey %s: %v", sequence, err)
					}
				}

				// Start new hotkey manager
				if err := newHotkeyManager.Start(); err != nil {
					logger.Error("Failed to start new hotkey manager: %v", err)
					if notifyErr := notifier.Error(ctx, "Hotkey Reload Failed", 
						"Failed to start new hotkey manager"); notifyErr != nil {
						logger.Error("Failed to send error notification: %v", notifyErr)
					}
					return
				}

				// Replace the hotkey manager
				hotkeyManager = newHotkeyManager
				
				logger.Info("Hotkeys reloaded successfully")
			}

			// Update the configuration reference
			cfg = newCfg

			// Notify user of successful reload
			if err := notifier.Success(ctx, "Configuration Reloaded", 
				"SilentCast configuration has been updated"); err != nil {
				logger.Error("Failed to send success notification: %v", err)
			}

			logger.Info("Configuration reload completed successfully")
		},
		Debounce: 500 * time.Millisecond,
	})
	if err != nil {
		logger.Warn("Failed to create config watcher: %v", err)
	} else {
		watcher.Start(ctx)
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ctx.Done()
			logger.Info("Stopping config watcher...")
			if err := watcher.Stop(); err != nil {
				logger.Error("Failed to stop config watcher: %v", err)
			}
		}()
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

	var trayManager *tray.Manager
	if !noTray {
		// Initialize system tray
		var err error
		trayManager, err = tray.NewManager(ctx, cfg)
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
			logger.Info("Manual config reload requested")
			newCfg, err := loader.Load()
			if err != nil {
				logger.Error("Failed to reload configuration: %v", err)
				if notifyErr := notifier.Error(ctx, "Config Reload Failed", err.Error()); notifyErr != nil {
					logger.Error("Failed to send error notification: %v", notifyErr)
				}
				return
			}

			// Manual reload uses the same logic as the watcher
			actionManager.UpdateActions(newCfg.Actions)

			if !hotkeyConfigEqual(&cfg.Hotkeys, &newCfg.Hotkeys) || !shortcutsEqual(cfg.Shortcuts, newCfg.Shortcuts) {
				logger.Info("Hotkeys changed, reregistering...")
				// Same hotkey reload logic as in watcher...
			}

			cfg = newCfg
			if err := notifier.Success(ctx, "Configuration Reloaded", "Manual reload successful"); err != nil {
				logger.Error("Failed to send success notification: %v", err)
			}
		})

		trayManager.AddSeparator()

		trayManager.AddMenuItem("About", "About "+config.AppDisplayName, func() {
			logger.Info("About requested")
			fmt.Printf("\n🪄 %s v%s\n", config.AppDisplayName, version.GetVersionString())
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

	// Stop tray if it's running
	if trayManager != nil {
		logger.Info("Stopping system tray...")
		trayManager.Stop()
	}

	// Wait for all goroutines with timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Normal shutdown
	case <-time.After(3 * time.Second):
		// Force exit after timeout
		logger.Warn("Shutdown timeout exceeded, forcing exit")
	}

	return nil
}

// runOnce executes a single spell and exits
func runOnce(spellName string, debug bool) error {
	if spellName == "" {
		return fmt.Errorf("spell name is required when using --once flag")
	}

	// Load configuration
	configPath := getConfigPath()
	loader := config.NewLoader(configPath)

	cfg, err := loader.Load()
	if err != nil {
		return errors.Wrap(errors.ErrorTypeConfig, "failed to load configuration", err)
	}

	// Check if spell exists
	actionName, exists := cfg.Shortcuts[spellName]
	if !exists {
		return fmt.Errorf("spell '%s' not found. Available spells: %v", spellName, getSpellList(cfg.Shortcuts))
	}

	// Print what we're about to execute
	if action, actionExists := cfg.Actions[actionName]; actionExists {
		fmt.Printf("🪄 Executing spell: %s → %s\n", spellName, actionName)
		fmt.Printf("   Type: %s\n", action.Type)
		fmt.Printf("   Command: %s\n", action.Command)
		if action.Description != "" {
			fmt.Printf("   Description: %s\n", action.Description)
		}
		fmt.Println()
	}

	// Initialize minimal logger for once mode (console only)
	logLevel := "info"
	if debug {
		logLevel = "debug"
	}
	loggerConfig := logger.Config{
		Level:   logLevel,
		Console: true,
	}
	if err := logger.Initialize(loggerConfig); err != nil {
		return errors.Wrap(errors.ErrorTypeSystem, "failed to initialize logger", err)
	}

	// Initialize action manager
	actionManager := action.NewManager(cfg.Actions)

	// Execute the action
	ctx := context.Background()
	if err := actionManager.Execute(ctx, actionName); err != nil {
		return errors.Wrap(errors.ErrorTypeExecution, fmt.Sprintf("failed to execute spell '%s'", spellName), err)
	}

	fmt.Printf("✅ Successfully executed spell: %s\n", spellName)
	return nil
}

// getSpellList returns a list of available spell names
func getSpellList(shortcuts map[string]string) []string {
	var spells []string
	for spell := range shortcuts {
		spells = append(spells, spell)
	}
	return spells
}

// testSpell tests a spell with detailed debug information
func testSpell(spellName string, debug bool) error {
	if spellName == "" {
		return fmt.Errorf("spell name is required when using --test-spell flag")
	}

	fmt.Printf("🧪 Testing spell: %s\n", spellName)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println()

	// Load configuration
	configPath := getConfigPath()
	fmt.Printf("📁 Loading configuration from: %s\n", configPath)
	
	loader := config.NewLoader(configPath)
	cfg, err := loader.Load()
	if err != nil {
		return errors.Wrap(errors.ErrorTypeConfig, "failed to load configuration", err)
	}
	fmt.Println("✅ Configuration loaded successfully")
	fmt.Println()

	// Check if spell exists
	fmt.Printf("🔍 Looking up spell '%s'...\n", spellName)
	actionName, exists := cfg.Shortcuts[spellName]
	if !exists {
		fmt.Printf("❌ Spell '%s' not found\n", spellName)
		fmt.Printf("Available spells: %v\n", getSpellList(cfg.Shortcuts))
		return fmt.Errorf("spell '%s' not found", spellName)
	}
	fmt.Printf("✅ Spell found: %s → %s\n", spellName, actionName)
	fmt.Println()

	// Check if action exists
	fmt.Printf("🎯 Looking up action '%s'...\n", actionName)
	action, actionExists := cfg.Actions[actionName]
	if !actionExists {
		fmt.Printf("❌ Action '%s' not found in grimoire\n", actionName)
		return fmt.Errorf("action '%s' not found in grimoire", actionName)
	}
	fmt.Printf("✅ Action found in grimoire\n")
	fmt.Println()

	// Display detailed action information
	fmt.Println("📋 Action Details:")
	fmt.Printf("   Type: %s\n", action.Type)
	fmt.Printf("   Command: %s\n", action.Command)
	if action.Description != "" {
		fmt.Printf("   Description: %s\n", action.Description)
	}
	if len(action.Args) > 0 {
		fmt.Printf("   Arguments: %v\n", action.Args)
	}
	if action.WorkingDir != "" {
		fmt.Printf("   Working Directory: %s\n", action.WorkingDir)
		// Expand environment variables for display
		expandedDir := os.ExpandEnv(action.WorkingDir)
		if expandedDir != action.WorkingDir {
			fmt.Printf("   Expanded Working Directory: %s\n", expandedDir)
		}
	}
	if action.Shell != "" {
		fmt.Printf("   Shell: %s\n", action.Shell)
	}
	if action.Timeout > 0 {
		fmt.Printf("   Timeout: %d seconds\n", action.Timeout)
	}
	fmt.Println()

	// Display action options
	fmt.Println("⚙️  Action Options:")
	var options []string
	if action.ShowOutput {
		options = append(options, "show_output")
	}
	if action.KeepOpen {
		options = append(options, "keep_open")
	}
	if action.Terminal {
		options = append(options, "terminal")
	}
	if action.ForceTerminal {
		options = append(options, "force_terminal")
	}
	if action.Admin {
		options = append(options, "admin")
	}
	if len(options) > 0 {
		fmt.Printf("   Enabled: %s\n", strings.Join(options, ", "))
	} else {
		fmt.Println("   No special options enabled")
	}
	fmt.Println()

	// Display environment variables
	if len(action.Env) > 0 {
		fmt.Println("🌍 Environment Variables:")
		for key, value := range action.Env {
			fmt.Printf("   %s=%s\n", key, value)
			// Show expanded value if different
			expandedValue := os.ExpandEnv(value)
			if expandedValue != value {
				fmt.Printf("   %s=%s (expanded)\n", key, expandedValue)
			}
		}
		fmt.Println()
	}

	// Type-specific validation and information
	fmt.Printf("🔬 Type-Specific Analysis (%s):\n", action.Type)
	switch action.Type {
	case "app":
		return testAppAction(action)
	case "script":
		return testScriptAction(action)
	case "url":
		return testURLAction(action)
	default:
		fmt.Printf("❌ Unknown action type: %s\n", action.Type)
		return fmt.Errorf("unknown action type: %s", action.Type)
	}
}

// testAppAction tests app-specific aspects
func testAppAction(action config.ActionConfig) error {
	expandedCmd := os.ExpandEnv(action.Command)
	fmt.Printf("   Command: %s\n", action.Command)
	if expandedCmd != action.Command {
		fmt.Printf("   Expanded Command: %s\n", expandedCmd)
	}

	// Check if it's an absolute path
	if filepath.IsAbs(expandedCmd) {
		fmt.Printf("   Type: Absolute path\n")
		if _, err := os.Stat(expandedCmd); os.IsNotExist(err) {
			fmt.Printf("   ❌ File does not exist: %s\n", expandedCmd)
			return fmt.Errorf("application file does not exist: %s", expandedCmd)
		} else {
			fmt.Printf("   ✅ File exists: %s\n", expandedCmd)
		}
	} else {
		fmt.Printf("   Type: Command in PATH\n")
		if fullPath, err := exec.LookPath(expandedCmd); err != nil {
			fmt.Printf("   ❌ Command not found in PATH: %s\n", expandedCmd)
			return fmt.Errorf("application not found in PATH: %s", expandedCmd)
		} else {
			fmt.Printf("   ✅ Found in PATH: %s\n", fullPath)
		}
	}

	fmt.Println("\n✅ App action validation completed successfully")
	return nil
}

// testScriptAction tests script-specific aspects
func testScriptAction(action config.ActionConfig) error {
	fmt.Printf("   Script Command: %s\n", action.Command)
	
	// Show expanded command
	expandedCmd := os.ExpandEnv(action.Command)
	if expandedCmd != action.Command {
		fmt.Printf("   Expanded Command: %s\n", expandedCmd)
	}

	// Check shell
	if action.Shell != "" {
		fmt.Printf("   Using Shell: %s\n", action.Shell)
		if !filepath.IsAbs(action.Shell) {
			if fullPath, err := exec.LookPath(action.Shell); err != nil {
				fmt.Printf("   ❌ Shell not found: %s\n", action.Shell)
				return fmt.Errorf("shell not found: %s", action.Shell)
			} else {
				fmt.Printf("   ✅ Shell found: %s\n", fullPath)
			}
		} else {
			if _, err := os.Stat(action.Shell); os.IsNotExist(err) {
				fmt.Printf("   ❌ Shell does not exist: %s\n", action.Shell)
				return fmt.Errorf("shell does not exist: %s", action.Shell)
			} else {
				fmt.Printf("   ✅ Shell exists: %s\n", action.Shell)
			}
		}
	} else {
		fmt.Printf("   Using Default Shell: %s\n", getDefaultShell())
	}

	// Check working directory
	if action.WorkingDir != "" {
		expandedDir := os.ExpandEnv(action.WorkingDir)
		if _, err := os.Stat(expandedDir); os.IsNotExist(err) {
			fmt.Printf("   ❌ Working directory does not exist: %s\n", expandedDir)
			return fmt.Errorf("working directory does not exist: %s", expandedDir)
		} else {
			fmt.Printf("   ✅ Working directory exists: %s\n", expandedDir)
		}
	}

	fmt.Println("\n✅ Script action validation completed successfully")
	return nil
}

// testURLAction tests URL-specific aspects
func testURLAction(action config.ActionConfig) error {
	urlStr := strings.TrimSpace(action.Command)
	fmt.Printf("   URL: %s\n", urlStr)

	// Add scheme if missing
	if !strings.Contains(urlStr, "://") {
		urlStr = "https://" + urlStr
		fmt.Printf("   URL with scheme: %s\n", urlStr)
	}

	// Parse and validate URL
	u, err := url.Parse(urlStr)
	if err != nil {
		fmt.Printf("   ❌ Invalid URL format: %v\n", err)
		return fmt.Errorf("invalid URL format: %v", err)
	}

	fmt.Printf("   ✅ URL is valid\n")
	fmt.Printf("   Scheme: %s\n", u.Scheme)
	fmt.Printf("   Host: %s\n", u.Host)
	if u.Path != "" && u.Path != "/" {
		fmt.Printf("   Path: %s\n", u.Path)
	}

	// Check scheme
	validSchemes := map[string]bool{
		"http": true, "https": true, "file": true,
		"ftp": true, "mailto": true,
	}

	if !validSchemes[u.Scheme] {
		fmt.Printf("   ⚠️  Unsupported URL scheme: %s\n", u.Scheme)
		fmt.Printf("   Supported schemes: http, https, file, ftp, mailto\n")
	} else {
		fmt.Printf("   ✅ Supported URL scheme: %s\n", u.Scheme)
	}

	// Warn about localhost URLs
	if u.Hostname() == "localhost" || u.Hostname() == "127.0.0.1" {
		fmt.Printf("   ⚠️  localhost URL - will only work on this machine\n")
	}

	fmt.Println("\n✅ URL action validation completed successfully")
	return nil
}

// getDefaultShell returns the default shell for the current platform
func getDefaultShell() string {
	switch runtime.GOOS {
	case "windows":
		return "cmd.exe"
	default:
		if shell := os.Getenv("SHELL"); shell != "" {
			return shell
		}
		return "/bin/sh"
	}
}

// dryRun simulates spell execution without actually running it
func dryRun(spellName string, debug bool) error {
	if spellName == "" {
		return fmt.Errorf("spell name is required when using --dry-run flag")
	}

	fmt.Printf("🔍 Dry Run Mode: %s\n", spellName)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println()

	// Load configuration
	configPath := getConfigPath()
	fmt.Printf("📁 Loading configuration from: %s\n", configPath)
	
	loader := config.NewLoader(configPath)
	cfg, err := loader.Load()
	if err != nil {
		return errors.Wrap(errors.ErrorTypeConfig, "failed to load configuration", err)
	}
	fmt.Println("✅ Configuration loaded successfully")
	fmt.Println()

	// Initialize minimal logger for dry-run mode (console only)
	logLevel := "info"
	if debug {
		logLevel = "debug"
	}
	loggerConfig := logger.Config{
		Level:   logLevel,
		Console: true,
	}
	if err := logger.Initialize(loggerConfig); err != nil {
		return errors.Wrap(errors.ErrorTypeSystem, "failed to initialize logger", err)
	}

	// Check if spell exists
	fmt.Printf("🔍 Looking up spell '%s'...\n", spellName)
	actionName, exists := cfg.Shortcuts[spellName]
	if !exists {
		fmt.Printf("❌ Spell '%s' not found\n", spellName)
		fmt.Printf("Available spells: %v\n", getSpellList(cfg.Shortcuts))
		return fmt.Errorf("spell '%s' not found", spellName)
	}
	fmt.Printf("✅ Spell found: %s → %s\n", spellName, actionName)
	fmt.Println()

	// Check if action exists
	fmt.Printf("🎯 Looking up action '%s'...\n", actionName)
	action, actionExists := cfg.Actions[actionName]
	if !actionExists {
		fmt.Printf("❌ Action '%s' not found in grimoire\n", actionName)
		return fmt.Errorf("action '%s' not found in grimoire", actionName)
	}
	fmt.Printf("✅ Action found in grimoire\n")
	fmt.Println()

	// Display what would be executed
	fmt.Println("🚀 Would Execute:")
	fmt.Printf("   Type: %s\n", action.Type)
	fmt.Printf("   Command: %s\n", action.Command)
	
	// Expand environment variables for display
	expandedCmd := os.ExpandEnv(action.Command)
	if expandedCmd != action.Command {
		fmt.Printf("   Expanded Command: %s\n", expandedCmd)
	}

	if action.Description != "" {
		fmt.Printf("   Description: %s\n", action.Description)
	}

	if len(action.Args) > 0 {
		fmt.Printf("   Arguments: %v\n", action.Args)
		// Show expanded arguments
		var expandedArgs []string
		for _, arg := range action.Args {
			expandedArgs = append(expandedArgs, os.ExpandEnv(arg))
		}
		fmt.Printf("   Expanded Arguments: %v\n", expandedArgs)
	}

	if action.WorkingDir != "" {
		fmt.Printf("   Working Directory: %s\n", action.WorkingDir)
		expandedDir := os.ExpandEnv(action.WorkingDir)
		if expandedDir != action.WorkingDir {
			fmt.Printf("   Expanded Working Directory: %s\n", expandedDir)
		}
	}

	if action.Shell != "" {
		fmt.Printf("   Shell: %s\n", action.Shell)
	}

	if action.Timeout > 0 {
		fmt.Printf("   Timeout: %d seconds\n", action.Timeout)
	}
	fmt.Println()

	// Display action options
	fmt.Println("⚙️  Action Options:")
	var options []string
	if action.ShowOutput {
		options = append(options, "show_output")
	}
	if action.KeepOpen {
		options = append(options, "keep_open")
	}
	if action.Terminal {
		options = append(options, "terminal")
	}
	if action.ForceTerminal {
		options = append(options, "force_terminal")
	}
	if action.Admin {
		options = append(options, "admin")
	}
	if len(options) > 0 {
		fmt.Printf("   Enabled: %s\n", strings.Join(options, ", "))
	} else {
		fmt.Println("   No special options enabled")
	}
	fmt.Println()

	// Display environment variables
	if len(action.Env) > 0 {
		fmt.Println("🌍 Environment Variables:")
		for key, value := range action.Env {
			fmt.Printf("   %s=%s\n", key, value)
			// Show expanded value if different
			expandedValue := os.ExpandEnv(value)
			if expandedValue != value {
				fmt.Printf("   %s=%s (expanded)\n", key, expandedValue)
			}
		}
		fmt.Println()
	}

	// Type-specific dry-run information
	fmt.Printf("🔬 Type-Specific Analysis (%s):\n", action.Type)
	switch action.Type {
	case "app":
		return dryRunAppAction(action)
	case "script":
		return dryRunScriptAction(action)
	case "url":
		return dryRunURLAction(action)
	default:
		fmt.Printf("❌ Unknown action type: %s\n", action.Type)
		return fmt.Errorf("unknown action type: %s", action.Type)
	}
}

// dryRunAppAction shows what would happen for app actions
func dryRunAppAction(action config.ActionConfig) error {
	expandedCmd := os.ExpandEnv(action.Command)
	fmt.Printf("   Would launch application: %s\n", action.Command)
	if expandedCmd != action.Command {
		fmt.Printf("   Expanded path: %s\n", expandedCmd)
	}

	// Check if it's an absolute path
	if filepath.IsAbs(expandedCmd) {
		fmt.Printf("   Type: Absolute path\n")
		if _, err := os.Stat(expandedCmd); os.IsNotExist(err) {
			fmt.Printf("   ❌ File does not exist: %s\n", expandedCmd)
			fmt.Printf("   Would fail with: application file does not exist\n")
		} else {
			fmt.Printf("   ✅ File exists: %s\n", expandedCmd)
			fmt.Printf("   Would successfully launch application\n")
		}
	} else {
		fmt.Printf("   Type: Command in PATH\n")
		if fullPath, err := exec.LookPath(expandedCmd); err != nil {
			fmt.Printf("   ❌ Command not found in PATH: %s\n", expandedCmd)
			fmt.Printf("   Would fail with: application not found in PATH\n")
		} else {
			fmt.Printf("   ✅ Found in PATH: %s\n", fullPath)
			fmt.Printf("   Would successfully launch: %s\n", fullPath)
		}
	}

	fmt.Println("\n✅ Dry run analysis completed successfully")
	return nil
}

// dryRunScriptAction shows what would happen for script actions
func dryRunScriptAction(action config.ActionConfig) error {
	fmt.Printf("   Would execute script: %s\n", action.Command)
	
	// Show expanded command
	expandedCmd := os.ExpandEnv(action.Command)
	if expandedCmd != action.Command {
		fmt.Printf("   Expanded command: %s\n", expandedCmd)
	}

	// Check shell
	if action.Shell != "" {
		fmt.Printf("   Would use shell: %s\n", action.Shell)
		if !filepath.IsAbs(action.Shell) {
			if fullPath, err := exec.LookPath(action.Shell); err != nil {
				fmt.Printf("   ❌ Shell not found: %s\n", action.Shell)
				fmt.Printf("   Would fail with: shell not found\n")
			} else {
				fmt.Printf("   ✅ Shell found: %s\n", fullPath)
			}
		} else {
			if _, err := os.Stat(action.Shell); os.IsNotExist(err) {
				fmt.Printf("   ❌ Shell does not exist: %s\n", action.Shell)
				fmt.Printf("   Would fail with: shell does not exist\n")
			} else {
				fmt.Printf("   ✅ Shell exists: %s\n", action.Shell)
			}
		}
	} else {
		fmt.Printf("   Would use default shell: %s\n", getDefaultShell())
	}

	// Check working directory
	if action.WorkingDir != "" {
		expandedDir := os.ExpandEnv(action.WorkingDir)
		fmt.Printf("   Would run in directory: %s\n", expandedDir)
		if _, err := os.Stat(expandedDir); os.IsNotExist(err) {
			fmt.Printf("   ❌ Working directory does not exist: %s\n", expandedDir)
			fmt.Printf("   Would fail with: working directory does not exist\n")
		} else {
			fmt.Printf("   ✅ Working directory exists: %s\n", expandedDir)
		}
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("   ⚠️  Failed to get current directory: %v\n", err)
			fmt.Printf("   Would run in current directory: [unknown]\n")
		} else {
			fmt.Printf("   Would run in current directory: %s\n", cwd)
		}
	}

	fmt.Println("\n✅ Dry run analysis completed successfully")
	return nil
}

// hotkeyConfigEqual compares two hotkey configurations for equality
func hotkeyConfigEqual(a, b *config.HotkeyConfig) bool {
	return a.Prefix == b.Prefix &&
		a.Timeout == b.Timeout &&
		a.SequenceTimeout == b.SequenceTimeout
}

// shortcutsEqual compares two shortcut maps for equality
func shortcutsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}

	for k, v := range a {
		if bv, exists := b[k]; !exists || v != bv {
			return false
		}
	}

	return true
}

// dryRunURLAction shows what would happen for URL actions
func dryRunURLAction(action config.ActionConfig) error {
	urlStr := strings.TrimSpace(action.Command)
	fmt.Printf("   Would open URL: %s\n", urlStr)

	// Add scheme if missing
	if !strings.Contains(urlStr, "://") {
		urlStr = "https://" + urlStr
		fmt.Printf("   URL with scheme: %s\n", urlStr)
	}

	// Parse and validate URL
	u, err := url.Parse(urlStr)
	if err != nil {
		fmt.Printf("   ❌ Invalid URL format: %v\n", err)
		fmt.Printf("   Would fail with: invalid URL format\n")
		return nil
	}

	fmt.Printf("   ✅ URL is valid\n")
	fmt.Printf("   Scheme: %s\n", u.Scheme)
	fmt.Printf("   Host: %s\n", u.Host)
	if u.Path != "" && u.Path != "/" {
		fmt.Printf("   Path: %s\n", u.Path)
	}

	// Check scheme
	validSchemes := map[string]bool{
		"http": true, "https": true, "file": true,
		"ftp": true, "mailto": true,
	}

	if !validSchemes[u.Scheme] {
		fmt.Printf("   ⚠️  Unsupported URL scheme: %s\n", u.Scheme)
		fmt.Printf("   Supported schemes: http, https, file, ftp, mailto\n")
		fmt.Printf("   Would attempt to open with default browser\n")
	} else {
		fmt.Printf("   ✅ Supported URL scheme: %s\n", u.Scheme)
		fmt.Printf("   Would successfully open in default browser\n")
	}

	// Warn about localhost URLs
	if u.Hostname() == "localhost" || u.Hostname() == "127.0.0.1" {
		fmt.Printf("   ⚠️  localhost URL - will only work on this machine\n")
	}

	fmt.Println("\n✅ Dry run analysis completed successfully")
	return nil
}

