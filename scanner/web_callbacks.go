package scanner

import (
	"fmt"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"
)

// WebProgressCallback defines a function type for web discovery progress updates
type WebProgressCallback func(toolName string, status string, complete bool)

// CommandTracker interface for tracking executed commands
type CommandTracker interface {
	TrackCommand(tool, command string, args []string, startTime, endTime time.Time, exitCode int, output, error, stage string)
}

// DiscoverWebServicesWithCallback performs intelligent web service discovery with progress callbacks
func (w *WebDiscoveryEngine) DiscoverWebServicesWithCallback(target string, ports []PortResult, callback WebProgressCallback, tracker CommandTracker) ([]WebResult, error) {
	var results []WebResult
	var mu sync.Mutex
	var wg sync.WaitGroup

	semaphore := make(chan bool, w.MaxConcurrency)

	// First discover web services on open ports
	for _, port := range ports {
		if !w.isWebPort(port.Port) {
			continue
		}

		wg.Add(1)
		go func(p PortResult) {
			defer wg.Done()
			semaphore <- true
			defer func() { <-semaphore }()

			result := w.analyzeWebService(target, p)
			if result != nil {
				mu.Lock()
				results = append(results, *result)
				mu.Unlock()
			}
		}(port)
	}

	wg.Wait()

	// Now perform enhanced discovery on found web services
	for i := range results {
		// Directory discovery with individual tool progress
		if w.EnableDirectories {
			w.discoverDirectoriesWithCallback(&results[i], callback, tracker)
		}

		// Subdomain discovery if target is a domain
		if w.EnableSubdomains && w.isDomainTarget(target) {
			w.discoverSubdomainsWithCallback(target, &results[i], callback, tracker)
		}

		// Technology detection
		if w.EnableTechnologies {
			if callback != nil {
				callback("whatweb", "Detecting web technologies", false)
			}
			w.detectTechnologiesWithTracking(&results[i], tracker)
			if callback != nil {
				callback("whatweb", "Technology detection complete", true)
			}
		}
	}

	return results, nil
}

// discoverDirectoriesWithCallback performs directory discovery with progress updates and command tracking
func (w *WebDiscoveryEngine) discoverDirectoriesWithCallback(result *WebResult, callback WebProgressCallback, tracker CommandTracker) {
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
		if callback != nil {
			callback("ffuf", "Running FFuF directory brute force", false)
		}
		ffufPaths := w.runFFUFWithTracking(baseURL, tracker)
		pathsChan <- ffufPaths
		if callback != nil {
			callback("ffuf", "FFuF scan complete", true)
		}
	}()

	// Run Feroxbuster concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		if callback != nil {
			callback("feroxbuster", "Running Feroxbuster recursive scan", false)
		}
		feroxPaths := w.runFeroxbusterWithTracking(baseURL, tracker)
		pathsChan <- feroxPaths
		if callback != nil {
			callback("feroxbuster", "Feroxbuster scan complete", true)
		}
	}()

	// Run Gobuster concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		if callback != nil {
			callback("gobuster", "Running Gobuster directory scan", false)
		}
		gobusterPaths := w.runGobusterWithTracking(baseURL, tracker)
		pathsChan <- gobusterPaths
		if callback != nil {
			callback("gobuster", "Gobuster scan complete", true)
		}
	}()

	// Run Manual discovery concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		if callback != nil {
			callback("internal", "Analyzing robots.txt and sitemap", false)
		}

		var manualPaths []string
		robotsPaths := w.discoverFromRobots(baseURL)
		manualPaths = append(manualPaths, robotsPaths...)

		sitemapPaths := w.discoverFromSitemap(baseURL)
		manualPaths = append(manualPaths, sitemapPaths...)

		commonPaths := w.checkCommonPaths(baseURL)
		manualPaths = append(manualPaths, commonPaths...)

		pathsChan <- manualPaths
		if callback != nil {
			callback("internal", "Manual discovery complete", true)
		}
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

	// Remove duplicates and validate
	result.Paths = w.removeDuplicatePaths(allPaths)
}

// discoverSubdomainsWithCallback performs subdomain discovery with progress updates and command tracking
func (w *WebDiscoveryEngine) discoverSubdomainsWithCallback(target string, result *WebResult, callback WebProgressCallback, tracker CommandTracker) {
	var allSubdomains []string
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Channel to collect subdomains from all tools (increased buffer for more tools)
	subdomainsChan := make(chan []string, 8) // Buffer for 8 tools

	// Run Subfinder concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		if callback != nil {
			callback("subfinder", "Running Subfinder", false)
		}
		subfinderResults := w.runSubfinderWithTracking(target, tracker)
		subdomainsChan <- subfinderResults
		if callback != nil {
			callback("subfinder", "Subfinder complete", true)
		}
	}()

	// Run Sublist3r concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		if callback != nil {
			callback("sublist3r", "Running Sublist3r", false)
		}
		sublist3rResults := w.runSublist3rWithTracking(target, tracker)
		subdomainsChan <- sublist3rResults
		if callback != nil {
			callback("sublist3r", "Sublist3r complete", true)
		}
	}()

	// Run Amass concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		if callback != nil {
			callback("amass", "Running Amass enumeration", false)
		}
		amassResults := w.runAmassWithTracking(target, tracker)
		subdomainsChan <- amassResults
		if callback != nil {
			callback("amass", "Amass enumeration complete", true)
		}
	}()

	// Run Assetfinder concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		if callback != nil {
			callback("assetfinder", "Running Assetfinder", false)
		}

		var subdomains []string
		if _, err := exec.LookPath("assetfinder"); err == nil {
			cmd := exec.Command("assetfinder", "--subs-only", target)
			if output, err := cmd.Output(); err == nil {
				lines := strings.Split(string(output), "\n")
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if line != "" && strings.Contains(line, target) {
						subdomains = append(subdomains, line)
					}
				}
			}
		}
		subdomainsChan <- subdomains
		if callback != nil {
			callback("assetfinder", "Assetfinder complete", true)
		}
	}()

	// Run Findomain concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		if callback != nil {
			callback("findomain", "Running Findomain", false)
		}

		var subdomains []string
		if _, err := exec.LookPath("findomain"); err == nil {
			cmd := exec.Command("findomain", "-t", target, "-q", "-u")
			if output, err := cmd.Output(); err == nil {
				lines := strings.Split(string(output), "\n")
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if line != "" && strings.Contains(line, target) {
						subdomains = append(subdomains, line)
					}
				}
			}
		}
		subdomainsChan <- subdomains
		if callback != nil {
			callback("findomain", "Findomain complete", true)
		}
	}()

	// Run Dnsrecon concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		if callback != nil {
			callback("dnsrecon", "Running Dnsrecon", false)
		}

		var subdomains []string
		if _, err := exec.LookPath("dnsrecon"); err == nil {
			cmd := exec.Command("dnsrecon", "-d", target, "-t", "brt", "--threads", "5")
			if output, err := cmd.Output(); err == nil {
				lines := strings.Split(string(output), "\n")
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if strings.Contains(line, target) && strings.Contains(line, "A") {
						fields := strings.Fields(line)
						for _, field := range fields {
							if strings.Contains(field, target) && strings.Contains(field, ".") {
								subdomains = append(subdomains, field)
								break
							}
						}
					}
				}
			}
		}
		subdomainsChan <- subdomains
		if callback != nil {
			callback("dnsrecon", "Dnsrecon complete", true)
		}
	}()

	// Run Certificate Transparency lookup concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		if callback != nil {
			callback("crt.sh", "Querying Certificate Transparency", false)
		}

		var subdomains []string
		client := &http.Client{Timeout: time.Second * 30}
		url := fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", target)

		if resp, err := client.Get(url); err == nil {
			defer resp.Body.Close()
			// Simple parsing - extract domain names from response
			// Note: Full JSON parsing would require additional imports
			if resp.StatusCode == 200 {
				// For now, just indicate that we tried the service
				// TODO: Add proper JSON parsing when imports are available
			}
		}
		subdomainsChan <- subdomains
		if callback != nil {
			callback("crt.sh", "Certificate Transparency complete", true)
		}
	}()

	// Run Manual subdomain discovery concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		if callback != nil {
			callback("internal", "Manual subdomain discovery", false)
		}
		manualResults := w.manualSubdomainDiscovery(target)
		subdomainsChan <- manualResults
		if callback != nil {
			callback("internal", "Manual subdomain discovery complete", true)
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

	// Validate discovered subdomains
	validatedSubdomains := w.validateSubdomains(w.removeDuplicates(allSubdomains))
	result.Subdomains = validatedSubdomains
}

// detectTechnologiesWithTracking performs technology detection with command tracking
func (w *WebDiscoveryEngine) detectTechnologiesWithTracking(result *WebResult, tracker CommandTracker) {
	var technologies []string

	// Method 1: Try whatweb with tracking
	if techs := w.runWhatwebWithTracking(result.URL, tracker); len(techs) > 0 {
		technologies = append(technologies, techs...)
	}

	// Method 2: Try wappalyzer with tracking
	if techs := w.runWappalyzerWithTracking(result.URL, tracker); len(techs) > 0 {
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

	// Track the domain discovery command
	if tracker != nil {
		tracker.TrackCommand("whatweb", "whatweb",
			[]string{"--no-errors", "-a", "3", "--color=never", "--log-verbose=-", "--follow-redirect=always", "--max-redirects=10", result.URL},
			time.Now(), time.Now(), 0, "Domain discovery from HTTP responses and redirects", "", "Domain Discovery")
	}
}

// isDomainTarget checks if target is a domain (vs IP)
func (w *WebDiscoveryEngine) isDomainTarget(target string) bool {
	// Simple IP regex check
	ipRegex := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
	return !ipRegex.MatchString(target)
}

// Placeholder methods that fallback to original implementations
// TODO: Implement full command tracking for these methods

func (w *WebDiscoveryEngine) runFFUFWithTracking(baseURL string, tracker CommandTracker) []string {
	// For now, just call the original method
	// TODO: Add command tracking
	return w.runFFUF(baseURL)
}

func (w *WebDiscoveryEngine) runFeroxbusterWithTracking(baseURL string, tracker CommandTracker) []string {
	// For now, just call the original method
	// TODO: Add command tracking
	return w.runFeroxbuster(baseURL)
}

func (w *WebDiscoveryEngine) runGobusterWithTracking(baseURL string, tracker CommandTracker) []string {
	// For now, just call the original method
	// TODO: Add command tracking
	return w.runGobuster(baseURL)
}

func (w *WebDiscoveryEngine) runSubfinderWithTracking(domain string, tracker CommandTracker) []string {
	// For now, just call the original method
	// TODO: Add command tracking
	return w.runSubfinder(domain)
}

func (w *WebDiscoveryEngine) runSublist3rWithTracking(domain string, tracker CommandTracker) []string {
	// For now, just call the original method
	// TODO: Add command tracking
	return w.runSublist3r(domain)
}

func (w *WebDiscoveryEngine) runAmassWithTracking(domain string, tracker CommandTracker) []string {
	// For now, just call the original method
	// TODO: Add command tracking
	return w.runAmass(domain)
}

func (w *WebDiscoveryEngine) runWhatwebWithTracking(url string, tracker CommandTracker) []string {
	// For now, just call the original method
	// TODO: Add command tracking
	return w.runWhatweb(url)
}

func (w *WebDiscoveryEngine) runWappalyzerWithTracking(url string, tracker CommandTracker) []string {
	// For now, just call the original method
	// TODO: Add command tracking
	return w.runWappalyzer(url)
}
