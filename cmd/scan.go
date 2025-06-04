package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"ipcrawler/internal/models"
	"ipcrawler/internal/scanner"
)

var scanCmd = &cobra.Command{
	Use:   "scan [IP]",
	Short: "Perform a direct scan without TUI",
	Long: `Perform a direct port scan and vulnerability analysis on the target IP address.
This command runs in non-interactive mode and outputs results to stdout or a file.

Examples:
  ipcrawler scan 10.10.10.1
  ipcrawler scan 192.168.1.100 --aggressive --output results.json
  ipcrawler scan 10.0.0.1 --verbose --format table`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetIP := args[0]

		// Get flags
		aggressive, _ := cmd.Flags().GetBool("aggressive")
		verbose, _ := cmd.Flags().GetBool("verbose")
		outputFile, _ := cmd.Flags().GetString("output")
		format, _ := cmd.Flags().GetString("format")

		// Create scanner
		engine := scanner.NewEngine(targetIP, aggressive, verbose, false) // false = non-interactive CLI mode

		if verbose {
			fmt.Printf("🎯 Starting scan of target: %s\n", targetIP)
			fmt.Printf("⚙️  Configuration: aggressive=%v, verbose=%v\n", aggressive, verbose)
			fmt.Printf("⏱️  Started at: %s\n\n", time.Now().Format("15:04:05"))
		}

		// Monitor progress
		go monitorProgress(engine.GetProgressChannel(), verbose)

		// Execute scan
		ctx := context.Background()
		result, err := engine.StartFullScan(ctx)
		if err != nil {
			fmt.Printf("❌ Scan failed: %v\n", err)
			os.Exit(1)
		}

		// Output results
		if err := outputResults(result, outputFile, format); err != nil {
			fmt.Printf("❌ Failed to output results: %v\n", err)
			os.Exit(1)
		}

		if verbose {
			fmt.Printf("\n✅ Scan completed successfully in %v\n", result.Duration.Truncate(time.Second))
		}
	},
}

func monitorProgress(progressChan <-chan models.ScanProgress, verbose bool) {
	for progress := range progressChan {
		if verbose {
			fmt.Printf("📊 [%s] %.1f%% - %s\n",
				progress.Phase,
				progress.Progress*100,
				progress.Message)
		} else {
			// Simple progress indicator
			fmt.Print(".")
		}
	}
	if !verbose {
		fmt.Println() // New line after dots
	}
}

func outputResults(result *models.ScanResult, outputFile, format string) error {
	var output string
	var err error

	switch format {
	case "json":
		output, err = formatJSON(result)
	case "table":
		output, err = formatTable(result)
	default:
		output, err = formatSummary(result)
	}

	if err != nil {
		return err
	}

	if outputFile != "" {
		return os.WriteFile(outputFile, []byte(output), 0644)
	}

	fmt.Println(output)
	return nil
}

func formatJSON(result *models.ScanResult) (string, error) {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func formatTable(result *models.ScanResult) (string, error) {
	output := fmt.Sprintf(`
🎯 SCAN RESULTS FOR %s
═══════════════════════════════════════════════════════════════

📊 SUMMARY
• Scan Duration: %v
• TCP Open Ports: %d
• UDP Open Ports: %d
• Total Vulnerabilities: %d

🔌 OPEN PORTS
`,
		result.TargetIP,
		result.Duration.Truncate(time.Second),
		len(result.TCPPorts),
		len(result.UDPPorts),
		len(result.Vulnerabilities),
	)

	// TCP Ports
	if len(result.TCPPorts) > 0 {
		output += "\nTCP Ports:\n"
		output += "Port\tProtocol\tService\t\tVersion\n"
		output += "----\t--------\t-------\t\t-------\n"
		for _, port := range result.TCPPorts {
			output += fmt.Sprintf("%d\t%s\t\t%s\t\t%s\n",
				port.Number, port.Protocol, port.Service, port.Version)
		}
	}

	// UDP Ports
	if len(result.UDPPorts) > 0 {
		output += "\nUDP Ports:\n"
		output += "Port\tProtocol\tService\t\tVersion\n"
		output += "----\t--------\t-------\t\t-------\n"
		for _, port := range result.UDPPorts {
			output += fmt.Sprintf("%d\t%s\t\t%s\t\t%s\n",
				port.Number, port.Protocol, port.Service, port.Version)
		}
	}

	// Vulnerabilities
	if len(result.Vulnerabilities) > 0 {
		output += "\n🚨 VULNERABILITIES\n"
		for _, vuln := range result.Vulnerabilities {
			output += fmt.Sprintf("\n[%s] %s:%d - %s\n",
				vuln.Severity, vuln.Service, vuln.Port, vuln.Type)
			output += fmt.Sprintf("Description: %s\n", vuln.Description)
			if len(vuln.Suggestions) > 0 {
				output += "Suggestions:\n"
				for _, suggestion := range vuln.Suggestions {
					output += fmt.Sprintf("  • %s\n", suggestion)
				}
			}
		}
	}

	return output, nil
}

func formatSummary(result *models.ScanResult) (string, error) {
	s := result.Summary
	output := fmt.Sprintf(`
🎯 IPCrawler Scan Results
═══════════════════════════════

Target: %s
Duration: %v
Started: %s
Completed: %s

📊 Port Summary:
• Total Open Ports: %d
• TCP Ports: %d
• UDP Ports: %d

🚨 Vulnerability Summary:
• Critical: %d
• High: %d
• Medium: %d
• Low: %d

🔌 Key Services Found:
`,
		result.TargetIP,
		result.Duration.Truncate(time.Second),
		result.StartTime.Format("15:04:05"),
		result.EndTime.Format("15:04:05"),
		s.TotalOpenPorts,
		s.TCPOpenPorts,
		s.UDPOpenPorts,
		s.CriticalVulns,
		s.HighVulns,
		s.MediumVulns,
		s.LowVulns,
	)

	// List high-value ports found
	for _, port := range result.TCPPorts {
		if port.IsHighValuePort() {
			output += fmt.Sprintf("• %s\n", port.GetPortDisplay())
		}
	}

	for _, port := range result.UDPPorts {
		if port.IsHighValuePort() {
			output += fmt.Sprintf("• %s\n", port.GetPortDisplay())
		}
	}

	if len(result.Vulnerabilities) > 0 {
		output += "\n🚨 Top Vulnerabilities:\n"
		count := 0
		for _, vuln := range result.Vulnerabilities {
			if count >= 5 { // Show top 5
				break
			}
			output += fmt.Sprintf("• [%s] %s on port %d\n",
				vuln.Severity, vuln.Type, vuln.Port)
			count++
		}
	}

	output += "\n💡 Next Steps:\n"
	output += "• Run 'ipcrawler tui' for interactive analysis\n"
	output += "• Use --format=json for detailed output\n"
	output += "• Check individual services for known exploits\n"

	return output, nil
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().BoolP("aggressive", "a", false, "Enable aggressive scanning")
	scanCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
	scanCmd.Flags().StringP("output", "o", "", "Output file for results")
	scanCmd.Flags().StringP("format", "f", "summary", "Output format (summary, table, json)")
}
