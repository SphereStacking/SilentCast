package updater

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestProgressReporter_Write(t *testing.T) {
	var buf bytes.Buffer
	pr := NewProgressReporter(1024, &buf)

	// Write some data
	n, err := pr.Write([]byte("test"))
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	if n != 4 {
		t.Errorf("Expected 4 bytes written, got %d", n)
	}

	// Check progress was reported
	output := buf.String()
	if !strings.Contains(output, "0.4%") {
		t.Errorf("Expected progress percentage, got: %s", output)
	}
}

func TestProgressReporter_UnknownTotal(t *testing.T) {
	var buf bytes.Buffer
	pr := NewProgressReporter(0, &buf)

	// Write some data
	_, err := pr.Write([]byte("test data"))
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Check output
	output := buf.String()
	if !strings.Contains(output, "Downloading...") {
		t.Errorf("Expected downloading message, got: %s", output)
	}
}

func TestProgressReporter_Done(t *testing.T) {
	var buf bytes.Buffer
	pr := NewProgressReporter(1024, &buf)

	// Write some data
	_, err := pr.Write([]byte("test"))
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Clear buffer
	buf.Reset()

	// Call done
	pr.Done()

	// Check final message
	output := buf.String()
	if !strings.Contains(output, "✅ Downloaded") {
		t.Errorf("Expected download complete message, got: %s", output)
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{0, "0 B"},
		{100, "100 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatBytes(tt.bytes)
			if result != tt.expected {
				t.Errorf("formatBytes(%d) = %s, want %s", tt.bytes, result, tt.expected)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{500 * time.Millisecond, "< 1s"},
		{5 * time.Second, "5s"},
		{65 * time.Second, "1m 5s"},
		{3665 * time.Second, "1h 1m"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatDuration(tt.duration)
			if result != tt.expected {
				t.Errorf("formatDuration(%v) = %s, want %s", tt.duration, result, tt.expected)
			}
		})
	}
}

func TestProgressWriter(t *testing.T) {
	var dataBuf bytes.Buffer
	var progressBuf bytes.Buffer

	pw := NewProgressWriter(&dataBuf, 100, &progressBuf)

	// Write some data
	testData := []byte("Hello, World!")
	n, err := pw.Write(testData)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	if n != len(testData) {
		t.Errorf("Expected %d bytes written, got %d", len(testData), n)
	}

	// Check data was written
	if dataBuf.String() != string(testData) {
		t.Errorf("Data not written correctly")
	}

	// Check progress was reported
	progressOutput := progressBuf.String()
	if !strings.Contains(progressOutput, "13.0%") {
		t.Errorf("Expected progress percentage, got: %s", progressOutput)
	}

	// Clear progress buffer
	progressBuf.Reset()

	// Call done
	pw.Done()

	// Check final message
	finalOutput := progressBuf.String()
	if !strings.Contains(finalOutput, "✅ Downloaded") {
		t.Errorf("Expected download complete message, got: %s", finalOutput)
	}
}
