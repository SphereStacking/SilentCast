# Logging Guide

SilentCast provides comprehensive logging capabilities to help you monitor, debug, and troubleshoot your automation workflows. This guide covers how to configure logging, understand log levels, and use logs effectively.

## Quick Start

By default, SilentCast logs to the console with `info` level. To enable file logging:

```yaml
# spellbook.yml
logger:
  level: info
  file: "~/.local/share/silentcast/silentcast.log"
```

View logs in real-time:
```bash
# Follow log file
tail -f ~/.local/share/silentcast/silentcast.log

# View with colors (if supported)
tail -f ~/.local/share/silentcast/silentcast.log | ccze -A
```

## Log Levels

SilentCast uses four log levels, from most to least verbose:

| Level | Usage | Example Output |
|-------|-------|----------------|
| `debug` | Detailed information for debugging | Key events, config loading, action execution details |
| `info` | General informational messages | Startup, spell casts, configuration changes |
| `warn` | Warning messages | Deprecated features, non-fatal errors |
| `error` | Error messages only | Failed actions, invalid configuration |

### Setting Log Level

```yaml
# In configuration
logger:
  level: debug  # debug, info, warn, error

# Via CLI
silentcast --log-level debug

# Via environment variable
SILENTCAST_LOG_LEVEL=debug silentcast
```

## Configuration Options

### Complete Logger Configuration

```yaml
logger:
  # Log level: debug, info, warn, error
  level: info
  
  # Log file path (empty = console only)
  # Supports ~ for home directory
  file: "~/.local/share/silentcast/silentcast.log"
  
  # Maximum size in MB before rotation
  max_size: 10
  
  # Number of old log files to keep
  max_backups: 3
  
  # Maximum age in days
  max_age: 7
  
  # Compress rotated log files
  compress: true
  
  # Include timestamps (console output)
  timestamps: true
  
  # Colorize output (console only)
  color: true
  
  # Log format: text, json
  format: text
```

### Log File Paths

Default log locations by platform:

::: code-group

```bash [macOS]
# Default log location
~/.local/share/silentcast/silentcast.log

# Alternative locations
~/.silentcast/logs/silentcast.log
/var/log/silentcast/silentcast.log  # System-wide installation
```

```powershell [Windows]
# Default log location
%LOCALAPPDATA%\SilentCast\logs\silentcast.log

# Alternative locations
%APPDATA%\SilentCast\silentcast.log
C:\ProgramData\SilentCast\logs\silentcast.log
```

:::

## Log Rotation

SilentCast automatically rotates logs based on size and age:

```yaml
logger:
  file: "~/logs/silentcast.log"
  max_size: 10        # Rotate when log reaches 10MB
  max_backups: 3      # Keep 3 old files
  max_age: 7          # Delete files older than 7 days
  compress: true      # Compress old logs to .gz
```

Rotated files are named:
```
silentcast.log          # Current log
silentcast.log.1        # Most recent rotation
silentcast.log.2        # Older rotation
silentcast.log.3.gz     # Oldest (compressed)
```

## Log Message Format

### Text Format (Default)

```
2024-01-15 10:23:45 [INFO]  SilentCast v1.0.0 started
2024-01-15 10:23:45 [INFO]  Loaded configuration from ~/.config/silentcast/spellbook.yml
2024-01-15 10:23:45 [INFO]  Registered 15 spells
2024-01-15 10:23:46 [DEBUG] Hotkey registered: alt+space
2024-01-15 10:24:12 [INFO]  Spell cast: git_status
2024-01-15 10:24:12 [DEBUG] Executing action: type=script, command=git status
2024-01-15 10:24:12 [INFO]  Action completed successfully
```

### JSON Format

```yaml
logger:
  format: json
```

```json
{"time":"2024-01-15T10:23:45Z","level":"INFO","msg":"SilentCast v1.0.0 started","version":"1.0.0"}
{"time":"2024-01-15T10:23:45Z","level":"INFO","msg":"Loaded configuration","file":"~/.config/silentcast/spellbook.yml","spells":15}
{"time":"2024-01-15T10:24:12Z","level":"INFO","msg":"Spell cast","spell":"git_status","action":"git_status"}
```

## Debug Logging

Enable debug logging to see detailed information:

```bash
# Start with debug logging
silentcast --debug

# Or set in configuration
logger:
  level: debug
```

Debug logs include:
- Configuration parsing details
- Key event detection
- Action execution steps
- Environment variable expansion
- Error stack traces

### Example Debug Output

```
[DEBUG] Loading configuration from ~/.config/silentcast/spellbook.yml
[DEBUG] Parsing YAML configuration
[DEBUG] Found 3 platform override files
[DEBUG] Merging platform-specific configuration: spellbook.mac.yml
[DEBUG] Validating spell definitions
[DEBUG] Spell 'git_status' -> action 'git_status' (type: script)
[DEBUG] Starting hotkey manager
[DEBUG] Registering hotkey: alt+space (prefix)
[DEBUG] Hotkey registration successful
[DEBUG] Key event: KeyDown alt
[DEBUG] Key event: KeyDown space
[DEBUG] Prefix activated, waiting for spell...
[DEBUG] Key event: KeyDown g
[DEBUG] Key event: KeyDown s
[DEBUG] Spell detected: g,s -> git_status
[DEBUG] Executing action: git_status
[DEBUG] Working directory: /home/user/project
[DEBUG] Environment: PATH=/usr/bin:/bin, USER=john
[DEBUG] Command: git status
[DEBUG] Action completed in 45ms
```

## Logging Best Practices

### 1. Development vs Production

```yaml
# Development configuration
logger:
  level: debug
  file: "./dev-silentcast.log"
  timestamps: true
  color: true

# Production configuration  
logger:
  level: info
  file: "/var/log/silentcast/silentcast.log"
  max_size: 50
  max_backups: 5
  compress: true
  format: json  # For log aggregation
```

### 2. Troubleshooting Configuration

Create a debug configuration for troubleshooting:

```yaml
# debug-spellbook.yml
logger:
  level: debug
  file: "./debug.log"
  timestamps: true
  
# Run with debug config
silentcast --config debug-spellbook.yml
```

### 3. Log Monitoring

Monitor logs in real-time during development:

```bash
# Basic monitoring
tail -f ~/.local/share/silentcast/silentcast.log

# Filter by level
tail -f silentcast.log | grep -E "\[ERROR\]|\[WARN\]"

# Watch for specific spells
tail -f silentcast.log | grep "Spell cast"

# Colorize output (install ccze first)
tail -f silentcast.log | ccze -A
```

### 4. Log Analysis

Analyze logs for patterns:

```bash
# Count spell usage
grep "Spell cast" silentcast.log | awk -F': ' '{print $NF}' | sort | uniq -c | sort -rn

# Find errors
grep "\[ERROR\]" silentcast.log

# Actions by hour
grep "Spell cast" silentcast.log | awk '{print $2}' | cut -d: -f1 | sort | uniq -c

# Failed actions
grep -B2 -A2 "Action failed" silentcast.log
```

## Script Output Logging

When scripts show output, it's included in logs:

```yaml
grimoire:
  system_check:
    type: script
    command: "df -h && free -h"
    show_output: true  # Output appears in logs
```

Log output:
```
[INFO]  Spell cast: system_check
[DEBUG] Executing script: df -h && free -h
[INFO]  Script output:
        Filesystem      Size  Used Avail Use% Mounted on
        /dev/sda1       100G   45G   55G  45% /
        Memory:         16Gi  8.2Gi  7.8Gi
[INFO]  Action completed successfully
```

## Performance Considerations

### Log Level Impact

| Level | Performance Impact | Disk Usage |
|-------|-------------------|------------|
| `error` | Minimal | Very Low |
| `warn` | Minimal | Low |
| `info` | Low | Moderate |
| `debug` | Moderate | High |

### Optimization Tips

1. **Production Settings**:
   ```yaml
   logger:
     level: info      # Not debug
     compress: true   # Compress old logs
     max_age: 7       # Don't keep forever
   ```

2. **High-Volume Environments**:
   ```yaml
   logger:
     format: json     # Faster parsing
     timestamps: false # Saved in JSON already
     color: false     # No ANSI codes
   ```

3. **Disk-Constrained Systems**:
   ```yaml
   logger:
     max_size: 5      # Smaller rotation size
     max_backups: 1   # Fewer backups
     compress: true   # Always compress
   ```

## Integration with External Tools

### Syslog Integration

Forward logs to syslog:

```yaml
grimoire:
  forward_to_syslog:
    type: script
    command: |
      tail -f ~/.local/share/silentcast/silentcast.log | logger -t silentcast
    keep_open: true
```

### Log Aggregation

For centralized logging:

```yaml
# JSON format for easy parsing
logger:
  format: json
  file: "/var/log/silentcast/silentcast.log"

# Ship with Filebeat, Fluentd, etc.
```

### Monitoring Integration

Create alerts based on log patterns:

```yaml
grimoire:
  check_errors:
    type: script
    command: |
      ERROR_COUNT=$(grep -c "\[ERROR\]" ~/.local/share/silentcast/silentcast.log || echo 0)
      if [ $ERROR_COUNT -gt 10 ]; then
        notify-send "SilentCast Errors" "$ERROR_COUNT errors in log"
      fi
```

## Troubleshooting

### Common Issues

<details>
<summary>Logs not appearing</summary>

1. Check file permissions:
   ```bash
   ls -la ~/.local/share/silentcast/
   ```

2. Verify configuration:
   ```bash
   silentcast --show-config | grep -A10 "logger:"
   ```

3. Try console output:
   ```yaml
   logger:
     file: ""  # Empty = console only
   ```

</details>

<details>
<summary>Log file growing too large</summary>

Adjust rotation settings:
```yaml
logger:
  max_size: 5        # Smaller files
  max_backups: 2     # Fewer backups
  max_age: 3         # Shorter retention
  compress: true     # Always compress
```

</details>

<details>
<summary>Can't see debug messages</summary>

Ensure debug level is set:
```bash
# Check current level
silentcast --show-config | grep "level:"

# Override via CLI
silentcast --log-level debug

# Or environment
SILENTCAST_LOG_LEVEL=debug silentcast
```

</details>

## Security Considerations

### Sensitive Information

Never log sensitive data:

```yaml
grimoire:
  # BAD: Password visible in logs
  bad_example:
    type: script
    command: "mysql -u root -pSecretPass123"
    
  # GOOD: Use environment variables
  good_example:
    type: script
    command: "mysql -u root -p$DB_PASS"
    env:
      DB_PASS: "${DB_PASSWORD}"  # Set securely
```

### Log File Permissions

Ensure proper permissions:

```bash
# Restrict log file access
chmod 600 ~/.local/share/silentcast/silentcast.log

# Restrict directory
chmod 700 ~/.local/share/silentcast/
```

## Next Steps

- Configure [Auto-start](/guide/auto-start) with proper logging
- Use [Environment Variables](/guide/env-vars) for dynamic log paths
- See [CLI Options](/api/) for runtime log control
- Check [Troubleshooting](/guide/troubleshooting) for common issues