//go:build linux

package browser

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type linuxDetector struct{}

func newPlatformDetector() Detector {
	return &linuxDetector{}
}

func (d *linuxDetector) DetectBrowsers(ctx context.Context) ([]Browser, error) {
	var browsers []Browser

	// Check desktop files
	browsers = append(browsers, d.detectFromDesktopFiles()...)

	// Check common binary locations
	browsers = append(browsers, d.detectFromPaths()...)

	// Remove duplicates
	browsers = d.deduplicateBrowsers(browsers)

	// Get default browser
	defaultBrowser, _ := d.GetDefaultBrowser(ctx) // Error handled by checking if defaultBrowser is nil
	if defaultBrowser != nil {
		for i := range browsers {
			if browsers[i].Executable == defaultBrowser.Executable {
				browsers[i].IsDefault = true
			}
		}
	}

	return browsers, nil
}

func (d *linuxDetector) GetDefaultBrowser(ctx context.Context) (*Browser, error) {
	// Try xdg-settings first
	cmd := exec.CommandContext(ctx, "xdg-settings", "get", "default-web-browser")
	output, err := cmd.Output()
	if err == nil {
		desktopFile := strings.TrimSpace(string(output))
		if desktopFile != "" {
			browser := d.getBrowserFromDesktopFile(desktopFile)
			if browser != nil {
				browser.IsDefault = true
				return browser, nil
			}
		}
	}

	// Try alternatives system
	cmd = exec.CommandContext(ctx, "update-alternatives", "--query", "x-www-browser")
	output, err = cmd.Output()
	if err == nil {
		scanner := bufio.NewScanner(strings.NewReader(string(output)))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "Value: ") {
				executable := strings.TrimPrefix(line, "Value: ")
				name := d.getBrowserNameFromExecutable(executable)
				return &Browser{
					Name:       name,
					Executable: executable,
					IsDefault:  true,
				}, nil
			}
		}
	}

	// Try BROWSER environment variable
	if browserEnv := os.Getenv("BROWSER"); browserEnv != "" {
		// BROWSER can be a colon-separated list
		browserList := strings.Split(browserEnv, ":")
		if len(browserList) > 0 && browserList[0] != "" {
			executable := browserList[0]
			if _, err := exec.LookPath(executable); err == nil {
				name := d.getBrowserNameFromExecutable(executable)
				return &Browser{
					Name:       name,
					Executable: executable,
					IsDefault:  true,
				}, nil
			}
		}
	}

	return nil, ErrNoBrowserFound
}

func (d *linuxDetector) FindBrowser(ctx context.Context, name string) (*Browser, error) {
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

func (d *linuxDetector) detectFromDesktopFiles() []Browser {
	var browsers []Browser

	// Desktop file locations
	desktopDirs := []string{
		"/usr/share/applications",
		"/usr/local/share/applications",
		filepath.Join(os.Getenv("HOME"), ".local", "share", "applications"),
		"/var/lib/snapd/desktop/applications",
		"/var/lib/flatpak/exports/share/applications",
		filepath.Join(os.Getenv("HOME"), ".local", "share", "flatpak", "exports", "share", "applications"),
	}

	// Browser-related desktop files
	browserPatterns := []string{
		"*chrome*.desktop",
		"*chromium*.desktop",
		"*firefox*.desktop",
		"*brave*.desktop",
		"*edge*.desktop",
		"*opera*.desktop",
		"*vivaldi*.desktop",
		"*web-browser*.desktop",
		"*webbrowser*.desktop",
	}

	for _, dir := range desktopDirs {
		if _, err := os.Stat(dir); err != nil {
			continue
		}

		for _, pattern := range browserPatterns {
			matches, err := filepath.Glob(filepath.Join(dir, pattern))
			if err != nil {
				continue
			}

			for _, desktopFile := range matches {
				browser := d.parseDesktopFile(desktopFile)
				if browser != nil {
					browsers = append(browsers, *browser)
				}
			}
		}
	}

	return browsers
}

func (d *linuxDetector) parseDesktopFile(path string) *Browser {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer func() { _ = file.Close() }()

	var name, execLine, categories string
	scanner := bufio.NewScanner(file)
	inDesktopEntry := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "[Desktop Entry]" {
			inDesktopEntry = true
			continue
		}

		if strings.HasPrefix(line, "[") && line != "[Desktop Entry]" {
			inDesktopEntry = false
			continue
		}

		if !inDesktopEntry {
			continue
		}

		switch {
		case strings.HasPrefix(line, "Name=") && name == "":
			name = strings.TrimPrefix(line, "Name=")
		case strings.HasPrefix(line, "Exec=") && execLine == "":
			execLine = strings.TrimPrefix(line, "Exec=")
		case strings.HasPrefix(line, "Categories="):
			categories = strings.TrimPrefix(line, "Categories=")
		}
	}

	// Check if it's a web browser
	if !strings.Contains(categories, "WebBrowser") &&
		!strings.Contains(categories, "Network") &&
		!strings.Contains(name, "Browser") &&
		!strings.Contains(name, "Web") {
		return nil
	}

	if name == "" || execLine == "" {
		return nil
	}

	// Extract executable from Exec line
	executable := d.extractExecutableFromExec(execLine)
	if executable == "" {
		return nil
	}

	// Resolve full path if needed
	if !filepath.IsAbs(executable) {
		if fullPath, err := exec.LookPath(executable); err == nil {
			executable = fullPath
		}
	}

	return &Browser{
		Name:       name,
		Executable: executable,
	}
}

func (d *linuxDetector) extractExecutableFromExec(execLine string) string {
	// Remove field codes (%f, %F, %u, %U, etc.)
	execLine = strings.ReplaceAll(execLine, "%f", "")
	execLine = strings.ReplaceAll(execLine, "%F", "")
	execLine = strings.ReplaceAll(execLine, "%u", "")
	execLine = strings.ReplaceAll(execLine, "%U", "")

	// Split and get first part
	parts := strings.Fields(execLine)
	if len(parts) == 0 {
		return ""
	}

	executable := parts[0]

	// Remove quotes if present
	executable = strings.Trim(executable, `"'`)

	return executable
}

func (d *linuxDetector) detectFromPaths() []Browser {
	var browsers []Browser

	// Common browser executables
	browserExecs := []struct {
		name  string
		execs []string
	}{
		{
			name:  "Google Chrome",
			execs: []string{"google-chrome", "google-chrome-stable", "chrome"},
		},
		{
			name:  "Chromium",
			execs: []string{"chromium", "chromium-browser"},
		},
		{
			name:  "Mozilla Firefox",
			execs: []string{"firefox", "firefox-esr"},
		},
		{
			name:  "Microsoft Edge",
			execs: []string{"microsoft-edge", "microsoft-edge-stable", "edge"},
		},
		{
			name:  "Opera",
			execs: []string{"opera"},
		},
		{
			name:  "Brave Browser",
			execs: []string{"brave-browser", "brave"},
		},
		{
			name:  "Vivaldi",
			execs: []string{"vivaldi", "vivaldi-stable"},
		},
	}

	for _, browser := range browserExecs {
		for _, execName := range browser.execs {
			if fullPath, err := exec.LookPath(execName); err == nil {
				browsers = append(browsers, Browser{
					Name:       browser.name,
					Executable: fullPath,
				})
				break // Found this browser, move to next
			}
		}
	}

	return browsers
}

func (d *linuxDetector) deduplicateBrowsers(browsers []Browser) []Browser {
	seen := make(map[string]bool)
	var result []Browser

	for _, browser := range browsers {
		key := browser.Executable
		// Normalize symlinks
		if resolved, err := filepath.EvalSymlinks(browser.Executable); err == nil {
			key = resolved
		}

		if !seen[key] {
			seen[key] = true
			result = append(result, browser)
		}
	}

	return result
}

func (d *linuxDetector) getBrowserFromDesktopFile(desktopFile string) *Browser {
	// Find desktop file in standard locations
	desktopDirs := []string{
		"/usr/share/applications",
		"/usr/local/share/applications",
		filepath.Join(os.Getenv("HOME"), ".local", "share", "applications"),
	}

	for _, dir := range desktopDirs {
		fullPath := filepath.Join(dir, desktopFile)
		if browser := d.parseDesktopFile(fullPath); browser != nil {
			return browser
		}
	}

	return nil
}

func (d *linuxDetector) getBrowserNameFromExecutable(executable string) string {
	base := filepath.Base(executable)
	base = strings.ToLower(base)

	switch {
	case strings.Contains(base, "chrome"):
		if strings.Contains(base, "chromium") {
			return "Chromium"
		}
		return "Google Chrome"
	case strings.Contains(base, "firefox"):
		return "Mozilla Firefox"
	case strings.Contains(base, "edge"):
		return "Microsoft Edge"
	case strings.Contains(base, "opera"):
		return "Opera"
	case strings.Contains(base, "brave"):
		return "Brave Browser"
	case strings.Contains(base, "vivaldi"):
		return "Vivaldi"
	default:
		return base
	}
}

func openURLFallback(ctx context.Context, url string) error {
	// Try xdg-open first
	if _, err := exec.LookPath("xdg-open"); err == nil {
		cmd := exec.CommandContext(ctx, "xdg-open", url)
		return cmd.Start()
	}

	// Try common desktop environment specific openers
	openers := []string{"gnome-open", "kde-open", "exo-open"}
	for _, opener := range openers {
		if _, err := exec.LookPath(opener); err == nil {
			cmd := exec.CommandContext(ctx, opener, url)
			return cmd.Start()
		}
	}

	return fmt.Errorf("no URL opener found")
}
