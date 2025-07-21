.PHONY: build build-dev build-stub build-snapshot test test-integration test-e2e test-all test-coverage test-integration-coverage test-e2e-coverage clean run install lint fmt release-test snapshot

# Build variables
BINARY_NAME=silentcast
MAIN_PATH=./cmd/silentcast
BUILD_DIR=./build
DIST_DIR=./dist
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "0.1.0-dev")
GIT_COMMIT=$(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(shell date -u '+%Y-%m-%d %H:%M:%S UTC')
LDFLAGS=-ldflags "\
	-X 'github.com/SphereStacking/silentcast/internal/version.Version=${VERSION}' \
	-X 'github.com/SphereStacking/silentcast/internal/version.GitCommit=${GIT_COMMIT}' \
	-X 'github.com/SphereStacking/silentcast/internal/version.BuildTime=${BUILD_TIME}' \
	-X 'github.com/SphereStacking/silentcast/internal/version.BuildTags=${BUILD_TAGS}' \
	-s -w"
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
	@if [ "$$OS" = "Windows_NT" ] || [ "$$(go env GOOS)" = "windows" ]; then \
		CGO_ENABLED=0 go build -ldflags "\
			-X 'github.com/SphereStacking/silentcast/internal/version.Version=${VERSION}' \
			-X 'github.com/SphereStacking/silentcast/internal/version.GitCommit=${GIT_COMMIT}' \
			-X 'github.com/SphereStacking/silentcast/internal/version.BuildTime=${BUILD_TIME}' \
			-X 'github.com/SphereStacking/silentcast/internal/version.BuildTags=nogohook,notray' \
			-s -w" -tags "nogohook notray" -o $(BUILD_DIR)/$(BINARY_NAME).exe $(MAIN_PATH); \
		echo "✅ Dev build complete: $(BUILD_DIR)/$(BINARY_NAME).exe"; \
	else \
		CGO_ENABLED=0 go build -ldflags "\
			-X 'github.com/SphereStacking/silentcast/internal/version.Version=${VERSION}' \
			-X 'github.com/SphereStacking/silentcast/internal/version.GitCommit=${GIT_COMMIT}' \
			-X 'github.com/SphereStacking/silentcast/internal/version.BuildTime=${BUILD_TIME}' \
			-X 'github.com/SphereStacking/silentcast/internal/version.BuildTags=nogohook,notray' \
			-s -w" -tags "nogohook notray" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH); \
		echo "✅ Dev build complete: $(BUILD_DIR)/$(BINARY_NAME)"; \
	fi
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

# TDD workflow commands
tdd-watch:
	@echo "👀 Running TDD watch mode..."
	@echo "Watching for changes and running tests automatically..."
	@if command -v entr > /dev/null 2>&1; then \
		find . -name "*.go" | entr -c go test -tags "nogohook notray" -v -race -cover ./...; \
	elif command -v fswatch > /dev/null 2>&1; then \
		fswatch -o . | xargs -n1 -I{} make test; \
	else \
		echo "❌ Need 'entr' or 'fswatch' for watch mode. Install with:"; \
		echo "  - Ubuntu/Debian: apt install entr"; \
		echo "  - macOS: brew install entr"; \
		echo "  - Or install fswatch"; \
		exit 1; \
	fi

tdd-red:
	@echo "🔴 TDD RED: Running tests (expecting failures)..."
	@go test -tags "nogohook notray" -v -failfast ./... || true
	@echo "RED phase complete. Write minimal code to make tests pass."

tdd-green:
	@echo "🟢 TDD GREEN: Running tests (should pass now)..."
	@go test -tags "nogohook notray" -v -race -cover ./...
	@echo "GREEN phase complete. Time to refactor!"

tdd-refactor:
	@echo "🔵 TDD REFACTOR: Running tests while refactoring..."
	@go test -tags "nogohook notray" -v -race -cover ./...
	@echo "REFACTOR phase complete. Tests should still be green!"

tdd-cycle:
	@echo "♻️ TDD CYCLE: Red -> Green -> Refactor"
	@echo "1. 🔴 RED: Write a failing test"
	@make tdd-red
	@echo ""
	@echo "2. 🟢 GREEN: Make the test pass"
	@read -p "Press Enter when ready to run GREEN phase..." dummy
	@make tdd-green  
	@echo ""
	@echo "3. 🔵 REFACTOR: Improve code while keeping tests green"
	@read -p "Press Enter when ready to run REFACTOR phase..." dummy
	@make tdd-refactor

tdd-metrics-start:
	@echo "📊 Starting TDD metrics collection..."
	@../scripts/tdd-metrics.sh start

tdd-metrics-red:
	@echo "📊 Recording RED phase metrics..."
	@../scripts/tdd-metrics.sh red

tdd-metrics-green:
	@echo "📊 Recording GREEN phase metrics..."
	@../scripts/tdd-metrics.sh green

tdd-metrics-refactor:
	@echo "📊 Recording REFACTOR phase metrics..."
	@../scripts/tdd-metrics.sh refactor

tdd-metrics-complete:
	@echo "📊 Completing TDD cycle metrics..."
	@../scripts/tdd-metrics.sh complete

tdd-metrics-report:
	@echo "📈 Generating TDD metrics report..."
	@../scripts/tdd-metrics.sh report

# Run tests with coverage report
test-coverage:
	@echo "📊 Generating coverage report..."
	@go test -tags "nogohook notray" -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: coverage.html"

# Run integration tests
test-integration:
	@echo "🔗 Running integration tests..."
	@go test -tags "integration nogohook notray" -v -timeout=10m ./test/integration/...

# Run end-to-end tests
test-e2e: build-stub
	@echo "🌐 Running E2E tests..."
	@echo "Building application for E2E testing..."
	@if [ ! -f $(BUILD_DIR)/$(BINARY_NAME) ]; then \
		echo "Building application binary..."; \
		$(MAKE) build-stub; \
	fi
	@go test -tags "e2e" -v -timeout=30m ./test/e2e/...

# Run E2E tests with coverage
test-e2e-coverage: build-stub
	@echo "🌐 Running E2E tests with coverage..."
	@if [ ! -f $(BUILD_DIR)/$(BINARY_NAME) ]; then \
		echo "Building application binary..."; \
		$(MAKE) build-stub; \
	fi
	@go test -tags "e2e" -v -timeout=30m -cover ./test/e2e/...

# Run all tests (unit + integration + e2e)
test-all:
	@echo "🧪 Running all tests..."
	@$(MAKE) test
	@$(MAKE) test-integration 
	@$(MAKE) test-e2e
	@echo "✅ All tests completed"

# Run integration tests with coverage
test-integration-coverage:
	@echo "📊 Running integration tests with coverage..."
	@go test -tags "integration nogohook notray" -v -timeout=10m -coverprofile=integration-coverage.out ./test/integration/...
	@go tool cover -html=integration-coverage.out -o integration-coverage.html
	@echo "✅ Integration coverage report: integration-coverage.html"

# Benchmark targets removed - use standard Go benchmark tools if needed

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
	@PATH="$(PATH):$$(go env GOPATH)/bin" golangci-lint run ./...

# Auto-fix linting issues
lint-fix:
	@echo "🔧 Auto-fixing lint issues..."
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "golangci-lint not found. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@PATH="$(PATH):$$(go env GOPATH)/bin" golangci-lint run --fix ./...
	@echo "✅ Auto-fixable issues resolved"

# Format and auto-fix all issues
fix-all: fmt lint-fix
	@echo "🎉 All auto-fixable issues resolved"

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
# Profiling targets removed - use standard Go profiling tools if needed

help:
	@echo "SilentCast Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build                    Build the application"
	@echo "  run                      Run the application"
	@echo "  test                     Run unit tests"
	@echo "  test-coverage            Run unit tests with coverage report"
	@echo "  test-integration         Run integration tests"
	@echo "  test-integration-coverage Run integration tests with coverage"
	@echo "  test-e2e                 Run end-to-end tests"
	@echo "  test-e2e-coverage        Run E2E tests with coverage"
	@echo "  test-all                 Run all tests (unit + integration + e2e)"
	@echo ""
	@echo "Other Targets:"
	@echo "  install                  Install the application"
	@echo "  fmt                      Format code"
	@echo "  lint                     Lint code"
	@echo "  lint-fix                 Auto-fix lint issues"
	@echo "  fix-all                  Format and auto-fix all issues"
	@echo "  clean                    Clean build artifacts"
	@echo "  deps                     Download dependencies"
	@echo "  build-all                Build for Windows and macOS"
	@echo "  dev                      Run in development mode"
	@echo "  help                     Show this help message"
