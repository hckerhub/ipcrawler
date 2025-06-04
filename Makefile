# IPCrawler Makefile v0.1
# Build system for the advanced IP scanner
# Author: hckerhub (https://github.com/hckerhub)

.PHONY: all build clean test run help deps crawler start tui uninstall

# Variables
BINARY_NAME=ipcrawler
BINARY_PATH=./$(BINARY_NAME)
VERSION=0.1
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Default target
all: crawler

# Build the application
build:
	@echo "🔨 Building IPCrawler..."
	@go build $(LDFLAGS) -o $(BINARY_NAME) .
	@echo "✅ Build complete: $(BINARY_NAME)"

# Download dependencies
deps:
	@echo "📦 Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies updated"

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test -v ./...
	@echo "✅ Tests completed"

# Run with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out coverage.html
	@echo "✅ Clean complete"

# Complete system setup - build, install, and make ready to use
crawler: deps build
	@echo "🚀 Setting up IPCrawler for system-wide use..."
	@sudo cp $(BINARY_NAME) /usr/local/bin/ 2>/dev/null || { \
		echo "⚠️  Could not install to system PATH (permission denied)"; \
		echo "📋 You can still use './crawler' or './$(BINARY_NAME)' directly"; \
	}
	@sudo chmod +x /usr/local/bin/$(BINARY_NAME) 2>/dev/null || true
	@echo "✅ IPCrawler ready! Use: 'ipcrawler', './crawler', or 'make start'"

# Uninstall from system PATH
uninstall:
	@echo "🗑️  Uninstalling IPCrawler..."
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "✅ IPCrawler uninstalled"

# Run the TUI directly
run-tui: build
	@echo "🚀 Starting IPCrawler TUI..."
	@$(BINARY_PATH) tui

# Quick start aliases
start: run-tui
tui: run-tui

# Run a test scan
run-scan: build
	@echo "🚀 Running test scan on scanme.nmap.org..."
	@$(BINARY_PATH) scan scanme.nmap.org --verbose

# Format code
fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted"

# Lint code
lint:
	@echo "🔍 Linting code..."
	@golangci-lint run
	@echo "✅ Linting complete"

# Build for multiple platforms
build-all: clean
	@echo "🔨 Building for multiple platforms..."
	@mkdir -p dist
	
	# Linux 64-bit
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 .
	
	# macOS 64-bit (Intel)
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 .
	
	# macOS ARM64 (Apple Silicon)
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 .
	
	# Windows 64-bit
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe .
	
	@echo "✅ Multi-platform build complete:"
	@ls -la dist/

# Create release packages
package: build-all
	@echo "📦 Creating release packages..."
	@cd dist && \
		tar -czf $(BINARY_NAME)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64 && \
		tar -czf $(BINARY_NAME)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64 && \
		tar -czf $(BINARY_NAME)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64 && \
		zip $(BINARY_NAME)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe
	@echo "✅ Release packages created in dist/"

# Development mode (live reload)
dev:
	@echo "🔧 Starting development mode with live reload..."
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	@air

# Check prerequisites
check-deps:
	@echo "🔍 Checking prerequisites..."
	@command -v go >/dev/null 2>&1 || { echo "❌ Go is not installed"; exit 1; }
	@command -v nmap >/dev/null 2>&1 || { echo "❌ NMAP is not installed"; exit 1; }
	@echo "✅ All prerequisites are installed"

# Security scan
security:
	@echo "🔒 Running security scan..."
	@which gosec > /dev/null || (echo "Installing gosec..." && go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest)
	@gosec ./...
	@echo "✅ Security scan complete"

# Generate documentation
docs:
	@echo "📚 Generating documentation..."
	@which godoc > /dev/null || (echo "Installing godoc..." && go install golang.org/x/tools/cmd/godoc@latest)
	@godoc -http=:6060 &
	@echo "✅ Documentation server started at http://localhost:6060"

# Docker build
docker-build:
	@echo "🐳 Building Docker image..."
	@docker build -t ipcrawler:latest .
	@echo "✅ Docker image built: ipcrawler:latest"

# Docker run
docker-run: docker-build
	@echo "🐳 Running IPCrawler in Docker..."
	@docker run -it --rm ipcrawler:latest

# Show version info
version:
	@echo "IPCrawler version: $(VERSION)"
	@echo "Build time: $(BUILD_TIME)"

# Help target
help:
	@echo "🎯 IPCrawler Build System v$(VERSION)"
	@echo ""
	@echo "🚀 MAIN COMMANDS:"
	@echo "  crawler       Complete setup - build, install, and make ready (RECOMMENDED)"
	@echo "  start         Quick start TUI mode"
	@echo "  clean         Clean all build artifacts"
	@echo ""
	@echo "🔧 DEVELOPMENT:"
	@echo "  build         Build the application only"
	@echo "  test          Run tests"
	@echo "  test-coverage Run tests with coverage report"
	@echo "  fmt           Format source code"
	@echo "  lint          Lint source code"
	@echo "  dev           Start development mode with live reload"
	@echo ""
	@echo "📦 DISTRIBUTION:"
	@echo "  build-all     Build for multiple platforms"
	@echo "  package       Create release packages"
	@echo "  uninstall     Remove from system PATH"
	@echo ""
	@echo "🛠️  UTILITIES:"
	@echo "  check-deps    Check prerequisites"
	@echo "  security      Run security scan"
	@echo "  docs          Generate and serve documentation"
	@echo "  version       Show version information"
	@echo ""
	@echo "Examples:"
	@echo "  make crawler        # Complete setup (build + install)"
	@echo "  make start          # Quick start TUI"
	@echo "  ./crawler           # Use quick launcher script"
	@echo "  ipcrawler --help    # After setup, use system command" 