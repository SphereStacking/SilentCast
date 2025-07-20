package benchmarks

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/SphereStacking/silentcast/internal/config"
	"github.com/SphereStacking/silentcast/internal/notify"
)

// BenchmarkFileWatcherSetup measures file watcher initialization performance
func BenchmarkFileWatcherSetup(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()

	if err := env.WriteConfigFile("spellbook.yml", GetBenchmarkConfig()); err != nil {
		b.Fatalf("Failed to write config: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		watcher, err := config.NewWatcher(config.WatcherConfig{
			ConfigPath: env.ConfigDir,
			OnChange:   func(*config.Config) {},
			Debounce:   100 * time.Millisecond,
		})
		if err != nil {
			b.Fatalf("Failed to create watcher: %v", err)
		}
		watcher.Stop()
	}
}

// BenchmarkConfigReload measures configuration reload performance
func BenchmarkConfigReload(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()

	if err := env.WriteConfigFile("spellbook.yml", GetBenchmarkConfig()); err != nil {
		b.Fatalf("Failed to write config: %v", err)
	}
	if err := env.LoadConfig(); err != nil {
		b.Fatalf("Failed to load config: %v", err)
	}
	if err := env.InitializeComponents(); err != nil {
		b.Fatalf("Failed to initialize components: %v", err)
	}

	// Alternative configs for reloading
	configs := []string{
		GetBenchmarkConfig(),
		`spells: {"r": "reload_test"}
grimoire:
  reload_test: {type: script, command: echo, args: ["reloaded"]}`,
		`spells: {"x": "exit_test"}
grimoire:
  exit_test: {type: script, command: echo, args: ["exit"]}`,
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		config := configs[i%len(configs)]
		
		// Write new config
		if err := env.WriteConfigFile("spellbook.yml", config); err != nil {
			b.Fatalf("Failed to write config: %v", err)
		}
		
		// Reload config
		if err := env.LoadConfig(); err != nil {
			b.Fatalf("Failed to reload config: %v", err)
		}
		
		// Update action manager
		env.ActionManager.UpdateActions(env.Config.Actions)
	}
}

// BenchmarkRapidConfigChanges measures rapid configuration file changes
func BenchmarkRapidConfigChanges(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()

	if err := env.WriteConfigFile("spellbook.yml", GetBenchmarkConfig()); err != nil {
		b.Fatalf("Failed to write config: %v", err)
	}
	if err := env.LoadConfig(); err != nil {
		b.Fatalf("Failed to load config: %v", err)
	}
	if err := env.InitializeComponents(); err != nil {
		b.Fatalf("Failed to initialize components: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Simulate rapid config changes
		for j := 0; j < 5; j++ {
			config := `spells: {"rapid` + string(rune('a'+j)) + `": "rapid_action"}
grimoire:
  rapid_action: {type: script, command: echo, args: ["rapid"]}`

			if err := env.WriteConfigFile("spellbook.yml", config); err != nil {
				b.Fatalf("Failed to write rapid config: %v", err)
			}
			
			// Small delay to simulate file system events
			time.Sleep(1 * time.Millisecond)
		}
		
		// Final reload
		if err := env.LoadConfig(); err != nil {
			b.Fatalf("Failed to reload after rapid changes: %v", err)
		}
	}
}

// BenchmarkFileWatcherEvents measures file watcher event processing
func BenchmarkFileWatcherEvents(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()

	if err := env.WriteConfigFile("spellbook.yml", GetBenchmarkConfig()); err != nil {
		b.Fatalf("Failed to write config: %v", err)
	}

	// Channel to count reload events
	reloadChan := make(chan struct{}, 100)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	watcher, err := config.NewWatcher(config.WatcherConfig{
		ConfigPath: env.ConfigDir,
		OnChange: func(*config.Config) {
			reloadChan <- struct{}{}
		},
		Debounce: 10 * time.Millisecond,
	})
	if err != nil {
		b.Fatalf("Failed to create watcher: %v", err)
	}
	defer watcher.Stop()

	// Start watching
	watcher.Start(ctx)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Trigger file change
		config := `spells: {"event` + string(rune('0'+i%10)) + `": "event_action"}
grimoire:
  event_action: {type: script, command: echo, args: ["event"]}`

		if err := env.WriteConfigFile("spellbook.yml", config); err != nil {
			b.Fatalf("Failed to write config for event: %v", err)
		}

		// Wait for watcher event (with timeout)
		select {
		case <-reloadChan:
			// Event received
		case <-time.After(100 * time.Millisecond):
			// Timeout - continue anyway
		}
	}
}

// BenchmarkMultipleConfigFiles measures watching multiple config files
func BenchmarkMultipleConfigFiles(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()

	// Create base config
	if err := env.WriteConfigFile("spellbook.yml", GetBenchmarkConfig()); err != nil {
		b.Fatalf("Failed to write base config: %v", err)
	}

	// Create platform-specific config
	platformFile := "spellbook." + runtime.GOOS + ".yml"
	platformConfig := `daemon: {log_level: debug}
grimoire:
  platform_action: {type: script, command: echo, args: ["platform"]}`

	if err := env.WriteConfigFile(platformFile, platformConfig); err != nil {
		b.Fatalf("Failed to write platform config: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		loader := config.NewLoader(env.ConfigDir)
		
		// Load all config files
		_, err := loader.Load()
		if err != nil {
			b.Fatalf("Failed to load multiple configs: %v", err)
		}

		// Modify base config
		if err := env.WriteConfigFile("spellbook.yml", GetBenchmarkConfig()); err != nil {
			b.Fatalf("Failed to update base config: %v", err)
		}

		// Modify platform config
		if err := env.WriteConfigFile(platformFile, platformConfig); err != nil {
			b.Fatalf("Failed to update platform config: %v", err)
		}
	}
}

// BenchmarkConfigWatcherMemoryUsage measures memory usage of file watcher
func BenchmarkConfigWatcherMemoryUsage(b *testing.B) {
	RunMemoryBenchmark(b, func() {
		env := SetupBenchmarkEnvironment(b)
		defer env.Cleanup()

		if err := env.WriteConfigFile("spellbook.yml", GetBenchmarkConfig()); err != nil {
			b.Fatalf("Failed to write config: %v", err)
		}

		watcher, err := config.NewWatcher(config.WatcherConfig{
			ConfigPath: env.ConfigDir,
			OnChange:   func(*config.Config) {},
			Debounce:   100 * time.Millisecond,
		})
		if err != nil {
			b.Fatalf("Failed to create watcher: %v", err)
		}
		defer watcher.Stop()

		// Simulate file changes
		for i := 0; i < 10; i++ {
			config := `spells: {"mem` + string(rune('0'+i)) + `": "mem_action"}
grimoire:
  mem_action: {type: script, command: echo, args: ["memory"]}`

			if err := env.WriteConfigFile("spellbook.yml", config); err != nil {
				b.Fatalf("Failed to write config: %v", err)
			}
			time.Sleep(10 * time.Millisecond) // Allow watcher to process
		}
	})
}

// BenchmarkConfigReloadNotification measures notification during config reload
func BenchmarkConfigReloadNotification(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()

	if err := env.WriteConfigFile("spellbook.yml", GetBenchmarkConfig()); err != nil {
		b.Fatalf("Failed to write config: %v", err)
	}
	if err := env.LoadConfig(); err != nil {
		b.Fatalf("Failed to load config: %v", err)
	}
	if err := env.InitializeComponents(); err != nil {
		b.Fatalf("Failed to initialize components: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Update config
		config := `spells: {"notify` + string(rune('0'+i%10)) + `": "notify_action"}
grimoire:
  notify_action: {type: script, command: echo, args: ["notify"]}`

		if err := env.WriteConfigFile("spellbook.yml", config); err != nil {
			b.Fatalf("Failed to write config: %v", err)
		}

		// Reload and notify
		if err := env.LoadConfig(); err != nil {
			b.Fatalf("Failed to reload config: %v", err)
		}

		// Simulate config reload notification
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		_ = env.NotifyManager.Notify(ctx, notify.Notification{
			Title:   "Config Reloaded",
			Message: "Configuration has been reloaded successfully",
			Level:   notify.LevelInfo,
		})
		cancel()
	}
}

// BenchmarkLargeConfigReload measures reload performance with large configs
func BenchmarkLargeConfigReload(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()

	largeConfig := GetLargeConfig()
	if err := env.WriteConfigFile("spellbook.yml", largeConfig); err != nil {
		b.Fatalf("Failed to write large config: %v", err)
	}
	if err := env.LoadConfig(); err != nil {
		b.Fatalf("Failed to load large config: %v", err)
	}
	if err := env.InitializeComponents(); err != nil {
		b.Fatalf("Failed to initialize components: %v", err)
	}

	// Create alternative large config
	altLargeConfig := `spells:`
	for i := 0; i < 500; i++ {
		char1 := 'a' + rune(i%26)
		char2 := 'a' + rune((i/26)%26)
		key := string(char1) + string(char2)
		altLargeConfig += "\n  \"" + key + "\": \"large_action_" + string(rune('0'+i%10)) + "\""
	}
	altLargeConfig += "\ngrimoire:"
	for i := 0; i < 10; i++ {
		altLargeConfig += "\n  large_action_" + string(rune('0'+i)) + ":"
		altLargeConfig += "\n    type: script"
		altLargeConfig += "\n    command: echo"
		altLargeConfig += "\n    args: [\"large_" + string(rune('0'+i)) + "\"]"
	}

	configs := []string{largeConfig, altLargeConfig}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		config := configs[i%len(configs)]
		
		if err := env.WriteConfigFile("spellbook.yml", config); err != nil {
			b.Fatalf("Failed to write large config: %v", err)
		}
		
		if err := env.LoadConfig(); err != nil {
			b.Fatalf("Failed to reload large config: %v", err)
		}
		
		env.ActionManager.UpdateActions(env.Config.Actions)
	}
}

// BenchmarkConfigWatcherStartupShutdown measures watcher lifecycle performance
func BenchmarkConfigWatcherStartupShutdown(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()

	if err := env.WriteConfigFile("spellbook.yml", GetBenchmarkConfig()); err != nil {
		b.Fatalf("Failed to write config: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Start watcher
		watcher, err := config.NewWatcher(config.WatcherConfig{
			ConfigPath: env.ConfigDir,
			OnChange:   func(*config.Config) {},
			Debounce:   100 * time.Millisecond,
		})
		if err != nil {
			b.Fatalf("Failed to create watcher: %v", err)
		}

		// Brief operation
		time.Sleep(1 * time.Millisecond)

		// Stop watcher
		watcher.Stop()
	}
}

// BenchmarkConcurrentConfigOperations measures concurrent config operations
func BenchmarkConcurrentConfigOperations(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()

	if err := env.WriteConfigFile("spellbook.yml", GetBenchmarkConfig()); err != nil {
		b.Fatalf("Failed to write config: %v", err)
	}
	if err := env.LoadConfig(); err != nil {
		b.Fatalf("Failed to load config: %v", err)
	}
	if err := env.InitializeComponents(); err != nil {
		b.Fatalf("Failed to initialize components: %v", err)
	}

	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			// Concurrent config reads
			_ = len(env.Config.Actions)
			_ = env.Config.Hotkeys.Prefix

			// Occasional config write
			if counter%10 == 0 {
				config := `spells: {"concurrent` + string(rune('0'+counter%10)) + `": "concurrent_action"}
grimoire:
  concurrent_action: {type: script, command: echo, args: ["concurrent"]}`

				env.WriteConfigFile("spellbook.yml", config)
			}
			counter++
		}
	})
}