# CLI Reference

Complete command-line interface reference for SilentCast, including all options, flags, and usage examples.

## 📋 Synopsis

```bash
silentcast [OPTIONS]
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
silentcast --test-hotkey          # Test key detection
silentcast --version              # Show version
silentcast --help                 # Show help
```

## 🎯 Core Options

### `--help`
Display help information and exit.

```bash
silentcast --help
```

### `--version`
Display version information and exit.

```bash
silentcast --version
```

**Output example:**
```
SilentCast v1.0.0
Build: 2024-01-15T10:30:00Z
Commit: abc123def
Go: 1.23.5
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

### `--debug`
Enable debug-level logging for detailed output.

```bash
silentcast --debug
```

**Debug output includes:**
- Configuration file paths searched
- Configuration loading and merging
- Hotkey registration details
- Key press events
- Spell matching logic
- Action execution steps
- Error details with stack traces


### `--log-level`
Set logging verbosity level.

```bash
silentcast --log-level=debug    # Most verbose
silentcast --log-level=info     # Default
silentcast --log-level=warn     # Warnings and errors
silentcast --log-level=error    # Errors only
```

## 📁 Configuration Commands


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

### `--list-spells`
List all configured spells.

```bash
# List all spells
silentcast --list-spells
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


## 🎯 Execution Modes



## 🧪 Testing & Debugging

### `--test-hotkey`
Interactive hotkey testing mode.

```bash
# Test until Ctrl+C
silentcast --test-hotkey
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


## 💡 Usage Examples

### Development Workflow

```bash
# Start development session
silentcast --debug --no-tray

# Validate after config changes
silentcast --validate-config

# Test hotkey detection
silentcast --test-hotkey
```

### Debugging Issues

```bash
# Full debug mode
silentcast --debug --log-level=debug 2>&1 | tee debug.log

# Test hotkey detection
silentcast --test-hotkey

# Check configuration paths
silentcast --show-config-path
```

## 🚦 Exit Codes

| Code | Meaning | Description |
|------|---------|-------------|
| 0 | Success | Operation completed successfully |
| 1 | General Error | Unspecified error occurred |

**Future Implementation**: Additional exit codes (2-9) are planned for implementation to provide more specific error reporting for different failure scenarios.

## 🔧 Troubleshooting

### Common Issues

**Hotkeys not working:**
```bash
# Test hotkey system
silentcast --test-hotkey
```

**Configuration not loading:**
```bash
# Show search paths
silentcast --show-config-path

# Validate syntax
silentcast --validate-config
```

**Spell not found:**
```bash
# List all spells
silentcast --list-spells
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