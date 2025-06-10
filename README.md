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

### 🪟 Windows
- [Docker Desktop](https://www.docker.com/products/docker-desktop)

### 🐧 Linux/macOS  
- **Python 3.8+**: [python.org](https://www.python.org/downloads/) or package manager
  ```bash
  # Ubuntu/Debian: sudo apt install python3 python3-pip python3-venv
  # CentOS/RHEL: sudo yum install python3 python3-pip  
  # Arch: sudo pacman -S python python-pip
  # macOS: brew install python3
  ```
- **make**: Usually pre-installed or `sudo apt install make`

## 🚀 Quick Start

### 🪟 Windows Users
1. **Install Docker Desktop** - [Download here](https://www.docker.com/products/docker-desktop)
2. **Run the launcher:**
```cmd
git clone https://github.com/hckerhub/ipcrawler.git
cd ipcrawler
ipcrawler-windows.bat
```

The launcher automatically builds Docker image and opens an interactive terminal with all tools ready.

### 🐧 Linux/macOS Users

**Prerequisites:** Python 3.8+, make ([see installation](#-prerequisites))

```bash
# One-liner setup
git clone https://github.com/hckerhub/ipcrawler.git && cd ipcrawler && make setup

# Or if you don't have make installed
./bootstrap.sh && make setup
```

### 🐳 Docker Setup (All Platforms)

**Bypasses all system requirements - just needs Docker!**

```bash
git clone https://github.com/hckerhub/ipcrawler.git && cd ipcrawler
make setup-docker && make docker-cmd

# Or manual Docker setup
docker build -t ipcrawler . && docker run -it --rm -v $(pwd)/results:/scans ipcrawler
```

### 🔧 Make Commands (Linux/macOS Only)

| Command | Description |
|---------|-------------|
| `make setup` | Install all tools and create virtual environment |
| `make setup-docker` | Build Docker image |
| `make docker-cmd` | Start interactive Docker container |
| `make clean` | Remove everything (preserves scan results) |
| `make update` | Update tools and Docker image |

**Windows users:** Use `ipcrawler-windows.bat` - no make commands needed!

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

## 💡 Quick Tips

- **Start early** - Launch on all targets while focusing on one
- **Use `-v`** - Shows discovered services in real-time  
- **Check `_manual_commands.txt`** - Additional tests to run
- **Results in `./results/`** - Preserved after cleanup

## 🔍 Verbosity: None | `-v` | `-vv` | `-vvv` (minimal to maximum output)

## 🏆 OSCP Success Stories

*"AutoRecon was invaluable during my OSCP exam... I would strongly recommend this utility."* **- b0ats** (5/5 hosts)

*"The strongest feature is the speed... in minutes I had all output waiting for me."* **- tr3mb0** (4/5 hosts)

*ipcrawler provides the same power with easier setup!*

## 🤝 Contributing

This project maintains compatibility with AutoRecon plugins and configurations. For core functionality improvements, consider contributing to the original [AutoRecon project](https://github.com/Tib3rius/AutoRecon).

## ⚠️ Disclaimer

ipcrawler performs **no automated exploitation** by default, keeping it OSCP exam compliant. The tool is for authorized testing only. Users are responsible for compliance with applicable laws and regulations.

---

**⭐ Star this repo if ipcrawler helps you ace your OSCP exam or CTF challenges!**

Made with ❤️ based on [AutoRecon](https://github.com/Tib3rius/AutoRecon) by [Tib3rius](https://github.com/Tib3rius)
