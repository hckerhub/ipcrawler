package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Color palette inspired by cyberpunk/hacker aesthetics
	primaryColor   = lipgloss.Color("#00FF00") // Matrix green
	secondaryColor = lipgloss.Color("#0080FF") // Cyber blue
	accentColor    = lipgloss.Color("#FF6600") // Orange
	errorColor     = lipgloss.Color("#FF0000") // Red
	warningColor   = lipgloss.Color("#FFFF00") // Yellow
	successColor   = lipgloss.Color("#00FF00") // Green
	mutedColor     = lipgloss.Color("#666666") // Gray
	backgroundColor = lipgloss.Color("#0A0A0A") // Dark background

	// Base styles
	baseStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(primaryColor).
		Padding(1).
		Margin(1)

	// Header styles
	titleStyle = lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Align(lipgloss.Center).
		MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
		Foreground(secondaryColor).
		Italic(true).
		Align(lipgloss.Center)

	// Panel styles
	panelStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(1).
		Margin(0, 1)

	activePanelStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(accentColor).
		Padding(1).
		Margin(0, 1).
		Bold(true)

	// Progress bar styles
	progressBarStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(secondaryColor).
		Padding(0, 1)

	progressFillStyle = lipgloss.NewStyle().
		Background(primaryColor).
		Foreground(backgroundColor)

	progressEmptyStyle = lipgloss.NewStyle().
		Background(mutedColor).
		Foreground(backgroundColor)

	// Status styles
	statusStyle = lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Align(lipgloss.Center).
		Margin(1, 0)

	// Port list styles
	portHeaderStyle = lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true).
		Underline(true).
		MarginBottom(1)

	openPortStyle = lipgloss.NewStyle().
		Foreground(successColor).
		Bold(true)

	highValuePortStyle = lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true).
		Italic(true)

	// Vulnerability styles
	vulnHeaderStyle = lipgloss.NewStyle().
		Foreground(errorColor).
		Bold(true).
		Underline(true).
		MarginBottom(1)

	criticalVulnStyle = lipgloss.NewStyle().
		Foreground(errorColor).
		Bold(true)

	highVulnStyle = lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true)

	mediumVulnStyle = lipgloss.NewStyle().
		Foreground(warningColor)

	lowVulnStyle = lipgloss.NewStyle().
		Foreground(primaryColor)

	// Summary styles
	summaryStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.DoubleBorder()).
		BorderForeground(accentColor).
		Padding(1).
		Margin(1, 0)

	summaryTitleStyle = lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true).
		Underline(true).
		MarginBottom(1)

	// Input styles
	inputStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(secondaryColor).
		Padding(0, 1).
		Margin(0, 1)

	focusedInputStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(accentColor).
		Padding(0, 1).
		Margin(0, 1)

	// Button styles
	buttonStyle = lipgloss.NewStyle().
		Background(primaryColor).
		Foreground(backgroundColor).
		Bold(true).
		Padding(0, 2).
		Margin(0, 1).
		Align(lipgloss.Center)

	activeButtonStyle = lipgloss.NewStyle().
		Background(accentColor).
		Foreground(backgroundColor).
		Bold(true).
		Padding(0, 2).
		Margin(0, 1).
		Align(lipgloss.Center)

	// Help styles
	helpStyle = lipgloss.NewStyle().
		Foreground(mutedColor).
		Italic(true).
		Margin(1, 0)

	// Error styles
	errorStyle = lipgloss.NewStyle().
		Foreground(errorColor).
		Bold(true).
		Margin(1, 0)

	// Layout styles
	columnStyle = lipgloss.NewStyle().
		Width(50).
		Align(lipgloss.Left)

	rightColumnStyle = lipgloss.NewStyle().
		Width(50).
		Align(lipgloss.Right)
)

// GetSeverityStyle returns the appropriate style for vulnerability severity
func GetSeverityStyle(severity string) lipgloss.Style {
	switch severity {
	case "Critical":
		return criticalVulnStyle
	case "High":
		return highVulnStyle
	case "Medium":
		return mediumVulnStyle
	case "Low":
		return lowVulnStyle
	default:
		return lipgloss.NewStyle()
	}
}

// FormatProgressBar creates a visual progress bar
func FormatProgressBar(progress float64, width int) string {
	if width <= 0 {
		width = 40
	}
	
	filled := int(progress * float64(width))
	empty := width - filled
	
	fillStr := progressFillStyle.Render(lipgloss.PlaceHorizontal(filled, lipgloss.Left, ""))
	emptyStr := progressEmptyStyle.Render(lipgloss.PlaceHorizontal(empty, lipgloss.Left, ""))
	
	return progressBarStyle.Render(fillStr + emptyStr)
}

// CenterText centers text within the given width
func CenterText(text string, width int) string {
	return lipgloss.Place(width, 1, lipgloss.Center, lipgloss.Center, text)
}

// FormatTwoColumns creates a two-column layout
func FormatTwoColumns(left, right string) string {
	leftCol := columnStyle.Render(left)
	rightCol := rightColumnStyle.Render(right)
	return lipgloss.JoinHorizontal(lipgloss.Top, leftCol, rightCol)
} 