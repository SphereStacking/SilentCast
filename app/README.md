# SilentCast Application

The core application code for SilentCast - a lightweight, system-wide hotkey driven application/script launcher written in Go.

## 📁 Project Structure

```
app/
├── cmd/silentcast/     # Main entry point
│   ├── main.go         # Application entry
│   ├── validate.go     # --validate-config implementation
│   └── show.go         # --show-config implementation
├── internal/           # Internal packages
│   ├── action/         # Action execution
│   │   ├── app.go      # Application launcher
│   │   ├── script.go   # Script executor with output capture
│   │   ├── url.go      # URL opener
│   │   └── elevated.go # Admin/elevated execution
│   ├── config/         # Configuration management
│   ├── hotkey/         # Hotkey detection
│   ├── notify/         # Notification system
│   │   ├── console.go  # Console output
│   │   ├── windows.go  # Windows notifications
│   │   ├── darwin.go   # macOS notifications
│   │   └── linux.go    # Linux notifications
│   ├── output/         # Output management
│   │   └── buffered.go # Output capture implementation
│   ├── permission/     # OS permissions
│   └── tray/           # System tray
├── pkg/logger/         # Logging utilities
└── Makefile            # Build configuration
```

## 🚀 Building

### Development Build (No CGO)
```bash
make build-stub
# Fast build without hotkey functionality
# All features work except global hotkeys
```

### Production Build
```bash
make build
# Full build with hotkey support
# Requires CGO and platform-specific libraries
```

### Cross-Platform Build
```bash
make build-all
# Builds for all supported platforms
# Used in CI/CD pipeline
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Test specific package
go test ./internal/action/... -v
```

## 🔧 CLI Options

```bash
# Basic usage
silentcast [options]

# Options:
--no-tray           # Disable system tray
--validate-config   # Validate configuration and exit
--show-config       # Show merged configuration
--format=json       # Output format (human, json, yaml)
--show-paths        # Show config search paths
--version           # Show version
```

## 📦 Dependencies

### Core Dependencies
- `github.com/robotn/gohook` - Global hotkey support
- `github.com/getlantern/systray` - System tray
- `gopkg.in/yaml.v3` - YAML configuration

### Build Dependencies
- Go 1.21+
- CGO (for production builds)
- Platform-specific:
  - Windows: MinGW-w64
  - macOS: Xcode Command Line Tools
  - Linux: gcc, libX11-dev

## 🏗️ Architecture

### Action System
- **Executor Interface**: Common interface for all action types
- **Action Types**:
  - `app`: Launch applications
  - `script`: Execute shell commands
  - `url`: Open URLs in browser
- **Features**:
  - Output capture and display
  - Terminal control (keep_open)
  - Timeout support
  - Admin/elevated execution
  - Custom shell selection

### Notification System
- **Manager**: Routes notifications to available notifiers
- **Notifiers**:
  - Console (always available)
  - System (platform-specific)
- **Fallback Chain**: System → Console

### Configuration
- **Cascade Loading**: Base → OS-specific override
- **File Watching**: Auto-reload on changes
- **Validation**: Type checking and constraint validation

## 🐛 Debugging

### Enable Debug Logging
```yaml
logger:
  level: debug
  file: silentcast.log
```

### View Key Events
```bash
# Debug logging shows all key events
tail -f silentcast.log | grep "Key pressed"
```

### Test Notifications
```yaml
grimoire:
  test_notify:
    type: script
    command: "echo 'Test notification'"
    show_output: true
```