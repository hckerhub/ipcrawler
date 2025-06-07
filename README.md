# IPCrawler - Intelligent Recon Flow Engine

A modern, enterprise-grade terminal UI for cybersecurity reconnaissance automation built with Python and Textual. Features comprehensive tool management, cross-platform deployment, and zero-dependency setup.

## ✨ Features

- 🎨 **Modern Terminal UI** - Beautiful interface with ASCII art branding and intuitive navigation
- 🛠️ **Automated Tool Management** - One-command installation of 15+ reconnaissance tools
- 🎯 **Smart Target Configuration** - IP addresses, domains, and CIDR ranges support
- 📊 **Reconnaissance Workflow** - Guided tool selection and execution summary
- ⚡ **Cross-Platform Deployment** - Works on Linux, macOS with automatic dependency resolution
- 🔒 **Enterprise Ready** - Portable setup, version-controlled installations, team collaboration
- ⌨️ **Intuitive Controls** - Arrow keys, space, enter for seamless navigation

## 🛠️ Reconnaissance Arsenal

### ProjectDiscovery Suite
- **naabu** - Lightning-fast port scanner with SYN/ACK techniques
- **subfinder** - Passive subdomain discovery using 40+ sources
- **dnsx** - DNS toolkit with bruteforcing and resolution validation
- **httpx** - HTTP probing with technology detection
- **katana** - Next-generation web crawler and spider
- **nuclei** - Vulnerability scanner with 4,000+ templates
- **shuffledns** - DNS resolver wrapper for mass DNS resolution
- **mapcidr** - CIDR manipulation and subnet enumeration
- **asnmap** - ASN mapping and organization discovery

### Additional Tools
- **nmap** - Network discovery and security auditing
- **amass** - Attack surface mapping and asset discovery
- **ffuf** - High-performance web fuzzer
- **feroxbuster** - Fast directory and file brute-forcer
- **massdns** - High-performance DNS resolver
- **puredns** - DNS bruteforcing with wildcard filtering
- **sublist3r** - Python-based subdomain enumeration

### Wordlist Collections
- **SecLists** - Comprehensive penetration testing wordlists
- **CommonSpeak2** - Content discovery wordlists from real data
- **OneListForAll** - Optimized web fuzzing wordlists
- **RockyOu** - Classic password dictionary
- **RAFT Wordlists** - Directory and file discovery lists
- **CTF Wordlists** - Capture The Flag specific collections

## 🚀 Quick Start

### Development Machine Setup
```bash
# Clone the repository
git clone <your-repo-url>
cd ipcrawler

# One-command installation (installs Go, Rust, tools, wordlists)
make install

# Start the application
python ipcrawler.py
```

### Test Machine Deployment
```bash
# Clone and immediately deploy
git clone <your-repo-url>
cd ipcrawler
make install     # Gets fresh installation of all tools
python ipcrawler.py  # Ready to use!
```

## 📦 What `make install` Does

The installation process is completely automated and includes:

### 🔧 **Language Runtimes**
- **Go 1.21.5** - Downloaded and installed from official source
- **Rust** - Latest stable toolchain via rustup
- **Python 3** - System package manager installation

### 🛠️ **System Dependencies**
- **Git** - Version control for tool repositories
- **Build tools** - gcc, make, build-essential
- **Network tools** - curl, wget, unzip
- **Development libraries** - Required for compilation

### 🎯 **Tool Installation Matrix**
- **Go tools** - Compiled from source using `go install`
- **Rust tools** - Built using `cargo build --release`
- **Python tools** - Cloned and configured with pip
- **C/C++ tools** - Compiled with make/cmake

### 📚 **Wordlist Management**
- **Automatic downloads** - Curated security wordlists
- **Git synchronization** - Always up-to-date collections
- **Organized structure** - Categorized by use case

## 🎮 Usage Guide

### Starting IPCrawler
```bash
# Activate Python environment (if using venv)
source venv/bin/activate  # Optional

# Launch the application
python ipcrawler.py
```

### Navigation Controls
- **Arrow Keys** ↑↓ - Navigate through options
- **Space Bar** ⌨️ - Select/deselect tools
- **Enter** ↵ - Continue to next screen
- **Escape** ⎋ - Go back to previous screen
- **A** - Select all tools
- **C** - Clear all selections
- **Ctrl+C** - Quit application

### Workflow Steps
1. **Welcome Screen** - Project introduction and branding
2. **Tool Selection** - Choose reconnaissance tools for your scan
3. **Target Configuration** - Input IP addresses, domains, or ranges
4. **Scan Summary** - Review selected tools and targets
5. **Execution** - Run the reconnaissance workflow

## 🏗️ Project Architecture

```
ipcrawler/
├── 📁 src/                    # Source code
│   ├── app.py                 # Main application logic
│   ├── config.py              # Configuration management
│   ├── styles.tcss            # Terminal UI styling
│   ├── utils.py               # Utility functions
│   └── screens/               # UI screens
│       ├── welcome.py         # Welcome screen
│       ├── tool_selection.py  # Tool selection interface
│       ├── target_input.py    # Target configuration
│       └── summary.py         # Scan summary
├── 📄 ipcrawler.py           # Application entry point
├── 📄 Makefile               # Installation automation
├── 📄 requirements.txt       # Python dependencies
├── 📄 .gitignore            # Git exclusion rules
├── 📁 tools/                # Downloaded tool repositories (gitignored)
├── 📁 wordlists/            # Security wordlists (gitignored)
└── 📁 bin/                  # Compiled binaries (gitignored)
```

## 🌐 Cross-Platform Support

### Supported Operating Systems
- **Linux** - Ubuntu, Debian, CentOS, RHEL, Fedora, Arch Linux
- **macOS** - Intel and Apple Silicon (M1/M2)
- **Package Managers** - apt, pacman, brew, yum, dnf

### Architecture Support
- **x86_64 (amd64)** - Standard 64-bit Intel/AMD
- **ARM64 (aarch64)** - Apple Silicon, AWS Graviton
- **ARM v6** - Raspberry Pi and embedded systems

## 🔒 Enterprise Features

### 🚀 **Portable Deployment**
- **Lightweight Repository** - Only source code committed
- **One-Command Setup** - `make install` handles everything
- **Consistent Environments** - Same tools, same versions everywhere

### 👥 **Team Collaboration**
- **Development Workflow** - Install locally, push only code
- **Test Machine Ready** - Clone and run immediately
- **Version Control** - Tool versions locked and reproducible

### 🛡️ **Security Considerations**
- **No Secrets in Git** - API keys and configs excluded
- **Local Tool Storage** - Tools installed outside repository
- **Audit Trail** - Installation process fully logged

## 🧪 Development Commands

```bash
# Install all tools and dependencies
make install

# List installed tools and their status
make list-tools

# Test tool installations
make test-tools

# Update all tools to latest versions
make update-tools

# Clean installation (remove all tools)
make clean

# Show help
make help
```

## 🔄 Update Workflow

```bash
# Update the application
git pull

# Update all reconnaissance tools
make update-tools

# Update Python dependencies
pip install -r requirements.txt --upgrade
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes and test thoroughly
4. Commit: `git commit -m "Add feature description"`
5. Push: `git push origin feature-name`
6. Submit a pull request

## 📋 Requirements

### Minimum System Requirements
- **OS**: Linux, macOS
- **RAM**: 2GB available
- **Storage**: 5GB for tools and wordlists
- **Network**: Internet connection for downloads
- **Python**: 3.8+ (automatically managed)

### Automatically Installed
- Go 1.21.5
- Rust (latest stable)
- Git and build tools
- All reconnaissance tools
- Security wordlists

## ⚠️ Current Status

IPCrawler is in active development. The UI framework and tool management system are complete. Tool execution integration is in progress.

**Completed:**
- ✅ Modern terminal UI with full navigation
- ✅ Comprehensive tool installation system
- ✅ Cross-platform deployment automation
- ✅ Enterprise-grade project structure

**In Development:**
- 🔄 Tool execution engine
- 🔄 Results processing and reporting
- 🔄 Configuration persistence
- 🔄 Advanced scanning workflows

## 📄 License

This project is open source. See LICENSE file for details.

## 👨‍💻 Developer

**Created by hckerhub**

Built with ❤️ for the cybersecurity community.

---

*IPCrawler - Making reconnaissance accessible, automated, and awesome.*
