package terminal

import (
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

func TestNewCommandBuilder(t *testing.T) {
	builder := NewCommandBuilder()

	switch runtime.GOOS {
	case "windows":
		if _, ok := builder.(*WindowsCommandBuilder); !ok {
			t.Errorf("Expected WindowsCommandBuilder on Windows, got %T", builder)
		}
	case "darwin":
		if _, ok := builder.(*MacOSCommandBuilder); !ok {
			t.Errorf("Expected MacOSCommandBuilder on macOS, got %T", builder)
		}
	case "linux":
		if _, ok := builder.(*LinuxCommandBuilder); !ok {
			t.Errorf("Expected LinuxCommandBuilder on Linux, got %T", builder)
		}
	}
}

func TestBuildCommandString(t *testing.T) {
	tests := []struct {
		name string
		cmd  *exec.Cmd
		want string
	}{
		{
			name: "simple command",
			cmd:  exec.Command("echo", "hello"),
			want: "echo hello",
		},
		{
			name: "command with spaces",
			cmd:  exec.Command("echo", "hello world"),
			want: `echo "hello world"`,
		},
		{
			name: "multiple args with spaces",
			cmd:  exec.Command("cmd", "/c", "echo hello", "world test"),
			want: `cmd /c "echo hello" "world test"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildCommandString(tt.cmd)
			if got != tt.want {
				t.Errorf("buildCommandString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestWindowsCommandBuilder_BuildCommand(t *testing.T) {
	builder := NewWindowsCommandBuilder()
	cmd := exec.Command("echo", "test")

	tests := []struct {
		name      string
		terminal  Terminal
		options   Options
		wantParts []string
		wantErr   bool
	}{
		{
			name:      "cmd.exe with keep open",
			terminal:  Terminal{Command: "cmd.exe"},
			options:   Options{KeepOpen: true, Title: "Test"},
			wantParts: []string{"/c", "start", "cmd", "/k"},
			wantErr:   false,
		},
		{
			name:      "cmd.exe without keep open",
			terminal:  Terminal{Command: "cmd.exe"},
			options:   Options{KeepOpen: false},
			wantParts: []string{"/c", "start", "cmd", "/c"},
			wantErr:   false,
		},
		{
			name:      "Windows Terminal with keep open",
			terminal:  Terminal{Command: "wt.exe"},
			options:   Options{KeepOpen: true, Title: "Test"},
			wantParts: []string{"--title", "Test", "cmd.exe", "/k"},
			wantErr:   false,
		},
		{
			name:      "PowerShell with keep open",
			terminal:  Terminal{Command: "powershell.exe"},
			options:   Options{KeepOpen: true},
			wantParts: []string{"-NoExit", "-Command"},
			wantErr:   false,
		},
		{
			name:     "unsupported terminal",
			terminal: Terminal{Command: "unknown.exe"},
			options:  Options{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := builder.BuildCommand(tt.terminal, cmd, &tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			// Check that expected parts are in the result
			gotStr := strings.Join(got, " ")
			for _, part := range tt.wantParts {
				if !strings.Contains(gotStr, part) {
					t.Errorf("BuildCommand() result %q missing expected part %q", gotStr, part)
				}
			}
		})
	}
}

func TestWindowsCommandBuilder_SupportsTerminal(t *testing.T) {
	builder := NewWindowsCommandBuilder()

	tests := []struct {
		terminal Terminal
		want     bool
	}{
		{Terminal{Command: "cmd.exe"}, true},
		{Terminal{Command: "wt.exe"}, true},
		{Terminal{Command: "powershell.exe"}, true},
		{Terminal{Command: "pwsh.exe"}, true},
		{Terminal{Command: "unknown.exe"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.terminal.Command, func(t *testing.T) {
			if got := builder.SupportsTerminal(tt.terminal); got != tt.want {
				t.Errorf("SupportsTerminal(%q) = %v, want %v", tt.terminal.Command, got, tt.want)
			}
		})
	}
}

func TestMacOSCommandBuilder_BuildCommand(t *testing.T) {
	builder := NewMacOSCommandBuilder()
	cmd := exec.Command("echo", "test")

	tests := []struct {
		name      string
		terminal  Terminal
		options   Options
		wantCheck func([]string) bool
		wantErr   bool
	}{
		{
			name:     "Terminal.app with keep open",
			terminal: Terminal{Command: "/System/Applications/Utilities/Terminal.app/Contents/MacOS/Terminal"},
			options:  Options{KeepOpen: true},
			wantCheck: func(args []string) bool {
				// Should have AppleScript with read command
				return len(args) == 2 && args[0] == "-e" &&
					strings.Contains(args[1], "tell application \"Terminal\"") &&
					strings.Contains(args[1], "read")
			},
			wantErr: false,
		},
		{
			name:     "alacritty with title",
			terminal: Terminal{Command: "alacritty"},
			options:  Options{Title: "Test", KeepOpen: false},
			wantCheck: func(args []string) bool {
				// Should have title and execute flags
				return contains(args, "--title") && contains(args, "Test") && contains(args, "-e")
			},
			wantErr: false,
		},
		{
			name:     "kitty with keep open",
			terminal: Terminal{Command: "kitty"},
			options:  Options{KeepOpen: true},
			wantCheck: func(args []string) bool {
				return contains(args, "--hold")
			},
			wantErr: false,
		},
		{
			name:      "unsupported terminal",
			terminal:  Terminal{Command: "unknown"},
			options:   Options{},
			wantCheck: func(args []string) bool { return false },
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := builder.BuildCommand(tt.terminal, cmd, &tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if !tt.wantCheck(got) {
				t.Errorf("BuildCommand() result does not match expected pattern: %v", got)
			}
		})
	}
}

func TestLinuxCommandBuilder_BuildCommand(t *testing.T) {
	builder := NewLinuxCommandBuilder()
	cmd := exec.Command("echo", "test")

	tests := []struct {
		name      string
		terminal  Terminal
		options   Options
		wantParts []string
		wantErr   bool
	}{
		{
			name:      "gnome-terminal with keep open",
			terminal:  Terminal{Command: "gnome-terminal"},
			options:   Options{KeepOpen: true, Title: "Test"},
			wantParts: []string{"--title", "Test", "--", "bash", "-c"},
			wantErr:   false,
		},
		{
			name:      "konsole with keep open",
			terminal:  Terminal{Command: "konsole"},
			options:   Options{KeepOpen: true},
			wantParts: []string{"--hold", "-e"},
			wantErr:   false,
		},
		{
			name:      "xterm with title",
			terminal:  Terminal{Command: "xterm"},
			options:   Options{Title: "Test", KeepOpen: false},
			wantParts: []string{"-title", "Test", "-e"},
			wantErr:   false,
		},
		{
			name:      "alacritty with keep open",
			terminal:  Terminal{Command: "alacritty"},
			options:   Options{KeepOpen: true},
			wantParts: []string{"-e", "sh", "-c", "read"},
			wantErr:   false,
		},
		{
			name:     "unsupported terminal",
			terminal: Terminal{Command: "unknown"},
			options:  Options{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := builder.BuildCommand(tt.terminal, cmd, &tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			// Check that expected parts are in the result
			gotStr := strings.Join(got, " ")
			for _, part := range tt.wantParts {
				if !strings.Contains(gotStr, part) {
					t.Errorf("BuildCommand() result %q missing expected part %q", gotStr, part)
				}
			}
		})
	}
}

func TestLinuxCommandBuilder_SupportsTerminal(t *testing.T) {
	builder := NewLinuxCommandBuilder()

	supportedTerminals := []string{
		"gnome-terminal", "konsole", "xfce4-terminal", "terminator",
		"alacritty", "kitty", "urxvt", "xterm",
	}

	for _, term := range supportedTerminals {
		t.Run(term, func(t *testing.T) {
			if !builder.SupportsTerminal(Terminal{Command: term}) {
				t.Errorf("SupportsTerminal(%q) = false, want true", term)
			}
		})
	}

	// Test unsupported
	if builder.SupportsTerminal(Terminal{Command: "unknown"}) {
		t.Error("SupportsTerminal(\"unknown\") = true, want false")
	}
}

// Helper function
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
