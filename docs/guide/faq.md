# Frequently Asked Questions

## General Questions

### What is SilentCast?

SilentCast is a lightweight, cross-platform hotkey task runner that lets you execute commands, launch applications, and open URLs using customizable keyboard shortcuts. It runs silently in the background and responds instantly to your configured key combinations.

### How is it different from other hotkey tools?

- **Developer-focused**: Designed specifically for developer workflows
- **Cross-platform**: Works on Windows and macOS with the same configuration
- **YAML configuration**: Simple, version-controllable configuration
- **Lightweight**: Minimal resource usage (~15MB RAM, <1% CPU when idle)
- **VS Code-style sequences**: Support for multi-key combinations like `g,s` for git status

### Is SilentCast free?

Yes, SilentCast is completely free and open source under the MIT license.

## Installation Questions

### Where should I install SilentCast?

- **Windows**: `C:\Program Files\SilentCast\` or anywhere in your PATH
- **macOS**: `/usr/local/bin/` or `/Applications/`
- **Linux**: `/usr/local/bin/` or `~/.local/bin/`

### Do I need administrator privileges?

- **Installation**: Admin rights recommended for system-wide installation
- **Running**: No admin rights needed, except for specific actions that require elevation
- **macOS**: Accessibility permissions required (not admin)

### How do I update SilentCast?

```bash
# If installed via package manager
brew upgrade silentcast  # macOS
scoop update silentcast  # Windows

# Manual update
# Download new version and replace the binary
```

## Configuration Questions

### Where is the configuration file located?

Default locations:
- **Windows**: `%APPDATA%\SilentCast\spellbook.yml`
- **macOS**: `~/.config/silentcast/spellbook.yml`
- **Linux**: `~/.config/silentcast/spellbook.yml`

You can also specify a custom location:
```bash
silentcast --config /path/to/spellbook.yml
```

### Can I use the same configuration across platforms?

Yes! Use the base `spellbook.yml` for common settings and platform-specific files for overrides:
- `spellbook.windows.yml`
- `spellbook.mac.yml`
- `spellbook.linux.yml`

### How do I share configurations between machines?

1. Store your configuration in a Git repository
2. Use symlinks to link to the repository:
   ```bash
   ln -s ~/dotfiles/silentcast/spellbook.yml ~/.config/silentcast/spellbook.yml
   ```

### Can I use environment variables in configuration?

Yes, environment variables are supported:
```yaml
grimoire:
  editor:
    command: "${EDITOR:-code}"
    working_dir: "${PROJECT_DIR}"
```

## Usage Questions

### Why aren't my hotkeys working?

Common causes:
1. **macOS**: Missing accessibility permissions
2. **Hotkey conflicts**: Another app is using the same combination
3. **Wrong format**: Use lowercase (e.g., `alt+space`, not `Alt+Space`)
4. **Not running**: Check if SilentCast is actually running

### Can I use the same hotkey for different actions?

No, each spell (hotkey combination) can only map to one action. However, you can:
- Use modifier keys for variations: `e` vs `shift+e`
- Use sequences for related actions: `g,s` for status, `g,c` for commit

### How do I disable SilentCast temporarily?

Options:
1. **Quit from system tray**: Right-click tray icon → Quit
2. **Kill the process**: `pkill silentcast` or Task Manager
3. **Disable auto-start**: Set `daemon.auto_start: false`

### Can I run multiple instances?

No, SilentCast prevents multiple instances to avoid conflicts. If you need different configurations, use different config files and run them separately.

## Hotkey Questions

### What key combinations are supported?

Modifiers:
- `ctrl`, `alt`, `shift`, `cmd` (macOS), `win` (Windows)

Keys:
- Letters: `a-z`
- Numbers: `0-9`
- Function keys: `f1-f24`
- Special: `space`, `enter`, `tab`, `escape`, `backspace`, `delete`
- Navigation: `up`, `down`, `left`, `right`, `home`, `end`, `pageup`, `pagedown`

### Can I use mouse buttons?

Not currently. SilentCast focuses on keyboard shortcuts only.

### How do sequences work?

Sequences are comma-separated keys pressed in order:
```yaml
spells:
  "g,s": "git_status"    # Press g, then s
  "d,o,c": "open_docs"   # Press d, then o, then c
```

You have a limited time to complete the sequence (configurable via `sequence_timeout`).

## Action Questions

### How do I run commands with arguments?

Use the `args` field:
```yaml
grimoire:
  vscode_project:
    type: app
    command: "code"
    args: ["~/projects/myproject", "--new-window"]
```

### Can I run commands as administrator/root?

Yes, use the `admin` field:
```yaml
grimoire:
  admin_task:
    type: app
    command: "systemctl"
    args: ["restart", "nginx"]
    admin: true  # Will prompt for elevation
```

### How do I set environment variables for actions?

```yaml
grimoire:
  dev_server:
    type: script
    command: "npm start"
    env:
      NODE_ENV: "development"
      PORT: "3000"
```

### Can I chain multiple actions?

Not directly, but you can create a script that runs multiple commands:
```yaml
grimoire:
  deploy:
    type: script
    command: |
      git pull &&
      npm install &&
      npm run build &&
      pm2 restart app
```

## Troubleshooting Questions

### How do I debug issues?

1. Enable debug logging:
   ```yaml
   daemon:
     log_level: debug
   logger:
     level: debug
   ```

2. Run in foreground:
   ```bash
   silentcast --no-tray --log-level debug
   ```

3. Check logs:
   ```bash
   tail -f ~/.local/share/silentcast/silentcast.log
   ```

### Why is SilentCast using high CPU?

Possible causes:
- Config file watching on network drive
- Too frequent update checks
- Logging at debug level

Solutions:
- Disable config watching: `config_watch: false`
- Increase update interval
- Set log level to `info` or `warn`

### How do I reset SilentCast?

1. Stop SilentCast
2. Delete configuration:
   ```bash
   # Backup first!
   rm -rf ~/.config/silentcast
   rm -rf ~/.local/share/silentcast
   ```
3. Reinstall or create new configuration

## Advanced Questions

### Can I integrate SilentCast with my IDE?

Yes, you can trigger SilentCast actions from your IDE:
```bash
# VS Code task
{
  "label": "Deploy",
  "type": "shell",
  "command": "silentcast --execute deploy"
}
```

### Can I control SilentCast programmatically?

Currently limited to:
- `--execute <spell>`: Run a specific spell
- `--validate-config`: Check configuration
- Process signals for reload/quit

### Is there an API?

Not yet, but it's on the roadmap. Currently, you can:
- Use command-line flags
- Modify config files and reload
- Send process signals

### Can I create custom action types?

Not directly, but you can:
1. Use script actions for custom behavior
2. Create wrapper scripts
3. Contribute new action types to the project

## Security Questions

### Is it safe to use?

Yes, SilentCast:
- Runs with user privileges only
- No network access (except update checks)
- No data collection or telemetry
- Open source for audit

### How are sensitive values handled?

Best practices:
1. Use environment variables for secrets
2. Set appropriate file permissions on config
3. Don't commit secrets to version control

### Can malicious configs harm my system?

SilentCast executes whatever commands you configure. Always:
- Review configurations before using
- Be cautious with scripts from others
- Use absolute paths to avoid PATH hijacking

## Still Have Questions?

- Check the [documentation](/guide/getting-started)
- Search [GitHub issues](https://github.com/SphereStacking/silentcast/issues)
- Open a [new issue](https://github.com/SphereStacking/silentcast/issues/new)

---

*Last updated: January 2025*