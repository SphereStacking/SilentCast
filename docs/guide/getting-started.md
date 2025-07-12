# Getting Started

:::warning Development Status
SilentCast is currently under active development. Features may change and bugs may exist. Please use at your own risk and report any issues on [GitHub](https://github.com/SphereStacking/silentcast/issues).
:::

This guide will help you set up SilentCast and start using keyboard shortcuts in less than 5 minutes.

## Prerequisites

Before we begin, make sure you have:
- A supported operating system (Windows 10+ or macOS 10.15+)
- Admin/sudo access for installation
- 5 minutes of your time

## Quick Install

<div class="install-tabs">

::: code-group

```bash [macOS]
# Quick install script
curl -sSL https://get.silentcast.dev | bash

# Or with Homebrew (macOS)
brew install silentcast

# Or download directly
curl -L https://github.com/SphereStacking/silentcast/releases/latest/download/silentcast-$(uname -s)-$(uname -m).tar.gz | tar xz
sudo mv silentcast /usr/local/bin/
```

```powershell [Windows]
# Quick install script (PowerShell as Admin)
iwr -useb https://get.silentcast.dev/install.ps1 | iex

# Or with Scoop
scoop install silentcast

# Or with Chocolatey
choco install silentcast

# Or download the .exe from releases
# https://github.com/SphereStacking/silentcast/releases
```

```bash [From Source]
# Clone and build
git clone https://github.com/SphereStacking/silentcast.git
cd silentcast
make build
sudo make install
```

:::

</div>

## First Run

### 1. Start SilentCast

```bash
# Start with system tray icon
silentcast

# Or start without tray (terminal mode)
silentcast --no-tray
```

<div class="tip custom-block">

**macOS Users**: On first run, you'll need to grant Accessibility permissions:
1. Go to System Preferences → Security & Privacy → Privacy → Accessibility
2. Click the lock and add SilentCast
3. Restart SilentCast

</div>

### 2. Verify Installation

Press `Alt+Space` (default prefix), then `h` to see the built-in help. You should see a notification or terminal output showing available shortcuts.

## Your First Spell

Let's create a simple configuration to open your favorite editor:

### 1. Create Configuration File

Create `~/.config/silentcast/spellbook.yml`:

```yaml
# My first spellbook
daemon:
  auto_start: false
  log_level: info

hotkeys:
  prefix: "alt+space"    # Your activation key
  timeout: 1000          # Milliseconds to wait

spells:
  e: "open_editor"       # Alt+Space, then E
  t: "open_terminal"     # Alt+Space, then T
  b: "open_browser"      # Alt+Space, then B

grimoire:
  open_editor:
    type: app
    command: "code"      # Change to your editor
    description: "Open VS Code"
  
  open_terminal:
    type: app
    command: "wt"        # Windows Terminal
    # command: "terminal" # macOS
    description: "Open Terminal"
  
  open_browser:
    type: app
    command: "chrome"
    description: "Open Chrome"
```

### 2. Test Your Spells

1. Save the configuration file
2. SilentCast will auto-reload (or restart it)
3. Press `Alt+Space`, then `e` - Your editor should open!
4. Press `Alt+Space`, then `t` - Your terminal appears!

## Level Up: Multi-Key Sequences

Add these VS Code-style sequences to your spellbook:

```yaml
spells:
  # Git shortcuts
  "g,s": "git_status"
  "g,p": "git_pull"
  "g,c": "git_commit"
  
  # Docker shortcuts
  "d,u": "docker_up"
  "d,d": "docker_down"
  "d,l": "docker_logs"

grimoire:
  git_status:
    type: script
    command: "git status"
    working_dir: "${PWD}"  # Use current directory
    
  docker_up:
    type: script
    command: "docker-compose up -d"
    
  docker_logs:
    type: script
    command: "docker-compose logs -f"
```

Now try `Alt+Space`, `g`, `s` to check git status!

## Platform-Specific Configuration

Create platform-specific overrides:

::: code-group

```yaml [spellbook.mac.yml]
# macOS specific overrides
grimoire:
  open_terminal:
    command: "Terminal"
  
  open_browser:
    command: "Safari"
  
  # macOS specific shortcut
  show_desktop:
    type: script
    command: "osascript -e 'tell application \"System Events\" to key code 103 using {command down, shift down}'"
```

```yaml [spellbook.windows.yml]
# Windows specific overrides
grimoire:
  open_terminal:
    command: "wt"  # Windows Terminal
  
  open_browser:
    command: "msedge"
  
  # Windows specific
  task_manager:
    type: app
    command: "taskmgr"
```

    command: "gnome-terminal"
    # command: "konsole"  # KDE
    # command: "xfce4-terminal"  # XFCE
  
  open_browser:
    command: "firefox"
```

:::

## Essential Tips

### 🎯 Choosing Good Shortcuts

- **Single letters** for most-used apps: `e` (editor), `t` (terminal)
- **Grouped sequences** for related tasks: `g,*` for git, `d,*` for docker
- **Memorable combinations**: `d,b` for "docker build"
- **Avoid conflicts**: Check your IDE shortcuts first

### ⚡ Performance Tips

1. **Start with system**: Enable auto-start for instant availability
2. **Minimal spells**: Only add shortcuts you'll actually use
3. **Fast sequences**: Set lower timeouts for faster response

### 🔍 Debugging

If shortcuts aren't working:

```bash
# Run in debug mode
silentcast --log-level debug

# Check logs
tail -f ~/.local/share/silentcast/silentcast.log

# Verify configuration
silentcast --validate-config
```

## What's Next?

Congratulations! You're now a SilentCast wizard! 🧙‍♂️ Here's where to go next:

<div class="next-cards">

<div class="card">

### 📦 [Installation Guide](/guide/installation)
Platform-specific installation details and troubleshooting

</div>

<div class="card">

### ⚙️ [Configuration Guide](/guide/configuration) 
Deep dive into all configuration options

</div>

<div class="card">

### 🎮 [Shortcuts Guide](/guide/shortcuts)
Advanced shortcut patterns and best practices

</div>

<div class="card">

### 🤖 [Automation Guide](/guide/scripts)
Script execution, environment variables, and more

</div>

</div>

## Join the Community

Having issues or want to share your spellbook?

- 💬 [Discord Community](https://discord.gg/silentcast)
- 🐛 [Report Issues](https://github.com/SphereStacking/silentcast/issues)
- 📚 [Example Configurations](https://github.com/SphereStacking/silentcast/tree/main/examples)

<style>
.install-tabs {
  margin: 2rem 0;
}

.next-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1rem;
  margin: 2rem 0;
}

.card {
  padding: 1.5rem;
  border: 1px solid var(--vp-c-divider);
  border-radius: 8px;
  transition: all 0.3s;
}

.card:hover {
  border-color: var(--vp-c-brand);
  transform: translateY(-2px);
}

.card h3 {
  margin: 0 0 0.5rem 0;
  font-size: 1.1rem;
}

.card p {
  margin: 0;
  color: var(--vp-c-text-2);
  font-size: 0.95rem;
}
</style>