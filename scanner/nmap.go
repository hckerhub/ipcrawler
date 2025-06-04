package scanner

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ProgressCallback defines a function type for progress updates
type ProgressCallback func(stage int, toolName string, status string, complete bool)

// NmapScanner handles nmap scanning operations
type NmapScanner struct {
	BinaryPath       string
	Timeout          time.Duration
	Verbose          bool
	ProgressCallback ProgressCallback
	ToolAvailability *ToolAvailability
}

// ScanResult represents the result of a network scan
type ScanResult struct {
	IP              string            `json:"ip"`
	Hostname        string            `json:"hostname"`
	Ports           []PortResult      `json:"ports"`
	OpenPorts       []PortResult      `json:"open_ports"`
	FilteredPorts   []PortResult      `json:"filtered_ports"`
	ClosedPorts     []PortResult      `json:"closed_ports"`
	Services        []ServiceResult   `json:"services"`
	WebResults      []WebResult       `json:"web_results"`
	Vulnerabilities []VulnResult      `json:"vulnerabilities"`
	Commands        []CommandInfo     `json:"commands"`
	OS              OSDetection       `json:"os"`
	Timestamp       time.Time         `json:"timestamp"`
	ScanType        string            `json:"scan_type"`
	Duration        time.Duration     `json:"duration"`
	Raw             string            `json:"raw_output"`
	Metadata        map[string]string `json:"metadata"`
}

// PortResult represents information about a scanned port
type PortResult struct {
	Port     int      `json:"port"`
	Protocol string   `json:"protocol"`
	State    string   `json:"state"`
	Service  string   `json:"service"`
	Version  string   `json:"version"`
	Banner   string   `json:"banner"`
	CPE      string   `json:"cpe"`
	Scripts  []string `json:"scripts"`
}

// ServiceResult represents detailed service information
type ServiceResult struct {
	Port       int    `json:"port"`
	Service    string `json:"service"`
	Version    string `json:"version"`
	Product    string `json:"product"`
	ExtraInfo  string `json:"extra_info"`
	OSType     string `json:"os_type"`
	Method     string `json:"method"`
	Confidence int    `json:"confidence"`
}

// WebResult represents web discovery results
type WebResult struct {
	URL             string            `json:"url"`
	Port            int               `json:"port"`
	Title           string            `json:"title"`
	Server          string            `json:"server"`
	Technologies    []string          `json:"technologies"`
	StatusCode      int               `json:"status_code"`
	ContentLength   int               `json:"content_length"`
	ResponseTime    int64             `json:"response_time_ms"`
	Paths           []string          `json:"paths"`
	Subdomains      []string          `json:"subdomains"`
	Headers         map[string]string `json:"headers"`
	Cookies         []string          `json:"cookies"`
	SecurityHeaders map[string]string `json:"security_headers"`
	Redirects       []string          `json:"redirects"`
	Forms           int               `json:"forms_count"`
	JavaScriptFiles []string          `json:"js_files"`
	CSSFiles        []string          `json:"css_files"`
	Images          []string          `json:"images"`
	ExternalLinks   []string          `json:"external_links"`
	EmailAddresses  []string          `json:"email_addresses"`
}

// VulnResult represents vulnerability information
type VulnResult struct {
	Port        int    `json:"port"`
	Service     string `json:"service"`
	CVE         string `json:"cve"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
	Script      string `json:"script"`
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

// CommandInfo represents information about executed commands
type CommandInfo struct {
	Tool      string    `json:"tool"`
	Command   string    `json:"command"`
	Args      []string  `json:"args"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Duration  string    `json:"duration"`
	ExitCode  int       `json:"exit_code"`
	Output    string    `json:"output"`
	Error     string    `json:"error"`
	Stage     string    `json:"stage"`
}

// ScanResultTracker wraps a scanner and result to implement CommandTracker
type ScanResultTracker struct {
	scanner *NmapScanner
	result  *ScanResult
}

// TrackCommand implements the CommandTracker interface
func (s *ScanResultTracker) TrackCommand(tool, command string, args []string, startTime, endTime time.Time, exitCode int, output, error, stage string) {
	s.scanner.trackCommand(s.result, tool, command, args, startTime, endTime, exitCode, output, error, stage)
}

// NewNmapScanner creates a new nmap scanner instance
func NewNmapScanner() *NmapScanner {
	return &NmapScanner{
		BinaryPath:       "/usr/bin/nmap",
		Timeout:          time.Minute * 10,
		Verbose:          false,
		ToolAvailability: NewToolAvailability(),
	}
}

// SetProgressCallback sets the progress callback function
func (n *NmapScanner) SetProgressCallback(callback ProgressCallback) {
	n.ProgressCallback = callback
}

// updateProgress sends a progress update if callback is set
func (n *NmapScanner) updateProgress(stage int, toolName string, status string, complete bool) {
	if n.ProgressCallback != nil {
		n.ProgressCallback(stage, toolName, status, complete)
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

// AggressiveScan performs aggressive scanning with max speed and all features
func (n *NmapScanner) AggressiveScan(target string, ports string) (*ScanResult, error) {
	var portSpec string
	if ports == "" {
		portSpec = "-p-" // All ports
	} else {
		portSpec = fmt.Sprintf("-p %s", ports)
	}

	options := ScanOptions{
		ScanType:    "aggressive",
		Ports:       portSpec,
		Aggressive:  true,
		ServiceScan: true,
		OSScan:      true,
		ScriptScan:  true,
		Timing:      5,
		CustomArgs:  []string{"--min-rate=10000", "-sV", "-sC", "-O", "--version-all"},
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
	scanStartTime := time.Now()

	// Build nmap command
	args := n.buildNmapArgs(target, options)

	// Initialize result early so we can track commands
	result := &ScanResult{
		IP:        target,
		Timestamp: scanStartTime,
		ScanType:  options.ScanType,
		Ports:     []PortResult{},
		Metadata:  make(map[string]string),
		Commands:  []CommandInfo{},
	}

	// Add tool availability information to metadata
	n.ToolAvailability.CheckAllTools()
	result.Metadata["available_tools"] = strings.Join(n.ToolAvailability.GetAvailableTools(), ",")
	result.Metadata["missing_tools"] = strings.Join(n.ToolAvailability.GetMissingTools(), ",")

	// Add capability information
	capabilities := n.ToolAvailability.GetCapabilities()
	for capability, available := range capabilities {
		result.Metadata["capability_"+capability] = fmt.Sprintf("%t", available)
	}

	// Debug: Store the command being executed
	result.Metadata["nmap_command"] = "nmap " + strings.Join(args, " ")

	// Execute nmap command with tracking
	nmapStartTime := time.Now()
	cmd := exec.Command("nmap", args...)
	output, err := cmd.Output()
	nmapEndTime := time.Now()

	// Track the nmap command execution
	exitCode := 0
	errorMsg := ""
	if err != nil {
		exitCode = 1
		errorMsg = err.Error()
	}

	n.trackCommand(result, "nmap", "nmap", args, nmapStartTime, nmapEndTime, exitCode, string(output), errorMsg, "Port Scanning")

	if err != nil {
		result.Duration = time.Since(scanStartTime)
		return result, fmt.Errorf("nmap execution failed: %v", err)
	}

	// Parse nmap output
	result.Raw = string(output)
	if err := n.parseNmapOutput(result, string(output)); err != nil {
		result.Duration = time.Since(scanStartTime)
		return result, fmt.Errorf("failed to parse nmap output: %v", err)
	}

	// If aggressive scan, perform additional discovery
	if options.Aggressive {
		n.performAdditionalDiscovery(result, target)
	}

	result.Duration = time.Since(scanStartTime)
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

	// Only show open ports for specific scan types to avoid missing filtered/closed ports
	// For HTB and CTF environments, we want to see filtered ports as they may indicate services
	if options.ScanType == "quick" || options.ScanType == "web" {
		// Keep --open for quick scans to reduce noise
		args = append(args, "--open")
	}
	// For other scan types (full, aggressive, etc.), show all port states

	// Target
	args = append(args, target)

	return args
}

// parseNmapOutput parses nmap output and populates the scan result
func (n *NmapScanner) parseNmapOutput(result *ScanResult, output string) error {
	scanner := bufio.NewScanner(strings.NewReader(output))

	// Debug: Store the raw output for troubleshooting
	if result.Metadata == nil {
		result.Metadata = make(map[string]string)
	}
	result.Metadata["raw_nmap_output"] = output

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Parse hostname
		if strings.Contains(line, "Nmap scan report for") {
			parts := strings.Split(line, " ")
			if len(parts) >= 4 {
				result.Hostname = parts[4]
			}
		}

		// Parse ports (both TCP and UDP, all states)
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

	// Debug: Add parsing statistics to metadata
	result.Metadata["total_ports_found"] = fmt.Sprintf("%d", len(result.Ports))

	// Count ports by state for debugging
	var openCount, filteredCount, closedCount int
	for _, port := range result.Ports {
		switch port.State {
		case "open":
			openCount++
		case "filtered":
			filteredCount++
		case "closed":
			closedCount++
		}
	}
	result.Metadata["open_ports_count"] = fmt.Sprintf("%d", openCount)
	result.Metadata["filtered_ports_count"] = fmt.Sprintf("%d", filteredCount)
	result.Metadata["closed_ports_count"] = fmt.Sprintf("%d", closedCount)

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

// performAdditionalDiscovery performs additional discovery for aggressive scans
func (n *NmapScanner) performAdditionalDiscovery(result *ScanResult, target string) {
	// Initialize additional fields
	result.OpenPorts = []PortResult{}
	result.FilteredPorts = []PortResult{}
	result.ClosedPorts = []PortResult{}
	result.Services = []ServiceResult{}
	result.WebResults = []WebResult{}
	result.Vulnerabilities = []VulnResult{}

	// Categorize ports by state
	for _, port := range result.Ports {
		switch port.State {
		case "open":
			result.OpenPorts = append(result.OpenPorts, port)
		case "filtered":
			result.FilteredPorts = append(result.FilteredPorts, port)
		case "closed":
			result.ClosedPorts = append(result.ClosedPorts, port)
		}
	}

	// Enhanced service detection for open ports
	n.updateProgress(2, "nmap", "Performing service detection", false)
	for _, port := range result.OpenPorts {
		if port.Service != "" {
			service := ServiceResult{
				Port:       port.Port,
				Service:    port.Service,
				Version:    port.Version,
				Product:    extractProduct(port.Version),
				ExtraInfo:  extractExtraInfo(port.Version),
				OSType:     extractOSInfo(port.Version),
				Method:     "nmap",
				Confidence: 80, // Default confidence
			}
			result.Services = append(result.Services, service)
		}
	}
	n.updateProgress(2, "nmap", "Service detection complete", true)

	// Reverse DNS lookup is now handled within the reconnaissance phase

	// Perform reconnaissance for both domains and IPs (before web discovery)
	n.updateProgress(1, "", "Starting reconnaissance phase", false)
	n.performReconnaissance(result, target)
	n.updateProgress(1, "", "Reconnaissance complete", true)

	// Web discovery for HTTP/HTTPS services
	// Include both open and filtered ports for web discovery (HTB machines often have filtered web ports)
	allRelevantPorts := append(result.OpenPorts, result.FilteredPorts...)
	webPorts := GetWebPorts(allRelevantPorts)
	if len(webPorts) > 0 {
		n.updateProgress(3, "", "Starting web discovery", false)
		webEngine := NewWebDiscoveryEngine()
		// Create a command tracker that has access to the current result
		commandTracker := &ScanResultTracker{scanner: n, result: result}
		webResults, err := webEngine.DiscoverWebServicesWithCallback(target, webPorts, n.updateWebProgress, commandTracker)
		if err == nil {
			result.WebResults = webResults

			// Enhanced web technology analysis for discovered URLs
			if len(webResults) > 0 {
				n.updateProgress(4, "", "Starting advanced web analysis", false)
				n.performAdvancedWebAnalysisWithCallback(result, webResults)
				n.updateProgress(4, "", "Advanced web analysis complete", true)
			}
		}
		n.updateProgress(3, "", "Web discovery complete", true)
	}

	// Basic vulnerability detection based on services
	n.updateProgress(5, "nmap", "Analyzing vulnerabilities", false)
	result.Vulnerabilities = n.detectBasicVulnerabilities(result.OpenPorts)
	n.updateProgress(5, "nmap", "Vulnerability analysis complete", true)
}

// updateWebProgress is a callback for web discovery progress
func (n *NmapScanner) updateWebProgress(toolName, status string, complete bool) {
	n.updateProgress(3, toolName, status, complete)
}

// updateAdvancedWebProgress is a callback for advanced web analysis progress
func (n *NmapScanner) updateAdvancedWebProgress(toolName, status string, complete bool) {
	n.updateProgress(4, toolName, status, complete)
}

// updateReconProgress is a callback for reconnaissance progress
func (n *NmapScanner) updateReconProgress(toolName, status string, complete bool) {
	n.updateProgress(1, toolName, status, complete)
}

// extractProduct extracts product name from version string
func extractProduct(version string) string {
	if version == "" {
		return ""
	}

	// Simple extraction - take first word that looks like a product
	parts := strings.Fields(version)
	for _, part := range parts {
		if len(part) > 2 && !strings.Contains(part, ".") {
			return part
		}
	}

	return ""
}

// extractExtraInfo extracts additional info from version string
func extractExtraInfo(version string) string {
	if version == "" {
		return ""
	}

	// Look for parenthetical information
	if strings.Contains(version, "(") && strings.Contains(version, ")") {
		start := strings.Index(version, "(")
		end := strings.Index(version, ")")
		if end > start {
			return version[start+1 : end]
		}
	}

	return ""
}

// extractOSInfo extracts OS information from version string
func extractOSInfo(version string) string {
	version = strings.ToLower(version)

	if strings.Contains(version, "ubuntu") {
		return "Ubuntu"
	}
	if strings.Contains(version, "debian") {
		return "Debian"
	}
	if strings.Contains(version, "centos") {
		return "CentOS"
	}
	if strings.Contains(version, "rhel") || strings.Contains(version, "red hat") {
		return "Red Hat"
	}
	if strings.Contains(version, "windows") {
		return "Windows"
	}
	if strings.Contains(version, "freebsd") {
		return "FreeBSD"
	}
	if strings.Contains(version, "openbsd") {
		return "OpenBSD"
	}

	return ""
}

// detectBasicVulnerabilities performs basic vulnerability detection
func (n *NmapScanner) detectBasicVulnerabilities(ports []PortResult) []VulnResult {
	var vulns []VulnResult

	for _, port := range ports {
		service := strings.ToLower(port.Service)
		version := strings.ToLower(port.Version)

		// Check for known vulnerable services
		switch port.Port {
		case 21: // FTP
			if strings.Contains(version, "vsftpd 2.3.4") {
				vulns = append(vulns, VulnResult{
					Port:        port.Port,
					Service:     port.Service,
					CVE:         "CVE-2011-2523",
					Severity:    "HIGH",
					Description: "vsftpd 2.3.4 backdoor command execution",
					Script:      "ftp-vsftpd-backdoor",
				})
			}
		case 22: // SSH
			if strings.Contains(version, "openssh") {
				// Check for common SSH vulnerabilities
				if strings.Contains(version, "6.6") {
					vulns = append(vulns, VulnResult{
						Port:        port.Port,
						Service:     port.Service,
						CVE:         "CVE-2016-0777",
						Severity:    "MEDIUM",
						Description: "OpenSSH client information leak",
						Script:      "ssh-auth-methods",
					})
				}
			}
		case 23: // Telnet
			vulns = append(vulns, VulnResult{
				Port:        port.Port,
				Service:     port.Service,
				CVE:         "",
				Severity:    "HIGH",
				Description: "Telnet service detected - unencrypted communication",
				Script:      "telnet-encryption",
			})
		case 139, 445: // SMB
			vulns = append(vulns, VulnResult{
				Port:        port.Port,
				Service:     port.Service,
				CVE:         "CVE-2017-0144",
				Severity:    "CRITICAL",
				Description: "Potential EternalBlue vulnerability (MS17-010)",
				Script:      "smb-vuln-ms17-010",
			})
		case 1433: // MSSQL
			vulns = append(vulns, VulnResult{
				Port:        port.Port,
				Service:     port.Service,
				CVE:         "",
				Severity:    "MEDIUM",
				Description: "MSSQL service exposed - check for weak authentication",
				Script:      "ms-sql-info",
			})
		case 3306: // MySQL
			vulns = append(vulns, VulnResult{
				Port:        port.Port,
				Service:     port.Service,
				CVE:         "",
				Severity:    "MEDIUM",
				Description: "MySQL service exposed - check for weak authentication",
				Script:      "mysql-info",
			})
		case 5432: // PostgreSQL
			vulns = append(vulns, VulnResult{
				Port:        port.Port,
				Service:     port.Service,
				CVE:         "",
				Severity:    "MEDIUM",
				Description: "PostgreSQL service exposed - check for weak authentication",
				Script:      "pgsql-brute",
			})
		}

		// Check for default credentials based on service
		if strings.Contains(service, "http") && (port.Port == 80 || port.Port == 443 || port.Port == 8080) {
			vulns = append(vulns, VulnResult{
				Port:        port.Port,
				Service:     port.Service,
				CVE:         "",
				Severity:    "LOW",
				Description: "Web service detected - check for default admin credentials",
				Script:      "http-default-accounts",
			})
		}
	}

	return vulns
}

// performAdvancedWebAnalysisWithCallback runs advanced web technology detection with progress callbacks
func (n *NmapScanner) performAdvancedWebAnalysisWithCallback(result *ScanResult, webResults []WebResult) {
	analyzer := NewAdvancedWebAnalyzer()

	// Extract URLs from web results
	var urls []string
	for _, web := range webResults {
		urls = append(urls, web.URL)
	}

	// Run advanced analysis with callback
	techResults, err := analyzer.AnalyzeWebTechnologiesWithCallback(urls, n.updateAdvancedWebProgress)
	if err != nil {
		return
	}

	// Merge technology information back into web results
	for i, web := range result.WebResults {
		for _, tech := range techResults {
			if tech.URL == web.URL {
				// Convert TechnologyInfo to strings for existing structure
				var techStrings []string
				for _, t := range tech.CombinedTechs {
					techStr := t.Name
					if t.Version != "" {
						techStr += " " + t.Version
					}
					techStrings = append(techStrings, techStr)
				}
				result.WebResults[i].Technologies = techStrings
				break
			}
		}
	}
}

// performReconnaissance runs reconnaissance tools for domain and IP targets
func (n *NmapScanner) performReconnaissance(result *ScanResult, target string) {
	reconEngine := NewReconEngine("./logs")

	// Load API keys from environment variables
	securityTrailsKey := os.Getenv("SECURITYTRAILS_API_KEY")
	censysID := os.Getenv("CENSYS_API_ID")
	censysSecret := os.Getenv("CENSYS_API_SECRET")

	reconEngine.SetAPIKeys(securityTrailsKey, censysID, censysSecret)

	// Get list of installed tools to report what's available
	installedTools := reconEngine.GetInstalledTools()

	// Track which tools we have available
	if result.Metadata == nil {
		result.Metadata = make(map[string]string)
	}
	result.Metadata["recon_tools_available"] = strings.Join(installedTools, ",")

	// Show API key status (without revealing keys)
	var apiStatus []string
	if securityTrailsKey != "" {
		apiStatus = append(apiStatus, "securitytrails")
	}
	if censysID != "" && censysSecret != "" {
		apiStatus = append(apiStatus, "censys")
	}
	result.Metadata["recon_api_keys"] = strings.Join(apiStatus, ",")

	// Create a command tracker that has access to the current result
	commandTracker := &ScanResultTracker{scanner: n, result: result}
	reconResult, err := reconEngine.PerformReconnaissanceWithCallbackAndTracker(target, n.updateReconProgress, commandTracker)
	if err != nil {
		// Don't fail the entire scan, but log the error
		result.Metadata["recon_error"] = fmt.Sprintf("Reconnaissance failed: %v", err)
		return
	}

	// Add reconnaissance data to metadata
	result.Metadata["recon_subdomains"] = fmt.Sprintf("%d", len(reconResult.Subdomains))
	result.Metadata["recon_certificates"] = fmt.Sprintf("%d", len(reconResult.Certificates))
	result.Metadata["recon_sources"] = "amass,recon-ng,crt.sh,securitytrails,censys,netcraft"
	result.Metadata["recon_total_hosts"] = fmt.Sprintf("%d", reconResult.TotalHosts)
}

// isDomain checks if target is a domain name
func (n *NmapScanner) isDomain(target string) bool {
	// Simple check for IP vs domain
	ipRegex := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
	return !ipRegex.MatchString(target)
}

// trackCommand adds a command to the result's command log
func (n *NmapScanner) trackCommand(result *ScanResult, tool, command string, args []string, startTime, endTime time.Time, exitCode int, output, error, stage string) {
	if result.Commands == nil {
		result.Commands = []CommandInfo{}
	}

	duration := endTime.Sub(startTime).Round(time.Millisecond).String()

	commandInfo := CommandInfo{
		Tool:      tool,
		Command:   command,
		Args:      args,
		StartTime: startTime,
		EndTime:   endTime,
		Duration:  duration,
		ExitCode:  exitCode,
		Output:    output,
		Error:     error,
		Stage:     stage,
	}

	result.Commands = append(result.Commands, commandInfo)
}
