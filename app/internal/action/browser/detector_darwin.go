//go:build darwin

package browser

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type darwinDetector struct{}

func newPlatformDetector() Detector {
	return &darwinDetector{}
}

func (d *darwinDetector) DetectBrowsers(ctx context.Context) ([]Browser, error) {
	var browsers []Browser

	// Check Applications folder
	browsers = append(browsers, d.detectFromApplications()...)

	// Get default browser
	defaultBrowser, _ := d.GetDefaultBrowser(ctx)
	if defaultBrowser != nil {
		for i := range browsers {
			if browsers[i].Executable == defaultBrowser.Executable {
				browsers[i].IsDefault = true
			}
		}
	}

	return browsers, nil
}

func (d *darwinDetector) GetDefaultBrowser(ctx context.Context) (*Browser, error) {
	// Use LSCopyDefaultHandlerForURLScheme to get default browser
	cmd := exec.CommandContext(ctx, "defaults", "read",
		"com.apple.LaunchServices/com.apple.launchservices.secure",
		"LSHandlers")
	_, err := cmd.Output()
	if err != nil {
		// Try alternative method using open command
		return d.getDefaultBrowserUsingOpen(ctx)
	}

	// Parse output to find HTTP handler
	// This is complex, so we'll use the simpler open method
	return d.getDefaultBrowserUsingOpen(ctx)
}

func (d *darwinDetector) getDefaultBrowserUsingOpen(ctx context.Context) (*Browser, error) {
	// First, open a URL to ensure browser is running
	openCmd := exec.CommandContext(ctx, "open", "-g", "https://")
	_ = openCmd.Run()

	// Get default browser bundle ID
	cmd := exec.CommandContext(ctx, "osascript", "-e",
		`tell application "System Events" to get bundle identifier of (first process whose frontmost is true)`)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get default browser: %w", err)
	}

	bundleID := strings.TrimSpace(string(output))

	// Find browser by bundle ID
	browsers, _ := d.DetectBrowsers(ctx)
	for i := range browsers {
		if d.getBundleID(browsers[i].Executable) == bundleID {
			browsers[i].IsDefault = true
			return &browsers[i], nil
		}
	}

	return nil, ErrNoBrowserFound
}

func (d *darwinDetector) FindBrowser(ctx context.Context, name string) (*Browser, error) {
	browsers, err := d.DetectBrowsers(ctx)
	if err != nil {
		return nil, err
	}

	for i := range browsers {
		if MatchBrowser(&browsers[i], name) {
			return &browsers[i], nil
		}
	}

	return nil, ErrNoBrowserFound
}

func (d *darwinDetector) detectFromApplications() []Browser {
	var browsers []Browser

	// Common browser apps
	apps := []struct {
		name     string
		appNames []string
	}{
		{
			name:     "Google Chrome",
			appNames: []string{"Google Chrome.app", "Chrome.app"},
		},
		{
			name:     "Mozilla Firefox",
			appNames: []string{"Firefox.app"},
		},
		{
			name:     "Safari",
			appNames: []string{"Safari.app"},
		},
		{
			name:     "Microsoft Edge",
			appNames: []string{"Microsoft Edge.app"},
		},
		{
			name:     "Opera",
			appNames: []string{"Opera.app"},
		},
		{
			name:     "Brave Browser",
			appNames: []string{"Brave Browser.app"},
		},
		{
			name:     "Vivaldi",
			appNames: []string{"Vivaldi.app"},
		},
		{
			name:     "Arc",
			appNames: []string{"Arc.app"},
		},
	}

	// Check multiple locations
	searchPaths := []string{
		"/Applications",
		filepath.Join(os.Getenv("HOME"), "Applications"),
		"/System/Applications",
	}

	for _, browser := range apps {
		for _, appName := range browser.appNames {
			for _, searchPath := range searchPaths {
				appPath := filepath.Join(searchPath, appName)
				if info, err := os.Stat(appPath); err == nil && info.IsDir() {
					// Get executable path
					executable := filepath.Join(appPath, "Contents", "MacOS",
						strings.TrimSuffix(appName, ".app"))

					// For some apps, the executable name might be different
					if _, err := os.Stat(executable); err != nil {
						// Try alternative names
						altNames := []string{
							strings.ReplaceAll(strings.TrimSuffix(appName, ".app"), " ", ""),
							strings.ToLower(strings.ReplaceAll(strings.TrimSuffix(appName, ".app"), " ", "")),
						}

						for _, altName := range altNames {
							altExec := filepath.Join(appPath, "Contents", "MacOS", altName)
							if _, err := os.Stat(altExec); err == nil {
								executable = altExec
								break
							}
						}
					}

					if _, err := os.Stat(executable); err == nil {
						// Get version if possible
						version := d.getAppVersion(appPath)

						browsers = append(browsers, Browser{
							Name:       browser.name,
							Executable: executable,
							Version:    version,
						})

						// Found this browser, move to next
						goto nextBrowser
					}
				}
			}
		}
	nextBrowser:
	}

	return browsers
}

func (d *darwinDetector) getAppVersion(appPath string) string {
	// Read version from Info.plist
	plistPath := filepath.Join(appPath, "Contents", "Info.plist")

	cmd := exec.Command("defaults", "read", plistPath, "CFBundleShortVersionString")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(output))
}

func (d *darwinDetector) getBundleID(executable string) string {
	// Get bundle ID from app path
	if strings.Contains(executable, ".app/") {
		parts := strings.Split(executable, ".app/")
		if len(parts) > 0 {
			appPath := parts[0] + ".app"
			plistPath := filepath.Join(appPath, "Contents", "Info.plist")

			cmd := exec.Command("defaults", "read", plistPath, "CFBundleIdentifier")
			output, err := cmd.Output()
			if err == nil {
				return strings.TrimSpace(string(output))
			}
		}
	}

	return ""
}

func openURLFallback(ctx context.Context, url string) error {
	// Use macOS open command
	cmd := exec.CommandContext(ctx, "open", url)
	return cmd.Start()
}
