# Update Notifications

SilentCast includes a comprehensive update notification system that keeps you informed about new releases and provides convenient ways to update your installation.

## Features

- **Automatic Background Checks**: Periodic checks for new versions without interrupting your workflow
- **Multiple Notification Channels**: Console, desktop notifications, and system tray integration
- **Rich Notification Content**: Version information, release notes, download size, and publication date
- **Interactive Actions**: Update now, view release notes, remind later, or dismiss notifications
- **Configurable Preferences**: Control when and how you receive update notifications
- **Smart Caching**: Prevents excessive API calls while staying current

## Commands

### Check Update Status

View detailed information about your current version and available updates:

```bash
./silentcast --update-status
```

With verbose information:

```bash
./silentcast --update-status --debug
```

Force a fresh check (ignore cache):

```bash
./silentcast --update-status --force
```

### Quick Update Check

Simple check for updates:

```bash
./silentcast --check-update
```

Force fresh check:

```bash
./silentcast --check-update --force-update-check
```

### Self-Update

Update to the latest version:

```bash
./silentcast --self-update
```

## Notification Types

### 1. Update Available Notifications

When a new version is detected, you'll receive a detailed notification including:

- Current and new version numbers
- Publication date
- Download size
- Release notes summary
- Available actions

**Console Example:**
```
[14:30:15] 🚀 SilentCast Update Available: Version v0.2.0 is now available (8.5 MB)
[14:30:15] Update Details:
  Current: v0.1.0-alpha.8
  Latest:  v0.2.0
  Published: 2024-01-20
  Size: 8.5 MB
[14:30:15] Release Notes:
        ✨ Add update notification system
        🐛 Fix hotkey detection on Linux
        📝 Improve documentation
[14:30:15] Available actions: update, view, remind, dismiss
  Run: ./silentcast --self-update
```

### 2. Update Status Notifications

Information about update checks and current status:

- Update check failures (network issues, API errors)
- "No updates available" confirmations
- Cache status and next check time

### 3. Update Progress Notifications

During the update process:

- Download started notifications
- Update completion confirmations
- Update failure alerts with error details

## Configuration

### Default Settings

The update notification system comes with sensible defaults:

```yaml
update_notifications:
  enabled: true
  check_interval: 24h          # Check daily
  show_on_startup: true        # Check on app start
  remind_interval: 168h        # Remind weekly
  auto_check: true             # Background checking
  include_prereleases: false   # Stable releases only
```

### Notification Preferences

You can configure how and when you receive notifications:

#### Check Frequency
- **Daily** (default): Check every 24 hours
- **Weekly**: Check every 7 days
- **Manual**: Only check when explicitly requested

#### Startup Behavior
- **Check on startup** (default): Brief check when app starts
- **Startup disabled**: No automatic checks on startup

#### Reminder Settings
- **Weekly reminders** (default): Remind about dismissed updates
- **Custom interval**: Set your own reminder frequency
- **No reminders**: Dismiss notifications permanently

### Disabling Notifications

To disable update notifications entirely:

```bash
# Disable via configuration
echo "update_notifications.enabled: false" >> spellbook.yml

# Or use manual-only mode
./silentcast --check-update  # Only when explicitly requested
```

## Notification Channels

### Console Notifications

Always available, shows detailed information with color coding:

- 🚀 **Update Available**: Blue/cyan for new versions
- ⚠️ **Check Failed**: Yellow for network issues
- ✅ **Success**: Green for completions
- ❌ **Error**: Red for failures

### Desktop Notifications

Platform-specific system notifications:

- **Windows**: Toast notifications with action buttons
- **macOS**: Notification Center integration
- **Linux**: Desktop notification daemon support

### System Tray Integration

When running with system tray:

- Update indicator in tray icon
- Context menu with update actions
- Quick access to update commands

## Update Actions

### Available Actions

1. **Update Now** (`update`)
   - Downloads and installs the update immediately
   - Shows progress during download
   - Automatically restarts the application

2. **View Release Notes** (`view`)
   - Displays full release notes and changelog
   - Shows what's new in the update
   - Helps decide whether to update

3. **Remind Later** (`remind`)
   - Snoozes the notification for the configured interval
   - Default: 7 days
   - Configurable reminder frequency

4. **Dismiss** (`dismiss`)
   - Dismisses the current update notification
   - Won't remind about this specific version
   - Will notify about future versions

### Action Examples

**Console Actions:**
```bash
# View current status
./silentcast --update-status

# Start update process
./silentcast --self-update

# Force update without confirmation
./silentcast --self-update --force-self-update
```

**Programmatic Actions:**
```go
// Handle notification actions
manager.HandleUpdateAction(ctx, notify.UpdateActionUpdate, updateInfo, updater)
manager.HandleUpdateAction(ctx, notify.UpdateActionView, updateInfo, nil)
manager.HandleUpdateAction(ctx, notify.UpdateActionRemind, updateInfo, nil)
manager.HandleUpdateAction(ctx, notify.UpdateActionDismiss, updateInfo, nil)
```

## Background Checking

### Automatic Checks

When enabled, SilentCast performs background update checks:

1. **Startup Check**: Brief check 30 seconds after startup
2. **Periodic Checks**: Based on configured interval (default: 24h)
3. **Smart Caching**: Respects cache to avoid excessive API calls
4. **Non-Intrusive**: Only shows notifications when updates are available

### Cache Behavior

The notification system includes intelligent caching:

- **Cache Duration**: 1 hour for GitHub API responses
- **Check Interval**: 24 hours between automatic checks
- **Force Refresh**: Available via `--force-update-check`
- **Cache Location**: User configuration directory

### Network Handling

Robust handling of network conditions:

- **Timeout Protection**: 30-second timeout for API calls
- **Error Recovery**: Graceful handling of network failures
- **Rate Limiting**: Respects GitHub API rate limits
- **Offline Mode**: No disruptive errors when offline

## Integration with Self-Update

The notification system works seamlessly with the self-update mechanism:

### Workflow Integration

1. **Detection**: Background process detects new version
2. **Notification**: User receives update notification
3. **Action**: User chooses to update via notification
4. **Download**: Progress-tracked download with verification
5. **Installation**: Atomic installation with rollback protection
6. **Completion**: Success notification and automatic restart

### Progress Tracking

During updates initiated from notifications:

```
⬇️ Downloading update...
[████████████████████████████████████████] 100.0% 8.5 MB/8.5 MB @ 2.1 MB/s
🔐 Verifying checksum...
✅ Update applied successfully!
🔄 Restarting SilentCast...
```

## Error Handling

### Common Scenarios

**Network Issues:**
- Graceful degradation when GitHub is unreachable
- Clear error messages for connectivity problems
- No disruptive notifications for temporary failures

**API Limitations:**
- Respects GitHub API rate limiting
- Intelligent backoff for 403/429 responses
- Fallback to cached information when appropriate

**Update Failures:**
- Detailed error reporting for download failures
- Automatic rollback for installation errors
- User-friendly error messages with suggested actions

### Troubleshooting

**No notifications appearing:**
```bash
# Check notification system status
./silentcast --update-status --debug

# Test notification system
./silentcast --check-update --force-update-check

# Verify configuration
./silentcast --show-config | grep -A5 update_notifications
```

**Notifications not working:**
```bash
# Check available notifiers
./silentcast --validate-config

# Test with console output
./silentcast --no-tray --debug --check-update
```

## Security Considerations

### Secure Updates

- **HTTPS Only**: All update checks use encrypted connections
- **Checksum Verification**: SHA256 verification of downloaded files
- **Signed Releases**: Future support for GPG signature verification
- **Atomic Installation**: Safe replacement with automatic rollback

### Privacy

- **Minimal Data**: Only version checks, no personal information
- **GitHub API**: Standard GitHub release API usage
- **No Tracking**: No analytics or usage tracking
- **Local Cache**: All caching is local to your machine

## Development Integration

### Custom Notification Handlers

```go
// Create custom update notification manager
notifier := notify.NewManager()
updateManager := notify.NewUpdateNotificationManager(notifier)

// Configure notification preferences
config := notify.UpdateNotificationConfig{
    Enabled:            true,
    CheckInterval:      12 * time.Hour,
    ShowOnStartup:      false,
    RemindInterval:     24 * time.Hour,
    AutoCheck:          true,
    IncludePreReleases: false,
}
updateManager.SetConfig(config)

// Start background checking
updateManager.StartPeriodicChecks(ctx, updater, currentVersion)
```

### Custom Notifiers

Implement the `UpdateNotifier` interface for custom notification handling:

```go
type CustomUpdateNotifier struct{}

func (n *CustomUpdateNotifier) ShowUpdateNotification(ctx context.Context, notification notify.UpdateNotification) error {
    // Custom notification display logic
    return nil
}

func (n *CustomUpdateNotifier) SupportsUpdateActions() bool {
    return true
}

func (n *CustomUpdateNotifier) OnUpdateAction(action notify.UpdateAction, updateInfo notify.UpdateNotification) error {
    // Handle user actions
    return nil
}
```

## Best Practices

### For Users

1. **Keep Notifications Enabled**: Stay informed about security updates
2. **Update Regularly**: Don't postpone important updates indefinitely
3. **Read Release Notes**: Understand what's changing before updating
4. **Backup Configuration**: Export your spellbook before major updates

### For Developers

1. **Respect User Preferences**: Honor notification settings
2. **Provide Context**: Include meaningful release notes
3. **Test Thoroughly**: Verify notifications on all platforms
4. **Handle Failures Gracefully**: Provide helpful error messages

### For System Administrators

1. **Configure Appropriately**: Set suitable check intervals for your environment
2. **Monitor Updates**: Keep track of available updates across installations
3. **Test Updates**: Validate updates in staging environments
4. **Automate When Appropriate**: Consider automated updates for security patches

---

The update notification system ensures you stay current with SilentCast improvements while respecting your workflow and preferences.