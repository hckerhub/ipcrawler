#!/usr/bin/env make

# IPCrawler Build Configuration
BINARY_NAME=ipcrawler
MAIN_PACKAGE_PATH=./
VERSION?=$(shell git describe --tags --dirty --always)
GIT_COMMIT?=$(shell git rev-parse HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-s -w -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME)

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
deps: ## Download and verify dependencies
	@echo -e "$(YELLOW)Downloading dependencies...$(NC)"
	go mod download
	go mod verify
	go mod tidy
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
install: build ## Install binary to local system
	@echo -e "$(YELLOW)Installing $(BINARY_NAME) to system...$(NC)"
	@if [ -w /usr/local/bin ]; then \
		cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/; \
		echo -e "$(GREEN)✓ Installed to /usr/local/bin/$(BINARY_NAME)$(NC)"; \
	else \
		echo "Installing to /usr/local/bin requires sudo:"; \
		sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/; \
		echo -e "$(GREEN)✓ Installed to /usr/local/bin/$(BINARY_NAME)$(NC)"; \
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