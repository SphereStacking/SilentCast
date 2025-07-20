# Backup and Restore

SilentCast provides built-in tools for backing up and restoring your configuration files.

## Overview

The backup and restore features allow you to:
- Export configurations for backup purposes
- Share configurations with other users
- Migrate settings between machines
- Create configuration templates
- Version control your settings

## Exporting Configuration

### Basic Export

Export your current configuration to a YAML file:

```bash
silentcast --export-config backup.yml
```

This creates a clean YAML file with your merged configuration.

### Export to Standard Output

Export to stdout for piping or redirection:

```bash
# View configuration
silentcast --export-config -

# Pipe to another tool
silentcast --export-config - | grep spell

# Redirect to file
silentcast --export-config - > my-config.yml
```

### Archive Export

Create a tar.gz archive with all configuration files:

```bash
silentcast --export-config backup.tar.gz --export-format tar.gz
```

The archive includes:
- Main configuration (`spellbook.yml`)
- OS-specific overrides (`spellbook.windows.yml`, etc.)
- Export metadata with timestamp

## Importing Configuration

### Basic Import

Import a configuration from a file:

```bash
silentcast --import-config backup.yml
```

This will:
1. Validate the configuration
2. Backup existing configuration (if any)
3. Import the new configuration
4. Show validation warnings if any

### Import from Standard Input

Import from stdin for scripting:

```bash
# From a file
cat my-config.yml | silentcast --import-config -

# From a command
curl https://example.com/config.yml | silentcast --import-config -

# From heredoc
silentcast --import-config - << EOF
spells:
  e: editor
grimoire:
  editor:
    type: app
    command: code
EOF
```

### Import from Archive

Import from a tar.gz archive:

```bash
silentcast --import-config backup.tar.gz
```

The import process:
1. Extracts all YAML files from the archive
2. Validates each configuration file
3. Backs up existing files
4. Imports valid configurations
5. Shows any validation errors

## Backup Strategy

### Automatic Backups

When importing, SilentCast automatically creates timestamped backups:

```
spellbook.yml.backup.20240719-143022
```

### Manual Backup Workflow

```bash
# 1. Export current configuration
silentcast --export-config ~/backups/silentcast-$(date +%Y%m%d).yml

# 2. Make changes to configuration
# ...

# 3. If needed, restore from backup
silentcast --import-config ~/backups/silentcast-20240719.yml
```

### Version Control

Keep your configuration in version control:

```bash
# Export and commit
silentcast --export-config ~/.config/silentcast-backup/spellbook.yml
cd ~/.config/silentcast-backup
git add spellbook.yml
git commit -m "Update SilentCast configuration"
```

## Sharing Configurations

### Team Configurations

Share a base configuration with your team:

```bash
# Export clean configuration
silentcast --export-config team-config.yml

# Team members import
silentcast --import-config team-config.yml
```

### Configuration Templates

Create templates for different use cases:

```bash
# Development environment
silentcast --export-config dev-env.yml

# Production environment
silentcast --export-config prod-env.yml

# Switch between them
silentcast --import-config dev-env.yml
```

### Online Sharing

Share configurations via URLs:

```bash
# Import from URL (using curl)
curl -sL https://gist.github.com/user/config.yml | silentcast --import-config -

# Or download first
wget https://example.com/configs/developer.yml
silentcast --import-config developer.yml
```

## Migration

### Between Machines

Migrate configuration to a new machine:

```bash
# On old machine
silentcast --export-config - > silentcast-config.yml
# Transfer file to new machine

# On new machine
silentcast --import-config silentcast-config.yml
```

### Platform Migration

When moving between platforms (Windows ↔ macOS):

```bash
# Export full archive with OS-specific configs
silentcast --export-config full-backup.tar.gz --export-format tar.gz

# Import on new platform
silentcast --import-config full-backup.tar.gz
# Platform-specific configs will be ignored if not applicable
```

## Validation

### Pre-Import Validation

Import validates configuration before applying:

```bash
# Invalid config will show errors
silentcast --import-config invalid.yml
# ⚠️  Warning: Imported configuration has validation errors:
#    configuration validation failed: grimoire.missing.command: command is required
```

### Post-Import Validation

Always validate after import:

```bash
silentcast --import-config new-config.yml
silentcast --validate-config
```

## Best Practices

### 1. Regular Backups

Create a backup script:

```bash
#!/bin/bash
# backup-silentcast.sh
BACKUP_DIR="$HOME/backups/silentcast"
mkdir -p "$BACKUP_DIR"
silentcast --export-config "$BACKUP_DIR/$(date +%Y%m%d-%H%M%S).yml"

# Keep only last 30 backups
ls -t "$BACKUP_DIR"/*.yml | tail -n +31 | xargs -r rm
```

### 2. Safe Import

Always review before importing:

```bash
# First, examine the configuration
cat new-config.yml

# Then validate without importing
silentcast --export-config current.yml  # Backup current
silentcast --import-config new-config.yml
silentcast --validate-config

# If issues, restore
silentcast --import-config current.yml
```

### 3. Configuration Merging

Merge configurations manually:

```bash
# Export both configurations
silentcast --export-config current.yml
silentcast --export-config - > default.yml

# Merge with your preferred tool
# Then import merged result
silentcast --import-config merged.yml
```

## Troubleshooting

### Import Fails

**"Failed to read input"**
- Check file exists and is readable
- Verify file is valid YAML

**"Invalid YAML configuration"**
- Use a YAML validator
- Check for syntax errors
- Ensure proper indentation

**"Configuration validation failed"**
- Review error messages
- Check required fields
- Verify command paths exist

### Backup Issues

**Can't find backup files**
- Check working directory
- Look for `.backup.` pattern
- Check file permissions

**Backup creation fails**
- Ensure write permissions
- Check disk space
- Verify directory exists

### Recovery

If configuration is corrupted:

```bash
# 1. List all backups
ls -la ~/.config/silentcast/*.backup.*

# 2. Restore most recent backup
cp ~/.config/silentcast/spellbook.yml.backup.20240719-143022 \
   ~/.config/silentcast/spellbook.yml

# 3. Validate restoration
silentcast --validate-config
```

## Advanced Usage

### Scripted Configuration

Generate configurations programmatically:

```python
#!/usr/bin/env python3
import yaml
import subprocess

config = {
    'spells': {
        'e': 'editor',
        't': 'terminal'
    },
    'grimoire': {
        'editor': {
            'type': 'app',
            'command': 'code'
        },
        'terminal': {
            'type': 'app',
            'command': 'wt' if os.name == 'nt' else 'terminal'
        }
    }
}

# Import via stdin
proc = subprocess.Popen(
    ['silentcast', '--import-config', '-'],
    stdin=subprocess.PIPE
)
proc.communicate(yaml.dump(config).encode())
```

### Configuration Diff

Compare configurations:

```bash
# Export two configs
silentcast --export-config config1.yml
# Make changes
silentcast --export-config config2.yml

# Compare
diff -u config1.yml config2.yml
```

## See Also

- [Configuration Guide](configuration.md)
- [CLI Reference](cli-reference.md)
- [Troubleshooting Guide](troubleshooting.md)