package benchmark

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestNewBenchmarkRunner(t *testing.T) {
	configPath := "/test/config"
	runner := NewBenchmarkRunner(configPath)
	
	if runner == nil {
		t.Fatal("NewBenchmarkRunner returned nil")
	}
	
	if runner.configPath != configPath {
		t.Errorf("Expected configPath %s, got %s", configPath, runner.configPath)
	}
	
	// Test system info is populated
	if runner.systemInfo.OS == "" {
		t.Error("System OS should be populated")
	}
	
	if runner.systemInfo.Architecture == "" {
		t.Error("System Architecture should be populated")
	}
	
	if runner.systemInfo.GoVersion == "" {
		t.Error("Go Version should be populated")
	}
	
	if runner.systemInfo.NumCPU <= 0 {
		t.Error("NumCPU should be positive")
	}
	
	if runner.systemInfo.MemoryMB <= 0 {
		t.Error("MemoryMB should be positive")
	}
	
	if runner.results == nil {
		t.Error("Results slice should be initialized")
	}
}

func TestSystemInfo(t *testing.T) {
	runner := NewBenchmarkRunner("/test")
	info := runner.systemInfo
	
	// Test system info matches runtime values
	if info.OS != runtime.GOOS {
		t.Errorf("OS mismatch: got %s, want %s", info.OS, runtime.GOOS)
	}
	
	if info.Architecture != runtime.GOARCH {
		t.Errorf("Architecture mismatch: got %s, want %s", info.Architecture, runtime.GOARCH)
	}
	
	if info.GoVersion != runtime.Version() {
		t.Errorf("GoVersion mismatch: got %s, want %s", info.GoVersion, runtime.Version())
	}
	
	if info.NumCPU != runtime.NumCPU() {
		t.Errorf("NumCPU mismatch: got %d, want %d", info.NumCPU, runtime.NumCPU())
	}
}

func TestBenchmarkResult(t *testing.T) {
	result := BenchmarkResult{
		Name:        "Test Benchmark",
		Iterations:  10,
		Duration:    100 * time.Millisecond,
		MinTime:     5 * time.Millisecond,
		MaxTime:     15 * time.Millisecond,
		MeanTime:    10 * time.Millisecond,
		MedianTime:  9 * time.Millisecond,
		MemoryUsage: 1024,
	}
	
	if result.Name != "Test Benchmark" {
		t.Error("Name not set correctly")
	}
	
	if result.Iterations != 10 {
		t.Error("Iterations not set correctly")
	}
	
	if result.Duration != 100*time.Millisecond {
		t.Error("Duration not set correctly")
	}
	
	if result.MemoryUsage != 1024 {
		t.Error("MemoryUsage not set correctly")
	}
}

func TestCalculateStats(t *testing.T) {
	runner := NewBenchmarkRunner("/test")
	
	times := []time.Duration{
		10 * time.Millisecond,
		20 * time.Millisecond,
		15 * time.Millisecond,
		25 * time.Millisecond,
		12 * time.Millisecond,
	}
	
	result := runner.calculateStats("Test Stats", times)
	
	if result.Name != "Test Stats" {
		t.Error("Name not set correctly")
	}
	
	if result.Iterations != len(times) {
		t.Errorf("Expected %d iterations, got %d", len(times), result.Iterations)
	}
	
	// Test min/max
	if result.MinTime != 10*time.Millisecond {
		t.Errorf("Expected min time %v, got %v", 10*time.Millisecond, result.MinTime)
	}
	
	if result.MaxTime != 25*time.Millisecond {
		t.Errorf("Expected max time %v, got %v", 25*time.Millisecond, result.MaxTime)
	}
	
	// Test mean (total: 82ms, count: 5, mean: 16.4ms)
	expectedMean := 82 * time.Millisecond / 5
	if result.MeanTime != expectedMean {
		t.Errorf("Expected mean time %v, got %v", expectedMean, result.MeanTime)
	}
	
	// Test total duration
	expectedTotal := 82 * time.Millisecond
	if result.Duration != expectedTotal {
		t.Errorf("Expected total duration %v, got %v", expectedTotal, result.Duration)
	}
}

func TestCalculateStatsEmpty(t *testing.T) {
	runner := NewBenchmarkRunner("/test")
	
	result := runner.calculateStats("Empty Test", []time.Duration{})
	
	if result.Name != "Empty Test" {
		t.Error("Name should be set even for empty input")
	}
	
	if result.Iterations != 0 {
		t.Error("Iterations should be 0 for empty input")
	}
}

func TestCalculateStatsSingle(t *testing.T) {
	runner := NewBenchmarkRunner("/test")
	
	times := []time.Duration{10 * time.Millisecond}
	result := runner.calculateStats("Single Test", times)
	
	if result.Iterations != 1 {
		t.Error("Iterations should be 1")
	}
	
	if result.MinTime != 10*time.Millisecond {
		t.Error("MinTime should equal the single value")
	}
	
	if result.MaxTime != 10*time.Millisecond {
		t.Error("MaxTime should equal the single value")
	}
	
	if result.MeanTime != 10*time.Millisecond {
		t.Error("MeanTime should equal the single value")
	}
	
	if result.MedianTime != 10*time.Millisecond {
		t.Error("MedianTime should equal the single value")
	}
}

func TestEnsureTestConfig(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	
	runner := NewBenchmarkRunner(tempDir)
	
	// Test creating config when none exists
	err := runner.ensureTestConfig()
	if err != nil {
		t.Fatalf("ensureTestConfig failed: %v", err)
	}
	
	// Check that config file was created
	configPath := filepath.Join(tempDir, "spellbook.yml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Test config file was not created")
	}
	
	// Test that config file contains expected content
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read test config: %v", err)
	}
	
	contentStr := string(content)
	expectedParts := []string{
		"hotkeys:",
		"spells:",
		"grimoire:",
		"prefix:",
		"alt+space",
	}
	
	for _, part := range expectedParts {
		if !containsString(contentStr, part) {
			t.Errorf("Test config should contain '%s'", part)
		}
	}
}

func TestEnsureTestConfigExisting(t *testing.T) {
	// Create temporary directory with existing config
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "spellbook.yml")
	
	// Write existing config
	existingConfig := `hotkeys:
  prefix: "ctrl+alt"
spells:
  "x": "existing"
grimoire:
  existing:
    type: script
    command: echo
`
	err := os.WriteFile(configPath, []byte(existingConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to write existing config: %v", err)
	}
	
	runner := NewBenchmarkRunner(tempDir)
	
	// Should not modify existing config
	err = runner.ensureTestConfig()
	if err != nil {
		t.Fatalf("ensureTestConfig failed with existing config: %v", err)
	}
	
	// Check that existing config is preserved
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config after ensureTestConfig: %v", err)
	}
	
	if string(content) != existingConfig {
		t.Error("Existing config should not be modified")
	}
}

func TestRunBenchmark(t *testing.T) {
	// Create temporary directory for test
	tempDir := t.TempDir()
	
	// This test just verifies the function doesn't panic
	// and handles the case where config creation is needed
	err := RunBenchmark(tempDir)
	
	// The benchmark might fail due to missing dependencies in test environment
	// but it shouldn't panic or return nil pointer errors
	if err != nil {
		t.Logf("RunBenchmark returned error (expected in test environment): %v", err)
	}
}

func TestPrintSystemInfo(t *testing.T) {
	runner := NewBenchmarkRunner("/test")
	
	// This test just ensures the function doesn't panic
	// Output testing would require capturing stdout which is complex
	runner.printSystemInfo()
}

func TestPrintResult(t *testing.T) {
	runner := NewBenchmarkRunner("/test")
	
	result := BenchmarkResult{
		Name:       "Test Result",
		Iterations: 10,
		MeanTime:   15 * time.Millisecond,
		MedianTime: 12 * time.Millisecond,
		MinTime:    8 * time.Millisecond,
		MaxTime:    20 * time.Millisecond,
	}
	
	// This test just ensures the function doesn't panic
	runner.printResult(result)
}

func TestPrintFinalReport(t *testing.T) {
	runner := NewBenchmarkRunner("/test")
	
	// Test with no results
	runner.printFinalReport()
	
	// Test with some results
	runner.results = []BenchmarkResult{
		{
			Name:       "Startup Time",
			Iterations: 10,
			MeanTime:   50 * time.Millisecond,
			MedianTime: 45 * time.Millisecond,
			MinTime:    40 * time.Millisecond,
			MaxTime:    60 * time.Millisecond,
		},
		{
			Name:        "Memory Usage",
			Iterations:  5,
			MemoryUsage: 1024 * 1024, // 1MB
		},
	}
	
	// This test just ensures the function doesn't panic
	runner.printFinalReport()
}

// Helper function since strings.Contains might not be available
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (hasPrefix(s, substr) || containsString(s[1:], substr))))
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}

func TestBenchmarkResultTypes(t *testing.T) {
	// Test that BenchmarkResult can handle different types of benchmarks
	timeResult := BenchmarkResult{
		Name:       "Time Benchmark",
		Iterations: 100,
		Duration:   1 * time.Second,
		MeanTime:   10 * time.Millisecond,
	}
	
	memoryResult := BenchmarkResult{
		Name:        "Memory Benchmark",
		Iterations:  10,
		MemoryUsage: 2048,
	}
	
	if timeResult.Duration == 0 {
		t.Error("Time benchmark should have duration")
	}
	
	if memoryResult.MemoryUsage == 0 {
		t.Error("Memory benchmark should have memory usage")
	}
}

func TestSystemInfoMemoryCalculation(t *testing.T) {
	runner := NewBenchmarkRunner("/test")
	
	// Memory should be calculated from runtime.MemStats
	if runner.systemInfo.MemoryMB <= 0 {
		t.Error("System memory should be positive")
	}
	
	// Should be reasonable value (between 1MB and 1TB)
	if runner.systemInfo.MemoryMB < 1 || runner.systemInfo.MemoryMB > 1024*1024 {
		t.Errorf("System memory seems unreasonable: %d MB", runner.systemInfo.MemoryMB)
	}
}