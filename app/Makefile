.PHONY: build test clean run install lint fmt

# Build variables
BINARY_NAME=silentcast
MAIN_PATH=./cmd/silentcast
BUILD_DIR=./build
DIST_DIR=./dist
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "0.1.0-dev")
LDFLAGS=-ldflags "-X main.version=${VERSION} -s -w"
BUILD_TAGS ?=

# Default target
all: build

# Build the application
build:
	@echo "🔨 Building silentcast..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 go build $(LDFLAGS) -tags "$(BUILD_TAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build with stub implementations (for testing)
build-stub:
	@echo "🔨 Building silentcast with stubs..."
	@$(MAKE) build BUILD_TAGS="nogohook notray"

# Run the application
run:
	@go run $(MAIN_PATH)

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test -v -race -cover -tags "nogohook notray" ./...

# Run tests with coverage report
test-coverage:
	@echo "📊 Generating coverage report..."
	@go test -v -race -coverprofile=coverage.out -tags "nogohook notray" ./...
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