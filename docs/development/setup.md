# Development Setup Guide

This guide will help you set up a development environment for contributing to SilentCast.

## 🚀 Quick Start

```bash
# Clone the repository
git clone https://github.com/SphereStacking/silentcast.git
cd silentcast

# Set up development environment
make setup

# Build and run (stub mode - no C dependencies)
make build-stub
./app/build/silentcast --no-tray --debug
```

## 📋 Prerequisites

### Required Tools

- **Go**: 1.21 or later
  ```bash
  # Check version
  go version
  
  # Install via package manager
  # macOS
  brew install go
  
  # Ubuntu/Debian
  sudo apt install golang-go
  
  # Windows (Scoop)
  scoop install go
  ```

- **Git**: For version control
  ```bash
  git --version
  ```

- **Make**: Build automation (optional but recommended)
  ```bash
  # macOS (already installed)
  # Linux
  sudo apt install make
  
  # Windows
  scoop install make
  ```

### Optional Tools

- **C Compiler**: Required only for full hotkey support
  - macOS: Xcode Command Line Tools
  - Linux: gcc/g++
  - Windows: MinGW or MSVC

- **golangci-lint**: For code quality checks
  ```bash
  # Install
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
  ```

- **Node.js**: For documentation development
  ```bash
  # For VitePress docs
  npm install -g pnpm
  ```

## 🏗️ Development Modes

### 1. Stub Mode (Recommended for Development)

Fast builds without C dependencies. All features work except global hotkeys.

```bash
# Build stub version
make build-stub

# Or directly with go
cd app
go build -tags "nogohook notray" -o build/silentcast ./cmd/silentcast

# Run
./build/silentcast --no-tray --debug
```

**Advantages:**
- ⚡ Fast compilation (no CGO)
- 🔧 Easy debugging
- 🖥️ Works in any environment
- ✅ All action execution works

**Limitations:**
- ❌ No global hotkey detection
- ❌ No system tray icon

### 2. Full Mode

Complete functionality with hotkey support.

```bash
# Build full version
make build

# Or with go
cd app
CGO_ENABLED=1 go build -o build/silentcast ./cmd/silentcast

# Run
./build/silentcast --debug
```

**Requirements:**
- C compiler installed
- Platform-specific dependencies
- Longer build times

## 🔧 Development Workflow

### 1. Initial Setup

```bash
# Fork and clone
git clone https://github.com/YOUR_USERNAME/silentcast.git
cd silentcast

# Add upstream
git remote add upstream https://github.com/SphereStacking/silentcast.git

# Install dependencies
make setup

# Verify setup
make test
```

### 2. Create Development Branch

```bash
# Update main branch
git checkout main
git pull upstream main

# Create feature branch
git checkout -b feature/your-feature-name
```

### 3. Development Cycle

```bash
# 1. Write tests first (TDD)
cd app
go test ./internal/your-package -run TestYourFeature

# 2. Implement feature
# Edit files...

# 3. Run tests
make test

# 4. Check code quality
make lint

# 5. Format code
make fmt

# 6. Build and test manually
make build-stub
./build/silentcast --no-tray --debug
```

### 4. Testing Your Changes

```bash
# Unit tests
make test

# With coverage
make test-coverage

# Specific package
cd app
go test ./internal/config/... -v

# Single test
go test -run TestConfigLoader_Load ./internal/config

# Integration tests
make test-integration

# Benchmarks
make benchmark
```

## 🧪 Test-Driven Development (TDD)

SilentCast follows TDD practices. Always write tests first!

### TDD Cycle (Red-Green-Refactor)

1. **🔴 Red**: Write a failing test
   ```go
   func TestNewFeature(t *testing.T) {
       result := NewFeature()
       assert.Equal(t, expected, result)
   }
   ```

2. **🟢 Green**: Make it pass with minimal code
   ```go
   func NewFeature() string {
       return "expected"
   }
   ```

3. **🔵 Refactor**: Improve while keeping tests green
   ```go
   func NewFeature() string {
       // Improved implementation
       return processFeature()
   }
   ```

### Testing Guidelines

- Use table-driven tests
- Mock external dependencies
- Test edge cases
- Maintain >80% coverage

## 🏃 Running SilentCast

### Development Mode

```bash
# Run with debug logging
./app/build/silentcast --debug --no-tray

# Test specific spell
./app/build/silentcast --test-spell --spell=e

# Validate configuration
./app/build/silentcast --validate-config

# Dry run
./app/build/silentcast --dry-run --spell=git_status
```

### With Test Configuration

```bash
# Create test spellbook
cat > test-spellbook.yml << EOF
hotkeys:
  prefix: "alt+space"

spells:
  t: "test"

grimoire:
  test:
    type: script
    command: "echo 'Test successful!'"
    show_output: true
EOF

# Run with test config
./app/build/silentcast --config=./test-spellbook.yml --debug
```

## 🐛 Debugging

### VS Code Configuration

Create `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug SilentCast",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/app/cmd/silentcast",
      "args": ["--debug", "--no-tray"],
      "buildFlags": "-tags 'nogohook notray'",
      "env": {
        "CGO_ENABLED": "0"
      }
    },
    {
      "name": "Debug Test",
      "type": "go",
      "request": "launch",
      "mode": "test",
      "program": "${workspaceFolder}/app/internal/config",
      "args": ["-test.run", "TestConfigLoader"]
    }
  ]
}
```

### Debugging Tips

1. **Enable verbose logging**:
   ```bash
   SILENTCAST_LOG_LEVEL=debug ./app/build/silentcast
   ```

2. **Use print debugging**:
   ```go
   log.Printf("DEBUG: variable = %+v", variable)
   ```

3. **Check goroutine leaks**:
   ```go
   defer goleak.VerifyNone(t)
   ```

4. **Profile performance**:
   ```bash
   go test -cpuprofile=cpu.prof -memprofile=mem.prof
   go tool pprof cpu.prof
   ```

## 📦 Project Structure

Understanding the codebase:

```
app/
├── cmd/silentcast/        # Entry point
│   └── main.go           # Main function, CLI parsing
│
├── internal/             # Private packages
│   ├── action/           # Action execution
│   │   ├── executor.go   # Main executor
│   │   ├── script/       # Script execution
│   │   └── app/          # App launching
│   │
│   ├── config/           # Configuration
│   │   ├── loader.go     # YAML loading
│   │   ├── types.go      # Config structs
│   │   └── validator.go  # Validation
│   │
│   ├── hotkey/           # Hotkey management
│   │   ├── manager.go    # Hotkey registration
│   │   └── parser.go     # Key parsing
│   │
│   └── notify/           # Notifications
│       ├── interface.go  # Notifier interface
│       └── console.go    # Console output
│
└── pkg/                  # Public packages
    └── logger/           # Logging utilities
```

## 🛠️ Common Development Tasks

### Adding a New Action Type

1. Define the action in `internal/config/types.go`
2. Create handler in `internal/action/your_type/`
3. Register in `internal/action/executor.go`
4. Add tests in `internal/action/your_type_test.go`
5. Update documentation

### Adding a New CLI Flag

1. Add flag in `cmd/silentcast/main.go`
2. Update help text
3. Implement functionality
4. Add tests
5. Update CLI reference docs

### Platform-Specific Code

Use build tags for platform code:

```go
// +build darwin

package hotkey

// macOS-specific implementation
```

File naming:
- `file_darwin.go` - macOS only
- `file_windows.go` - Windows only
- `file_linux.go` - Linux only

## 🔍 Code Quality

### Before Committing

```bash
# Format code
make fmt

# Run linter
make lint

# Run tests
make test

# Check everything
make check
```

### Code Standards

- Follow Go idioms
- Use meaningful names
- Add comments for exports
- Keep functions small
- Handle errors properly

### Commit Messages

Use conventional commits:

```
feat: add new spell type
fix: resolve config loading issue
docs: update installation guide
test: add config validation tests
refactor: simplify action executor
```

## 📚 Resources

### Internal Documentation

- [Architecture Overview](../api/architecture.md)
- [Testing Guide](./testing.md)
- [Error Handling](./error-handling.md)
- [TDD Best Practices](./tdd-best-practices.md)

### External Resources

- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Testify Documentation](https://github.com/stretchr/testify)

## 🤝 Getting Help

- **GitHub Issues**: For bugs and features
- **Discussions**: For questions
- **Discord**: Real-time chat (if available)

### Common Issues

**Build fails with CGO errors**:
```bash
# Use stub build instead
make build-stub
```

**Tests timeout in CI**:
```bash
# Increase timeout
go test -timeout 30s ./...
```

**Hotkeys not working in development**:
- This is normal in stub mode
- Use `--test-spell` for testing
- Or build full version with `make build`

## 🎯 Next Steps

1. Read [Contributing Guide](../../CONTRIBUTING.md)
2. Check [open issues](https://github.com/SphereStacking/silentcast/issues)
3. Join discussions
4. Start coding!

---

<div align="center">
  <p><strong>Happy coding! 🚀</strong></p>
</div>