package url

import (
	"context"
	"fmt"
	"net/url"
	"os/exec"
	"runtime"
	"strings"

	"github.com/SphereStacking/silentcast/internal/config"
	"github.com/SphereStacking/silentcast/pkg/logger"
)

// URLExecutor opens URLs in the default browser
type URLExecutor struct {
	config config.ActionConfig
}

// NewURLExecutor creates a new URL executor
func NewURLExecutor(cfg *config.ActionConfig) *URLExecutor {
	return &URLExecutor{
		config: *cfg,
	}
}

// Execute opens the URL in the default browser
func (e *URLExecutor) Execute(ctx context.Context) error {
	// Validate URL
	urlStr := strings.TrimSpace(e.config.Command)
	if urlStr == "" {
		return fmt.Errorf("empty URL")
	}

	// Parse and validate URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// Ensure URL has a scheme
	if parsedURL.Scheme == "" {
		// Default to https
		urlStr = "https://" + urlStr
		parsedURL, err = url.Parse(urlStr)
		if err != nil {
			return fmt.Errorf("invalid URL after adding scheme: %w", err)
		}
	}

	// Validate scheme
	validSchemes := []string{"http", "https", "file", "ftp", "mailto"}
	schemeValid := false
	for _, scheme := range validSchemes {
		if parsedURL.Scheme == scheme {
			schemeValid = true
			break
		}
	}
	if !schemeValid {
		return fmt.Errorf("unsupported URL scheme: %s", parsedURL.Scheme)
	}

	logger.Debug("Opening URL: %s", urlStr)

	// Open URL based on platform
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.CommandContext(ctx, "cmd", "/c", "start", "", urlStr)
	case "darwin":
		cmd = exec.CommandContext(ctx, "open", urlStr)
	case "linux":
		// Try xdg-open first, then fallback to specific browsers
		cmd = e.getLinuxCommand(ctx, urlStr)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	// Execute command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open URL: %w", err)
	}

	// Don't wait for browser to close
	if err := cmd.Process.Release(); err != nil {
		// Non-fatal error
		logger.Debug("Failed to release browser process: %v", err)
	}

	return nil
}

// String returns a string representation of the action
func (e *URLExecutor) String() string {
	if e.config.Description != "" {
		return e.config.Description
	}
	return fmt.Sprintf("Open URL: %s", e.config.Command)
}

// getLinuxCommand returns the appropriate command for Linux
func (e *URLExecutor) getLinuxCommand(ctx context.Context, urlStr string) *exec.Cmd {
	// Try xdg-open first
	if _, err := exec.LookPath("xdg-open"); err == nil {
		return exec.CommandContext(ctx, "xdg-open", urlStr)
	}

	// Try common browsers
	browsers := []string{
		"firefox",
		"google-chrome",
		"chromium",
		"brave-browser",
		"opera",
		"vivaldi",
	}

	for _, browser := range browsers {
		if _, err := exec.LookPath(browser); err == nil {
			return exec.CommandContext(ctx, browser, urlStr)
		}
	}

	// Fallback to xdg-open even if not found
	return exec.CommandContext(ctx, "xdg-open", urlStr)
}
