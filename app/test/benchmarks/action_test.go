package benchmarks

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/SphereStacking/silentcast/internal/config"
)

// BenchmarkActionExecution measures action execution performance by type
func BenchmarkActionExecution(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()
	
	// Setup environment
	if err := env.WriteConfigFile("spellbook.yml", GetBenchmarkConfig()); err != nil {
		b.Fatalf("Failed to write config: %v", err)
	}
	if err := env.LoadConfig(); err != nil {
		b.Fatalf("Failed to load config: %v", err)
	}
	if err := env.InitializeComponents(); err != nil {
		b.Fatalf("Failed to initialize components: %v", err)
	}
	
	b.Run("AppAction", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			err := env.ExecuteAction("editor")
			if err != nil {
				b.Fatalf("App action failed: %v", err)
			}
		}
	})
	
	b.Run("ScriptAction", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			err := env.ExecuteAction("git_status")
			if err != nil {
				b.Fatalf("Script action failed: %v", err)
			}
		}
	})
	
	b.Run("URLAction", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			err := env.ExecuteAction("browser")
			// URL actions may fail in test environment, that's okay
			_ = err
		}
	})
}

// BenchmarkActionExecutionWithOutput measures script actions with output display
func BenchmarkActionExecutionWithOutput(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()
	
	config := `
spells:
  output_test: "output_action"
grimoire:
  output_action:
    type: script
    command: "echo"
    args: ["benchmark output test"]
    show_output: true
    description: "Output benchmark test"
`
	
	if err := env.WriteConfigFile("spellbook.yml", config); err != nil {
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
		err := env.ExecuteAction("output_action")
		if err != nil {
			b.Fatalf("Output action failed: %v", err)
		}
	}
}

// BenchmarkActionExecutionWithEnvironment measures actions with environment variables
func BenchmarkActionExecutionWithEnvironment(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()
	
	var envCommand string
	var envArgs string
	
	switch runtime.GOOS {
	case "windows":
		envCommand = "cmd"
		envArgs = `["/c", "echo %BENCH_VAR%"]`
	default:
		envCommand = "sh"
		envArgs = `["-c", "echo $BENCH_VAR"]`
	}
	
	config := `
spells:
  env_test: "env_action"
grimoire:
  env_action:
    type: script
    command: "` + envCommand + `"
    args: ` + envArgs + `
    env:
      BENCH_VAR: "benchmark_value"
    show_output: true
    description: "Environment benchmark test"
`
	
	if err := env.WriteConfigFile("spellbook.yml", config); err != nil {
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
		err := env.ExecuteAction("env_action")
		if err != nil {
			b.Fatalf("Environment action failed: %v", err)
		}
	}
}

// BenchmarkActionExecutionWithTimeout measures script actions with timeouts
func BenchmarkActionExecutionWithTimeout(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()
	
	config := `
spells:
  timeout_test: "timeout_action"
grimoire:
  timeout_action:
    type: script
    command: "echo"
    args: ["quick timeout test"]
    timeout: 5
    description: "Timeout benchmark test"
`
	
	if err := env.WriteConfigFile("spellbook.yml", config); err != nil {
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
		err := env.ExecuteAction("timeout_action")
		if err != nil {
			b.Fatalf("Timeout action failed: %v", err)
		}
	}
}

// BenchmarkConcurrentActionExecution measures concurrent action execution
func BenchmarkConcurrentActionExecution(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()
	
	config := `
spells:
  concurrent1: "action1"
  concurrent2: "action2"
  concurrent3: "action3"
grimoire:
  action1:
    type: script
    command: "echo"
    args: ["concurrent 1"]
  action2:
    type: script
    command: "echo"
    args: ["concurrent 2"]
  action3:
    type: script
    command: "echo"
    args: ["concurrent 3"]
`
	
	if err := env.WriteConfigFile("spellbook.yml", config); err != nil {
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
		actions := []string{"action1", "action2", "action3"}
		counter := 0
		
		for pb.Next() {
			action := actions[counter%len(actions)]
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			err := env.ActionManager.Execute(ctx, action)
			cancel()
			
			if err != nil {
				b.Fatalf("Concurrent action %s failed: %v", action, err)
			}
			counter++
		}
	})
}

// BenchmarkActionLookup measures action lookup performance
func BenchmarkActionLookup(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()
	
	// Use large config for lookup testing
	if err := env.WriteConfigFile("spellbook.yml", GetLargeConfig()); err != nil {
		b.Fatalf("Failed to write large config: %v", err)
	}
	if err := env.LoadConfig(); err != nil {
		b.Fatalf("Failed to load config: %v", err)
	}
	if err := env.InitializeComponents(); err != nil {
		b.Fatalf("Failed to initialize components: %v", err)
	}
	
	// Get all action names for lookup testing
	actionNames := make([]string, 0, len(env.Config.Actions))
	for name := range env.Config.Actions {
		actionNames = append(actionNames, name)
	}
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		// Look up random action
		actionName := actionNames[i%len(actionNames)]
		err := env.ExecuteAction(actionName)
		if err != nil {
			b.Fatalf("Action lookup failed for %s: %v", actionName, err)
		}
	}
}

// BenchmarkActionMemoryUsage measures memory usage during action execution
func BenchmarkActionMemoryUsage(b *testing.B) {
	RunMemoryBenchmark(b, func() {
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
		
		// Execute various actions
		actions := []string{"editor", "terminal", "git_status", "browser"}
		for _, action := range actions {
			_ = env.ExecuteAction(action)
		}
	})
}

// BenchmarkActionValidation measures action validation performance
func BenchmarkActionValidation(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()
	
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

// BenchmarkActionManagerUpdate measures action manager update performance
func BenchmarkActionManagerUpdate(b *testing.B) {
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
	
	// New actions for updating
	newConfig := `
spells:
  new1: "new_action1"
  new2: "new_action2"
grimoire:
  new_action1:
    type: script
    command: "echo"
    args: ["new action 1"]
  new_action2:
    type: script
    command: "echo"
    args: ["new action 2"]
`
	
	if err := env.WriteConfigFile("new_spellbook.yml", newConfig); err != nil {
		b.Fatalf("Failed to write new config: %v", err)
	}
	
	// Load new config
	newLoader := config.NewLoader(env.ConfigDir)
	newCfg, err := newLoader.Load()
	if err != nil {
		b.Fatalf("Failed to load new config: %v", err)
	}
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		env.ActionManager.UpdateActions(newCfg.Actions)
	}
}

// BenchmarkActionContextCreation measures context creation overhead
func BenchmarkActionContextCreation(b *testing.B) {
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		cancel() // Immediate cleanup for benchmark
	}
}

// BenchmarkActionErrorHandling measures error handling performance
func BenchmarkActionErrorHandling(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()
	
	config := `
spells:
  error_test: "error_action"
grimoire:
  error_action:
    type: script
    command: "nonexistent_command_benchmark"
    description: "Error benchmark test"
`
	
	if err := env.WriteConfigFile("spellbook.yml", config); err != nil {
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
		// This will fail, but we're measuring error handling performance
		_ = env.ExecuteAction("error_action")
	}
}