package scanner

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"
)

// AdvancedWebAnalyzer extends web discovery with additional tools
type AdvancedWebAnalyzer struct {
	Timeout        time.Duration
	MaxConcurrency int
	Verbose        bool
	UserAgent      string
}

// WebTechResult represents comprehensive web technology analysis
type WebTechResult struct {
	URL               string              `json:"url"`
	WhatWebResults    []WhatWebTechnology `json:"whatweb_results"`
	WappalyzerResults []WappalyzerTech    `json:"wappalyzer_results"`
	CombinedTechs     []TechnologyInfo    `json:"combined_technologies"`
	ResponseHeaders   map[string]string   `json:"response_headers"`
	StatusCode        int                 `json:"status_code"`
	Server            string              `json:"server"`
	ContentType       string              `json:"content_type"`
	Title             string              `json:"title"`
	AnalysisTime      time.Duration       `json:"analysis_time"`
	Timestamp         time.Time           `json:"timestamp"`
}

// WhatWebTechnology represents whatweb output
type WhatWebTechnology struct {
	Name       string                 `json:"name"`
	Version    string                 `json:"version"`
	Categories []string               `json:"categories"`
	Certainty  string                 `json:"certainty"`
	Details    map[string]interface{} `json:"details"`
}

// WappalyzerTech represents wappalyzer output
type WappalyzerTech struct {
	Name       string   `json:"name"`
	Categories []string `json:"categories"`
	Version    string   `json:"version"`
	Website    string   `json:"website"`
	Confidence string   `json:"confidence"`
}

// TechnologyInfo represents normalized technology information
type TechnologyInfo struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Category    string   `json:"category"`
	Confidence  int      `json:"confidence"`
	Sources     []string `json:"sources"`
	Description string   `json:"description"`
}

// NewAdvancedWebAnalyzer creates a new advanced web analyzer
func NewAdvancedWebAnalyzer() *AdvancedWebAnalyzer {
	return &AdvancedWebAnalyzer{
		Timeout:        time.Second * 30,
		MaxConcurrency: 3,
		Verbose:        true,
		UserAgent:      "Mozilla/5.0 (compatible; IPCrawler/2.0; +https://github.com/ipcrawler)",
	}
}

// AnalyzeWebTechnologies performs comprehensive web technology analysis
func (a *AdvancedWebAnalyzer) AnalyzeWebTechnologies(urls []string) ([]WebTechResult, error) {
	var results []WebTechResult
	var wg sync.WaitGroup
	semaphore := make(chan bool, a.MaxConcurrency)
	resultsChan := make(chan WebTechResult, len(urls))

	for _, url := range urls {
		wg.Add(1)
		go func(targetURL string) {
			defer wg.Done()
			semaphore <- true
			defer func() { <-semaphore }()

			result := a.analyzeURL(targetURL)
			if result != nil {
				resultsChan <- *result
			}
		}(url)
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

// analyzeURL performs comprehensive analysis of a single URL
func (a *AdvancedWebAnalyzer) analyzeURL(url string) *WebTechResult {
	startTime := time.Now()

	result := &WebTechResult{
		URL:       url,
		Timestamp: time.Now(),
	}

	// Run whatweb analysis
	if a.isToolInstalled("whatweb") {
		result.WhatWebResults = a.runWhatWeb(url)
	}

	// Run wappalyzer analysis
	if a.isToolInstalled("wappalyzer") {
		result.WappalyzerResults = a.runWappalyzer(url)
	}

	// Combine and normalize results
	result.CombinedTechs = a.combineAndNormalizeTechnologies(result)

	result.AnalysisTime = time.Since(startTime)
	return result
}

// runWhatWeb executes whatweb for technology detection
func (a *AdvancedWebAnalyzer) runWhatWeb(url string) []WhatWebTechnology {
	cmd := exec.Command("whatweb", "--color=never", "--no-errors", "--format=json", url)
	output, err := cmd.Output()
	if err != nil {
		if a.Verbose {
			fmt.Printf("WhatWeb failed for %s: %v\n", url, err)
		}
		return nil
	}

	return a.parseWhatWebOutput(string(output))
}

// runWappalyzer executes wappalyzer CLI for technology detection
func (a *AdvancedWebAnalyzer) runWappalyzer(url string) []WappalyzerTech {
	cmd := exec.Command("wappalyzer", url, "--pretty")
	output, err := cmd.Output()
	if err != nil {
		if a.Verbose {
			fmt.Printf("Wappalyzer failed for %s: %v\n", url, err)
		}
		return nil
	}

	return a.parseWappalyzerOutput(string(output))
}

// parseWhatWebOutput parses whatweb JSON output
func (a *AdvancedWebAnalyzer) parseWhatWebOutput(output string) []WhatWebTechnology {
	var technologies []WhatWebTechnology

	// WhatWeb outputs one JSON object per line
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		var result map[string]interface{}
		if err := json.Unmarshal([]byte(line), &result); err != nil {
			continue
		}

		if plugins, ok := result["plugins"].(map[string]interface{}); ok {
			for pluginName, pluginData := range plugins {
				tech := WhatWebTechnology{
					Name: pluginName,
				}

				if data, ok := pluginData.(map[string]interface{}); ok {
					// Extract version
					if version, ok := data["version"].([]interface{}); ok && len(version) > 0 {
						if v, ok := version[0].(string); ok {
							tech.Version = v
						}
					}

					// Extract categories
					if categories, ok := data["category"].([]interface{}); ok {
						for _, cat := range categories {
							if c, ok := cat.(string); ok {
								tech.Categories = append(tech.Categories, c)
							}
						}
					}

					// Extract certainty
					if certainty, ok := data["certainty"].(string); ok {
						tech.Certainty = certainty
					}

					tech.Details = data
				}

				technologies = append(technologies, tech)
			}
		}
	}

	return technologies
}

// parseWappalyzerOutput parses wappalyzer output
func (a *AdvancedWebAnalyzer) parseWappalyzerOutput(output string) []WappalyzerTech {
	var technologies []WappalyzerTech

	// Try to parse as JSON first
	var wappData map[string]interface{}
	if err := json.Unmarshal([]byte(output), &wappData); err == nil {
		if technologies_data, ok := wappData["technologies"].([]interface{}); ok {
			for _, tech_data := range technologies_data {
				if tech, ok := tech_data.(map[string]interface{}); ok {
					wappTech := WappalyzerTech{}

					if name, ok := tech["name"].(string); ok {
						wappTech.Name = name
					}

					if version, ok := tech["version"].(string); ok {
						wappTech.Version = version
					}

					if website, ok := tech["website"].(string); ok {
						wappTech.Website = website
					}

					if confidence, ok := tech["confidence"].(string); ok {
						wappTech.Confidence = confidence
					}

					if categories, ok := tech["categories"].([]interface{}); ok {
						for _, cat := range categories {
							if c, ok := cat.(map[string]interface{}); ok {
								if name, ok := c["name"].(string); ok {
									wappTech.Categories = append(wappTech.Categories, name)
								}
							}
						}
					}

					technologies = append(technologies, wappTech)
				}
			}
		}
	} else {
		// Fallback to text parsing
		technologies = a.parseWappalyzerTextOutput(output)
	}

	return technologies
}

// parseWappalyzerTextOutput parses wappalyzer text output as fallback
func (a *AdvancedWebAnalyzer) parseWappalyzerTextOutput(output string) []WappalyzerTech {
	var technologies []WappalyzerTech

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "http") {
			continue
		}

		// Parse format like: "Technology Name version (category)"
		re := regexp.MustCompile(`^([^(]+?)(?:\s+([0-9][^\s]*))?\s*(?:\(([^)]+)\))?$`)
		matches := re.FindStringSubmatch(line)

		if len(matches) > 1 {
			tech := WappalyzerTech{
				Name: strings.TrimSpace(matches[1]),
			}

			if len(matches) > 2 && matches[2] != "" {
				tech.Version = strings.TrimSpace(matches[2])
			}

			if len(matches) > 3 && matches[3] != "" {
				tech.Categories = []string{strings.TrimSpace(matches[3])}
			}

			technologies = append(technologies, tech)
		}
	}

	return technologies
}

// combineAndNormalizeTechnologies combines results from different tools
func (a *AdvancedWebAnalyzer) combineAndNormalizeTechnologies(result *WebTechResult) []TechnologyInfo {
	techMap := make(map[string]*TechnologyInfo)

	// Process WhatWeb results
	for _, tech := range result.WhatWebResults {
		key := strings.ToLower(tech.Name)
		if existing, exists := techMap[key]; exists {
			// Merge sources
			existing.Sources = append(existing.Sources, "whatweb")
			if tech.Version != "" && existing.Version == "" {
				existing.Version = tech.Version
			}
		} else {
			confidence := a.convertCertaintyToConfidence(tech.Certainty)
			category := ""
			if len(tech.Categories) > 0 {
				category = tech.Categories[0]
			}

			techMap[key] = &TechnologyInfo{
				Name:       tech.Name,
				Version:    tech.Version,
				Category:   category,
				Confidence: confidence,
				Sources:    []string{"whatweb"},
			}
		}
	}

	// Process Wappalyzer results
	for _, tech := range result.WappalyzerResults {
		key := strings.ToLower(tech.Name)
		if existing, exists := techMap[key]; exists {
			// Merge sources
			existing.Sources = append(existing.Sources, "wappalyzer")
			if tech.Version != "" && existing.Version == "" {
				existing.Version = tech.Version
			}
			// Increase confidence if detected by multiple tools
			existing.Confidence = min(100, existing.Confidence+20)
		} else {
			confidence := a.convertWappalyzerConfidenceToInt(tech.Confidence)
			category := ""
			if len(tech.Categories) > 0 {
				category = tech.Categories[0]
			}

			techMap[key] = &TechnologyInfo{
				Name:       tech.Name,
				Version:    tech.Version,
				Category:   category,
				Confidence: confidence,
				Sources:    []string{"wappalyzer"},
			}
		}
	}

	// Convert map to slice
	var technologies []TechnologyInfo
	for _, tech := range techMap {
		technologies = append(technologies, *tech)
	}

	return technologies
}

// convertCertaintyToConfidence converts whatweb certainty to confidence percentage
func (a *AdvancedWebAnalyzer) convertCertaintyToConfidence(certainty string) int {
	switch strings.ToLower(certainty) {
	case "certain":
		return 100
	case "firm":
		return 80
	case "tentative":
		return 60
	case "maybe":
		return 40
	default:
		return 50
	}
}

// convertWappalyzerConfidenceToInt converts wappalyzer confidence to int
func (a *AdvancedWebAnalyzer) convertWappalyzerConfidenceToInt(confidence string) int {
	if confidence == "" {
		return 50
	}

	// Remove % sign if present
	confidence = strings.TrimSuffix(confidence, "%")

	// Try to parse as number
	val := 0
	if n, _ := fmt.Sscanf(confidence, "%d", &val); n == 1 {
		return val
	}

	return 50
}

// isToolInstalled checks if a tool is installed
func (a *AdvancedWebAnalyzer) isToolInstalled(tool string) bool {
	_, err := exec.LookPath(tool)
	return err == nil
}

// GetInstalledWebTools returns list of available web analysis tools
func (a *AdvancedWebAnalyzer) GetInstalledWebTools() []string {
	tools := []string{"whatweb", "wappalyzer"}
	var installed []string

	for _, tool := range tools {
		if a.isToolInstalled(tool) {
			installed = append(installed, tool)
		}
	}

	return installed
}

// GetWebToolsInstallCommands provides installation commands for missing web tools
func (a *AdvancedWebAnalyzer) GetWebToolsInstallCommands() map[string]string {
	recommendations := map[string]string{
		"whatweb":    "gem install whatweb",
		"wappalyzer": "npm install -g wappalyzer",
	}

	missing := make(map[string]string)
	for tool, command := range recommendations {
		if !a.isToolInstalled(tool) {
			missing[tool] = command
		}
	}

	return missing
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// AnalyzeSingleURL provides a simple interface for analyzing one URL
func (a *AdvancedWebAnalyzer) AnalyzeSingleURL(url string) (*WebTechResult, error) {
	results, err := a.AnalyzeWebTechnologies([]string{url})
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results for URL: %s", url)
	}

	return &results[0], nil
}

// GetTechnologySummary provides a summary of detected technologies
func (a *AdvancedWebAnalyzer) GetTechnologySummary(results []WebTechResult) map[string]int {
	summary := make(map[string]int)

	for _, result := range results {
		for _, tech := range result.CombinedTechs {
			summary[tech.Name]++
		}
	}

	return summary
}
