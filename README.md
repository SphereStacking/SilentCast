# SilentCast

<div align="center">
  <img src="https://spherestacking.github.io/SilentCast/logo.svg" alt="SilentCast Logo" width="200" height="200">
  
  <h3>🪄 Cast spells, run tasks</h3>
  
  <p>A silent hotkey-driven task runner that lets you execute tasks with simple keyboard spells</p>
</div>

<p align="center">
  <a href="https://github.com/SphereStacking/silentcast/actions/workflows/ci.yml"><img src="https://github.com/SphereStacking/silentcast/actions/workflows/ci.yml/badge.svg" alt="CI"></a>
  <a href="https://github.com/SphereStacking/silentcast/releases"><img src="https://img.shields.io/github/v/release/SphereStacking/silentcast" alt="Release"></a>
  <a href="https://goreportcard.com/report/github.com/SphereStacking/silentcast"><img src="https://goreportcard.com/badge/github.com/SphereStacking/silentcast" alt="Go Report Card"></a>
  <a href="LICENSE"><img src="https://img.shields.io/github/license/SphereStacking/silentcast" alt="License"></a>
  <a href="https://pkg.go.dev/github.com/SphereStacking/silentcast"><img src="https://pkg.go.dev/badge/github.com/SphereStacking/silentcast.svg" alt="Go Reference"></a>
</p>

<p align="center">
  <a href="README.md">English</a> | <a href="README.ja.md">日本語</a> | <a href="https://spherestacking.github.io/SilentCast/">Documentation</a>
</p>

---

## 🌟 What is SilentCast?

SilentCast is a lightweight, cross-platform application that runs silently in the background, waiting for your keyboard spells to execute predefined tasks. Whether you're a developer, system administrator, or power user, SilentCast helps you automate repetitive tasks with simple keyboard shortcuts.

### ✨ Key Features

- **🎯 Global Hotkeys** - Works anywhere, no window focus required
- **⚡ Lightning Fast** - Instant task execution with minimal resource usage
- **🔮 Magic Terminology** - Uses spells (shortcuts) and grimoire (actions)
- **🎹 VS Code-style Sequences** - Multi-key combinations like `g,s` for git status
- **🌍 Cross-Platform** - Native support for Windows, macOS, and Linux
- **🔄 Live Configuration** - Changes apply instantly without restart
- **📊 Smart Output** - Show command results in notifications or terminal
- **🔐 Elevated Execution** - Run tasks with admin privileges when needed
- **🧪 Developer Friendly** - Comprehensive CLI tools for testing and debugging

## 🚀 Quick Start

### Installation

#### Using Package Manager

```bash
# macOS (Homebrew)
brew install spherestacking/tap/silentcast

# Windows (Scoop)
scoop bucket add spherestacking https://github.com/spherestacking/scoop-bucket
scoop install silentcast

# Linux (Snap)
sudo snap install silentcast
```

#### Direct Download

Download the latest binary for your platform from the [releases page](https://github.com/SphereStacking/silentcast/releases).

```bash
# Example for macOS (Apple Silicon)
curl -L https://github.com/SphereStacking/silentcast/releases/latest/download/silentcast-darwin-arm64.tar.gz | tar xz
sudo mv silentcast /usr/local/bin/
```

### Your First Spell

1. Create a configuration file `spellbook.yml`:

```yaml
# Basic spellbook configuration
hotkeys:
  prefix: "alt+space"      # Your activation key

spells:
  e: "editor"              # Alt+Space, then E
  t: "terminal"            # Alt+Space, then T
  "g,s": "git_status"      # Alt+Space, then G, then S

grimoire:
  editor:
    type: app
    command: "code"        # Opens VS Code
    
  terminal:
    type: app
    command: "wt"          # Opens Windows Terminal
    
  git_status:
    type: script
    command: "git status"
    show_output: true      # Shows result in notification
```

2. Start SilentCast:

```bash
silentcast
```

3. Cast your first spell:
   - Press `Alt+Space` (your prefix key)
   - Press `e` to open your editor
   - Or press `g`, then `s` to see git status

## 🎮 Usage Examples

### Basic Commands

```bash
# Start SilentCast
silentcast                           # Run with system tray
silentcast --no-tray                 # Run without system tray
silentcast --debug                   # Enable debug logging

# Configuration Management
silentcast --validate-config         # Check config syntax
silentcast --show-config             # Display merged configuration
silentcast --list-spells             # Show all available spells

# Testing & Debugging
silentcast --test-hotkey             # Test hotkey detection

# Help
silentcast --help                    # Show help
```

### Advanced Configuration

```yaml
grimoire:
  # Show command output
  docker_ps:
    type: script
    command: "docker ps"
    show_output: true
    description: "List Docker containers"
    
  # Keep terminal open after execution
  python_shell:
    type: script
    command: "python"
    terminal: true
    keep_open: true
    description: "Interactive Python shell"
    
  # Run with elevated privileges
  system_update:
    type: script
    command: "apt update && apt upgrade -y"
    admin: true
    terminal: true
    description: "Update system packages"
    
  # Open URLs
  github_profile:
    type: url
    command: "https://github.com/{{.Username}}"
    description: "Open GitHub profile"
    
  # Custom shell and timeout
  long_process:
    type: script
    command: "./backup.sh"
    shell: "bash"
    timeout: 300
    show_output: true
    description: "Run backup with 5 minute timeout"
```

## 📚 Documentation

### User Guides
- [Getting Started](https://spherestacking.github.io/SilentCast/guide/getting-started)
- [Configuration Guide](https://spherestacking.github.io/SilentCast/guide/configuration)
- [Spells & Grimoire](https://spherestacking.github.io/SilentCast/guide/spells)
- [Platform Setup](https://spherestacking.github.io/SilentCast/guide/platforms)

### Reference
- [CLI Reference](https://spherestacking.github.io/SilentCast/guide/cli-reference)
- [Configuration Schema](https://spherestacking.github.io/SilentCast/config/)
- [Troubleshooting](https://spherestacking.github.io/SilentCast/troubleshooting/)

### Developer Resources
- [API Documentation](https://pkg.go.dev/github.com/SphereStacking/silentcast)
- [Architecture Guide](https://spherestacking.github.io/SilentCast/api/architecture)
- [Contributing](https://spherestacking.github.io/SilentCast/contributing)

## 💻 Platform Support

| Platform | Global Hotkeys | System Tray | Notifications | Admin/Sudo | Auto-start |
|----------|----------------|-------------|---------------|------------|------------|
| Windows  | ✅ | ✅ | ✅ (Native) | ✅ | ✅ |
| macOS    | ✅* | ✅ | ✅ (Native) | ✅ | ✅ |
| Linux    | ✅ | ✅** | ✅ (Multiple) | ✅ | ✅ |

\* macOS requires accessibility permissions on first run  
\** Linux requires `libappindicator3-1` for system tray

## 🔧 Development

### Prerequisites
- Go 1.23 or later
- Make (optional but recommended)
- C compiler (for production builds with hotkey support)

### Building from Source

```bash
# Clone repository
git clone https://github.com/SphereStacking/silentcast.git
cd silentcast

# Setup development environment
make setup

# Build options
make build-dev      # Fast build without hotkey support (for development)
make build          # Production build with full features
make build-all      # Build for all platforms

# Run tests
make test           # Unit tests
make test-all       # All tests including integration
make benchmark      # Performance benchmarks

# Development workflow
make lint           # Run linters
make fmt            # Format code
make docs-dev       # Start documentation server
```

### Project Structure

```
SilentCast/
├── app/                    # Application source code
│   ├── cmd/                # Main entry point
│   ├── internal/           # Internal packages
│   │   ├── action/         # Action execution
│   │   ├── config/         # Configuration management
│   │   ├── hotkey/         # Hotkey detection
│   │   └── notify/         # Notification system
│   └── pkg/                # Public packages
├── docs/                   # Documentation (VitePress)
├── examples/               # Example configurations
└── .ticket/                # Ticket-based development system
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Quick Contribution Guide

1. Check existing [issues](https://github.com/SphereStacking/silentcast/issues) and [tickets](.ticket/README.md)
2. Fork the repository
3. Create a feature branch (`git checkout -b feature/amazing-spell`)
4. Follow our coding standards and use magic terminology
5. Write tests for new features
6. Submit a pull request

### Development Philosophy

- **Magic Terminology**: We use spells, grimoire, and spellbook consistently
- **Test-Driven Development**: Write tests first, then implementation
- **Clean Architecture**: Clear separation of concerns
- **User Experience First**: Simple for users, powerful for developers

## 📊 Performance

SilentCast is designed to be lightweight and efficient:

- **Memory Usage**: ~15MB idle, ~25MB active
- **CPU Usage**: <0.1% idle, <1% during execution
- **Startup Time**: <100ms
- **Hotkey Response**: <10ms

See our [Performance Guide](docs/performance/README.md) for optimization tips.

## 🔒 Security

- No network connections except for self-update checks (optional)
- All configurations are local
- Admin/sudo execution requires explicit configuration
- No telemetry or data collection

Report security issues to: security@spherestacking.com

## 📄 License

SilentCast is open source software licensed under the [MIT License](LICENSE).

## 🙏 Acknowledgments

Built with these excellent libraries:
- [gohook](https://github.com/robotn/gohook) - Cross-platform hotkey support
- [systray](https://github.com/getlantern/systray) - System tray integration
- [fsnotify](https://github.com/fsnotify/fsnotify) - File watching
- [lumberjack](https://github.com/natefinch/lumberjack) - Log rotation
- [yaml.v3](https://github.com/go-yaml/yaml) - YAML configuration

<!-- ## 🌟 Star History

[![Star History Chart](https://api.star-history.com/svg?repos=SphereStacking/silentcast&type=Date)](https://star-history.com/#SphereStacking/silentcast&Date) -->

---

<div align="center">
  <p>Made with 🪄 by developers who ❤️ keyboard magic</p>
  
  <p>
    <a href="https://github.com/SphereStacking/silentcast/issues/new?labels=bug">Report Bug</a>
    ·
    <a href="https://github.com/SphereStacking/silentcast/issues/new?labels=enhancement">Request Feature</a>
    ·
    <a href="https://spherestacking.github.io/SilentCast/">Read Docs</a>
  </p>
</div>
