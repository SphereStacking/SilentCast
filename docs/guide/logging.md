# Logging Configuration

SilentCast provides comprehensive logging functionality to track application behavior and troubleshoot issues.

## Log File Locations

By default, if no log file path is specified in the configuration, SilentCast will create a log file in the same directory as your configuration files:

### Windows
- **Config Directory**: `%APPDATA%\SilentCast\`
- **Default Log File**: `%APPDATA%\SilentCast\silentcast.log`

Example:
```
C:\Users\YourUsername\AppData\Roaming\SilentCast\silentcast.log
```

### macOS
- **Config Directory**: `~/Library/Application Support/SilentCast/`
- **Default Log File**: `~/Library/Application Support/SilentCast/silentcast.log`

### Linux
- **Config Directory**: `~/.config/SilentCast/`
- **Default Log File**: `~/.config/SilentCast/silentcast.log`

### Console Output
If no log file is configured and the default cannot be created, logs will be output to the console (stderr) only.

## Configuration Options

You can customize logging behavior in your `spellbook.yml`:

```yaml
logger:
  # Log level: debug, info, warn, error
  level: info
  
  # Log file path (uses default if empty)
  file: ""
  
  # Maximum log file size in MB
  max_size: 10
  
  # Number of old log files to keep
  max_backups: 3
  
  # Maximum age of log files in days
  max_age: 7
  
  # Compress old log files
  compress: true
```

## Custom Log File Path

To save logs to a specific location, use the `file` setting:

```yaml
logger:
  file: "C:\\Logs\\silentcast.log"  # Windows
  # file: "/var/log/silentcast.log"  # Linux/macOS
```

## Log Levels

- **debug**: Records all information (development/debugging)
- **info**: Records normal operation information (default)
- **warn**: Records only warning messages
- **error**: Records only error messages

## Viewing Log Files

### Windows PowerShell
```powershell
# View last 20 lines
Get-Content "$env:APPDATA\SilentCast\silentcast.log" -Tail 20

# Monitor logs in real-time
Get-Content "$env:APPDATA\SilentCast\silentcast.log" -Wait -Tail 10
```

### Windows Command Prompt
```cmd
# View log file
type "%APPDATA%\SilentCast\silentcast.log"

# Open in Notepad
notepad "%APPDATA%\SilentCast\silentcast.log"
```

### macOS/Linux
```bash
# View last 20 lines
tail -n 20 ~/Library/Application\ Support/SilentCast/silentcast.log  # macOS
tail -n 20 ~/.config/SilentCast/silentcast.log                       # Linux

# Monitor logs in real-time
tail -f ~/Library/Application\ Support/SilentCast/silentcast.log     # macOS
tail -f ~/.config/SilentCast/silentcast.log                          # Linux
```

## Troubleshooting

### Log File Not Found

1. Verify the application is running
2. Check if the config directory exists:
   ```cmd
   # Windows
   dir "%APPDATA%\SilentCast"
   
   # macOS
   ls ~/Library/Application\ Support/SilentCast/
   
   # Linux
   ls ~/.config/SilentCast/
   ```

3. Check if a custom path is specified in your configuration

### No Log Output

1. Verify log level is appropriate (setting to `error` will hide detailed information)
2. Check file write permissions
3. Verify available disk space

## Log Rotation

SilentCast automatically rotates log files:

- New file created when size reaches `max_size`
- Old files are numbered (`.1`, `.2`, etc.)
- Files exceeding `max_backups` are deleted
- With `compress: true`, old files are compressed as `.gz`

Example:
```
silentcast.log        # Current log
silentcast.log.1.gz   # Previous log (compressed)
silentcast.log.2.gz   # Older log (compressed)
```