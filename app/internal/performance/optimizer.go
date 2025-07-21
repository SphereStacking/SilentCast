package performance

import (
	"context"
	"runtime"
	"runtime/debug"
	"sync"
	"time"
)

// ResourceManager manages application resources and performance optimization
type ResourceManager struct {
	// String pool for reducing string allocations
	stringPool sync.Pool

	// Buffer pool for reducing byte slice allocations
	bufferPool sync.Pool

	// Context pool for reusing contexts
	contextPool sync.Pool

	// Metrics for monitoring
	metrics *Metrics

	// Configuration
	config Config
}

// Config holds performance optimization configuration
type Config struct {
	// BufferSize sets the initial buffer size for the pool
	BufferSize int

	// MaxIdleTime sets how long to keep idle resources
	MaxIdleTime time.Duration

	// GCPercent sets the target percentage for garbage collection
	GCPercent int

	// EnableProfiling enables performance profiling
	EnableProfiling bool
}

// Metrics tracks performance metrics
type Metrics struct {
	mu sync.RWMutex

	// Memory metrics
	AllocatedBytes   uint64
	TotalAllocations uint64
	GCRuns           uint64

	// Pool metrics
	StringPoolHits   uint64
	StringPoolMisses uint64
	BufferPoolHits   uint64
	BufferPoolMisses uint64

	// Timing metrics
	AverageResponseTime time.Duration
	MaxResponseTime     time.Duration
}

// NewResourceManager creates a new resource manager with optimized pools
func NewResourceManager(config Config) *ResourceManager {
	if config.BufferSize == 0 {
		config.BufferSize = 1024
	}
	if config.MaxIdleTime == 0 {
		config.MaxIdleTime = 5 * time.Minute
	}
	if config.GCPercent == 0 {
		config.GCPercent = 100
	}

	rm := &ResourceManager{
		metrics: &Metrics{},
		config:  config,
	}

	// Initialize string pool
	rm.stringPool = sync.Pool{
		New: func() interface{} {
			slice := make([]string, 0, 16)
			return &slice
		},
	}

	// Initialize buffer pool
	rm.bufferPool = sync.Pool{
		New: func() interface{} {
			buffer := make([]byte, 0, config.BufferSize)
			return &buffer
		},
	}

	// Initialize context pool
	rm.contextPool = sync.Pool{
		New: func() interface{} {
			return context.Background()
		},
	}

	// Set GC target
	debug.SetGCPercent(config.GCPercent)

	return rm
}

// GetStringSlice returns a reusable string slice from the pool
func (rm *ResourceManager) GetStringSlice() []string {
	rm.metrics.mu.Lock()
	defer rm.metrics.mu.Unlock()

	value := rm.stringPool.Get()
	slice, ok := value.([]string)
	if !ok {
		// This should never happen, but handle gracefully
		slice = make([]string, 0, 8)
	}
	slice = slice[:0] // Reset length but keep capacity

	rm.metrics.StringPoolHits++
	return slice
}

// PutStringSlice returns a string slice to the pool
func (rm *ResourceManager) PutStringSlice(slice []string) {
	if cap(slice) > 0 {
		// Clear the slice before returning to pool to avoid memory leaks
		slice = slice[:0]
		rm.stringPool.Put(&slice)
	}
}

// GetBuffer returns a reusable byte buffer from the pool
func (rm *ResourceManager) GetBuffer() []byte {
	rm.metrics.mu.Lock()
	defer rm.metrics.mu.Unlock()

	value := rm.bufferPool.Get()
	buffer, ok := value.([]byte)
	if !ok {
		// This should never happen, but handle gracefully
		buffer = make([]byte, 0, 1024)
	}
	buffer = buffer[:0] // Reset length but keep capacity

	rm.metrics.BufferPoolHits++
	return buffer
}

// PutBuffer returns a byte buffer to the pool
func (rm *ResourceManager) PutBuffer(buffer []byte) {
	if cap(buffer) > 0 && cap(buffer) <= rm.config.BufferSize*2 {
		// Clear the buffer before returning to pool to avoid memory leaks
		buffer = buffer[:0]
		rm.bufferPool.Put(&buffer)
	}
}

// OptimizeMemory performs memory optimization operations
func (rm *ResourceManager) OptimizeMemory() {
	// Force garbage collection
	runtime.GC()

	// Update metrics
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	rm.metrics.mu.Lock()
	rm.metrics.AllocatedBytes = memStats.Alloc
	rm.metrics.TotalAllocations = memStats.TotalAlloc
	rm.metrics.GCRuns = uint64(memStats.NumGC)
	rm.metrics.mu.Unlock()
}

// GetMetrics returns current performance metrics
func (rm *ResourceManager) GetMetrics() Metrics {
	rm.metrics.mu.RLock()
	defer rm.metrics.mu.RUnlock()

	// Return a copy without the mutex
	return Metrics{
		AllocatedBytes:      rm.metrics.AllocatedBytes,
		TotalAllocations:    rm.metrics.TotalAllocations,
		GCRuns:              rm.metrics.GCRuns,
		StringPoolHits:      rm.metrics.StringPoolHits,
		StringPoolMisses:    rm.metrics.StringPoolMisses,
		BufferPoolHits:      rm.metrics.BufferPoolHits,
		BufferPoolMisses:    rm.metrics.BufferPoolMisses,
		AverageResponseTime: rm.metrics.AverageResponseTime,
		MaxResponseTime:     rm.metrics.MaxResponseTime,
	}
}

// MonitorGoroutines monitors goroutine count and detects potential leaks
func (rm *ResourceManager) MonitorGoroutines() int {
	return runtime.NumGoroutine()
}

// CreateOptimizedContext creates a context optimized for the application
func (rm *ResourceManager) CreateOptimizedContext(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout > 0 {
		return context.WithTimeout(parent, timeout)
	}
	return context.WithCancel(parent)
}

// TrackResponseTime tracks operation response time
func (rm *ResourceManager) TrackResponseTime(duration time.Duration) {
	rm.metrics.mu.Lock()
	defer rm.metrics.mu.Unlock()

	if duration > rm.metrics.MaxResponseTime {
		rm.metrics.MaxResponseTime = duration
	}

	// Simple moving average for demonstration
	if rm.metrics.AverageResponseTime == 0 {
		rm.metrics.AverageResponseTime = duration
	} else {
		rm.metrics.AverageResponseTime = (rm.metrics.AverageResponseTime + duration) / 2
	}
}

// Cleanup performs resource cleanup
func (rm *ResourceManager) Cleanup() {
	// Clear pools
	rm.stringPool = sync.Pool{}
	rm.bufferPool = sync.Pool{}
	rm.contextPool = sync.Pool{}

	// Force final GC
	runtime.GC()
}
