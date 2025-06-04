package models

import (
	"fmt"
	"time"
)

// ScanResult represents the complete scan result for a target
type ScanResult struct {
	TargetIP       string        `json:"target_ip"`
	StartTime      time.Time     `json:"start_time"`
	EndTime        time.Time     `json:"end_time"`
	Duration       time.Duration `json:"duration"`
	TCPPorts       []Port        `json:"tcp_ports"`
	UDPPorts       []Port        `json:"udp_ports"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
	Summary        ScanSummary   `json:"summary"`
}

// Port represents an open port with its details
type Port struct {
	Number      uint16 `json:"number"`
	Protocol    string `json:"protocol"` // TCP or UDP
	State       string `json:"state"`    // open, closed, filtered
	Service     string `json:"service"`
	Version     string `json:"version"`
	Banner      string `json:"banner"`
	Fingerprint string `json:"fingerprint"`
}

// Vulnerability represents a potential security issue
type Vulnerability struct {
	Port        uint16   `json:"port"`
	Protocol    string   `json:"protocol"`
	Service     string   `json:"service"`
	Severity    string   `json:"severity"`    // Critical, High, Medium, Low
	Type        string   `json:"type"`        // e.g., "Anonymous FTP", "Default Credentials"
	Description string   `json:"description"`
	References  []string `json:"references"`
	Suggestions []string `json:"suggestions"`
}

// ScanSummary provides a high-level overview of scan results
type ScanSummary struct {
	TotalOpenPorts      int `json:"total_open_ports"`
	TCPOpenPorts       int `json:"tcp_open_ports"`
	UDPOpenPorts       int `json:"udp_open_ports"`
	CriticalVulns      int `json:"critical_vulns"`
	HighVulns          int `json:"high_vulns"`
	MediumVulns        int `json:"medium_vulns"`
	LowVulns           int `json:"low_vulns"`
	InterestingServices []string `json:"interesting_services"`
}

// ScanProgress represents the current scanning progress
type ScanProgress struct {
	Phase         string  `json:"phase"`          // "tcp_scan", "udp_scan", "vuln_scan", "complete"
	Progress      float64 `json:"progress"`       // 0.0 to 1.0
	CurrentPort   uint16  `json:"current_port"`
	PortsScanned  int     `json:"ports_scanned"`
	TotalPorts    int     `json:"total_ports"`
	Message       string  `json:"message"`
	ElapsedTime   time.Duration `json:"elapsed_time"`
}

// GetSeverityColor returns a color for vulnerability severity
func (v Vulnerability) GetSeverityColor() string {
	switch v.Severity {
	case "Critical":
		return "#FF0000" // Red
	case "High":
		return "#FF6600" // Orange
	case "Medium":
		return "#FFFF00" // Yellow
	case "Low":
		return "#00FF00" // Green
	default:
		return "#FFFFFF" // White
	}
}

// GetPortDisplay returns a formatted string for port display
func (p Port) GetPortDisplay() string {
	if p.Service != "" {
		if p.Version != "" {
			return fmt.Sprintf("%d/%s (%s %s)", p.Number, p.Protocol, p.Service, p.Version)
		}
		return fmt.Sprintf("%d/%s (%s)", p.Number, p.Protocol, p.Service)
	}
	return fmt.Sprintf("%d/%s", p.Number, p.Protocol)
}

// IsHighValuePort determines if a port is particularly interesting for HTB
func (p Port) IsHighValuePort() bool {
	highValuePorts := map[uint16]bool{
		21:   true, // FTP
		22:   true, // SSH
		23:   true, // Telnet
		25:   true, // SMTP
		53:   true, // DNS
		80:   true, // HTTP
		110:  true, // POP3
		135:  true, // RPC
		139:  true, // NetBIOS
		143:  true, // IMAP
		389:  true, // LDAP
		443:  true, // HTTPS
		445:  true, // SMB
		993:  true, // IMAPS
		995:  true, // POP3S
		1433: true, // MSSQL
		3306: true, // MySQL
		3389: true, // RDP
		5432: true, // PostgreSQL
		5985: true, // WinRM HTTP
		5986: true, // WinRM HTTPS
		8080: true, // HTTP Alt
		8443: true, // HTTPS Alt
	}
	return highValuePorts[p.Number]
} 