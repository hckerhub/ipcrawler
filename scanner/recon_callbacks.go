package scanner

import (
	"sync"
	"time"
)

// ReconProgressCallback defines a function type for reconnaissance progress updates
type ReconProgressCallback func(toolName string, status string, complete bool)

// PerformReconnaissanceWithCallback runs reconnaissance tools with progress callbacks
func (r *ReconEngine) PerformReconnaissanceWithCallback(target string, callback ReconProgressCallback) (*ReconResult, error) {
	return r.PerformReconnaissanceWithCallbackAndTracker(target, callback, nil)
}

// PerformReconnaissanceWithCallbackAndTracker runs reconnaissance tools with progress callbacks and command tracking
func (r *ReconEngine) PerformReconnaissanceWithCallbackAndTracker(target string, callback ReconProgressCallback, tracker CommandTracker) (*ReconResult, error) {
	result := &ReconResult{
		Target:    target,
		Type:      r.determineTargetType(target),
		Timestamp: time.Now(),
	}

	// Set the tracker for this reconnaissance session
	r.CommandTracker = tracker

	var wg sync.WaitGroup

	// For IP targets, run reverse DNS lookup first
	if r.determineTargetType(target) == "IP" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if callback != nil {
				callback("dig", "Performing reverse DNS lookup", false)
			}

			r.runReverseDNSLookup(target, result)

			if callback != nil {
				callback("dig", "Reverse DNS lookup complete", true)
			}
		}()
	}

	// Run all reconnaissance tools concurrently
	// Amass subdomain enumeration
	wg.Add(1)
	go func() {
		defer wg.Done()
		if callback != nil {
			callback("amass", "Running Amass subdomain enumeration", false)
		}

		var err error
		if r.isToolInstalled("amass") {
			err = r.runAmass(target, result)
			if err != nil && r.Verbose {
				// Log error but continue
			}
		} else {
			// Track command even if tool not installed
			if tracker != nil {
				tracker.TrackCommand("amass", "amass", []string{"enum", "-d", target, "-passive", "-json"}, time.Now(), time.Now(), 1, "", "amass not installed", "Reconnaissance")
			}
		}

		if callback != nil {
			callback("amass", "Amass enumeration complete", true)
		}
	}()

	// Recon-ng framework
	wg.Add(1)
	go func() {
		defer wg.Done()
		if callback != nil {
			callback("recon-ng", "Running Recon-ng modules", false)
		}

		err := r.runReconNG(target, result)
		if err != nil && r.Verbose {
			// Log error but continue
		}

		if callback != nil {
			callback("recon-ng", "Recon-ng analysis complete", true)
		}
	}()

	// Certificate transparency (crt.sh)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if callback != nil {
			callback("crt.sh", "Querying certificate transparency logs", false)
		}

		err := r.runCrtSh(target, result)
		if err != nil && r.Verbose {
			// Log error but continue
		}

		if callback != nil {
			callback("crt.sh", "Certificate transparency search complete", true)
		}
	}()

	// SecurityTrails API
	wg.Add(1)
	go func() {
		defer wg.Done()
		if callback != nil {
			callback("securitytrails", "Querying SecurityTrails API", false)
		}

		err := r.runSecurityTrails(target, result)
		if err != nil && r.Verbose {
			// Log error but continue
		}

		if callback != nil {
			callback("securitytrails", "SecurityTrails analysis complete", true)
		}
	}()

	// Censys API
	wg.Add(1)
	go func() {
		defer wg.Done()
		if callback != nil {
			callback("censys", "Querying Censys database", false)
		}

		err := r.runCensys(target, result)
		if err != nil && r.Verbose {
			// Log error but continue
		}

		if callback != nil {
			callback("censys", "Censys analysis complete", true)
		}
	}()

	// Netcraft reconnaissance
	wg.Add(1)
	go func() {
		defer wg.Done()
		if callback != nil {
			callback("netcraft", "Running Netcraft reconnaissance", false)
		}

		err := r.runNetcraft(target, result)
		if err != nil && r.Verbose {
			// Log error but continue
		}

		if callback != nil {
			callback("netcraft", "Netcraft analysis complete", true)
		}
	}()

	// Wait for all reconnaissance tools to complete
	wg.Wait()

	// Post-process and deduplicate results
	r.postProcessResults(result)

	return result, nil
}
