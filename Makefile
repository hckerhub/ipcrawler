# ipcrawler - Advanced Reconnaissance Toolchain Makefile
# Author: AI Assistant
# Description: Automated setup for cybersecurity reconnaissance tools

# ===== CONFIGURATION =====
PROJECT_NAME := ipcrawler
TOOLS_DIR := $(CURDIR)/tools
WORDLISTS_DIR := $(CURDIR)/wordlists
BIN_DIR := $(HOME)/.local/bin
LOCAL_BIN_DIR := $(CURDIR)/bin
GO_VERSION := 1.21.5
GO_INSTALL_DIR := $(HOME)/.local/go

# Colors for output
RED := \033[31m
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
MAGENTA := \033[35m
CYAN := \033[36m
RESET := \033[0m

# ===== HELPER FUNCTIONS =====
define print_status
	@echo "$(CYAN)➡️  $(1)$(RESET)"
endef

define print_success
	@echo "$(GREEN)✅ $(1)$(RESET)"
endef

define print_warning
	@echo "$(YELLOW)⚠️  $(1)$(RESET)"
endef

define print_error
	@echo "$(RED)❌ $(1)$(RESET)"
endef

# ===== MAIN TARGETS =====
.PHONY: all install clean help check-deps setup-dirs install-deps install-tools install-wordlists

all: install

help:
	@echo "$(MAGENTA)🔍 ipcrawler - Advanced Reconnaissance Toolchain$(RESET)"
	@echo ""
	@echo "$(CYAN)Available commands:$(RESET)"
	@echo "  $(GREEN)make install$(RESET)     - Install all tools and dependencies"
	@echo "  $(GREEN)make clean$(RESET)       - Remove all installed tools and data"
	@echo "  $(GREEN)make check-deps$(RESET)  - Check system dependencies"
	@echo "  $(GREEN)make help$(RESET)        - Show this help message"
	@echo ""
	@echo "$(CYAN)Directory structure:$(RESET)"
	@echo "  $(CURDIR)/tools/      - Cloned tool repositories"
	@echo "  $(CURDIR)/bin/        - Locally built tool binaries"
	@echo "  $(CURDIR)/wordlists/  - Downloaded wordlists"
	@echo "  $(HOME)/.local/bin    - Global tool binaries (fallback)"

install: check-installation-status check-deps setup-dirs install-deps install-tools install-wordlists
	$(call print_success,"ipcrawler installation completed!")
	@echo ""
	@echo "$(CYAN)📋 Next steps:$(RESET)"
	@echo "  1. Add $(HOME)/.local/bin to your PATH if not already done:"
	@echo "     export PATH=\"$(HOME)/.local/bin:\$$PATH\""
	@echo "  2. Tools are ready to use!"
	@echo "  3. Run 'make help' for available commands"

# ===== DEPENDENCY CHECKING =====
check-installation-status:
	@echo "$(MAGENTA)🔍 Checking current installation status...$(RESET)"
	@echo ""
	@echo "$(CYAN)System Dependencies:$(RESET)"
	@command -v git >/dev/null 2>&1 && echo "  $(GREEN)✅ git$(RESET)" || echo "  $(RED)❌ git$(RESET) - Required"
	@command -v curl >/dev/null 2>&1 && echo "  $(GREEN)✅ curl$(RESET)" || echo "  $(RED)❌ curl$(RESET)"
	@command -v wget >/dev/null 2>&1 && echo "  $(GREEN)✅ wget$(RESET)" || echo "  $(RED)❌ wget$(RESET)"
	@command -v unzip >/dev/null 2>&1 && echo "  $(GREEN)✅ unzip$(RESET)" || echo "  $(RED)❌ unzip$(RESET)"
	@command -v python3 >/dev/null 2>&1 && echo "  $(GREEN)✅ python3$(RESET)" || echo "  $(RED)❌ python3$(RESET)"
	@command -v pip3 >/dev/null 2>&1 && echo "  $(GREEN)✅ pip3$(RESET)" || echo "  $(RED)❌ pip3$(RESET)"
	@echo ""
	@echo "$(CYAN)Development Dependencies:$(RESET)"
	@if command -v go >/dev/null 2>&1; then \
		GO_CURRENT=$$(go version | cut -d' ' -f3 | sed 's/go//'); \
		echo "  $(GREEN)✅ go $$GO_CURRENT$(RESET)"; \
	else \
		echo "  $(RED)❌ go$(RESET) - Will be installed"; \
	fi
	@command -v cargo >/dev/null 2>&1 && echo "  $(GREEN)✅ rust/cargo$(RESET)" || echo "  $(RED)❌ rust/cargo$(RESET) - Will be installed"
	@echo ""
	@echo "$(CYAN)Reconnaissance Tools:$(RESET)"
	@(command -v naabu >/dev/null 2>&1 || [ -f "$(LOCAL_BIN_DIR)/naabu" ]) && echo "  $(GREEN)✅ naabu$(RESET)" || echo "  $(YELLOW)⏳ naabu$(RESET) - Will be installed"
	@(command -v httpx >/dev/null 2>&1 || [ -f "$(LOCAL_BIN_DIR)/httpx" ]) && echo "  $(GREEN)✅ httpx$(RESET)" || echo "  $(YELLOW)⏳ httpx$(RESET) - Will be installed"
	@(command -v subfinder >/dev/null 2>&1 || [ -f "$(LOCAL_BIN_DIR)/subfinder" ]) && echo "  $(GREEN)✅ subfinder$(RESET)" || echo "  $(YELLOW)⏳ subfinder$(RESET) - Will be installed"
	@(command -v nuclei >/dev/null 2>&1 || [ -f "$(LOCAL_BIN_DIR)/nuclei" ]) && echo "  $(GREEN)✅ nuclei$(RESET)" || echo "  $(YELLOW)⏳ nuclei$(RESET) - Will be installed"
	@(command -v amass >/dev/null 2>&1 || [ -f "$(LOCAL_BIN_DIR)/amass" ]) && echo "  $(GREEN)✅ amass$(RESET)" || echo "  $(YELLOW)⏳ amass$(RESET) - Will be installed"
	@(command -v ffuf >/dev/null 2>&1 || [ -f "$(LOCAL_BIN_DIR)/ffuf" ]) && echo "  $(GREEN)✅ ffuf$(RESET)" || echo "  $(YELLOW)⏳ ffuf$(RESET) - Will be installed"
	@(command -v feroxbuster >/dev/null 2>&1 || [ -f "$(LOCAL_BIN_DIR)/feroxbuster" ]) && echo "  $(GREEN)✅ feroxbuster$(RESET)" || echo "  $(YELLOW)⏳ feroxbuster$(RESET) - Will be installed"
	@echo ""
	@echo "$(CYAN)Wordlists:$(RESET)"
	@[ -d "$(WORDLISTS_DIR)/SecLists" ] && echo "  $(GREEN)✅ SecLists$(RESET)" || echo "  $(YELLOW)⏳ SecLists$(RESET) - Will be downloaded"
	@[ -f "$(WORDLISTS_DIR)/rockyou.txt" ] && echo "  $(GREEN)✅ rockyou.txt$(RESET)" || echo "  $(YELLOW)⏳ rockyou.txt$(RESET) - Will be downloaded"
	@[ -d "$(WORDLISTS_DIR)/commonspeak2-wordlists" ] && echo "  $(GREEN)✅ commonspeak2$(RESET)" || echo "  $(YELLOW)⏳ commonspeak2$(RESET) - Will be downloaded"
	@echo ""

check-deps:
	$(call print_status,"Checking system dependencies...")
	@command -v git >/dev/null 2>&1 || ($(call print_error,"git is required but not installed") && exit 1)
	$(call print_success,"Dependencies check passed")

# ===== DIRECTORY SETUP =====
setup-dirs:
	$(call print_status,"Setting up directories...")
	@mkdir -p $(TOOLS_DIR)
	@mkdir -p $(WORDLISTS_DIR)
	@mkdir -p $(HOME)/.local/bin 2>/dev/null || mkdir -p $(LOCAL_BIN_DIR)
	$(call print_success,"Directories created")

# ===== DEPENDENCY INSTALLATION =====
install-deps: install-go install-system-deps install-rust
	$(call print_success,"All dependencies installed")

# Go installation from official source
install-go:
	$(call print_status,"Checking Go installation...")
	@if command -v go >/dev/null 2>&1; then \
		echo "$(GREEN)✅ Go is already installed$(RESET)"; \
	else \
		echo "$(CYAN)➡️  Go not found, installing Go $(GO_VERSION)...$(RESET)"; \
		$(MAKE) download-go; \
	fi

# Download and install Go from official source
download-go:
	$(call print_status,"Downloading Go $(GO_VERSION)...")
	@# Detect architecture and OS
	@if [ "$$(uname -m)" = "x86_64" ]; then \
		ARCH="amd64"; \
	elif [ "$$(uname -m)" = "aarch64" ] || [ "$$(uname -m)" = "arm64" ]; then \
		ARCH="arm64"; \
	elif [ "$$(uname -m)" = "armv6l" ]; then \
		ARCH="armv6l"; \
	else \
		$(call print_error,"Unsupported architecture: $$(uname -m)"); \
		exit 1; \
	fi; \
	if [ "$$(uname -s)" = "Linux" ]; then \
		OS="linux"; \
	elif [ "$$(uname -s)" = "Darwin" ]; then \
		OS="darwin"; \
	else \
		$(call print_error,"Unsupported operating system: $$(uname -s)"); \
		exit 1; \
	fi; \
	GO_ARCHIVE="go$(GO_VERSION).$$OS-$$ARCH.tar.gz"; \
	$(call print_status,"Downloading $$GO_ARCHIVE..."); \
	mkdir -p $(HOME)/.local; \
	cd $(HOME)/.local && \
	wget -q "https://golang.org/dl/$$GO_ARCHIVE" && \
	$(call print_status,"Extracting Go..."); \
	tar -xzf "$$GO_ARCHIVE" && \
	rm "$$GO_ARCHIVE"; \
	$(call print_status,"Setting up Go environment..."); \
	mkdir -p $(HOME)/.local/bin; \
	if [ ! -L "$(HOME)/.local/bin/go" ]; then \
		ln -sf $(HOME)/.local/go/bin/go $(HOME)/.local/bin/go; \
	fi; \
	if [ ! -L "$(HOME)/.local/bin/gofmt" ]; then \
		ln -sf $(HOME)/.local/go/bin/gofmt $(HOME)/.local/bin/gofmt; \
	fi; \
	$(call print_success,"Go $(GO_VERSION) installed successfully"); \
	$(call print_warning,"Please ensure $(HOME)/.local/bin is in your PATH"); \
	export GOROOT=$(HOME)/.local/go; \
	export PATH=$(HOME)/.local/go/bin:$$PATH

# System dependencies installation
install-system-deps:
	$(call print_status,"Installing system dependencies...")
	@if command -v apt >/dev/null 2>&1; then \
		echo "$(CYAN)➡️  Detected apt package manager (Debian/Ubuntu)$(RESET)"; \
		sudo apt update && \
		sudo apt install -y git build-essential curl wget unzip python3-pip; \
	elif command -v pacman >/dev/null 2>&1; then \
		echo "$(CYAN)➡️  Detected pacman package manager (Arch Linux)$(RESET)"; \
		sudo pacman -Sy --noconfirm git base-devel curl wget unzip python python-pip; \
	elif command -v brew >/dev/null 2>&1; then \
		echo "$(CYAN)➡️  Detected brew package manager (macOS)$(RESET)"; \
		brew install git curl wget unzip python3; \
	elif command -v yum >/dev/null 2>&1; then \
		echo "$(CYAN)➡️  Detected yum package manager (CentOS/RHEL)$(RESET)"; \
		sudo yum update -y && \
		sudo yum install -y git gcc make curl wget unzip python3-pip; \
	elif command -v dnf >/dev/null 2>&1; then \
		echo "$(CYAN)➡️  Detected dnf package manager (Fedora)$(RESET)"; \
		sudo dnf update -y && \
		sudo dnf install -y git gcc make curl wget unzip python3-pip; \
	else \
		echo "$(YELLOW)⚠️  No supported package manager found, attempting manual installation...$(RESET)"; \
		echo "$(CYAN)➡️  Please ensure git, curl, wget, unzip are installed manually$(RESET)"; \
	fi
	$(call print_success,"System dependencies installed")

# Rust installation
install-rust:
	$(call print_status,"Checking Rust installation...")
	@if ! command -v cargo >/dev/null 2>&1; then \
		echo "$(CYAN)➡️  Installing Rust toolchain...$(RESET)"; \
		curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y; \
		echo "$(GREEN)✅ Rust installed successfully$(RESET)"; \
		echo "$(YELLOW)⚠️  Please restart your terminal or run: source $$HOME/.cargo/env$(RESET)"; \
	else \
		echo "$(GREEN)✅ Rust is already installed$(RESET)"; \
	fi

# ===== TOOL INSTALLATION =====
install-tools: install-go-tools install-rust-tools install-python-tools install-other-tools

# Go-based tools
install-go-tools:
	$(call print_status,"Installing Go-based reconnaissance tools...")
	
	# ProjectDiscovery tools - Build locally
	@if [ ! -f "$(LOCAL_BIN_DIR)/naabu" ]; then \
		echo "$(CYAN)➡️  Building naabu (fast port scanner)...$(RESET)"; \
		if [ -d "$(TOOLS_DIR)/naabu" ]; then \
			cd $(TOOLS_DIR)/naabu && git pull; \
		else \
			cd $(TOOLS_DIR) && git clone https://github.com/projectdiscovery/naabu.git; \
		fi; \
		cd $(TOOLS_DIR)/naabu/cmd/naabu && go build -o $(LOCAL_BIN_DIR)/naabu .; \
	fi
	
	@if [ ! -f "$(LOCAL_BIN_DIR)/httpx" ]; then \
		echo "$(CYAN)➡️  Building httpx (HTTP probing tool)...$(RESET)"; \
		if [ -d "$(TOOLS_DIR)/httpx" ]; then \
			cd $(TOOLS_DIR)/httpx && git pull; \
		else \
			cd $(TOOLS_DIR) && git clone https://github.com/projectdiscovery/httpx.git; \
		fi; \
		cd $(TOOLS_DIR)/httpx/cmd/httpx && go build -o $(LOCAL_BIN_DIR)/httpx .; \
	fi
	
	@if [ ! -f "$(LOCAL_BIN_DIR)/subfinder" ]; then \
		echo "$(CYAN)➡️  Building subfinder (subdomain discovery)...$(RESET)"; \
		if [ -d "$(TOOLS_DIR)/subfinder" ]; then \
			cd $(TOOLS_DIR)/subfinder && git pull; \
		else \
			cd $(TOOLS_DIR) && git clone https://github.com/projectdiscovery/subfinder.git; \
		fi; \
		cd $(TOOLS_DIR)/subfinder/v2/cmd/subfinder && go build -o $(LOCAL_BIN_DIR)/subfinder .; \
	fi
	
	@if [ ! -f "$(LOCAL_BIN_DIR)/dnsx" ]; then \
		echo "$(CYAN)➡️  Building dnsx (DNS toolkit)...$(RESET)"; \
		if [ -d "$(TOOLS_DIR)/dnsx" ]; then \
			cd $(TOOLS_DIR)/dnsx && git pull; \
		else \
			cd $(TOOLS_DIR) && git clone https://github.com/projectdiscovery/dnsx.git; \
		fi; \
		cd $(TOOLS_DIR)/dnsx/cmd/dnsx && go build -o $(LOCAL_BIN_DIR)/dnsx .; \
	fi
	
	@if [ ! -f "$(LOCAL_BIN_DIR)/nuclei" ]; then \
		echo "$(CYAN)➡️  Building nuclei (vulnerability scanner)...$(RESET)"; \
		if [ -d "$(TOOLS_DIR)/nuclei" ]; then \
			cd $(TOOLS_DIR)/nuclei && git pull; \
		else \
			cd $(TOOLS_DIR) && git clone https://github.com/projectdiscovery/nuclei.git; \
		fi; \
		cd $(TOOLS_DIR)/nuclei/cmd/nuclei && go build -o $(LOCAL_BIN_DIR)/nuclei .; \
	fi
	
	@if [ ! -f "$(LOCAL_BIN_DIR)/katana" ]; then \
		echo "$(CYAN)➡️  Building katana (web crawler)...$(RESET)"; \
		if [ -d "$(TOOLS_DIR)/katana" ]; then \
			cd $(TOOLS_DIR)/katana && git pull; \
		else \
			cd $(TOOLS_DIR) && git clone https://github.com/projectdiscovery/katana.git; \
		fi; \
		cd $(TOOLS_DIR)/katana/cmd/katana && go build -o $(LOCAL_BIN_DIR)/katana .; \
	fi
	
	@if [ ! -f "$(LOCAL_BIN_DIR)/shuffledns" ]; then \
		echo "$(CYAN)➡️  Building shuffledns (DNS bruteforcer)...$(RESET)"; \
		if [ -d "$(TOOLS_DIR)/shuffledns" ]; then \
			cd $(TOOLS_DIR)/shuffledns && git pull; \
		else \
			cd $(TOOLS_DIR) && git clone https://github.com/projectdiscovery/shuffledns.git; \
		fi; \
		cd $(TOOLS_DIR)/shuffledns/cmd/shuffledns && go build -o $(LOCAL_BIN_DIR)/shuffledns .; \
	fi
	
	@if [ ! -f "$(LOCAL_BIN_DIR)/mapcidr" ]; then \
		echo "$(CYAN)➡️  Building mapcidr (CIDR manipulation)...$(RESET)"; \
		if [ -d "$(TOOLS_DIR)/mapcidr" ]; then \
			cd $(TOOLS_DIR)/mapcidr && git pull; \
		else \
			cd $(TOOLS_DIR) && git clone https://github.com/projectdiscovery/mapcidr.git; \
		fi; \
		cd $(TOOLS_DIR)/mapcidr/cmd/mapcidr && go build -o $(LOCAL_BIN_DIR)/mapcidr .; \
	fi
	
	@if [ ! -f "$(LOCAL_BIN_DIR)/asnmap" ]; then \
		echo "$(CYAN)➡️  Building asnmap (ASN mapping)...$(RESET)"; \
		if [ -d "$(TOOLS_DIR)/asnmap" ]; then \
			cd $(TOOLS_DIR)/asnmap && git pull; \
		else \
			cd $(TOOLS_DIR) && git clone https://github.com/projectdiscovery/asnmap.git; \
		fi; \
		cd $(TOOLS_DIR)/asnmap/cmd/asnmap && go build -o $(LOCAL_BIN_DIR)/asnmap .; \
	fi
	
	# Other Go tools
	@if [ ! -f "$(LOCAL_BIN_DIR)/amass" ]; then \
		echo "$(CYAN)➡️  Building amass (subdomain enumeration)...$(RESET)"; \
		if [ -d "$(TOOLS_DIR)/amass" ]; then \
			cd $(TOOLS_DIR)/amass && git pull; \
		else \
			cd $(TOOLS_DIR) && git clone https://github.com/owasp-amass/amass.git; \
		fi; \
		cd $(TOOLS_DIR)/amass/cmd/amass && go build -o $(LOCAL_BIN_DIR)/amass .; \
	fi
	
	@if [ ! -f "$(LOCAL_BIN_DIR)/ffuf" ]; then \
		echo "$(CYAN)➡️  Building ffuf (web fuzzer)...$(RESET)"; \
		if [ -d "$(TOOLS_DIR)/ffuf" ]; then \
			cd $(TOOLS_DIR)/ffuf && git pull; \
		else \
			cd $(TOOLS_DIR) && git clone https://github.com/ffuf/ffuf.git; \
		fi; \
		cd $(TOOLS_DIR)/ffuf && go build -o $(LOCAL_BIN_DIR)/ffuf .; \
	fi
	
	@if [ ! -f "$(LOCAL_BIN_DIR)/puredns" ]; then \
		echo "$(CYAN)➡️  Building puredns (DNS resolver)...$(RESET)"; \
		if [ -d "$(TOOLS_DIR)/puredns" ]; then \
			cd $(TOOLS_DIR)/puredns && git pull; \
		else \
			cd $(TOOLS_DIR) && git clone https://github.com/d3mondev/puredns.git; \
		fi; \
		cd $(TOOLS_DIR)/puredns/v2 && go build -o $(LOCAL_BIN_DIR)/puredns .; \
	fi
	
	$(call print_success,"Go-based tools built locally")

# Rust-based tools
install-rust-tools:
	$(call print_status,"Installing Rust-based reconnaissance tools...")
	
	@if [ ! -f "$(LOCAL_BIN_DIR)/feroxbuster" ]; then \
		echo "$(CYAN)➡️  Building feroxbuster (directory buster)...$(RESET)"; \
		if [ -d "$(TOOLS_DIR)/feroxbuster" ]; then \
			cd $(TOOLS_DIR)/feroxbuster && git pull; \
		else \
			cd $(TOOLS_DIR) && git clone https://github.com/epi052/feroxbuster.git; \
		fi; \
		cd $(TOOLS_DIR)/feroxbuster && cargo build --release; \
		cp target/release/feroxbuster $(LOCAL_BIN_DIR)/; \
	fi
	
	$(call print_success,"Rust-based tools built locally")

# Python-based tools
install-python-tools:
	$(call print_status,"Installing Python-based reconnaissance tools...")
	
	@if [ ! -f "$(LOCAL_BIN_DIR)/sublist3r" ]; then \
		echo "$(CYAN)➡️  Installing Sublist3r (subdomain enumeration)...$(RESET)"; \
		if [ -d "$(TOOLS_DIR)/Sublist3r" ]; then \
			cd $(TOOLS_DIR)/Sublist3r && git pull; \
		else \
			cd $(TOOLS_DIR) && git clone https://github.com/aboul3la/Sublist3r.git; \
		fi; \
		cd $(TOOLS_DIR)/Sublist3r && pip3 install -r requirements.txt; \
		echo '#!/bin/bash' > $(LOCAL_BIN_DIR)/sublist3r; \
		echo "python3 $(TOOLS_DIR)/Sublist3r/sublist3r.py \"\$$@\"" >> $(LOCAL_BIN_DIR)/sublist3r; \
		chmod +x $(LOCAL_BIN_DIR)/sublist3r; \
	fi
	
	$(call print_success,"Python-based tools installed locally")

# Other tools (C/C++, etc.)
install-other-tools:
	$(call print_status,"Installing other reconnaissance tools...")
	
	@if ! command -v massdns >/dev/null 2>&1; then \
		echo "$(CYAN)➡️  Installing massdns (high-performance DNS resolver)...$(RESET)"; \
		if command -v brew >/dev/null 2>&1; then \
			echo "$(CYAN)➡️  Installing massdns via Homebrew...$(RESET)"; \
			brew install massdns; \
		else \
			echo "$(CYAN)➡️  Installing massdns from source...$(RESET)"; \
			if [ -d "$(TOOLS_DIR)/massdns" ]; then \
				cd $(TOOLS_DIR)/massdns && git pull; \
			else \
				cd $(TOOLS_DIR) && git clone https://github.com/blechschmidt/massdns.git; \
			fi; \
			cd $(TOOLS_DIR)/massdns && make; \
			cp bin/massdns $(HOME)/.local/bin/ 2>/dev/null || cp bin/massdns $(LOCAL_BIN_DIR)/; \
		fi; \
	fi
	
	@if ! command -v cname-permutator >/dev/null 2>&1; then \
		echo "$(CYAN)➡️  Creating cname-permutator script (subdomain takeover detection)...$(RESET)"; \
		mkdir -p $(HOME)/.local/bin 2>/dev/null || mkdir -p $(LOCAL_BIN_DIR); \
		echo '#!/bin/bash' > $(HOME)/.local/bin/cname-permutator 2>/dev/null || echo '#!/bin/bash' > $(LOCAL_BIN_DIR)/cname-permutator; \
		echo '# Simple CNAME permutator for subdomain takeover detection' >> $(HOME)/.local/bin/cname-permutator 2>/dev/null || echo '# Simple CNAME permutator for subdomain takeover detection' >> $(LOCAL_BIN_DIR)/cname-permutator; \
		echo 'echo "CNAME permutator - checking subdomain takeover possibilities"' >> $(HOME)/.local/bin/cname-permutator 2>/dev/null || echo 'echo "CNAME permutator - checking subdomain takeover possibilities"' >> $(LOCAL_BIN_DIR)/cname-permutator; \
		echo 'for domain in "$$@"; do' >> $(HOME)/.local/bin/cname-permutator 2>/dev/null || echo 'for domain in "$$@"; do' >> $(LOCAL_BIN_DIR)/cname-permutator; \
		echo '  echo "Checking $$domain"' >> $(HOME)/.local/bin/cname-permutator 2>/dev/null || echo '  echo "Checking $$domain"' >> $(LOCAL_BIN_DIR)/cname-permutator; \
		echo '  dig +short CNAME "$$domain" | head -5' >> $(HOME)/.local/bin/cname-permutator 2>/dev/null || echo '  dig +short CNAME "$$domain" | head -5' >> $(LOCAL_BIN_DIR)/cname-permutator; \
		echo 'done' >> $(HOME)/.local/bin/cname-permutator 2>/dev/null || echo 'done' >> $(LOCAL_BIN_DIR)/cname-permutator; \
		chmod +x $(HOME)/.local/bin/cname-permutator 2>/dev/null || chmod +x $(LOCAL_BIN_DIR)/cname-permutator; \
	fi
	
	$(call print_success,"Other tools installed")

# ===== WORDLIST INSTALLATION =====
install-wordlists:
	$(call print_status,"Downloading reconnaissance wordlists...")
	
	@if [ ! -d "$(WORDLISTS_DIR)/SecLists" ]; then \
		$(call print_status,"Downloading SecLists..."); \
		cd $(WORDLISTS_DIR) && git clone https://github.com/danielmiessler/SecLists.git; \
	else \
		$(call print_status,"Updating SecLists..."); \
		cd $(WORDLISTS_DIR)/SecLists && git pull; \
	fi
	
	@if [ ! -f "$(WORDLISTS_DIR)/rockyou.txt" ]; then \
		$(call print_status,"Downloading rockyou.txt..."); \
		cd $(WORDLISTS_DIR) && \
		wget -q https://github.com/brannondorsey/naive-hashcat/releases/download/data/rockyou.txt || \
		(wget -q https://gitlab.com/kalilinux/packages/wordlists/-/raw/kali/master/rockyou.txt.gz && gunzip rockyou.txt.gz); \
	fi
	
	@if [ ! -d "$(WORDLISTS_DIR)/commonspeak2-wordlists" ]; then \
		$(call print_status,"Downloading commonspeak2 wordlists..."); \
		cd $(WORDLISTS_DIR) && git clone https://github.com/assetnote/commonspeak2-wordlists.git; \
	else \
		$(call print_status,"Updating commonspeak2 wordlists..."); \
		cd $(WORDLISTS_DIR)/commonspeak2-wordlists && git pull; \
	fi
	
	@if [ ! -d "$(WORDLISTS_DIR)/OneListForAll" ]; then \
		$(call print_status,"Downloading OneListForAll (optimized web fuzzing wordlists)..."); \
		cd $(WORDLISTS_DIR) && git clone https://github.com/six2dez/OneListForAll.git; \
	else \
		$(call print_status,"Updating OneListForAll..."); \
		cd $(WORDLISTS_DIR)/OneListForAll && git pull; \
	fi
	
	@if [ ! -d "$(WORDLISTS_DIR)/dirbuster-wordlists" ]; then \
		$(call print_status,"Downloading classic DirBuster wordlists..."); \
		mkdir -p $(WORDLISTS_DIR)/dirbuster-wordlists; \
		cd $(WORDLISTS_DIR)/dirbuster-wordlists && \
		wget -q https://raw.githubusercontent.com/daviddias/node-dirbuster/master/lists/directory-list-2.3-small.txt && \
		wget -q https://raw.githubusercontent.com/daviddias/node-dirbuster/master/lists/directory-list-2.3-medium.txt && \
		wget -q https://raw.githubusercontent.com/danielmiessler/SecLists/master/Discovery/Web-Content/directory-list-2.3-big.txt && \
		wget -q https://raw.githubusercontent.com/digination/dirbuster-ng/master/wordlists/common.txt; \
	fi
	
	@if [ ! -d "$(WORDLISTS_DIR)/raft-wordlists" ]; then \
		$(call print_status,"Downloading RAFT wordlists (penetration testing focused)..."); \
		mkdir -p $(WORDLISTS_DIR)/raft-wordlists; \
		cd $(WORDLISTS_DIR)/raft-wordlists && \
		wget -q https://raw.githubusercontent.com/danielmiessler/SecLists/master/Discovery/Web-Content/raft-large-directories.txt && \
		wget -q https://raw.githubusercontent.com/danielmiessler/SecLists/master/Discovery/Web-Content/raft-large-files.txt && \
		wget -q https://raw.githubusercontent.com/danielmiessler/SecLists/master/Discovery/Web-Content/raft-large-words.txt && \
		wget -q https://raw.githubusercontent.com/danielmiessler/SecLists/master/Discovery/Web-Content/raft-medium-directories.txt && \
		wget -q https://raw.githubusercontent.com/danielmiessler/SecLists/master/Discovery/Web-Content/raft-medium-files.txt; \
	fi
	
	@if [ ! -d "$(WORDLISTS_DIR)/ctf-wordlists" ]; then \
		$(call print_status,"Downloading CTF-specific wordlists..."); \
		cd $(WORDLISTS_DIR) && git clone https://github.com/bryanmcnulty/ctf-wordlists.git; \
	else \
		$(call print_status,"Updating CTF wordlists..."); \
		cd $(WORDLISTS_DIR)/ctf-wordlists && git pull; \
	fi
	
	@if [ ! -d "$(WORDLISTS_DIR)/awesome-wordlists" ]; then \
		$(call print_status,"Downloading awesome-wordlists collection..."); \
		cd $(WORDLISTS_DIR) && git clone https://github.com/gmelodie/awesome-wordlists.git; \
	else \
		$(call print_status,"Updating awesome-wordlists..."); \
		cd $(WORDLISTS_DIR)/awesome-wordlists && git pull; \
	fi
	
	$(call print_success,"Wordlists downloaded and updated")

# ===== CLEANUP =====
clean:
	$(call print_status,"Cleaning up ipcrawler installation...")
	
	@# Remove tools directory
	@if [ -d "$(TOOLS_DIR)" ]; then \
		$(call print_status,"Removing tools directory..."); \
		rm -rf $(TOOLS_DIR); \
	fi
	
	@# Remove wordlists directory
	@if [ -d "$(WORDLISTS_DIR)" ]; then \
		$(call print_status,"Removing wordlists directory..."); \
		rm -rf $(WORDLISTS_DIR); \
	fi
	
	@# Remove local bin directory if it exists
	@if [ -d "$(LOCAL_BIN_DIR)" ]; then \
		$(call print_status,"Removing local bin directory..."); \
		rm -rf $(LOCAL_BIN_DIR); \
	fi
	
	@# Remove installed binaries from ~/.local/bin
	@$(call print_status,"Removing installed tool binaries...")
	@for tool in naabu httpx subfinder dnsx nuclei katana shuffledns mapcidr asnmap amass ffuf feroxbuster sublist3r massdns puredns cname-permutator; do \
		if [ -f "$(HOME)/.local/bin/$$tool" ]; then \
			rm -f $(HOME)/.local/bin/$$tool; \
			$(call print_status,"Removed $$tool"); \
		fi; \
	done
	
	$(call print_success,"Cleanup completed!")

# ===== UTILITY TARGETS =====
.PHONY: list-tools test-tools update-tools

list-tools:
	@echo "$(MAGENTA)📋 Installed reconnaissance tools:$(RESET)"
	@echo ""
	@echo "$(CYAN)ProjectDiscovery Suite:$(RESET)"
	@command -v naabu >/dev/null 2>&1 && echo "  $(GREEN)✅ naabu$(RESET) - Fast port scanner" || echo "  $(RED)❌ naabu$(RESET) - Not installed"
	@command -v httpx >/dev/null 2>&1 && echo "  $(GREEN)✅ httpx$(RESET) - HTTP probing tool" || echo "  $(RED)❌ httpx$(RESET) - Not installed"
	@command -v subfinder >/dev/null 2>&1 && echo "  $(GREEN)✅ subfinder$(RESET) - Subdomain discovery" || echo "  $(RED)❌ subfinder$(RESET) - Not installed"
	@command -v dnsx >/dev/null 2>&1 && echo "  $(GREEN)✅ dnsx$(RESET) - DNS toolkit" || echo "  $(RED)❌ dnsx$(RESET) - Not installed"
	@command -v nuclei >/dev/null 2>&1 && echo "  $(GREEN)✅ nuclei$(RESET) - Vulnerability scanner" || echo "  $(RED)❌ nuclei$(RESET) - Not installed"
	@command -v katana >/dev/null 2>&1 && echo "  $(GREEN)✅ katana$(RESET) - Web crawler" || echo "  $(RED)❌ katana$(RESET) - Not installed"
	@command -v shuffledns >/dev/null 2>&1 && echo "  $(GREEN)✅ shuffledns$(RESET) - DNS bruteforcer" || echo "  $(RED)❌ shuffledns$(RESET) - Not installed"
	@command -v mapcidr >/dev/null 2>&1 && echo "  $(GREEN)✅ mapcidr$(RESET) - CIDR manipulation" || echo "  $(RED)❌ mapcidr$(RESET) - Not installed"
	@command -v asnmap >/dev/null 2>&1 && echo "  $(GREEN)✅ asnmap$(RESET) - ASN mapping" || echo "  $(RED)❌ asnmap$(RESET) - Not installed"
	@echo ""
	@echo "$(CYAN)Other Tools:$(RESET)"
	@command -v amass >/dev/null 2>&1 && echo "  $(GREEN)✅ amass$(RESET) - Subdomain enumeration" || echo "  $(RED)❌ amass$(RESET) - Not installed"
	@command -v ffuf >/dev/null 2>&1 && echo "  $(GREEN)✅ ffuf$(RESET) - Web fuzzer" || echo "  $(RED)❌ ffuf$(RESET) - Not installed"
	@command -v feroxbuster >/dev/null 2>&1 && echo "  $(GREEN)✅ feroxbuster$(RESET) - Directory buster" || echo "  $(RED)❌ feroxbuster$(RESET) - Not installed"
	@command -v massdns >/dev/null 2>&1 && echo "  $(GREEN)✅ massdns$(RESET) - High-performance DNS resolver" || echo "  $(RED)❌ massdns$(RESET) - Not installed"
	@command -v puredns >/dev/null 2>&1 && echo "  $(GREEN)✅ puredns$(RESET) - DNS resolver using massdns" || echo "  $(RED)❌ puredns$(RESET) - Not installed"
	@command -v sublist3r >/dev/null 2>&1 && echo "  $(GREEN)✅ sublist3r$(RESET) - Subdomain bruteforce" || echo "  $(RED)❌ sublist3r$(RESET) - Not installed"
	@command -v cname-permutator >/dev/null 2>&1 && echo "  $(GREEN)✅ cname-permutator$(RESET) - Subdomain takeover detection" || echo "  $(RED)❌ cname-permutator$(RESET) - Not installed"

test-tools:
	$(call print_status,"Testing installed tools...")
	@echo "$(CYAN)Running basic tests on installed tools:$(RESET)"
	@command -v naabu >/dev/null 2>&1 && echo "  $(GREEN)✅ naabu version: $$(naabu -version 2>&1 | head -1)$(RESET)" || echo "  $(RED)❌ naabu not found$(RESET)"
	@command -v httpx >/dev/null 2>&1 && echo "  $(GREEN)✅ httpx version: $$(httpx -version 2>&1 | head -1)$(RESET)" || echo "  $(RED)❌ httpx not found$(RESET)"
	@command -v subfinder >/dev/null 2>&1 && echo "  $(GREEN)✅ subfinder version: $$(subfinder -version 2>&1 | head -1)$(RESET)" || echo "  $(RED)❌ subfinder not found$(RESET)"
	@command -v nuclei >/dev/null 2>&1 && echo "  $(GREEN)✅ nuclei version: $$(nuclei -version 2>&1 | head -1)$(RESET)" || echo "  $(RED)❌ nuclei not found$(RESET)"

update-tools:
	$(call print_status,"Updating all tools...")
	@# Update Go tools
	@go install -v github.com/projectdiscovery/naabu/v2/cmd/naabu@latest
	@go install -v github.com/projectdiscovery/httpx/cmd/httpx@latest
	@go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest
	@go install -v github.com/projectdiscovery/dnsx/cmd/dnsx@latest
	@go install -v github.com/projectdiscovery/nuclei/v3/cmd/nuclei@latest
	@go install github.com/projectdiscovery/katana/cmd/katana@latest
	@go install -v github.com/projectdiscovery/shuffledns/cmd/shuffledns@latest
	@go install -v github.com/projectdiscovery/mapcidr/cmd/mapcidr@latest
	@go install github.com/projectdiscovery/asnmap/cmd/asnmap@latest
	@go install -v github.com/owasp-amass/amass/v4/...@master
	@go install github.com/ffuf/ffuf/v2@latest
	@go install github.com/d3mondev/puredns/v2@latest
	@# Update Git repositories
	@if [ -d "$(TOOLS_DIR)/feroxbuster" ]; then cd $(TOOLS_DIR)/feroxbuster && git pull && cargo build --release; fi
	@if [ -d "$(TOOLS_DIR)/Sublist3r" ]; then cd $(TOOLS_DIR)/Sublist3r && git pull; fi
	@if [ -d "$(TOOLS_DIR)/massdns" ]; then cd $(TOOLS_DIR)/massdns && git pull && make; fi
	@if [ -d "$(TOOLS_DIR)/cname-permutator" ]; then cd $(TOOLS_DIR)/cname-permutator && git pull && cargo build --release; fi
	@# Update wordlists
	@if [ -d "$(WORDLISTS_DIR)/SecLists" ]; then cd $(WORDLISTS_DIR)/SecLists && git pull; fi
	@if [ -d "$(WORDLISTS_DIR)/commonspeak2-wordlists" ]; then cd $(WORDLISTS_DIR)/commonspeak2-wordlists && git pull; fi
	@if [ -d "$(WORDLISTS_DIR)/OneListForAll" ]; then cd $(WORDLISTS_DIR)/OneListForAll && git pull; fi
	@if [ -d "$(WORDLISTS_DIR)/ctf-wordlists" ]; then cd $(WORDLISTS_DIR)/ctf-wordlists && git pull; fi
	@if [ -d "$(WORDLISTS_DIR)/awesome-wordlists" ]; then cd $(WORDLISTS_DIR)/awesome-wordlists && git pull; fi
	$(call print_success,"All tools updated!")
