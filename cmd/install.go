//go:build install

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"ipcrawler/pkg/platform"
)

func main() {
	fmt.Println("🔧 IPCrawler Smart Installer")

	// Create platform detector
	detector := platform.NewDetector()
	osInfo, err := detector.Detect()
	if err != nil {
		fmt.Printf("❌ Failed to detect platform: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("📋 Detected: %s %s (%s)\n", osInfo.OS, osInfo.Version, osInfo.Architecture)

	// Find current binary
	currentBinary, err := os.Executable()
	if err != nil {
		fmt.Printf("❌ Failed to find current binary: %v\n", err)
		os.Exit(1)
	}

	// Get preferred install path
	installPath, err := detector.GetPreferredInstallPath()
	if err != nil {
		fmt.Printf("❌ Failed to determine install path: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("📦 Installing to: %s\n", installPath)

	// Create install directory if it doesn't exist
	if err := os.MkdirAll(installPath, 0755); err != nil {
		fmt.Printf("❌ Failed to create install directory: %v\n", err)
		os.Exit(1)
	}

	// Determine if we need sudo
	needsSudo := false
	if installPath == "/usr/local/bin" || installPath == "/usr/bin" {
		if _, err := os.Stat(installPath); err == nil {
			// Directory exists, check if we can write to it
			testFile := filepath.Join(installPath, ".test_write")
			if f, err := os.Create(testFile); err != nil {
				needsSudo = true
			} else {
				f.Close()
				os.Remove(testFile)
			}
		} else {
			needsSudo = true
		}
	}

	// Install the binary
	binaryName := "ipcrawler"
	if osInfo.OS == "windows" {
		binaryName += ".exe"
	}
	targetPath := filepath.Join(installPath, binaryName)

	if needsSudo {
		fmt.Println("🔐 Installing with elevated privileges...")
		if err := installWithSudo(currentBinary, targetPath); err != nil {
			fmt.Printf("❌ Failed to install with sudo: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("📋 Installing with user privileges...")
		if err := copyFile(currentBinary, targetPath); err != nil {
			fmt.Printf("❌ Failed to copy binary: %v\n", err)
			os.Exit(1)
		}
		if err := os.Chmod(targetPath, 0755); err != nil {
			fmt.Printf("❌ Failed to set permissions: %v\n", err)
			os.Exit(1)
		}
	}

	// Test installation
	if _, err := exec.LookPath(binaryName); err != nil {
		fmt.Printf("⚠️  Warning: %s not found in PATH. You may need to restart your terminal.\n", binaryName)
		fmt.Printf("   Or add %s to your PATH manually.\n", installPath)
	} else {
		fmt.Printf("✅ %s is now available system-wide!\n", binaryName)
	}

	// Show completion message
	fmt.Println("\n🎉 Installation complete!")
	fmt.Printf("Now you can run: %s\n", binaryName)
	fmt.Printf("Key bindings: %s (copy), %s (paste), %s (quit)\n",
		osInfo.KeyBindings.Copy, osInfo.KeyBindings.Paste, osInfo.KeyBindings.Quit)
}

func installWithSudo(src, dst string) error {
	cmd := exec.Command("sudo", "cp", src, dst)
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("sudo", "chmod", "755", dst)
	return cmd.Run()
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, input, 0755)
}
