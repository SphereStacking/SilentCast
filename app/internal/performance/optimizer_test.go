package performance

import (
	"context"
	"runtime"
	"testing"
	"time"
)

// TestResourceManager tests the resource manager functionality
func TestResourceManager(t *testing.T) {
	config := Config{
		BufferSize:  1024,
		MaxIdleTime: time.Minute,
		GCPercent:   100,
	}
	
	rm := NewResourceManager(config)
	defer rm.Cleanup()
	
	// Test string pool
	t.Run("string pool", func(t *testing.T) {
		slice1 := rm.GetStringSlice()
		if len(slice1) != 0 {
			t.Errorf("expected empty slice, got length %d", len(slice1))
		}
		
		slice1 = append(slice1, "test")
		rm.PutStringSlice(slice1)
		
		slice2 := rm.GetStringSlice()
		if cap(slice2) == 0 {
			t.Error("expected reused slice with capacity")
		}
		rm.PutStringSlice(slice2)
	})
	
	// Test buffer pool
	t.Run("buffer pool", func(t *testing.T) {
		buffer1 := rm.GetBuffer()
		if len(buffer1) != 0 {
			t.Errorf("expected empty buffer, got length %d", len(buffer1))
		}
		
		buffer1 = append(buffer1, []byte("test data")...)
		rm.PutBuffer(buffer1)
		
		buffer2 := rm.GetBuffer()
		if cap(buffer2) == 0 {
			t.Error("expected reused buffer with capacity")
		}
		rm.PutBuffer(buffer2)
	})
	
	// Test metrics
	t.Run("metrics", func(t *testing.T) {
		metrics := rm.GetMetrics()
		if metrics.StringPoolHits == 0 {
			t.Error("expected string pool hits to be recorded")
		}
		if metrics.BufferPoolHits == 0 {
			t.Error("expected buffer pool hits to be recorded")
		}
	})
}

// TestMemoryOptimization tests memory optimization features
func TestMemoryOptimization(t *testing.T) {
	config := Config{
		BufferSize: 512,
		GCPercent:  50, // More aggressive GC for testing
	}
	
	rm := NewResourceManager(config)
	defer rm.Cleanup()
	
	// Get initial memory stats
	var initialStats runtime.MemStats
	runtime.ReadMemStats(&initialStats)
	
	// Allocate some memory
	buffers := make([][]byte, 100)
	for i := range buffers {
		buffers[i] = rm.GetBuffer()
		for j := 0; j < 100; j++ {
			buffers[i] = append(buffers[i], byte(j))
		}
	}
	
	// Return buffers to pool
	for _, buffer := range buffers {
		rm.PutBuffer(buffer)
	}
	
	// Optimize memory
	rm.OptimizeMemory()
	
	// Check that metrics were updated
	metrics := rm.GetMetrics()
	if metrics.AllocatedBytes == 0 {
		t.Error("expected allocated bytes to be tracked")
	}
	if metrics.GCRuns == 0 {
		t.Error("expected GC runs to be tracked")
	}
}

// TestGoroutineMonitoring tests goroutine leak detection
func TestGoroutineMonitoring(t *testing.T) {
	rm := NewResourceManager(Config{})
	defer rm.Cleanup()
	
	initialCount := rm.MonitorGoroutines()
	
	// Start some goroutines
	done := make(chan bool, 3)
	for i := 0; i < 3; i++ {
		go func() {
			time.Sleep(10 * time.Millisecond)
			done <- true
		}()
	}
	
	// Check increased count
	midCount := rm.MonitorGoroutines()
	if midCount <= initialCount {
		t.Errorf("expected goroutine count to increase, got %d -> %d", initialCount, midCount)
	}
	
	// Wait for goroutines to finish
	for i := 0; i < 3; i++ {
		<-done
	}
	
	// Give runtime time to clean up
	time.Sleep(50 * time.Millisecond)
	finalCount := rm.MonitorGoroutines()
	
	// Count should return to approximately initial level
	if finalCount > midCount {
		t.Errorf("possible goroutine leak: %d -> %d -> %d", initialCount, midCount, finalCount)
	}
}

// TestResponseTimeTracking tests response time monitoring
func TestResponseTimeTracking(t *testing.T) {
	rm := NewResourceManager(Config{})
	defer rm.Cleanup()
	
	// Track some response times
	rm.TrackResponseTime(10 * time.Millisecond)
	rm.TrackResponseTime(20 * time.Millisecond)
	rm.TrackResponseTime(30 * time.Millisecond)
	
	metrics := rm.GetMetrics()
	if metrics.MaxResponseTime != 30*time.Millisecond {
		t.Errorf("expected max response time 30ms, got %v", metrics.MaxResponseTime)
	}
	if metrics.AverageResponseTime == 0 {
		t.Error("expected average response time to be calculated")
	}
}

// BenchmarkStringPool benchmarks string pool performance
func BenchmarkStringPool(b *testing.B) {
	rm := NewResourceManager(Config{})
	defer rm.Cleanup()
	
	b.ResetTimer()
	
	b.Run("with pool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			slice := rm.GetStringSlice()
			slice = append(slice, "test", "data", "benchmark")
			rm.PutStringSlice(slice)
		}
	})
	
	b.Run("without pool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			slice := make([]string, 0, 16)
			_ = append(slice, "test", "data", "benchmark")
			// No reuse
		}
	})
}

// BenchmarkBufferPool benchmarks buffer pool performance
func BenchmarkBufferPool(b *testing.B) {
	rm := NewResourceManager(Config{BufferSize: 1024})
	defer rm.Cleanup()
	
	data := []byte("test data for benchmark")
	
	b.ResetTimer()
	
	b.Run("with pool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buffer := rm.GetBuffer()
			buffer = append(buffer, data...)
			rm.PutBuffer(buffer)
		}
	})
	
	b.Run("without pool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buffer := make([]byte, 0, 1024)
			_ = append(buffer, data...)
			// No reuse
		}
	})
}

// BenchmarkContextCreation benchmarks optimized context creation
func BenchmarkContextCreation(b *testing.B) {
	rm := NewResourceManager(Config{})
	defer rm.Cleanup()
	
	b.ResetTimer()
	
	b.Run("optimized context", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ctx, cancel := rm.CreateOptimizedContext(context.Background(), time.Second)
			cancel()
			_ = ctx
		}
	})
	
	b.Run("standard context", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			cancel()
			_ = ctx
		}
	})
}