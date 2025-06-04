package platform

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// OSInfo contains detected operating system information
type OSInfo struct {
	OS           string            `json:"os"`           // windows, darwin, linux
	Architecture string            `json:"architecture"` // amd64, arm64, etc.
	Distribution string            `json:"distribution"` // ubuntu, fedora, etc. (Linux only)
	Version      string            `json:"version"`      // OS version
	Shell        string            `json:"shell"`        // bash, zsh, powershell, etc.
	Commands     map[string]string `json:"commands"`     // Available commands and their paths
	KeyBindings  KeyBindings       `json:"key_bindings"` // Platform-specific key bindings
}

// KeyBindings defines platform-specific keyboard shortcuts
type KeyBindings struct {
	Copy      string `json:"copy"`      // Ctrl+C vs Cmd+C
	Paste     string `json:"paste"`     // Ctrl+V vs Cmd+V
	Quit      string `json:"quit"`      // Ctrl+C vs Cmd+Q
	Interrupt string `json:"interrupt"` // Ctrl+C vs Cmd+.
	Tab       string `json:"tab"`       // Tab completion
	Enter     string `json:"enter"`     // Enter/Return
	Escape    string `json:"escape"`    // Escape key
}

// Detector provides OS and platform detection capabilities
type Detector struct {
	info *OSInfo
}

// NewDetector creates a new platform detector
func NewDetector() *Detector {
	return &Detector{}
}

// Detect performs comprehensive OS and platform detection
func (d *Detector) Detect() (*OSInfo, error) {
	if d.info != nil {
		return d.info, nil // Return cached result
	}

	info := &OSInfo{
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		Commands:     make(map[string]string),
	}

	// Detect OS-specific information
	switch runtime.GOOS {
	case "windows":
		d.detectWindows(info)
	case "darwin":
		d.detectMacOS(info)
	case "linux":
		d.detectLinux(info)
	default:
		d.detectGenericUnix(info)
	}

	// Detect available commands
	d.detectCommands(info)

	// Set key bindings
	d.setKeyBindings(info)

	d.info = info
	return info, nil
}

// detectWindows detects Windows-specific information
func (d *Detector) detectWindows(info *OSInfo) {
	// Detect Windows version
	if output, err := exec.Command("ver").Output(); err == nil {
		info.Version = strings.TrimSpace(string(output))
	}

	// Detect shell
	if os.Getenv("PSModulePath") != "" {
		info.Shell = "powershell"
	} else {
		info.Shell = "cmd"
	}

	// Check for WSL
	if os.Getenv("WSL_DISTRO_NAME") != "" {
		info.Distribution = "WSL"
	}
}

// detectMacOS detects macOS-specific information
func (d *Detector) detectMacOS(info *OSInfo) {
	// Detect macOS version
	if output, err := exec.Command("sw_vers", "-productVersion").Output(); err == nil {
		info.Version = strings.TrimSpace(string(output))
	}

	// Detect shell
	shell := os.Getenv("SHELL")
	if shell != "" {
		parts := strings.Split(shell, "/")
		info.Shell = parts[len(parts)-1]
	} else {
		info.Shell = "bash" // Default
	}

	// Check for Homebrew
	if _, err := exec.LookPath("brew"); err == nil {
		info.Commands["brew"] = "available"
	}
}

// detectLinux detects Linux-specific information and distribution
func (d *Detector) detectLinux(info *OSInfo) {
	// Detect distribution
	d.detectLinuxDistribution(info)

	// Detect shell
	shell := os.Getenv("SHELL")
	if shell != "" {
		parts := strings.Split(shell, "/")
		info.Shell = parts[len(parts)-1]
	} else {
		info.Shell = "bash" // Default
	}

	// Check for package managers
	packageManagers := map[string]string{
		"apt":    "apt-get",
		"yum":    "yum",
		"dnf":    "dnf",
		"pacman": "pacman",
		"zypper": "zypper",
		"emerge": "emerge",
		"apk":    "apk",
	}

	for pm, cmd := range packageManagers {
		if _, err := exec.LookPath(cmd); err == nil {
			info.Commands[pm] = cmd
		}
	}
}

// detectLinuxDistribution identifies the Linux distribution
func (d *Detector) detectLinuxDistribution(info *OSInfo) {
	// Try /etc/os-release first (modern standard)
	if data, err := os.ReadFile("/etc/os-release"); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "ID=") {
				info.Distribution = strings.Trim(strings.TrimPrefix(line, "ID="), "\"")
				return
			}
		}
	}

	// Fallback to other methods
	distroFiles := map[string]string{
		"/etc/redhat-release": "redhat",
		"/etc/debian_version": "debian",
		"/etc/arch-release":   "arch",
		"/etc/alpine-release": "alpine",
		"/etc/gentoo-release": "gentoo",
	}

	for file, distro := range distroFiles {
		if _, err := os.Stat(file); err == nil {
			info.Distribution = distro
			return
		}
	}

	info.Distribution = "unknown"
}

// detectGenericUnix detects information for other Unix-like systems
func (d *Detector) detectGenericUnix(info *OSInfo) {
	// Detect shell
	shell := os.Getenv("SHELL")
	if shell != "" {
		parts := strings.Split(shell, "/")
		info.Shell = parts[len(parts)-1]
	} else {
		info.Shell = "sh" // Most basic shell
	}
}

// detectCommands checks for availability of important commands
func (d *Detector) detectCommands(info *OSInfo) {
	// Essential commands to check
	commands := []string{
		"nmap", "curl", "wget", "git", "make", "sudo",
		"grep", "awk", "sed", "find", "tar", "zip",
		"python", "python3", "go", "node", "docker",
	}

	for _, cmd := range commands {
		if path, err := exec.LookPath(cmd); err == nil {
			info.Commands[cmd] = path
		}
	}

	// Check for privilege escalation commands
	privCommands := []string{"sudo", "doas", "su"}
	for _, cmd := range privCommands {
		if path, err := exec.LookPath(cmd); err == nil {
			info.Commands["privilege"] = cmd
			info.Commands["privilege_path"] = path
			break
		}
	}
}

// setKeyBindings sets appropriate key bindings for the platform
func (d *Detector) setKeyBindings(info *OSInfo) {
	switch info.OS {
	case "darwin": // macOS
		info.KeyBindings = KeyBindings{
			Copy:      "Cmd+C",
			Paste:     "Cmd+V",
			Quit:      "Cmd+Q",
			Interrupt: "Cmd+.",
			Tab:       "Tab",
			Enter:     "Return",
			Escape:    "Esc",
		}
	case "windows":
		info.KeyBindings = KeyBindings{
			Copy:      "Ctrl+C",
			Paste:     "Ctrl+V",
			Quit:      "Alt+F4",
			Interrupt: "Ctrl+C",
			Tab:       "Tab",
			Enter:     "Enter",
			Escape:    "Esc",
		}
	default: // Linux and other Unix-like
		info.KeyBindings = KeyBindings{
			Copy:      "Ctrl+Shift+C",
			Paste:     "Ctrl+Shift+V",
			Quit:      "Ctrl+C",
			Interrupt: "Ctrl+C",
			Tab:       "Tab",
			Enter:     "Enter",
			Escape:    "Esc",
		}
	}
}

// GetInstallCommand returns the appropriate install command for the platform
func (d *Detector) GetInstallCommand(packageName string) (string, error) {
	info, err := d.Detect()
	if err != nil {
		return "", err
	}

	switch info.OS {
	case "darwin":
		if _, exists := info.Commands["brew"]; exists {
			return fmt.Sprintf("brew install %s", packageName), nil
		}
		return "", fmt.Errorf("homebrew not found on macOS")

	case "linux":
		// Try package managers in order of preference
		if _, exists := info.Commands["apt"]; exists {
			return fmt.Sprintf("apt-get install -y %s", packageName), nil
		}
		if _, exists := info.Commands["dnf"]; exists {
			return fmt.Sprintf("dnf install -y %s", packageName), nil
		}
		if _, exists := info.Commands["yum"]; exists {
			return fmt.Sprintf("yum install -y %s", packageName), nil
		}
		if _, exists := info.Commands["pacman"]; exists {
			return fmt.Sprintf("pacman -S --noconfirm %s", packageName), nil
		}
		if _, exists := info.Commands["apk"]; exists {
			return fmt.Sprintf("apk add %s", packageName), nil
		}
		return "", fmt.Errorf("no supported package manager found")

	case "windows":
		// Check for package managers on Windows
		if _, err := exec.LookPath("choco"); err == nil {
			return fmt.Sprintf("choco install %s", packageName), nil
		}
		if _, err := exec.LookPath("winget"); err == nil {
			return fmt.Sprintf("winget install %s", packageName), nil
		}
		return "", fmt.Errorf("no supported package manager found on Windows")

	default:
		return "", fmt.Errorf("unsupported operating system: %s", info.OS)
	}
}

// GetSystemPaths returns appropriate system paths for binary installation
func (d *Detector) GetSystemPaths() ([]string, error) {
	info, err := d.Detect()
	if err != nil {
		return nil, err
	}

	switch info.OS {
	case "darwin", "linux":
		return []string{
			"/usr/local/bin",
			"/usr/bin",
			fmt.Sprintf("%s/bin", os.Getenv("HOME")),
		}, nil
	case "windows":
		return []string{
			"C:\\Program Files\\ipcrawler",
			fmt.Sprintf("%s\\AppData\\Local\\ipcrawler", os.Getenv("USERPROFILE")),
		}, nil
	default:
		return []string{"/usr/local/bin"}, nil
	}
}

// IsInPath checks if a directory is in the system PATH
func (d *Detector) IsInPath(dir string) bool {
	pathEnv := os.Getenv("PATH")
	paths := strings.Split(pathEnv, string(os.PathListSeparator))

	for _, path := range paths {
		if path == dir {
			return true
		}
	}
	return false
}

// GetPreferredInstallPath returns the best path for installing the binary
func (d *Detector) GetPreferredInstallPath() (string, error) {
	paths, err := d.GetSystemPaths()
	if err != nil {
		return "", err
	}

	// Find the first path that exists and is in PATH
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil && d.IsInPath(path) {
			return path, nil
		}
	}

	// If none exist, return the first one (will be created)
	if len(paths) > 0 {
		return paths[0], nil
	}

	return "", fmt.Errorf("no suitable installation path found")
}
