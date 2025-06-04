# IPCrawler

Advanced IP Analysis & Penetration Testing Tool with Modern TUI Interface

## Overview

IPCrawler is a comprehensive IP reconnaissance and network scanning tool built with Go, featuring a modern Terminal User Interface (TUI) powered by Bubble Tea and Lipgloss. Designed specifically for penetration testing and Hack the Box CTFs, it provides an intuitive interface for network discovery, service enumeration, and vulnerability assessment.

## Features

### 🔍 **Advanced Scanning Capabilities**
- **Quick Scan**: Fast reconnaissance of top 1000 ports
- **Full Port Scan**: Comprehensive scan of all 65535 TCP ports
- **Stealth Scan**: TCP SYN stealth scanning for evasion
- **Service Detection**: Detailed service version identification
- **OS Detection**: Operating system fingerprinting
- **Vulnerability Scan**: Automated vulnerability assessment
- **Web Scan**: Specialized HTTP/HTTPS analysis (ports 80/443)
- **SSH Scan**: Dedicated SSH service analysis (port 22)

### 🖥️ **Modern TUI Interface**
- Beautiful, responsive terminal interface
- Intuitive navigation with keyboard shortcuts
- Real-time scan progress and results
- Toggleable sidebar for easy navigation
- Syntax-highlighted output and results
- Modern color scheme optimized for dark terminals

### 🛠️ **Integrated Tools**
- **nmap**: Network discovery and security scanning
- **whatweb**: Web technology identification
- **curl**: HTTP requests and web analysis
- **dig**: DNS lookup utilities
- Additional pentesting tools integration

### 🎯 **Hack The Box Optimized**
- Pre-configured scan profiles for CTF environments
- Quick target analysis workflows
- Service enumeration best practices
- Common vulnerability checking patterns

## Installation

### Prerequisites
- Go 1.21 or higher
- Unix-like operating system (Linux, macOS)
- Terminal with 256-color support

### Quick Install
```bash
# Clone the repository
git clone <repository-url>
cd ipcrawler

# Run the installer (checks dependencies and installs tools)
./ipcrawler-installer

# Run the application
./ipcrawler
```

### Manual Installation
```bash
# Install Go dependencies
go mod tidy

# Build the application
go build -o bin/ipcrawler .

# Run the application
./bin/ipcrawler
```

## Usage

### Basic Commands
```bash
# Run the TUI interface
./ipcrawler

# Install dependencies and tools
./ipcrawler-installer
```

### TUI Navigation
- **Tab**: Toggle sidebar
- **Enter**: Select/Activate
- **↑/↓ or j/k**: Navigate up/down
- **←/→ or h/l**: Navigate left/right
- **Esc**: Back/Cancel
- **q**: Quit application
- **?**: Show help

### Quick Start Guide
1. Launch IPCrawler: `./ipcrawler`
2. Navigate to "IP Scanner" using arrow keys and Enter
3. Enter target IP address
4. Select scan type (Quick Scan recommended for start)
5. Press Enter to begin scanning
6. View results in the "Results" section

## Scan Types

### Quick Scan
Fast reconnaissance targeting the most common 1000 ports. Ideal for initial target assessment.

### Full Port Scan
Comprehensive scan of all 65535 TCP ports. Use when thoroughness is required.

### Stealth Scan
TCP SYN scan designed to be less detectable by intrusion detection systems.

### Service Detection
Identifies service versions running on open ports. Essential for vulnerability research.

### OS Detection
Attempts to identify the target operating system using various fingerprinting techniques.

### Vulnerability Scan
Automated vulnerability assessment using NSE scripts. Highlights potential security issues.

### Web Scan
Specialized scan for web services, including HTTP header analysis, directory enumeration, and technology identification.

### SSH Scan
Focused analysis of SSH services, including authentication methods and key algorithms.

## Architecture

```
ipcrawler/
├── cmd/           # Cobra CLI commands
├── ui/            # TUI components
│   ├── models/    # Bubble Tea models
│   ├── views/     # UI view functions
│   └── styles/    # Lipgloss styling
├── scanner/       # Scanning logic
├── tools/         # External tool integrations
├── main.go        # Application entry point
├── ipcrawler      # Executable script
└── ipcrawler-installer # Installation script
```

## Configuration

### Tool Paths
The application automatically detects installed tools. Custom paths can be configured in the settings.

### Scan Options
- Timing templates (T0-T5)
- Thread count adjustment
- Custom nmap arguments
- Output format selection

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## Security Notice

This tool is intended for authorized security testing only. Users are responsible for ensuring they have proper authorization before scanning any networks or systems.

## License

[Add your license here]

## Changelog

### v1.0.0
- Initial release
- Basic nmap integration
- Modern TUI interface
- Multiple scan types
- Real-time progress tracking

## Support

For issues, questions, or contributions, please refer to the project's issue tracker.

---

**Disclaimer**: This tool is for educational and authorized testing purposes only. Always ensure you have explicit permission before scanning networks or systems you do not own. 