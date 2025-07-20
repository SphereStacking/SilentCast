package benchmarks

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/SphereStacking/silentcast/internal/action"
	"github.com/SphereStacking/silentcast/internal/notify"
	"github.com/SphereStacking/silentcast/internal/output"
)

// BenchmarkMemoryAllocation measures memory allocation patterns
func BenchmarkMemoryAllocation(b *testing.B) {
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		env := SetupBenchmarkEnvironment(b)
		
		// Simulate full application lifecycle
		if err := env.WriteConfigFile("spellbook.yml", GetBenchmarkConfig()); err != nil {
			b.Fatalf("Failed to write config: %v", err)
		}
		if err := env.LoadConfig(); err != nil {
			b.Fatalf("Failed to load config: %v", err)
		}
		if err := env.InitializeComponents(); err != nil {
			b.Fatalf("Failed to initialize components: %v", err)
		}
		
		// Execute some actions
		_ = env.ExecuteAction("editor")
		_ = env.ExecuteAction("git_status")
		
		env.Cleanup()
	}
}

// BenchmarkMemoryUsageOverTime measures memory usage patterns over time
func BenchmarkMemoryUsageOverTime(b *testing.B) {
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
	
	for i := 0; i < b.N; i++ {
		// Simulate sustained usage
		for j := 0; j < 10; j++ {
			action := actions[j%len(actions)]
			_ = env.ExecuteAction(action)
		}
		
		// Force garbage collection periodically
		if i%100 == 0 {
			runtime.GC()
		}
	}
}

// BenchmarkMemoryLeakDetection measures potential memory leaks
func BenchmarkMemoryLeakDetection(b *testing.B) {
	var initialStats, finalStats runtime.MemStats
	
	// Get baseline memory stats
	runtime.GC()
	runtime.ReadMemStats(&initialStats)
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		env := SetupBenchmarkEnvironment(b)
		
		if err := env.WriteConfigFile("spellbook.yml", GetBenchmarkConfig()); err != nil {
			b.Fatalf("Failed to write config: %v", err)
		}
		if err := env.LoadConfig(); err != nil {
			b.Fatalf("Failed to load config: %v", err)
		}
		if err := env.InitializeComponents(); err != nil {
			b.Fatalf("Failed to initialize components: %v", err)
		}
		
		// Execute multiple actions
		for j := 0; j < 5; j++ {
			_ = env.ExecuteAction("git_status")
		}
		
		env.Cleanup()
		
		// Force GC every 50 iterations
		if i%50 == 0 {
			runtime.GC()
		}
	}
	
	b.StopTimer()
	
	// Check for memory leaks
	runtime.GC()
	runtime.ReadMemStats(&finalStats)
	
	allocDiff := finalStats.TotalAlloc - initialStats.TotalAlloc
	b.ReportMetric(float64(allocDiff)/float64(b.N), "leaked-B/op")
}

// BenchmarkGoroutineUsage measures goroutine creation and cleanup
func BenchmarkGoroutineUsage(b *testing.B) {
	initialGoroutines := runtime.NumGoroutine()
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		env := SetupBenchmarkEnvironment(b)
		
		if err := env.WriteConfigFile("spellbook.yml", GetBenchmarkConfig()); err != nil {
			b.Fatalf("Failed to write config: %v", err)
		}
		if err := env.LoadConfig(); err != nil {
			b.Fatalf("Failed to load config: %v", err)
		}
		if err := env.InitializeComponents(); err != nil {
			b.Fatalf("Failed to initialize components: %v", err)
		}
		
		// Execute concurrent actions
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		go func() {
			_ = env.ActionManager.Execute(ctx, "git_status")
		}()
		
		cancel()
		env.Cleanup()
	}
	
	b.StopTimer()
	
	// Allow goroutines to cleanup
	time.Sleep(100 * time.Millisecond)
	finalGoroutines := runtime.NumGoroutine()
	
	if finalGoroutines > initialGoroutines+10 { // Allow some tolerance
		b.Errorf("Potential goroutine leak: started with %d, ended with %d", initialGoroutines, finalGoroutines)
	}
}

// BenchmarkMemoryPressure measures performance under memory pressure
func BenchmarkMemoryPressure(b *testing.B) {
	// Allocate memory to simulate pressure
	pressure := make([][]byte, 0, 1000)
	for i := 0; i < 500; i++ {
		pressure = append(pressure, make([]byte, 1024*1024)) // 1MB each
	}
	
	defer func() {
		pressure = nil
		runtime.GC()
	}()
	
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
		// Force some memory pressure operations
		_ = env.ExecuteAction("git_status")
		
		// Allocate and release memory
		temp := make([]byte, 64*1024) // 64KB
		_ = temp
		temp = nil
	}
}

// BenchmarkMemoryPooling measures memory pooling efficiency
func BenchmarkMemoryPooling(b *testing.B) {
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
		// Use output manager multiple times to test pooling
		writer := env.OutputManager.StartCapture()
		_, _ = writer.Write([]byte("test output"))
		_ = env.OutputManager.GetOutput()
		env.OutputManager.Clear()
	}
}

// BenchmarkStringAllocation measures string allocation overhead
func BenchmarkStringAllocation(b *testing.B) {
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
		// Test common string operations
		for actionName := range env.Config.Actions {
			_ = "action: " + actionName
			_ = actionName + " executed"
		}
		
		for spellName := range env.Config.Shortcuts {
			_ = "spell: " + spellName
		}
	}
}

// BenchmarkStructAllocation measures struct allocation patterns
func BenchmarkStructAllocation(b *testing.B) {
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		// Simulate creating various structs
		_, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		
		type TestStruct struct {
			Name        string
			Description string
			Args        []string
			Env         map[string]string
		}
		
		test := &TestStruct{
			Name:        "benchmark_test",
			Description: "Benchmark test structure",
			Args:        []string{"arg1", "arg2", "arg3"},
			Env:         map[string]string{"KEY": "value"},
		}
		
		_ = test
		cancel()
	}
}

// BenchmarkMapOperations measures map allocation and operations
func BenchmarkMapOperations(b *testing.B) {
	env := SetupBenchmarkEnvironment(b)
	defer env.Cleanup()
	
	if err := env.WriteConfigFile("spellbook.yml", GetLargeConfig()); err != nil {
		b.Fatalf("Failed to write large config: %v", err)
	}
	if err := env.LoadConfig(); err != nil {
		b.Fatalf("Failed to load large config: %v", err)
	}
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		// Test map operations
		actionMap := make(map[string]string)
		for name := range env.Config.Actions {
			actionMap[name] = "processed"
		}
		
		// Lookup operations
		for name := range actionMap {
			_ = actionMap[name]
		}
		
		// Delete operations
		for name := range actionMap {
			delete(actionMap, name)
		}
	}
}

// BenchmarkSliceOperations measures slice allocation and operations
func BenchmarkSliceOperations(b *testing.B) {
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		// Test slice operations
		var items []string
		
		// Append operations
		for j := 0; j < 100; j++ {
			items = append(items, "item_"+string(rune('a'+j%26)))
		}
		
		// Access operations
		for _, item := range items {
			_ = item
		}
		
		// Slice operations
		if len(items) > 50 {
			_ = items[0:50]
			_ = items[50:]
		}
	}
}

// BenchmarkInterfaceAllocation measures interface allocation overhead
func BenchmarkInterfaceAllocation(b *testing.B) {
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
		// Test interface boxing/unboxing
		var outputManager interface{} = env.OutputManager
		var notifyManager interface{} = env.NotifyManager
		var actionManager interface{} = env.ActionManager
		
		// Type assertions
		_ = outputManager.(output.Manager)
		_ = notifyManager.(*notify.Manager)
		_ = actionManager.(*action.Manager)
	}
}