package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"ipcrawler/internal/scanner"
	"ipcrawler/internal/ui"
)

const Version = "0.1"

var (
	targetIP   string
	outputFile string
	verbose    bool
	aggressive bool
)

const asciiArt = `
██╗██████╗  ██████╗██████╗  █████╗ ██╗    ██╗██╗     ███████╗██████╗ 
██║██╔══██╗██╔════╝██╔══██╗██╔══██╗██║    ██║██║     ██╔════╝██╔══██╗
██║██████╔╝██║     ██████╔╝███████║██║ █╗ ██║██║     █████╗  ██████╔╝
██║██╔═══╝ ██║     ██╔══██╗██╔══██║██║███╗██║██║     ██╔══╝  ██╔══██╗
██║██║     ╚██████╗██║  ██║██║  ██║╚███╔███╔╝███████╗███████╗██║  ██║
╚═╝╚═╝      ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝ ╚══╝╚══╝ ╚══════╝╚══════╝╚═╝  ╚═╝
                                                                      
        🎯 Advanced IP Scanner & Vulnerability Hunter 🎯
             Hack The Box Edition - Built with ❤️  in Go
                              v` + Version + ` by hckerhub
`

var rootCmd = &cobra.Command{
	Use:     "ipcrawler",
	Version: Version,
	Short:   "Advanced IP scanner and vulnerability hunter for Hack The Box",
	Long: asciiArt + `

IPCrawler is an advanced terminal user interface for IP scanning and vulnerability discovery.
It intelligently combines NMAP scanning with automated vulnerability assessment to help
you discover attack vectors on Hack The Box machines.

Features:
  • Interactive TUI with beautiful ASCII art
  • Intelligent NMAP TCP & UDP scanning
  • Automated vulnerability discovery
  • Smart port-based service detection
  • Beautiful terminal reports
  • Export results to multiple formats

Author: hckerhub (https://github.com/hckerhub)`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check privileges BEFORE starting TUI to avoid stdin conflicts
		fmt.Printf("🔍 Checking system privileges for enhanced scanning...\n")
		privilegeLevel, err := scanner.CheckPrivileges(true) // true = interactive
		if err != nil {
			fmt.Printf("⚠️  Privilege check failed: %v\n", err)
			fmt.Printf("📋 Continuing with limited scanning capabilities.\n\n")
			privilegeLevel = scanner.UserDeclined
		}

		if targetIP == "" && len(args) == 0 {
			// Start interactive TUI mode
			ui.StartTUI(privilegeLevel)
			return
		}

		// Use provided IP from argument or flag
		ip := targetIP
		if len(args) > 0 {
			ip = args[0]
		}

		if ip == "" {
			fmt.Println("❌ Error: Please provide a target IP address")
			cmd.Help()
			os.Exit(1)
		}

		// Start TUI with target IP
		ui.StartTUIWithTarget(ip, outputFile, verbose, aggressive, privilegeLevel)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(&targetIP, "target", "t", "", "Target IP address to scan")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file for scan results")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.Flags().BoolVarP(&aggressive, "aggressive", "a", false, "Enable aggressive scanning mode")

	// Set version template
	rootCmd.SetVersionTemplate(`IPCrawler v{{.Version}} - Advanced IP Scanner & Vulnerability Hunter
Built with ❤️ for Hack The Box by hckerhub
`)
}
