package output

import (
	"fmt"
	"io"
	"sync"
)

// StreamingManager streams output in real-time to multiple destinations
type StreamingManager struct {
	mu           sync.RWMutex
	destinations []io.Writer
	options      Options
	writer       *streamingWriter
	stopped      bool
}

// streamingWriter implements io.Writer for real-time streaming
type streamingWriter struct {
	manager *StreamingManager
}

// NewStreamingManager creates a new streaming output manager
func NewStreamingManager(options Options) *StreamingManager {
	// Apply defaults if options are empty
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

	return &StreamingManager{
		destinations: make([]io.Writer, 0),
		options:      options,
	}
}

// StartCapture begins capturing output and returns a writer
func (m *StreamingManager) StartCapture() io.Writer {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.stopped {
		return nil
	}

	m.writer = &streamingWriter{manager: m}
	return m.writer
}

// GetOutput returns empty string since streaming doesn't buffer
func (m *StreamingManager) GetOutput() string {
	// Streaming manager doesn't store output, return empty string
	return ""
}

// Stream adds a writer to the list of streaming destinations
func (m *StreamingManager) Stream(writer io.Writer) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.stopped {
		return fmt.Errorf("output manager is stopped")
	}

	// Add the writer to destinations
	m.destinations = append(m.destinations, writer)
	return nil
}

// Stop stops capturing and cleans up resources
func (m *StreamingManager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.stopped = true
	m.writer = nil
	m.destinations = nil
	return nil
}

// Clear clears the destinations list
func (m *StreamingManager) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.destinations = m.destinations[:0]
}

// AddDestination adds a new streaming destination
func (m *StreamingManager) AddDestination(writer io.Writer) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.stopped {
		m.destinations = append(m.destinations, writer)
	}
}

// RemoveDestination removes a streaming destination
func (m *StreamingManager) RemoveDestination(writer io.Writer) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, dest := range m.destinations {
		if dest == writer {
			// Remove by replacing with last element and shrinking slice
			m.destinations[i] = m.destinations[len(m.destinations)-1]
			m.destinations = m.destinations[:len(m.destinations)-1]
			break
		}
	}
}

// Write implements io.Writer interface for real-time streaming
func (w *streamingWriter) Write(p []byte) (n int, err error) {
	w.manager.mu.RLock()
	defer w.manager.mu.RUnlock()

	if w.manager.stopped {
		return 0, fmt.Errorf("output manager is stopped")
	}

	// Stream to all destinations simultaneously
	var lastErr error
	writtenLen := len(p)

	for _, dest := range w.manager.destinations {
		if dest != nil {
			_, writeErr := dest.Write(p)
			if writeErr != nil {
				lastErr = writeErr
			}
		}
	}

	// Return the original length even if some destinations failed
	// This prevents the source from thinking the write failed
	return writtenLen, lastErr
}
