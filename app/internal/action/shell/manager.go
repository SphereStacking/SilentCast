package shell

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// Manager manages shell detection and execution
type Manager interface {
	// GetShell finds and validates a shell by name or path
	GetShell(ctx context.Context, nameOrPath string) (*Shell, error)

	// GetDefaultShell returns the system default shell
	GetDefaultShell(ctx context.Context) (*Shell, error)

	// ListShells returns all available shells
	ListShells(ctx context.Context) ([]Shell, error)

	// CreateCommand creates an exec.Cmd for running a script in a shell
	CreateCommand(ctx context.Context, shell *Shell, script string, env []string) (*exec.Cmd, error)

	// CreateInterpreterCommand creates an exec.Cmd for running code directly in an interpreter
	CreateInterpreterCommand(ctx context.Context, interpreter *Shell, code string, args []string, env []string) (*exec.Cmd, error)

	// GetShellForScript determines the best shell for a script
	GetShellForScript(ctx context.Context, script, preferredShell string) (*Shell, error)

	// GetInterpreterForScript determines the best interpreter for direct execution
	GetInterpreterForScript(ctx context.Context, script string, preferredInterpreter string) (*Shell, error)
}

// manager implements the Manager interface
type manager struct {
	detector Detector
	mu       sync.RWMutex
	cache    map[string]*Shell
}

// NewManager creates a new shell manager
func NewManager() Manager {
	return &manager{
		detector: NewDetector(),
		cache:    make(map[string]*Shell),
	}
}

// GetShell finds and validates a shell by name or path
func (m *manager) GetShell(ctx context.Context, nameOrPath string) (*Shell, error) {
	if nameOrPath == "" {
		return m.GetDefaultShell(ctx)
	}

	// Check cache
	m.mu.RLock()
	if cached, ok := m.cache[nameOrPath]; ok {
		m.mu.RUnlock()
		return cached, nil
	}
	m.mu.RUnlock()

	// Find shell
	shell, err := m.detector.FindShell(ctx, nameOrPath)
	if err != nil {
		return nil, err
	}

	// Cache result
	m.mu.Lock()
	m.cache[nameOrPath] = shell
	m.mu.Unlock()

	return shell, nil
}

// GetDefaultShell returns the system default shell
func (m *manager) GetDefaultShell(ctx context.Context) (*Shell, error) {
	return m.detector.GetDefaultShell(ctx)
}

// ListShells returns all available shells
func (m *manager) ListShells(ctx context.Context) ([]Shell, error) {
	return m.detector.DetectShells(ctx)
}

// CreateCommand creates an exec.Cmd for running a script in a shell
func (m *manager) CreateCommand(ctx context.Context, shell *Shell, script string, env []string) (*exec.Cmd, error) {
	if shell == nil {
		var err error
		shell, err = m.GetDefaultShell(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get default shell: %w", err)
		}
	}

	// Get command and arguments
	shellCmd, args := GetShellCommand(shell, script)

	// Create command
	cmd := exec.CommandContext(ctx, shellCmd, args...)

	// Set environment
	if len(env) > 0 {
		cmd.Env = append(os.Environ(), env...)
	} else {
		cmd.Env = os.Environ()
	}

	return cmd, nil
}

// CreateInterpreterCommand creates an exec.Cmd for running code directly in an interpreter
func (m *manager) CreateInterpreterCommand(ctx context.Context, interpreter *Shell, code string, args, env []string) (*exec.Cmd, error) {
	if interpreter == nil {
		return nil, fmt.Errorf("interpreter is required")
	}

	if !IsInterpreter(interpreter) {
		return nil, fmt.Errorf("shell %s is not an interpreter", interpreter.Name)
	}

	// Build command arguments for direct execution
	cmdArgs := make([]string, 0)

	// Add interpreter arguments (e.g., -c for Python, -e for Node)
	cmdArgs = append(cmdArgs, interpreter.Args...)

	// Add the code
	cmdArgs = append(cmdArgs, code)

	// Add any additional arguments
	cmdArgs = append(cmdArgs, args...)

	// Create command
	// nosec G204: interpreter.Executable is from trusted shell detection
	cmd := exec.CommandContext(ctx, interpreter.Executable, cmdArgs...)

	// Set environment
	if len(env) > 0 {
		cmd.Env = append(os.Environ(), env...)
	} else {
		cmd.Env = os.Environ()
	}

	return cmd, nil
}

// GetInterpreterForScript determines the best interpreter for direct execution
func (m *manager) GetInterpreterForScript(ctx context.Context, script, preferredInterpreter string) (*Shell, error) {
	// If preferred interpreter is specified, try to use it
	if preferredInterpreter != "" {
		shell, err := m.GetShell(ctx, preferredInterpreter)
		if err == nil && IsInterpreter(shell) {
			return shell, nil
		}
		// Log warning but continue
		fmt.Fprintf(os.Stderr, "Warning: preferred interpreter %q not found or not an interpreter\n", preferredInterpreter)
	}

	// Try to detect interpreter from shebang
	if shell := m.detectShebang(ctx, script); shell != nil && IsInterpreter(shell) {
		return shell, nil
	}

	// Try to detect from script content
	if shell := m.detectInterpreterFromContent(ctx, script); shell != nil {
		return shell, nil
	}

	return nil, fmt.Errorf("no suitable interpreter found for script")
}

// GetShellForScript determines the best shell for a script
func (m *manager) GetShellForScript(ctx context.Context, script, preferredShell string) (*Shell, error) {
	// If preferred shell is specified, try to use it
	if preferredShell != "" {
		shell, err := m.GetShell(ctx, preferredShell)
		if err == nil {
			return shell, nil
		}
		// Log warning but continue
		fmt.Fprintf(os.Stderr, "Warning: preferred shell %q not found, using default\n", preferredShell)
	}

	// Try to detect shell from shebang
	if shell := m.detectShebang(ctx, script); shell != nil {
		return shell, nil
	}

	// Use default shell
	return m.GetDefaultShell(ctx)
}

// detectShebang tries to detect shell from script shebang
func (m *manager) detectShebang(ctx context.Context, script string) *Shell {
	script = strings.TrimSpace(script)
	if !strings.HasPrefix(script, "#!") {
		return nil
	}

	// Extract shebang line
	lines := strings.SplitN(script, "\n", 2)
	if len(lines) == 0 {
		return nil
	}

	shebang := strings.TrimPrefix(lines[0], "#!")
	shebang = strings.TrimSpace(shebang)

	// Split shebang into command and args
	parts := strings.Fields(shebang)
	if len(parts) == 0 {
		return nil
	}

	// Handle /usr/bin/env
	if len(parts) >= 2 && strings.HasSuffix(parts[0], "/env") {
		// #! /usr/bin/env python3
		shell, err := m.detector.FindShell(ctx, parts[1])
		if err == nil {
			return shell
		}
	}

	// Direct path
	// #! /bin/bash
	shell, err := m.detector.ValidateShell(parts[0])
	if err == nil {
		return shell
	}

	return nil
}

// detectInterpreterFromContent tries to detect interpreter from script content
func (m *manager) detectInterpreterFromContent(ctx context.Context, script string) *Shell {
	script = strings.TrimSpace(script)
	if script == "" {
		return nil
	}

	// Look for interpreter-specific patterns
	patterns := []struct {
		interpreterName string
		patterns        []string
	}{
		{
			"ruby",
			[]string{
				"puts ",
				"end",
				".each",
				"||",
				"require ",
			},
		},
		{
			"python",
			[]string{
				"import ",
				"from ",
				"def ",
				"class ",
				"print(",
				"if __name__",
			},
		},
		{
			"node",
			[]string{
				"console.log",
				"require(",
				"module.exports",
				"const ",
				"let ",
				"var ",
				"process.",
				"Buffer.",
			},
		},
		{
			"perl",
			[]string{
				"use ",
				"my ",
				"our ",
				"sub ",
				"$_",
				"@",
				"%",
			},
		},
	}

	// Check each pattern
	for _, patternGroup := range patterns {
		matches := 0
		for _, pattern := range patternGroup.patterns {
			if strings.Contains(script, pattern) {
				matches++
			}
		}

		// If we find enough matches, try to get that interpreter
		if matches >= 2 {
			if shell, err := m.detector.FindShell(ctx, patternGroup.interpreterName); err == nil {
				return shell
			}
		}
	}

	return nil
}

// FormatShellInfo formats shell information for display
func FormatShellInfo(shell *Shell) string {
	if shell == nil {
		return "No shell"
	}

	info := shell.Name

	// Add type
	typeStr := ""
	switch shell.Type {
	case ShellTypeBourne:
		typeStr = "Bourne-compatible"
	case ShellTypeCsh:
		typeStr = "C shell"
	case ShellTypeFish:
		typeStr = "Fish"
	case ShellTypePowerShell:
		typeStr = "PowerShell"
	case ShellTypeCmd:
		typeStr = "Windows CMD"
	case ShellTypeInterpreter:
		typeStr = "Interpreter"
	}

	if typeStr != "" {
		info += fmt.Sprintf(" (%s)", typeStr)
	}

	// Add version if available
	if shell.Version != "" {
		// Truncate long version strings
		version := shell.Version
		if len(version) > 50 {
			version = version[:47] + "..."
		}
		info += fmt.Sprintf(" - %s", version)
	}

	// Add default indicator
	if shell.IsDefault {
		info += " [default]"
	}

	return info
}

// GetShellEnvironment returns environment variables for shell execution
func GetShellEnvironment(shell *Shell, additionalEnv map[string]string) []string {
	env := os.Environ()

	// Add shell-specific environment variables
	if shell != nil {
		switch shell.Type {
		case ShellTypePowerShell:
			// PowerShell specific
			env = append(env, "PSModulePath=")
		case ShellTypeInterpreter:
			// Python/Node/Ruby might need specific paths
			switch shell.Name {
			case "python", "python2", "python3":
				// Ensure Python doesn't buffer output
				env = append(env, "PYTHONUNBUFFERED=1")
			case "node", "nodejs":
				// Node.js specific
				env = append(env, "NODE_NO_WARNINGS=1")
			}
		}
	}

	// Add additional environment variables
	for k, v := range additionalEnv {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	return env
}
