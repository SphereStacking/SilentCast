package terminal

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// baseDetector provides common detection functionality
type baseDetector struct {
	// Platform-specific terminal definitions
	terminals []Terminal
}

// DetectTerminals finds all available terminal emulators
func (d *baseDetector) DetectTerminals() []Terminal {
	var available []Terminal

	for _, term := range d.terminals {
		if d.isTerminalAvailable(term) {
			available = append(available, term)
		}
	}

	return available
}

// FindTerminal searches for a specific terminal by name or command
func (d *baseDetector) FindTerminal(nameOrCommand string) (Terminal, error) {
	nameOrCommand = strings.ToLower(nameOrCommand)

	for _, term := range d.terminals {
		if strings.EqualFold(term.Name, nameOrCommand) ||
			strings.EqualFold(term.Command, nameOrCommand) {
			if d.isTerminalAvailable(term) {
				return term, nil
			}
		}
	}

	return Terminal{}, fmt.Errorf("terminal '%s' not found or not available", nameOrCommand)
}

// isTerminalAvailable checks if a terminal command exists in PATH
func (d *baseDetector) isTerminalAvailable(terminal Terminal) bool {
	_, err := exec.LookPath(terminal.Command)
	return err == nil
}

// Platform-specific detector implementations

// WindowsDetector detects terminals on Windows
type WindowsDetector struct {
	baseDetector
}

// NewWindowsDetector creates a detector for Windows terminals
func NewWindowsDetector() Detector {
	return &WindowsDetector{
		baseDetector: baseDetector{
			terminals: []Terminal{
				{
					Name:      "Windows Terminal",
					Command:   "wt.exe",
					Priority:  100,
					IsDefault: false,
					SupportedFeatures: TerminalFeatures{
						WindowSize:     true,
						WindowPosition: true,
						FontSize:       false,
						ColorScheme:    true,
						WindowState:    true,
						AlwaysOnTop:    false,
					},
				},
				{
					Name:      "Command Prompt",
					Command:   "cmd.exe",
					Priority:  50,
					IsDefault: true,
					SupportedFeatures: TerminalFeatures{
						WindowSize:     true,
						WindowPosition: true,
						FontSize:       false,
						ColorScheme:    false,
						WindowState:    false,
						AlwaysOnTop:    false,
					},
				},
				{
					Name:      "PowerShell",
					Command:   "powershell.exe",
					Priority:  80,
					IsDefault: false,
					SupportedFeatures: TerminalFeatures{
						WindowSize:     true,
						WindowPosition: true,
						FontSize:       false,
						ColorScheme:    false,
						WindowState:    false,
						AlwaysOnTop:    false,
					},
				},
				{
					Name:      "PowerShell Core",
					Command:   "pwsh.exe",
					Priority:  90,
					IsDefault: false,
					SupportedFeatures: TerminalFeatures{
						WindowSize:     true,
						WindowPosition: true,
						FontSize:       false,
						ColorScheme:    false,
						WindowState:    false,
						AlwaysOnTop:    false,
					},
				},
			},
		},
	}
}

// MacOSDetector detects terminals on macOS
type MacOSDetector struct {
	baseDetector
}

// NewMacOSDetector creates a detector for macOS terminals
func NewMacOSDetector() Detector {
	return &MacOSDetector{
		baseDetector: baseDetector{
			terminals: []Terminal{
				{
					Name:      "Terminal",
					Command:   "/System/Applications/Utilities/Terminal.app/Contents/MacOS/Terminal",
					Priority:  100,
					IsDefault: true,
					SupportedFeatures: TerminalFeatures{
						WindowSize:     true,
						WindowPosition: true,
						FontSize:       true,
						ColorScheme:    true,
						WindowState:    false,
						AlwaysOnTop:    false,
					},
				},
				{
					Name:      "iTerm2",
					Command:   "/Applications/iTerm.app/Contents/MacOS/iTerm2",
					Priority:  90,
					IsDefault: false,
					SupportedFeatures: TerminalFeatures{
						WindowSize:     true,
						WindowPosition: true,
						FontSize:       true,
						ColorScheme:    true,
						WindowState:    true,
						AlwaysOnTop:    true,
					},
				},
				{
					Name:      "Alacritty",
					Command:   "alacritty",
					Priority:  80,
					IsDefault: false,
					SupportedFeatures: TerminalFeatures{
						WindowSize:     true,
						WindowPosition: true,
						FontSize:       true,
						ColorScheme:    true,
						WindowState:    true,
						AlwaysOnTop:    false,
					},
				},
				{
					Name:      "kitty",
					Command:   "kitty",
					Priority:  75,
					IsDefault: false,
					SupportedFeatures: TerminalFeatures{
						WindowSize:     true,
						WindowPosition: true,
						FontSize:       true,
						ColorScheme:    true,
						WindowState:    true,
						AlwaysOnTop:    false,
					},
				},
			},
		},
	}
}

// LinuxDetector detects terminals on Linux
type LinuxDetector struct {
	baseDetector
}

// NewLinuxDetector creates a detector for Linux terminals
func NewLinuxDetector() Detector {
	// Check desktop environment to set appropriate defaults
	desktopEnv := strings.ToLower(os.Getenv("XDG_CURRENT_DESKTOP"))

	terminals := []Terminal{
		{
			Name:      "GNOME Terminal",
			Command:   "gnome-terminal",
			Priority:  100,
			IsDefault: strings.Contains(desktopEnv, "gnome"),
			SupportedFeatures: TerminalFeatures{
				WindowSize:     true,
				WindowPosition: true,
				FontSize:       true,
				ColorScheme:    true,
				WindowState:    true,
				AlwaysOnTop:    false,
			},
		},
		{
			Name:      "Konsole",
			Command:   "konsole",
			Priority:  100,
			IsDefault: strings.Contains(desktopEnv, "kde"),
			SupportedFeatures: TerminalFeatures{
				WindowSize:     true,
				WindowPosition: true,
				FontSize:       true,
				ColorScheme:    true,
				WindowState:    true,
				AlwaysOnTop:    false,
			},
		},
		{
			Name:      "xfce4-terminal",
			Command:   "xfce4-terminal",
			Priority:  100,
			IsDefault: strings.Contains(desktopEnv, "xfce"),
			SupportedFeatures: TerminalFeatures{
				WindowSize:     true,
				WindowPosition: true,
				FontSize:       true,
				ColorScheme:    false,
				WindowState:    true,
				AlwaysOnTop:    false,
			},
		},
		{
			Name:      "Terminator",
			Command:   "terminator",
			Priority:  80,
			IsDefault: false,
			SupportedFeatures: TerminalFeatures{
				WindowSize:     true,
				WindowPosition: true,
				FontSize:       true,
				ColorScheme:    true,
				WindowState:    true,
				AlwaysOnTop:    false,
			},
		},
		{
			Name:      "Alacritty",
			Command:   "alacritty",
			Priority:  85,
			IsDefault: false,
			SupportedFeatures: TerminalFeatures{
				WindowSize:     true,
				WindowPosition: true,
				FontSize:       true,
				ColorScheme:    true,
				WindowState:    true,
				AlwaysOnTop:    false,
			},
		},
		{
			Name:      "kitty",
			Command:   "kitty",
			Priority:  85,
			IsDefault: false,
			SupportedFeatures: TerminalFeatures{
				WindowSize:     true,
				WindowPosition: true,
				FontSize:       true,
				ColorScheme:    true,
				WindowState:    true,
				AlwaysOnTop:    false,
			},
		},
		{
			Name:      "urxvt",
			Command:   "urxvt",
			Priority:  70,
			IsDefault: false,
			SupportedFeatures: TerminalFeatures{
				WindowSize:     true,
				WindowPosition: true,
				FontSize:       true,
				ColorScheme:    false,
				WindowState:    false,
				AlwaysOnTop:    false,
			},
		},
		{
			Name:      "xterm",
			Command:   "xterm",
			Priority:  50,
			IsDefault: false,
			SupportedFeatures: TerminalFeatures{
				WindowSize:     true,
				WindowPosition: true,
				FontSize:       true,
				ColorScheme:    false,
				WindowState:    false,
				AlwaysOnTop:    false,
			},
		},
	}

	return &LinuxDetector{
		baseDetector: baseDetector{
			terminals: terminals,
		},
	}
}

// NewDetector creates a platform-specific detector
func NewDetector() Detector {
	switch runtime.GOOS {
	case "windows":
		return NewWindowsDetector()
	case "darwin":
		return NewMacOSDetector()
	case "linux":
		return NewLinuxDetector()
	default:
		// Return a detector with no terminals for unsupported platforms
		return &baseDetector{terminals: []Terminal{}}
	}
}
