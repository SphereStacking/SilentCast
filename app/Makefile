.PHONY: build build-dev build-stub build-snapshot test clean run install lint fmt release-test snapshot

# Build variables
BINARY_NAME=silentcast
MAIN_PATH=./cmd/silentcast
BUILD_DIR=./build
DIST_DIR=./dist
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "0.1.0-dev")
LDFLAGS=-ldflags "-X main.Version=${VERSION} -s -w"
BUILD_TAGS ?=

# Default target
all: build

# Build the application
build:
	@echo "🔨 Building silentcast..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=1 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Development build (fast, no C dependencies)
build-dev:
	@echo "🚀 Building development version (fast)..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 go build $(LDFLAGS) -tags "nogohook notray" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ Dev build complete: $(BUILD_DIR)/$(BINARY_NAME)"
	@echo "💡 Note: Hotkeys disabled in dev build"

# Alias for backward compatibility
build-stub: build-dev

# Run the application
run:
	@go run $(MAIN_PATH)

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test -tags "nogohook notray" -v -race -cover ./...

# Run tests with coverage report
test-coverage:
	@echo "📊 Generating coverage report..."
	@go test -tags "nogohook notray" -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: coverage.html"

# Install the application
install:
	@echo "📦 Installing silentcast..."
	@go install $(LDFLAGS) $(MAIN_PATH)
	@echo "✅ Installed to $(GOPATH)/bin/$(BINARY_NAME)"

# Format code
fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted"

# Lint code
lint:
	@echo "🔍 Linting code..."
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "golangci-lint not found. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@golangci-lint run ./...

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	@rm -rf $(BUILD_DIR) $(DIST_DIR)
	@rm -f coverage.out coverage.html
	@echo "✅ Clean complete"

# Snapshot build for all platforms (test release)
build-snapshot:
	@echo "📦 Building snapshot for all platforms..."
	@if command -v goreleaser > /dev/null 2>&1; then \
		cd .. && goreleaser release --snapshot --clean; \
	elif [ -f ~/go/bin/goreleaser ]; then \
		cd .. && ~/go/bin/goreleaser release --snapshot --clean; \
	else \
		echo "❌ goreleaser not found. Install with: go install github.com/goreleaser/goreleaser/v2@latest"; \
		exit 1; \
	fi
	@echo "✅ Snapshot build complete! Check app/dist/"

# Aliases for backward compatibility
snapshot: build-snapshot
release-test: build-snapshot

# Download dependencies
deps:
	@echo "📥 Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies ready"

# Build for multiple platforms
build-all: build-darwin build-windows

build-darwin:
	@echo "🍎 Building for macOS..."
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)

build-windows:
	@echo "🪟 Building for Windows..."
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)

# Development mode
dev: build-stub
	@echo "🚀 Running in development mode..."
	@$(BUILD_DIR)/$(BINARY_NAME) --no-tray

# Build release for all platforms
release:
	@echo "📦 Building release $(VERSION)..."
	@VERSION=$(VERSION) $(BUILD_DIR)/build.sh

# Help
help:
	@echo "SilentCast Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build          Build the application"
	@echo "  run            Run the application"
	@echo "  test           Run tests"
	@echo "  test-coverage  Run tests with coverage report"
	@echo "  install        Install the application"
	@echo "  fmt            Format code"
	@echo "  lint           Lint code"
	@echo "  clean          Clean build artifacts"
	@echo "  deps           Download dependencies"
	@echo "  build-all      Build for Windows and macOS"
	@echo "  dev            Run in development mode with hot reload"
	@echo "  help           Show this help message"
