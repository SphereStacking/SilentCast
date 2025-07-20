# Force Terminal Execution

The `force_terminal` option allows scripts to be executed in a terminal window even when SilentCast is running in GUI or tray mode. This is particularly useful for interactive scripts, installer scripts, or scripts that display important output that users need to see.

## Overview

Force terminal execution:
- Opens a terminal window regardless of SilentCast's execution mode
- Enables user interaction with scripts
- Supports all existing terminal options (`keep_open`, `show_output`)
- Works with output capture for notifications
- Maintains compatibility with existing terminal functionality

## When to Use Force Terminal

Use `force_terminal: true` when:
- Scripts require user input (interactive scripts)
- Scripts display important information users must see
- Running installation or setup scripts
- Scripts need a terminal environment for proper display
- Users need to monitor script progress in real-time

**Examples of good candidates:**
- Package installers
- Interactive configuration scripts
- System maintenance scripts
- Scripts that prompt for passwords or confirmations
- Long-running processes with progress indicators

## Configuration

### Basic Force Terminal

```yaml
grimoire:
  interactive_setup:
    type: script
    command: |
      echo "Welcome to the setup wizard"
      echo "Please enter your configuration..."
      read -p "Enter your name: " name
      echo "Hello, $name!"
    force_terminal: true
    description: "Interactive Setup"
```

### Force Terminal with Output Capture

Combine `force_terminal` with `show_output` to both display in terminal AND capture output for notifications:

```yaml
grimoire:
  system_check:
    type: script
    command: |
      echo "Running system diagnostics..."
      df -h
      free -h
      uptime
      echo "Diagnostics complete"
    force_terminal: true
    show_output: true
    description: "System Diagnostics"
```

### Force Terminal with Keep Open

Keep the terminal open after completion for user review:

```yaml
grimoire:
  log_analysis:
    type: script
    command: |
      echo "=== Recent Error Logs ==="
      tail -20 /var/log/syslog | grep -i error
      echo ""
      echo "Review complete. Press any key to close."
    force_terminal: true
    keep_open: true
    description: "Error Log Analysis"
```

## Comparison with Other Terminal Options

### Regular Terminal Mode (`terminal: true`)

```yaml
# Only opens terminal when SilentCast runs in console mode
terminal_script:
  type: script
  command: "echo 'May or may not show in terminal'"
  terminal: true
```

### Force Terminal Mode (`force_terminal: true`)

```yaml
# Always opens terminal, regardless of SilentCast execution mode
force_terminal_script:
  type: script
  command: "echo 'Always shows in terminal'"
  force_terminal: true
```

### Combined Options

```yaml
# All terminal-related options can be combined
comprehensive_terminal:
  type: script
  command: "echo 'Full terminal features'"
  terminal: true        # Request terminal when appropriate
  force_terminal: true  # Force terminal even in GUI mode
  keep_open: true       # Keep terminal open after completion
  show_output: true     # Also capture output for notifications
```

## Output Handling

### Force Terminal Only

```yaml
# Output goes to terminal only
terminal_only:
  type: script
  command: "echo 'Terminal output only'"
  force_terminal: true
```

### Force Terminal + Notifications

```yaml
# Output goes to terminal AND notifications
terminal_and_notifications:
  type: script
  command: "echo 'Both terminal and notification'"
  force_terminal: true
  show_output: true
```

The system automatically handles output routing:
- `force_terminal: true` alone → Terminal display only
- `force_terminal: true` + `show_output: true` → Terminal display + notification capture

## Advanced Examples

### Interactive Package Manager

```yaml
spells:
  p,i: package_install
  p,u: package_update
  p,s: package_search

grimoire:
  package_install:
    type: script
    command: |
      echo "Package Installation Wizard"
      echo "=========================="
      echo ""
      read -p "Enter package name: " package
      if [ -z "$package" ]; then
        echo "Error: Package name required"
        exit 1
      fi
      
      echo "Installing $package..."
      sudo apt update && sudo apt install "$package"
      
      if [ $? -eq 0 ]; then
        echo "✓ Package $package installed successfully"
      else
        echo "✗ Failed to install $package"
      fi
      
      echo ""
      echo "Press any key to close..."
      read -n 1
    force_terminal: true
    keep_open: true
    description: "Interactive Package Installer"

  package_update:
    type: script
    command: |
      echo "System Update"
      echo "============="
      sudo apt update
      echo ""
      echo "Available updates:"
      apt list --upgradable
      echo ""
      read -p "Proceed with update? (y/N): " confirm
      if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
        sudo apt upgrade
        echo "✓ System updated successfully"
      else
        echo "Update cancelled"
      fi
    force_terminal: true
    keep_open: true
    description: "System Update"

  package_search:
    type: script
    command: |
      read -p "Search for package: " query
      if [ -n "$query" ]; then
        apt search "$query"
      else
        echo "No search query provided"
      fi
    force_terminal: true
    keep_open: true
    description: "Package Search"
```

### Development Environment Setup

```yaml
grimoire:
  dev_setup:
    type: script
    command: |
      echo "Development Environment Setup"
      echo "============================="
      echo ""
      
      # Check requirements
      echo "Checking requirements..."
      command -v git >/dev/null || echo "⚠ Git not found"
      command -v node >/dev/null || echo "⚠ Node.js not found"
      command -v python3 >/dev/null || echo "⚠ Python3 not found"
      
      echo ""
      read -p "Enter project name: " project
      read -p "Enter git repository URL (optional): " repo
      
      if [ -n "$project" ]; then
        mkdir -p ~/Projects/"$project"
        cd ~/Projects/"$project"
        
        if [ -n "$repo" ]; then
          git clone "$repo" .
        else
          git init
          echo "# $project" > README.md
          git add README.md
          git commit -m "Initial commit"
        fi
        
        echo "✓ Project $project created in ~/Projects/$project"
        echo "Opening in default editor..."
        code . 2>/dev/null || echo "VS Code not found, skipping editor launch"
      fi
    force_terminal: true
    show_output: true
    description: "Development Environment Setup"
```

### System Maintenance Script

```yaml
grimoire:
  system_maintenance:
    type: script
    command: |
      echo "System Maintenance Script"
      echo "========================"
      echo ""
      
      echo "1. Cleaning package cache..."
      sudo apt autoremove -y && sudo apt autoclean
      
      echo ""
      echo "2. Updating locate database..."
      sudo updatedb
      
      echo ""
      echo "3. Checking disk usage..."
      df -h | head -1
      df -h | grep -E "/$|/home"
      
      echo ""
      echo "4. Memory usage:"
      free -h
      
      echo ""
      echo "5. System uptime:"
      uptime
      
      echo ""
      echo "✓ Maintenance complete!"
      echo "Press any key to close..."
      read -n 1
    force_terminal: true
    keep_open: true
    timeout: 300  # 5 minutes timeout
    description: "System Maintenance"
```

## Best Practices

### 1. Use for Interactive Scripts

```yaml
# Good: Interactive script that needs user input
interactive_config:
  type: script
  command: |
    read -p "Enter configuration value: " value
    echo "Configuration saved: $value"
  force_terminal: true
```

### 2. Combine with Appropriate Options

```yaml
# Good: Combine force_terminal with relevant options
installer_script:
  type: script
  command: "sudo ./install.sh"
  force_terminal: true
  keep_open: true      # Review installation results
  show_output: true    # Get notification when done
  timeout: 600         # 10 minute timeout for long installs
```

### 3. Provide Clear Descriptions

```yaml
# Good: Clear description helps users understand what will happen
database_backup:
  type: script
  command: "./backup-database.sh"
  force_terminal: true
  description: "Interactive Database Backup (requires password)"
```

### 4. Handle Errors Gracefully

```yaml
# Good: Script handles errors and provides feedback
network_diagnostic:
  type: script
  command: |
    echo "Network Diagnostic Tool"
    echo "====================="
    
    if ! command -v ping >/dev/null; then
      echo "Error: ping command not found"
      exit 1
    fi
    
    read -p "Enter hostname or IP to test: " target
    if [ -z "$target" ]; then
      echo "Error: No target specified"
      exit 1
    fi
    
    echo "Testing connection to $target..."
    ping -c 4 "$target"
    
    if [ $? -eq 0 ]; then
      echo "✓ Connection successful"
    else
      echo "✗ Connection failed"
    fi
  force_terminal: true
  keep_open: true
  description: "Network Connectivity Test"
```

## Troubleshooting

### Terminal Not Opening

If terminal doesn't open:

1. **Check terminal availability:**
   ```bash
   # Check if terminal emulators are installed
   which gnome-terminal xterm konsole
   ```

2. **Verify display environment:**
   ```bash
   echo $DISPLAY
   ```

3. **Test in console mode:**
   ```bash
   # Run SilentCast in console mode to test
   ./silentcast --no-tray
   ```

### Output Not Captured

If output capture isn't working with force terminal:

```yaml
# Ensure show_output is enabled
script_with_capture:
  type: script
  command: "echo 'This will be captured'"
  force_terminal: true
  show_output: true  # Required for output capture
```

### Terminal Closes Too Quickly

Use `keep_open` to prevent immediate closure:

```yaml
# Terminal stays open for review
quick_script:
  type: script
  command: "echo 'Quick output'"
  force_terminal: true
  keep_open: true
```

## Platform Differences

### Linux

- Uses default terminal emulator (gnome-terminal, xterm, etc.)
- Supports all terminal features
- Requires X11/Wayland display for GUI mode

### macOS

- Uses Terminal.app or iTerm2
- Full feature support
- Works in GUI mode

### Windows

- Uses Windows Terminal, PowerShell, or Command Prompt
- Feature support varies by terminal
- Some advanced features may be limited

## Migration from Regular Terminal

To migrate existing scripts to force terminal:

```yaml
# Before: Regular terminal mode
old_script:
  type: script
  command: "interactive-script.sh"
  terminal: true

# After: Force terminal mode
new_script:
  type: script
  command: "interactive-script.sh"
  force_terminal: true  # Always opens terminal
  terminal: true        # Keep for backward compatibility
```

The `force_terminal` option is additive - it works alongside existing terminal options for maximum compatibility.