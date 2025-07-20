# Performance Optimization Guide

This guide provides comprehensive information about SilentCast's performance characteristics and optimization strategies.

## Overview

SilentCast is designed to be a lightweight, responsive hotkey-driven task runner. Performance is critical since users expect instant response times when triggering spells.

## Performance Targets

| Component | Target | Measurement |
|-----------|--------|-------------|
| Startup Time | < 100ms | Application initialization |
| Spell Trigger | < 50ms | Hotkey detection to action start |
| Config Reload | < 200ms | Live configuration reload |
| Memory Usage | < 50MB | Idle state memory footprint |
| CPU Usage | < 1% | Idle state CPU usage |

## Benchmarking

### Running Benchmarks

```bash
# Run all performance benchmarks
make benchmark

# Run specific benchmark categories
make benchmark-startup
make benchmark-hotkey
make benchmark-action
make benchmark-config
make benchmark-memory
make benchmark-stress

# Generate profiling data
make benchmark-cpu    # CPU profiling
make benchmark-mem    # Memory profiling
```

### Benchmark Categories

#### 1. Startup Benchmarks
- Application initialization time
- Component setup overhead
- Configuration loading performance
- System resource allocation

#### 2. Hotkey Benchmarks
- Key sequence parsing performance
- Hotkey registration overhead
- Event processing latency
- Multi-key sequence handling

#### 3. Action Benchmarks
- Action execution startup time
- Different action type performance (app, script, URL)
- Argument processing overhead
- Environment variable expansion

#### 4. Configuration Benchmarks
- YAML parsing performance
- Configuration validation speed
- File watching overhead
- Configuration reload time

#### 5. Memory Benchmarks
- Memory allocation patterns
- Garbage collection impact
- Memory leak detection
- Long-running session stability

#### 6. Stress Tests
- High-frequency action execution
- Massive configuration files
- Concurrent user simulation
- Resource exhaustion scenarios

## Performance Analysis

### Profiling Tools

```bash
# CPU profiling
go test -bench=BenchmarkStartup -cpuprofile=cpu.prof ./test/benchmarks/
go tool pprof cpu.prof

# Memory profiling
go test -bench=BenchmarkMemory -memprofile=mem.prof ./test/benchmarks/
go tool pprof mem.prof

# Trace analysis
go test -bench=BenchmarkAction -trace=trace.out ./test/benchmarks/
go tool trace trace.out
```

### Key Metrics

#### Latency Metrics
- **P50**: Median response time
- **P95**: 95th percentile response time
- **P99**: 99th percentile response time

#### Throughput Metrics
- **Actions/second**: Sustainable action execution rate
- **Config reloads/minute**: Configuration reload frequency

#### Resource Metrics
- **Memory allocation rate**: Bytes allocated per operation
- **GC pressure**: Garbage collection frequency and duration
- **CPU utilization**: Processor usage patterns

## Optimization Strategies

### 1. Startup Optimization

#### Lazy Initialization
```go
// Good: Lazy component initialization
type Manager struct {
    notifier atomic.Value // Loaded on first use
}

func (m *Manager) getNotifier() Notifier {
    if n := m.notifier.Load(); n != nil {
        return n.(Notifier)
    }
    // Initialize only when needed
    n := createNotifier()
    m.notifier.Store(n)
    return n
}
```

#### Configuration Caching
```go
// Cache resolved paths to avoid repeated filesystem operations
type PathCache struct {
    cache sync.Map
}

func (c *PathCache) ResolvePath(app string) string {
    if path, ok := c.cache.Load(app); ok {
        return path.(string)
    }
    resolved := resolvePath(app)
    c.cache.Store(app, resolved)
    return resolved
}
```

### 2. Hotkey Performance

#### Efficient Key Parsing
```go
// Good: Pre-compile key sequences
type KeySequence struct {
    keys     []Key
    compiled bool
}

func (ks *KeySequence) Compile() error {
    if ks.compiled {
        return nil
    }
    // One-time compilation
    ks.keys = parseKeys(ks.raw)
    ks.compiled = true
    return nil
}
```

#### Minimal Event Processing
```go
// Good: Fast path for common operations
func (m *Manager) processKeyEvent(event KeyEvent) {
    // Quick rejection for irrelevant events
    if !event.IsRelevant() {
        return
    }
    // Process only when necessary
    m.handleRelevantEvent(event)
}
```

### 3. Action Execution

#### Process Pool
```go
// Use process pools for script execution
type ProcessPool struct {
    workers chan *exec.Cmd
}

func (p *ProcessPool) Execute(command string) error {
    select {
    case worker := <-p.workers:
        defer func() { p.workers <- worker }()
        return worker.Run()
    default:
        // Fallback to new process
        return exec.Command(command).Run()
    }
}
```

#### Resource Limits
```go
// Set resource limits for child processes
func configureProcess(cmd *exec.Cmd) {
    cmd.SysProcAttr = &syscall.SysProcAttr{
        Setpgid: true,
        // Platform-specific resource limits
    }
}
```

### 4. Memory Management

#### Object Pooling
```go
// Pool frequently used objects
var notificationPool = sync.Pool{
    New: func() interface{} {
        return &Notification{}
    },
}

func GetNotification() *Notification {
    return notificationPool.Get().(*Notification)
}

func PutNotification(n *Notification) {
    n.Reset()
    notificationPool.Put(n)
}
```

#### String Interning
```go
// Intern frequently used strings
type StringInterner struct {
    strings sync.Map
}

func (si *StringInterner) Intern(s string) string {
    if interned, ok := si.strings.Load(s); ok {
        return interned.(string)
    }
    si.strings.Store(s, s)
    return s
}
```

### 5. Configuration Optimization

#### Incremental Parsing
```go
// Parse only changed sections
type ConfigDiff struct {
    AddedActions   map[string]ActionConfig
    RemovedActions []string
    ModifiedActions map[string]ActionConfig
}

func (c *Config) ApplyDiff(diff ConfigDiff) {
    // Update only changed parts
    for name, action := range diff.AddedActions {
        c.Actions[name] = action
    }
    // ... handle other changes
}
```

#### Configuration Validation Caching
```go
// Cache validation results
type ValidationCache struct {
    cache map[string]ValidationResult
    mutex sync.RWMutex
}

func (vc *ValidationCache) Validate(config Config) ValidationResult {
    hash := config.Hash()
    vc.mutex.RLock()
    if result, ok := vc.cache[hash]; ok {
        vc.mutex.RUnlock()
        return result
    }
    vc.mutex.RUnlock()
    
    result := performValidation(config)
    vc.mutex.Lock()
    vc.cache[hash] = result
    vc.mutex.Unlock()
    return result
}
```

## Performance Monitoring

### Runtime Metrics

```go
// Collect runtime metrics
func collectMetrics() Metrics {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    return Metrics{
        HeapAlloc:    m.HeapAlloc,
        HeapObjects:  m.HeapObjects,
        NumGoroutine: runtime.NumGoroutine(),
        NumGC:        m.NumGC,
    }
}
```

### Application Metrics

```go
// Track application-specific metrics
type AppMetrics struct {
    ActionExecutions    counter
    ConfigReloads      counter
    AverageActionTime  histogram
    ErrorRate          gauge
}

func (m *AppMetrics) RecordActionExecution(duration time.Duration) {
    m.ActionExecutions.Inc()
    m.AverageActionTime.Observe(duration.Seconds())
}
```

## CI/CD Integration

### Automated Performance Testing

The CI/CD pipeline automatically runs performance benchmarks:

1. **Pull Request Analysis**: Compare performance against main branch
2. **Regression Detection**: Alert on significant performance degradations
3. **Historical Tracking**: Maintain performance baseline over time
4. **Cross-Platform Testing**: Ensure consistent performance across OS platforms

### Performance Gates

```yaml
# Example performance gate configuration
performance_gates:
  startup_time: 100ms
  memory_usage: 50MB
  action_latency_p95: 50ms
  regression_threshold: 20%
```

## Troubleshooting Performance Issues

### Common Performance Problems

#### 1. Slow Startup
- **Symptoms**: Application takes >200ms to start
- **Causes**: Heavy initialization, configuration loading
- **Solutions**: Lazy loading, configuration caching

#### 2. High Memory Usage
- **Symptoms**: Memory usage >100MB during normal operation
- **Causes**: Memory leaks, excessive object allocation
- **Solutions**: Object pooling, garbage collection tuning

#### 3. Sluggish Hotkey Response
- **Symptoms**: >100ms delay between keypress and action
- **Causes**: Complex key parsing, inefficient event handling
- **Solutions**: Key sequence pre-compilation, event filtering

#### 4. Configuration Reload Delays
- **Symptoms**: >500ms for configuration reload
- **Causes**: Full re-parsing, validation overhead
- **Solutions**: Incremental updates, validation caching

### Diagnostic Tools

```bash
# Memory leak detection
go test -run=TestMemoryLeaks -memprofile=leak.prof
go tool pprof -alloc_space leak.prof

# CPU bottleneck analysis
go test -bench=. -cpuprofile=cpu.prof
go tool pprof -top cpu.prof

# Goroutine leak detection
go test -run=TestGoroutineLeaks
```

### Performance Debugging

```go
// Add performance logging
func logPerformance(operation string, fn func()) {
    start := time.Now()
    defer func() {
        duration := time.Since(start)
        if duration > threshold {
            log.Printf("SLOW: %s took %v", operation, duration)
        }
    }()
    fn()
}
```

## Best Practices

### Development Guidelines

1. **Profile Before Optimizing**: Always measure before making performance changes
2. **Optimize Hot Paths**: Focus on frequently executed code paths
3. **Minimize Allocations**: Reduce garbage collection pressure
4. **Use Appropriate Data Structures**: Choose efficient algorithms and data structures
5. **Avoid Premature Optimization**: Don't optimize until you have evidence of problems

### Code Review Checklist

- [ ] Are there any obvious performance bottlenecks?
- [ ] Is memory allocation minimized in hot paths?
- [ ] Are expensive operations cached when possible?
- [ ] Is error handling efficient and non-blocking?
- [ ] Are goroutines cleaned up properly?

### Testing Requirements

- [ ] Performance tests added for new features
- [ ] Benchmark comparison against previous version
- [ ] Memory leak tests for long-running operations
- [ ] Stress tests for resource-intensive features

## Advanced Topics

### Garbage Collection Tuning

```bash
# Tune GC for lower latency
export GOGC=50  # More frequent, shorter GC cycles

# Tune GC for throughput
export GOGC=200 # Less frequent, longer GC cycles
```

### Platform-Specific Optimizations

#### Windows
- Use Windows-specific APIs for better integration
- Optimize for Windows task scheduler behavior

#### macOS
- Leverage Cocoa APIs for system integration
- Optimize for macOS security model

#### Linux
- Use epoll for efficient event handling
- Optimize for various desktop environments

### Memory Pool Strategies

```go
// Different pool strategies for different use cases
type PoolStrategy interface {
    Get() interface{}
    Put(interface{})
}

// Fixed-size pool for predictable workloads
type FixedPool struct {
    pool chan interface{}
}

// Dynamic pool for variable workloads
type DynamicPool struct {
    sync.Pool
}
```

## See Also

- [Performance Documentation](../performance/README.md)
- [Development Setup](../development/setup.md)
- [Troubleshooting Performance](../troubleshooting/performance.md)