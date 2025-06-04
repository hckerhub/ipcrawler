# 🎯 IPCrawler

**Advanced IP Scanner & Vulnerability Hunter for Hack The Box**

IPCrawler is a powerful, interactive terminal-based IP scanning and vulnerability assessment tool designed specifically for Hack The Box machines. Built with Go, it combines the power of NMAP with intelligent vulnerability analysis and a beautiful cyberpunk-themed terminal user interface.

```
 ██╗██████╗  ██████╗██████╗  █████╗ ██╗    ██╗██╗     ███████╗██████╗ 
 ██║██╔══██╗██╔════╝██╔══██╗██╔══██╗██║    ██║██║     ██╔════╝██╔══██╗
 ██║██████╔╝██║     ██████╔╝███████║██║ █╗ ██║██║     █████╗  ██████╔╝
 ██║██╔═══╝ ██║     ██╔══██╗██╔══██║██║███╗██║██║     ██╔══╝  ██╔══██╗
 ██║██║     ╚██████╗██║  ██║██║  ██║╚███╔███╔╝███████╗███████╗██║  ██║
 ╚═╝╚═╝      ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝ ╚══╝╚══╝ ╚══════╝╚══════╝╚═╝  ╚═╝

     🎯 Advanced IP Scanner & Vulnerability Hunter 🎯
          Hack The Box Edition - Built with ❤️ in Go
```

## ✨ Features

- **🔍 Advanced Scanning**: Full TCP/UDP port scanning with service detection
- **🎨 Beautiful TUI**: Interactive cyberpunk-themed terminal interface
- **🚨 Smart Vulnerability Detection**: HTB-focused vulnerability analysis with severity ratings
- **📊 Multiple Output Formats**: Interactive TUI, JSON, and table formats
- **🔐 Smart Privilege Management**: Automatic sudo handling for UDP scans

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- NMAP installed
- Linux/macOS/Windows terminal

### Installation

```bash
# Clone and setup
git clone https://github.com/hckerhub/ipcrawler.git
cd ipcrawler

# One command does everything - build, install, ready to use!
make crawler

# Quick start options
ipcrawler                  # System command (after setup)
./crawler                  # Smart launcher script  
./crawler 10.10.10.1      # Start with target IP
./crawler --aggressive    # Start with aggressive mode
```

## 📖 Usage

### Interactive TUI Mode
```bash
# Launch TUI
./ipcrawler tui

# Start with target
./ipcrawler tui 10.10.10.1

# Aggressive scan
./ipcrawler tui 192.168.1.100 --aggressive
```

### Command Line Mode
```bash
# Basic scan
./ipcrawler scan 10.10.10.1

# Aggressive scan with output
./ipcrawler scan 10.10.10.1 --aggressive --output results.json
```

### Quick Commands
```bash
make crawler              # Complete setup (build + install)
./crawler                 # Smart launcher (builds if needed)
ipcrawler                 # System command (after setup)
make start                # Quick TUI start
```

## 🎮 Controls

**Input Screen:**
- **Enter**: Start scan
- **Tab**: Toggle aggressive mode
- **Ctrl+V**: Toggle verbose
- **Ctrl+C**: Quit

**Results Screen:**
- **← →**: Navigate tabs
- **R**: New scan
- **S**: Save results
- **Q**: Quit

## 🎯 HTB-Optimized Features

### High-Value Ports
- **21** (FTP) - Anonymous access detection
- **22** (SSH) - Weak credential analysis
- **80/443** (HTTP/HTTPS) - Web vulnerability scanning
- **445** (SMB) - Share enumeration and exploits
- **1433** (MSSQL) - Default credentials and xp_cmdshell

### Smart Analysis
- Service-specific vulnerability detection
- Severity classification (Critical, High, Medium, Low)
- Actionable next-step suggestions
- HTB-focused attack vectors

## 🔐 Privilege Management

**TCP Scanning** (No sudo required):
- ✅ Full port range (1-65535)
- ✅ Service detection

**UDP Scanning** (Requires sudo):
- 🔐 Interactive privilege prompts
- ❌ Graceful degradation if declined
- ✅ Enhanced detection with privileges

## 📊 Output Formats

- **Interactive TUI**: Real-time results with navigation
- **JSON**: Machine-readable structured data
- **Table**: Human-readable tabular reports
- **Summary**: Quick overview with key findings

## 🔧 Configuration

### Environment Variables
```bash
export IPCRAWLER_AGGRESSIVE=true
export IPCRAWLER_VERBOSE=true
export IPCRAWLER_OUTPUT_DIR=./scans
```

### Config File (`~/.ipcrawler/config.yaml`)
```yaml
default:
  aggressive: false
  verbose: false
  timeout: 300

scan:
  tcp_ports: "1-65535"
  udp_ports: "53,67,123,135,161,445,1434"
  threads: 50
```

## 🤝 Contributing

```bash
# Development setup
git clone https://github.com/hckerhub/ipcrawler.git
cd ipcrawler
go mod download
go test ./...
```

## 📞 Support & Author

**Author**: [hckerhub](https://github.com/hckerhub)

- 🐛 **Issues**: [GitHub Issues](https://github.com/hckerhub/ipcrawler/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/hckerhub/ipcrawler/discussions)
- 🐦 **Twitter/X**: [@hckerhub](https://x.com/hckerhub)
- ☕ **Buy me a coffee**: [buymeacoffee.com/hckerhub](https://www.buymeacoffee.com/hckerhub)

## 📝 License

MIT License - see [LICENSE](LICENSE) file for details.

## ⚠️ Disclaimer

For ethical hacking and educational purposes only. Use only on systems you own or have permission to test.

## 🙏 Acknowledgments

- [NMAP](https://nmap.org/) - Network exploration tool
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- Hack The Box community

---

**Made with ❤️ for the Hack The Box community by [hckerhub](https://github.com/hckerhub)** 