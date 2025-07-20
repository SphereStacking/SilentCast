package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestBufferedManager_BasicCapture(t *testing.T) {
	manager := NewBufferedManager(DefaultOptions())

	// Start capture
	writer := manager.StartCapture()
	if writer == nil {
		t.Fatal("StartCapture returned nil writer")
	}

	// Write some data
	testData := "Hello, World!"
	n, err := writer.Write([]byte(testData))
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	if n != len(testData) {
		t.Fatalf("Write returned %d, expected %d", n, len(testData))
	}

	// Get output
	output := manager.GetOutput()
	if output != testData {
		t.Fatalf("GetOutput returned %q, expected %q", output, testData)
	}
}

func TestBufferedManager_MultipleWrites(t *testing.T) {
	manager := NewBufferedManager(DefaultOptions())
	writer := manager.StartCapture()

	// Write multiple times
	writes := []string{"Line 1\n", "Line 2\n", "Line 3\n"}
	for _, data := range writes {
		_, err := writer.Write([]byte(data))
		if err != nil {
			t.Fatalf("Write failed: %v", err)
		}
	}

	// Check output
	expected := strings.Join(writes, "")
	output := manager.GetOutput()
	if output != expected {
		t.Fatalf("GetOutput returned %q, expected %q", output, expected)
	}
}

func TestBufferedManager_MaxSize(t *testing.T) {
	options := Options{
		Type:            TypeBuffered,
		MaxSize:         10,
		TruncateMessage: "...",
		Destinations:    []Destination{DestConsole},
	}

	manager := NewBufferedManager(options)
	writer := manager.StartCapture()

	// Write more than max size
	longData := "This is a very long string that exceeds the limit"
	n, err := writer.Write([]byte(longData))
	t.Logf("Write returned n=%d, err=%v", n, err)

	// Output should be truncated
	output := manager.GetOutput()
	t.Logf("Output: %q (length: %d)", output, len(output))
	if !strings.HasSuffix(output, "...") {
		t.Fatalf("Output not truncated: %q", output)
	}
	// The output should be exactly MaxSize (10) + TruncateMessage length (3) = 13
	expectedLen := options.MaxSize + len(options.TruncateMessage)
	if len(output) != expectedLen {
		t.Fatalf("Output wrong length: got %d chars, expected %d", len(output), expectedLen)
	}
}

func TestBufferedManager_Clear(t *testing.T) {
	manager := NewBufferedManager(DefaultOptions())
	writer := manager.StartCapture()

	// Write data
	writer.Write([]byte("Test data"))

	// Clear
	manager.Clear()

	// Output should be empty
	output := manager.GetOutput()
	if output != "" {
		t.Fatalf("GetOutput after Clear returned %q, expected empty", output)
	}
}

func TestBufferedManager_Stop(t *testing.T) {
	manager := NewBufferedManager(DefaultOptions())
	writer := manager.StartCapture()

	// Write data
	writer.Write([]byte("Before stop"))

	// Stop
	err := manager.Stop()
	if err != nil {
		t.Fatalf("Stop failed: %v", err)
	}

	// Writing after stop should fail
	_, err = writer.Write([]byte("After stop"))
	if err == nil {
		t.Fatal("Write after Stop should have failed")
	}

	// StartCapture after stop should return nil
	writer2 := manager.StartCapture()
	if writer2 != nil {
		t.Fatal("StartCapture after Stop should return nil")
	}
}

func TestBufferedManager_Stream(t *testing.T) {
	manager := NewBufferedManager(DefaultOptions())
	writer := manager.StartCapture()

	// Write data
	testData := "Stream test data"
	writer.Write([]byte(testData))

	// Stream to buffer
	var buf bytes.Buffer
	err := manager.Stream(&buf)
	if err != nil {
		t.Fatalf("Stream failed: %v", err)
	}

	// Check streamed data
	if buf.String() != testData {
		t.Fatalf("Streamed data %q doesn't match %q", buf.String(), testData)
	}
}

func TestBufferedManager_ConcurrentAccess(t *testing.T) {
	manager := NewBufferedManager(DefaultOptions())
	writer := manager.StartCapture()

	// Concurrent writes
	done := make(chan bool, 3)

	go func() {
		for i := 0; i < 100; i++ {
			writer.Write([]byte("A"))
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			writer.Write([]byte("B"))
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			_ = manager.GetOutput()
		}
		done <- true
	}()

	// Wait for all goroutines
	for i := 0; i < 3; i++ {
		<-done
	}

	// Check final output has both A and B
	output := manager.GetOutput()
	if !strings.Contains(output, "A") || !strings.Contains(output, "B") {
		t.Fatal("Concurrent writes failed")
	}
}
