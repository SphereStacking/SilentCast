//go:build linux

package shell

import (
	"context"
	"os"
	"path/filepath"
	"strings"
)

// detectPlatformShells detects Linux-specific shells
func (d *detector) detectPlatformShells(ctx context.Context) []Shell {
	var shells []Shell

	// Check system shells from /etc/shells
	shells = append(shells, d.detectSystemShells(ctx)...)

	// Check snap packages
	shells = append(shells, d.detectSnapShells(ctx)...)

	// Check flatpak packages
	shells = append(shells, d.detectFlatpakShells(ctx)...)

	// Check user-installed shells
	shells = append(shells, d.detectUserShells(ctx)...)

	return shells
}

// getPlatformDefaultShell returns the Linux default shell
func (d *detector) getPlatformDefaultShell(ctx context.Context) *Shell {
	// Check SHELL environment variable first
	if envShell := os.Getenv("SHELL"); envShell != "" {
		if shell, err := d.ValidateShell(envShell); err == nil {
			shell.IsDefault = true
			return shell
		}
	}

	// Try to get from /etc/passwd
	if shell := d.getShellFromPasswd(); shell != nil {
		shell.IsDefault = true
		return shell
	}

	// Try common defaults
	defaults := []string{"/bin/bash", "/bin/sh", "/usr/bin/bash", "/usr/bin/sh"}
	for _, path := range defaults {
		if shell, err := d.ValidateShell(path); err == nil {
			shell.IsDefault = true
			return shell
		}
	}

	// Last resort
	return &Shell{
		Name:       "sh",
		Executable: "/bin/sh",
		Type:       ShellTypeBourne,
		Args:       []string{"-c"},
		IsDefault:  true,
	}
}

// detectSystemShells detects system-provided shells
func (d *detector) detectSystemShells(ctx context.Context) []Shell {
	var shells []Shell

	// Read /etc/shells
	if data, err := os.ReadFile("/etc/shells"); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
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

// detectSnapShells detects shells installed via Snap
func (d *detector) detectSnapShells(_ context.Context) []Shell {
	var shells []Shell

	snapBin := "/snap/bin"
	if entries, err := os.ReadDir(snapBin); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				name := entry.Name()
				// Check if it's a known shell
				for _, shellDef := range CommonShells {
					for _, shellName := range shellDef.Names {
						if name == shellName {
							shellPath := filepath.Join(snapBin, name)
							if shell, err := d.ValidateShell(shellPath); err == nil {
								shell.Name = "snap-" + shell.Name
								shells = append(shells, *shell)
							}
							break
						}
					}
				}
			}
		}
	}

	return shells
}

// detectFlatpakShells detects shells installed via Flatpak
func (d *detector) detectFlatpakShells(_ context.Context) []Shell {
	var shells []Shell

	// Flatpak typically doesn't install shells directly
	// But check common locations anyway
	flatpakDirs := []string{
		"/var/lib/flatpak/exports/bin",
		filepath.Join(os.Getenv("HOME"), ".local", "share", "flatpak", "exports", "bin"), //nolint:gocritic // Path components are separate constants
	}

	for _, dir := range flatpakDirs {
		if entries, err := os.ReadDir(dir); err == nil {
			for _, entry := range entries {
				if !entry.IsDir() {
					name := entry.Name()
					// Check if it's a known shell
					for _, shellDef := range CommonShells {
						for _, shellName := range shellDef.Names {
							if name == shellName {
								shellPath := filepath.Join(dir, name)
								if shell, err := d.ValidateShell(shellPath); err == nil {
									shell.Name = "flatpak-" + shell.Name
									shells = append(shells, *shell)
								}
								break
							}
						}
					}
				}
			}
		}
	}

	return shells
}

// detectUserShells detects user-installed shells
func (d *detector) detectUserShells(_ context.Context) []Shell {
	var shells []Shell

	// Check common user directories
	userDirs := []string{
		filepath.Join(os.Getenv("HOME"), ".local", "bin"), //nolint:gocritic // Path components are separate constants
		filepath.Join(os.Getenv("HOME"), "bin"),
		"/usr/local/bin",
		"/opt/bin",
	}

	shellNames := []string{"bash", "zsh", "fish", "tcsh", "ksh", "dash"}

	for _, dir := range userDirs {
		for _, shellName := range shellNames {
			shellPath := filepath.Join(dir, shellName)
			if shell, err := d.ValidateShell(shellPath); err == nil {
				// Check if it's different from system version
				systemPath := filepath.Join("/bin", shellName) //nolint:gocritic // "/bin" is a standard path constant
				if shellPath != systemPath {
					shell.Name = "user-" + shell.Name
				}
				shells = append(shells, *shell)
			}
		}
	}

	return shells
}

// getShellFromPasswd tries to get the user's shell from /etc/passwd
func (d *detector) getShellFromPasswd() *Shell {
	username := os.Getenv("USER")
	if username == "" {
		return nil
	}

	data, err := os.ReadFile("/etc/passwd")
	if err != nil {
		return nil
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) >= 7 && fields[0] == username {
			shellPath := fields[6]
			if shell, err := d.ValidateShell(shellPath); err == nil {
				return shell
			}
		}
	}

	return nil
}
