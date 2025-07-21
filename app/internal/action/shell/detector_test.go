package shell

import (
	"context"
	"runtime"
	"strings"
	"testing"
)

func TestNewDetector(t *testing.T) {
	detector := NewDetector()
	if detector == nil {
		t.Fatal("NewDetector() returned nil")
	}
}

func TestDetectShells(t *testing.T) {
	detector := NewDetector()
	ctx := context.Background()

	shells, err := detector.DetectShells(ctx)
	if err != nil {
		t.Fatalf("DetectShells() error: %v", err)
	}

	// Should find at least one shell on any system
	if len(shells) == 0 {
		t.Error("DetectShells() returned no shells")
	}

	t.Logf("Found %d shells:", len(shells))
	for _, shell := range shells {
		t.Logf("  %s", FormatShellInfo(&shell))
		t.Logf("    Path: %s", shell.Executable)
		t.Logf("    Args: %v", shell.Args)
	}

	// Check for duplicates
	seen := make(map[string]bool)
	for _, shell := range shells {
		if seen[shell.Executable] {
			t.Errorf("Duplicate shell found: %s", shell.Executable)
		}
		seen[shell.Executable] = true
	}
}

func TestGetDefaultShell(t *testing.T) {
	detector := NewDetector()
	ctx := context.Background()

	shell, err := detector.GetDefaultShell(ctx)
	if err != nil {
		t.Fatalf("GetDefaultShell() error: %v", err)
	}

	if shell == nil {
		t.Fatal("GetDefaultShell() returned nil")
	}

	if !shell.IsDefault {
		t.Error("Default shell should have IsDefault=true")
	}

	t.Logf("Default shell: %s", FormatShellInfo(shell))
	t.Logf("  Path: %s", shell.Executable)
	t.Logf("  Args: %v", shell.Args)

	// Verify shell is executable
	if _, err := detector.ValidateShell(shell.Executable); err != nil {
		t.Errorf("Default shell validation failed: %v", err)
	}
}

func TestFindShell(t *testing.T) {
	detector := NewDetector()
	ctx := context.Background()

	tests := []struct {
		name      string
		shellName string
		wantErr   bool
	}{
		{"Find sh", "sh", false},
		{"Find nonexistent", "nonexistentshell123", true},
	}

	// Add platform-specific tests
	switch runtime.GOOS {
	case "windows":
		tests = append(tests, struct { //nolint:gocritic // Platform-specific test cases need separate append
			name      string
			shellName string
			wantErr   bool
		}{"Find cmd", "cmd", false})
		tests = append(tests, struct {
			name      string
			shellName string
			wantErr   bool
		}{"Find powershell", "powershell", false})
	case "darwin", "linux":
		tests = append(tests, struct {
			name      string
			shellName string
			wantErr   bool
		}{"Find bash", "bash", false})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shell, err := detector.FindShell(ctx, tt.shellName)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindShell(%q) error = %v, wantErr %v", tt.shellName, err, tt.wantErr)
				return
			}

			if !tt.wantErr && shell != nil {
				t.Logf("Found %s: %s at %s", tt.shellName, shell.Name, shell.Executable)
			}
		})
	}
}

func TestValidateShell(t *testing.T) {
	detector := NewDetector()

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"Invalid path", "/nonexistent/shell", true},
		{"Directory", "/tmp", true},
		{"Empty path", "", true},
	}

	// Add valid shell test based on platform
	defaultShell, err := detector.GetDefaultShell(context.Background())
	if err == nil && defaultShell != nil {
		tests = append(tests, struct {
			name    string
			path    string
			wantErr bool
		}{"Valid shell", defaultShell.Executable, false})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shell, err := detector.ValidateShell(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateShell(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
				return
			}

			if !tt.wantErr && shell != nil {
				if shell.Executable != tt.path {
					t.Errorf("ValidateShell() returned wrong path: got %q, want %q", shell.Executable, tt.path)
				}
			}
		})
	}
}

func TestGetShellType(t *testing.T) {
	detector := &detector{}

	tests := []struct {
		name     string
		expected ShellType
	}{
		{"bash", ShellTypeBourne},
		{"zsh", ShellTypeBourne},
		{"sh", ShellTypeBourne},
		{"dash", ShellTypeBourne},
		{"ksh", ShellTypeBourne},
		{"csh", ShellTypeCsh},
		{"tcsh", ShellTypeCsh},
		{"fish", ShellTypeFish},
		{"powershell", ShellTypePowerShell},
		{"pwsh", ShellTypePowerShell},
		{"cmd", ShellTypeCmd},
		{"cmd.exe", ShellTypeCmd},
		{"python", ShellTypeInterpreter},
		{"python3", ShellTypeInterpreter},
		{"node", ShellTypeInterpreter},
		{"ruby", ShellTypeInterpreter},
		{"unknown", ShellTypeUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detector.getShellType(tt.name)
			if result != tt.expected {
				t.Errorf("getShellType(%q) = %v, want %v", tt.name, result, tt.expected)
			}
		})
	}
}

func TestGetShellArgs(t *testing.T) {
	detector := &detector{}

	tests := []struct {
		name      string
		shellType ShellType
		expected  []string
	}{
		{"bash", ShellTypeBourne, []string{"-c"}},
		{"zsh", ShellTypeBourne, []string{"-c"}},
		{"fish", ShellTypeFish, []string{"-c"}},
		{"csh", ShellTypeCsh, []string{"-c"}},
		{"powershell", ShellTypePowerShell, []string{"-Command"}},
		{"pwsh", ShellTypePowerShell, []string{"-c"}},
		{"cmd", ShellTypeCmd, []string{"/c"}},
		{"python", ShellTypeInterpreter, []string{"-c"}},
		{"node", ShellTypeInterpreter, []string{"-e"}},
		{"ruby", ShellTypeInterpreter, []string{"-e"}},
		{"perl", ShellTypeInterpreter, []string{"-e"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detector.getShellArgs(tt.name, tt.shellType)
			if len(result) != len(tt.expected) {
				t.Errorf("getShellArgs(%q, %v) returned %d args, want %d", tt.name, tt.shellType, len(result), len(tt.expected))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("getShellArgs(%q, %v)[%d] = %q, want %q", tt.name, tt.shellType, i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestGetShellCommand(t *testing.T) {
	tests := []struct {
		name   string
		shell  *Shell
		script string
	}{
		{
			name: "Bash shell",
			shell: &Shell{
				Name:       "bash",
				Executable: "/bin/bash",
				Type:       ShellTypeBourne,
				Args:       []string{"-c"},
			},
			script: "echo hello",
		},
		{
			name: "PowerShell",
			shell: &Shell{
				Name:       "powershell",
				Executable: "powershell.exe",
				Type:       ShellTypePowerShell,
				Args:       []string{"-Command"},
			},
			script: "Write-Host 'hello'",
		},
		{
			name:   "Nil shell (use default)",
			shell:  nil,
			script: "echo hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, args := GetShellCommand(tt.shell, tt.script)

			if tt.shell != nil {
				if cmd != tt.shell.Executable {
					t.Errorf("GetShellCommand() cmd = %q, want %q", cmd, tt.shell.Executable)
				}

				expectedArgs := append(tt.shell.Args, tt.script) //nolint:gocritic // Test needs to compare against expected args
				if len(args) != len(expectedArgs) {
					t.Errorf("GetShellCommand() returned %d args, want %d", len(args), len(expectedArgs))
				}
			} else {
				// Should use default shell
				if cmd == "" {
					t.Error("GetShellCommand() with nil shell returned empty command")
				}
				if len(args) < 2 {
					t.Error("GetShellCommand() with nil shell returned too few args")
				}
			}
		})
	}
}

func TestIsInterpreter(t *testing.T) {
	tests := []struct {
		name     string
		shell    *Shell
		expected bool
	}{
		{
			name:     "Nil shell",
			shell:    nil,
			expected: false,
		},
		{
			name: "Python interpreter",
			shell: &Shell{
				Name: "python",
				Type: ShellTypeInterpreter,
			},
			expected: true,
		},
		{
			name: "Bash shell",
			shell: &Shell{
				Name: "bash",
				Type: ShellTypeBourne,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsInterpreter(tt.shell)
			if result != tt.expected {
				t.Errorf("IsInterpreter() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCaching(t *testing.T) {
	detector := NewDetector()
	ctx := context.Background()

	// First call should detect
	shell1, err1 := detector.GetDefaultShell(ctx)
	if err1 != nil {
		t.Fatalf("First GetDefaultShell() error: %v", err1)
	}

	// Second call should use cache
	shell2, err2 := detector.GetDefaultShell(ctx)
	if err2 != nil {
		t.Fatalf("Second GetDefaultShell() error: %v", err2)
	}

	// Should return same shell
	if shell1.Executable != shell2.Executable {
		t.Error("GetDefaultShell() not using cache")
	}
}

func TestShellVersionDetection(t *testing.T) {
	detector := &detector{}

	// This test is informational - version detection may fail
	shells := []string{"/bin/bash", "/bin/sh", "python", "node"}

	for _, shellPath := range shells {
		version := detector.getShellVersion(shellPath)
		if version != "" {
			t.Logf("%s version: %s", shellPath, version)

			// Basic validation
			if strings.Contains(version, "\n") {
				t.Errorf("Version for %s contains newline", shellPath)
			}
		}
	}
}
