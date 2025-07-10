package action

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/SphereStacking/silentcast/internal/config"
)

// ScriptExecutor executes script/command actions
type ScriptExecutor struct {
	config config.ActionConfig
}

// NewScriptExecutor creates a new script executor
func NewScriptExecutor(cfg config.ActionConfig) *ScriptExecutor {
	return &ScriptExecutor{
		config: cfg,
	}
}

// Execute runs the script or command
func (e *ScriptExecutor) Execute(ctx context.Context) error {
	// Expand environment variables in command
	command := os.ExpandEnv(e.config.Command)
	
	// Check for empty command
	if strings.TrimSpace(command) == "" {
		return fmt.Errorf("empty command")
	}
	
	// Get platform-specific shell executor
	shellExec := GetShellExecutor()
	shell, shellFlag := shellExec.GetShell()
	
	// Create command
	var cmd *exec.Cmd
	if len(e.config.Args) > 0 {
		// If args are provided, don't use shell, execute directly
		parts := strings.Fields(command)
		if len(parts) > 0 {
			cmd = exec.CommandContext(ctx, parts[0], append(parts[1:], e.config.Args...)...)
		} else {
			return fmt.Errorf("empty command")
		}
	} else {
		// Use shell to execute the command
		cmd = exec.CommandContext(ctx, shell, shellFlag, command)
	}
	
	// Set working directory if specified
	if e.config.WorkingDir != "" {
		cmd.Dir = os.ExpandEnv(e.config.WorkingDir)
	} else {
		// Default to user's home directory
		if home, err := os.UserHomeDir(); err == nil {
			cmd.Dir = home
		}
	}
	
	// Set environment variables
	cmd.Env = os.Environ()
	if len(e.config.Env) > 0 {
		for k, v := range e.config.Env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, os.ExpandEnv(v)))
		}
	}
	
	// For scripts, we'll run them in a new terminal window
	if shellExec.IsInteractiveCommand(command) {
		cmd = shellExec.WrapInTerminal(ctx, cmd)
	}
	
	// Start the script
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start script: %w", err)
	}
	
	// For background scripts, detach from the process
	if !e.shouldWaitForCompletion(command) {
		if err := cmd.Process.Release(); err != nil {
			// Non-fatal error
		}
		return nil
	}
	
	// Wait for completion if needed
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("script execution failed: %w", err)
	}
	
	return nil
}

// String returns a string representation of the action
func (e *ScriptExecutor) String() string {
	if e.config.Description != "" {
		return e.config.Description
	}
	return fmt.Sprintf("Run script: %s", e.config.Command)
}

// shouldWaitForCompletion determines if we should wait for the script to complete
func (e *ScriptExecutor) shouldWaitForCompletion(command string) bool {
	// Quick commands that should complete
	quickCommands := []string{"git", "echo", "date", "pwd", "ls", "true", "false"}
	for _, cmd := range quickCommands {
		if strings.HasPrefix(command, cmd) {
			return true
		}
	}
	return false
}