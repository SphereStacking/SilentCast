# CLI Reference

Complete command-line interface reference for SilentCast, including all options, flags, and usage examples.

## 📋 Synopsis

```bash
silentcast [OPTIONS]
silentcast --once --spell=SPELL [OPTIONS]
silentcast --test-spell --spell=SPELL [OPTIONS]
silentcast --validate-config [OPTIONS]
```

## 🚀 Quick Reference

```bash
# Common commands
silentcast                        # Start with system tray
silentcast --no-tray              # Start without tray
silentcast --debug                # Debug mode
silentcast --validate-config      # Validate spellbook
silentcast --list-spells          # Show all spells
silentcast --once --spell=e       # Execute spell once
silentcast --dry-run --spell=g,s  # Preview execution
silentcast --test-hotkey          # Test key detection
silentcast --version              # Show version
silentcast --help                 # Show help
```

## 🎯 Core Options

### `--help, -h`
Display help information and exit.

```bash
silentcast --help
silentcast -h
```

### `--version, -v`
Display version information and exit.

```bash
silentcast --version
silentcast -v
```

**Output example:**
```
SilentCast v1.0.0
Build: 2024-01-15T10:30:00Z
Commit: abc123def
Go: 1.21.5
```

### `--no-tray`
Run without system tray icon. Useful for:
- Headless environments
- Testing and debugging
- Server deployments
- When system tray is unavailable

```bash
silentcast --no-tray
silentcast --no-tray --debug
```

### `--debug, -d`
Enable debug-level logging for detailed output.

```bash
silentcast --debug
silentcast -d
```

**Debug output includes:**
- Configuration file paths searched
- Configuration loading and merging
- Hotkey registration details
- Key press events
- Spell matching logic
- Action execution steps
- Error details with stack traces

### `--quiet, -q`
Suppress all non-error output.

```bash
silentcast --quiet
silentcast -q --once --spell=backup
```

### `--log-level`
Set logging verbosity level.

```bash
silentcast --log-level=debug    # Most verbose
silentcast --log-level=info     # Default
silentcast --log-level=warn     # Warnings and errors
silentcast --log-level=error    # Errors only
```

## 📁 Configuration Commands

### `--config, -c`
Use a specific configuration file instead of the default search paths.

```bash
silentcast --config=/path/to/spellbook.yml
silentcast -c ./test-config.yml
```

### `--validate-config`
Validate configuration syntax and references.

```bash
silentcast --validate-config
silentcast --validate-config --config=./test.yml
```

**Validation checks:**
- YAML syntax validity
- Required fields presence
- Type validation for all fields
- Spell → grimoire action references
- Command/application existence
- Working directory accessibility
- Environment variable expansion

**Output example:**
```
✓ Configuration syntax valid
✓ Found 15 spells
✓ Found 15 grimoire entries
✓ All spell references valid
✓ All commands accessible
Configuration is valid!
```

### `--show-config`
Display the fully merged configuration.

```bash
# Default format (pretty-printed)
silentcast --show-config

# JSON format (for parsing)
silentcast --show-config --format=json

# YAML format
silentcast --show-config --format=yaml

# Include file paths
silentcast --show-config --show-paths
```

**Options:**
- `--format=pretty` (default) - Human-readable format
- `--format=json` - JSON output for scripts
- `--format=yaml` - YAML format
- `--show-paths` - Include configuration file paths

### `--show-config-path`
Show configuration search paths and found files.

```bash
silentcast --show-config-path
```

**Output example:**
```
Configuration search paths (in order):
1. ✓ /home/user/.config/silentcast/spellbook.yml (found)
2. ✗ /home/user/.config/silentcast/spellbook.linux.yml (not found)
3. ✗ /home/user/.silentcast/spellbook.yml (not found)
4. ✗ /etc/silentcast/spellbook.yml (not found)

Using: /home/user/.config/silentcast/spellbook.yml
```

## 🪄 Spell Management

### `--list-spells, -l`
List all configured spells.

```bash
# List all spells
silentcast --list-spells
silentcast -l

# Filter by keyword
silentcast --list-spells --filter=git
silentcast --list-spells --filter="docker"

# Compact format
silentcast --list-spells --compact

# JSON format
silentcast --list-spells --format=json
```

**Output formats:**
- Default: Table with spell, action, and description
- `--compact`: One spell per line
- `--format=json`: JSON array
- `--format=csv`: CSV format

**Example output:**
```
Spells:
┌─────────┬─────────────┬─────────────────────────────┐
│ Spell   │ Action      │ Description                 │
├─────────┼─────────────┼─────────────────────────────┤
│ e       │ editor      │ Open VS Code                │
│ t       │ terminal    │ Open terminal               │
│ g,s     │ git_status  │ Show git status             │
│ g,p     │ git_pull    │ Pull from origin            │
└─────────┴─────────────┴─────────────────────────────┘
```

### `--test-spell`
Test a specific spell configuration.

```bash
# Test single-key spell
silentcast --test-spell --spell=e

# Test sequence spell
silentcast --test-spell --spell=g,s

# Test with debug output
silentcast --test-spell --spell=docker,up --debug
```

**Test includes:**
- Spell existence
- Action mapping
- Command validation
- Permission checks
- Environment resolution
- Working directory access

### `--dry-run`
Preview spell execution without running.

```bash
# Preview single spell
silentcast --dry-run --spell=deploy

# Preview with environment
silentcast --dry-run --spell=backup --debug

# Show expanded command
silentcast --dry-run --spell=e --show-expanded
```

**Output includes:**
- Action type and name
- Full command to execute
- Working directory
- Environment variables
- Shell to use
- Action options

**Example output:**
```
Dry run for spell 'g,s':
  Action: git_status
  Type: script
  Command: git status --short
  Working Dir: /home/user/project
  Shell: /bin/bash
  Options:
    - show_output: true
    - timeout: 30s
```

## 🎯 Execution Modes

### `--once, -o`
Execute a single spell and exit.

```bash
# Execute and exit
silentcast --once --spell=editor
silentcast -o --spell=g,s

# With output capture
silentcast --once --spell=backup > backup.log

# Silent execution
silentcast --once --spell=deploy --quiet
```

**Exit codes:**
- `0` - Success
- `1` - General error
- `2` - Spell not found
- `3` - Action execution failed
- `4` - Permission denied

### `--spell, -s`
Specify spell for execution (used with --once, --dry-run, --test-spell).

```bash
silentcast --once --spell=e
silentcast --dry-run --spell=g,s
silentcast --test-spell --spell=d,u
```

**Spell formats:**
- Single key: `e`, `t`, `b`
- Sequence: `g,s`, `d,u`, `w,m`
- Special chars: `ctrl+s`, `shift+f1`

## 🧪 Testing & Debugging

### `--test-hotkey`
Interactive hotkey testing mode.

```bash
# Test until Ctrl+C
silentcast --test-hotkey

# Test for specific duration
silentcast --test-hotkey --duration=30

# Test with specific prefix
silentcast --test-hotkey --prefix="ctrl+alt+space"
```

**Shows:**
- Raw key events
- Detected modifiers
- Key combinations
- Sequence building
- Timeout handling

**Example output:**
```
Hotkey test mode (press Ctrl+C to exit)
Prefix key: alt+space

[10:30:45] Key pressed: alt
[10:30:45] Key pressed: space
[10:30:45] Prefix detected! Waiting for spell...
[10:30:46] Key pressed: g
[10:30:46] Building sequence: g
[10:30:47] Key pressed: s
[10:30:47] Sequence complete: g,s
[10:30:47] Would execute spell: git_status
```

### `--benchmark`
Run performance benchmarks.

```bash
# Run all benchmarks
silentcast --benchmark

# Benchmark specific operations
silentcast --benchmark --type=config-load
silentcast --benchmark --type=spell-lookup
silentcast --benchmark --type=action-exec
```

### `--diagnose`
Run system diagnostics.

```bash
silentcast --diagnose
```

**Checks:**
- Configuration validity
- File permissions
- Command availability
- System tray support
- Hotkey system status
- Platform-specific requirements

## 🌍 Environment Variables

### `SILENTCAST_CONFIG`
Override configuration file or directory.

```bash
# Use specific file
SILENTCAST_CONFIG=/path/to/spellbook.yml silentcast

# Use directory
SILENTCAST_CONFIG=/etc/myconfig silentcast
```

### `SILENTCAST_LOG_LEVEL`
Set default log level.

```bash
SILENTCAST_LOG_LEVEL=debug silentcast
SILENTCAST_LOG_LEVEL=error silentcast --no-tray
```

**Values:** `debug`, `info`, `warn`, `error`

### `SILENTCAST_LOG_FILE`
Override log file location.

```bash
SILENTCAST_LOG_FILE=/var/log/silentcast.log silentcast
SILENTCAST_LOG_FILE=/tmp/debug.log silentcast --debug
```

### `SILENTCAST_NO_COLOR`
Disable colored output.

```bash
SILENTCAST_NO_COLOR=1 silentcast --list-spells
```

### `SILENTCAST_PREFIX_KEY`
Override configured prefix key.

```bash
SILENTCAST_PREFIX_KEY="ctrl+shift+space" silentcast
```

## 📊 Advanced Options

### `--format`
Output format for various commands.

```bash
# Configuration output
silentcast --show-config --format=json
silentcast --show-config --format=yaml

# Spell listing
silentcast --list-spells --format=json
silentcast --list-spells --format=csv
silentcast --list-spells --format=table
```

### `--filter`
Filter output for list commands.

```bash
# Filter spells
silentcast --list-spells --filter=git
silentcast --list-spells --filter="^d,"  # Regex

# Case-insensitive
silentcast --list-spells --filter=GIT --ignore-case
```

### `--timeout`
Set operation timeout.

```bash
# Set global timeout
silentcast --timeout=60s

# For specific operations
silentcast --once --spell=backup --timeout=5m
silentcast --test-spell --spell=deploy --timeout=30s
```

### `--workers`
Set number of concurrent workers.

```bash
silentcast --workers=4
silentcast --workers=1 --debug  # Sequential for debugging
```

## 💡 Usage Examples

### Development Workflow

```bash
# Start development session
silentcast --debug --no-tray

# Validate after config changes
silentcast --validate-config

# Test new spell
silentcast --test-spell --spell=new_spell
silentcast --dry-run --spell=new_spell

# Execute once for testing
silentcast --once --spell=new_spell --debug
```

### CI/CD Integration

```bash
#!/bin/bash
# CI/CD script

# Validate configuration
if ! silentcast --validate-config; then
    echo "Invalid configuration"
    exit 1
fi

# Run tests
silentcast --once --spell=test || exit 1

# Build
silentcast --once --spell=build || exit 1

# Deploy with timeout
silentcast --once --spell=deploy --timeout=10m || exit 1
```

### Automation Scripts

```bash
# Backup script
#!/bin/bash
LOG_FILE="/var/log/backup.log"

# Run backup spell
if silentcast --once --spell=backup >> "$LOG_FILE" 2>&1; then
    echo "Backup completed successfully"
else
    echo "Backup failed with exit code $?"
    exit 1
fi
```

### Debugging Issues

```bash
# Full debug mode
silentcast --debug --log-level=debug 2>&1 | tee debug.log

# Test hotkey detection
silentcast --test-hotkey --debug

# Diagnose system
silentcast --diagnose > diagnosis.txt

# Test specific spell with full output
silentcast --test-spell --spell=problematic --debug
```

## 🚦 Exit Codes

| Code | Meaning | Description |
|------|---------|-------------|
| 0 | Success | Operation completed successfully |
| 1 | General Error | Unspecified error occurred |
| 2 | Config Error | Configuration file invalid or not found |
| 3 | Spell Error | Spell not found or invalid |
| 4 | Permission Error | Insufficient permissions |
| 5 | Hotkey Error | Failed to register hotkeys |
| 6 | Execution Error | Action execution failed |
| 7 | Timeout | Operation timed out |
| 8 | Cancelled | Operation cancelled by user |

## 🔧 Troubleshooting

### Common Issues

**Hotkeys not working:**
```bash
# Test hotkey system
silentcast --test-hotkey --debug

# Check permissions (macOS)
silentcast --diagnose | grep "Accessibility"
```

**Configuration not loading:**
```bash
# Show search paths
silentcast --show-config-path

# Validate syntax
silentcast --validate-config --debug
```

**Spell not found:**
```bash
# List all spells
silentcast --list-spells

# Test specific spell
silentcast --test-spell --spell=myspell --debug
```

## 📚 See Also

- [Configuration Guide](./configuration.md) - Spellbook configuration
- [Spell Patterns](./spells.md) - Keyboard combinations
- [Scripting Guide](./scripts.md) - Advanced automation
- [Platform Guide](./platforms.md) - OS-specific features

---

<div align="center">
  <p><strong>Need more help? Run `silentcast --help` 🪄</strong></p>
</div>