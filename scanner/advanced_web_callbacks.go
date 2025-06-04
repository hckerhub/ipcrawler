package scanner

import (
	"time"
)

// AdvancedWebProgressCallback defines a function type for advanced web analysis progress updates
type AdvancedWebProgressCallback func(toolName string, status string, complete bool)

// AnalyzeWebTechnologiesWithCallback performs technology analysis with progress callbacks
func (a *AdvancedWebAnalyzer) AnalyzeWebTechnologiesWithCallback(urls []string, callback AdvancedWebProgressCallback) ([]WebTechResult, error) {
	var results []WebTechResult

	for _, url := range urls {
		if callback != nil {
			callback("whatweb", "Analyzing "+url+" with WhatWeb", false)
		}

		var whatwebTechs []WhatWebTechnology
		if a.isToolInstalled("whatweb") {
			whatwebTechs = a.runWhatWeb(url)
		}

		if callback != nil {
			callback("whatweb", "WhatWeb analysis complete for "+url, true)
		}

		if callback != nil {
			callback("wappalyzer", "Analyzing "+url+" with Wappalyzer", false)
		}

		var wappalyzerTechs []WappalyzerTech
		if a.isToolInstalled("wappalyzer") {
			wappalyzerTechs = a.runWappalyzer(url)
		}

		if callback != nil {
			callback("wappalyzer", "Wappalyzer analysis complete for "+url, true)
		}

		if callback != nil {
			callback("internal", "Performing manual analysis for "+url, false)
		}

		// Create result and combine technologies
		result := &WebTechResult{
			URL:               url,
			WhatWebResults:    whatwebTechs,
			WappalyzerResults: wappalyzerTechs,
			Timestamp:         time.Now(),
		}

		// Combine and normalize results
		result.CombinedTechs = a.combineAndNormalizeTechnologies(result)

		results = append(results, *result)

		if callback != nil {
			callback("internal", "Manual analysis complete for "+url, true)
		}
	}

	return results, nil
}
