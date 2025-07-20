package output

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Formatter provides methods to format output for display
type Formatter struct {
	// MaxLength is the maximum length of formatted output
	MaxLength int
	// RemoveANSI determines if ANSI escape codes should be removed
	RemoveANSI bool
	// NormalizeWhitespace determines if whitespace should be normalized
	NormalizeWhitespace bool
	// HighlightErrors determines if error patterns should be highlighted
	HighlightErrors bool
}

// DefaultFormatter returns a formatter with sensible defaults
func DefaultFormatter() *Formatter {
	return &Formatter{
		MaxLength:           2000,
		RemoveANSI:          true,
		NormalizeWhitespace: true,
		HighlightErrors:     true,
	}
}

// Format applies all formatting rules to the input string
func (f *Formatter) Format(input string) string {
	if input == "" {
		return ""
	}

	output := input

	// Remove ANSI escape codes if requested
	if f.RemoveANSI {
		output = f.stripANSI(output)
	}

	// Normalize whitespace if requested
	if f.NormalizeWhitespace {
		output = f.normalizeWhitespace(output)
	}

	// Highlight errors if requested
	if f.HighlightErrors {
		output = f.highlightErrors(output)
	}

	// Truncate if needed
	if f.MaxLength > 0 && utf8.RuneCountInString(output) > f.MaxLength {
		output = f.truncate(output, f.MaxLength)
	}

	return output
}

// stripANSI removes ANSI escape sequences from the string
func (f *Formatter) stripANSI(s string) string {
	// Match ANSI escape sequences
	// ESC [ ... m (SGR - Select Graphic Rendition)
	// ESC [ ... J (ED - Erase Display)
	// ESC [ ... K (EL - Erase Line)
	// ESC [ ... H (CUP - Cursor Position)
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[mJKH]`)
	return ansiRegex.ReplaceAllString(s, "")
}

// normalizeWhitespace normalizes whitespace in the string
func (f *Formatter) normalizeWhitespace(s string) string {
	// Replace tabs with spaces
	s = strings.ReplaceAll(s, "\t", "    ")

	// Replace multiple consecutive newlines with double newline
	s = regexp.MustCompile(`\n{3,}`).ReplaceAllString(s, "\n\n")

	// Trim trailing whitespace from each line
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRightFunc(line, unicode.IsSpace)
	}
	s = strings.Join(lines, "\n")

	// Trim leading and trailing whitespace
	s = strings.TrimSpace(s)

	return s
}

// highlightErrors adds markers for common error patterns
func (f *Formatter) highlightErrors(s string) string {
	// Common error patterns
	errorPatterns := []struct {
		pattern *regexp.Regexp
		prefix  string
	}{
		{regexp.MustCompile(`(?i)^error:`), "❌ "},
		{regexp.MustCompile(`(?i)^error\s`), "❌ "},
		{regexp.MustCompile(`(?i)^fatal:`), "💀 "},
		{regexp.MustCompile(`(?i)^warning:`), "⚠️  "},
		{regexp.MustCompile(`(?i)^failed to`), "❌ "},
		{regexp.MustCompile(`(?i)^cannot\s`), "❌ "},
		{regexp.MustCompile(`(?i)^unable to`), "❌ "},
	}

	lines := strings.Split(s, "\n")
	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		for _, ep := range errorPatterns {
			if ep.pattern.MatchString(trimmedLine) {
				// Only add prefix if not already present
				if !strings.HasPrefix(trimmedLine, ep.prefix) {
					lines[i] = strings.Replace(line, trimmedLine, ep.prefix+trimmedLine, 1)
				}
				break
			}
		}
	}

	return strings.Join(lines, "\n")
}

// truncate truncates the string to the specified length
func (f *Formatter) truncate(s string, maxLength int) string {
	if maxLength <= 0 {
		return s
	}

	// Count runes, not bytes
	runes := []rune(s)
	if len(runes) <= maxLength {
		return s
	}

	// Reserve space for truncation message
	truncationMsg := "\n\n... (output truncated)"
	reservedLength := utf8.RuneCountInString(truncationMsg)

	if maxLength <= reservedLength {
		return truncationMsg
	}

	// Truncate at a reasonable boundary (preferably at a newline)
	targetLength := maxLength - reservedLength
	truncateAt := targetLength

	// Look for a newline within the last 10% of the target length
	searchStart := targetLength - (targetLength / 10)
	if searchStart < 0 {
		searchStart = 0
	}

	for i := targetLength - 1; i >= searchStart; i-- {
		if runes[i] == '\n' {
			truncateAt = i
			break
		}
	}

	return string(runes[:truncateAt]) + truncationMsg
}

// FormatForNotification formats output specifically for system notifications
func (f *Formatter) FormatForNotification(input string) string {
	// Use stricter limits for notifications
	notificationFormatter := &Formatter{
		MaxLength:           500, // Much shorter for notifications
		RemoveANSI:          true,
		NormalizeWhitespace: true,
		HighlightErrors:     false, // Emojis might not work in all notification systems
	}

	output := notificationFormatter.Format(input)

	// Additional formatting for notifications
	// Replace multiple newlines with single space for better display
	output = regexp.MustCompile(`\n+`).ReplaceAllString(output, " ")

	// Replace multiple spaces with single space
	output = regexp.MustCompile(`\s+`).ReplaceAllString(output, " ")

	// Trim to ensure clean display
	output = strings.TrimSpace(output)

	return output
}

// FormatError formats an error message for display
func (f *Formatter) FormatError(err error) string {
	if err == nil {
		return ""
	}

	msg := err.Error()

	// Add error prefix if not already present
	if !strings.HasPrefix(strings.ToLower(msg), "error") {
		msg = "Error: " + msg
	}

	return f.Format(msg)
}
