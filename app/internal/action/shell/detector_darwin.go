//go:build darwin

package shell

import (
	"context"
	"os"
	"path/filepath"
	"strings"
)

// detectPlatformShells detects macOS-specific shells
func (d *detector) detectPlatformShells(ctx context.Context) []Shell {
	var shells []Shell

	// Check for Homebrew shells
	shells = append(shells, d.detectHomebrewShells(ctx)...)

	// Check for MacPorts shells
	shells = append(shells, d.detectMacPortsShells(ctx)...)

	// Check system shells
	shells = append(shells, d.detectSystemShells(ctx)...)

	return shells
}

// getPlatformDefaultShell returns the macOS default shell
func (d *detector) getPlatformDefaultShell(ctx context.Context) *Shell {
	// Check SHELL environment variable first
	if envShell := os.Getenv("SHELL"); envShell != "" {
		if shell, err := d.ValidateShell(envShell); err == nil {
			shell.IsDefault = true
			return shell
		}
	}

	// Modern macOS uses zsh by default
	if shell, err := d.ValidateShell("/bin/zsh"); err == nil {
		shell.IsDefault = true
		return shell
	}

	// Fallback to bash
	if shell, err := d.ValidateShell("/bin/bash"); err == nil {
		shell.IsDefault = true
		return shell
	}

	// Last resort: sh
	return &Shell{
		Name:       "sh",
		Executable: "/bin/sh",
		Type:       ShellTypeBourne,
		Args:       []string{"-c"},
		IsDefault:  true,
	}
}

// detectHomebrewShells detects shells installed via Homebrew
func (d *detector) detectHomebrewShells(ctx context.Context) []Shell {
	var shells []Shell

	// Check both Intel and Apple Silicon paths
	brewPaths := []string{
		"/usr/local/bin",    // Intel Macs
		"/opt/homebrew/bin", // Apple Silicon
		filepath.Join(os.Getenv("HOME"), ".brew/bin"), // User installation
	}

	shellNames := []string{"bash", "zsh", "fish", "tcsh", "ksh"}

	for _, brewPath := range brewPaths {
		for _, shellName := range shellNames {
			shellPath := filepath.Join(brewPath, shellName)
			if shell, err := d.ValidateShell(shellPath); err == nil {
				// Mark as Homebrew version
				shell.Name = "brew-" + shell.Name
				shells = append(shells, *shell)
			}
		}
	}

	return shells
}

// detectMacPortsShells detects shells installed via MacPorts
func (d *detector) detectMacPortsShells(ctx context.Context) []Shell {
	var shells []Shell

	// MacPorts default location
	portPath := "/opt/local/bin"

	shellNames := []string{"bash", "zsh", "fish", "tcsh", "ksh"}

	for _, shellName := range shellNames {
		shellPath := filepath.Join(portPath, shellName)
		if shell, err := d.ValidateShell(shellPath); err == nil {
			// Mark as MacPorts version
			shell.Name = "port-" + shell.Name
			shells = append(shells, *shell)
		}
	}

	return shells
}

// detectSystemShells detects system-provided shells
func (d *detector) detectSystemShells(ctx context.Context) []Shell {
	var shells []Shell

	// Check /etc/shells for available shells
	shellsFile := "/etc/shells"
	if data, err := os.ReadFile(shellsFile); err == nil {
		lines := string(data)
		for _, line := range strings.Split(lines, "\n") {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "#") {
				if shell, err := d.ValidateShell(line); err == nil {
					shells = append(shells, *shell)
				}
			}
		}
	}

	return shells
}
