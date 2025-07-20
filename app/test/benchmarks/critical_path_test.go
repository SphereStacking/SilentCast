package benchmarks

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/SphereStacking/silentcast/internal/action/script"
	"github.com/SphereStacking/silentcast/internal/config"
	"github.com/SphereStacking/silentcast/internal/notify"
	"github.com/SphereStacking/silentcast/internal/output"
)

// Critical path benchmarks focus on the most performance-sensitive operations

func BenchmarkCriticalActionExecution(b *testing.B) {
	b.Run("SimpleScript", func(b *testing.B) {
		RunWithProfiling(b, "action_simple_script", func(b *testing.B) {
			// env := SetupBenchmarkEnvironment(b) // Not used in this simplified version
			
			cfg := &config.ActionConfig{
				Type:    "script",
				Command: "echo 'test'",
			}
			
			executor := script.NewScriptExecutor(cfg)
			ctx := context.Background()
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = executor.Execute(ctx)
			}
		})
	})

	b.Run("ScriptWithOutput", func(b *testing.B) {
		RunWithProfiling(b, "action_script_output", func(b *testing.B) {
			// env := SetupBenchmarkEnvironment(b) // Not used in this simplified version
			
			cfg := &config.ActionConfig{
				Type:       "script",
				Command:    "echo 'test output line 1\ntest output line 2\ntest output line 3'",
				ShowOutput: true,
			}
			
			executor := script.NewScriptExecutor(cfg)
			ctx := context.Background()
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = executor.Execute(ctx)
			}
		})
	})

	b.Run("ParallelExecution", func(b *testing.B) {
		RunWithProfiling(b, "action_parallel", func(b *testing.B) {
			// env := SetupBenchmarkEnvironment(b) // Not used in this simplified version
			
			cfg := &config.ActionConfig{
				Type:    "script",
				Command: "echo 'parallel test'",
			}
			
			ctx := context.Background()
			
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					executor := script.NewScriptExecutor(cfg)
					_ = executor.Execute(ctx)
				}
			})
		})
	})
}

func BenchmarkCriticalConfigLoading(b *testing.B) {
	b.Run("SmallConfig", func(b *testing.B) {
		RunWithProfiling(b, "config_small", func(b *testing.B) {
			env := SetupBenchmarkEnvironment(b)
			configContent := GetBenchmarkConfig()
			env.WriteConfigFile("spellbook.yml", configContent)
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				loader := config.NewLoader(env.ConfigDir)
				_, _ = loader.Load()
			}
		})
	})

	b.Run("LargeConfig", func(b *testing.B) {
		RunWithProfiling(b, "config_large", func(b *testing.B) {
			env := SetupBenchmarkEnvironment(b)
			configContent := GetLargeConfig()
			env.WriteConfigFile("spellbook.yml", configContent)
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				loader := config.NewLoader(env.ConfigDir)
				_, _ = loader.Load()
			}
		})
	})

	b.Run("ConfigValidation", func(b *testing.B) {
		RunWithProfiling(b, "config_validation", func(b *testing.B) {
			env := SetupBenchmarkEnvironment(b)
			configContent := GetBenchmarkConfig()
			env.WriteConfigFile("spellbook.yml", configContent)
			
			loader := config.NewLoader(env.ConfigDir)
			cfg, _ := loader.Load()
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				validator := config.NewValidator()
				_ = validator.Validate(cfg)
			}
		})
	})
}

func BenchmarkNotificationSystem(b *testing.B) {
	b.Run("SimpleNotification", func(b *testing.B) {
		RunWithProfiling(b, "notify_simple", func(b *testing.B) {
			mgr := notify.NewManager()
			ctx := context.Background()
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				mgr.Info(ctx, "Test", "Message")
			}
		})
	})

	b.Run("ConcurrentNotifications", func(b *testing.B) {
		RunWithProfiling(b, "notify_concurrent", func(b *testing.B) {
			mgr := notify.NewManager()
			ctx := context.Background()
			
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				i := 0
				for pb.Next() {
					mgr.Info(ctx, fmt.Sprintf("Test %d", i), "Concurrent message")
					i++
				}
			})
		})
	})

	b.Run("NotificationWithOutput", func(b *testing.B) {
		RunWithProfiling(b, "notify_output", func(b *testing.B) {
			mgr := notify.NewManager()
			ctx := context.Background()
			
			longOutput := "Line 1\nLine 2\nLine 3\nLine 4\nLine 5\n"
			for i := 0; i < 10; i++ {
				longOutput += fmt.Sprintf("Additional line %d with some content\n", i)
			}
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				mgr.Success(ctx, "Long Output", longOutput)
			}
		})
	})
}

func BenchmarkOutputManager(b *testing.B) {
	b.Run("BufferedCapture", func(b *testing.B) {
		RunWithProfiling(b, "output_buffered", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				mgr := output.NewBufferedManager(output.DefaultOptions())
				writer := mgr.StartCapture()
				
				// Write typical command output
				for j := 0; j < 100; j++ {
					fmt.Fprintf(writer, "Line %d: Some command output with timestamp %s\n", 
						j, time.Now().Format(time.RFC3339))
				}
				
				mgr.Stop()
				_ = mgr.GetOutput()
			}
		})
	})

	b.Run("StreamingCapture", func(b *testing.B) {
		RunWithProfiling(b, "output_streaming", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				mgr := output.NewStreamingManager(output.Options{})
				writer := mgr.StartCapture()
				
				// Write typical command output
				for j := 0; j < 100; j++ {
					fmt.Fprintf(writer, "Line %d: Streaming output\n", j)
				}
				
				mgr.Stop()
			}
		})
	})
}

func BenchmarkCriticalMemoryAllocation(b *testing.B) {
	b.Run("ActionCreation", func(b *testing.B) {
		RunMemoryBenchmark(b, func() {
			cfg := &config.ActionConfig{
				Type:        "script",
				Command:     "echo test",
				Description: "Test action",
				Args:        []string{"arg1", "arg2"},
				Env:         map[string]string{"KEY": "value"},
			}
			_ = script.NewScriptExecutor(cfg)
		})
	})

	b.Run("NotificationCreation", func(b *testing.B) {
		RunMemoryBenchmark(b, func() {
			mgr := notify.NewManager()
			_ = mgr
		})
	})

	b.Run("ConfigParsing", func(b *testing.B) {
		RunMemoryBenchmark(b, func() {
			// Simulate config parsing by creating a new config
			cfg := &config.Config{
				Daemon: config.DaemonConfig{
					AutoStart: false,
					LogLevel:  "error",
				},
			}
			_ = cfg
		})
	})
}

func BenchmarkGoroutineManagement(b *testing.B) {
	b.Run("ActionExecutionGoroutines", func(b *testing.B) {
		baseline := runtime.NumGoroutine()
		
		RunWithProfiling(b, "goroutine_actions", func(b *testing.B) {
			cfg := &config.ActionConfig{
				Type:    "script",
				Command: "echo test",
			}
			
			ctx := context.Background()
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				executor := script.NewScriptExecutor(cfg)
				_ = executor.Execute(ctx)
			}
			b.StopTimer()
			
			CheckGoroutineLeaks(b, baseline, 5)
		})
	})

	b.Run("ConcurrentActions", func(b *testing.B) {
		baseline := runtime.NumGoroutine()
		
		RunWithProfiling(b, "goroutine_concurrent", func(b *testing.B) {
			cfg := &config.ActionConfig{
				Type:    "script",
				Command: "echo concurrent",
			}
			
			ctx := context.Background()
			
			b.ResetTimer()
			var wg sync.WaitGroup
			for i := 0; i < b.N; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					executor := script.NewScriptExecutor(cfg)
					_ = executor.Execute(ctx)
				}()
				
				// Limit concurrent goroutines
				if i%100 == 0 {
					wg.Wait()
				}
			}
			wg.Wait()
			b.StopTimer()
			
			CheckGoroutineLeaks(b, baseline, 10)
		})
	})
}

func BenchmarkResourceUsage(b *testing.B) {
	b.Run("LongRunningAction", func(b *testing.B) {
		before := CaptureMemorySnapshot()
		
		RunWithProfiling(b, "resource_longrunning", func(b *testing.B) {
			cfg := &config.ActionConfig{
				Type:       "script",
				Command:    "sleep 0.1 && echo done",
				ShowOutput: true,
				Timeout:    1,
			}
			
			ctx := context.Background()
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				executor := script.NewScriptExecutor(cfg)
				_ = executor.Execute(ctx)
			}
			b.StopTimer()
			
			after := CaptureMemorySnapshot()
			b.Log(CompareSnapshots(before, after))
		})
	})

	b.Run("HighFrequencyActions", func(b *testing.B) {
		before := CaptureMemorySnapshot()
		
		RunWithProfiling(b, "resource_highfreq", func(b *testing.B) {
			env := SetupBenchmarkEnvironment(b)
			env.WriteConfigFile("spellbook.yml", GetLargeConfig())
			env.LoadConfig()
			env.InitializeComponents()
			
			// ctx := context.Background() // Not used
			actions := []string{"action_a0", "action_b1", "action_c2", "action_d3"}
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				action := actions[i%len(actions)]
				_ = env.ExecuteAction(action)
			}
			b.StopTimer()
			
			after := CaptureMemorySnapshot()
			b.Log(CompareSnapshots(before, after))
		})
	})
}