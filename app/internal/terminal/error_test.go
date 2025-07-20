package terminal

import (
	"errors"
	"testing"
)

func TestError(t *testing.T) {
	tests := []struct {
		name     string
		err      *Error
		expected string
	}{
		{
			name: "error with message",
			err: &Error{
				Op:  "test",
				Err: errors.New("underlying error"),
				Msg: "additional context",
			},
			expected: "test: additional context: underlying error",
		},
		{
			name: "error without message",
			err: &Error{
				Op:  "test",
				Err: errors.New("underlying error"),
			},
			expected: "test: underlying error",
		},
		{
			name: "error with empty message",
			err: &Error{
				Op:  "test",
				Err: errors.New("underlying error"),
				Msg: "",
			},
			expected: "test: underlying error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestError_Unwrap(t *testing.T) {
	underlying := errors.New("underlying error")
	err := &Error{
		Op:  "test",
		Err: underlying,
		Msg: "context",
	}

	if unwrapped := err.Unwrap(); unwrapped != underlying {
		t.Errorf("Error.Unwrap() = %v, want %v", unwrapped, underlying)
	}
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		op   string
		msg  string
	}{
		{
			name: "ErrNoTerminalFound",
			err:  ErrNoTerminalFound,
			op:   "detect",
			msg:  "no terminal emulator found",
		},
		{
			name: "ErrTerminalNotSupported",
			err:  ErrTerminalNotSupported,
			op:   "build",
			msg:  "terminal not supported",
		},
		{
			name: "ErrCommandFailed",
			err:  ErrCommandFailed,
			op:   "execute",
			msg:  "command execution failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Op != tt.op {
				t.Errorf("Error.Op = %v, want %v", tt.err.Op, tt.op)
			}
			if tt.err.Msg != tt.msg {
				t.Errorf("Error.Msg = %v, want %v", tt.err.Msg, tt.msg)
			}
		})
	}
}

func TestTerminalFeatures(t *testing.T) {
	// Test that TerminalFeatures struct can be created with all features
	features := TerminalFeatures{
		WindowSize:     true,
		WindowPosition: true,
		FontSize:       true,
		ColorScheme:    true,
		WindowState:    true,
		AlwaysOnTop:    true,
	}

	if !features.WindowSize {
		t.Error("WindowSize should be true")
	}
	if !features.WindowPosition {
		t.Error("WindowPosition should be true")
	}
	if !features.FontSize {
		t.Error("FontSize should be true")
	}
	if !features.ColorScheme {
		t.Error("ColorScheme should be true")
	}
	if !features.WindowState {
		t.Error("WindowState should be true")
	}
	if !features.AlwaysOnTop {
		t.Error("AlwaysOnTop should be true")
	}
}

func TestTerminal(t *testing.T) {
	// Test Terminal struct creation and properties
	terminal := Terminal{
		Name:              "Test Terminal",
		Command:           "test-term",
		Priority:          100,
		IsDefault:         true,
		SupportedFeatures: TerminalFeatures{WindowSize: true},
	}

	if terminal.Name != "Test Terminal" {
		t.Errorf("Terminal.Name = %v, want %v", terminal.Name, "Test Terminal")
	}
	if terminal.Command != "test-term" {
		t.Errorf("Terminal.Command = %v, want %v", terminal.Command, "test-term")
	}
	if terminal.Priority != 100 {
		t.Errorf("Terminal.Priority = %v, want %v", terminal.Priority, 100)
	}
	if !terminal.IsDefault {
		t.Error("Terminal.IsDefault should be true")
	}
	if !terminal.SupportedFeatures.WindowSize {
		t.Error("Terminal.SupportedFeatures.WindowSize should be true")
	}
}

func TestCustomization(t *testing.T) {
	// Test Customization struct creation and properties
	custom := Customization{
		Width:       80,
		Height:      24,
		X:           100,
		Y:           200,
		FontSize:    12,
		Theme:       "dark",
		Background:  "#000000",
		Foreground:  "#ffffff",
		Fullscreen:  true,
		Maximized:   false,
		AlwaysOnTop: true,
	}

	if custom.Width != 80 {
		t.Errorf("Customization.Width = %v, want %v", custom.Width, 80)
	}
	if custom.Height != 24 {
		t.Errorf("Customization.Height = %v, want %v", custom.Height, 24)
	}
	if custom.X != 100 {
		t.Errorf("Customization.X = %v, want %v", custom.X, 100)
	}
	if custom.Y != 200 {
		t.Errorf("Customization.Y = %v, want %v", custom.Y, 200)
	}
	if custom.FontSize != 12 {
		t.Errorf("Customization.FontSize = %v, want %v", custom.FontSize, 12)
	}
	if custom.Theme != "dark" {
		t.Errorf("Customization.Theme = %v, want %v", custom.Theme, "dark")
	}
	if custom.Background != "#000000" {
		t.Errorf("Customization.Background = %v, want %v", custom.Background, "#000000")
	}
	if custom.Foreground != "#ffffff" {
		t.Errorf("Customization.Foreground = %v, want %v", custom.Foreground, "#ffffff")
	}
	if !custom.Fullscreen {
		t.Error("Customization.Fullscreen should be true")
	}
	if custom.Maximized {
		t.Error("Customization.Maximized should be false")
	}
	if !custom.AlwaysOnTop {
		t.Error("Customization.AlwaysOnTop should be true")
	}
}

func TestOptions(t *testing.T) {
	// Test Options struct creation and properties
	customization := &Customization{
		Width:  80,
		Height: 24,
	}

	options := Options{
		KeepOpen:          true,
		WorkingDir:        "/tmp",
		Title:             "Test Terminal",
		PreferredTerminal: Terminal{Name: "Test", Command: "test"},
		ForceTerminal:     true,
		Customization:     customization,
	}

	if !options.KeepOpen {
		t.Error("Options.KeepOpen should be true")
	}
	if options.WorkingDir != "/tmp" {
		t.Errorf("Options.WorkingDir = %v, want %v", options.WorkingDir, "/tmp")
	}
	if options.Title != "Test Terminal" {
		t.Errorf("Options.Title = %v, want %v", options.Title, "Test Terminal")
	}
	if options.PreferredTerminal.Name != "Test" {
		t.Errorf("Options.PreferredTerminal.Name = %v, want %v", options.PreferredTerminal.Name, "Test")
	}
	if !options.ForceTerminal {
		t.Error("Options.ForceTerminal should be true")
	}
	if options.Customization == nil {
		t.Error("Options.Customization should not be nil")
	}
	if options.Customization.Width != 80 {
		t.Errorf("Options.Customization.Width = %v, want %v", options.Customization.Width, 80)
	}
}