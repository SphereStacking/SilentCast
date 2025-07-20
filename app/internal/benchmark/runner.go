package benchmark

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/SphereStacking/silentcast/internal/action"
	"github.com/SphereStacking/silentcast/internal/config"
	"github.com/SphereStacking/silentcast/internal/hotkey"
	"github.com/SphereStacking/silentcast/internal/notify"
)

// BenchmarkResult holds the results of a benchmark run
type BenchmarkResult struct {
	Name        string
	Iterations  int
	Duration    time.Duration
	MinTime     time.Duration
	MaxTime     time.Duration
	MeanTime    time.Duration
	MedianTime  time.Duration
	MemoryUsage int64
}

// SystemInfo holds system information for context
type SystemInfo struct {
	OS           string
	Architecture string
	GoVersion    string
	NumCPU       int
	MemoryMB     int64
}

// BenchmarkRunner manages performance benchmarks
type BenchmarkRunner struct {
	configPath  string
	systemInfo  SystemInfo
	results     []BenchmarkResult
}

// NewBenchmarkRunner creates a new benchmark runner
func NewBenchmarkRunner(configPath string) *BenchmarkRunner {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return &BenchmarkRunner{
		configPath: configPath,
		systemInfo: SystemInfo{
			OS:           runtime.GOOS,
			Architecture: runtime.GOARCH,
			GoVersion:    runtime.Version(),
			NumCPU:       runtime.NumCPU(),
			MemoryMB:     int64(m.Sys / 1024 / 1024),
		},
		results: make([]BenchmarkResult, 0),
	}
}

// RunAll executes all benchmark tests
func (br *BenchmarkRunner) RunAll() error {
	fmt.Println("🏃 SilentCast Performance Benchmarks")
	fmt.Println("====================================")
	fmt.Println()
	
	br.printSystemInfo()
	fmt.Println()

	// Run startup benchmark
	if err := br.benchmarkStartup(); err != nil {
		return fmt.Errorf("startup benchmark failed: %w", err)
	}

	// Run configuration loading benchmark
	if err := br.benchmarkConfigLoading(); err != nil {
		return fmt.Errorf("config loading benchmark failed: %w", err)
	}

	// Run hotkey detection benchmark
	if err := br.benchmarkHotkeyDetection(); err != nil {
		return fmt.Errorf("hotkey detection benchmark failed: %w", err)
	}

	// Run action execution benchmark
	if err := br.benchmarkActionExecution(); err != nil {
		return fmt.Errorf("action execution benchmark failed: %w", err)
	}

	// Run memory usage benchmark
	if err := br.benchmarkMemoryUsage(); err != nil {
		return fmt.Errorf("memory usage benchmark failed: %w", err)
	}

	// Print final report
	br.printFinalReport()

	return nil
}

// printSystemInfo prints system information
func (br *BenchmarkRunner) printSystemInfo() {
	fmt.Printf("System Information:\n")
	fmt.Printf("  OS:           %s\n", br.systemInfo.OS)
	fmt.Printf("  Architecture: %s\n", br.systemInfo.Architecture)
	fmt.Printf("  Go Version:   %s\n", br.systemInfo.GoVersion)
	fmt.Printf("  CPU Cores:    %d\n", br.systemInfo.NumCPU)
	fmt.Printf("  System Memory: %dMB\n", br.systemInfo.MemoryMB)
}

// benchmarkStartup measures startup time
func (br *BenchmarkRunner) benchmarkStartup() error {
	fmt.Println("📈 Benchmarking startup time...")
	
	const iterations = 10
	times := make([]time.Duration, iterations)

	for i := 0; i < iterations; i++ {
		start := time.Now()
		
		// Load configuration
		loader := config.NewLoader(br.configPath)
		cfg, err := loader.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Initialize action manager
		actionManager := action.NewManager(cfg.Actions)
		_ = actionManager

		// Initialize notification manager
		notifyManager := notify.NewManager()
		_ = notifyManager

		// Initialize hotkey manager (stub version)
		hotkeyManager := hotkey.NewMockManager()
		_ = hotkeyManager

		times[i] = time.Since(start)
	}

	result := br.calculateStats("Startup Time", times)
	br.results = append(br.results, result)
	br.printResult(result)

	return nil
}

// benchmarkConfigLoading measures configuration loading performance
func (br *BenchmarkRunner) benchmarkConfigLoading() error {
	fmt.Println("📈 Benchmarking configuration loading...")
	
	const iterations = 20
	times := make([]time.Duration, iterations)

	for i := 0; i < iterations; i++ {
		start := time.Now()
		
		loader := config.NewLoader(br.configPath)
		_, err := loader.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		times[i] = time.Since(start)
	}

	result := br.calculateStats("Config Loading", times)
	br.results = append(br.results, result)
	br.printResult(result)

	return nil
}

// benchmarkHotkeyDetection measures hotkey detection performance
func (br *BenchmarkRunner) benchmarkHotkeyDetection() error {
	fmt.Println("📈 Benchmarking hotkey detection...")
	
	const iterations = 100
	times := make([]time.Duration, iterations)

	manager := hotkey.NewMockManager()

	// Register some test sequences
	testSequences := []string{"e", "t", "g,s", "ctrl+c", "alt+space"}
	for _, seq := range testSequences {
		if err := manager.Register(seq, "test_action"); err != nil {
			return fmt.Errorf("failed to register sequence %s: %w", seq, err)
		}
	}

	for i := 0; i < iterations; i++ {
		start := time.Now()
		
		// Simulate key parsing and detection by triggering key simulation
		sequence := testSequences[i%len(testSequences)]
		_ = manager.SimulateKeyPress(sequence)

		times[i] = time.Since(start)
	}

	result := br.calculateStats("Hotkey Detection", times)
	br.results = append(br.results, result)
	br.printResult(result)

	return nil
}

// benchmarkActionExecution measures action execution performance
func (br *BenchmarkRunner) benchmarkActionExecution() error {
	fmt.Println("📈 Benchmarking action execution...")
	
	const iterations = 10
	times := make([]time.Duration, iterations)

	// Load configuration to get real actions
	loader := config.NewLoader(br.configPath)
	cfg, err := loader.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	actionManager := action.NewManager(cfg.Actions)

	// Create a simple test action if no actions are configured
	if len(cfg.Actions) == 0 {
		testActions := map[string]config.ActionConfig{
			"test_echo": {
				Type:        "script",
				Command:     "echo",
				Args:        []string{"benchmark test"},
				Description: "Test action for benchmarking",
			},
		}
		actionManager = action.NewManager(testActions)
	}

	// Get the first action for testing
	var testAction string
	for name := range cfg.Actions {
		testAction = name
		break
	}
	if testAction == "" {
		testAction = "test_echo"
	}

	for i := 0; i < iterations; i++ {
		start := time.Now()
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_ = actionManager.Execute(ctx, testAction)
		cancel()

		times[i] = time.Since(start)
	}

	result := br.calculateStats("Action Execution", times)
	br.results = append(br.results, result)
	br.printResult(result)

	return nil
}

// benchmarkMemoryUsage measures memory usage patterns
func (br *BenchmarkRunner) benchmarkMemoryUsage() error {
	fmt.Println("📈 Benchmarking memory usage...")

	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	// Simulate application lifecycle
	const iterations = 5
	for i := 0; i < iterations; i++ {
		// Load configuration
		loader := config.NewLoader(br.configPath)
		cfg, err := loader.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Initialize managers
		actionManager := action.NewManager(cfg.Actions)
		notifyManager := notify.NewManager()
		hotkeyManager := hotkey.NewMockManager()

		// Register some hotkeys
		for j := 0; j < 10; j++ {
			_ = hotkeyManager.Register(fmt.Sprintf("test_%d", j), "test_action")
		}

		// Execute some operations
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		_ = notifyManager.Info(ctx, "Test", "Memory benchmark")
		cancel()

		// Cleanup
		_ = actionManager
		_ = notifyManager
		_ = hotkeyManager
	}

	runtime.GC()
	runtime.ReadMemStats(&m2)

	memoryUsed := int64(m2.TotalAlloc - m1.TotalAlloc)
	result := BenchmarkResult{
		Name:        "Memory Usage",
		Iterations:  iterations,
		MemoryUsage: memoryUsed,
	}

	br.results = append(br.results, result)
	
	fmt.Printf("  Memory allocated: %d bytes (%.2f KB)\n", memoryUsed, float64(memoryUsed)/1024)
	fmt.Printf("  Memory per operation: %d bytes\n", memoryUsed/int64(iterations))
	fmt.Println()

	return nil
}

// calculateStats calculates statistical measures from timing data
func (br *BenchmarkRunner) calculateStats(name string, times []time.Duration) BenchmarkResult {
	if len(times) == 0 {
		return BenchmarkResult{Name: name}
	}

	// Sort for median calculation
	sortedTimes := make([]time.Duration, len(times))
	copy(sortedTimes, times)
	for i := 0; i < len(sortedTimes); i++ {
		for j := i + 1; j < len(sortedTimes); j++ {
			if sortedTimes[i] > sortedTimes[j] {
				sortedTimes[i], sortedTimes[j] = sortedTimes[j], sortedTimes[i]
			}
		}
	}

	// Calculate statistics
	var total time.Duration
	minTime := times[0]
	maxTime := times[0]

	for _, t := range times {
		total += t
		if t < minTime {
			minTime = t
		}
		if t > maxTime {
			maxTime = t
		}
	}

	meanTime := total / time.Duration(len(times))
	medianTime := sortedTimes[len(sortedTimes)/2]

	return BenchmarkResult{
		Name:       name,
		Iterations: len(times),
		Duration:   total,
		MinTime:    minTime,
		MaxTime:    maxTime,
		MeanTime:   meanTime,
		MedianTime: medianTime,
	}
}

// printResult prints a benchmark result
func (br *BenchmarkRunner) printResult(result BenchmarkResult) {
	fmt.Printf("  Iterations: %d\n", result.Iterations)
	fmt.Printf("  Mean time:  %v\n", result.MeanTime)
	fmt.Printf("  Median:     %v\n", result.MedianTime)
	fmt.Printf("  Min:        %v\n", result.MinTime)
	fmt.Printf("  Max:        %v\n", result.MaxTime)
	fmt.Println()
}

// printFinalReport prints a summary of all benchmarks
func (br *BenchmarkRunner) printFinalReport() {
	fmt.Println("📊 Final Benchmark Report")
	fmt.Println("=========================")
	fmt.Println()

	if len(br.results) == 0 {
		fmt.Println("No benchmark results to display.")
		return
	}

	// Find the longest name for alignment
	maxNameLen := 0
	for _, result := range br.results {
		if len(result.Name) > maxNameLen {
			maxNameLen = len(result.Name)
		}
	}

	fmt.Printf("%-*s | %10s | %10s | %10s | %10s\n", 
		maxNameLen, "Benchmark", "Mean", "Median", "Min", "Max")
	fmt.Printf("%s\n", strings.Repeat("-", maxNameLen+50))

	for _, result := range br.results {
		if result.Name == "Memory Usage" {
			fmt.Printf("%-*s | %10s | %10s | %10s | %10s\n", 
				maxNameLen, result.Name, 
				fmt.Sprintf("%d bytes", result.MemoryUsage),
				"-", "-", "-")
		} else {
			fmt.Printf("%-*s | %10v | %10v | %10v | %10v\n", 
				maxNameLen, result.Name, 
				result.MeanTime, result.MedianTime, result.MinTime, result.MaxTime)
		}
	}

	fmt.Println()
	fmt.Println("💡 Performance Tips:")
	
	// Analyze results and provide tips
	for _, result := range br.results {
		switch result.Name {
		case "Startup Time":
			if result.MeanTime > 200*time.Millisecond {
				fmt.Println("  ⚠️  Startup time is high - consider lazy loading components")
			} else {
				fmt.Println("  ✅ Startup time is good")
			}
		case "Config Loading":
			if result.MeanTime > 50*time.Millisecond {
				fmt.Println("  ⚠️  Config loading is slow - consider caching")
			} else {
				fmt.Println("  ✅ Config loading performance is good")
			}
		case "Action Execution":
			if result.MeanTime > 1*time.Second {
				fmt.Println("  ⚠️  Action execution is slow - check command performance")
			} else {
				fmt.Println("  ✅ Action execution performance is good")
			}
		case "Memory Usage":
			if result.MemoryUsage > 10*1024*1024 { // 10MB
				fmt.Println("  ⚠️  High memory usage detected - check for memory leaks")
			} else {
				fmt.Println("  ✅ Memory usage is reasonable")
			}
		}
	}

	fmt.Println()
	fmt.Printf("Benchmark completed at %s\n", time.Now().Format("2006-01-02 15:04:05"))
}

// RunBenchmark runs a benchmark and returns the result
func RunBenchmark(configPath string) error {
	runner := NewBenchmarkRunner(configPath)
	
	// Ensure we have a valid configuration for benchmarking
	if err := runner.ensureTestConfig(); err != nil {
		return fmt.Errorf("failed to setup test config: %w", err)
	}
	
	return runner.RunAll()
}

// ensureTestConfig creates a basic test configuration if none exists
func (br *BenchmarkRunner) ensureTestConfig() error {
	loader := config.NewLoader(br.configPath)
	_, err := loader.Load()
	if err != nil {
		// Create a temporary test configuration
		testConfig := `daemon:
  auto_start: false
  log_level: error

hotkeys:
  prefix: "alt+space"
  timeout: 1000

spells:
  "e": "editor"
  "t": "terminal"
  "g,s": "git_status"
  "b": "browser"

grimoire:
  editor:
    type: script
    command: "echo"
    args: ["launching editor"]
    description: "Open text editor"
  
  terminal:
    type: script
    command: "echo"
    args: ["opening terminal"]
    description: "Open terminal"
  
  git_status:
    type: script
    command: "echo"
    args: ["git status"]
    description: "Show git status"
  
  browser:
    type: script
    command: "echo"
    args: ["opening browser"]
    description: "Open web browser"
`
		
		// Write temporary config
		tempConfigPath := filepath.Join(br.configPath, "spellbook.yml")
		if err := os.MkdirAll(br.configPath, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
		
		if err := os.WriteFile(tempConfigPath, []byte(testConfig), 0644); err != nil {
			return fmt.Errorf("failed to write test config: %w", err)
		}
		
		fmt.Printf("ℹ️  Created temporary test configuration at %s\n", tempConfigPath)
		fmt.Printf("   (This is only used for benchmarking purposes)\n\n")
	}
	
	return nil
}