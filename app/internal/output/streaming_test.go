package output

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestNewStreamingManager(t *testing.T) {
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
				Type:            TypeStreaming,
				MaxSize:         512,
				TruncateMessage: "truncated",
				Destinations:    []Destination{DestNotification},
			},
			want: Options{
				Type:            TypeStreaming,
				MaxSize:         512,
				TruncateMessage: "truncated",
				Destinations:    []Destination{DestNotification},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewStreamingManager(tt.options)

			if manager == nil {
				t.Fatal("NewStreamingManager returned nil")
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

func TestStreamingManager_StartCapture(t *testing.T) {
	manager := NewStreamingManager(Options{})

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
	if stopErr := manager.Stop(); stopErr != nil {
		t.Fatalf("Stop failed: %v", stopErr)
	}
	writer2 := manager.StartCapture()
	if writer2 != nil {
		t.Error("StartCapture should return nil after Stop")
	}
}

func TestStreamingManager_GetOutput(t *testing.T) {
	manager := NewStreamingManager(Options{})

	// GetOutput should always return empty string for streaming manager
	output := manager.GetOutput()
	if output != "" {
		t.Errorf("GetOutput = %q, want empty string", output)
	}

	// Even after writing data, GetOutput should return empty
	writer := manager.StartCapture()
	_, err := writer.Write([]byte("test data"))
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	output = manager.GetOutput()
	if output != "" {
		t.Errorf("GetOutput = %q, want empty string after writing", output)
	}
}

func TestStreamingManager_Stream(t *testing.T) {
	manager := NewStreamingManager(Options{})

	// Add a destination
	var buf bytes.Buffer
	err := manager.Stream(&buf)
	if err != nil {
		t.Errorf("Stream failed: %v", err)
	}

	// Start capturing and write data
	writer := manager.StartCapture()
	testData := "Hello, streaming!"
	_, err = writer.Write([]byte(testData))
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Verify data was streamed to destination
	if buf.String() != testData {
		t.Errorf("Streamed data = %q, want %q", buf.String(), testData)
	}

	// Test that Stream fails after Stop
	if stopErr := manager.Stop(); stopErr != nil {
		t.Fatalf("Stop failed: %v", stopErr)
	}
	var buf2 bytes.Buffer
	err = manager.Stream(&buf2)
	if err == nil {
		t.Error("Stream should fail after Stop")
	}
}

func TestStreamingManager_MultipleDestinations(t *testing.T) {
	manager := NewStreamingManager(Options{})

	// Add multiple destinations
	var buf1, buf2, buf3 bytes.Buffer
	err := manager.Stream(&buf1)
	if err != nil {
		t.Fatalf("Stream failed: %v", err)
	}
	manager.Stream(&buf2)
	manager.Stream(&buf3)

	// Write data
	writer := manager.StartCapture()
	testData := "Multi-destination test"
	_, err = writer.Write([]byte(testData))
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Verify all destinations received the data
	destinations := []*bytes.Buffer{&buf1, &buf2, &buf3}
	for i, buf := range destinations {
		if buf.String() != testData {
			t.Errorf("Destination %d got %q, want %q", i, buf.String(), testData)
		}
	}
}

func TestStreamingManager_AddRemoveDestination(t *testing.T) {
	manager := NewStreamingManager(Options{})

	var buf1, buf2 bytes.Buffer

	// Add destinations
	manager.AddDestination(&buf1)
	manager.AddDestination(&buf2)

	// Write data
	writer := manager.StartCapture()
	writer.Write([]byte("test1"))

	// Both should have received data
	if buf1.String() != "test1" {
		t.Errorf("buf1 = %q, want %q", buf1.String(), "test1")
	}
	if buf2.String() != "test1" {
		t.Errorf("buf2 = %q, want %q", buf2.String(), "test1")
	}

	// Remove one destination
	manager.RemoveDestination(&buf1)

	// Write more data
	writer.Write([]byte("test2"))

	// Only buf2 should have the new data
	if buf1.String() != "test1" {
		t.Errorf("buf1 should not have received new data, got %q", buf1.String())
	}
	if buf2.String() != "test1test2" {
		t.Errorf("buf2 = %q, want %q", buf2.String(), "test1test2")
	}
}

func TestStreamingManager_Stop(t *testing.T) {
	manager := NewStreamingManager(Options{})

	var buf bytes.Buffer
	manager.Stream(&buf)

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

	// Only the data before stop should be in buffer
	if buf.String() != "before stop" {
		t.Errorf("Buffer = %q, want %q", buf.String(), "before stop")
	}
}

func TestStreamingManager_Clear(t *testing.T) {
	manager := NewStreamingManager(Options{})

	var buf bytes.Buffer
	manager.Stream(&buf)

	// Clear should remove all destinations
	manager.Clear()

	// Write data after clear
	writer := manager.StartCapture()
	writer.Write([]byte("after clear"))

	// Buffer should be empty since destination was removed
	if buf.String() != "" {
		t.Errorf("Buffer should be empty after Clear, got %q", buf.String())
	}
}

func TestStreamingManager_ConcurrentAccess(t *testing.T) {
	manager := NewStreamingManager(Options{})

	var buf bytes.Buffer
	manager.Stream(&buf)

	writer := manager.StartCapture()

	// Simulate concurrent writes
	done := make(chan bool, 2)

	go func() {
		for i := 0; i < 100; i++ {
			writer.Write([]byte("A"))
			time.Sleep(time.Microsecond)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			writer.Write([]byte("B"))
			time.Sleep(time.Microsecond)
		}
		done <- true
	}()

	// Wait for both goroutines to complete
	<-done
	<-done

	output := buf.String()
	countA := strings.Count(output, "A")
	countB := strings.Count(output, "B")

	if countA != 100 {
		t.Errorf("Expected 100 'A's, got %d", countA)
	}
	if countB != 100 {
		t.Errorf("Expected 100 'B's, got %d", countB)
	}
}

func TestStreamingManager_WriteErrorHandling(t *testing.T) {
	manager := NewStreamingManager(Options{})

	// Create a writer that always fails
	failingWriter := &failingWriter{}
	manager.AddDestination(failingWriter)

	// Also add a working writer
	var buf bytes.Buffer
	manager.AddDestination(&buf)

	writer := manager.StartCapture()

	// Write should return error from failing writer but still write to working one
	_, err := writer.Write([]byte("test"))
	if err == nil {
		t.Error("Expected error from failing writer")
	}

	// Working writer should still have received the data
	if buf.String() != "test" {
		t.Errorf("Working writer should have data, got %q", buf.String())
	}
}

// failingWriter is a test helper that always returns an error on Write
type failingWriter struct{}

func (w *failingWriter) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("intentional write failure")
}
