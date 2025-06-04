package tools

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// ToolInstaller manages installation of external tools
type ToolInstaller struct {
	Verbose bool
	DryRun  bool
}

// Tool represents an external tool with installation information
type Tool struct {
	Name        string
	Description string
	CheckCmd    string
	InstallCmds map[string][]string // OS -> commands
	Required    bool
	Category    string
}

// NewToolInstaller creates a new tool installer
func NewToolInstaller(verbose, dryRun bool) *ToolInstaller {
	return &ToolInstaller{
		Verbose: verbose,
		DryRun:  dryRun,
	}
}

// GetRequiredTools returns list of tools needed for IPCrawler
func (ti *ToolInstaller) GetRequiredTools() []Tool {
	return []Tool{
		{
			Name:        "nmap",
			Description: "Network discovery and security auditing tool",
			CheckCmd:    "nmap --version",
			Required:    true,
			Category:    "core",
			InstallCmds: map[string][]string{
				"darwin": {"brew install nmap"},
				"linux":  {"sudo apt-get update && sudo apt-get install -y nmap", "sudo yum install -y nmap", "sudo pacman -S nmap"},
			},
		},

		{
			Name:        "amass",
			Description: "In-depth attack surface mapping and asset discovery",
			CheckCmd:    "amass enum --help",
			Required:    false,
			Category:    "reconnaissance",
			InstallCmds: map[string][]string{
				"darwin": {"brew install amass"},
				"linux":  {"sudo snap install amass", "go install -v github.com/owasp-amass/amass/v4/...@master"},
			},
		},
		{
			Name:        "recon-ng",
			Description: "Web reconnaissance framework",
			CheckCmd:    "recon-ng --version",
			Required:    false,
			Category:    "reconnaissance",
			InstallCmds: map[string][]string{
				"darwin": {"pip3 install recon-ng"},
				"linux":  {"pip3 install recon-ng", "sudo pip3 install recon-ng"},
			},
		},
		{
			Name:        "censys",
			Description: "Censys CLI for internet-wide scanning data",
			CheckCmd:    "censys --version",
			Required:    false,
			Category:    "reconnaissance",
			InstallCmds: map[string][]string{
				"darwin": {"pip3 install censys"},
				"linux":  {"pip3 install censys", "sudo pip3 install censys"},
			},
		},
		{
			Name:        "whatweb",
			Description: "Web application fingerprinting tool",
			CheckCmd:    "whatweb --version",
			Required:    false,
			Category:    "web_analysis",
			InstallCmds: map[string][]string{
				"darwin": {"brew install whatweb"},
				"linux":  {"sudo apt-get install -y whatweb", "sudo yum install -y whatweb", "gem install whatweb"},
			},
		},
		{
			Name:        "wappalyzer",
			Description: "Technology profiler for websites",
			CheckCmd:    "wappalyzer --version",
			Required:    false,
			Category:    "web_analysis",
			InstallCmds: map[string][]string{
				"darwin": {"npm install -g wappalyzer"},
				"linux":  {"npm install -g wappalyzer", "sudo npm install -g wappalyzer"},
			},
		},
		{
			Name:        "curl",
			Description: "Command line tool for transferring data",
			CheckCmd:    "curl --version",
			Required:    true,
			Category:    "core",
			InstallCmds: map[string][]string{
				"darwin": {"curl is pre-installed on macOS"},
				"linux":  {"sudo apt-get install -y curl", "sudo yum install -y curl"},
			},
		},
		{
			Name:        "ffuf",
			Description: "Fast web fuzzer",
			CheckCmd:    "ffuf -V",
			Required:    false,
			Category:    "web_analysis",
			InstallCmds: map[string][]string{
				"darwin": {"brew install ffuf"},
				"linux":  {"go install github.com/ffuf/ffuf/v2@latest", "sudo snap install ffuf"},
			},
		},
		{
			Name:        "feroxbuster",
			Description: "Fast, simple, recursive content discovery tool",
			CheckCmd:    "feroxbuster --version",
			Required:    false,
			Category:    "web_analysis",
			InstallCmds: map[string][]string{
				"darwin": {"brew install feroxbuster"},
				"linux":  {"sudo snap install feroxbuster", "cargo install feroxbuster"},
			},
		},
		{
			Name:        "gobuster",
			Description: "Directory/file, DNS and VHost busting tool",
			CheckCmd:    "gobuster version",
			Required:    false,
			Category:    "web_analysis",
			InstallCmds: map[string][]string{
				"darwin": {"brew install gobuster"},
				"linux":  {"sudo apt-get install -y gobuster", "go install github.com/OJ/gobuster/v3@latest"},
			},
		},
		{
			Name:        "subfinder",
			Description: "Subdomain discovery tool",
			CheckCmd:    "subfinder -version",
			Required:    false,
			Category:    "reconnaissance",
			InstallCmds: map[string][]string{
				"darwin": {"brew install subfinder"},
				"linux":  {"go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest"},
			},
		},
		{
			Name:        "sublist3r",
			Description: "Fast subdomains enumeration tool",
			CheckCmd:    "sublist3r --help",
			Required:    false,
			Category:    "reconnaissance",
			InstallCmds: map[string][]string{
				"darwin": {"pip3 install sublist3r"},
				"linux":  {"pip3 install sublist3r", "sudo pip3 install sublist3r"},
			},
		},
	}
}

// CheckToolStatus checks if a tool is installed
func (ti *ToolInstaller) CheckToolStatus(tool Tool) (bool, string) {
	cmd := exec.Command("sh", "-c", tool.CheckCmd+" >/dev/null 2>&1")
	err := cmd.Run()

	if err != nil {
		return false, fmt.Sprintf("❌ %s - Not installed", tool.Name)
	}

	return true, fmt.Sprintf("✅ %s - Installed", tool.Name)
}

// GetSystemOS detects the operating system
func (ti *ToolInstaller) GetSystemOS() string {
	switch runtime.GOOS {
	case "darwin":
		return "darwin"
	case "linux":
		return "linux"
	case "windows":
		return "windows"
	default:
		return "unknown"
	}
}

// InstallTool attempts to install a specific tool
func (ti *ToolInstaller) InstallTool(tool Tool) error {
	osType := ti.GetSystemOS()
	commands, exists := tool.InstallCmds[osType]

	if !exists {
		return fmt.Errorf("no installation commands available for %s on %s", tool.Name, osType)
	}

	fmt.Printf("📦 Installing %s...\n", tool.Name)

	var lastErr error
	for _, cmdStr := range commands {
		if ti.Verbose {
			fmt.Printf("   Running: %s\n", cmdStr)
		}

		if ti.DryRun {
			fmt.Printf("   [DRY RUN] Would execute: %s\n", cmdStr)
			continue
		}

		// Try each command until one succeeds
		cmd := exec.Command("sh", "-c", cmdStr)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err == nil {
			fmt.Printf("✅ Successfully installed %s\n", tool.Name)
			return nil
		}

		lastErr = err
		if ti.Verbose {
			fmt.Printf("   ⚠️  Command failed: %v\n", err)
		}
	}

	return fmt.Errorf("failed to install %s: %v", tool.Name, lastErr)
}

// CheckAllTools checks status of all tools
func (ti *ToolInstaller) CheckAllTools() map[string]bool {
	tools := ti.GetRequiredTools()
	status := make(map[string]bool)

	fmt.Println("🔍 Checking tool availability...")
	fmt.Println()

	categories := make(map[string][]Tool)
	for _, tool := range tools {
		categories[tool.Category] = append(categories[tool.Category], tool)
	}

	for category, catTools := range categories {
		fmt.Printf("📂 %s Tools:\n", strings.Title(category))

		for _, tool := range catTools {
			installed, message := ti.CheckToolStatus(tool)
			status[tool.Name] = installed

			fmt.Printf("   %s", message)
			if tool.Required && !installed {
				fmt.Printf(" (REQUIRED)")
			}
			fmt.Println()
		}
		fmt.Println()
	}

	return status
}

// InstallMissingTools installs all missing tools
func (ti *ToolInstaller) InstallMissingTools() error {
	tools := ti.GetRequiredTools()
	status := ti.CheckAllTools()

	var missingRequired []Tool
	var missingOptional []Tool

	for _, tool := range tools {
		if !status[tool.Name] {
			if tool.Required {
				missingRequired = append(missingRequired, tool)
			} else {
				missingOptional = append(missingOptional, tool)
			}
		}
	}

	if len(missingRequired) == 0 && len(missingOptional) == 0 {
		fmt.Println("🎉 All tools are already installed!")
		return nil
	}

	// Install required tools first
	if len(missingRequired) > 0 {
		fmt.Printf("🚨 Installing %d required tools...\n\n", len(missingRequired))
		for _, tool := range missingRequired {
			if err := ti.InstallTool(tool); err != nil {
				return fmt.Errorf("failed to install required tool %s: %v", tool.Name, err)
			}
		}
	}

	// Install optional tools
	if len(missingOptional) > 0 {
		fmt.Printf("🔧 Installing %d optional tools for enhanced functionality...\n\n", len(missingOptional))
		for _, tool := range missingOptional {
			if err := ti.InstallTool(tool); err != nil {
				fmt.Printf("⚠️  Failed to install optional tool %s: %v\n", tool.Name, err)
				fmt.Printf("   IPCrawler will work without %s, but some features may be limited.\n\n", tool.Name)
			}
		}
	}

	fmt.Println("✨ Tool installation complete!")
	return nil
}

// GenerateInstallScript creates a shell script for manual installation
func (ti *ToolInstaller) GenerateInstallScript(filename string) error {
	tools := ti.GetRequiredTools()
	osType := ti.GetSystemOS()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString("#!/bin/bash\n")
	file.WriteString("# IPCrawler Tool Installation Script\n")
	file.WriteString("# Generated automatically - review before running\n\n")
	file.WriteString("set -e\n\n")

	file.WriteString("echo \"🚀 Installing IPCrawler dependencies...\"\n\n")

	for _, tool := range tools {
		if commands, exists := tool.InstallCmds[osType]; exists {
			file.WriteString(fmt.Sprintf("# Installing %s - %s\n", tool.Name, tool.Description))
			file.WriteString(fmt.Sprintf("echo \"📦 Installing %s...\"\n", tool.Name))

			for i, cmd := range commands {
				if i == 0 {
					file.WriteString(fmt.Sprintf("if ! command -v %s &> /dev/null; then\n", tool.Name))
					file.WriteString(fmt.Sprintf("    %s || {\n", cmd))
				} else {
					file.WriteString(fmt.Sprintf("        echo \"Trying alternative installation method...\"\n"))
					file.WriteString(fmt.Sprintf("        %s || {\n", cmd))
				}
			}

			file.WriteString("            echo \"Failed to install " + tool.Name + "\"\n")
			if tool.Required {
				file.WriteString("            exit 1\n")
			} else {
				file.WriteString("            echo \"Continuing without " + tool.Name + " (optional tool)\"\n")
			}

			for range tool.InstallCmds[osType] {
				file.WriteString("        }\n")
			}
			file.WriteString("else\n")
			file.WriteString(fmt.Sprintf("    echo \"✅ %s is already installed\"\n", tool.Name))
			file.WriteString("fi\n\n")
		}
	}

	file.WriteString("echo \"✨ Installation complete!\"\n")
	file.WriteString("echo \"Run './ipcrawler' to start using IPCrawler\"\n")

	// Make script executable
	err = os.Chmod(filename, 0755)
	if err != nil {
		return err
	}

	fmt.Printf("📝 Installation script saved to %s\n", filename)
	fmt.Printf("   Review the script and run: ./%s\n", filename)

	return nil
}

// GetMissingToolsReport generates a report of missing tools with installation instructions
func (ti *ToolInstaller) GetMissingToolsReport() string {
	tools := ti.GetRequiredTools()
	status := make(map[string]bool)

	// Quick check without output
	for _, tool := range tools {
		cmd := exec.Command("sh", "-c", tool.CheckCmd+" >/dev/null 2>&1")
		status[tool.Name] = cmd.Run() == nil
	}

	var report strings.Builder
	osType := ti.GetSystemOS()

	report.WriteString("IPCrawler Tool Status Report\n")
	report.WriteString("============================\n\n")

	// Group by category
	categories := make(map[string][]Tool)
	for _, tool := range tools {
		categories[tool.Category] = append(categories[tool.Category], tool)
	}

	for category, catTools := range categories {
		report.WriteString(fmt.Sprintf("%s Tools:\n", strings.Title(category)))
		report.WriteString(strings.Repeat("-", len(category)+7) + "\n")

		for _, tool := range catTools {
			installed := status[tool.Name]
			statusIcon := "❌"
			if installed {
				statusIcon = "✅"
			}

			report.WriteString(fmt.Sprintf("%s %s - %s", statusIcon, tool.Name, tool.Description))
			if tool.Required && !installed {
				report.WriteString(" (REQUIRED)")
			}
			report.WriteString("\n")

			if !installed {
				if commands, exists := tool.InstallCmds[osType]; exists {
					report.WriteString("   Installation: ")
					for i, cmd := range commands {
						if i > 0 {
							report.WriteString(" OR ")
						}
						report.WriteString(cmd)
					}
					report.WriteString("\n")
				}
			}
		}
		report.WriteString("\n")
	}

	// Summary
	installedCount := 0
	requiredMissing := 0
	for _, tool := range tools {
		if status[tool.Name] {
			installedCount++
		} else if tool.Required {
			requiredMissing++
		}
	}

	report.WriteString("Summary:\n")
	report.WriteString("--------\n")
	report.WriteString(fmt.Sprintf("Installed: %d/%d tools\n", installedCount, len(tools)))
	if requiredMissing > 0 {
		report.WriteString(fmt.Sprintf("Missing required tools: %d\n", requiredMissing))
		report.WriteString("⚠️  Some required tools are missing. IPCrawler may not function properly.\n")
	} else {
		report.WriteString("✅ All required tools are available.\n")
	}

	return report.String()
}
