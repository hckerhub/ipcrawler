package scanner

import (
	"fmt"
	"os/exec"
	"strings"
)

// ToolAvailability tracks which tools are available on the system
type ToolAvailability struct {
	Available map[string]bool
	Checked   bool
}

// NewToolAvailability creates a new tool availability checker
func NewToolAvailability() *ToolAvailability {
	return &ToolAvailability{
		Available: make(map[string]bool),
		Checked:   false,
	}
}

// CheckAllTools checks the availability of all tools used by ipcrawler
func (ta *ToolAvailability) CheckAllTools() {
	if ta.Checked {
		return
	}

	tools := []string{
		// Core tools
		"nmap", "curl", "dig",

		// Web analysis tools
		"whatweb", "wappalyzer", "ffuf", "feroxbuster", "gobuster",

		// Reconnaissance tools
		"subfinder", "sublist3r", "amass", "assetfinder", "findomain",
		"dnsrecon", "recon-ng", "censys",
	}

	for _, tool := range tools {
		ta.Available[tool] = ta.isToolAvailable(tool)
	}

	ta.Checked = true
}

// isToolAvailable checks if a specific tool is available
func (ta *ToolAvailability) isToolAvailable(toolName string) bool {
	_, err := exec.LookPath(toolName)
	return err == nil
}

// IsAvailable returns whether a specific tool is available
func (ta *ToolAvailability) IsAvailable(toolName string) bool {
	if !ta.Checked {
		ta.CheckAllTools()
	}
	return ta.Available[toolName]
}

// GetAvailableTools returns a list of available tools
func (ta *ToolAvailability) GetAvailableTools() []string {
	if !ta.Checked {
		ta.CheckAllTools()
	}

	var available []string
	for tool, isAvailable := range ta.Available {
		if isAvailable {
			available = append(available, tool)
		}
	}
	return available
}

// GetMissingTools returns a list of missing tools
func (ta *ToolAvailability) GetMissingTools() []string {
	if !ta.Checked {
		ta.CheckAllTools()
	}

	var missing []string
	for tool, isAvailable := range ta.Available {
		if !isAvailable {
			missing = append(missing, tool)
		}
	}
	return missing
}

// GetToolStatus returns a detailed status report
func (ta *ToolAvailability) GetToolStatus() string {
	if !ta.Checked {
		ta.CheckAllTools()
	}

	var report strings.Builder

	// Core tools
	coreTools := []string{"nmap", "curl", "dig"}
	report.WriteString("🔧 Core Tools:\n")
	for _, tool := range coreTools {
		status := "❌ Missing"
		if ta.Available[tool] {
			status = "✅ Available"
		}
		report.WriteString(fmt.Sprintf("   %s: %s\n", tool, status))
	}

	// Web analysis tools
	webTools := []string{"whatweb", "wappalyzer", "ffuf", "feroxbuster", "gobuster"}
	report.WriteString("\n🕸️  Web Analysis Tools:\n")
	for _, tool := range webTools {
		status := "❌ Missing"
		if ta.Available[tool] {
			status = "✅ Available"
		}
		report.WriteString(fmt.Sprintf("   %s: %s\n", tool, status))
	}

	// Reconnaissance tools
	reconTools := []string{"subfinder", "sublist3r", "amass", "assetfinder", "findomain", "dnsrecon", "recon-ng", "censys"}
	report.WriteString("\n🔍 Reconnaissance Tools:\n")
	for _, tool := range reconTools {
		status := "❌ Missing"
		if ta.Available[tool] {
			status = "✅ Available"
		}
		report.WriteString(fmt.Sprintf("   %s: %s\n", tool, status))
	}

	// Summary
	availableCount := len(ta.GetAvailableTools())
	totalCount := len(ta.Available)
	missingCount := totalCount - availableCount

	report.WriteString(fmt.Sprintf("\n📊 Summary: %d/%d tools available", availableCount, totalCount))
	if missingCount > 0 {
		report.WriteString(fmt.Sprintf(" (%d missing)", missingCount))
		report.WriteString("\n\n💡 To install missing tools, run: ./ipcrawler-installer")
	}

	return report.String()
}

// GetCapabilities returns what the scanner can do with available tools
func (ta *ToolAvailability) GetCapabilities() map[string]bool {
	if !ta.Checked {
		ta.CheckAllTools()
	}

	capabilities := map[string]bool{
		"port_scanning":        ta.Available["nmap"],
		"basic_web_requests":   ta.Available["curl"],
		"dns_lookups":          ta.Available["dig"],
		"web_tech_detection":   ta.Available["whatweb"] || ta.Available["wappalyzer"],
		"directory_bruteforce": ta.Available["ffuf"] || ta.Available["feroxbuster"] || ta.Available["gobuster"],
		"subdomain_discovery":  ta.Available["subfinder"] || ta.Available["sublist3r"] || ta.Available["amass"],
		"advanced_recon":       ta.Available["amass"] || ta.Available["recon-ng"],
	}

	return capabilities
}

// GetMissingCapabilities returns a user-friendly explanation of what won't work
func (ta *ToolAvailability) GetMissingCapabilities() []string {
	capabilities := ta.GetCapabilities()
	var missing []string

	if !capabilities["port_scanning"] {
		missing = append(missing, "⚠️  Port scanning disabled (nmap not found)")
	}
	if !capabilities["web_tech_detection"] {
		missing = append(missing, "⚠️  Web technology detection limited (whatweb/wappalyzer not found)")
	}
	if !capabilities["directory_bruteforce"] {
		missing = append(missing, "⚠️  Directory bruteforcing disabled (ffuf/feroxbuster/gobuster not found)")
	}
	if !capabilities["subdomain_discovery"] {
		missing = append(missing, "⚠️  Subdomain discovery limited (subfinder/sublist3r/amass not found)")
	}
	if !capabilities["advanced_recon"] {
		missing = append(missing, "⚠️  Advanced reconnaissance limited (amass/recon-ng not found)")
	}

	return missing
}
