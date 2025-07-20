package output

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestNewTeeManager(t *testing.T) {
	tests := []struct {
		name    string
		options Options
		want    Options
	}{
		{
			name:    "empty options should use defaults",
			options: Options{},
			want: Options{
				MaxSize:         1024 * 1024,
				TruncateMessage: "\n... (output truncated)",
				Destinations:    []Destination{DestConsole},
			},
		},
		{
			name: "custom options should be preserved",
			options: Options{
				Type:            TypeTee,
				MaxSize:         512,
				TruncateMessage: "truncated",
				Destinations:    []Destination{DestNotification},
			},
			want: Options{
				Type:            TypeTee,
				MaxSize:         512,
				TruncateMessage: "truncated",
				Destinations:    []Destination{DestNotification},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewTeeManager(tt.options)

			if manager == nil {
				t.Fatal("NewTeeManager returned nil")
			}

			if manager.options.MaxSize != tt.want.MaxSize {
				t.Errorf("MaxSize = %d, want %d", manager.options.MaxSize, tt.want.MaxSize)
			}

			if manager.options.TruncateMessage != tt.want.TruncateMessage {
				t.Errorf("TruncateMessage = %q, want %q", manager.options.TruncateMessage, tt.want.TruncateMessage)
			}
		})
	}
}

func TestTeeManager_StartCapture(t *testing.T) {
	manager := NewTeeManager(Options{})

	writer := manager.StartCapture()
	if writer == nil {
		t.Fatal("StartCapture returned nil writer")
	}

	// Test that we can write to the writer
	data := []byte("test data")
	n, err := writer.Write(data)
	if err != nil {
		t.Errorf("Write failed: %v", err)
	}
	if n != len(data) {
		t.Errorf("Write returned %d bytes, want %d", n, len(data))
	}

	// Stop the manager and verify StartCapture returns nil
	manager.Stop()
	writer2 := manager.StartCapture()
	if writer2 != nil {
		t.Error("StartCapture should return nil after Stop")
	}
}

func TestTeeManager_BufferAndStream(t *testing.T) {
	manager := NewTeeManager(Options{})

	// Add destinations
	var buf1, buf2 bytes.Buffer
	manager.AddDestination(&buf1)
	manager.AddDestination(&buf2)

	// Start capturing and write data
	writer := manager.StartCapture()
	testData := "Hello, Tee!"
	writer.Write([]byte(testData))

	// Verify data was buffered
	buffered := manager.GetOutput()
	if buffered != testData {
		t.Errorf("Buffered data = %q, want %q", buffered, testData)
	}

	// Verify data was streamed to destinations
	if buf1.String() != testData {
		t.Errorf("Destination 1 = %q, want %q", buf1.String(), testData)
	}
	if buf2.String() != testData {
		t.Errorf("Destination 2 = %q, want %q", buf2.String(), testData)
	}
}

func TestTeeManager_Stream(t *testing.T) {
	manager := NewTeeManager(Options{})

	// Write some initial data
	writer := manager.StartCapture()
	initialData := "Initial data"
	writer.Write([]byte(initialData))

	// Add a destination via Stream method
	var buf bytes.Buffer
	err := manager.Stream(&buf)
	if err != nil {
		t.Errorf("Stream failed: %v", err)
	}

	// Verify existing data was streamed
	if buf.String() != initialData {
		t.Errorf("Existing data not streamed, got %q, want %q", buf.String(), initialData)
	}

	// Write more data
	additionalData := " and more"
	writer.Write([]byte(additionalData))

	// Verify complete data in destination
	expected := initialData + additionalData
	if buf.String() != expected {
		t.Errorf("Complete data = %q, want %q", buf.String(), expected)
	}

	// Test Stream after Stop
	manager.Stop()
	var buf2 bytes.Buffer
	err = manager.Stream(&buf2)
	if err == nil {
		t.Error("Stream should fail after Stop")
	}
}

func TestTeeManager_AddDestination(t *testing.T) {
	manager := NewTeeManager(Options{})

	// Write initial data
	writer := manager.StartCapture()
	initialData := "Before add"
	writer.Write([]byte(initialData))

	// Add destination
	var buf bytes.Buffer
	err := manager.AddDestination(&buf)
	if err != nil {
		t.Errorf("AddDestination failed: %v", err)
	}

	// Verify existing data was synced
	if buf.String() != initialData {
		t.Errorf("Existing data not synced to new destination, got %q", buf.String())
	}

	// Write more data
	newData := " | After add"
	writer.Write([]byte(newData))

	// Verify complete data
	expected := initialData + newData
	if buf.String() != expected {
		t.Errorf("Buffer = %q, want %q", buf.String(), expected)
	}

	// Test AddDestination after Stop
	manager.Stop()
	var buf2 bytes.Buffer
	err = manager.AddDestination(&buf2)
	if err == nil {
		t.Error("AddDestination should fail after Stop")
	}
}

func TestTeeManager_RemoveDestination(t *testing.T) {
	manager := NewTeeManager(Options{})

	var buf1, buf2, buf3 bytes.Buffer
	manager.AddDestination(&buf1)
	manager.AddDestination(&buf2)
	manager.AddDestination(&buf3)

	// Initial count
	if count := manager.GetDestinationCount(); count != 3 {
		t.Errorf("Initial destination count = %d, want 3", count)
	}

	// Write data
	writer := manager.StartCapture()
	writer.Write([]byte("test1"))

	// All should have data
	for i, buf := range []*bytes.Buffer{&buf1, &buf2, &buf3} {
		if buf.String() != "test1" {
			t.Errorf("Buffer %d = %q, want %q", i+1, buf.String(), "test1")
		}
	}

	// Remove middle destination
	manager.RemoveDestination(&buf2)

	if count := manager.GetDestinationCount(); count != 2 {
		t.Errorf("After removal, destination count = %d, want 2", count)
	}

	// Write more data
	writer.Write([]byte("test2"))

	// buf1 and buf3 should have new data, buf2 should not
	if buf1.String() != "test1test2" {
		t.Errorf("buf1 = %q, want %q", buf1.String(), "test1test2")
	}
	if buf2.String() != "test1" {
		t.Errorf("buf2 should not receive new data, got %q", buf2.String())
	}
	if buf3.String() != "test1test2" {
		t.Errorf("buf3 = %q, want %q", buf3.String(), "test1test2")
	}
}

func TestTeeManager_MaxSize(t *testing.T) {
	manager := NewTeeManager(Options{
		MaxSize: 10, // Very small for testing
	})

	var buf bytes.Buffer
	manager.AddDestination(&buf)

	writer := manager.StartCapture()

	// Write data that exceeds max size
	longData := "This is a very long string that exceeds max size"
	writer.Write([]byte(longData))

	// Buffer should be truncated
	buffered := manager.GetOutput()
	if len(buffered) > 10 {
		t.Errorf("Buffer size = %d, want <= 10", len(buffered))
	}

	// But destination should have full data
	if buf.String() != longData {
		t.Errorf("Destination should have full data, got %q", buf.String())
	}
}

func TestTeeManager_Clear(t *testing.T) {
	manager := NewTeeManager(Options{})

	var buf bytes.Buffer
	manager.AddDestination(&buf)

	writer := manager.StartCapture()
	writer.Write([]byte("test data"))

	// Verify data exists
	if manager.GetOutput() == "" {
		t.Error("Buffer should have data before clear")
	}
	if manager.GetDestinationCount() == 0 {
		t.Error("Should have destinations before clear")
	}

	// Clear
	manager.Clear()

	// Verify cleared
	if manager.GetOutput() != "" {
		t.Error("Buffer should be empty after clear")
	}
	if manager.GetDestinationCount() != 0 {
		t.Error("Destinations should be cleared")
	}

	// Writing should not affect cleared destination
	writer.Write([]byte("after clear"))
	if buf.String() != "test data" {
		t.Errorf("Cleared destination should not receive new data, got %q", buf.String())
	}
}

func TestTeeManager_ConcurrentAccess(t *testing.T) {
	manager := NewTeeManager(Options{})

	var buffers [5]bytes.Buffer
	for i := range buffers {
		manager.AddDestination(&buffers[i])
	}

	writer := manager.StartCapture()

	// Concurrent writes
	var wg sync.WaitGroup
	iterations := 100

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				data := fmt.Sprintf("W%d-%d|", id, j)
				writer.Write([]byte(data))
				time.Sleep(time.Microsecond)
			}
		}(i)
	}

	// Concurrent destination management
	wg.Add(1)
	go func() {
		defer wg.Done()
		var tempBuf bytes.Buffer
		for i := 0; i < 10; i++ {
			manager.AddDestination(&tempBuf)
			time.Sleep(time.Millisecond)
			manager.RemoveDestination(&tempBuf)
		}
	}()

	wg.Wait()

	// Verify all original destinations have data
	for i, buf := range buffers {
		if buf.Len() == 0 {
			t.Errorf("Buffer %d has no data", i)
		}
	}

	// Verify buffer has data
	output := manager.GetOutput()
	if output == "" {
		t.Error("Manager buffer has no data")
	}

	// Verify expected patterns in output
	for i := 0; i < 3; i++ {
		pattern := fmt.Sprintf("W%d-", i)
		if !strings.Contains(output, pattern) {
			t.Errorf("Output missing data from writer %d", i)
		}
	}
}

func TestTeeManager_DestinationError(t *testing.T) {
	manager := NewTeeManager(Options{})

	// Add a working destination
	var goodBuf bytes.Buffer
	manager.AddDestination(&goodBuf)

	// Add a failing destination
	failWriter := &teeFailingWriter{shouldFail: true}
	manager.AddDestination(failWriter)

	// Write should succeed even with one failing destination
	writer := manager.StartCapture()
	testData := "test data"
	n, err := writer.Write([]byte(testData))

	if err != nil {
		t.Errorf("Write should succeed with partial destination failure, got error: %v", err)
	}
	if n != len(testData) {
		t.Errorf("Write returned %d, want %d", n, len(testData))
	}

	// Good destination should have data
	if goodBuf.String() != testData {
		t.Errorf("Good destination = %q, want %q", goodBuf.String(), testData)
	}

	// Buffer should have data
	if manager.GetOutput() != testData {
		t.Errorf("Buffer = %q, want %q", manager.GetOutput(), testData)
	}
}

func TestTeeManager_AllDestinationsFail(t *testing.T) {
	manager := NewTeeManager(Options{})

	// Add only failing destinations
	for i := 0; i < 3; i++ {
		manager.AddDestination(&teeFailingWriter{shouldFail: true})
	}

	// Write should still succeed (buffer still works)
	writer := manager.StartCapture()
	testData := "test data"
	n, err := writer.Write([]byte(testData))

	// Write should report success to avoid breaking the writer
	if err != nil {
		t.Errorf("Write should succeed even if all destinations fail, got: %v", err)
	}
	if n != len(testData) {
		t.Errorf("Write returned %d, want %d", n, len(testData))
	}

	// Buffer should still have data
	if manager.GetOutput() != testData {
		t.Errorf("Buffer should still work, got %q, want %q", manager.GetOutput(), testData)
	}
}

func TestTeeManager_Stop(t *testing.T) {
	manager := NewTeeManager(Options{})

	var buf bytes.Buffer
	manager.AddDestination(&buf)

	writer := manager.StartCapture()
	writer.Write([]byte("before stop"))

	// Stop the manager
	err := manager.Stop()
	if err != nil {
		t.Errorf("Stop failed: %v", err)
	}

	// Writing after stop should fail
	_, err = writer.Write([]byte("after stop"))
	if err == nil {
		t.Error("Write should fail after Stop")
	}

	// GetOutput should still work after stop
	output := manager.GetOutput()
	if output != "before stop" {
		t.Errorf("GetOutput after stop = %q, want %q", output, "before stop")
	}
}

func TestTeeManager_AddDestinationSync(t *testing.T) {
	manager := NewTeeManager(Options{})

	// Write substantial initial data
	writer := manager.StartCapture()
	initialData := strings.Repeat("Initial data line\n", 100)
	writer.Write([]byte(initialData))

	// Add destination with failing initial sync
	failOnce := &teeFailingWriter{failCount: 1} // Fail only first write
	err := manager.AddDestination(failOnce)
	if err == nil {
		t.Error("AddDestination should fail if initial sync fails")
	}

	// Destination should not be added
	if manager.GetDestinationCount() != 0 {
		t.Error("Failed destination should not be added")
	}
}

// teeFailingWriter is a test writer that can be configured to fail
type teeFailingWriter struct {
	shouldFail bool
	failCount  int
	writeCount int
}

func (w *teeFailingWriter) Write(p []byte) (n int, err error) {
	w.writeCount++
	if w.shouldFail || (w.failCount > 0 && w.writeCount <= w.failCount) {
		return 0, fmt.Errorf("intentional write failure")
	}
	return len(p), nil
}
