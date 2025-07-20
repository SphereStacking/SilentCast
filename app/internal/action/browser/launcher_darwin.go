//go:build darwin

package browser

import (
	"context"
	"os/exec"
	"strings"
)

// launchBrowser launches a browser with the given arguments on macOS
func launchBrowser(ctx context.Context, executable string, args []string) error {
	// On macOS, we need to handle .app bundles differently
	if strings.Contains(executable, ".app/") {
		// Extract app path
		parts := strings.Split(executable, ".app/")
		if len(parts) > 0 {
			appPath := parts[0] + ".app"

			// Use open command with -a flag
			openArgs := []string{"-a", appPath}

			// Add --args if we have arguments
			if len(args) > 0 {
				openArgs = append(openArgs, "--args")
				openArgs = append(openArgs, args...)
			}

			cmd := exec.CommandContext(ctx, "open", openArgs...)
			return cmd.Start()
		}
	}

	// For non-.app executables, launch directly
	cmd := exec.CommandContext(ctx, executable, args...)
	return cmd.Start()
}

// launchURLFallback opens a URL using the system default handler on macOS
func launchURLFallback(ctx context.Context, url string) error {
	// Validate URL first
	validatedURL, err := ValidateURL(url)
	if err != nil {
		return err
	}

	// Use open command
	cmd := exec.CommandContext(ctx, "open", validatedURL)
	return cmd.Start()
}
