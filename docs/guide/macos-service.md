# macOS Service (LaunchAgent)

SilentCast can be installed as a macOS LaunchAgent to run automatically at login and provide system-wide hotkey functionality.

## Overview

When running as a LaunchAgent:
- Starts automatically when you log in
- Runs in the background without dock icon
- Provides hotkeys across all applications
- Managed through launchctl
- Logs to standard macOS locations

## LaunchAgent vs LaunchDaemon

SilentCast uses **LaunchAgent** (user-level) by default:
- **LaunchAgent**: Runs when user logs in, has access to user's GUI session
- **LaunchDaemon**: Runs at boot, no GUI access (not suitable for hotkeys)

## Requirements

- macOS 10.12 (Sierra) or later
- Accessibility permissions (for global hotkeys)
- No special privileges needed for user-level installation

## Installation

### 1. Grant Accessibility Permission

Before installing as a service, ensure SilentCast has accessibility permission:

1. Open System Preferences → Security & Privacy → Privacy
2. Select "Accessibility" from the left panel
3. Click the lock to make changes
4. Add SilentCast to the list (or check if already present)

### 2. Install LaunchAgent

```bash
silentcast --service-install
```

This will:
- Create a LaunchAgent plist in `~/Library/LaunchAgents/`
- Configure automatic start at login
- Set up crash recovery
- Configure logging to `~/Library/Logs/`

### 3. Start the Service

The service will start automatically on next login, or start it now:

```bash
silentcast --service-start
```

## Management Commands

### Check Status
```bash
silentcast --service-status
```

Shows:
- Installation status
- Running state
- Plist file location

### Stop Service
```bash
silentcast --service-stop
```

### Start Service
```bash
silentcast --service-start
```

### Uninstall Service
```bash
silentcast --service-uninstall
```

This will:
- Stop the service if running
- Remove the LaunchAgent plist
- Clean up launchctl registration

## Configuration

### Plist Location

The LaunchAgent plist is installed at:
```
~/Library/LaunchAgents/com.spherestacking.silentcast.plist
```

### Plist Structure

The generated plist includes:
```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" 
  "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.spherestacking.silentcast</string>
    
    <key>ProgramArguments</key>
    <array>
        <string>/path/to/silentcast</string>
        <string>--no-tray</string>
    </array>
    
    <key>RunAtLoad</key>
    <true/>
    
    <key>KeepAlive</key>
    <dict>
        <key>SuccessfulExit</key>
        <false/>
        <key>Crashed</key>
        <true/>
    </dict>
    
    <key>ProcessType</key>
    <string>Interactive</string>
    
    <key>StandardOutPath</key>
    <string>~/Library/Logs/silentcast.log</string>
    
    <key>StandardErrorPath</key>
    <string>~/Library/Logs/silentcast.error.log</string>
</dict>
</plist>
```

### Configuration Files

The service uses the same configuration as the desktop version:
- `~/.config/silentcast/spellbook.yml` - User configuration
- `/etc/silentcast/spellbook.yml` - System-wide configuration (if exists)

### Logging

Service logs are written to:
- **Standard output**: `~/Library/Logs/silentcast.log`
- **Error output**: `~/Library/Logs/silentcast.error.log`
- **Application logs**: As configured in spellbook.yml

View logs:
```bash
# Real-time log monitoring
tail -f ~/Library/Logs/silentcast.log

# View errors
tail -f ~/Library/Logs/silentcast.error.log

# Using Console.app
open -a Console ~/Library/Logs/silentcast.log
```

## Troubleshooting

### Service Won't Install

**Error: "failed to load service"**
```bash
# Check if already loaded
launchctl list | grep silentcast

# If found, unload first
launchctl unload ~/Library/LaunchAgents/com.spherestacking.silentcast.plist
```

### Service Won't Start

**Check accessibility permissions:**
```bash
# Run SilentCast directly to test
silentcast --debug --no-tray
```

**Check plist syntax:**
```bash
plutil ~/Library/LaunchAgents/com.spherestacking.silentcast.plist
```

**View launchctl errors:**
```bash
# Check service status
launchctl print gui/$(id -u)/com.spherestacking.silentcast
```

### Hotkeys Not Working

1. **Verify Accessibility Permission**:
   - System Preferences → Security & Privacy → Privacy → Accessibility
   - Ensure SilentCast is checked

2. **Check Process is Running**:
   ```bash
   ps aux | grep silentcast
   ```

3. **Test Hotkeys Directly**:
   ```bash
   # Stop service and run in debug mode
   silentcast --service-stop
   silentcast --debug --test-hotkey
   ```

### Service Keeps Restarting

Check crash logs:
```bash
# View crash reports
ls ~/Library/Logs/DiagnosticReports/ | grep silentcast

# Check launchd logs
log show --predicate 'process == "launchd"' --last 1h | grep silentcast
```

## Best Practices

### Security

1. **Permissions**: Only grant necessary permissions
2. **Configuration**: Secure your spellbook.yml with appropriate file permissions
3. **Scripts**: Be cautious with script actions that run with your user privileges

### Performance

1. **Startup**: Service adds minimal overhead to login time
2. **Memory**: Monitor memory usage if using many file watchers
3. **CPU**: Should remain under 1% during idle

### Updates

When updating SilentCast:
1. Stop the service: `silentcast --service-stop`
2. Replace the binary
3. Start the service: `silentcast --service-start`

Or reinstall:
```bash
silentcast --service-uninstall
# Update binary
silentcast --service-install
```

## Advanced Configuration

### Custom Plist Modifications

To add custom settings, edit the plist after installation:

```bash
# Edit plist
nano ~/Library/LaunchAgents/com.spherestacking.silentcast.plist

# Reload service
launchctl unload ~/Library/LaunchAgents/com.spherestacking.silentcast.plist
launchctl load ~/Library/LaunchAgents/com.spherestacking.silentcast.plist
```

### Environment Variables

Add environment variables to the plist:
```xml
<key>EnvironmentVariables</key>
<dict>
    <key>SILENTCAST_CONFIG</key>
    <string>/custom/path/to/config</string>
</dict>
```

### Resource Limits

Add resource constraints:
```xml
<key>HardResourceLimits</key>
<dict>
    <key>NumberOfFiles</key>
    <integer>1024</integer>
</dict>
```

## Migration from Login Items

If you previously added SilentCast to Login Items:

1. Remove from Login Items:
   - System Preferences → Users & Groups → Login Items
   - Select SilentCast and click "-"

2. Install as LaunchAgent:
   ```bash
   silentcast --service-install
   ```

## Integration with Other Tools

### Alfred/Raycast

LaunchAgent runs independently of launcher tools:
- No conflicts with Alfred, Raycast, etc.
- Can use both simultaneously
- Hotkeys remain available

### Homebrew Services

If installed via Homebrew:
```bash
# Use Homebrew services instead
brew services start silentcast
brew services stop silentcast
```

## See Also

- [Configuration Guide](configuration.md)
- [Troubleshooting Guide](troubleshooting.md)
- [Platform-Specific Features](platforms.md)