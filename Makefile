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
	@echo "  $(CURDIR)/wordlists/  - Downloaded wordlists"
	@echo "  $(HOME)/.local/bin    - Installed tool binaries"

install: check-deps setup-dirs install-deps install-tools install-wordlists
	$(call print_success,"ipcrawler installation completed!")
	@echo ""
	@echo "$(CYAN)📋 Next steps:$(RESET)"
	@echo "  1. Add $(HOME)/.local/bin to your PATH if not already done:"
	@echo "     export PATH=\"$(HOME)/.local/bin:\$$PATH\""
	@echo "  2. Tools are ready to use!"
	@echo "  3. Run 'make help' for available commands"

# ===== DEPENDENCY CHECKING =====
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
		GO_CURRENT_VERSION=$$(go version | grep -o 'go[0-9]\+\.[0-9]\+\.[0-9]\+' | head -1 | sed 's/go//'); \
		if [ "$$(printf '%s\n' "$(GO_VERSION)" "$$GO_CURRENT_VERSION" | sort -V | head -n1)" = "$(GO_VERSION)" ]; then \
			$(call print_success,"Go $$GO_CURRENT_VERSION is already installed and meets minimum requirement"); \
		else \
			$(call print_warning,"Go $$GO_CURRENT_VERSION is outdated, installing Go $(GO_VERSION)..."); \
			$(MAKE) download-go; \
		fi; \
	else \
		$(call print_status,"Go not found, installing Go $(GO_VERSION)..."); \
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
		$(call print_status,"Detected apt package manager (Debian/Ubuntu)"); \
		sudo apt update && \
		sudo apt install -y git build-essential curl wget unzip python3-pip; \
	elif command -v pacman >/dev/null 2>&1; then \
		$(call print_status,"Detected pacman package manager (Arch Linux)"); \
		sudo pacman -Sy --noconfirm git base-devel curl wget unzip python python-pip; \
	elif command -v brew >/dev/null 2>&1; then \
		$(call print_status,"Detected brew package manager (macOS)"); \
		brew install git curl wget unzip python3; \
	elif command -v yum >/dev/null 2>&1; then \
		$(call print_status,"Detected yum package manager (CentOS/RHEL)"); \
		sudo yum update -y && \
		sudo yum install -y git gcc make curl wget unzip python3-pip; \
	elif command -v dnf >/dev/null 2>&1; then \
		$(call print_status,"Detected dnf package manager (Fedora)"); \
		sudo dnf update -y && \
		sudo dnf install -y git gcc make curl wget unzip python3-pip; \
	else \
		$(call print_warning,"No supported package manager found, attempting manual installation..."); \
		$(call print_status,"Please ensure git, curl, wget, unzip are installed manually"); \
	fi
	$(call print_success,"System dependencies installed")

# Rust installation
install-rust:
	$(call print_status,"Checking Rust installation...")
	@if ! command -v cargo >/dev/null 2>&1; then \
		$(call print_status,"Installing Rust toolchain..."); \
		curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y; \
		$(call print_success,"Rust installed successfully"); \
		$(call print_warning,"Please restart your terminal or run: source $$HOME/.cargo/env"); \
	else \
		$(call print_success,"Rust is already installed"); \
	fi

# ===== TOOL INSTALLATION =====
install-tools: install-go-tools install-rust-tools install-python-tools install-other-tools

# Go-based tools
install-go-tools:
	$(call print_status,"Installing Go-based reconnaissance tools...")
	
	# ProjectDiscovery tools
	@if ! command -v naabu >/dev/null 2>&1; then \
		$(call print_status,"Installing naabu (fast port scanner)..."); \
		go install -v github.com/projectdiscovery/naabu/v2/cmd/naabu@latest; \
	fi
	
	@if ! command -v httpx >/dev/null 2>&1; then \
		$(call print_status,"Installing httpx (HTTP probing tool)..."); \
		go install -v github.com/projectdiscovery/httpx/cmd/httpx@latest; \
	fi
	
	@if ! command -v subfinder >/dev/null 2>&1; then \
		$(call print_status,"Installing subfinder (subdomain discovery)..."); \
		go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest; \
	fi
	
	@if ! command -v dnsx >/dev/null 2>&1; then \
		$(call print_status,"Installing dnsx (DNS toolkit)..."); \
		go install -v github.com/projectdiscovery/dnsx/cmd/dnsx@latest; \
	fi
	
	@if ! command -v nuclei >/dev/null 2>&1; then \
		$(call print_status,"Installing nuclei (vulnerability scanner)..."); \
		go install -v github.com/projectdiscovery/nuclei/v3/cmd/nuclei@latest; \
	fi
	
	@if ! command -v katana >/dev/null 2>&1; then \
		$(call print_status,"Installing katana (web crawler)..."); \
		go install github.com/projectdiscovery/katana/cmd/katana@latest; \
	fi
	
	@if ! command -v shuffledns >/dev/null 2>&1; then \
		$(call print_status,"Installing shuffledns (DNS bruteforcer)..."); \
		go install -v github.com/projectdiscovery/shuffledns/cmd/shuffledns@latest; \
	fi
	
	@if ! command -v mapcidr >/dev/null 2>&1; then \
		$(call print_status,"Installing mapcidr (CIDR manipulation)..."); \
		go install -v github.com/projectdiscovery/mapcidr/cmd/mapcidr@latest; \
	fi
	
	@if ! command -v asnmap >/dev/null 2>&1; then \
		$(call print_status,"Installing asnmap (ASN mapping)..."); \
		go install github.com/projectdiscovery/asnmap/cmd/asnmap@latest; \
	fi
	
	# Other Go tools
	@if ! command -v amass >/dev/null 2>&1; then \
		$(call print_status,"Installing amass (subdomain enumeration)..."); \
		go install -v github.com/owasp-amass/amass/v4/...@master; \
	fi
	
	@if ! command -v ffuf >/dev/null 2>&1; then \
		$(call print_status,"Installing ffuf (web fuzzer)..."); \
		go install github.com/ffuf/ffuf/v2@latest; \
	fi
	
	@if ! command -v puredns >/dev/null 2>&1; then \
		$(call print_status,"Installing puredns (DNS resolver)..."); \
		go install github.com/d3mondev/puredns/v2@latest; \
	fi
	
	$(call print_success,"Go-based tools installed")

# Rust-based tools
install-rust-tools:
	$(call print_status,"Installing Rust-based reconnaissance tools...")
	
	@if ! command -v feroxbuster >/dev/null 2>&1; then \
		$(call print_status,"Installing feroxbuster (directory buster)..."); \
		if [ -d "$(TOOLS_DIR)/feroxbuster" ]; then \
			cd $(TOOLS_DIR)/feroxbuster && git pull; \
		else \
			cd $(TOOLS_DIR) && git clone https://github.com/epi052/feroxbuster.git; \
		fi; \
		cd $(TOOLS_DIR)/feroxbuster && cargo build --release; \
		cp target/release/feroxbuster $(HOME)/.local/bin/ 2>/dev/null || cp target/release/feroxbuster $(LOCAL_BIN_DIR)/; \
	fi
	
	$(call print_success,"Rust-based tools installed")

# Python-based tools
install-python-tools:
	$(call print_status,"Installing Python-based reconnaissance tools...")
	
	@if ! command -v sublist3r >/dev/null 2>&1; then \
		$(call print_status,"Installing Sublist3r (subdomain enumeration)..."); \
		if [ -d "$(TOOLS_DIR)/Sublist3r" ]; then \
			cd $(TOOLS_DIR)/Sublist3r && git pull; \
		else \
			cd $(TOOLS_DIR) && git clone https://github.com/aboul3la/Sublist3r.git; \
		fi; \
		cd $(TOOLS_DIR)/Sublist3r && pip3 install -r requirements.txt; \
		echo '#!/bin/bash\npython3 $(TOOLS_DIR)/Sublist3r/sublist3r.py "$$@"' > $(HOME)/.local/bin/sublist3r 2>/dev/null || echo '#!/bin/bash\npython3 $(TOOLS_DIR)/Sublist3r/sublist3r.py "$$@"' > $(LOCAL_BIN_DIR)/sublist3r; \
		chmod +x $(HOME)/.local/bin/sublist3r 2>/dev/null || chmod +x $(LOCAL_BIN_DIR)/sublist3r; \
	fi
	
	$(call print_success,"Python-based tools installed")

# Other tools (C/C++, etc.)
install-other-tools:
	$(call print_status,"Installing other reconnaissance tools...")
	
	@if ! command -v massdns >/dev/null 2>&1; then \
		$(call print_status,"Installing massdns (high-performance DNS resolver)..."); \
		if [ -d "$(TOOLS_DIR)/massdns" ]; then \
			cd $(TOOLS_DIR)/massdns && git pull; \
		else \
			cd $(TOOLS_DIR) && git clone https://github.com/blechschmidt/massdns.git; \
		fi; \
		cd $(TOOLS_DIR)/massdns && make; \
		cp bin/massdns $(HOME)/.local/bin/ 2>/dev/null || cp bin/massdns $(LOCAL_BIN_DIR)/; \
	fi
	
	@if ! command -v cname-permutator >/dev/null 2>&1; then \
		$(call print_status,"Installing cname-permutator (subdomain takeover detection)..."); \
		if [ -d "$(TOOLS_DIR)/cname-permutator" ]; then \
			cd $(TOOLS_DIR)/cname-permutator && git pull; \
		else \
			cd $(TOOLS_DIR) && git clone https://github.com/Sh1yo/cname-permutator.git; \
		fi; \
		cd $(TOOLS_DIR)/cname-permutator && cargo build --release; \
		cp target/release/cname-permutator $(HOME)/.local/bin/ 2>/dev/null || cp target/release/cname-permutator $(LOCAL_BIN_DIR)/; \
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
