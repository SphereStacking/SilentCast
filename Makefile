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
	@echo "  build         Build the application"
	@echo "  build-stub    Build with stub mode (for development)"
	@echo "  test          Run tests"
	@echo "  clean         Remove build artifacts"
	@echo ""
	@echo "Documentation:"
	@echo "  docs-dev      Start documentation development server"
	@echo "  docs-build    Build documentation"
	@echo ""
	@echo "Project Management:"
	@echo "  setup         Setup development environment"

# Application build
build:
	@$(MAKE) -C app build

build-stub:
	@$(MAKE) -C app build-stub

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