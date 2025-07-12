# Troubleshooting

This guide helps you resolve common issues with SilentCast.

## Common Issues

### Hotkeys Not Working

<details>
<summary>macOS: "SilentCast" wants to control this computer</summary>

**Problem**: Hotkeys don't work and you see an accessibility permission prompt.

**Solution**:
1. Open System Preferences → Security & Privacy → Privacy → Accessibility
2. Click the lock to make changes
3. Add SilentCast to the list (or check the box if already present)
4. Restart SilentCast

</details>

<details>
<summary>Windows: Hotkeys not registering</summary>

**Problem**: Pressing the prefix key does nothing.

**Solution**:
1. Check if another application is using the same hotkey
2. Try running SilentCast as administrator
3. Verify the hotkey format in your configuration (e.g., `alt+space`, not `Alt+Space`)

</details>

<details>
<summary>Hotkey conflicts with other applications</summary>

**Problem**: Your chosen prefix key is already used by another application.

**Solution**:
Change the prefix key in your configuration:

```yaml
hotkeys:
  prefix: "ctrl+shift+space"  # or another combination
```

Common alternatives:
- `ctrl+shift+space`
- `alt+shift+p`
- `ctrl+alt+s`

</details>

### Configuration Issues

<details>
<summary>Configuration file not found</summary>

**Problem**: SilentCast can't find your spellbook.yml file.

**Solution**:
1. Check the default locations:
   - macOS/Linux: `~/.config/silentcast/spellbook.yml`
   - Windows: `%APPDATA%\SilentCast\spellbook.yml`

2. Specify a custom location:
   ```bash
   silentcast --config /path/to/spellbook.yml
   ```

3. Set the environment variable:
   ```bash
   export SILENTCAST_CONFIG=/path/to/config/directory
   ```

📖 **See also**: [Windows configuration setup guide](/guide/faq#i-get-hotkeys-prefix-is-required-error-on-windows-what-do-i-do) for detailed step-by-step instructions.

</details>

<details>
<summary>Invalid configuration syntax</summary>

**Problem**: YAML syntax errors in your configuration.

**Solution**:
1. Validate your configuration:
   ```bash
   silentcast --validate-config
   ```

2. Common YAML issues:
   - Use spaces, not tabs
   - Proper indentation (2 spaces)
   - Quote strings with special characters
   - Escape backslashes in Windows paths

3. Use a YAML validator: https://www.yamllint.com/

</details>

<details>
<summary>Changes not taking effect</summary>

**Problem**: Configuration changes aren't applied.

**Solution**:
1. Ensure config watching is enabled:
   ```yaml
   daemon:
     config_watch: true
   ```

2. Manually reload:
   - System tray → Reload Config
   - Or restart SilentCast

3. Check logs for configuration errors:
   ```bash
   tail -f ~/.local/share/silentcast/silentcast.log
   ```

</details>

### Application Launch Issues

<details>
<summary>Application not found</summary>

**Problem**: "command not found" or "application not found" errors.

**Solution**:
1. Use full paths for applications:
   ```yaml
   grimoire:
     editor:
       type: app
       command: "/Applications/Visual Studio Code.app/Contents/Resources/app/bin/code"  # macOS
       # command: "C:\\Program Files\\Microsoft VS Code\\Code.exe"  # Windows
   ```

2. Add application directories to PATH
3. For macOS apps, use the `.app` bundle path

</details>

<details>
<summary>Applications launch but don't focus</summary>

**Problem**: Applications open in the background.

**Solution**:
Add focus arguments:
```yaml
grimoire:
  editor:
    type: app
    command: "code"
    args: ["--new-window", "--foreground"]
```

</details>

### Script Execution Issues

<details>
<summary>Scripts fail with permission denied</summary>

**Problem**: Shell scripts can't execute.

**Solution**:
1. Make scripts executable:
   ```bash
   chmod +x /path/to/script.sh
   ```

2. Specify the shell explicitly:
   ```yaml
   grimoire:
     my_script:
       type: script
       command: "/path/to/script.sh"
       shell: "/bin/bash"
   ```

</details>

<details>
<summary>Environment variables not available</summary>

**Problem**: Scripts can't access expected environment variables.

**Solution**:
1. Explicitly set required variables:
   ```yaml
   grimoire:
     build:
       type: script
       command: "npm run build"
       env:
         NODE_ENV: "production"
         PATH: "${PATH}:/usr/local/bin"
   ```

2. Source profile in scripts:
   ```yaml
   command: "source ~/.bashrc && npm run build"
   ```

</details>

### Performance Issues

<details>
<summary>High CPU usage</summary>

**Problem**: SilentCast uses excessive CPU.

**Solution**:
1. Disable config watching if not needed:
   ```yaml
   daemon:
     config_watch: false
   ```

2. Increase update check interval:
   ```yaml
   daemon:
     update_interval: 168  # Weekly instead of daily
   ```

3. Check logs for errors causing loops

</details>

<details>
<summary>Slow response time</summary>

**Problem**: Delays between pressing hotkeys and action execution.

**Solution**:
1. Reduce timeouts:
   ```yaml
   hotkeys:
     timeout: 500          # Reduce from 1000
     sequence_timeout: 1000  # Reduce from 2000
   ```

2. Use simpler commands (avoid complex scripts)
3. Check system resources

</details>

## Debugging

### Enable Debug Logging

```yaml
daemon:
  log_level: debug

logger:
  level: debug
  file: "~/.local/share/silentcast/debug.log"
```

### View Logs

```bash
# Real-time log viewing
tail -f ~/.local/share/silentcast/silentcast.log

# Search for errors
grep ERROR ~/.local/share/silentcast/silentcast.log

# View with timestamp
less +F ~/.local/share/silentcast/silentcast.log
```

### Test Mode

Run SilentCast in test mode to see what would happen:

```bash
# Dry run mode
silentcast --dry-run

# Verbose output
silentcast --log-level debug --no-tray
```

## Getting Help

If you're still experiencing issues:

1. **Check existing issues**: [GitHub Issues](https://github.com/SphereStacking/silentcast/issues)
2. **Report a bug**: Include:
   - SilentCast version (`silentcast --version`)
   - Operating system and version
   - Configuration file (remove sensitive data)
   - Debug logs
   - Steps to reproduce

## Common Error Messages

| Error | Meaning | Solution |
|-------|---------|----------|
| `Failed to register hotkey` | Hotkey already in use | Change prefix key |
| `Permission denied: accessibility` | macOS accessibility needed | Grant permission in System Preferences |
| `Config file not found` | Missing spellbook.yml | Create configuration file |
| `Invalid YAML syntax` | Configuration parse error | Fix YAML formatting |
| `Command not found` | Application path incorrect | Use full path or fix PATH |
| `Failed to watch config` | File system permissions | Check file permissions |

## Platform-Specific Issues

### macOS

- **Gatekeeper blocks execution**: Right-click and select "Open" or `xattr -d com.apple.quarantine silentcast`
- **Hotkeys in secure input**: Some apps block global hotkeys (like password fields)
- **Mission Control conflicts**: Avoid F3, F4, and gesture-based shortcuts

### Windows

- **Antivirus interference**: Add SilentCast to exceptions
- **UAC prompts**: Run as administrator for system-wide hotkeys
- **Corporate policies**: May block hotkey registration

## Next Steps

- [Configuration Guide](/guide/configuration) - Optimize your setup
- [Platform Support](/guide/platforms) - OS-specific tips
- [FAQ](/guide/faq) - Frequently asked questions