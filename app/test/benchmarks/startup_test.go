package benchmarks

import (
	"os"
	"runtime"
	"testing"

	"github.com/SphereStacking/silentcast/internal/config"
	"github.com/SphereStacking/silentcast/internal/action"
	"github.com/SphereStacking/silentcast/internal/hotkey"
	"github.com/SphereStacking/silentcast/internal/notify"
	"github.com/SphereStacking/silentcast/internal/output"
	"github.com/SphereStacking/silentcast/internal/permission"
)

// BenchmarkApplicationStartup measures the complete application startup time
func BenchmarkApplicationStartup(b *testing.B) {
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		env := SetupBenchmarkEnvironment(b)
		
		// Write configuration
		if err := env.WriteConfigFile("spellbook.yml", GetBenchmarkConfig()); err != nil {
			b.Fatalf("Failed to write config: %v", err)
		}
		
		// Load configuration
		if err := env.LoadConfig(); err != nil {
			b.Fatalf("Failed to load config: %v", err)
		}
		
		// Initialize all components
		if err := env.InitializeComponents(); err != nil {
			b.Fatalf("Failed to initialize components: %v", err)
		}
		
		env.Cleanup()
	}
}

// BenchmarkConfigLoading measures configuration loading performance
func BenchmarkConfigLoading(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()
	
	// Write configuration file once
	if err := env.WriteConfigFile("spellbook.yml", GetBenchmarkConfig()); err != nil {
		b.Fatalf("Failed to write config: %v", err)
	}
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		loader := config.NewLoader(env.ConfigDir)
		_, err := loader.Load()
		if err != nil {
			b.Fatalf("Failed to load config: %v", err)
		}
	}
}

// BenchmarkConfigLoadingLarge measures loading performance with large configurations
func BenchmarkConfigLoadingLarge(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()
	
	// Write large configuration file
	if err := env.WriteConfigFile("spellbook.yml", GetLargeConfig()); err != nil {
		b.Fatalf("Failed to write large config: %v", err)
	}
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		loader := config.NewLoader(env.ConfigDir)
		_, err := loader.Load()
		if err != nil {
			b.Fatalf("Failed to load large config: %v", err)
		}
	}
}

// BenchmarkStartupConfigValidation measures configuration validation performance during startup
func BenchmarkStartupConfigValidation(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()
	
	// Load configuration once
	if err := env.WriteConfigFile("spellbook.yml", GetBenchmarkConfig()); err != nil {
		b.Fatalf("Failed to write config: %v", err)
	}
	if err := env.LoadConfig(); err != nil {
		b.Fatalf("Failed to load config: %v", err)
	}
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		validator := config.NewValidator()
		_ = validator.Validate(env.Config)
	}
}

// BenchmarkComponentInitialization measures individual component initialization
func BenchmarkComponentInitialization(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()
	
	// Load configuration once
	if err := env.WriteConfigFile("spellbook.yml", GetBenchmarkConfig()); err != nil {
		b.Fatalf("Failed to write config: %v", err)
	}
	if err := env.LoadConfig(); err != nil {
		b.Fatalf("Failed to load config: %v", err)
	}
	
	b.Run("OutputManager", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = output.NewBufferedManager(output.DefaultOptions())
		}
	})
	
	b.Run("NotificationManager", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = notify.NewManager()
		}
	})
	
	b.Run("PermissionManager", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = permission.NewManager()
		}
	})
	
	b.Run("ActionManager", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = action.NewManager(env.Config.Actions)
		}
	})
	
	b.Run("HotkeyManager", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = hotkey.NewMockManager()
		}
	})
}

// BenchmarkMemoryUsage measures memory usage during startup
func BenchmarkMemoryUsage(b *testing.B) {
	var startStats runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&startStats)
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		env := SetupBenchmarkEnvironment(b)
		
		// Complete startup sequence
		if err := env.WriteConfigFile("spellbook.yml", GetBenchmarkConfig()); err != nil {
			b.Fatalf("Failed to write config: %v", err)
		}
		if err := env.LoadConfig(); err != nil {
			b.Fatalf("Failed to load config: %v", err)
		}
		if err := env.InitializeComponents(); err != nil {
			b.Fatalf("Failed to initialize components: %v", err)
		}
		
		// Measure memory usage
		var stats runtime.MemStats
		runtime.ReadMemStats(&stats)
		
		env.Cleanup()
	}
}

// BenchmarkConfigCascade measures platform-specific configuration cascading
func BenchmarkConfigCascade(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()
	
	// Write base configuration
	baseConfig := GetBenchmarkConfig()
	if err := env.WriteConfigFile("spellbook.yml", baseConfig); err != nil {
		b.Fatalf("Failed to write base config: %v", err)
	}
	
	// Write platform-specific configuration
	platformConfig := `
grimoire:
  editor:
    command: "platform-specific-editor"
    args: ["--platform-flag"]
  new_platform_action:
    type: script
    command: "platform-command"
`
	
	platformFile := "spellbook." + runtime.GOOS + ".yml"
	if err := env.WriteConfigFile(platformFile, platformConfig); err != nil {
		b.Fatalf("Failed to write platform config: %v", err)
	}
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		loader := config.NewLoader(env.ConfigDir)
		_, err := loader.Load()
		if err != nil {
			b.Fatalf("Failed to load cascaded config: %v", err)
		}
	}
}

// BenchmarkStartupWithWatcher measures startup time including file watcher
func BenchmarkStartupWithWatcher(b *testing.B) {
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		env := SetupBenchmarkEnvironment(b)
		
		// Write configuration
		if err := env.WriteConfigFile("spellbook.yml", GetBenchmarkConfig()); err != nil {
			b.Fatalf("Failed to write config: %v", err)
		}
		
		// Load configuration
		if err := env.LoadConfig(); err != nil {
			b.Fatalf("Failed to load config: %v", err)
		}
		
		// Initialize components
		if err := env.InitializeComponents(); err != nil {
			b.Fatalf("Failed to initialize components: %v", err)
		}
		
		// Setup file watcher
		watcherConfig := config.WatcherConfig{
			ConfigPath: env.ConfigDir,
			OnChange: func(cfg *config.Config) {
				// Simulated config change handler
			},
		}
		
		watcher, err := config.NewWatcher(watcherConfig)
		if err != nil {
			b.Fatalf("Failed to create watcher: %v", err)
		}
		
		// Start and stop watcher
		watcher.Start(env.ctx)
		_ = watcher.Stop()
		
		env.Cleanup()
	}
}

// BenchmarkConfigReload measures configuration reload performance
func BenchmarkStartupConfigReload(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()
	
	// Initial setup
	if err := env.WriteConfigFile("spellbook.yml", GetBenchmarkConfig()); err != nil {
		b.Fatalf("Failed to write config: %v", err)
	}
	if err := env.LoadConfig(); err != nil {
		b.Fatalf("Failed to load config: %v", err)
	}
	if err := env.InitializeComponents(); err != nil {
		b.Fatalf("Failed to initialize components: %v", err)
	}
	
	// New configuration for reloading
	newConfig := `
spells:
  r: "reload_test"
grimoire:
  reload_test:
    type: script
    command: "echo"
    args: ["reloaded"]
`
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		// Write new configuration
		if err := env.WriteConfigFile("spellbook.yml", newConfig); err != nil {
			b.Fatalf("Failed to write new config: %v", err)
		}
		
		// Reload configuration
		if err := env.LoadConfig(); err != nil {
			b.Fatalf("Failed to reload config: %v", err)
		}
		
		// Update action manager
		env.ActionManager.UpdateActions(env.Config.Actions)
	}
}

// BenchmarkEnvironmentVariableExpansion measures env var expansion performance
func BenchmarkStartupEnvironmentVariableExpansion(b *testing.B) {
	// Set test environment variables
	os.Setenv("BENCH_EDITOR", "echo")  // Use echo command which exists in PATH
	os.Setenv("BENCH_HOME", "/benchmark/home")
	defer func() {
		os.Unsetenv("BENCH_EDITOR")
		os.Unsetenv("BENCH_HOME")
	}()
	
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()
	
	configWithEnvVars := `
spells:
  e: "editor"
grimoire:
  editor:
    type: app
    command: "${BENCH_EDITOR}"
    working_dir: "${BENCH_HOME}/projects"
    env:
      EDITOR_CONFIG: "${BENCH_HOME}/.editor"
`
	
	if err := env.WriteConfigFile("spellbook.yml", configWithEnvVars); err != nil {
		b.Fatalf("Failed to write config with env vars: %v", err)
	}
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		loader := config.NewLoader(env.ConfigDir)
		_, err := loader.Load()
		if err != nil {
			b.Fatalf("Failed to load config with env vars: %v", err)
		}
	}
}