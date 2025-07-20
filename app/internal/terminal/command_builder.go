package terminal

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// baseCommandBuilder provides common command building functionality
type baseCommandBuilder struct {
	platformBuilders map[string]CommandBuilder
}

// NewCommandBuilder creates a platform-specific command builder
func NewCommandBuilder() CommandBuilder {
	switch runtime.GOOS {
	case "windows":
		return NewWindowsCommandBuilder()
	case "darwin":
		return NewMacOSCommandBuilder()
	case "linux":
		return NewLinuxCommandBuilder()
	default:
		return &baseCommandBuilder{}
	}
}

// BuildCommand returns an error for unsupported platforms
func (b *baseCommandBuilder) BuildCommand(terminal Terminal, cmd *exec.Cmd, options Options) ([]string, error) {
	return nil, fmt.Errorf("command building not supported on this platform")
}

// SupportsTerminal returns false for unsupported platforms
func (b *baseCommandBuilder) SupportsTerminal(terminal Terminal) bool {
	return false
}

// WindowsCommandBuilder builds commands for Windows terminals
type WindowsCommandBuilder struct{}

// NewWindowsCommandBuilder creates a command builder for Windows
func NewWindowsCommandBuilder() CommandBuilder {
	return &WindowsCommandBuilder{}
}

// BuildCommand creates Windows terminal commands
func (b *WindowsCommandBuilder) BuildCommand(terminal Terminal, cmd *exec.Cmd, options Options) ([]string, error) {
	// Build the command string
	cmdStr := buildCommandString(cmd)

	switch strings.ToLower(terminal.Command) {
	case "wt.exe":
		// Windows Terminal
		args := []string{}
		if options.Title != "" {
			args = append(args, "--title", options.Title)
		}
		
		// Apply customization if provided and supported
		if options.Customization != nil && terminal.SupportedFeatures.WindowSize {
			if options.Customization.Width > 0 && options.Customization.Height > 0 {
				args = append(args, "--size", fmt.Sprintf("%d,%d", options.Customization.Width, options.Customization.Height))
			}
		}
		
		if options.Customization != nil && terminal.SupportedFeatures.WindowPosition {
			if options.Customization.X >= 0 && options.Customization.Y >= 0 {
				args = append(args, "--pos", fmt.Sprintf("%d,%d", options.Customization.X, options.Customization.Y))
			}
		}
		
		if options.Customization != nil && terminal.SupportedFeatures.WindowState {
			if options.Customization.Maximized {
				args = append(args, "--maximized")
			} else if options.Customization.Fullscreen {
				args = append(args, "--fullscreen")
			}
		}
		
		if options.Customization != nil && terminal.SupportedFeatures.ColorScheme {
			if options.Customization.Theme != "" {
				args = append(args, "--colorScheme", options.Customization.Theme)
			}
		}

		// Use cmd.exe to execute the command
		if options.KeepOpen {
			args = append(args, "cmd.exe", "/k", cmdStr)
		} else {
			args = append(args, "cmd.exe", "/c", cmdStr)
		}
		return args, nil

	case "cmd.exe":
		// Command Prompt
		startArgs := []string{"/c", "start"}
		
		// Add title
		startArgs = append(startArgs, fmt.Sprintf("\"%s\"", options.Title))
		
		// Apply window size and position if supported
		if options.Customization != nil && terminal.SupportedFeatures.WindowSize {
			if options.Customization.Width > 0 && options.Customization.Height > 0 {
				startArgs = append(startArgs, "/SIZE", fmt.Sprintf("%d,%d", options.Customization.Width, options.Customization.Height))
			}
		}
		
		if options.Customization != nil && terminal.SupportedFeatures.WindowPosition {
			if options.Customization.X >= 0 && options.Customization.Y >= 0 {
				startArgs = append(startArgs, "/POS", fmt.Sprintf("%d,%d", options.Customization.X, options.Customization.Y))
			}
		}
		
		startArgs = append(startArgs, "cmd")
		if options.KeepOpen {
			startArgs = append(startArgs, "/k", cmdStr)
		} else {
			startArgs = append(startArgs, "/c", cmdStr)
		}
		return startArgs, nil

	case "powershell.exe", "pwsh.exe":
		// PowerShell
		psCmd := cmdStr
		if options.KeepOpen {
			psCmd = fmt.Sprintf("%s; Read-Host 'Press Enter to exit'", cmdStr)
		}

		// Use -NoExit for PowerShell to keep window open
		if options.KeepOpen {
			return []string{"-NoExit", "-Command", cmdStr}, nil
		}
		return []string{"-Command", psCmd}, nil

	default:
		return nil, fmt.Errorf("unsupported Windows terminal: %s", terminal.Name)
	}
}

// SupportsTerminal checks if this builder supports the given terminal
func (b *WindowsCommandBuilder) SupportsTerminal(terminal Terminal) bool {
	switch strings.ToLower(terminal.Command) {
	case "wt.exe", "cmd.exe", "powershell.exe", "pwsh.exe":
		return true
	default:
		return false
	}
}

// MacOSCommandBuilder builds commands for macOS terminals
type MacOSCommandBuilder struct{}

// NewMacOSCommandBuilder creates a command builder for macOS
func NewMacOSCommandBuilder() CommandBuilder {
	return &MacOSCommandBuilder{}
}

// BuildCommand creates macOS terminal commands
func (b *MacOSCommandBuilder) BuildCommand(terminal Terminal, cmd *exec.Cmd, options Options) ([]string, error) {
	cmdStr := buildCommandString(cmd)

	switch {
	case strings.Contains(terminal.Command, "Terminal.app"):
		// macOS Terminal.app - use AppleScript
		script := cmdStr
		if options.KeepOpen {
			script = fmt.Sprintf("%s; echo ''; echo 'Press Enter to close...'; read", cmdStr)
		}

		// Escape quotes in the command
		script = strings.ReplaceAll(script, `"`, `\"`)
		
		// Build AppleScript with customization
		appleScript := `tell application "Terminal"`
		appleScript += fmt.Sprintf(`
			set newTab to do script "%s"`, script)
		
		// Apply customizations if supported
		if options.Customization != nil {
			if terminal.SupportedFeatures.WindowSize && options.Customization.Width > 0 && options.Customization.Height > 0 {
				appleScript += fmt.Sprintf(`
			set bounds of window 1 to {0, 0, %d, %d}`, options.Customization.Width*8, options.Customization.Height*16)
			}
			
			if terminal.SupportedFeatures.WindowPosition && options.Customization.X >= 0 && options.Customization.Y >= 0 {
				appleScript += fmt.Sprintf(`
			set position of window 1 to {%d, %d}`, options.Customization.X, options.Customization.Y)
			}
			
			if terminal.SupportedFeatures.FontSize && options.Customization.FontSize > 0 {
				appleScript += fmt.Sprintf(`
			set font size of window 1 to %d`, options.Customization.FontSize)
			}
		}
		
		appleScript += `
		end tell`

		// We need to use osascript to run AppleScript
		// Return args for: osascript -e '<script>'
		return []string{"-e", appleScript}, nil

	case strings.Contains(terminal.Command, "iTerm"):
		// iTerm2 - use AppleScript
		script := cmdStr
		if options.KeepOpen {
			script = fmt.Sprintf("%s; echo ''; echo 'Press Enter to close...'; read", cmdStr)
		}

		script = strings.ReplaceAll(script, `"`, `\"`)
		appleScript := fmt.Sprintf(`tell application "iTerm2"
			create window with default profile
			tell current session of current window
				write text "%s"
			end tell
		end tell`, script)

		return []string{"-e", appleScript}, nil

	case terminal.Command == "alacritty":
		// Alacritty
		args := []string{}
		if options.Title != "" {
			args = append(args, "--title", options.Title)
		}
		
		// Apply customizations if supported
		if options.Customization != nil {
			if terminal.SupportedFeatures.WindowSize && options.Customization.Width > 0 && options.Customization.Height > 0 {
				args = append(args, "--dimensions", fmt.Sprintf("%dx%d", options.Customization.Width, options.Customization.Height))
			}
			
			if terminal.SupportedFeatures.WindowPosition && options.Customization.X >= 0 && options.Customization.Y >= 0 {
				args = append(args, "--position", fmt.Sprintf("%d,%d", options.Customization.X, options.Customization.Y))
			}
			
			if terminal.SupportedFeatures.ColorScheme && options.Customization.Theme != "" {
				args = append(args, "--config-file", fmt.Sprintf("/path/to/themes/%s.yml", options.Customization.Theme))
			}
		}

		if options.KeepOpen {
			// Use shell to add a pause
			args = append(args, "-e", "sh", "-c", fmt.Sprintf("%s; echo ''; echo 'Press Enter to close...'; read", cmdStr))
		} else {
			args = append(args, "-e", "sh", "-c", cmdStr)
		}
		return args, nil

	case terminal.Command == "kitty":
		// kitty
		args := []string{}
		if options.Title != "" {
			args = append(args, "--title", options.Title)
		}

		if options.KeepOpen {
			args = append(args, "--hold")
		}
		args = append(args, "sh", "-c", cmdStr)
		return args, nil

	default:
		return nil, fmt.Errorf("unsupported macOS terminal: %s", terminal.Name)
	}
}

// SupportsTerminal checks if this builder supports the given terminal
func (b *MacOSCommandBuilder) SupportsTerminal(terminal Terminal) bool {
	switch {
	case strings.Contains(terminal.Command, "Terminal.app"),
		strings.Contains(terminal.Command, "iTerm"),
		terminal.Command == "alacritty",
		terminal.Command == "kitty":
		return true
	default:
		return false
	}
}

// LinuxCommandBuilder builds commands for Linux terminals
type LinuxCommandBuilder struct{}

// NewLinuxCommandBuilder creates a command builder for Linux
func NewLinuxCommandBuilder() CommandBuilder {
	return &LinuxCommandBuilder{}
}

// BuildCommand creates Linux terminal commands
func (b *LinuxCommandBuilder) BuildCommand(terminal Terminal, cmd *exec.Cmd, options Options) ([]string, error) {
	cmdStr := buildCommandString(cmd)

	switch terminal.Command {
	case "gnome-terminal":
		args := []string{}
		if options.Title != "" {
			args = append(args, "--title", options.Title)
		}
		
		// Apply customizations if supported
		if options.Customization != nil {
			if terminal.SupportedFeatures.WindowSize && options.Customization.Width > 0 && options.Customization.Height > 0 {
				args = append(args, "--geometry", fmt.Sprintf("%dx%d", options.Customization.Width, options.Customization.Height))
			}
			
			if terminal.SupportedFeatures.WindowState {
				if options.Customization.Maximized {
					args = append(args, "--maximize")
				} else if options.Customization.Fullscreen {
					args = append(args, "--full-screen")
				}
			}
		}

		if options.KeepOpen {
			// Use bash with read to keep terminal open
			args = append(args, "--", "bash", "-c", fmt.Sprintf("%s; echo ''; echo 'Press Enter to close...'; read", cmdStr))
		} else {
			args = append(args, "--", "sh", "-c", cmdStr)
		}
		return args, nil

	case "konsole":
		args := []string{}
		if options.Title != "" {
			args = append(args, "--title", options.Title)
		}
		
		// Apply customizations if supported
		if options.Customization != nil {
			if terminal.SupportedFeatures.WindowSize && options.Customization.Width > 0 && options.Customization.Height > 0 {
				args = append(args, "--geometry", fmt.Sprintf("%dx%d", options.Customization.Width, options.Customization.Height))
			}
			
			if terminal.SupportedFeatures.ColorScheme && options.Customization.Theme != "" {
				args = append(args, "--profile", options.Customization.Theme)
			}
		}

		if options.KeepOpen {
			args = append(args, "--hold")
		}
		args = append(args, "-e", "sh", "-c", cmdStr)
		return args, nil

	case "xfce4-terminal":
		args := []string{}
		if options.Title != "" {
			args = append(args, "--title", options.Title)
		}

		if options.KeepOpen {
			args = append(args, "--hold")
		}
		args = append(args, "-e", fmt.Sprintf("sh -c '%s'", cmdStr))
		return args, nil

	case "terminator":
		args := []string{}
		if options.Title != "" {
			args = append(args, "-T", options.Title)
		}

		if options.KeepOpen {
			args = append(args, "-e", fmt.Sprintf("bash -c '%s; echo \"\"; echo \"Press Enter to close...\"; read'", cmdStr))
		} else {
			args = append(args, "-e", fmt.Sprintf("sh -c '%s'", cmdStr))
		}
		return args, nil

	case "alacritty":
		args := []string{}
		if options.Title != "" {
			args = append(args, "--title", options.Title)
		}
		
		// Apply customizations if supported (same as macOS version)
		if options.Customization != nil {
			if terminal.SupportedFeatures.WindowSize && options.Customization.Width > 0 && options.Customization.Height > 0 {
				args = append(args, "--dimensions", fmt.Sprintf("%dx%d", options.Customization.Width, options.Customization.Height))
			}
			
			if terminal.SupportedFeatures.WindowPosition && options.Customization.X >= 0 && options.Customization.Y >= 0 {
				args = append(args, "--position", fmt.Sprintf("%d,%d", options.Customization.X, options.Customization.Y))
			}
			
			if terminal.SupportedFeatures.ColorScheme && options.Customization.Theme != "" {
				args = append(args, "--config-file", fmt.Sprintf("/path/to/themes/%s.yml", options.Customization.Theme))
			}
		}

		if options.KeepOpen {
			args = append(args, "-e", "sh", "-c", fmt.Sprintf("%s; echo ''; echo 'Press Enter to close...'; read", cmdStr))
		} else {
			args = append(args, "-e", "sh", "-c", cmdStr)
		}
		return args, nil

	case "kitty":
		args := []string{}
		if options.Title != "" {
			args = append(args, "--title", options.Title)
		}

		if options.KeepOpen {
			args = append(args, "--hold")
		}
		args = append(args, "sh", "-c", cmdStr)
		return args, nil

	case "urxvt":
		args := []string{}
		if options.Title != "" {
			args = append(args, "-title", options.Title)
		}

		if options.KeepOpen {
			args = append(args, "-hold", "-e", "sh", "-c", cmdStr)
		} else {
			args = append(args, "-e", "sh", "-c", cmdStr)
		}
		return args, nil

	case "xterm":
		args := []string{}
		if options.Title != "" {
			args = append(args, "-title", options.Title)
		}
		
		// Apply customizations if supported
		if options.Customization != nil {
			if terminal.SupportedFeatures.WindowSize && options.Customization.Width > 0 && options.Customization.Height > 0 {
				args = append(args, "-geometry", fmt.Sprintf("%dx%d", options.Customization.Width, options.Customization.Height))
			}
			
			if terminal.SupportedFeatures.FontSize && options.Customization.FontSize > 0 {
				args = append(args, "-fs", fmt.Sprintf("%d", options.Customization.FontSize))
			}
		}

		if options.KeepOpen {
			args = append(args, "-hold", "-e", cmdStr)
		} else {
			args = append(args, "-e", cmdStr)
		}
		return args, nil

	default:
		return nil, fmt.Errorf("unsupported Linux terminal: %s", terminal.Name)
	}
}

// SupportsTerminal checks if this builder supports the given terminal
func (b *LinuxCommandBuilder) SupportsTerminal(terminal Terminal) bool {
	switch terminal.Command {
	case "gnome-terminal", "konsole", "xfce4-terminal", "terminator",
		"alacritty", "kitty", "urxvt", "xterm":
		return true
	default:
		return false
	}
}

// buildCommandString constructs a command string from exec.Cmd
func buildCommandString(cmd *exec.Cmd) string {
	// Use Args which includes the command name as the first element
	parts := make([]string, len(cmd.Args))
	copy(parts, cmd.Args)

	// Quote parts that contain spaces
	for i, part := range parts {
		if strings.Contains(part, " ") {
			parts[i] = fmt.Sprintf(`"%s"`, part)
		}
	}

	return strings.Join(parts, " ")
}
