package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"ipcrawler/scanner"
	"ipcrawler/ui/styles"
)

// AppState represents the current state of the application
type AppState int

const (
	StateIPInput AppState = iota
	StateScanning
	StateResults
)

// ScanType represents different scan types
type ScanType int

const (
	ScanTypeFull ScanType = iota
	ScanTypeCustom
)

// ResultTab represents different result view tabs
type ResultTab int

const (
	TabPorts ResultTab = iota
	TabServices
	TabWeb
	TabVulns
	TabCommands // New tab for showing executed commands
	TabSummary
)

// MainModel represents the main application state
type MainModel struct {
	state  AppState
	width  int
	height int
	ready  bool

	// IP Input
	ipInput   textinput.Model
	portInput textinput.Model
	targetIP  string

	// Scan selection
	selectedScan   ScanType
	customPorts    string
	aggressiveMode bool // Default to true (aggressive mode enabled by default)

	// Scan state
	scanning     bool
	scanResult   *UIScanResult
	scanError    string
	scanStatus   string
	scanProgress *ScanStatus

	// Results navigation
	activeTab ResultTab
	tabs      []TabInfo

	// Export status
	exportMessage string
	exportTime    time.Time

	// Scanner
	nmapScanner *scanner.NmapScanner

	// Key bindings
	keyMap KeyMap

	// Wordlist prompt
	wordlistPrompt          string
	showingWordlistPrompt   bool
	wordlistResponseChannel chan wordlistResponseMsg
}

// TabInfo represents tab information
type TabInfo struct {
	Name        string
	Description string
	ID          ResultTab
}

// UIScanResult represents scan results for the UI
type UIScanResult struct {
	IP              string
	ScanType        string
	AggressiveMode  bool
	Ports           []UIPortInfo
	Services        []UIServiceInfo
	WebResults      []UIWebInfo
	Vulnerabilities []UIVulnInfo
	Commands        []UICommandInfo // New field for tracking commands
	Timestamp       string
	Duration        string
	Status          string
	Summary         UISummaryInfo
}

// UIPortInfo represents port information
type UIPortInfo struct {
	Port     int
	Protocol string
	State    string
	Service  string
	Version  string
	Scripts  []string
}

// UIServiceInfo represents detailed service information
type UIServiceInfo struct {
	Port       int
	Service    string
	Version    string
	Product    string
	ExtraInfo  string
	OSType     string
	Method     string
	Confidence int
}

// UIWebInfo represents web discovery information for display
type UIWebInfo struct {
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

// UIVulnInfo represents vulnerability information
type UIVulnInfo struct {
	Port        int
	Service     string
	CVE         string
	Severity    string
	Description string
	Script      string
}

// UISummaryInfo represents scan summary
type UISummaryInfo struct {
	TotalPorts      int
	OpenPorts       int
	FilteredPorts   int
	ClosedPorts     int
	WebServices     int
	Vulnerabilities int
	ScanTime        string
	AggressiveMode  bool
}

// UICommandInfo represents executed command information
type UICommandInfo struct {
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

// KeyMap defines the key bindings for the application
type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Enter  key.Binding
	Quit   key.Binding
	Help   key.Binding
	Escape key.Binding
	Space  key.Binding
}

// Message types for tea.Cmd communication
type wordlistPromptMsg struct{}
type wordlistResponseMsg struct {
	continue_ bool
	wordlist  string
}

// DefaultKeyMap returns the default key bindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "move left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "move right"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "start scan/confirm"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Escape: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Space: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "toggle"),
		),
	}
}

// NewMainModel creates a new main model
func NewMainModel() *MainModel {
	// Create IP input
	ipInput := textinput.New()
	ipInput.Placeholder = "Enter target IP address (e.g., 192.168.1.1)"
	ipInput.Focus()
	ipInput.CharLimit = 50
	ipInput.Width = 40

	// Create port input
	portInput := textinput.New()
	portInput.Placeholder = "Enter port range (e.g., 1-1000 or 22,80,443)"
	portInput.CharLimit = 100
	portInput.Width = 40

	// Define tabs for results
	tabs := []TabInfo{
		{Name: "Ports", Description: "Open ports and services", ID: TabPorts},
		{Name: "Services", Description: "Detailed service information", ID: TabServices},
		{Name: "Web", Description: "Web discovery results", ID: TabWeb},
		{Name: "Vulns", Description: "Vulnerability assessment", ID: TabVulns},
		{Name: "Commands", Description: "Executed commands and tools", ID: TabCommands},
		{Name: "Summary", Description: "Scan summary and statistics", ID: TabSummary},
	}

	return &MainModel{
		state:                   StateIPInput,
		ipInput:                 ipInput,
		portInput:               portInput,
		selectedScan:            ScanTypeFull,
		aggressiveMode:          true, // Aggressive mode enabled by default
		activeTab:               TabPorts,
		tabs:                    tabs,
		nmapScanner:             scanner.NewNmapScanner(),
		keyMap:                  DefaultKeyMap(),
		wordlistResponseChannel: make(chan wordlistResponseMsg),
	}
}

// Init initializes the model
func (m *MainModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages and updates the model
func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

	case tea.KeyMsg:
		switch m.state {
		case StateIPInput:
			return m.updateIPInput(msg)
		case StateScanning:
			// Only allow quit during scanning
			if key.Matches(msg, m.keyMap.Quit) {
				return m, tea.Quit
			}
		case StateResults:
			return m.updateResults(msg)
		}

	case ScanStartedMsg:
		m.scanning = true
		m.scanStatus = "Starting aggressive scan..."
		m.scanProgress = NewScanStatus(msg.IP, msg.ScanType, msg.Aggressive)
		return m, tea.Batch(
			m.performScan(msg.ScanType, msg.IP, msg.Ports, msg.Aggressive),
			GetSpinnerTickCmd(),
		)

	case ScanCompletedMsg:
		m.scanning = false
		m.scanResult = msg.Result
		m.scanError = msg.Error
		m.state = StateResults
		m.scanProgress = nil
		return m, nil

	case ScanStatusMsg:
		m.scanStatus = msg.Status
		if m.scanProgress != nil {
			// Update progress based on status message
			if strings.Contains(msg.Status, "nmap") {
				m.scanProgress.UpdateStage(ProgressNmapScan, msg.Status)
			} else if strings.Contains(msg.Status, "web") {
				m.scanProgress.UpdateStage(ProgressWebDiscovery, msg.Status)
			} else if strings.Contains(msg.Status, "vuln") {
				m.scanProgress.UpdateStage(ProgressVulnAnalysis, msg.Status)
			}
		}
		return m, nil

	case SpinnerTickMsg:
		if m.scanning && m.scanProgress != nil {
			m.scanProgress.AdvanceAnimation()
			return m, GetSpinnerTickCmd()
		}
		return m, nil

	case ExportCompletedMsg:
		if msg.Error != "" {
			m.exportMessage = fmt.Sprintf("❌ Export failed: %s", msg.Error)
		} else {
			m.exportMessage = fmt.Sprintf("💾 Scan results saved to: %s", msg.Filepath)
		}
		m.exportTime = time.Now()

		// Clear the message after 5 seconds
		return m, tea.Tick(time.Second*5, func(t time.Time) tea.Msg {
			return ClearExportMsg{}
		})

	case ClearExportMsg:
		m.exportMessage = ""
		return m, nil

	case wordlistPromptMsg:
		return m.showWordlistPrompt()
	case wordlistResponseMsg:
		if msg.continue_ {
			return m.startWordlistDiscovery(msg.wordlist)
		}
	}

	// Note: Input updates are now handled in updateIPInput

	return m, tea.Batch(cmds...)
}

// updateIPInput handles updates when in IP input state
func (m *MainModel) updateIPInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Handle special keys first
	switch {
	case key.Matches(msg, m.keyMap.Quit):
		return m, tea.Quit

	case key.Matches(msg, m.keyMap.Up), key.Matches(msg, m.keyMap.Down):
		// Toggle between Full and Custom scan types
		if m.selectedScan == ScanTypeFull {
			m.selectedScan = ScanTypeCustom
			m.portInput.Focus()
			m.ipInput.Blur()
		} else {
			m.selectedScan = ScanTypeFull
			m.ipInput.Focus()
			m.portInput.Blur()
		}
		return m, nil

	case key.Matches(msg, m.keyMap.Space):
		// Toggle aggressive mode
		m.aggressiveMode = !m.aggressiveMode
		return m, nil

	case key.Matches(msg, m.keyMap.Enter):
		if m.ipInput.Value() != "" {
			m.targetIP = m.ipInput.Value()
			m.customPorts = m.portInput.Value()
			m.state = StateScanning

			// Start the appropriate scan
			var scanType string
			var ports string
			if m.selectedScan == ScanTypeFull {
				scanType = "full"
				ports = "" // Full scan doesn't need port specification
			} else {
				scanType = "custom"
				ports = strings.TrimSpace(m.customPorts)
				// For custom scan, ensure we have a port range
				// If empty, use default range instead of full scan
				if ports == "" {
					ports = "1-1000"
				}
			}

			return m, func() tea.Msg {
				return ScanStartedMsg{
					ScanType:   scanType,
					IP:         m.targetIP,
					Ports:      ports,
					Aggressive: m.aggressiveMode,
				}
			}
		}
		return m, nil
	}

	// Update the focused input with the key message
	if m.selectedScan == ScanTypeFull {
		// IP input is focused
		m.ipInput, cmd = m.ipInput.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		// Port input is focused
		m.portInput, cmd = m.portInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// updateResults handles updates when viewing results
func (m *MainModel) updateResults(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keyMap.Quit):
		return m, tea.Quit

	case msg.String() == "w":
		// Trigger wordlist discovery prompt if web results exist
		if m.scanResult != nil && len(m.scanResult.WebResults) > 0 {
			return m, func() tea.Msg {
				return wordlistPromptMsg{}
			}
		}

	case key.Matches(msg, m.keyMap.Escape):
		// Go back to IP input
		m.state = StateIPInput
		m.ipInput.SetValue("")
		m.ipInput.Focus()
		m.portInput.SetValue("")
		m.portInput.Blur()
		m.scanning = false
		m.scanResult = nil
		m.scanError = ""
		m.scanStatus = ""
		m.activeTab = TabPorts // Reset to first tab
		return m, nil

	case key.Matches(msg, m.keyMap.Left):
		// Navigate to previous tab
		if m.activeTab > 0 {
			m.activeTab--
		} else {
			m.activeTab = ResultTab(len(m.tabs) - 1) // Wrap to last tab
		}
		return m, nil

	case key.Matches(msg, m.keyMap.Right):
		// Navigate to next tab
		if int(m.activeTab) < len(m.tabs)-1 {
			m.activeTab++
		} else {
			m.activeTab = 0 // Wrap to first tab
		}
		return m, nil

	case msg.String() == "s":
		// Save results to file
		return m, func() tea.Msg {
			filepath, err := m.ExportResults()
			if err != nil {
				return ExportCompletedMsg{Error: err.Error()}
			}
			return ExportCompletedMsg{Filepath: filepath}
		}
	}

	return m, nil
}

// performScan executes the actual nmap scan with proper progress updates
func (m *MainModel) performScan(scanType, target, ports string, aggressive bool) tea.Cmd {
	return tea.Sequence(
		// Stage 1: Initialize and start port scanning
		func() tea.Msg {
			return ScanStatusMsg{Status: "Initializing nmap scan parameters..."}
		},
		func() tea.Msg {
			return ScanStatusMsg{Status: "Starting nmap port discovery..."}
		},

		// Stage 2: Actual scanning
		func() tea.Msg {
			// Check if nmap is available
			if !m.nmapScanner.IsInstalled() {
				return ScanCompletedMsg{
					Error: "nmap is not installed or not found in PATH. Please install nmap first.",
				}
			}

			startTime := time.Now()

			// Set up progress callback to update the UI
			m.nmapScanner.SetProgressCallback(func(stage int, toolName string, status string, complete bool) {
				// Convert stage numbers to ScanProgress constants
				var scanStage ScanProgress
				switch stage {
				case 1:
					scanStage = ProgressReconnaissance
				case 2:
					scanStage = ProgressNmapScan
				case 3:
					scanStage = ProgressWebDiscovery
				case 4:
					scanStage = ProgressAdvancedWebAnalysis
				case 5:
					scanStage = ProgressVulnAnalysis
				default:
					scanStage = ProgressNmapScan
				}

				// Update the progress in the UI
				if m.scanProgress != nil {
					if toolName != "" {
						m.scanProgress.UpdateSubStage(scanStage, toolName, status, complete)
					} else {
						m.scanProgress.UpdateStage(scanStage, status)
					}
				}
			})

			var result *scanner.ScanResult
			var err error

			// Perform the actual scan based on type and aggressive mode
			if aggressive {
				result, err = m.nmapScanner.AggressiveScan(target, ports)
			} else {
				switch scanType {
				case "full":
					result, err = m.nmapScanner.FullPortScan(target)
				case "custom":
					// For custom scan, always use ServiceScan with specified ports
					// If no ports specified, use a reasonable default range
					if ports == "" {
						ports = "1-1000"
					}
					result, err = m.nmapScanner.ServiceScan(target, ports)
				default:
					result, err = m.nmapScanner.QuickScan(target)
				}
			}

			if err != nil {
				return ScanCompletedMsg{
					Error: fmt.Sprintf("Port scan failed: %v", err),
				}
			}

			// Convert scanner result to UI result
			uiResult := m.convertScanResult(result, target, scanType, aggressive, startTime)

			return ScanCompletedMsg{
				Result: uiResult,
			}
		},
	)
}

// convertScanResult converts scanner.ScanResult to UIScanResult
func (m *MainModel) convertScanResult(result *scanner.ScanResult, target, scanType string, aggressive bool, startTime time.Time) *UIScanResult {
	// Convert scanner result to UI result
	uiResult := &UIScanResult{
		IP:              target,
		ScanType:        scanType,
		AggressiveMode:  aggressive,
		Timestamp:       startTime.Format("2006-01-02 15:04:05"),
		Duration:        time.Since(startTime).Round(time.Second).String(),
		Status:          "completed",
		Ports:           make([]UIPortInfo, len(result.Ports)),
		Services:        make([]UIServiceInfo, len(result.Services)),
		WebResults:      make([]UIWebInfo, len(result.WebResults)),
		Vulnerabilities: make([]UIVulnInfo, len(result.Vulnerabilities)),
		Commands:        make([]UICommandInfo, len(result.Commands)),
		Summary: UISummaryInfo{
			TotalPorts:      len(result.Ports),
			OpenPorts:       len(result.OpenPorts),
			FilteredPorts:   len(result.FilteredPorts),
			ClosedPorts:     len(result.ClosedPorts),
			WebServices:     len(result.WebResults),
			Vulnerabilities: len(result.Vulnerabilities),
			ScanTime:        time.Since(startTime).String(),
			AggressiveMode:  aggressive,
		},
	}

	// Convert ports
	for i, port := range result.Ports {
		uiResult.Ports[i] = UIPortInfo{
			Port:     port.Port,
			Protocol: port.Protocol,
			State:    port.State,
			Service:  port.Service,
			Version:  port.Version,
			Scripts:  port.Scripts,
		}
	}

	// Convert services
	for i, service := range result.Services {
		uiResult.Services[i] = UIServiceInfo{
			Port:       service.Port,
			Service:    service.Service,
			Version:    service.Version,
			Product:    service.Product,
			ExtraInfo:  service.ExtraInfo,
			OSType:     service.OSType,
			Method:     service.Method,
			Confidence: service.Confidence,
		}
	}

	// Convert web results
	for i, web := range result.WebResults {
		uiResult.WebResults[i] = UIWebInfo{
			URL:             web.URL,
			Port:            web.Port,
			Title:           web.Title,
			Server:          web.Server,
			Technologies:    web.Technologies,
			StatusCode:      web.StatusCode,
			ContentLength:   web.ContentLength,
			ResponseTime:    web.ResponseTime,
			Paths:           web.Paths,
			Subdomains:      web.Subdomains,
			Headers:         web.Headers,
			Cookies:         web.Cookies,
			SecurityHeaders: web.SecurityHeaders,
			Redirects:       web.Redirects,
			Forms:           web.Forms,
			JavaScriptFiles: web.JavaScriptFiles,
			CSSFiles:        web.CSSFiles,
			Images:          web.Images,
			ExternalLinks:   web.ExternalLinks,
			EmailAddresses:  web.EmailAddresses,
		}
	}

	// Convert vulnerabilities
	for i, vuln := range result.Vulnerabilities {
		uiResult.Vulnerabilities[i] = UIVulnInfo{
			Port:        vuln.Port,
			Service:     vuln.Service,
			CVE:         vuln.CVE,
			Severity:    vuln.Severity,
			Description: vuln.Description,
			Script:      vuln.Script,
		}
	}

	// Convert commands
	for i, command := range result.Commands {
		uiResult.Commands[i] = UICommandInfo{
			Tool:      command.Tool,
			Command:   command.Command,
			Args:      command.Args,
			StartTime: command.StartTime,
			EndTime:   command.EndTime,
			Duration:  command.Duration,
			ExitCode:  command.ExitCode,
			Output:    command.Output,
			Error:     command.Error,
			Stage:     command.Stage,
		}
	}

	return uiResult
}

// Message types
type ScanCompletedMsg struct {
	Result *UIScanResult
	Error  string
}

type ScanStatusMsg struct {
	Status string
}

// ClearExportMsg represents clearing the export message
type ClearExportMsg struct{}

// View renders the model
func (m *MainModel) View() string {
	if !m.ready {
		return "Initializing IPCrawler..."
	}

	if m.showingWordlistPrompt {
		return m.renderWordlistPrompt()
	}

	switch m.state {
	case StateIPInput:
		return m.renderIPInput()
	case StateScanning:
		return m.renderScanning()
	case StateResults:
		return m.renderResults()
	}

	return "Unknown state"
}

// renderIPInput renders the IP input screen
func (m *MainModel) renderIPInput() string {
	var b strings.Builder

	// Title
	title := styles.RenderTitle("IPCrawler")
	subtitle := styles.SubtitleStyle.Render("Advanced IP Analysis & Red Teaming tools")

	// Center the content
	availableHeight := m.height - 15
	topPadding := availableHeight / 2

	// Add vertical padding
	for i := 0; i < topPadding; i++ {
		b.WriteString("\n")
	}

	// IP Input
	ipLabel := "Target IP Address:"
	ipInputStyle := styles.InputFocusedStyle

	// Scan type options
	fullScanLabel := "Full Port Scan (all 65535 ports)"
	customScanLabel := "Custom Port Range"

	if m.selectedScan == ScanTypeFull {
		fullScanLabel = styles.SelectedStyle.Render("► " + fullScanLabel)
		customScanLabel = styles.UnselectedStyle.Render("  " + customScanLabel)
	} else {
		fullScanLabel = styles.UnselectedStyle.Render("  " + fullScanLabel)
		customScanLabel = styles.SelectedStyle.Render("► " + customScanLabel)
	}

	// Port input
	portInputStyle := styles.InputBlurredStyle
	portLabel := "Port Range (for custom scan):"
	if m.selectedScan == ScanTypeCustom {
		portInputStyle = styles.InputFocusedStyle
		portLabel = styles.HighlightStyle.Render("► " + portLabel)
		// Show helpful hint about default range
		if m.portInput.Value() == "" {
			portLabel += " " + styles.RenderInfo("(default: 1-1000)")
		}
	} else {
		portLabel = "  " + portLabel
	}

	// Aggressive mode toggle
	aggressiveLabel := "🚀 Aggressive Mode (min-rate 10k, optimal for CTF)"
	if m.aggressiveMode {
		aggressiveLabel = styles.SelectedStyle.Render("☑ " + aggressiveLabel)
	} else {
		aggressiveLabel = styles.UnselectedStyle.Render("☐ " + aggressiveLabel)
	}

	// Main content
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		subtitle,
		"",
		"",
		ipLabel,
		ipInputStyle.Render(m.ipInput.View()),
		"",
		"",
		styles.RenderInfo("Select scan type:"),
		"",
		fullScanLabel,
		customScanLabel,
		"",
		"",
		portLabel,
		portInputStyle.Render(m.portInput.View()),
		"",
		"",
		aggressiveLabel,
		"",
		"",
		lipgloss.NewStyle().Foreground(styles.MutedColor).Render("↑/↓: Select scan type • Space: Toggle aggressive mode • Enter: Start scan • q: Quit"),
	)

	// Center horizontally
	centeredContent := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render(content)

	b.WriteString(centeredContent)

	return b.String()
}

// renderScanning renders the scanning screen with modern animation
func (m *MainModel) renderScanning() string {
	if m.scanProgress != nil {
		return m.scanProgress.RenderLoadingScreen(m.width, m.height)
	}

	// Fallback to simple loading if progress tracker not available
	var b strings.Builder

	// Center the content
	availableHeight := m.height - 10
	topPadding := availableHeight / 2

	// Add vertical padding
	for i := 0; i < topPadding; i++ {
		b.WriteString("\n")
	}

	scanTypeDisplay := "Full Port Scan"
	if m.selectedScan == ScanTypeCustom {
		scanTypeDisplay = fmt.Sprintf("Custom Scan (%s)", m.customPorts)
	}
	if m.aggressiveMode {
		scanTypeDisplay += " [AGGRESSIVE MODE]"
	}

	// Main content
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		styles.RenderTitle("Scanning in Progress"),
		"",
		styles.RenderInfo(fmt.Sprintf("Target: %s", m.targetIP)),
		styles.RenderInfo(fmt.Sprintf("Scan Type: %s", scanTypeDisplay)),
		"",
		"",
		styles.LoadingStyle.Render("⟳ Running nmap scan..."),
		"",
		styles.LoadingStyle.Render(m.scanStatus),
		"",
		"",
		lipgloss.NewStyle().Foreground(styles.MutedColor).Render("Please wait... • q: Quit"),
	)

	// Center horizontally
	centeredContent := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render(content)

	b.WriteString(centeredContent)

	return b.String()
}

// renderResults renders the results screen with tabbed interface
func (m *MainModel) renderResults() string {
	var b strings.Builder

	// Header with scan info
	header := styles.RenderTitle(fmt.Sprintf("🎯 %s", m.targetIP))
	b.WriteString(header)
	b.WriteString("\n")

	// Scan info line
	var scanInfo string
	if m.scanResult != nil {
		scanInfo = fmt.Sprintf("Scan: %s", m.scanResult.ScanType)
		if m.scanResult.AggressiveMode {
			scanInfo += " [🚀 AGGRESSIVE]"
		}
		scanInfo += fmt.Sprintf(" • Duration: %s • Completed: %s",
			m.scanResult.Duration, m.scanResult.Timestamp)
	} else {
		scanInfo = "No scan data available"
	}

	b.WriteString(styles.RenderInfo(scanInfo))
	b.WriteString("\n\n")

	// Check for errors
	if m.scanError != "" {
		errorContent := styles.ErrorStyle.Render("✗ Scan Failed")
		errorContent += "\n\n" + m.scanError
		errorContent += "\n\n" + styles.RenderInfo("Common issues:")
		errorContent += "\n• Check if nmap is installed: sudo apt install nmap (Linux) or brew install nmap (macOS)"
		errorContent += "\n• Verify the target IP is reachable"
		errorContent += "\n• Some scans require sudo privileges"

		cardContent := styles.ContentStyle.
			Width(m.width - 4).
			Height(m.height - 8).
			Render(errorContent)

		b.WriteString(cardContent)
	} else if m.scanResult != nil {
		// Render tabs
		b.WriteString(m.renderTabs())
		b.WriteString("\n")

		// Render active tab content
		b.WriteString(m.renderActiveTabContent())
	} else {
		b.WriteString("No scan results available.")
	}

	return b.String()
}

// renderTabs renders the tab navigation header
func (m *MainModel) renderTabs() string {
	var tabs []string

	for i, tab := range m.tabs {
		style := styles.UnselectedStyle
		if ResultTab(i) == m.activeTab {
			style = styles.SelectedStyle
		}

		// Add count indicators
		count := m.getTabCount(tab.ID)
		tabText := tab.Name
		if count > 0 {
			tabText += fmt.Sprintf(" (%d)", count)
		}

		if ResultTab(i) == m.activeTab {
			tabs = append(tabs, style.Render(fmt.Sprintf("┌─ %s ─┐", tabText)))
		} else {
			tabs = append(tabs, style.Render(fmt.Sprintf("  %s  ", tabText)))
		}
	}

	tabsLine := strings.Join(tabs, "  ")

	// Add navigation hints and export message
	nav := lipgloss.NewStyle().
		Foreground(styles.MutedColor).
		Render(" ← → Navigate tabs • s: Save to file • w: Wordlist discovery • Esc: New scan • q: Quit")

	result := tabsLine + "\n" + nav

	// Add export message if present
	if m.exportMessage != "" {
		exportStyle := styles.SuccessStyle
		if strings.Contains(m.exportMessage, "failed") {
			exportStyle = styles.ErrorStyle
		}
		result += "\n" + exportStyle.Render(m.exportMessage)
	}

	return result
}

// renderActiveTabContent renders the content for the currently active tab
func (m *MainModel) renderActiveTabContent() string {
	if m.scanResult == nil {
		return "No scan results available"
	}

	contentHeight := m.height - 12 // Account for header, tabs, and footer
	contentWidth := m.width - 4

	var content string

	switch m.activeTab {
	case TabPorts:
		content = m.renderPortsTab()
	case TabServices:
		content = m.renderServicesTab()
	case TabWeb:
		content = m.renderWebTab()
	case TabVulns:
		content = m.renderVulnsTab()
	case TabCommands:
		content = m.renderCommandsTab()
	case TabSummary:
		content = m.renderSummaryTab()
	default:
		content = "Unknown tab"
	}

	return styles.ContentStyle.
		Width(contentWidth).
		Height(contentHeight).
		Render(content)
}

// getTabCount returns the count of items for a tab
func (m *MainModel) getTabCount(tabID ResultTab) int {
	if m.scanResult == nil {
		return 0
	}

	switch tabID {
	case TabPorts:
		return len(m.scanResult.Ports)
	case TabServices:
		return len(m.scanResult.Services)
	case TabWeb:
		return len(m.scanResult.WebResults)
	case TabVulns:
		return len(m.scanResult.Vulnerabilities)
	case TabCommands:
		return len(m.scanResult.Commands)
	case TabSummary:
		return 1 // Always show summary
	default:
		return 0
	}
}

// renderPortsTab renders the ports tab content
func (m *MainModel) renderPortsTab() string {
	var b strings.Builder

	b.WriteString(styles.RenderHeader("📡 Port Scan Results"))
	b.WriteString("\n\n")

	if len(m.scanResult.Ports) == 0 {
		b.WriteString("No open ports found.\n")
		b.WriteString("\nThis could mean:\n")
		b.WriteString("• All ports are closed or filtered\n")
		b.WriteString("• The target is not responding\n")
		b.WriteString("• A firewall is blocking the scan\n")
		return b.String()
	}

	// Table header
	header := fmt.Sprintf("%-8s %-8s %-8s %-15s %-20s %s",
		"PORT", "PROTO", "STATE", "SERVICE", "VERSION", "SCRIPTS")
	b.WriteString(styles.TableHeaderStyle.Render(header))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", 80))
	b.WriteString("\n")

	// Table rows
	for _, port := range m.scanResult.Ports {
		portLine := fmt.Sprintf("%-8d %-8s %-8s %-15s %-20s %s",
			port.Port,
			port.Protocol,
			port.State,
			port.Service,
			port.Version,
			strings.Join(port.Scripts, ", "),
		)

		// Highlight important ports
		if port.Port == 22 || port.Port == 80 || port.Port == 443 || port.Port == 21 || port.Port == 25 || port.Port == 3389 {
			portLine = styles.HighlightStyle.Render(portLine)
		}

		b.WriteString(fmt.Sprintf("%s\n", portLine))
	}

	return b.String()
}

// renderServicesTab renders the services tab content
func (m *MainModel) renderServicesTab() string {
	var b strings.Builder

	b.WriteString(styles.RenderHeader("🔧 Service Details"))
	b.WriteString("\n\n")

	if len(m.scanResult.Services) == 0 {
		b.WriteString("No detailed service information available.\n")
		b.WriteString("Try running with aggressive mode for more service details.\n")
		return b.String()
	}

	for _, service := range m.scanResult.Services {
		b.WriteString(styles.HighlightStyle.Render(fmt.Sprintf("Port %d: %s", service.Port, service.Service)))
		b.WriteString("\n")

		if service.Product != "" {
			b.WriteString(fmt.Sprintf("  Product: %s\n", service.Product))
		}
		if service.Version != "" {
			b.WriteString(fmt.Sprintf("  Version: %s\n", service.Version))
		}
		if service.ExtraInfo != "" {
			b.WriteString(fmt.Sprintf("  Extra Info: %s\n", service.ExtraInfo))
		}
		if service.OSType != "" {
			b.WriteString(fmt.Sprintf("  OS Type: %s\n", service.OSType))
		}
		if service.Confidence > 0 {
			b.WriteString(fmt.Sprintf("  Confidence: %d%%\n", service.Confidence))
		}
		b.WriteString("\n")
	}

	return b.String()
}

// renderWebTab renders the web discovery tab content
func (m *MainModel) renderWebTab() string {
	var b strings.Builder

	b.WriteString(styles.RenderHeader("🌐 Web Discovery Results"))
	b.WriteString("\n\n")

	if len(m.scanResult.WebResults) == 0 {
		b.WriteString("No web services discovered.\n")
		b.WriteString("Web discovery runs automatically when HTTP/HTTPS ports are found.\n")
		b.WriteString("\n")
		b.WriteString(styles.InfoStyle.Render("💡 Tip: "))
		b.WriteString("Try running an aggressive scan to enable comprehensive web discovery.\n")
		return b.String()
	}

	for i, web := range m.scanResult.WebResults {
		// Header with URL and basic info
		b.WriteString(styles.HighlightStyle.Render(fmt.Sprintf("%d. %s", i+1, web.URL)))
		b.WriteString("\n")

		// Basic information
		if web.Title != "" {
			b.WriteString(fmt.Sprintf("   📄 Title:        %s\n", web.Title))
		}
		if web.Server != "" {
			b.WriteString(fmt.Sprintf("   🖥️  Server:       %s\n", web.Server))
		}
		if web.StatusCode > 0 {
			statusStyle := styles.SuccessStyle // Default green
			if web.StatusCode >= 400 {
				statusStyle = styles.ErrorStyle
			} else if web.StatusCode >= 300 {
				statusStyle = styles.WarningStyle
			}
			b.WriteString(fmt.Sprintf("   📊 Status Code:  %s", statusStyle.Render(fmt.Sprintf("%d", web.StatusCode))))

			// Add response time and content length if available
			if web.ResponseTime > 0 {
				b.WriteString(fmt.Sprintf(" (%dms)", web.ResponseTime))
			}
			if web.ContentLength > 0 {
				b.WriteString(fmt.Sprintf(" [%d bytes]", web.ContentLength))
			}
			b.WriteString("\n")
		}

		// Technologies detected
		if len(web.Technologies) > 0 {
			b.WriteString(fmt.Sprintf("   🔧 Technologies: %s\n", styles.InfoStyle.Render(strings.Join(web.Technologies, ", "))))
		}

		// Security Headers Analysis
		if len(web.SecurityHeaders) > 0 {
			b.WriteString("   🛡️  Security Headers:\n")
			securityScore := 0
			totalChecks := 6 // HSTS, CSP, X-Frame-Options, X-Content-Type-Options, X-XSS-Protection, Referrer-Policy

			importantHeaders := []string{
				"Strict-Transport-Security", "Content-Security-Policy", "X-Frame-Options",
				"X-Content-Type-Options", "X-XSS-Protection", "Referrer-Policy",
			}

			for _, header := range importantHeaders {
				if value, exists := web.SecurityHeaders[header]; exists && value != "" {
					securityScore++
					b.WriteString(fmt.Sprintf("     ✅ %s: %s\n", header, value))
				} else {
					b.WriteString(fmt.Sprintf("     ❌ %s: Missing\n", header))
				}
			}

			scoreStyle := styles.ErrorStyle
			if securityScore >= 4 {
				scoreStyle = styles.SuccessStyle
			} else if securityScore >= 2 {
				scoreStyle = styles.WarningStyle
			}
			b.WriteString(fmt.Sprintf("     🏆 Security Score: %s\n", scoreStyle.Render(fmt.Sprintf("%d/%d", securityScore, totalChecks))))
		}

		// Cookies information
		if len(web.Cookies) > 0 {
			b.WriteString(fmt.Sprintf("   🍪 Cookies: %d found", len(web.Cookies)))
			// Show first few cookies for brevity
			if len(web.Cookies) <= 3 {
				b.WriteString(fmt.Sprintf(" (%s)", strings.Join(web.Cookies, ", ")))
			}
			b.WriteString("\n")
		}

		// Forms detected
		if web.Forms > 0 {
			formsStyle := styles.WarningStyle
			if web.Forms > 5 {
				formsStyle = styles.ErrorStyle
			}
			b.WriteString(fmt.Sprintf("   📝 Forms: %s\n", formsStyle.Render(fmt.Sprintf("%d forms detected", web.Forms))))
		}

		// Discovered paths - show them in a more organized way
		if len(web.Paths) > 0 {
			b.WriteString("   🔍 Discovered Paths:\n")

			// Group paths by type for better readability
			adminPaths := []string{}
			apiPaths := []string{}
			configPaths := []string{}
			otherPaths := []string{}

			for _, path := range web.Paths {
				pathLower := strings.ToLower(path)
				if strings.Contains(pathLower, "admin") || strings.Contains(pathLower, "manage") || strings.Contains(pathLower, "dashboard") {
					adminPaths = append(adminPaths, path)
				} else if strings.Contains(pathLower, "api") || strings.Contains(pathLower, "rest") || strings.Contains(pathLower, "graphql") || strings.Contains(pathLower, "endpoint") {
					apiPaths = append(apiPaths, path)
				} else if strings.Contains(pathLower, "config") || strings.Contains(pathLower, "settings") || strings.Contains(pathLower, ".php") || strings.Contains(pathLower, ".txt") || strings.Contains(pathLower, ".xml") {
					configPaths = append(configPaths, path)
				} else {
					otherPaths = append(otherPaths, path)
				}
			}

			if len(adminPaths) > 0 {
				adminStyle := styles.ErrorStyle
				b.WriteString(fmt.Sprintf("     🔐 Admin/Management: %s\n", adminStyle.Render(strings.Join(adminPaths, ", "))))
			}
			if len(apiPaths) > 0 {
				b.WriteString(fmt.Sprintf("     📡 APIs/Endpoints:   %s\n", styles.InfoStyle.Render(strings.Join(apiPaths, ", "))))
			}
			if len(configPaths) > 0 {
				configStyle := styles.WarningStyle
				b.WriteString(fmt.Sprintf("     ⚙️  Config/Files:     %s\n", configStyle.Render(strings.Join(configPaths, ", "))))
			}
			if len(otherPaths) > 0 {
				// Limit display of other paths to avoid clutter
				if len(otherPaths) > 8 {
					displayPaths := otherPaths[:8]
					b.WriteString(fmt.Sprintf("     📁 Other Paths:      %s... (+%d more)\n",
						strings.Join(displayPaths, ", "), len(otherPaths)-8))
				} else {
					b.WriteString(fmt.Sprintf("     📁 Other Paths:      %s\n", strings.Join(otherPaths, ", ")))
				}
			}
		}

		// JavaScript and CSS files
		if len(web.JavaScriptFiles) > 0 || len(web.CSSFiles) > 0 {
			b.WriteString("   📂 Static Resources:\n")
			if len(web.JavaScriptFiles) > 0 {
				jsCount := len(web.JavaScriptFiles)
				b.WriteString(fmt.Sprintf("     📜 JavaScript: %d files", jsCount))
				if jsCount <= 3 {
					b.WriteString(fmt.Sprintf(" (%s)", strings.Join(web.JavaScriptFiles, ", ")))
				}
				b.WriteString("\n")
			}
			if len(web.CSSFiles) > 0 {
				cssCount := len(web.CSSFiles)
				b.WriteString(fmt.Sprintf("     🎨 CSS: %d files", cssCount))
				if cssCount <= 3 {
					b.WriteString(fmt.Sprintf(" (%s)", strings.Join(web.CSSFiles, ", ")))
				}
				b.WriteString("\n")
			}
		}

		// External links and email addresses
		if len(web.ExternalLinks) > 0 {
			b.WriteString(fmt.Sprintf("   🔗 External Links: %d found", len(web.ExternalLinks)))
			if len(web.ExternalLinks) <= 3 {
				b.WriteString(fmt.Sprintf(" (%s)", strings.Join(web.ExternalLinks, ", ")))
			}
			b.WriteString("\n")
		}

		if len(web.EmailAddresses) > 0 {
			emailStyle := styles.InfoStyle
			if len(web.EmailAddresses) > 5 {
				emailStyle = styles.WarningStyle
			}
			b.WriteString(fmt.Sprintf("   📧 Email Addresses: %s\n",
				emailStyle.Render(fmt.Sprintf("%d found (%s)", len(web.EmailAddresses), strings.Join(web.EmailAddresses, ", ")))))
		}

		// Subdomains
		if len(web.Subdomains) > 0 {
			if len(web.Subdomains) > 5 {
				displaySubs := web.Subdomains[:5]
				b.WriteString(fmt.Sprintf("   🌐 Subdomains:   %s... (+%d more)\n",
					styles.InfoStyle.Render(strings.Join(displaySubs, ", ")), len(web.Subdomains)-5))
			} else {
				b.WriteString(fmt.Sprintf("   🌐 Subdomains:   %s\n", styles.InfoStyle.Render(strings.Join(web.Subdomains, ", "))))
			}
		}

		// Add some spacing between entries
		b.WriteString("\n")
	}

	// Show discovery method information
	b.WriteString(styles.InfoStyle.Render("🔍 Discovery Methods Used:"))
	b.WriteString("\n")
	b.WriteString("• robots.txt and sitemap.xml analysis\n")
	b.WriteString("• HTML link extraction and technology detection\n")
	b.WriteString("• Directory brute forcing (ffuf, feroxbuster, gobuster)\n")
	b.WriteString("• Security header analysis\n")
	b.WriteString("• Form and input field detection\n")
	b.WriteString("• Static resource enumeration\n")
	b.WriteString("\n")
	b.WriteString(styles.InfoStyle.Render("💡 Tip: "))
	b.WriteString("Press 'w' for wordlist-based discovery, 's' to save detailed results\n")

	return b.String()
}

// renderVulnsTab renders the vulnerabilities tab content
func (m *MainModel) renderVulnsTab() string {
	var b strings.Builder

	b.WriteString(styles.RenderHeader("🛡️ Vulnerability Assessment"))
	b.WriteString("\n\n")

	if len(m.scanResult.Vulnerabilities) == 0 {
		b.WriteString("No vulnerabilities detected.\n")
		b.WriteString("Use aggressive mode to run vulnerability scripts.\n")
		return b.String()
	}

	for _, vuln := range m.scanResult.Vulnerabilities {
		severity := vuln.Severity
		if severity == "HIGH" || severity == "CRITICAL" {
			severity = styles.ErrorStyle.Render(severity)
		} else if severity == "MEDIUM" {
			severity = styles.WarningStyle.Render(severity)
		} else {
			severity = styles.InfoStyle.Render(severity)
		}

		b.WriteString(fmt.Sprintf("🚨 Port %d (%s) - %s", vuln.Port, vuln.Service, severity))
		b.WriteString("\n")

		if vuln.CVE != "" {
			b.WriteString(fmt.Sprintf("  CVE: %s\n", vuln.CVE))
		}
		if vuln.Description != "" {
			b.WriteString(fmt.Sprintf("  Description: %s\n", vuln.Description))
		}
		if vuln.Script != "" {
			b.WriteString(fmt.Sprintf("  Script: %s\n", vuln.Script))
		}
		b.WriteString("\n")
	}

	return b.String()
}

// renderCommandsTab renders the commands/tools execution tab content
func (m *MainModel) renderCommandsTab() string {
	var b strings.Builder

	b.WriteString(styles.RenderHeader("⚙️ Executed Commands & Tools"))
	b.WriteString("\n\n")

	if len(m.scanResult.Commands) == 0 {
		b.WriteString("No command execution data available.\n")
		b.WriteString("Command tracking is enabled for aggressive scans and tool execution.\n")
		b.WriteString("\n")
		b.WriteString(styles.InfoStyle.Render("💡 Note: "))
		b.WriteString("Commands are logged when tools like nmap, ffuf, amass, etc. are executed.\n")
		return b.String()
	}

	// Group commands by stage for better organization
	stageCommands := make(map[string][]UICommandInfo)
	for _, cmd := range m.scanResult.Commands {
		stage := cmd.Stage
		if stage == "" {
			stage = "Other"
		}
		stageCommands[stage] = append(stageCommands[stage], cmd)
	}

	// Display commands by stage
	stageOrder := []string{"Reconnaissance", "Port Scanning", "Web Discovery", "Advanced Analysis", "Vulnerability Scanning", "Other"}

	for _, stage := range stageOrder {
		commands, exists := stageCommands[stage]
		if !exists || len(commands) == 0 {
			continue
		}

		b.WriteString(styles.HighlightStyle.Render(fmt.Sprintf("📋 %s", stage)))
		b.WriteString("\n")

		for i, cmd := range commands {
			// Command header with tool name and status
			statusIcon := "✅"
			if cmd.ExitCode != 0 {
				statusIcon = "❌"
			}

			b.WriteString(fmt.Sprintf("  %s %d. %s", statusIcon, i+1, styles.InfoStyle.Render(cmd.Tool)))
			if cmd.Duration != "" {
				b.WriteString(fmt.Sprintf(" (%s)", cmd.Duration))
			}
			b.WriteString("\n")

			// Command details
			fullCommand := cmd.Command
			if len(cmd.Args) > 0 {
				fullCommand += " " + strings.Join(cmd.Args, " ")
			}

			// Truncate very long commands for display
			if len(fullCommand) > 80 {
				fullCommand = fullCommand[:77] + "..."
			}

			b.WriteString(fmt.Sprintf("     💻 Command: %s\n", fullCommand))

			if cmd.ExitCode != 0 {
				b.WriteString(fmt.Sprintf("     ⚠️  Exit Code: %s\n", styles.ErrorStyle.Render(fmt.Sprintf("%d", cmd.ExitCode))))
			}

			if cmd.Error != "" {
				errorMsg := cmd.Error
				if len(errorMsg) > 100 {
					errorMsg = errorMsg[:97] + "..."
				}
				b.WriteString(fmt.Sprintf("     🚫 Error: %s\n", styles.ErrorStyle.Render(errorMsg)))
			}

			// Show output summary for successful commands
			if cmd.ExitCode == 0 && cmd.Output != "" {
				outputLines := strings.Split(cmd.Output, "\n")
				nonEmptyLines := 0
				for _, line := range outputLines {
					if strings.TrimSpace(line) != "" {
						nonEmptyLines++
					}
				}
				if nonEmptyLines > 0 {
					b.WriteString(fmt.Sprintf("     📄 Output: %d lines of data\n", nonEmptyLines))
				}
			}

			b.WriteString("\n")
		}
	}

	// Summary section
	totalCommands := len(m.scanResult.Commands)
	successfulCommands := 0
	failedCommands := 0
	totalDuration := time.Duration(0)

	for _, cmd := range m.scanResult.Commands {
		if cmd.ExitCode == 0 {
			successfulCommands++
		} else {
			failedCommands++
		}

		if !cmd.StartTime.IsZero() && !cmd.EndTime.IsZero() {
			totalDuration += cmd.EndTime.Sub(cmd.StartTime)
		}
	}

	b.WriteString(styles.HighlightStyle.Render("📊 Execution Summary"))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("Total Commands: %d\n", totalCommands))
	b.WriteString(fmt.Sprintf("Successful: %s\n", styles.SuccessStyle.Render(fmt.Sprintf("%d", successfulCommands))))
	if failedCommands > 0 {
		b.WriteString(fmt.Sprintf("Failed: %s\n", styles.ErrorStyle.Render(fmt.Sprintf("%d", failedCommands))))
	}
	if totalDuration > 0 {
		b.WriteString(fmt.Sprintf("Total Tool Execution Time: %s\n", totalDuration.Round(time.Second)))
	}

	b.WriteString("\n")
	b.WriteString(styles.InfoStyle.Render("💡 Tip: "))
	b.WriteString("Use 's' to save detailed command logs with full output to file\n")

	return b.String()
}

// renderSummaryTab renders the summary tab content
func (m *MainModel) renderSummaryTab() string {
	var b strings.Builder

	b.WriteString(styles.RenderHeader("📊 Scan Summary"))
	b.WriteString("\n\n")

	summary := m.scanResult.Summary

	// Scan details
	b.WriteString(styles.HighlightStyle.Render("Scan Information"))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("Target: %s\n", m.scanResult.IP))
	b.WriteString(fmt.Sprintf("Scan Type: %s", m.scanResult.ScanType))
	if m.scanResult.AggressiveMode {
		b.WriteString(" (Aggressive Mode)")
	}
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("Duration: %s\n", m.scanResult.Duration))
	b.WriteString(fmt.Sprintf("Completed: %s\n", m.scanResult.Timestamp))
	b.WriteString("\n")

	// Port statistics
	b.WriteString(styles.HighlightStyle.Render("Port Statistics"))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("Total Ports Scanned: %d\n", summary.TotalPorts))
	b.WriteString(fmt.Sprintf("Open Ports: %s\n", styles.SuccessStyle.Render(fmt.Sprintf("%d", summary.OpenPorts))))
	b.WriteString(fmt.Sprintf("Filtered Ports: %s\n", styles.WarningStyle.Render(fmt.Sprintf("%d", summary.FilteredPorts))))
	b.WriteString(fmt.Sprintf("Closed Ports: %d\n", summary.ClosedPorts))
	b.WriteString("\n")

	// Service summary
	b.WriteString(styles.HighlightStyle.Render("Discovery Summary"))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("Services Identified: %d\n", len(m.scanResult.Services)))
	b.WriteString(fmt.Sprintf("Web Services: %s\n", styles.InfoStyle.Render(fmt.Sprintf("%d", summary.WebServices))))
	b.WriteString(fmt.Sprintf("Vulnerabilities: %s\n", styles.ErrorStyle.Render(fmt.Sprintf("%d", summary.Vulnerabilities))))
	b.WriteString("\n")

	// Next steps
	b.WriteString(styles.HighlightStyle.Render("Recommended Next Steps"))
	b.WriteString("\n")
	if summary.WebServices > 0 {
		b.WriteString("• Investigate web services with directory busting\n")
		b.WriteString("• Check for default credentials\n")
		b.WriteString("• Look for web vulnerabilities\n")
	}
	if summary.OpenPorts > 0 {
		b.WriteString("• Test open services for misconfigurations\n")
		b.WriteString("• Attempt banner grabbing for version info\n")
	}
	if summary.Vulnerabilities > 0 {
		b.WriteString("• Investigate identified vulnerabilities\n")
		b.WriteString("• Check exploit databases\n")
	}

	return b.String()
}

// showWordlistPrompt shows a prompt asking if user wants to continue with wordlist-based discovery
func (m *MainModel) showWordlistPrompt() (tea.Model, tea.Cmd) {
	// Create a simple prompt model for wordlist discovery
	prompt := `
🔍 Intelligent URL discovery completed!

Would you like to continue with wordlist-based discovery?
This will use dictionary files to find additional paths.

Common wordlist locations:
📁 Kali Linux:     /usr/share/wordlists/dirb/common.txt
                   /usr/share/wordlists/dirbuster/directory-list-2.3-medium.txt
                   /usr/share/seclists/Discovery/Web-Content/big.txt

🎯 Hack The Box:   /opt/useful/SecLists/Discovery/Web-Content/directory-list-2.3-medium.txt
                   /opt/useful/SecLists/Discovery/Web-Content/big.txt
                   /opt/useful/SecLists/Discovery/Web-Content/common.txt

Press 'y' to continue with wordlist discovery
Press 'n' to skip
Press 'c' to specify custom wordlist path
`

	m.wordlistPrompt = prompt
	m.showingWordlistPrompt = true

	return m, m.waitForWordlistResponse()
}

// waitForWordlistResponse waits for user response to wordlist prompt
func (m *MainModel) waitForWordlistResponse() tea.Cmd {
	return func() tea.Msg {
		// This is a simplified version - in practice you'd want a proper input handler
		// For now, we'll just return a response to continue the flow
		return wordlistResponseMsg{continue_: false, wordlist: ""}
	}
}

// startWordlistDiscovery begins wordlist-based directory discovery
func (m *MainModel) startWordlistDiscovery(wordlistPath string) (tea.Model, tea.Cmd) {
	m.showingWordlistPrompt = false
	m.scanStatus = "Starting wordlist-based discovery..."

	// Start wordlist discovery in background
	return m, func() tea.Msg {
		// Here you would implement the actual wordlist discovery
		// For now, we'll just simulate it
		time.Sleep(time.Second * 2)
		return ScanStatusMsg{Status: "Wordlist discovery complete"}
	}
}

// renderWordlistPrompt renders the wordlist prompt
func (m *MainModel) renderWordlistPrompt() string {
	return m.wordlistPrompt
}
