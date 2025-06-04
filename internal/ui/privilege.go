package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"ipcrawler/internal/scanner"
)

// PrivilegeModel handles privilege checking in the TUI
type PrivilegeModel struct {
	choice    textinput.Model
	confirmed bool
	denied    bool
	message   string
}

// PrivilegeMsg represents privilege decision messages
type PrivilegeMsg struct {
	Level scanner.PrivilegeLevel
	Error error
}

// NewPrivilegeModel creates a new privilege checking model
func NewPrivilegeModel() PrivilegeModel {
	input := textinput.New()
	input.Placeholder = "y/N"
	input.Focus()
	input.CharLimit = 3
	input.Width = 5

	return PrivilegeModel{
		choice: input,
		message: buildPrivilegeMessage(),
	}
}

func buildPrivilegeMessage() string {
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(warningColor).
		Padding(1).
		Margin(1)

	content := `🔐 UDP Scanning Requires Elevated Privileges

📋 Without sudo privileges:
   ✅ TCP scanning will work normally
   ❌ UDP scanning will be skipped
   ❌ Limited service detection

📋 With sudo privileges:
   ✅ Full TCP and UDP scanning
   ✅ Comprehensive service detection
   ✅ OS fingerprinting capabilities

Would you like to provide sudo privileges for comprehensive scanning?`

	return style.Render(content)
}

// Init implements tea.Model
func (m PrivilegeModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements tea.Model
func (m PrivilegeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			response := m.choice.Value()
			if response == "y" || response == "yes" || response == "Y" {
				// User wants sudo privileges
				m.confirmed = true
				return m, func() tea.Msg {
					// Check if we can get sudo access
					level, err := scanner.CheckPrivileges(false) // Non-interactive check
					if err != nil || level != scanner.Privileged {
						// Try to prompt for sudo (this will happen outside TUI)
						return PrivilegeMsg{Level: scanner.UserDeclined, Error: err}
					}
					return PrivilegeMsg{Level: scanner.Privileged, Error: nil}
				}
			} else {
				// User declined sudo privileges
				m.denied = true
				return m, func() tea.Msg {
					return PrivilegeMsg{Level: scanner.UserDeclined, Error: nil}
				}
			}
		}
	}

	m.choice, cmd = m.choice.Update(msg)
	return m, cmd
}

// View implements tea.Model
func (m PrivilegeModel) View() string {
	if m.confirmed {
		return statusStyle.Render("🔑 Checking sudo access...")
	}
	
	if m.denied {
		return statusStyle.Render("📋 Continuing with TCP-only scanning...")
	}

	return m.message + "\n\n" + 
		inputStyle.Render(m.choice.View()) + "\n\n" +
		helpStyle.Render("Enter your choice and press Enter • Ctrl+C to quit")
} 