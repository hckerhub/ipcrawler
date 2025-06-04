package views

import (
	"fmt"
	"strings"

	"ipcrawler/ui/styles"
)

// RenderScanPreview renders the scan screen preview
func RenderScanPreview() string {
	content := strings.Join([]string{
		styles.RenderHeader("IP Scanner"),
		"",
		"🔍 Advanced network scanning and reconnaissance",
		"• Multiple scan types (Quick, Stealth, Full)",
		"• Service and OS detection",
		"• Vulnerability assessment",
		"• Custom port ranges",
		"",
		styles.RenderInfo("Click to enter full scanner interface"),
	}, "\n")

	return styles.RenderCard(content)
}

// RenderResultsPreview renders the results screen preview
func RenderResultsPreview() string {
	content := strings.Join([]string{
		styles.RenderHeader("Scan Results"),
		"",
		"📊 View and analyze scan results",
		"• Detailed port information",
		"• Service version detection",
		"• Vulnerability findings",
		"• Export capabilities",
		"",
		styles.RenderInfo("Click to view detailed results"),
	}, "\n")

	return styles.RenderCard(content)
}

// RenderSettingsPreview renders the settings screen preview
func RenderSettingsPreview() string {
	content := strings.Join([]string{
		styles.RenderHeader("Settings"),
		"",
		"⚙️ Configure tools and preferences",
		"• Nmap scan options",
		"• Tool paths and configurations",
		"• Output format settings",
		"• Performance tuning",
		"",
		styles.RenderInfo("Click to configure settings"),
	}, "\n")

	return styles.RenderCard(content)
}

// RenderHelpPreview renders the help screen preview
func RenderHelpPreview() string {
	content := strings.Join([]string{
		styles.RenderHeader("Help & Documentation"),
		"",
		"❓ Get help and learn shortcuts",
		"• Keyboard shortcuts",
		"• Scanning techniques",
		"• Tool documentation",
		"• Best practices",
		"",
		styles.RenderInfo("Click for detailed help"),
	}, "\n")

	return styles.RenderCard(content)
}

// RenderSettingsFull renders the full settings screen
func RenderSettingsFull(width, height int) string {
	title := styles.RenderTitle("Settings & Configuration")

	sections := []string{
		styles.RenderHeader("Tool Configuration"),
		"",
		"📁 Tool Paths:",
		fmt.Sprintf("  %-15s %s", "nmap:", "/usr/bin/nmap"),
		fmt.Sprintf("  %-15s %s", "whatweb:", "/usr/bin/whatweb"),
		fmt.Sprintf("  %-15s %s", "curl:", "/usr/bin/curl"),
		"",
		styles.RenderHeader("Scan Options"),
		"",
		"🔧 Default Settings:",
		"  ├ Timeout: 30 seconds",
		"  ├ Threads: 50",
		"  ├ Rate limit: 1000 packets/sec",
		"  └ Output format: JSON",
		"",
		styles.RenderHeader("Advanced Options"),
		"",
		"⚡ Performance:",
		"  ├ Enable aggressive scanning: No",
		"  ├ Use TCP SYN scan: Yes",
		"  ├ Skip host discovery: No",
		"  └ Fragment packets: No",
		"",
		styles.RenderInfo("Use arrow keys to navigate • Enter to edit • Esc to go back"),
	}

	content := strings.Join(sections, "\n")
	return styles.ContentStyle.Width(width).Height(height).Render(title + "\n\n" + content)
}

// RenderHelpFull renders the full help screen
func RenderHelpFull(width, height int) string {
	title := styles.RenderTitle("Help & Documentation")

	sections := []string{
		styles.RenderHeader("Keyboard Shortcuts"),
		"",
		"Navigation:",
		"  ├ Tab/Shift+Tab: Switch panels",
		"  ├ ↑/↓ or j/k: Move up/down",
		"  ├ ←/→ or h/l: Move left/right",
		"  ├ Enter: Select/Activate",
		"  ├ Esc: Back/Cancel",
		"  └ q: Quit application",
		"",
		styles.RenderHeader("Scanning Guide"),
		"",
		"🎯 Quick Start:",
		"  1. Enter target IP address",
		"  2. Select scan type",
		"  3. Configure port range (optional)",
		"  4. Press Enter to start scan",
		"",
		"🔍 Scan Types:",
		"  ├ Quick Scan: Top 1000 ports",
		"  ├ Full Port Scan: All 65535 ports",
		"  ├ Stealth Scan: SYN scan (stealthy)",
		"  ├ Service Detection: Identify services",
		"  ├ OS Detection: Detect operating system",
		"  ├ Vulnerability Scan: Check for vulns",
		"  ├ Web Scan: HTTP/HTTPS specific",
		"  └ SSH Scan: SSH specific analysis",
		"",
		styles.RenderHeader("Tips for Hack The Box"),
		"",
		"💡 Best Practices:",
		"  ├ Start with Quick Scan for overview",
		"  ├ Use Service Detection on open ports",
		"  ├ Check web services with Web Scan",
		"  ├ Analyze SSH services carefully",
		"  └ Look for unusual service versions",
		"",
		styles.RenderInfo("Press Esc to go back to the main menu"),
	}

	content := strings.Join(sections, "\n")
	return styles.ContentStyle.Width(width).Height(height).Render(title + "\n\n" + content)
}
