package e2e

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

// TestEnvironment manages the complete SilentCast application environment for E2E testing
type TestEnvironment struct {
	t              *testing.T
	TempDir        string
	ConfigDir      string
	AppBinary      string
	LogFile        string
	PidFile        string
	ctx            context.Context
	cancel         context.CancelFunc
	appProcess     *exec.Cmd
	startupTimeout time.Duration
}

// NewTestEnvironment creates a new isolated test environment for E2E testing
func NewTestEnvironment(t *testing.T) *TestEnvironment {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")

	// Create config directory
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	// Determine binary path based on build
	binaryName := "silentcast"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	// Look for binary in build directory (relative to current working directory)
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Try multiple possible locations for the binary
	possiblePaths := []string{
		filepath.Join(wd, "build", binaryName),               // ./build/silentcast
		filepath.Join(wd, "..", "build", binaryName),         // ../build/silentcast
		filepath.Join(wd, "..", "..", "build", binaryName),   // ../../build/silentcast (from test/e2e)
		filepath.Join(wd, binaryName),                        // ./silentcast (current directory)
		filepath.Join(filepath.Dir(wd), "build", binaryName), // parent/build/silentcast
	}

	var appBinary string
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			appBinary = path
			break
		}
	}

	if appBinary == "" {
		// Default to the most likely location relative to app directory
		appDir := wd
		if filepath.Base(wd) == "e2e" {
			appDir = filepath.Join(wd, "..", "..")
		} else if filepath.Base(wd) == "test" {
			appDir = filepath.Join(wd, "..")
		}
		appBinary = filepath.Join(appDir, "build", binaryName)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &TestEnvironment{
		t:              t,
		TempDir:        tempDir,
		ConfigDir:      configDir,
		AppBinary:      appBinary,
		LogFile:        filepath.Join(tempDir, "silentcast.log"),
		PidFile:        filepath.Join(tempDir, "silentcast.pid"),
		ctx:            ctx,
		cancel:         cancel,
		startupTimeout: 30 * time.Second,
	}
}

// SetupSpellbook creates a test spellbook configuration
func (env *TestEnvironment) SetupSpellbook(spellbookContent string) error {
	spellbookPath := filepath.Join(env.ConfigDir, "spellbook.yml")
	return os.WriteFile(spellbookPath, []byte(spellbookContent), 0o600)
}

// StartApplication starts the SilentCast application in the test environment
func (env *TestEnvironment) StartApplication(args ...string) error {
	// Check if binary exists
	if _, err := os.Stat(env.AppBinary); os.IsNotExist(err) {
		return fmt.Errorf("application binary not found at %s", env.AppBinary)
	}

	// Default arguments for E2E testing
	defaultArgs := []string{
		"--config", env.ConfigDir,
		"--log-file", env.LogFile,
		"--no-tray", // Run without system tray for testing
		"--debug",   // Enable debug logging
	}

	allArgs := append([]string{}, defaultArgs...)
	allArgs = append(allArgs, args...)

	// nosec G204: AppBinary is controlled by test framework
	env.appProcess = exec.CommandContext(env.ctx, env.AppBinary, allArgs...)
	env.appProcess.Dir = env.TempDir

	// Set environment variables
	env.appProcess.Env = append(os.Environ(),
		fmt.Sprintf("SILENTCAST_CONFIG_DIR=%s", env.ConfigDir),
		"SILENTCAST_LOG_LEVEL=debug",
	)

	// Start the application
	if err := env.appProcess.Start(); err != nil {
		return fmt.Errorf("failed to start application: %w", err)
	}

	// Write PID file for monitoring
	pidContent := fmt.Sprintf("%d\n", env.appProcess.Process.Pid)
	if err := os.WriteFile(env.PidFile, []byte(pidContent), 0o600); err != nil {
		env.t.Logf("Warning: Failed to write PID file: %v", err)
	}

	return nil
}

// WaitForStartup waits for the application to fully start up
func (env *TestEnvironment) WaitForStartup() error {
	ctx, cancel := context.WithTimeout(env.ctx, env.startupTimeout)
	defer cancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("application startup timeout after %v", env.startupTimeout)
		case <-ticker.C:
			// Check if process is still running
			if env.appProcess.ProcessState != nil && env.appProcess.ProcessState.Exited() {
				return fmt.Errorf("application exited during startup")
			}

			// Check if log file indicates successful startup
			if env.isApplicationReady() {
				return nil
			}
		}
	}
}

// isApplicationReady checks if the application has completed startup
func (env *TestEnvironment) isApplicationReady() bool {
	// Check log file for startup completion indicators
	logContent, err := os.ReadFile(env.LogFile)
	if err != nil {
		return false
	}

	logStr := string(logContent)

	// Look for indicators that the application is ready
	readyIndicators := []string{
		"Application started",
		"Hotkey manager initialized",
		"Configuration loaded",
		"Ready to receive hotkeys",
	}

	for _, indicator := range readyIndicators {
		if contains(logStr, indicator) {
			return true
		}
	}

	return false
}

// StopApplication gracefully stops the SilentCast application
func (env *TestEnvironment) StopApplication() error {
	if env.appProcess == nil || env.appProcess.Process == nil {
		return nil
	}

	// Send termination signal
	if err := env.appProcess.Process.Signal(os.Interrupt); err != nil {
		env.t.Logf("Failed to send interrupt signal: %v", err)
		// Force kill if interrupt fails
		return env.appProcess.Process.Kill()
	}

	// Wait for graceful shutdown with timeout
	done := make(chan error, 1)
	go func() {
		done <- env.appProcess.Wait()
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(5 * time.Second):
		env.t.Log("Graceful shutdown timeout, force killing application")
		return env.appProcess.Process.Kill()
	}
}

// GetLogs returns the current application logs
func (env *TestEnvironment) GetLogs() (string, error) {
	content, err := os.ReadFile(env.LogFile)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// Cleanup performs cleanup operations for the test environment
func (env *TestEnvironment) Cleanup() {
	if env.cancel != nil {
		env.cancel()
	}

	if env.appProcess != nil {
		_ = env.StopApplication()
	}

	// Clean up PID file
	if env.PidFile != "" {
		_ = os.Remove(env.PidFile)
	}
}

// SimulateHotkey simulates a hotkey press sequence
func (env *TestEnvironment) SimulateHotkey(sequence string) error {
	// For E2E testing, we'll send a signal to the application
	// In a real implementation, this might use platform-specific APIs
	// For now, we'll use a test endpoint or signal

	// This is a placeholder - in real implementation, you'd:
	// 1. Use platform-specific key simulation APIs
	// 2. Send test signals to the application
	// 3. Use automation frameworks like Selenium or similar

	env.t.Logf("Simulating hotkey sequence: %s", sequence)

	// For testing, we can create a test file that the application monitors
	testSignalFile := filepath.Join(env.TempDir, "test_signal")
	return os.WriteFile(testSignalFile, []byte(sequence), 0o600)
}

// WaitForAction waits for a specific action to complete
func (env *TestEnvironment) WaitForAction(actionName string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(env.ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for action '%s' after %v", actionName, timeout)
		case <-ticker.C:
			logs, err := env.GetLogs()
			if err != nil {
				continue
			}

			// Look for action completion in logs
			if contains(logs, fmt.Sprintf("Action completed: %s", actionName)) ||
				contains(logs, fmt.Sprintf("Executed action: %s", actionName)) {
				return nil
			}
		}
	}
}

// AssertNoErrors checks that no errors occurred during execution
func (env *TestEnvironment) AssertNoErrors() {
	logs, err := env.GetLogs()
	if err != nil {
		env.t.Fatalf("Failed to read logs: %v", err)
	}

	// Check for error indicators in logs
	errorIndicators := []string{
		"ERROR",
		"FATAL",
		"panic:",
		"runtime error",
	}

	for _, indicator := range errorIndicators {
		if contains(logs, indicator) {
			env.t.Errorf("Found error in logs: %s", indicator)
		}
	}
}

// Helper function to check if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsMiddle(s, substr))))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestSpellbookTemplate provides a basic spellbook for testing
const TestSpellbookTemplate = `
spells:
  e: editor
  t: terminal
  b: browser
  g,s: git-status

grimoire:
  editor:
    type: app
    command: echo
    description: "Test editor action"
    
  terminal:
    type: script
    command: echo "Terminal opened"
    description: "Test terminal action"
    
  browser:
    type: url
    command: "https://example.com"
    description: "Test browser action"
    
  git-status:
    type: script
    command: git status
    description: "Git status check"
    show_output: true
`
