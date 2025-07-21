# Building SilentCast

This guide covers building SilentCast from source for development and distribution.

## Prerequisites

### Required Tools

- **Go 1.23+** - [Download](https://golang.org/dl/)
- **Git** - For cloning the repository
- **Make** - For build automation (optional but recommended)

### Platform-Specific Requirements

#### Windows
- **MinGW-w64** or **MSYS2** for CGO support
- **Windows SDK** for native APIs

#### macOS
- **Xcode Command Line Tools**
  ```bash
  xcode-select --install
  ```

#### Linux
- **GCC** and development headers
  ```bash
  # Ubuntu/Debian
  sudo apt-get install build-essential libx11-dev libxtst-dev
  
  # Fedora/RHEL
  sudo dnf install gcc libX11-devel libXtst-devel
  ```

## Quick Build

### Using Make

```bash
# Clone repository
git clone https://github.com/SphereStacking/silentcast.git
cd silentcast

# Build for current platform
make build

# Run
./app/build/silentcast
```

### Using Go Directly

```bash
cd app
go build -o ../build/silentcast ./cmd/silentcast
```

## Build Modes

### Development Build (Stub)

No C dependencies, faster builds, but no hotkey functionality:

```bash
make build-stub
# or
go build -tags "nogohook notray" ./cmd/silentcast
```

### Production Build

Full functionality with hotkey support:

```bash
make build
# or
CGO_ENABLED=1 go build ./cmd/silentcast
```

### Static Build

No external dependencies:

```bash
make build-static
# or
CGO_ENABLED=0 go build -ldflags="-s -w" ./cmd/silentcast
```

## Cross-Platform Building

### Build for All Platforms

```bash
make build-all
```

This creates:
- `silentcast-windows-amd64.exe`
- `silentcast-darwin-amd64`
- `silentcast-darwin-arm64`
- `silentcast-linux-amd64`

### Manual Cross-Compilation

```bash
# Windows from macOS/Linux
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 \
  CC=x86_64-w64-mingw32-gcc \
  go build -o silentcast.exe ./cmd/silentcast

# macOS from Linux (requires osxcross)
GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 \
  CC=o64-clang \
  go build -o silentcast-mac ./cmd/silentcast

# Linux from Windows/macOS
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
  go build -o silentcast-linux ./cmd/silentcast
```

## Build Tags

Control feature inclusion with build tags:

| Tag | Description | Effect |
|-----|-------------|--------|
| `nogohook` | Disable hotkey support | No CGO required |
| `notray` | Disable system tray | Smaller binary |
| `noautoupdate` | Disable auto-updater | No update checks |

Example:
```bash
go build -tags "nogohook notray" ./cmd/silentcast
```

## Build Configuration

### Version Information

Set version during build:

```bash
VERSION=1.2.3
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD)

go build -ldflags "\
  -X main.Version=$VERSION \
  -X main.BuildTime=$BUILD_TIME \
  -X main.GitCommit=$GIT_COMMIT" \
  ./cmd/silentcast
```

### Optimization Flags

```bash
# Reduce binary size
go build -ldflags="-s -w" ./cmd/silentcast

# Enable optimizations
go build -gcflags="-m -l" ./cmd/silentcast
```

## Build Artifacts

### Directory Structure

```
app/
├── build/              # Build output
│   ├── silentcast     # Binary
│   └── *.log          # Build logs
├── dist/              # Distribution packages
│   ├── windows/       # Windows installer
│   ├── macos/         # macOS .app bundle
│   └── linux/         # Linux packages
└── release/           # GitHub release artifacts
```

### Creating Distribution Packages

#### Windows Installer

```bash
make build-windows-installer
# Creates: dist/windows/SilentCast-Setup.exe
```

#### macOS App Bundle

```bash
make build-macos-app
# Creates: dist/macos/SilentCast.app
```

#### Linux Packages

```bash
# Debian/Ubuntu
make build-deb
# Creates: dist/linux/silentcast_1.0.0_amd64.deb

# Red Hat/Fedora
make build-rpm
# Creates: dist/linux/silentcast-1.0.0.x86_64.rpm
```

## CI/CD Building

### GitHub Actions

```yaml
# .github/workflows/build.yml
name: Build

on: [push, pull_request]

jobs:
  build:
    strategy:
      matrix:
        os: [ubuntu-latest, windows-latest, macos-latest]
        go: ['1.23']
    
    runs-on: ${{ matrix.os }}
    
    steps:
    - uses: actions/checkout@v4
    
    - uses: actions/setup-go@v5
      with:
        go-version: ${{ matrix.go }}
    
    - name: Build
      run: make build
    
    - name: Upload artifact
      uses: actions/upload-artifact@v4
      with:
        name: silentcast-${{ matrix.os }}
        path: app/build/silentcast*
```

### Local CI Testing

```bash
# Test CI build locally with act
act -j build
```

## Troubleshooting Builds

### Common Issues

#### CGO Errors

```bash
# Error: C compiler not found
# Solution: Install build tools
# Windows: Install MinGW-w64
# macOS: xcode-select --install
# Linux: sudo apt-get install build-essential
```

#### Missing Dependencies

```bash
# Update dependencies
go mod download

# Verify dependencies
go mod verify

# Clean module cache
go clean -modcache
```

#### Cross-Compilation Failures

```bash
# Use stub build for cross-platform
GOOS=windows go build -tags "nogohook notray" ./cmd/silentcast
```

### Build Debugging

```bash
# Verbose output
go build -v -x ./cmd/silentcast

# Show build decisions
go build -gcflags="-m" ./cmd/silentcast

# Check build constraints
go list -f '{{.GoFiles}}' ./cmd/silentcast
```

## Performance Profiling

### CPU Profiling

```bash
# Build with profiling
go build -o silentcast.prof ./cmd/silentcast

# Run with profiling
./silentcast.prof -cpuprofile=cpu.prof

# Analyze
go tool pprof cpu.prof
```

### Memory Profiling

```bash
# Build with memory profiling
go build -gcflags="-m" ./cmd/silentcast

# Run with memory profiling
./silentcast -memprofile=mem.prof

# Analyze
go tool pprof mem.prof
```

## Security Considerations

### Code Signing

#### Windows
```powershell
# Sign with certificate
signtool sign /a /t http://timestamp.digicert.com silentcast.exe
```

#### macOS
```bash
# Sign for distribution
codesign --force --deep --sign "Developer ID Application: Your Name" SilentCast.app

# Notarize
xcrun altool --notarize-app --file SilentCast.zip
```

### Reproducible Builds

```bash
# Set consistent build environment
export CGO_ENABLED=1
export GOFLAGS="-trimpath"
export SOURCE_DATE_EPOCH=$(date +%s)

# Build
go build -buildmode=pie \
  -ldflags="-s -w -buildid=" \
  ./cmd/silentcast
```

## Best Practices

1. **Always test builds** on target platforms
2. **Use CI/CD** for consistent builds
3. **Sign binaries** for distribution
4. **Include version info** in builds
5. **Document dependencies** clearly
6. **Test cross-platform** builds regularly

## Next Steps

- [Architecture](/api/architecture) - Understanding the codebase
- [Testing](/api/testing) - Testing your builds
- [Contributing](/api/contributing) - Contributing guidelines