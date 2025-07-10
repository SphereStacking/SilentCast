# Installation

This guide covers all installation methods for SilentCast across different platforms. Choose the method that best suits your needs.

## System Requirements

### Minimum Requirements
- **OS**: Windows 10+ or macOS 10.15+
- **Memory**: 512MB RAM
- **Storage**: 50MB available space
- **Permissions**: Admin/sudo access for global hotkeys

### Supported Architectures
- x86_64 (AMD64)
- ARM64 (Apple Silicon)

## Quick Install Scripts

The fastest way to get started:

::: code-group

```bash [macOS]
# Universal install script
curl -sSL https://get.silentcast.dev | bash

# With custom installation directory
curl -sSL https://get.silentcast.dev | bash -s -- --prefix=/opt/silentcast

# Install specific version
curl -sSL https://get.silentcast.dev | bash -s -- --version=1.0.0
```

```powershell [Windows]
# Run PowerShell as Administrator
Set-ExecutionPolicy Bypass -Scope Process -Force
iwr -useb https://get.silentcast.dev/install.ps1 | iex

# Install to custom location
iwr -useb https://get.silentcast.dev/install.ps1 | iex -InstallDir "C:\Tools\SilentCast"

# Install specific version
iwr -useb https://get.silentcast.dev/install.ps1 | iex -Version "1.0.0"
```

:::

<div class="warning custom-block">

**Security Note**: Always review install scripts before running them. You can view the scripts at:
- [install.sh](https://get.silentcast.dev/install.sh)
- [install.ps1](https://get.silentcast.dev/install.ps1)

</div>

## Package Managers

### macOS (Homebrew)

```bash
# Install
brew tap spherestacking/silentcast
brew install silentcast

# Upgrade
brew upgrade silentcast

# Start as service
brew services start silentcast
```

### Windows

#### Scoop

```powershell
# Add bucket
scoop bucket add silentcast https://github.com/SphereStacking/scoop-silentcast

# Install
scoop install silentcast

# Update
scoop update silentcast
```

#### Chocolatey

```powershell
# Install (Admin PowerShell)
choco install silentcast

# Upgrade
choco upgrade silentcast
```

#### WinGet

```powershell
# Install
winget install SphereStacking.SilentCast

# Upgrade
winget upgrade SphereStacking.SilentCast
```

## Manual Installation

### Download Binaries

1. Visit the [releases page](https://github.com/SphereStacking/silentcast/releases)
2. Download the appropriate archive for your platform:
   - `silentcast-windows-amd64.zip` - Windows 64-bit
   - `silentcast-darwin-amd64.tar.gz` - macOS Intel
   - `silentcast-darwin-arm64.tar.gz` - macOS Apple Silicon

### Installation Steps

::: code-group

```bash [macOS]
# Extract archive
tar -xzf silentcast-*.tar.gz

# Move to PATH
sudo mv silentcast /usr/local/bin/

# Make executable
sudo chmod +x /usr/local/bin/silentcast

# Verify installation
silentcast --version
```

```powershell [Windows]
# Extract ZIP (PowerShell)
Expand-Archive silentcast-windows-amd64.zip -DestinationPath C:\Tools\SilentCast

# Add to PATH
[Environment]::SetEnvironmentVariable(
    "Path",
    $env:Path + ";C:\Tools\SilentCast",
    [EnvironmentVariableTarget]::User
)

# Refresh PATH
$env:Path = [Environment]::GetEnvironmentVariable("Path", "User")

# Verify installation
silentcast --version
```

:::

## Building from Source

### Prerequisites

- Go 1.21 or later
- Git
- Make (optional but recommended)

### Build Steps

```bash
# Clone repository
git clone https://github.com/SphereStacking/silentcast.git
cd silentcast

# Build for current platform
make build

# Build for all platforms
make build-all

# Install locally
sudo make install

# Run tests
make test
```

### Custom Build Options

```bash
# Build without C dependencies (no system tray)
make build-stub

# Build with specific version
make build VERSION=1.2.3

# Build for specific platform
GOOS=darwin GOARCH=arm64 make build

# Static build (no CGO)
CGO_ENABLED=0 make build
```

## Post-Installation Setup

### 1. Configuration Directory

SilentCast looks for configuration in these locations:

::: code-group

```bash [macOS]
# Create config directory
mkdir -p ~/.config/silentcast

# Copy example configuration
curl -o ~/.config/silentcast/spellbook.yml \
  https://raw.githubusercontent.com/SphereStacking/silentcast/main/examples/spellbook.yml
```

```powershell [Windows]
# Create config directory
New-Item -ItemType Directory -Force -Path "$env:APPDATA\SilentCast"

# Download example configuration
Invoke-WebRequest `
  -Uri "https://raw.githubusercontent.com/SphereStacking/silentcast/main/examples/spellbook.yml" `
  -OutFile "$env:APPDATA\SilentCast\spellbook.yml"
```

:::

### 2. Platform-Specific Setup

#### macOS - Accessibility Permissions

1. Run SilentCast once: `silentcast`
2. When prompted, open System Preferences
3. Go to Security & Privacy → Privacy → Accessibility
4. Click the lock and authenticate
5. Add SilentCast to the list
6. Restart SilentCast

### 3. Auto-Start Configuration

::: code-group

```bash [macOS - launchd]
# Create launch agent
cat > ~/Library/LaunchAgents/dev.silentcast.plist << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>dev.silentcast</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/silentcast</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
</dict>
</plist>
EOF

# Load the agent
launchctl load ~/Library/LaunchAgents/dev.silentcast.plist
```

```powershell [Windows - Task Scheduler]
# Create scheduled task
$Action = New-ScheduledTaskAction -Execute "C:\Tools\SilentCast\silentcast.exe"
$Trigger = New-ScheduledTaskTrigger -AtLogon
$Principal = New-ScheduledTaskPrincipal -UserId $env:USERNAME -RunLevel Highest
$Settings = New-ScheduledTaskSettingsSet -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries

Register-ScheduledTask `
    -TaskName "SilentCast" `
    -Action $Action `
    -Trigger $Trigger `
    -Principal $Principal `
    -Settings $Settings `
    -Description "SilentCast - Silent Hotkey Task Runner"
```

:::

## Verification

After installation, verify everything is working:

```bash
# Check version
silentcast --version

# Validate configuration
silentcast --validate-config

# Run in debug mode
silentcast --log-level debug

# Test a simple spell
# Press your prefix key (e.g., Alt+Space) then 'h' for help
```

## Troubleshooting Installation

### Common Issues

<details>
<summary>Command not found</summary>

The binary is not in your PATH. Either:
- Add the installation directory to PATH
- Move the binary to a directory already in PATH
- Use the full path to run: `/usr/local/bin/silentcast`

</details>

<details>
<summary>Permission denied</summary>

Make the binary executable:
```bash
chmod +x /path/to/silentcast
```

For global hotkeys, you may need to run with elevated permissions initially.

</details>

<details>
<summary>macOS: "Cannot be opened because developer cannot be verified"</summary>

Clear the quarantine attribute:
```bash
xattr -d com.apple.quarantine /usr/local/bin/silentcast
```

Or right-click the app and select "Open" once.

</details>

## Next Steps

Now that SilentCast is installed:

1. [Configure your spellbook](/guide/configuration) - Set up your shortcuts
2. [Learn about shortcuts](/guide/shortcuts) - Master the magic
3. [Explore automation](/guide/scripts) - Advanced scripting features

Need help? Join our [Discord community](https://discord.gg/silentcast) or check the [troubleshooting guide](/guide/troubleshooting).