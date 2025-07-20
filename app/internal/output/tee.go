package output

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

// TeeOutputManager duplicates output to multiple destinations
// It combines features of BufferedManager and StreamingManager to provide
// both real-time streaming and buffered output retrieval.
//
// Example usage:
//
//	// Create a tee manager that streams to console and buffers for later
//	options := Options{Type: TypeTee}
//	manager := NewTeeManager(options)
//
//	// Add destinations
//	manager.AddDestination(os.Stdout)     // Stream to console
//	manager.AddDestination(&logFile)      // Stream to log file
//
//	// Start capturing
//	writer := manager.StartCapture()
//	cmd.Stdout = writer
//	cmd.Run()
//
//	// Get buffered output
//	output := manager.GetOutput()
type TeeManager struct {
	mu           sync.RWMutex
	destinations []io.Writer
	buffer       *bytes.Buffer
	options      Options
	writer       *teeWriter
	stopped      bool
}

// teeWriter implements io.Writer for the TeeManager
type teeWriter struct {
	manager *TeeManager
}

// NewTeeManager creates a new tee output manager
func NewTeeManager(options Options) *TeeManager {
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

	return &TeeManager{
		destinations: make([]io.Writer, 0),
		buffer:       &bytes.Buffer{},
		options:      options,
	}
}

// StartCapture begins capturing output and returns a writer
func (m *TeeManager) StartCapture() io.Writer {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.stopped {
		return nil
	}

	m.writer = &teeWriter{manager: m}
	return m.writer
}

// GetOutput returns all captured output as a string
func (m *TeeManager) GetOutput() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.buffer.String()
}

// Stream adds a writer to the list of streaming destinations
func (m *TeeManager) Stream(writer io.Writer) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.stopped {
		return fmt.Errorf("output manager is stopped")
	}

	// Add the writer to destinations
	m.destinations = append(m.destinations, writer)

	// Also stream any existing buffered content
	if m.buffer.Len() > 0 {
		_, err := io.Copy(writer, bytes.NewReader(m.buffer.Bytes()))
		if err != nil {
			return fmt.Errorf("failed to stream existing buffer: %w", err)
		}
	}

	return nil
}

// Stop stops capturing and cleans up resources
func (m *TeeManager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.stopped = true
	m.writer = nil
	// Keep destinations and buffer for GetOutput() calls
	return nil
}

// Clear clears both the buffer and destinations list
func (m *TeeManager) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.buffer.Reset()
	m.destinations = m.destinations[:0]
}

// AddDestination adds a new output destination
// This allows dynamic addition of destinations during capture
func (m *TeeManager) AddDestination(writer io.Writer) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.stopped {
		return fmt.Errorf("cannot add destination: output manager is stopped")
	}

	m.destinations = append(m.destinations, writer)

	// Stream existing buffer to new destination
	if m.buffer.Len() > 0 {
		_, err := io.Copy(writer, bytes.NewReader(m.buffer.Bytes()))
		if err != nil {
			// Remove the destination if initial sync fails
			m.destinations = m.destinations[:len(m.destinations)-1]
			return fmt.Errorf("failed to sync existing buffer to new destination: %w", err)
		}
	}

	return nil
}

// RemoveDestination removes an output destination
func (m *TeeManager) RemoveDestination(writer io.Writer) {
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

// GetDestinationCount returns the number of active destinations
func (m *TeeManager) GetDestinationCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.destinations)
}

// Write implements io.Writer interface
func (w *teeWriter) Write(p []byte) (n int, err error) {
	w.manager.mu.Lock()
	defer w.manager.mu.Unlock()

	if w.manager.stopped {
		return 0, fmt.Errorf("output manager is stopped")
	}

	// Always write to buffer first (within size limits)
	bufferN, bufferErr := w.writeToBuffer(p)

	// Then write to all destinations
	streamErr := w.streamToDestinations(p)

	// Return the original length to indicate success to the writer
	// even if some destinations failed
	n = len(p)

	// Return error only if buffer write failed
	if bufferErr != nil {
		return bufferN, bufferErr
	}

	// For streaming errors, we return success but could log them
	_ = streamErr

	return n, nil
}

// writeToBuffer writes data to the internal buffer with size limits
func (w *teeWriter) writeToBuffer(p []byte) (int, error) {
	// Check size limit
	if w.manager.options.MaxSize > 0 {
		currentSize := w.manager.buffer.Len()
		if currentSize >= w.manager.options.MaxSize {
			// Already at max size, don't write more
			return len(p), nil // Pretend we wrote it
		}

		// Calculate how much we can write
		remaining := w.manager.options.MaxSize - currentSize
		if len(p) > remaining {
			// Write what we can
			return w.manager.buffer.Write(p[:remaining])
		}
	}

	return w.manager.buffer.Write(p)
}

// streamToDestinations writes data to all registered destinations
func (w *teeWriter) streamToDestinations(p []byte) error {
	var lastErr error
	successCount := 0

	// Create a copy of destinations to avoid holding lock during IO
	destinations := make([]io.Writer, len(w.manager.destinations))
	copy(destinations, w.manager.destinations)

	// Release lock before doing IO operations
	w.manager.mu.Unlock()
	defer w.manager.mu.Lock()

	for _, dest := range destinations {
		if dest != nil {
			_, err := dest.Write(p)
			if err != nil {
				lastErr = err
			} else {
				successCount++
			}
		}
	}

	// Return error only if all destinations failed
	if successCount == 0 && len(destinations) > 0 {
		return fmt.Errorf("all destinations failed, last error: %w", lastErr)
	}

	return nil
}
