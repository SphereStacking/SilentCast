package output

import (
	"io"
)

// Type represents the type of output manager
type Type int

const (
	// TypeBuffered stores all output in memory
	TypeBuffered Type = iota
	// TypeStreaming streams output in real-time
	TypeStreaming
	// TypeTee duplicates output to multiple destinations
	TypeTee
)

// Destination represents where output should be sent
type Destination int

const (
	// DestConsole outputs to console/terminal
	DestConsole Destination = iota
	// DestNotification sends output as system notification
	DestNotification
	// DestFile writes output to a file
	DestFile
	// DestWindow displays output in a window
	DestWindow
)

// Manager handles command output capture and display
type Manager interface {
	// StartCapture begins capturing output and returns a writer
	StartCapture() io.Writer

	// GetOutput returns all captured output as a string
	GetOutput() string

	// Stream sends output to the specified writer in real-time
	Stream(writer io.Writer) error

	// Stop stops capturing and cleans up resources
	Stop() error

	// Clear clears any buffered output
	Clear()
}

// Options configures output behavior
type Options struct {
	// Type specifies the output manager type
	Type Type

	// MaxSize limits the amount of output stored (0 = unlimited)
	MaxSize int

	// TruncateMessage is appended when output is truncated
	TruncateMessage string

	// Destinations specifies where output should be sent
	Destinations []Destination
}

// DefaultOptions returns default output options
func DefaultOptions() Options {
	return Options{
		Type:            TypeBuffered,
		MaxSize:         1024 * 1024, // 1MB default
		TruncateMessage: "\n... (output truncated)",
		Destinations:    []Destination{DestConsole},
	}
}
