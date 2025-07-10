# Configuration Samples

This page provides complete, working configuration examples for SilentCast. You can copy these files directly to your configuration directory and customize them for your needs.

## Base Configuration

This is a comprehensive base configuration that works across all platforms. Save this as `spellbook.yml` in your configuration directory.

::: details spellbook.yml - Click to expand
```yaml
# SilentCast Configuration Example
# This is the main configuration file that works across all platforms
# Platform-specific settings should go in spellbook.{platform}.yml files

# Daemon settings control how SilentCast runs in the background
daemon:
  # Automatically start when system boots
  auto_start: false
  
  # Logging verbosity: debug, info, warn, error
  log_level: info
  
  # Watch configuration files for changes and auto-reload
  config_watch: true
  
  # Check for updates automatically
  check_updates: true
  
  # How often to check for updates (in hours)
  update_interval: 24

# Logger configuration for detailed logging and debugging
logger:
  # Log level: debug, info, warn, error
  level: info
  
  # Log file location (leave empty for console only)
  # Use ~ for home directory, supports environment variables
  file: "~/.local/share/silentcast/silentcast.log"
  
  # Maximum log file size in MB before rotation
  max_size: 10
  
  # Number of old log files to keep
  max_backups: 3
  
  # Maximum age of log files in days
  max_age: 7
  
  # Compress old log files to save space
  compress: true

# Hotkey settings define how keyboard shortcuts work
hotkeys:
  # The magic activation key - press this first, then your spell
  # Common options: "alt+space", "ctrl+shift+space", "cmd+k"
  prefix: "alt+space"
  
  # Time in milliseconds to wait for spell after pressing prefix
  timeout: 1000
  
  # Total time for multi-key sequences (like "g,s" for git status)
  sequence_timeout: 2000
  
  # Show notification when prefix is pressed
  show_notification: true
  
  # Play sound on successful spell cast
  play_sound: false

# Spells are your keyboard shortcuts
# Format: "keys": "action_name"
# Single keys: "e", "t", "b"
# Sequences: "g,s", "g,p", "d,b"
spells:
  # === Quick Launch ===
  e: "editor"              # Alt+Space, E -> Launch code editor
  t: "terminal"            # Alt+Space, T -> Open terminal
  b: "browser"             # Alt+Space, B -> Web browser
  f: "file_manager"        # Alt+Space, F -> File explorer
  s: "settings"            # Alt+Space, S -> System settings
  
  # === Development Tools ===
  v: "vscode"              # Alt+Space, V -> VS Code
  i: "ide"                 # Alt+Space, I -> Your IDE
  d: "docker"              # Alt+Space, D -> Docker Desktop
  p: "postman"             # Alt+Space, P -> API testing
  
  # === Git Workflows ===
  "g,s": "git_status"      # Alt+Space, G, S -> git status
  "g,a": "git_add_all"     # Alt+Space, G, A -> git add -A
  "g,c": "git_commit"      # Alt+Space, G, C -> git commit
  "g,p": "git_push"        # Alt+Space, G, P -> git push
  "g,l": "git_pull"        # Alt+Space, G, L -> git pull
  "g,d": "git_diff"        # Alt+Space, G, D -> git diff
  "g,b": "git_branch"      # Alt+Space, G, B -> git branch list
  "g,h": "git_history"     # Alt+Space, G, H -> git log
  
  # === Docker Commands ===
  "d,u": "docker_up"       # Alt+Space, D, U -> docker-compose up
  "d,d": "docker_down"     # Alt+Space, D, D -> docker-compose down
  "d,l": "docker_logs"     # Alt+Space, D, L -> docker logs
  "d,p": "docker_ps"       # Alt+Space, D, P -> docker ps
  "d,i": "docker_images"   # Alt+Space, D, I -> docker images
  "d,c": "docker_clean"    # Alt+Space, D, C -> docker cleanup
  
  # === Project Shortcuts ===
  "p,b": "project_build"   # Alt+Space, P, B -> Build project
  "p,t": "project_test"    # Alt+Space, P, T -> Run tests
  "p,r": "project_run"     # Alt+Space, P, R -> Run project
  "p,c": "project_clean"   # Alt+Space, P, C -> Clean build
  
  # === System Control ===
  "s,l": "system_lock"     # Alt+Space, S, L -> Lock screen
  "s,s": "system_sleep"    # Alt+Space, S, S -> Sleep mode
  "s,r": "system_restart"  # Alt+Space, S, R -> Restart
  
  # === Quick Access ===
  "q,n": "quick_note"      # Alt+Space, Q, N -> Quick note
  "q,t": "quick_todo"      # Alt+Space, Q, T -> Todo list
  "q,c": "quick_calc"      # Alt+Space, Q, C -> Calculator
  "q,s": "quick_search"    # Alt+Space, Q, S -> Search

# Grimoire contains the actual commands/applications to run
# Each entry has a type: app, script, or url
grimoire:
  # === Applications ===
  editor:
    type: app
    command: "code"  # Will be overridden per platform
    description: "Open code editor"
    args: ["--new-window"]
    
  terminal:
    type: app
    command: "wt"  # Windows Terminal, override for other platforms
    description: "Open terminal emulator"
    
  browser:
    type: app
    command: "chrome"
    description: "Open web browser"
    args: ["--new-window"]
    
  file_manager:
    type: app
    command: "explorer"  # Windows Explorer, override for other platforms
    description: "Open file manager"
    
  vscode:
    type: app
    command: "code"
    description: "Visual Studio Code"
    args: ["--new-window", "${CURRENT_PROJECT:-~/projects}"]
    
  # === Git Commands ===
  git_status:
    type: script
    command: "git status"
    description: "Show git repository status"
    working_dir: "${PWD}"
    show_output: true
    
  git_add_all:
    type: script
    command: "git add -A"
    description: "Stage all changes"
    working_dir: "${PWD}"
    
  git_commit:
    type: script
    command: |
      echo "Enter commit message:"
      read -r message
      git commit -m "$message"
    description: "Commit staged changes"
    working_dir: "${PWD}"
    keep_open: true
    
  git_push:
    type: script
    command: "git push origin $(git branch --show-current)"
    description: "Push to current branch"
    working_dir: "${PWD}"
    
  git_pull:
    type: script
    command: "git pull origin $(git branch --show-current)"
    description: "Pull from current branch"
    working_dir: "${PWD}"
    
  git_diff:
    type: script
    command: "git diff"
    description: "Show unstaged changes"
    working_dir: "${PWD}"
    keep_open: true
    
  git_branch:
    type: script
    command: "git branch -a"
    description: "List all branches"
    working_dir: "${PWD}"
    show_output: true
    
  git_history:
    type: script
    command: "git log --oneline --graph --decorate -10"
    description: "Show commit history"
    working_dir: "${PWD}"
    keep_open: true
    
  # === Docker Commands ===
  docker_up:
    type: script
    command: "docker-compose up -d"
    description: "Start Docker containers"
    working_dir: "${PWD}"
    
  docker_down:
    type: script
    command: "docker-compose down"
    description: "Stop Docker containers"
    working_dir: "${PWD}"
    
  docker_logs:
    type: script
    command: "docker-compose logs -f --tail=100"
    description: "Show Docker logs"
    working_dir: "${PWD}"
    keep_open: true
    
  docker_ps:
    type: script
    command: "docker ps -a"
    description: "List Docker containers"
    show_output: true
    
  docker_images:
    type: script
    command: "docker images"
    description: "List Docker images"
    show_output: true
    
  docker_clean:
    type: script
    command: |
      docker system prune -f
      docker volume prune -f
      echo "Docker cleanup completed!"
    description: "Clean unused Docker resources"
    show_output: true
    
  # === Project Commands ===
  project_build:
    type: script
    command: "${BUILD_CMD:-make build}"
    description: "Build the project"
    working_dir: "${PROJECT_ROOT:-${PWD}}"
    keep_open: true
    
  project_test:
    type: script
    command: "${TEST_CMD:-make test}"
    description: "Run tests"
    working_dir: "${PROJECT_ROOT:-${PWD}}"
    keep_open: true
    
  project_run:
    type: script
    command: "${RUN_CMD:-make run}"
    description: "Run the project"
    working_dir: "${PROJECT_ROOT:-${PWD}}"
    keep_open: true
    
  project_clean:
    type: script
    command: "${CLEAN_CMD:-make clean}"
    description: "Clean build artifacts"
    working_dir: "${PROJECT_ROOT:-${PWD}}"
    
  # === System Commands ===
  system_lock:
    type: script
    command: "rundll32.exe user32.dll,LockWorkStation"  # Windows, override for other platforms
    description: "Lock the screen"
    
  system_sleep:
    type: script
    command: "rundll32.exe powrprof.dll,SetSuspendState 0,1,0"  # Windows, override for other platforms
    description: "Put system to sleep"
    
  # === Quick Tools ===
  quick_note:
    type: app
    command: "notepad"  # Override with your favorite note app
    description: "Quick note taking"
    
  quick_calc:
    type: app
    command: "calc"  # Windows calculator, override for other platforms
    description: "Calculator"
    
  quick_search:
    type: url
    command: "https://www.google.com"
    description: "Web search"

# Environment variables available to all scripts
env:
  # Custom environment variables for your scripts
  EDITOR: "code"
  BROWSER: "chrome"
  TERMINAL: "wt"
```
:::

## Windows Configuration

Platform-specific overrides for Windows. Save this as `spellbook.windows.yml` in your configuration directory.

::: details spellbook.windows.yml - Click to expand
```yaml
# Windows-specific overrides for SilentCast
# This file only contains settings that differ from the base spellbook.yml

# Windows-specific grimoire overrides
grimoire:
  # === Windows Applications ===
  terminal:
    type: app
    command: "wt"  # Windows Terminal
    description: "Windows Terminal"
    args: ["-w", "_quake"]  # Quake mode
    
  file_manager:
    type: app
    command: "explorer"
    description: "Windows Explorer"
    
  settings:
    type: app
    command: "ms-settings:"
    description: "Windows Settings"
    
  # === Windows-specific tools ===
  powershell:
    type: app
    command: "pwsh"  # PowerShell Core
    description: "PowerShell Core"
    admin: false
    
  powershell_admin:
    type: app
    command: "pwsh"
    description: "PowerShell (Admin)"
    admin: true
    
  cmd:
    type: app
    command: "cmd"
    description: "Command Prompt"
    
  cmd_admin:
    type: app
    command: "cmd"
    description: "Command Prompt (Admin)"
    admin: true
    
  wsl:
    type: app
    command: "wsl"
    description: "Windows Subsystem for Linux"
    terminal: true
    
  # === System Tools ===
  task_manager:
    type: app
    command: "taskmgr"
    description: "Task Manager"
    
  device_manager:
    type: app
    command: "devmgmt.msc"
    description: "Device Manager"
    
  services:
    type: app
    command: "services.msc"
    description: "Services"
    
  registry:
    type: app
    command: "regedit"
    description: "Registry Editor"
    admin: true
    
  event_viewer:
    type: app
    command: "eventvwr"
    description: "Event Viewer"
    
  # === Development Tools ===
  vs2022:
    type: app
    command: "C:\\Program Files\\Microsoft Visual Studio\\2022\\Community\\Common7\\IDE\\devenv.exe"
    description: "Visual Studio 2022"
    
  vs2019:
    type: app
    command: "C:\\Program Files (x86)\\Microsoft Visual Studio\\2019\\Community\\Common7\\IDE\\devenv.exe"
    description: "Visual Studio 2019"
    
  # === Windows Terminal profiles ===
  terminal_split:
    type: app
    command: "wt"
    args: ["split-pane", "-H", "pwsh", ";", "split-pane", "-V", "cmd"]
    description: "Terminal with split panes"
    
  terminal_tabs:
    type: app
    command: "wt"
    args: ["new-tab", "-p", "PowerShell", ";", "new-tab", "-p", "Command Prompt"]
    description: "Terminal with multiple tabs"
    
  # === System Commands Override ===
  system_lock:
    type: script
    command: "rundll32.exe user32.dll,LockWorkStation"
    description: "Lock Windows"
    
  system_sleep:
    type: script
    command: "rundll32.exe powrprof.dll,SetSuspendState 0,1,0"
    description: "Sleep mode"
    
  system_hibernate:
    type: script
    command: "shutdown /h"
    description: "Hibernate"
    admin: true
    
  system_restart:
    type: script
    command: "shutdown /r /t 0"
    description: "Restart Windows"
    admin: true
    
  system_shutdown:
    type: script
    command: "shutdown /s /t 0"
    description: "Shutdown Windows"
    admin: true
    
  # === Quick Access ===
  quick_note:
    type: app
    command: "notepad"
    description: "Notepad"
    
  quick_calc:
    type: app
    command: "calc"
    description: "Calculator"
    
  snipping_tool:
    type: app
    command: "SnippingTool"
    description: "Snipping Tool"
    
  # === PowerShell Scripts ===
  ps_profile:
    type: script
    shell: "pwsh"
    command: "code $PROFILE"
    description: "Edit PowerShell profile"
    
  ps_update_help:
    type: script
    shell: "pwsh"
    command: "Update-Help -Force"
    description: "Update PowerShell help"
    admin: true
    
  # === System Info ===
  system_info:
    type: script
    shell: "pwsh"
    command: |
      Get-ComputerInfo | Select-Object CsName, WindowsVersion, OsArchitecture, CsProcessors, CsTotalPhysicalMemory
    description: "System information"
    show_output: true
    
  # === Network Tools ===
  flush_dns:
    type: script
    command: "ipconfig /flushdns"
    description: "Flush DNS cache"
    admin: true
    
  network_reset:
    type: script
    command: |
      netsh winsock reset
      netsh int ip reset
      ipconfig /release
      ipconfig /renew
      ipconfig /flushdns
    description: "Reset network stack"
    admin: true
    keep_open: true
    
  # === Windows Store Apps ===
  store:
    type: app
    command: "ms-windows-store:"
    description: "Microsoft Store"
    
  xbox:
    type: app
    command: "ms-xbox:"
    description: "Xbox App"
    
  photos:
    type: app
    command: "ms-photos:"
    description: "Photos App"

# Windows-specific spells
spells:
  # === Windows Key Shortcuts ===
  "win+e": "file_manager"
  "win+r": "run_dialog"
  "win+x": "power_menu"
  "win+t": "task_manager"
  
  # === Admin Tools ===
  "a,p": "powershell_admin"     # Alt+Space, A, P -> Admin PowerShell
  "a,c": "cmd_admin"            # Alt+Space, A, C -> Admin CMD
  "a,r": "registry"             # Alt+Space, A, R -> Registry Editor
  "a,s": "services"             # Alt+Space, A, S -> Services
  
  # === Windows Specific ===
  "w,s": "wsl"                  # Alt+Space, W, S -> WSL
  "w,t": "terminal_split"       # Alt+Space, W, T -> Terminal split
  "w,n": "snipping_tool"        # Alt+Space, W, N -> Screenshot

# Windows environment variables
env:
  # Common Windows paths
  USERPROFILE: "${USERPROFILE}"
  LOCALAPPDATA: "${LOCALAPPDATA}"
  APPDATA: "${APPDATA}"
  PROGRAMFILES: "${PROGRAMFILES}"
  PROGRAMFILES_X86: "${PROGRAMFILES(X86)}"
  
  # Tools
  TERMINAL: "wt"
  SHELL: "pwsh"
```
:::

## macOS Configuration

Platform-specific overrides for macOS. Save this as `spellbook.mac.yml` in your configuration directory.

::: details spellbook.mac.yml - Click to expand
```yaml
# macOS-specific overrides for SilentCast
# This file only contains settings that differ from the base spellbook.yml

# macOS-specific grimoire overrides
grimoire:
  # === macOS Applications ===
  terminal:
    type: app
    command: "Terminal"
    description: "Terminal.app"
    
  iterm:
    type: app
    command: "iTerm"
    description: "iTerm2"
    
  file_manager:
    type: app
    command: "Finder"
    description: "Finder"
    
  settings:
    type: app
    command: "System Preferences"
    description: "System Preferences"
    
  # === macOS-specific tools ===
  safari:
    type: app
    command: "Safari"
    description: "Safari browser"
    
  preview:
    type: app
    command: "Preview"
    description: "Preview app"
    
  xcode:
    type: app
    command: "Xcode"
    description: "Xcode IDE"
    
  simulator:
    type: app
    command: "Simulator"
    description: "iOS Simulator"
    
  # === System Tools ===
  activity_monitor:
    type: app
    command: "Activity Monitor"
    description: "Activity Monitor"
    
  disk_utility:
    type: app
    command: "Disk Utility"
    description: "Disk Utility"
    
  keychain:
    type: app
    command: "Keychain Access"
    description: "Keychain Access"
    
  console:
    type: app
    command: "Console"
    description: "Console logs"
    
  # === Creative Apps ===
  music:
    type: app
    command: "Music"
    description: "Apple Music"
    
  photos:
    type: app
    command: "Photos"
    description: "Photos app"
    
  # === Development Tools ===
  homebrew_update:
    type: script
    command: "brew update && brew upgrade"
    description: "Update Homebrew"
    keep_open: true
    
  homebrew_cleanup:
    type: script
    command: "brew cleanup -s && brew doctor"
    description: "Cleanup Homebrew"
    show_output: true
    
  # === macOS Commands ===
  spotlight:
    type: script
    command: "open -a 'Spotlight'"
    description: "Open Spotlight"
    
  mission_control:
    type: script
    command: |
      osascript -e 'tell application "Mission Control" to launch'
    description: "Mission Control"
    
  show_desktop:
    type: script
    command: |
      osascript -e 'tell application "System Events" to key code 103 using {command down, option down}'
    description: "Show Desktop"
    
  # === System Commands Override ===
  system_lock:
    type: script
    command: |
      osascript -e 'tell application "System Events" to keystroke "q" using {command down, control down}'
    description: "Lock screen"
    
  system_sleep:
    type: script
    command: "pmset sleepnow"
    description: "Sleep mode"
    admin: true
    
  system_restart:
    type: script
    command: |
      osascript -e 'tell app "System Events" to restart'
    description: "Restart macOS"
    
  system_shutdown:
    type: script
    command: |
      osascript -e 'tell app "System Events" to shut down'
    description: "Shutdown macOS"
    
  # === Quick Access ===
  quick_note:
    type: app
    command: "Notes"
    description: "Apple Notes"
    
  quick_calc:
    type: app
    command: "Calculator"
    description: "Calculator"
    
  # === Screenshots ===
  screenshot_full:
    type: script
    command: "screencapture ~/Desktop/screenshot.png"
    description: "Full screenshot"
    
  screenshot_area:
    type: script
    command: "screencapture -i ~/Desktop/screenshot.png"
    description: "Area screenshot"
    
  screenshot_window:
    type: script
    command: "screencapture -w ~/Desktop/screenshot.png"
    description: "Window screenshot"
    
  # === AppleScript Actions ===
  toggle_dark_mode:
    type: script
    command: |
      osascript -e 'tell app "System Events" to tell appearance preferences to set dark mode to not dark mode'
    description: "Toggle dark mode"
    
  notification_test:
    type: script
    command: |
      osascript -e 'display notification "SilentCast is working!" with title "Test Notification" sound name "Glass"'
    description: "Test notification"
    
  # === File Operations ===
  reveal_in_finder:
    type: script
    command: "open -R '${FILE}'"
    description: "Reveal in Finder"
    
  quick_look:
    type: script
    command: "qlmanage -p '${FILE}' &>/dev/null &"
    description: "Quick Look"
    
  # === Network Tools ===
  wifi_toggle:
    type: script
    command: |
      networksetup -getairportpower en0 | grep -q "On" && \
      networksetup -setairportpower en0 off || \
      networksetup -setairportpower en0 on
    description: "Toggle WiFi"
    admin: true
    
  flush_dns:
    type: script
    command: "sudo dscacheutil -flushcache && sudo killall -HUP mDNSResponder"
    description: "Flush DNS cache"
    admin: true
    
  # === Development ===
  xcode_clean:
    type: script
    command: "rm -rf ~/Library/Developer/Xcode/DerivedData/*"
    description: "Clean Xcode derived data"
    
  simulator_reset:
    type: script
    command: "xcrun simctl shutdown all && xcrun simctl erase all"
    description: "Reset all simulators"

# macOS-specific spells
spells:
  # === macOS Key Shortcuts ===
  "cmd+space": "spotlight"
  "cmd+tab": "app_switcher"
  "cmd+shift+3": "screenshot_full"
  "cmd+shift+4": "screenshot_area"
  "cmd+shift+5": "screenshot_window"
  
  # === macOS Specific ===
  "m,s": "mission_control"      # Alt+Space, M, S -> Mission Control
  "m,d": "show_desktop"         # Alt+Space, M, D -> Show Desktop
  "m,n": "notification_test"    # Alt+Space, M, N -> Test notification
  "m,t": "toggle_dark_mode"     # Alt+Space, M, T -> Toggle dark mode
  
  # === Development ===
  "x,c": "xcode"                # Alt+Space, X, C -> Xcode
  "x,s": "simulator"            # Alt+Space, X, S -> iOS Simulator
  "x,r": "xcode_clean"          # Alt+Space, X, R -> Clean Xcode
  
  # === Homebrew ===
  "h,u": "homebrew_update"      # Alt+Space, H, U -> Update Homebrew
  "h,c": "homebrew_cleanup"     # Alt+Space, H, C -> Cleanup Homebrew

# macOS environment variables
env:
  # Common macOS paths
  HOME: "${HOME}"
  USER: "${USER}"
  
  # Tools
  TERMINAL: "Terminal"
  SHELL: "/bin/zsh"
  EDITOR: "code"
```
:::

## Test Configuration

A minimal test configuration for testing and debugging. Save this as `spellbook.test.yml`.

::: details spellbook.test.yml - Click to expand
```yaml
# Minimal test configuration
daemon:
  log_level: debug
  config_watch: false

logger:
  level: debug
  file: "./test.log"
  max_size: 1
  max_backups: 1

hotkeys:
  prefix: "alt+shift+t"
  timeout: 500

spells:
  t: "test"
  e: "echo"

grimoire:
  test:
    type: script
    command: "echo 'Test successful!'"
    show_output: true
    
  echo:
    type: script
    command: "echo '${MESSAGE:-Hello from SilentCast}'"
    show_output: true
```
:::

## Custom Key Names Example

Example showing how to use custom YAML key names for better organization.

::: details example_custom_keys.yml - Click to expand
```yaml
# Example with custom YAML keys
# Shows how underscores, hyphens, and other naming conventions work

# Custom daemon settings with different naming styles
daemon-settings:
  auto_start: true
  log-level: info
  config_watch: true

# Alternative hotkey naming
hotkey_config:
  prefix_key: "alt+space"
  timeout_ms: 1000
  sequence_timeout_ms: 2000

# Custom spell naming
my_shortcuts:
  e: "my_editor"
  t: "my_terminal"
  "g,s": "git_status_custom"

# Custom grimoire naming
my_actions:
  my_editor:
    type: app
    command: "code"
    
  my_terminal:
    type: app
    command: "wt"
    
  git_status_custom:
    type: script
    command: "git status --short"
    show_output: true
```
:::

## Usage Instructions

### How to Use These Samples

1. **Choose your base configuration**
   - Copy the base `spellbook.yml` to your configuration directory
   - This provides a comprehensive starting point

2. **Add platform-specific overrides**
   - If on Windows, also copy `spellbook.windows.yml`
   - If on macOS, also copy `spellbook.mac.yml`
   - These files only need to contain settings that differ from the base

3. **Customize for your needs**
   - Edit the file paths to match your installed applications
   - Adjust keyboard shortcuts to your preference
   - Add or remove actions based on your workflow

### Configuration Locations

::: code-group
```powershell [Windows]
# Configuration directory
%APPDATA%\SilentCast\

# Copy files
copy spellbook.yml %APPDATA%\SilentCast\
copy spellbook.windows.yml %APPDATA%\SilentCast\
```

```bash [macOS]
# Configuration directory
~/.config/silentcast/

# Copy files
cp spellbook.yml ~/.config/silentcast/
cp spellbook.mac.yml ~/.config/silentcast/
```
:::

### Testing Your Configuration

After setting up your configuration:

1. **Validate the configuration**
   ```bash
   silentcast --validate-config
   ```

2. **Run in debug mode**
   ```bash
   silentcast --log-level debug
   ```

3. **Test a simple spell**
   - Press your prefix key (e.g., `Alt+Space`)
   - Press a spell key (e.g., `e` for editor)
   - Verify the action executes correctly

## Tips and Best Practices

### Organization Tips

1. **Group related shortcuts**
   - Use prefixes for categories (g for git, d for docker, etc.)
   - Keep similar actions together in the configuration

2. **Use descriptive names**
   - Action names should clearly indicate what they do
   - Add descriptions to all grimoire entries

3. **Start simple**
   - Begin with a few essential shortcuts
   - Add more as you get comfortable

### Performance Tips

1. **Keep commands fast**
   - Avoid long-running scripts in shortcuts
   - Use `keep_open: true` for commands that need time

2. **Use environment variables**
   - Define common paths in the env section
   - Reference them with `${VARIABLE_NAME}`

3. **Platform-specific optimization**
   - Use native commands for each platform
   - Avoid unnecessary shell wrapping

### Security Considerations

1. **Never store secrets**
   - Don't put passwords or API keys in configuration
   - Use environment variables from your system

2. **Be careful with admin/sudo**
   - Only use elevated privileges when necessary
   - Understand what commands do before adding them

3. **Validate scripts**
   - Test scripts manually first
   - Be cautious with destructive commands

## Next Steps

- Review the [Configuration Guide](/guide/configuration) for detailed explanations
- Learn about [Advanced Shortcuts](/guide/shortcuts)
- Explore [Script Execution](/guide/scripts) features
- Check [Platform-Specific Features](/guide/platforms)