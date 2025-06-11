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
- [Docker Desktop](https://www.docker.com/products/docker-desktop) - **Manual installation required**

### 🐧 Linux/macOS  
- **Python 3.8+**: [python.org](https://www.python.org/downloads/) or package manager
  ```bash
  # Ubuntu/Debian: sudo apt install python3 python3-pip python3-venv
  # CentOS/RHEL: sudo yum install python3 python3-pip  
  # Arch: sudo pacman -S python python-pip
  # macOS: brew install python3
  ```
- **make**: Usually pre-installed or `sudo apt install make` or `brew install make`
- **seclists**: `sudo apt install seclists` (Required for macOS only)

## 🚀 Quick Start Guide

### 📖 Step-by-Step Setup Instructions

#### 🪟 **For Windows Users (Docker Recommended)**

**Step 1: Install Docker Desktop**
1. Download [Docker Desktop for Windows](https://www.docker.com/products/docker-desktop)
2. Install and start Docker Desktop
3. Verify Docker is running: open PowerShell/CMD and run `docker --version`

**Step 2: Get ipcrawler**
```cmd
git clone https://github.com/hckerhub/ipcrawler.git
cd ipcrawler
```

**Step 3: Run ipcrawler**
```cmd
# Double-click this file OR run from command line:
ipcrawler-windows.bat
```

**What the `.bat` file does:**
- ✅ Checks if Docker is installed and running
- 🔨 Builds the ipcrawler Docker image automatically (first time only)
- 🚀 Opens an interactive terminal with all security tools pre-installed
- 💾 Maps your `results/` folder so scan results persist after closing

**Inside the Docker container, you can run:**
```bash
ipcrawler --help           # Show beautiful help with examples
ipcrawler 127.0.0.1        # Test scan localhost
ipcrawler 10.10.10.1       # Scan your target
exit                       # Close container
```

#### 🐧 **For Linux/macOS Users (Local or Docker)**

**Option 1: Local Installation (Recommended for Linux)**
```bash
# One-liner setup
git clone https://github.com/hckerhub/ipcrawler.git && cd ipcrawler && make setup

# If 'make' is not installed, run bootstrap first:
./bootstrap.sh && make setup
```

**What `bootstrap.sh` does:**
- 🔍 Detects your operating system (Ubuntu, Kali, macOS, Arch, etc.)
- 📦 Installs `make` automatically using your system's package manager
- 🎯 Provides specific instructions for your OS if automatic installation fails

**What `make setup` does:**
- 🐍 Creates a Python virtual environment in `./venv/`
- 📦 Installs all Python dependencies (Rich, colorama, etc.)
- 🔧 Downloads and installs security tools (nmap, nikto, dirb, etc.)
- 🔗 Creates global `ipcrawler` command that you can run from anywhere
- ✅ Verifies everything is working

**Option 2: Docker Setup (Works on all platforms)**
```bash
# Prerequisites: Install Docker first for your system:
# Ubuntu/Debian: sudo apt install docker.io
# macOS: brew install --cask docker (or Docker Desktop)
# Then:

git clone https://github.com/hckerhub/ipcrawler.git && cd ipcrawler
make setup-docker
```

**What `make setup-docker` does:**
- ✅ Checks if Docker is installed and running
- 🔨 Builds the ipcrawler Docker image with all tools
- 🚀 Opens an interactive Docker terminal
- 💾 Maps your `results/` folder for persistent storage

### 🔧 Understanding the Tools

#### `ipcrawler-cmd` (Local Installation Only)
- **Purpose**: Global command wrapper created during `make setup`
- **Location**: Installed to `/usr/local/bin/ipcrawler` (or similar system PATH)
- **What it does**: Activates the virtual environment and runs ipcrawler
- **Usage**: Allows you to run `ipcrawler target` from any directory
- **Do you need it?** Yes for local installation, not needed for Docker

#### Make Commands Reference

| Command | When to Use | What It Does |
|---------|------------|--------------|
| `./bootstrap.sh` | First time, if `make` is missing | Installs `make` for your OS |
| `make setup` | Local installation | Full local setup with tools |
| `make setup-docker` | Docker installation | Docker build + interactive session |
| `make docker-cmd` | Docker users | Start additional Docker sessions |
| `make clean` | When uninstalling | Removes everything (preserves results) |
| `make update` | Periodic updates | Updates tools and Docker image |

**Windows users:** Only use `ipcrawler-windows.bat` - no make commands needed!

## ⚙️ Configuration Files

ipcrawler uses TOML configuration files for customization. After running setup, config files are created:

### 📍 Configuration Locations

**Local Installation:**
- Main config: `~/.config/ipcrawler/config.toml`
- Global settings: `~/.config/ipcrawler/global.toml`

**Docker Installation:**
- Config files are inside the container at `/opt/ipcrawler/ipcrawler/`
- Changes persist only during the session

### 📝 `config.toml` - Main Configuration

Controls ipcrawler behavior and plugin settings:

```toml
# Basic ipcrawler options
# nmap-append = '-T3'          # Add flags to nmap commands
# verbose = 1                  # Default verbosity level
# max-scans = 30              # Maximum concurrent scans

# Global options
# [global]
# username-wordlist = '/usr/share/seclists/Usernames/cirt-default-usernames.txt'

# Plugin-specific options
# [dirbuster]
# threads = 50
# wordlist = [
#     '/usr/share/seclists/Discovery/Web-Content/common.txt',
#     '/usr/share/seclists/Discovery/Web-Content/big.txt'
# ]
```

**Common configurations:**
```toml
# For faster scans
nmap-append = '-T4'
max-scans = 20

# For OSCP exam (conservative)
nmap-append = '-T3'
max-scans = 10

# Verbose by default
verbose = 1
```

### 🌐 `global.toml` - Global Settings

Defines global options and pattern matching:

```toml
# Default wordlists
[global.username-wordlist]
default = '/usr/share/seclists/Usernames/top-usernames-shortlist.txt'

[global.password-wordlist]  
default = '/usr/share/seclists/Passwords/darkweb2017-top100.txt'

[global.domain]
help = 'Domain for DNS/AD enumeration'

# Pattern matching for interesting findings
[[pattern]]
description = 'CVE Identified: {match}'
pattern = '(CVE-\d{4}-\d{4,7})'

[[pattern]]
description = 'Potential vulnerability: {match}'
pattern = 'State: (?:(?:LIKELY\_?)?VULNERABLE)'
```

**How to customize:**
1. **Uncomment lines** by removing the `#` symbol
2. **Edit values** to match your preferences
3. **Add new sections** for additional plugins
4. **Save and restart** ipcrawler for changes to take effect

## 🔥 Key Features

- **🎯 Smart Enumeration**: Automatically launches appropriate tools based on discovered services
- **⚡ Multi-threading**: Scan multiple targets concurrently
- **📁 Organized Output**: Clean directory structure for results
- **📋 Rich Summary Reports**: Beautiful HTML reports with collapsible sections
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
    │   ├── Full_Report.html     # 🌟 Rich HTML summary
    │   ├── local.txt            # Local flag
    │   ├── proof.txt            # Proof flag
    │   └── screenshots/         # Screenshots
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

# Multiple targets
ipcrawler 10.10.10.1 10.10.10.2 10.10.10.3

# Custom nmap timing
ipcrawler --nmap-append '-T3' 10.10.10.1
```

### Results Location
- **Local installation**: `./results/` in the directory where you ran ipcrawler
- **Docker installation**: `./results/` on your host machine (mapped from container)

### Cleanup
```bash
make clean
```

**What `make clean` removes:**
- ✅ Virtual environment and local installation
- ✅ Docker images and containers (if they exist)
- ✅ Empty results directories
- ✅ Global command installation

**What it preserves:**
- 🛡️ Results directories containing scan data
- 🛡️ Configuration files in `~/.config/ipcrawler/`
- 🛡️ Non-ipcrawler Docker resources

## 🎓 Perfect for OSCP & CTFs

ipcrawler excels in time-constrained environments:
- **OSCP Exam**: Run against all targets while focusing on one
- **HTB/VulnHub**: Quick initial enumeration 
- **CTF Events**: Rapid service discovery and enumeration

## 💡 Pro Tips

1. **Beautiful Help**: Run `ipcrawler --help` to see the enhanced help with examples and organized sections
2. **Rich Summary Report**: Check `Full_Report.html` in each target's report directory for a comprehensive HTML summary
3. **Start Early**: Launch ipcrawler on all targets at the beginning
4. **Use Rich Verbosity**: `-v` shows visual progress, `-vv` shows timing, `-vvv` shows live output
5. **Check Manual Commands**: Review `_manual_commands.txt` for additional tests
6. **Organized Results**: The directory structure keeps everything organized
7. **Multiple Sessions**: Run different scan types in parallel
8. **Easy Cleanup**: Use `make clean` for complete removal when done
9. **Safe Results**: Cleanup preserves scan data in results directories

## 🔍 Enhanced Verbosity Levels

| Flag | Output Level | Rich Enhancements |
|------|-------------|-------------------|
| (none) | Minimal - start/end announcements | Standard output |
| [`-v`] | Plugin starts, discoveries | 🔍 **Visual icons**, colored progress indicators |
| [`-vv`] | Commands, timing, patterns | ✅ **Completion status**, timing info, pattern highlights |
| [`-vvv`] | Live command output | │ **Subtle formatting** to avoid overwhelming |

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
- **🔤 Visual Flag Formatting**: [`-v`] [`--verbose`] style makes options easy to scan
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
- **Use [`-v`]** - Shows visual progress with Rich icons and colors
- **Check `_manual_commands.txt`** - Additional tests to run
- **Results in `./results/`** - Preserved after cleanup

## 🔍 Verbosity: None | [`-v`] | [`-vv`] | [`-vvv`] (minimal to maximum output)

## 🤝 Contributing

This project maintains compatibility with AutoRecon plugins and configurations. For core functionality improvements, consider contributing to the original [AutoRecon project](https://github.com/Tib3rius/AutoRecon).

## ⚠️ Disclaimer

ipcrawler performs **no automated exploitation** by default, keeping it OSCP exam compliant. The tool is for authorized testing only. Users are responsible for compliance with applicable laws and regulations.

---

**⭐ Star this repo if ipcrawler helps you ace your OSCP exam or CTF challenges!**

Made with ❤️ based on [AutoRecon](https://github.com/Tib3rius/AutoRecon) by [Tib3rius](https://github.com/Tib3rius)
