package scanner

import (
	"os/exec"
	"strings"

	"ipcrawler/tools"
)

// SimpleToolChecker implements ToolAvailabilityChecker using the tools package
type SimpleToolChecker struct {
	installer *tools.ToolInstaller
}

// NewSimpleToolChecker creates a new tool checker
func NewSimpleToolChecker() *SimpleToolChecker {
	return &SimpleToolChecker{
		installer: tools.NewToolInstaller(false, false), // not verbose, not dry run
	}
}

// CheckToolStatus checks if a specific tool is installed
func (tc *SimpleToolChecker) CheckToolStatus(toolName string) (bool, string) {
	// Get tool info from the tools package
	requiredTools := tc.installer.GetRequiredTools()

	for _, tool := range requiredTools {
		if tool.Name == toolName {
			return tc.installer.CheckToolStatus(tool)
		}
	}

	// Fallback to simple command check
	cmd := exec.Command("which", toolName)
	err := cmd.Run()
	if err != nil {
		return false, toolName + " not found in PATH"
	}

	return true, toolName + " is available"
}

// GetMissingToolsForCategory returns missing tools for a specific category
func (tc *SimpleToolChecker) GetMissingToolsForCategory(category string) []string {
	requiredTools := tc.installer.GetRequiredTools()
	var missing []string

	for _, tool := range requiredTools {
		if tool.Category == category {
			if installed, _ := tc.installer.CheckToolStatus(tool); !installed {
				missing = append(missing, tool.Name)
			}
		}
	}

	return missing
}

// GetToolInstallationReport generates a comprehensive tool status report
func (tc *SimpleToolChecker) GetToolInstallationReport() string {
	return tc.installer.GetMissingToolsReport()
}

// GetInstalledReconTools returns a list of installed reconnaissance tools
func (tc *SimpleToolChecker) GetInstalledReconTools() []string {
	requiredTools := tc.installer.GetRequiredTools()
	var installed []string

	for _, tool := range requiredTools {
		if strings.Contains(tool.Category, "reconnaissance") {
			if available, _ := tc.installer.CheckToolStatus(tool); available {
				installed = append(installed, tool.Name)
			}
		}
	}

	return installed
}

// GetInstalledWebTools returns a list of installed web analysis tools
func (tc *SimpleToolChecker) GetInstalledWebTools() []string {
	requiredTools := tc.installer.GetRequiredTools()
	var installed []string

	for _, tool := range requiredTools {
		if strings.Contains(tool.Category, "web_analysis") {
			if available, _ := tc.installer.CheckToolStatus(tool); available {
				installed = append(installed, tool.Name)
			}
		}
	}

	return installed
}

// CheckCoreTools ensures all required tools are available
func (tc *SimpleToolChecker) CheckCoreTools() (bool, []string) {
	requiredTools := tc.installer.GetRequiredTools()
	var missing []string

	for _, tool := range requiredTools {
		if tool.Required {
			if installed, _ := tc.installer.CheckToolStatus(tool); !installed {
				missing = append(missing, tool.Name)
			}
		}
	}

	return len(missing) == 0, missing
}
