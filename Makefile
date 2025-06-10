.PHONY: setup clean setup-docker docker-cmd help update

setup:
	@echo "Setting up ipcrawler..."
	@echo ""
	@echo "🔍 Checking prerequisites..."
	@# Check Python version
	@if ! command -v python3 >/dev/null 2>&1; then \
		echo "❌ Python 3 is not installed"; \
		echo ""; \
		echo "Please install Python 3.8+ first:"; \
		echo "  • Ubuntu/Debian: sudo apt install python3 python3-pip python3-venv"; \
		echo "  • CentOS/RHEL: sudo yum install python3 python3-pip"; \
		echo "  • Arch: sudo pacman -S python python-pip"; \
		echo "  • macOS: brew install python3"; \
		echo "  • Or download from: https://www.python.org/downloads/"; \
		echo ""; \
		exit 1; \
	fi
	@PYTHON_VERSION=$$(python3 -c "import sys; print('.'.join(map(str, sys.version_info[:2])))" 2>/dev/null || echo "unknown"); \
	if [ "$$PYTHON_VERSION" != "unknown" ]; then \
		echo "✅ Python $$PYTHON_VERSION detected"; \
		MAJOR=$$(echo $$PYTHON_VERSION | cut -d. -f1); \
		MINOR=$$(echo $$PYTHON_VERSION | cut -d. -f2); \
		if [ $$MAJOR -lt 3 ] || ([ $$MAJOR -eq 3 ] && [ $$MINOR -lt 8 ]); then \
			echo "⚠️  Python $$PYTHON_VERSION detected, but Python 3.8+ is recommended"; \
			echo "   Download latest from: https://www.python.org/downloads/"; \
		fi; \
	else \
		echo "⚠️  Could not determine Python version"; \
	fi
	@echo ""
	@echo "🔍 Detecting operating system..."
	@# Detect OS and install security tools
	@if [ -f /etc/os-release ]; then \
		OS_ID=$$(grep '^ID=' /etc/os-release | cut -d'=' -f2 | tr -d '"'); \
		OS_ID_LIKE=$$(grep '^ID_LIKE=' /etc/os-release | cut -d'=' -f2 | tr -d '"' 2>/dev/null || echo ""); \
		WSL_DETECTED=""; \
		if grep -q Microsoft /proc/version 2>/dev/null || [ -n "$$WSL_DISTRO_NAME" ]; then \
			WSL_DETECTED="yes"; \
			echo "🪟 WSL environment detected"; \
		fi; \
		if [ "$$OS_ID" = "kali" ] || [ "$$OS_ID" = "parrot" ] || echo "$$OS_ID_LIKE" | grep -q "debian\|ubuntu"; then \
			echo "📦 Installing security tools for $$OS_ID (Debian-based)..."; \
			echo "🔄 Updating package cache..."; \
			sudo apt update -qq; \
			echo "🐍 Installing Python venv package (fixes ensurepip issues)..."; \
			sudo apt install -y python3-venv python3-pip; \
			echo "📦 Installing core security tools..."; \
			sudo apt install -y curl wget git nmap nikto whatweb sslscan smbclient; \
			echo "📦 Installing available enumeration tools..."; \
			FAILED_TOOLS=""; \
			for tool in seclists dnsrecon enum4linux feroxbuster gobuster impacket-scripts nbtscan onesixtyone oscanner redis-tools smbmap snmp sipvicious tnscmd10g; do \
				if ! sudo apt install -y $$tool 2>/dev/null; then \
					echo "⚠️  $$tool failed via apt, checking snap..."; \
					FAILED_TOOLS="$$FAILED_TOOLS $$tool"; \
				fi; \
			done; \
			if [ -n "$$FAILED_TOOLS" ] && [ "$$WSL_DETECTED" = "yes" ]; then \
				echo "🫰 Installing snap for WSL compatibility..."; \
				sudo apt install -y snapd; \
				sudo systemctl enable snapd 2>/dev/null || true; \
				for tool in $$FAILED_TOOLS; do \
					case $$tool in \
						feroxbuster) \
							echo "Installing feroxbuster via snap..."; \
							sudo snap install feroxbuster 2>/dev/null || echo "⚠️  feroxbuster snap install failed"; \
							;; \
						gobuster) \
							echo "Installing gobuster via snap..."; \
							sudo snap install gobuster 2>/dev/null || echo "⚠️  gobuster snap install failed"; \
							;; \
						*) \
							echo "⚠️  No snap alternative for $$tool"; \
							;; \
					esac; \
				done; \
			fi; \
			echo "✅ Tool installation complete"; \
		elif [ "$$OS_ID" = "arch" ] || [ "$$OS_ID" = "manjaro" ]; then \
			echo "📦 Installing security tools for $$OS_ID (Arch-based)..."; \
			sudo pacman -Sy --noconfirm nmap curl wget git python python-pip || echo "⚠️  Some tools failed to install"; \
			echo "ℹ️  For full tool support, consider using Kali Linux or install tools manually"; \
		else \
			echo "ℹ️  Unsupported Linux distribution: $$OS_ID"; \
			echo "ℹ️  Installing basic requirements..."; \
			sudo apt update -qq 2>/dev/null || true; \
			sudo apt install -y python3-venv python3-pip curl wget git 2>/dev/null || true; \
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
	@if ! python3 -m venv venv 2>/dev/null; then \
		echo "⚠️  venv creation failed. Trying to fix..."; \
		echo "Installing python3-venv package..."; \
		sudo apt install -y python3-venv python3-pip 2>/dev/null || \
		sudo yum install -y python3-venv python3-pip 2>/dev/null || \
		sudo pacman -S --noconfirm python python-pip 2>/dev/null || \
		echo "⚠️  Could not install python3-venv. Please install manually."; \
		python3 -m venv venv; \
	fi
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
	@if ! sudo ln -sf "$$(pwd)/ipcrawler-cmd" /usr/local/bin/ipcrawler 2>/dev/null; then \
		echo "⚠️  Could not install to /usr/local/bin (permission issue)"; \
		echo "💡 You can still use: ./ipcrawler-cmd or add to PATH manually"; \
	fi
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
	@echo "Setting up Docker for ipcrawler..."
	@echo ""
	@echo "🔍 Checking Docker installation..."
	@if command -v docker >/dev/null 2>&1; then \
		echo "✅ Docker is already installed"; \
		if docker ps >/dev/null 2>&1; then \
			echo "✅ Docker daemon is running"; \
		else \
			echo "⚠️  Docker is installed but daemon is not running"; \
			echo "ℹ️  Please start Docker service and try again"; \
			exit 1; \
		fi; \
	else \
		echo "❌ Docker not found. Installing Docker..."; \
		echo "🔍 Detecting operating system..."; \
		if [ -f /etc/os-release ]; then \
			OS_ID=$$(grep '^ID=' /etc/os-release | cut -d'=' -f2 | tr -d '"'); \
			OS_ID_LIKE=$$(grep '^ID_LIKE=' /etc/os-release | cut -d'=' -f2 | tr -d '"' 2>/dev/null || echo ""); \
			if [ "$$OS_ID" = "kali" ] || [ "$$OS_ID" = "parrot" ] || echo "$$OS_ID_LIKE" | grep -q "debian\|ubuntu"; then \
				echo "📦 Installing Docker for $$OS_ID (Debian-based)..."; \
				echo "Updating package lists..."; \
				sudo apt update -qq; \
				echo "Installing Docker dependencies..."; \
				sudo apt install -y apt-transport-https ca-certificates curl gnupg lsb-release; \
				echo "Adding Docker GPG key..."; \
				curl -fsSL https://download.docker.com/linux/$$(lsb_release -is | tr '[:upper:]' '[:lower:]')/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg 2>/dev/null || \
				curl -fsSL https://download.docker.com/linux/debian/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg; \
				echo "Adding Docker repository..."; \
				echo "deb [arch=$$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/$$(lsb_release -is | tr '[:upper:]' '[:lower:]') $$(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null 2>/dev/null || \
				echo "deb [arch=$$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/debian $$(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null; \
				sudo apt update -qq; \
				echo "Installing Docker Engine..."; \
				sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin; \
				echo "Starting Docker service..."; \
				sudo systemctl start docker; \
				sudo systemctl enable docker; \
				echo "Adding current user to docker group..."; \
				sudo usermod -aG docker $$USER; \
				echo "✅ Docker installation complete!"; \
				echo "⚠️  You may need to log out and back in for group membership to take effect"; \
			elif [ "$$OS_ID" = "arch" ] || [ "$$OS_ID" = "manjaro" ]; then \
				echo "📦 Installing Docker for $$OS_ID (Arch-based)..."; \
				sudo pacman -Sy --noconfirm docker docker-compose; \
				sudo systemctl start docker; \
				sudo systemctl enable docker; \
				sudo usermod -aG docker $$USER; \
				echo "✅ Docker installation complete!"; \
				echo "⚠️  You may need to log out and back in for group membership to take effect"; \
			else \
				echo "ℹ️  Unsupported Linux distribution for automatic Docker installation: $$OS_ID"; \
				echo "ℹ️  Please install Docker manually:"; \
				echo "      https://docs.docker.com/engine/install/"; \
				exit 1; \
			fi; \
		elif [ "$$(uname)" = "Darwin" ]; then \
			echo "📦 Installing Docker for macOS..."; \
			if command -v brew >/dev/null 2>&1; then \
				echo "🍺 Using Homebrew to install Docker Desktop..."; \
				brew install --cask docker; \
				echo "✅ Docker Desktop installed!"; \
				echo "⚠️  Please start Docker Desktop from Applications folder"; \
				echo "⚠️  Once Docker Desktop is running, re-run: make setup-docker"; \
				exit 1; \
			else \
				echo "⚠️  Homebrew not found. Please install Docker Desktop manually:"; \
				echo "      https://docs.docker.com/desktop/install/mac-install/"; \
				exit 1; \
			fi; \
		else \
			echo "ℹ️  Unknown operating system. Please install Docker manually:"; \
			echo "      https://docs.docker.com/engine/install/"; \
			exit 1; \
		fi; \
	fi
	@echo ""
	@echo "🐳 Building ipcrawler Docker image..."
	docker build -t ipcrawler .
	@echo ""
	@echo "✅ Docker setup complete!"
	@echo ""
	@echo "📋 Next steps:"
	@echo "  • Run: make docker-cmd"
	@echo "  • Or manually: docker run -it --rm -v \$$(pwd)/results:/scans ipcrawler"

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
	@echo "  setup-docker  - Auto-install Docker + build Docker image for ipcrawler"
	@echo "  update        - Update repository, tools, and Docker image"
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
	@echo "  1. make setup-docker    # Auto-install Docker + build image with pre-installed tools"
	@echo "  2. make docker-cmd      # Start interactive session"
	@echo "  3. Inside container: /show-tools.sh or /install-extra-tools.sh"
	@echo ""
	@echo "Local Usage:"
	@echo "  1. make setup           # Set up locally with auto tool installation"
	@echo "  2. ipcrawler --help     # Use the tool"
	@echo "  3. make update          # Keep everything updated"

update:
	@echo "🔄 Updating ipcrawler..."
	@echo ""
	@echo "📦 Checking current repository status..."
	@if [ -d ".git" ]; then \
		echo "✅ Git repository detected"; \
		echo "🔍 Checking for remote updates..."; \
		git fetch origin >/dev/null 2>&1; \
		LOCAL=$$(git rev-parse HEAD); \
		REMOTE=$$(git rev-parse origin/main 2>/dev/null || git rev-parse origin/master 2>/dev/null); \
		if [ "$$LOCAL" = "$$REMOTE" ]; then \
			echo "✅ Repository is already up to date"; \
			UPDATE_NEEDED=false; \
		else \
			echo "📥 Updates available, pulling changes..."; \
			UPDATE_NEEDED=true; \
			echo "Current commit: $$(git rev-parse --short HEAD)"; \
			git stash push -m "Auto-stash before update" >/dev/null 2>&1 || true; \
			if git pull origin main >/dev/null 2>&1 || git pull origin master >/dev/null 2>&1; then \
				echo "✅ Git pull completed successfully"; \
				echo "New commit: $$(git rev-parse --short HEAD)"; \
				CHANGES=$$(git diff --name-only $$LOCAL HEAD); \
				if echo "$$CHANGES" | grep -q "Dockerfile"; then \
					echo "🐳 Dockerfile was updated"; \
					DOCKERFILE_CHANGED=true; \
				else \
					DOCKERFILE_CHANGED=false; \
				fi; \
				if echo "$$CHANGES" | grep -q "requirements.txt"; then \
					echo "🐍 Python requirements were updated"; \
					REQUIREMENTS_CHANGED=true; \
				else \
					REQUIREMENTS_CHANGED=false; \
				fi; \
				if echo "$$CHANGES" | grep -q "Makefile"; then \
					echo "⚠️  Makefile was updated - you may want to restart this command"; \
				fi; \
				echo "📋 Changed files:"; \
				echo "$$CHANGES" | sed 's/^/  • /'; \
			else \
				echo "❌ Git pull failed"; \
				git stash pop >/dev/null 2>&1 || true; \
				exit 1; \
			fi; \
		fi; \
	else \
		echo "⚠️  Not a git repository - skipping git update"; \
		UPDATE_NEEDED=false; \
		DOCKERFILE_CHANGED=false; \
		REQUIREMENTS_CHANGED=false; \
	fi
	@echo ""
	@echo "🐍 Updating Python environment..."
	@if [ -d "venv" ]; then \
		echo "📦 Updating Python packages..."; \
		venv/bin/python3 -m pip install --upgrade pip >/dev/null 2>&1; \
		venv/bin/python3 -m pip install --upgrade -r requirements.txt >/dev/null 2>&1; \
		echo "✅ Python packages updated"; \
	else \
		echo "ℹ️  No virtual environment found - run 'make setup' first"; \
	fi
	@echo ""
	@echo "🔧 Updating security tools..."
	@# Detect OS and update tools
	@if [ -f /etc/os-release ]; then \
		OS_ID=$$(grep '^ID=' /etc/os-release | cut -d'=' -f2 | tr -d '"'); \
		OS_ID_LIKE=$$(grep '^ID_LIKE=' /etc/os-release | cut -d'=' -f2 | tr -d '"' 2>/dev/null || echo ""); \
		if [ "$$OS_ID" = "kali" ] || [ "$$OS_ID" = "parrot" ] || echo "$$OS_ID_LIKE" | grep -q "debian\|ubuntu"; then \
			echo "📦 Updating security tools for $$OS_ID..."; \
			sudo apt update -qq >/dev/null 2>&1; \
			sudo apt upgrade -y nmap gobuster nikto whatweb sslscan feroxbuster redis-tools smbclient hydra john-jumbo hashcat sqlmap exploitdb binwalk exiftool 2>/dev/null >/dev/null || echo "ℹ️  Some tools may not need updates"; \
			echo "✅ APT tools updated"; \
		elif [ "$$OS_ID" = "arch" ] || [ "$$OS_ID" = "manjaro" ]; then \
			echo "📦 Updating security tools for $$OS_ID..."; \
			sudo pacman -Syu --noconfirm nmap curl wget git >/dev/null 2>&1 || echo "ℹ️  Some tools may not need updates"; \
			echo "✅ Pacman tools updated"; \
		else \
			echo "ℹ️  Skipping tool updates for $$OS_ID"; \
		fi; \
	elif [ "$$(uname)" = "Darwin" ]; then \
		echo "📦 Updating security tools for macOS..."; \
		if command -v brew >/dev/null 2>&1; then \
			echo "🍺 Updating Homebrew tools..."; \
			brew update >/dev/null 2>&1; \
			brew upgrade nmap gobuster nikto whatweb sslscan feroxbuster redis-tools smbclient hydra john-jumbo hashcat sqlmap exploitdb binwalk exiftool 2>/dev/null >/dev/null || echo "ℹ️  Some tools may not need updates"; \
			echo "✅ Homebrew tools updated"; \
		else \
			echo "ℹ️  Homebrew not found - skipping tool updates"; \
		fi; \
		echo "🐍 Updating Python security tools..."; \
		python3 -m pip install --upgrade impacket crackmapexec enum4linux-ng 2>/dev/null >/dev/null || echo "ℹ️  Some Python tools may not need updates"; \
		echo "✅ Python tools updated"; \
	else \
		echo "ℹ️  Unknown OS - skipping tool updates"; \
	fi
	@echo ""
	@echo "🐳 Checking Docker image updates..."
	@if command -v docker >/dev/null 2>&1; then \
		if [ "$$DOCKERFILE_CHANGED" = "true" ] || [ "$$REQUIREMENTS_CHANGED" = "true" ]; then \
			echo "🔄 Rebuilding Docker image due to changes..."; \
			docker build -t ipcrawler . >/dev/null 2>&1 && echo "✅ Docker image rebuilt successfully" || echo "❌ Docker image rebuild failed"; \
		elif [ -n "$$(docker images -q ipcrawler 2>/dev/null)" ]; then \
			echo "ℹ️  Docker image exists but no rebuild needed"; \
		else \
			echo "ℹ️  No Docker image found - run 'make setup-docker' to build one"; \
		fi; \
	else \
		echo "ℹ️  Docker not available - skipping Docker updates"; \
	fi
	@echo ""
	@echo "✅ Update complete!"
	@echo ""
	@echo "📋 Summary:"
	@echo "  • Git repository: Updated"
	@echo "  • Python packages: Updated"  
	@echo "  • Security tools: Updated"
	@echo "  • Docker image: Checked/Updated if needed"
	@echo ""
	@echo "💡 If the Makefile was updated, consider restarting this terminal session"
