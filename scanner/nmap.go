package scanner

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// NmapScanner handles nmap scanning operations
type NmapScanner struct {
	BinaryPath string
	Timeout    time.Duration
	Verbose    bool
}

// ScanResult represents the result of a network scan
type ScanResult struct {
	IP        string            `json:"ip"`
	Hostname  string            `json:"hostname"`
	Ports     []PortResult      `json:"ports"`
	OS        OSDetection       `json:"os"`
	Timestamp time.Time         `json:"timestamp"`
	ScanType  string            `json:"scan_type"`
	Duration  time.Duration     `json:"duration"`
	Raw       string            `json:"raw_output"`
	Metadata  map[string]string `json:"metadata"`
}

// PortResult represents information about a scanned port
type PortResult struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	State    string `json:"state"`
	Service  string `json:"service"`
	Version  string `json:"version"`
	Banner   string `json:"banner"`
	CPE      string `json:"cpe"`
}

// OSDetection represents OS detection results
type OSDetection struct {
	Name       string `json:"name"`
	Accuracy   int    `json:"accuracy"`
	Type       string `json:"type"`
	Vendor     string `json:"vendor"`
	Family     string `json:"family"`
	Generation string `json:"generation"`
}

// ScanOptions defines options for different scan types
type ScanOptions struct {
	ScanType    string
	Ports       string
	Timing      int
	Aggressive  bool
	ServiceScan bool
	OSScan      bool
	ScriptScan  bool
	UDP         bool
	Stealth     bool
	FragScan    bool
	Threads     int
	CustomArgs  []string
}

// NewNmapScanner creates a new nmap scanner instance
func NewNmapScanner() *NmapScanner {
	return &NmapScanner{
		BinaryPath: "/usr/bin/nmap",
		Timeout:    time.Minute * 10,
		Verbose:    false,
	}
}

// IsInstalled checks if nmap is installed and accessible
func (n *NmapScanner) IsInstalled() bool {
	_, err := exec.LookPath("nmap")
	return err == nil
}

// QuickScan performs a quick scan of the top 1000 ports
func (n *NmapScanner) QuickScan(target string) (*ScanResult, error) {
	options := ScanOptions{
		ScanType: "quick",
		Ports:    "--top-ports=1000",
		Timing:   4,
	}
	return n.Scan(target, options)
}

// FullPortScan scans all 65535 TCP ports
func (n *NmapScanner) FullPortScan(target string) (*ScanResult, error) {
	options := ScanOptions{
		ScanType: "full",
		Ports:    "-p-",
		Timing:   4,
	}
	return n.Scan(target, options)
}

// StealthScan performs a TCP SYN stealth scan
func (n *NmapScanner) StealthScan(target string) (*ScanResult, error) {
	options := ScanOptions{
		ScanType: "stealth",
		Ports:    "--top-ports=1000",
		Timing:   3,
		Stealth:  true,
	}
	return n.Scan(target, options)
}

// ServiceScan performs service version detection
func (n *NmapScanner) ServiceScan(target string, ports string) (*ScanResult, error) {
	options := ScanOptions{
		ScanType:    "service",
		Ports:       fmt.Sprintf("-p %s", ports),
		ServiceScan: true,
		Timing:      4,
	}
	return n.Scan(target, options)
}

// OSScan performs OS detection
func (n *NmapScanner) OSScan(target string) (*ScanResult, error) {
	options := ScanOptions{
		ScanType: "os",
		Ports:    "--top-ports=1000",
		OSScan:   true,
		Timing:   4,
	}
	return n.Scan(target, options)
}

// VulnScan performs vulnerability scanning using NSE scripts
func (n *NmapScanner) VulnScan(target string) (*ScanResult, error) {
	options := ScanOptions{
		ScanType:   "vuln",
		Ports:      "--top-ports=1000",
		ScriptScan: true,
		Timing:     3,
		CustomArgs: []string{"--script=vuln"},
	}
	return n.Scan(target, options)
}

// WebScan scans common web ports with HTTP-specific checks
func (n *NmapScanner) WebScan(target string) (*ScanResult, error) {
	options := ScanOptions{
		ScanType:    "web",
		Ports:       "-p 80,443,8080,8443,8000,8888,9000",
		ServiceScan: true,
		ScriptScan:  true,
		Timing:      4,
		CustomArgs:  []string{"--script=http-enum,http-headers,http-methods,http-title"},
	}
	return n.Scan(target, options)
}

// SSHScan scans and analyzes SSH services
func (n *NmapScanner) SSHScan(target string) (*ScanResult, error) {
	options := ScanOptions{
		ScanType:    "ssh",
		Ports:       "-p 22",
		ServiceScan: true,
		ScriptScan:  true,
		Timing:      4,
		CustomArgs:  []string{"--script=ssh-auth-methods,ssh-hostkey,ssh2-enum-algos"},
	}
	return n.Scan(target, options)
}

// Scan performs a network scan with the given options
func (n *NmapScanner) Scan(target string, options ScanOptions) (*ScanResult, error) {
	startTime := time.Now()

	// Build nmap command
	args := n.buildNmapArgs(target, options)

	// Execute nmap command
	cmd := exec.Command("nmap", args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("nmap execution failed: %v", err)
	}

	// Parse results
	result := &ScanResult{
		IP:        target,
		Timestamp: startTime,
		ScanType:  options.ScanType,
		Duration:  time.Since(startTime),
		Raw:       string(output),
		Ports:     []PortResult{},
		Metadata:  make(map[string]string),
	}

	// Parse nmap output
	if err := n.parseNmapOutput(result, string(output)); err != nil {
		return result, fmt.Errorf("failed to parse nmap output: %v", err)
	}

	return result, nil
}

// buildNmapArgs constructs the command line arguments for nmap
func (n *NmapScanner) buildNmapArgs(target string, options ScanOptions) []string {
	args := []string{}

	// Basic scan options
	if options.Stealth {
		args = append(args, "-sS") // TCP SYN scan
	}

	if options.UDP {
		args = append(args, "-sU") // UDP scan
	}

	if options.ServiceScan {
		args = append(args, "-sV") // Version detection
	}

	if options.OSScan {
		args = append(args, "-O") // OS detection
	}

	if options.ScriptScan {
		args = append(args, "-sC") // Default scripts
	}

	if options.Aggressive {
		args = append(args, "-A") // Aggressive scan
	}

	if options.FragScan {
		args = append(args, "-f") // Fragment packets
	}

	// Timing template
	if options.Timing > 0 && options.Timing <= 5 {
		args = append(args, fmt.Sprintf("-T%d", options.Timing))
	}

	// Port specification
	if options.Ports != "" {
		if strings.HasPrefix(options.Ports, "-p") {
			args = append(args, options.Ports)
		} else if strings.HasPrefix(options.Ports, "--") {
			args = append(args, options.Ports)
		} else {
			args = append(args, "-p", options.Ports)
		}
	}

	// Thread count
	if options.Threads > 0 {
		args = append(args, fmt.Sprintf("--max-parallelism=%d", options.Threads))
	}

	// Custom arguments
	args = append(args, options.CustomArgs...)

	// Output options
	args = append(args, "-oN", "-") // Normal output to stdout
	args = append(args, "--open")   // Only show open ports

	// Target
	args = append(args, target)

	return args
}

// parseNmapOutput parses nmap output and populates the scan result
func (n *NmapScanner) parseNmapOutput(result *ScanResult, output string) error {
	scanner := bufio.NewScanner(strings.NewReader(output))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Parse hostname
		if strings.Contains(line, "Nmap scan report for") {
			parts := strings.Split(line, " ")
			if len(parts) >= 4 {
				result.Hostname = parts[4]
			}
		}

		// Parse open ports
		if strings.Contains(line, "/tcp") || strings.Contains(line, "/udp") {
			port := n.parsePortLine(line)
			if port != nil {
				result.Ports = append(result.Ports, *port)
			}
		}

		// Parse OS information
		if strings.Contains(line, "OS details:") {
			osInfo := strings.TrimPrefix(line, "OS details: ")
			result.OS.Name = osInfo
		}

		// Parse service info
		if strings.Contains(line, "Service Info:") {
			// Extract additional service information
			serviceInfo := strings.TrimPrefix(line, "Service Info: ")
			result.Metadata["service_info"] = serviceInfo
		}
	}

	return nil
}

// parsePortLine parses a single port line from nmap output
func (n *NmapScanner) parsePortLine(line string) *PortResult {
	// Example: "22/tcp   open  ssh     OpenSSH 8.2p1 Ubuntu 4ubuntu0.5"
	fields := strings.Fields(line)
	if len(fields) < 3 {
		return nil
	}

	// Parse port and protocol
	portProto := strings.Split(fields[0], "/")
	if len(portProto) != 2 {
		return nil
	}

	port, err := strconv.Atoi(portProto[0])
	if err != nil {
		return nil
	}

	result := &PortResult{
		Port:     port,
		Protocol: portProto[1],
		State:    fields[1],
	}

	// Parse service name
	if len(fields) > 2 {
		result.Service = fields[2]
	}

	// Parse version information
	if len(fields) > 3 {
		result.Version = strings.Join(fields[3:], " ")
	}

	return result
}

// GetScanTypes returns available scan types
func (n *NmapScanner) GetScanTypes() []string {
	return []string{
		"Quick Scan",
		"Full Port Scan",
		"Stealth Scan",
		"Service Detection",
		"OS Detection",
		"Vulnerability Scan",
		"Web Scan (Port 80/443)",
		"SSH Scan (Port 22)",
	}
}

// ValidateTarget checks if the target is a valid IP address or hostname
func (n *NmapScanner) ValidateTarget(target string) bool {
	// Basic validation - could be enhanced
	return target != "" && !strings.Contains(target, " ")
}
