//go:build windows

package browser

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/sys/windows/registry"
)

type windowsDetector struct{}

func newPlatformDetector() Detector {
	return &windowsDetector{}
}

func (d *windowsDetector) DetectBrowsers(ctx context.Context) ([]Browser, error) {
	var browsers []Browser

	// Check registry for installed browsers
	browsers = append(browsers, d.detectFromRegistry()...)

	// Check common installation paths
	browsers = append(browsers, d.detectFromPaths()...)

	// Remove duplicates
	browsers = d.deduplicateBrowsers(browsers)

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

func (d *windowsDetector) GetDefaultBrowser(ctx context.Context) (*Browser, error) {
	// Try to get default browser from registry
	key, err := registry.OpenKey(registry.CURRENT_USER,
		`Software\Microsoft\Windows\Shell\Associations\UrlAssociations\http\UserChoice`,
		registry.QUERY_VALUE)
	if err != nil {
		return nil, fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	progID, _, err := key.GetStringValue("ProgId")
	if err != nil {
		return nil, fmt.Errorf("failed to get ProgId: %w", err)
	}

	// Get command from ProgID
	cmdKey, err := registry.OpenKey(registry.CLASSES_ROOT,
		fmt.Sprintf(`%s\shell\open\command`, progID),
		registry.QUERY_VALUE)
	if err != nil {
		return nil, fmt.Errorf("failed to open command key: %w", err)
	}
	defer cmdKey.Close()

	command, _, err := cmdKey.GetStringValue("")
	if err != nil {
		return nil, fmt.Errorf("failed to get command: %w", err)
	}

	// Extract executable from command
	executable := extractExecutable(command)
	if executable == "" {
		return nil, fmt.Errorf("failed to extract executable from command")
	}

	// Get browser name from ProgID
	name := d.getBrowserNameFromProgID(progID)

	return &Browser{
		Name:       name,
		Executable: executable,
		IsDefault:  true,
	}, nil
}

func (d *windowsDetector) FindBrowser(ctx context.Context, name string) (*Browser, error) {
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

func (d *windowsDetector) detectFromRegistry() []Browser {
	var browsers []Browser

	// Check both HKLM and HKCU
	roots := []registry.Key{registry.LOCAL_MACHINE, registry.CURRENT_USER}

	for _, root := range roots {
		// Check StartMenuInternet entries
		key, err := registry.OpenKey(root,
			`SOFTWARE\Clients\StartMenuInternet`,
			registry.ENUMERATE_SUB_KEYS)
		if err != nil {
			continue
		}
		defer key.Close()

		subkeys, err := key.ReadSubKeyNames(-1)
		if err != nil {
			continue
		}

		for _, subkey := range subkeys {
			browser := d.getBrowserFromStartMenu(root, subkey)
			if browser != nil {
				browsers = append(browsers, *browser)
			}
		}
	}

	return browsers
}

func (d *windowsDetector) getBrowserFromStartMenu(root registry.Key, browserKey string) *Browser {
	key, err := registry.OpenKey(root,
		fmt.Sprintf(`SOFTWARE\Clients\StartMenuInternet\%s`, browserKey),
		registry.QUERY_VALUE)
	if err != nil {
		return nil
	}
	defer key.Close()

	// Get display name
	name, _, _ := key.GetStringValue("")
	if name == "" {
		name = browserKey
	}

	// Get command
	cmdKey, err := registry.OpenKey(root,
		fmt.Sprintf(`SOFTWARE\Clients\StartMenuInternet\%s\shell\open\command`, browserKey),
		registry.QUERY_VALUE)
	if err != nil {
		return nil
	}
	defer cmdKey.Close()

	command, _, err := cmdKey.GetStringValue("")
	if err != nil {
		return nil
	}

	executable := extractExecutable(command)
	if executable == "" {
		return nil
	}

	return &Browser{
		Name:       name,
		Executable: executable,
	}
}

func (d *windowsDetector) detectFromPaths() []Browser {
	var browsers []Browser

	// Common browser paths
	paths := []struct {
		name  string
		paths []string
	}{
		{
			name: "Google Chrome",
			paths: []string{
				`C:\Program Files\Google\Chrome\Application\chrome.exe`,
				`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
				filepath.Join(os.Getenv("LOCALAPPDATA"), `Google\Chrome\Application\chrome.exe`),
			},
		},
		{
			name: "Mozilla Firefox",
			paths: []string{
				`C:\Program Files\Mozilla Firefox\firefox.exe`,
				`C:\Program Files (x86)\Mozilla Firefox\firefox.exe`,
			},
		},
		{
			name: "Microsoft Edge",
			paths: []string{
				`C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`,
				`C:\Program Files\Microsoft\Edge\Application\msedge.exe`,
			},
		},
		{
			name: "Opera",
			paths: []string{
				`C:\Program Files\Opera\launcher.exe`,
				`C:\Program Files (x86)\Opera\launcher.exe`,
				filepath.Join(os.Getenv("LOCALAPPDATA"), `Programs\Opera\launcher.exe`),
			},
		},
		{
			name: "Brave",
			paths: []string{
				`C:\Program Files\BraveSoftware\Brave-Browser\Application\brave.exe`,
				`C:\Program Files (x86)\BraveSoftware\Brave-Browser\Application\brave.exe`,
			},
		},
	}

	for _, browser := range paths {
		for _, path := range browser.paths {
			if _, err := os.Stat(path); err == nil {
				browsers = append(browsers, Browser{
					Name:       browser.name,
					Executable: path,
				})
				break
			}
		}
	}

	return browsers
}

func (d *windowsDetector) deduplicateBrowsers(browsers []Browser) []Browser {
	seen := make(map[string]bool)
	var result []Browser

	for _, browser := range browsers {
		if !seen[browser.Executable] {
			seen[browser.Executable] = true
			result = append(result, browser)
		}
	}

	return result
}

func (d *windowsDetector) getBrowserNameFromProgID(progID string) string {
	progID = strings.ToLower(progID)

	switch {
	case strings.Contains(progID, "chrome"):
		return "Google Chrome"
	case strings.Contains(progID, "firefox"):
		return "Mozilla Firefox"
	case strings.Contains(progID, "edge"):
		return "Microsoft Edge"
	case strings.Contains(progID, "opera"):
		return "Opera"
	case strings.Contains(progID, "brave"):
		return "Brave"
	default:
		return progID
	}
}

func extractExecutable(command string) string {
	// Remove quotes and extract executable path
	command = strings.TrimSpace(command)

	if strings.HasPrefix(command, `"`) {
		// Quoted path
		end := strings.Index(command[1:], `"`)
		if end > 0 {
			return command[1 : end+1]
		}
	} else {
		// Unquoted path
		parts := strings.Fields(command)
		if len(parts) > 0 {
			return parts[0]
		}
	}

	return ""
}

func openURLFallback(ctx context.Context, url string) error {
	// Use Windows shell to open URL
	cmd := exec.CommandContext(ctx, "cmd", "/c", "start", url)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Start()
}
