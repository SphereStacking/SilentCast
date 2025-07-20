package benchmarks

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"testing"
	"time"
)

// ProfilingOptions configures profiling behavior
type ProfilingOptions struct {
	CPUProfile    bool
	MemProfile    bool
	TraceProfile  bool
	GoroutineProf bool
	BlockProfile  bool
	MutexProfile  bool
	ProfileDir    string
}

// DefaultProfilingOptions returns default profiling options
func DefaultProfilingOptions() ProfilingOptions {
	return ProfilingOptions{
		CPUProfile:    true,
		MemProfile:    true,
		TraceProfile:  false, // Trace is more expensive
		GoroutineProf: true,
		BlockProfile:  false, // Block profiling has overhead
		MutexProfile:  false, // Mutex profiling has overhead
		ProfileDir:    "profiles",
	}
}

// ProfileManager manages profiling for benchmarks
type ProfileManager struct {
	opts       ProfilingOptions
	cpuFile    *os.File
	traceFile  *os.File
	profileDir string
	testName   string
}

// NewProfileManager creates a new profile manager
func NewProfileManager(testName string, opts ProfilingOptions) (*ProfileManager, error) {
	// Create profile directory
	profileDir := filepath.Join(opts.ProfileDir, time.Now().Format("2006-01-02-15-04-05"))
	if err := os.MkdirAll(profileDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create profile directory: %w", err)
	}

	return &ProfileManager{
		opts:       opts,
		profileDir: profileDir,
		testName:   testName,
	}, nil
}

// Start begins profiling
func (pm *ProfileManager) Start(b *testing.B) error {
	// Enable mutex and block profiling if requested
	if pm.opts.MutexProfile {
		runtime.SetMutexProfileFraction(1)
	}
	if pm.opts.BlockProfile {
		runtime.SetBlockProfileRate(1)
	}

	// Start CPU profiling
	if pm.opts.CPUProfile {
		cpuPath := filepath.Join(pm.profileDir, fmt.Sprintf("%s-cpu.prof", pm.testName))
		cpuFile, err := os.Create(cpuPath)
		if err != nil {
			return fmt.Errorf("failed to create CPU profile: %w", err)
		}
		pm.cpuFile = cpuFile
		
		if err := pprof.StartCPUProfile(cpuFile); err != nil {
			cpuFile.Close()
			return fmt.Errorf("failed to start CPU profile: %w", err)
		}
		b.Logf("CPU profiling started: %s", cpuPath)
	}

	// Start trace profiling
	if pm.opts.TraceProfile {
		tracePath := filepath.Join(pm.profileDir, fmt.Sprintf("%s-trace.out", pm.testName))
		traceFile, err := os.Create(tracePath)
		if err != nil {
			return fmt.Errorf("failed to create trace file: %w", err)
		}
		pm.traceFile = traceFile
		
		if err := trace.Start(traceFile); err != nil {
			traceFile.Close()
			return fmt.Errorf("failed to start trace: %w", err)
		}
		b.Logf("Trace profiling started: %s", tracePath)
	}

	return nil
}

// Stop stops profiling and writes profiles
func (pm *ProfileManager) Stop(b *testing.B) error {
	// Stop CPU profiling
	if pm.opts.CPUProfile && pm.cpuFile != nil {
		pprof.StopCPUProfile()
		pm.cpuFile.Close()
		b.Logf("CPU profile saved")
	}

	// Stop trace profiling
	if pm.opts.TraceProfile && pm.traceFile != nil {
		trace.Stop()
		pm.traceFile.Close()
		b.Logf("Trace profile saved")
	}

	// Write memory profile
	if pm.opts.MemProfile {
		memPath := filepath.Join(pm.profileDir, fmt.Sprintf("%s-mem.prof", pm.testName))
		memFile, err := os.Create(memPath)
		if err != nil {
			return fmt.Errorf("failed to create memory profile: %w", err)
		}
		defer memFile.Close()

		runtime.GC() // Force GC for accurate memory profile
		if err := pprof.WriteHeapProfile(memFile); err != nil {
			return fmt.Errorf("failed to write memory profile: %w", err)
		}
		b.Logf("Memory profile saved: %s", memPath)
	}

	// Write goroutine profile
	if pm.opts.GoroutineProf {
		goroutinePath := filepath.Join(pm.profileDir, fmt.Sprintf("%s-goroutine.prof", pm.testName))
		goroutineFile, err := os.Create(goroutinePath)
		if err != nil {
			return fmt.Errorf("failed to create goroutine profile: %w", err)
		}
		defer goroutineFile.Close()

		if err := pprof.Lookup("goroutine").WriteTo(goroutineFile, 0); err != nil {
			return fmt.Errorf("failed to write goroutine profile: %w", err)
		}
		b.Logf("Goroutine profile saved: %s", goroutinePath)
	}

	// Write block profile
	if pm.opts.BlockProfile {
		blockPath := filepath.Join(pm.profileDir, fmt.Sprintf("%s-block.prof", pm.testName))
		blockFile, err := os.Create(blockPath)
		if err != nil {
			return fmt.Errorf("failed to create block profile: %w", err)
		}
		defer blockFile.Close()

		if err := pprof.Lookup("block").WriteTo(blockFile, 0); err != nil {
			return fmt.Errorf("failed to write block profile: %w", err)
		}
		b.Logf("Block profile saved: %s", blockPath)
		runtime.SetBlockProfileRate(0) // Disable for next test
	}

	// Write mutex profile
	if pm.opts.MutexProfile {
		mutexPath := filepath.Join(pm.profileDir, fmt.Sprintf("%s-mutex.prof", pm.testName))
		mutexFile, err := os.Create(mutexPath)
		if err != nil {
			return fmt.Errorf("failed to create mutex profile: %w", err)
		}
		defer mutexFile.Close()

		if err := pprof.Lookup("mutex").WriteTo(mutexFile, 0); err != nil {
			return fmt.Errorf("failed to write mutex profile: %w", err)
		}
		b.Logf("Mutex profile saved: %s", mutexPath)
		runtime.SetMutexProfileFraction(0) // Disable for next test
	}

	b.Logf("All profiles saved to: %s", pm.profileDir)
	return nil
}

// RunWithProfiling runs a benchmark with profiling enabled
func RunWithProfiling(b *testing.B, name string, fn func(b *testing.B)) {
	// Check if profiling is requested via environment variable
	if os.Getenv("ENABLE_PROFILING") != "1" {
		// Run without profiling
		fn(b)
		return
	}

	opts := DefaultProfilingOptions()
	
	// Allow customization via environment variables
	if os.Getenv("PROFILE_TRACE") == "1" {
		opts.TraceProfile = true
	}
	if os.Getenv("PROFILE_BLOCK") == "1" {
		opts.BlockProfile = true
	}
	if os.Getenv("PROFILE_MUTEX") == "1" {
		opts.MutexProfile = true
	}

	pm, err := NewProfileManager(name, opts)
	if err != nil {
		b.Fatalf("Failed to create profile manager: %v", err)
	}

	if err := pm.Start(b); err != nil {
		b.Fatalf("Failed to start profiling: %v", err)
	}

	// Run the actual benchmark
	fn(b)

	if err := pm.Stop(b); err != nil {
		b.Errorf("Failed to stop profiling: %v", err)
	}
}

// MemorySnapshot captures memory statistics at a point in time
type MemorySnapshot struct {
	Timestamp    time.Time
	Alloc        uint64
	TotalAlloc   uint64
	Sys          uint64
	Mallocs      uint64
	Frees        uint64
	HeapAlloc    uint64
	HeapSys      uint64
	HeapIdle     uint64
	HeapInuse    uint64
	HeapReleased uint64
	HeapObjects  uint64
	StackInuse   uint64
	StackSys     uint64
	NumGC        uint32
	NumGoroutine int
}

// CaptureMemorySnapshot captures current memory statistics
func CaptureMemorySnapshot() MemorySnapshot {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	return MemorySnapshot{
		Timestamp:    time.Now(),
		Alloc:        m.Alloc,
		TotalAlloc:   m.TotalAlloc,
		Sys:          m.Sys,
		Mallocs:      m.Mallocs,
		Frees:        m.Frees,
		HeapAlloc:    m.HeapAlloc,
		HeapSys:      m.HeapSys,
		HeapIdle:     m.HeapIdle,
		HeapInuse:    m.HeapInuse,
		HeapReleased: m.HeapReleased,
		HeapObjects:  m.HeapObjects,
		StackInuse:   m.StackInuse,
		StackSys:     m.StackSys,
		NumGC:        m.NumGC,
		NumGoroutine: runtime.NumGoroutine(),
	}
}

// CompareSnapshots compares two memory snapshots and returns the differences
func CompareSnapshots(before, after MemorySnapshot) string {
	return fmt.Sprintf(`Memory Usage Change:
  Heap Alloc:    %+d bytes (%.2f MB)
  Total Alloc:   %+d bytes (%.2f MB)
  Heap Objects:  %+d
  Goroutines:    %+d
  GC Runs:       %+d
  Duration:      %v`,
		int64(after.HeapAlloc-before.HeapAlloc),
		float64(after.HeapAlloc-before.HeapAlloc)/1024/1024,
		int64(after.TotalAlloc-before.TotalAlloc),
		float64(after.TotalAlloc-before.TotalAlloc)/1024/1024,
		int64(after.HeapObjects-before.HeapObjects),
		after.NumGoroutine-before.NumGoroutine,
		after.NumGC-before.NumGC,
		after.Timestamp.Sub(before.Timestamp),
	)
}

// CheckGoroutineLeaks checks for goroutine leaks
func CheckGoroutineLeaks(b *testing.B, baseline int, tolerance int) {
	// Give goroutines time to clean up
	time.Sleep(100 * time.Millisecond)
	runtime.GC()
	
	current := runtime.NumGoroutine()
	leaked := current - baseline
	
	if leaked > tolerance {
		b.Errorf("Goroutine leak detected: baseline=%d, current=%d, leaked=%d (tolerance=%d)",
			baseline, current, leaked, tolerance)
		
		// Print goroutine stack traces for debugging
		buf := make([]byte, 1<<20)
		stackLen := runtime.Stack(buf, true)
		b.Logf("Goroutine stack traces:\n%s", buf[:stackLen])
	}
}