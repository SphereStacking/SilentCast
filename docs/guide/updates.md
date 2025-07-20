# Updates and Version Management

SilentCast includes built-in functionality to check for updates and notify you when new versions are available. This guide covers how to use the update checking features and configure automatic update checks.

## Manual Update Check

You can manually check for updates at any time using the command line:

```bash
# Check for updates
silentcast --check-update

# Force a fresh check (ignores cache)
silentcast --check-update --force-update-check
```

### Example Output

When a new version is available:
```
🔍 Checking for updates...

🎉 New version available!
  Current version: v1.0.0
  Latest version:  v1.1.0
  Published:       2024-01-15
  Download size:   15.2 MB

📋 Release Notes:
  - Added new spell management features
  - Improved performance
  - Fixed various bugs

💡 To update, visit: https://github.com/SphereStacking/SilentCast/releases
```

When you're up to date:
```
🔍 Checking for updates...
✅ You're running the latest version (v1.1.0)
```

## Update Checking Mechanism

### How It Works

1. **GitHub Releases API**: SilentCast checks the GitHub releases API for the latest version
2. **Version Comparison**: Compares your current version with the latest release
3. **Platform Detection**: Identifies the appropriate binary for your platform
4. **Caching**: Results are cached to avoid excessive API calls

### Cache Behavior

- Update check results are cached for 1 hour by default
- Cache prevents repeated API calls during the cache period
- Use `--force-update-check` to bypass the cache
- Cache files are stored in `~/.config/silentcast/cache/`

## Automatic Update Checks

SilentCast can be configured to automatically check for updates at regular intervals.

### Configuration

Add update settings to your `spellbook.yml`:

```yaml
# Update configuration
updates:
  # Enable automatic update checks
  auto_check: true
  
  # Check interval (duration format: 24h, 7d, etc.)
  check_interval: 24h
  
  # Show notifications for available updates
  notify: true
  
  # Cache duration for update checks
  cache_duration: 1h
```

### Startup Check

When configured, SilentCast will:
1. Wait 1 minute after startup before first check
2. Check for updates based on the configured interval
3. Display a notification if updates are available

## Version Information

### Viewing Version Details

Get detailed version information:

```bash
# Basic version
silentcast --version

# Detailed version info (JSON format)
silentcast --version --version-format json

# Compact version
silentcast --version --version-format compact
```

### Version Formats

**Human-readable (default):**
```
SilentCast v1.0.0
Built: 2024-01-15T10:30:00Z
Commit: abc123def456
Go: 1.21.0
Platform: linux/amd64
```

**JSON format:**
```json
{
  "version": "v1.0.0",
  "git_commit": "abc123def456",
  "build_time": "2024-01-15T10:30:00Z",
  "go_version": "1.21.0",
  "platform": "linux/amd64",
  "cgo_enabled": true
}
```

**Compact format:**
```
v1.0.0 (abc123d)
```

## Update Process

### Manual Update Steps

1. Check for updates: `silentcast --check-update`
2. Visit the releases page: https://github.com/SphereStacking/SilentCast/releases
3. Download the appropriate binary for your platform
4. Replace the existing binary
5. Verify the update: `silentcast --version`

### Platform-Specific Binaries

SilentCast automatically detects the correct binary for your platform:
- `silentcast-windows-amd64.exe` - Windows 64-bit
- `silentcast-darwin-amd64` - macOS Intel
- `silentcast-darwin-arm64` - macOS Apple Silicon
- `silentcast-linux-amd64` - Linux 64-bit
- `silentcast-linux-arm64` - Linux ARM64

## Troubleshooting

### Update Check Fails

**Error: "GitHub API returned status: 404 Not Found"**
- The repository may be private or the URL is incorrect
- Check your internet connection

**Error: "no asset found for platform"**
- The release may not include a binary for your platform
- Check the releases page manually

### Cache Issues

Clear the update cache:
```bash
# Remove cache directory
rm -rf ~/.config/silentcast/cache

# Or force a fresh check
silentcast --check-update --force-update-check
```

### Rate Limiting

GitHub API has rate limits:
- Anonymous: 60 requests per hour
- Authenticated: 5000 requests per hour

The cache mechanism helps avoid hitting these limits.

## Security Considerations

### Checksum Verification

Future versions will support checksum verification:
- SHA256 checksums for each release
- Automatic verification during updates
- Protection against corrupted downloads

### HTTPS Only

All update checks use HTTPS:
- Secure communication with GitHub
- No downgrade to HTTP
- Certificate validation

## Configuration Reference

### Update Settings

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| `updates.auto_check` | bool | false | Enable automatic update checks |
| `updates.check_interval` | duration | 24h | How often to check for updates |
| `updates.notify` | bool | true | Show notifications for updates |
| `updates.cache_duration` | duration | 1h | How long to cache update results |

### Environment Variables

| Variable | Description |
|----------|-------------|
| `SILENTCAST_NO_UPDATE_CHECK` | Disable all update checks |
| `SILENTCAST_UPDATE_CACHE_DIR` | Custom cache directory |

## Future Features

The following features are planned for future releases:

1. **Self-Update Mechanism** (T055)
   - Automatic binary replacement
   - Rollback on failure
   - Progress indication

2. **Update Notifications** (T056)
   - System tray notifications
   - Configurable notification preferences
   - Release notes preview

3. **Beta Channel Support**
   - Opt-in to pre-release versions
   - Separate update channels
   - Channel switching

## Related Documentation

- [Installation Guide](installation.md)
- [Configuration Reference](../config/)
- [Version History](../CHANGELOG.md)
- [Contributing Guide](../contributing.md)