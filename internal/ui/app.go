package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"ipcrawler/internal/models"
	"ipcrawler/internal/scanner"
)

// Application states
type appState int

const (
	stateInput appState = iota
	stateScanning
	stateResults
	stateError
)

// Message types for Bubble Tea
type scanProgressMsg models.ScanProgress
type scanCompleteMsg *models.ScanResult
type scanErrorMsg error

// Model represents the application state
type Model struct {
	state  appState
	width  int
	height int

	// Input components
	ipInput        textinput.Model
	focused        bool
	aggressive     bool
	verbose        bool
	privilegeLevel scanner.PrivilegeLevel

	// Scanning components
	progress   progress.Model
	scanResult *models.ScanResult
	scanEngine *scanner.Engine

	// Navigation
	currentTab int
	tabNames   []string

	// Error handling
	err error

	// ASCII art
	asciiArt string
}

// Initialize creates a new model
func NewModel(privilegeLevel scanner.PrivilegeLevel) Model {
	// Create IP input
	input := textinput.New()
	input.Placeholder = "Enter target IP address (e.g., 10.10.10.1)"
	input.Focus()
	input.CharLimit = 15
	input.Width = 30

	// Create progress bar
	prog := progress.New(progress.WithDefaultGradient())

	// ASCII art for the header
	art := `
 ██╗██████╗  ██████╗██████╗  █████╗ ██╗    ██╗██╗     ███████╗██████╗ 
 ██║██╔══██╗██╔════╝██╔══██╗██╔══██╗██║    ██║██║     ██╔════╝██╔══██╗
 ██║██████╔╝██║     ██████╔╝███████║██║ █╗ ██║██║     █████╗  ██████╔╝
 ██║██╔═══╝ ██║     ██╔══██╗██╔══██║██║███╗██║██║     ██╔══╝  ██╔══██╗
 ██║██║     ╚██████╗██║  ██║██║  ██║╚███╔███╔╝███████╗███████╗██║  ██║
 ╚═╝╚═╝      ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝ ╚══╝╚══╝ ╚══════╝╚══════╝╚═╝  ╚═╝
                                                                       
     🎯 Advanced IP Scanner & Vulnerability Hunter 🎯
          Hack The Box Edition - Built with ❤️ in Go
                              v0.1 by hckerhub`

	return Model{
		state:          stateInput,
		ipInput:        input,
		focused:        true,
		progress:       prog,
		privilegeLevel: privilegeLevel,
		tabNames:       []string{"📋 Summary", "🔌 Ports", "🚨 Vulnerabilities", "📊 Details"},
		asciiArt:       art,
	}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case scanProgressMsg:
		return m.handleScanProgress(models.ScanProgress(msg))

	case scanCompleteMsg:
		return m.handleScanComplete(msg)

	case scanErrorMsg:
		return m.handleScanError(error(msg))

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}

	// Update input if we're in input state
	if m.state == stateInput {
		m.ipInput, cmd = m.ipInput.Update(msg)
	}

	return m, cmd
}

// View implements tea.Model
func (m Model) View() string {
	switch m.state {
	case stateInput:
		return m.viewInput()
	case stateScanning:
		return m.viewScanning()
	case stateResults:
		return m.viewResults()
	case stateError:
		return m.viewError()
	default:
		return "Unknown state"
	}
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.state {
	case stateInput:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.ipInput.Value() != "" {
				return m.startScan()
			}
		case "tab":
			m.aggressive = !m.aggressive
		case "ctrl+v":
			m.verbose = !m.verbose
		default:
			var cmd tea.Cmd
			m.ipInput, cmd = m.ipInput.Update(msg)
			return m, cmd
		}

	case stateResults:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "r":
			m.state = stateInput
			m.scanResult = nil
			m.ipInput.SetValue("")
			m.ipInput.Focus()
			return m, textinput.Blink
		case "left", "h":
			if m.currentTab > 0 {
				m.currentTab--
			}
		case "right", "l":
			if m.currentTab < len(m.tabNames)-1 {
				m.currentTab++
			}
		case "s":
			// TODO: Save results
			return m, nil
		}

	case stateScanning:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}

	case stateError:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "r", "enter":
			m.state = stateInput
			m.err = nil
			m.ipInput.Focus()
			return m, textinput.Blink
		}
	}

	return m, nil
}

func (m Model) startScan() (Model, tea.Cmd) {
	m.state = stateScanning
	m.scanEngine = scanner.NewEngine(m.ipInput.Value(), m.aggressive, m.verbose, false) // false = non-interactive (privileges already checked)
	// Set the privilege level we determined earlier
	m.scanEngine.SetPrivilegeLevel(m.privilegeLevel)

	return m, tea.Batch(
		m.progress.Init(),
		m.listenForProgress(),
		m.performScan(),
	)
}

func (m Model) listenForProgress() tea.Cmd {
	return func() tea.Msg {
		select {
		case prog := <-m.scanEngine.GetProgressChannel():
			return scanProgressMsg(prog)
		case <-time.After(100 * time.Millisecond):
			return nil
		}
	}
}

func (m Model) performScan() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(time.Time) tea.Msg {
		ctx := context.Background()
		result, err := m.scanEngine.StartFullScan(ctx)
		if err != nil {
			return scanErrorMsg(err)
		}
		return scanCompleteMsg(result)
	})
}

func (m Model) handleScanProgress(prog models.ScanProgress) (Model, tea.Cmd) {
	cmd := m.progress.SetPercent(prog.Progress)
	return m, tea.Batch(cmd, m.listenForProgress())
}

func (m Model) handleScanComplete(result *models.ScanResult) (Model, tea.Cmd) {
	m.state = stateResults
	m.scanResult = result
	return m, nil
}

func (m Model) handleScanError(err error) (Model, tea.Cmd) {
	m.state = stateError
	m.err = err
	return m, nil
}

func (m Model) viewInput() string {
	title := titleStyle.Render(m.asciiArt)

	var sections []string
	sections = append(sections, title)

	// Input section
	inputSection := fmt.Sprintf("%s\n%s\n\n%s",
		subtitleStyle.Render("Enter the target IP address to begin scanning"),
		m.ipInput.View(),
		m.renderOptions(),
	)
	sections = append(sections, panelStyle.Render(inputSection))

	// Help section
	helpText := helpStyle.Render("Press Enter to start scan • Tab to toggle aggressive mode • Ctrl+V for verbose • Ctrl+C to quit")
	sections = append(sections, helpText)

	return lipgloss.JoinVertical(lipgloss.Center, sections...)
}

func (m Model) renderOptions() string {
	aggressiveStatus := "❌"
	if m.aggressive {
		aggressiveStatus = "✅"
	}

	verboseStatus := "❌"
	if m.verbose {
		verboseStatus = "✅"
	}

	return fmt.Sprintf("Options: %s Aggressive Mode (Tab)  %s Verbose Output (Ctrl+V)",
		aggressiveStatus, verboseStatus)
}

func (m Model) viewScanning() string {
	title := titleStyle.Render("🔍 Scanning in Progress")

	progressBar := m.progress.View()
	status := statusStyle.Render(fmt.Sprintf("Scanning target: %s", m.ipInput.Value()))

	var sections []string
	sections = append(sections, title)
	sections = append(sections, status)
	sections = append(sections, progressBar)

	helpText := helpStyle.Render("Press Ctrl+C to cancel scan")
	sections = append(sections, helpText)

	return lipgloss.JoinVertical(lipgloss.Center, sections...)
}

func (m Model) viewResults() string {
	if m.scanResult == nil {
		return errorStyle.Render("No scan results available")
	}

	title := titleStyle.Render("📊 Scan Results")
	target := subtitleStyle.Render(fmt.Sprintf("Target: %s", m.scanResult.TargetIP))

	// Render tabs
	tabs := m.renderTabs()

	// Render content based on current tab
	var content string
	switch m.currentTab {
	case 0:
		content = m.renderSummary()
	case 1:
		content = m.renderPorts()
	case 2:
		content = m.renderVulnerabilities()
	case 3:
		content = m.renderDetails()
	}

	helpText := helpStyle.Render("← → Navigate tabs • R to rescan • S to save • Q to quit")

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		target,
		tabs,
		content,
		helpText,
	)
}

func (m Model) renderTabs() string {
	var tabs []string

	for i, name := range m.tabNames {
		if i == m.currentTab {
			tabs = append(tabs, activeButtonStyle.Render(name))
		} else {
			tabs = append(tabs, buttonStyle.Render(name))
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}

func (m Model) renderSummary() string {
	s := m.scanResult.Summary

	udpStatus := fmt.Sprintf("%d UDP", s.UDPOpenPorts)
	if s.UDPOpenPorts == 0 && len(m.scanResult.UDPPorts) == 0 {
		udpStatus = "UDP skipped (privileges required)"
	}

	content := fmt.Sprintf(`%s

📊 Scan Statistics:
• Total Open Ports: %d (%d TCP, %s)
• Scan Duration: %v
• Vulnerabilities Found: %d

🚨 Vulnerability Breakdown:
• Critical: %d
• High: %d  
• Medium: %d
• Low: %d

🎯 Interesting Services:
%s`,
		summaryTitleStyle.Render("Scan Summary"),
		s.TotalOpenPorts, s.TCPOpenPorts, udpStatus,
		m.scanResult.Duration.Truncate(time.Second),
		s.CriticalVulns+s.HighVulns+s.MediumVulns+s.LowVulns,
		s.CriticalVulns, s.HighVulns, s.MediumVulns, s.LowVulns,
		strings.Join(s.InterestingServices, ", "),
	)

	return summaryStyle.Render(content)
}

func (m Model) renderPorts() string {
	var sections []string

	// TCP Ports
	if len(m.scanResult.TCPPorts) > 0 {
		tcpContent := portHeaderStyle.Render("🔌 TCP Ports")
		for _, port := range m.scanResult.TCPPorts {
			style := openPortStyle
			if port.IsHighValuePort() {
				style = highValuePortStyle
			}
			tcpContent += "\n" + style.Render(fmt.Sprintf("• %s", port.GetPortDisplay()))
		}
		sections = append(sections, panelStyle.Render(tcpContent))
	}

	// UDP Ports
	if len(m.scanResult.UDPPorts) > 0 {
		udpContent := portHeaderStyle.Render("📡 UDP Ports")
		for _, port := range m.scanResult.UDPPorts {
			style := openPortStyle
			if port.IsHighValuePort() {
				style = highValuePortStyle
			}
			udpContent += "\n" + style.Render(fmt.Sprintf("• %s", port.GetPortDisplay()))
		}
		sections = append(sections, panelStyle.Render(udpContent))
	}

	if len(sections) == 0 {
		sections = append(sections, panelStyle.Render("No open ports found"))
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m Model) renderVulnerabilities() string {
	if len(m.scanResult.Vulnerabilities) == 0 {
		return panelStyle.Render("✅ No vulnerabilities detected")
	}

	var sections []string
	vulnContent := vulnHeaderStyle.Render("🚨 Vulnerabilities")

	for _, vuln := range m.scanResult.Vulnerabilities {
		style := GetSeverityStyle(vuln.Severity)
		vulnLine := fmt.Sprintf("• [%s] %s:%d - %s",
			vuln.Severity, vuln.Service, vuln.Port, vuln.Type)
		vulnContent += "\n" + style.Render(vulnLine)

		if len(vuln.Suggestions) > 0 {
			for _, suggestion := range vuln.Suggestions {
				vulnContent += "\n  " + lipgloss.NewStyle().Foreground(mutedColor).Render("→ "+suggestion)
			}
		}
		vulnContent += "\n"
	}

	sections = append(sections, panelStyle.Render(vulnContent))
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m Model) renderDetails() string {
	content := fmt.Sprintf(`%s

🎯 Target Information:
• IP Address: %s
• Scan Start: %s
• Scan End: %s
• Total Duration: %v

🔧 Scan Configuration:
• Aggressive Mode: %v
• Verbose Output: %v
• TCP Port Range: 1-65535
• UDP Top Ports: 17 most common

📈 Performance Metrics:
• Ports Scanned: %d
• Services Identified: %d
• Vulnerabilities Analyzed: %d`,
		summaryTitleStyle.Render("Detailed Information"),
		m.scanResult.TargetIP,
		m.scanResult.StartTime.Format("15:04:05"),
		m.scanResult.EndTime.Format("15:04:05"),
		m.scanResult.Duration.Truncate(time.Second),
		m.aggressive,
		m.verbose,
		65535+17, // TCP + UDP ports
		len(m.scanResult.TCPPorts)+len(m.scanResult.UDPPorts),
		len(m.scanResult.Vulnerabilities),
	)

	return summaryStyle.Render(content)
}

func (m Model) viewError() string {
	title := titleStyle.Render("❌ Error Occurred")
	errorMsg := errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	helpText := helpStyle.Render("Press R or Enter to retry • Q to quit")

	return lipgloss.JoinVertical(lipgloss.Center, title, errorMsg, helpText)
}

// StartTUI starts the TUI application
func StartTUI(privilegeLevel scanner.PrivilegeLevel) {
	model := NewModel(privilegeLevel)
	program := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
	}
}

// StartTUIWithTarget starts the TUI with a pre-filled target
func StartTUIWithTarget(target, outputFile string, verbose, aggressive bool, privilegeLevel scanner.PrivilegeLevel) {
	model := NewModel(privilegeLevel)
	model.ipInput.SetValue(target)
	model.verbose = verbose
	model.aggressive = aggressive

	// Auto-start scan
	model.state = stateScanning
	model.scanEngine = scanner.NewEngine(target, aggressive, verbose, false) // false = non-interactive
	model.scanEngine.SetPrivilegeLevel(privilegeLevel)

	program := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
	}
}
