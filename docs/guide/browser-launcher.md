# Browser Launcher

The browser launcher provides a robust, cross-platform way to open URLs in specific browsers with advanced options.

## Features

- Cross-platform support (Windows, macOS, Linux)
- Automatic browser detection
- Browser preference lists
- Incognito/private browsing mode
- New window/tab control
- URL validation and security
- Graceful fallback mechanisms

## Basic Usage

### URL Action Type

The simplest way to open a URL:

```yaml
grimoire:
  open_website:
    type: url
    url: "https://example.com"
    description: "Open example website"
```

This will open the URL in the system default browser.

### Specific Browser

Open in a specific browser:

```yaml
grimoire:
  open_in_firefox:
    type: url
    url: "https://mozilla.org"
    browser: "firefox"
    description: "Open in Firefox"
```

### Browser Preferences

Try browsers in order of preference:

```yaml
grimoire:
  open_with_preference:
    type: url
    url: "https://example.com"
    browser_preference:
      - "firefox"      # Try first
      - "chrome"       # Try second
      - "edge"         # Try third
    description: "Open with browser preference"
```

## Advanced Options

### Incognito/Private Mode

Open in private browsing mode:

```yaml
grimoire:
  open_private:
    type: url
    url: "https://example.com"
    browser: "chrome"
    incognito: true
    description: "Open in Chrome incognito"
```

Browser-specific private mode flags:
- Chrome/Edge/Brave: `--incognito`
- Firefox: `--private-window`
- Safari: Not supported via command line

### New Window

Force opening in a new window:

```yaml
grimoire:
  open_new_window:
    type: url
    url: "https://example.com"
    new_window: true
    description: "Open in new window"
```

### Custom Browser Arguments

Pass additional arguments to the browser:

```yaml
grimoire:
  open_with_args:
    type: url
    url: "https://example.com"
    browser: "chrome"
    browser_args:
      - "--disable-extensions"
      - "--disable-plugins"
      - "--no-sandbox"
    description: "Open with custom args"
```

## URL Handling

### Automatic HTTPS

URLs without a scheme automatically get HTTPS:

```yaml
# These are equivalent:
url: "example.com"        # Becomes https://example.com
url: "https://example.com"
```

### Local Development

Local URLs automatically get HTTP:

```yaml
# These automatically use HTTP:
url: "localhost:3000"     # Becomes http://localhost:3000
url: "127.0.0.1:8080"     # Becomes http://127.0.0.1:8080
url: "192.168.1.100"      # Becomes http://192.168.1.100
```

### Special URLs

Supported special URL schemes:

```yaml
# Browser-specific pages
url: "about:blank"
url: "chrome://settings"
url: "edge://flags"

# File URLs
url: "file:///home/user/index.html"
```

### Security

Dangerous URL schemes are blocked:

```yaml
# These will be rejected:
url: "javascript:alert(1)"    # XSS risk
url: "data:text/html,..."     # Data URLs
url: "vbscript:..."           # VBScript
```

## Platform-Specific Behavior

### Windows

- Uses registry to detect browsers
- Launches browsers directly via executable path
- Falls back to `cmd /c start` for default browser
- Hides console windows automatically

### macOS

- Detects browsers in Applications folders
- Uses `open -a` for .app bundles
- Falls back to `open` command for default browser
- Supports macOS-specific browser arguments

### Linux

- Reads .desktop files for browser information
- Checks common binary locations
- Uses `xdg-open` for default browser
- Supports various desktop environments

## Examples

### Development Workflow

```yaml
spells:
  d,l: dev_local        # Local development
  d,s: dev_staging      # Staging server
  d,p: dev_production   # Production site

grimoire:
  dev_local:
    type: url
    url: "localhost:3000"
    browser: "chrome"
    browser_args: ["--disable-web-security"]
    description: "Local dev server"
    
  dev_staging:
    type: url
    url: "https://staging.example.com"
    browser: "firefox"
    incognito: true
    description: "Staging server (private)"
    
  dev_production:
    type: url
    url: "https://example.com"
    new_window: true
    description: "Production site"
```

### Multi-Browser Testing

```yaml
grimoire:
  test_all_browsers:
    type: script
    command: |
      # Open same URL in multiple browsers for testing
      echo "Opening test page in all browsers..."
    description: "Cross-browser testing"
    
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
```

### Quick Access Bookmarks

```yaml
spells:
  g,h: github
  g,m: gmail
  g,c: calendar
  g,d: drive

grimoire:
  github:
    type: url
    url: "https://github.com"
    browser_preference: ["chrome", "firefox"]
    
  gmail:
    type: url
    url: "https://mail.google.com"
    browser: "chrome"  # Use Chrome for Google services
    
  calendar:
    type: url
    url: "https://calendar.google.com"
    browser: "chrome"
    
  drive:
    type: url
    url: "https://drive.google.com"
    browser: "chrome"
```

## Error Handling

The launcher handles various error scenarios:

1. **Browser Not Found**: Falls back to browser preference list or default browser
2. **Invalid URL**: Returns clear error message
3. **Launch Failure**: Attempts fallback methods
4. **No Browser Available**: Uses system URL opener

## Best Practices

1. **Use Browser Preferences**: Instead of hardcoding a single browser, provide a preference list
2. **Validate URLs**: The launcher validates and normalizes URLs automatically
3. **Handle Local Development**: Use appropriate schemes for localhost
4. **Consider Privacy**: Use incognito mode for sensitive sites
5. **Test Cross-Platform**: Test your configuration on different operating systems

## Troubleshooting

### Browser Not Launching

1. Check if browser is installed and in PATH
2. Verify browser name matches expected format
3. Check system logs for launch errors
4. Try using full browser path

### URL Not Opening

1. Verify URL is valid and accessible
2. Check for typos in URL
3. Ensure scheme is supported
4. Test URL in browser manually

### Wrong Browser Opens

1. Verify browser preference order
2. Check system default browser settings
3. Ensure specified browser is detected
4. Use explicit browser path if needed