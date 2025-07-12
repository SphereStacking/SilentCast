# Exit Codes

SilentCast uses standard exit codes to indicate the result of operations. Understanding these codes helps with scripting and troubleshooting.

## Standard Exit Codes

| Code | Constant | Description |
|------|----------|-------------|
| 0 | `EXIT_SUCCESS` | Successful execution |
| 1 | `EXIT_ERROR` | General error |
| 2 | `EXIT_USAGE` | Invalid command line usage |
| 3 | `EXIT_CONFIG` | Configuration error |
| 4 | `EXIT_PERMISSION` | Permission denied |
| 5 | `EXIT_NOT_FOUND` | File or command not found |
| 6 | `EXIT_TIMEOUT` | Operation timed out |
| 7 | `EXIT_INTERRUPT` | User interrupted (Ctrl+C) |
| 8 | `EXIT_HOTKEY` | Hotkey registration failed |
| 9 | `EXIT_ALREADY_RUNNING` | Another instance is running |

## Exit Code Details

### 0 - Success

Normal successful completion.

```bash
silentcast --version
echo $?  # 0
```

### 1 - General Error

Unspecified error occurred.

```bash
# Example: Unknown internal error
silentcast
echo $?  # 1
```

### 2 - Usage Error

Invalid command line arguments or options.

```bash
silentcast --invalid-option
echo $?  # 2

silentcast --config  # Missing argument
echo $?  # 2
```

### 3 - Configuration Error

Problems with configuration files.

```bash
# Invalid YAML syntax
silentcast --config invalid.yml
echo $?  # 3

# Missing required configuration
silentcast --validate-config
echo $?  # 3
```

### 4 - Permission Error

Insufficient permissions for operation.

```bash
# macOS: Missing accessibility permissions
silentcast
echo $?  # 4

# Cannot write to log file
silentcast --log-file /root/silentcast.log
echo $?  # 4
```

### 5 - Not Found

Required file or command not found.

```bash
# Configuration file doesn't exist
silentcast --config /path/to/missing.yml
echo $?  # 5

# Spell references non-existent action
# When executing a spell with missing grimoire entry
echo $?  # 5
```

### 6 - Timeout

Operation exceeded time limit.

```bash
# Script execution timeout
# When a script action exceeds its timeout
echo $?  # 6
```

### 7 - Interrupt

User interrupted execution.

```bash
# User pressed Ctrl+C
silentcast
^C
echo $?  # 7
```

### 8 - Hotkey Error

Failed to register global hotkeys.

```bash
# Hotkey already in use by another application
silentcast
echo $?  # 8

# Invalid hotkey format in configuration
echo $?  # 8
```

### 9 - Already Running

Another instance is already running.

```bash
# First instance
silentcast &

# Second instance
silentcast
echo $?  # 9
```

## Using Exit Codes in Scripts

### Bash Scripts

```bash
#!/bin/bash

# Start SilentCast and check result
if silentcast --validate-config; then
    echo "Configuration is valid"
    silentcast
else
    case $? in
        3)
            echo "Configuration error - check spellbook.yml"
            ;;
        5)
            echo "Configuration file not found"
            ;;
        *)
            echo "Unknown error: $?"
            ;;
    esac
    exit 1
fi
```

### PowerShell Scripts

```powershell
# Start SilentCast
& silentcast --validate-config

if ($LASTEXITCODE -eq 0) {
    Write-Host "Configuration is valid"
    & silentcast
} else {
    switch ($LASTEXITCODE) {
        3 { Write-Error "Configuration error" }
        5 { Write-Error "Configuration file not found" }
        default { Write-Error "Unknown error: $LASTEXITCODE" }
    }
    exit $LASTEXITCODE
}
```

### Python Scripts

```python
import subprocess
import sys

# Run SilentCast
result = subprocess.run(['silentcast', '--validate-config'])

if result.returncode == 0:
    print("Configuration is valid")
    subprocess.run(['silentcast'])
else:
    errors = {
        3: "Configuration error",
        5: "Configuration file not found",
        8: "Hotkey registration failed"
    }
    
    error_msg = errors.get(result.returncode, f"Unknown error: {result.returncode}")
    print(f"Error: {error_msg}", file=sys.stderr)
    sys.exit(result.returncode)
```

## Action Exit Codes

When executing actions, SilentCast returns the exit code from the action:

### Application Launch

```yaml
grimoire:
  editor:
    type: app
    command: "code"
    # Returns 0 if launched successfully
    # Returns app's exit code if it fails
```

### Script Execution

```yaml
grimoire:
  build:
    type: script
    command: "make build"
    # Returns the script's exit code
```

## Checking Exit Codes

### In Spellbook

```yaml
grimoire:
  deploy:
    type: script
    command: |
      if make test; then
        make deploy
      else
        echo "Tests failed with code $?"
        exit 1
      fi
```

### In Terminal

```bash
# Run action and check result
silentcast --execute git_status
if [ $? -eq 0 ]; then
    echo "Git status completed successfully"
fi
```

## Exit Code Handling

### Retry on Failure

```bash
#!/bin/bash

MAX_RETRIES=3
RETRY_COUNT=0

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    silentcast
    EXIT_CODE=$?
    
    case $EXIT_CODE in
        0)
            echo "SilentCast started successfully"
            break
            ;;
        8)
            echo "Hotkey conflict, retrying in 5 seconds..."
            sleep 5
            RETRY_COUNT=$((RETRY_COUNT + 1))
            ;;
        9)
            echo "Already running"
            exit 0
            ;;
        *)
            echo "Failed with exit code: $EXIT_CODE"
            exit $EXIT_CODE
            ;;
    esac
done
```

### Logging Exit Codes

```bash
# Log exit codes for debugging
silentcast 2>&1 | tee silentcast.log
EXIT_CODE=${PIPESTATUS[0]}
echo "$(date): SilentCast exited with code $EXIT_CODE" >> exit_codes.log
```

## Best Practices

1. **Always check exit codes** in scripts and automation
2. **Use specific codes** to handle different error scenarios
3. **Document custom codes** if extending SilentCast
4. **Log exit codes** for troubleshooting
5. **Test error paths** to ensure proper handling

## Custom Exit Codes

For scripts executed by SilentCast:

```bash
#!/bin/bash
# Custom script with meaningful exit codes

# Define custom codes
readonly SUCCESS=0
readonly ERROR_DB_CONNECTION=10
readonly ERROR_API_TIMEOUT=11
readonly ERROR_INVALID_DATA=12

# Use in script
if ! ping -c 1 database.local > /dev/null; then
    echo "Database unreachable"
    exit $ERROR_DB_CONNECTION
fi

# SilentCast will return these codes
```

## Next Steps

- [CLI Options](/api/) - Command line interface reference
- [Troubleshooting](/guide/troubleshooting) - Common issues and solutions
- [Testing](/api/testing) - Testing SilentCast configurations