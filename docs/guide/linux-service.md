# Linux Service Management

SilentCast provides flexible service management options for Linux users, supporting both systemd services and XDG autostart. This guide covers installation, configuration, and management of SilentCast as a Linux service.

## Overview

SilentCast can be configured to start automatically on Linux using:
- **Systemd user service**: For systemd-based distributions (Ubuntu, Fedora, Arch, etc.)
- **XDG autostart**: For desktop environments that support the XDG autostart specification

Both methods are installed simultaneously to ensure maximum compatibility.

## Quick Start

### Install Service

```bash
# Install as user service (recommended)
silentcast --service-install

# The following will be created:
# - ~/.config/systemd/user/silentcast.service
# - ~/.config/autostart/silentcast.desktop
```

### Manage Service

```bash
# Start the service
silentcast --service-start

# Stop the service
silentcast --service-stop

# Check service status
silentcast --service-status

# Uninstall the service
silentcast --service-uninstall
```

## Systemd Service

### Service File Location

The systemd service file is installed to:
- **User service**: `~/.config/systemd/user/silentcast.service`
- **System service**: `/etc/systemd/system/silentcast.service` (requires root)

### Service Configuration

The generated systemd service file includes:

```ini
[Unit]
Description=Silent hotkey-driven task runner
After=graphical-session.target

[Service]
Type=simple
ExecStart=/path/to/silentcast --no-tray
Restart=on-failure
RestartSec=5

# Security settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=read-only
ReadWritePaths=~/.config/silentcast ~/.local/share/silentcast

[Install]
WantedBy=default.target
```

### Manual Systemd Management

If you prefer manual control:

```bash
# Reload systemd after manual changes
systemctl --user daemon-reload

# Enable service to start at boot
systemctl --user enable silentcast

# Start service immediately
systemctl --user start silentcast

# Check service logs
journalctl --user -u silentcast -f
```

## XDG Autostart

### Desktop Entry Location

The XDG autostart file is installed to:
- `~/.config/autostart/silentcast.desktop`

### Desktop Entry Configuration

The generated desktop entry includes:

```ini
[Desktop Entry]
Type=Application
Name=SilentCast
Comment=Silent hotkey-driven task runner
Exec=/path/to/silentcast --no-tray
Terminal=false
Categories=Utility;System;
StartupNotify=false
X-GNOME-Autostart-enabled=true
Hidden=false
```

### Desktop Environment Support

XDG autostart is supported by most desktop environments:
- GNOME
- KDE Plasma
- XFCE
- LXDE/LXQt
- MATE
- Cinnamon

## Configuration

### Service Options

When running as a service, SilentCast automatically:
- Runs without system tray (`--no-tray`)
- Loads configuration from standard locations
- Logs to systemd journal (for systemd service)

### Environment Variables

You can set environment variables in the service file:

```bash
# Edit the service file
systemctl --user edit silentcast

# Add environment variables
[Service]
Environment="SILENTCAST_CONFIG=/custom/path/spellbook.yml"
Environment="SILENTCAST_DEBUG=true"
```

## Troubleshooting

### Service Won't Start

1. Check service status:
   ```bash
   silentcast --service-status
   systemctl --user status silentcast
   ```

2. View service logs:
   ```bash
   journalctl --user -u silentcast -n 50
   ```

3. Test manual execution:
   ```bash
   silentcast --no-tray --debug
   ```

### Permission Issues

1. Ensure executable has correct permissions:
   ```bash
   chmod +x /path/to/silentcast
   ```

2. Check hotkey permissions:
   - Some distributions require additional permissions for global hotkeys
   - See [Permissions Guide](./permissions.md)

### Autostart Not Working

1. Verify desktop entry is enabled:
   ```bash
   cat ~/.config/autostart/silentcast.desktop | grep Hidden
   # Should show: Hidden=false
   ```

2. Check desktop environment compatibility:
   ```bash
   echo $XDG_CURRENT_DESKTOP
   ```

3. Test desktop entry manually:
   ```bash
   dex ~/.config/autostart/silentcast.desktop
   ```

## Advanced Configuration

### Custom Service File

Create a custom systemd service:

```bash
# Create service directory
mkdir -p ~/.config/systemd/user

# Create custom service file
cat > ~/.config/systemd/user/silentcast-custom.service << EOF
[Unit]
Description=SilentCast Custom Configuration
After=graphical-session.target

[Service]
Type=simple
ExecStart=/usr/local/bin/silentcast --config /home/user/custom/spellbook.yml
Restart=always
RestartSec=10

[Install]
WantedBy=default.target
EOF

# Enable and start
systemctl --user enable silentcast-custom
systemctl --user start silentcast-custom
```

### Multiple Instances

Run multiple SilentCast instances with different configurations:

```bash
# Create instance service
systemctl --user edit --force silentcast@work.service

# Add configuration
[Unit]
Description=SilentCast - Work Profile

[Service]
ExecStart=/usr/local/bin/silentcast --config %h/.config/silentcast/work.yml

# Enable instance
systemctl --user enable silentcast@work
systemctl --user start silentcast@work
```

## Security Considerations

The systemd service includes security hardening:

- **NoNewPrivileges**: Prevents privilege escalation
- **PrivateTmp**: Isolates temporary files
- **ProtectSystem**: Makes system directories read-only
- **ProtectHome**: Restricts home directory access
- **ReadWritePaths**: Explicitly allows config and log directories

## Distribution-Specific Notes

### Ubuntu/Debian

```bash
# Ensure user systemd is enabled
loginctl enable-linger $USER
```

### Fedora/RHEL

```bash
# SELinux may require additional configuration
setsebool -P user_execmod on
```

### Arch Linux

```bash
# User services are enabled by default
# Just use standard commands
```

## Best Practices

1. **Use user services**: Prefer user systemd services over system-wide installation
2. **Check logs regularly**: Monitor service health through journalctl
3. **Test configuration**: Validate spellbook before enabling service
4. **Backup before updates**: Save configuration before upgrading

## Related Documentation

- [Installation Guide](installation.md)
- [Configuration Reference](../config/)
- [Permissions Guide](./permissions.md)
- [Troubleshooting Guide](../troubleshooting.md)