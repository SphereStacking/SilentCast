# Timeout Notifications

SilentCast provides comprehensive timeout notification support to keep you informed when scripts exceed their configured time limits.

## Overview

When a script execution times out, SilentCast can:
- Send notifications when timeouts occur
- Provide warnings before timeouts
- Show partial output captured before termination
- Indicate whether the process was terminated gracefully or forcefully

## Configuration

### Global Notification Settings

Configure notification behavior in your `spellbook.yml`:

```yaml
notification:
  enable_timeout: true       # Enable timeout notifications (default: true)
  enable_warning: true       # Enable warning before timeout (default: true)
  sound: true               # Play sound for notifications (default: true)
  max_output_length: 1024   # Max output length in notifications (default: 1024)
```

### Per-Action Timeout Configuration

Each grimoire entry can have its own timeout settings:

```yaml
grimoire:
  my_script:
    type: script
    command: "long-running-command"
    timeout: 30             # Timeout after 30 seconds
    timeout_warning: 10     # Warn 10 seconds before timeout
    grace_period: 5         # Allow 5 seconds for graceful shutdown
    show_output: true       # Required for output notifications
```

## Timeout Notification Types

### 1. Timeout Warning

Sent before the actual timeout occurs:

```yaml
timeout_warning: 10  # Warn 10 seconds before timeout
```

The warning notification will:
- Have a "Warning" level (yellow/orange color)
- Show remaining time until timeout
- Display the command that will be terminated

### 2. Timeout Notification

Sent when a script is terminated due to timeout:

```yaml
timeout: 30  # Timeout after 30 seconds
```

The timeout notification includes:
- Action name and description
- Configured timeout duration
- Actual elapsed time
- Whether termination was graceful or forced
- Partial output (if `show_output: true`)

## Graceful Shutdown

SilentCast attempts to terminate processes gracefully:

```yaml
grace_period: 5  # Allow 5 seconds for graceful termination
```

The shutdown process:
1. Send SIGTERM (Unix) or taskkill (Windows) when timeout is reached
2. Wait for the grace period
3. Force kill if process hasn't exited

Scripts can handle graceful shutdown:

```bash
#!/bin/bash
# Handle SIGTERM for cleanup
trap 'echo "Cleaning up..."; exit 0' TERM

# Long-running process
while true; do
    # Do work...
    sleep 1
done
```

## Examples

### Basic Timeout with Notification

```yaml
grimoire:
  backup:
    type: script
    command: "rsync -av /source/ /backup/"
    description: "Backup files"
    timeout: 300          # 5 minutes
    show_output: true     # Enable notifications
```

### Timeout with Warning

```yaml
grimoire:
  deploy:
    type: script
    command: "./deploy.sh production"
    description: "Deploy to production"
    timeout: 600          # 10 minutes
    timeout_warning: 60   # Warn 1 minute before
    show_output: true
```

### Graceful Database Backup

```yaml
grimoire:
  db_backup:
    type: script
    command: "pg_dump mydb > backup.sql"
    description: "Database backup"
    timeout: 1800         # 30 minutes
    timeout_warning: 300  # Warn 5 minutes before
    grace_period: 60      # Allow 1 minute for cleanup
    show_output: true
```

### Silent Timeout (No Notification)

```yaml
grimoire:
  background_task:
    type: script
    command: "./background-task.sh"
    timeout: 60
    # No show_output means no timeout notification
```

## Platform-Specific Behavior

### Windows
- Uses `taskkill /T /PID` for process termination
- Terminates entire process tree
- Grace period allows child processes to exit

### macOS/Linux
- Uses SIGTERM for graceful termination
- Sends signal to process group
- Falls back to SIGKILL after grace period

## Best Practices

1. **Set Reasonable Timeouts**: Allow enough time for normal completion but prevent runaway processes

2. **Use Warnings for Long Tasks**: Give yourself time to check on long-running tasks:
   ```yaml
   timeout: 3600         # 1 hour
   timeout_warning: 300  # Warn 5 minutes before
   ```

3. **Handle Signals in Scripts**: Make your scripts timeout-friendly:
   ```bash
   trap 'cleanup; exit 0' TERM INT
   ```

4. **Enable Output for Important Tasks**: Always use `show_output: true` for tasks where you need timeout feedback

5. **Adjust Grace Period**: Give databases and services time to shutdown cleanly:
   ```yaml
   grace_period: 30  # 30 seconds for database shutdown
   ```

## Troubleshooting

### No Timeout Notification Received

Check that:
- `show_output: true` is set for the action
- Global `notification.enable_timeout` is not set to `false`
- The script actually exceeded the timeout

### Process Not Terminating

Increase the grace period or check if the process is handling signals:
```yaml
grace_period: 10  # Increase from default 5 seconds
```

### Warning Not Appearing

Ensure `timeout_warning` is less than `timeout`:
```yaml
timeout: 60
timeout_warning: 10  # Must be less than timeout
```