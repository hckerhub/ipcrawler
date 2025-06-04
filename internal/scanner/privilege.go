package scanner

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// PrivilegeLevel represents the current privilege level
type PrivilegeLevel int

const (
	Unprivileged PrivilegeLevel = iota
	Privileged
	UserDeclined
)

// CheckPrivileges determines current privilege level and asks user if needed
func CheckPrivileges(interactive bool) (PrivilegeLevel, error) {
	// Check if already running as root/sudo
	if isRunningAsRoot() {
		return Privileged, nil
	}

	// If not interactive (CLI mode), return unprivileged
	if !interactive {
		return Unprivileged, nil
	}

	// Ask user if they want to provide sudo privileges
	fmt.Printf("\n🔐 UDP scanning requires elevated privileges for accurate results.\n")
	fmt.Printf("📋 Without sudo:\n")
	fmt.Printf("   ✅ TCP scanning will work normally\n")
	fmt.Printf("   ❌ UDP scanning will be skipped\n\n")
	fmt.Printf("📋 With sudo:\n")
	fmt.Printf("   ✅ Full TCP and UDP scanning\n")
	fmt.Printf("   ✅ More accurate service detection\n")
	fmt.Printf("   ✅ OS fingerprinting (when available)\n\n")

	for {
		fmt.Printf("Would you like to provide sudo privileges? [y/N]: ")
		
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return Unprivileged, fmt.Errorf("failed to read user input: %w", err)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		
		switch response {
		case "y", "yes":
			// Test sudo access
			if canUseSudo() {
				fmt.Printf("✅ Sudo access confirmed. Full scanning enabled.\n\n")
				return Privileged, nil
			} else {
				fmt.Printf("❌ Unable to obtain sudo access. Continuing with limited scanning.\n\n")
				return UserDeclined, nil
			}
		case "n", "no", "":
			fmt.Printf("📋 Continuing with TCP-only scanning.\n\n")
			return UserDeclined, nil
		default:
			fmt.Printf("Please enter 'y' for yes or 'n' for no.\n")
		}
	}
}

// isRunningAsRoot checks if the current process has root privileges
func isRunningAsRoot() bool {
	return os.Geteuid() == 0
}

// canUseSudo tests if sudo is available and working
func canUseSudo() bool {
	cmd := exec.Command("sudo", "-n", "true")
	err := cmd.Run()
	if err == nil {
		return true
	}

	// Try with password prompt
	fmt.Printf("🔑 Please enter your password for sudo access:\n")
	cmd = exec.Command("sudo", "true")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	err = cmd.Run()
	return err == nil
}

// GetScannerWithPrivileges creates an NMAP scanner with appropriate privileges
func GetScannerWithPrivileges(privilegeLevel PrivilegeLevel) (bool, error) {
	switch privilegeLevel {
	case Privileged:
		// Running with full privileges
		return true, nil
	case Unprivileged, UserDeclined:
		// Running without privileges
		return false, nil
	default:
		return false, fmt.Errorf("unknown privilege level")
	}
}

// NeedsPrivileges checks if a specific scan type needs elevated privileges
func NeedsPrivileges(scanType string) bool {
	switch scanType {
	case "udp":
		return true
	case "syn":
		return true
	case "os_detection":
		return true
	default:
		return false
	}
}

// GetEffectiveUID returns the effective user ID
func GetEffectiveUID() int {
	return syscall.Geteuid()
} 