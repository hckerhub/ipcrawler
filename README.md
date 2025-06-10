# 🕷️ ipcrawler

> *"It's like bowling with bumpers."* - [@ippsec](https://twitter.com/ippsec)

A simplified, streamlined version of **AutoRecon** - the multi-threaded network reconnaissance tool that performs automated enumeration of services for CTFs, OSCP, and penetration testing environments.

## 🙏 Credits

**ipcrawler** is a fork of [**AutoRecon**](https://github.com/Tib3rius/AutoRecon) by [**Tib3rius**](https://github.com/Tib3rius). All core functionality, plugins, and the brilliant multi-threaded architecture are thanks to his incredible work. This fork simply provides a cleaner setup experience while maintaining all the powerful features of the original tool.

## ✨ What's New

**ipcrawler** takes AutoRecon's powerful enumeration capabilities and makes setup effortless:

| **Before (AutoRecon)** | **After (ipcrawler)** |
|---|---|
| `pipx install git+https://github.com/Tib3rius/AutoRecon.git` | `make setup` |
| `sudo env "PATH=$PATH" autorecon target` | `ipcrawler target` |
| Complex dependency management | Automatic virtual environment |
| Manual uninstallation | `make clean` |

## 🚀 Quick Start

### Prerequisites
```bash
# Update package cache
sudo apt update

# Install required tools (Kali Linux recommended)
sudo apt install seclists curl dnsrecon enum4linux feroxbuster gobuster impacket-scripts nbtscan nikto nmap onesixtyone oscanner redis-tools smbclient smbmap snmp sslscan sipvicious tnscmd10g whatweb
```

### Installation
```bash
# Clone and setup (one command!)
git clone https://github.com/yourusername/ipcrawler.git
cd ipcrawler
make setup
```

### Usage
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

### Cleanup
```bash
# Remove everything (including global command)
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

## 🔍 Verbosity Levels

| Flag | Output Level |
|------|-------------|
| (none) | Minimal - start/end announcements |
| `-v` | Verbose - plugin starts, open ports, services |
| `-vv` | Very verbose - commands executed, pattern matches |
| `-vvv` | Maximum - live output from all commands |

## 🏆 Testimonials

> *"ipcrawler was invaluable during my OSCP exam... I would strongly recommend this utility for anyone in the PWK labs, the OSCP exam, or other environments such as VulnHub or HTB."*
> 
> **- b0ats** (rooted 5/5 exam hosts)

> *"The strongest feature of ipcrawler is the speed; on the OSCP exam I left the tool running in the background while I started with another target, and in a matter of minutes I had all of the output waiting for me."*
> 
> **- tr3mb0** (rooted 4/5 exam hosts)

> *"Being introduced to ipcrawler was a complete game changer for me while taking the OSCP... After running ipcrawler on my OSCP exam hosts, I was given a treasure chest full of information that helped me to start on each host and pass on my first try."*
> 
> **- rufy** (rooted 4/5 exam hosts)

*[More testimonials in the original AutoRecon repository]*

## 📋 Requirements

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
