package scanner

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Ullaakut/nmap/v3"

	"ipcrawler/internal/models"
)

type Engine struct {
	target        string
	aggressive    bool
	verbose       bool
	interactive   bool
	privilegeLevel PrivilegeLevel
	progress      chan models.ScanProgress
}

func NewEngine(target string, aggressive, verbose, interactive bool) *Engine {
	return &Engine{
		target:      target,
		aggressive:  aggressive,
		verbose:     verbose,
		interactive: interactive,
		progress:    make(chan models.ScanProgress, 100),
	}
}

func (e *Engine) GetProgressChannel() <-chan models.ScanProgress {
	return e.progress
}

// SetPrivilegeLevel sets the privilege level for the engine
func (e *Engine) SetPrivilegeLevel(level PrivilegeLevel) {
	e.privilegeLevel = level
}

func (e *Engine) StartFullScan(ctx context.Context) (*models.ScanResult, error) {
	startTime := time.Now()
	result := &models.ScanResult{
		TargetIP:  e.target,
		StartTime: startTime,
	}

	// Check privileges for UDP scanning (only if not already set)
	if e.privilegeLevel == 0 { // Not set yet
		var err error
		e.privilegeLevel, err = CheckPrivileges(e.interactive)
		if err != nil {
			return nil, fmt.Errorf("privilege check failed: %w", err)
		}
	}

	// Phase 1: TCP Scan
	e.sendProgress("tcp_scan", 0.0, "Starting TCP port scan...")
	tcpPorts, err := e.scanTCPPorts(ctx)
	if err != nil {
		return nil, fmt.Errorf("TCP scan failed: %w", err)
	}
	result.TCPPorts = tcpPorts

	// Phase 2: UDP Scan (conditional based on privileges)
	var udpPorts []models.Port
	if e.privilegeLevel == Privileged {
		e.sendProgress("udp_scan", 0.33, "Starting UDP port scan...")
		udpPorts, err = e.scanUDPPorts(ctx)
		if err != nil {
			// Don't fail completely if UDP scan fails, just log and continue
			if e.verbose {
				e.sendProgress("udp_scan", 0.5, fmt.Sprintf("UDP scan failed: %v", err))
			}
			udpPorts = []models.Port{} // Empty UDP results
		}
	} else {
		e.sendProgress("udp_scan", 0.5, "Skipping UDP scan (requires elevated privileges)")
		udpPorts = []models.Port{} // Empty UDP results
	}
	result.UDPPorts = udpPorts

	// Phase 3: Vulnerability Analysis
	e.sendProgress("vuln_scan", 0.66, "Analyzing vulnerabilities...")
	vulns := e.analyzeVulnerabilities(result.TCPPorts, result.UDPPorts)
	result.Vulnerabilities = vulns

	// Complete scan
	endTime := time.Now()
	result.EndTime = endTime
	result.Duration = endTime.Sub(startTime)
	result.Summary = e.generateSummary(result)

	e.sendProgress("complete", 1.0, "Scan complete!")
	close(e.progress)

	return result, nil
}

func (e *Engine) scanTCPPorts(ctx context.Context) ([]models.Port, error) {
	var ports []models.Port
	
	// Build options for NMAP scanner
	options := []nmap.Option{
		nmap.WithTargets(e.target),
		nmap.WithPorts("1-65535"),
		nmap.WithServiceInfo(),
	}
	
	// Add privilege-based options
	if e.privilegeLevel == Privileged {
		// Running with privileges - can use SYN scan and OS detection
		options = append(options, nmap.WithSYNScan())
		if e.aggressive {
			options = append(options, nmap.WithOSDetection())
		}
	} else {
		// Running without privileges - use connect scan
		options = append(options, nmap.WithConnectScan())
	}
	
	if e.aggressive {
		// Add aggressive scan options
		options = append(options, nmap.WithAggressiveScan())
	}
	
	// Build NMAP scanner for TCP
	scanner, err := nmap.NewScanner(ctx, options...)
	
	if err != nil {
		return nil, fmt.Errorf("unable to create nmap scanner: %w", err)
	}

	// Execute scan
	result, warnings, err := scanner.Run()
	if err != nil {
		return nil, fmt.Errorf("nmap scan failed: %w", err)
	}

	if e.verbose && len(*warnings) > 0 {
		for _, warning := range *warnings {
			e.sendProgress("tcp_scan", 0.2, fmt.Sprintf("Warning: %s", warning))
		}
	}

	// Process results
	for _, host := range result.Hosts {
		for _, port := range host.Ports {
			if port.State.State == "open" {
				p := models.Port{
					Number:      port.ID,
					Protocol:    strings.ToUpper(port.Protocol),
					State:       port.State.State,
					Service:     port.Service.Name,
					Version:     port.Service.Version,
					Banner:      port.Service.Product,
					Fingerprint: port.Service.ExtraInfo,
				}
				ports = append(ports, p)
			}
		}
	}

	return ports, nil
}

func (e *Engine) scanUDPPorts(ctx context.Context) ([]models.Port, error) {
	var ports []models.Port
	
	// Scan only top UDP ports for speed
	topUDPPorts := "53,67,68,69,123,135,137,138,139,161,162,445,631,1434,1900,4500,5353"
	
	scanner, err := nmap.NewScanner(
		ctx,
		nmap.WithTargets(e.target),
		nmap.WithPorts(topUDPPorts),
		nmap.WithUDPScan(),
		nmap.WithServiceInfo(),
	)
	
	if err != nil {
		return nil, fmt.Errorf("unable to create UDP nmap scanner: %w", err)
	}

	// Execute scan
	result, warnings, err := scanner.Run()
	if err != nil {
		return nil, fmt.Errorf("UDP nmap scan failed: %w", err)
	}

	if e.verbose && len(*warnings) > 0 {
		for _, warning := range *warnings {
			e.sendProgress("udp_scan", 0.5, fmt.Sprintf("Warning: %s", warning))
		}
	}

	// Process results
	for _, host := range result.Hosts {
		for _, port := range host.Ports {
			if port.State.State == "open" || port.State.State == "open|filtered" {
				p := models.Port{
					Number:      port.ID,
					Protocol:    "UDP",
					State:       port.State.State,
					Service:     port.Service.Name,
					Version:     port.Service.Version,
					Banner:      port.Service.Product,
					Fingerprint: port.Service.ExtraInfo,
				}
				ports = append(ports, p)
			}
		}
	}

	return ports, nil
}

func (e *Engine) analyzeVulnerabilities(tcpPorts, udpPorts []models.Port) []models.Vulnerability {
	var vulns []models.Vulnerability
	
	allPorts := append(tcpPorts, udpPorts...)
	
	for _, port := range allPorts {
		// Analyze each port for common vulnerabilities
		portVulns := e.analyzePortVulnerabilities(port)
		vulns = append(vulns, portVulns...)
	}
	
	return vulns
}

func (e *Engine) analyzePortVulnerabilities(port models.Port) []models.Vulnerability {
	var vulns []models.Vulnerability
	
	switch port.Number {
	case 21: // FTP
		vulns = append(vulns, models.Vulnerability{
			Port:        port.Number,
			Protocol:    port.Protocol,
			Service:     port.Service,
			Severity:    "High",
			Type:        "Anonymous FTP",
			Description: "FTP service detected. Check for anonymous login and weak credentials.",
			References:  []string{"CVE-2010-4221", "CWE-287"},
			Suggestions: []string{"Try anonymous login", "Test common credentials", "Check for directory traversal"},
		})
		
	case 22: // SSH
		vulns = append(vulns, models.Vulnerability{
			Port:        port.Number,
			Protocol:    port.Protocol,
			Service:     port.Service,
			Severity:    "Medium",
			Type:        "SSH Service",
			Description: "SSH service detected. Potential for brute force attacks.",
			References:  []string{"CWE-307"},
			Suggestions: []string{"Check SSH version", "Test weak passwords", "Look for key-based auth"},
		})
		
	case 80, 8080, 8000, 8888: // HTTP
		vulns = append(vulns, models.Vulnerability{
			Port:        port.Number,
			Protocol:    port.Protocol,
			Service:     port.Service,
			Severity:    "Medium",
			Type:        "HTTP Service",
			Description: "HTTP web service detected. Potential for web-based attacks.",
			References:  []string{"OWASP-2021"},
			Suggestions: []string{"Directory enumeration", "Check for admin panels", "Test for SQLi/XSS", "Scan with dirb/gobuster"},
		})
		
	case 443, 8443: // HTTPS
		vulns = append(vulns, models.Vulnerability{
			Port:        port.Number,
			Protocol:    port.Protocol,
			Service:     port.Service,
			Severity:    "Medium",
			Type:        "HTTPS Service",
			Description: "HTTPS web service detected. Check SSL configuration and web vulnerabilities.",
			References:  []string{"OWASP-2021"},
			Suggestions: []string{"SSL/TLS scan", "Certificate analysis", "Directory enumeration", "Web app testing"},
		})
		
	case 445: // SMB
		vulns = append(vulns, models.Vulnerability{
			Port:        port.Number,
			Protocol:    port.Protocol,
			Service:     port.Service,
			Severity:    "High",
			Type:        "SMB Service",
			Description: "SMB service detected. High potential for lateral movement and privilege escalation.",
			References:  []string{"MS17-010", "CVE-2017-0144"},
			Suggestions: []string{"Check for EternalBlue", "Enumerate shares", "Test null sessions", "SMB version scan"},
		})
		
	case 139: // NetBIOS
		vulns = append(vulns, models.Vulnerability{
			Port:        port.Number,
			Protocol:    port.Protocol,
			Service:     port.Service,
			Severity:    "Medium",
			Type:        "NetBIOS Service",
			Description: "NetBIOS service detected. Information disclosure and enumeration possible.",
			References:  []string{"CWE-200"},
			Suggestions: []string{"NetBIOS enumeration", "Check for null sessions", "Enumerate users/shares"},
		})
		
	case 3389: // RDP
		vulns = append(vulns, models.Vulnerability{
			Port:        port.Number,
			Protocol:    port.Protocol,
			Service:     port.Service,
			Severity:    "High",
			Type:        "RDP Service",
			Description: "RDP service detected. High value target for credential attacks.",
			References:  []string{"CVE-2019-0708", "BlueKeep"},
			Suggestions: []string{"Check for BlueKeep", "Test weak credentials", "RDP enumeration"},
		})
		
	case 1433: // MSSQL
		vulns = append(vulns, models.Vulnerability{
			Port:        port.Number,
			Protocol:    port.Protocol,
			Service:     port.Service,
			Severity:    "High",
			Type:        "MSSQL Service",
			Description: "Microsoft SQL Server detected. Database access and privilege escalation potential.",
			References:  []string{"CWE-89"},
			Suggestions: []string{"Test default credentials", "SQL injection", "xp_cmdshell", "Database enumeration"},
		})
		
	case 3306: // MySQL
		vulns = append(vulns, models.Vulnerability{
			Port:        port.Number,
			Protocol:    port.Protocol,
			Service:     port.Service,
			Severity:    "High",
			Type:        "MySQL Service",
			Description: "MySQL database service detected. Potential for data access and code execution.",
			References:  []string{"CWE-89"},
			Suggestions: []string{"Test weak credentials", "MySQL enumeration", "UDF exploitation"},
		})
		
	case 53: // DNS
		vulns = append(vulns, models.Vulnerability{
			Port:        port.Number,
			Protocol:    port.Protocol,
			Service:     port.Service,
			Severity:    "Low",
			Type:        "DNS Service",
			Description: "DNS service detected. Information gathering and zone transfer possibilities.",
			References:  []string{"CWE-200"},
			Suggestions: []string{"Zone transfer test", "DNS enumeration", "Subdomain discovery"},
		})
	}
	
	return vulns
}

func (e *Engine) generateSummary(result *models.ScanResult) models.ScanSummary {
	summary := models.ScanSummary{
		TCPOpenPorts: len(result.TCPPorts),
		UDPOpenPorts: len(result.UDPPorts),
	}
	
	summary.TotalOpenPorts = summary.TCPOpenPorts + summary.UDPOpenPorts
	
	// Count vulnerabilities by severity
	for _, vuln := range result.Vulnerabilities {
		switch vuln.Severity {
		case "Critical":
			summary.CriticalVulns++
		case "High":
			summary.HighVulns++
		case "Medium":
			summary.MediumVulns++
		case "Low":
			summary.LowVulns++
		}
	}
	
	// Identify interesting services
	var services []string
	for _, port := range result.TCPPorts {
		if port.IsHighValuePort() && port.Service != "" {
			services = append(services, fmt.Sprintf("%s:%d", port.Service, port.Number))
		}
	}
	summary.InterestingServices = services
	
	return summary
}

func (e *Engine) sendProgress(phase string, progress float64, message string) {
	select {
	case e.progress <- models.ScanProgress{
		Phase:     phase,
		Progress:  progress,
		Message:   message,
		ElapsedTime: time.Since(time.Now()), // This will be updated by the caller
	}:
	default:
		// Non-blocking send
	}
} 