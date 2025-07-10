# CLI Options Reference

SilentCast provides a comprehensive command-line interface for controlling its behavior. This reference covers all available options and their usage.

## Synopsis

```bash
silentcast [OPTIONS]
```

## Basic Usage

```bash
# Start SilentCast with default settings
silentcast

# Start without system tray
silentcast --no-tray

# Start with custom configuration
silentcast --config ~/my-spellbook.yml

# Show version and exit
silentcast --version

# Show help
silentcast --help
```

## Global Options

### `--help`, `-h`

Show help message and exit.

```bash
silentcast --help
```

### `--version`, `-v`

Display version information and exit.

```bash
silentcast --version
# Output: SilentCast v1.0.0 (git: abc123)
```

## Configuration Options

### `--config`, `-c`

Specify a custom configuration file path.

```bash
# Use specific configuration file
silentcast --config ~/dotfiles/silentcast/spellbook.yml

# Use configuration from current directory
silentcast --config ./spellbook.yml

# Multiple config files (last wins)
silentcast --config base.yml --config overrides.yml
```

**Default**: `~/.config/silentcast/spellbook.yml`

### `--validate-config`

Validate configuration file and exit. Useful for testing configurations.

```bash
# Validate default configuration
silentcast --validate-config

# Validate specific file
silentcast --config test.yml --validate-config

# Output on success
✓ Configuration is valid

# Output on error
✗ Configuration error in spellbook.yml:
  Line 15: Invalid action type "scritp" (did you mean "script"?)
```

### `--show-config`

Display the resolved configuration after all overrides and exit.

```bash
silentcast --show-config

# With custom config
silentcast --config custom.yml --show-config

# Output format
Resolved configuration:
━━━━━━━━━━━━━━━━━━━━━
daemon:
  auto_start: false
  log_level: info
  config_watch: true
hotkeys:
  prefix: alt+space
  timeout: 1000
...
```

### `--show-config-path`

Show where SilentCast looks for configuration files.

```bash
silentcast --show-config-path

# Output
Configuration search paths (in order):
1. ~/.config/silentcast/spellbook.yml ✓ (found)
2. ~/.silentcast/spellbook.yml ✗ (not found)
3. /etc/silentcast/spellbook.yml ✗ (not found)
```

## Runtime Options

### `--no-tray`

Run without system tray icon. Useful for terminal-only environments or debugging.

```bash
silentcast --no-tray
```

### `--dry-run`

Show what would be executed without actually running commands.

```bash
silentcast --dry-run

# When triggering shortcuts:
[DRY RUN] Would execute: code --new-window
[DRY RUN] Working directory: /home/user/projects
[DRY RUN] Environment: NODE_ENV=development
```

### `--once`

Run once and exit. Useful for testing single commands.

```bash
# Execute a specific spell and exit
silentcast --once --spell "git_status"
```

### `--spell`, `-s`

Execute a specific spell immediately (requires `--once`).

```bash
# Run a spell directly
silentcast --once --spell "git_status"

# Run with arguments
silentcast --once --spell "open_file" --args "README.md"
```

## Logging Options

### `--log-level`, `-l`

Set logging verbosity.

```bash
# Set log level
silentcast --log-level debug
silentcast -l warn

# Available levels
# - debug: Everything including key events
# - info:  General information (default)
# - warn:  Warnings only
# - error: Errors only
```

### `--log-file`

Specify log file location. Overrides configuration.

```bash
# Log to specific file
silentcast --log-file ~/silentcast-debug.log

# Log to stdout
silentcast --log-file -
```

### `--debug`, `-d`

Shortcut for `--log-level debug`.

```bash
silentcast --debug
# Equivalent to: silentcast --log-level debug
```

### `--quiet`, `-q`

Suppress all output except errors.

```bash
silentcast --quiet
# Equivalent to: silentcast --log-level error
```

### `--verbose`

Enable verbose output. Can be specified multiple times.

```bash
silentcast -v      # info level
silentcast -vv     # debug level
silentcast -vvv    # trace level (very verbose)
```

## Service Management

### `--install-service`

Install SilentCast as a system service (requires admin/root).

# macOS (launchd)
silentcast --install-service

# Windows (Task Scheduler)
silentcast --install-service
```

### `--uninstall-service`

Remove SilentCast system service.

# macOS/Windows
silentcast --uninstall-service
```

### `--service-status`

Check service installation status.

```bash
silentcast --service-status

# Output
Service Status:
━━━━━━━━━━━━━━
Installed: ✓ Yes
Running: ✓ Active
Auto-start: ✓ Enabled
```

## Development Options

### `--test-hotkey`

Test hotkey detection. Useful for debugging key combinations.

```bash
silentcast --test-hotkey

# Output
Hotkey Test Mode
━━━━━━━━━━━━━━━
Press any key combination (Ctrl+C to exit)

You pressed: alt+space
You pressed: g
You pressed: s
Detected sequence: g,s
```

### `--list-spells`

List all configured spells and their descriptions.

```bash
silentcast --list-spells

# Output
Available Spells:
━━━━━━━━━━━━━━━━
e         → editor         Open Visual Studio Code
t         → terminal       Open terminal emulator
b         → browser        Open web browser
g,s       → git_status     Show git repository status
g,p       → git_pull       Pull latest changes
d,u       → docker_up      Start Docker containers
```

### `--test-spell`

Test a spell without triggering via hotkey.

```bash
# Test specific spell
silentcast --test-spell git_status

# Test with dry-run
silentcast --test-spell git_status --dry-run
```

### `--benchmark`

Run performance benchmarks.

```bash
silentcast --benchmark

# Output
Performance Benchmark:
━━━━━━━━━━━━━━━━━━━━
Config load: 2.3ms
Hotkey init: 15.7ms
Tray init: 8.2ms
Total startup: 26.2ms
Memory usage: 14.8MB
```

## Environment Variables

### `SILENTCAST_CONFIG`

Override default configuration directory.

```bash
SILENTCAST_CONFIG=/custom/path silentcast
```

### `SILENTCAST_LOG_LEVEL`

Set log level via environment.

```bash
SILENTCAST_LOG_LEVEL=debug silentcast
```

### `SILENTCAST_NO_COLOR`

Disable colored output.

```bash
SILENTCAST_NO_COLOR=1 silentcast
```

### `SILENTCAST_HOME`

Override SilentCast home directory.

```bash
SILENTCAST_HOME=/opt/silentcast silentcast
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Configuration error |
| 3 | Permission error |
| 4 | Hotkey registration failed |
| 5 | Already running |
| 64 | Command line usage error |
| 77 | Permission denied |

## Examples

### Debug Configuration Issues

```bash
# Validate and show config
silentcast --validate-config --show-config

# Test with verbose logging
silentcast --debug --dry-run
```

### Production Setup

```bash
# Install as service with custom config
sudo silentcast --install-service --config /etc/silentcast/prod.yml

# Check everything is working
silentcast --service-status --list-spells
```

### Development Testing

```bash
# Test new configuration
silentcast --config dev-spellbook.yml --validate-config

# Run with debug output
silentcast --config dev-spellbook.yml --debug --no-tray

# Test specific spell
silentcast --once --spell "new_feature" --dry-run
```

### Troubleshooting

```bash
# Maximum verbosity
silentcast -vvv --log-file debug.log

# Test hotkey detection
silentcast --test-hotkey

# Check configuration paths
silentcast --show-config-path --verbose
```

## Platform-Specific Options

### Windows Only

`--console`
: Show console window (hidden by default)

`--startup-folder`
: Add to Windows startup folder

### macOS Only

`--accessibility-check`
: Check accessibility permissions

`--login-item`
: Add to login items

: Specify X11 display

`--desktop-file`
: Install desktop file for app launchers

## See Also

- [Configuration Reference](/config/) - Detailed configuration options
- [Key Names Reference](/api/key-names) - Valid key names for shortcuts
- [Environment Variables](/api/env-vars) - All environment variables
- [Exit Codes](/api/exit-codes) - Detailed exit code explanations