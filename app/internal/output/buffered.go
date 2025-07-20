package output

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

// BufferedManager stores all output in memory
type BufferedManager struct {
	mu        sync.RWMutex
	buffer    *bytes.Buffer
	options   Options
	writer    *bufferedWriter
	stopped   bool
	truncated bool
}

// bufferedWriter implements io.Writer for capturing output
type bufferedWriter struct {
	manager *BufferedManager
}

// NewBufferedManager creates a new buffered output manager
func NewBufferedManager(options Options) *BufferedManager {
	// Don't override the entire options struct if Type is already set correctly
	// Only use defaults if options are completely empty
	if options.MaxSize == 0 && options.TruncateMessage == "" {
		defaults := DefaultOptions()
		if options.MaxSize == 0 {
			options.MaxSize = defaults.MaxSize
		}
		if options.TruncateMessage == "" {
			options.TruncateMessage = defaults.TruncateMessage
		}
		if len(options.Destinations) == 0 {
			options.Destinations = defaults.Destinations
		}
	}

	return &BufferedManager{
		buffer:  &bytes.Buffer{},
		options: options,
	}
}

// StartCapture begins capturing output and returns a writer
func (m *BufferedManager) StartCapture() io.Writer {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.stopped {
		return nil
	}

	m.writer = &bufferedWriter{manager: m}
	return m.writer
}

// GetOutput returns all captured output as a string
func (m *BufferedManager) GetOutput() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	output := m.buffer.String()
	if m.truncated {
		output += m.options.TruncateMessage
	}
	return output
}

// Stream sends output to the specified writer
func (m *BufferedManager) Stream(writer io.Writer) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.stopped {
		return fmt.Errorf("output manager is stopped")
	}

	_, err := io.Copy(writer, bytes.NewReader(m.buffer.Bytes()))
	return err
}

// Stop stops capturing and cleans up resources
func (m *BufferedManager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.stopped = true
	m.writer = nil
	return nil
}

// Clear clears any buffered output
func (m *BufferedManager) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.buffer.Reset()
	m.truncated = false
}

// Write implements io.Writer interface
func (w *bufferedWriter) Write(p []byte) (n int, err error) {
	w.manager.mu.Lock()
	defer w.manager.mu.Unlock()

	if w.manager.stopped {
		return 0, fmt.Errorf("output manager is stopped")
	}

	// Check size limit
	if w.manager.options.MaxSize > 0 {
		currentSize := w.manager.buffer.Len()
		if currentSize >= w.manager.options.MaxSize {
			// Already at max size, don't write more
			w.manager.truncated = true
			return len(p), nil // Pretend we wrote it to avoid errors
		}

		// Calculate how much we can write
		remaining := w.manager.options.MaxSize - currentSize
		if len(p) > remaining {
			// Write what we can
			w.manager.buffer.Write(p[:remaining])
			w.manager.truncated = true
			return len(p), nil // Return original length to avoid errors
		}
	}

	return w.manager.buffer.Write(p)
}
