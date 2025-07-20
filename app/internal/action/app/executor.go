package app

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/SphereStacking/silentcast/internal/config"
)

// AppExecutor executes application launch actions
type AppExecutor struct {
	config config.ActionConfig
}

// NewAppExecutor creates a new application executor
func NewAppExecutor(cfg *config.ActionConfig) *AppExecutor {
	return &AppExecutor{
		config: *cfg,
	}
}

// Execute launches the application
func (e *AppExecutor) Execute(ctx context.Context) error {
	// Expand environment variables in command
	path := os.ExpandEnv(e.config.Command)

	// Get platform-specific app launcher
	launcher := GetAppLauncher()

	// Check if the application exists (skip for certain special cases)
	if !launcher.IsSpecialPath(path) {
		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				// Try to find in PATH
				if _, lookupErr := exec.LookPath(path); lookupErr != nil {
					return fmt.Errorf("application not found: %s", path)
				}
			} else {
				return fmt.Errorf("failed to check application: %w", err)
			}
		}
	}

	// Create command using platform-specific launcher
	cmd := launcher.PrepareCommand(ctx, path, e.config.Args)

	// Set working directory if specified
	if e.config.WorkingDir != "" {
		cmd.Dir = os.ExpandEnv(e.config.WorkingDir)
	}

	// Set environment variables
	if len(e.config.Env) > 0 {
		cmd.Env = os.Environ()
		for k, v := range e.config.Env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, os.ExpandEnv(v)))
		}
	}

	// Start the application without waiting
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start application: %w", err)
	}

	// Detach from the process
	if err := cmd.Process.Release(); err != nil {
		// Non-fatal error, log it but don't fail
		// In production, this would be logged
		_ = err // Explicitly ignore
	}

	return nil
}

// String returns a string representation of the action
func (e *AppExecutor) String() string {
	if e.config.Description != "" {
		return e.config.Description
	}
	return fmt.Sprintf("Launch %s", e.config.Command)
}
