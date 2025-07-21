//go:build linux

package browser

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

// launchBrowser launches a browser with the given arguments on Linux
func launchBrowser(ctx context.Context, executable string, args []string) error {
	// Create command
	cmd := exec.CommandContext(ctx, executable, args...)

	// Set up environment
	cmd.Env = os.Environ()

	// Detach from terminal to prevent blocking
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil

	// Start the browser
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to launch browser: %w", err)
	}

	// Detach the process
	return cmd.Process.Release()
}

// launchURLFallback opens a URL using the system default handler on Linux
func launchURLFallback(ctx context.Context, url string) error {
	// Validate URL first
	validatedURL, err := ValidateURL(url)
	if err != nil {
		return err
	}

	// Try different URL openers in order of preference
	openers := []string{
		"xdg-open",         // Standard freedesktop.org opener
		"gnome-open",       // GNOME
		"kde-open",         // KDE
		"exo-open",         // Xfce
		"sensible-browser", // Debian/Ubuntu fallback
	}

	for _, opener := range openers {
		if _, err := exec.LookPath(opener); err != nil {
			continue
		}

		cmd := exec.CommandContext(ctx, opener, validatedURL)

		// Detach from terminal
		cmd.Stdin = nil
		cmd.Stdout = nil
		cmd.Stderr = nil

		if err := cmd.Start(); err == nil {
			// Successfully started, detach process
			return cmd.Process.Release()
		}
	}

	// If no opener found, try BROWSER environment variable
	if browser := os.Getenv("BROWSER"); browser != "" {
		cmd := exec.CommandContext(ctx, browser, validatedURL)
		cmd.Stdin = nil
		cmd.Stdout = nil
		cmd.Stderr = nil

		if err := cmd.Start(); err == nil {
			return cmd.Process.Release()
		}
	}

	return fmt.Errorf("no URL opener found")
}
