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
	selectedScan ScanType
	customPorts  string

	// Scan state
	scanning   bool
	scanResult *UIScanResult
	scanError  string
	scanStatus string

	// Scanner
	nmapScanner *scanner.NmapScanner

	// Key bindings
	keyMap KeyMap
}

// UIScanResult represents scan results for display
type UIScanResult struct {
	IP        string
	ScanType  string
	Ports     []UIPortInfo
	Timestamp string
	Duration  string
	Status    string
	Details   string
}

// UIPortInfo represents port information
type UIPortInfo struct {
	Port     int
	Protocol string
	State    string
	Service  string
	Version  string
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

	return &MainModel{
		state:        StateIPInput,
		ipInput:      ipInput,
		portInput:    portInput,
		selectedScan: ScanTypeFull,
		nmapScanner:  scanner.NewNmapScanner(),
		keyMap:       DefaultKeyMap(),
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
		m.scanStatus = "Starting scan..."
		return m, m.performScan(msg.ScanType, msg.IP, msg.Ports)

	case ScanCompletedMsg:
		m.scanning = false
		m.scanResult = msg.Result
		m.scanError = msg.Error
		m.state = StateResults
		return m, nil

	case ScanStatusMsg:
		m.scanStatus = msg.Status
		return m, nil
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
					ScanType: scanType,
					IP:       m.targetIP,
					Ports:    ports,
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
		return m, nil
	}

	return m, nil
}

// performScan executes the actual nmap scan
func (m *MainModel) performScan(scanType, target, ports string) tea.Cmd {
	return func() tea.Msg {
		startTime := time.Now()

		// Update status
		tea.Printf("Starting %s scan on %s...\n", scanType, target)

		var result *scanner.ScanResult
		var err error

		// Check if nmap is available
		if !m.nmapScanner.IsInstalled() {
			return ScanCompletedMsg{
				Error: "nmap is not installed or not found in PATH. Please install nmap first.",
			}
		}

		// Perform the actual scan
		switch scanType {
		case "full":
			result, err = m.nmapScanner.FullPortScan(target)
		case "custom":
			// For custom scan, always use ServiceScan with specified ports
			// If no ports specified, use a reasonable default range
			if ports == "" {
				ports = "1-2000"
			}
			result, err = m.nmapScanner.ServiceScan(target, ports)
		default:
			result, err = m.nmapScanner.QuickScan(target)
		}

		if err != nil {
			return ScanCompletedMsg{
				Error: fmt.Sprintf("Scan failed: %v", err),
			}
		}

		// Convert scanner result to UI result
		uiResult := &UIScanResult{
			IP:        target,
			ScanType:  scanType,
			Timestamp: startTime.Format("2006-01-02 15:04:05"),
			Duration:  time.Since(startTime).Round(time.Second).String(),
			Status:    "completed",
			Ports:     make([]UIPortInfo, len(result.Ports)),
		}

		// Convert ports
		for i, port := range result.Ports {
			uiResult.Ports[i] = UIPortInfo{
				Port:     port.Port,
				Protocol: port.Protocol,
				State:    port.State,
				Service:  port.Service,
				Version:  port.Version,
			}
		}

		return ScanCompletedMsg{
			Result: uiResult,
		}
	}
}

// Message types
type ScanCompletedMsg struct {
	Result *UIScanResult
	Error  string
}

type ScanStatusMsg struct {
	Status string
}

// View renders the model
func (m *MainModel) View() string {
	if !m.ready {
		return "Initializing IPCrawler..."
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
	subtitle := styles.SubtitleStyle.Render("Advanced IP Analysis & Penetration Testing")

	// Center the content
	availableHeight := m.height - 15
	topPadding := availableHeight / 2

	// Add vertical padding
	for i := 0; i < topPadding; i++ {
		b.WriteString("\n")
	}

	// IP Input
	ipLabel := "Target IP Address:"
	if m.selectedScan == ScanTypeFull {
		ipLabel = styles.HighlightStyle.Render("► " + ipLabel)
	} else {
		ipLabel = "  " + ipLabel
	}

	ipInputStyle := styles.InputBlurredStyle
	if m.selectedScan == ScanTypeFull {
		ipInputStyle = styles.InputFocusedStyle
	}

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
		lipgloss.NewStyle().Foreground(styles.MutedColor).Render("↑/↓: Select scan type • Enter: Start scan • q: Quit"),
	)

	// Center horizontally
	centeredContent := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render(content)

	b.WriteString(centeredContent)

	return b.String()
}

// renderScanning renders the scanning screen
func (m *MainModel) renderScanning() string {
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

// renderResults renders the results screen
func (m *MainModel) renderResults() string {
	var b strings.Builder

	// Header
	header := styles.RenderHeader(fmt.Sprintf("Scan Results: %s", m.targetIP))
	b.WriteString(header)
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
		// Show results
		result := m.scanResult

		content := styles.SuccessStyle.Render("✓ Scan completed")
		content += "\n\n"
		content += fmt.Sprintf("Scan Type: %s\n", result.ScanType)
		content += fmt.Sprintf("Duration: %s\n", result.Duration)
		content += fmt.Sprintf("Ports found: %d\n\n", len(result.Ports))

		// Show port details
		if len(result.Ports) > 0 {
			content += styles.RenderHeader("Open Ports:")
			content += "\n\n"

			// Table header
			header := fmt.Sprintf("%-8s %-10s %-8s %-15s %s",
				"PORT", "PROTOCOL", "STATE", "SERVICE", "VERSION")
			content += styles.TableHeaderStyle.Render(header)
			content += "\n"

			// Table rows
			for i, port := range result.Ports {
				if i >= 20 { // Limit display
					content += fmt.Sprintf("... and %d more ports\n", len(result.Ports)-20)
					break
				}

				portLine := fmt.Sprintf("%-8d %-10s %-8s %-15s %s",
					port.Port,
					port.Protocol,
					port.State,
					port.Service,
					port.Version,
				)

				// Highlight important ports
				if port.Port == 22 || port.Port == 80 || port.Port == 443 || port.Port == 21 || port.Port == 25 {
					portLine = styles.HighlightStyle.Render(portLine)
				}

				content += fmt.Sprintf("  %s\n", portLine)
			}
		} else {
			content += "No open ports found.\n"
			content += "\nThis could mean:\n"
			content += "• All ports are closed or filtered\n"
			content += "• The target is not responding\n"
			content += "• A firewall is blocking the scan\n"
		}

		cardContent := styles.ContentStyle.
			Width(m.width - 4).
			Height(m.height - 8).
			Render(content)

		b.WriteString(cardContent)
	}

	// Footer
	b.WriteString("\n\n")
	controls := "Esc: New scan • q: Quit"
	footer := lipgloss.NewStyle().
		Foreground(styles.MutedColor).
		Italic(true).
		Width(m.width).
		Align(lipgloss.Center).
		Render(controls)
	b.WriteString(footer)

	return b.String()
}
