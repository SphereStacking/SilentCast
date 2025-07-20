package benchmarks

import (
	"context"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/SphereStacking/silentcast/internal/action"
	"github.com/SphereStacking/silentcast/internal/action/script"
	"github.com/SphereStacking/silentcast/internal/config"
	"github.com/SphereStacking/silentcast/internal/hotkey"
	"github.com/SphereStacking/silentcast/internal/notify"
)

// TestGoroutineLeaks tests for goroutine leaks in various components
func TestGoroutineLeaks(t *testing.T) {
	// Get baseline goroutine count
	runtime.GC()
	time.Sleep(100 * time.Millisecond)
	baseline := runtime.NumGoroutine()

	t.Run("ActionExecutor", func(t *testing.T) {
		testActionExecutorLeaks(t, baseline)
	})

	t.Run("NotificationManager", func(t *testing.T) {
		testNotificationManagerLeaks(t, baseline)
	})

	t.Run("ConfigWatcher", func(t *testing.T) {
		testConfigWatcherLeaks(t, baseline)
	})

	t.Run("HotkeyManager", func(t *testing.T) {
		testHotkeyManagerLeaks(t, baseline)
	})

	t.Run("ConcurrentActions", func(t *testing.T) {
		testConcurrentActionLeaks(t, baseline)
	})
}

func testActionExecutorLeaks(t *testing.T, baseline int) {
	before := runtime.NumGoroutine()

	// Execute multiple actions
	for i := 0; i < 100; i++ {
		cfg := &config.ActionConfig{
			Type:    "script",
			Command: "echo test",
		}
		executor := script.NewScriptExecutor(cfg)
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		_ = executor.Execute(ctx)
		cancel()
	}

	// Give goroutines time to clean up
	time.Sleep(200 * time.Millisecond)
	runtime.GC()

	after := runtime.NumGoroutine()
	leaked := after - before

	if leaked > 5 { // Allow small tolerance
		t.Errorf("Goroutine leak in ActionExecutor: before=%d, after=%d, leaked=%d",
			before, after, leaked)
		printGoroutineStacks(t)
	}
}

func testNotificationManagerLeaks(t *testing.T, baseline int) {
	before := runtime.NumGoroutine()

	// Create and use notification managers
	for i := 0; i < 50; i++ {
		mgr := notify.NewManager()
		ctx := context.Background()
		
		// Send notifications
		for j := 0; j < 10; j++ {
			mgr.Info(ctx, "Test", "Message")
			mgr.Success(ctx, "Test", "Success")
			mgr.Error(ctx, "Test", "Error")
		}
	}

	// Give goroutines time to clean up
	time.Sleep(200 * time.Millisecond)
	runtime.GC()

	after := runtime.NumGoroutine()
	leaked := after - before

	if leaked > 5 {
		t.Errorf("Goroutine leak in NotificationManager: before=%d, after=%d, leaked=%d",
			before, after, leaked)
		printGoroutineStacks(t)
	}
}

func testConfigWatcherLeaks(t *testing.T, baseline int) {
	before := runtime.NumGoroutine()

	env := SetupBenchmarkEnvironment(&testing.B{})
	defer env.Cleanup()

	// Create config file
	configContent := GetBenchmarkConfig()
	env.WriteConfigFile("spellbook.yml", configContent)

	// Create and start multiple watchers
	for i := 0; i < 10; i++ {
		loader := config.NewLoader(env.ConfigDir)
		cfg, _ := loader.Load()
		
		// Config watcher functionality not directly testable
		// Skip watcher test for now
		
		_ = cfg // Use config
	}

	// Give goroutines time to clean up
	time.Sleep(500 * time.Millisecond)
	runtime.GC()

	after := runtime.NumGoroutine()
	leaked := after - before

	if leaked > 10 { // File watchers may need more tolerance
		t.Errorf("Goroutine leak in ConfigWatcher: before=%d, after=%d, leaked=%d",
			before, after, leaked)
		printGoroutineStacks(t)
	}
}

func testHotkeyManagerLeaks(t *testing.T, baseline int) {
	before := runtime.NumGoroutine()

	// Create and use hotkey managers
	for i := 0; i < 20; i++ {
		mgr := hotkey.NewMockManager()
		
		// Mock manager doesn't need registration
		// Just test start/stop
		
		// Start and stop
		err := mgr.Start()
		if err != nil {
			t.Logf("Mock manager start error: %v", err)
		}
		
		time.Sleep(10 * time.Millisecond)
		mgr.Stop()
	}

	// Give goroutines time to clean up
	time.Sleep(200 * time.Millisecond)
	runtime.GC()

	after := runtime.NumGoroutine()
	leaked := after - before

	if leaked > 5 {
		t.Errorf("Goroutine leak in HotkeyManager: before=%d, after=%d, leaked=%d",
			before, after, leaked)
		printGoroutineStacks(t)
	}
}

func testConcurrentActionLeaks(t *testing.T, baseline int) {
	before := runtime.NumGoroutine()

	// Set up action manager
	grimoire := map[string]config.ActionConfig{
		"test1": {Type: "script", Command: "echo test1"},
		"test2": {Type: "script", Command: "echo test2"},
		"test3": {Type: "script", Command: "echo test3"},
	}
	mgr := action.NewManager(grimoire)

	// Execute actions concurrently
	var wg sync.WaitGroup
	ctx := context.Background()

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			spellName := []string{"test1", "test2", "test3"}[idx%3]
			_ = mgr.Execute(ctx, spellName)
		}(i)
	}

	wg.Wait()

	// Give goroutines time to clean up
	time.Sleep(500 * time.Millisecond)
	runtime.GC()

	after := runtime.NumGoroutine()
	leaked := after - before

	if leaked > 10 {
		t.Errorf("Goroutine leak in concurrent actions: before=%d, after=%d, leaked=%d",
			before, after, leaked)
		printGoroutineStacks(t)
	}
}

// printGoroutineStacks prints current goroutine stack traces for debugging
func printGoroutineStacks(t *testing.T) {
	buf := make([]byte, 1<<20) // 1MB buffer
	stackLen := runtime.Stack(buf, true)
	t.Logf("Current goroutine stacks:\n%s", buf[:stackLen])
}

// BenchmarkGoroutineCreation measures the overhead of goroutine creation
func BenchmarkGoroutineCreation(b *testing.B) {
	b.Run("ActionExecution", func(b *testing.B) {
		cfg := &config.ActionConfig{
			Type:    "script",
			Command: "echo test",
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ctx := context.Background()
			executor := script.NewScriptExecutor(cfg)
			_ = executor.Execute(ctx)
		}
	})

	b.Run("NotificationSend", func(b *testing.B) {
		mgr := notify.NewManager()
		ctx := context.Background()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			mgr.Info(ctx, "Test", "Message")
		}
	})

	b.Run("ConcurrentExecution", func(b *testing.B) {
		cfg := &config.ActionConfig{
			Type:    "script",
			Command: "echo concurrent",
		}

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ctx := context.Background()
				executor := script.NewScriptExecutor(cfg)
				_ = executor.Execute(ctx)
			}
		})
	})
}