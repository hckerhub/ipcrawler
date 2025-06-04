package models

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ExportResult exports scan results to a formatted text file
func (m *MainModel) ExportResults() (string, error) {
	if m.scanResult == nil {
		return "", fmt.Errorf("no scan results to export")
	}

	// Ensure logs directory exists
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create logs directory: %v", err)
	}

	// Generate filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("ipcrawler_scan_%s_%s.txt",
		strings.ReplaceAll(m.scanResult.IP, ".", "_"),
		timestamp)
	filepath := filepath.Join(logsDir, filename)

	// Generate formatted content
	content := m.generateReportContent()

	// Write to file
	if err := os.WriteFile(filepath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %v", err)
	}

	return filepath, nil
}

// generateReportContent creates a formatted text report
func (m *MainModel) generateReportContent() string {
	var b strings.Builder
	result := m.scanResult

	// Header
	b.WriteString("==================================================\n")
	b.WriteString("           IPCrawler Scan Report\n")
	b.WriteString("==================================================\n\n")

	// Scan Information
	b.WriteString("SCAN INFORMATION\n")
	b.WriteString("================\n")
	b.WriteString(fmt.Sprintf("Target IP:      %s\n", result.IP))
	b.WriteString(fmt.Sprintf("Scan Type:      %s", result.ScanType))
	if result.AggressiveMode {
		b.WriteString(" (Aggressive Mode)")
	}
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("Timestamp:      %s\n", result.Timestamp))
	b.WriteString(fmt.Sprintf("Duration:       %s\n", result.Duration))
	b.WriteString(fmt.Sprintf("Status:         %s\n\n", result.Status))

	// Summary Statistics
	b.WriteString("SCAN SUMMARY\n")
	b.WriteString("============\n")
	summary := result.Summary
	b.WriteString(fmt.Sprintf("Total Ports Scanned:    %d\n", summary.TotalPorts))
	b.WriteString(fmt.Sprintf("Open Ports:             %d\n", summary.OpenPorts))
	b.WriteString(fmt.Sprintf("Filtered Ports:         %d\n", summary.FilteredPorts))
	b.WriteString(fmt.Sprintf("Closed Ports:           %d\n", summary.ClosedPorts))
	b.WriteString(fmt.Sprintf("Web Services:           %d\n", summary.WebServices))
	b.WriteString(fmt.Sprintf("Vulnerabilities Found:  %d\n\n", summary.Vulnerabilities))

	// Port Details
	if len(result.Ports) > 0 {
		b.WriteString("PORT SCAN RESULTS\n")
		b.WriteString("=================\n")
		b.WriteString(fmt.Sprintf("%-8s %-8s %-8s %-15s %-20s %s\n",
			"PORT", "PROTOCOL", "STATE", "SERVICE", "VERSION", "SCRIPTS"))
		b.WriteString(strings.Repeat("-", 80))
		b.WriteString("\n")

		for _, port := range result.Ports {
			scripts := strings.Join(port.Scripts, ", ")
			if len(scripts) > 30 {
				scripts = scripts[:27] + "..."
			}
			b.WriteString(fmt.Sprintf("%-8d %-8s %-8s %-15s %-20s %s\n",
				port.Port, port.Protocol, port.State, port.Service, port.Version, scripts))
		}
		b.WriteString("\n")
	}

	// Detailed Service Information
	if len(result.Services) > 0 {
		b.WriteString("DETAILED SERVICE INFORMATION\n")
		b.WriteString("============================\n")
		for i, service := range result.Services {
			b.WriteString(fmt.Sprintf("%d. Port %d - %s\n", i+1, service.Port, service.Service))
			if service.Product != "" {
				b.WriteString(fmt.Sprintf("   Product:    %s\n", service.Product))
			}
			if service.Version != "" {
				b.WriteString(fmt.Sprintf("   Version:    %s\n", service.Version))
			}
			if service.ExtraInfo != "" {
				b.WriteString(fmt.Sprintf("   Extra Info: %s\n", service.ExtraInfo))
			}
			if service.OSType != "" {
				b.WriteString(fmt.Sprintf("   OS Type:    %s\n", service.OSType))
			}
			if service.Confidence > 0 {
				b.WriteString(fmt.Sprintf("   Confidence: %d%%\n", service.Confidence))
			}
			b.WriteString("\n")
		}
	}

	// Web Discovery Results
	if len(result.WebResults) > 0 {
		b.WriteString("WEB DISCOVERY RESULTS\n")
		b.WriteString("=====================\n")
		for i, web := range result.WebResults {
			b.WriteString(fmt.Sprintf("%d. %s\n", i+1, web.URL))
			if web.Title != "" {
				b.WriteString(fmt.Sprintf("   Title:        %s\n", web.Title))
			}
			if web.Server != "" {
				b.WriteString(fmt.Sprintf("   Server:       %s\n", web.Server))
			}
			if web.StatusCode > 0 {
				b.WriteString(fmt.Sprintf("   Status Code:  %d\n", web.StatusCode))
			}
			if len(web.Technologies) > 0 {
				b.WriteString(fmt.Sprintf("   Technologies: %s\n", strings.Join(web.Technologies, ", ")))
			}
			if len(web.Paths) > 0 {
				b.WriteString(fmt.Sprintf("   Discovered Paths:\n"))
				for _, path := range web.Paths {
					b.WriteString(fmt.Sprintf("     - %s\n", path))
				}
			}
			if len(web.Subdomains) > 0 {
				b.WriteString(fmt.Sprintf("   Subdomains:\n"))
				for _, subdomain := range web.Subdomains {
					b.WriteString(fmt.Sprintf("     - %s\n", subdomain))
				}
			}
			b.WriteString("\n")
		}
	}

	// Vulnerability Assessment
	if len(result.Vulnerabilities) > 0 {
		b.WriteString("VULNERABILITY ASSESSMENT\n")
		b.WriteString("========================\n")

		// Group by severity
		critical := []UIVulnInfo{}
		high := []UIVulnInfo{}
		medium := []UIVulnInfo{}
		low := []UIVulnInfo{}

		for _, vuln := range result.Vulnerabilities {
			switch vuln.Severity {
			case "CRITICAL":
				critical = append(critical, vuln)
			case "HIGH":
				high = append(high, vuln)
			case "MEDIUM":
				medium = append(medium, vuln)
			case "LOW":
				low = append(low, vuln)
			}
		}

		// Report by severity
		severityGroups := []struct {
			name   string
			vulns  []UIVulnInfo
			symbol string
		}{
			{"CRITICAL", critical, "🔴"},
			{"HIGH", high, "🟠"},
			{"MEDIUM", medium, "🟡"},
			{"LOW", low, "🔵"},
		}

		for _, group := range severityGroups {
			if len(group.vulns) > 0 {
				b.WriteString(fmt.Sprintf("%s %s SEVERITY (%d findings)\n",
					group.symbol, group.name, len(group.vulns)))
				b.WriteString(strings.Repeat("-", 40))
				b.WriteString("\n")

				for i, vuln := range group.vulns {
					b.WriteString(fmt.Sprintf("%d. Port %d (%s)\n", i+1, vuln.Port, vuln.Service))
					if vuln.CVE != "" {
						b.WriteString(fmt.Sprintf("   CVE:         %s\n", vuln.CVE))
					}
					if vuln.Description != "" {
						b.WriteString(fmt.Sprintf("   Description: %s\n", vuln.Description))
					}
					if vuln.Script != "" {
						b.WriteString(fmt.Sprintf("   Script:      %s\n", vuln.Script))
					}
					b.WriteString("\n")
				}
			}
		}
	}

	// Recommendations
	b.WriteString("RECOMMENDED NEXT STEPS\n")
	b.WriteString("======================\n")

	if summary.WebServices > 0 {
		b.WriteString("Web Services Detected:\n")
		b.WriteString("• Run directory enumeration (gobuster, dirb, dirbuster)\n")
		b.WriteString("• Check for default credentials on web interfaces\n")
		b.WriteString("• Test for common web vulnerabilities (SQLi, XSS, etc.)\n")
		b.WriteString("• Examine robots.txt and sitemap.xml\n\n")
	}

	if summary.OpenPorts > 0 {
		b.WriteString("Open Services:\n")
		b.WriteString("• Test services for default/weak credentials\n")
		b.WriteString("• Perform banner grabbing for version information\n")
		b.WriteString("• Check for service-specific vulnerabilities\n")
		b.WriteString("• Attempt service enumeration\n\n")
	}

	if summary.Vulnerabilities > 0 {
		b.WriteString("Vulnerabilities Found:\n")
		b.WriteString("• Research CVEs in exploit databases\n")
		b.WriteString("• Test for exploitability\n")
		b.WriteString("• Check for public exploits (Metasploit, ExploitDB)\n")
		b.WriteString("• Verify patch levels\n\n")
	}

	// Footer
	b.WriteString("==================================================\n")
	b.WriteString("Report generated by IPCrawler\n")
	b.WriteString(fmt.Sprintf("Generated on: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	b.WriteString("==================================================\n")

	return b.String()
}

// SaveResultsMsg represents a save results message
type SaveResultsMsg struct{}

// ExportCompletedMsg represents a completed export
type ExportCompletedMsg struct {
	Filepath string
	Error    string
}
