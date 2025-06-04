package scanner

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"
)

// ReconEngine handles comprehensive reconnaissance operations
type ReconEngine struct {
	SecurityTrailsKey string
	CensysAPIID       string
	CensysAPISecret   string
	MaxConcurrency    int
	Timeout           time.Duration
	Verbose           bool
	OutputDir         string
	CommandTracker    CommandTracker
	ToolChecker       ToolAvailabilityChecker
}

// ToolAvailabilityChecker interface for checking tool availability
type ToolAvailabilityChecker interface {
	CheckToolStatus(toolName string) (bool, string)
	GetMissingToolsForCategory(category string) []string
}

// ReconResult represents comprehensive reconnaissance results
type ReconResult struct {
	Target          string               `json:"target"`
	Type            string               `json:"type"` // IP or domain
	IPs             []IPInfo             `json:"ips"`
	Domains         []DomainInfo         `json:"domains"`
	Subdomains      []SubdomainInfo      `json:"subdomains"`
	Certificates    []CertificateInfo    `json:"certificates"`
	CensysResults   []CensysResult       `json:"censys_results"`
	SecurityTrails  SecurityTrailsResult `json:"securitytrails"`
	NetcraftResults []NetcraftResult     `json:"netcraft_results"`
	PassiveDNS      []PassiveDNSRecord   `json:"passive_dns"`
	RelatedIPs      []string             `json:"related_ips"`
	TotalHosts      int                  `json:"total_hosts"`
	Timestamp       time.Time            `json:"timestamp"`
}

// IPInfo represents detailed IP information
type IPInfo struct {
	IP           string              `json:"ip"`
	ISP          string              `json:"isp"`
	Organization string              `json:"organization"`
	Country      string              `json:"country"`
	City         string              `json:"city"`
	ASN          string              `json:"asn"`
	Ports        []PortInfo          `json:"ports"`
	Services     []ServiceInfo       `json:"services"`
	Vulns        []VulnerabilityInfo `json:"vulnerabilities"`
	Source       string              `json:"source"`
	LastSeen     time.Time           `json:"last_seen"`
}

// DomainInfo represents domain information
type DomainInfo struct {
	Domain      string    `json:"domain"`
	Registrar   string    `json:"registrar"`
	CreatedDate time.Time `json:"created_date"`
	ExpiryDate  time.Time `json:"expiry_date"`
	NameServers []string  `json:"nameservers"`
	Status      string    `json:"status"`
	Source      string    `json:"source"`
}

// CertificateInfo represents SSL certificate information
type CertificateInfo struct {
	CommonName   string    `json:"common_name"`
	Issuer       string    `json:"issuer"`
	NotBefore    time.Time `json:"not_before"`
	NotAfter     time.Time `json:"not_after"`
	SANs         []string  `json:"sans"`
	SerialNumber string    `json:"serial_number"`
	Source       string    `json:"source"`
}

// CensysResult represents Censys scan results
type CensysResult struct {
	IP       string                 `json:"ip"`
	Services []CensysService        `json:"services"`
	Location CensysLocation         `json:"location"`
	ASN      CensysASN              `json:"asn"`
	Metadata map[string]interface{} `json:"metadata"`
}

// CensysService represents service information from Censys
type CensysService struct {
	Port              int                `json:"port"`
	ServiceName       string             `json:"service_name"`
	TransportProtocol string             `json:"transport_protocol"`
	Certificate       *CensusCertificate `json:"certificate,omitempty"`
	Banner            string             `json:"banner"`
	Software          []CensusSoftware   `json:"software"`
}

// CensysLocation represents location data from Censys
type CensysLocation struct {
	Country     string     `json:"country"`
	CountryCode string     `json:"country_code"`
	City        string     `json:"city"`
	Province    string     `json:"province"`
	Timezone    string     `json:"timezone"`
	Coordinates [2]float64 `json:"coordinates"`
}

// CensysASN represents ASN information from Censys
type CensysASN struct {
	ASN          int    `json:"asn"`
	Description  string `json:"description"`
	CountryCode  string `json:"country_code"`
	Organization string `json:"organization"`
}

// CensusCertificate represents certificate data from Censys
type CensusCertificate struct {
	Parsed CensusCertParsed `json:"parsed"`
}

// CensusCertParsed represents parsed certificate data
type CensusCertParsed struct {
	Subject            CensusCertSubject   `json:"subject"`
	Issuer             CensusCertIssuer    `json:"issuer"`
	SubjectAltName     CensusCertSANs      `json:"subject_alt_name"`
	Validity           CensusCertValidity  `json:"validity"`
	SerialNumber       string              `json:"serial_number"`
	SignatureAlgorithm CensusCertSignature `json:"signature_algorithm"`
}

// CensusCertSubject represents certificate subject
type CensusCertSubject struct {
	CommonName []string `json:"common_name"`
}

// CensusCertIssuer represents certificate issuer
type CensusCertIssuer struct {
	CommonName []string `json:"common_name"`
}

// CensusCertSANs represents Subject Alternative Names
type CensusCertSANs struct {
	DNSNames []string `json:"dns_names"`
}

// CensusCertValidity represents certificate validity period
type CensusCertValidity struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// CensusCertSignature represents signature algorithm
type CensusCertSignature struct {
	Name string `json:"name"`
}

// CensusSoftware represents software information
type CensusSoftware struct {
	Product string `json:"product"`
	Vendor  string `json:"vendor"`
	Version string `json:"version"`
}

// SecurityTrailsResult represents SecurityTrails API results
type SecurityTrailsResult struct {
	Subdomains     []string           `json:"subdomains"`
	PassiveDNS     []PassiveDNSRecord `json:"passive_dns"`
	WhoisHistory   []WhoisRecord      `json:"whois_history"`
	SSLCertHistory []SSLCertRecord    `json:"ssl_cert_history"`
}

// PassiveDNSRecord represents passive DNS information
type PassiveDNSRecord struct {
	Hostname  string    `json:"hostname"`
	IP        string    `json:"ip"`
	Type      string    `json:"type"`
	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
	Source    string    `json:"source"`
}

// WhoisRecord represents WHOIS information
type WhoisRecord struct {
	Registrar   string    `json:"registrar"`
	CreatedDate time.Time `json:"created_date"`
	ExpiryDate  time.Time `json:"expiry_date"`
	NameServers []string  `json:"nameservers"`
}

// SSLCertRecord represents SSL certificate record
type SSLCertRecord struct {
	SHA1        string    `json:"sha1"`
	FirstSeen   time.Time `json:"first_seen"`
	LastSeen    time.Time `json:"last_seen"`
	CommonNames []string  `json:"common_names"`
}

// NetcraftResult represents Netcraft results
type NetcraftResult struct {
	Domain       string    `json:"domain"`
	IP           string    `json:"ip"`
	NetBlock     string    `json:"netblock"`
	ISP          string    `json:"isp"`
	Country      string    `json:"country"`
	Server       string    `json:"server"`
	Technologies []string  `json:"technologies"`
	Risk         string    `json:"risk"`
	LastSeen     time.Time `json:"last_seen"`
}

// PortInfo represents port information
type PortInfo struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Service  string `json:"service"`
	Banner   string `json:"banner"`
	State    string `json:"state"`
}

// ServiceInfo represents service information
type ServiceInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Product string `json:"product"`
	CPE     string `json:"cpe"`
}

// VulnerabilityInfo represents vulnerability information
type VulnerabilityInfo struct {
	CVE        string   `json:"cve"`
	Summary    string   `json:"summary"`
	CVSS       float64  `json:"cvss"`
	References []string `json:"references"`
	Verified   bool     `json:"verified"`
}

// NewReconEngine creates a new reconnaissance engine
func NewReconEngine(outputDir string) *ReconEngine {
	return &ReconEngine{
		MaxConcurrency: 5,
		Timeout:        time.Minute * 2,
		Verbose:        true,
		OutputDir:      outputDir,
		ToolChecker:    NewSimpleToolChecker(),
	}
}

// SetAPIKeys configures API keys for various services
func (r *ReconEngine) SetAPIKeys(securityTrailsKey, censysID, censysSecret string) {
	r.SecurityTrailsKey = securityTrailsKey
	r.CensysAPIID = censysID
	r.CensysAPISecret = censysSecret
}

// PerformReconnaissance conducts comprehensive reconnaissance
func (r *ReconEngine) PerformReconnaissance(target string) (*ReconResult, error) {
	result := &ReconResult{
		Target:    target,
		Type:      r.determineTargetType(target),
		Timestamp: time.Now(),
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, 10)

	// Run reconnaissance tools in parallel - each tool runs independently
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := r.runAmass(target, result); err != nil && r.Verbose {
			mu.Lock()
			errChan <- fmt.Errorf("Amass: %v", err)
			mu.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := r.runReconNG(target, result); err != nil && r.Verbose {
			mu.Lock()
			errChan <- fmt.Errorf("Recon-ng: %v", err)
			mu.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := r.runCrtSh(target, result); err != nil && r.Verbose {
			mu.Lock()
			errChan <- fmt.Errorf("crt.sh: %v", err)
			mu.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := r.runSecurityTrails(target, result); err != nil && r.Verbose {
			mu.Lock()
			errChan <- fmt.Errorf("SecurityTrails: %v", err)
			mu.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := r.runCensys(target, result); err != nil && r.Verbose {
			mu.Lock()
			errChan <- fmt.Errorf("Censys: %v", err)
			mu.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := r.runNetcraft(target, result); err != nil && r.Verbose {
			mu.Lock()
			errChan <- fmt.Errorf("Netcraft: %v", err)
			mu.Unlock()
		}
	}()

	// Wait for all reconnaissance tools to complete
	wg.Wait()
	close(errChan)

	// Collect any errors (for logging, don't fail the scan)
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	// Post-process and deduplicate results
	r.postProcessResults(result)

	return result, nil
}

// determineTargetType determines if target is IP or domain
func (r *ReconEngine) determineTargetType(target string) string {
	// Simple IP regex check
	ipRegex := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
	if ipRegex.MatchString(target) {
		return "ip"
	}
	return "domain"
}

// runAmass executes Amass for subdomain enumeration
func (r *ReconEngine) runAmass(target string, result *ReconResult) error {
	// Check if amass is available using the tool checker
	if r.ToolChecker != nil {
		if available, message := r.ToolChecker.CheckToolStatus("amass"); !available {
			// Track the command as not available
			if r.CommandTracker != nil {
				r.CommandTracker.TrackCommand("amass", "amass", []string{"enum", "-d", target, "-passive", "-json"}, time.Now(), time.Now(), 1, "", fmt.Sprintf("Tool not available: %s", message), "Reconnaissance")
			}
			return fmt.Errorf("amass not available: %s", message)
		}
	} else if !r.isToolInstalled("amass") {
		return fmt.Errorf("amass not installed")
	}

	if result.Type != "domain" {
		// Track command for IP targets to show it was skipped
		if r.CommandTracker != nil {
			r.CommandTracker.TrackCommand("amass", "amass", []string{"enum", "-d", target, "-passive", "-json"}, time.Now(), time.Now(), 0, "", "Skipped for IP target (domains only)", "Reconnaissance")
		}
		return nil
	}

	args := []string{"enum", "-d", target, "-passive", "-json"}
	cmd := exec.Command("amass", args...)

	startTime := time.Now()
	output, err := cmd.Output()
	endTime := time.Now()

	exitCode := 0
	errorOutput := ""
	if err != nil {
		exitCode = 1
		errorOutput = err.Error()
	}

	// Track command execution
	if r.CommandTracker != nil {
		r.CommandTracker.TrackCommand("amass", "amass", args, startTime, endTime, exitCode, string(output), errorOutput, "Reconnaissance")
	}

	if err != nil {
		return err
	}

	r.parseAmassOutput(string(output), result)
	return nil
}

// runReconNG executes Recon-ng modules
func (r *ReconEngine) runReconNG(target string, result *ReconResult) error {
	// Check if recon-ng is available using the tool checker
	if r.ToolChecker != nil {
		if available, message := r.ToolChecker.CheckToolStatus("recon-ng"); !available {
			// Track the command as not available
			if r.CommandTracker != nil {
				r.CommandTracker.TrackCommand("recon-ng", "recon-ng", []string{"-w", "ipcrawler", "-x", "workspaces"}, time.Now(), time.Now(), 1, "", fmt.Sprintf("Tool not available: %s", message), "Reconnaissance")
			}
			return fmt.Errorf("recon-ng not available: %s", message)
		}
	} else if !r.isToolInstalled("recon-ng") {
		return fmt.Errorf("recon-ng not installed")
	}

	// Create temporary workspace and run modules
	workspace := fmt.Sprintf("ipcrawler_%d", time.Now().Unix())

	// Create workspace
	cmd := exec.Command("recon-ng", "-w", workspace, "-m", "workspaces", "-x", fmt.Sprintf("workspaces create %s", workspace))
	if err := cmd.Run(); err != nil {
		return err
	}

	// Add domain to workspace
	if result.Type == "domain" {
		cmd = exec.Command("recon-ng", "-w", workspace, "-x", fmt.Sprintf("db insert domains %s", target))
		cmd.Run()

		// Run subdomain modules
		modules := []string{
			"recon/domains-subdomains/certificate_transparency",
			"recon/domains-subdomains/threatcrowd",
		}

		for _, module := range modules {
			cmd = exec.Command("recon-ng", "-w", workspace, "-m", module, "-x", "run")
			output, _ := cmd.Output()
			r.parseReconNGOutput(string(output), result)
		}
	}

	return nil
}

// runCrtSh queries crt.sh for certificate transparency data
func (r *ReconEngine) runCrtSh(target string, result *ReconResult) error {
	if result.Type != "domain" {
		// Track command for IP targets to show it was skipped
		if r.CommandTracker != nil {
			r.CommandTracker.TrackCommand("crt.sh", "curl", []string{fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", target)}, time.Now(), time.Now(), 0, "", "Skipped for IP target (domains only)", "Reconnaissance")
		}
		return nil
	}

	url := fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", target)

	client := &http.Client{Timeout: r.Timeout}

	startTime := time.Now()
	resp, err := client.Get(url)
	endTime := time.Now()

	exitCode := 0
	errorOutput := ""
	output := ""

	if err != nil {
		exitCode = 1
		errorOutput = err.Error()
	} else {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			exitCode = 1
			errorOutput = err.Error()
		} else {
			output = string(body)
			r.parseCrtShOutput(output, result)
		}
	}

	// Track command execution
	if r.CommandTracker != nil {
		r.CommandTracker.TrackCommand("crt.sh", "curl", []string{url}, startTime, endTime, exitCode, output, errorOutput, "Reconnaissance")
	}

	return err
}

// runSecurityTrails queries SecurityTrails API
func (r *ReconEngine) runSecurityTrails(target string, result *ReconResult) error {
	if r.SecurityTrailsKey == "" {
		return fmt.Errorf("SecurityTrails API key not configured")
	}

	if result.Type != "domain" {
		return nil
	}

	client := &http.Client{Timeout: r.Timeout}

	// Get subdomains
	url := fmt.Sprintf("https://api.securitytrails.com/v1/domain/%s/subdomains", target)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("APIKEY", r.SecurityTrailsKey)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	r.parseSecurityTrailsOutput(string(body), result)
	return nil
}

// runCensys queries Censys API
func (r *ReconEngine) runCensys(target string, result *ReconResult) error {
	if r.CensysAPIID == "" || r.CensysAPISecret == "" {
		return fmt.Errorf("Censys API credentials not configured")
	}

	// Check if censys CLI is available
	if !r.isToolInstalled("censys") {
		return fmt.Errorf("censys CLI not installed")
	}

	var cmd *exec.Cmd
	if result.Type == "ip" {
		cmd = exec.Command("censys", "search", "hosts", target)
	} else {
		cmd = exec.Command("censys", "search", "certificates", fmt.Sprintf("names:%s", target))
	}

	cmd.Env = append(os.Environ(),
		fmt.Sprintf("CENSYS_API_ID=%s", r.CensysAPIID),
		fmt.Sprintf("CENSYS_API_SECRET=%s", r.CensysAPISecret))

	output, err := cmd.Output()
	if err != nil {
		return err
	}

	r.parseCensysOutput(string(output), result)
	return nil
}

// runNetcraft performs Netcraft reconnaissance
func (r *ReconEngine) runNetcraft(target string, result *ReconResult) error {
	// Netcraft doesn't have official CLI, so we'll use curl with their search
	url := fmt.Sprintf("https://searchdns.netcraft.com/?restriction=site+contains&host=%s", target)

	cmd := exec.Command("curl", "-s", "-A", "Mozilla/5.0", url)
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	r.parseNetcraftOutput(string(output), result)
	return nil
}

// isToolInstalled checks if a tool is installed
func (r *ReconEngine) isToolInstalled(tool string) bool {
	_, err := exec.LookPath(tool)
	return err == nil
}

// Parser functions for each tool's output

func (r *ReconEngine) parseAmassOutput(output string, result *ReconResult) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && strings.Contains(line, ".") {
			subdomain := SubdomainInfo{
				Subdomain:   line,
				LastChecked: time.Now(),
			}
			result.Subdomains = append(result.Subdomains, subdomain)
		}
	}
}

func (r *ReconEngine) parseReconNGOutput(output string, result *ReconResult) {
	// Parse recon-ng output for subdomains and other data
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "[+]") && strings.Contains(line, ".") {
			// Extract subdomain from recon-ng output
			parts := strings.Fields(line)
			for _, part := range parts {
				if strings.Contains(part, ".") && !strings.Contains(part, "[") {
					subdomain := SubdomainInfo{
						Subdomain:   part,
						LastChecked: time.Now(),
					}
					result.Subdomains = append(result.Subdomains, subdomain)
				}
			}
		}
	}
}

func (r *ReconEngine) parseCrtShOutput(output string, result *ReconResult) {
	var crtData []map[string]interface{}
	if err := json.Unmarshal([]byte(output), &crtData); err != nil {
		return
	}

	for _, cert := range crtData {
		if nameValue, ok := cert["name_value"].(string); ok {
			domains := strings.Split(nameValue, "\n")
			for _, domain := range domains {
				domain = strings.TrimSpace(domain)
				if domain != "" {
					certInfo := CertificateInfo{
						CommonName: domain,
						Source:     "crt.sh",
					}
					if issuer, ok := cert["issuer_name"].(string); ok {
						certInfo.Issuer = issuer
					}
					result.Certificates = append(result.Certificates, certInfo)

					// Also add as subdomain if it's a subdomain
					if strings.Contains(domain, ".") && domain != result.Target {
						subdomain := SubdomainInfo{
							Subdomain:   domain,
							LastChecked: time.Now(),
						}
						result.Subdomains = append(result.Subdomains, subdomain)
					}
				}
			}
		}
	}
}

func (r *ReconEngine) parseSecurityTrailsOutput(output string, result *ReconResult) {
	var stData map[string]interface{}
	if err := json.Unmarshal([]byte(output), &stData); err != nil {
		return
	}

	if subdomains, ok := stData["subdomains"].([]interface{}); ok {
		for _, sub := range subdomains {
			if subdomain, ok := sub.(string); ok {
				fullDomain := fmt.Sprintf("%s.%s", subdomain, result.Target)
				subInfo := SubdomainInfo{
					Subdomain:   fullDomain,
					LastChecked: time.Now(),
				}
				result.Subdomains = append(result.Subdomains, subInfo)
			}
		}
	}
}

func (r *ReconEngine) parseCensysOutput(output string, result *ReconResult) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "IP:") {
			// Parse Censys output for IP information
			censysResult := CensysResult{}
			// Basic parsing logic
			if r.extractCensysData(line, &censysResult) {
				result.CensysResults = append(result.CensysResults, censysResult)
			}
		}
	}
}

func (r *ReconEngine) parseNetcraftOutput(output string, result *ReconResult) {
	// Parse HTML output from Netcraft search
	re := regexp.MustCompile(`href="/site_report\?url=([^"]+)"`)
	matches := re.FindAllStringSubmatch(output, -1)

	for _, match := range matches {
		if len(match) > 1 {
			netcraftResult := NetcraftResult{
				Domain:   match[1],
				LastSeen: time.Now(),
			}
			result.NetcraftResults = append(result.NetcraftResults, netcraftResult)
		}
	}
}

// Helper functions for data extraction

func (r *ReconEngine) extractCensysData(line string, result *CensysResult) bool {
	// Basic implementation - would need more sophisticated parsing
	if strings.Contains(line, "IP:") {
		parts := strings.Split(line, " ")
		for i, part := range parts {
			if part == "IP:" && i+1 < len(parts) {
				result.IP = parts[i+1]
				return true
			}
		}
	}
	return false
}

// postProcessResults deduplicates and organizes results
func (r *ReconEngine) postProcessResults(result *ReconResult) {
	// Deduplicate subdomains
	seen := make(map[string]bool)
	var uniqueSubdomains []SubdomainInfo
	for _, sub := range result.Subdomains {
		if !seen[sub.Subdomain] {
			seen[sub.Subdomain] = true
			uniqueSubdomains = append(uniqueSubdomains, sub)
		}
	}
	result.Subdomains = uniqueSubdomains

	// Count total hosts
	result.TotalHosts = len(result.IPs) + len(result.Subdomains)
}

// GetInstalledTools returns list of available reconnaissance tools
func (r *ReconEngine) GetInstalledTools() []string {
	if r.ToolChecker != nil {
		return r.ToolChecker.(*SimpleToolChecker).GetInstalledReconTools()
	}

	// Fallback to old method
	tools := []string{"amass", "recon-ng", "censys", "curl"}
	var installed []string

	for _, tool := range tools {
		if r.isToolInstalled(tool) {
			installed = append(installed, tool)
		}
	}

	return installed
}

// InstallRecommendations provides installation recommendations for missing tools
func (r *ReconEngine) InstallRecommendations() map[string]string {
	recommendations := map[string]string{
		"amass":    "go install -v github.com/owasp-amass/amass/v4/...@master",
		"recon-ng": "pip install recon-ng",
		"censys":   "pip install censys",
		"curl":     "Install via package manager (curl is usually pre-installed)",
	}

	missing := make(map[string]string)
	for tool, command := range recommendations {
		if !r.isToolInstalled(tool) {
			missing[tool] = command
		}
	}

	return missing
}

// GetToolInstallationReport returns a comprehensive tool installation report
func (r *ReconEngine) GetToolInstallationReport() string {
	if r.ToolChecker != nil {
		return r.ToolChecker.(*SimpleToolChecker).GetToolInstallationReport()
	}

	return "Tool checker not available - using basic tool detection"
}

// CheckRequiredTools returns true if all required tools are available
func (r *ReconEngine) CheckRequiredTools() (bool, []string) {
	if r.ToolChecker != nil {
		return r.ToolChecker.(*SimpleToolChecker).CheckCoreTools()
	}

	// Fallback check for core tools
	coreTools := []string{"curl", "nmap"}
	var missing []string

	for _, tool := range coreTools {
		if !r.isToolInstalled(tool) {
			missing = append(missing, tool)
		}
	}

	return len(missing) == 0, missing
}

// runReverseDNSLookup performs reverse DNS discovery for IP addresses
func (r *ReconEngine) runReverseDNSLookup(target string, result *ReconResult) error {
	// Method 1: Standard reverse DNS lookup using dig
	startTime := time.Now()
	cmd := exec.Command("dig", "-x", target, "+short")
	output, err := cmd.Output()
	endTime := time.Now()

	exitCode := 0
	errorMsg := ""
	if err != nil {
		exitCode = 1
		errorMsg = err.Error()
	}

	// Track the dig command
	if r.CommandTracker != nil {
		r.CommandTracker.TrackCommand("dig", "dig", []string{"-x", target, "+short"}, startTime, endTime, exitCode, string(output), errorMsg, "Reconnaissance")
	}

	if err == nil && len(output) > 0 {
		// Process the reverse DNS results
		domains := strings.Fields(strings.TrimSpace(string(output)))
		for _, domain := range domains {
			domain = strings.TrimSuffix(strings.TrimSpace(domain), ".")
			if domain != "" {
				// Add as a discovered domain
				domainInfo := DomainInfo{
					Domain: domain,
					Source: "reverse_dns",
				}
				result.Domains = append(result.Domains, domainInfo)
			}
		}
	}

	// Method 2: PTR record lookup
	startTime = time.Now()
	reversedIP := r.reverseIPAddress(target)
	ptrQuery := reversedIP + ".in-addr.arpa"
	cmd = exec.Command("dig", "-t", "PTR", ptrQuery, "+short")
	output, err = cmd.Output()
	endTime = time.Now()

	exitCode = 0
	errorMsg = ""
	if err != nil {
		exitCode = 1
		errorMsg = err.Error()
	}

	// Track the PTR lookup command
	if r.CommandTracker != nil {
		r.CommandTracker.TrackCommand("dig", "dig", []string{"-t", "PTR", ptrQuery, "+short"}, startTime, endTime, exitCode, string(output), errorMsg, "Reconnaissance")
	}

	if err == nil && len(output) > 0 {
		// Process PTR record results
		ptrDomains := strings.Fields(strings.TrimSpace(string(output)))
		for _, domain := range ptrDomains {
			domain = strings.TrimSuffix(strings.TrimSpace(domain), ".")
			if domain != "" {
				// Check if we already have this domain
				found := false
				for _, existing := range result.Domains {
					if existing.Domain == domain {
						found = true
						break
					}
				}
				if !found {
					domainInfo := DomainInfo{
						Domain: domain,
						Source: "ptr_record",
					}
					result.Domains = append(result.Domains, domainInfo)
				}
			}
		}
	}

	return nil
}

// reverseIPAddress reverses an IP address for PTR lookups
func (r *ReconEngine) reverseIPAddress(ip string) string {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return ip
	}
	return fmt.Sprintf("%s.%s.%s.%s", parts[3], parts[2], parts[1], parts[0])
}
