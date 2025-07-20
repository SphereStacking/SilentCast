package output

import (
	"errors"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestFormatter_stripANSI(t *testing.T) {
	f := &Formatter{}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no ANSI codes",
			input:    "Hello, world!",
			expected: "Hello, world!",
		},
		{
			name:     "color codes",
			input:    "\x1b[31mRed text\x1b[0m",
			expected: "Red text",
		},
		{
			name:     "multiple codes",
			input:    "\x1b[1m\x1b[32mBold green\x1b[0m normal \x1b[4munderline\x1b[0m",
			expected: "Bold green normal underline",
		},
		{
			name:     "cursor movement",
			input:    "Text\x1b[2J\x1b[H",
			expected: "Text",
		},
		{
			name:     "complex sequence",
			input:    "\x1b[2;37;41mWhite on red\x1b[0m",
			expected: "White on red",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := f.stripANSI(tt.input)
			if result != tt.expected {
				t.Errorf("stripANSI() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFormatter_normalizeWhitespace(t *testing.T) {
	f := &Formatter{}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "tabs to spaces",
			input:    "Hello\tworld",
			expected: "Hello    world",
		},
		{
			name:     "multiple newlines",
			input:    "Line1\n\n\n\nLine2",
			expected: "Line1\n\nLine2",
		},
		{
			name:     "trailing spaces",
			input:    "Line with spaces   \nAnother line  ",
			expected: "Line with spaces\nAnother line",
		},
		{
			name:     "leading and trailing whitespace",
			input:    "\n\n  Content  \n\n",
			expected: "Content",
		},
		{
			name:     "mixed whitespace",
			input:    "\tIndented\t\n\n\n  Spaced  \n",
			expected: "Indented\n\n  Spaced",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := f.normalizeWhitespace(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeWhitespace() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFormatter_highlightErrors(t *testing.T) {
	f := &Formatter{}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "error prefix",
			input:    "Error: Something went wrong",
			expected: "❌ Error: Something went wrong",
		},
		{
			name:     "fatal prefix",
			input:    "Fatal: System crash",
			expected: "💀 Fatal: System crash",
		},
		{
			name:     "warning prefix",
			input:    "Warning: Low memory",
			expected: "⚠️  Warning: Low memory",
		},
		{
			name:     "failed to pattern",
			input:    "Failed to connect",
			expected: "❌ Failed to connect",
		},
		{
			name:     "cannot pattern",
			input:    "Cannot open file",
			expected: "❌ Cannot open file",
		},
		{
			name:     "unable to pattern",
			input:    "Unable to parse",
			expected: "❌ Unable to parse",
		},
		{
			name:     "case insensitive",
			input:    "ERROR: uppercase",
			expected: "❌ ERROR: uppercase",
		},
		{
			name:     "indented error",
			input:    "  Error: Indented",
			expected: "  ❌ Error: Indented",
		},
		{
			name:     "no error pattern",
			input:    "Success: All good",
			expected: "Success: All good",
		},
		{
			name:     "already has emoji",
			input:    "❌ Error: Already marked",
			expected: "❌ Error: Already marked",
		},
		{
			name:     "multiple lines with errors",
			input:    "Line 1\nError: Line 2\nLine 3\nWarning: Line 4",
			expected: "Line 1\n❌ Error: Line 2\nLine 3\n⚠️  Warning: Line 4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := f.highlightErrors(tt.input)
			if result != tt.expected {
				t.Errorf("highlightErrors() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFormatter_truncate(t *testing.T) {
	f := &Formatter{}

	tests := []struct {
		name        string
		input       string
		maxLength   int
		checkFunc   func(string) bool
		expectExact string
	}{
		{
			name:        "no truncation needed",
			input:       "Short text",
			maxLength:   20,
			expectExact: "Short text",
		},
		{
			name:        "exact length",
			input:       "12345",
			maxLength:   5,
			expectExact: "12345",
		},
		{
			name:      "simple truncation",
			input:     "This is a long text that needs truncation because it's too long",
			maxLength: 40,
			checkFunc: func(result string) bool {
				return strings.HasSuffix(result, "\n\n... (output truncated)") &&
					utf8.RuneCountInString(result) == 40
			},
		},
		{
			name:      "truncate at newline",
			input:     "Line 1\nLine 2\nLine 3\nLine 4\nLine 5\nLine 6",
			maxLength: 35,
			checkFunc: func(result string) bool {
				return strings.HasSuffix(result, "\n\n... (output truncated)") &&
					utf8.RuneCountInString(result) == 35 &&
					strings.HasPrefix(result, "Line 1")
			},
		},
		{
			name:      "unicode support",
			input:     "Hello 世界 🌍 Test with more text to make it longer",
			maxLength: 30,
			checkFunc: func(result string) bool {
				return strings.HasSuffix(result, "\n\n... (output truncated)") &&
					utf8.RuneCountInString(result) == 30
			},
		},
		{
			name:        "very small max length",
			input:       "Long text",
			maxLength:   5,
			expectExact: "\n\n... (output truncated)",
		},
		{
			name:        "zero max length",
			input:       "Text",
			maxLength:   0,
			expectExact: "Text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := f.truncate(tt.input, tt.maxLength)
			if tt.expectExact != "" {
				if result != tt.expectExact {
					t.Errorf("truncate() = %q, want %q", result, tt.expectExact)
				}
			} else if tt.checkFunc != nil {
				if !tt.checkFunc(result) {
					t.Errorf("truncate() = %q, check function failed", result)
				}
			}
		})
	}
}

func TestFormatter_Format(t *testing.T) {
	tests := []struct {
		name      string
		formatter *Formatter
		input     string
		expected  string
	}{
		{
			name:      "default formatter with ANSI",
			formatter: DefaultFormatter(),
			input:     "\x1b[32mSuccess\x1b[0m: Operation completed\n\n\n",
			expected:  "Success: Operation completed",
		},
		{
			name:      "error highlighting",
			formatter: DefaultFormatter(),
			input:     "Error: Connection failed\nRetrying...",
			expected:  "❌ Error: Connection failed\nRetrying...",
		},
		{
			name: "no ANSI removal",
			formatter: &Formatter{
				RemoveANSI:          false,
				NormalizeWhitespace: true,
				HighlightErrors:     false,
				MaxLength:           100,
			},
			input:    "\x1b[31mRed text\x1b[0m",
			expected: "\x1b[31mRed text\x1b[0m",
		},
		{
			name:      "empty input",
			formatter: DefaultFormatter(),
			input:     "",
			expected:  "",
		},
		{
			name:      "complex formatting",
			formatter: DefaultFormatter(),
			input:     "\x1b[1mBold\x1b[0m\tTabbed\n\n\nMultiple lines\nError: Test error",
			expected:  "Bold    Tabbed\n\nMultiple lines\n❌ Error: Test error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.formatter.Format(tt.input)
			if result != tt.expected {
				t.Errorf("Format() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFormatter_FormatForNotification(t *testing.T) {
	f := &Formatter{}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "multiline to single line",
			input:    "Line 1\nLine 2\nLine 3",
			expected: "Line 1 Line 2 Line 3",
		},
		{
			name:  "long text truncation",
			input: strings.Repeat("Very long text ", 100),
			expected: func() string {
				// Build expected output by formatting with 500 char limit
				longText := strings.Repeat("Very long text ", 100)
				f := &Formatter{
					MaxLength:           500,
					RemoveANSI:          true,
					NormalizeWhitespace: true,
					HighlightErrors:     false,
				}
				formatted := f.Format(longText)
				// Replace newlines with spaces and collapse multiple spaces
				result := strings.ReplaceAll(formatted, "\n", " ")
				result = strings.TrimSpace(result)
				// Collapse multiple spaces
				for strings.Contains(result, "  ") {
					result = strings.ReplaceAll(result, "  ", " ")
				}
				return result
			}(),
		},
		{
			name:     "ANSI removal",
			input:    "\x1b[32mColored\x1b[0m output",
			expected: "Colored output",
		},
		{
			name:     "whitespace normalization",
			input:    "  Text  \n\n  with   spaces  ",
			expected: "Text with spaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := f.FormatForNotification(tt.input)
			if result != tt.expected {
				t.Errorf("FormatForNotification() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFormatter_FormatError(t *testing.T) {
	f := DefaultFormatter()

	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: "",
		},
		{
			name:     "simple error",
			err:      errors.New("something went wrong"),
			expected: "❌ Error: something went wrong",
		},
		{
			name:     "error with Error prefix",
			err:      errors.New("Error: already has prefix"),
			expected: "❌ Error: already has prefix",
		},
		{
			name:     "error with error lowercase",
			err:      errors.New("error connecting to server"),
			expected: "❌ error connecting to server",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := f.FormatError(tt.err)
			if result != tt.expected {
				t.Errorf("FormatError() = %q, want %q", result, tt.expected)
			}
		})
	}
}
