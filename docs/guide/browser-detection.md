# Browser Detection

SilentCast can automatically detect installed browsers on your system and open URLs in your preferred browser.

## Overview

The browser detection feature:
- Automatically detects installed browsers on Windows, macOS, and Linux
- Identifies the system default browser
- Allows opening URLs in specific browsers
- Supports browser preference lists
- Provides fallback mechanisms

## Configuration

### Basic URL Action

Open a URL in the default browser:

```yaml
grimoire:
  open_github:
    type: url
    url: "https://github.com"
    description: "Open GitHub"
```

### Specific Browser

Open a URL in a specific browser:

```yaml
grimoire:
  open_in_chrome:
    type: url
    url: "https://www.google.com"
    browser: "chrome"
    description: "Open in Chrome"
```

### Browser Preference List

Try browsers in order of preference:

```yaml
grimoire:
  open_with_fallback:
    type: url
    url: "https://example.com"
    browser_preference:
      - "firefox"    # Try Firefox first
      - "chrome"     # Fall back to Chrome
      - "edge"       # Then Edge
    description: "Open with browser preferences"
```

## Supported Browsers

SilentCast detects the following browsers:

### Windows
- Google Chrome
- Mozilla Firefox
- Microsoft Edge
- Opera
- Brave Browser
- Vivaldi

### macOS
- Safari
- Google Chrome
- Mozilla Firefox
- Microsoft Edge
- Opera
- Brave Browser
- Vivaldi
- Arc

### Linux
- Google Chrome / Chromium
- Mozilla Firefox
- Microsoft Edge
- Opera
- Brave Browser
- Vivaldi

## Browser Names

Use these names in your configuration:

| Browser | Config Name |
|---------|------------|
| Google Chrome | `chrome` |
| Mozilla Firefox | `firefox` |
| Microsoft Edge | `edge` |
| Safari | `safari` |
| Opera | `opera` |
| Brave | `brave` |
| Vivaldi | `vivaldi` |
| Chromium | `chromium` |

## Advanced Options

### Browser Arguments

Pass command-line arguments to browsers:

```yaml
grimoire:
  open_incognito:
    type: url
    url: "https://example.com"
    browser: "chrome"
    browser_args: ["--incognito"]
    
  open_new_window:
    type: url
    url: "https://example.com"
    browser: "firefox"
    browser_args: ["--new-window"]
```

### Platform-Specific Configuration

Use platform-specific configuration files for different browser paths:

```yaml
# spellbook.windows.yml
grimoire:
  open_edge:
    type: url
    url: "https://microsoft.com"
    browser_path: "C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe"
```

## Detection Process

### Windows
1. Checks Windows Registry for installed browsers
2. Looks in common installation directories
3. Identifies default browser from registry

### macOS
1. Scans Applications folders
2. Reads Info.plist files for browser information
3. Uses Launch Services to find default browser

### Linux
1. Parses desktop entry files (.desktop)
2. Checks common binary locations
3. Uses xdg-settings for default browser
4. Falls back to BROWSER environment variable

## Troubleshooting

### Browser Not Detected

If a browser isn't detected:

1. **Check Installation Path**: Ensure the browser is installed in a standard location
2. **Use Full Path**: Specify the full path to the browser executable:
   ```yaml
   browser_path: "/opt/custom/browser/bin/browser"
   ```

### Wrong Default Browser

If the wrong browser is detected as default:

1. **Windows**: Check Default Apps in Settings
2. **macOS**: Check General > Default web browser in System Preferences
3. **Linux**: Use `xdg-settings set default-web-browser browser.desktop`

### Browser Fails to Launch

If a browser fails to launch:

1. **Check Permissions**: Ensure SilentCast has permission to launch applications
2. **Test Manually**: Try launching the browser from command line
3. **Check Logs**: Look for error messages in SilentCast logs

## Examples

### Open Different Sites in Different Browsers

```yaml
spells:
  g,w: work_browser    # Work sites in Chrome
  g,p: personal_browser # Personal sites in Firefox

grimoire:
  work_browser:
    type: url
    url: "https://workspace.google.com"
    browser: "chrome"
    description: "Open work apps in Chrome"
    
  personal_browser:
    type: url
    url: "https://reddit.com"
    browser: "firefox"
    description: "Open personal sites in Firefox"
```

### Development Workflow

```yaml
grimoire:
  open_localhost:
    type: url
    url: "http://localhost:3000"
    browser_preference:
      - "chrome"      # Prefer Chrome for dev tools
      - "firefox"     # Firefox as backup
    browser_args: ["--disable-web-security"]  # For development
    description: "Open local development server"
```

### Testing Across Browsers

```yaml
spells:
  t,c: test_chrome
  t,f: test_firefox
  t,e: test_edge
  t,a: test_all

grimoire:
  test_chrome:
    type: url
    url: "http://localhost:8080/test"
    browser: "chrome"
    
  test_firefox:
    type: url
    url: "http://localhost:8080/test"
    browser: "firefox"
    
  test_edge:
    type: url
    url: "http://localhost:8080/test"
    browser: "edge"
    
  test_all:
    type: script
    command: |
      # Open in all browsers for testing
      xdg-open "http://localhost:8080/test" &
      firefox "http://localhost:8080/test" &
      google-chrome "http://localhost:8080/test" &
    description: "Open in all browsers"
```