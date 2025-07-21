package browser

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// LaunchOptions contains options for launching a browser
type LaunchOptions struct {
	// Browser to use (nil for default)
	Browser *Browser

	// URL to open
	URL string

	// Additional command-line arguments
	Args []string

	// Open in incognito/private mode
	Incognito bool

	// Open in new window
	NewWindow bool

	// Open in new tab (default behavior for most browsers)
	NewTab bool

	// Browser preference list (try in order)
	BrowserPreference []string
}

// Launcher is the interface for launching browsers
type Launcher interface {
	// Launch opens a URL with the specified options
	Launch(ctx context.Context, opts LaunchOptions) error

	// LaunchDefault opens a URL with the system default browser
	LaunchDefault(ctx context.Context, url string) error

	// GetIncognitoArgs returns browser-specific incognito mode arguments
	GetIncognitoArgs(browserName string) []string

	// GetNewWindowArgs returns browser-specific new window arguments
	GetNewWindowArgs(browserName string) []string
}

// launcher implements the Launcher interface
type launcher struct {
	detector Detector
}

// NewLauncher creates a new browser launcher
func NewLauncher() Launcher {
	return &launcher{
		detector: NewDetector(),
	}
}

// Launch opens a URL with the specified options
func (l *launcher) Launch(ctx context.Context, opts LaunchOptions) error {
	// Validate URL
	if opts.URL == "" {
		return errors.New("URL is required")
	}

	// Parse URL to ensure it's valid
	parsedURL, err := url.Parse(opts.URL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// Ensure URL has a scheme
	if parsedURL.Scheme == "" {
		// Default to https for URLs without scheme
		if strings.Contains(opts.URL, "localhost") || strings.HasPrefix(opts.URL, "127.0.0.1") {
			opts.URL = "http://" + opts.URL
		} else {
			opts.URL = "https://" + opts.URL
		}
	}

	// Determine which browser to use
	var browser *Browser

	if opts.Browser != nil {
		// Use specified browser
		browser = opts.Browser
	} else if len(opts.BrowserPreference) > 0 {
		// Try browsers in preference order
		browsers, detectErr := l.detector.DetectBrowsers(ctx)
		if detectErr != nil {
			return fmt.Errorf("failed to detect browsers: %w", detectErr)
		}

		browser = GetBrowserByPreference(browsers, opts.BrowserPreference)
		if browser == nil {
			// Fall back to default if no preference matches
			browser, _ = l.detector.GetDefaultBrowser(ctx) // Error handled by checking if browser is nil
		}
	} else {
		// Use default browser
		browser, err = l.detector.GetDefaultBrowser(ctx)
		if err != nil {
			// Fall back to platform-specific URL opener
			return l.LaunchDefault(ctx, opts.URL)
		}
	}

	if browser == nil {
		return ErrNoBrowserFound
	}

	// Build command arguments
	args := make([]string, 0)

	// Add incognito mode if requested
	if opts.Incognito {
		args = append(args, l.GetIncognitoArgs(browser.Name)...)
	}

	// Add new window if requested
	if opts.NewWindow {
		args = append(args, l.GetNewWindowArgs(browser.Name)...)
	}

	// Add custom arguments
	args = append(args, opts.Args...)

	// Add URL as last argument
	args = append(args, opts.URL)

	// Launch browser
	return launchBrowser(ctx, browser.Executable, args)
}

// LaunchDefault opens a URL with the system default browser
func (l *launcher) LaunchDefault(ctx context.Context, url string) error {
	if url == "" {
		return errors.New("URL is required")
	}

	return launchURLFallback(ctx, url)
}

// GetIncognitoArgs returns browser-specific incognito mode arguments
func (l *launcher) GetIncognitoArgs(browserName string) []string {
	normalized := NormalizeBrowserName(browserName)

	switch {
	case strings.Contains(normalized, "chrome") ||
		normalized == "chromium" ||
		normalized == "edge" ||
		normalized == "brave" ||
		normalized == "opera" ||
		normalized == "vivaldi":
		return []string{"--incognito"}
	case strings.Contains(normalized, "firefox"):
		return []string{"--private-window"}
	case normalized == "safari":
		// Safari doesn't support command-line private mode
		return []string{}
	default:
		// Try Chrome-style first as it's more common
		return []string{"--incognito"}
	}
}

// GetNewWindowArgs returns browser-specific new window arguments
func (l *launcher) GetNewWindowArgs(browserName string) []string {
	normalized := NormalizeBrowserName(browserName)

	switch normalized {
	case "chrome", "chromium", "edge", "brave", "opera", "vivaldi":
		return []string{"--new-window"}
	case "firefox":
		return []string{"--new-window"}
	case "safari":
		// Safari uses different mechanism
		return []string{}
	default:
		return []string{"--new-window"}
	}
}

// ValidateURL validates and normalizes a URL
func ValidateURL(rawURL string) (string, error) {
	if rawURL == "" {
		return "", errors.New("empty URL")
	}

	// Check if it has a scheme or is a special URL
	if strings.Contains(rawURL, "://") || strings.HasPrefix(rawURL, "about:") {
		// Parse URL with scheme
		parsedURL, err := url.Parse(rawURL)
		if err != nil {
			return "", fmt.Errorf("invalid URL: %w", err)
		}

		// Validate scheme
		switch parsedURL.Scheme {
		case "http", "https", "file", "about", "chrome", "edge":
			// Valid schemes
			return parsedURL.String(), nil
		case "javascript", "vbscript", "data":
			// Dangerous schemes
			return "", fmt.Errorf("dangerous URL scheme: %s", parsedURL.Scheme)
		default:
			return "", fmt.Errorf("unsupported URL scheme: %s", parsedURL.Scheme)
		}
	}

	// No scheme - need to add one
	// Handle localhost and IPs specially
	if strings.HasPrefix(rawURL, "localhost") ||
		strings.HasPrefix(rawURL, "127.0.0.1") ||
		strings.HasPrefix(rawURL, "192.168.") ||
		strings.HasPrefix(rawURL, "10.") ||
		strings.HasPrefix(rawURL, "172.") {
		rawURL = "http://" + rawURL
	} else {
		// Default to https for everything else
		rawURL = "https://" + rawURL
	}

	// Parse with added scheme
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL after adding scheme: %w", err)
	}

	return parsedURL.String(), nil
}

// EscapeURLForShell escapes a URL for safe use in shell commands
func EscapeURLForShell(url string) string {
	// For safety, quote the entire URL
	// Most browsers handle quoted URLs correctly
	if strings.Contains(url, "'") {
		// If URL contains single quotes, use double quotes and escape any double quotes
		return `"` + strings.ReplaceAll(url, `"`, `\"`) + `"`
	}
	if strings.Contains(url, `"`) {
		// If URL contains double quotes, escape them and use double quotes
		return `"` + strings.ReplaceAll(url, `"`, `\"`) + `"`
	}
	// Otherwise use single quotes for better safety
	return "'" + url + "'"
}
