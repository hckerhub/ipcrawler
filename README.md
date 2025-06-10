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

## 🚀 Quick Start

### 🔧 Don't have `make` installed?

**No problem!** Our bootstrap script automatically installs `make` on any system:

```bash
# Clone the repository
git clone https://github.com/hckerhub/ipcrawler.git
cd ipcrawler

# Run the bootstrap script (installs make automatically)
./bootstrap.sh

# Then use standard make commands
make setup      # or make setup-docker
```

**What the bootstrap script does:**
- ✅ **Detects your operating system** automatically
- ✅ **Installs make** using the appropriate package manager
- ✅ **Supports all major platforms**: Linux, macOS, Windows (WSL)
- ✅ **One command solution** - no manual steps needed

**Supported systems:**
- **Linux**: Kali, Ubuntu, Debian, Arch, CentOS, RHEL, Fedora, openSUSE, Alpine
- **macOS**: Via Homebrew or Xcode Command Line Tools
- **Windows**: WSL, Chocolatey, or Scoop

### 🪟 Windows Users

**For Windows, the bootstrap script recommends:**

1. **WSL (Windows Subsystem for Linux)** - Best option
   ```cmd
   wsl --install
   # Then clone and run bootstrap inside WSL
   ```

2. **Package managers** - If you prefer native Windows
   ```cmd
   # Via Chocolatey
   choco install make
   
   # Via Scoop  
   scoop install make
   ```

3. **Docker Desktop** - Cross-platform solution
   ```cmd
   # Install Docker Desktop, then use make setup-docker
   ```

### 🐳 Docker Setup (Recommended)
No dependencies needed - everything runs in a container!

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

## 🐳 Docker Details

### Available Make Commands
```bash
make help          # Show all available commands
make setup-docker  # Build the Docker image (one time)
make docker-cmd    # Start interactive container
make clean         # Complete cleanup (local + Docker resources)
```

### Docker Features
- ✅ **No Dependencies**: Works on any system with Docker
- ✅ **Isolated Environment**: No impact on host system
- ✅ **Persistent Results**: Scans saved to `results/` directory
- ✅ **Pre-installed Tools**: Includes nmap and essential tools
- ✅ **Expandable**: Run `/install-tools.sh` for additional tools
- ✅ **Smart Cleanup**: `make clean` removes everything safely

### Complete Workflow
```bash
# 1. Setup (one time)
make setup-docker

# 2. Scan targets
make docker-cmd
# Inside container: ipcrawler target.com

# 3. Check results (on host)
ls results/

# 4. Complete cleanup (when done)
make clean
```

### Manual Docker Commands
```bash
# Build image
docker build -t ipcrawler .

# Run interactively with results mounted
docker run -it --rm -v $(pwd)/results:/scans ipcrawler

# For local network scanning
docker run -it --rm --network host -v $(pwd)/results:/scans ipcrawler
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
