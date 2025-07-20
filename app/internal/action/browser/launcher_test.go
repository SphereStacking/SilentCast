package browser

import (
	"context"
	"os"
	"testing"
)

func TestNewLauncher(t *testing.T) {
	launcher := NewLauncher()
	if launcher == nil {
		t.Fatal("NewLauncher() returned nil")
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		wantErr  bool
	}{
		// Valid URLs
		{"https://example.com", "https://example.com", false},
		{"http://localhost:3000", "http://localhost:3000", false},
		{"file:///home/user/index.html", "file:///home/user/index.html", false},

		// URLs without scheme
		{"example.com", "https://example.com", false},
		{"localhost:8080", "http://localhost:8080", false},
		{"127.0.0.1:3000", "http://127.0.0.1:3000", false},
		{"192.168.1.1", "http://192.168.1.1", false},
		{"10.0.0.1:8080", "http://10.0.0.1:8080", false},

		// Special schemes
		{"about:blank", "about:blank", false},
		{"chrome://settings", "chrome://settings", false},

		// Invalid URLs
		{"", "", true},
		{"ftp://example.com", "", true},   // Unsupported scheme
		{"javascript:alert(1)", "", true}, // Dangerous scheme
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ValidateURL(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateURL(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("ValidateURL(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestEscapeURLForShell(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"https://example.com", "'https://example.com'"},
		{"https://example.com/path?query=value", "'https://example.com/path?query=value'"},
		{"https://example.com/path?q=it's", `"https://example.com/path?q=it's"`},
		{"https://example.com/\"quotes\"", `"https://example.com/\"quotes\""`},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := EscapeURLForShell(tt.input)
			if result != tt.expected {
				t.Errorf("EscapeURLForShell(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetIncognitoArgs(t *testing.T) {
	launcher := NewLauncher()

	tests := []struct {
		browser  string
		expected []string
	}{
		{"chrome", []string{"--incognito"}},
		{"Google Chrome", []string{"--incognito"}},
		{"firefox", []string{"--private-window"}},
		{"Mozilla Firefox", []string{"--private-window"}},
		{"edge", []string{"--incognito"}},
		{"safari", []string{}},               // Safari doesn't support CLI private mode
		{"unknown", []string{"--incognito"}}, // Default to Chrome style
	}

	for _, tt := range tests {
		t.Run(tt.browser, func(t *testing.T) {
			result := launcher.GetIncognitoArgs(tt.browser)
			if len(result) != len(tt.expected) {
				t.Errorf("GetIncognitoArgs(%q) returned %v, want %v", tt.browser, result, tt.expected)
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("GetIncognitoArgs(%q)[%d] = %q, want %q", tt.browser, i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestGetNewWindowArgs(t *testing.T) {
	launcher := NewLauncher()

	tests := []struct {
		browser  string
		expected []string
	}{
		{"chrome", []string{"--new-window"}},
		{"firefox", []string{"--new-window"}},
		{"edge", []string{"--new-window"}},
		{"safari", []string{}}, // Safari uses different mechanism
		{"unknown", []string{"--new-window"}},
	}

	for _, tt := range tests {
		t.Run(tt.browser, func(t *testing.T) {
			result := launcher.GetNewWindowArgs(tt.browser)
			if len(result) != len(tt.expected) {
				t.Errorf("GetNewWindowArgs(%q) returned %v, want %v", tt.browser, result, tt.expected)
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("GetNewWindowArgs(%q)[%d] = %q, want %q", tt.browser, i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestLaunchDefault(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping browser launch test in short mode")
	}

	// Skip in CI environment
	if _, ok := os.LookupEnv("CI"); ok {
		t.Skip("Skipping browser launch test in CI")
	}

	launcher := NewLauncher()
	ctx := context.Background()

	// Test with about:blank (safe URL that won't make network requests)
	err := launcher.LaunchDefault(ctx, "about:blank")
	if err != nil {
		t.Logf("LaunchDefault() error: %v (this might be expected in some environments)", err)
	}

	// Test empty URL
	err = launcher.LaunchDefault(ctx, "")
	if err == nil {
		t.Error("LaunchDefault() with empty URL should return error")
	}
}

func TestLaunch(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping browser launch test in short mode")
	}

	// Skip in CI environment
	if _, ok := os.LookupEnv("CI"); ok {
		t.Skip("Skipping browser launch test in CI")
	}

	launcher := NewLauncher()
	ctx := context.Background()

	// Test launching with default browser
	t.Run("Default browser", func(t *testing.T) {
		opts := LaunchOptions{
			URL: "about:blank",
		}

		err := launcher.Launch(ctx, opts)
		if err != nil {
			t.Logf("Launch() error: %v (this might be expected)", err)
		}
	})

	// Test with browser preference
	t.Run("Browser preference", func(t *testing.T) {
		opts := LaunchOptions{
			URL:               "about:blank",
			BrowserPreference: []string{"firefox", "chrome"},
		}

		err := launcher.Launch(ctx, opts)
		if err != nil {
			t.Logf("Launch() with preference error: %v", err)
		}
	})

	// Test incognito mode
	t.Run("Incognito mode", func(t *testing.T) {
		opts := LaunchOptions{
			URL:       "about:blank",
			Incognito: true,
		}

		err := launcher.Launch(ctx, opts)
		if err != nil {
			t.Logf("Launch() with incognito error: %v", err)
		}
	})

	// Test invalid URL
	t.Run("Invalid URL", func(t *testing.T) {
		opts := LaunchOptions{
			URL: "",
		}

		err := launcher.Launch(ctx, opts)
		if err == nil {
			t.Error("Launch() with empty URL should return error")
		}
	})
}

func TestLaunchWithSpecificBrowser(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping browser launch test in short mode")
	}

	// Skip in CI environment
	if _, ok := os.LookupEnv("CI"); ok {
		t.Skip("Skipping browser launch test in CI")
	}

	launcher := NewLauncher()
	ctx := context.Background()
	detector := NewDetector()

	// Try to find Chrome for testing
	chrome, err := detector.FindBrowser(ctx, "chrome")
	if err != nil {
		t.Skip("Chrome not found, skipping test")
	}

	opts := LaunchOptions{
		URL:       "about:blank",
		Browser:   chrome,
		NewWindow: true,
	}

	err = launcher.Launch(ctx, opts)
	if err != nil {
		t.Logf("Launch() with specific browser error: %v", err)
	}
}

func TestLaunchSafety(t *testing.T) {
	launcher := NewLauncher()
	ctx := context.Background()

	// Test that dangerous URLs are rejected
	dangerousURLs := []string{
		"javascript:alert(1)",
		"data:text/html,<script>alert(1)</script>",
		"vbscript:msgbox(1)",
	}

	for _, url := range dangerousURLs {
		t.Run(url, func(t *testing.T) {
			opts := LaunchOptions{URL: url}
			err := launcher.Launch(ctx, opts)
			if err == nil {
				t.Errorf("Launch() should reject dangerous URL: %s", url)
			}
		})
	}
}

func BenchmarkValidateURL(b *testing.B) {
	urls := []string{
		"https://example.com",
		"example.com",
		"localhost:3000",
		"192.168.1.1:8080",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, url := range urls {
			_, _ = ValidateURL(url)
		}
	}
}

func ExampleLauncher_Launch() {
	launcher := NewLauncher()
	ctx := context.Background()

	// Open URL in default browser
	opts := LaunchOptions{
		URL: "https://example.com",
	}
	_ = launcher.Launch(ctx, opts)

	// Open URL in Firefox in private mode
	opts = LaunchOptions{
		URL:               "https://example.com",
		BrowserPreference: []string{"firefox"},
		Incognito:         true,
	}
	_ = launcher.Launch(ctx, opts)
}
