# 📋 Changelog

All notable changes to **ipcrawler** will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [1.1.4] - 2025-01-12 🛠️

### 🛠️ Plugin & Backup Improvements
- **Fixed Plugin Loading**: Resolved "datetime not a subclass" error by filtering out imported classes from plugin validation
- **Enhanced Backup Location**: VHost /etc/hosts backups now stored in target scan directory instead of /etc/ (survives terminal closure)
- **Improved Plugin Validation**: Plugin loader now only validates classes defined in plugin files, not imported dependencies
- **Better Backup Persistence**: Backup files stored in `results/IP/hosts.backup.timestamp` for easy access and restoration

---

## [1.1.3] - 2025-01-12 🔥

### 🔥 HOTFIX - Critical IP Extraction Bug
- **CRITICAL**: Fixed IP extraction logic to use parent directory instead of scan directory name
- **Root Cause**: VHost post-processor was extracting IP from "scans" directory instead of parent target directory
- **Impact**: Now correctly adds `10.10.11.68 planning.htb` instead of `scans planning.htb` to /etc/hosts
- **Directory Structure**: Properly handles ipcrawler's structure: `results/IP/scans/` → extracts IP from parent directory

---

## [1.1.2] - 2025-01-12 🚨

### 🚨 Critical Bug Fixes
- **CRITICAL**: Fixed incorrect IP address extraction in VHost post-processor (was using directory name "scans" instead of target IP)
- **CRITICAL**: Fixed invisible input during VHost interactive prompts with enhanced terminal handling
- **Enhanced Input System**: Added character echoing and proper terminal mode management for better user experience
- **Improved Error Handling**: Better fallback options for input failures with multiple retry attempts
- **Visual Feedback**: Enhanced interactive VHost management sessions with clearer prompts and status messages

---

## [1.1.1] - 2025-06-12 🔧

### 🐛 Fixed
- **VHost Discovery Priority**: VHost Redirect Hunter now runs before Virtual Host Enumeration (priority 10 vs 5)
- **Post-Processing Path Resolution**: Fixed post-processor to correctly scan all target directories instead of single 'scans' directory
- **File Detection Logic**: Enhanced VHost file discovery to properly locate `vhost_redirects_*.txt` files in target subdirectories
- **Plugin Execution Order**: Ensured VHost Redirect Hunter takes priority over other VHost enumeration methods

### 🔧 Enhanced
- **Error Handling**: Improved error messages and debugging information for VHost discovery issues
- **Directory Scanning**: More robust directory traversal for multi-target scenarios

---

## [1.1.0] - 2025-06-12 🌐

### 🚀 VHost Discovery Release
**ipcrawler** v1.1.0 introduces comprehensive virtual host discovery capabilities with interactive `/etc/hosts` management, making it even easier to handle complex web applications during penetration testing.

### ✨ Added
- **🌐 VHost Discovery System**
  - **VHost Redirect Hunter Plugin**: Automatic discovery of hostnames via HTTP redirect analysis
  - **Interactive Post-Processing**: Smart `/etc/hosts` management with Y/N/S options for user control
  - **Safety Features**: Automatic backup of `/etc/hosts` before modifications with timestamped files
  - **Configuration-Driven**: Full control via `[vhost_discovery]` section in config.toml
  - **Intelligent Detection**: Duplicate prevention and existing entry detection
  - **Manual Fallback**: Copy-paste commands generated when sudo privileges unavailable
  - **Beautiful Output**: Emoji-rich progress indicators and summary tables
  - **Integration**: Seamless integration with existing HTML reporting system

### 🔧 Enhanced
- **Post-Scan Processing**: Automatic VHost management runs after all scans complete
- **Configuration System**: New VHost-specific settings (timeout, user-agent, auto-add behavior)
- **Error Handling**: Graceful fallbacks and comprehensive error management
- **Reporting Integration**: VHost discoveries automatically included in Rich HTML reports

---

## [1.0.0] - 2025-06-10 🎉

### 🚀 Initial Release
**ipcrawler** v1.0.0 marks the first stable release of our simplified AutoRecon fork, designed to make network reconnaissance accessible for CTFs, OSCP, and penetration testing.

### ✨ Added
- **🎨 Enhanced User Experience**
  - Beautiful Rich-formatted `--help` with organized sections, examples, and pro tips
  - Rich-colored output for verbosity levels (`-v`, `-vv`, `-vvv`) with visual progress indicators
  - Professional command-line interface with improved readability
  
- **🔧 Streamlined Setup Process**
  - One-command setup with `make setup` for local installation
  - Docker support with `make setup-docker` for cross-platform compatibility
  - Automatic Docker detection and smart installation workflow
  - `bootstrap.sh` script for OS-specific dependency management
  - Windows `.bat` file for seamless Docker integration

- **📦 Platform Support**
  - Native Linux/macOS support with local installation
  - Full Windows compatibility via Docker
  - HTB machine compatibility with optimized tool paths
  - Cross-platform wordlist and dependency management

- **🛠️ Build System & Automation**
  - Comprehensive Makefile with setup, docker, update, and cleanup commands
  - Automated virtual environment creation and management
  - Docker image building with security tools pre-installed
  - Global command wrapper installation (`ipcrawler-cmd`)

- **📚 Documentation & User Guides**
  - Step-by-step setup instructions for all platforms
  - Video tutorial integration for HTB and macOS setups
  - Comprehensive README with troubleshooting guides
  - Configuration file documentation with TOML examples

- **📊 Enhanced Reporting System**
  - Rich HTML Summary Reports with interactive collapsible sections
  - Key findings extraction and executive summary generation
  - Beautiful CSS styling with responsive design and dark code themes
  - Combined multi-target reporting capabilities
  - Automatic discovery and presentation of URLs, domains, vulnerabilities
  - Technology stack detection and credential extraction

- **🕷️ VHost Redirect Hunter (NEW)**
  - Automatic discovery of hostnames via HTTP redirect analysis
  - Post-scan interactive prompt to add discovered VHosts to `/etc/hosts`
  - Smart privilege elevation - prompts for `sudo` only when needed
  - Beautiful emoji-rich output with summary tables
  - Automatic backup of `/etc/hosts` before modifications
  - Duplicate detection and manual command generation
  - Configurable timeouts and user agents
  - File-based results organization with intelligent filtering
  - Print-friendly layouts and mobile-responsive design

### 🔄 Changed (from AutoRecon)
- **Simplified Installation**: From complex `pipx` commands to simple `make setup`
- **Enhanced CLI**: Rich formatting replaces basic terminal output
- **Better Error Handling**: Improved error messages and dependency checking
- **Modern Configuration**: TOML-based config files with better organization
- **Streamlined Workflow**: Automated setup replaces manual dependency management

### 🔧 Technical Improvements
- **Dependencies**: Added Rich library for enhanced terminal output
- **Python Requirements**: Maintained Python 3.8+ compatibility
- **Plugin System**: Inherited 80+ reconnaissance plugins from AutoRecon
- **Architecture**: Multi-threaded scanning with concurrent execution
- **Configuration**: Enhanced TOML configuration with user-friendly defaults

### 🗂️ Core Features (Inherited from AutoRecon)
- **Port Scanning**: Comprehensive TCP/UDP port discovery
- **Service Enumeration**: 80+ specialized plugins for service analysis
- **Web Application Testing**: Directory busting, vulnerability scanning
- **SMB Enumeration**: Share discovery, user enumeration, vulnerability checks
- **DNS Reconnaissance**: Zone transfers, subdomain enumeration
- **Database Testing**: MySQL, MSSQL, Oracle, MongoDB scanning
- **Brute Force Testing**: SSH, FTP, RDP, HTTP authentication
- **Reporting**: Multiple output formats (Markdown, CherryTree, Rich HTML)

### 🐛 Fixed
- Platform-specific installation issues resolved via Docker option
- Dependency conflicts addressed through virtual environment isolation
- Path resolution problems fixed for various operating systems
- WordList location issues resolved with automatic detection

---

## 🚀 Future Goal

### 🎯 [1.2.0] - YOLO Mode & Enhanced Automation
**Target Date**: Q1 2025

#### 🚀 Full YOLO Mode
- **Auto-Execute Recommended Commands**: Automatically run all commands from `_manual_commands.txt`
- **Smart Command Prioritization**: Execute high-value commands first (credential dumping, exploit attempts)
- **Interactive Confirmation**: Optional prompts for destructive commands
- **Parallel Execution**: Run multiple manual commands simultaneously where safe

#### 🛠️ Additional Tools & Templates
- **Extended Tool Integration**: Add more reconnaissance tools (subfinder, amass, feroxbuster, etc.)
- **CTF-Focused Templates**: Optimized plugin configurations for common CTF scenarios
- **Custom Wordlists**: Curated wordlists for specific environments (HTB, THM, OSCP)
- **Quick Exploit Modules**: Built-in common exploit execution (EternalBlue, PrintNightmare, etc.)

#### 📋 Smart Templates
- **Environment Detection**: Auto-detect HTB/THM/OSCP environments and adjust scanning approach
- **Target Profiling**: Automatically select optimal plugin sets based on discovered services
- **Time-Bounded Modes**: OSCP exam mode with strict time limits and prioritized scanning
- **Stealth Mode**: Reduced noise scanning for more realistic penetration testing scenarios

---

## 📝 Versioning Scheme

- **Major Version** (x.0.0): Breaking changes, major feature additions
- **Minor Version** (1.x.0): New features, backward compatible
- **Patch Version** (1.0.x): Bug fixes, security updates

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details on:
- Bug reports and feature requests
- Plugin development
- Documentation improvements
- Code contributions

## 📞 Support

- **GitHub Issues**: [Report bugs or request features](https://github.com/neur0map/ipcrawler/issues)
- **Documentation**: [Full documentation](https://github.com/neur0map/ipcrawler/wiki)
- **Community**: [Discord server](https://discord.gg/ipcrawler) for discussions and support

---

*Built with ❤️ for the cybersecurity community* 