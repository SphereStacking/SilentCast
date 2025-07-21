package terminal

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
)

// mockFullDetector is a complete mock detector for testing
type mockFullDetector struct {
	terminals []Terminal
}

func (d *mockFullDetector) DetectTerminals() []Terminal {
	// Return all terminals as available
	return d.terminals
}

func (d *mockFullDetector) FindTerminal(nameOrCommand string) (Terminal, error) {
	nameOrCommand = strings.ToLower(nameOrCommand)
	for _, term := range d.terminals {
		if strings.EqualFold(term.Name, nameOrCommand) ||
			strings.EqualFold(term.Command, nameOrCommand) {
			return term, nil
		}
	}
	return Terminal{}, fmt.Errorf("terminal '%s' not found or not available", nameOrCommand)
}

func TestBaseDetector_FindTerminal(t *testing.T) {
	detector := &mockFullDetector{
		terminals: []Terminal{
			{Name: "Terminal App", Command: "terminal"},
			{Name: "XTerm", Command: "xterm"},
		},
	}

	tests := []struct {
		name        string
		search      string
		wantCommand string
		wantErr     bool
	}{
		{
			name:        "find by name",
			search:      "Terminal App",
			wantCommand: "terminal",
			wantErr:     false,
		},
		{
			name:        "find by name case insensitive",
			search:      "terminal app",
			wantCommand: "terminal",
			wantErr:     false,
		},
		{
			name:        "find by command",
			search:      "xterm",
			wantCommand: "xterm",
			wantErr:     false,
		},
		{
			name:        "not found",
			search:      "notexist",
			wantCommand: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := detector.FindTerminal(tt.search)
			if tt.wantErr {
				if err == nil {
					t.Errorf("FindTerminal() error = nil, wantErr %v", tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("FindTerminal() unexpected error = %v", err)
				return
			}
			if got.Command != tt.wantCommand {
				t.Errorf("FindTerminal() = %v, want %v", got.Command, tt.wantCommand)
			}
		})
	}
}

func TestNewDetector(t *testing.T) {
	detector := NewDetector()

	switch runtime.GOOS {
	case "windows":
		if _, ok := detector.(*WindowsDetector); !ok {
			t.Errorf("Expected WindowsDetector on Windows, got %T", detector)
		}
	case "darwin":
		if _, ok := detector.(*MacOSDetector); !ok {
			t.Errorf("Expected MacOSDetector on macOS, got %T", detector)
		}
	case "linux":
		if _, ok := detector.(*LinuxDetector); !ok {
			t.Errorf("Expected LinuxDetector on Linux, got %T", detector)
		}
	}
}

func TestWindowsDetector_Terminals(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows-specific test")
	}

	detector := NewWindowsDetector()
	terminals := detector.DetectTerminals()

	// Check that cmd.exe is always available on Windows
	foundCmd := false
	for _, term := range terminals {
		if term.Command == "cmd.exe" {
			foundCmd = true
			if !term.IsDefault {
				t.Error("cmd.exe should be marked as default on Windows")
			}
			break
		}
	}
	if !foundCmd {
		t.Error("cmd.exe should always be available on Windows")
	}
}

func TestMacOSDetector_Terminals(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("macOS-specific test")
	}

	detector := NewMacOSDetector()
	bd := detector.(*MacOSDetector).baseDetector

	// Check terminal definitions
	foundTerminalApp := false
	for _, term := range bd.terminals {
		if strings.Contains(term.Command, "Terminal.app") {
			foundTerminalApp = true
			if !term.IsDefault {
				t.Error("Terminal.app should be marked as default on macOS")
			}
			if term.Priority != 100 {
				t.Error("Terminal.app should have highest priority")
			}
			break
		}
	}
	if !foundTerminalApp {
		t.Error("Terminal.app should be defined for macOS")
	}
}

func TestLinuxDetector_Terminals(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Linux-specific test")
	}

	detector := NewLinuxDetector()
	bd := detector.(*LinuxDetector).baseDetector

	// Check that common terminals are defined
	expectedTerminals := []string{"gnome-terminal", "konsole", "xterm", "xfce4-terminal"}
	for _, expected := range expectedTerminals {
		found := false
		for _, term := range bd.terminals {
			if term.Command == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected terminal %s not found in Linux detector", expected)
		}
	}

	// Check that xterm has lower priority as fallback
	for _, term := range bd.terminals {
		if term.Command == "xterm" && term.Priority >= 80 {
			t.Error("xterm should have lower priority as it's a fallback option")
		}
	}
}

func TestBaseDetector_FindTerminal_RealImplementation(t *testing.T) {
	// Create a baseDetector with test terminals
	detector := &baseDetector{
		terminals: []Terminal{
			{Name: "Test Terminal", Command: "sh", Priority: 10}, // Using 'sh' as it exists
			{Name: "XTerm", Command: "nonexistent-xterm", Priority: 5},
		},
	}

	tests := []struct {
		name          string
		searchTerm    string
		expectFound   bool
		expectCommand string
	}{
		{
			name:          "find by name - available terminal",
			searchTerm:    "Test Terminal",
			expectFound:   true,
			expectCommand: "sh",
		},
		{
			name:          "find by command - available terminal",
			searchTerm:    "sh",
			expectFound:   true,
			expectCommand: "sh",
		},
		{
			name:          "find by name case insensitive",
			searchTerm:    "test terminal",
			expectFound:   true,
			expectCommand: "sh",
		},
		{
			name:        "terminal not available in PATH",
			searchTerm:  "XTerm",
			expectFound: false,
		},
		{
			name:        "terminal not found",
			searchTerm:  "Unknown Terminal",
			expectFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			term, err := detector.FindTerminal(tt.searchTerm)

			if tt.expectFound {
				if err != nil {
					t.Errorf("FindTerminal() error = %v, want no error", err)
				}
				if term.Command != tt.expectCommand {
					t.Errorf("FindTerminal() command = %v, want %v", term.Command, tt.expectCommand)
				}
			} else {
				if err == nil {
					t.Error("FindTerminal() expected error but got none")
				}
				if !strings.Contains(err.Error(), "not found") && !strings.Contains(err.Error(), "not available") {
					t.Errorf("FindTerminal() error = %v, expected 'not found' or 'not available'", err)
				}
			}
		})
	}
}
