package benchmarks

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/SphereStacking/silentcast/internal/config"
	"github.com/SphereStacking/silentcast/internal/action"
	"github.com/SphereStacking/silentcast/internal/hotkey"
	"github.com/SphereStacking/silentcast/internal/notify"
	"github.com/SphereStacking/silentcast/internal/output"
	"github.com/SphereStacking/silentcast/internal/permission"
)

// BenchmarkEnvironment provides a standardized environment for performance testing
type BenchmarkEnvironment struct {
	TempDir       string
	ConfigDir     string
	Config        *config.Config
	ConfigLoader  *config.Loader
	ActionManager *action.Manager
	HotkeyManager hotkey.Manager
	NotifyManager *notify.Manager
	OutputManager output.Manager
	PermManager   permission.Manager
	ctx           context.Context
	cancel        context.CancelFunc
}

// SetupBenchmarkEnvironment creates a benchmark environment with all components
func SetupBenchmarkEnvironment(b *testing.B) *BenchmarkEnvironment {
	tempDir, err := os.MkdirTemp("", "silentcast-benchmark-*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	
	configDir := filepath.Join(tempDir, "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		b.Fatalf("Failed to create config dir: %v", err)
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	
	env := &BenchmarkEnvironment{
		TempDir:   tempDir,
		ConfigDir: configDir,
		ctx:       ctx,
		cancel:    cancel,
	}
	
	// Setup cleanup
	b.Cleanup(func() {
		env.Cleanup()
	})
	
	return env
}

// WriteConfigFile writes a configuration file to the benchmark environment
func (e *BenchmarkEnvironment) WriteConfigFile(filename, content string) error {
	configPath := filepath.Join(e.ConfigDir, filename)
	return os.WriteFile(configPath, []byte(content), 0644)
}

// LoadConfig loads configuration from the benchmark environment
func (e *BenchmarkEnvironment) LoadConfig() error {
	e.ConfigLoader = config.NewLoader(e.ConfigDir)
	cfg, err := e.ConfigLoader.Load()
	if err != nil {
		return err
	}
	e.Config = cfg
	return nil
}

// InitializeComponents initializes all components with the loaded configuration
func (e *BenchmarkEnvironment) InitializeComponents() error {
	if e.Config == nil {
		return fmt.Errorf("configuration must be loaded before initializing components")
	}
	
	// Initialize output manager
	e.OutputManager = output.NewBufferedManager(output.DefaultOptions())
	
	// Initialize notification manager
	e.NotifyManager = notify.NewManager()
	
	// Initialize permission manager
	permManager, err := permission.NewManager()
	if err != nil {
		return err
	}
	e.PermManager = permManager
	
	// Initialize action manager
	e.ActionManager = action.NewManager(e.Config.Actions)
	
	// Initialize hotkey manager (using mock for benchmarking)
	e.HotkeyManager = hotkey.NewMockManager()
	
	return nil
}

// ExecuteAction executes an action by name
func (e *BenchmarkEnvironment) ExecuteAction(actionName string) error {
	return e.ActionManager.Execute(e.ctx, actionName)
}

// Cleanup cleans up the benchmark environment
func (e *BenchmarkEnvironment) Cleanup() {
	if e.cancel != nil {
		e.cancel()
	}
	if e.TempDir != "" {
		os.RemoveAll(e.TempDir)
	}
}

// GetBenchmarkConfig returns a standard configuration for benchmarking
func GetBenchmarkConfig() string {
	return `
daemon:
  auto_start: false
  log_level: error  # Reduce logging overhead for benchmarks
  config_watch: false

logger:
  level: error
  file: ""

hotkeys:
  prefix: "alt+space"
  timeout: 1000
  sequence_timeout: 2000

spells:
  # Single key spells
  e: "editor"
  t: "terminal"
  b: "browser"
  f: "file_manager"
  c: "calculator"
  
  # Multi-key sequences
  "g,s": "git_status"
  "g,p": "git_pull"
  "g,c": "git_commit"
  "d,l": "docker_logs"
  "d,s": "docker_status"
  
  # Complex sequences
  "s,y,s": "system_info"
  "d,e,v": "dev_tools"

grimoire:
  editor:
    type: app
    command: "echo"
    args: ["editor"]
    description: "Text editor"
    
  terminal:
    type: app
    command: "echo"
    args: ["terminal"]
    description: "Terminal"
    
  browser:
    type: url
    command: "https://example.com"
    description: "Web browser"
    
  file_manager:
    type: app
    command: "echo"
    args: ["file manager"]
    description: "File manager"
    
  calculator:
    type: app
    command: "echo"
    args: ["calculator"]
    description: "Calculator"
    
  git_status:
    type: script
    command: "echo"
    args: ["git status output"]
    show_output: true
    description: "Git status"
    
  git_pull:
    type: script
    command: "echo"
    args: ["git pull completed"]
    show_output: true
    description: "Git pull"
    
  git_commit:
    type: script
    command: "echo"
    args: ["commit created"]
    show_output: true
    description: "Git commit"
    
  docker_logs:
    type: script
    command: "echo"
    args: ["docker logs output"]
    show_output: true
    description: "Docker logs"
    
  docker_status:
    type: script
    command: "echo"
    args: ["docker containers running"]
    show_output: true
    description: "Docker status"
    
  system_info:
    type: script
    command: "echo"
    args: ["system information"]
    show_output: true
    description: "System information"
    
  dev_tools:
    type: app
    command: "echo"
    args: ["development tools"]
    description: "Development tools"
`
}

// GetLargeConfig returns a configuration with many entries for stress testing
func GetLargeConfig() string {
	base := `
daemon:
  auto_start: false
  log_level: error
  config_watch: false

hotkeys:
  prefix: "alt+space"
  timeout: 1000

spells:
`
	
	// Add many spell mappings
	for i := 0; i < 100; i++ {
		base += "  \"s" + string(rune('a'+i%26)) + string(rune('0'+i%10)) + "\": \"action_" + string(rune('a'+i%26)) + string(rune('0'+i%10)) + "\"\n"
	}
	
	base += "\ngrimoire:\n"
	
	// Add many actions
	for i := 0; i < 100; i++ {
		actionName := "action_" + string(rune('a'+i%26)) + string(rune('0'+i%10))
		base += "  " + actionName + ":\n"
		base += "    type: script\n"
		base += "    command: \"echo\"\n"
		base += "    args: [\"" + actionName + " output\"]\n"
		base += "    description: \"Benchmark action " + string(rune('a'+i%26)) + string(rune('0'+i%10)) + "\"\n"
	}
	
	return base
}

// GetMemoryStats returns current memory usage statistics
func GetMemoryStats() runtime.MemStats {
	var stats runtime.MemStats
	runtime.GC() // Force garbage collection for consistent measurements
	runtime.ReadMemStats(&stats)
	return stats
}

// BenchmarkResult represents the result of a benchmark run
type BenchmarkResult struct {
	Duration    time.Duration
	Allocations int64
	AllocBytes  int64
	Name        string
	Platform    string
}

// FormatBenchmarkResult formats a benchmark result for display
func FormatBenchmarkResult(result BenchmarkResult) string {
	return "Benchmark: " + result.Name + 
		" | Duration: " + result.Duration.String() + 
		" | Allocs: " + strconv.FormatInt(result.Allocations, 10) + 
		" | Memory: " + strconv.FormatInt(result.AllocBytes, 10) + " bytes" +
		" | Platform: " + result.Platform
}

// RunMemoryBenchmark runs a function and measures memory allocation
func RunMemoryBenchmark(b *testing.B, fn func()) {
	b.ReportAllocs()
	
	var startStats, endStats runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&startStats)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fn()
	}
	b.StopTimer()
	
	runtime.GC()
	runtime.ReadMemStats(&endStats)
	
	allocations := endStats.Mallocs - startStats.Mallocs
	allocBytes := endStats.TotalAlloc - startStats.TotalAlloc
	
	b.ReportMetric(float64(allocations)/float64(b.N), "allocs/op")
	b.ReportMetric(float64(allocBytes)/float64(b.N), "B/op")
}