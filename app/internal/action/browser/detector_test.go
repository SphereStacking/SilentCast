package browser

import (
	"context"
	"errors"
	"os"
	"runtime"
	"testing"
)

func TestNewDetector(t *testing.T) {
	detector := NewDetector()
	if detector == nil {
		t.Fatal("NewDetector() returned nil")
	}

	// Just verify we got a detector - the specific type depends on build tags
	t.Logf("Got detector for platform: %s", runtime.GOOS)
}

func TestDetectBrowsers(t *testing.T) {
	detector := NewDetector()
	ctx := context.Background()

	browsers, err := detector.DetectBrowsers(ctx)
	if err != nil {
		t.Logf("DetectBrowsers() error: %v", err)
	}

	// We should find at least one browser on most systems
	// But in CI environments, there might be none
	t.Logf("Found %d browsers", len(browsers))

	for _, browser := range browsers {
		t.Logf("Browser: %s", FormatBrowserInfo(&browser))
		t.Logf("  Executable: %s", browser.Executable)
		if browser.Version != "" {
			t.Logf("  Version: %s", browser.Version)
		}
	}
}

func TestGetDefaultBrowser(t *testing.T) {
	detector := NewDetector()
	ctx := context.Background()

	browser, err := detector.GetDefaultBrowser(ctx)
	if err != nil {
		t.Logf("GetDefaultBrowser() error: %v (this is expected in CI)", err)
		return
	}

	if browser == nil {
		t.Error("GetDefaultBrowser() returned nil browser with no error")
		return
	}

	t.Logf("Default browser: %s", FormatBrowserInfo(browser))
	t.Logf("  Executable: %s", browser.Executable)

	if !browser.IsDefault {
		t.Error("Default browser should have IsDefault=true")
	}
}

func TestNormalizeBrowserName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Google Chrome", "chrome"},
		{"google-chrome", "chrome"},
		{"Microsoft Edge", "edge"},
		{"msedge", "edge"},
		{"Chromium Browser", "chromium"},
		{"firefox", "firefox"},
		{"Mozilla Firefox", "mozilla-firefox"},
		{"Brave_Browser", "brave-browser"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := NormalizeBrowserName(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeBrowserName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMatchBrowser(t *testing.T) {
	tests := []struct {
		browser  *Browser
		name     string
		expected bool
	}{
		{
			&Browser{Name: "Google Chrome", Executable: "/usr/bin/google-chrome"},
			"chrome",
			true,
		},
		{
			&Browser{Name: "Mozilla Firefox", Executable: "/usr/bin/firefox"},
			"firefox",
			true,
		},
		{
			&Browser{Name: "Microsoft Edge", Executable: "C:\\Program Files\\Microsoft\\Edge\\Application\\msedge.exe"},
			"edge",
			true,
		},
		{
			&Browser{Name: "Safari", Executable: "/Applications/Safari.app/Contents/MacOS/Safari"},
			"chrome",
			false,
		},
		{
			nil,
			"chrome",
			false,
		},
	}

	for _, tt := range tests {
		name := tt.name
		if tt.browser != nil {
			name = tt.browser.Name + " vs " + tt.name
		}
		t.Run(name, func(t *testing.T) {
			result := MatchBrowser(tt.browser, tt.name)
			if result != tt.expected {
				t.Errorf("MatchBrowser(%v, %q) = %v, want %v", tt.browser, tt.name, result, tt.expected)
			}
		})
	}
}

func TestGetBrowserByPreference(t *testing.T) {
	browsers := []Browser{
		{Name: "Mozilla Firefox", Executable: "/usr/bin/firefox"},
		{Name: "Google Chrome", Executable: "/usr/bin/chrome"},
		{Name: "Safari", Executable: "/Applications/Safari.app/Contents/MacOS/Safari"},
	}

	tests := []struct {
		name        string
		preferences []string
		expected    string
	}{
		{
			"Prefer Chrome",
			[]string{"chrome", "firefox"},
			"Google Chrome",
		},
		{
			"Prefer Firefox",
			[]string{"firefox", "chrome"},
			"Mozilla Firefox",
		},
		{
			"Prefer unavailable browser",
			[]string{"edge", "opera"},
			"Mozilla Firefox", // First in list
		},
		{
			"Empty preferences",
			[]string{},
			"Mozilla Firefox", // First in list
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetBrowserByPreference(browsers, tt.preferences)
			if result == nil {
				t.Error("GetBrowserByPreference returned nil")
				return
			}
			if result.Name != tt.expected {
				t.Errorf("GetBrowserByPreference() = %q, want %q", result.Name, tt.expected)
			}
		})
	}

	// Test with empty browser list
	result := GetBrowserByPreference([]Browser{}, []string{"chrome"})
	if result != nil {
		t.Error("GetBrowserByPreference with empty browser list should return nil")
	}
}

func TestFormatBrowserInfo(t *testing.T) {
	tests := []struct {
		browser  *Browser
		expected string
	}{
		{
			nil,
			"No browser",
		},
		{
			&Browser{Name: "Google Chrome"},
			"Google Chrome",
		},
		{
			&Browser{Name: "Firefox", Version: "120.0"},
			"Firefox (120.0)",
		},
		{
			&Browser{Name: "Safari", Version: "17.0", IsDefault: true},
			"Safari (17.0) [default]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := FormatBrowserInfo(tt.browser)
			if result != tt.expected {
				t.Errorf("FormatBrowserInfo() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFindBrowser(t *testing.T) {
	detector := NewDetector()
	ctx := context.Background()

	// Try to find common browsers
	browserNames := []string{"chrome", "firefox", "safari", "edge"}

	for _, name := range browserNames {
		browser, err := detector.FindBrowser(ctx, name)
		if errors.Is(err, ErrNoBrowserFound) {
			t.Logf("Browser %q not found (expected in CI)", name)
			continue
		}
		if err != nil {
			t.Errorf("FindBrowser(%q) error: %v", name, err)
			continue
		}

		t.Logf("Found %q: %s", name, FormatBrowserInfo(browser))
	}
}

func TestOpenURLWithDefault(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping browser launch test in short mode")
	}

	// This test would actually open a browser, so we skip it in CI
	if _, ok := os.LookupEnv("CI"); ok {
		t.Skip("Skipping browser launch test in CI")
	}

	ctx := context.Background()

	// Test with a harmless URL
	err := OpenURLWithDefault(ctx, "about:blank")
	if err != nil {
		t.Logf("OpenURLWithDefault() error: %v (this might be expected)", err)
	}
}
