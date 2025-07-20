# Installation

This guide covers all installation methods for SilentCast across different platforms. Choose the method that best suits your needs.

## 📋 System Requirements

### Minimum Requirements
- **OS**: Windows 10/11, macOS 10.15+, or Linux (Ubuntu 20.04+)
- **Memory**: 64MB RAM
- **Storage**: 20MB available space
- **Permissions**: Admin/sudo access for global hotkeys

### Supported Architectures
- x86_64 (AMD64)
- ARM64 (Apple Silicon, ARM Linux)

## 🚀 Quick Install

The fastest way to install SilentCast:

### Using Package Managers (Recommended)

::: code-group

```bash [macOS]
# Install with Homebrew
brew install spherestacking/tap/silentcast

# Start SilentCast
silentcast
```

```powershell [Windows]
# Install with Scoop
scoop bucket add spherestacking https://github.com/spherestacking/scoop-bucket
scoop install silentcast

# Start SilentCast
silentcast
```

```bash [Linux]
# Install with Snap
sudo snap install silentcast

# Start SilentCast
silentcast
```

:::

### Direct Download

Download the latest binary for your platform:

::: code-group

```bash [macOS/Linux]
# Download and install (replace VERSION and PLATFORM)
curl -L https://github.com/SphereStacking/silentcast/releases/latest/download/silentcast-{PLATFORM}-{ARCH}.tar.gz | tar xz
sudo mv silentcast /usr/local/bin/
silentcast --version
```

```powershell [Windows]
# Download latest release
$url = "https://github.com/SphereStacking/silentcast/releases/latest/download/silentcast-windows-amd64.zip"
Invoke-WebRequest -Uri $url -OutFile silentcast.zip

# Extract and add to PATH
Expand-Archive silentcast.zip -DestinationPath C:\Tools\SilentCast
$env:Path += ";C:\Tools\SilentCast"
silentcast --version
```

:::

## 📦 Package Managers

### macOS - Homebrew

```bash
# Add tap and install
brew install spherestacking/tap/silentcast

# Or tap first, then install
brew tap spherestacking/tap
brew install silentcast

# Upgrade to latest version
brew upgrade silentcast

# Start as background service
brew services start silentcast

# Stop service
brew services stop silentcast

# Restart service
brew services restart silentcast
```

### Windows - Multiple Options

#### Scoop (Recommended)

```powershell
# Add the SphereStacking bucket
scoop bucket add spherestacking https://github.com/spherestacking/scoop-bucket

# Install SilentCast
scoop install silentcast

# Update to latest version
scoop update silentcast

# Show installed version
scoop info silentcast
```

#### Chocolatey

```powershell
# Run as Administrator
choco install silentcast -y

# Upgrade to latest
choco upgrade silentcast -y

# List version
choco list silentcast
```

#### WinGet

```powershell
# Search for SilentCast
winget search silentcast

# Install
winget install SphereStacking.SilentCast

# Upgrade
winget upgrade SphereStacking.SilentCast

# Show info
winget show SphereStacking.SilentCast
```

### Linux - Package Managers

#### Snap

```bash
# Install from Snap Store
sudo snap install silentcast

# Install specific channel
sudo snap install silentcast --channel=edge

# Update
sudo snap refresh silentcast

# List installed snaps
snap list silentcast
```

#### APT (Debian/Ubuntu)

```bash
# Add repository
curl -fsSL https://pkg.silentcast.dev/gpg | sudo apt-key add -
echo "deb https://pkg.silentcast.dev/apt stable main" | sudo tee /etc/apt/sources.list.d/silentcast.list

# Install
sudo apt update
sudo apt install silentcast

# Upgrade
sudo apt upgrade silentcast
```

#### YUM/DNF (RHEL/Fedora)

```bash
# Add repository
sudo dnf config-manager --add-repo https://pkg.silentcast.dev/rpm/silentcast.repo

# Install
sudo dnf install silentcast

# Upgrade
sudo dnf upgrade silentcast
```

## 🔧 Manual Installation

### Download Binaries

1. Visit the [releases page](https://github.com/SphereStacking/silentcast/releases)
2. Download the appropriate archive for your platform:
   - **Windows**: `silentcast-windows-amd64.zip`
   - **macOS Intel**: `silentcast-darwin-amd64.tar.gz`
   - **macOS Apple Silicon**: `silentcast-darwin-arm64.tar.gz`
   - **Linux x64**: `silentcast-linux-amd64.tar.gz`
   - **Linux ARM64**: `silentcast-linux-arm64.tar.gz`

### Installation Steps

::: code-group

```bash [macOS/Linux]
# Extract archive
tar -xzf silentcast-*.tar.gz

# Move to PATH (choose one)
sudo mv silentcast /usr/local/bin/      # System-wide
mv silentcast ~/.local/bin/             # User only

# Make executable
chmod +x /usr/local/bin/silentcast      # or ~/.local/bin/silentcast

# Verify installation
silentcast --version

# Test basic functionality
silentcast --validate-config
```

```powershell [Windows]
# Create installation directory
New-Item -ItemType Directory -Force -Path "C:\Program Files\SilentCast"

# Extract ZIP
Expand-Archive silentcast-windows-amd64.zip -DestinationPath "C:\Program Files\SilentCast"

# Add to system PATH (requires admin)
[Environment]::SetEnvironmentVariable(
    "Path",
    $env:Path + ";C:\Program Files\SilentCast",
    [EnvironmentVariableTarget]::Machine
)

# Or add to user PATH (no admin required)
[Environment]::SetEnvironmentVariable(
    "Path",
    $env:Path + ";C:\Program Files\SilentCast",
    [EnvironmentVariableTarget]::User
)

# Refresh PATH in current session
$env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")

# Verify installation
silentcast --version
```

:::

## 🛠️ Building from Source

### Prerequisites

- Go 1.21 or later
- Git
- Make (optional but recommended)
- C compiler (for full features, optional)

### Quick Build

```bash
# Clone repository
git clone https://github.com/SphereStacking/silentcast.git
cd silentcast

# Setup development environment
make setup

# Build for current platform
make build

# Install to system
sudo make install
```

### Build Options

#### Development Build (No C Dependencies)

```bash
# Fast build without hotkey support (perfect for development)
make build-stub

# Run directly
./app/build/silentcast --no-tray
```

#### Production Build

```bash
# Full build with all features
make build

# Build for all platforms
make build-all

# Build with version info
make build VERSION=1.2.3
```

#### Cross-Platform Build

```bash
# macOS (Intel)
GOOS=darwin GOARCH=amd64 make build

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 make build

# Windows
GOOS=windows GOARCH=amd64 make build

# Linux
GOOS=linux GOARCH=amd64 make build

# Linux (ARM64)
GOOS=linux GOARCH=arm64 make build
```

#### Static Build (No External Dependencies)

```bash
# Build with static linking
CGO_ENABLED=0 make build

# Build minimal binary
make build-minimal
```

## ⚙️ Post-Installation Setup

### 1. Create Your Spellbook

SilentCast needs a configuration file to define your spells. Configuration locations:

::: code-group

```bash [macOS/Linux]
# Create config directory
mkdir -p ~/.config/silentcast

# Create a basic spellbook
cat > ~/.config/silentcast/spellbook.yml << 'EOF'
# SilentCast Spellbook
hotkeys:
  prefix: "alt+space"      # Your activation key

spells:
  e: "editor"              # Alt+Space, then E
  t: "terminal"            # Alt+Space, then T
  "g,s": "git_status"      # Alt+Space, then G, then S

grimoire:
  editor:
    type: app
    command: "code"        # VS Code
    description: "Open VS Code"
    
  terminal:
    type: app
    command: "open -a Terminal"  # macOS Terminal
    description: "Open Terminal"
    
  git_status:
    type: script
    command: "git status"
    show_output: true
    description: "Show git status"
EOF

# Or download example configuration
curl -o ~/.config/silentcast/spellbook.yml \
  https://raw.githubusercontent.com/SphereStacking/silentcast/main/examples/config/basic_spellbook.yml
```

```powershell [Windows]
# Create config directory
New-Item -ItemType Directory -Force -Path "$env:APPDATA\SilentCast"

# Create a basic spellbook
@'
# SilentCast Spellbook
hotkeys:
  prefix: "alt+space"      # Your activation key

spells:
  e: "editor"              # Alt+Space, then E
  t: "terminal"            # Alt+Space, then T
  "g,s": "git_status"      # Alt+Space, then G, then S

grimoire:
  editor:
    type: app
    command: "code"        # VS Code
    description: "Open VS Code"
    
  terminal:
    type: app
    command: "wt"          # Windows Terminal
    description: "Open Windows Terminal"
    
  git_status:
    type: script
    command: "git status"
    show_output: true
    description: "Show git status"
'@ | Out-File -FilePath "$env:APPDATA\SilentCast\spellbook.yml" -Encoding UTF8

# Or download example configuration
Invoke-WebRequest `
  -Uri "https://raw.githubusercontent.com/SphereStacking/silentcast/main/examples/config/basic_spellbook.yml" `
  -OutFile "$env:APPDATA\SilentCast\spellbook.yml"
```

:::

### 2. Platform-Specific Setup

#### macOS - Accessibility Permissions

SilentCast needs accessibility permissions to capture global hotkeys:

1. **First Run**: Start SilentCast
   ```bash
   silentcast
   ```

2. **Permission Dialog**: macOS will show a permission request
   - Click "Open System Preferences" when prompted

3. **Grant Permission**:
   - Navigate to: **System Preferences → Security & Privacy → Privacy → Accessibility**
   - Click the lock 🔒 and authenticate
   - Check the box next to SilentCast ✅
   - Click the lock again to save

4. **Restart SilentCast**:
   ```bash
   # Stop if running
   pkill silentcast
   
   # Start again
   silentcast
   ```

#### Windows - First Run

No special permissions needed, but you may see:

1. **Windows Defender SmartScreen**: Click "More info" → "Run anyway"
2. **Firewall prompt**: Not needed unless using network features
3. **Antivirus warnings**: Add exception for silentcast.exe

#### Linux - Additional Setup

Depending on your distribution:

```bash
# Ubuntu/Debian - Install dependencies
sudo apt install libappindicator3-1 libgtk-3-0

# Fedora
sudo dnf install libappindicator-gtk3

# Arch
sudo pacman -S libappindicator-gtk3

# For Wayland users - might need X11 compatibility
sudo apt install libxkbcommon-x11-0
```

### 3. Auto-Start Configuration

::: code-group

```bash [macOS - Homebrew]
# If installed via Homebrew
brew services start silentcast

# Check status
brew services list | grep silentcast

# Stop/restart
brew services stop silentcast
brew services restart silentcast
```

```bash [macOS - Manual]
# Create launch agent
cat > ~/Library/LaunchAgents/com.spherestacking.silentcast.plist << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.spherestacking.silentcast</string>
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
EOF

# Load the agent
launchctl load ~/Library/LaunchAgents/com.spherestacking.silentcast.plist

# Verify it's running
launchctl list | grep silentcast
```

```powershell [Windows - Startup]
# Option 1: Add to Startup folder (simplest)
$startupPath = "$env:APPDATA\Microsoft\Windows\Start Menu\Programs\Startup"
$silentcastPath = (Get-Command silentcast).Path
$shortcut = "$startupPath\SilentCast.lnk"

$WshShell = New-Object -comObject WScript.Shell
$Shortcut = $WshShell.CreateShortcut($shortcut)
$Shortcut.TargetPath = $silentcastPath
$Shortcut.Save()

Write-Host "SilentCast added to startup"
```

```powershell [Windows - Task Scheduler]
# Option 2: Create scheduled task (more control)
$taskName = "SilentCast"
$exePath = (Get-Command silentcast).Path

# Create the scheduled task
$action = New-ScheduledTaskAction -Execute $exePath
$trigger = New-ScheduledTaskTrigger -AtLogon -User $env:USERNAME
$settings = New-ScheduledTaskSettingsSet `
    -AllowStartIfOnBatteries `
    -DontStopIfGoingOnBatteries `
    -StartWhenAvailable

Register-ScheduledTask `
    -TaskName $taskName `
    -Action $action `
    -Trigger $trigger `
    -Settings $settings `
    -Description "SilentCast - Hotkey-driven task runner"

# Start the task immediately
Start-ScheduledTask -TaskName $taskName

# Verify it's running
Get-ScheduledTask -TaskName $taskName
```

```bash [Linux - systemd]
# Create user service
mkdir -p ~/.config/systemd/user/

cat > ~/.config/systemd/user/silentcast.service << EOF
[Unit]
Description=SilentCast - Hotkey-driven task runner
After=graphical-session.target

[Service]
Type=simple
ExecStart=/usr/local/bin/silentcast
Restart=always
RestartSec=10

[Install]
WantedBy=default.target
EOF

# Enable and start the service
systemctl --user daemon-reload
systemctl --user enable silentcast
systemctl --user start silentcast

# Check status
systemctl --user status silentcast
```

```bash [Linux - Desktop Entry]
# Create autostart entry (works on most desktop environments)
mkdir -p ~/.config/autostart/

cat > ~/.config/autostart/silentcast.desktop << EOF
[Desktop Entry]
Type=Application
Name=SilentCast
Comment=Hotkey-driven task runner
Exec=/usr/local/bin/silentcast
Hidden=false
NoDisplay=false
X-GNOME-Autostart-enabled=true
EOF

# Make it executable
chmod +x ~/.config/autostart/silentcast.desktop
```

:::

## ✅ Verification

After installation, verify everything is working:

```bash
# Check version
silentcast --version

# Validate configuration
silentcast --validate-config

# List configured spells
silentcast --list-spells

# Test hotkey detection
silentcast --test-hotkey

# Run in foreground with debug logging
silentcast --no-tray --debug
```

### Test Your First Spell

1. Start SilentCast
2. Press your prefix key (default: `Alt+Space`)
3. Press `e` to open your editor
4. Try a sequence: `Alt+Space`, then `g`, then `s` for git status

If spells aren't working:
- Check the system tray icon (it should be visible)
- Review logs: `tail -f ~/.config/silentcast/silentcast.log`
- Ensure you have the required permissions (especially on macOS)

## 🔍 Troubleshooting

### Common Installation Issues

<details>
<summary><strong>Command not found</strong></summary>

The binary is not in your PATH. Solutions:

```bash
# Find where silentcast was installed
which silentcast || find / -name silentcast 2>/dev/null

# Add to PATH (bash/zsh)
echo 'export PATH="$PATH:/path/to/silentcast/directory"' >> ~/.bashrc
source ~/.bashrc

# Or move to standard location
sudo mv /current/path/silentcast /usr/local/bin/
```
</details>

<details>
<summary><strong>Permission denied when running</strong></summary>

```bash
# Make executable
chmod +x /path/to/silentcast

# If installed system-wide
sudo chmod +x /usr/local/bin/silentcast

# Check file permissions
ls -la $(which silentcast)
```
</details>

<details>
<summary><strong>macOS: "Cannot be opened because developer cannot be verified"</strong></summary>

This is macOS Gatekeeper protection. Solutions:

```bash
# Option 1: Remove quarantine attribute
xattr -d com.apple.quarantine /usr/local/bin/silentcast

# Option 2: Allow in System Preferences
# Go to System Preferences → Security & Privacy → General
# Click "Open Anyway" next to the SilentCast message

# Option 3: Temporarily disable Gatekeeper (not recommended)
sudo spctl --master-disable
# Run silentcast
sudo spctl --master-enable
```
</details>

<details>
<summary><strong>Windows: "Windows protected your PC"</strong></summary>

Windows Defender SmartScreen warning:

1. Click **"More info"**
2. Click **"Run anyway"**

To prevent future warnings:
1. Right-click `silentcast.exe`
2. Properties → General
3. Check "Unblock"
4. Apply
</details>

<details>
<summary><strong>Linux: System tray icon missing</strong></summary>

Install required dependencies:

```bash
# Ubuntu/Debian
sudo apt install libappindicator3-1 gir1.2-appindicator3-0.1

# Fedora
sudo dnf install libappindicator-gtk3

# Arch
sudo pacman -S libappindicator-gtk3

# For KDE Plasma
sudo apt install plasma-systray-legacy
```
</details>

<details>
<summary><strong>Hotkeys not working</strong></summary>

1. **Check permissions**:
   - macOS: Accessibility permissions granted?
   - Linux: Running under X11 or Wayland?
   
2. **Test hotkey detection**:
   ```bash
   silentcast --test-hotkey
   ```
   
3. **Check for conflicts**:
   - Other apps using same hotkeys?
   - Try different prefix key
   
4. **Review logs**:
   ```bash
   # macOS/Linux
   tail -f ~/.config/silentcast/silentcast.log
   
   # Windows
   Get-Content "$env:APPDATA\SilentCast\silentcast.log" -Wait
   ```
</details>

## 🎯 Next Steps

Now that SilentCast is installed:

1. **[Configure Your Spellbook](./configuration.md)** - Customize your spells and grimoire
2. **[Learn Spell Patterns](./spells.md)** - Master keyboard combinations
3. **[Explore Scripts](./scripts.md)** - Advanced automation techniques
4. **[Platform Guide](./platforms.md)** - OS-specific features and tips

### Quick Tips

- 🎹 Press `Alt+Space` (or your prefix) + `?` to see available spells
- 📝 Edit your spellbook and changes apply instantly (live reload)
- 🔍 Use `silentcast --dry-run --spell=xxx` to test without executing
- 📊 Run `silentcast --show-config` to see your merged configuration

### Getting Help

- 📖 **Documentation**: [spherestacking.github.io/SilentCast](https://spherestacking.github.io/SilentCast/)
- 🐛 **Issues**: [GitHub Issues](https://github.com/SphereStacking/silentcast/issues)
- 💬 **Community**: [Discord Server](https://discord.gg/silentcast)
- 📧 **Email**: support@spherestacking.com

---

<div align="center">
  <p><strong>Happy spell casting! 🪄</strong></p>
</div>