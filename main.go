package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"

	"ipcrawler/cmd"
	"ipcrawler/ui/styles"
)

func main() {
	// Set color profile for consistent styling
	lipgloss.SetColorProfile(termenv.TrueColor)

	// Initialize styles
	styles.Init()

	// Execute the root command
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
