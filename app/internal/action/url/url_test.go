package url

import (
	"context"
	"strings"
	"testing"

	"github.com/SphereStacking/silentcast/internal/config"
)

func TestURLExecutor_Execute(t *testing.T) {
	tests := []struct {
		name    string
		config  config.ActionConfig
		wantErr bool
	}{
		{
			name: "Valid HTTPS URL",
			config: config.ActionConfig{
				Type:    "url",
				Command: "https://github.com/SphereStacking/SilentCast",
			},
			wantErr: false,
		},
		{
			name: "Valid HTTP URL",
			config: config.ActionConfig{
				Type:    "url",
				Command: "http://example.com",
			},
			wantErr: false,
		},
		{
			name: "URL without scheme (should add https)",
			config: config.ActionConfig{
				Type:    "url",
				Command: "github.com/SphereStacking/SilentCast",
			},
			wantErr: false,
		},
		{
			name: "File URL",
			config: config.ActionConfig{
				Type:    "url",
				Command: "file:///home/user/document.html",
			},
			wantErr: false,
		},
		{
			name: "Mailto URL",
			config: config.ActionConfig{
				Type:    "url",
				Command: "mailto:test@example.com",
			},
			wantErr: false,
		},
		{
			name: "Empty URL",
			config: config.ActionConfig{
				Type:    "url",
				Command: "",
			},
			wantErr: true,
		},
		{
			name: "Invalid URL scheme",
			config: config.ActionConfig{
				Type:    "url",
				Command: "javascript:alert('test')",
			},
			wantErr: true,
		},
		{
			name: "URL with spaces (should trim)",
			config: config.ActionConfig{
				Type:    "url",
				Command: "  https://example.com  ",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := NewURLExecutor(&tt.config)
			ctx := context.Background()

			// For testing, we don't actually want to open browsers
			// So we'll just check if the URL validation works
			err := executor.Execute(ctx)

			// In CI/test environment, browser opening might fail
			// but URL validation should still work
			if err != nil && !tt.wantErr {
				// Check if it's a browser launch error (expected in tests)
				if strings.Contains(err.Error(), "failed to open URL") ||
					strings.Contains(err.Error(), "executable file not found") {
					t.Logf("Browser launch failed (expected in test environment): %v", err)
					return
				}
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			} else if err == nil && tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestURLExecutor_String(t *testing.T) {
	tests := []struct {
		name     string
		config   config.ActionConfig
		expected string
	}{
		{
			name: "With description",
			config: config.ActionConfig{
				Type:        "url",
				Command:     "https://github.com",
				Description: "Open GitHub",
			},
			expected: "Open GitHub",
		},
		{
			name: "Without description",
			config: config.ActionConfig{
				Type:    "url",
				Command: "https://example.com",
			},
			expected: "Open URL: https://example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := NewURLExecutor(&tt.config)
			if got := executor.String(); got != tt.expected {
				t.Errorf("String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestURLValidation(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"Valid HTTPS", "https://example.com", false},
		{"Valid HTTP", "http://example.com", false},
		{"No scheme", "example.com", false}, // Should add https://
		{"With path", "https://example.com/path/to/page", false},
		{"With query", "https://example.com?q=test", false},
		{"With fragment", "https://example.com#section", false},
		{"FTP URL", "ftp://ftp.example.com", false},
		{"File URL", "file:///C:/test.html", false},
		{"Mailto", "mailto:test@example.com", false},
		{"Invalid scheme", "javascript:void(0)", true},
		{"Data URL", "data:text/html,<h1>Test</h1>", true},
		{"Empty", "", true},
		{"Just spaces", "   ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := config.ActionConfig{
				Type:    "url",
				Command: tt.url,
			}
			executor := NewURLExecutor(&config)
			err := executor.Execute(context.Background())

			// Check URL validation specifically
			if tt.wantErr {
				if err == nil {
					t.Error("Expected error for invalid URL")
				} else if !strings.Contains(err.Error(), "URL") &&
					!strings.Contains(err.Error(), "scheme") {
					// Should be a URL-related error
					t.Errorf("Expected URL validation error, got: %v", err)
				}
			}
		})
	}
}