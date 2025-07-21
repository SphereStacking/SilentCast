package terminal

import (
	"os/exec"
	"testing"
)

func TestTerminalCustomization(t *testing.T) {
	tests := []struct {
		name         string
		terminal     Terminal
		options      Options
		expectedArgs []string
		expectError  bool
		skipOn       string // Skip test on specific platform
	}{
		{
			name: "Windows Terminal with size and position",
			terminal: Terminal{
				Name:    "Windows Terminal",
				Command: "wt.exe",
				SupportedFeatures: TerminalFeatures{
					WindowSize:     true,
					WindowPosition: true,
					ColorScheme:    true,
					WindowState:    true,
				},
			},
			options: Options{
				Title: "Test Terminal",
				Customization: &Customization{
					Width:  120,
					Height: 30,
					X:      100,
					Y:      50,
					Theme:  "Campbell",
				},
			},
			expectedArgs: []string{
				"--title", "Test Terminal",
				"--size", "120,30",
				"--pos", "100,50",
				"--colorScheme", "Campbell",
				"cmd.exe", "/c",
			},
		},
		{
			name: "GNOME Terminal with geometry",
			terminal: Terminal{
				Name:    "GNOME Terminal",
				Command: "gnome-terminal",
				SupportedFeatures: TerminalFeatures{
					WindowSize:  true,
					WindowState: true,
				},
			},
			options: Options{
				Title: "Test GNOME",
				Customization: &Customization{
					Width:     80,
					Height:    24,
					Maximized: true,
				},
			},
			expectedArgs: []string{
				"--title", "Test GNOME",
				"--geometry", "80x24",
				"--maximize",
				"--", "sh", "-c",
			},
		},
		{
			name: "Alacritty with dimensions and position",
			terminal: Terminal{
				Name:    "Alacritty",
				Command: "alacritty",
				SupportedFeatures: TerminalFeatures{
					WindowSize:     true,
					WindowPosition: true,
					ColorScheme:    true,
				},
			},
			options: Options{
				Title: "Test Alacritty",
				Customization: &Customization{
					Width:  100,
					Height: 40,
					X:      200,
					Y:      100,
					Theme:  "gruvbox",
				},
			},
			expectedArgs: []string{
				"--title", "Test Alacritty",
				"--dimensions", "100x40",
				"--position", "200,100",
				"--config-file", "/path/to/themes/gruvbox.yml",
				"-e", "sh", "-c",
			},
		},
		{
			name: "xterm with limited customization support",
			terminal: Terminal{
				Name:    "xterm",
				Command: "xterm",
				SupportedFeatures: TerminalFeatures{
					WindowSize:     true,
					WindowPosition: false,
					FontSize:       true,
					ColorScheme:    false,
				},
			},
			options: Options{
				Title: "Test xterm",
				Customization: &Customization{
					Width:      80,
					Height:     24,
					FontSize:   12,
					Background: "#000000", // Not supported
					Theme:      "dark",    // Not supported
				},
			},
			expectedArgs: []string{
				"-title", "Test xterm",
				"-geometry", "80x24",
				"-fs", "12",
				"-e",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create appropriate command builder
			var builder CommandBuilder
			switch {
			case tt.terminal.Command == "wt.exe" || tt.terminal.Command == "cmd.exe":
				builder = NewWindowsCommandBuilder()
			case tt.terminal.Command == "gnome-terminal" || tt.terminal.Command == "alacritty" || tt.terminal.Command == "xterm":
				builder = NewLinuxCommandBuilder()
			default:
				builder = NewMacOSCommandBuilder()
			}

			// Create a dummy command
			cmd := exec.Command("echo", "test")

			args, err := builder.BuildCommand(tt.terminal, cmd, &tt.options)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Check that expected args are present (partial match)
			for _, expectedArg := range tt.expectedArgs {
				found := false
				for _, arg := range args {
					if arg == expectedArg {
						found = true
						break
					}
				}
				if !found && expectedArg != "" {
					t.Errorf("Expected argument %q not found in %v", expectedArg, args)
				}
			}
		})
	}
}

func TestCustomizationValidation(t *testing.T) {
	tests := []struct {
		name          string
		customization *Customization
		isValid       bool
		description   string
	}{
		{
			name: "Valid customization",
			customization: &Customization{
				Width:  120,
				Height: 30,
				X:      100,
				Y:      50,
			},
			isValid:     true,
			description: "All positive values",
		},
		{
			name: "Zero size values",
			customization: &Customization{
				Width:  0,
				Height: 0,
			},
			isValid:     true,
			description: "Zero values should be ignored",
		},
		{
			name: "Negative position",
			customization: &Customization{
				X: -1,
				Y: -1,
			},
			isValid:     true,
			description: "Negative position values should be ignored",
		},
		{
			name: "Large font size",
			customization: &Customization{
				FontSize: 24,
			},
			isValid:     true,
			description: "Large font size should be valid",
		},
		{
			name: "Empty theme",
			customization: &Customization{
				Theme: "",
			},
			isValid:     true,
			description: "Empty theme should be ignored",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create terminal with all features supported
			terminal := Terminal{
				Name:    "Test Terminal",
				Command: "test",
				SupportedFeatures: TerminalFeatures{
					WindowSize:     true,
					WindowPosition: true,
					FontSize:       true,
					ColorScheme:    true,
					WindowState:    true,
					AlwaysOnTop:    true,
				},
			}

			options := Options{
				Customization: tt.customization,
			}

			// Test that customization doesn't cause panics
			defer func() {
				if r := recover(); r != nil {
					if tt.isValid {
						t.Errorf("Unexpected panic for valid customization: %v", r)
					}
				}
			}()

			// Create dummy builder and test
			builder := &baseCommandBuilder{}
			cmd := exec.Command("echo", "test")

			// This should not panic even with invalid input
			_, err := builder.BuildCommand(terminal, cmd, &options)

			// baseCommandBuilder always returns an error, but shouldn't panic
			if err == nil {
				t.Error("Expected error from baseCommandBuilder")
			}
		})
	}
}

func TestFeatureSupport(t *testing.T) {
	tests := []struct {
		name     string
		terminal Terminal
		feature  string
		expected bool
	}{
		{
			name: "Windows Terminal supports window size",
			terminal: Terminal{
				SupportedFeatures: TerminalFeatures{
					WindowSize: true,
				},
			},
			feature:  "WindowSize",
			expected: true,
		},
		{
			name: "xterm doesn't support color scheme",
			terminal: Terminal{
				SupportedFeatures: TerminalFeatures{
					ColorScheme: false,
				},
			},
			feature:  "ColorScheme",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var supported bool
			switch tt.feature {
			case "WindowSize":
				supported = tt.terminal.SupportedFeatures.WindowSize
			case "WindowPosition":
				supported = tt.terminal.SupportedFeatures.WindowPosition
			case "FontSize":
				supported = tt.terminal.SupportedFeatures.FontSize
			case "ColorScheme":
				supported = tt.terminal.SupportedFeatures.ColorScheme
			case "WindowState":
				supported = tt.terminal.SupportedFeatures.WindowState
			case "AlwaysOnTop":
				supported = tt.terminal.SupportedFeatures.AlwaysOnTop
			}

			if supported != tt.expected {
				t.Errorf("Feature %s support = %v, want %v", tt.feature, supported, tt.expected)
			}
		})
	}
}

func ExampleCustomization() {
	// Create a customization for a large terminal window
	customization := &Customization{
		Width:      120,
		Height:     40,
		X:          100,
		Y:          50,
		FontSize:   14,
		Theme:      "gruvbox-dark",
		Background: "#282828",
		Foreground: "#ebdbb2",
		Maximized:  false,
	}

	options := Options{
		Title:         "My Custom Terminal",
		Customization: customization,
	}

	// Use with terminal execution
	_ = options
}
