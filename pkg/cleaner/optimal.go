package cleaner

import (
	"os/exec"
	"fmt"
)

// SetOptimalWindowsSettings applies recommended Windows settings for best stability and compatibility.
// Currently, it disables Fast Boot. Extend this function to add more tweaks as needed.
func SetOptimalWindowsSettings(verbose bool) error {
	var input string

	// Wizard: Ask user for power plan preference
	fmt.Println("Choose a power plan to apply:")
	fmt.Println("1. High Performance")
	fmt.Println("2. Balanced (Recommended)")
	fmt.Print("Enter your choice (1 or 2): ")

	var choice int
	_, err := fmt.Scanln(&choice)
	if err != nil {
		return fmt.Errorf("failed to read input: %v", err)
	}

	var powerPlanCmd *exec.Cmd
	switch choice {
	case 1:
		// Set to High Performance
		if verbose {
			fmt.Println("[VERBOSE] Running command: powershell -Command powercfg /setactive SCHEME_MIN")
		}
		powerPlanCmd = exec.Command("powershell", "-Command", "powercfg /setactive SCHEME_MIN")
	case 2:
		// Set to Balanced
		if verbose {
			fmt.Println("[VERBOSE] Running command: powershell -Command powercfg /setactive SCHEME_BALANCED")
		}
		powerPlanCmd = exec.Command("powershell", "-Command", "powercfg /setactive SCHEME_BALANCED")
	default:
		fmt.Println("Invalid choice. Skipping power plan change.")
	}

	if powerPlanCmd != nil {
		output, err := powerPlanCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to set power plan: %v\nOutput: %s", err, string(output))
		}
		fmt.Println("Power plan applied successfully.")
	}

	if verbose {
		fmt.Println("[VERBOSE] Running command: powershell -Command Set-ItemProperty -Path 'HKLM:\\SYSTEM\\CurrentControlSet\\Control\\Session Manager\\Power' -Name 'HiberbootEnabled' -Value 0")
	}
	cmd := exec.Command("powershell", "-Command", "Set-ItemProperty -Path 'HKLM:\\SYSTEM\\CurrentControlSet\\Control\\Session Manager\\Power' -Name 'HiberbootEnabled' -Value 0")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to disable Fast Boot: %v\nOutput: %s", err, string(output))
	}

	// 1. Adjust Visual Effects for Best Performance
	fmt.Println("\n1. Adjust Visual Effects for Best Performance:")
	fmt.Println("   Disables most Windows animations and effects to improve speed.")
	fmt.Print("   [e]nable / [d]isable / [s]kip? ")
	fmt.Scanln(&input)
	if input == "e" || input == "E" {
		cmd := exec.Command("powershell", "-Command", "Set-ItemProperty -Path 'HKCU:\\Software\\Microsoft\\Windows\\CurrentVersion\\Explorer\\VisualEffects' -Name 'VisualFXSetting' -Value 2")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Failed to adjust visual effects: %v\nOutput: %s\n", err, string(output))
		} else {
			fmt.Println("Visual effects set for best performance.")
		}
	} else if input == "d" || input == "D" {
		cmd := exec.Command("powershell", "-Command", "Set-ItemProperty -Path 'HKCU:\\Software\\Microsoft\\Windows\\CurrentVersion\\Explorer\\VisualEffects' -Name 'VisualFXSetting' -Value 0")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Failed to revert visual effects: %v\nOutput: %s\n", err, string(output))
		} else {
			fmt.Println("Visual effects reverted to default.")
		}
	}

	// 2. Disable Transparency Effects
	fmt.Println("\n2. Disable Transparency Effects:")
	fmt.Println("   Turns off window transparency to reduce GPU usage.")
	fmt.Print("   [e]nable / [d]isable / [s]kip? ")
	fmt.Scanln(&input)
	if input == "e" || input == "E" {
		cmd := exec.Command("powershell", "-Command", "Set-ItemProperty -Path 'HKCU:\\Software\\Microsoft\\Windows\\CurrentVersion\\Themes\\Personalize' -Name 'EnableTransparency' -Value 0")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Failed to disable transparency: %v\nOutput: %s\n", err, string(output))
		} else {
			fmt.Println("Transparency effects disabled.")
		}
	} else if input == "d" || input == "D" {
		cmd := exec.Command("powershell", "-Command", "Set-ItemProperty -Path 'HKCU:\\Software\\Microsoft\\Windows\\CurrentVersion\\Themes\\Personalize' -Name 'EnableTransparency' -Value 1")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Failed to enable transparency: %v\nOutput: %s\n", err, string(output))
		} else {
			fmt.Println("Transparency effects enabled.")
		}
	}

	// 3. Enable Storage Sense
	fmt.Println("\n3. Enable Storage Sense:")
	fmt.Println("   Automatically frees up disk space by deleting unnecessary files.")
	fmt.Print("   [e]nable / [d]isable / [s]kip? ")
	fmt.Scanln(&input)
	if input == "e" || input == "E" {
		cmd := exec.Command("powershell", "-Command", "Set-ItemProperty -Path 'HKCU:\\Software\\Microsoft\\Windows\\CurrentVersion\\StorageSense\\Parameters\\StoragePolicy' -Name '01' -Value 1")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Failed to enable Storage Sense: %v\nOutput: %s\n", err, string(output))
		} else {
			fmt.Println("Storage Sense enabled.")
		}
	} else if input == "d" || input == "D" {
		cmd := exec.Command("powershell", "-Command", "Set-ItemProperty -Path 'HKCU:\\Software\\Microsoft\\Windows\\CurrentVersion\\StorageSense\\Parameters\\StoragePolicy' -Name '01' -Value 0")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Failed to disable Storage Sense: %v\nOutput: %s\n", err, string(output))
		} else {
			fmt.Println("Storage Sense disabled.")
		}
	}

	// 4. Disable Startup Delay
	fmt.Println("\n4. Disable Startup Delay:")
	fmt.Println("   Speeds up startup for apps in the Startup folder.")
	fmt.Print("   [e]nable / [d]isable / [s]kip? ")
	fmt.Scanln(&input)
	if input == "e" || input == "E" {
		// Ensure the Serialize key exists before setting the property
		exec.Command("powershell", "-Command", "if (-not (Test-Path 'HKCU:\\Software\\Microsoft\\Windows\\CurrentVersion\\Explorer\\Serialize')) { New-Item -Path 'HKCU:\\Software\\Microsoft\\Windows\\CurrentVersion\\Explorer\\Serialize' | Out-Null }").Run()
		cmd := exec.Command("powershell", "-Command", "Set-ItemProperty -Path 'HKCU:\\Software\\Microsoft\\Windows\\CurrentVersion\\Explorer\\Serialize' -Name 'StartupDelayInMSec' -Value 0")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Failed to disable startup delay: %v\nOutput: %s\n", err, string(output))
		} else {
			fmt.Println("Startup delay disabled.")
		}
	} else if input == "d" || input == "D" {
		cmd := exec.Command("powershell", "-Command", "Remove-ItemProperty -Path 'HKCU:\\Software\\Microsoft\\Windows\\CurrentVersion\\Explorer\\Serialize' -Name 'StartupDelayInMSec' -ErrorAction SilentlyContinue")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Failed to restore startup delay: %v\nOutput: %s\n", err, string(output))
		} else {
			fmt.Println("Startup delay restored to default.")
		}
	}

	// Add more optimal settings here as needed

	return nil
} 