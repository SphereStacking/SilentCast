# CLAUDE.md

This file provides guidance for Claude Code (claude.ai/code) when working with this repository's code.

## 🎯 Project Overview

**SilentCast** is a hotkey-driven task runner that executes tasks via keyboard shortcuts.
- Press **prefix key** (default: Alt+Space) followed by configured **spells**
- Supports VS Code-style sequential key input (e.g., `g,s` for git status)
- Cross-platform support (Windows/macOS)

### Magic Terminology
- **Spells**: Keyboard shortcuts (`e` = editor, `g,s` = git status)
- **Grimoire**: Action definitions (what spells execute)
- **Spellbook**: Configuration file (`spellbook.yml`)

## 📁 Project Structure

```
SilentCast/
├── app/                    # Application core (Go code)
│   ├── cmd/silentcast/     # Entry point (main.go)
│   ├── internal/           # Internal packages
│   │   ├── action/         # Action execution (app launch, script run)
│   │   ├── config/         # Configuration management (YAML loading, cascade)
│   │   ├── hotkey/         # Hotkey detection (gohook integration)
│   │   ├── permission/     # OS-specific permission management
│   │   └── tray/           # System tray
│   ├── pkg/logger/         # Logging utilities
│   └── Makefile            # Build configuration
├── docs/                   # VitePress documentation
├── examples/config/        # Configuration file samples
└── Makefile                # Top-level commands
```

## 🚀 Quick Start

```bash
# Setup development environment
make setup

# Development build (no hotkey functionality, no C dependencies)
make build-stub

# Run
./app/build/silentcast --no-tray

# Run tests
make test
```

## 🔧 Development Commands

### Build Commands
```bash
# Development build (stub mode) - recommended
make build-stub              # No C language libraries required, all features work except hotkeys

# Production build
make build                   # Full version with hotkey functionality

# All platforms build
make build-all               # For CI/CD, static binaries

# Direct run (no build)
cd app && go run cmd/silentcast/main.go
```

### Testing & Quality Control
```bash
# Run tests (using stub implementation)
make test                    # Run all tests
make test-coverage           # Generate coverage report

# Test specific packages
cd app
go test ./internal/config/... -v
go test ./internal/action/... -v
go test ./internal/hotkey/... -v

# Run single test
go test -run TestExecutor_Execute ./internal/action/...

# Code quality
make fmt                     # Format
make lint                    # Lint (golangci-lint)
```

### Documentation
```bash
make docs-dev                # View docs at http://localhost:5173
make docs-build              # Generate static site
```

## 🏗️ Architecture Details

### Core Components

#### 1. **Configuration System** (`internal/config/`)
- **Cascade loading**: `spellbook.yml` → `spellbook.{os}.yml`
- **File watching**: Auto-reload on config changes (`watcher.go`)
- **Path resolution**: OS-specific application path resolution (`resolver_{os}.go`)
- **Note**: `Config.UnmarshalYAML` tracks explicit empty string settings

#### 2. **Hotkey Management** (`internal/hotkey/`)
- **Implementation switching**: Controlled by build tags
  - `manager.go`: Uses gohook (requires CGO)
  - `manager_stub.go`: Mock implementation (`//go:build nogohook`)
- **Key parsing**: Parses "alt+space", "g,s" formats
- **Timeouts**: 1000ms after prefix, 2000ms for entire sequence

#### 3. **Action Execution** (`internal/action/`)
- **Executor**: Routes by type (app/script)
- **App Launcher**: OS-specific implementation (`launcher_{os}.go`)
- **Script Runner**: Shell command execution, environment variable expansion

#### 4. **Permission Management** (`internal/permission/`)
- **Interface design**: `Manager` interface
- **OS-specific requirements**:
  - macOS: Accessibility permissions
  - Windows: No permissions required

### Build System

#### Build Tags
- `nogohook`: Disables gohook (no hotkey functionality)
- `notray`: Disables systray (no tray functionality)

#### CGO Configuration
```bash
# Normal build (CGO required)
CGO_ENABLED=1 go build

# Stub build (no CGO)
go build -tags "nogohook notray"
```

## 🐛 Troubleshooting

### Common Issues

1. **Build error: CGO related**
   ```bash
   # Solution: Use stub build
   make build-stub
   ```

2. **Test failure: `TestIsNewerVersion`**
   - Bug in version comparison logic in `internal/updater/updater_test.go`
   - v1.9.0 < v1.10.0 comparison fails (due to string comparison)

3. **Hotkeys not working**
   - Hotkey functionality doesn't work in stub builds
   - Full version can be built via GitHub Actions

### Debugging
```bash
# Check logs
tail -f spellbook.log

# Run in debug mode
./app/build/silentcast --debug --no-tray

# Validate configuration
./app/build/silentcast --validate-config
```

## 📚 Important Implementation Patterns

### Error Handling
```go
// Always wrap with context
if err != nil {
    return fmt.Errorf("failed to load config: %w", err)
}
```

### Table-Driven Tests
```go
tests := []struct {
    name    string
    input   string
    want    string
    wantErr bool
}{
    {"valid key", "alt+space", "alt+space", false},
    // ...
}
```

### Interface Design
```go
// Always interface OS-dependent parts
type Manager interface {
    CheckPermission(PermissionType) Status
    RequestPermission(PermissionType) error
}
```

## 🔗 Related Documentation

- **User Guide**: `docs/guide/`
- **Configuration Reference**: `docs/config/`
- **API Documentation**: `docs/api/`
- **Planning Documents**: `docs/planning/`

## 💡 Development Tips

1. **Develop in stub mode**: All features work except hotkeys, fast builds
2. **Use GitHub Actions**: After push, download full binaries from Artifacts
3. **Test configurations**: Use samples in `examples/config/`
4. **Platform-specific implementation**: Separate with `_darwin.go`, `_windows.go`

---

**Module Path**: `github.com/SphereStacking/silentcast`


## Git Operation Rules

### Branch Strategy
- **Main branch**: `main` - Always deployable to production
- **Feature branches**: `feature/feature-name` or `claude/issue-number` - For new features
- **Fix branches**: `fix/fix-description` - For bug fixes

### How to Incorporate Main Branch Changes
To keep history clean, use `rebase`:

```bash
# 1. Commit current changes
git add .
git commit -m "work description"

# 2. Fetch latest main
git fetch origin main

# 3. Rebase to incorporate main branch changes
git rebase origin/main

# 4. Resolve conflicts if any and continue
git add .
git rebase --continue

# 5. Push to remote (first time or when history changed)
git push --force-with-lease origin branch-name
```

### Pull Request Guidelines
- Never push directly to main branch
- Always rebase with latest main before creating PR
- Keep commit history logical (use `git rebase -i` to clean up if needed)

### Git Commit Rules

This project uses emoji-prefixed English commit messages:

#### Common Emojis
- ✨ Add new feature
- 🎨 Improve UI/styling
- 📝 Add/update documentation
- 🐛 Fix bug
- ♻️ Refactor code
- 🔧 Change configuration files
- ➕ Add dependency
- ➖ Remove dependency
- 🚚 Move/rename files
- 🔥 Remove code or files
- ⚡ Improve performance
- 🔒 Fix security issues
- 🚧 Work in progress
- ✅ Add/update tests

#### Commit Message Examples
```
✨ Add shadcn-vue integration
🎨 Improve responsive design
📝 Add Git branch strategy documentation
```