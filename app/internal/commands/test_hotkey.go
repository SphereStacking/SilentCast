package commands

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/SphereStacking/silentcast/internal/config"
	"github.com/SphereStacking/silentcast/internal/hotkey"
	"github.com/SphereStacking/silentcast/pkg/logger"
)

// TestHotkeyCommand tests hotkey detection
type TestHotkeyCommand struct {
	getConfigPath func() string
}

// NewTestHotkeyCommand creates a new test hotkey command
func NewTestHotkeyCommand(getConfigPath func() string) Command {
	return &TestHotkeyCommand{
		getConfigPath: getConfigPath,
	}
}

// Name returns the command name
func (c *TestHotkeyCommand) Name() string {
	return "Test Hotkey"
}

// Description returns the command description
func (c *TestHotkeyCommand) Description() string {
	return "Test hotkey detection"
}

// FlagName returns the flag name
func (c *TestHotkeyCommand) FlagName() string {
	return "test-hotkey"
}

// IsActive checks if the command should run
func (c *TestHotkeyCommand) IsActive(flags interface{}) bool {
	f, ok := flags.(*Flags)
	if !ok {
		return false
	}
	return f.TestHotkey
}

// Execute runs the command
func (c *TestHotkeyCommand) Execute(flags interface{}) error {
	f, ok := flags.(*Flags)
	if !ok {
		return fmt.Errorf("invalid flags type")
	}

	fmt.Println("🔍 SilentCast Hotkey Test Mode")
	fmt.Println("================================")
	fmt.Println()

	// Initialize logger for debugging
	loggerConfig := logger.Config{
		Level:   "debug",
		Console: true,
	}
	if err := logger.Initialize(loggerConfig); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Load configuration to get hotkey settings
	loader := config.NewLoader(c.getConfigPath())
	cfg, err := loader.Load()
	if err != nil {
		// Use default settings if config fails
		fmt.Println("⚠️  Could not load config, using defaults")
		cfg = &config.Config{
			Hotkeys: config.HotkeyConfig{
				Prefix:          "alt+space",
				Timeout:         1000,
				SequenceTimeout: 2000,
			},
		}
	}

	fmt.Printf("📋 Configuration:\n")
	fmt.Printf("   Prefix key: %s\n", cfg.Hotkeys.Prefix)
	fmt.Printf("   Timeout: %dms\n", cfg.Hotkeys.Timeout)
	fmt.Printf("   Sequence timeout: %dms\n", cfg.Hotkeys.SequenceTimeout)
	fmt.Println()

	// Create hotkey manager
	manager, err := hotkey.NewManager(&cfg.Hotkeys)
	if err != nil {
		return fmt.Errorf("failed to create hotkey manager: %w", err)
	}

	// Track key events
	type keyEvent struct {
		timestamp time.Time
		sequence  string
		spellName string
		keys      []hotkey.Key
	}

	var events []keyEvent

	// Set up test handler
	manager.SetHandler(hotkey.HandlerFunc(func(event hotkey.Event) error {
		ke := keyEvent{
			timestamp: event.Timestamp,
			sequence:  event.Sequence.String(),
			spellName: event.SpellName,
			keys:      event.Sequence.Keys,
		}
		events = append(events, ke)

		// Display event in real-time
		fmt.Printf("\n🎯 Key Event Detected:\n")
		fmt.Printf("   Time: %s\n", ke.timestamp.Format("15:04:05.000"))
		fmt.Printf("   Sequence: %s\n", ke.sequence)

		// Show details for each key in the sequence
		for i, key := range ke.keys {
			fmt.Printf("   Key %d: %s", i+1, key.Name)
			if len(key.Modifiers) > 0 {
				fmt.Printf(" (modifiers: %s)", strings.Join(key.Modifiers, "+"))
			}
			fmt.Printf(" [code: 0x%04X]\n", key.Code)
		}

		// Show if this would trigger a spell
		if event.SpellName != "" {
			fmt.Printf("   ✨ Would trigger spell: %s\n", event.SpellName)
		}

		return nil
	}))

	// Register a test spell
	if err := manager.Register("test", "test_spell"); err != nil {
		logger.Warn("Could not register test spell: %v", err)
	}

	// Start the manager
	if err := manager.Start(); err != nil {
		return fmt.Errorf("failed to start hotkey manager: %w", err)
	}
	defer func() {
		if err := manager.Stop(); err != nil {
			logger.Warn("Failed to stop hotkey manager: %v", err)
		}
	}()

	fmt.Println("🎮 Hotkey test mode active!")
	fmt.Println()
	fmt.Println("Instructions:")
	fmt.Printf("1. Press %s to activate prefix mode\n", cfg.Hotkeys.Prefix)
	fmt.Println("2. Then press any key to see the full sequence")
	fmt.Println("3. Try multi-key sequences like 'g,s' or 't,e,s,t'")
	fmt.Println("4. Press Ctrl+C to exit")
	fmt.Println()
	fmt.Println("Listening for key events...")
	fmt.Println()

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Set up duration timer if specified
	var timer *time.Timer
	if f.TestDuration > 0 {
		timer = time.NewTimer(time.Duration(f.TestDuration) * time.Second)
		fmt.Printf("⏱️  Test will run for %d seconds\n\n", f.TestDuration)
	}

	// Wait for exit
	select {
	case <-sigChan:
		fmt.Println("\n\n🛑 Interrupted by user")
	case <-func() <-chan time.Time {
		if timer != nil {
			return timer.C
		}
		return make(chan time.Time) // Never fires
	}():
		fmt.Printf("\n\n⏱️  Test duration of %d seconds completed\n", f.TestDuration)
	}

	// Show summary
	fmt.Println("\n📊 Test Summary:")
	fmt.Printf("   Total events captured: %d\n", len(events))

	if len(events) > 0 {
		fmt.Println("\n   Event History:")
		for i, evt := range events {
			fmt.Printf("   %d. %s - %s",
				i+1,
				evt.timestamp.Format("15:04:05"),
				evt.sequence,
			)
			if evt.spellName != "" {
				fmt.Printf(" → %s", evt.spellName)
			}
			fmt.Println()
		}

		// Analyze common issues
		fmt.Println("\n💡 Analysis:")

		// Check if prefix key was detected
		prefixDetected := false
		for _, evt := range events {
			if strings.Contains(evt.sequence, cfg.Hotkeys.Prefix) {
				prefixDetected = true
				break
			}
		}

		if !prefixDetected {
			fmt.Printf("   ⚠️  Prefix key '%s' was never detected.\n", cfg.Hotkeys.Prefix)
			fmt.Println("      This could mean:")
			fmt.Println("      - The key combination is not supported on your system")
			fmt.Println("      - Another application is intercepting this hotkey")
			fmt.Println("      - You need to grant accessibility permissions")
		} else {
			fmt.Println("   ✅ Prefix key detection is working correctly")
		}

		// Check for timing issues
		var timingIssues []string
		for i := 1; i < len(events); i++ {
			gap := events[i].timestamp.Sub(events[i-1].timestamp)
			if gap > time.Duration(cfg.Hotkeys.Timeout)*time.Millisecond {
				timingIssues = append(timingIssues,
					fmt.Sprintf("Gap of %v between events %d and %d", gap, i, i+1))
			}
		}

		if len(timingIssues) > 0 {
			fmt.Println("\n   ⏱️  Timing issues detected:")
			for _, issue := range timingIssues {
				fmt.Printf("      - %s\n", issue)
			}
			fmt.Printf("      Consider increasing timeout (currently %dms)\n", cfg.Hotkeys.Timeout)
		}
	} else {
		fmt.Println("   ❌ No key events were captured!")
		fmt.Println("\n   Possible issues:")
		fmt.Println("   - Hotkey detection is not working (check permissions)")
		fmt.Println("   - The application needs to run with elevated privileges")
		fmt.Println("   - Another application may be blocking global hotkeys")
	}

	return nil
}

// Group returns the command group
func (c *TestHotkeyCommand) Group() string {
	return "debug"
}

// HasOptions returns if this command has additional options
func (c *TestHotkeyCommand) HasOptions() bool {
	return true
}
