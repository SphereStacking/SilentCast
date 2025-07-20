package terminal

import (
	"context"
	"os/exec"
)

// Manager defines the interface for managing terminal operations
type Manager interface {
	// ExecuteInTerminal executes a command in a new terminal window
	ExecuteInTerminal(ctx context.Context, cmd *exec.Cmd, options Options) error

	// GetAvailableTerminals returns a list of available terminal emulators
	GetAvailableTerminals() []Terminal

	// GetDefaultTerminal returns the default terminal for the platform
	GetDefaultTerminal() (Terminal, error)

	// IsTerminalAvailable checks if a specific terminal is available
	IsTerminalAvailable(terminal Terminal) bool
}

// Options represents options for terminal execution
type Options struct {
	// KeepOpen determines if the terminal should stay open after command completion
	KeepOpen bool

	// WorkingDir sets the working directory for the command
	WorkingDir string

	// Title sets the terminal window title (if supported)
	Title string

	// PreferredTerminal specifies a preferred terminal emulator
	PreferredTerminal Terminal

	// ForceTerminal forces terminal execution even in GUI/tray mode
	ForceTerminal bool
	
	// Customization contains visual customization options for the terminal window
	Customization *Customization
}

// Customization represents visual customization options for terminal windows
type Customization struct {
	// Window size
	Width  int `yaml:"width,omitempty"`  // Window width in columns
	Height int `yaml:"height,omitempty"` // Window height in rows
	
	// Window position
	X int `yaml:"x,omitempty"` // X position in pixels (if supported)
	Y int `yaml:"y,omitempty"` // Y position in pixels (if supported)
	
	// Appearance
	FontSize   int    `yaml:"font_size,omitempty"`   // Font size (if supported)
	Theme      string `yaml:"theme,omitempty"`       // Color theme/profile name
	Background string `yaml:"background,omitempty"`  // Background color (hex or name)
	Foreground string `yaml:"foreground,omitempty"`  // Text color (hex or name)
	
	// Behavior
	Fullscreen bool `yaml:"fullscreen,omitempty"` // Start in fullscreen mode
	Maximized  bool `yaml:"maximized,omitempty"`  // Start maximized
	AlwaysOnTop bool `yaml:"always_on_top,omitempty"` // Keep window always on top
}

// Terminal represents a terminal emulator
type Terminal struct {
	// Name is the display name of the terminal
	Name string

	// Command is the executable name/path
	Command string

	// Priority determines the preference order (higher is preferred)
	Priority int

	// IsDefault indicates if this is the system default terminal
	IsDefault bool
	
	// SupportedFeatures indicates which customization features this terminal supports
	SupportedFeatures TerminalFeatures
}

// TerminalFeatures represents which customization features a terminal supports
type TerminalFeatures struct {
	WindowSize     bool // Supports setting window size
	WindowPosition bool // Supports setting window position  
	FontSize       bool // Supports setting font size
	ColorScheme    bool // Supports color customization
	WindowState    bool // Supports fullscreen/maximized states
	AlwaysOnTop    bool // Supports always on top
}

// CommandBuilder creates terminal-specific command arguments
type CommandBuilder interface {
	// BuildCommand creates the command arguments for executing in a terminal
	BuildCommand(terminal Terminal, cmd *exec.Cmd, options Options) ([]string, error)

	// SupportsTerminal checks if this builder supports the given terminal
	SupportsTerminal(terminal Terminal) bool
}

// Detector finds available terminal emulators on the system
type Detector interface {
	// DetectTerminals finds all available terminal emulators
	DetectTerminals() []Terminal

	// FindTerminal searches for a specific terminal by name or command
	FindTerminal(nameOrCommand string) (Terminal, error)
}

// Error types for terminal operations
type Error struct {
	Op  string // Operation that failed
	Err error  // Underlying error
	Msg string // Additional context
}

func (e *Error) Error() string {
	if e.Msg != "" {
		return e.Op + ": " + e.Msg + ": " + e.Err.Error()
	}
	return e.Op + ": " + e.Err.Error()
}

func (e *Error) Unwrap() error {
	return e.Err
}

// Common error variables
var (
	// ErrNoTerminalFound indicates no suitable terminal emulator was found
	ErrNoTerminalFound = &Error{Op: "detect", Msg: "no terminal emulator found"}

	// ErrTerminalNotSupported indicates the terminal is not supported
	ErrTerminalNotSupported = &Error{Op: "build", Msg: "terminal not supported"}

	// ErrCommandFailed indicates the terminal command failed to execute
	ErrCommandFailed = &Error{Op: "execute", Msg: "command execution failed"}
)
