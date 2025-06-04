package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"ipcrawler/internal/scanner"
	"ipcrawler/internal/ui"
)

var tuiCmd = &cobra.Command{
	Use:   "tui [IP]",
	Short: "Start the interactive Terminal User Interface",
	Long: `Launch the interactive Terminal User Interface (TUI) for IPCrawler.
The TUI provides a beautiful, interactive way to perform IP scanning and
vulnerability analysis with real-time progress updates and tabbed results.

Features:
  • Interactive IP input with validation
  • Real-time scanning progress visualization
  • Tabbed results view (Summary, Ports, Vulnerabilities, Details)
  • Keyboard shortcuts for navigation
  • Beautiful ASCII art and cyberpunk-themed styling

Examples:
  ipcrawler tui                    # Start TUI in interactive mode
  ipcrawler tui 10.10.10.1         # Start TUI with pre-filled target
  ipcrawler tui 192.168.1.100 -a   # Start with aggressive scanning enabled`,
	Run: func(cmd *cobra.Command, args []string) {
		aggressive, _ := cmd.Flags().GetBool("aggressive")
		verbose, _ := cmd.Flags().GetBool("verbose")
		outputFile, _ := cmd.Flags().GetString("output")
		
		// Check privileges BEFORE starting TUI to avoid stdin conflicts
		fmt.Printf("🔍 Checking system privileges for enhanced scanning...\n")
		privilegeLevel, err := scanner.CheckPrivileges(true) // true = interactive
		if err != nil {
			fmt.Printf("⚠️  Privilege check failed: %v\n", err)
			fmt.Printf("📋 Continuing with limited scanning capabilities.\n\n")
			privilegeLevel = scanner.UserDeclined
		}
		
		if len(args) > 0 {
			// Start TUI with target pre-filled
			ui.StartTUIWithTarget(args[0], outputFile, verbose, aggressive, privilegeLevel)
		} else {
			// Start TUI in interactive mode
			ui.StartTUI(privilegeLevel)
		}
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
	
	tuiCmd.Flags().BoolP("aggressive", "a", false, "Enable aggressive scanning by default")
	tuiCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output by default")
	tuiCmd.Flags().StringP("output", "o", "", "Default output file for results")
} 