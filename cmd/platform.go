package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"ipcrawler/pkg/platform"
)

var platformCmd = &cobra.Command{
	Use:   "platform",
	Short: "Show platform and system information",
	Long: `Display detailed information about the detected operating system,
available commands, and platform-specific key bindings.

This information is used by IPCrawler to optimize its behavior
for your specific system and provide appropriate instructions.`,
	Run: func(cmd *cobra.Command, args []string) {
		showPlatformInfo(cmd)
	},
}

func showPlatformInfo(cmd *cobra.Command) {
	// Create platform detector
	detector := platform.NewDetector()
	osInfo, err := detector.Detect()
	if err != nil {
		fmt.Printf("❌ Failed to detect platform: %v\n", err)
		return
	}

	// Check if JSON output requested
	jsonOutput, _ := cmd.Flags().GetBool("json")
	if jsonOutput {
		data, err := json.MarshalIndent(osInfo, "", "  ")
		if err != nil {
			fmt.Printf("❌ Failed to marshal JSON: %v\n", err)
			return
		}
		fmt.Println(string(data))
		return
	}

	// Pretty print platform information
	fmt.Println("🖥️  Platform Information")
	fmt.Println("========================")
	fmt.Printf("Operating System: %s\n", osInfo.OS)
	fmt.Printf("Architecture:     %s\n", osInfo.Architecture)
	if osInfo.Distribution != "" {
		fmt.Printf("Distribution:     %s\n", osInfo.Distribution)
	}
	if osInfo.Version != "" {
		fmt.Printf("Version:          %s\n", osInfo.Version)
	}
	fmt.Printf("Shell:            %s\n", osInfo.Shell)

	fmt.Println("\n⌨️  Key Bindings")
	fmt.Println("================")
	fmt.Printf("Copy:      %s\n", osInfo.KeyBindings.Copy)
	fmt.Printf("Paste:     %s\n", osInfo.KeyBindings.Paste)
	fmt.Printf("Quit:      %s\n", osInfo.KeyBindings.Quit)
	fmt.Printf("Interrupt: %s\n", osInfo.KeyBindings.Interrupt)

	fmt.Println("\n🛠️  Available Commands")
	fmt.Println("======================")
	
	// Essential commands
	essentialCommands := []string{"nmap", "sudo", "git", "go", "curl", "wget"}
	for _, cmd := range essentialCommands {
		if path, exists := osInfo.Commands[cmd]; exists {
			fmt.Printf("✅ %-8s %s\n", cmd, path)
		} else {
			fmt.Printf("❌ %-8s Not found\n", cmd)
		}
	}

	// Package managers
	fmt.Println("\n📦 Package Managers")
	fmt.Println("===================")
	packageManagers := []string{"brew", "apt", "dnf", "yum", "pacman", "apk"}
	found := false
	for _, pm := range packageManagers {
		if path, exists := osInfo.Commands[pm]; exists {
			fmt.Printf("✅ %-8s %s\n", pm, path)
			found = true
		}
	}
	if !found {
		fmt.Println("❌ No package managers found")
	}

	// Installation paths
	fmt.Println("\n📁 Installation Paths")
	fmt.Println("=====================")
	if paths, err := detector.GetSystemPaths(); err == nil {
		for i, path := range paths {
			inPath := "❌"
			if detector.IsInPath(path) {
				inPath = "✅"
			}
			priority := ""
			if i == 0 {
				priority = " (preferred)"
			}
			fmt.Printf("%s %s%s\n", inPath, path, priority)
		}
	}

	// Show installation suggestions
	fmt.Println("\n💡 Installation Suggestions")
	fmt.Println("============================")
	
	if preferredPath, err := detector.GetPreferredInstallPath(); err == nil {
		fmt.Printf("Recommended install path: %s\n", preferredPath)
		
		// Show install commands
		switch osInfo.OS {
		case "darwin":
			fmt.Println("\nTo install IPCrawler system-wide:")
			fmt.Println("  ./install.sh                    # Use installer script")
			fmt.Println("  make install                    # Use Makefile")
		case "linux":
			fmt.Println("\nTo install IPCrawler system-wide:")
			fmt.Println("  ./install.sh                    # Use installer script")
			fmt.Println("  make install                    # Use Makefile")
			if _, exists := osInfo.Commands["sudo"]; !exists {
				fmt.Println("  Note: sudo not found - may need manual installation")
			}
		case "windows":
			fmt.Println("\nTo install IPCrawler system-wide:")
			fmt.Println("  Use Windows Subsystem for Linux (WSL)")
			fmt.Println("  Or install Go and build manually")
		}
	}

	fmt.Println("\n🎯 Optimization Notes")
	fmt.Println("=====================")
	
	// Platform-specific optimizations
	switch osInfo.OS {
	case "darwin":
		fmt.Println("• macOS detected - using Cmd key shortcuts")
		fmt.Println("• Terminal app supports full color display")
		if _, exists := osInfo.Commands["brew"]; exists {
			fmt.Println("• Homebrew available for easy dependency management")
		} else {
			fmt.Println("• Consider installing Homebrew for easier package management")
		}
	case "linux":
		fmt.Println("• Linux detected - using Ctrl+Shift shortcuts in terminal")
		fmt.Println("• Full NMAP functionality available")
		if osInfo.Distribution != "unknown" {
			fmt.Printf("• %s-specific optimizations enabled\n", osInfo.Distribution)
		}
	case "windows":
		fmt.Println("• Windows detected - limited NMAP functionality")
		fmt.Println("• Consider using WSL for better compatibility")
		if osInfo.Distribution == "WSL" {
			fmt.Println("• WSL detected - Linux-style operations available")
		}
	}
}

func init() {
	rootCmd.AddCommand(platformCmd)
	
	platformCmd.Flags().BoolP("json", "j", false, "Output platform information as JSON")
} 