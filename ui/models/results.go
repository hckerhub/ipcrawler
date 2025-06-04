package models

import (
	"ipcrawler/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
)

// ResultsModel handles the display of scan results
type ResultsModel struct {
	width  int
	height int

	// Results data
	results    []ScanResult
	selected   int
	showDetail bool
}

// ScanResult represents a single scan result
type ScanResult struct {
	IP        string
	Port      int
	Service   string
	Version   string
	Status    string
	Timestamp string
	Details   map[string]string
}

// NewResultsModel creates a new results model
func NewResultsModel() *ResultsModel {
	return &ResultsModel{
		results:    []ScanResult{},
		selected:   0,
		showDetail: false,
	}
}

// Init initializes the results model
func (m *ResultsModel) Init() tea.Cmd {
	return nil
}

// SetSize updates the model dimensions
func (m *ResultsModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// Update handles messages for the results model
func (m *ResultsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if m.selected < len(m.results)-1 {
				m.selected++
			}
		case "enter":
			m.showDetail = !m.showDetail
		case "esc":
			m.showDetail = false
		}
	}

	return m, nil
}

// View renders the results model
func (m *ResultsModel) View() string {
	title := styles.RenderTitle("Scan Results")

	if len(m.results) == 0 {
		content := "No scan results available yet.\n\nRun a scan from the IP Scanner to see results here."
		return styles.ContentStyle.Width(m.width).Height(m.height).Render(title + "\n\n" + content)
	}

	// TODO: Implement results table view
	content := "Results will be displayed here after scanning."
	return styles.ContentStyle.Width(m.width).Height(m.height).Render(title + "\n\n" + content)
}
