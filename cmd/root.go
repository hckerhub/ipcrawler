package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"ipcrawler/ui/models"
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
			// Check and escalate privileges if needed
			if !hasRequiredPrivileges() {
				escalatePrivileges()
				return
			}

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

// hasRequiredPrivileges checks if the program has sufficient privileges
func hasRequiredPrivileges() bool {
	// On Windows, this is more complex, but for Unix-like systems:
	if runtime.GOOS == "windows" {
		// For Windows, we'll assume privileges are sufficient for now
		// In a production environment, you'd want to check for admin privileges
		return true
	}

	// Check if running as root (UID 0)
	return os.Geteuid() == 0
}

// escalatePrivileges re-executes the program with sudo
func escalatePrivileges() {
	fmt.Println("🔐 IPCrawler requires elevated privileges for advanced scanning features.")
	fmt.Println("   (SYN scans, OS detection, aggressive timing, etc.)")
	fmt.Println()
	fmt.Println("🚀 Automatically escalating privileges with sudo...")
	fmt.Println()

	// Get the current executable path
	executable, err := os.Executable()
	if err != nil {
		fmt.Printf("❌ Error getting executable path: %v\n", err)
		fmt.Println("   Please run: sudo ./ipcrawler")
		os.Exit(1)
	}

	// Prepare the sudo command with all current arguments
	args := append([]string{executable}, os.Args[1:]...)

	// Create the sudo command
	cmd := exec.Command("sudo", args...)

	// Set up the command to use the current terminal
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute with sudo
	err = cmd.Run()
	if err != nil {
		// Check if it's an exit error to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			}
		}
		fmt.Printf("❌ Error running with sudo: %v\n", err)
		fmt.Println("   Please ensure sudo is available and try: sudo ./ipcrawler")
		os.Exit(1)
	}

	// If we reach here, the sudo command completed successfully
	os.Exit(0)
}

// startTUI initializes and runs the Bubble Tea TUI
func startTUI() {
	// Display privilege confirmation
	if hasRequiredPrivileges() {
		fmt.Println("✅ Running with elevated privileges - all features enabled!")
		fmt.Println()
	}

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
