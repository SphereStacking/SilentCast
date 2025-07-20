# Troubleshooting Guide

This guide helps you solve common issues with SilentCast.

## Common Issues

### Spells Not Working

#### Hotkeys Not Triggering
- **Check if SilentCast is running**: Look for the system tray icon
- **Verify permissions**: On macOS, grant Accessibility permissions
- **Check prefix key conflicts**: Ensure no other app uses your prefix key
- **Test with simple spells**: Try basic single-key spells first

#### "Permission denied" Errors
- **macOS**: Grant Accessibility and Input Monitoring permissions in System Preferences
- **Windows**: Run as administrator if accessing system files
- **Linux**: Check file permissions and group membership

### Output and Display Issues

#### Script Output Not Visible
The `show_output` feature displays command output in system notifications:

```yaml
grimoire:
  git_status:
    type: script
    command: "git status --short"
    show_output: true  # Shows output in notification
    description: "Git status"
```

#### Long-Running Commands
For interactive or long-running commands, use terminal mode:

```yaml
grimoire:
  interactive_python:
    type: script
    command: "python3"
    terminal: true      # Opens in terminal
    keep_open: true     # Keeps terminal open
    description: "Python REPL"
```

### Platform-Specific Issues

#### Windows

##### PowerShell Execution Policy
If scripts fail to run in PowerShell:

```powershell
# Check current policy
Get-ExecutionPolicy

# Set to allow local scripts (as Administrator)
Set-ExecutionPolicy RemoteSigned
```

##### Command Prompt vs PowerShell
Specify the shell explicitly for better compatibility:

```yaml
grimoire:
  cmd_command:
    type: script
    command: "dir"
    shell: "cmd"        # Use Command Prompt
    show_output: true
    
  ps_command:
    type: script
    command: "Get-ChildItem"
    shell: "powershell" # Use PowerShell
    show_output: true
```

##### Windows Terminal Integration
For better terminal experience:

```yaml
grimoire:
  dev_terminal:
    type: script
    command: "wt"       # Opens Windows Terminal
    args: ["-d", "C:\\dev"]
    description: "Open dev terminal"
```

#### macOS

##### Accessibility Permissions
SilentCast requires Accessibility permissions to capture global hotkeys:

1. Open **System Preferences** → **Security & Privacy** → **Privacy**
2. Select **Accessibility** from the list
3. Click the lock icon and enter your password
4. Add SilentCast to the list or check the existing entry

##### Application Not Found
If apps aren't launching:

```yaml
grimoire:
  vscode:
    type: app
    command: "/Applications/Visual Studio Code.app/Contents/Resources/app/bin/code"
    description: "VS Code (full path)"
```

#### Linux

##### Display Server Issues
On Wayland, some hotkey functionality may be limited. Consider using X11 mode.

##### Missing Dependencies
Install required packages:

```bash
# Ubuntu/Debian
sudo apt install libx11-dev libxtst-dev

# Fedora
sudo dnf install libX11-devel libXtst-devel

# Arch
sudo pacman -S libx11 libxtst
```

### Configuration Issues

#### Configuration Not Loading
1. **Check file location**: Ensure `spellbook.yml` is in the expected location
2. **Validate YAML syntax**: Use a YAML validator online
3. **Check file permissions**: Ensure SilentCast can read the file
4. **Review logs**: Check the log file for specific errors

#### Spells Not Responding
1. **Check spell syntax**: Ensure proper YAML formatting
2. **Verify grimoire entries**: Make sure referenced actions exist
3. **Test individual components**: Test spells and grimoire entries separately

### Performance Issues

#### High CPU Usage
- **Check for infinite loops** in scripts
- **Reduce update frequency** for auto-updating features
- **Use `show_output: false`** for background tasks

#### Memory Usage
- **Clear output buffers** periodically for long-running instances
- **Limit script output size** with timeouts
- **Restart SilentCast** periodically for long-running sessions

## Debugging Tips

### Enable Debug Logging
```yaml
daemon:
  log_level: debug  # Increases log verbosity
```

### Test Configuration
```bash
# Validate configuration
silentcast --validate-config

# Test specific spell
silentcast --test-spell git_status

# Run without tray for console output
silentcast --no-tray --debug
```

### Check System Compatibility
```bash
# Check platform support
silentcast --version

# Test hotkey registration
silentcast --test-hotkey "alt+space"
```

## Advanced Workarounds

### Custom Terminal Commands
For platforms where terminal integration is limited:

```yaml
grimoire:
  custom_terminal:
    type: script
    command: "start cmd /k echo 'Custom terminal' && pause"
    description: "Windows: Custom terminal window"
    
  unix_terminal:
    type: script
    command: "gnome-terminal -- bash -c 'echo Custom terminal; read'"
    description: "Linux: Custom terminal window"
```

### Environment-Specific Configurations
Use platform-specific configuration files:

- `spellbook.windows.yml` - Windows-specific settings
- `spellbook.macos.yml` - macOS-specific settings  
- `spellbook.linux.yml` - Linux-specific settings

These files override settings in the main `spellbook.yml`.

## Getting Help

If you're still experiencing issues:

1. **Check the logs**: Look in the log file for specific error messages
2. **Search existing issues**: Visit the [GitHub Issues](https://github.com/SphereStacking/silentcast/issues)
3. **Create a new issue**: Include:
   - Operating system and version
   - SilentCast version
   - Configuration file (remove sensitive data)
   - Log output with debug enabled
   - Steps to reproduce the issue

## Known Limitations

- **Wayland support**: Limited hotkey functionality on Wayland display servers
- **Snap/Flatpak**: Sandboxed applications may have limited file system access
- **Remote desktop**: Hotkeys may not work in remote desktop sessions
- **Virtual machines**: Some virtualization software may interfere with global hotkeys

## Best Practices

1. **Start simple**: Begin with basic spells and gradually add complexity
2. **Test incrementally**: Test each spell individually before combining
3. **Use version control**: Keep your configuration in version control
4. **Document custom spells**: Add clear descriptions to your grimoire entries
5. **Regular backups**: Backup your configuration files regularly