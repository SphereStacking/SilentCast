# SilentCast - Root Makefile

.PHONY: help build test docs clean

# Default target
help:
	@echo "SilentCast Project"
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Application Build:"
	@echo "  build-dev     Build for development (fast, no hotkeys)"
	@echo "  build         Build production version (requires C libs)"
	@echo "  build-snapshot Test release build for all platforms"
	@echo "  test          Run tests"
	@echo "  clean         Remove build artifacts"
	@echo ""
	@echo "Documentation:"
	@echo "  docs-dev      Start documentation development server"
	@echo "  docs-build    Build documentation"
	@echo ""
	@echo "Project Management:"
	@echo "  setup         Setup development environment"
	@echo "  lint          Run linting checks"

# Application build
# Production build (requires C libraries)
build:
	@$(MAKE) -C app build

# Development build (fast, no dependencies)
build-dev:
	@$(MAKE) -C app build-dev

# Snapshot build for all platforms
build-snapshot:
	@$(MAKE) -C app build-snapshot

test:
	@$(MAKE) -C app test

clean:
	@$(MAKE) -C app clean
	@rm -rf docs/node_modules docs/.vitepress/dist

# Documentation
docs-dev:
	@echo "📚 Starting documentation server..."
	@cd docs && npm install && npm run docs:dev

docs-build:
	@echo "📦 Building documentation..."
	@cd docs && npm install && npm run docs:build

# Development environment setup
setup:
	@echo "🔧 Setting up development environment..."
	@cd app && go mod download
	@cd docs && npm install
	@echo "✅ Setup complete!"

# Run linting
lint:
	@echo "🔍 Running linters..."
	@cd app && golangci-lint run ./...
