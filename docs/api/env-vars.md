# Environment Variables

SilentCast uses environment variables for configuration and provides built-in variables for use in your spellbook.

## Configuration Variables

These environment variables control SilentCast's behavior:

### `SILENTCAST_CONFIG`

Sets the configuration directory path.

```bash
# Linux/macOS
export SILENTCAST_CONFIG=/home/user/my-config

# Windows
set SILENTCAST_CONFIG=C:\Users\user\my-config
```

Default locations:
- Linux: `~/.config/silentcast`
- macOS: `~/Library/Application Support/silentcast`
- Windows: `%APPDATA%\SilentCast`

### `SILENTCAST_LOG_LEVEL`

Override the log level without modifying configuration.

```bash
export SILENTCAST_LOG_LEVEL=debug
```

Values: `debug`, `info`, `warn`, `error`

### `SILENTCAST_NO_TRAY`

Disable system tray integration.

```bash
export SILENTCAST_NO_TRAY=1
```

### `SILENTCAST_DRY_RUN`

Run in test mode without executing actions.

```bash
export SILENTCAST_DRY_RUN=1
```

## Built-in Variables

These variables are available for use in your configuration:

### System Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `${HOME}` | User home directory | `/home/user` |
| `${USER}` | Current username | `john` |
| `${HOSTNAME}` | Machine hostname | `desktop-01` |
| `${PWD}` | Current working directory | `/home/user/projects` |
| `${TMPDIR}` | Temporary directory | `/tmp` |

### SilentCast Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `${SILENTCAST_VERSION}` | Current version | `1.0.0` |
| `${SILENTCAST_CONFIG}` | Config directory | `~/.config/silentcast` |
| `${SILENTCAST_DATA}` | Data directory | `~/.local/share/silentcast` |
| `${SILENTCAST_CACHE}` | Cache directory | `~/.cache/silentcast` |
| `${SILENTCAST_SPELL}` | Current spell being cast | `git_status` |

## Using Variables in Configuration

### In Commands

```yaml
grimoire:
  open_home:
    type: app
    command: "${FILE_MANAGER:-explorer}"
    args: ["${HOME}"]
    
  edit_config:
    type: app
    command: "${EDITOR:-notepad}"
    args: ["${SILENTCAST_CONFIG}/spellbook.yml"]
```

### In Working Directories

```yaml
grimoire:
  git_status:
    type: script
    command: "git status"
    working_dir: "${PWD}"  # Use current directory
    
  build_project:
    type: script
    command: "make build"
    working_dir: "${PROJECT_ROOT:-${HOME}/projects}"
```

### In Environment Variables

```yaml
grimoire:
  development_server:
    type: script
    command: "npm start"
    env:
      NODE_ENV: "${NODE_ENV:-development}"
      PORT: "${PORT:-3000}"
      API_URL: "${API_URL:-http://localhost:8080}"
```

## Variable Expansion

### Default Values

Use `:-` to provide defaults:

```yaml
command: "${EDITOR:-code}"  # Use $EDITOR or 'code' if not set
```

### Required Variables

Use `:?` to fail if not set:

```yaml
env:
  API_KEY: "${API_KEY:?API_KEY must be set}"
```

### Nested Expansion

Variables can reference other variables:

```yaml
env:
  PROJECT_DIR: "${HOME}/projects"
  BUILD_DIR: "${PROJECT_DIR}/build"
```

## Platform-Specific Variables

### Windows

| Variable | Description |
|----------|-------------|
| `%APPDATA%` | Application data folder |
| `%LOCALAPPDATA%` | Local application data |
| `%PROGRAMFILES%` | Program Files directory |
| `%SYSTEMROOT%` | Windows directory |
| `%TEMP%` | Temporary files directory |

### macOS

| Variable | Description |
|----------|-------------|
| `$HOME/Library` | User library directory |
| `$HOME/Applications` | User applications |
| `/Applications` | System applications |

### Linux

| Variable | Description |
|----------|-------------|
| `$XDG_CONFIG_HOME` | User config directory |
| `$XDG_DATA_HOME` | User data directory |
| `$XDG_CACHE_HOME` | User cache directory |

## Examples

### Development Environment

```yaml
grimoire:
  start_dev:
    type: script
    command: "${HOME}/scripts/start-dev.sh"
    env:
      NODE_ENV: "development"
      DATABASE_URL: "${DATABASE_URL:-postgresql://localhost/dev}"
      REDIS_URL: "${REDIS_URL:-redis://localhost:6379}"
      LOG_LEVEL: "${LOG_LEVEL:-debug}"
```

### Dynamic Paths

```yaml
grimoire:
  open_project:
    type: app
    command: "${EDITOR}"
    args: ["${CURRENT_PROJECT:-${HOME}/projects/default}"]
    
  backup_config:
    type: script
    command: "cp -r ${SILENTCAST_CONFIG} ${BACKUP_DIR:-${HOME}/backups}/silentcast-$(date +%Y%m%d)"
```

### Conditional Execution

```yaml
grimoire:
  deploy:
    type: script
    command: |
      if [ "${ENVIRONMENT}" = "production" ]; then
        ${HOME}/scripts/deploy-prod.sh
      else
        ${HOME}/scripts/deploy-dev.sh
      fi
    env:
      ENVIRONMENT: "${ENVIRONMENT:-development}"
```

## Best Practices

1. **Use defaults** - Always provide sensible defaults with `:-`
2. **Document requirements** - List required variables in comments
3. **Validate early** - Use `:?` for critical variables
4. **Cross-platform paths** - Use forward slashes or variables
5. **Avoid hardcoding** - Use variables for paths and values that might change

## Debugging Variables

View expanded configuration:

```bash
# Show all environment variables
silentcast --show-env

# Test variable expansion
silentcast --dry-run --log-level debug

# Export for testing
export SILENTCAST_DEBUG=1
export API_KEY=test-key
silentcast --validate-config
```

## Security Considerations

1. **Never commit secrets** - Keep sensitive variables in environment
2. **Use secure storage** - Consider using password managers for API keys
3. **Limit scope** - Only set variables where needed
4. **Rotate regularly** - Change sensitive values periodically

## Next Steps

- [Configuration Guide](/guide/configuration) - Full configuration reference
- [Scripts Guide](/guide/scripts) - Using variables in scripts
- [Platform Support](/guide/platforms) - Platform-specific variables