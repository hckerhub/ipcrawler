# 🕷️ ipcrawler

> *"It's like bowling with bumpers."* - [@ippsec](https://twitter.com/ippsec)

A simplified, streamlined version of **AutoRecon** - the multi-threaded network reconnaissance tool that performs automated enumeration of services for CTFs, OSCP, and penetration testing environments.

## 🙏 Credits

**ipcrawler** is a fork of [**AutoRecon**](https://github.com/Tib3rius/AutoRecon) by [**Tib3rius**](https://github.com/Tib3rius). All core functionality, plugins, and the brilliant multi-threaded architecture are thanks to his incredible work. This fork simply provides a cleaner setup experience while maintaining all the powerful features of the original tool.

## ✨ What's New

**ipcrawler** takes AutoRecon's powerful enumeration capabilities and makes setup effortless:

### 🤔 Docker vs Local Setup

| Feature | Docker | Local |
|---------|--------|-------|
| **Setup Time** | 2 minutes | 5-10 minutes |
| **Dependencies** | None (isolated) | Manual tool installation |
| **Platform Support** | Windows, macOS, Linux | Linux/Unix only |
| **Resource Usage** | Higher (container overhead) | Lower (native) |
| **Tool Updates** | Rebuild image | Manual updates |
| **Cleanup** | `make clean` (removes everything) | `make clean` (removes everything) |
| **Recommended For** | Most users, beginners | Advanced users, performance |

| **Before (AutoRecon)** | **After (ipcrawler)** |
|---|---|
| `pipx install git+https://github.com/Tib3rius/AutoRecon.git` | `make setup` or `make setup-docker` |
| `sudo env "PATH=$PATH" autorecon target` | `ipcrawler target` |
| Complex dependency management | Automatic virtual environment or Docker |
| Manual uninstallation | `make clean` |
| Platform-specific installation issues | Docker works everywhere |

## 📋 Prerequisites

### 🪟 Windows Users
- **Docker Desktop for Windows** - [Download here](https://www.docker.com/products/docker-desktop)
- That's it! The `ipcrawler-windows.bat` handles everything else.

### 🐧 Linux Users
- **Python 3.8+** - [Download from python.org](https://www.python.org/downloads/) or install via package manager:
  ```bash
  # Ubuntu/Debian
  sudo apt update && sudo apt install python3 python3-pip python3-venv
  
  # CentOS/RHEL/Fedora
  sudo yum install python3 python3-pip  # or dnf install
  
  # Arch Linux
  sudo pacman -S python python-pip
  ```

- **make** (for convenience commands):
  ```bash
  # Ubuntu/Debian/Kali
  sudo apt install make
  
  # CentOS/RHEL/Fedora
  sudo yum install make  # or dnf install make
  
  # Arch Linux
  sudo pacman -S make
  ```

### 🍎 macOS Users
- **Python 3.8+** - [Download from python.org](https://www.python.org/downloads/) or use Homebrew:
  ```bash
  brew install python3
  ```

- **make** (usually pre-installed, or via Xcode tools):
  ```bash
  # Install Xcode Command Line Tools
  xcode-select --install
  
  # Or via Homebrew
  brew install make
  ```

## 🚀 Quick Start

### 🪟 Windows Users (Docker Only)

**The simplest way to run ipcrawler on Windows is with Docker:**

#### **🚀 One-Click Setup**
1. **Install Docker Desktop for Windows** from [docker.com](https://www.docker.com/products/docker-desktop)
2. **Double-click to run:**
```cmd
git clone https://github.com/hckerhub/ipcrawler.git
cd ipcrawler
ipcrawler-windows.bat
```

**What the launcher does:**
- ✅ Checks Docker is installed and running
- ✅ Builds Docker image automatically (first time only)  
- ✅ Opens interactive Docker terminal with all tools
- ✅ Mounts results folder to your Windows filesystem
- ✅ Shows helpful commands and usage tips

**Inside the container:**
```bash
# Available immediately:
ipcrawler --help              # Show help
ipcrawler 127.0.0.1           # Test scan  
ipcrawler target.com          # Scan target
ipcrawler -v target.com       # Verbose scan
ls /scans                     # View results
exit                          # Leave container
```

**Results automatically saved to:** `your-folder\results\`

### 🐧 Linux/macOS Users

#### **⚡ One-liner (if you have make):**
```bash
git clone https://github.com/hckerhub/ipcrawler.git && cd ipcrawler && make setup
```

#### **📦 Don't have make? Auto-install it:**
```bash
git clone https://github.com/hckerhub/ipcrawler.git && cd ipcrawler && ./bootstrap.sh && make setup
```

### 🐧 Linux/macOS Users

**📋 First, ensure you have the [prerequisites](#-prerequisites) installed (Python 3.8+ and make)**

#### **⚡ One-liner (if you have Python and make):**
```bash
git clone https://github.com/hckerhub/ipcrawler.git && cd ipcrawler && make setup
```

#### **📦 Don't have make? Auto-install it:**
```bash
git clone https://github.com/hckerhub/ipcrawler.git && cd ipcrawler && ./bootstrap.sh && make setup
```

#### **🔍 Quick prerequisite check:**
```bash
python3 --version  # Should show 3.8 or higher
make --version     # Should show GNU Make
```

### 🔧 Missing Make Commands (Docker Users)

If you're using Docker only (Windows Path 1), you won't have access to these convenience commands:

| Command | Description |
|---------|-------------|
| `make setup` | Auto-install all security tools (nmap, masscan, nikto, etc.) |
| `make setup-docker` | Auto-install Docker and Docker Compose |
| `make clean` | Remove all installed tools and virtual environment |
| `make update` | Update ipcrawler, tools, and rebuild Docker image |
| `make test` | Run unit tests and linting |
| `make lint` | Check code formatting and style |
| `make format` | Auto-format Python code |

**💡 Docker alternative:** You can update the Docker image with:
```bash
# Update ipcrawler
git pull

# Rebuild Docker image with latest changes
docker build -t ipcrawler .
```

### 🔧 Advanced Make Commands (For Developers)

**Windows users:** You have the `ipcrawler-windows.bat` launcher - no make commands needed!

**Linux/macOS users** get these additional convenience commands:

| Command | Description |
|---------|-------------|
| `make setup` | Auto-install all security tools (nmap, masscan, nikto, etc.) |
| `make setup-docker` | Auto-install Docker and Docker Compose |
| `make clean` | Remove all installed tools and virtual environment |
| `make update` | Update ipcrawler, tools, and rebuild Docker image |
| `make test` | Run unit tests and linting |
| `make lint` | Check code formatting and style |
| `make format` | Auto-format Python code |

**💡 Windows alternatives:**
```cmd
# Update ipcrawler and rebuild
git pull
docker build -t ipcrawler .

# View Docker images
docker images

# Remove old images
docker rmi ipcrawler
```

### 🐳 Docker Setup (Recommended)

**🎯 Bypasses all system requirements - just needs Docker!**

No Python, make, or security tools needed on your host system!

```bash
# Clone and setup
git clone https://github.com/hckerhub/ipcrawler.git
cd ipcrawler

# Build Docker image (one time)
make setup-docker

# Start interactive container
make docker-cmd

# Inside container - scan away!
ipcrawler 10.10.10.1
ipcrawler -v target.com
exit  # Leave container when done
```

**Alternative Docker setup** (without make):
```bash
# Manual Docker setup
docker build -t ipcrawler .
docker run -it --rm -v $(pwd)/results:/scans ipcrawler
```

### 🖥️ Local Installation

#### Prerequisites
```bash
# Update package cache
sudo apt update

# Install required tools (Kali Linux recommended)
sudo apt install seclists curl dnsrecon enum4linux feroxbuster gobuster impacket-scripts nbtscan nikto nmap onesixtyone oscanner redis-tools smbclient smbmap snmp sslscan sipvicious tnscmd10g whatweb
```

#### Installation
```bash
# Clone and setup (one command!)
git clone https://github.com/hckerhub/ipcrawler.git
cd ipcrawler
make setup
```

#### Usage
```bash
# Scan single target
ipcrawler 10.10.10.1

# Scan multiple targets
ipcrawler 10.10.10.1 10.10.10.2 192.168.1.0/24

# Scan from file
ipcrawler -t targets.txt

# Custom ports
ipcrawler -p 80,443,8080 10.10.10.1
```

#### Cleanup
```bash
# Remove everything (local setup, Docker resources, empty directories)
make clean
```

## 🔥 Key Features

- **🎯 Smart Enumeration**: Automatically launches appropriate tools based on discovered services
- **⚡ Multi-threading**: Scan multiple targets concurrently
- **📁 Organized Output**: Clean directory structure for results
- **🔧 Highly Configurable**: Customizable via config files and command-line options
- **🏷️ Plugin System**: Extensive plugin ecosystem for different services
- **⏱️ Time Management**: Global and per-target timeouts
- **🎨 Clean Interface**: Color-coded output with multiple verbosity levels

## 📊 Example Output Structure

```
results/
└── 10.10.10.1/
    ├── exploit/          # Exploit code and payloads
    ├── loot/            # Credentials, hashes, files
    ├── report/          # Flags, notes, screenshots
    └── scans/           # All scan results
        ├── _commands.log         # Commands executed
        ├── _manual_commands.txt  # Suggested manual commands
        ├── tcp80/               # HTTP enumeration
        ├── tcp22/               # SSH enumeration
        └── xml/                 # Raw XML outputs
```

## 🛠️ Common Usage Examples

```bash
# Basic scan with verbose output
ipcrawler -v 10.10.10.1

# Fast scan (top 1000 ports only)
ipcrawler -p 1-1000 10.10.10.1

# Scan specific services
ipcrawler --force-services tcp/80/http tcp/443/https 10.10.10.1

# Exclude certain plugins
ipcrawler --exclude-tags bruteforce 10.10.10.1

# Time-limited scan (60 minutes max)
ipcrawler --timeout 60 10.10.10.1
```



### Results Location
- **Host machine**: `./results/` (persistent after container exits)
- **Inside container**: `/scans/` (mounted volume)

### Cleanup
The `make clean` command provides intelligent cleanup for both local and Docker installations:

```bash
make clean
```

**What it removes:**
- ✅ Virtual environment and local installation
- ✅ Docker images and containers (if they exist)
- ✅ Empty results directories
- ✅ Global command installation

**What it preserves:**
- 🛡️ Results directories containing scan data
- 🛡️ Non-ipcrawler Docker resources

## ⚙️ Configuration

ipcrawler uses the same configuration system as AutoRecon. Config files are located at:
- `~/.config/ipcrawler/config.toml` - Main configuration
- `~/.config/ipcrawler/global.toml` - Global settings

## 🎓 Perfect for OSCP & CTFs

ipcrawler excels in time-constrained environments:
- **OSCP Exam**: Run against all targets while focusing on one
- **HTB/VulnHub**: Quick initial enumeration 
- **CTF Events**: Rapid service discovery and enumeration

## 💡 Pro Tips

1. **Start Early**: Launch ipcrawler on all targets at the beginning
2. **Use Verbosity**: `-v` shows discovered services in real-time
3. **Check Manual Commands**: Review `_manual_commands.txt` for additional tests
4. **Organized Results**: The directory structure keeps everything organized
5. **Multiple Sessions**: Run different scan types in parallel
6. **Easy Cleanup**: Use `make clean` for complete removal when done
7. **Safe Results**: Cleanup preserves scan data in results directories

## 🔍 Verbosity Levels

| Flag | Output Level |
|------|-------------|
| (none) | Minimal - start/end announcements |
| `-v` | Verbose - plugin starts, open ports, services |
| `-vv` | Very verbose - commands executed, pattern matches |
| `-vvv` | Maximum - live output from all commands |

## 🏆 What Users Say About AutoRecon (ipcrawler's foundation)

> *"AutoRecon was invaluable during my OSCP exam... I would strongly recommend this utility for anyone in the PWK labs, the OSCP exam, or other environments such as VulnHub or HTB."*
> 
> **- b0ats** (rooted 5/5 exam hosts)

> *"The strongest feature of AutoRecon is the speed; on the OSCP exam I left the tool running in the background while I started with another target, and in a matter of minutes I had all of the output waiting for me."*
> 
> **- tr3mb0** (rooted 4/5 exam hosts)

> *"Being introduced to AutoRecon was a complete game changer for me while taking the OSCP... After running AutoRecon on my OSCP exam hosts, I was given a treasure chest full of information that helped me to start on each host and pass on my first try."*
> 
> **- rufy** (rooted 4/5 exam hosts)

*ipcrawler provides the same powerful enumeration with easier setup - just `make setup` and you're ready to go!*

## 📋 Requirements

### Docker Setup (Recommended)
- **Docker Desktop** or **Docker Engine**
- **Any operating system** (Windows, macOS, Linux)

### Local Setup
- **Python 3.8+**
- **Linux/Unix environment** (Kali Linux recommended)
- **Network enumeration tools** (listed in prerequisites)
- **SecLists wordlists** (`sudo apt install seclists`)

## 🤝 Contributing

This project maintains compatibility with AutoRecon plugins and configurations. For core functionality improvements, consider contributing to the original [AutoRecon project](https://github.com/Tib3rius/AutoRecon).

## ⚠️ Disclaimer

ipcrawler performs **no automated exploitation** by default, keeping it OSCP exam compliant. The tool is for authorized testing only. Users are responsible for compliance with applicable laws and regulations.

---

**⭐ Star this repo if ipcrawler helps you ace your OSCP exam or CTF challenges!**

Made with ❤️ based on [AutoRecon](https://github.com/Tib3rius/AutoRecon) by [Tib3rius](https://github.com/Tib3rius)
