package shell

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// Shell represents a detected shell
type Shell struct {
	Name       string    // e.g., "bash", "zsh", "powershell"
	Executable string    // Full path to shell executable
	Version    string    // Version string (if available)
	Type       ShellType // Type of shell
	Args       []string  // Default arguments for the shell
	IsDefault  bool      // Whether this is the system default shell
}

// ShellType represents the type of shell
type ShellType int

const (
	ShellTypeUnknown     ShellType = iota
	ShellTypeBourne                // sh, bash, zsh, etc.
	ShellTypeCsh                   // csh, tcsh
	ShellTypeFish                  // fish
	ShellTypePowerShell            // PowerShell
	ShellTypeCmd                   // Windows cmd
	ShellTypeInterpreter           // Python, Ruby, Node.js, etc.
)

// Detector interface for shell detection
type Detector interface {
	// DetectShells returns a list of available shells
	DetectShells(ctx context.Context) ([]Shell, error)

	// GetDefaultShell returns the system default shell
	GetDefaultShell(ctx context.Context) (*Shell, error)

	// FindShell finds a specific shell by name
	FindShell(ctx context.Context, name string) (*Shell, error)

	// ValidateShell checks if a shell path is valid and executable
	ValidateShell(path string) (*Shell, error)
}

// detector implements the Detector interface
type detector struct {
	mu    sync.RWMutex
	cache map[string]*Shell
}

// NewDetector creates a new shell detector
func NewDetector() Detector {
	return &detector{
		cache: make(map[string]*Shell),
	}
}

// CommonShells contains common shell definitions
var CommonShells = []struct {
	Name  string
	Type  ShellType
	Names []string // Alternative names/commands
}{
	{"bash", ShellTypeBourne, []string{"bash", "sh"}},
	{"zsh", ShellTypeBourne, []string{"zsh"}},
	{"fish", ShellTypeFish, []string{"fish"}},
	{"sh", ShellTypeBourne, []string{"sh"}},
	{"dash", ShellTypeBourne, []string{"dash"}},
	{"ksh", ShellTypeBourne, []string{"ksh", "ksh93"}},
	{"csh", ShellTypeCsh, []string{"csh"}},
	{"tcsh", ShellTypeCsh, []string{"tcsh"}},
	{"powershell", ShellTypePowerShell, []string{"powershell", "pwsh", "powershell.exe", "pwsh.exe"}},
	{"cmd", ShellTypeCmd, []string{"cmd", "cmd.exe"}},
	{"python", ShellTypeInterpreter, []string{"python", "python3", "python2"}},
	{"node", ShellTypeInterpreter, []string{"node", "nodejs"}},
	{"ruby", ShellTypeInterpreter, []string{"ruby"}},
	{"perl", ShellTypeInterpreter, []string{"perl"}},
}

// DetectShells returns a list of available shells
func (d *detector) DetectShells(ctx context.Context) ([]Shell, error) {
	var shells []Shell

	// Detect common shells
	for _, shellDef := range CommonShells {
		for _, cmd := range shellDef.Names {
			if shell := d.detectShellCommand(ctx, cmd, shellDef.Name, shellDef.Type); shell != nil {
				shells = append(shells, *shell)
				break // Found this shell, no need to check other names
			}
		}
	}

	// Add platform-specific shells
	shells = append(shells, d.detectPlatformShells(ctx)...)

	// Remove duplicates
	shells = d.deduplicateShells(shells)

	// Mark default shell
	if defaultShell, err := d.GetDefaultShell(ctx); err == nil && defaultShell != nil {
		for i := range shells {
			if shells[i].Executable == defaultShell.Executable {
				shells[i].IsDefault = true
			}
		}
	}

	return shells, nil
}

// GetDefaultShell returns the system default shell
func (d *detector) GetDefaultShell(ctx context.Context) (*Shell, error) {
	// Check cache first
	d.mu.RLock()
	if cached, ok := d.cache["default"]; ok {
		d.mu.RUnlock()
		return cached, nil
	}
	d.mu.RUnlock()

	var shell *Shell

	// Try environment variable first
	if envShell := os.Getenv("SHELL"); envShell != "" {
		if s, err := d.ValidateShell(envShell); err == nil {
			shell = s
			shell.IsDefault = true
		}
	}

	// Platform-specific fallback
	if shell == nil {
		shell = d.getPlatformDefaultShell(ctx)
	}

	if shell == nil {
		return nil, fmt.Errorf("no default shell found")
	}

	// Cache the result
	d.mu.Lock()
	d.cache["default"] = shell
	d.mu.Unlock()

	return shell, nil
}

// FindShell finds a specific shell by name
func (d *detector) FindShell(ctx context.Context, name string) (*Shell, error) {
	// Check cache first
	cacheKey := "find:" + name
	d.mu.RLock()
	if cached, ok := d.cache[cacheKey]; ok {
		d.mu.RUnlock()
		return cached, nil
	}
	d.mu.RUnlock()

	// Normalize name
	name = strings.ToLower(strings.TrimSpace(name))

	// Check if it's a full path
	if filepath.IsAbs(name) || strings.Contains(name, string(os.PathSeparator)) {
		return d.ValidateShell(name)
	}

	// Search for shell
	for _, shellDef := range CommonShells {
		if shellDef.Name == name {
			for _, cmd := range shellDef.Names {
				if shell := d.detectShellCommand(ctx, cmd, shellDef.Name, shellDef.Type); shell != nil {
					// Cache the result
					d.mu.Lock()
					d.cache[cacheKey] = shell
					d.mu.Unlock()
					return shell, nil
				}
			}
		}
	}

	// Try as-is
	if shell := d.detectShellCommand(ctx, name, name, ShellTypeUnknown); shell != nil {
		// Cache the result
		d.mu.Lock()
		d.cache[cacheKey] = shell
		d.mu.Unlock()
		return shell, nil
	}

	return nil, fmt.Errorf("shell not found: %s", name)
}

// ValidateShell checks if a shell path is valid and executable
func (d *detector) ValidateShell(path string) (*Shell, error) {
	// Check if file exists and is executable
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("shell not found: %w", err)
	}

	if info.IsDir() {
		return nil, fmt.Errorf("shell path is a directory: %s", path)
	}

	// Check if executable (on Unix)
	if runtime.GOOS != "windows" {
		if info.Mode()&0o111 == 0 {
			return nil, fmt.Errorf("shell is not executable: %s", path)
		}
	}

	// Determine shell type and name
	name := filepath.Base(path)
	shellType := d.getShellType(name)

	// Try to get version
	version := d.getShellVersion(path)

	return &Shell{
		Name:       name,
		Executable: path,
		Version:    version,
		Type:       shellType,
		Args:       d.getShellArgs(name, shellType),
	}, nil
}

// detectShellCommand detects a shell by command name
func (d *detector) detectShellCommand(_ context.Context, cmd, name string, shellType ShellType) *Shell {
	// Use 'which' or 'where' to find the shell
	path, err := exec.LookPath(cmd)
	if err != nil {
		return nil
	}

	// Validate the shell
	shell, err := d.ValidateShell(path)
	if err != nil {
		return nil
	}

	// Override name and type if provided
	if name != "" {
		shell.Name = name
	}
	if shellType != ShellTypeUnknown {
		shell.Type = shellType
	}

	return shell
}

// getShellType determines the shell type from its name
func (d *detector) getShellType(name string) ShellType {
	name = strings.ToLower(filepath.Base(name))

	// Remove .exe suffix on Windows
	name = strings.TrimSuffix(name, ".exe")

	for _, shellDef := range CommonShells {
		for _, shellName := range shellDef.Names {
			if name == shellName {
				return shellDef.Type
			}
		}
	}

	// Check by suffix patterns
	switch {
	case strings.HasSuffix(name, "sh"):
		return ShellTypeBourne
	case strings.Contains(name, "python"):
		return ShellTypeInterpreter
	case strings.Contains(name, "node"):
		return ShellTypeInterpreter
	case strings.Contains(name, "ruby"):
		return ShellTypeInterpreter
	default:
		return ShellTypeUnknown
	}
}

// getShellArgs returns default arguments for a shell
func (d *detector) getShellArgs(name string, shellType ShellType) []string {
	name = strings.ToLower(filepath.Base(name))
	name = strings.TrimSuffix(name, ".exe")

	switch shellType {
	case ShellTypeBourne, ShellTypeCsh, ShellTypeFish:
		return []string{"-c"}
	case ShellTypePowerShell:
		if name == "powershell" {
			return []string{"-Command"}
		}
		// pwsh (PowerShell Core)
		return []string{"-c"}
	case ShellTypeCmd:
		return []string{"/c"}
	case ShellTypeInterpreter:
		switch name {
		case "python", "python2", "python3":
			return []string{"-c"}
		case "node", "nodejs":
			return []string{"-e"}
		case "ruby":
			return []string{"-e"}
		case "perl":
			return []string{"-e"}
		default:
			return []string{}
		}
	default:
		return []string{"-c"} // Default to -c
	}
}

// getShellVersion attempts to get the shell version
func (d *detector) getShellVersion(path string) string {
	// Try common version flags
	versionFlags := []string{"--version", "-version", "/version", "-v"}

	for _, flag := range versionFlags {
		cmd := exec.Command(path, flag)
		output, err := cmd.Output()
		if err == nil && len(output) > 0 {
			// Extract first line
			lines := strings.Split(string(output), "\n")
			if len(lines) > 0 {
				return strings.TrimSpace(lines[0])
			}
		}
	}

	return ""
}

// deduplicateShells removes duplicate shells
func (d *detector) deduplicateShells(shells []Shell) []Shell {
	seen := make(map[string]bool)
	var result []Shell

	for _, shell := range shells {
		if !seen[shell.Executable] {
			seen[shell.Executable] = true
			result = append(result, shell)
		}
	}

	return result
}

// GetShellCommand returns the command and arguments to execute a script
func GetShellCommand(shell *Shell, script string) (string, []string) {
	if shell == nil {
		// Use default shell
		detector := NewDetector()
		defaultShell, err := detector.GetDefaultShell(context.Background())
		if err != nil {
			// Fallback to basic shell
			if runtime.GOOS == "windows" {
				return "cmd", []string{"/c", script}
			}
			return "sh", []string{"-c", script}
		}
		shell = defaultShell
	}

	args := make([]string, len(shell.Args))
	copy(args, shell.Args)
	args = append(args, script)

	return shell.Executable, args
}

// IsInterpreter returns true if the shell is an interpreter (Python, Node, etc.)
func IsInterpreter(shell *Shell) bool {
	return shell != nil && shell.Type == ShellTypeInterpreter
}
