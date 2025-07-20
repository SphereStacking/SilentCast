//go:build windows

package shell

import (
	"context"
	"os"
	"path/filepath"
)

// detectPlatformShells detects Windows-specific shells
func (d *detector) detectPlatformShells(ctx context.Context) []Shell {
	var shells []Shell

	// Check for Windows Terminal shells
	if shell := d.detectWindowsTerminal(ctx); shell != nil {
		shells = append(shells, *shell)
	}

	// Check for Git Bash
	if shell := d.detectGitBash(ctx); shell != nil {
		shells = append(shells, *shell)
	}

	// Check for WSL
	if shell := d.detectWSL(ctx); shell != nil {
		shells = append(shells, *shell)
	}

	// Check for Cygwin
	if shell := d.detectCygwin(ctx); shell != nil {
		shells = append(shells, *shell)
	}

	return shells
}

// getPlatformDefaultShell returns the Windows default shell
func (d *detector) getPlatformDefaultShell(ctx context.Context) *Shell {
	// On Windows, cmd.exe is the default
	comspec := os.Getenv("COMSPEC")
	if comspec == "" {
		comspec = "C:\\Windows\\System32\\cmd.exe"
	}

	shell, err := d.ValidateShell(comspec)
	if err != nil {
		// Fallback to basic cmd
		return &Shell{
			Name:       "cmd",
			Executable: "cmd.exe",
			Type:       ShellTypeCmd,
			Args:       []string{"/c"},
			IsDefault:  true,
		}
	}

	shell.IsDefault = true
	return shell
}

// detectWindowsTerminal detects Windows Terminal if available
func (d *detector) detectWindowsTerminal(ctx context.Context) *Shell {
	// Windows Terminal doesn't have its own shell, skip
	return nil
}

// detectGitBash detects Git Bash installation
func (d *detector) detectGitBash(ctx context.Context) *Shell {
	// Common Git Bash locations
	paths := []string{
		`C:\Program Files\Git\bin\bash.exe`,
		`C:\Program Files (x86)\Git\bin\bash.exe`,
		filepath.Join(os.Getenv("PROGRAMFILES"), `Git\bin\bash.exe`),
		filepath.Join(os.Getenv("PROGRAMFILES(X86)"), `Git\bin\bash.exe`),
		filepath.Join(os.Getenv("LOCALAPPDATA"), `Programs\Git\bin\bash.exe`),
	}

	for _, path := range paths {
		if shell, err := d.ValidateShell(path); err == nil {
			shell.Name = "git-bash"
			shell.Type = ShellTypeBourne
			shell.Args = []string{"-c"}
			return shell
		}
	}

	return nil
}

// detectWSL detects Windows Subsystem for Linux
func (d *detector) detectWSL(ctx context.Context) *Shell {
	// Check if wsl.exe exists
	wslPath := filepath.Join(os.Getenv("WINDIR"), `System32\wsl.exe`)
	if _, err := os.Stat(wslPath); err == nil {
		return &Shell{
			Name:       "wsl",
			Executable: wslPath,
			Type:       ShellTypeBourne,
			Args:       []string{"--", "sh", "-c"},
		}
	}

	return nil
}

// detectCygwin detects Cygwin installation
func (d *detector) detectCygwin(ctx context.Context) *Shell {
	// Common Cygwin locations
	paths := []string{
		`C:\cygwin64\bin\bash.exe`,
		`C:\cygwin\bin\bash.exe`,
		filepath.Join(os.Getenv("PROGRAMFILES"), `cygwin\bin\bash.exe`),
	}

	for _, path := range paths {
		if shell, err := d.ValidateShell(path); err == nil {
			shell.Name = "cygwin-bash"
			shell.Type = ShellTypeBourne
			shell.Args = []string{"-c"}
			return shell
		}
	}

	return nil
}
