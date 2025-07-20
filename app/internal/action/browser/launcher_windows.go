//go:build windows

package browser

import (
	"context"
	"os/exec"
	"strings"
	"syscall"
)

// launchBrowser launches a browser with the given arguments on Windows
func launchBrowser(ctx context.Context, executable string, args []string) error {
	// Use exec.CommandContext for proper context handling
	cmd := exec.CommandContext(ctx, executable, args...)

	// Hide console window
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	// Start the browser
	return cmd.Start()
}

// launchURLFallback opens a URL using the system default handler on Windows
func launchURLFallback(ctx context.Context, url string) error {
	// Validate URL first
	validatedURL, err := ValidateURL(url)
	if err != nil {
		return err
	}

	// Use cmd /c start to open URL
	// This is safer than ShellExecute for URLs
	args := []string{"/c", "start", ""}

	// For URLs with special characters, we need to escape them
	if strings.ContainsAny(validatedURL, "&^%<>|") {
		// Use explicit browser call instead
		return openURLWithExplorer(ctx, validatedURL)
	}

	args = append(args, validatedURL)

	cmd := exec.CommandContext(ctx, "cmd", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	return cmd.Start()
}

// openURLWithExplorer opens URL using Windows Explorer
func openURLWithExplorer(ctx context.Context, url string) error {
	// Use explorer.exe to open URLs safely
	cmd := exec.CommandContext(ctx, "explorer.exe", url)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	return cmd.Start()
}
