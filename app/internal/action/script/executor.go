package script

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/SphereStacking/silentcast/internal/action/shell"
	"github.com/SphereStacking/silentcast/internal/config"
	"github.com/SphereStacking/silentcast/internal/errors"
	"github.com/SphereStacking/silentcast/internal/notify"
	"github.com/SphereStacking/silentcast/internal/output"
)

// ScriptExecutor executes script/command actions
type ScriptExecutor struct {
	config   config.ActionConfig
	notifier *notify.Manager
}

// NewScriptExecutor creates a new script executor
func NewScriptExecutor(cfg *config.ActionConfig) *ScriptExecutor {
	return &ScriptExecutor{
		config:   *cfg,
		notifier: notify.NewManager(),
	}
}

// Execute runs the script or command
func (e *ScriptExecutor) Execute(ctx context.Context) error {
	// Apply timeout if configured
	if e.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(e.config.Timeout)*time.Second)
		defer cancel()
	}
	// Expand environment variables in command
	command := os.ExpandEnv(e.config.Command)

	// Check for empty command
	if strings.TrimSpace(command) == "" {
		return errors.New(errors.ErrorTypeConfig, "empty command").
			WithContext("command", e.config.Command).
			WithContext("action_type", "script").
			WithContext("error_type", "empty_command").
			WithContext("suggested_action", "provide a valid command in spellbook.yml")
	}

	// Get platform-specific shell executor
	shellExec := shell.GetShellExecutor()
	shell, shellFlag := shellExec.GetShell()
	
	// Override shell if custom one is specified
	if e.config.Shell != "" {
		shell = e.config.Shell
		// Determine the appropriate flag for the custom shell
		switch strings.ToLower(e.config.Shell) {
		case "powershell", "pwsh":
			shellFlag = "-Command"
		case "cmd", "cmd.exe":
			shellFlag = "/c"
		default:
			// Most Unix-like shells use -c
			shellFlag = "-c"
		}
	}

	// Create command
	var cmd *exec.Cmd
	if len(e.config.Args) > 0 {
		// If args are provided, don't use shell, execute directly
		parts := strings.Fields(command)
		if len(parts) > 0 {
			cmd = exec.CommandContext(ctx, parts[0], append(parts[1:], e.config.Args...)...) // #nosec G204 - Command is from trusted config file
		} else {
			return errors.New(errors.ErrorTypeConfig, "empty command").
				WithContext("command", command).
				WithContext("action_type", "script").
				WithContext("error_type", "empty_command").
				WithContext("suggested_action", "provide a valid command in spellbook.yml")
		}
	} else {
		// Use shell to execute the command
		cmd = exec.CommandContext(ctx, shell, shellFlag, command) // #nosec G204 - Command is from trusted config file
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

	// Setup output capture if ShowOutput is enabled
	var outputManager output.Manager
	if e.config.ShowOutput {
		outputManager = output.NewBufferedManager(output.DefaultOptions())
		writer := outputManager.StartCapture()
		
		// Redirect both stdout and stderr to our output manager
		cmd.Stdout = writer
		cmd.Stderr = writer
	}

	// Handle terminal execution based on config
	needsTerminal := e.config.Terminal || e.config.KeepOpen || shellExec.IsInteractiveCommand(command)
	if needsTerminal {
		// Use WrapInTerminalWithOptions when available
		cmd = shellExec.WrapInTerminalWithOptions(ctx, cmd, e.config.KeepOpen)
	}


	// Start the script
	if err := cmd.Start(); err != nil {
		return errors.Wrap(errors.ErrorTypeSystem, "failed to start script", err).
			WithContext("command", e.config.Command).
			WithContext("action_type", "script").
			WithContext("working_dir", cmd.Dir).
			WithContext("error_type", "start_failed").
			WithContext("suggested_action", "check if command exists and has execute permissions")
	}

	// For background scripts, detach from the process
	if !e.shouldWaitForCompletion(command) && !e.config.ShowOutput {
		if err := cmd.Process.Release(); err != nil {
			// Non-fatal error
			_ = err // Explicitly ignore
		}
		return nil
	}

	// Wait for completion if needed
	err := cmd.Wait()
	
	// Send notification with output if ShowOutput is enabled
	if e.config.ShowOutput && outputManager != nil {
		capturedOutput := outputManager.GetOutput()
		
		// Prepare notification
		title := e.config.Description
		if title == "" {
			title = fmt.Sprintf("Script: %s", e.config.Command)
		}
		
		// Determine notification level based on error
		if err != nil {
			_ = e.notifier.Error(ctx, title, fmt.Sprintf("Failed: %v\n\nOutput:\n%s", err, capturedOutput))
		} else if capturedOutput != "" {
			_ = e.notifier.Success(ctx, title, capturedOutput)
		} else {
			_ = e.notifier.Info(ctx, title, "Command completed with no output")
		}
		
		// Clean up output manager
		_ = outputManager.Stop()
	}
	
	if err != nil {
		return errors.Wrap(errors.ErrorTypeSystem, "script execution failed", err).
			WithContext("command", e.config.Command).
			WithContext("action_type", "script").
			WithContext("working_dir", cmd.Dir).
			WithContext("error_type", "execution_failed").
			WithContext("exit_code", cmd.ProcessState.ExitCode()).
			WithContext("suggested_action", "check command syntax and arguments")
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
