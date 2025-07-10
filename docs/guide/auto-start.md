# Auto-start Guide

Configure SilentCast to start automatically when your system boots, ensuring your keyboard shortcuts are always available. This guide covers platform-specific setup and best practices for auto-start configuration.

## Overview

SilentCast can be configured to start automatically in several ways:

1. **Built-in auto-start** - Using the `daemon.auto_start` configuration
2. **System service** - Installing as a system service
3. **User login** - Starting when you log in
4. **Desktop environment** - Using your DE's autostart features

## Quick Setup

### Enable Auto-start in Configuration

```yaml
# spellbook.yml
daemon:
  auto_start: true  # Enable auto-start
```

Then install the service:

```bash
# Install auto-start service
silentcast --install-service

# Verify installation
silentcast --service-status
```

## Platform-Specific Setup

### Windows

#### Method 1: Task Scheduler (Recommended)

```powershell
# Run as Administrator
silentcast --install-service

# Or manually create scheduled task
$Action = New-ScheduledTaskAction -Execute "$env:LOCALAPPDATA\SilentCast\silentcast.exe"
$Trigger = New-ScheduledTaskTrigger -AtLogon -User $env:USERNAME
$Settings = New-ScheduledTaskSettingsSet -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries

Register-ScheduledTask -TaskName "SilentCast" `
    -Action $Action `
    -Trigger $Trigger `
    -Settings $Settings `
    -Description "SilentCast - Silent Hotkey Task Runner"
```

#### Method 2: Startup Folder

```powershell
# Create shortcut in startup folder
$WshShell = New-Object -comObject WScript.Shell
$Shortcut = $WshShell.CreateShortcut("$env:APPDATA\Microsoft\Windows\Start Menu\Programs\Startup\SilentCast.lnk")
$Shortcut.TargetPath = "$env:LOCALAPPDATA\SilentCast\silentcast.exe"
$Shortcut.Arguments = "--config `"$env:APPDATA\SilentCast\spellbook.yml`""
$Shortcut.Save()
```

#### Method 3: Registry (Advanced)

```powershell
# Add to registry (requires admin)
New-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Run" `
    -Name "SilentCast" `
    -Value "`"$env:LOCALAPPDATA\SilentCast\silentcast.exe`"" `
    -PropertyType String -Force
```

### macOS

#### Method 1: launchd Service (Recommended)

```bash
# Install service
silentcast --install-service

# Or manually create launch agent
cat > ~/Library/LaunchAgents/com.silentcast.plist << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.silentcast</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/silentcast</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <dict>
        <key>SuccessfulExit</key>
        <false/>
    </dict>
    <key>StandardErrorPath</key>
    <string>/tmp/silentcast.err</string>
    <key>StandardOutPath</key>
    <string>/tmp/silentcast.out</string>
    <key>EnvironmentVariables</key>
    <dict>
        <key>PATH</key>
        <string>/usr/local/bin:/usr/bin:/bin</string>
    </dict>
</dict>
</plist>
EOF

# Load the service
launchctl load ~/Library/LaunchAgents/com.silentcast.plist

# Enable auto-start
launchctl enable user/$(id -u)/com.silentcast
```

#### Method 2: Login Items

```bash
# Add to login items via osascript
osascript -e 'tell application "System Events" to make login item at end with properties {path:"/usr/local/bin/silentcast", hidden:false}'

# Or use the GUI:
# System Preferences → Users & Groups → Login Items → Add SilentCast
```


## Configuration Options

### Auto-start Settings

```yaml
daemon:
  # Enable auto-start installation
  auto_start: true
  
  # Start minimized to tray
  start_minimized: true
  
  # Delay startup (seconds)
  startup_delay: 5
  
  # Check if already running
  single_instance: true
```

### Conditional Auto-start

```yaml
# Platform-specific auto-start
# spellbook.windows.yml
daemon:
  auto_start: true
  startup_delay: 10  # Windows needs more time

# spellbook.mac.yml
daemon:
  auto_start: true
  startup_delay: 5   # macOS is faster

```

## Managing Auto-start

### Check Status

```bash
# Check if auto-start is enabled
silentcast --service-status

# Output
Service Status:
━━━━━━━━━━━━━━
Installed: ✓ Yes
Running: ✓ Active
Auto-start: ✓ Enabled
Start time: 2024-01-15 08:00:00
```

### Disable Auto-start

```bash
# Remove auto-start
silentcast --uninstall-service

# Or disable temporarily
# Windows
Disable-ScheduledTask -TaskName "SilentCast"

# macOS
launchctl unload ~/Library/LaunchAgents/com.silentcast.plist

```

### Update Auto-start Configuration

After changing configuration:

```bash
# Reload service
# Windows
Restart-Service SilentCast

# macOS
launchctl unload ~/Library/LaunchAgents/com.silentcast.plist
launchctl load ~/Library/LaunchAgents/com.silentcast.plist

```

## Best Practices

### 1. Startup Delay

Add a delay to ensure system resources are ready:

```yaml
daemon:
  auto_start: true
  startup_delay: 10  # Wait 10 seconds
```

### 2. Logging for Auto-start

Enable logging to troubleshoot startup issues:

```yaml
logger:
  level: info
  file: "~/.local/share/silentcast/startup.log"
  max_size: 5
  max_backups: 2
```

### 3. Environment Variables

Ensure required environment variables are set:

# macOS launchd
<key>EnvironmentVariables</key>
<dict>
    <key>PATH</key>
    <string>/usr/local/bin:/usr/bin:/bin</string>
</dict>
```

### 4. Error Recovery

Configure automatic restart on failure:

```yaml
# In service definitions
Restart=on-failure
RestartSec=5
StartLimitIntervalSec=60
StartLimitBurst=3
```

## Troubleshooting

### Auto-start Not Working

<details>
<summary>Windows: Task not running</summary>

1. Check Task Scheduler:
   ```powershell
   Get-ScheduledTask -TaskName "SilentCast"
   ```

2. Check event logs:
   ```powershell
   Get-WinEvent -LogName System | Where-Object {$_.Message -like "*SilentCast*"}
   ```

3. Run manually to test:
   ```powershell
   Start-ScheduledTask -TaskName "SilentCast"
   ```

</details>

<details>
<summary>macOS: Accessibility permissions</summary>

Auto-start may fail if accessibility permissions aren't granted:

1. Check permissions:
   ```bash
   sqlite3 /Library/Application\ Support/com.apple.TCC/TCC.db \
     "SELECT * FROM access WHERE service='kTCCServiceAccessibility'"
   ```

2. Grant permissions manually:
   - System Preferences → Security & Privacy → Privacy → Accessibility
   - Add SilentCast

3. Reset permissions if needed:
   ```bash
   tccutil reset Accessibility com.silentcast
   ```

</details>


### Multiple Instances

Prevent multiple instances:

```yaml
daemon:
  single_instance: true
  pid_file: "~/.local/share/silentcast/silentcast.pid"
```

Check for existing instance:
```bash
# Check if already running
pgrep -f silentcast || echo "Not running"

# Kill existing instance
pkill -f silentcast
```

### Delayed Startup Issues

If SilentCast starts too early:

1. Increase startup delay:
   ```yaml
   daemon:
     startup_delay: 30  # 30 seconds
   ```

2. Add dependencies (systemd):
   ```ini
   [Unit]
   After=graphical-session.target network-online.target
   Wants=network-online.target
   ```

3. Check system readiness:
   ```yaml
   grimoire:
     startup_check:
       type: script
       command: |
         # Wait for network
         while ! ping -c 1 google.com &> /dev/null; do
           sleep 1
         done
         # Start SilentCast
         exec silentcast
   ```

## Security Considerations

### Permissions

- Run with user privileges, not root/admin
- Store configuration in user directory
- Avoid storing sensitive data in auto-start scripts

### Startup Scripts

```yaml
# Good: Use environment variables
grimoire:
  secure_startup:
    type: script
    command: "source ~/.env && silentcast"

# Bad: Hardcoded credentials
grimoire:
  insecure_startup:
    type: script
    command: "API_KEY=secret123 silentcast"
```

## Platform-Specific Notes

### Windows
- Requires "Log on as a batch job" permission for Task Scheduler
- UAC may interfere with some shortcuts
- Consider using highest available privileges

### macOS
- Requires accessibility permissions before auto-start
- May need to disable Gatekeeper for unsigned binaries
- Login items are user-visible


## Next Steps

- Configure [Logging](/guide/logging) for startup debugging
- Set up [Environment Variables](/guide/env-vars) for auto-start
- Review [Platform Guide](/guide/platforms) for OS-specific details
- Check [Troubleshooting](/guide/troubleshooting) for common issues