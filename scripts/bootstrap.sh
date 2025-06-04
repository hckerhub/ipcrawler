#!/bin/bash

# IPCrawler Complete Bootstrap Script
# Comprehensive system setup for all dependencies and tools
# This script prepares your system for current and future IPCrawler needs

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Smart dependency detection functions
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check if a package is installed via package manager
package_installed() {
    local package="$1"
    case $OS in
        "macos")
            brew list "$package" &> /dev/null
            ;;
        "linux")
            case $PACKAGE_MANAGER in
                "apt")
                    dpkg -l | grep -q "^ii.*$package " 2>/dev/null
                    ;;
                "dnf"|"yum")
                    rpm -qa | grep -q "$package" 2>/dev/null
                    ;;
                "pacman")
                    pacman -Qi "$package" &> /dev/null
                    ;;
                "apk")
                    apk info | grep -q "$package" 2>/dev/null
                    ;;
                *)
                    command_exists "$package"
                    ;;
            esac
            ;;
        *)
            command_exists "$package"
            ;;
    esac
}

# Check if Go tool is installed
go_tool_installed() {
    local tool_name="$1"
    local tool_path="$TOOLS_DIR/go-tools/bin/$tool_name"
    [[ -f "$tool_path" ]] && [[ -x "$tool_path" ]]
}

# Check if Python package is installed
python_package_installed() {
    local package="$1"
    pipx list | grep -q "$package" 2>/dev/null || pip3 show "$package" &>/dev/null
}

# Configuration
BINARY_NAME="ipcrawler"
MISC_DIR="misc"
TOOLS_DIR="tools"
CACHE_DIR=".cache"

# ASCII Art
cat << 'EOF'
 ██╗██████╗  ██████╗██████╗  █████╗ ██╗    ██╗██╗     ███████╗██████╗ 
 ██║██╔══██╗██╔════╝██╔══██╗██╔══██╗██║    ██║██║     ██╔════╝██╔══██╗
 ██║██████╔╝██║     ██████╔╝███████║██║ █╗ ██║██║     █████╗  ██████╔╝
 ██║██╔═══╝ ██║     ██╔══██╗██╔══██║██║███╗██║██║     ██╔══╝  ██╔══██╗
 ██║██║     ╚██████╗██║  ██║██║  ██║╚███╔███╔╝███████╗███████╗██║  ██║
 ╚═╝╚═╝      ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝ ╚══╝╚══╝ ╚══════╝╚══════╝╚═╝  ╚═╝
                                                                       
        🚀 COMPLETE SYSTEM BOOTSTRAP & CRAWLER SETUP 🚀
        Installing everything you need for current & future use

EOF

echo -e "${BLUE}Starting IPCrawler Complete System Setup...${NC}\n"

# Function to print status messages
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_header() {
    echo -e "\n${CYAN}============================================${NC}"
    echo -e "${CYAN} $1${NC}"
    echo -e "${CYAN}============================================${NC}"
}

# Detect operating system and architecture
detect_system() {
    print_header "SYSTEM DETECTION"
    
    if [[ "$OSTYPE" == "darwin"* ]]; then
        OS="macos"
        ARCH=$(uname -m)
        DISTRO="macOS $(sw_vers -productVersion 2>/dev/null || echo 'Unknown')"
    elif [[ "$OSTYPE" == "linux-gnu"* ]] || [[ "$OSTYPE" == "linux"* ]]; then
        OS="linux"
        ARCH=$(uname -m)
        
        # Detect Linux distribution
        if [[ -f /etc/os-release ]]; then
            DISTRO=$(grep '^PRETTY_NAME=' /etc/os-release | cut -d'"' -f2)
        elif [[ -f /etc/redhat-release ]]; then
            DISTRO=$(cat /etc/redhat-release)
        elif [[ -f /etc/debian_version ]]; then
            DISTRO="Debian $(cat /etc/debian_version)"
        else
            DISTRO="Unknown Linux"
        fi
    elif [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "cygwin" ]]; then
        OS="windows"
        ARCH=$(uname -m)
        DISTRO="Windows"
    else
        OS="unknown"
        ARCH="unknown"
        DISTRO="Unknown OS"
    fi
    
    # Normalize architecture
    case $ARCH in
        x86_64|amd64) ARCH="amd64" ;;
        arm64|aarch64) ARCH="arm64" ;;
        armv7l) ARCH="arm" ;;
        i386|i686) ARCH="386" ;;
    esac
    
    print_status "Operating System: $OS"
    print_status "Architecture: $ARCH"
    print_status "Distribution: $DISTRO"
    print_status "Shell: $SHELL"
    
    # Detect available package managers
    detect_package_managers
}

# Detect available package managers
detect_package_managers() {
    print_status "Detecting package managers..."
    
    PACKAGE_MANAGERS=()
    
    case $OS in
        "macos")
            if command -v brew &> /dev/null; then
                PACKAGE_MANAGERS+=("brew")
                print_status "✅ Homebrew found"
            else
                print_warning "❌ Homebrew not found - will install"
                NEED_HOMEBREW=true
            fi
            ;;
        "linux")
            # Check for various Linux package managers
            if command -v apt-get &> /dev/null; then
                PACKAGE_MANAGERS+=("apt")
                print_status "✅ APT found"
            fi
            if command -v dnf &> /dev/null; then
                PACKAGE_MANAGERS+=("dnf")
                print_status "✅ DNF found"
            fi
            if command -v yum &> /dev/null; then
                PACKAGE_MANAGERS+=("yum")
                print_status "✅ YUM found"
            fi
            if command -v pacman &> /dev/null; then
                PACKAGE_MANAGERS+=("pacman")
                print_status "✅ Pacman found"
            fi
            if command -v apk &> /dev/null; then
                PACKAGE_MANAGERS+=("apk")
                print_status "✅ APK found"
            fi
            if command -v zypper &> /dev/null; then
                PACKAGE_MANAGERS+=("zypper")
                print_status "✅ Zypper found"
            fi
            ;;
        "windows")
            if command -v choco &> /dev/null; then
                PACKAGE_MANAGERS+=("choco")
                print_status "✅ Chocolatey found"
            fi
            if command -v winget &> /dev/null; then
                PACKAGE_MANAGERS+=("winget")
                print_status "✅ Winget found"
            fi
            ;;
    esac
    
    if [[ ${#PACKAGE_MANAGERS[@]} -eq 0 ]]; then
        print_warning "No package managers found - manual installation may be required"
    fi
}

# Create directory structure
setup_directories() {
    print_header "DIRECTORY ORGANIZATION"
    
    print_status "Creating organized directory structure..."
    
    # Create directories
    mkdir -p "$MISC_DIR"/{build-artifacts,logs,temp,old-versions}
    mkdir -p "$TOOLS_DIR"/{external,scripts,configs}
    mkdir -p "$CACHE_DIR"/{downloads,builds}
    mkdir -p "docs/generated"
    mkdir -p "config/profiles"
    
    print_status "✅ Created misc/ for non-essential files"
    print_status "✅ Created tools/ for external tools"
    print_status "✅ Created .cache/ for temporary data"
    print_status "✅ Created docs/generated/ for documentation"
    print_status "✅ Created config/profiles/ for user configurations"
}

# Install package manager if needed
install_package_manager() {
    if [[ "$NEED_HOMEBREW" == true ]] && [[ "$OS" == "macos" ]]; then
        print_header "INSTALLING HOMEBREW"
        print_status "Installing Homebrew package manager..."
        
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
        
        # Add Homebrew to PATH for this session
        if [[ -f "/opt/homebrew/bin/brew" ]]; then
            export PATH="/opt/homebrew/bin:$PATH"
        elif [[ -f "/usr/local/bin/brew" ]]; then
            export PATH="/usr/local/bin:$PATH"
        fi
        
        if command -v brew &> /dev/null; then
            print_success "✅ Homebrew installed successfully"
            PACKAGE_MANAGERS+=("brew")
        else
            print_error "❌ Homebrew installation failed"
            return 1
        fi
    fi
}

# Install core dependencies
install_core_dependencies() {
    print_header "CORE DEPENDENCIES INSTALLATION"
    
    # Install based on OS and package manager
    case $OS in
        "macos")
            install_macos_dependencies
            ;;
        "linux")
            install_linux_dependencies
            ;;
        "windows")
            install_windows_dependencies
            ;;
    esac
    
    # Install Python packages
    install_python_packages
    
    # Install Go tools
    install_go_tools
}

# Install macOS dependencies with smart detection
install_macos_dependencies() {
    if command -v brew &> /dev/null; then
        print_status "Checking and installing core dependencies via Homebrew..."
        
        # Define packages with their command names for verification
        local packages=(
            "git:git"
            "curl:curl" 
            "wget:wget"
            "go:go"
            "nmap:nmap"
            "masscan:masscan"
            "nikto:nikto"
            "gobuster:gobuster"
            "netcat:nc"
            "tcpdump:tcpdump"
            "httpie:http"
            "jq:jq"
            "sqlite3:sqlite3"
            "zip:zip"
            "unzip:unzip"
            "python3:python3"
        )
        
        local installed_count=0
        local skipped_count=0
        local failed_packages=()
        
        for package_pair in "${packages[@]}"; do
            local package="${package_pair%:*}"
            local cmd="${package_pair#*:}"
            
            # Check if command already exists
            if command_exists "$cmd"; then
                print_success "✅ $package already installed"
                skipped_count=$((skipped_count + 1))
                continue
            fi
            
            # Check if package is installed via brew but command missing (edge case)
            if package_installed "$package"; then
                print_success "✅ $package already installed via Homebrew"
                skipped_count=$((skipped_count + 1))
                continue
            fi
            
            # Install the package
            print_status "Installing $package..."
            if brew install "$package" 2>/dev/null; then
                print_success "✅ $package installed"
                installed_count=$((installed_count + 1))
            else
                print_warning "⚠️  Failed to install $package (might not be available)"
                failed_packages+=("$package")
            fi
        done
        
        # Summary
        print_status "📊 Dependencies Summary:"
        print_status "   • Already installed: $skipped_count"
        print_status "   • Newly installed: $installed_count" 
        
        if [[ ${#failed_packages[@]} -gt 0 ]]; then
            print_warning "   • Failed to install: ${failed_packages[*]}"
            print_status "You can install them manually later if needed"
        fi
        
        print_success "✅ Core dependencies installation completed"
    else
        print_error "❌ Homebrew not found - please install Homebrew first"
        print_status "Install Homebrew: /bin/bash -c \"\$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
        return 1
    fi
}

# Install Linux dependencies
install_linux_dependencies() {
    for pm in "${PACKAGE_MANAGERS[@]}"; do
        case $pm in
            "apt")
                print_status "Installing dependencies via APT..."
                sudo apt-get update
                sudo apt-get install -y git curl wget golang-go nmap masscan nikto dirb gobuster netcat-openbsd tcpdump httpie jq sqlite3 zip unzip python3 python3-pip vim nano
                ;;
            "dnf")
                print_status "Installing dependencies via DNF..."
                sudo dnf install -y git curl wget golang nmap masscan nikto dirb gobuster nc tcpdump httpie jq sqlite zip unzip python3 python3-pip vim nano
                ;;
            "yum")
                print_status "Installing dependencies via YUM..."
                sudo yum install -y git curl wget golang nmap masscan nikto dirb gobuster nc tcpdump jq sqlite zip unzip python3 python3-pip vim nano
                ;;
            "pacman")
                print_status "Installing dependencies via Pacman..."
                sudo pacman -S --noconfirm git curl wget go nmap masscan nikto dirb gobuster netcat tcpdump httpie jq sqlite zip unzip python python-pip vim nano
                ;;
            "apk")
                print_status "Installing dependencies via APK..."
                sudo apk add git curl wget go nmap masscan nikto dirb gobuster netcat-openbsd tcpdump httpie jq sqlite zip unzip python3 py3-pip vim nano
                ;;
        esac
        break # Use first available package manager
    done
}

# Install Windows dependencies
install_windows_dependencies() {
    if command -v choco &> /dev/null; then
        print_status "Installing dependencies via Chocolatey..."
        choco install -y git curl wget golang nmap sqlite zip unzip python3 vim nano
    elif command -v winget &> /dev/null; then
        print_status "Installing dependencies via Winget..."
        winget install Git.Git
        winget install GoLang.Go
        winget install Insecure.Nmap
        winget install Python.Python.3
    else
        print_warning "No package manager found for Windows - manual installation required"
    fi
}

# Install Python security packages with smart detection
install_python_packages() {
    print_header "PYTHON SECURITY PACKAGES"
    
    if command -v pip3 &> /dev/null; then
        print_status "Checking and installing Python security packages..."
        
        local installed_count=0
        local skipped_count=0
        local failed_tools=()
        
        # Try using pipx first (better for externally managed environments)
        if command -v pipx &> /dev/null || brew install pipx &>/dev/null; then
            print_status "Using pipx for Python package management..."
            local python_tools=("dirsearch" "sqlmap")
            
            for tool in "${python_tools[@]}"; do
                # Check if already installed
                if python_package_installed "$tool"; then
                    print_success "✅ $tool already installed"
                    skipped_count=$((skipped_count + 1))
                    continue
                fi
                
                # Install the package
                print_status "Installing $tool via pipx..."
                if pipx install "$tool" 2>/dev/null; then
                    print_success "✅ $tool installed via pipx"
                    installed_count=$((installed_count + 1))
                else
                    print_warning "⚠️  $tool installation failed"
                    failed_tools+=("$tool")
                fi
            done
        else
            # Fallback: use pip with --break-system-packages
            print_status "Using pip3 with --break-system-packages flag..."
            local pip_packages=("requests" "beautifulsoup4" "pwntools" "scapy")
            
            for package in "${pip_packages[@]}"; do
                if python_package_installed "$package"; then
                    print_success "✅ $package already installed"
                    skipped_count=$((skipped_count + 1))
                    continue
                fi
                
                print_status "Installing $package..."
                if pip3 install --user --break-system-packages "$package" 2>/dev/null; then
                    print_success "✅ $package installed"
                    installed_count=$((installed_count + 1))
                else
                    print_warning "⚠️  $package installation failed"
                    failed_tools+=("$package")
                fi
            done
        fi
        
        # Summary
        print_status "📊 Python Packages Summary:"
        print_status "   • Already installed: $skipped_count"
        print_status "   • Newly installed: $installed_count"
        
        if [[ ${#failed_tools[@]} -gt 0 ]]; then
            print_warning "   • Failed to install: ${failed_tools[*]}"
            print_status "You can install them manually later if needed"
        fi
        
        print_success "✅ Python security packages installation completed"
    else
        print_warning "pip3 not found - skipping Python packages"
    fi
}

# Install Go security tools with smart detection
install_go_tools() {
    print_header "GO SECURITY TOOLS"
    
    if command -v go &> /dev/null; then
        print_status "Checking and installing Go security tools..."
        
        # Create tools directory for Go binaries
        mkdir -p "$TOOLS_DIR/go-tools/bin"
        
        # Set GOPATH and GOBIN for tools
        export GOPATH="$PWD/$TOOLS_DIR/go-tools"
        export GOBIN="$PWD/$TOOLS_DIR/go-tools/bin"
        
        # Define Go tools with their binary names
        local go_tools=(
            "github.com/ffuf/ffuf@latest:ffuf"
            "github.com/tomnomnom/httprobe@latest:httprobe"
            "github.com/tomnomnom/assetfinder@latest:assetfinder"
            "github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest:subfinder"
            "github.com/projectdiscovery/httpx/cmd/httpx@latest:httpx"
            "github.com/projectdiscovery/nuclei/v2/cmd/nuclei@latest:nuclei"
            "github.com/owasp-amass/amass/v3/...@master:amass"
        )
        
        local installed_count=0
        local skipped_count=0
        local failed_tools=()
        
        for tool_pair in "${go_tools[@]}"; do
            local tool_url="${tool_pair%:*}"
            local tool_name="${tool_pair#*:}"
            
            # Check if tool is already installed
            if go_tool_installed "$tool_name"; then
                print_success "✅ $tool_name already installed"
                skipped_count=$((skipped_count + 1))
                continue
            fi
            
            # Install the tool
            print_status "Installing $tool_name..."
            if go install "$tool_url" 2>/dev/null; then
                print_success "✅ $tool_name installed"
                installed_count=$((installed_count + 1))
            else
                print_warning "⚠️  $tool_name installation failed"
                failed_tools+=("$tool_name")
            fi
        done
        
        # Summary
        print_status "📊 Go Tools Summary:"
        print_status "   • Already installed: $skipped_count"
        print_status "   • Newly installed: $installed_count"
        
        if [[ ${#failed_tools[@]} -gt 0 ]]; then
            print_warning "   • Failed to install: ${failed_tools[*]}"
            print_status "You can install them manually later if needed"
        fi
        
        print_success "✅ Go security tools installation completed"
        print_success "✅ Tools installed to: $TOOLS_DIR/go-tools/bin"
        
        # Add to PATH suggestion
        print_status "💡 Add to your PATH: export PATH=\"\$PATH:$PWD/$TOOLS_DIR/go-tools/bin\""
    else
        print_warning "Go not found - skipping Go tools"
    fi
}

# Build and install IPCrawler
build_and_install() {
    print_header "BUILDING & INSTALLING IPCRAWLER"
    
    # Reset Go environment for clean build
    unset GOPATH
    unset GOBIN
    export GOCACHE="$PWD/.cache/builds"
    export GOMODCACHE="$PWD/.cache/modules"
    
    # Clean Go cache to avoid module conflicts
    go clean -modcache 2>/dev/null || true
    go clean -cache 2>/dev/null || true
    
    # Move old builds to misc
    if [[ -f "$BINARY_NAME" ]]; then
        print_status "Moving old binary to misc/old-versions/"
        mv "$BINARY_NAME" "$MISC_DIR/old-versions/${BINARY_NAME}-$(date +%Y%m%d-%H%M%S)"
    fi
    
    # Temporarily move tools directory completely outside project root to avoid module conflicts
    if [[ -d "$TOOLS_DIR" ]]; then
        TOOLS_BACKUP="$HOME/.cache/ipcrawler-tools-backup-$$"
        mkdir -p "$(dirname "$TOOLS_BACKUP")"
        mv "$TOOLS_DIR" "$TOOLS_BACKUP"
        print_status "Temporarily moved tools directory outside project root for clean build"
    fi
    
    # Download dependencies  
    print_status "Downloading Go dependencies..."
    
    # Make sure we're in a clean state - remove any old backup directories
    rm -rf tools.backup tools.temp 2>/dev/null || true
    find . -name "tools.*" -type d -exec rm -rf {} + 2>/dev/null || true
    
    go mod download
    go mod tidy
    
    # Build with enhanced build info
    print_status "Building IPCrawler with complete build info..."
    BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
    VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev-$(date +%Y%m%d)")
    COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
    
    LDFLAGS="-ldflags \"-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.Commit=$COMMIT -X main.OS=$OS -X main.Arch=$ARCH\""
    
    eval "go build $LDFLAGS -o $BINARY_NAME ."
    
    if [[ ! -f "$BINARY_NAME" ]]; then
        print_error "❌ Build failed"
        return 1
    fi
    
    print_success "✅ Build completed successfully"
    
    # Restore tools directory
    if [[ -d "$TOOLS_BACKUP" ]]; then
        mv "$TOOLS_BACKUP" "$TOOLS_DIR"
        print_status "Restored tools directory"
    fi
    
    # Clean up any leftover backup directories
    rm -rf /tmp/ipcrawler-tools-backup-* "$HOME/.cache/ipcrawler-tools-backup-"* 2>/dev/null || true
    
    # Install to system PATH
    install_to_system
}

# Install binary to system PATH
install_to_system() {
    print_status "Installing to system PATH..."
    
    # Determine installation path
    case $OS in
        "macos"|"linux")
            if [[ -w "/usr/local/bin" ]] || command -v sudo &> /dev/null; then
                INSTALL_PATH="/usr/local/bin"
                if [[ ! -w "/usr/local/bin" ]]; then
                    sudo cp "$BINARY_NAME" "$INSTALL_PATH/"
                    sudo chmod +x "$INSTALL_PATH/$BINARY_NAME"
                else
                    cp "$BINARY_NAME" "$INSTALL_PATH/"
                    chmod +x "$INSTALL_PATH/$BINARY_NAME"
                fi
            else
                # Fallback to user bin
                INSTALL_PATH="$HOME/bin"
                mkdir -p "$INSTALL_PATH"
                cp "$BINARY_NAME" "$INSTALL_PATH/"
                chmod +x "$INSTALL_PATH/$BINARY_NAME"
                
                # Add to PATH if not already there
                add_to_user_path "$INSTALL_PATH"
            fi
            ;;
        "windows")
            INSTALL_PATH="$USERPROFILE/bin"
            mkdir -p "$INSTALL_PATH"
            cp "$BINARY_NAME" "$INSTALL_PATH/"
            ;;
    esac
    
    print_success "✅ Installed to $INSTALL_PATH"
    
    # Test installation
    if command -v "$BINARY_NAME" &> /dev/null; then
        print_success "✅ $BINARY_NAME is now available system-wide!"
    else
        print_warning "⚠️  $BINARY_NAME may not be in PATH. Restart terminal or source shell config."
    fi
}

# Add directory to user PATH
add_to_user_path() {
    local dir="$1"
    local shell_config=""
    
    case $SHELL in
        */bash)
            shell_config="$HOME/.bashrc"
            [[ ! -f "$shell_config" ]] && touch "$shell_config"
            ;;
        */zsh)
            shell_config="$HOME/.zshrc"
            ;;
        */fish)
            shell_config="$HOME/.config/fish/config.fish"
            mkdir -p "$(dirname "$shell_config")"
            ;;
        *)
            shell_config="$HOME/.profile"
            ;;
    esac
    
    if [[ -n "$shell_config" ]] && ! grep -q "export PATH.*$dir" "$shell_config" 2>/dev/null; then
        echo "" >> "$shell_config"
        echo "# Added by IPCrawler bootstrap" >> "$shell_config"
        echo "export PATH=\"$dir:\$PATH\"" >> "$shell_config"
        print_status "Added $dir to PATH in $shell_config"
    fi
}

# Clean up and organize files
cleanup_and_organize() {
    print_header "CLEANUP & ORGANIZATION"
    
    # Move build artifacts to misc
    print_status "Organizing build artifacts..."
    
    # Move coverage files
    [[ -f "coverage.out" ]] && mv coverage.out "$MISC_DIR/build-artifacts/"
    [[ -f "coverage.html" ]] && mv coverage.html "$MISC_DIR/build-artifacts/"
    
    # Move any leftover binaries
    find . -maxdepth 1 -name "${BINARY_NAME}-*" -type f -exec mv {} "$MISC_DIR/old-versions/" \; 2>/dev/null || true
    
    # Move logs if any
    find . -maxdepth 1 -name "*.log" -type f -exec mv {} "$MISC_DIR/logs/" \; 2>/dev/null || true
    
    # Move temporary files
    find . -maxdepth 1 -name "*.tmp" -type f -exec mv {} "$MISC_DIR/temp/" \; 2>/dev/null || true
    
    # Create useful scripts in tools/scripts/
    create_utility_scripts
    
    print_success "✅ Files organized and cleaned up"
}

# Create utility scripts
create_utility_scripts() {
    print_status "Creating utility scripts..."
    
    # Quick update script
    cat > "$TOOLS_DIR/scripts/update-ipcrawler.sh" << 'EOF'
#!/bin/bash
# Quick update script for IPCrawler
cd "$(dirname "$0")/../.."
git pull
make crawler
EOF
    chmod +x "$TOOLS_DIR/scripts/update-ipcrawler.sh"
    
    # System info script
    cat > "$TOOLS_DIR/scripts/system-info.sh" << 'EOF'
#!/bin/bash
# Show system information for IPCrawler
ipcrawler platform
EOF
    chmod +x "$TOOLS_DIR/scripts/system-info.sh"
    
    # Dependency check script
    cat > "$TOOLS_DIR/scripts/check-deps.sh" << 'EOF'
#!/bin/bash
# Check all dependencies for IPCrawler
echo "🔍 Checking IPCrawler Dependencies"
echo "=================================="

# Core tools
for tool in git go nmap curl wget; do
    if command -v "$tool" &> /dev/null; then
        echo "✅ $tool: $(which $tool)"
    else
        echo "❌ $tool: Not found"
    fi
done

# Security tools
echo -e "\n🔒 Security Tools:"
for tool in nikto dirb gobuster masscan; do
    if command -v "$tool" &> /dev/null; then
        echo "✅ $tool: $(which $tool)"
    else
        echo "❌ $tool: Not found"
    fi
done
EOF
    chmod +x "$TOOLS_DIR/scripts/check-deps.sh"
    
    print_status "✅ Created utility scripts in tools/scripts/"
}

# Generate documentation
generate_documentation() {
    print_header "DOCUMENTATION GENERATION"
    
    # Generate README for new directory structure
    cat > "docs/generated/DIRECTORY_STRUCTURE.md" << EOF
# IPCrawler Directory Structure

This document describes the organized directory structure created by the bootstrap process.

## Directory Layout

\`\`\`
ipcrawler/
├── misc/                     # Non-essential files
│   ├── build-artifacts/      # Coverage reports, build logs
│   ├── logs/                 # Application logs
│   ├── temp/                 # Temporary files
│   └── old-versions/         # Previous binary versions
├── tools/                    # External tools and utilities
│   ├── external/             # Downloaded external tools
│   ├── scripts/              # Utility scripts
│   ├── configs/              # Tool configurations
│   └── go-tools/             # Go security tools
├── .cache/                   # Temporary cache data
│   ├── downloads/            # Downloaded files
│   └── builds/               # Build cache
├── docs/generated/           # Generated documentation
└── config/profiles/          # User configuration profiles
\`\`\`

## Installed Tools

### Core Dependencies
- Go programming language
- Git version control
- NMAP network scanner
- curl/wget for HTTP requests

### Security Tools
- nikto - Web server scanner
- dirb - Directory brute forcer
- gobuster - Directory/file/DNS busting tool
- masscan - Fast port scanner

### Go Tools (in tools/go-tools/bin/)
- ffuf - Fast web fuzzer
- httprobe - HTTP probe
- subfinder - Subdomain finder
- httpx - HTTP toolkit
- nuclei - Vulnerability scanner
- assetfinder - Asset discovery

### Python Tools
- requests, beautifulsoup4, scrapy
- pwntools, scapy, impacket
- sqlmap, dirsearch

## Usage

After bootstrap completion, you can:

1. Run IPCrawler: \`ipcrawler\`
2. Check system info: \`ipcrawler platform\`
3. Update tools: \`./tools/scripts/update-ipcrawler.sh\`
4. Check dependencies: \`./tools/scripts/check-deps.sh\`

## Future Additions

The bootstrap script is designed to be extensible. New tools and dependencies can be added by modifying the \`DEPENDENCIES\` array in the bootstrap script.
EOF

    print_success "✅ Documentation generated in docs/generated/"
}

# Show completion summary
show_completion_summary() {
    print_header "BOOTSTRAP COMPLETION SUMMARY"
    
    echo -e "${GREEN}🎉 IPCrawler Complete System Bootstrap Finished!${NC}\n"
    
    echo -e "${CYAN}📋 System Information:${NC}"
    echo "  OS: $DISTRO ($OS $ARCH)"
    echo "  Shell: $SHELL"
    echo "  Package Managers: ${PACKAGE_MANAGERS[*]}"
    
    echo -e "\n${CYAN}📦 Installation Summary:${NC}"
    echo "  Binary: $(which $BINARY_NAME 2>/dev/null || echo 'Local build only')"
    echo "  Go Tools: $TOOLS_DIR/go-tools/bin/"
    echo "  Utility Scripts: $TOOLS_DIR/scripts/"
    echo "  Documentation: docs/generated/"
    
    echo -e "\n${CYAN}🚀 Quick Start Commands:${NC}"
    echo "  ${PURPLE}ipcrawler${NC}                           # Start IPCrawler"
    echo "  ${PURPLE}ipcrawler platform${NC}                  # Show system info"
    echo "  ${PURPLE}./tools/scripts/check-deps.sh${NC}       # Check dependencies"
    echo "  ${PURPLE}./tools/scripts/update-ipcrawler.sh${NC} # Update IPCrawler"
    
    echo -e "\n${CYAN}📁 Directory Organization:${NC}"
    echo "  ${PURPLE}misc/${NC}        # Non-essential files (logs, old versions)"
    echo "  ${PURPLE}tools/${NC}       # External tools and utilities" 
    echo "  ${PURPLE}.cache/${NC}      # Temporary cache data"
    echo "  ${PURPLE}docs/generated/${NC} # Generated documentation"
    
    echo -e "\n${CYAN}🎯 Key Bindings for $OS:${NC}"
    case $OS in
        "macos")
            echo "  Copy: ${PURPLE}Cmd+C${NC}, Paste: ${PURPLE}Cmd+V${NC}, Quit: ${PURPLE}Cmd+Q${NC}"
            ;;
        "windows")
            echo "  Copy: ${PURPLE}Ctrl+C${NC}, Paste: ${PURPLE}Ctrl+V${NC}, Quit: ${PURPLE}Alt+F4${NC}"
            ;;
        *)
            echo "  Copy: ${PURPLE}Ctrl+Shift+C${NC}, Paste: ${PURPLE}Ctrl+Shift+V${NC}, Quit: ${PURPLE}Ctrl+C${NC}"
            ;;
    esac
    
    echo -e "\n${GREEN}✅ Your system is now fully prepared for current and future IPCrawler needs!${NC}"
    echo -e "${YELLOW}⚠️  Remember: Only use these tools on systems you own or have permission to test!${NC}"
    echo ""
    echo -e "${GREEN}Happy hacking! 🎯${NC}"
}

# Main bootstrap function
main() {
    detect_system
    setup_directories
    install_package_manager
    install_core_dependencies
    build_and_install
    cleanup_and_organize
    generate_documentation
    show_completion_summary
}

# Run main function
main "$@" 