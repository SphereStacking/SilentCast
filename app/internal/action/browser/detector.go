package browser

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// Browser represents a detected browser
type Browser struct {
	Name       string // e.g., "Google Chrome", "Firefox", "Safari"
	Executable string // Full path to executable
	Version    string // Version string (if available)
	IsDefault  bool   // Whether this is the system default browser
}

// Detector is the interface for browser detection
type Detector interface {
	// DetectBrowsers returns a list of installed browsers
	DetectBrowsers(ctx context.Context) ([]Browser, error)

	// GetDefaultBrowser returns the system default browser
	GetDefaultBrowser(ctx context.Context) (*Browser, error)

	// FindBrowser finds a specific browser by name
	FindBrowser(ctx context.Context, name string) (*Browser, error)
}

// ErrNoBrowserFound is returned when no browsers are found
var ErrNoBrowserFound = errors.New("no browser found")

// CommonBrowserNames contains common browser names for searching
var CommonBrowserNames = []string{
	"chrome",
	"google-chrome",
	"chromium",
	"firefox",
	"safari",
	"edge",
	"microsoft-edge",
	"opera",
	"brave",
	"vivaldi",
}

// NewDetector creates a platform-specific browser detector
func NewDetector() Detector {
	return newPlatformDetector()
}

// OpenURL opens a URL with the specified browser
func OpenURL(ctx context.Context, browser *Browser, url string) error {
	if browser == nil {
		return errors.New("browser is nil")
	}

	if browser.Executable == "" {
		return errors.New("browser executable not specified")
	}

	cmd := exec.CommandContext(ctx, browser.Executable, url)
	return cmd.Start()
}

// OpenURLWithDefault opens a URL with the system default browser
func OpenURLWithDefault(ctx context.Context, url string) error {
	launcher := NewLauncher()
	return launcher.LaunchDefault(ctx, url)
}

// NormalizeBrowserName normalizes browser names for comparison
func NormalizeBrowserName(name string) string {
	normalized := strings.ToLower(name)
	normalized = strings.ReplaceAll(normalized, " ", "-")
	normalized = strings.ReplaceAll(normalized, "_", "-")

	// Handle common variations
	replacements := map[string]string{
		"google-chrome":    "chrome",
		"chromium-browser": "chromium",
		"microsoft-edge":   "edge",
		"msedge":           "edge",
	}

	if replacement, ok := replacements[normalized]; ok {
		return replacement
	}

	return normalized
}

// MatchBrowser checks if a browser matches the given name
func MatchBrowser(browser *Browser, name string) bool {
	if browser == nil {
		return false
	}

	normalizedName := NormalizeBrowserName(name)
	normalizedBrowserName := NormalizeBrowserName(browser.Name)

	// Check exact match
	if normalizedBrowserName == normalizedName {
		return true
	}

	// Check if browser name contains the search name
	if strings.Contains(normalizedBrowserName, normalizedName) {
		return true
	}

	// Check executable name
	execName := strings.ToLower(browser.Executable)
	return strings.Contains(execName, normalizedName)
}

// GetBrowserByPreference returns a browser based on preference list
func GetBrowserByPreference(browsers []Browser, preferences []string) *Browser {
	for _, pref := range preferences {
		for i := range browsers {
			if MatchBrowser(&browsers[i], pref) {
				return &browsers[i]
			}
		}
	}

	// Return first browser if no preference matches
	if len(browsers) > 0 {
		return &browsers[0]
	}

	return nil
}

// FormatBrowserInfo formats browser information for display
func FormatBrowserInfo(browser *Browser) string {
	if browser == nil {
		return "No browser"
	}

	info := browser.Name
	if browser.Version != "" {
		info += fmt.Sprintf(" (%s)", browser.Version)
	}
	if browser.IsDefault {
		info += " [default]"
	}

	return info
}
