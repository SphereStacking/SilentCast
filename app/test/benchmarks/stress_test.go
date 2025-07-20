package benchmarks

import (
	"context"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/SphereStacking/silentcast/internal/hotkey"
	"github.com/SphereStacking/silentcast/internal/notify"
)

// BenchmarkStressTest_HighFrequencyActions measures performance under high action frequency
func BenchmarkStressTest_HighFrequencyActions(b *testing.B) {
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
	
	actions := []string{"editor", "terminal", "git_status", "browser"}
	
	b.ReportAllocs()
	b.ResetTimer()
	
	// Execute actions rapidly
	for i := 0; i < b.N; i++ {
		action := actions[i%len(actions)]
		err := env.ExecuteAction(action)
		if err != nil {
			b.Fatalf("High frequency action failed: %v", err)
		}
	}
}

// BenchmarkStressTest_MassiveConfiguration measures performance with very large configs
func BenchmarkStressTest_MassiveConfiguration(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()
	
	// Create massive configuration
	massiveConfig := `
daemon:
  auto_start: false
  log_level: error
  config_watch: false

hotkeys:
  prefix: "alt+space"
  timeout: 1000

spells:
`
	
	// Add 1000 spell mappings
	for i := 0; i < 1000; i++ {
		char1 := 'a' + rune(i%26)
		char2 := 'a' + rune((i/26)%26)
		num := i % 10
		spellKey := string(char1) + string(char2) + string(rune('0'+num))
		actionName := "massive_action_" + spellKey
		massiveConfig += "  \"" + spellKey + "\": \"" + actionName + "\"\n"
	}
	
	massiveConfig += "\ngrimoire:\n"
	
	// Add 1000 actions
	for i := 0; i < 1000; i++ {
		char1 := 'a' + rune(i%26)
		char2 := 'a' + rune((i/26)%26)
		num := i % 10
		spellKey := string(char1) + string(char2) + string(rune('0'+num))
		actionName := "massive_action_" + spellKey
		
		massiveConfig += "  " + actionName + ":\n"
		massiveConfig += "    type: script\n"
		massiveConfig += "    command: \"echo\"\n"
		massiveConfig += "    args: [\"" + actionName + "\"]\n"
		massiveConfig += "    description: \"Massive config action " + spellKey + "\"\n"
	}
	
	if err := env.WriteConfigFile("spellbook.yml", massiveConfig); err != nil {
		b.Fatalf("Failed to write massive config: %v", err)
	}
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		if err := env.LoadConfig(); err != nil {
			b.Fatalf("Failed to load massive config: %v", err)
		}
		if err := env.InitializeComponents(); err != nil {
			b.Fatalf("Failed to initialize with massive config: %v", err)
		}
	}
}

// BenchmarkStressTest_ConcurrentUsers measures concurrent user simulation
func BenchmarkStressTest_ConcurrentUsers(b *testing.B) {
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
	
	actions := []string{"editor", "terminal", "git_status", "browser"}
	
	b.ReportAllocs()
	
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			action := actions[counter%len(actions)]
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			err := env.ActionManager.Execute(ctx, action)
			cancel()
			
			if err != nil {
				b.Fatalf("Concurrent user action failed: %v", err)
			}
			counter++
		}
	})
}

// BenchmarkStressTest_RapidConfigReloads measures rapid configuration reloading
func BenchmarkStressTest_RapidConfigReloads(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()
	
	configs := []string{
		GetBenchmarkConfig(),
		`spells: {"r1": "reload1"}
grimoire:
  reload1: {type: script, command: echo, args: ["reload1"]}`,
		`spells: {"r2": "reload2"}
grimoire:
  reload2: {type: script, command: echo, args: ["reload2"]}`,
		`spells: {"r3": "reload3"}
grimoire:
  reload3: {type: script, command: echo, args: ["reload3"]}`,
	}
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		config := configs[i%len(configs)]
		
		if err := env.WriteConfigFile("spellbook.yml", config); err != nil {
			b.Fatalf("Failed to write config: %v", err)
		}
		if err := env.LoadConfig(); err != nil {
			b.Fatalf("Failed to reload config: %v", err)
		}
		if err := env.InitializeComponents(); err != nil {
			b.Fatalf("Failed to reinitialize: %v", err)
		}
	}
}

// BenchmarkStressTest_MemoryExhaustion measures behavior under memory pressure
func BenchmarkStressTest_MemoryExhaustion(b *testing.B) {
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
	
	// Create memory pressure
	var memoryPressure [][]byte
	for i := 0; i < 100; i++ {
		memoryPressure = append(memoryPressure, make([]byte, 1024*1024)) // 1MB each
	}
	defer func() {
		memoryPressure = nil
	}()
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		// Execute action under memory pressure
		err := env.ExecuteAction("git_status")
		if err != nil {
			b.Fatalf("Action under memory pressure failed: %v", err)
		}
		
		// Add more pressure periodically
		if i%10 == 0 {
			memoryPressure = append(memoryPressure, make([]byte, 512*1024))
		}
	}
}

// BenchmarkStressTest_LongRunningSession measures long-running session performance
func BenchmarkStressTest_LongRunningSession(b *testing.B) {
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
	
	actions := []string{"editor", "terminal", "git_status", "browser"}
	
	b.ReportAllocs()
	b.ResetTimer()
	
	// Simulate sustained usage over time
	for i := 0; i < b.N; i++ {
		// Execute multiple actions in sequence
		for j := 0; j < 10; j++ {
			action := actions[j%len(actions)]
			err := env.ExecuteAction(action)
			if err != nil {
				b.Fatalf("Long session action failed: %v", err)
			}
		}
		
		// Occasional config reload
		if i%50 == 0 {
			if err := env.LoadConfig(); err != nil {
				b.Fatalf("Session config reload failed: %v", err)
			}
		}
		
		// Periodic memory cleanup
		if i%100 == 0 {
			runtime.GC()
		}
	}
}

// BenchmarkStressTest_ErrorRecovery measures recovery from error conditions
func BenchmarkStressTest_ErrorRecovery(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()
	
	// Config with both valid and invalid actions
	mixedConfig := `
spells:
  valid: "valid_action"
  invalid: "invalid_action"
  error: "error_action"
grimoire:
  valid_action:
    type: script
    command: "echo"
    args: ["valid"]
  invalid_action:
    type: script
    command: "nonexistent_command_stress_test"
  error_action:
    type: script
    command: "sh"
    args: ["-c", "exit 1"]
`
	
	if err := env.WriteConfigFile("spellbook.yml", mixedConfig); err != nil {
		b.Fatalf("Failed to write mixed config: %v", err)
	}
	if err := env.LoadConfig(); err != nil {
		b.Fatalf("Failed to load mixed config: %v", err)
	}
	if err := env.InitializeComponents(); err != nil {
		b.Fatalf("Failed to initialize with mixed config: %v", err)
	}
	
	actions := []string{"valid", "invalid", "error", "valid"}
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		action := actions[i%len(actions)]
		// Execute action, errors are expected
		_ = env.ExecuteAction(action)
	}
}

// BenchmarkStressTest_RapidHotkeyRegistration measures rapid hotkey changes
func BenchmarkStressTest_RapidHotkeyRegistration(b *testing.B) {
	manager := hotkey.NewMockManager()
	
	sequences := make([]string, 100)
	for i := range sequences {
		sequences[i] = "stress_" + string(rune('a'+i%26)) + string(rune('0'+i%10))
	}
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		// Register all sequences
		for j, seq := range sequences {
			spellName := "spell_" + string(rune(j+48))
			err := manager.Register(seq, spellName)
			if err != nil {
				b.Fatalf("Rapid registration failed for %s: %v", seq, err)
			}
		}
		
		// Unregister all sequences
		for _, seq := range sequences {
			err := manager.Unregister(seq)
			if err != nil {
				b.Fatalf("Rapid unregistration failed for %s: %v", seq, err)
			}
		}
	}
}

// BenchmarkStressTest_MassiveNotifications measures notification system under load
func BenchmarkStressTest_MassiveNotifications(b *testing.B) {
	manager := notify.NewManager()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		// Send many notifications rapidly
		for j := 0; j < 50; j++ {
			notification := notify.Notification{
				Title:   "Stress Test",
				Message: "Massive notification test " + string(rune(j%10+48)),
				Level:   notify.Level(j % 4),
			}
			
			err := manager.Notify(ctx, notification)
			if err != nil {
				b.Fatalf("Massive notification failed: %v", err)
			}
		}
	}
}

// BenchmarkStressTest_ConcurrentOperations measures mixed concurrent operations
func BenchmarkStressTest_ConcurrentOperations(b *testing.B) {
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
	
	var wg sync.WaitGroup
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		// Concurrent action execution
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = env.ExecuteAction("git_status")
		}()
		
		// Concurrent notification
		wg.Add(1)
		go func() {
			defer wg.Done()
			notification := notify.Notification{
				Title:   "Concurrent Test",
				Message: "Concurrent operation",
				Level:   notify.LevelInfo,
			}
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			_ = env.NotifyManager.Notify(ctx, notification)
			cancel()
		}()
		
		// Concurrent hotkey operation
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			seq := "concurrent_" + string(rune(index%26+'a'))
			_ = env.HotkeyManager.Register(seq, "test_spell")
			_ = env.HotkeyManager.Unregister(seq)
		}(i)
		
		wg.Wait()
	}
}

// BenchmarkStressTest_ResourceExhaustion measures behavior when resources are limited
func BenchmarkStressTest_ResourceExhaustion(b *testing.B) {
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
	
	// Create file descriptor pressure by opening many files
	var files []*os.File
	for i := 0; i < 50; i++ {
		if file, err := os.CreateTemp("", "stress-test-*"); err == nil {
			files = append(files, file)
		}
	}
	defer func() {
		for _, file := range files {
			file.Close()
			os.Remove(file.Name())
		}
	}()
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		// Execute action under resource pressure
		err := env.ExecuteAction("git_status")
		if err != nil {
			b.Fatalf("Action under resource pressure failed: %v", err)
		}
	}
}