package cmd

import (
	"fmt"
	"os"

	"ipcrawler/ui/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	version = "1.0.0"
	rootCmd = &cobra.Command{
		Use:   "ipcrawler",
		Short: "IPCrawler - Advanced IP Analysis & Penetration Testing Tool",
		Long: `
╔══════════════════════════════════════════════════════════════╗
║                        IPCrawler                             ║
║          Advanced IP Analysis & Penetration Testing         ║
║                                                              ║
║  A modern TUI tool for comprehensive IP reconnaissance,      ║
║  network scanning, and penetration testing. Designed        ║
║  specifically for Hack the Box CTFs and security research.  ║
║                                                              ║
║  Features:                                                   ║
║  • Advanced nmap scanning with custom profiles              ║
║  • Web technology identification with whatweb               ║
║  • Port 22 (SSH) and 80 (HTTP) specialized analysis        ║
║  • URL discovery and enumeration                           ║
║  • CMS and database detection                               ║
║  • Vulnerability scanning                                   ║
║  • Modern, intuitive TUI interface                         ║
╚══════════════════════════════════════════════════════════════╝`,
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			// Start the TUI application
			startTUI()
		},
	}
)

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Suppress output except errors")
	rootCmd.PersistentFlags().StringP("output", "o", "", "Output file for results")
	rootCmd.PersistentFlags().StringP("format", "f", "json", "Output format (json, xml, txt)")

	// Version template
	rootCmd.SetVersionTemplate(`
╔══════════════════════════════════════════╗
║               IPCrawler                  ║
║            Version {{.Version}}                ║
║                                          ║
║  Advanced IP Analysis & PenTest Tool     ║
╚══════════════════════════════════════════╝
`)
}

// startTUI initializes and runs the Bubble Tea TUI
func startTUI() {
	// Create the main model
	m := models.NewMainModel()

	// Create the Bubble Tea program
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	// Start the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
