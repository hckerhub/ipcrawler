package scanner

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

// WebDiscoveryEngine handles comprehensive web service discovery
type WebDiscoveryEngine struct {
	UserAgent          string
	Timeout            time.Duration
	MaxConcurrency     int
	CustomWordlists    []string
	EnableSubdomains   bool
	EnableDirectories  bool
	EnableTechnologies bool
	EnableScreenshots  bool
	Verbose            bool
}

// AdvancedWebResult contains comprehensive web discovery results
type AdvancedWebResult struct {
	URL             string            `json:"url"`
	Port            int               `json:"port"`
	Title           string            `json:"title"`
	Server          string            `json:"server"`
	Technologies    []TechInfo        `json:"technologies"`
	StatusCode      int               `json:"status_code"`
	ContentLength   int               `json:"content_length"`
	ResponseTime    int64             `json:"response_time_ms"`
	Paths           []PathInfo        `json:"paths"`
	Subdomains      []SubdomainInfo   `json:"subdomains"`
	Headers         map[string]string `json:"headers"`
	Cookies         []string          `json:"cookies"`
	Forms           []FormInfo        `json:"forms"`
	JSFiles         []string          `json:"js_files"`
	CSSFiles        []string          `json:"css_files"`
	Images          []string          `json:"images"`
	ExternalLinks   []string          `json:"external_links"`
	EmailAddresses  []string          `json:"email_addresses"`
	SecurityHeaders SecurityHeaders   `json:"security_headers"`
	SSLInfo         SSLInfo           `json:"ssl_info"`
	Vulnerabilities []string          `json:"vulnerabilities"`
	Screenshot      string            `json:"screenshot_path"`
	LastChecked     time.Time         `json:"last_checked"`
}

// TechInfo represents detected technology information
type TechInfo struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	Category   string `json:"category"`
	Confidence int    `json:"confidence"`
	Source     string `json:"source"`
}

// PathInfo represents discovered path information
type PathInfo struct {
	Path         string `json:"path"`
	StatusCode   int    `json:"status_code"`
	Size         int    `json:"size"`
	ResponseTime int64  `json:"response_time_ms"`
	ContentType  string `json:"content_type"`
	Title        string `json:"title"`
	Interesting  bool   `json:"interesting"`
	Reason       string `json:"reason"`
}

// SubdomainInfo represents subdomain discovery results
type SubdomainInfo struct {
	Subdomain   string    `json:"subdomain"`
	IP          string    `json:"ip"`
	Status      string    `json:"status"`
	Title       string    `json:"title"`
	Server      string    `json:"server"`
	LastChecked time.Time `json:"last_checked"`
}

// FormInfo represents discovered form information
type FormInfo struct {
	Action   string            `json:"action"`
	Method   string            `json:"method"`
	Fields   map[string]string `json:"fields"`
	Endpoint string            `json:"endpoint"`
}

// SecurityHeaders represents security-related HTTP headers
type SecurityHeaders struct {
	StrictTransportSecurity string   `json:"strict_transport_security"`
	ContentSecurityPolicy   string   `json:"content_security_policy"`
	XFrameOptions           string   `json:"x_frame_options"`
	XContentTypeOptions     string   `json:"x_content_type_options"`
	XSSProtection           string   `json:"xss_protection"`
	ReferrerPolicy          string   `json:"referrer_policy"`
	Missing                 []string `json:"missing"`
}

// SSLInfo represents SSL/TLS information
type SSLInfo struct {
	Valid       bool     `json:"valid"`
	Issuer      string   `json:"issuer"`
	Subject     string   `json:"subject"`
	ExpiryDate  string   `json:"expiry_date"`
	SANs        []string `json:"sans"`
	Certificate string   `json:"certificate"`
}

// NewWebDiscoveryEngine creates a new enhanced web discovery engine
func NewWebDiscoveryEngine() *WebDiscoveryEngine {
	return &WebDiscoveryEngine{
		UserAgent:          "Mozilla/5.0 (compatible; IPCrawler/2.0; +https://github.com/ipcrawler)",
		Timeout:            time.Second * 5, // Shorter timeout for faster scans
		MaxConcurrency:     3,               // Reduced concurrency for stability
		EnableSubdomains:   false,           // Disable for basic discovery
		EnableDirectories:  true,
		EnableTechnologies: true,
		EnableScreenshots:  false,
		Verbose:            true, // Enable verbose for debugging
	}
}

// DiscoverWebServices performs comprehensive web service discovery
func (w *WebDiscoveryEngine) DiscoverWebServices(target string, ports []PortResult) ([]WebResult, error) {
	var results []WebResult
	var wg sync.WaitGroup
	semaphore := make(chan bool, w.MaxConcurrency)
	resultsChan := make(chan WebResult, len(ports))

	for _, port := range ports {
		if w.isWebPort(port.Port) {
			wg.Add(1)
			go func(p PortResult) {
				defer wg.Done()
				semaphore <- true
				defer func() { <-semaphore }()

				result := w.analyzeWebService(target, p)
				if result != nil {
					resultsChan <- *result
				}
			}(port)
		}
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	for result := range resultsChan {
		results = append(results, result)
	}

	return results, nil
}

// isWebPort checks if a port is commonly used for web services
func (w *WebDiscoveryEngine) isWebPort(port int) bool {
	webPorts := []int{80, 443, 8000, 8008, 8080, 8443, 8888, 9000, 9001, 9080, 9090, 9443, 3000, 3001, 4000, 4001, 5000, 5001, 7001, 7002}
	for _, webPort := range webPorts {
		if port == webPort {
			return true
		}
	}
	return false
}

// analyzeWebService performs comprehensive analysis of a web service
func (w *WebDiscoveryEngine) analyzeWebService(target string, port PortResult) *WebResult {
	baseURL := w.buildURL(target, port.Port)

	result := &WebResult{
		URL:  baseURL,
		Port: port.Port,
	}

	// Basic HTTP analysis
	if err := w.performBasicAnalysis(result); err != nil {
		if w.Verbose {
			// Silent - basic analysis failed
		}
		return nil
	}

	// Technology detection
	if w.EnableTechnologies {
		w.detectTechnologies(result)
	}

	// Directory discovery
	if w.EnableDirectories {
		w.discoverDirectories(result)
	}

	// Subdomain discovery
	if w.EnableSubdomains {
		w.discoverSubdomains(target, result)
	}

	return result
}

// buildURL constructs the appropriate URL scheme
func (w *WebDiscoveryEngine) buildURL(target string, port int) string {
	scheme := "http"
	if port == 443 || port == 8443 || port == 9443 {
		scheme = "https"
	}

	// Don't add port if it's default for the scheme
	if (scheme == "http" && port == 80) || (scheme == "https" && port == 443) {
		return fmt.Sprintf("%s://%s", scheme, target)
	}

	return fmt.Sprintf("%s://%s:%d", scheme, target, port)
}

// performBasicAnalysis conducts basic HTTP analysis
func (w *WebDiscoveryEngine) performBasicAnalysis(result *WebResult) error {
	client := &http.Client{
		Timeout: w.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Store redirect chain for domain analysis
			if len(via) < 10 { // Prevent infinite redirects
				result.Redirects = append(result.Redirects, req.URL.String())
			}
			return nil
		},
	}

	startTime := time.Now()
	req, err := http.NewRequest("GET", result.URL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", w.UserAgent)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	result.ResponseTime = time.Since(startTime).Nanoseconds() / 1000000

	// Extract basic information
	result.Server = resp.Header.Get("Server")
	result.ContentLength = int(resp.ContentLength)

	// Store all headers for analysis
	result.Headers = make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			result.Headers[key] = values[0]
		}
	}

	// Extract cookies
	for _, cookie := range resp.Cookies() {
		result.Cookies = append(result.Cookies, cookie.String())
	}

	// Critical Domain Discovery Analysis
	w.extractDomainsFromResponse(resp, result)

	// Read and analyze body content for domain extraction
	body := make([]byte, 100*1024) // Read first 100KB for analysis
	n, _ := resp.Body.Read(body)
	bodyContent := string(body[:n])

	// Extract title
	if titleMatch := regexp.MustCompile(`<title[^>]*>([^<]+)</title>`).FindStringSubmatch(bodyContent); len(titleMatch) > 1 {
		result.Title = strings.TrimSpace(titleMatch[1])
	}

	// Extract domains from HTML content
	w.extractDomainsFromHTML(bodyContent, result)

	// Extract other web assets
	w.extractWebAssets(bodyContent, result)

	// Always add basic paths for immediate discovery
	result.Paths = w.checkCommonPaths(result.URL)

	// Basic technology detection from headers and content
	result.Technologies = w.manualTechDetection(result)

	return nil
}

// extractDomainsFromResponse extracts domains from HTTP response headers and redirects
func (w *WebDiscoveryEngine) extractDomainsFromResponse(resp *http.Response, result *WebResult) {
	var discoveredDomains []string

	// Extract domains from Location header (redirects)
	if location := resp.Header.Get("Location"); location != "" {
		if domain := w.extractDomainFromURL(location); domain != "" {
			discoveredDomains = append(discoveredDomains, domain)
		}
	}

	// Extract domains from Host header
	if host := resp.Header.Get("Host"); host != "" {
		discoveredDomains = append(discoveredDomains, host)
	}

	// Extract domains from Server header (sometimes contains domain info)
	if server := resp.Header.Get("Server"); server != "" {
		if matches := regexp.MustCompile(`([a-zA-Z0-9-]+\.(?:htb|thm|local|box))`).FindAllString(server, -1); len(matches) > 0 {
			discoveredDomains = append(discoveredDomains, matches...)
		}
	}

	// Extract domains from other common headers
	headerFields := []string{
		"X-Forwarded-Host", "X-Original-Host", "X-Host", "Origin",
		"Referer", "Content-Location", "X-Frame-Options",
	}

	for _, header := range headerFields {
		if value := resp.Header.Get(header); value != "" {
			if domain := w.extractDomainFromURL(value); domain != "" {
				discoveredDomains = append(discoveredDomains, domain)
			}
		}
	}

	// Filter and validate CTF/HTB domains
	validDomains := w.filterCTFDomains(discoveredDomains)
	result.Subdomains = append(result.Subdomains, validDomains...)
}

// extractDomainsFromHTML extracts domains from HTML content
func (w *WebDiscoveryEngine) extractDomainsFromHTML(content string, result *WebResult) {
	var discoveredDomains []string

	// Extract domains from href attributes
	hrefRegex := regexp.MustCompile(`href=['"]([^'"]+)['"]`)
	for _, match := range hrefRegex.FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			if domain := w.extractDomainFromURL(match[1]); domain != "" {
				discoveredDomains = append(discoveredDomains, domain)
			}
		}
	}

	// Extract domains from src attributes
	srcRegex := regexp.MustCompile(`src=['"]([^'"]+)['"]`)
	for _, match := range srcRegex.FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			if domain := w.extractDomainFromURL(match[1]); domain != "" {
				discoveredDomains = append(discoveredDomains, domain)
			}
		}
	}

	// Extract domains from JavaScript
	jsRegex := regexp.MustCompile(`(?:window\.location|document\.location|location\.href)\s*=\s*['"]([^'"]+)['"]`)
	for _, match := range jsRegex.FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			if domain := w.extractDomainFromURL(match[1]); domain != "" {
				discoveredDomains = append(discoveredDomains, domain)
			}
		}
	}

	// Extract domains from forms
	formRegex := regexp.MustCompile(`action=['"]([^'"]+)['"]`)
	for _, match := range formRegex.FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			if domain := w.extractDomainFromURL(match[1]); domain != "" {
				discoveredDomains = append(discoveredDomains, domain)
			}
		}
	}

	// Extract domains from meta tags
	metaRegex := regexp.MustCompile(`<meta[^>]+content=['"]([^'"]*(?:htb|thm|local|box)[^'"]*)['"][^>]*>`)
	for _, match := range metaRegex.FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			if domain := w.extractDomainFromURL(match[1]); domain != "" {
				discoveredDomains = append(discoveredDomains, domain)
			}
		}
	}

	// Filter and validate CTF/HTB domains
	validDomains := w.filterCTFDomains(discoveredDomains)
	result.Subdomains = append(result.Subdomains, validDomains...)
}

// extractDomainFromURL extracts domain from a URL or string
func (w *WebDiscoveryEngine) extractDomainFromURL(input string) string {
	// Handle relative URLs and malformed input
	if input == "" || input == "/" || strings.HasPrefix(input, "/") {
		return ""
	}

	// Parse URL
	if parsedURL, err := url.Parse(input); err == nil && parsedURL.Host != "" {
		return parsedURL.Host
	}

	// Direct domain pattern matching for CTF/HTB domains
	ctfDomainRegex := regexp.MustCompile(`([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*\.(?:htb|thm|local|box))`)
	if matches := ctfDomainRegex.FindStringSubmatch(input); len(matches) > 1 {
		return matches[1]
	}

	return ""
}

// filterCTFDomains filters and validates CTF/HTB specific domains
func (w *WebDiscoveryEngine) filterCTFDomains(domains []string) []string {
	var validDomains []string
	seen := make(map[string]bool)

	ctfPatterns := []string{".htb", ".thm", ".local", ".box"}

	for _, domain := range domains {
		// Skip empty or duplicate domains
		if domain == "" || seen[domain] {
			continue
		}

		// Convert to lowercase for consistency
		domain = strings.ToLower(strings.TrimSpace(domain))

		// Check if it matches CTF patterns
		isValidCTF := false
		for _, pattern := range ctfPatterns {
			if strings.HasSuffix(domain, pattern) {
				isValidCTF = true
				break
			}
		}

		// Skip common public domains
		publicDomains := []string{".com", ".org", ".net", ".edu", ".gov", ".mil", ".int"}
		isPublic := false
		for _, publicDomain := range publicDomains {
			if strings.HasSuffix(domain, publicDomain) {
				isPublic = true
				break
			}
		}

		// Only include CTF domains and exclude public ones
		if isValidCTF && !isPublic {
			seen[domain] = true
			validDomains = append(validDomains, domain)
		}
	}

	return validDomains
}

// extractWebAssets extracts JavaScript, CSS, images, and external links from HTML content
func (w *WebDiscoveryEngine) extractWebAssets(content string, result *WebResult) {
	// Extract JavaScript files
	jsRegex := regexp.MustCompile(`src=['"]([^'"]*\.js[^'"]*)['"]`)
	for _, match := range jsRegex.FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			result.JavaScriptFiles = append(result.JavaScriptFiles, match[1])
		}
	}

	// Extract CSS files
	cssRegex := regexp.MustCompile(`(?:href=['"]([^'"]*\.css[^'"]*)['"]|@import\s+['"]([^'"]*\.css[^'"]*)['"])`)
	for _, match := range cssRegex.FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			for i := 1; i < len(match); i++ {
				if match[i] != "" {
					result.CSSFiles = append(result.CSSFiles, match[i])
				}
			}
		}
	}

	// Extract image files
	imgRegex := regexp.MustCompile(`src=['"]([^'"]*\.(?:jpg|jpeg|png|gif|svg|webp|ico)[^'"]*)['"]`)
	for _, match := range imgRegex.FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			result.Images = append(result.Images, match[1])
		}
	}

	// Extract external links
	linkRegex := regexp.MustCompile(`href=['"]https?://([^'"]+)['"]`)
	for _, match := range linkRegex.FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			result.ExternalLinks = append(result.ExternalLinks, match[1])
		}
	}

	// Extract email addresses
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	for _, match := range emailRegex.FindAllString(content, -1) {
		result.EmailAddresses = append(result.EmailAddresses, match)
	}

	// Count forms
	formRegex := regexp.MustCompile(`<form[^>]*>`)
	result.Forms = len(formRegex.FindAllString(content, -1))

	// Remove duplicates from all slices
	result.JavaScriptFiles = w.removeDuplicates(result.JavaScriptFiles)
	result.CSSFiles = w.removeDuplicates(result.CSSFiles)
	result.Images = w.removeDuplicates(result.Images)
	result.ExternalLinks = w.removeDuplicates(result.ExternalLinks)
	result.EmailAddresses = w.removeDuplicates(result.EmailAddresses)
}

// detectTechnologies uses multiple methods to detect web technologies
func (w *WebDiscoveryEngine) detectTechnologies(result *WebResult) {
	var technologies []string

	// Method 1: Try whatweb
	if techs := w.runWhatweb(result.URL); len(techs) > 0 {
		technologies = append(technologies, techs...)
	}

	// Method 2: Try wappalyzer
	if techs := w.runWappalyzer(result.URL); len(techs) > 0 {
		technologies = append(technologies, techs...)
	}

	// Method 3: Manual detection
	if techs := w.manualTechDetection(result); len(techs) > 0 {
		technologies = append(technologies, techs...)
	}

	// Remove duplicates
	result.Technologies = w.removeDuplicates(technologies)

	// Enhanced: Domain discovery from WhatWeb redirects and headers
	w.runWhatwebForDomains(result.URL, result)
}

// runWhatweb executes whatweb for technology detection
func (w *WebDiscoveryEngine) runWhatweb(url string) []string {
	var technologies []string

	cmd := exec.Command("whatweb", "--no-errors", "-a", "3", "--color=never", url)
	output, err := cmd.Output()
	if err != nil {
		return technologies
	}

	// Parse whatweb output
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, url) {
			// Extract technologies from whatweb output
			parts := strings.Split(line, " ")
			for _, part := range parts {
				if strings.Contains(part, "[") && strings.Contains(part, "]") {
					tech := strings.Trim(part, "[]")
					if tech != "" && !strings.Contains(tech, "Status:") {
						technologies = append(technologies, tech)
					}
				}
			}
		}
	}

	return technologies
}

// runWhatwebForDomains executes whatweb specifically for domain discovery from redirects
func (w *WebDiscoveryEngine) runWhatwebForDomains(url string, result *WebResult) {
	// Check if whatweb is available
	if _, err := exec.LookPath("whatweb"); err != nil {
		return
	}

	cmd := exec.Command("whatweb",
		"--no-errors",
		"-a", "3",
		"--color=never",
		"--log-verbose=-",
		"--follow-redirect=always",
		"--max-redirects=10",
		url)

	output, err := cmd.Output()
	if err != nil {
		return
	}

	// Extract domains from whatweb output including redirects
	var discoveredDomains []string
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		// Look for redirect information in whatweb output
		if strings.Contains(line, "RedirectLocation") || strings.Contains(line, "Location:") || strings.Contains(line, "Redirect:") {
			// Extract domain from redirect location
			if domain := w.extractDomainFromWhatwebLine(line); domain != "" {
				discoveredDomains = append(discoveredDomains, domain)
			}
		}

		// Look for Host headers
		if strings.Contains(line, "Host:") {
			if domain := w.extractDomainFromWhatwebLine(line); domain != "" {
				discoveredDomains = append(discoveredDomains, domain)
			}
		}

		// Look for any CTF domain patterns in the output
		ctfDomainRegex := regexp.MustCompile(`([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*\.(?:htb|thm|local|box))`)
		if matches := ctfDomainRegex.FindAllString(line, -1); len(matches) > 0 {
			discoveredDomains = append(discoveredDomains, matches...)
		}
	}

	// Filter and add CTF domains
	validDomains := w.filterCTFDomains(discoveredDomains)
	result.Subdomains = append(result.Subdomains, validDomains...)
}

// extractDomainFromWhatwebLine extracts domain from whatweb output line
func (w *WebDiscoveryEngine) extractDomainFromWhatwebLine(line string) string {
	// Look for URL patterns in whatweb output
	urlRegex := regexp.MustCompile(`https?://([^/\s]+)`)
	if matches := urlRegex.FindStringSubmatch(line); len(matches) > 1 {
		return w.extractDomainFromURL("http://" + matches[1])
	}

	// Look for domain patterns directly
	domainRegex := regexp.MustCompile(`([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*\.(?:htb|thm|local|box))`)
	if matches := domainRegex.FindStringSubmatch(line); len(matches) > 1 {
		return matches[1]
	}

	return ""
}

// runWappalyzer executes wappalyzer for technology detection
func (w *WebDiscoveryEngine) runWappalyzer(url string) []string {
	var technologies []string

	// Try wappalyzer CLI if available
	cmd := exec.Command("wappalyzer", url)
	output, err := cmd.Output()
	if err != nil {
		return technologies
	}

	// Parse JSON output from wappalyzer
	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err == nil {
		if apps, ok := result["applications"].([]interface{}); ok {
			for _, app := range apps {
				if appMap, ok := app.(map[string]interface{}); ok {
					if name, ok := appMap["name"].(string); ok {
						technologies = append(technologies, name)
					}
				}
			}
		}
	}

	return technologies
}

// manualTechDetection performs manual technology detection
func (w *WebDiscoveryEngine) manualTechDetection(result *WebResult) []string {
	var technologies []string

	client := &http.Client{Timeout: w.Timeout}
	req, err := http.NewRequest("GET", result.URL, nil)
	if err != nil {
		return technologies
	}

	req.Header.Set("User-Agent", w.UserAgent)
	resp, err := client.Do(req)
	if err != nil {
		return technologies
	}
	defer resp.Body.Close()

	body := make([]byte, 16384) // Read first 16KB
	n, _ := resp.Body.Read(body)
	bodyStr := strings.ToLower(string(body[:n]))
	headers := resp.Header

	// Server header analysis
	server := strings.ToLower(result.Server)
	if strings.Contains(server, "apache") {
		technologies = append(technologies, "Apache")
	}
	if strings.Contains(server, "nginx") {
		technologies = append(technologies, "Nginx")
	}
	if strings.Contains(server, "iis") {
		technologies = append(technologies, "IIS")
	}

	// X-Powered-By header
	if powered := headers.Get("X-Powered-By"); powered != "" {
		technologies = append(technologies, powered)
	}

	// Technology-specific patterns in HTML
	patterns := map[string][]string{
		"WordPress": {
			"/wp-content/", "/wp-includes/", "wp-json",
		},
		"Joomla": {
			"/components/com_", "/modules/mod_", "joomla",
		},
		"Drupal": {
			"/sites/default/files/", "drupal", "/core/",
		},
		"React": {
			"react", "_react", "reactdom",
		},
		"Angular": {
			"angular", "ng-", "@angular",
		},
		"Vue.js": {
			"vue", "vue.js", "__vue__",
		},
		"jQuery": {
			"jquery", "$.fn.jquery",
		},
		"Bootstrap": {
			"bootstrap", "bs-", "btn btn-",
		},
		"PHP": {
			"<?php", ".php", "phpsessid",
		},
		"ASP.NET": {
			"__viewstate", "asp.net", "aspnet_sessionid",
		},
		"Laravel": {
			"laravel", "_token", "laravel_session",
		},
		"Django": {
			"django", "csrfmiddlewaretoken",
		},
		"Flask": {
			"flask", "session=.",
		},
		"Spring": {
			"spring", "jsessionid",
		},
	}

	for tech, patternsSlice := range patterns {
		for _, pattern := range patternsSlice {
			if strings.Contains(bodyStr, pattern) {
				technologies = append(technologies, tech)
				break
			}
		}
	}

	return technologies
}

// discoverDirectories performs comprehensive directory discovery
func (w *WebDiscoveryEngine) discoverDirectories(result *WebResult) {
	var allPaths []string
	var mu sync.Mutex
	var wg sync.WaitGroup
	baseURL := result.URL

	// Channel to collect paths from all tools
	pathsChan := make(chan []string, 4) // Buffer for 4 tools

	// Run FFuF concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		ffufPaths := w.runFFUF(baseURL)
		pathsChan <- ffufPaths
	}()

	// Run Feroxbuster concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		feroxPaths := w.runFeroxbuster(baseURL)
		pathsChan <- feroxPaths
	}()

	// Run Gobuster concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		gobusterPaths := w.runGobuster(baseURL)
		pathsChan <- gobusterPaths
	}()

	// Run manual discovery methods concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		var manualPaths []string

		robotsPaths := w.discoverFromRobots(baseURL)
		manualPaths = append(manualPaths, robotsPaths...)

		sitemapPaths := w.discoverFromSitemap(baseURL)
		manualPaths = append(manualPaths, sitemapPaths...)

		commonPaths := w.checkCommonPaths(baseURL)
		manualPaths = append(manualPaths, commonPaths...)

		pathsChan <- manualPaths
	}()

	// Collect results in a separate goroutine
	go func() {
		wg.Wait()
		close(pathsChan)
	}()

	// Collect all paths from all tools
	for paths := range pathsChan {
		mu.Lock()
		allPaths = append(allPaths, paths...)
		mu.Unlock()
	}

	// Remove duplicates
	result.Paths = w.removeDuplicatePaths(allPaths)
}

// runFFUF executes ffuf for directory discovery
func (w *WebDiscoveryEngine) runFFUF(baseURL string) []string {
	var paths []string

	// Use multiple wordlists for comprehensive discovery
	wordlists := []string{
		"/usr/share/seclists/Discovery/Web-Content/common.txt",
		"/usr/share/seclists/Discovery/Web-Content/big.txt",
		"/usr/share/wordlists/dirb/common.txt",
		"/usr/share/wordlists/dirbuster/directory-list-2.3-medium.txt",
	}

	for _, wordlist := range wordlists {
		if _, err := exec.LookPath("ffuf"); err != nil {
			continue
		}

		cmd := exec.Command("ffuf",
			"-u", baseURL+"/FUZZ",
			"-w", wordlist,
			"-ac", // Auto-calibration
			"-c",  // Colorize output
			"-sf", // Stop on spurious errors
			"-H", "User-Agent: "+w.UserAgent,
			"-t", "10", // 10 threads
			"-timeout", "5", // 5 second timeout
			"-mc", "200,301,302,403", // Match these codes
			"-fs", "0", // Filter size 0
			"-o", "-", // Output to stdout
			"-of", "json", // JSON format
		)

		output, err := cmd.Output()
		if err != nil {
			continue
		}

		// Parse ffuf JSON output
		var results map[string]interface{}
		if err := json.Unmarshal(output, &results); err == nil {
			if resultsArray, ok := results["results"].([]interface{}); ok {
				for _, item := range resultsArray {
					if resultMap, ok := item.(map[string]interface{}); ok {
						if input, ok := resultMap["input"].(map[string]interface{}); ok {
							if fuzz, ok := input["FUZZ"].(string); ok {
								paths = append(paths, "/"+fuzz)
							}
						}
					}
				}
			}
		}

		// Limit to prevent excessive results
		if len(paths) > 100 {
			break
		}
	}

	return paths
}

// runFeroxbuster executes feroxbuster for directory discovery
func (w *WebDiscoveryEngine) runFeroxbuster(baseURL string) []string {
	var paths []string

	if _, err := exec.LookPath("feroxbuster"); err != nil {
		return paths
	}

	cmd := exec.Command("feroxbuster",
		"-u", baseURL,
		"-w", "/usr/share/seclists/Discovery/Web-Content/common.txt",
		"-t", "10", // 10 threads
		"-d", "2", // Depth 2
		"-H", "User-Agent: "+w.UserAgent,
		"--timeout", "5", // 5 second timeout
		"--status-codes", "200,301,302,403,401",
		"--silent", // Silent mode
		"--json",   // JSON output
		"-o", "-",  // Output to stdout
	)

	output, err := cmd.Output()
	if err != nil {
		return paths
	}

	// Parse feroxbuster output (line-delimited JSON)
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		var result map[string]interface{}
		if err := json.Unmarshal(scanner.Bytes(), &result); err == nil {
			if urlStr, ok := result["url"].(string); ok {
				if parsedURL, err := url.Parse(urlStr); err == nil {
					paths = append(paths, parsedURL.Path)
				}
			}
		}
	}

	return paths
}

// runGobuster executes gobuster for directory discovery
func (w *WebDiscoveryEngine) runGobuster(baseURL string) []string {
	var paths []string

	if _, err := exec.LookPath("gobuster"); err != nil {
		return paths
	}

	cmd := exec.Command("gobuster", "dir",
		"-u", baseURL,
		"-w", "/usr/share/wordlists/dirb/common.txt",
		"-t", "10", // 10 threads
		"-a", w.UserAgent, // User agent
		"--timeout", "5s", // 5 second timeout
		"-s", "200,301,302,403", // Status codes
		"-q",         // Quiet
		"--no-error", // Don't display errors
	)

	output, err := cmd.Output()
	if err != nil {
		return paths
	}

	// Parse gobuster output
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "/") {
			// Extract path from gobuster output
			parts := strings.Fields(line)
			if len(parts) > 0 {
				paths = append(paths, parts[0])
			}
		}
	}

	return paths
}

// checkCommonPaths performs intelligent URL discovery without wordlists
func (w *WebDiscoveryEngine) checkCommonPaths(baseURL string) []string {
	var foundPaths []string

	if w.Verbose {
		// Starting intelligent URL discovery
	}

	// Phase 1: Extract information from discoverable sources
	foundPaths = append(foundPaths, w.discoverFromRobots(baseURL)...)
	foundPaths = append(foundPaths, w.discoverFromSitemap(baseURL)...)
	foundPaths = append(foundPaths, w.extractLinksFromPage(baseURL)...)

	// Phase 2: Technology-specific discovery
	foundPaths = append(foundPaths, w.discoverTechnologyPaths(baseURL)...)

	// Phase 3: Common paths that often exist
	commonPaths := []string{
		// Administrative interfaces (high priority)
		"/admin", "/admin/", "/administrator", "/admin.php", "/admin.html",
		"/wp-admin", "/wp-admin/", "/phpmyadmin", "/phpmyadmin/",
		"/admin/login", "/admin/index.php", "/adminpanel", "/control",
		"/cpanel", "/controlpanel", "/admin-console", "/management",

		// Authentication endpoints
		"/login", "/login.php", "/login.html", "/signin", "/sign-in",
		"/auth", "/authenticate", "/session", "/oauth", "/sso",
		"/logout", "/logoff", "/account", "/profile",

		// API endpoints (crucial for modern apps)
		"/api", "/api/", "/api/v1", "/api/v2", "/api/v3",
		"/rest", "/rest/", "/graphql", "/graphql/", "/soap",
		"/service", "/services", "/ws", "/webservice", "/endpoint",

		// Configuration and system info
		"/config", "/configuration", "/settings", "/setup",
		"/install", "/installation", "/wizard", "/configure",
		"/system", "/sys", "/status", "/health", "/info",
		"/version", "/about", "/debug", "/diagnostics",
		"/server-status", "/server-info", "/phpinfo.php", "/info.php",

		// File and content management
		"/upload", "/uploads", "/files", "/data", "/backup",
		"/backups", "/download", "/downloads", "/documents",
		"/media", "/assets", "/static", "/public", "/content",

		// Development and testing
		"/test", "/testing", "/dev", "/development", "/staging",
		"/beta", "/alpha", "/demo", "/sample", "/example",
		"/temp", "/tmp", "/cache", "/logs", "/log",

		// Common web directories
		"/images", "/img", "/pics", "/pictures", "/photos",
		"/js", "/javascript", "/css", "/stylesheets", "/fonts",
		"/includes", "/lib", "/library", "/vendor", "/node_modules",

		// Documentation and help
		"/docs", "/documentation", "/help", "/support",
		"/manual", "/guide", "/faq", "/wiki", "/readme",

		// Special files and directories
		"/robots.txt", "/sitemap.xml", "/.htaccess", "/.env",
		"/.git", "/.svn", "/.DS_Store", "/web.config",
		"/crossdomain.xml", "/clientaccesspolicy.xml",

		// Database and monitoring
		"/database", "/db", "/mysql", "/postgresql", "/mongo",
		"/monitor", "/monitoring", "/metrics", "/stats",
		"/analytics", "/reports", "/dashboard",
	}

	foundPaths = append(foundPaths, w.checkPathsExistence(baseURL, commonPaths)...)

	// Silent - discovery complete

	return w.removeDuplicatePaths(foundPaths)
}

// discoverFromRobots extracts paths from robots.txt
func (w *WebDiscoveryEngine) discoverFromRobots(baseURL string) []string {
	var paths []string

	robotsURL := baseURL + "/robots.txt"
	if w.Verbose {
		// Checking robots.txt
	}

	client := &http.Client{Timeout: w.Timeout}
	resp, err := client.Get(robotsURL)
	if err != nil {
		return paths
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		body := make([]byte, 10240) // Read up to 10KB
		n, _ := resp.Body.Read(body)
		content := string(body[:n])

		// Extract paths from Disallow and Allow directives
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "Disallow:") || strings.HasPrefix(line, "Allow:") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					path := strings.TrimSpace(parts[1])
					if path != "" && path != "/" && !strings.Contains(path, "*") {
						paths = append(paths, path)
						// Found path in robots.txt
					}
				}
			}
		}
	}

	return paths
}

// discoverFromSitemap extracts URLs from sitemap.xml
func (w *WebDiscoveryEngine) discoverFromSitemap(baseURL string) []string {
	var paths []string

	sitemapURLs := []string{
		baseURL + "/sitemap.xml",
		baseURL + "/sitemap_index.xml",
		baseURL + "/sitemap.txt",
	}

	for _, sitemapURL := range sitemapURLs {
		if w.Verbose {
			// Checking sitemap
		}

		client := &http.Client{Timeout: w.Timeout}
		resp, err := client.Get(sitemapURL)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			body := make([]byte, 50240) // Read up to 50KB
			n, _ := resp.Body.Read(body)
			content := string(body[:n])

			// Extract URLs from XML or text format
			if strings.Contains(content, "<loc>") {
				// XML format
				locRegex := regexp.MustCompile(`<loc>([^<]+)</loc>`)
				matches := locRegex.FindAllStringSubmatch(content, -1)
				for _, match := range matches {
					if len(match) > 1 {
						fullURL := match[1]
						// Extract path from full URL
						if u, err := url.Parse(fullURL); err == nil {
							if u.Path != "" && u.Path != "/" {
								paths = append(paths, u.Path)
								// Found path in sitemap
							}
						}
					}
				}
			}
			break // Found a valid sitemap
		}
	}

	return paths
}

// extractLinksFromPage extracts links from the main page HTML
func (w *WebDiscoveryEngine) extractLinksFromPage(baseURL string) []string {
	var paths []string

	if w.Verbose {
		// Extracting links from main page
	}

	client := &http.Client{Timeout: w.Timeout}
	resp, err := client.Get(baseURL)
	if err != nil {
		return paths
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return paths
	}

	body := make([]byte, 102400) // Read up to 100KB
	n, _ := resp.Body.Read(body)
	content := string(body[:n])

	// Extract href attributes from links
	hrefRegex := regexp.MustCompile(`href=["']([^"']+)["']`)
	matches := hrefRegex.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) > 1 {
			link := match[1]

			// Skip external links, anchors, and non-http protocols
			if strings.HasPrefix(link, "http") && !strings.Contains(link, baseURL) {
				continue
			}
			if strings.HasPrefix(link, "#") || strings.HasPrefix(link, "javascript:") {
				continue
			}
			if strings.HasPrefix(link, "mailto:") || strings.HasPrefix(link, "tel:") {
				continue
			}

			// Normalize the path
			if strings.HasPrefix(link, "/") {
				paths = append(paths, link)
			} else if !strings.Contains(link, "://") {
				// Relative path
				if !strings.HasPrefix(link, "/") {
					link = "/" + link
				}
				paths = append(paths, link)
			}
		}
	}

	// Also extract src attributes for resources
	srcRegex := regexp.MustCompile(`src=["']([^"']+)["']`)
	srcMatches := srcRegex.FindAllStringSubmatch(content, -1)

	for _, match := range srcMatches {
		if len(match) > 1 {
			src := match[1]
			if strings.HasPrefix(src, "/") && !strings.Contains(src, "data:") {
				paths = append(paths, src)
			}
		}
	}

	if w.Verbose && len(paths) > 0 {
		// Extracted links from page
	}

	return paths
}

// checkPathsExistence verifies which paths actually exist and are accessible
func (w *WebDiscoveryEngine) checkPathsExistence(baseURL string, paths []string) []string {
	var foundPaths []string

	client := &http.Client{
		Timeout: w.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Don't follow redirects for discovery
		},
	}

	semaphore := make(chan bool, 5) // Limit concurrent requests
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, path := range paths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			semaphore <- true
			defer func() { <-semaphore }()

			fullURL := baseURL + p

			req, err := http.NewRequest("GET", fullURL, nil)
			if err != nil {
				return
			}

			req.Header.Set("User-Agent", w.UserAgent)

			resp, err := client.Do(req)
			if err != nil {
				return
			}
			resp.Body.Close()

			// Consider interesting if not 404 or 405
			if resp.StatusCode != 404 && resp.StatusCode != 405 {
				mu.Lock()
				foundPaths = append(foundPaths, p)
				// Silent mode - path found
				mu.Unlock()
			}
		}(path)
	}

	wg.Wait()
	return foundPaths
}

// discoverTechnologyPaths discovers paths based on detected technologies
func (w *WebDiscoveryEngine) discoverTechnologyPaths(baseURL string) []string {
	var paths []string

	// Silent mode - technology discovery

	// Get the main page to detect technologies
	client := &http.Client{Timeout: w.Timeout}
	resp, err := client.Get(baseURL)
	if err != nil {
		return paths
	}
	defer resp.Body.Close()

	body := make([]byte, 10240)
	n, _ := resp.Body.Read(body)
	content := strings.ToLower(string(body[:n]))
	headers := resp.Header

	// WordPress detection
	if strings.Contains(content, "wp-content") || strings.Contains(content, "wordpress") {
		// WordPress detected
		wpPaths := []string{
			"/wp-admin/", "/wp-content/", "/wp-includes/",
			"/wp-login.php", "/wp-config.php", "/xmlrpc.php",
			"/wp-json/", "/wp-json/wp/v2/users", "/wp-json/wp/v2/posts",
		}
		paths = append(paths, w.checkPathsExistence(baseURL, wpPaths)...)
	}

	// Drupal detection
	if strings.Contains(content, "drupal") || strings.Contains(content, "sites/default") {
		// Drupal detected
		drupalPaths := []string{
			"/admin/", "/user/", "/node/", "/sites/default/",
			"/modules/", "/themes/", "/misc/", "/core/",
		}
		paths = append(paths, w.checkPathsExistence(baseURL, drupalPaths)...)
	}

	// Joomla detection
	if strings.Contains(content, "joomla") || strings.Contains(content, "administrator/index.php") {
		// Joomla detected
		joomlaPaths := []string{
			"/administrator/", "/components/", "/modules/",
			"/templates/", "/libraries/", "/cache/",
		}
		paths = append(paths, w.checkPathsExistence(baseURL, joomlaPaths)...)
	}

	// PHP detection
	if headers.Get("X-Powered-By") != "" && strings.Contains(headers.Get("X-Powered-By"), "PHP") {
		// PHP detected
		phpPaths := []string{
			"/index.php", "/config.php", "/info.php",
			"/phpinfo.php", "/test.php", "/admin.php",
		}
		paths = append(paths, w.checkPathsExistence(baseURL, phpPaths)...)
	}

	// Apache detection
	if strings.Contains(headers.Get("Server"), "Apache") {
		// Apache detected
		apachePaths := []string{
			"/server-status", "/server-info", "/cgi-bin/",
		}
		paths = append(paths, w.checkPathsExistence(baseURL, apachePaths)...)
	}

	return paths
}

// discoverSubdomains performs subdomain discovery
func (w *WebDiscoveryEngine) discoverSubdomains(target string, result *WebResult) {
	var allSubdomains []string
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Channel to collect subdomains from all tools
	subdomainsChan := make(chan []string, 4) // Buffer for 4 tools

	// Method 1: Try subfinder (ProjectDiscovery) - run concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		if subs := w.runSubfinder(target); len(subs) > 0 {
			subdomainsChan <- subs
		} else {
			subdomainsChan <- []string{}
		}
	}()

	// Method 2: Try sublist3r - run concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		if subs := w.runSublist3r(target); len(subs) > 0 {
			subdomainsChan <- subs
		} else {
			subdomainsChan <- []string{}
		}
	}()

	// Method 3: Try amass - run concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		if subs := w.runAmass(target); len(subs) > 0 {
			subdomainsChan <- subs
		} else {
			subdomainsChan <- []string{}
		}
	}()

	// Method 4: Manual discovery via certificate transparency, DNS, etc. - run concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		if subs := w.manualSubdomainDiscovery(target); len(subs) > 0 {
			subdomainsChan <- subs
		} else {
			subdomainsChan <- []string{}
		}
	}()

	// Collect results in a separate goroutine
	go func() {
		wg.Wait()
		close(subdomainsChan)
	}()

	// Collect all subdomains from all tools
	for subdomains := range subdomainsChan {
		mu.Lock()
		allSubdomains = append(allSubdomains, subdomains...)
		mu.Unlock()
	}

	// Remove duplicates and validate
	result.Subdomains = w.validateSubdomains(w.removeDuplicates(allSubdomains))
}

// runSubfinder executes subfinder for subdomain discovery
func (w *WebDiscoveryEngine) runSubfinder(domain string) []string {
	var subdomains []string

	if _, err := exec.LookPath("subfinder"); err != nil {
		return subdomains
	}

	cmd := exec.Command("subfinder",
		"-d", domain,
		"-silent",
		"-o", "-", // Output to stdout
	)

	output, err := cmd.Output()
	if err != nil {
		return subdomains
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && strings.Contains(line, domain) {
			subdomains = append(subdomains, line)
		}
	}

	return subdomains
}

// runSublist3r executes sublist3r for subdomain discovery
func (w *WebDiscoveryEngine) runSublist3r(domain string) []string {
	var subdomains []string

	if _, err := exec.LookPath("sublist3r"); err != nil {
		return subdomains
	}

	cmd := exec.Command("sublist3r",
		"-d", domain,
		"-o", "-", // Output to stdout
	)

	output, err := cmd.Output()
	if err != nil {
		return subdomains
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && strings.Contains(line, domain) && !strings.Contains(line, "Starting") {
			subdomains = append(subdomains, line)
		}
	}

	return subdomains
}

// runAmass executes amass for subdomain discovery
func (w *WebDiscoveryEngine) runAmass(domain string) []string {
	var subdomains []string

	if _, err := exec.LookPath("amass"); err != nil {
		return subdomains
	}

	cmd := exec.Command("amass", "enum",
		"-d", domain,
		"-passive",
		"-silent",
	)

	output, err := cmd.Output()
	if err != nil {
		return subdomains
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && strings.Contains(line, domain) {
			subdomains = append(subdomains, line)
		}
	}

	return subdomains
}

// manualSubdomainDiscovery performs manual subdomain discovery
func (w *WebDiscoveryEngine) manualSubdomainDiscovery(domain string) []string {
	var subdomains []string

	// Common subdomain patterns
	commonSubs := []string{
		"www", "mail", "email", "webmail", "secure", "docs", "blog",
		"admin", "administrator", "login", "signin", "auth",
		"api", "rest", "graphql", "ws", "wss",
		"dev", "development", "test", "testing", "stage", "staging",
		"demo", "beta", "alpha", "preview",
		"app", "apps", "application", "portal",
		"static", "assets", "cdn", "media", "img", "images",
		"ftp", "sftp", "ssh", "vpn", "remote",
		"monitoring", "status", "health", "metrics",
		"support", "help", "kb", "wiki", "forum",
		"shop", "store", "cart", "payment", "billing",
		"mobile", "m", "wap", "touch",
	}

	client := &http.Client{Timeout: time.Second * 3}
	var wg sync.WaitGroup
	var mu sync.Mutex
	semaphore := make(chan bool, 10)

	for _, sub := range commonSubs {
		wg.Add(1)
		go func(subdomain string) {
			defer wg.Done()
			semaphore <- true
			defer func() { <-semaphore }()

			testDomain := subdomain + "." + domain
			testURL := "http://" + testDomain

			req, err := http.NewRequest("GET", testURL, nil)
			if err != nil {
				return
			}

			req.Header.Set("User-Agent", w.UserAgent)
			resp, err := client.Do(req)
			if err != nil {
				// Try HTTPS
				testURL = "https://" + testDomain
				req, err = http.NewRequest("GET", testURL, nil)
				if err != nil {
					return
				}
				req.Header.Set("User-Agent", w.UserAgent)
				resp, err = client.Do(req)
				if err != nil {
					return
				}
			}
			resp.Body.Close()

			mu.Lock()
			subdomains = append(subdomains, testDomain)
			mu.Unlock()
		}(sub)
	}

	wg.Wait()
	return subdomains
}

// validateSubdomains validates discovered subdomains
func (w *WebDiscoveryEngine) validateSubdomains(subdomains []string) []string {
	var validated []string

	client := &http.Client{
		Timeout: time.Second * 5,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	semaphore := make(chan bool, 5)

	for _, subdomain := range subdomains {
		wg.Add(1)
		go func(sub string) {
			defer wg.Done()
			semaphore <- true
			defer func() { <-semaphore }()

			for _, scheme := range []string{"https", "http"} {
				testURL := scheme + "://" + sub
				req, err := http.NewRequest("GET", testURL, nil)
				if err != nil {
					continue
				}

				req.Header.Set("User-Agent", w.UserAgent)
				resp, err := client.Do(req)
				if err != nil {
					continue
				}
				resp.Body.Close()

				if resp.StatusCode < 500 {
					mu.Lock()
					validated = append(validated, sub)
					mu.Unlock()
					break
				}
			}
		}(subdomain)
	}

	wg.Wait()
	return validated
}

// Helper functions

func (w *WebDiscoveryEngine) removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	var result []string

	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}

	sort.Strings(result)
	return result
}

func (w *WebDiscoveryEngine) removeDuplicatePaths(paths []string) []string {
	return w.removeDuplicates(paths)
}

// GetWebPorts filters ports for web services
func GetWebPorts(ports []PortResult) []PortResult {
	var webPorts []PortResult
	webPortNumbers := map[int]bool{
		80: true, 443: true, 8000: true, 8008: true, 8080: true,
		8443: true, 8888: true, 9000: true, 9001: true, 9080: true,
		9090: true, 9443: true, 3000: true, 3001: true, 4000: true,
		4001: true, 5000: true, 5001: true, 7001: true, 7002: true,
	}

	for _, port := range ports {
		if webPortNumbers[port.Port] || strings.Contains(strings.ToLower(port.Service), "http") {
			webPorts = append(webPorts, port)
		}
	}

	return webPorts
}
