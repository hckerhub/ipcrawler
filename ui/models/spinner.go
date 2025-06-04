package models

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"ipcrawler/ui/styles"
)

// ScanProgress represents different stages of scanning
type ScanProgress int

const (
	ProgressInitializing ScanProgress = iota
	ProgressReconnaissance
	ProgressNmapScan
	ProgressWebDiscovery
	ProgressAdvancedWebAnalysis
	ProgressVulnAnalysis
	ProgressFinishing
	ProgressComplete
)

// ScanStage represents a scanning stage with details
type ScanStage struct {
	Name        string
	Description string
	Progress    ScanProgress
	Active      bool
	Complete    bool
	Duration    time.Duration
	SubStages   []SubStage
}

// SubStage represents individual tool execution within a stage
type SubStage struct {
	Name        string
	Tool        string
	Description string
	Active      bool
	Complete    bool
	Duration    time.Duration
	Status      string
}

// ScanStatus holds the current scanning status
type ScanStatus struct {
	CurrentStage   ScanProgress
	Stages         []ScanStage
	StartTime      time.Time
	Message        string
	SpinnerFrame   int
	AnimationTick  int
	EstimatedTotal time.Duration
	TargetIP       string
	ScanType       string
	AggressiveMode bool
}

// SpinnerFrames contains different spinner animations
var SpinnerFrames = []string{
	"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏",
}

var DotsFrames = []string{
	"⠁", "⠂", "⠄", "⡀", "⢀", "⠠", "⠐", "⠈",
}

var CircleFrames = []string{
	"◐", "◓", "◑", "◒",
}

var ArrowFrames = []string{
	"→", "↘", "↓", "↙", "←", "↖", "↑", "↗",
}

// NewScanStatus creates a new scan status tracker
func NewScanStatus(targetIP, scanType string, aggressive bool) *ScanStatus {
	stages := []ScanStage{
		{
			Name:        "Initializing",
			Description: "Preparing scan parameters",
			Progress:    ProgressInitializing,
			Active:      true,
			SubStages:   []SubStage{},
		},
		{
			Name:        "Reconnaissance",
			Description: "OSINT and subdomain discovery",
			Progress:    ProgressReconnaissance,
			SubStages: []SubStage{
				{Name: "Shodan", Tool: "shodan", Description: "Internet-wide host scanning"},
				{Name: "Amass", Tool: "amass", Description: "Subdomain enumeration"},
				{Name: "Recon-ng", Tool: "recon-ng", Description: "Web reconnaissance framework"},
				{Name: "Censys", Tool: "censys", Description: "Certificate & host discovery"},
				{Name: "crt.sh", Tool: "crt.sh", Description: "Certificate transparency logs"},
				{Name: "SecurityTrails", Tool: "securitytrails", Description: "DNS history & analysis"},
			},
		},
		{
			Name:        "Port Scanning",
			Description: "Network port discovery",
			Progress:    ProgressNmapScan,
			SubStages: []SubStage{
				{Name: "Nmap Discovery", Tool: "nmap", Description: "TCP port scanning"},
				{Name: "Service Detection", Tool: "nmap", Description: "Service version detection"},
				{Name: "OS Detection", Tool: "nmap", Description: "Operating system fingerprinting"},
			},
		},
		{
			Name:        "Web Discovery",
			Description: "Web service enumeration",
			Progress:    ProgressWebDiscovery,
			SubStages: []SubStage{
				{Name: "Directory Brute", Tool: "ffuf", Description: "Directory enumeration"},
				{Name: "Path Discovery", Tool: "feroxbuster", Description: "Recursive path scanning"},
				{Name: "Gobuster", Tool: "gobuster", Description: "URI and DNS brute-forcing"},
				{Name: "Subfinder", Tool: "subfinder", Description: "ProjectDiscovery subdomain finder"},
				{Name: "Sublist3r", Tool: "sublist3r", Description: "Python subdomain enumerator"},
				{Name: "Assetfinder", Tool: "assetfinder", Description: "Find subdomains with OSINT"},
				{Name: "Findomain", Tool: "findomain", Description: "Fast subdomain enumeration"},
				{Name: "DNSRecon", Tool: "dnsrecon", Description: "DNS brute force enumeration"},
				{Name: "Cert Transparency", Tool: "crt.sh", Description: "Certificate transparency logs"},
				{Name: "Manual Discovery", Tool: "internal", Description: "Robots.txt & sitemap analysis"},
			},
		},
		{
			Name:        "Advanced Web Analysis",
			Description: "Technology detection & fingerprinting",
			Progress:    ProgressAdvancedWebAnalysis,
			SubStages: []SubStage{
				{Name: "WhatWeb", Tool: "whatweb", Description: "Ruby-based web fingerprinting"},
				{Name: "Wappalyzer", Tool: "wappalyzer", Description: "Technology stack detection"},
				{Name: "Manual Analysis", Tool: "internal", Description: "Header & banner analysis"},
			},
		},
		{
			Name:        "Vulnerability Analysis",
			Description: "Security assessment",
			Progress:    ProgressVulnAnalysis,
			SubStages: []SubStage{
				{Name: "Service Vulns", Tool: "nmap", Description: "Service-based vulnerability detection"},
				{Name: "Web Vulns", Tool: "internal", Description: "Web application security checks"},
			},
		},
		{
			Name:        "Finalizing",
			Description: "Processing results",
			Progress:    ProgressFinishing,
			SubStages:   []SubStage{},
		},
	}

	// Estimate total time based on scan type and aggressiveness
	estimatedTime := time.Minute * 3 // Default with reconnaissance
	if aggressive {
		if scanType == "full" {
			estimatedTime = time.Minute * 8 // Full aggressive scan with recon
		} else {
			estimatedTime = time.Minute * 5 // Custom aggressive scan with recon
		}
	} else {
		if scanType == "full" {
			estimatedTime = time.Minute * 12 // Full normal scan with recon
		} else {
			estimatedTime = time.Minute * 6 // Custom normal scan with recon
		}
	}

	return &ScanStatus{
		CurrentStage:   ProgressInitializing,
		Stages:         stages,
		StartTime:      time.Now(),
		Message:        "Starting enhanced scan...",
		EstimatedTotal: estimatedTime,
		TargetIP:       targetIP,
		ScanType:       scanType,
		AggressiveMode: aggressive,
	}
}

// UpdateStage updates the current scanning stage
func (s *ScanStatus) UpdateStage(stage ScanProgress, message string) {
	s.CurrentStage = stage
	s.Message = message

	// Update stage statuses
	for i := range s.Stages {
		if s.Stages[i].Progress == stage {
			s.Stages[i].Active = true
			s.Stages[i].Complete = false
		} else if s.Stages[i].Progress < stage {
			s.Stages[i].Active = false
			s.Stages[i].Complete = true
			if s.Stages[i].Duration == 0 {
				s.Stages[i].Duration = time.Since(s.StartTime)
			}
		} else {
			s.Stages[i].Active = false
			s.Stages[i].Complete = false
		}
	}
}

// UpdateSubStage updates a specific tool's status within a stage
func (s *ScanStatus) UpdateSubStage(stage ScanProgress, toolName string, status string, complete bool) {
	for i := range s.Stages {
		if s.Stages[i].Progress == stage {
			for j := range s.Stages[i].SubStages {
				if s.Stages[i].SubStages[j].Tool == toolName {
					s.Stages[i].SubStages[j].Status = status
					s.Stages[i].SubStages[j].Active = !complete
					s.Stages[i].SubStages[j].Complete = complete
					if complete && s.Stages[i].SubStages[j].Duration == 0 {
						s.Stages[i].SubStages[j].Duration = time.Since(s.StartTime)
					}
					break
				}
			}
			break
		}
	}
	s.Message = status
}

// AdvanceAnimation advances the spinner animation
func (s *ScanStatus) AdvanceAnimation() {
	s.SpinnerFrame = (s.SpinnerFrame + 1) % len(SpinnerFrames)
	s.AnimationTick++
}

// GetElapsedTime returns the elapsed time since scan start
func (s *ScanStatus) GetElapsedTime() time.Duration {
	return time.Since(s.StartTime)
}

// GetProgressPercentage estimates the overall progress percentage
func (s *ScanStatus) GetProgressPercentage() float64 {
	elapsed := s.GetElapsedTime()
	if s.EstimatedTotal <= 0 {
		return 0
	}

	percentage := float64(elapsed) / float64(s.EstimatedTotal) * 100
	if percentage > 95 {
		percentage = 95 // Cap at 95% until actually complete
	}

	return percentage
}

// RenderLoadingScreen renders the modern loading interface
func (s *ScanStatus) RenderLoadingScreen(width, height int) string {
	var b strings.Builder

	// Header section
	title := styles.RenderTitle("🎯 IPCrawler Active Scan")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Target info
	targetInfo := fmt.Sprintf("Target: %s", s.TargetIP)
	scanInfo := fmt.Sprintf("Scan Type: %s", strings.Title(s.ScanType))
	if s.AggressiveMode {
		scanInfo += " [🚀 AGGRESSIVE]"
	}

	b.WriteString(styles.RenderInfo(targetInfo))
	b.WriteString("\n")
	b.WriteString(styles.RenderInfo(scanInfo))
	b.WriteString("\n\n")

	// Progress overview
	elapsed := s.GetElapsedTime()
	progress := s.GetProgressPercentage()

	progressHeader := fmt.Sprintf("Progress: %.1f%% • Elapsed: %s • ETA: %s",
		progress,
		formatDuration(elapsed),
		formatDuration(s.EstimatedTotal-elapsed))

	b.WriteString(styles.HighlightStyle.Render(progressHeader))
	b.WriteString("\n\n")

	// Progress bar
	progressBar := s.renderProgressBar(60, progress)
	b.WriteString(progressBar)
	b.WriteString("\n\n")

	// Current stage indicator
	currentStageInfo := s.renderCurrentStage()
	b.WriteString(currentStageInfo)
	b.WriteString("\n\n")

	// Stage list
	stageList := s.renderStageList()
	b.WriteString(stageList)
	b.WriteString("\n\n")

	// Current activity with animated spinner
	activityLine := s.renderActivityLine()
	b.WriteString(activityLine)
	b.WriteString("\n\n")

	// Controls
	controls := "⏹ Press 'q' to cancel scan"
	b.WriteString(lipgloss.NewStyle().
		Foreground(styles.MutedColor).
		Italic(true).
		Render(controls))

	// Center content
	content := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Align(lipgloss.Center).
		Render(b.String())

	return content
}

// renderProgressBar creates a visual progress bar
func (s *ScanStatus) renderProgressBar(width int, percentage float64) string {
	filled := int(percentage * float64(width) / 100)
	if filled > width {
		filled = width
	}

	// Color the progress bar
	progressStyle := lipgloss.NewStyle().Foreground(styles.PrimaryColor)
	emptyStyle := lipgloss.NewStyle().Foreground(styles.MutedColor)

	coloredBar := progressStyle.Render(strings.Repeat("█", filled)) +
		emptyStyle.Render(strings.Repeat("░", width-filled))

	return fmt.Sprintf("[%s] %.1f%%", coloredBar, percentage)
}

// renderCurrentStage shows the current active stage with substages
func (s *ScanStatus) renderCurrentStage() string {
	for _, stage := range s.Stages {
		if stage.Active {
			icon := SpinnerFrames[s.SpinnerFrame]
			stageText := fmt.Sprintf("%s %s", icon, stage.Name)

			result := styles.SuccessStyle.Render(stageText) + "\n" +
				styles.RenderInfo(stage.Description)

			// Show active sub-stages if any
			if len(stage.SubStages) > 0 {
				result += "\n\n"
				for _, subStage := range stage.SubStages {
					var subIcon string
					var subStyle lipgloss.Style

					if subStage.Complete {
						subIcon = "✓"
						subStyle = styles.SuccessStyle
						result += subStyle.Render(fmt.Sprintf("  %s %s", subIcon, subStage.Name)) + "\n"
					} else if subStage.Active {
						subIcon = ArrowFrames[s.AnimationTick%len(ArrowFrames)]
						subStyle = styles.HighlightStyle
						result += subStyle.Render(fmt.Sprintf("  %s %s", subIcon, subStage.Name))
						if subStage.Status != "" {
							result += " - " + styles.RenderInfo(subStage.Status)
						}
						result += "\n"
					} else {
						subIcon = "○"
						subStyle = styles.UnselectedStyle
						result += subStyle.Render(fmt.Sprintf("  %s %s", subIcon, subStage.Name)) + "\n"
					}
				}
			}

			return result
		}
	}
	return ""
}

// renderStageList shows all stages with their status
func (s *ScanStatus) renderStageList() string {
	var stages []string

	for _, stage := range s.Stages {
		var icon, statusText string
		var style lipgloss.Style

		if stage.Complete {
			icon = "✓"
			style = styles.SuccessStyle
			duration := formatDuration(stage.Duration)
			statusText = fmt.Sprintf("%s %s (%s)", icon, stage.Name, duration)
		} else if stage.Active {
			icon = ArrowFrames[s.AnimationTick%len(ArrowFrames)]
			style = styles.HighlightStyle
			statusText = fmt.Sprintf("%s %s", icon, stage.Name)
		} else {
			icon = "○"
			style = styles.UnselectedStyle
			statusText = fmt.Sprintf("%s %s", icon, stage.Name)
		}

		stages = append(stages, style.Render(statusText))
	}

	return strings.Join(stages, "\n")
}

// renderActivityLine shows current activity with spinner
func (s *ScanStatus) renderActivityLine() string {
	spinner := SpinnerFrames[s.SpinnerFrame]
	dots := strings.Repeat(".", (s.AnimationTick/3)%4)

	activity := fmt.Sprintf("%s %s%s", spinner, s.Message, dots)

	return styles.LoadingStyle.Render(activity)
}

// formatDuration formats a duration for display
func formatDuration(d time.Duration) string {
	if d < 0 {
		return "0s"
	}

	if d < time.Minute {
		return fmt.Sprintf("%.0fs", d.Seconds())
	}

	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60

	if minutes > 0 {
		return fmt.Sprintf("%dm%ds", minutes, seconds)
	}

	return fmt.Sprintf("%.0fs", d.Seconds())
}

// Animation tick message for continuous updates
type SpinnerTickMsg struct{}

// GetSpinnerTickCmd returns a command for spinner animation
func GetSpinnerTickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return SpinnerTickMsg{}
	})
}
