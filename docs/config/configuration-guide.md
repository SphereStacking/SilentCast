# SilentCast Configuration Guide 🤫

## Overview

SilentCast uses a configuration system based on YAML files. The configuration follows a cascading model where OS-specific settings override common settings.

## Configuration Files

1. **`spellbook.yml`** - Common configuration loaded on all platforms
2. **`spellbook.mac.yml`** - macOS-specific overrides
3. **`spellbook.windows.yml`** - Windows-specific overrides
4. **`spellbook.linux.yml`** - Linux-specific overrides

## Configuration Structure

### Daemon Settings

```yaml
daemon:
  auto_start: false      # Start with system
  log_level: info       # debug, info, warn, error
  config_watch: true    # Auto-reload on changes
```

### Hotkey Settings

```yaml
hotkeys:
  prefix: "alt+space"        # Activation key combo
  timeout: 1000             # MS after prefix key
  sequence_timeout: 2000    # MS for full sequence
```

### Spells (Shortcuts)

```yaml
spells:
  # Single key
  e: "editor"
  
  # Multi-key sequence
  "g,s": "git_status"
```

### Grimoire (Actions)

```yaml
grimoire:
  editor:
    type: app                    # app or script
    command: "/path/to/app"      # Path or command
    args: ["-n", "--new-window"] # Optional
    env:                         # Optional
      TERM: "xterm-256color"
    working_dir: "~/projects"    # Optional
    description: "Text editor"   # Optional
```

## Key Notation

### Modifier Keys
- **Windows**: `win`, `ctrl`, `alt`, `shift`
- **macOS**: `cmd`, `ctrl`, `opt`, `shift`
- **Linux**: `super`, `ctrl`, `alt`, `shift`

### Special Keys
- Function keys: `f1`-`f12`
- Navigation: `up`, `down`, `left`, `right`
- Other: `space`, `enter`, `tab`, `esc`, `backspace`

### Examples
- `"alt+space"` - Alt and Space together
- `"cmd+shift+p"` - Command+Shift+P (macOS)
- `"g,s"` - G followed by S (sequential)
- `"editor"` - Full word sequence

## Action Types

### App Type
Launches applications with optional arguments:

```yaml
vscode:
  type: app
  command: "/usr/bin/code"
  args: ["-n", "~/projects"]
  description: "Open VS Code in projects folder"
```

### Script Type
Executes shell commands:

```yaml
git_status:
  type: script
  command: "git status --short"
  working_dir: "."
  description: "Show git status"
```

## Tips

1. **Test First**: Test your spells with simple actions before complex ones
2. **Use Descriptions**: Help your future self remember what spells do
3. **OS Paths**: Use full paths for applications to avoid PATH issues
4. **Environment Variables**: Supported in `command` with `$VAR` or `${VAR}`
5. **Escape Sequences**: Use quotes for special characters in YAML

## Example Custom Spell

Add to your OS-specific config:

```yaml
# Custom development setup
spells:
  "d,e,v": "dev_setup"

grimoire:
  dev_setup:
    type: script
    command: |
      cd ~/projects &&
      code . &&
      gnome-terminal --tab --title="Server" -- npm start &&
      firefox http://localhost:3000
    description: "Full dev environment setup"
```

## Debugging

If a spell doesn't work:

1. Check the logs: `~/.silentcast/spellbook.log`
2. Verify the path exists: `which <command>` or `ls <path>`
3. Test the command manually in terminal
4. Check for typos in spell names
5. Ensure proper YAML formatting (use spaces, not tabs)

## Magic Words 🎭

Remember: You're not just configuring software, you're crafting spells! Have fun with your spell names and descriptions.