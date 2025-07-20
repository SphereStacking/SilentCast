package benchmarks

import (
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/SphereStacking/silentcast/internal/config"
	"gopkg.in/yaml.v3"
)

// BenchmarkYAMLParsing measures YAML parsing performance
func BenchmarkYAMLParsing(b *testing.B) {
	configContent := GetBenchmarkConfig()
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		var cfg config.Config
		err := yaml.Unmarshal([]byte(configContent), &cfg)
		if err != nil {
			b.Fatalf("Failed to parse YAML: %v", err)
		}
	}
}

// BenchmarkYAMLParsingLarge measures YAML parsing with large configs
func BenchmarkYAMLParsingLarge(b *testing.B) {
	configContent := GetLargeConfig()
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		var cfg config.Config
		err := yaml.Unmarshal([]byte(configContent), &cfg)
		if err != nil {
			b.Fatalf("Failed to parse large YAML: %v", err)
		}
	}
}

// BenchmarkConfigResolution measures platform resolver performance
func BenchmarkConfigResolution(b *testing.B) {
	resolver := config.GetPlatformResolver()
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		// Test platform-specific operations
		_ = resolver.GetDefaultConfigPath()
		_ = resolver.GetPlatformConfigFile()
	}
}

// BenchmarkEnvironmentVariableExpansion measures env var expansion performance
func BenchmarkEnvironmentVariableExpansion(b *testing.B) {
	// Set test environment variables
	os.Setenv("BENCH_HOME", "/benchmark/home")
	os.Setenv("BENCH_USER", "benchmark_user")
	os.Setenv("BENCH_PATH", "/benchmark/bin:/usr/bin")
	defer func() {
		os.Unsetenv("BENCH_HOME")
		os.Unsetenv("BENCH_USER")
		os.Unsetenv("BENCH_PATH")
	}()
	
	testStrings := []string{
		"${BENCH_HOME}/projects",
		"${BENCH_USER}@example.com",
		"PATH=${BENCH_PATH}:${BENCH_HOME}/bin",
		"${BENCH_HOME}/.config/${BENCH_USER}/settings",
		"No variables here",
		"${UNDEFINED_VAR}",
		"${BENCH_HOME}/deep/nested/${BENCH_USER}/path",
	}
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		for _, str := range testStrings {
			_ = os.ExpandEnv(str)
		}
	}
}

// BenchmarkConfigValidation measures configuration validation performance
func BenchmarkConfigValidation(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()
	
	// Load valid configuration
	if err := env.WriteConfigFile("spellbook.yml", GetBenchmarkConfig()); err != nil {
		b.Fatalf("Failed to write config: %v", err)
	}
	if err := env.LoadConfig(); err != nil {
		b.Fatalf("Failed to load config: %v", err)
	}
	
	validator := config.NewValidator()
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_ = validator.Validate(env.Config)
	}
}

// BenchmarkConfigValidationErrors measures validation with errors
func BenchmarkConfigValidationErrors(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()
	
	// Create invalid configuration
	invalidConfig := `
spells:
  e: "nonexistent_action"
  t: "another_missing"
grimoire:
  valid_action:
    type: script
    command: "echo"
  invalid_action:
    type: "invalid_type"
    command: "echo"
  missing_command:
    type: "script"
    description: "No command field"
`
	
	if err := env.WriteConfigFile("spellbook.yml", invalidConfig); err != nil {
		b.Fatalf("Failed to write invalid config: %v", err)
	}
	if err := env.LoadConfig(); err != nil {
		b.Fatalf("Failed to load invalid config: %v", err)
	}
	
	validator := config.NewValidator()
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_ = validator.Validate(env.Config)
	}
}

// BenchmarkConfigFileSearch measures config file search performance
func BenchmarkConfigFileSearch(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()
	
	// Create multiple config files
	if err := env.WriteConfigFile("spellbook.yml", GetBenchmarkConfig()); err != nil {
		b.Fatalf("Failed to write base config: %v", err)
	}
	
	platformFile := "spellbook." + runtime.GOOS + ".yml"
	if err := env.WriteConfigFile(platformFile, "# Platform override"); err != nil {
		b.Fatalf("Failed to write platform config: %v", err)
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

// BenchmarkConfigCascading measures configuration cascading performance
func BenchmarkConfigCascading(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()
	
	// Base configuration
	baseConfig := GetBenchmarkConfig()
	if err := env.WriteConfigFile("spellbook.yml", baseConfig); err != nil {
		b.Fatalf("Failed to write base config: %v", err)
	}
	
	// Platform-specific configuration
	platformConfig := `
daemon:
  log_level: debug

grimoire:
  editor:
    command: "platform-specific-editor"
    args: ["--platform-flag"]
  platform_action:
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

// BenchmarkConfigSerialization measures config serialization back to YAML
func BenchmarkConfigSerialization(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()
	
	if err := env.WriteConfigFile("spellbook.yml", GetBenchmarkConfig()); err != nil {
		b.Fatalf("Failed to write config: %v", err)
	}
	if err := env.LoadConfig(); err != nil {
		b.Fatalf("Failed to load config: %v", err)
	}
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := yaml.Marshal(env.Config)
		if err != nil {
			b.Fatalf("Failed to serialize config: %v", err)
		}
	}
}

// BenchmarkConfigMemoryUsage measures memory usage during config operations
func BenchmarkConfigMemoryUsage(b *testing.B) {
	RunMemoryBenchmark(b, func() {
		env := SetupBenchmarkEnvironment(b)
		defer env.Cleanup()
		
		// Load large configuration
		if err := env.WriteConfigFile("spellbook.yml", GetLargeConfig()); err != nil {
			b.Fatalf("Failed to write large config: %v", err)
		}
		if err := env.LoadConfig(); err != nil {
			b.Fatalf("Failed to load large config: %v", err)
		}
		
		// Perform various operations
		validator := config.NewValidator()
		_ = validator.Validate(env.Config)
		
		_, _ = yaml.Marshal(env.Config)
	})
}

// BenchmarkConfigStringProcessing measures string processing during config loading
func BenchmarkConfigStringProcessing(b *testing.B) {
	testStrings := []string{
		"simple_string",
		"string with spaces",
		"string-with-dashes",
		"string_with_underscores",
		"MixedCaseString",
		"string.with.dots",
		"string/with/slashes",
		"string\\with\\backslashes",
		"string@with#special$chars%",
		"veryLongStringWithLotsOfCharactersThatMightCausePerformanceIssues",
	}
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		for _, str := range testStrings {
			// Simulate common string operations during config processing
			_ = strings.ToLower(str)
			_ = strings.TrimSpace(str)
			_ = strings.Contains(str, "_")
			parts := strings.Split(str, "_")
			_ = strings.Join(parts, "-")
		}
	}
}

// BenchmarkConfigDeepCopy measures deep copying configuration objects
func BenchmarkConfigDeepCopy(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()
	
	if err := env.WriteConfigFile("spellbook.yml", GetBenchmarkConfig()); err != nil {
		b.Fatalf("Failed to write config: %v", err)
	}
	if err := env.LoadConfig(); err != nil {
		b.Fatalf("Failed to load config: %v", err)
	}
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		// Simulate deep copy by marshaling and unmarshaling
		data, err := yaml.Marshal(env.Config)
		if err != nil {
			b.Fatalf("Failed to marshal config: %v", err)
		}
		
		var newConfig config.Config
		err = yaml.Unmarshal(data, &newConfig)
		if err != nil {
			b.Fatalf("Failed to unmarshal config: %v", err)
		}
	}
}

// BenchmarkConfigUpdate measures configuration update performance
func BenchmarkConfigUpdate(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()
	
	// Initial configuration
	if err := env.WriteConfigFile("spellbook.yml", GetBenchmarkConfig()); err != nil {
		b.Fatalf("Failed to write initial config: %v", err)
	}
	if err := env.LoadConfig(); err != nil {
		b.Fatalf("Failed to load initial config: %v", err)
	}
	if err := env.InitializeComponents(); err != nil {
		b.Fatalf("Failed to initialize components: %v", err)
	}
	
	// Updated configuration
	updatedConfig := `
spells:
  updated: "updated_action"
grimoire:
  updated_action:
    type: script
    command: "echo"
    args: ["updated"]
`
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		// Write updated config
		if err := env.WriteConfigFile("spellbook.yml", updatedConfig); err != nil {
			b.Fatalf("Failed to write updated config: %v", err)
		}
		
		// Reload config
		if err := env.LoadConfig(); err != nil {
			b.Fatalf("Failed to reload config: %v", err)
		}
		
		// Update components
		env.ActionManager.UpdateActions(env.Config.Actions)
	}
}

// BenchmarkConfigConcurrentAccess measures concurrent config access
func BenchmarkConfigConcurrentAccess(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()
	
	if err := env.WriteConfigFile("spellbook.yml", GetBenchmarkConfig()); err != nil {
		b.Fatalf("Failed to write config: %v", err)
	}
	if err := env.LoadConfig(); err != nil {
		b.Fatalf("Failed to load config: %v", err)
	}
	
	b.ReportAllocs()
	
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Simulate concurrent reads
			_ = len(env.Config.Actions)
			_ = len(env.Config.Shortcuts)
			_ = env.Config.Hotkeys.Prefix
			_ = env.Config.Daemon.LogLevel
		}
	})
}