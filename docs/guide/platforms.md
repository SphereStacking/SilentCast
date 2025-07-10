# Platform Support Guide

SilentCast is designed to work seamlessly across Windows and macOS. This guide covers platform-specific features, limitations, and optimizations.

## Platform Overview

| Feature | Windows | macOS |
|---------|---------|--------|
| Global Hotkeys | ✅ Full | ✅ Full |
| System Tray | ✅ Native | ✅ Native |
| Auto-start | ✅ Task Scheduler | ✅ launchd |
| Permissions | ✅ None | ⚠️ Accessibility |
| Shell | PowerShell/CMD | zsh/bash |

## Windows

### System Requirements

- **Version**: Windows 10 version 1903+ or Windows 11
- **Architecture**: x64 or ARM64
- **Dependencies**: None (self-contained)

### Installation Methods

```powershell
# Recommended: Scoop
scoop bucket add silentcast https://github.com/SphereStacking/scoop-silentcast
scoop install silentcast

# Alternative: Chocolatey
choco install silentcast

# Alternative: WinGet
winget install SphereStacking.SilentCast
```

### Windows-Specific Features

#### 1. PowerShell Integration

```yaml
# spellbook.windows.yml
grimoire:
  system_info:
    type: script
    shell: "pwsh"  # PowerShell Core
    command: |
      Get-ComputerInfo | Select-Object CsName, WindowsVersion, OsArchitecture
    show_output: true
    
  classic_powershell:
    type: script
    shell: "powershell"  # Windows PowerShell
    command: "$PSVersionTable"
```

#### 2. Windows Terminal Support

```yaml
grimoire:
  terminal:
    type: app
    command: "wt"
    args: ["-w", "_quake"]  # Quake mode
    
  terminal_split:
    type: app
    command: "wt"
    args: ["split-pane", "-H", "pwsh", ";", "split-pane", "-V", "cmd"]
```

#### 3. Admin Elevation

```yaml
grimoire:
  edit_hosts:
    type: app
    command: "notepad"
    args: ["C:\\Windows\\System32\\drivers\\etc\\hosts"]
    admin: true  # Request UAC elevation
    
  restart_service:
    type: script
    shell: "pwsh"
    command: "Restart-Service -Name 'MySQL' -Force"
    admin: true
```

#### 4. Windows-Specific Shortcuts

```yaml
spells:
  # Windows key combinations
  "win+e": "explorer"
  "win+r": "run_dialog"
  "win+x": "power_menu"
  
grimoire:
  # Quick access to Windows tools
  device_manager:
    type: app
    command: "devmgmt.msc"
    
  services:
    type: app
    command: "services.msc"
    
  registry:
    type: app
    command: "regedit"
    admin: true
```

### Windows Tips

#### Auto-start Setup

```powershell
# Create scheduled task for auto-start
$Action = New-ScheduledTaskAction -Execute "$env:LOCALAPPDATA\SilentCast\silentcast.exe"
$Trigger = New-ScheduledTaskTrigger -AtLogon -User $env:USERNAME
$Settings = New-ScheduledTaskSettingsSet -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries -StartWhenAvailable

Register-ScheduledTask -TaskName "SilentCast" -Action $Action -Trigger $Trigger -Settings $Settings
```

#### Defender Exclusion

```powershell
# Add to Windows Defender exclusions (Admin required)
Add-MpPreference -ExclusionPath "$env:LOCALAPPDATA\SilentCast"
```

### Windows Limitations

1. **UAC Prompts**: Cannot bypass UAC for admin operations
2. **Secure Desktop**: Hotkeys don't work on UAC/login screens
3. **Game Mode**: Some games may block global hotkeys

## macOS

### System Requirements

- **Version**: macOS 10.15 Catalina or later
- **Architecture**: Intel (x64) or Apple Silicon (arm64)
- **Permissions**: Accessibility access required

### Installation Methods

```bash
# Recommended: Homebrew
brew tap spherestacking/silentcast
brew install silentcast

# Alternative: Direct download
curl -L https://github.com/SphereStacking/silentcast/releases/latest/download/silentcast-darwin-$(uname -m).tar.gz | tar xz
sudo mv silentcast /usr/local/bin/
```

### macOS-Specific Features

#### 1. Accessibility Permissions

First run requires granting accessibility permissions:

1. Run SilentCast: `silentcast`
2. When prompted, open System Preferences
3. Go to Security & Privacy → Privacy → Accessibility
4. Click the lock and authenticate
5. Add SilentCast to the list
6. Restart SilentCast

#### 2. AppleScript Integration

```yaml
# spellbook.mac.yml
grimoire:
  toggle_dark_mode:
    type: script
    command: |
      osascript -e 'tell app "System Events" to tell appearance preferences to set dark mode to not dark mode'
      
  notification:
    type: script
    command: |
      osascript -e 'display notification "Task completed!" with title "SilentCast" sound name "Glass"'
      
  finder_selection:
    type: script
    command: |
      osascript -e 'tell application "Finder" to get selection as alias list'
    show_output: true
```

#### 3. macOS App Shortcuts

```yaml
grimoire:
  # Open apps with specific documents
  preview_pdf:
    type: app
    command: "open"
    args: ["-a", "Preview", "document.pdf"]
    
  # Reveal in Finder
  show_in_finder:
    type: app
    command: "open"
    args: ["-R", "${FILE}"]
    
  # Quick Look
  quick_look:
    type: script
    command: "qlmanage -p '${FILE}' &>/dev/null &"
```

#### 4. System Integration

```yaml
grimoire:
  # Spotlight search
  spotlight:
    type: script
    command: "open -a Spotlight"
    
  # Mission Control
  mission_control:
    type: script
    command: |
      osascript -e 'tell application "Mission Control" to launch'
      
  # Screenshot with options
  screenshot:
    type: script
    command: "screencapture -i -c"  # Interactive, to clipboard
```

### macOS Tips

#### Auto-start with launchd

Create `~/Library/LaunchAgents/com.silentcast.plist`:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.silentcast</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/silentcast</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardErrorPath</key>
    <string>/tmp/silentcast.err</string>
    <key>StandardOutPath</key>
    <string>/tmp/silentcast.out</string>
</dict>
</plist>
```

Load it:
```bash
launchctl load ~/Library/LaunchAgents/com.silentcast.plist
```

#### Keyboard Shortcuts

```yaml
spells:
  # macOS standard shortcuts
  "cmd+space": "spotlight_enhanced"
  "cmd+tab": "app_switcher_enhanced"
  "cmd+shift+3": "screenshot_full"
  "cmd+shift+4": "screenshot_area"
```

### macOS Limitations

1. **Accessibility Required**: Must grant permissions for global hotkeys
2. **Secure Input**: Some apps (password managers) block hotkeys
3. **Full Disk Access**: May need for certain file operations

## Cross-Platform Best Practices

### 1. Use Platform Overrides

```yaml
# spellbook.yml (base config)
grimoire:
  terminal:
    type: app
    command: "terminal"  # Generic placeholder
    
# spellbook.windows.yml
grimoire:
  terminal:
    command: "wt"  # Windows Terminal
    
# spellbook.mac.yml
grimoire:
  terminal:
    command: "Terminal"  # macOS Terminal
```

### 2. Environment Detection

```yaml
grimoire:
  smart_browser:
    type: script
    command: |
      # Detect and use default browser
      if command -v open &> /dev/null; then
        open "https://example.com"       # macOS
      elif command -v start &> /dev/null; then
        start "https://example.com"      # Windows
      fi
```

### 3. Path Handling

```yaml
grimoire:
  open_config:
    type: app
    command: "${EDITOR:-code}"
    args: ["${SILENTCAST_CONFIG}/spellbook.yml"]
    # SILENTCAST_CONFIG is automatically set to the right path
```

### 4. Shell Compatibility

```yaml
grimoire:
  cross_platform_script:
    type: script
    # Use sh for maximum compatibility
    shell: "sh"
    command: |
      echo "This works everywhere"
      echo "Home: $HOME"
      echo "User: $USER"
```

## Platform-Specific Troubleshooting

### Windows Issues

<details>
<summary>Hotkeys not working in some applications</summary>

Some applications (especially games) may capture all keyboard input. Try:
- Running SilentCast as administrator
- Using windowed mode in games
- Configuring game-specific exceptions

</details>

<details>
<summary>PowerShell scripts blocked</summary>

If PowerShell scripts are blocked:
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

</details>

### macOS Issues

<details>
<summary>Accessibility permission keeps resetting</summary>

1. Fully quit SilentCast
2. Remove from accessibility list
3. Delete the app's preferences: `defaults delete com.silentcast`
4. Re-add to accessibility
5. Restart SilentCast

</details>

<details>
<summary>Hotkeys not working after macOS update</summary>

macOS updates often reset permissions:
1. Re-grant accessibility permissions
2. Check for SilentCast updates
3. Restart your Mac

</details>

## Next Steps

- Explore [Script Execution](/guide/scripts) for platform-specific automation
- Configure [Environment Variables](/guide/env-vars)
- Check [CLI Reference](/api/) for platform-specific options
- Join platform-specific discussions on [Discord](https://discord.gg/silentcast)