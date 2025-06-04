package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Color palette
	PrimaryColor    = lipgloss.Color("#00D4AA")
	SecondaryColor  = lipgloss.Color("#7C3AED")
	AccentColor     = lipgloss.Color("#F59E0B")
	SuccessColor    = lipgloss.Color("#10B981")
	ErrorColor      = lipgloss.Color("#EF4444")
	WarningColor    = lipgloss.Color("#F59E0B")
	InfoColor       = lipgloss.Color("#3B82F6")
	TextColor       = lipgloss.Color("#E5E7EB")
	MutedColor      = lipgloss.Color("#9CA3AF")
	BackgroundColor = lipgloss.Color("#1F2937")
	BorderColor     = lipgloss.Color("#374151")

	// Base styles
	BaseStyle       lipgloss.Style
	TitleStyle      lipgloss.Style
	SubtitleStyle   lipgloss.Style
	HeaderStyle     lipgloss.Style
	FooterStyle     lipgloss.Style
	BorderStyle     lipgloss.Style
	FocusedStyle    lipgloss.Style
	BlurredStyle    lipgloss.Style
	SelectedStyle   lipgloss.Style
	UnselectedStyle lipgloss.Style

	// Status styles
	SuccessStyle lipgloss.Style
	ErrorStyle   lipgloss.Style
	WarningStyle lipgloss.Style
	InfoStyle    lipgloss.Style

	// Layout styles
	ContainerStyle lipgloss.Style
	SidebarStyle   lipgloss.Style
	ContentStyle   lipgloss.Style
	PanelStyle     lipgloss.Style
	CardStyle      lipgloss.Style
	ListItemStyle  lipgloss.Style
	HighlightStyle lipgloss.Style

	// Input styles
	InputStyle        lipgloss.Style
	InputFocusedStyle lipgloss.Style
	InputBlurredStyle lipgloss.Style
	ButtonStyle       lipgloss.Style
	ButtonActiveStyle lipgloss.Style

	// Progress and loading styles
	ProgressStyle lipgloss.Style
	SpinnerStyle  lipgloss.Style
	LoadingStyle  lipgloss.Style

	// Table styles
	TableHeaderStyle lipgloss.Style
	TableCellStyle   lipgloss.Style
	TableRowStyle    lipgloss.Style
	TableBorderStyle lipgloss.Style
)

// Init initializes all the styles
func Init() {
	// Base style
	BaseStyle = lipgloss.NewStyle().
		Foreground(TextColor).
		Background(BackgroundColor)

	// Title styles
	TitleStyle = lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		Padding(0, 1).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(PrimaryColor).
		Align(lipgloss.Center)

	SubtitleStyle = lipgloss.NewStyle().
		Foreground(SecondaryColor).
		Bold(true).
		Italic(true)

	HeaderStyle = lipgloss.NewStyle().
		Foreground(TextColor).
		Background(SecondaryColor).
		Bold(true).
		Padding(0, 1).
		Margin(1, 0)

	FooterStyle = lipgloss.NewStyle().
		Foreground(MutedColor).
		Italic(true).
		Padding(1, 0).
		Align(lipgloss.Center)

	// Border and focus styles
	BorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor).
		Padding(1, 2)

	FocusedStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor).
		Padding(1, 2)

	BlurredStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(MutedColor).
		Padding(1, 2)

	SelectedStyle = lipgloss.NewStyle().
		Foreground(BackgroundColor).
		Background(PrimaryColor).
		Bold(true).
		Padding(0, 1)

	UnselectedStyle = lipgloss.NewStyle().
		Foreground(TextColor).
		Padding(0, 1)

	// Status styles
	SuccessStyle = lipgloss.NewStyle().
		Foreground(SuccessColor).
		Bold(true)

	ErrorStyle = lipgloss.NewStyle().
		Foreground(ErrorColor).
		Bold(true)

	WarningStyle = lipgloss.NewStyle().
		Foreground(WarningColor).
		Bold(true)

	InfoStyle = lipgloss.NewStyle().
		Foreground(InfoColor).
		Bold(true)

	// Layout styles
	ContainerStyle = lipgloss.NewStyle().
		Padding(1, 2).
		Margin(1, 0).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor)

	SidebarStyle = lipgloss.NewStyle().
		Width(30).
		Height(20).
		Padding(1, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(SecondaryColor)

	ContentStyle = lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor)

	PanelStyle = lipgloss.NewStyle().
		Padding(1, 2).
		Margin(1, 0).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(AccentColor)

	CardStyle = lipgloss.NewStyle().
		Padding(1, 2).
		Margin(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor).
		Background(lipgloss.Color("#111827"))

	ListItemStyle = lipgloss.NewStyle().
		Padding(0, 1).
		Margin(0, 0)

	HighlightStyle = lipgloss.NewStyle().
		Foreground(AccentColor).
		Bold(true)

	// Input styles
	InputStyle = lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor)

	InputFocusedStyle = lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor)

	InputBlurredStyle = lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(MutedColor)

	ButtonStyle = lipgloss.NewStyle().
		Foreground(TextColor).
		Background(SecondaryColor).
		Padding(0, 2).
		Margin(0, 1).
		Bold(true)

	ButtonActiveStyle = lipgloss.NewStyle().
		Foreground(BackgroundColor).
		Background(PrimaryColor).
		Padding(0, 2).
		Margin(0, 1).
		Bold(true)

	// Progress styles
	ProgressStyle = lipgloss.NewStyle().
		Foreground(PrimaryColor)

	SpinnerStyle = lipgloss.NewStyle().
		Foreground(AccentColor)

	LoadingStyle = lipgloss.NewStyle().
		Foreground(InfoColor).
		Italic(true)

	// Table styles
	TableHeaderStyle = lipgloss.NewStyle().
		Foreground(BackgroundColor).
		Background(PrimaryColor).
		Bold(true).
		Padding(0, 1).
		Align(lipgloss.Center)

	TableCellStyle = lipgloss.NewStyle().
		Padding(0, 1).
		Align(lipgloss.Left)

	TableRowStyle = lipgloss.NewStyle().
		Border(lipgloss.Border{Bottom: "─"}).
		BorderForeground(BorderColor)

	TableBorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor)
}

// Helper functions for common styling patterns
func RenderTitle(text string) string {
	return TitleStyle.Render(text)
}

func RenderHeader(text string) string {
	return HeaderStyle.Render(text)
}

func RenderSuccess(text string) string {
	return SuccessStyle.Render("✓ " + text)
}

func RenderError(text string) string {
	return ErrorStyle.Render("✗ " + text)
}

func RenderWarning(text string) string {
	return WarningStyle.Render("⚠ " + text)
}

func RenderInfo(text string) string {
	return InfoStyle.Render("ℹ " + text)
}

func RenderHighlight(text string) string {
	return HighlightStyle.Render(text)
}

func RenderCard(content string) string {
	return CardStyle.Render(content)
}

func RenderPanel(content string) string {
	return PanelStyle.Render(content)
}
