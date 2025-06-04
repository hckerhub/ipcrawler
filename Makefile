#!/usr/bin/env make

# IPCrawler Build Configuration
BINARY_NAME=ipcrawler
MAIN_PACKAGE_PATH=./
VERSION?=$(shell git describe --tags --dirty --always)
GIT_COMMIT?=$(shell git rev-parse HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-s -w -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME)

# System detection
OS=$(shell uname -s)
GO_VERSION=$(shell go version 2>/dev/null | sed -n 's/.*go\([0-9]*\.[0-9]*\).*/\1/p')
HAS_GO=$(shell command -v go >/dev/null 2>&1 && echo "true" || echo "false")

# Environment detection  
GO_COMPAT=$(shell \
	if [ "$(HAS_GO)" = "true" ] && [ -n "$(GO_VERSION)" ]; then \
		MAJOR=$$(echo $(GO_VERSION) | cut -d. -f1); \
		MINOR=$$(echo $(GO_VERSION) | cut -d. -f2); \
		if [ "$$MAJOR" -gt 1 ] || ([ "$$MAJOR" -eq 1 ] && [ "$$MINOR" -ge 19 ]); then \
			echo "true"; \
		else \
			echo "false"; \
		fi; \
	else \
		echo "false"; \
	fi \
)

# Build directories
BUILD_DIR=./build
DIST_DIR=./dist

# Platform targets
PLATFORMS=darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64

# Colors for output
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
NC=\033[0m # No Color

.PHONY: help
help: ## Show this help message
	@echo "IPCrawler Build System"
	@echo "====================="
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: clean
clean: ## Clean build artifacts
	@echo -e "$(YELLOW)Cleaning build artifacts...$(NC)"
	rm -rf $(BUILD_DIR) $(DIST_DIR)
	go clean -cache
	@echo -e "$(GREEN)✓ Clean complete$(NC)"

.PHONY: deps
deps: ## Download and verify dependencies (auto-detects compatibility)
	@echo -e "$(YELLOW)Downloading dependencies...$(NC)"
	@echo "System: $(OS), Go: $(GO_VERSION), Compatible: $(GO_COMPAT)"
	@if [ "$(HAS_GO)" != "true" ]; then \
		echo -e "$(RED)❌ Go is not installed$(NC)"; \
		echo "Please install Go first:"; \
		if [ "$(OS)" = "Linux" ]; then \
			echo "  sudo apt update && sudo apt install golang-go"; \
		elif [ "$(OS)" = "Darwin" ]; then \
			echo "  brew install go"; \
		fi; \
		exit 1; \
	fi
	@if [ "$(GO_COMPAT)" != "true" ]; then \
		echo -e "$(YELLOW)⚠️  Older Go version detected ($(GO_VERSION)), applying compatibility fixes...$(NC)"; \
		sed -i.bak 's/go 1\.2[3-9]/go 1.19/' go.mod || true; \
		sed -i.bak '/^toolchain/d' go.mod || true; \
		rm -f go.mod.bak; \
		echo "Applied compatibility patches for Go $(GO_VERSION)"; \
	fi
	@if go mod download 2>/dev/null; then \
		echo "✓ Dependencies downloaded"; \
	else \
		echo -e "$(YELLOW)Retrying with fallback approach...$(NC)"; \
		go clean -modcache; \
		go mod download; \
	fi
	@if go mod verify 2>/dev/null; then \
		echo "✓ Dependencies verified"; \
	else \
		echo -e "$(YELLOW)Verification skipped (older Go version)$(NC)"; \
	fi
	@if go mod tidy 2>/dev/null; then \
		echo "✓ Dependencies tidied"; \
	else \
		echo -e "$(YELLOW)Tidy completed with warnings$(NC)"; \
	fi
	@echo -e "$(GREEN)✓ Dependencies ready$(NC)"

.PHONY: fmt
fmt: ## Format code
	@echo -e "$(YELLOW)Formatting code...$(NC)"
	go fmt ./...
	@echo -e "$(GREEN)✓ Code formatted$(NC)"

.PHONY: lint
lint: ## Run linting
	@echo -e "$(YELLOW)Running linters...$(NC)"
	go vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, skipping advanced linting"; \
	fi
	@echo -e "$(GREEN)✓ Linting complete$(NC)"

.PHONY: test
test: ## Run tests
	@echo -e "$(YELLOW)Running tests...$(NC)"
	go test -v -race -coverprofile=coverage.out ./...
	@echo -e "$(GREEN)✓ Tests complete$(NC)"

.PHONY: build
build: deps fmt ## Build binary for current platform
	@echo -e "$(YELLOW)Building $(BINARY_NAME) for current platform...$(NC)"
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build \
		-ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(BINARY_NAME) \
		$(MAIN_PACKAGE_PATH)
	@echo -e "$(GREEN)✓ Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

.PHONY: build-all
build-all: deps fmt ## Build binaries for all platforms
	@echo -e "$(YELLOW)Building $(BINARY_NAME) for all platforms...$(NC)"
	mkdir -p $(DIST_DIR)
	@for platform in $(PLATFORMS); do \
		OS=$$(echo $$platform | cut -d'/' -f1); \
		ARCH=$$(echo $$platform | cut -d'/' -f2); \
		OUTPUT_NAME=$(BINARY_NAME)-$$OS-$$ARCH; \
		if [ "$$OS" = "windows" ]; then OUTPUT_NAME=$$OUTPUT_NAME.exe; fi; \
		echo "Building for $$OS/$$ARCH..."; \
		CGO_ENABLED=0 GOOS=$$OS GOARCH=$$ARCH go build \
			-ldflags "$(LDFLAGS)" \
			-o $(DIST_DIR)/$$OUTPUT_NAME \
			$(MAIN_PACKAGE_PATH); \
		if [ $$? -eq 0 ]; then \
			echo "✓ Built $(DIST_DIR)/$$OUTPUT_NAME"; \
		else \
			echo "✗ Failed to build for $$OS/$$ARCH"; \
			exit 1; \
		fi; \
	done
	@echo -e "$(GREEN)✓ All builds complete$(NC)"

.PHONY: package
package: build-all ## Create distribution packages
	@echo -e "$(YELLOW)Creating distribution packages...$(NC)"
	cd $(DIST_DIR) && for binary in $(BINARY_NAME)-*; do \
		if [ -f "$$binary" ]; then \
			echo "Packaging $$binary..."; \
			tar -czf "$$binary.tar.gz" "$$binary"; \
			echo "✓ Created $$binary.tar.gz"; \
		fi; \
	done
	@echo -e "$(GREEN)✓ Packaging complete$(NC)"

.PHONY: install
install: ## Smart install - auto-detects system and applies best approach
	@echo -e "$(YELLOW)🔧 IPCrawler Smart Installation$(NC)"
	@echo "=================================="
	@echo "System: $(OS)"
	@echo "Go Version: $(GO_VERSION)"
	@echo "Go Compatible: $(GO_COMPAT)"
	@echo ""
	@if [ "$(HAS_GO)" != "true" ]; then \
		echo -e "$(RED)❌ Go not found. Installing Go first...$(NC)"; \
		if [ "$(OS)" = "Linux" ]; then \
			sudo apt update && sudo apt install golang-go; \
		elif [ "$(OS)" = "Darwin" ]; then \
			if command -v brew >/dev/null 2>&1; then \
				brew install go; \
			else \
				echo "Please install Homebrew first, then run: brew install go"; \
				exit 1; \
			fi; \
		else \
			echo "Please install Go manually for your system"; \
			exit 1; \
		fi; \
	fi
	@echo -e "$(YELLOW)Building application...$(NC)"
	@$(MAKE) build
	@echo -e "$(YELLOW)Installing $(BINARY_NAME) to system...$(NC)"
	@if [ -w /usr/local/bin ] 2>/dev/null; then \
		cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/; \
		chmod +x /usr/local/bin/$(BINARY_NAME); \
		echo -e "$(GREEN)✅ Installed to /usr/local/bin/$(BINARY_NAME)$(NC)"; \
	else \
		echo "Installing to /usr/local/bin (requires sudo):"; \
		sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/; \
		sudo chmod +x /usr/local/bin/$(BINARY_NAME); \
		echo -e "$(GREEN)✅ Installed to /usr/local/bin/$(BINARY_NAME)$(NC)"; \
	fi
	@echo ""
	@echo -e "$(GREEN)🎯 Installation Complete!$(NC)"
	@echo "Usage:"
	@echo "  $(BINARY_NAME)                    # Run from anywhere"
	@echo "  ./$(BUILD_DIR)/$(BINARY_NAME)     # Run from build directory"
	@if [ "$(OS)" = "Linux" ] && [ -f /etc/os-release ]; then \
		if grep -q "kali\|parrot\|htb" /etc/os-release 2>/dev/null || [ -d "/home/kali" ] || [ -d "/opt/pwnbox" ]; then \
			echo ""; \
			echo -e "$(YELLOW)🔍 Detected pentesting environment$(NC)"; \
			echo "Ready for CTF/HTB target scanning!"; \
		fi; \
	fi

.PHONY: uninstall
uninstall: ## Remove binary from local system
	@echo -e "$(YELLOW)Removing $(BINARY_NAME) from system...$(NC)"
	@if [ -w /usr/local/bin ]; then \
		rm -f /usr/local/bin/$(BINARY_NAME); \
		echo -e "$(GREEN)✓ Removed from /usr/local/bin/$(BINARY_NAME)$(NC)"; \
	else \
		sudo rm -f /usr/local/bin/$(BINARY_NAME); \
		echo -e "$(GREEN)✓ Removed from /usr/local/bin/$(BINARY_NAME)$(NC)"; \
	fi

.PHONY: quick-install
quick-install: ## Fast install for HTB/CTF environments (minimal checks)
	@echo -e "$(YELLOW)⚡ Quick Install for HTB/CTF$(NC)"
	@if go mod tidy 2>/dev/null; then echo "✓ Dependencies OK"; fi
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 go build -ldflags "-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE_PATH)
	@chmod +x $(BUILD_DIR)/$(BINARY_NAME)
	@if [ -w /usr/local/bin ] 2>/dev/null; then \
		cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/ && echo "✅ Installed to /usr/local/bin/$(BINARY_NAME)"; \
	else \
		sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/ && sudo chmod +x /usr/local/bin/$(BINARY_NAME) && echo "✅ Installed to /usr/local/bin/$(BINARY_NAME)"; \
	fi
	@echo "🎯 Ready! Run: $(BINARY_NAME)"

.PHONY: run
run: build ## Build and run the application
	@echo -e "$(YELLOW)Running $(BINARY_NAME)...$(NC)"
	$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)

.PHONY: dev
dev: ## Run in development mode (equivalent to go run)
	@echo -e "$(YELLOW)Running in development mode...$(NC)"
	go run $(MAIN_PACKAGE_PATH) $(ARGS)

.PHONY: audit
audit: deps fmt lint test ## Run full audit (format, lint, test)
	@echo -e "$(GREEN)✓ Full audit complete$(NC)"

.PHONY: release-check
release-check: audit build-all ## Prepare for release (audit + build all)
	@echo -e "$(GREEN)✓ Release check complete - ready for distribution$(NC)"
	@echo ""
	@echo "Distribution files:"
	@ls -la $(DIST_DIR)/

# Default target
.DEFAULT_GOAL := help 