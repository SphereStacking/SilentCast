# Key Names Reference

This reference lists all valid key names that can be used in SilentCast configurations. Key names are case-insensitive but we recommend using lowercase for consistency.

## Quick Reference

### Most Common Keys

| Key | Config Name | Example Usage |
|-----|-------------|---------------|
| Letters | `a`-`z` | `e: "editor"` |
| Numbers | `0`-`9` | `1: "workspace_1"` |
| Space | `space` | `space: "launcher"` |
| Enter | `enter` | `enter: "confirm"` |
| Escape | `esc` | `esc: "cancel"` |

### Common Modifiers

| Modifier | Config Name | Platform Notes |
|----------|-------------|----------------|
| Alt | `alt` | Option on macOS |
| Control | `ctrl` | All platforms |
| Shift | `shift` | All platforms |
| Command/Windows | `cmd`, `meta`, `win` | OS-specific |

## Detailed Key Reference

### Alphabetic Keys

All letters A-Z can be used. They are case-insensitive in configuration.

```yaml
spells:
  a: "action_a"
  b: "action_b"
  # ... through z
  z: "action_z"
```

### Numeric Keys

#### Top Row Numbers
```yaml
spells:
  1: "workspace_1"
  2: "workspace_2"
  3: "workspace_3"
  4: "workspace_4"
  5: "workspace_5"
  6: "workspace_6"
  7: "workspace_7"
  8: "workspace_8"
  9: "workspace_9"
  0: "workspace_10"
```

#### Numpad Numbers
```yaml
spells:
  numpad0: "calc_0"
  numpad1: "calc_1"
  numpad2: "calc_2"
  numpad3: "calc_3"
  numpad4: "calc_4"
  numpad5: "calc_5"
  numpad6: "calc_6"
  numpad7: "calc_7"
  numpad8: "calc_8"
  numpad9: "calc_9"
```

### Function Keys

```yaml
spells:
  f1: "help"
  f2: "rename"
  f3: "search"
  f4: "close"
  f5: "refresh"
  f6: "focus_address"
  f7: "spell_check"
  f8: "debug"
  f9: "compile"
  f10: "menu"
  f11: "fullscreen"
  f12: "dev_tools"
```

### Navigation Keys

```yaml
spells:
  # Arrow keys
  up: "navigate_up"
  down: "navigate_down"
  left: "navigate_left"
  right: "navigate_right"
  
  # Page navigation
  pageup: "page_up"
  pagedown: "page_down"
  home: "go_home"
  end: "go_end"
  
  # Text navigation
  insert: "insert_mode"
  delete: "delete_forward"
  backspace: "delete_backward"
```

### Special Keys

```yaml
spells:
  # Common special keys
  space: "quick_launcher"
  enter: "confirm"        # Also: return
  tab: "next_field"
  escape: "cancel"        # Also: esc
  
  # Punctuation and symbols
  minus: "zoom_out"       # -
  plus: "zoom_in"         # + (shift+equals)
  equals: "equals"        # =
  comma: "previous"       # ,
  period: "next"          # .
  slash: "search"         # /
  backslash: "split"      # \
  semicolon: "command"    # ;
  apostrophe: "quote"     # '
  grave: "console"        # ` (backtick)
  
  # Brackets
  leftbracket: "decrease"    # [
  rightbracket: "increase"   # ]
  
  # Lock keys
  capslock: "caps_toggle"
  numlock: "num_toggle"
  scrolllock: "scroll_toggle"
```

### Modifier Keys

Modifiers can be used alone or in combination with other keys.

#### Single Modifiers
```yaml
hotkeys:
  prefix: "alt+space"       # Alt + Space
  prefix: "ctrl+shift+p"    # Ctrl + Shift + P
  prefix: "cmd+k"          # Command + K (macOS)

spells:
  "shift+e": "editor_admin"
  "ctrl+t": "new_tab"
  "alt+f4": "close_window"
```

#### Platform-Specific Modifiers
```yaml
# Windows
spells:
  "win+e": "explorer"       # Windows key
  "meta+r": "run_dialog"    # Meta key (same as Windows key)
  
# macOS
spells:
  "cmd+space": "spotlight"  # Command key
  "meta+tab": "app_switch"  # Meta key (same as Command)
  "option+a": "alternate"   # Option key (same as Alt)
```

### Media Keys

Support varies by platform and keyboard.

```yaml
spells:
  # Playback control
  play: "media_play"
  pause: "media_pause"
  playpause: "media_play_pause"
  stop: "media_stop"
  next: "media_next"
  previous: "media_previous"
  
  # Volume control
  volumeup: "volume_increase"
  volumedown: "volume_decrease"
  mute: "volume_mute"
  
  # Other media keys
  eject: "eject_media"
  record: "start_recording"
```

### Numpad Keys

```yaml
spells:
  # Numpad operators
  numpadplus: "add"         # Numpad +
  numpadminus: "subtract"   # Numpad -
  numpadmultiply: "multiply" # Numpad *
  numpaddivide: "divide"    # Numpad /
  numpadenter: "calculate"  # Numpad Enter
  numpaddecimal: "decimal"  # Numpad .
  
  # Numpad navigation (when NumLock is off)
  numpadhome: "nav_home"    # Numpad 7
  numpadup: "nav_up"        # Numpad 8
  numpadpageup: "nav_pgup"  # Numpad 9
  numpadleft: "nav_left"    # Numpad 4
  numpadcenter: "nav_center" # Numpad 5
  numpadright: "nav_right"  # Numpad 6
  numpadend: "nav_end"      # Numpad 1
  numpaddown: "nav_down"    # Numpad 2
  numpadpagedown: "nav_pgdn" # Numpad 3
  numpadinsert: "nav_insert" # Numpad 0
  numpaddelete: "nav_delete" # Numpad .
```

### Platform-Specific Keys

#### Windows
```yaml
spells:
  printscreen: "screenshot"  # Print Screen
  pause: "pause_break"      # Pause/Break
  apps: "context_menu"      # Application/Menu key
```

#### macOS
```yaml
spells:
  fn: "function_modifier"   # Fn key
  clear: "clear_display"    # Clear key (on some keyboards)
  help: "help_key"         # Help key (on older keyboards)
```

  super: "activities"       # Super key (Windows key)
  menu: "app_menu"         # Menu key
  compose: "compose_key"   # Compose key
```

## Key Combination Syntax

### Modifier + Key
```yaml
spells:
  # Single modifier
  "ctrl+s": "save"
  "alt+f": "file_menu"
  "shift+tab": "previous_field"
  
  # Multiple modifiers
  "ctrl+shift+s": "save_as"
  "ctrl+alt+delete": "task_manager"
  "cmd+shift+4": "screenshot_area"  # macOS
```

### Sequential Keys (Multi-key)
```yaml
spells:
  # Two-key sequences
  "g,s": "git_status"
  "g,p": "git_pull"
  
  # Three-key sequences
  "w,i,n": "window_management"
  "d,e,v": "development_mode"
  
  # Numbers in sequences
  "p,1": "project_1"
  "p,2": "project_2"
```

## Limitations and Notes

### Case Sensitivity
- Key names are **case-insensitive**: `"a"`, `"A"`, and `"A"` are treated the same
- We recommend using lowercase for consistency

### Invalid Combinations
Some combinations may not work due to OS limitations:
- `ctrl+alt+del` on Windows (system reserved)
- `cmd+tab` on macOS (system app switcher)
- `alt+tab` on Windows (window switcher)

### Sequence Limitations
- Sequences cannot include modifier keys: ❌ `"ctrl,s"`
- Use modified keys in single shortcuts only: ✅ `"ctrl+s"`

### Platform Differences
- **Windows**: `win` key might be intercepted by Windows
- **macOS**: Some `cmd` combinations are reserved by the system

## Testing Keys

Use the test hotkey feature to verify key names:

```bash
silentcast --test-hotkey
```

This will show you the exact key name SilentCast detects when you press a key.

## Examples

### Developer Layout
```yaml
spells:
  # Editor controls
  e: "editor"
  "shift+e": "editor_admin"
  
  # Git shortcuts
  "g,s": "git_status"
  "g,c": "git_commit"
  "g,p": "git_push"
  
  # Build shortcuts
  "b,b": "build"
  "b,t": "build_test"
  "b,c": "build_clean"
  
  # Function keys
  f5: "run"
  f9: "debug"
  f10: "step_over"
  f11: "step_into"
```

### System Admin Layout
```yaml
spells:
  # System monitoring
  "s,t": "system_top"
  "s,m": "system_monitor"
  "s,l": "system_logs"
  
  # Service control
  "r,n": "restart_nginx"
  "r,d": "restart_database"
  "r,a": "restart_all"
  
  # Quick access
  1: "server_1"
  2: "server_2"
  3: "server_3"
```

### Creative Workflow
```yaml
spells:
  # Tool selection
  b: "brush"
  e: "eraser"
  i: "eyedropper"
  t: "text"
  
  # File operations
  "ctrl+s": "save"
  "ctrl+shift+s": "save_as"
  "ctrl+e": "export"
  
  # View controls
  space: "pan"
  z: "zoom"
  "shift+z": "zoom_out"
```

## Troubleshooting

### Key Not Recognized

1. Test the key:
   ```bash
   silentcast --test-hotkey
   ```

2. Check platform-specific names (e.g., `cmd` vs `win`)

3. Verify the key isn't reserved by the OS

4. Try alternative names (e.g., `return` vs `enter`)

### Modifier Issues

- Ensure modifiers come before the key: `ctrl+s` not `s+ctrl`
- Use `+` to separate modifiers and keys
- Don't use modifiers in sequences

### Special Characters

For special characters, use their descriptive names:
- Use `minus` not `-`
- Use `plus` not `+`
- Use `slash` not `/`

## See Also

- [Spells Guide](/guide/spells) - How to create effective spells
- [Configuration Reference](/config/) - Full configuration options
- [Platform Guide](/guide/platforms) - Platform-specific considerations