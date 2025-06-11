# 🕷️ ipcrawler

> *"It's like bowling with bumpers."* - [@ippsec](https://twitter.com/ippsec)

A simplified, streamlined version of **AutoRecon** - the multi-threaded network reconnaissance tool that performs automated enumeration of services for CTFs, OSCP, and penetration testing environments.

## 🙏 Credits

**ipcrawler** is a fork of [**AutoRecon**](https://github.com/Tib3rius/AutoRecon) by [**Tib3rius**](https://github.com/Tib3rius). All core functionality, plugins, and the brilliant multi-threaded architecture are thanks to his incredible work. This fork simply provides a cleaner setup experience while maintaining all the powerful features of the original tool.

## ✨ What's New

**ipcrawler** takes AutoRecon's powerful enumeration capabilities and makes setup effortless with modern enhancements:

### 🎨 Enhanced User Experience
- **🎯 Beginner-Friendly Help**: Beautiful Rich-formatted `--help` with organized sections, examples, and pro tips
- **📊 Enhanced Verbosity**: Rich-colored output for `-v`, `-vv`, `-vvv` levels with visual progress indicators
- **🔧 Smart Setup**: Automatic Docker detection with manual installation workflow
- **📋 Clear Documentation**: Step-by-step guides with troubleshooting

### 🤔 Docker vs Local Setup

| Feature | Docker | Local |
|---------|--------|-------|
| **Setup Time** | 3-5 minutes | 5-10 minutes |
| **Dependencies** | Manual Docker install | Manual tool installation |
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
- **make**: Usually pre-installed or `sudo apt install make` or `brew install make`

## 🚀 Quick Start

### 🪟 Windows Users
1. **Install Docker Desktop** - [Download here](https://www.docker.com/products/docker-desktop)
2. **Run the launcher:**
```cmd
git clone https://github.com/hckerhub/ipcrawler.git
cd ipcrawler
ipcrawler-windows.bat
```

The launcher checks for Docker, builds the image (if needed), and opens an interactive terminal with all tools ready.

### 🐧 Linux/macOS Users

**Prerequisites:** Python 3.8+, make ([see installation](#-prerequisites))

```bash
# One-liner setup
git clone https://github.com/hckerhub/ipcrawler.git && cd ipcrawler && make setup

# Or if you don't have make installed
./bootstrap.sh && make setup
```

### 🐳 Docker Setup (All Platforms)

**Requires manual Docker installation, then simple setup!**

```bash
# 1. Install Docker first (platform-specific):
# Windows: Docker Desktop from docker.com
# macOS: brew install --cask docker or Docker Desktop
# Linux: sudo apt install docker.io (or equivalent)

# 2. Clone and setup ipcrawler
git clone https://github.com/hckerhub/ipcrawler.git && cd ipcrawler
make setup-docker

# 3. Start additional sessions
make docker-cmd
```

### 🔧 Make Commands (Linux/macOS Only)

| Command | Description |
|---------|-------------|
| `make setup` | Install all tools and create virtual environment |
| `make setup-docker` | Check Docker + build image + open interactive terminal |
| `make docker-cmd` | Start additional interactive Docker sessions |
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
- **🎨 Rich Interface**: Beautiful help system and enhanced verbosity output
- **👤 Beginner-Friendly**: Organized help with examples, pro tips, and visual feedback

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

1. **Beautiful Help**: Run `ipcrawler --help` to see the enhanced help with examples and organized sections
2. **Start Early**: Launch ipcrawler on all targets at the beginning
3. **Use Rich Verbosity**: `-v` shows visual progress, `-vv` shows timing, `-vvv` shows live output
4. **Check Manual Commands**: Review `_manual_commands.txt` for additional tests
5. **Organized Results**: The directory structure keeps everything organized
6. **Multiple Sessions**: Run different scan types in parallel
7. **Easy Cleanup**: Use `make clean` for complete removal when done
8. **Safe Results**: Cleanup preserves scan data in results directories

## 🔍 Enhanced Verbosity Levels

| Flag | Output Level | Rich Enhancements |
|------|-------------|-------------------|
| (none) | Minimal - start/end announcements | Standard output |
| `-v` | Plugin starts, discoveries | 🔍 **Visual icons**, colored progress indicators |
| `-vv` | Commands, timing, patterns | ✅ **Completion status**, timing info, pattern highlights |
| `-vvv` | Live command output | │ **Subtle formatting** to avoid overwhelming |

**Example Rich Output:**
```
[🔍 PORT] TCP Fast Scan (nmap-fast-tcp) → 10.10.10.1
🎯 DISCOVERED tcp/22 on 10.10.10.1
✅ COMPLETED TCP Fast Scan (nmap-fast-tcp) on 10.10.10.1 in 2 minutes, 15 seconds
```

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

## 🎨 Rich User Interface

**ipcrawler** now features a beautiful, beginner-friendly interface powered by the Rich library:

### Enhanced Help System
- **🎯 Organized Sections**: Essential Options, Advanced Options, Port Syntax Examples
- **📋 Clear Examples**: Real command examples for every feature
- **🔤 Visual Flag Formatting**: `[-v] [--verbose]` style makes options easy to scan
- **💡 Pro Tips**: OSCP-specific advice and best practices
- **🎨 Consistent Styling**: Professional color scheme throughout

### Enhanced Verbosity Output
- **Level 1 (-v)**: Visual icons and progress indicators for plugin starts and discoveries
- **Level 2 (-vv)**: Completion status with timing and enhanced pattern matches
- **Level 3 (-vvv)**: Subtle formatting for live output without overwhelming

### Automatic Fallback
If Rich is not available, ipcrawler seamlessly falls back to standard colorama output with no loss of functionality.

## 📋 Requirements

### Docker Setup (Recommended)
- **Docker Desktop** or **Docker Engine** (manual installation required)
- **Any operating system** (Windows, macOS, Linux)
- Scripts check for Docker and provide installation guidance if missing

### Local Setup
- **Python 3.8+**
- **Rich library** (auto-installed for enhanced help and verbosity)
- **Linux/Unix environment** (Kali Linux recommended)
- **Network enumeration tools** (auto-installed via scripts)
- **SecLists wordlists** (`sudo apt install seclists`)

## 💡 Quick Tips

- **Try `ipcrawler --help`** - Beautiful help with examples and organized sections
- **Start early** - Launch on all targets while focusing on one
- **Use `-v`** - Shows visual progress with Rich icons and colors
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
