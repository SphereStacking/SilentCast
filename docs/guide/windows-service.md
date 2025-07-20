# Windows Service

SilentCast can be installed as a Windows service to run automatically at system startup and be available for all users.

## Overview

When running as a Windows service:
- SilentCast starts automatically when Windows boots
- Runs in the background without requiring user login
- Available system-wide for all users
- Managed through Windows Service Manager
- Logs events to Windows Event Log

## Requirements

- Windows 7 or later
- Administrator privileges for installation/management
- .NET Framework 4.5 or later (usually pre-installed)

## Installation

### 1. Install as Service

Open Command Prompt or PowerShell as Administrator and run:

```cmd
silentcast.exe --service-install
```

This will:
- Install SilentCast as a Windows service named "SilentCast"
- Set the service to start automatically at boot
- Configure recovery actions for automatic restart on failure

### 2. Start the Service

After installation, start the service:

```cmd
silentcast.exe --service-start
```

Or use Windows Service Manager:
1. Press `Win+R`, type `services.msc`, press Enter
2. Find "SilentCast Hotkey Service"
3. Right-click and select "Start"

## Management Commands

All service commands require Administrator privileges:

### Check Status
```cmd
silentcast.exe --service-status
```

Shows:
- Installation status
- Running state
- Start type (automatic/manual/disabled)

### Stop Service
```cmd
silentcast.exe --service-stop
```

### Start Service
```cmd
silentcast.exe --service-start
```

### Uninstall Service
```cmd
silentcast.exe --service-uninstall
```

This will:
- Stop the service if running
- Remove the service registration
- Clean up Event Log entries

## Configuration

### Service Properties

The service is configured with:
- **Name**: SilentCast
- **Display Name**: SilentCast Hotkey Service
- **Start Type**: Automatic (starts at boot)
- **Recovery**: Restart on failure (3 attempts)

### Configuration Files

The service uses the same configuration files as the desktop version:
- `%APPDATA%\silentcast\spellbook.yml` - User configuration
- `%PROGRAMDATA%\silentcast\spellbook.yml` - System-wide configuration

### Logging

Service logs are written to:
1. **Windows Event Log**: System events and errors
   - View in Event Viewer under "Applications and Services Logs"
2. **Application Log**: Detailed application logs
   - Default: `%PROGRAMDATA%\silentcast\silentcast.log`

## Permissions

### Service Account

By default, the service runs as:
- **Local System**: Full system access
- Can interact with desktop sessions
- Access to all user profiles

### File Access

The service has access to:
- System directories
- All user profile directories
- Network resources (if configured)

## Troubleshooting

### Service Won't Install

**Error: "Failed to connect to service manager"**
- Ensure running as Administrator
- Check Windows Service Manager is running

**Error: "Service already exists"**
- Uninstall existing service first:
  ```cmd
  silentcast.exe --service-uninstall
  ```

### Service Won't Start

**Error: "The service did not respond in a timely fashion"**
1. Check configuration file is valid:
   ```cmd
   silentcast.exe --validate-config
   ```
2. Check for missing dependencies
3. Review Windows Event Log for errors

**Error: "Access denied"**
- Ensure service account has necessary permissions
- Check file/folder permissions

### Hotkeys Not Working

When running as a service:
1. Service must have "Allow service to interact with desktop" enabled
2. Some hotkeys may not work across user sessions
3. UAC may block certain operations

To enable desktop interaction:
1. Open Services (services.msc)
2. Right-click "SilentCast" → Properties
3. Log On tab → Check "Allow service to interact with desktop"
4. Restart the service

### View Logs

**Windows Event Log:**
1. Open Event Viewer (eventvwr.msc)
2. Navigate to Windows Logs → Application
3. Filter by Source: "SilentCast"

**Application Log:**
```cmd
type %PROGRAMDATA%\silentcast\silentcast.log
```

## Best Practices

### Security

1. **Limit Service Permissions**: Consider running as a dedicated service account instead of Local System
2. **Restrict Configuration Access**: Set appropriate file permissions on configuration files
3. **Audit Actions**: Enable detailed logging for security-sensitive operations

### Performance

1. **Minimize Startup Actions**: Service startup affects boot time
2. **Configure Appropriate Timeouts**: Adjust hotkey timeouts for service mode
3. **Monitor Resource Usage**: Check CPU/memory usage in Task Manager

### Configuration

1. **Test Before Service Mode**: Validate configuration in desktop mode first
2. **Use System-Wide Config**: Place shared configuration in `%PROGRAMDATA%`
3. **Document Hotkeys**: Maintain documentation of configured hotkeys

## Example Workflows

### Basic Setup

```cmd
# 1. Install and validate configuration
silentcast.exe --validate-config

# 2. Install as service
silentcast.exe --service-install

# 3. Start service
silentcast.exe --service-start

# 4. Verify status
silentcast.exe --service-status
```

### Maintenance

```cmd
# Stop service for maintenance
silentcast.exe --service-stop

# Update configuration
notepad %APPDATA%\silentcast\spellbook.yml

# Validate changes
silentcast.exe --validate-config

# Restart service
silentcast.exe --service-start
```

### Troubleshooting Workflow

```cmd
# 1. Check service status
silentcast.exe --service-status

# 2. Stop service
silentcast.exe --service-stop

# 3. Run in debug mode
silentcast.exe --debug --no-tray

# 4. Fix issues, then restart service
silentcast.exe --service-start
```

## Advanced Configuration

### Custom Service Name

To install with a custom service name (for multiple instances):
```cmd
# Not implemented in current version
# Future feature for running multiple configurations
```

### Service Dependencies

Configure service dependencies in Windows:
1. Open Services (services.msc)
2. Right-click "SilentCast" → Properties
3. Dependencies tab → Add required services

### Group Policy

Deploy SilentCast via Group Policy:
1. Create MSI installer (see deployment guide)
2. Deploy via Group Policy Software Installation
3. Configure service settings via Group Policy Preferences

## Migration from Desktop Mode

To migrate from desktop to service mode:

1. **Export Current Configuration**:
   ```cmd
   silentcast.exe --export-config backup.yml
   ```

2. **Copy to System Location**:
   ```cmd
   copy backup.yml %PROGRAMDATA%\silentcast\spellbook.yml
   ```

3. **Install and Start Service**:
   ```cmd
   silentcast.exe --service-install
   silentcast.exe --service-start
   ```

4. **Remove from Startup** (if added):
   - Remove from user startup folder
   - Disable any scheduled tasks

## See Also

- [Configuration Guide](configuration.md)
- [Troubleshooting Guide](troubleshooting.md)
- [Platform-Specific Documentation](../troubleshooting/platform-specific.md)