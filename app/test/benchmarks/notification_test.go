package benchmarks

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/SphereStacking/silentcast/internal/notify"
)

// BenchmarkNotificationCreation measures notification object creation
func BenchmarkNotificationCreation(b *testing.B) {
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		notification := notify.Notification{
			Title:   "Benchmark Test",
			Message: "This is a benchmark test notification",
			Level:   notify.LevelInfo,
		}
		_ = notification
	}
}

// BenchmarkNotificationSending measures notification sending performance
func BenchmarkNotificationSending(b *testing.B) {
	manager := notify.NewManager()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	notification := notify.Notification{
		Title:   "Benchmark Test",
		Message: "Benchmark notification",
		Level:   notify.LevelInfo,
	}
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		err := manager.Notify(ctx, notification)
		if err != nil {
			b.Fatalf("Notification failed: %v", err)
		}
	}
}

// BenchmarkNotificationLevels measures different notification levels
func BenchmarkNotificationLevels(b *testing.B) {
	manager := notify.NewManager()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	levels := []notify.Level{
		notify.LevelInfo,
		notify.LevelWarning,
		notify.LevelError,
		notify.LevelSuccess,
	}
	
	for _, level := range levels {
		b.Run(level.String(), func(b *testing.B) {
			notification := notify.Notification{
				Title:   "Level Test",
				Message: "Testing " + level.String() + " level",
				Level:   level,
			}
			
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				err := manager.Notify(ctx, notification)
				if err != nil {
					b.Fatalf("Notification failed: %v", err)
				}
			}
		})
	}
}

// BenchmarkNotificationWithLongContent measures performance with large content
func BenchmarkNotificationWithLongContent(b *testing.B) {
	manager := notify.NewManager()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// Create progressively larger content
	sizes := []int{100, 500, 1000, 5000, 10000}
	
	for _, size := range sizes {
		b.Run("Size"+string(rune(size/1000+48))+"K", func(b *testing.B) {
			content := strings.Repeat("This is a test notification with long content. ", size/45)
			
			notification := notify.Notification{
				Title:   "Long Content Test",
				Message: content,
				Level:   notify.LevelInfo,
			}
			
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				err := manager.Notify(ctx, notification)
				// Large notifications may fail, that's okay for benchmarking
				_ = err
			}
		})
	}
}

// BenchmarkConcurrentNotifications measures concurrent notification performance
func BenchmarkConcurrentNotifications(b *testing.B) {
	manager := notify.NewManager()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	b.ReportAllocs()
	
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			notification := notify.Notification{
				Title:   "Concurrent Test",
				Message: "Concurrent notification " + string(rune(counter%10+48)),
				Level:   notify.Level(counter % 4),
			}
			
			err := manager.Notify(ctx, notification)
			if err != nil {
				b.Fatalf("Concurrent notification failed: %v", err)
			}
			counter++
		}
	})
}

// BenchmarkNotificationFormatting measures message formatting performance
func BenchmarkNotificationFormatting(b *testing.B) {
	testMessages := []string{
		"Simple message",
		"Message with\nmultiple\nlines",
		"Message with special characters: !@#$%^&*()",
		"Message with unicode: 🎯⚡🌐🎨🛠️",
		"Very long message " + strings.Repeat("that keeps going ", 50),
		"Message with\ttabs\tand    spaces",
		"Message with ANSI codes: \033[31mRed\033[0m \033[32mGreen\033[0m",
	}
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		for j, msg := range testMessages {
			// Simulate notification formatting
			formatted := strings.ReplaceAll(msg, "\n", " ")
			formatted = strings.ReplaceAll(formatted, "\t", " ")
			
			if len(formatted) > 200 {
				formatted = formatted[:200] + "..."
			}
			
			notification := notify.Notification{
				Title:   "Format Test " + string(rune(j+48)),
				Message: formatted,
				Level:   notify.LevelInfo,
			}
			_ = notification
		}
	}
}

// BenchmarkNotificationTimeout measures timeout handling
func BenchmarkNotificationTimeout(b *testing.B) {
	manager := notify.NewManager()
	
	notification := notify.Notification{
		Title:   "Timeout Test",
		Message: "Testing timeout behavior",
		Level:   notify.LevelInfo,
	}
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		// Very short timeout to test timeout handling
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		_ = manager.Notify(ctx, notification)
		cancel()
	}
}

// BenchmarkNotificationMemoryUsage measures memory usage of notifications
func BenchmarkNotificationMemoryUsage(b *testing.B) {
	RunMemoryBenchmark(b, func() {
		manager := notify.NewManager()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		// Send multiple notifications
		for i := 0; i < 10; i++ {
			notification := notify.Notification{
				Title:   "Memory Test " + string(rune(i+48)),
				Message: "Testing memory usage of notifications",
				Level:   notify.Level(i % 4),
			}
			_ = manager.Notify(ctx, notification)
		}
	})
}

// BenchmarkNotificationManagerCreation measures manager creation overhead
func BenchmarkNotificationManagerCreation(b *testing.B) {
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		manager := notify.NewManager()
		_ = manager
	}
}

// BenchmarkConsoleNotifier measures console notifier performance
func BenchmarkConsoleNotifier(b *testing.B) {
	notifier := notify.NewConsoleNotifier()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	notification := notify.Notification{
		Title:   "Console Test",
		Message: "Testing console notifier performance",
		Level:   notify.LevelInfo,
	}
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		err := notifier.Notify(ctx, notification)
		if err != nil {
			b.Fatalf("Console notification failed: %v", err)
		}
	}
}

// BenchmarkSystemNotifier measures system notifier performance
func BenchmarkSystemNotifier(b *testing.B) {
	notifier := notify.GetSystemNotifier()
	if notifier == nil {
		b.Skip("System notifier not available")
		return
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	notification := notify.Notification{
		Title:   "System Test",
		Message: "Testing system notifier performance",
		Level:   notify.LevelInfo,
	}
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		err := notifier.Notify(ctx, notification)
		// System notifications may fail in test environment
		_ = err
	}
}

// BenchmarkNotificationQueue measures notification queuing performance
func BenchmarkNotificationQueue(b *testing.B) {
	manager := notify.NewManager()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// Create many notifications to queue
	notifications := make([]notify.Notification, 100)
	for i := range notifications {
		notifications[i] = notify.Notification{
			Title:   "Queue Test " + string(rune(i%10+48)),
			Message: "Testing notification queuing",
			Level:   notify.Level(i % 4),
		}
	}
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		for _, notification := range notifications {
			err := manager.Notify(ctx, notification)
			if err != nil {
				b.Fatalf("Queued notification failed: %v", err)
			}
		}
	}
}

// BenchmarkNotificationStringOperations measures string operations in notifications
func BenchmarkNotificationStringOperations(b *testing.B) {
	testStrings := []string{
		"Simple string",
		"String with\nmultiple\nlines\nand\nmore\ncontent",
		"String with special characters !@#$%^&*()",
		strings.Repeat("Long string content ", 100),
		"String with\ttabs\tand    multiple    spaces",
	}
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		for _, str := range testStrings {
			// Common string operations in notification processing
			cleaned := strings.ReplaceAll(str, "\n", " ")
			cleaned = strings.ReplaceAll(cleaned, "\t", " ")
			cleaned = strings.TrimSpace(cleaned)
			
			if len(cleaned) > 100 {
				cleaned = cleaned[:100] + "..."
			}
			
			_ = cleaned
		}
	}
}

// BenchmarkNotificationContextOperations measures context operations
func BenchmarkNotificationContextOperations(b *testing.B) {
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		// Test context creation and cancellation patterns
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		
		// Simulate context usage
		select {
		case <-ctx.Done():
			// Context expired
		default:
			// Context still valid
		}
		
		cancel()
	}
}