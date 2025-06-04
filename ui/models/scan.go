package models

import (
	"fmt"
	"strings"

	"ipcrawler/ui/styles"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ScanModel handles the IP scanning interface
type ScanModel struct {
	width  int
	height int

	// Input fields
	ipInput   textinput.Model
	portInput textinput.Model

	// Scan options
	scanType       string
	scanTypes      []string
	selectedOption int

	// State
	scanning    bool
	results     []string
	currentStep string
	progress    float64

	// Focus state
	focused int
	inputs  []textinput.Model
}

// Init initializes the scan model
func (m *ScanModel) Init() tea.Cmd {
	return nil
}

// NewScanModel creates a new scan model
func NewScanModel() *ScanModel {
	// Create input fields
	ipInput := textinput.New()
	ipInput.Placeholder = "Enter IP address (e.g., 192.168.1.1)"
	ipInput.Focus()
	ipInput.CharLimit = 50
	ipInput.Width = 30

	portInput := textinput.New()
	portInput.Placeholder = "Port range (e.g., 1-1000 or 22,80,443)"
	portInput.CharLimit = 100
	portInput.Width = 30

	scanTypes := []string{
		"Quick Scan",
		"Full Port Scan",
		"Stealth Scan",
		"Service Detection",
		"OS Detection",
		"Vulnerability Scan",
		"Web Scan (Port 80/443)",
		"SSH Scan (Port 22)",
	}

	return &ScanModel{
		ipInput:   ipInput,
		portInput: portInput,
		inputs:    []textinput.Model{ipInput, portInput},
		scanTypes: scanTypes,
		scanType:  scanTypes[0],
		focused:   0,
	}
}

// SetSize updates the model dimensions
func (m *ScanModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// Update handles messages for the scan model
func (m *ScanModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" {
				if m.focused < len(m.inputs) {
					// Move to next input or start scan
					if m.focused == len(m.inputs)-1 {
						return m, m.startScan()
					}
					m.focused++
					m.updateFocus()
				} else {
					// Handle scan type selection or start scan
					if m.selectedOption < len(m.scanTypes)-1 {
						m.selectedOption++
					} else {
						return m, m.startScan()
					}
				}
				return m, nil
			}

			if s == "up" || s == "shift+tab" {
				if m.focused > 0 {
					m.focused--
					m.updateFocus()
				} else if m.selectedOption > 0 {
					m.selectedOption--
				}
			} else if s == "down" || s == "tab" {
				if m.focused < len(m.inputs)-1 {
					m.focused++
					m.updateFocus()
				} else if m.selectedOption < len(m.scanTypes)-1 {
					m.selectedOption++
				}
			}
		}
	}

	// Update inputs
	for i := range m.inputs {
		m.inputs[i], cmd = m.inputs[i].Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// updateFocus updates the focus state of inputs
func (m *ScanModel) updateFocus() {
	for i := 0; i < len(m.inputs); i++ {
		if i == m.focused {
			m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
}

// startScan initiates the scanning process
func (m *ScanModel) startScan() tea.Cmd {
	m.scanning = true
	m.currentStep = "Initializing scan..."
	m.progress = 0.0

	return func() tea.Msg {
		// This would trigger the actual scanning logic
		return ScanStartedMsg{
			IP:       m.ipInput.Value(),
			Ports:    m.portInput.Value(),
			ScanType: m.scanTypes[m.selectedOption],
		}
	}
}

// ScanStartedMsg represents a scan start message
type ScanStartedMsg struct {
	IP         string
	Ports      string
	ScanType   string
	Aggressive bool
}

// View renders the scan model
func (m *ScanModel) View() string {
	var b strings.Builder

	// Title
	title := styles.RenderTitle("IP Scanner")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Input section
	b.WriteString(styles.RenderHeader("Target Configuration"))
	b.WriteString("\n\n")

	// IP input
	ipLabel := "Target IP Address:"
	if m.focused == 0 {
		ipLabel = styles.HighlightStyle.Render("► " + ipLabel)
	} else {
		ipLabel = "  " + ipLabel
	}
	b.WriteString(ipLabel)
	b.WriteString("\n")

	var ipStyle lipgloss.Style
	if m.focused == 0 {
		ipStyle = styles.InputFocusedStyle
	} else {
		ipStyle = styles.InputBlurredStyle
	}
	b.WriteString(ipStyle.Render(m.ipInput.View()))
	b.WriteString("\n\n")

	// Port input
	portLabel := "Port Range:"
	if m.focused == 1 {
		portLabel = styles.HighlightStyle.Render("► " + portLabel)
	} else {
		portLabel = "  " + portLabel
	}
	b.WriteString(portLabel)
	b.WriteString("\n")

	var portStyle lipgloss.Style
	if m.focused == 1 {
		portStyle = styles.InputFocusedStyle
	} else {
		portStyle = styles.InputBlurredStyle
	}
	b.WriteString(portStyle.Render(m.portInput.View()))
	b.WriteString("\n\n")

	// Scan type selection
	b.WriteString(styles.RenderHeader("Scan Type"))
	b.WriteString("\n\n")

	for i, scanType := range m.scanTypes {
		if i == m.selectedOption {
			b.WriteString(styles.SelectedStyle.Render("► " + scanType))
		} else {
			b.WriteString(styles.UnselectedStyle.Render("  " + scanType))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")

	// Scan status
	if m.scanning {
		b.WriteString(styles.RenderHeader("Scan Status"))
		b.WriteString("\n\n")
		b.WriteString(styles.LoadingStyle.Render(m.currentStep))
		b.WriteString("\n")

		// Progress bar (simple text-based)
		progressWidth := 40
		filled := int(m.progress * float64(progressWidth))
		bar := strings.Repeat("█", filled) + strings.Repeat("░", progressWidth-filled)
		b.WriteString(fmt.Sprintf("[%s] %.1f%%", bar, m.progress*100))
		b.WriteString("\n\n")
	}

	// Instructions
	instructions := []string{
		"Tab/Shift+Tab: Navigate",
		"Enter: Next/Start Scan",
		"↑/↓: Select scan type",
		"Esc: Back to menu",
		"q: Quit",
	}

	b.WriteString(styles.RenderInfo("Controls:"))
	b.WriteString("\n")
	for _, instruction := range instructions {
		b.WriteString(fmt.Sprintf("  %s\n", instruction))
	}

	content := b.String()
	return styles.ContentStyle.Width(m.width).Height(m.height).Render(content)
}
