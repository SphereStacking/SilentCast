# SilentCast

<div align="center">
  <img src="https://spherestacking.github.io/SilentCast/logo.svg" alt="SilentCast Logo" width="200" height="200">
</div>

<p align="center">
  <a href="https://github.com/SphereStacking/silentcast/actions/workflows/ci.yml"><img src="https://github.com/SphereStacking/silentcast/actions/workflows/ci.yml/badge.svg" alt="CI"></a>
  <a href="https://github.com/SphereStacking/silentcast/releases"><img src="https://img.shields.io/github/v/release/SphereStacking/silentcast" alt="Release"></a>
  <a href="LICENSE"><img src="https://img.shields.io/github/license/SphereStacking/silentcast" alt="License"></a>
</p>

<p align="center">
  <a href="README.md">English</a> | <a href="README.ja.md">日本語</a>
</p>

> ⚠️ **Development Status**: This project is currently under active development. Features may change and bugs may exist. Use at your own risk.

<p align="center">
  SilentCast is a silent hotkey-driven task runner that lets you execute tasks with simple keyboard shortcuts.<br>
  Perfect for developers who want to streamline their workflow without leaving the keyboard.
</p>

## ✨ Features

- 🎯 **Global Hotkeys** - Works system-wide, no need to focus the app
- 🏃 **Fast Execution** - Instant application launching and script execution
- 📝 **Sequential Keys** - VS Code-style key combinations (e.g., `g,s` for git status)
- 🌐 **Cross-platform** - Works on Windows and macOS
- 🎨 **Intuitive** - Simple configuration with spells and actions
- 🧪 **Lightweight** - Minimal CPU and memory usage (~15MB)
- 🔄 **Auto-reload** - Configuration changes are applied automatically
- 📋 **System Tray** - Unobtrusive system tray integration
- 📊 **Structured Logging** - Comprehensive logs with rotation support

## 📦 Installation

### Download Binary

Download the latest release for your platform from the [releases page](https://github.com/SphereStacking/spellbook/releases).

#### macOS
```bash
# Intel
curl -L https://github.com/SphereStacking/silentcast/releases/latest/download/silentcast-darwin-amd64.tar.gz | tar xz
sudo mv silentcast/silentcast /usr/local/bin/

# Apple Silicon
curl -L https://github.com/SphereStacking/silentcast/releases/latest/download/silentcast-darwin-arm64.tar.gz | tar xz
sudo mv silentcast/silentcast /usr/local/bin/
```

#### Windows

Download the ZIP file from the releases page and extract it to a directory in your PATH.

### From Source

```bash
# Clone the repository
git clone https://github.com/SphereStacking/silentcast.git
cd silentcast

# Build and install
make install

# Or build for current platform only
make build
```

## 🔮 Configuration

SilentCast uses YAML configuration files:

- **Spells** - Your keyboard shortcuts
- **Grimoire** - The actual commands/applications to run
- **Command** - The command or path to execute

### Basic Configuration

Create a `spellbook.yml` file:

```yaml
# SilentCast Configuration
daemon:
  auto_start: false
  log_level: info
  config_watch: true

logger:
  level: info
  file: ""                 # Empty = console only
  max_size: 10             # MB
  max_backups: 3
  max_age: 7               # days
  compress: false

hotkeys:
  prefix: "alt+space"      # Your magic key
  timeout: 1000            # ms to wait after prefix
  sequence_timeout: 2000   # ms for full sequence

spells:
  # Single key spells
  e: "editor"
  t: "terminal"
  b: "browser"

  # Multi-key sequences (VS Code style)
  "g,s": "git_status"
  "g,p": "git_pull"
  "g,c": "git_commit"

grimoire:
  editor:
    type: app
    command: "code"    # VS Code
    description: "Open VS Code"

  terminal:
    type: app
    command: "wt"      # Windows Terminal
    description: "Open Terminal"

  git_status:
    type: script
    command: "git status"
    description: "Show git status"
```

### Platform-specific Configuration

SilentCast supports platform-specific overrides:

- `spellbook.yml` - Base configuration (loaded first)
- `spellbook.mac.yml` - macOS overrides
- `spellbook.windows.yml` - Windows overrides

Example `spellbook.mac.yml`:
```yaml
grimoire:
  terminal:
    type: app
    command: "Terminal"

  browser:
    type: app
    command: "Safari"
```

### Example Configurations

Example configuration files can be found in the `examples/config/` directory:
- `spellbook.yml` - Full-featured example with common shortcuts
- `spellbook.windows.yml` - Windows-specific overrides
- `spellbook.mac.yml` - macOS-specific overrides

See [CONFIG.md](CONFIG.md) for detailed configuration guide.

## 🎮 Usage

### Starting SilentCast

```bash
# Run with default configuration
silentcast

# Run without system tray
silentcast --no-tray

# Run with custom config location
SILENTCAST_CONFIG=/path/to/config silentcast
```

### Casting Spells

1. Press your prefix key (default: `Alt+Space`)
2. Cast a spell:
   - **Single key**: Press `e` to open your editor
   - **Sequence**: Press `g`, then `s` for git status
   - **Long sequence**: Press `d`, then `o`, then `c` for documentation

### System Tray

When running with system tray support:
- **Show Hotkeys** - Display all configured shortcuts
- **Reload Config** - Manually reload configuration
- **About** - Show version information
- **Quit** - Exit SilentCast

## 💻 Platform Support

### Windows
- ✅ Full global hotkey support
- ✅ System tray integration
- ✅ No special permissions required

### macOS
- ✅ Full global hotkey support
- ⚠️ Requires accessibility permissions
- ✅ System tray integration
- 📝 First run: Grant permissions in System Preferences > Security & Privacy > Accessibility


## 🔧 Development

### Prerequisites

- Go 1.21 or later
- Make (optional but recommended)
- Python 3.x (for pre-commit hooks)

## 📁 Project Structure

```
SilentCast/
├── app/              # Application source code
│   ├── cmd/          # Main entry point
│   ├── internal/     # Internal packages
│   ├── pkg/          # Public packages
│   └── Makefile      # Build configuration
├── docs/             # VitePress documentation
│   ├── guide/        # User guide
│   ├── config/       # Configuration reference
│   └── api/          # Developer documentation
├── examples/         # Example configurations
└── README.md         # This file
```

### Quick Start

```bash
# Clone the repository
git clone https://github.com/SphereStacking/silentcast.git
cd silentcast

# Setup development environment
make setup

# Install pre-commit hooks (recommended)
make pre-commit

# Build for development (fast)
make build-dev

# Run
./app/build/silentcast --no-tray
```

### Full Documentation

```bash
# Start documentation server
make docs-dev
# Open http://localhost:5173
```

See [docs/api/build.md](docs/api/build.md) for detailed build instructions.

### Available Commands

```bash
# Application
make build-dev     # Development build (fast, no hotkeys)
make build         # Production build (requires C libs)
make build-snapshot # Snapshot build for all platforms
make lint          # Run linting checks
make test          # Run tests
make clean         # Clean build artifacts

# Documentation
make docs-dev      # Start docs dev server
make docs-build    # Build documentation

# Setup
make setup         # Setup development environment
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test -v ./internal/config/
```

### Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-spell`)
3. Commit your changes (`git commit -m 'Add amazing spell'`)
4. Push to the branch (`git push origin feature/amazing-spell`)
5. Open a Pull Request

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [gohook](https://github.com/robotn/gohook) - Global hotkey support
- [systray](https://github.com/getlantern/systray) - System tray integration
- [fsnotify](https://github.com/fsnotify/fsnotify) - File system notifications
- [lumberjack](https://github.com/natefinch/lumberjack) - Log rotation

---

<p align="center">
  Made with 🪄 by developers who ❤️ keyboard shortcuts
</p>
