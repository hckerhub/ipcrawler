.PHONY: setup clean setup-docker docker-cmd help

setup:
	@echo "Setting up ipcrawler..."
	@echo ""
	@echo "🔍 Detecting operating system..."
	@# Detect OS and install security tools
	@if [ -f /etc/os-release ]; then \
		OS_ID=$$(grep '^ID=' /etc/os-release | cut -d'=' -f2 | tr -d '"'); \
		OS_ID_LIKE=$$(grep '^ID_LIKE=' /etc/os-release | cut -d'=' -f2 | tr -d '"' 2>/dev/null || echo ""); \
		if [ "$$OS_ID" = "kali" ] || [ "$$OS_ID" = "parrot" ] || echo "$$OS_ID_LIKE" | grep -q "debian\|ubuntu"; then \
			echo "📦 Installing security tools for $$OS_ID (Debian-based)..."; \
			sudo apt update -qq; \
			sudo apt install -y seclists curl dnsrecon enum4linux feroxbuster gobuster impacket-scripts nbtscan nikto nmap onesixtyone oscanner redis-tools smbclient smbmap snmp sslscan sipvicious tnscmd10g whatweb 2>/dev/null || \
			echo "⚠️  Some tools failed to install - continuing anyway"; \
		elif [ "$$OS_ID" = "arch" ] || [ "$$OS_ID" = "manjaro" ]; then \
			echo "📦 Installing security tools for $$OS_ID (Arch-based)..."; \
			sudo pacman -Sy --noconfirm nmap curl wget git || echo "⚠️  Some tools failed to install"; \
			echo "ℹ️  For full tool support, consider using Kali Linux or install tools manually"; \
		else \
			echo "ℹ️  Unsupported Linux distribution: $$OS_ID"; \
			echo "ℹ️  Please install security tools manually or use Docker setup"; \
		fi; \
	elif [ "$$(uname)" = "Darwin" ]; then \
		echo "📦 Installing comprehensive security toolkit for macOS..."; \
		if command -v brew >/dev/null 2>&1; then \
			echo "🍺 Using Homebrew to install security tools..."; \
			echo "Installing core network tools..."; \
			brew install nmap curl wget git gobuster nikto whatweb sslscan || echo "⚠️  Some core tools failed"; \
			echo "Installing enumeration tools..."; \
			brew install feroxbuster redis-tools smbclient || echo "⚠️  Some enum tools failed"; \
			echo "Installing additional security tools..."; \
			brew install hydra john-jumbo hashcat sqlmap exploitdb binwalk exiftool || echo "⚠️  Some additional tools failed"; \
			echo "Installing Python security tools via pip..."; \
			python3 -m pip install impacket crackmapexec enum4linux-ng 2>/dev/null || echo "⚠️  Some Python tools failed"; \
			echo "✅ macOS security toolkit installation complete!"; \
			echo "📋 Installed tools: nmap, gobuster, nikto, whatweb, sslscan,"; \
			echo "    feroxbuster, redis-tools, smbclient, hydra, john-jumbo,"; \
			echo "    hashcat, sqlmap, exploitdb, binwalk, exiftool, impacket, crackmapexec"; \
		else \
			echo "⚠️  Homebrew not found. Please install Homebrew first:"; \
			echo "    /bin/bash -c \"\$$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""; \
			echo "ℹ️  Or use Docker setup for full tool support"; \
		fi; \
	else \
		echo "ℹ️  Unknown operating system. Please install security tools manually or use Docker setup"; \
	fi
	@echo ""
	@echo "🐍 Setting up Python environment..."
	python3 -m venv venv
	venv/bin/python3 -m pip install --upgrade pip
	venv/bin/python3 -m pip install -r requirements.txt
	@echo "Creating ipcrawler command..."
	@rm -f ipcrawler-cmd
	@echo '#!/bin/bash' > ipcrawler-cmd
	@echo '# Resolve the real path of the script (follow symlinks)' >> ipcrawler-cmd
	@echo 'SCRIPT_PATH="$$(realpath "$${BASH_SOURCE[0]}")"' >> ipcrawler-cmd
	@echo 'DIR="$$(cd "$$(dirname "$$SCRIPT_PATH")" && pwd)"' >> ipcrawler-cmd
	@echo 'source "$$DIR/venv/bin/activate" && python3 "$$DIR/ipcrawler.py" "$$@"' >> ipcrawler-cmd
	@chmod +x ipcrawler-cmd
	@echo "Installing ipcrawler command to /usr/local/bin..."
	@sudo ln -sf "$$(pwd)/ipcrawler-cmd" /usr/local/bin/ipcrawler
	@echo ""
	@echo "✅ Setup complete!"
	@echo ""
	@echo "📋 Next steps:"
	@echo "  • Run: ipcrawler --help"
	@echo "  • Test with: ipcrawler 127.0.0.1"
	@echo "  • For full tool support on non-Kali systems, consider: make setup-docker"

clean:
	@echo "Cleaning up ipcrawler installation..."
	@echo ""
	@echo "🧹 Removing virtual environment and command..."
	@VENV_ACTIVE=""; \
	if [ -n "$$VIRTUAL_ENV" ]; then \
		echo "⚠️  You are currently in a virtual environment"; \
		VENV_ACTIVE="yes"; \
	fi; \
	export VENV_WAS_ACTIVE="$$VENV_ACTIVE"
	rm -rf venv .venv
	rm -f ipcrawler-cmd
	@echo "Removing ipcrawler from /usr/local/bin..."
	@sudo rm -f /usr/local/bin/ipcrawler
	@echo ""
	@echo "🗑️  Removing installed security tools..."
	@# Detect OS and remove tools that were installed by setup
	@if [ -f /etc/os-release ]; then \
		OS_ID=$$(grep '^ID=' /etc/os-release | cut -d'=' -f2 | tr -d '"'); \
		OS_ID_LIKE=$$(grep '^ID_LIKE=' /etc/os-release | cut -d'=' -f2 | tr -d '"' 2>/dev/null || echo ""); \
		if [ "$$OS_ID" = "kali" ] || [ "$$OS_ID" = "parrot" ] || echo "$$OS_ID_LIKE" | grep -q "debian\|ubuntu"; then \
			echo "Removing security tools for $$OS_ID..."; \
			echo "⚠️  Note: This will remove security tools that may be used by other applications"; \
			if [ -t 0 ]; then \
				read -p "Remove security tools? [y/N]: " -n 1 -r; \
				echo; \
			else \
				echo "Reading confirmation from input..."; \
				read -r REPLY; \
			fi; \
			if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
				echo "Removing apt-installed security tools..."; \
				for tool in seclists dnsrecon enum4linux feroxbuster gobuster impacket-scripts nbtscan nikto onesixtyone oscanner redis-tools smbclient smbmap snmp sslscan sipvicious tnscmd10g whatweb; do \
					if dpkg -l | grep -q "^ii.*$$tool" 2>/dev/null; then \
						echo "Removing $$tool..."; \
						sudo apt remove -y $$tool 2>/dev/null || echo "Failed to remove $$tool"; \
					fi; \
				done; \
				echo "Core tools (nmap, curl, wget, git) left installed."; \
			else \
				echo "Skipping tool removal."; \
			fi; \
		elif [ "$$OS_ID" = "arch" ] || [ "$$OS_ID" = "manjaro" ]; then \
			echo "ℹ️  Arch-based system detected. Basic tools (nmap, curl, wget, git) left installed."; \
		else \
			echo "ℹ️  No tools to remove for $$OS_ID."; \
		fi; \
	elif [ "$$(uname)" = "Darwin" ]; then \
		echo "Removing security tools for macOS..."; \
		echo "⚠️  Note: This will remove security tools that may be used by other applications"; \
		if [ -t 0 ]; then \
			read -p "Remove security tools? [y/N]: " -n 1 -r; \
			echo; \
		else \
			echo "Reading confirmation from input..."; \
			read -r REPLY; \
		fi; \
		if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
			echo "Removing Homebrew security tools..."; \
			for tool in nmap gobuster nikto whatweb sslscan feroxbuster redis-tools smbclient hydra john-jumbo hashcat sqlmap exploitdb binwalk exiftool; do \
				if brew list | grep -q "^$$tool$$" 2>/dev/null; then \
					echo "Removing $$tool..."; \
					brew uninstall --ignore-dependencies $$tool 2>/dev/null || echo "Failed to remove $$tool"; \
				fi; \
			done; \
			echo "Removing Python security tools..."; \
			python3 -m pip uninstall -y impacket crackmapexec enum4linux-ng 2>/dev/null || echo "Some Python tools couldn't be removed"; \
			echo "Core tools (curl, wget, git) left installed."; \
		else \
			echo "Skipping tool removal."; \
		fi; \
	else \
		echo "ℹ️  Unknown OS - no tools to remove."; \
	fi
	@echo ""
	@echo "🐳 Cleaning up Docker resources..."
	@if [ -n "$$(docker images -q ipcrawler 2>/dev/null)" ]; then \
		echo "Stopping any running ipcrawler containers..."; \
		docker ps -aq --filter ancestor=ipcrawler 2>/dev/null | xargs -r docker stop >/dev/null 2>&1 || true; \
		docker ps -aq --filter ancestor=ipcrawler 2>/dev/null | xargs -r docker rm >/dev/null 2>&1 || true; \
		echo "Removing ipcrawler Docker image..."; \
		docker rmi ipcrawler >/dev/null 2>&1 || true; \
		echo "Docker image removed."; \
	else \
		echo "No ipcrawler Docker image found."; \
	fi
	@echo "Cleaning up results directory..."
	@if [ -d "results" ] && [ -z "$$(ls -A results 2>/dev/null)" ]; then \
		rm -rf results; \
		echo "Empty results directory removed."; \
	elif [ -d "results" ]; then \
		echo "Results directory contains files - keeping it."; \
	else \
		echo "No results directory found."; \
	fi
	@echo ""
	@echo "✅ Clean complete!"
	@echo "Virtual environment, Docker image, and empty directories removed."
	@if [ -n "$$VIRTUAL_ENV" ]; then \
		echo ""; \
		echo "┌─────────────────────────────────────────────────────────────┐"; \
		echo "│  ⚠️  WARNING: IMPORTANT FINAL STEP                          │"; \
		echo "│                                                             │"; \
		echo "│  You are still in a virtual environment!                   │"; \
		echo "│  Please run the following command:                         │"; \
		echo "│                                                             │"; \
		echo "│      deactivate                                             │"; \
		echo "│                                                             │"; \
		echo "│  This will restore your normal terminal prompt.            │"; \
		echo "└─────────────────────────────────────────────────────────────┘"; \
		echo ""; \
	fi

setup-docker:
	@echo "Building ipcrawler Docker image..."
	docker build -t ipcrawler .
	@echo ""
	@echo "✓ Docker setup complete!"
	@echo "Now you can run: make docker-cmd"
	@echo "Or manually: docker run -it --rm -v \$$(pwd)/results:/scans ipcrawler"

docker-cmd:
	@echo "Starting ipcrawler Docker container..."
	@echo "Results will be saved to: $$(pwd)/results"
	@echo "Type 'exit' to leave the container"
	@echo ""
	docker run -it --rm -v "$$(pwd)/results:/scans" ipcrawler || true

help:
	@echo "Available make commands:"
	@echo ""
	@echo "  setup         - Set up local Python virtual environment + install security tools"
	@echo "  clean         - Remove local setup, virtual environment, and Docker resources"
	@echo "  setup-docker  - Build Docker image for ipcrawler"
	@echo "  docker-cmd    - Run interactive Docker container"
	@echo "  help          - Show this help message"
	@echo ""
	@echo "Supported Operating Systems:"
	@echo "  • Kali Linux       - Full tool installation (20+ security tools)"
	@echo "  • Parrot OS        - Full tool installation (20+ security tools)"
	@echo "  • Ubuntu/Debian    - Full tool installation (20+ security tools)"
	@echo "  • macOS (Homebrew) - Comprehensive toolkit (15+ security tools)"
	@echo "  • Arch/Manjaro     - Basic tools (nmap, curl, wget, git)"
	@echo "  • Other systems    - Python setup only (use Docker for full features)"
	@echo ""
	@echo "Docker Usage (Recommended for non-Kali systems):"
	@echo "  1. make setup-docker    # Build image with pre-installed tools"
	@echo "  2. make docker-cmd      # Start interactive session"
	@echo "  3. Inside container: /show-tools.sh or /install-extra-tools.sh"
	@echo ""
	@echo "Local Usage:"
	@echo "  1. make setup           # Set up locally with auto tool installation"
	@echo "  2. ipcrawler --help     # Use the tool"
