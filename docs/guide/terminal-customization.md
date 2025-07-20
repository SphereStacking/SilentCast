# Terminal Window Customization

SilentCast allows you to customize the appearance and behavior of terminal windows when executing scripts. This includes controlling window size, position, colors, fonts, and other visual properties.

## Overview

Terminal customization features:
- Window size (columns and rows)
- Window position (X, Y coordinates)
- Font size adjustment
- Color themes and schemes
- Window states (fullscreen, maximized)
- Always-on-top behavior
- Platform-specific features

## Configuration

Add terminal customization options using the `terminal_customization` field in your action configuration:

```yaml
grimoire:
  my_script:
    type: script
    command: echo "Hello World"
    force_terminal: true
    terminal_customization:
      width: 120          # Window width in columns
      height: 40          # Window height in rows
      x: 100              # X position in pixels
      y: 50               # Y position in pixels
      font_size: 14       # Font size in points
      theme: "gruvbox"    # Color theme name
      maximized: false    # Start maximized
      fullscreen: false   # Start fullscreen
```

## Customization Options

### Window Size

Control the terminal window dimensions:

```yaml
terminal_customization:
  width: 120   # Number of columns (characters wide)
  height: 40   # Number of rows (lines tall)
```

**Common sizes:**
- `80x24` - Classic terminal size
- `120x30` - Wide terminal for modern screens
- `100x50` - Tall terminal for logs
- `60x20` - Compact terminal

### Window Position

Set the initial position of the terminal window:

```yaml
terminal_customization:
  x: 100   # X coordinate in pixels from left edge
  y: 50    # Y coordinate in pixels from top edge
```

**Position examples:**
- `x: 0, y: 0` - Top-left corner
- `x: 100, y: 100` - Offset from corner
- Use negative values to position from right/bottom (platform dependent)

### Font Size

Adjust the terminal font size:

```yaml
terminal_customization:
  font_size: 14   # Font size in points
```

**Font size guidelines:**
- `10-12` - Small, for dense information
- `14-16` - Standard, comfortable reading
- `18-20` - Large, for presentations
- `24+` - Extra large, for accessibility

### Colors and Themes

Customize terminal colors:

```yaml
terminal_customization:
  theme: "gruvbox_dark"        # Named color scheme
  background: "#1e1e1e"        # Background color (hex)
  foreground: "#ffffff"        # Text color (hex)
```

**Popular themes:**
- `gruvbox`, `gruvbox_dark`, `gruvbox_light`
- `solarized`, `solarized_dark`, `solarized_light`
- `monokai`, `dracula`, `nord`
- `campbell` (Windows Terminal default)

### Window States

Control how the window appears:

```yaml
terminal_customization:
  maximized: true      # Start maximized
  fullscreen: false    # Start in fullscreen mode
  always_on_top: true  # Keep window on top
```

## Platform Support

Different terminals support different features:

### Windows

**Windows Terminal (wt.exe):**
- ✅ Window size, position
- ✅ Color schemes
- ✅ Maximized, fullscreen
- ❌ Font size, always-on-top

**Command Prompt (cmd.exe):**
- ✅ Window size, position
- ❌ Colors, themes
- ❌ Advanced features

**PowerShell:**
- ✅ Window size, position
- ❌ Colors, themes
- ❌ Advanced features

### macOS

**Terminal.app:**
- ✅ Window size, position
- ✅ Font size, colors
- ❌ Window states

**iTerm2:**
- ✅ All features supported
- ✅ Advanced customization
- ✅ Window management

### Linux

**GNOME Terminal:**
- ✅ Window size (geometry)
- ✅ Maximized, fullscreen
- ✅ Color profiles
- ❌ Precise positioning

**KDE Konsole:**
- ✅ Window size, profiles
- ✅ Color schemes
- ✅ Font settings

**Alacritty:**
- ✅ Window size, position
- ✅ Config-based theming
- ✅ Font configuration

## Examples

### Basic Terminal Customization

```yaml
spells:
  t,b: basic_terminal

grimoire:
  basic_terminal:
    type: script
    command: |
      echo "Basic customized terminal"
      echo "Size: 100x30"
      date
    force_terminal: true
    keep_open: true
    description: "Basic Terminal"
    terminal_customization:
      width: 100
      height: 30
```

### Development Terminal

```yaml
dev_terminal:
  type: script
  command: |
    echo "Development Environment"
    echo "Large terminal for coding"
    pwd
    ls -la
  force_terminal: true
  keep_open: true
  description: "Development Terminal"
  terminal_customization:
    width: 140
    height: 50
    font_size: 14
    theme: "gruvbox_dark"
    x: 200
    y: 100
```

### Monitoring Terminal

```yaml
monitor_terminal:
  type: script
  command: |
    echo "System Monitoring"
    top -n 1 | head -20
  force_terminal: true
  keep_open: true
  description: "System Monitor"
  terminal_customization:
    width: 120
    height: 40
    fullscreen: false
    always_on_top: true
    theme: "monokai"
```

### Presentation Mode

```yaml
presentation_terminal:
  type: script
  command: |
    echo "PRESENTATION MODE"
    echo "Large font for visibility"
    echo "Running demo script..."
  force_terminal: true
  keep_open: true
  description: "Presentation Terminal"
  terminal_customization:
    width: 80
    height: 25
    font_size: 24
    fullscreen: true
    background: "#000000"
    foreground: "#00ff00"
```

### Multi-Monitor Setup

```yaml
# Primary monitor
primary_terminal:
  type: script
  command: echo "Primary monitor terminal"
  force_terminal: true
  terminal_customization:
    width: 100
    height: 30
    x: 100
    y: 100

# Secondary monitor  
secondary_terminal:
  type: script
  command: echo "Secondary monitor terminal"
  force_terminal: true
  terminal_customization:
    width: 80
    height: 25
    x: 1920    # Second monitor offset
    y: 100
```

## Platform-Specific Examples

### Windows Terminal Profiles

```yaml
windows_dev:
  type: script
  command: |
    echo "Windows development terminal"
    dir
  force_terminal: true
  terminal_customization:
    width: 120
    height: 30
    theme: "Campbell Powershell"  # Windows Terminal profile
    maximized: false
```

### macOS Terminal Styles

```yaml
macos_terminal:
  type: script
  command: |
    echo "macOS styled terminal"
    ls -la
  force_terminal: true
  terminal_customization:
    width: 100
    height: 35
    font_size: 14
    background: "#001122"
    foreground: "#ffffff"
```

### Linux Desktop Integration

```yaml
gnome_terminal:
  type: script
  command: |
    echo "GNOME Terminal with profile"
    neofetch
  force_terminal: true
  terminal_customization:
    width: 80
    height: 24
    maximized: false
    theme: "Tango"  # GNOME profile name
```

## Advanced Customization

### Conditional Customization

You can create different customizations for different scenarios:

```yaml
# Small terminal for quick commands
quick_cmd:
  type: script
  command: date
  force_terminal: true
  terminal_customization:
    width: 60
    height: 15
    font_size: 12

# Large terminal for complex tasks
complex_task:
  type: script
  command: |
    echo "Complex development task"
    # Long-running command here
  force_terminal: true
  terminal_customization:
    width: 140
    height: 50
    font_size: 14
    fullscreen: true
```

### Theme-Based Configurations

```yaml
# Dark theme for night work
night_terminal:
  type: script
  command: echo "Night mode terminal"
  force_terminal: true
  terminal_customization:
    background: "#000000"
    foreground: "#00ff00"
    theme: "gruvbox_dark"

# Light theme for day work
day_terminal:
  type: script
  command: echo "Day mode terminal"
  force_terminal: true
  terminal_customization:
    background: "#ffffff"
    foreground: "#000000"
    theme: "solarized_light"
```

## Fallback Behavior

When a terminal doesn't support a feature:

1. **Unsupported features are ignored** - No error occurs
2. **Partial support** - Supported features are applied
3. **Graceful degradation** - Terminal opens with default settings

Example of fallback:
```yaml
# This works on all terminals, with varying support levels
universal_terminal:
  type: script
  command: echo "Universal terminal"
  force_terminal: true
  terminal_customization:
    width: 80        # Supported by most terminals
    height: 24       # Supported by most terminals
    font_size: 14    # May be ignored on some terminals
    theme: "dark"    # May be ignored if not available
```

## Troubleshooting

### Terminal Not Opening

1. **Check terminal availability:**
   ```bash
   # Check if terminal exists
   which gnome-terminal
   which wt.exe
   ```

2. **Verify feature support:**
   - Some terminals ignore certain options
   - Check terminal documentation for supported features

### Customization Not Applied

1. **Feature not supported:**
   - Check if your terminal supports the feature
   - Some options require specific terminal versions

2. **Invalid values:**
   - Ensure positive values for width/height
   - Check color format (hex: #ffffff)
   - Verify theme names exist

3. **Position issues:**
   - Some terminals ignore position on certain platforms
   - Check for multi-monitor offset requirements

### Performance Issues

1. **Large terminals:**
   - Very large terminal sizes may impact performance
   - Consider reasonable limits (width < 200, height < 100)

2. **Font size:**
   - Extremely large fonts may cause issues
   - Test with standard sizes first

## Best Practices

### 1. Use Reasonable Sizes

```yaml
# Good: Standard sizes
terminal_customization:
  width: 120    # Reasonable width
  height: 40    # Reasonable height

# Avoid: Extreme sizes
terminal_customization:
  width: 500    # Too wide
  height: 200   # Too tall
```

### 2. Platform-Appropriate Themes

```yaml
# Good: Standard theme names
terminal_customization:
  theme: "gruvbox"      # Widely supported
  
# Avoid: Platform-specific themes in general configs
terminal_customization:
  theme: "iTerm2 Custom Profile"  # Only works on macOS iTerm2
```

### 3. Test Across Terminals

Always test your configurations with different terminals:
- Windows: cmd, PowerShell, Windows Terminal
- macOS: Terminal.app, iTerm2
- Linux: gnome-terminal, konsole, xterm

### 4. Provide Fallbacks

```yaml
# Good: Works with and without customization support
script_with_fallback:
  type: script
  command: echo "Works everywhere"
  force_terminal: true
  # Customization enhances experience but isn't required
  terminal_customization:
    width: 100
    height: 30
```

## Migration Guide

### From Basic Terminal Usage

**Before:**
```yaml
my_script:
  type: script
  command: echo "Hello"
  terminal: true
```

**After:**
```yaml
my_script:
  type: script
  command: echo "Hello"
  force_terminal: true
  terminal_customization:
    width: 120
    height: 30
    font_size: 14
```

### Adding Customization to Existing Scripts

1. **Start simple:** Add basic size customization
2. **Test thoroughly:** Verify on your target platforms
3. **Add features gradually:** Colors, positioning, etc.
4. **Document requirements:** Note which terminals work best

## Reference

### All Customization Options

```yaml
terminal_customization:
  # Size
  width: 120              # Window width in columns
  height: 40              # Window height in rows
  
  # Position  
  x: 100                  # X position in pixels
  y: 50                   # Y position in pixels
  
  # Appearance
  font_size: 14           # Font size in points
  theme: "gruvbox"        # Color theme name
  background: "#1e1e1e"   # Background color (hex)
  foreground: "#ffffff"   # Text color (hex)
  
  # Behavior
  fullscreen: false       # Start in fullscreen
  maximized: false        # Start maximized
  always_on_top: false    # Keep window on top
```

### Feature Support Matrix

| Feature        | Windows Terminal | cmd.exe | Terminal.app | iTerm2 | GNOME Terminal | Konsole | Alacritty |
|----------------|------------------|---------|--------------|--------|----------------|---------|-----------|
| Window Size    | ✅               | ✅      | ✅           | ✅     | ✅             | ✅      | ✅        |
| Window Position| ✅               | ✅      | ✅           | ✅     | ❌             | ❌      | ✅        |
| Font Size      | ❌               | ❌      | ✅           | ✅     | ✅             | ✅      | ✅        |
| Color Scheme   | ✅               | ❌      | ✅           | ✅     | ✅             | ✅      | ✅        |
| Window State   | ✅               | ❌      | ❌           | ✅     | ✅             | ✅      | ✅        |
| Always On Top  | ❌               | ❌      | ❌           | ✅     | ❌             | ❌      | ❌        |